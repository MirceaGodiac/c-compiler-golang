package main

const eof byte = 0

// Lexer converts source text into a stream of tokens.
type Lexer struct {
	input        string
	position     int  // index of current character in input
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

// peekChar returns the next character without advancing the lexer.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return eof
	}
	return l.input[l.readPosition]
}

// skipWhitespace consumes spaces, tabs, and newlines.
// It also updates the current line counter when a newline is consumed.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

// makeToken builds a single-character token from the current character.
func (l *Lexer) makeToken(tokenType TokenType) Token {
	return Token{Type: tokenType, Lexeme: string(l.ch), Line: l.line}
}

// makeTokenFromPair builds a two-character token from the current and next character.
// The lexer is advanced once to consume the second character.
func (l *Lexer) makeTokenFromPair(tokenType TokenType) Token {
	first := l.ch
	l.readChar()
	return Token{Type: tokenType, Lexeme: string(first) + string(l.ch), Line: l.line}
}

// nextToken scans and returns the next token in the input.
func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok := l.makeTokenFromPair(TokenEq)
			l.readChar()
			return tok
		}

		tok := l.makeToken(TokenAssign)
		l.readChar()
		return tok

	case '+':
		if l.peekChar() == '+' {
			tok := l.makeTokenFromPair(TokenInc)
			l.readChar()
			return tok
		}

		tok := l.makeToken(TokenPlus)
		l.readChar()
		return tok

	case ';':
		tok := l.makeToken(TokenSemi)
		l.readChar()
		return tok

	case eof:
		return Token{Type: TokenEOF, Lexeme: "", Line: l.line}

	default:
		if isLetter(l.ch) {
			lexeme := l.readIdentifier()
			return Token{Type: LookupIdent(lexeme), Lexeme: lexeme, Line: l.line}
		}

		if isDigit(l.ch) {
			lexeme := l.readNumber()
			return Token{Type: TokenInt, Lexeme: lexeme, Line: l.line}
		}

		tok := l.makeToken(TokenIllegal)
		l.readChar()
		return tok
	}
}

// readIdentifier consumes an identifier: [A-Za-z_][A-Za-z0-9_]*
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber consumes a decimal integer literal: [0-9]+
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter reports whether ch is a valid identifier letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit reports whether ch is an ASCII decimal digit.
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
