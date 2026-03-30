package main

import "fmt"

type TokenType int

const (
	TokenIllegal TokenType = iota
	TokenEOF
	TokenIdent   // identifier (non-keyword)
	TokenInt     // e.g. 123
	TokenFloat   // e.g. 3.14
	TokenString  // e.g. "hello"
	TokenDouble  // e.g. 3.14e-10
	TokenAssign  // =
	TokenEq      // ==
	TokenPlus    // +
	TokenMinus   // -
	TokenInc     // ++
	TokenDec     // --
	TokenSemiColon // ;
	// C keywords
	TokenAuto     // 'auto'
	TokenBreak    // 'break'
	TokenCase     // 'case'
	TokenChar     // 'char'
	TokenConst    // 'const'
	TokenContinue // 'continue'
	TokenDefault  // 'default'
	TokenDo       // 'do'
	TokenDoubleKw // 'double'
	TokenElse     // 'else'
	TokenEnum     // 'enum'
	TokenExtern   // 'extern'
	TokenFloatKw  // 'float'
	TokenFor      // 'for'
	TokenGoto     // 'goto'
	TokenIf       // 'if'
	TokenInline   // 'inline'
	TokenIntKw    // 'int'
	TokenLong     // 'long'
	TokenRegister // 'register'
	TokenRestrict // 'restrict'
	TokenReturn   // 'return'
	TokenShort    // 'short'
	TokenSigned   // 'signed'
	TokenSizeof   // 'sizeof'
	TokenStatic   // 'static'
	TokenStruct   // 'struct'
	TokenSwitch   // 'switch'
	TokenTypedef  // 'typedef'
	TokenUnion    // 'union'
	TokenUnsigned // 'unsigned'
	TokenVoid     // 'void'
	TokenVolatile // 'volatile'
	TokenWhile    // 'while'
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

var keywords = map[string]TokenType{
	"auto":     TokenAuto,
	"break":    TokenBreak,
	"case":     TokenCase,
	"char":     TokenChar,
	"const":    TokenConst,
	"continue": TokenContinue,
	"default":  TokenDefault,
	"do":       TokenDo,
	"double":   TokenDoubleKw,
	"else":     TokenElse,
	"enum":     TokenEnum,
	"extern":   TokenExtern,
	"float":    TokenFloatKw,
	"for":      TokenFor,
	"goto":     TokenGoto,
	"if":       TokenIf,
	"inline":   TokenInline,
	"int":      TokenIntKw,
	"long":     TokenLong,
	"register": TokenRegister,
	"restrict": TokenRestrict,
	"return":   TokenReturn,
	"short":    TokenShort,
	"signed":   TokenSigned,
	"sizeof":   TokenSizeof,
	"static":   TokenStatic,
	"struct":   TokenStruct,
	"switch":   TokenSwitch,
	"typedef":  TokenTypedef,
	"union":    TokenUnion,
	"unsigned": TokenUnsigned,
	"void":     TokenVoid,
	"volatile": TokenVolatile,
	"while":    TokenWhile,
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
