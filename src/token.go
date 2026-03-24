package main

import "fmt"

type TokenType int

const (
	TokenIllegal TokenType = iota
	TokenEOF
	TokenIdent
	TokenInt
	TokenAssign // =
	TokenEq     // ==
	TokenPlus   // +
	TokenInc    // ++
	TokenSemi   // ;
	TokenIntKw  // 'int' keyword
	TokenReturn // 'return' keyword
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

var keywords = map[string]TokenType{
	"int":    TokenIntKw,
	"return": TokenReturn,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TokenIdent
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %d, lexeme: %s, line: %d}", t.Type, t.Lexeme, t.Line)
}
