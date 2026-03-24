package main

const eof byte = 0

// ---------------------------------------------------------------------------
// DFA definitions
// ---------------------------------------------------------------------------

// State represents a DFA state.
type State int

const (
	StateStart     State = iota // initial state
	StateAssign                 // seen '='
	StateEq                     // seen '=='
	StatePlus                   // seen '+'
	StateMinus                  // seen '-'
	StateInc                    // seen '++'
	StateDec                    // seen '--'
	StateSemiColon              // seen ';'
	StateIdent                  // inside an identifier
	StateInt                    // inside an integer literal
	stateCount                  // sentinel – total number of states
)

// NoTransition signals that no valid transition exists.
const NoTransition State = -1

// CharClass groups input bytes into equivalence classes for the DFA.
type CharClass int

const (
	ClassLetter    CharClass = iota // a-z A-Z _
	ClassDigit                      // 0-9
	ClassEquals                     // =
	ClassPlus                       // +
	ClassMinus                      // -
	ClassSemiColon                  // ;
	ClassOther                      // everything else
	classCount                      // sentinel
)

// transition[state][charClass] → next state (NoTransition if none).
var transition [stateCount][classCount]State

// acceptToken[state] → token type emitted when this state is accepting.
// A zero value (TokenIllegal) means the state is not accepting.
var acceptToken [stateCount]TokenType

func init() {
	// Default every cell to NoTransition.
	for s := range transition {
		for c := range transition[s] {
			transition[s][c] = NoTransition
		}
	}

	// --- transitions from StateStart ---
	transition[StateStart][ClassEquals] = StateAssign
	transition[StateStart][ClassPlus] = StatePlus
	transition[StateStart][ClassMinus] = StateMinus
	transition[StateStart][ClassSemiColon] = StateSemiColon
	transition[StateStart][ClassLetter] = StateIdent
	transition[StateStart][ClassDigit] = StateInt

	// --- multi-character operators ---
	transition[StateAssign][ClassEquals] = StateEq // '=' then '=' → '=='
	transition[StatePlus][ClassPlus] = StateInc    // '+' then '+' → '++'
	transition[StateMinus][ClassMinus] = StateDec  // '-' then '-' → '--'

	// --- identifiers: [A-Za-z_][A-Za-z0-9_]* ---
	transition[StateIdent][ClassLetter] = StateIdent
	transition[StateIdent][ClassDigit] = StateIdent

	// --- integer literals: [0-9]+ ---
	transition[StateInt][ClassDigit] = StateInt

	// --- accepting states ---
	acceptToken[StateAssign] = TokenAssign
	acceptToken[StateEq] = TokenEq
	acceptToken[StatePlus] = TokenPlus
	acceptToken[StateMinus] = TokenMinus
	acceptToken[StateInc] = TokenInc
	acceptToken[StateDec] = TokenDec
	acceptToken[StateSemiColon] = TokenSemiColon
	acceptToken[StateIdent] = TokenIdent
	acceptToken[StateInt] = TokenInt
}

// classify returns the character class for ch.
func classify(ch byte) CharClass {
	switch {
	case isLetter(ch):
		return ClassLetter
	case isDigit(ch):
		return ClassDigit
	case ch == '=':
		return ClassEquals
	case ch == '+':
		return ClassPlus
	case ch == '-':
		return ClassMinus
	case ch == ';':
		return ClassSemiColon
	default:
		return ClassOther
	}
}

// ---------------------------------------------------------------------------
// Lexer
// ---------------------------------------------------------------------------

// Lexer converts source text into a stream of tokens using a table-driven DFA.
type Lexer struct {
	input        string
	position     int  // index of current character
	readPosition int  // index of next character to read
	ch           byte // current character under examination
	line         int
}

// NewLexer creates a lexer and positions it on the first character.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

// readChar advances the lexer by one character.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = eof
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// rewind repositions the lexer so that the next character to be examined is
// at index pos.
func (l *Lexer) rewind(pos int) {
	if pos >= len(l.input) {
		l.ch = eof
		l.position = len(l.input)
		l.readPosition = len(l.input)
	} else {
		l.ch = l.input[pos]
		l.position = pos
		l.readPosition = pos + 1
	}
}

// skipWhitespace consumes spaces, tabs, and newlines, updating the line
// counter when a newline is consumed.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

// nextToken runs the DFA to scan and return the next token.
func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	// Handle EOF before entering the DFA.
	if l.ch == eof {
		return Token{Type: TokenEOF, Lexeme: "", Line: l.line}
	}

	state := StateStart
	startPos := l.position

	// Track the last accepting state so we can implement longest-match.
	lastAcceptPos := -1
	var lastAcceptType TokenType

	for {
		class := classify(l.ch)
		next := transition[state][class]
		if next == NoTransition {
			break
		}

		state = next
		if acceptToken[state] != TokenIllegal {
			lastAcceptPos = l.position
			lastAcceptType = acceptToken[state]
		}
		l.readChar()
	}

	// No accepting state was ever reached — illegal character.
	if lastAcceptPos == -1 {
		tok := Token{Type: TokenIllegal, Lexeme: string(l.input[startPos]), Line: l.line}
		l.readChar()
		return tok
	}

	// Rewind if the DFA consumed characters past the last accept.
	if l.position > lastAcceptPos+1 {
		l.rewind(lastAcceptPos + 1)
	}

	lexeme := l.input[startPos : lastAcceptPos+1]

	// Resolve keywords for identifiers.
	if lastAcceptType == TokenIdent {
		lastAcceptType = LookupIdent(lexeme)
	}

	return Token{Type: lastAcceptType, Lexeme: lexeme, Line: l.line}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Tokenize processes the whole input and returns all tokens, including EOF.
func Tokenize(input string) []Token {
	l := NewLexer(input)
	var tokens []Token
	for {
		tok := l.nextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}
