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
	TokenMinus  // -
	TokenDec    // --
	TokenArrow  // ->
	TokenStar   // *
	TokenSlash  // /
	TokenPercent // %
	TokenLt     // <
	TokenLtEq   // <=
	TokenGt     // >
	TokenGtEq   // >=
	TokenBang   // !
	TokenNotEq  // !=
	TokenAmp    // &
	TokenAnd    // &&
	TokenPipe   // |
	TokenOr     // ||
	TokenSemi   // ;
	TokenComma  // ,
	TokenDot    // .
	TokenLParen // (
	TokenRParen // )
	TokenLBrace // {
	TokenRBrace // }
	TokenLBrack // [
	TokenRBrack // ]
	// keywords
	TokenIntKw    // 'int'
	TokenReturn   // 'return'
	TokenIf       // 'if'
	TokenElse     // 'else'
	TokenWhile    // 'while'
	TokenFor      // 'for'
	TokenDo       // 'do'
	TokenVoid     // 'void'
	TokenChar     // 'char'
	TokenFloat    // 'float'
	TokenDouble   // 'double'
	TokenLong     // 'long'
	TokenShort    // 'short'
	TokenBreak    // 'break'
	TokenContinue // 'continue'
	TokenStruct   // 'struct'
	TokenSwitch   // 'switch'
	TokenCase     // 'case'
	TokenDefault  // 'default'
	TokenSizeof   // 'sizeof'
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

var keywords = map[string]TokenType{
	"int":      TokenIntKw,
	"return":   TokenReturn,
	"if":       TokenIf,
	"else":     TokenElse,
	"while":    TokenWhile,
	"for":      TokenFor,
	"do":       TokenDo,
	"void":     TokenVoid,
	"char":     TokenChar,
	"float":    TokenFloat,
	"double":   TokenDouble,
	"long":     TokenLong,
	"short":    TokenShort,
	"break":    TokenBreak,
	"continue": TokenContinue,
	"struct":   TokenStruct,
	"switch":   TokenSwitch,
	"case":     TokenCase,
	"default":  TokenDefault,
	"sizeof":   TokenSizeof,
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
