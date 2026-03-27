package main

import "fmt"

type TokenType int

const (
	TokenIllegal TokenType = iota
	TokenEOF
	TokenIdent
	TokenInt
	TokenAssign           // =
	TokenEq               // ==
	TokenPlus             // +
	TokenInc              // ++
	TokenSemi             // ;
	TokenIntKw            // 'int' keyword
	TokenReturn           // 'return' keyword
	TokenIf               // 'if' keyword
	TokenLParen           // '('
	TokenRParen           // ')'
	TokenLBrace           // '{'
	TokenRBrace           // '}'
	TokenLT               // <
	TokenGT               // >
	TokenLE               // <=
	TokenGE               // >=
	TokenNE               // !=
	TokenAnd              // &&
	TokenOr               // ||
	TokenNot              // !
	TokenComma            // ,
	TokenDot              // .
	TokenColon            // :
	TokenWhile            // 'while' keyword
	TokenDo               // 'do' keyword
	TokenFor              // 'for' keyword
	TokenChar             // 'char' keyword
	TokenCharLit          // character literal (e.g. 'a', '\n')
	TokenStar             // *
	TokenMinus            // -
	TokenDec              // --

)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

var keywords = map[string]TokenType{
	"int":    TokenIntKw,
	"return": TokenReturn,
	"if":     TokenIf,
	"while":  TokenWhile,
	"for":    TokenFor,
	"char":   TokenChar,
	"do":     TokenDo,
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
