package main

import "fmt"

type TokenType int

const (
	TokenIllegal TokenType = iota
	TokenEOF
	TokenIdent         // identifier (non-keyword)
	TokenInt           // e.g. 123
	TokenFloat         // e.g. 3.14
	TokenString        // e.g. "hello"
	TokenDouble        // e.g. 3.14e-10
	TokenAssign        // =
	TokenEq            // ==
	TokenPlus          // +
	TokenMinus         // -
	TokenInc           // ++
	TokenDec           // --
	TokenSemiColon     // ;
	TokenIntKeyword    // 'int' keyword
	TokenFloatKeyword  // 'float' keyword
	TokenStringKeyword // 'string' keyword
	TokenVarKeyword    // 'var' keyword
	TokendoubleKeyword // 'double' keyword
	TokenIf            // 'if' keyword
	TokenReturn        // 'return' keyword
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

var keywords = map[string]TokenType{
	"int":    TokenIntKeyword,
	"float":  TokenFloatKeyword,
	"double": TokendoubleKeyword,
	"if":     TokenIf,
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
