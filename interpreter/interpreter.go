package interpreter

import (
	"errors"
	"fmt"
	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
)

var BreakError = errors.New("unexpected 'break' outside of loop")
var ContinueError = errors.New("unexpected 'continue' outside of loop")

type ReturnError struct {
	object.RuntimeError
	Value object.Object
}

func (re *ReturnError) Error() string {
	return fmt.Sprintf("[Return error] at %s - %s", re.RuntimeError.Token.Lexeme, re.Value.String())
}

func NewReturnValue(keyword *lexer.Token, value object.Object) *ReturnError {
	return &ReturnError{RuntimeError: object.RuntimeError{Token: keyword}, Value: value}
}

type Interpreter struct {
	file    string
	stmts   []object.Stmt
	local   map[object.Expr]int
	env     *object.Environment
	globals *object.Environment
}

func New(statements []object.Stmt, depthMap map[object.Expr]int, file string) *Interpreter {
	env := object.NewEnvironment(file)
	glob := object.GetGlobal()
	return &Interpreter{file: file, stmts: statements, local: depthMap, env: env, globals: glob}
}

func (inter *Interpreter) Interpret() error {
	for _, stmt := range inter.stmts {
		err := inter.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (inter *Interpreter) execute(stmt object.Stmt) error {
	return stmt.Accept(inter)
}

func (inter *Interpreter) evaluate(expr object.Expr) (object.Object, error) {
	return expr.Accept(inter)
}

func (inter *Interpreter) executeBlock(stmts []object.Stmt, env *object.Environment) error {
	prevEnv := inter.env
	inter.env = env

	var err error
	for _, stmt := range stmts {
		err := inter.execute(stmt)
		if err != nil {
			break
		}
	}

	inter.env = prevEnv
	return err
}

func isTruthy(obj object.Object) bool {
	if obj == nil || obj.Type() == object.Null {
		return false
	}

	if obj.Type() == object.Boolean {
		return obj.(*Boolean).Value
	}

	return true
}

func (inter *Interpreter) VisitBlockStmt(stmt *object.BlockStmt) error {
	return inter.executeBlock(stmt.Statements, object.NewEnclosedEnvironment(inter.env))
}

func (inter *Interpreter) VisitBreakStmt(stmt *object.BreakStmt) error       { return BreakError }
func (inter *Interpreter) VisitClassStmt(stmt *object.ClassStmt) error       { return nil }
func (inter *Interpreter) VisitContinueStmt(stmt *object.ContinueStmt) error { return ContinueError }

func (inter *Interpreter) VisitExpressionStmt(stmt *object.ExpressionStmt) error {
	_, err := inter.evaluate(stmt.Expression)
	return err
}

func (inter *Interpreter) VisitFunctionStmt(stmt *object.FunctionStmt) error { return nil }

func (inter *Interpreter) VisitIfStmt(stmt *object.IfStmt) error {
	cond, err := inter.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	if isTruthy(cond) {
		return inter.execute(stmt.Then)
	}

	if stmt.Else != nil {
		return inter.execute(stmt.Else)
	}

	return nil
}

func (inter *Interpreter) VisitForStmt(stmt *object.ForStmt) error {
	prev := inter.env
	inter.env = object.NewEnclosedEnvironment(prev)

	if stmt.Initializer != nil {
		err := inter.execute(stmt.Initializer)
		if err != nil {
			inter.env = prev
			return err
		}
	}

	if stmt.Keyword.Type == lexer.Do {
		err := inter.execute(stmt.Body)
		if err != nil {
			if err == BreakError {
				inter.env = prev
				return nil
			}
			// If Continue error then don't return
			if err != ContinueError {
				inter.env = prev
				return err
			}
		}
	}

	cond, err := inter.evaluate(stmt.Condition)
	for err == nil && isTruthy(cond) {
		err = inter.execute(stmt.Body)
		if err == BreakError {
			err = nil
			break
		} else if err == ContinueError {
			err = nil
		} else if err != nil {
			break
		}

		if stmt.Increment != nil {
			_, err := inter.evaluate(stmt.Increment)
			if err != nil {
				break
			}
		}

		cond, err = inter.evaluate(stmt.Condition)
	}

	inter.env = prev
	return err
}

func (inter *Interpreter) VisitReturnStmt(stmt *object.ReturnStmt) error {
	var value object.Object
	var err error

	if stmt.Value != nil {
		value, err = inter.evaluate(stmt.Value)
		if err != nil {
			return err
		}
	}

	return NewReturnValue(stmt.Keyword, value)
}

func (inter *Interpreter) VisitVarStmt(stmt *object.VarStmt) error {
	var value object.Object
	var err error

	if stmt.Value != nil {
		value, err = inter.evaluate(stmt.Value)
		if err != nil {
			return err
		}
	}

	inter.env.Define(stmt.Name, value)
	return nil
}

func (inter *Interpreter) VisitAssignExpr(expr *object.AssignExpr) (object.Object, error) {
	value, err := inter.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	if dist, ok := inter.local[expr]; ok {
		err = inter.env.AssignAt(dist, expr.Name, value)
	} else {
		err = inter.env.Assign(expr.Name, value)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (inter *Interpreter) VisitBinaryExpr(expr *object.BinaryExpr) (object.Object, error) {
	left, err := inter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := inter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case lexer.Greater, lexer.GreaterEq, lexer.Less, lexer.LessEq, lexer.EqualEq, lexer.BangEq:
		return numberComparisonOperation(expr.Operator, left, right)
	case lexer.Minus, lexer.Star, lexer.Slash, lexer.TildSlash, lexer.Percent:
		return numberMathOperation(expr.Operator, left, right)
	case lexer.Plus:
		if left.Type() == object.Number && right.Type() == object.Number {
			return numberMathOperation(expr.Operator, left, right)
		}
		if left.Type() == object.String && right.Type() == object.String {
			l := left.(*String)
			r := right.(*String)
			return &String{Value: l.Value + r.Value}, nil
		}
	}

	return nil, object.NewRuntimeError(expr.Operator, fmt.Sprintf("No known operations for %s %s %s", left.Type(), expr.Operator.Lexeme, right.Type()))
}

func numberComparisonOperation(oper *lexer.Token, left, right object.Object) (*Boolean, error) {
	if left.Type() != object.Number || right.Type() != object.Number {
		return nil, object.NewRuntimeError(oper, "Operands must be numbers.")
	}

	l := left.(*Number)
	r := right.(*Number)

	var value bool
	if l.IsInt && r.IsInt {
		value = intComparison(oper.Lexeme, l.Int, r.Int)
	} else if l.IsInt {
		value = floatComparison(oper.Lexeme, float64(l.Int), r.Float)
	} else if r.IsInt {
		value = floatComparison(oper.Lexeme, l.Float, float64(r.Int))
	} else {
		value = floatComparison(oper.Lexeme, l.Float, r.Float)
	}

	return &Boolean{Value: value}, nil
}

func intComparison(oper string, left, right int) bool {
	switch oper {
	case ">":
		return left > right
	case ">=":
		return left >= right
	case "<":
		return left < right
	case "<=":
		return left <= right
	case "==":
		return left == right
	case "!=":
		return left != right
	}

	return false
}

func floatComparison(oper string, left, right float64) bool {
	switch oper {
	case ">":
		return left > right
	case ">=":
		return left >= right
	case "<":
		return left < right
	case "<=":
		return left <= right
	case "==":
		return left == right
	case "!=":
		return left != right
	}

	return false
}

func numberMathOperation(oper *lexer.Token, left, right object.Object) (*Number, error) {
	if left.Type() != object.Number || right.Type() != object.Number {
		return nil, object.NewRuntimeError(oper, "Operands must be numbers.")
	}

	l := left.(*Number)
	r := right.(*Number)

	number := &Number{}
	if l.IsInt && r.IsInt {
		if oper.Lexeme == "/" {
			if l.Int%r.Int == 0 {
				number.IsInt = true
				number.Int = l.Int / r.Int
			} else {
				number.Float = float64(l.Int) / float64(r.Int)
			}
			return number, nil
		}
		number.IsInt = true
		number.Int = intOperation(oper.Lexeme, l.Int, r.Int)
		return number, nil
	}

	if oper.Lexeme == "%" {
		return nil, object.NewRuntimeError(oper, "Operands must both be integer values.")
	}

	if l.IsInt {
		number.Float = floatOperation(oper.Lexeme, float64(l.Int), r.Float)
	} else if r.IsInt {
		number.Float = floatOperation(oper.Lexeme, l.Float, float64(r.Int))
	} else {
		number.Float = floatOperation(oper.Lexeme, l.Float, r.Float)
	}

	if oper.Lexeme == "~/" {
		number.IsInt = true
		number.Int = int(number.Float)
	}

	return number, nil
}

func intOperation(oper string, left, right int) int {
	switch oper {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "~/":
		return left / right
	case "%":
		return left % right
	}
	return 0
}

func floatOperation(oper string, left, right float64) float64 {
	switch oper {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right
	}
	return 0
}

func (inter *Interpreter) VisitBooleanExpr(expr *object.BooleanExpr) (object.Object, error) {
	return nil, nil
}
func (inter *Interpreter) VisitCallExpr(expr *object.CallExpr) (object.Object, error) { return nil, nil }
func (inter *Interpreter) VisitGetExpr(expr *object.GetExpr) (object.Object, error)   { return nil, nil }
func (inter *Interpreter) VisitGroupingExpr(expr *object.GroupingExpr) (object.Object, error) {
	return nil, nil
}
func (inter *Interpreter) VisitIndexExpr(expr *object.IndexExpr) (object.Object, error) {
	return nil, nil
}
func (inter *Interpreter) VisitListExpr(expr *object.ListExpr) (object.Object, error) { return nil, nil }
func (inter *Interpreter) VisitLogicalExpr(expr *object.LogicalExpr) (object.Object, error) {
	return nil, nil
}
func (inter *Interpreter) VisitNumberExpr(expr *object.NumberExpr) (object.Object, error) {
	return nil, nil
}
func (inter *Interpreter) VisitNullExpr(expr *object.NullExpr) (object.Object, error) { return nil, nil }
func (inter *Interpreter) VisitSetExpr(expr *object.SetExpr) (object.Object, error)   { return nil, nil }
func (inter *Interpreter) VisitStringExpr(expr *object.StringExpr) (object.Object, error) {
	return nil, nil
}
func (inter *Interpreter) VisitSuperExpr(expr *object.SuperExpr) (object.Object, error) {
	return nil, nil
}
func (inter *Interpreter) VisitThisExpr(expr *object.ThisExpr) (object.Object, error) { return nil, nil }
func (inter *Interpreter) VisitUnaryExpr(expr *object.UnaryExpr) (object.Object, error) {
	return nil, nil
}
func (inter *Interpreter) VisitVariableExpr(expr *object.VariableExpr) (object.Object, error) {
	return nil, nil
}
