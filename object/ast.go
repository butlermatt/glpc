package object

import "github.com/butlermatt/glpc/lexer"

// Expr is an AST expression which returns a value of type Object or an error.
type Expr interface {
	Accept(ExprVisitor) (Object, error)
}

// Stmt is an AST statement which returns no value but may produce an error.
type Stmt interface {
	Accept(StmtVisitor) error
}

// AssignExpr is a Expr of a Assign
type AssignExpr struct {
	Name  *lexer.Token
	Value Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (a *AssignExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitAssignExpr(a) }

// BinaryExpr is a Expr of a Binary
type BinaryExpr struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (b *BinaryExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitBinaryExpr(b) }

// NumberExpr is a Expr of a Number
type NumberExpr struct {
	Token *lexer.Token
	Float float64
	Int   int
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (n *NumberExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitNumberExpr(n) }

// UnaryExpr is a Expr of a Unary
type UnaryExpr struct {
	Operator *lexer.Token
	Right    Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (u *UnaryExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitUnaryExpr(u) }

// BooleanExpr is a Expr of a Boolean
type BooleanExpr struct {
	Token *lexer.Token
	Value bool
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (b *BooleanExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitBooleanExpr(b) }

// GroupingExpr is a Expr of a Grouping
type GroupingExpr struct {
	Expression Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (g *GroupingExpr) Accept(visitor ExprVisitor) (Object, error) {
	return visitor.VisitGroupingExpr(g)
}

// ListExpr is a Expr of a List
type ListExpr struct {
	Values []Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (l *ListExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitListExpr(l) }

// NullExpr is a Expr of a Null
type NullExpr struct {
	Token *lexer.Token
	Value interface{}
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (n *NullExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitNullExpr(n) }

// StringExpr is a Expr of a String
type StringExpr struct {
	Token *lexer.Token
	Value string
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (s *StringExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitStringExpr(s) }

// VariableExpr is a Expr of a Variable
type VariableExpr struct {
	Name *lexer.Token
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (v *VariableExpr) Accept(visitor ExprVisitor) (Object, error) {
	return visitor.VisitVariableExpr(v)
}

// ExprVisitor will visit Expr objects and must receive calls to their applicable methods.
type ExprVisitor interface {
	VisitAssignExpr(expr *AssignExpr) (Object, error)
	VisitBinaryExpr(expr *BinaryExpr) (Object, error)
	VisitNumberExpr(expr *NumberExpr) (Object, error)
	VisitUnaryExpr(expr *UnaryExpr) (Object, error)
	VisitBooleanExpr(expr *BooleanExpr) (Object, error)
	VisitGroupingExpr(expr *GroupingExpr) (Object, error)
	VisitListExpr(expr *ListExpr) (Object, error)
	VisitNullExpr(expr *NullExpr) (Object, error)
	VisitStringExpr(expr *StringExpr) (Object, error)
	VisitVariableExpr(expr *VariableExpr) (Object, error)
}

// ExpressionStmt is a Stmt of a Expression
type ExpressionStmt struct {
	Expression Expr
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (e *ExpressionStmt) Accept(visitor StmtVisitor) error { return visitor.VisitExpressionStmt(e) }

// VarStmt is a Stmt of a Var
type VarStmt struct {
	Name  *lexer.Token
	Value Expr
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (v *VarStmt) Accept(visitor StmtVisitor) error { return visitor.VisitVarStmt(v) }

// StmtVisitor will visit Stmt objects and must receive calls to their applicable methods.
type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) error
	VisitVarStmt(stmt *VarStmt) error
}
