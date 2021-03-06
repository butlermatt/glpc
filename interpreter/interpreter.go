package interpreter

import (
	"errors"
	"fmt"
	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
	"github.com/butlermatt/glpc/parser"
	"io/ioutil"
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
	local   map[object.Expr]int
	env     *object.Environment
	globals *object.Environment
}

func New() *Interpreter {
	glob := object.GetGlobal()
	glob = SetupGlobal(glob)
	return &Interpreter{globals: glob}
}

func (inter *Interpreter) RunMain(env *object.Environment) error {
	inter.env = env

	mnFn := inter.env.GetString("main")
	if mnFn == nil {
		return fmt.Errorf("Unable to locate main function.")
	}

	if mnFn.Type() != object.Function {
		return fmt.Errorf("Found main, but it was not a function.")
	}

	mainFunc := mnFn.(*Function)
	_, err := mainFunc.Call(inter, nil)

	return err
}

func (inter *Interpreter) Interpret(parser *parser.Parser, filename string) (*object.Environment, error) {
	stmts, depth := parser.Parse()
	errs := parser.Errors()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Printf("[Syntax error] %+v\n", err)
		}
		return nil, fmt.Errorf("%d syntax errors found.", len(errs))
	}

	inter.local = depth
	inter.env = object.NewEnvironment(filename)
	for _, stmt := range stmts {
		err := inter.execute(stmt)
		if err != nil {
			return nil, err
		}
	}

	return inter.env, nil
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
		err = inter.execute(stmt)
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
		return obj == True
	}

	return true
}

func (inter *Interpreter) VisitBlockStmt(stmt *object.BlockStmt) error {
	return inter.executeBlock(stmt.Statements, object.NewEnclosedEnvironment(inter.env))
}

func (inter *Interpreter) VisitBreakStmt(stmt *object.BreakStmt) error { return BreakError }

func (inter *Interpreter) VisitClassStmt(stmt *object.ClassStmt) error {
	inter.env.Define(stmt.Name, nil)

	var superClass *Class
	var prevEnv *object.Environment
	if stmt.Super != nil {
		sc, err := inter.evaluate(stmt.Super)
		if err != nil {
			return err
		}
		if sc.Type() != object.Class {
			return object.NewRuntimeError(stmt.Super.Name, "Superclass must be a class.")
		}
		superClass = sc.(*Class)
		prevEnv = inter.env
		inter.env = object.NewEnclosedEnvironment(inter.env)
		inter.env.DefineString("super", superClass)
	}

	var methods = make(map[string]*Function)
	for _, meth := range stmt.Methods {
		methods[meth.Name.Lexeme] = NewFunction(meth, inter.env, meth.Name.Lexeme == "init")
	}

	klass := &Class{Name: stmt.Name.Lexeme, superclass: superClass, methods: methods}
	if prevEnv != nil {
		inter.env = prevEnv
	}

	inter.env.Assign(stmt.Name, klass)
	return nil
}

func (inter *Interpreter) VisitContinueStmt(stmt *object.ContinueStmt) error { return ContinueError }

func (inter *Interpreter) VisitExpressionStmt(stmt *object.ExpressionStmt) error {
	_, err := inter.evaluate(stmt.Expression)
	return err
}

func (inter *Interpreter) VisitFunctionStmt(stmt *object.FunctionStmt) error {
	fun := NewFunction(stmt, inter.env, false)
	inter.env.Define(stmt.Name, fun)
	return nil
}

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

func (inter *Interpreter) VisitImportStmt(stmt *object.ImportStmt) error {
	str, ok := stmt.Other.(*object.StringExpr)
	if !ok {
		return object.NewRuntimeError(stmt.Keyword, "Other filename unexpected type.")
	}

	if str.Value == "" {
		return object.NewRuntimeError(stmt.Keyword, "imported filename cannot be empty.")
	}

	oEnv := object.GetFileEnvironment(str.Value)
	if oEnv != nil {
		inter.env.Copy(oEnv)
		return nil
	}

	file, err := ioutil.ReadFile(str.Value)
	if err != nil {
		return object.NewRuntimeError(stmt.Keyword, "error reading file "+str.Value+". Error: "+err.Error())
	}

	l := lexer.New(file, str.Value)
	p := parser.New(l)
	interpreter := New()
	oEnv, err = interpreter.Interpret(p, str.Value)

	inter.env.Copy(oEnv)

	return err
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
	case lexer.Greater, lexer.GreaterEq, lexer.Less, lexer.LessEq:
		return numberComparisonOperation(expr.Operator, left, right)
	case lexer.Minus, lexer.Star, lexer.Slash, lexer.TildSlash, lexer.Percent:
		return numberMathOperation(expr.Operator, left, right)
	case lexer.EqualEq:
		return isEqual(expr.Operator, left, right)
	case lexer.BangEq:
		b, err := isEqual(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		b.Value = !b.Value
		return b, nil
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

	if value {
		return True, nil
	}

	return False, nil
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

func isEqual(oper *lexer.Token, left, right object.Object) (*Boolean, error) {
	if left.Type() == object.Null && right.Type() == object.Null {
		return True, nil
	}

	if left.Type() != right.Type() {
		return False, nil
	}

	switch left.Type() {
	case object.Number:
		return numberComparisonOperation(oper, left, right)
	case object.Boolean:
		if left == right {
			return True, nil
		}
		return False, nil
	case object.String:
		l := left.(*String)
		r := right.(*String)
		if l.Value == r.Value {
			return True, nil
		}
		return False, nil
	}

	// Shouldn't ever reach here
	return nil, object.NewRuntimeError(oper, fmt.Sprintf("No known operations for %s %s %s", left.Type(), oper.Lexeme, right.Type()))
}

func (inter *Interpreter) VisitBooleanExpr(expr *object.BooleanExpr) (object.Object, error) {
	if expr.Value {
		return True, nil
	}
	return False, nil
}

func (inter *Interpreter) VisitCallExpr(expr *object.CallExpr) (object.Object, error) {
	callee, err := inter.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	if callee.Type() != object.Function && callee.Type() != object.Class && callee.Type() != object.BuiltIn {
		return nil, object.NewRuntimeError(expr.Paren, "Can only call functions and classes.")
	}

	var args []object.Object
	for _, arg := range expr.Args {
		a, err := inter.evaluate(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}

	function := callee.(Callable)

	// Allow any number of args.
	if function.Arity() == -1 {
		return function.Call(inter, args)
	}

	if len(args) != function.Arity() {
		return nil, object.NewRuntimeError(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(args)))
	}

	return function.Call(inter, args)
}

func (inter *Interpreter) VisitGetExpr(expr *object.GetExpr) (object.Object, error) {
	obj, err := inter.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if obj.Type() != object.Instance {
		return nil, object.NewRuntimeError(expr.Name, "Only instances have properties.")
	}

	inst := obj.(*Instance)
	return inst.Get(expr.Name)
}

func (inter *Interpreter) VisitGroupingExpr(expr *object.GroupingExpr) (object.Object, error) {
	return inter.evaluate(expr.Expression)
}

func (inter *Interpreter) VisitIndexExpr(expr *object.IndexExpr) (object.Object, error) {
	left, err := inter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := inter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	if left.Type() != object.List {
		return nil, object.NewRuntimeError(expr.Operator, "Cannot perform index lookup on anything except a list.")
	}

	if right.Type() != object.Number {
		return nil, object.NewRuntimeError(expr.Operator, "Index operand must be a number.")
	}

	l := left.(*List)
	r := right.(*Number)
	var ind int
	if r.IsInt {
		ind = r.Int
	} else {
		ind = int(r.Float)
	}

	if ind >= len(l.Elements) {
		return nil, object.NewRuntimeError(expr.Operator, "Index out of range.")
	}
	return l.Elements[ind], nil

}

func (inter *Interpreter) VisitListExpr(expr *object.ListExpr) (object.Object, error) {
	listLen := len(expr.Values)
	list := &List{Elements: make([]object.Object, listLen)}

	for i := 0; i < listLen; i++ {
		value, err := inter.evaluate(expr.Values[i])
		if err != nil {
			return nil, err
		}
		list.Elements[i] = value
	}

	return list, nil
}

func (inter *Interpreter) VisitLogicalExpr(expr *object.LogicalExpr) (object.Object, error) {
	left, err := inter.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.Type == lexer.Or {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return inter.evaluate(expr.Right)
}

func (inter *Interpreter) VisitNumberExpr(expr *object.NumberExpr) (object.Object, error) {
	n := &Number{}
	if expr.Token.Type == lexer.NumberI {
		n.IsInt = true
		n.Int = expr.Int
	} else {
		n.Float = expr.Float
	}

	return n, nil
}

func (inter *Interpreter) VisitNullExpr(expr *object.NullExpr) (object.Object, error) {
	return NullOb, nil
}

func (inter *Interpreter) VisitSetExpr(expr *object.SetExpr) (object.Object, error) {
	if expr.IsIndex {
		return inter.setIndexValue(expr)
	}

	obj, err := inter.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if obj.Type() != object.Instance {
		return nil, object.NewRuntimeError(expr.Name, "Only instances have fields.")
	}

	inst := obj.(*Instance)
	value, err := inter.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	inst.Set(expr.Name, value)
	return value, nil
}

func (inter *Interpreter) setIndexValue(expr *object.SetExpr) (object.Object, error) {
	ie := expr.Object.(*object.IndexExpr)
	li, err := inter.evaluate(ie.Left)
	if err != nil {
		return nil, err
	}
	if li.Type() != object.List {
		return nil, object.NewRuntimeError(ie.Operator, "Cannot perform index lookup on anything except a list.")
	}
	list := li.(*List)
	ind, err := inter.evaluate(ie.Right)
	if err != nil {
		return nil, err
	}
	if ind.Type() != object.Number {
		return nil, object.NewRuntimeError(ie.Operator, "Operand must be a number.")
	}
	in := ind.(*Number)
	var index int

	if in.IsInt {
		index = in.Int
	} else {
		index = int(in.Float)
	}

	if index >= len(list.Elements) {
		return nil, object.NewRuntimeError(ie.Operator, "Index out of range.")
	}
	value, err := inter.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	list.Elements[index] = value
	return value, nil
}

func (inter *Interpreter) VisitStringExpr(expr *object.StringExpr) (object.Object, error) {
	return &String{Value: expr.Value}, nil
}

func (inter *Interpreter) VisitSuperExpr(expr *object.SuperExpr) (object.Object, error) {
	var superClass *Class
	var obj *Instance

	dist := inter.local[expr]
	sc, err := inter.env.GetAt(dist, expr.Keyword)
	if err != nil {
		return nil, err
	}

	if sc.Type() != object.Class {
		return nil, object.NewRuntimeError(expr.Keyword, "Superclass was not a Class.")
	}

	superClass = sc.(*Class)
	thisTok := lexer.NewToken(lexer.This, "this", expr.Keyword.Filename, expr.Keyword.Line)
	o, err := inter.env.GetAt(dist-1, thisTok)
	if err != nil {
		return nil, err
	}

	if o.Type() != object.Instance {
		return nil, object.NewRuntimeError(expr.Keyword, "this was not an Instance of a class.")
	}
	obj = o.(*Instance)

	method := superClass.findMethod(obj, expr.Method.Lexeme)
	if method == nil {
		return nil, object.NewRuntimeError(expr.Method, "Undefined property on "+superClass.Name+".")
	}
	return method, nil
}

func (inter *Interpreter) VisitThisExpr(expr *object.ThisExpr) (object.Object, error) {
	return inter.lookupVariable(expr.Keyword, expr)
}

func (inter *Interpreter) VisitUnaryExpr(expr *object.UnaryExpr) (object.Object, error) {
	right, err := inter.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case lexer.Minus:
		if right.Type() != object.Number {
			return nil, object.NewRuntimeError(expr.Operator, "Operand must be a number.")
		}
		r := right.(*Number)
		if r.IsInt {
			r.Int = -r.Int
		} else {
			r.Float = -r.Float
		}
		return r, nil
	case lexer.Bang:
		b := !isTruthy(right)
		if b {
			return True, nil
		}
		return False, nil
	}

	// should never reach here.
	return nil, nil
}

func (inter *Interpreter) VisitVariableExpr(expr *object.VariableExpr) (object.Object, error) {
	return inter.lookupVariable(expr.Name, expr)
}

func (inter *Interpreter) lookupVariable(name *lexer.Token, expr object.Expr) (object.Object, error) {
	dist := inter.local[expr]
	obj, _ := inter.env.GetAt(dist, name)

	if obj != nil {
		return obj, nil
	}

	return inter.globals.Get(name)
}
