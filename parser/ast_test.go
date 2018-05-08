package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
)

type printerObj struct {
	value string
}

func (po printerObj) Type() object.Type { return object.Printer }
func (po printerObj) String() string    { return po.value }

type AstPrinter struct{}

func (p *AstPrinter) Print(expr object.Expr) string {
	po, _ := expr.Accept(p)
	return po.String()
}

func (p *AstPrinter) VisitAssignExpr(expr *object.AssignExpr) (object.Object, error) {
	return p.parenthesize("= "+expr.Name.Lexeme, expr.Value), nil
}

func (p *AstPrinter) VisitBinaryExpr(expr *object.BinaryExpr) (object.Object, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p *AstPrinter) VisitCallExpr(expr *object.CallExpr) (object.Object, error) {
	all := []object.Expr{expr.Callee}
	all = append(all, expr.Args...)
	return p.parenthesize("call", all...), nil
}

func (p *AstPrinter) VisitNumberExpr(expr *object.NumberExpr) (object.Object, error) {
	var s string
	if expr.Token.Type == lexer.NumberF {
		s = fmt.Sprintf("%.2f", expr.Float)
	} else {
		s = fmt.Sprintf("%d", expr.Int)
	}
	return printerObj{value: s}, nil
}

func (p *AstPrinter) VisitUnaryExpr(expr *object.UnaryExpr) (object.Object, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (p *AstPrinter) VisitBooleanExpr(expr *object.BooleanExpr) (object.Object, error) {
	return printerObj{value: fmt.Sprintf("%t", expr.Value)}, nil
}

func (p *AstPrinter) VisitGetExpr(expr *object.GetExpr) (object.Object, error) {
	return p.parenthesize("."+expr.Name.Lexeme, expr.Object), nil
}

func (p *AstPrinter) VisitGroupingExpr(expr *object.GroupingExpr) (object.Object, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p *AstPrinter) VisitIndexExpr(expr *object.IndexExpr) (object.Object, error) {
	return p.parenthesize("[]", expr.Left, expr.Right), nil
}

func (p *AstPrinter) VisitListExpr(expr *object.ListExpr) (object.Object, error) {
	var b bytes.Buffer

	b.WriteByte('[')
	if len(expr.Values) > 0 {
		b.WriteString(fmt.Sprintf("%v", expr.Values[0]))
		if len(expr.Values) > 1 {
			for i := 1; i < len(expr.Values); i++ {
				b.WriteString(fmt.Sprintf(", %v", expr.Values[i]))
			}
		}
	}

	b.WriteByte(']')
	return printerObj{value: b.String()}, nil
}
func (p *AstPrinter) VisitLogicalExpr(expr *object.LogicalExpr) (object.Object, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p *AstPrinter) VisitNullExpr(expr *object.NullExpr) (object.Object, error) {
	return printerObj{"null"}, nil
}

func (p *AstPrinter) VisitSetExpr(expr *object.SetExpr) (object.Object, error) {
	return p.parenthesize(expr.Name.Lexeme, expr.Object, expr.Value), nil
}

func (p *AstPrinter) VisitStringExpr(expr *object.StringExpr) (object.Object, error) {
	return printerObj{expr.Value}, nil
}

func (p *AstPrinter) VisitSuperExpr(expr *object.SuperExpr) (object.Object, error) {
	return printerObj{"super." + expr.Method.Lexeme}, nil
}

func (p *AstPrinter) VisitThisExpr(expr *object.ThisExpr) (object.Object, error) {
	return printerObj{value: "this"}, nil
}

func (p *AstPrinter) VisitVariableExpr(expr *object.VariableExpr) (object.Object, error) {
	return printerObj{value: expr.Name.Lexeme}, nil
}

func (p *AstPrinter) parenthesize(name string, exprs ...object.Expr) printerObj {
	var b bytes.Buffer

	b.WriteByte('(')
	b.WriteString(name)

	for _, e := range exprs {
		b.WriteByte(' ')
		s, _ := e.Accept(p)
		b.WriteString(s.String())
	}

	b.WriteByte(')')
	return printerObj{value: b.String()}
}

func TestASTGrouping(t *testing.T) {
	printer := &AstPrinter{}

	tests := []struct {
		input  string
		expect string
	}{
		{"4 + 5;", "(+ 4 5)"},
		{"4 + 5 + 6;", "(+ (+ 4 5) 6)"},
		{"-a * b;", "(* (- a) b)"},
		{"a + b * c;", "(+ a (* b c))"},
		{"a * b + c;", "(+ (* a b) c)"},
		{"(a + b) * c;", "(* (group (+ a b)) c)"},
		{"a + b * c + d / e - f;", "(- (+ (+ a (* b c)) (/ d e)) f)"},
		{"5 == a;", "(== 5 a)"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5;", "(== (+ 3 (* 4 5)) (+ (* 3 1) (* 4 5)))"},
		{"a += b;", "(= a (+ a b))"},
		{"true or false;", "(or true false)"},
		{"a[b + c];", "([] a (+ b c))"},
		{"a.b;", "(.b a)"},
	}

	for i, tt := range tests {
		l := lexer.New([]byte(tt.input), "testfile.gpc")
		p := New(l)
		stmts := p.Parse()

		if len(p.Errors()) != 0 {
			t.Errorf("test %d Parser encountered errors: %v", i+1, p.Errors())
			t.Fatalf("Test line was: %s", tt.input)
		}

		s := stmts[0].(*object.ExpressionStmt)
		out := printer.Print(s.Expression)

		if out != tt.expect {
			t.Errorf("test %d: output incorrect. expected=%q, got=%q", i+1, tt.expect, out)
		}
	}
}
