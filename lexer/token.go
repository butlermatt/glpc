package lexer

// TokenType identifies the type of token from the constants
type TokenType string

// Token represents an individual segment of source code.
type Token struct {
	// Type is the TokenType of this token.
	Type TokenType
	// Lexeme is the raw string from the input
	Lexeme string
	// Line is the line this token was found on.
	Line int
}

func NewToken(ty TokenType, lex string, line int) Token {
	return Token{ty, lex, line}
}

// TokenTypes of the various tokens
const (
	// Single Character tokens
	Colon     = ":"
	Comma     = ","
	Dot       = "."
	LBrace    = "{"
	RBrace    = "}"
	LBracket  = "["
	RBracket  = "]"
	LParen    = "("
	RParen    = ")"
	Minus     = "-"
	Plus      = "+"
	Semicolon = ";"
	Slash     = "/"
	Star      = "*"

	// Single or two character tokens.
	Bang      = "!"
	BangEq    = "!="
	Equal     = "="
	EqualEq   = "=="
	Greater   = ">"
	GreaterEq = ">="
	Less      = "<"
	LessEq    = "<="

	// Literals
	Ident    = "IDENT"
	String   = "STRING"
	UTSTRING = "UNTERMINATED STRING"
	NumberF  = "FLOAT NUMBER"
	NumberI  = "INT NUMBER"

	// Keywords
	And      = "AND"
	Break    = "BREAK"
	Class    = "CLASS"
	Continue = "CONTINUE"
	Else     = "ELSE"
	False    = "FALSE"
	Fn       = "FN"
	For      = "FOR"
	If       = "IF"
	Null     = "NULL"
	Or       = "OR"
	Print    = "PRINT"
	Return   = "RETURN"
	Super    = "SUPER"
	This     = "THIS"
	True     = "TRUE"
	Var      = "VAR"
	While    = "WHILE"

	Illegal = "ILLEGAL"
	EOF     = "EOF"
)
