package parser

import (
	"testing"

	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
)

func TestIfStatement(t *testing.T) {
	input := `if (a == 1) 
  x = true;
else
  y = false;
`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	ifStmt, ok := stmts[0].(*object.IfStmt)
	if !ok {
		t.Fatalf("Statement wrong type. expected=*object.IfStmt, got=%T", stmts[0])
	}

	testBinaryExpression(t, ifStmt.Condition, "a", "==", 1)

	es := ifStmt.Then.(*object.ExpressionStmt)
	ae := es.Expression.(*object.AssignExpr)
	if ae.Name.Lexeme != "x" {
		t.Errorf("wrong then branch. name expectd=%q, got=%q", "x", ae.Name.Lexeme)
	}
	be := ae.Value.(*object.BooleanExpr)
	if be.Value != true {
		t.Errorf("wrong value assigned. expected=%t, got=%t", true, be.Value)
	}

	es = ifStmt.Else.(*object.ExpressionStmt)
	ae = es.Expression.(*object.AssignExpr)
	if ae.Name.Lexeme != "y" {
		t.Errorf("wrong else branch. name expected=%q, got=%q", "y", ae.Name.Lexeme)
	}

	be = ae.Value.(*object.BooleanExpr)
	if be.Value != false {
		t.Errorf("wrong value assigned. expected=%t, got=%t", false, be.Value)
	}
}

func TestVarStatement(t *testing.T) {
	tests := []struct {
		input string
		ident string
		value interface{}
	}{
		{"var x = 5;", "x", int(5)},
		{"var hello = 5.23;", "hello", float64(5.23)},
		{"var y = true;", "y", true},
		{"var x = y;", "x", "y"},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("test %d: incorrect number of statements. expected=%d, got=%d", i, 1, len(stmts))
			continue
		}

		if !testVariable(t, stmts[0], tt.ident) {
			continue
		}

		val := stmts[0].(*object.VarStmt)
		if !testLiteralExpression(t, val.Value, tt.value) {
			t.Errorf("failed on test %d", i)
			continue
		}
	}
}

func TestAssignExpression(t *testing.T) {
	tests := []struct {
		input string
		name  string
		value interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = 2.25;", "y", 2.25},
		{"test = true;", "test", true},
		{"list = [1, 2, 3];", "list", []int{1, 2, 3}},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: incorrect number of statements. expected=%d, got=%d", i, 1, len(stmts))
		}

		s := stmts[0].(*object.ExpressionStmt)
		ae, ok := s.Expression.(*object.AssignExpr)
		if !ok {
			t.Fatalf("test %d: expression is wrong type. expected=*object.AssignExpr, got=%T", i, s.Expression)
		}

		if ae.Name.Lexeme != tt.name {
			t.Errorf("test %d: name does not match. expected=%q, got=%q", i, tt.name, ae.Name.Lexeme)
		}

		if !testLiteralExpression(t, ae.Value, tt.value) {
			t.Errorf("last error in test %d", i)
		}
	}
}

func TestCompoundAssignExpressions(t *testing.T) {
	type bin struct {
		left  string
		oper  string
		value interface{}
	}

	tests := []struct {
		input string
		name  string
		value bin
	}{
		{"x += 1;", "x", bin{"x", "+", 1}},
		{"y -= 2.2;", "y", bin{"y", "-", 2.2}},
		{"test *= 3;", "test", bin{"test", "*", 3}},
		{"z /= 4;", "z", bin{"z", "/", 4}},
		{"a ~/= 5;", "a", bin{"a", "~/", 5}},
		{"b %= 6;", "b", bin{"b", "%", 6}},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("test %d: incorrect statements length. expected=%d, got=%d", i+1, 1, len(stmts))
			continue
		}

		s, ok := stmts[0].(*object.ExpressionStmt)
		if !ok {
			t.Errorf("test %d: statement incorrect type. expected=*object.ExpressionStmt, got=%T", i+1, stmts[0])
			continue
		}

		ae, ok := s.Expression.(*object.AssignExpr)
		if !ok {
			t.Errorf("test %d: expression wrong type. expected=*object.AssignExpr, got=%T", i+1, s.Expression)
			continue
		}

		if ae.Name.Lexeme != tt.name {
			t.Errorf("test %d: name incorrect. expected=%q, got=%q", i+1, tt.name, ae.Name.Lexeme)
		}

		if !testBinaryExpression(t, ae.Value, tt.value.left, tt.value.oper, tt.value.value) {
			t.Errorf("last error occured in test %d", i+1)
		}
	}
}

func TestBinaryExpressions(t *testing.T) {
	tests := []struct {
		input string
		left  interface{}
		oper  string
		right interface{}
	}{
		{"5 + 4;", 5, "+", 4},
		{"5 - 4;", 5, "-", 4},
		{"5 * 4;", 5, "*", 4},
		{"5 / 4;", 5, "/", 4},
		{"5 % 4;", 5, "%", 4},
		{"5 ~/ 4;", 5, "~/", 4},
		{"5 < 4;", 5, "<", 4},
		{"5 > 4;", 5, ">", 4},
		{"5 <= 4;", 5, "<=", 4},
		{"5.5 == 4.4;", 5.5, "==", 4.4},
		{"5 != 4;", 5, "!=", 4},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("on test %d: incorrect number of statements. expected=1, got=%d", i+1, len(stmts))
			continue
		}

		s, ok := stmts[0].(*object.ExpressionStmt)
		if !ok {
			t.Errorf("on test %d: Statement wrong type. expected=*object.ExpressionStmt, got=%T", i+1, stmts[0])
			continue
		}

		testBinaryExpression(t, s.Expression, tt.left, tt.oper, tt.right)
	}
}

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

func TestGroupingExpression(t *testing.T) {
	input := "(5);"

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("statements incorrect length. expected=%d, got=%d", 1, len(stmts))
	}

	stmt := stmts[0].(*object.ExpressionStmt)

	gr, ok := stmt.Expression.(*object.GroupingExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.GroupingExpr, got=%T", stmt.Expression)
	}

	testNumberLiteral(t, gr.Expression, 5)

}

func TestIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		left  interface{}
		right interface{}
	}{
		{"a[5];", "a", 5},
		{"some[thing];", "some", "thing"},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("on test %d: incorrect number of statements. expected=1, got=%d", i+1, len(stmts))
			continue
		}

		s, ok := stmts[0].(*object.ExpressionStmt)
		if !ok {
			t.Errorf("on test %d: Statement wrong type. expected=*object.ExpressionStmt, got=%T", i+1, stmts[0])
			continue
		}

		ie, ok := s.Expression.(*object.IndexExpr)
		if !ok {
			t.Errorf("on test %d: Expression wrong type. expected=*object.IndexExpr, got=%T", i+1, s.Expression)
			continue
		}

		testLiteralExpression(t, ie.Left, tt.left)
		testLiteralExpression(t, ie.Right, tt.right)
	}
}

func TestListExpression(t *testing.T) {
	input := "[0, 'one', true];"

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("statements is incorrect length. expected=%d, got=%d", 1, len(stmts))
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

func TestLogicalExpressions(t *testing.T) {
	tests := []struct {
		input string
		left  interface{}
		oper  string
		right interface{}
	}{
		{"true and true;", true, "and", true},
		{"true or false;", true, "or", false},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("on test %d: incorrect number of statements. expected=1, got=%d", i+1, len(stmts))
			continue
		}

		s, ok := stmts[0].(*object.ExpressionStmt)
		if !ok {
			t.Errorf("on test %d: Statement wrong type. expected=*object.ExpressionStmt, got=%T", i+1, stmts[0])
			continue
		}

		testLogicalExpression(t, s.Expression, tt.left, tt.oper, tt.right)
	}
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

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input   string
		numErrs int
		where   string
		msg     string
	}{
		{`"hello world;`, 2, `"hello world;`, "Unterminated string."},
		{"7 = x;", 1, "=", "Invalid assignment target."},
		{"(-x;", 2, ";", "Expect ')' after expression."},
		{"[1, 2, 3", 2, "at end", "Expect ']' after list values."},
		{"var ;", 1, ";", "Expect variable name."},
		{"var x", 1, "at end", "Expect ';' after variable declaration."},
		{"x = true", 1, "at end", "Expect ';' after value."},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		if !testParseErrors(t, p, tt.numErrs, tt.where, tt.msg) {
			t.Errorf("last error occured in test %d", i+1)
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
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

func TestUnaryExpression(t *testing.T) {
	input := `!true;
!1;
!!true;
-10;
-2.25;`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 5 {
		t.Fatalf("stmts is wrong length. expected=5, got=%d", len(stmts))
	}

	s := stmts[0].(*object.ExpressionStmt)
	ue, ok := s.Expression.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Expression)
	}

	if ue.Operator.Type != lexer.Bang {
		t.Errorf("wrong operator, expected=%s, got=%s", lexer.Bang, ue.Operator.Type)
	}
	be, ok := ue.Right.(*object.BooleanExpr)
	if !ok {
		t.Fatalf("Right value wrong type. expected=*object.BooleanExpr, got=%T", ue.Right)
	}
	if be.Value != true {
		t.Errorf("Right value incorrect type. expected=%t, got=%t", true, be.Value)
	}

	s = stmts[1].(*object.ExpressionStmt)
	ue, ok = s.Expression.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Expression)
	}

	if ue.Operator.Type != lexer.Bang {
		t.Errorf("wrong operator, expected=%s, got=%s", lexer.Bang, ue.Operator.Type)
	}
	ne, ok := ue.Right.(*object.NumberExpr)
	if !ok {
		t.Fatalf("Right value wrong type. expected=*object.NumberExpr, got=%T", ue.Right)
	}
	if ne.Int != 1 {
		t.Errorf("Right value incorrect type. expected=%d, got=%d", 1, ne.Int)
	}

	s = stmts[2].(*object.ExpressionStmt)
	ue, ok = s.Expression.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Expression)
	}

	if ue.Operator.Type != lexer.Bang {
		t.Errorf("wrong operator, expected=%s, got=%s", lexer.Bang, ue.Operator.Type)
	}
	ue2, ok := ue.Right.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("Right value incorrect type. expected=*object.UnaryExpr, got=%T. %[1]v", ue.Right)
	}
	if ue2.Operator.Type != lexer.Bang {
		t.Errorf("wrong operator, expected=%s, got=%s", lexer.Bang, ue2.Operator.Type)
	}

	be, ok = ue2.Right.(*object.BooleanExpr)
	if !ok {
		t.Fatalf("Right value wrong type. expected=*object.BooleanExpr, got=%T", ue2.Right)
	}
	if be.Value != true {
		t.Errorf("Right value incorrect type. expected=%t, got=%t", true, be.Value)
	}

	s = stmts[3].(*object.ExpressionStmt)
	ue, ok = s.Expression.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Expression)
	}

	if ue.Operator.Type != lexer.Minus {
		t.Errorf("wrong operator, expected=%s, got=%s", lexer.Minus, ue.Operator.Type)
	}
	ne, ok = ue.Right.(*object.NumberExpr)
	if !ok {
		t.Fatalf("Right value wrong type. expected=*object.NumberExpr, got=%T", ue.Right)
	}
	if ne.Int != 10 {
		t.Errorf("Right value incorrect type. expected=%d, got=%d", 10, ne.Int)
	}

	s = stmts[4].(*object.ExpressionStmt)
	ue, ok = s.Expression.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Expression)
	}

	if ue.Operator.Type != lexer.Minus {
		t.Errorf("wrong operator, expected=%s, got=%s", lexer.Minus, ue.Operator.Type)
	}
	ne, ok = ue.Right.(*object.NumberExpr)
	if !ok {
		t.Fatalf("Right value wrong type. expected=*object.NumberExpr, got=%T", ue.Right)
	}
	if ne.Float != 2.25 {
		t.Errorf("Right value incorrect type. expected=%v, got=%v", 2.25, ne.Float)
	}

}

func TestVariableExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x;", "x"},
		{"hello;", "hello"},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("test %d: wrong number of statements. expected=%d, got=%d", i, 1, len(stmts))
			continue
		}

		es := stmts[0].(*object.ExpressionStmt)
		ve, ok := es.Expression.(*object.VariableExpr)

		if !ok {
			t.Errorf("test %d: Expression wrong type. expected=*object.VariableExpr, got=%T", i, es.Expression)
			continue
		}

		if ve.Name.Lexeme != tt.expected {
			t.Errorf("test %d: Name does not match. expected=%q, got=%q", i, tt.expected, ve.Name.Lexeme)
		}
	}
}

func testLiteralExpression(t *testing.T, expr object.Expr, expected interface{}) bool {
	switch v := expected.(type) {
	case float64, int:
		return testNumberLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	case []int:
		return testListLiteral(t, expr, v)
	}

	t.Errorf("type of expr not handled. got=%T", expected)
	return false
}

func testBinaryExpression(t *testing.T, expr object.Expr, left interface{}, oper string, right interface{}) bool {
	be, ok := expr.(*object.BinaryExpr)
	if !ok {
		t.Errorf("expr is wrong type. expected=*object.BinaryExpr, got=%T", expr)
		return false
	}

	if !testLiteralExpression(t, be.Left, left) {
		return false
	}

	if be.Operator.Lexeme != oper {
		t.Errorf("Operator incorrect. expected=%q, got=%q", oper, be.Operator.Lexeme)
		return false
	}

	if !testLiteralExpression(t, be.Right, right) {
		return false
	}

	return true
}

func testLogicalExpression(t *testing.T, expr object.Expr, left interface{}, oper string, right interface{}) bool {
	be, ok := expr.(*object.LogicalExpr)
	if !ok {
		t.Errorf("expr is wrong type. expected=*object.Logical, got=%T", expr)
		return false
	}

	if !testLiteralExpression(t, be.Left, left) {
		return false
	}

	if be.Operator.Lexeme != oper {
		t.Errorf("Operator incorrect. expected=%q, got=%q", oper, be.Operator.Lexeme)
		return false
	}

	if !testLiteralExpression(t, be.Right, right) {
		return false
	}

	return true
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

func testIdentifier(t *testing.T, expr object.Expr, value string) bool {
	ident, ok := expr.(*object.VariableExpr)
	if !ok {
		t.Errorf("expr wrong type. expected=*object.VariableExpr, got=%T", expr)
		return false
	}

	if ident.Name.Lexeme != value {
		t.Errorf("name wrong value. expected=%q, got=%q", value, ident.Name.Lexeme)
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

func testListLiteral(t *testing.T, expr object.Expr, values []int) bool {
	le, ok := expr.(*object.ListExpr)
	if !ok {
		t.Errorf("expr not correct type. expected=*object.ListExpr, got=%T", expr)
		return false
	}

	if len(values) != len(le.Values) {
		t.Errorf("list does not contain correct number of values. expected=%d, got=%d", len(values), len(le.Values))
		return false
	}

	for i, v := range values {
		ne, ok := le.Values[i].(*object.NumberExpr)
		if !ok {
			t.Errorf("Value unexpected type. expected=*object.NumberExpr, got=%T", le.Values[i])
			return false
		}

		if ne.Int != v {
			t.Errorf("Value is incorrect. expected=%d, got=%d (%v)", v, ne.Int, ne.Float)
			return false
		}
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

func testVariable(t *testing.T, stmt object.Stmt, ident string) bool {
	s, ok := stmt.(*object.VarStmt)
	if !ok {
		t.Errorf("Stmt wrong type, expected=*object.VarStmt, got=%T", stmt)
		return false
	}

	if s.Name.Lexeme != ident {
		t.Errorf("Name is wrong value. expected=%q, got=%q", ident, s.Name.Lexeme)
		return false
	}

	return true
}

func testParseErrors(t *testing.T, p *Parser, numErrs int, where, msg string) bool {
	stmts := p.Parse()

	if len(stmts) != 0 {
		t.Errorf("wrong number of statements. expected=0, got=%d", len(stmts))
		return false
	}

	errs := p.Errors()
	// 2 because of Unterminated string, then missing semicolon
	if len(errs) != numErrs {
		t.Errorf("wrong number of errors. expected=1, got=%d", len(errs))
		return false
	}

	e := errs[0]
	if e.Line != 1 {
		t.Errorf("error on wrong line, expected=1, got=%d", e.Line)
		return false
	}

	if e.Where != where {
		t.Errorf("error at wrong location. expected=%q, got=%q", where, e.Where)
		return false
	}

	if e.Msg != msg {
		t.Errorf("wrong error message. expected=%q, got=%q", msg, e.Msg)
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
		t.Errorf("parser error on line: %d at '%s': %s", msg.Line, msg.Where, msg.Msg)
	}
	t.FailNow()
}
