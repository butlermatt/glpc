package lexer

type Lexer struct {
	file    string
	input   []byte
	line    int
	start   int
	current int

	tokens []Token // Scanned tokens
	index  int     // Index in token list.
}

// New returns a new Lexer populated with the specified input program.
func New(input []byte, filename string) *Lexer {
	l := &Lexer{file: filename, input: input, line: 1}
	return l
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.input)
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
