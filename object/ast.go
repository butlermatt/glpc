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

// ListExpr is a Expr of a List
type ListExpr struct {
	Values []Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (l *ListExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitListExpr(l) }

// NumberExpr is a Expr of a Number
type NumberExpr struct {
	Token *lexer.Token
	Float float64
	Int   int
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (n *NumberExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitNumberExpr(n) }

// ExprVisitor will visit Expr objects and must receive calls to their applicable methods.
type ExprVisitor interface {
	VisitBooleanExpr(expr *BooleanExpr) (Object, error)
	VisitListExpr(expr *ListExpr) (Object, error)
	VisitNumberExpr(expr *NumberExpr) (Object, error)
	VisitNullExpr(expr *NullExpr) (Object, error)
	VisitStringExpr(expr *StringExpr) (Object, error)
	VisitUnaryExpr(expr *UnaryExpr) (Object, error)
}

// ExpressionStmt is a Stmt of a Expression
type ExpressionStmt struct {
	Expression Expr
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (e *ExpressionStmt) Accept(visitor StmtVisitor) error { return visitor.VisitExpressionStmt(e) }

// StmtVisitor will visit Stmt objects and must receive calls to their applicable methods.
type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) error
}
