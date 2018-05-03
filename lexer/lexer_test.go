package lexer

import "testing"

func TestLexer_ScanTokensCount(t *testing.T) {
	// Should be a token for each token in the input, plus one for EOF
	tests := []struct {
		input string
		count int
	}{
		{`+- =[]()*/{}!`, 13},
		{`+
// Ignore comment
-
{}[]`, 7},
		{`>=,==,<=,!= // Test 2 character tokens`, 8},
		{`<>%:.`, 6},
		{`~`, 2},           // Illegal and EOF
		{`"something"`, 2}, // Double quote String
		{`'something'`, 2}, // Single Quote String
		{`"something`, 2},  // Unterminated String
		{"`A fancy\nMultiline\nString`", 2},
		{"`Unterminated\nMultiline\n", 2},
		{`2342.2323`, 2},
		{`ident;if;`, 5}, // Ident, Semicolon, If keyword, semicolon EOF
	}

	for i, tt := range tests {
		l := New([]byte(tt.input), "testFile")
		l.ScanTokens()

		if len(l.tokens) != tt.count {
			t.Errorf("test %d: did not generate enough tokens. expected=%d, got=%d", i+1, tt.count, len(l.tokens))
			t.Errorf("%+v\n", l.tokens)
		}
	}
}

func TestLexer_NextToken(t *testing.T) {
	input := `+-
{}[]
42 123.456
// Ignore this
!= == <= >=
"a string" 'another "string"'
`
	input += "`A\nMultiline\nString`\n"
	input += `!
,
.
fn something true false
if and or else for while
class null this super return;
break continue
= +=-=%=*= /= ~/=
`

	expected := []struct {
		ty      TokenType
		literal string
		line    int
	}{
		{Plus, "+", 1},
		{Minus, "-", 1},
		{LBrace, "{", 2},
		{RBrace, "}", 2},
		{LBracket, "[", 2},
		{RBracket, "]", 2},
		{NumberI, "42", 3},
		{NumberF, "123.456", 3},
		{BangEq, "!=", 5},
		{EqualEq, "==", 5},
		{LessEq, "<=", 5},
		{GreaterEq, ">=", 5},
		{String, `a string`, 6},
		{String, `another "string"`, 6},
		{RawString, "A\nMultiline\nString", 7},
		{Bang, "!", 10},
		{Comma, ",", 11},
		{Dot, ".", 12},
		{Fn, "fn", 13},
		{Ident, "something", 13},
		{True, "true", 13},
		{False, "false", 13},
		{If, "if", 14},
		{And, "and", 14},
		{Or, "or", 14},
		{Else, "else", 14},
		{For, "for", 14},
		{While, "while", 14},
		{Class, "class", 15},
		{Null, "null", 15},
		{This, "this", 15},
		{Super, "super", 15},
		{Return, "return", 15},
		{Semicolon, ";", 15},
		{Break, "break", 16},
		{Continue, "continue", 16},
		{Equal, "=", 17},
		{PlusEq, "+=", 17},
		{MinusEq, "-=", 17},
		{PercentEq, "%=", 17},
		{StarEq, "*=", 17},
		{SlashEq, "/=", 17},
		{TildSlashEq, "~/=", 17},
		{EOF, "", 18},
	}

	l := New([]byte(input), "test.lpc")
	l.ScanTokens()

	for i, expect := range expected {
		tok := l.NextToken()
		if tok == nil {
			t.Fatalf("test %d: unexpected missing token. expected=%q", i, expect.ty)
		}

		if tok.Type != expect.ty {
			t.Errorf("test %d: unexpected token. expected=%q, got=%q", i, expect.ty, tok.Type)
		}

		if tok.Lexeme != expect.literal {
			t.Errorf("test %d: unexpected lexeme. expected=%q, got=%q", i, expect.literal, tok.Lexeme)
		}

		if tok.Line != expect.line {
			t.Errorf("test %d: unexpected line. expected=%d, got=%d", i, expect.line, tok.Line)
		}
	}
}
