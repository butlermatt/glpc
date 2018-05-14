package parser

import (
	"testing"

	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
)

func TestBlockStatement(t *testing.T) {
	input := `fn test() {
{
  var x = 1;
  x += 1;
  if (x == 2) {
    count + 2;
    count - 2;
  } else {
    error = "Oh oh";
    error += "broken";
  }
}
}`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("unexpected number of program statements. expected=%d, got=%d", 1, len(stmts))
	}

	fn := stmts[0].(*object.FunctionStmt)
	if len(fn.Body) != 1 {
		t.Fatalf("unexpected number of program statements. expected=%d, got=%d", 1, len(fn.Body))
	}

	bs, ok := fn.Body[0].(*object.BlockStmt)
	if !ok {
		t.Fatalf("unexpected statement type. expected=*object.BlockStmt, got=%T", stmts[0])
	}

	if len(bs.Statements) != 3 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 3, len(bs.Statements))
	}

	ifstmt, ok := bs.Statements[2].(*object.IfStmt)
	if !ok {
		t.Fatalf("3rd statement wrong type. expected=*object.IfStmt, got=%T", bs.Statements[2])
	}

	bs, ok = ifstmt.Then.(*object.BlockStmt)
	if !ok {
		t.Fatalf("incorrect statement type. expected=*object.BlockStmt, got=%T", ifstmt.Then)
	}

	if len(bs.Statements) != 2 {
		t.Errorf("incorrect number of statements. expected=%d, got=%d", 2, len(bs.Statements))
	}

	bs, ok = ifstmt.Else.(*object.BlockStmt)
	if !ok {
		t.Fatalf("incorrect statement type. expected=*object.BlockStmt, got=%T", ifstmt.Then)
	}

	if len(bs.Statements) != 2 {
		t.Errorf("incorrect number of statements. expected=%d, got=%d", 2, len(bs.Statements))
	}
}

func TestBreakStatement(t *testing.T) {
	tests := []string{
		`fn test() { for (var i = 0; i < 10; i += 1) { break; }}`,
		`fn test() { while (true) { break; } }`,
		`fn test() { do { break; } while(true); }`,
	}

	for _, tt := range tests {
		l := lexer.New([]byte(tt), "testfile.gpc")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
		}

		fn := stmts[0].(*object.FunctionStmt)

		if len(fn.Body) != 1 {
			t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
		}

		s, ok := fn.Body[0].(*object.ForStmt)
		if !ok {
			t.Fatalf("Statement wrong type. expected=*object.ForStmt, got=%T", stmts[0])
		}
		bl, ok := s.Body.(*object.BlockStmt)
		if !ok {
			t.Fatalf("Body wrong type, expected=*object.BlockStmt, got=%T", s.Body)
		}

		if len(bl.Statements) != 1 {
			t.Fatalf("wrong number of statements. expected=%d, got=%d", 1, len(bl.Statements))
		}

		br, ok := bl.Statements[0].(*object.BreakStmt)
		if !ok {
			t.Fatalf("Body wrong type. expected=*object.BreakStmt, got=%T", s.Body)
		}

		if br.Keyword.Lexeme != "break" {
			t.Fatalf("lexeme wrong value. expected=%q, got=%q", "break", br.Keyword.Lexeme)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `fn test() { add(1, 2 * 3, 4 + 5); }`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("wrong number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fn := stmts[0].(*object.FunctionStmt)

	if len(fn.Body) != 1 {
		t.Fatalf("wrong number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	st, ok := fn.Body[0].(*object.ExpressionStmt)
	if !ok {
		t.Fatalf("statement is wrong type. expected=*object.ExpressionStmt, got=%T", stmts[0])
	}

	ce, ok := st.Expression.(*object.CallExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.CallExpr, got=%T", st.Expression)
	}

	testIdentifier(t, ce.Callee, "add")

	if len(ce.Args) != 3 {
		t.Fatalf("wrong number of arguments. expected=%d, got=%d", 3, len(ce.Args))
	}

	testLiteralExpression(t, ce.Args[0], 1)
	testBinaryExpression(t, ce.Args[1], 2, "*", 3)
	testBinaryExpression(t, ce.Args[2], 4, "+", 5)
}

func TestClassStatement(t *testing.T) {
	input := `class test : origin { add(x, y) { return x + y; } }`
	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	cs, ok := stmts[0].(*object.ClassStmt)
	if !ok {
		t.Fatalf("statement wrong type. expected=*object.ClassStmt, got=%T", stmts[0])
	}

	if cs.Name.Lexeme != "test" {
		t.Errorf("class name wrong. expected=%q, got=%q", "test", cs.Name.Lexeme)
	}

	if cs.Super.Name.Lexeme != "origin" {
		t.Errorf("Superclass name wrong. expected=%q, got=%q", "origin", cs.Super.Name.Lexeme)
	}

	if len(cs.Methods) != 1 {
		t.Fatalf("wrong number of methods. expected=%d, got=%d", 1, len(cs.Methods))
	}

	fn := cs.Methods[0]
	if fn.Name.Lexeme != "add" {
		t.Errorf("wrong method name. expected=%q, got=%q", "add", fn.Name.Lexeme)
	}

	if len(fn.Parameters) != 2 {
		t.Errorf("wrong number of method parameters. expected=%d, got=%d", 2, len(fn.Parameters))
	}
}

func TestContinueStatement(t *testing.T) {
	input := `fn test() { while (true) continue; }`
	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fn := stmts[0].(*object.FunctionStmt)

	if len(fn.Body) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	s, ok := fn.Body[0].(*object.ForStmt)
	if !ok {
		t.Fatalf("Statement wrong type. expected=*object.ForStmt, got=%T", stmts[0])
	}

	br, ok := s.Body.(*object.ContinueStmt)
	if !ok {
		t.Fatalf("Body wrong type. expected=*object.BreakStmt, got=%T", s.Body)
	}

	if br.Keyword.Lexeme != "continue" {
		t.Fatalf("lexeme wrong value. expected=%q, got=%q", "continue", br.Keyword.Lexeme)
	}
}

func TestIfStatement(t *testing.T) {
	input := `fn test () { if (a == 1) 
  x = true;
else
  y = false;
}`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fn := stmts[0].(*object.FunctionStmt)

	if len(fn.Body) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	ifStmt, ok := fn.Body[0].(*object.IfStmt)
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

func TestFunctions(t *testing.T) {
	input := `fn test(x, y) { var temp = x; x = y; y = temp; }`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fn, ok := stmts[0].(*object.FunctionStmt)
	if !ok {
		t.Fatalf("statement wrong type. expected=*object.FunctionStmt, got=%T", stmts[0])
	}

	if fn.Name.Lexeme != "test" {
		t.Errorf("wrong function name. expected=%q, got=%q", "test", fn.Name.Lexeme)
	}

	if len(fn.Parameters) != 2 {
		t.Errorf("wrong number of parameters. expected=%d, got=%d", 2, len(fn.Parameters))
	}

	if fn.Parameters[0].Lexeme != "x" {
		t.Errorf("wrong parameter name. expected=%q, got=%q", "x", fn.Parameters[0].Lexeme)
	}

	if fn.Parameters[1].Lexeme != "y" {
		t.Errorf("wrong parameter name. expected=%q, got=%q", "y", fn.Parameters[1].Lexeme)
	}

	if len(fn.Body) != 3 {
		t.Errorf("wrong number of body statements. expected=%d, got=%d", 3, len(fn.Body))
	}

	testVariable(t, fn.Body[0], "temp")
	val := fn.Body[0].(*object.VarStmt)
	testLiteralExpression(t, val.Value, "x")

	es, ok := fn.Body[1].(*object.ExpressionStmt)
	if !ok {
		t.Fatalf("wrong type for statement 2. expected=*object.ExpressionStmt, got=%T", fn.Body[1])
	}

	ae, ok := es.Expression.(*object.AssignExpr)
	if !ok {
		t.Fatalf("expression is wrong type. expected=*object.AssignExpr, got=%T", es.Expression)
	}

	if ae.Name.Lexeme != "x" {
		t.Errorf("name does not match. expected=%q, got=%q", "x", ae.Name.Lexeme)
	}

	testLiteralExpression(t, ae.Value, "y")

	es, ok = fn.Body[2].(*object.ExpressionStmt)
	if !ok {
		t.Fatalf("wrong type for statement 2. expected=*object.ExpressionStmt, got=%T", fn.Body[2])
	}

	ae, ok = es.Expression.(*object.AssignExpr)
	if !ok {
		t.Fatalf("expression is wrong type. expected=*object.AssignExpr, got=%T", es.Expression)
	}

	if ae.Name.Lexeme != "y" {
		t.Errorf("name does not match. expected=%q, got=%q", "y", ae.Name.Lexeme)
	}

	testLiteralExpression(t, ae.Value, "temp")
}

func TestForStatement(t *testing.T) {
	input := `fn test() { for (var i = 0; i < 10; i += 1) {
    error += "Hello";
    Another = good;
  }
}`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fn := stmts[0].(*object.FunctionStmt)

	if len(fn.Body) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fs, ok := fn.Body[0].(*object.ForStmt)
	if !ok {
		t.Fatalf("statement wrong type. expected=*object.ForStmt, got=%T", stmts[0])
	}

	if fs.Keyword.Type != lexer.For {
		t.Errorf("for statement token wrong type. expected=%q, got=%q", lexer.For, fs.Keyword.Type)
	}

	is, ok := fs.Initializer.(*object.VarStmt)
	if !ok {
		t.Fatalf("Initializer incorrect type. expected=*object.VarStmt, got=%T", fs.Initializer)
	}

	testVariable(t, is, "i")
	testNumberLiteral(t, is.Value, 0)
	testBinaryExpression(t, fs.Condition, "i", "<", 10)

	bl, ok := fs.Body.(*object.BlockStmt)
	if !ok {
		t.Fatalf("body wrong type. expected=*object.BlockStmt, got=%T", fs.Body)
	}

	if len(bl.Statements) != 2 {
		t.Errorf("wrong number of statements in body. expected=%d, got=%d", 2, len(bl.Statements))
	}

	ae, ok := fs.Increment.(*object.AssignExpr)
	if !ok {
		t.Fatalf("Increment wrong type. expected=*object.AssignExpr, got=%T", fs.Increment)
	}

	if ae.Name.Lexeme != "i" {
		t.Errorf("initializer name incorrect. expected=%q, got=%q", "i", ae.Name.Lexeme)
	}

	testBinaryExpression(t, ae.Value, "i", "+", 1)
}

func TestDoWhileStatement(t *testing.T) {
	input := `fn test() { do {
  error += "Hello";
  Another = good;
} while(i < 10); }`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fn := stmts[0].(*object.FunctionStmt)

	if len(fn.Body) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fs, ok := fn.Body[0].(*object.ForStmt)
	if !ok {
		t.Fatalf("statement wrong type. expected=*object.ForStmt, got=%T", stmts[0])
	}

	if fs.Keyword.Type != lexer.Do {
		t.Errorf("for statement token wrong type. expected=%q, got=%q", lexer.Do, fs.Keyword.Type)
	}

	if fs.Initializer != nil {
		t.Errorf("Inititalizer wrong value. expected=<nil> got=%T", fs.Initializer)
	}

	testBinaryExpression(t, fs.Condition, "i", "<", 10)

	bl, ok := fs.Body.(*object.BlockStmt)
	if !ok {
		t.Fatalf("body wrong type. expected=*object.BlockStmt, got=%T", fs.Body)
	}

	if len(bl.Statements) != 2 {
		t.Errorf("wrong number of statements in body. expected=%d, got=%d", 2, len(bl.Statements))
	}

	if fs.Increment != nil {
		t.Errorf("Increment wrong type. expected=<nil>, got=%T", fs.Increment)
	}
}

func TestWhileStatement(t *testing.T) {
	input := `fn test() { while(i < 10) {
    error += "Hello";
    Another = good;
  }
}`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fn := stmts[0].(*object.FunctionStmt)

	if len(fn.Body) != 1 {
		t.Fatalf("incorrect number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	fs, ok := fn.Body[0].(*object.ForStmt)
	if !ok {
		t.Fatalf("statement wrong type. expected=*object.ForStmt, got=%T", stmts[0])
	}

	if fs.Keyword.Type != lexer.While {
		t.Errorf("for statement token wrong type. expected=%q, got=%q", lexer.While, fs.Keyword.Type)
	}

	if fs.Initializer != nil {
		t.Errorf("Inititalizer wrong value. expected=<nil> got=%T", fs.Initializer)
	}

	testBinaryExpression(t, fs.Condition, "i", "<", 10)

	bl, ok := fs.Body.(*object.BlockStmt)
	if !ok {
		t.Fatalf("body wrong type. expected=*object.BlockStmt, got=%T", fs.Body)
	}

	if len(bl.Statements) != 2 {
		t.Errorf("wrong number of statements in body. expected=%d, got=%d", 2, len(bl.Statements))
	}

	if fs.Increment != nil {
		t.Errorf("Increment wrong type. expected=<nil>, got=%T", fs.Increment)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{"fn x() { return; }", nil},
		{"fn x() { return true; }", true},
		{"fn x() { return 42; }", 42},
		{"fn x() { return y; }", "y"},
		{"fn x() { return null; }", nil},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gcp")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		fn, ok := stmts[0].(*object.FunctionStmt)
		if !ok {
			t.Errorf("test %d: unexpected statement. expected=*object.FunctionStatement, got=%T", i, stmts[0])
			continue
		}

		if fn.Name.Lexeme != "x" {
			t.Errorf("test %d: unexpected name. expected=%q, got=%q", i, "x", fn.Name.Lexeme)
		}

		if len(fn.Body) != 1 {
			t.Errorf("test %d: wrong number of body statements. expected=1, got=%d", i, len(fn.Body))
			continue
		}

		rt, ok := fn.Body[0].(*object.ReturnStmt)
		if !ok {
			t.Errorf("test %d: wrong statement type. expected=*object.ReturnStmt, got=%T", i, fn.Body[0])
			continue
		}

		testLiteralExpression(t, rt.Value, tt.value)
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
		stmts, _ := p.Parse()
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
		{"fn test() { x = 5; }", "x", 5},
		{"fn test() { y = 2.25; }", "y", 2.25},
		{"fn test() { test = true; }", "test", true},
		{"fn test() { list = [1, 2, 3]; }", "list", []int{1, 2, 3}},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: incorrect number of statements. expected=%d, got=%d", i, 1, len(stmts))
		}

		fn := stmts[0].(*object.FunctionStmt)

		if len(fn.Body) != 1 {
			t.Fatalf("test %d: incorrect number of statements. expected=%d, got=%d", i, 1, len(stmts))
		}

		s := fn.Body[0].(*object.ExpressionStmt)
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
		{"fn test() { x += 1; }", "x", bin{"x", "+", 1}},
		{"fn test() { y -= 2.2; }", "y", bin{"y", "-", 2.2}},
		{"fn test() { test *= 3; }", "test", bin{"test", "*", 3}},
		{"fn test() { z /= 4; }", "z", bin{"z", "/", 4}},
		{"fn test() { a ~/= 5; }", "a", bin{"a", "~/", 5}},
		{"fn test() { b %= 6; }", "b", bin{"b", "%", 6}},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("test %d: incorrect statements length. expected=%d, got=%d", i+1, 1, len(stmts))
			continue
		}

		fn := stmts[0].(*object.FunctionStmt)

		if len(fn.Body) != 1 {
			t.Errorf("test %d: incorrect statements length. expected=%d, got=%d", i+1, 1, len(stmts))
			continue
		}

		s, ok := fn.Body[0].(*object.ExpressionStmt)
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
		{"var x = 5 + 4;", 5, "+", 4},
		{"var x = 5 - 4;", 5, "-", 4},
		{"var x = 5 * 4;", 5, "*", 4},
		{"var x = 5 / 4;", 5, "/", 4},
		{"var x = 5 % 4;", 5, "%", 4},
		{"var x = 5 ~/ 4;", 5, "~/", 4},
		{"var x = 5 < 4;", 5, "<", 4},
		{"var x = 5 > 4;", 5, ">", 4},
		{"var x = 5 <= 4;", 5, "<=", 4},
		{"var x = 5.5 == 4.4;", 5.5, "==", 4.4},
		{"var x = 5 != 4;", 5, "!=", 4},
		{"var x = true == true;", true, "==", true},
		{"var x = true != false;", true, "!=", false},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("on test %d: incorrect number of statements. expected=1, got=%d", i+1, len(stmts))
			continue
		}

		s, ok := stmts[0].(*object.VarStmt)
		if !ok {
			t.Errorf("on test %d: Statement wrong type. expected=*object.VarStmt, got=%T", i+1, stmts[0])
			continue
		}

		testBinaryExpression(t, s.Value, tt.left, tt.oper, tt.right)
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"var x = true;", true},
		{"var x = false;", false},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)

		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: incorrect number of statements. expected=%d, got=%d", i, 1, len(stmts))
		}

		s, ok := stmts[0].(*object.VarStmt)
		if !ok {
			t.Fatalf("test %d: statement wrong type. expected=*object.VarStmt, got=%T", i, stmts[0])
		}

		if !testBooleanLiteral(t, s.Value, tt.value) {
			t.Errorf("last error occurred on test %d", i)
		}
	}
}

func TestGetExpression(t *testing.T) {
	input := `var y = test.x;`
	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("wrong number of statements. expected=%d, got=%d", 1, len(stmts))
	}

	es, ok := stmts[0].(*object.VarStmt)
	if !ok {
		t.Fatalf("statement wrong type. expected=*object.ExpressionStmt, got=%T", stmts[0])
	}

	ge, ok := es.Value.(*object.GetExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.GetExpr, got=%T", es.Value)
	}

	if ge.Name.Lexeme != "x" {
		t.Errorf("wrong name. expected=%q, got=%q", "x", ge.Name.Lexeme)
	}

	ve, ok := ge.Object.(*object.VariableExpr)
	if !ok {
		t.Fatalf("wrong object type. expected=*object.VariableExpr, got=%T", ge.Object)
	}

	if ve.Name.Lexeme != "test" {
		t.Errorf("wrong object name. expected=%q, got=%q", "test", ve.Name.Lexeme)
	}
}

func TestGroupingExpression(t *testing.T) {
	input := "var x = (5);"

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("statements incorrect length. expected=%d, got=%d", 1, len(stmts))
	}

	stmt := stmts[0].(*object.VarStmt)

	gr, ok := stmt.Value.(*object.GroupingExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.GroupingExpr, got=%T", stmt.Value)
	}

	testNumberLiteral(t, gr.Expression, 5)

}

func TestIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		left  interface{}
		right interface{}
	}{
		{"var x = a[5];", "a", 5},
		{"var x = some[thing];", "some", "thing"},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("on test %d: incorrect number of statements. expected=1, got=%d", i+1, len(stmts))
			continue
		}

		s, ok := stmts[0].(*object.VarStmt)
		if !ok {
			t.Errorf("on test %d: Statement wrong type. expected=*object.VarStmt, got=%T", i+1, stmts[0])
			continue
		}

		ie, ok := s.Value.(*object.IndexExpr)
		if !ok {
			t.Errorf("on test %d: Expression wrong type. expected=*object.IndexExpr, got=%T", i+1, s.Value)
			continue
		}

		testLiteralExpression(t, ie.Left, tt.left)
		testLiteralExpression(t, ie.Right, tt.right)
	}
}

func TestListExpression(t *testing.T) {
	input := "var x = [0, 'one', true];"

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("statements is incorrect length. expected=%d, got=%d", 1, len(stmts))
	}

	stmt := stmts[0].(*object.VarStmt)

	list, ok := stmt.Value.(*object.ListExpr)
	if !ok {
		t.Fatalf("expr is wrong type. expected=*object.ListExpr, got=%T", stmt.Value)
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
		{"var x = true and true;", true, "and", true},
		{"var x = true or false;", true, "or", false},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("on test %d: incorrect number of statements. expected=1, got=%d", i+1, len(stmts))
			continue
		}

		s, ok := stmts[0].(*object.VarStmt)
		if !ok {
			t.Errorf("on test %d: Statement wrong type. expected=*object.VarStmt, got=%T", i+1, stmts[0])
			continue
		}

		testLogicalExpression(t, s.Value, tt.left, tt.oper, tt.right)
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{"var x = 5;", 5},
		{"var x = 10;", 10},
		{"var x = 123.456;", 123.456},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)

		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: incorrect number of statements. exepected=%d, got=%d", i, 1, len(stmts))
		}

		s, ok := stmts[0].(*object.VarStmt)
		if !ok {
			t.Fatalf("test %d: statement[0] wrong type. expected=*object.VarStmt, got=%T", i, stmts[0])
		}

		if !testNumberLiteral(t, s.Value, tt.value) {
			t.Errorf("last test tha failed was %d", i)
		}
	}
}

func TestNullLiteralExpression(t *testing.T) {
	input := "var x = null;"

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)

	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 1 {
		t.Fatalf("incorrect number of statements. exepected=%d, got=%d", 1, len(stmts))
	}

	s, ok := stmts[0].(*object.VarStmt)
	if !ok {
		t.Fatalf("statement[0] wrong type. expected=*object.ExpressionStmt, got=%T", stmts[0])
	}

	testNullLiteral(t, s.Value)
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input   string
		numErrs int
		where   string
		msg     string
	}{
		{`fn test() { "hello world; }`, 3, `"hello world; }`, "Unterminated string."},
		{"fn test() { 7 = x; }", 2, "=", "Invalid assignment target."},
		{"fn test() { (-x; }", 3, ";", "Expect ')' after expression."},
		{"fn test() { [1, 2, 3 }", 3, "}", "Expect ']' after list values."},
		{"var ;", 1, ";", "Expect variable name."},
		{"var x", 1, "at end", "Expect ';' after variable declaration."},
		{"fn test() { x = true }", 2, "}", "Expect ';' after value."},
		{"fn test() { x[2 }", 3, "}", "Expect ']' after index."},
		{"fn test() { if y x = 7; }", 1, "y", "Expect '(' after 'if'."},
		{"fn test() { x = 2;", 1, "at end", "Expect '}' after block."}, //10

		{"fn test() { for x = 2; }", 1, "x", "Expect '(' after 'for'."},
		{"fn test() { for (;) }", 3, ")", "Expect expression."},
		{"fn test() { for (; x < 2) }", 2, ")", "Expect ';' after loop condition."},
		{"fn test() { for (; x < 2; }", 3, "}", "Expect expression."},
		{"fn test() { for (; x < 2; x += 2 {} }", 2, "{", "Expect ')' after for clauses."},
		{"fn test() { while true { } }", 2, "true", "Expect '(' after 'while'."},
		{"fn test() { while (true { } }", 2, "{", "Expect ')' after while condition."},
		{"fn test() { do {} ; }", 1, ";", "Expect 'while' after do-while body."},
		{"fn test() { do {} while; }", 1, ";", "Expect '(' after 'while'."},
		{"fn test() { do {} while(x == 2; }", 1, ";", "Expect ')' after while condition."}, //20

		{"fn test() { do {} while(x == 2) }", 2, "}", "Expect ';' after ')'."},
		{"fn test() { while(true) { break } }", 3, "}", "Expect ';' after 'break'."},
		{"fn test() { while(true) { continue } }", 3, "}", "Expect ';' after 'continue'."},
		{"fn test() { if (i == 5) { break; } }", 3, "break", "Cannot use 'break' outside of a loop."},
		{"fn test() { if (i == 5) { continue; } }", 3, "continue", "Cannot use 'continue' outside of a loop."},
		{"fn(x, y) {}", 1, "(", "Expect function name."},
		{"fn test {}", 1, "{", "Expect '(' after function name."},
		{"fn test(7){}", 1, "7", "Expect parameter name."},
		{"fn test(x, ){}", 1, ")", "Expect parameter name."},
		{"fn test(a, b {}", 1, "{", "Expect ')' after parameters."}, // 30

		{"fn test(a, b) x = 10; }", 2, "x", "Expect '{' before function body."},
		{"fn test() { return 1 }", 2, "}", "Expect ';' after return value."},
		{"fn test() { return true }", 2, "}", "Expect ';' after return value."},
		{"return;", 1, "return", "Only classes, functions and variables may be used in top-level."},
		{"class { }", 1, "{", "Expect class name."},
		{"class test : { }", 1, "{", "Expect superclass name."},
		{"class test }", 1, "}", "Expect '{' before class body."},
		{"class test {", 1, "at end", "Expect '}' after class body."},
		{"fn test() { this.x = true; }", 2, "this", "Cannot use 'this' outside of a class."},
		{"fn test() { super = true; }", 3, "=", "Expect '.' after 'super'."}, // 40

		{"fn test() { super. = true; }", 3, "=", "Expect superclass method name."},
		{"fn test() { super.x = true; }", 3, "super", "Cannot use 'super' outside of a class."},
		{"fn x() { var i = i + 1; }", 1, "i", "Cannot read local variable in its own initializer."},
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
		{`var x = "hello world";`, "hello world"},
		{`var x = 'hello world';`, "hello world"},
		{"var x = `hello world`;", "hello world"},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Fatalf("test %d: incorrect number of statements. expected=%d, got=%d", i, 1, len(stmts))
		}

		stmt := stmts[0].(*object.VarStmt)
		lit, ok := stmt.Value.(*object.StringExpr)
		if !ok {
			t.Fatalf("test %d: expression wrong type. expeted=*object.StringExpr, got=%T", i, stmt.Value)
		}

		if lit.Value != tt.expected {
			t.Errorf("test %d: value is wrong. expected=%q, got==%q", i, tt.expected, lit.Value)
		}
	}
}

func TestUnaryExpression(t *testing.T) {
	input := `var x = !true;
var x = !1;
var x = !!true;
var x = -10;
var x = -2.25;`

	l := lexer.New([]byte(input), "testfile.gpc")
	p := New(l)
	stmts, _ := p.Parse()
	checkParseErrors(t, p)

	if len(stmts) != 5 {
		t.Fatalf("stmts is wrong length. expected=5, got=%d", len(stmts))
	}

	s := stmts[0].(*object.VarStmt)
	ue, ok := s.Value.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Value)
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

	s = stmts[1].(*object.VarStmt)
	ue, ok = s.Value.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Value)
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

	s = stmts[2].(*object.VarStmt)
	ue, ok = s.Value.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Value)
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

	s = stmts[3].(*object.VarStmt)
	ue, ok = s.Value.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Value)
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

	s = stmts[4].(*object.VarStmt)
	ue, ok = s.Value.(*object.UnaryExpr)
	if !ok {
		t.Fatalf("expression wrong type. expected=*object.UnaryExpr, got=%T", s.Value)
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
		{"fn test() { x; }", "x"},
		{"fn test() { hello; }", "hello"},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts, _ := p.Parse()
		checkParseErrors(t, p)

		if len(stmts) != 1 {
			t.Errorf("test %d: wrong number of statements. expected=%d, got=%d", i, 1, len(stmts))
			continue
		}

		fn := stmts[0].(*object.FunctionStmt)

		if len(fn.Body) != 1 {
			t.Errorf("test %d: wrong number of statements. expected=%d, got=%d", i, 1, len(stmts))
			continue
		}

		es := fn.Body[0].(*object.ExpressionStmt)
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
	if expected == nil {
		return testNullLiteral(t, expr)
	}

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

func testNullLiteral(t *testing.T, expr object.Expr) bool {
	var value interface{}

	ne, ok := expr.(*object.NullExpr)
	if !ok {
		t.Errorf("expr not correct type. expected=*object.NullExpr, got=%T", expr)
		return false
	}

	if ne.Value != value {
		t.Errorf("null value incorrect. expected=%v, got=%v", value, ne.Value)
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
	stmts, _ := p.Parse()

	if len(stmts) > 1 {
		t.Errorf("wrong number of statements. expected=0, got=%d", len(stmts))
		return false
	}

	errs := p.Errors()
	// 2 because of Unterminated string, then missing semicolon
	if len(errs) != numErrs {
		t.Errorf("wrong number of errors. expected=%d, got=%d", numErrs, len(errs))
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
