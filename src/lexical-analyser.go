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

	case '(':
		tok := l.makeToken(TokenLParen)
		l.readChar()
		return tok

	case ')':
		tok := l.makeToken(TokenRParen)
		l.readChar()
		return tok

	case '{':
		tok := l.makeToken(TokenLBrace)
		l.readChar()
		return tok

	case '}':
		tok := l.makeToken(TokenRBrace)
		l.readChar()
		return tok

	case '-':
		if l.peekChar() == '-' {
			tok := l.makeTokenFromPair(TokenDec)
			l.readChar()
			return tok
		}

		tok := l.makeToken(TokenMinus)
		l.readChar()
		return tok

	case '*':
		tok := l.makeToken(TokenStar)
		l.readChar()
		return tok

	case '<':
		if l.peekChar() == '=' {
			tok := l.makeTokenFromPair(TokenLE)
			l.readChar()
			return tok
		}
		tok := l.makeToken(TokenLT)
		l.readChar()
		return tok

	case '>':
		if l.peekChar() == '=' {
			tok := l.makeTokenFromPair(TokenGE)
			l.readChar()
			return tok
		}
		tok := l.makeToken(TokenGT)
		l.readChar()
		return tok

	case '!':
		if l.peekChar() == '=' {
			tok := l.makeTokenFromPair(TokenNE)
			l.readChar()
			return tok
		}
		tok := l.makeToken(TokenNot)
		l.readChar()
		return tok

	case '\'':
		tok := l.readCharLiteral()
		return tok

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

func (l *Lexer) readCharLiteral() Token {
	startPos := l.position // at opening quote '
	line := l.line

	// consume opening quote
	l.readChar()

	// empty / newline / eof right after opening quote
	if l.ch == eof || l.ch == '\n' || l.ch == '\r' || l.ch == '\'' {
		if l.ch == '\'' {
			l.readChar()
		}
		return Token{Type: TokenIllegal, Lexeme: l.input[startPos:l.position], Line: line}
	}

	// consume literal body (one char or one escape sequence)
	if l.ch == '\\' {
		l.readChar() // consume backslash
		if l.ch == eof || l.ch == '\n' || l.ch == '\r' {
			return Token{Type: TokenIllegal, Lexeme: l.input[startPos:l.position], Line: line}
		}
		l.readChar() // consume escaped char (e.g. n, ', \)
	} else {
		l.readChar() // consume normal char
	}

	// must close with '
	if l.ch != '\'' {
		for l.ch != eof && l.ch != '\n' && l.ch != '\r' && l.ch != '\'' {
			l.readChar()
		}
		if l.ch == '\'' {
			l.readChar()
		}
		return Token{Type: TokenIllegal, Lexeme: l.input[startPos:l.position], Line: line}
	}

	// consume closing quote
	l.readChar()

	// include quotes in lexeme
	return Token{Type: TokenCharLit, Lexeme: l.input[startPos:l.position], Line: line}
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
