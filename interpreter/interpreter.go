package interpreter

import (
	"github.com/butlermatt/glpc/object"
	"errors"
	"fmt"
	"github.com/butlermatt/glpc/lexer"
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

func (inter *Interpreter) VisitBreakStmt(stmt *object.BreakStmt) error {return BreakError }
func (inter *Interpreter) VisitClassStmt(stmt *object.ClassStmt) error {return nil}
func (inter *Interpreter) VisitContinueStmt(stmt *object.ContinueStmt) error {return ContinueError }

func (inter *Interpreter) VisitExpressionStmt(stmt *object.ExpressionStmt) error {
	_, err := inter.evaluate(stmt.Expression)
	return err
}

func (inter *Interpreter) VisitFunctionStmt(stmt *object.FunctionStmt) error {return nil}

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
		}else if err == ContinueError {
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

func (inter *Interpreter) VisitAssignExpr(expr *object.AssignExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitBinaryExpr(expr *object.BinaryExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitBooleanExpr(expr *object.BooleanExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitCallExpr(expr *object.CallExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitGetExpr(expr *object.GetExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitGroupingExpr(expr *object.GroupingExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitIndexExpr(expr *object.IndexExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitListExpr(expr *object.ListExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitLogicalExpr(expr *object.LogicalExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitNumberExpr(expr *object.NumberExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitNullExpr(expr *object.NullExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitSetExpr(expr *object.SetExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitStringExpr(expr *object.StringExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitSuperExpr(expr *object.SuperExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitThisExpr(expr *object.ThisExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitUnaryExpr(expr *object.UnaryExpr) (object.Object, error) {return nil, nil}
func (inter *Interpreter) VisitVariableExpr(expr *object.VariableExpr) (object.Object, error) {return nil, nil}
