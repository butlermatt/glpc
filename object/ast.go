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

// NumberExpr is a Expr of a Number
type NumberExpr struct {
	Token *lexer.Token
	Float float64
	Int   int
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (n *NumberExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitNumberExpr(n)
}

// ExprVisitor will visit Expr objects and must receive calls to their applicable methods.
type ExprVisitor interface {
	VisitNumberExpr(expr *NumberExpr) (Object, error)
}

// StmtVisitor will visit Stmt objects and must receive calls to their applicable methods.
type StmtVisitor interface {
}