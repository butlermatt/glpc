package lexer

// TokenType identifies the type of token from the constants
type TokenType string

// Token represents an individual segment of source code.
type Token struct {
	// Type is the TokenType of this token.
	Type TokenType
	// Lexeme is the raw string from the input
	Lexeme string
	// Filename is the name of the file the token was found in.
	Filename string
	// Line is the line this token was found on.
	Line int
}

func NewToken(ty TokenType, lex string, filename string, line int) *Token {
	return &Token{ty, lex, filename, line}
}

// TokenTypes of the various tokens
const (
	// Single Character tokens
	Colon     TokenType = ":"
	Comma     TokenType = ","
	Dot       TokenType = "."
	LBrace    TokenType = "{"
	RBrace    TokenType = "}"
	LBracket  TokenType = "["
	RBracket  TokenType = "]"
	LParen    TokenType = "("
	RParen    TokenType = ")"
	Minus     TokenType = "-"
	Plus      TokenType = "+"
	Semicolon TokenType = ";"
	Slash     TokenType = "/"
	Star      TokenType = "*"
	Percent   TokenType = "%"

	// Single or two character tokens.
	TildSlash TokenType = "~/"
	Bang      TokenType = "!"
	BangEq    TokenType = "!="
	Equal     TokenType = "="
	EqualEq   TokenType = "=="
	Greater   TokenType = ">"
	GreaterEq TokenType = ">="
	Less      TokenType = "<"
	LessEq    TokenType = "<="

	// Literals
	Ident    TokenType = "IDENT"
	String   TokenType = "STRING"
	UTString TokenType = "UNTERMINATED STRING"
	NumberF  TokenType = "FLOAT NUMBER"
	NumberI  TokenType = "INT NUMBER"

	// Keywords
	And      TokenType = "AND"
	Break    TokenType = "BREAK"
	Class    TokenType = "CLASS"
	Continue TokenType = "CONTINUE"
	Else     TokenType = "ELSE"
	False    TokenType = "FALSE"
	Fn       TokenType = "FN"
	For      TokenType = "FOR"
	If       TokenType = "IF"
	Null     TokenType = "NULL"
	Or       TokenType = "OR"
	Print    TokenType = "PRINT"
	Return   TokenType = "RETURN"
	Super    TokenType = "SUPER"
	This     TokenType = "THIS"
	True     TokenType = "TRUE"
	Var      TokenType = "VAR"
	While    TokenType = "WHILE"

	Illegal TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
)

var keywords = map[string]TokenType{
	"and":      And,
	"break":    Break,
	"class":    Class,
	"continue": Continue,
	"else":     Else,
	"false":    False,
	"fn":       Fn,
	"for":      For,
	"if":       If,
	"null":     Null,
	"or":       Or,
	"print":    Print,
	"return":   Return,
	"super":    Super,
	"this":     This,
	"true":     True,
	"var":      Var,
	"while":    While,
}