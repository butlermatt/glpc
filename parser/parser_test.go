package parser

import (
	"testing"

	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
)

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)

		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: incorrect number of statements. expected=%d, got=%d", i, 1, len(stmts))
		}

		s, ok := stmts[0].(*object.ExpressionStmt)
		if !ok {
			t.Fatalf("test %d: statement wrong type. expected=*object.ExpressionStmt, got=%T", i, stmts[0])
		}

		if !testBooleanLiteral(t, s.Expression, tt.value) {
			t.Errorf("last error occurred on test %d", i)
		}
	}
}

func TestListExpression(t *testing.T) {
		input := "[0, 'one', true];"

		l := lexer.New([]byte(input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: statements is incorrect length. expected=%d, got=%d", 1, len(stmts))
		}

		stmt := stmts[0].(*object.ExpressionStmt)

		list, ok := stmt.Expression.(*object.ListExpr)
		if !ok {
			t.Fatalf("expr is wrong type. expected=*object.ListExpr, got=%T", stmt.Expression)
		}

		if len(list.Values) != 3 {
			t.Fatalf("list contains incorrect number of values. expected=%d, got=%d", 3, len(list.Values))
		}

		testNumberLiteral(t, list.Values[0], 0)
		testStringLiteral(t, list.Values[1], "one")
		testBooleanLiteral(t, list.Values[2], true)
}

func TestNumberLiteralExpression(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{"5;", 5},
		{"10;", 10},
		{"123.456;", 123.456},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)

		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: incorrect number of statements. exepected=%d, got=%d", i, 1, len(stmts))
		}

		s, ok := stmts[0].(*object.ExpressionStmt)
		if !ok {
			t.Fatalf("test %d: statement[0] wrong type. expected=*object.ExpressionStmt, got=%T", i, stmts[0])
		}

		if !testNumberLiteral(t, s.Expression, tt.value) {
			t.Errorf("last test tha failed was %d", i)
		}
	}
}

func TestNullLiteralExpression(t *testing.T) {
	var value interface{}
	input := "null;"

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)

	stmts := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. exepected=%d, got=%d", 1, len(stmts))
	}

	s, ok := stmts[0].(*object.ExpressionStmt)
	if !ok {
		t.Fatalf("statement[0] wrong type. expected=*object.ExpressionStmt, got=%T", stmts[0])
	}

	ne, ok := s.Expression.(*object.NullExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.NullExpression, got=%T", s.Expression)
	}

	if ne.Value != value {
		t.Fatalf("null value incorrect. expected=%v, got=%v", value, ne.Value)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct{
		input string
		expected string
	}{
		{`"hello world";`, "hello world"},
		{`'hello world';`, "hello world"},
		{"`hello world`;", "hello world"},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: incorrect number of statements. expected=%d, got=%d", i, 1, len(stmts))
		}

		stmt := stmts[0].(*object.ExpressionStmt)
		lit, ok := stmt.Expression.(*object.StringExpr)
		if !ok {
			t.Fatalf("test %d: expression wrong type. expeted=*object.StringExpr, got=%T", i, stmt.Expression)
		}

		if lit.Value != tt.expected {
			t.Errorf("test %d: value is wrong. expected=%q, got==%q", i, tt.expected, lit.Value)
		}
	}
}

func TestUnterminatedString(t *testing.T) {
	input := `"hello world;`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts := p.Parse()

	if len(stmts) != 0 {
		t.Fatalf("wrong number of statements. expected=0, got=%d", len(stmts))
	}

	errs := p.Errors()
	// 2 because of Unterminated string, then missing semicolon
	if len(errs) != 2 {
		t.Fatalf("wrong number of errors. expected=1, got=%d", len(errs))
	}

	e := errs[0]
	if e.Line != 1 {
		t.Errorf("error on wrong line, expected=1, got=%d", e.Line)
	}

	if e.Where != `"hello world;` {
		t.Errorf("error at wrong location. expected=%q, got=%q", `"hello world;`, e.Where)
	}

	if e.Msg != "unterminated string" {
		t.Errorf("wrong error message. expected=%q, got=%q", "unterminated string", e.Msg)
	}
}

func testBooleanLiteral(t *testing.T, expr object.Expr, value bool) bool {
	be, ok := expr.(*object.BooleanExpr)
	if !ok {
		t.Errorf("expr not correct type. expected=*object.BooleanExpr, got=%T", expr)
		return false
	}

	if be.Value != value {
		t.Errorf("value did not match. expected=%t, got=%t", value, be.Value)
		return false
	}

	return true
}

func testNumberLiteral(t *testing.T, expr object.Expr, value interface{}) bool {
	ne, ok := expr.(*object.NumberExpr)
	if !ok {
		t.Errorf("expr not correct type. expected=*object.NumberExpr, got=%T", expr)
		return false
	}

	var match bool
	switch v := value.(type) {
	case int:
		match = v == ne.Int
	case float64:
		match = v == ne.Float
	default:
		t.Errorf("unknown value to for number literal: %T", v)
		return false
	}

	if !match {
		t.Errorf("NumberExpression value did not match. expected=%v, got=%v or %v", value, ne.Int, ne.Float)
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, expr object.Expr, value string) bool {
	se, ok := expr.(*object.StringExpr)
	if !ok {
		t.Errorf("expr not correct type. expected=*object.StringExpr, got=%T", expr)
		return false
	}

	if se.Value != value {
		t.Errorf("value did not match. expected=%q, got=%q", value, se.Value)
		return false
	}

	return true
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error on line: %d at %s: %s", msg.Line, msg.Where, msg.Msg)
	}
	t.FailNow()
}
