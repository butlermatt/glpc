package lexer

type Lexer struct {
	file    string
	input   []byte
	line    int
	start   int
	current int

	tokens []*Token // Scanned tokens
	index  int      // Index in token list.
}

// New returns a new Lexer populated with the specified input program.
func New(input []byte, filename string) *Lexer {
	l := &Lexer{file: filename, input: input, line: 1}
	return l
}

// NextToken provides the next token that is in the list of tokens
func (l *Lexer) NextToken() *Token {
	if l.index >= len(l.tokens) {
		return nil
	}
	tok := l.tokens[l.index]
	l.index += 1
	return tok
}

// ScanTokens will scan all the input and generate a slice of tokens.
func (l *Lexer) ScanTokens() {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}

	l.addToken(NewToken(EOF, "", l.file, l.line))
}

func (l *Lexer) scanToken() {
	c := l.readChar()

	switch c {
	case ' ', '\t', '\r': // Ignore whitespace.
		break
	case '\n': // Whitespace but add line.
		l.line += 1
	case ':':
		l.addTokenType(Colon)
	case ',':
		l.addTokenType(Comma)
	case '.':
		l.addTokenType(Dot)
	case ';':
		l.addTokenType(Semicolon)
	case '(':
		l.addTokenType(LParen)
	case ')':
		l.addTokenType(RParen)
	case '{':
		l.addTokenType(LBrace)
	case '}':
		l.addTokenType(RBrace)
	case '[':
		l.addTokenType(LBracket)
	case ']':
		l.addTokenType(RBracket)
	case '-':
		l.addTokenType(Minus)
	case '+':
		l.addTokenType(Plus)
	case '%':
		l.addTokenType(Percent)
	case '*':
		l.addTokenType(Star)
	case '/':
		if l.match('/') {
			// Consume comments to end of line (or file)
			for l.peek() != '\n' && l.peek() != 0 {
				l.readChar()
			}
		} else {
			l.addTokenType(Slash)
		}
	case '=':
		if l.match('=') {
			l.addTokenType(EqualEq)
		} else {
			l.addTokenType(Equal)
		}
	case '>':
		if l.match('=') {
			l.addTokenType(GreaterEq)
		} else {
			l.addTokenType(Greater)
		}
	case '<':
		if l.match('=') {
			l.addTokenType(LessEq)
		} else {
			l.addTokenType(Less)
		}
	case '!':
		if l.match('=') {
			l.addTokenType(BangEq)
		} else {
			l.addTokenType(Bang)
		}
	case '"', '\'':
		l.string(c)
	case '`':
		l.multilineString()
	default:
		if isDigit(c) {
			l.number()
		} else if isAlpha(c) {
			l.identifier()
		} else {
			l.addTokenType(Illegal)
		}
	}
}

func (l *Lexer) addTokenType(ty TokenType) {
	l.tokens = append(l.tokens, NewToken(ty, string(l.input[l.start:l.current]), l.file, l.line))
}

func (l *Lexer) addToken(token *Token) {
	l.tokens = append(l.tokens, token)
}

func (l *Lexer) match(char byte) bool {
	if l.isAtEnd() || l.input[l.current] != char {
		return false
	}

	l.current += 1
	return true
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.input[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.input) {
		return 0
	}

	return l.input[l.current+1]
}

func (l *Lexer) readChar() byte {
	ch := l.input[l.current]
	l.current += 1
	return ch
}

func (l *Lexer) identifier() {
	for isAlphaNumeric(l.peek()) {
		l.readChar()
	}

	id := string(l.input[l.start:l.current])
	if tokenType, ok := keywords[id]; ok {
		l.addTokenType(tokenType)
	} else {
		l.addTokenType(Ident)
	}
}

func (l *Lexer) number() {
	tokType := NumberI
	for isDigit(l.peek()) {
		l.readChar()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		tokType = NumberF
		l.readChar()

		for isDigit(l.peek()) {
			l.readChar()
		}
	}

	l.addTokenType(tokType)
}

func (l *Lexer) string(endChar byte) {
	for !l.isAtEnd() && l.peek() != endChar && l.peek() != '\n' {
		l.readChar()
	}

	if l.isAtEnd() || l.peek() == '\n' {
		l.addToken(NewToken(UTString, string(l.input[l.start:l.current]), l.file, l.line))
		return
	}

	l.readChar() // consume last quote
	l.addToken(NewToken(String, string(l.input[l.start+1:l.current-1]), l.file, l.line))
}

func (l *Lexer) multilineString() {
	line := l.line
	for !l.isAtEnd() && l.peek() != '`' {
		if l.peek() == '\n' {
			l.line += 1
		}
		l.readChar()
	}

	if l.isAtEnd() {
		l.addToken(NewToken(UTString, string(l.input[l.start:l.current]), l.file, line))
		return
	}

	l.readChar()
	l.addToken(NewToken(String, string(l.input[l.start+1:l.current-1]), l.file, line))
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.input)
}

// Helper functions.

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
