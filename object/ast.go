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

// BooleanExpr is a Expr of a Boolean
type BooleanExpr struct {
	Token *lexer.Token
	Value bool
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (b *BooleanExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitBooleanExpr(b) }

// CallExpr is a Expr of a Call
type CallExpr struct {
	Callee Expr
	Paren  *lexer.Token
	Args   []Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (c *CallExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitCallExpr(c) }

// GetExpr is a Expr of a Get
type GetExpr struct {
	Object Expr
	Name   *lexer.Token
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (g *GetExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitGetExpr(g) }

// GroupingExpr is a Expr of a Grouping
type GroupingExpr struct {
	Expression Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (g *GroupingExpr) Accept(visitor ExprVisitor) (Object, error) {
	return visitor.VisitGroupingExpr(g)
}

// IndexExpr is a Expr of a Index
type IndexExpr struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (i *IndexExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitIndexExpr(i) }

// ListExpr is a Expr of a List
type ListExpr struct {
	Values []Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (l *ListExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitListExpr(l) }

// LogicalExpr is a Expr of a Logical
type LogicalExpr struct {
	Left     Expr
	Operator *lexer.Token
	Right    Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (l *LogicalExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitLogicalExpr(l) }

// NumberExpr is a Expr of a Number
type NumberExpr struct {
	Token *lexer.Token
	Float float64
	Int   int
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (n *NumberExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitNumberExpr(n) }

// NullExpr is a Expr of a Null
type NullExpr struct {
	Token *lexer.Token
	Value interface{}
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (n *NullExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitNullExpr(n) }

// SetExpr is a Expr of a Set
type SetExpr struct {
	Object  Expr
	Name    *lexer.Token
	Value   Expr
	IsIndex bool
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (s *SetExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitSetExpr(s) }

// StringExpr is a Expr of a String
type StringExpr struct {
	Token *lexer.Token
	Value string
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (s *StringExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitStringExpr(s) }

// SuperExpr is a Expr of a Super
type SuperExpr struct {
	Keyword *lexer.Token
	Method  *lexer.Token
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (s *SuperExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitSuperExpr(s) }

// ThisExpr is a Expr of a This
type ThisExpr struct {
	Keyword *lexer.Token
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (t *ThisExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitThisExpr(t) }

// UnaryExpr is a Expr of a Unary
type UnaryExpr struct {
	Operator *lexer.Token
	Right    Expr
}

// Accept calls the correct visit method on ExprVisitor, passing a reference to itself as a value
func (u *UnaryExpr) Accept(visitor ExprVisitor) (Object, error) { return visitor.VisitUnaryExpr(u) }

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
	VisitBooleanExpr(expr *BooleanExpr) (Object, error)
	VisitCallExpr(expr *CallExpr) (Object, error)
	VisitGetExpr(expr *GetExpr) (Object, error)
	VisitGroupingExpr(expr *GroupingExpr) (Object, error)
	VisitIndexExpr(expr *IndexExpr) (Object, error)
	VisitListExpr(expr *ListExpr) (Object, error)
	VisitLogicalExpr(expr *LogicalExpr) (Object, error)
	VisitNumberExpr(expr *NumberExpr) (Object, error)
	VisitNullExpr(expr *NullExpr) (Object, error)
	VisitSetExpr(expr *SetExpr) (Object, error)
	VisitStringExpr(expr *StringExpr) (Object, error)
	VisitSuperExpr(expr *SuperExpr) (Object, error)
	VisitThisExpr(expr *ThisExpr) (Object, error)
	VisitUnaryExpr(expr *UnaryExpr) (Object, error)
	VisitVariableExpr(expr *VariableExpr) (Object, error)
}

// BlockStmt is a Stmt of a Block
type BlockStmt struct {
	Statements []Stmt
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (b *BlockStmt) Accept(visitor StmtVisitor) error { return visitor.VisitBlockStmt(b) }

// BreakStmt is a Stmt of a Break
type BreakStmt struct {
	Keyword *lexer.Token
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (b *BreakStmt) Accept(visitor StmtVisitor) error { return visitor.VisitBreakStmt(b) }

// ClassStmt is a Stmt of a Class
type ClassStmt struct {
	Name    *lexer.Token
	Super   *VariableExpr
	Methods []*FunctionStmt
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (c *ClassStmt) Accept(visitor StmtVisitor) error { return visitor.VisitClassStmt(c) }

// ContinueStmt is a Stmt of a Continue
type ContinueStmt struct {
	Keyword *lexer.Token
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (c *ContinueStmt) Accept(visitor StmtVisitor) error { return visitor.VisitContinueStmt(c) }

// ExpressionStmt is a Stmt of a Expression
type ExpressionStmt struct {
	Expression Expr
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (e *ExpressionStmt) Accept(visitor StmtVisitor) error { return visitor.VisitExpressionStmt(e) }

// FunctionStmt is a Stmt of a Function
type FunctionStmt struct {
	Name       *lexer.Token
	Parameters []*lexer.Token
	Body       []Stmt
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (f *FunctionStmt) Accept(visitor StmtVisitor) error { return visitor.VisitFunctionStmt(f) }

// IfStmt is a Stmt of a If
type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (i *IfStmt) Accept(visitor StmtVisitor) error { return visitor.VisitIfStmt(i) }

// ImportStmt is a Stmt of a Import
type ImportStmt struct {
	Keyword *lexer.Token
	Other   Expr
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (i *ImportStmt) Accept(visitor StmtVisitor) error { return visitor.VisitImportStmt(i) }

// ForStmt is a Stmt of a For
type ForStmt struct {
	Keyword     *lexer.Token
	Initializer Stmt
	Condition   Expr
	Body        Stmt
	Increment   Expr
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (f *ForStmt) Accept(visitor StmtVisitor) error { return visitor.VisitForStmt(f) }

// ReturnStmt is a Stmt of a Return
type ReturnStmt struct {
	Keyword *lexer.Token
	Value   Expr
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (r *ReturnStmt) Accept(visitor StmtVisitor) error { return visitor.VisitReturnStmt(r) }

// VarStmt is a Stmt of a Var
type VarStmt struct {
	Name  *lexer.Token
	Value Expr
}

// Accept calls the correct visit method on StmtVisitor, passing a reference to itself as a value
func (v *VarStmt) Accept(visitor StmtVisitor) error { return visitor.VisitVarStmt(v) }

// StmtVisitor will visit Stmt objects and must receive calls to their applicable methods.
type StmtVisitor interface {
	VisitBlockStmt(stmt *BlockStmt) error
	VisitBreakStmt(stmt *BreakStmt) error
	VisitClassStmt(stmt *ClassStmt) error
	VisitContinueStmt(stmt *ContinueStmt) error
	VisitExpressionStmt(stmt *ExpressionStmt) error
	VisitFunctionStmt(stmt *FunctionStmt) error
	VisitIfStmt(stmt *IfStmt) error
	VisitImportStmt(stmt *ImportStmt) error
	VisitForStmt(stmt *ForStmt) error
	VisitReturnStmt(stmt *ReturnStmt) error
	VisitVarStmt(stmt *VarStmt) error
}
