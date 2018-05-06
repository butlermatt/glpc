package parser

import (
	"fmt"
	"strconv"

	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
)

// ParseError represents compile time syntax errors that the Parser discovers. It will try to recover from them
// when possible.
type ParseError struct {
	// Line is the line number the syntax error occurred on.
	Line int
	// Where will try to provide a hint to which token generated the syntax error.
	Where string
	// Msg will attempt to describe the error encountered.
	Msg string
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("On line %d: %s - %s", pe.Line, pe.Where, pe.Msg)
}

// Parser iterates through the tokens scanned by the lexer and generates the correct AST.
type Parser struct {
	l       *lexer.Lexer
	curTok  *lexer.Token
	prevTok *lexer.Token
	errors  []ParseError
	errLen  int
}

// New will return a new Parser initialized with the tokens from lexer. This will call ScanTokens on the lexer. Do not
// scan tokens prior to passing to the Parser.
func New(lexer *lexer.Lexer) *Parser {
	lexer.ScanTokens()
	p := &Parser{l: lexer}
	p.nextToken()

	return p
}

// Errors returns a slice of ParseErrors that were encountered during parsing.
func (p *Parser) Errors() []ParseError {
	return p.errors
}

// Parse will parse the tokens provided by the lexer and return a slice of statements that comprise the program.
func (p *Parser) Parse() []object.Stmt {
	var stmts []object.Stmt
	for p.curTok.Type != lexer.EOF {
		s := p.declaration()
		if s != nil {
			stmts = append(stmts, s)
		}
	}

	return stmts
}

func (p *Parser) addError(token *lexer.Token, msg string) {
	if token.Type == lexer.EOF {
		p.errors = append(p.errors, ParseError{Line: token.Line, Where: "at end", Msg: msg})
	} else {
		p.errors = append(p.errors, ParseError{Line: token.Line, Where: token.Lexeme, Msg: msg})
	}
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.curTok.Type == lexer.EOF {
		return false
	}
	return p.curTok.Type == tokenType
}

func (p *Parser) consume(tokenType lexer.TokenType, errMsg string) bool {
	if p.check(tokenType) {
		p.nextToken()
		return true
	}

	p.addError(p.curTok, errMsg)
	return false
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, tt := range types {
		if p.check(tt) {
			p.nextToken()
			return true
		}
	}

	return false
}

func (p *Parser) nextToken() {
	if p.curTok == nil || p.curTok.Type != lexer.EOF {
		p.prevTok = p.curTok
		p.curTok = p.l.NextToken()
	}
}

func (p *Parser) declaration() object.Stmt {
	var stmt object.Stmt

	switch {
	case p.match(lexer.Var):
		stmt = p.varDeclaration()
	default:
		stmt = p.statement()
	}

	if len(p.errors) > p.errLen {
		p.errLen = len(p.errors)
		p.synchronize()
		return nil
	}

	return stmt
}

func (p *Parser) varDeclaration() object.Stmt {
	if !p.consume(lexer.Ident, "Expect variable name.") {
		return nil
	}

	name := p.prevTok

	var init object.Expr
	if p.match(lexer.Equal) {
		init = p.expression()
	}

	p.consume(lexer.Semicolon, "Expect ';' after variable declaration.")
	return &object.VarStmt{Name: name, Value: init}
}

func (p *Parser) statement() object.Stmt {
	switch {
	case p.match(lexer.LBrace):
		return &object.BlockStmt{Statements: p.block()}
	case p.match(lexer.Break):
		return p.breakStatement()
	case p.match(lexer.Continue):
		return p.continueStatement()
	case p.match(lexer.Do):
		return p.doWhileStatement()
	case p.match(lexer.For):
		return p.forStatement()
	case p.match(lexer.If):
		return p.ifStatement()
	case p.match(lexer.While):
		return p.whileStatement()
		// TODO: Cases for break, continue, If, Print, Return, While, For, and LBrace
	}

	return p.expressionStatement()
}

func (p *Parser) block() []object.Stmt {
	var stmts []object.Stmt

	for !p.check(lexer.RBrace) && p.curTok.Type != lexer.EOF {
		s := p.declaration()
		if s != nil {
			stmts = append(stmts, s)
		}
	}

	p.consume(lexer.RBrace, "Expect '}' after block.")
	return stmts
}

func (p *Parser) breakStatement() object.Stmt {
	// TODO: add check for in loop and return error
	keyword := p.prevTok
	if !p.consume(lexer.Semicolon, "Expect ';' after 'break'.") {
		return nil
	}

	return &object.BreakStmt{Keyword: keyword}
}

func (p *Parser) continueStatement() object.Stmt {
	// TODO: add check for in loop and return error
	keyword := p.prevTok
	if !p.consume(lexer.Semicolon, "Expect ';' after 'continue'.") {
		return nil
	}

	return &object.ContinueStmt{Keyword: keyword}
}

func (p *Parser) doWhileStatement() object.Stmt {
	body := p.statement()

	if !p.consume(lexer.While, "Expect 'while' after do-while body.") {
		return nil
	}

	if !p.consume(lexer.LParen, "Expect '(' after 'while'.") {
		return nil
	}

	cond := p.expression()
	if !p.consume(lexer.RParen, "Expect ')' after while condition.") {
		return nil
	}
	if !p.consume(lexer.Semicolon, "Expect ';' after ')'.") {
		return nil
	}

	stmts := []object.Stmt{body, &object.ForStmt{Condition: cond, Body: body}}
	return &object.BlockStmt{Statements: stmts}
}

func (p *Parser) forStatement() object.Stmt {
	if !p.consume(lexer.LParen, "Expect '(' after 'for'.") {
		return nil
	}

	var init object.Stmt
	if p.match(lexer.Semicolon) {
		init = nil
	} else if p.match(lexer.Var) {
		init = p.varDeclaration()
	} else {
		init = p.expressionStatement()
	}

	var cond object.Expr
	if !p.check(lexer.Semicolon) {
		cond = p.expression()
	}
	if !p.consume(lexer.Semicolon, "Expect ';' after loop condition.") {
		return nil
	}

	var increment object.Expr
	if !p.check(lexer.RParen) {
		increment = p.expression()
	}
	if !p.consume(lexer.RParen, "Expect ')' after for clauses.") {
		return nil
	}

	body := p.statement()

	return &object.ForStmt{Initializer: init, Condition: cond, Body: body, Increment: increment}
}

func (p *Parser) ifStatement() object.Stmt {
	if !p.consume(lexer.LParen, "Expect '(' after 'if'.") {
		return nil
	}

	cond := p.expression()

	if !p.consume(lexer.RParen, "Expect ')' after if condition.") {
		return nil
	}

	thenBranch := p.statement()
	var elseBranch object.Stmt

	if p.match(lexer.Else) {
		elseBranch = p.statement()
	}

	return &object.IfStmt{Condition: cond, Then: thenBranch, Else: elseBranch}
}

func (p *Parser) whileStatement() object.Stmt {
	if !p.consume(lexer.LParen, "Expect '(' after 'while'.") {
		return nil
	}

	cond := p.expression()
	if !p.consume(lexer.RParen, "Expect ')' after while condition.") {
		return nil
	}

	body := p.statement()

	return &object.ForStmt{Condition: cond, Body: body}
}

func (p *Parser) expressionStatement() object.Stmt {
	expr := p.expression()
	if !p.consume(lexer.Semicolon, "Expect ';' after value.") {
		return nil
	}
	return &object.ExpressionStmt{Expression: expr}
}

func (p *Parser) expression() object.Expr {
	return p.assignment()
}

func (p *Parser) assignment() object.Expr {
	expr := p.or()

	if p.match(lexer.Equal) {
		equals := p.prevTok
		value := p.assignment()

		switch e := expr.(type) {
		case *object.VariableExpr:
			return &object.AssignExpr{Name: e.Name, Value: value}
			// TODO: Add case for GetExpr and IndexExpr
		}

		p.addError(equals, "Invalid assignment target.")
		return nil
	}

	if p.match(lexer.MinusEq, lexer.PlusEq, lexer.StarEq, lexer.SlashEq, lexer.PercentEq, lexer.TildSlashEq) {
		equals := p.prevTok
		value := p.assignment()

		var oper *lexer.Token
		switch equals.Type {
		case lexer.MinusEq:
			oper = lexer.NewToken(lexer.Minus, "-", equals.Filename, equals.Line)
		case lexer.PlusEq:
			oper = lexer.NewToken(lexer.Plus, "+", equals.Filename, equals.Line)
		case lexer.StarEq:
			oper = lexer.NewToken(lexer.Star, "*", equals.Filename, equals.Line)
		case lexer.SlashEq:
			oper = lexer.NewToken(lexer.Slash, "/", equals.Filename, equals.Line)
		case lexer.PercentEq:
			oper = lexer.NewToken(lexer.Percent, "%", equals.Filename, equals.Line)
		case lexer.TildSlashEq:
			oper = lexer.NewToken(lexer.TildSlash, "~/", equals.Filename, equals.Line)
		default:
			return nil
		}

		be := &object.BinaryExpr{Left: expr, Operator: oper, Right: value}
		switch e := expr.(type) {
		case *object.VariableExpr:
			return &object.AssignExpr{Name: e.Name, Value: be}
		}
	}

	return expr
}

func (p *Parser) or() object.Expr {
	expr := p.and()

	for p.match(lexer.Or) {
		oper := p.prevTok
		right := p.and()
		expr = &object.LogicalExpr{Left: expr, Operator: oper, Right: right}
	}

	return expr
}

func (p *Parser) and() object.Expr {
	expr := p.equality()

	for p.match(lexer.And) {
		oper := p.prevTok
		right := p.equality()
		expr = &object.LogicalExpr{Left: expr, Operator: oper, Right: right}
	}

	return expr
}

func (p *Parser) equality() object.Expr {
	expr := p.comparison()

	for p.match(lexer.BangEq, lexer.EqualEq) {
		if expr == nil {
			return nil
		}

		oper := p.prevTok

		right := p.comparison()
		if right == nil {
			return nil
		}

		expr = &object.BinaryExpr{Left: expr, Operator: oper, Right: right}
	}

	return expr
}

func (p *Parser) comparison() object.Expr {
	expr := p.addition()

	for p.match(lexer.Greater, lexer.GreaterEq, lexer.Less, lexer.LessEq) {
		if expr == nil {
			return nil
		}

		oper := p.prevTok
		right := p.addition()
		if right == nil {
			return nil
		}

		expr = &object.BinaryExpr{Left: expr, Operator: oper, Right: right}
	}

	return expr
}

func (p *Parser) addition() object.Expr {
	expr := p.multiplication()

	for p.match(lexer.Plus, lexer.Minus) {
		if expr == nil {
			return nil
		}

		oper := p.prevTok
		right := p.multiplication()
		if right == nil {
			return nil
		}

		expr = &object.BinaryExpr{Left: expr, Operator: oper, Right: right}
	}

	return expr
}

func (p *Parser) multiplication() object.Expr {
	expr := p.unary()

	for p.match(lexer.Star, lexer.Slash, lexer.Percent, lexer.TildSlash) {
		if expr == nil {
			return nil
		}

		oper := p.prevTok
		right := p.unary()
		if right == nil {
			return nil
		}

		expr = &object.BinaryExpr{Left: expr, Operator: oper, Right: right}
	}

	return expr
}

func (p *Parser) unary() object.Expr {
	if p.match(lexer.Bang, lexer.Minus) {
		oper := p.prevTok
		right := p.unary()
		if right == nil {
			return nil
		}
		return &object.UnaryExpr{Operator: oper, Right: right}
	}

	return p.index()
}

func (p *Parser) index() object.Expr {
	expr := p.primary()

	for p.match(lexer.LBracket) {
		if expr == nil {
			return nil
		}

		oper := p.prevTok
		right := p.expression()
		if !p.consume(lexer.RBracket, "Expect ']' after index.") {
			return nil
		}

		expr = &object.IndexExpr{Left: expr, Operator: oper, Right: right}
	}

	return expr
}

func (p *Parser) primary() object.Expr {
	switch {
	case p.match(lexer.False, lexer.True):
		return &object.BooleanExpr{Token: p.prevTok, Value: p.prevTok.Type == lexer.True}
	case p.match(lexer.Ident):
		return &object.VariableExpr{Name: p.prevTok}
	case p.match(lexer.Null):
		return &object.NullExpr{Token: p.prevTok, Value: nil}
	case p.match(lexer.NumberF, lexer.NumberI):
		return p.parseNumber()
	case p.match(lexer.String, lexer.RawString):
		return &object.StringExpr{Token: p.prevTok, Value: p.prevTok.Lexeme}
	case p.match(lexer.UTString):
		p.addError(p.prevTok, "Unterminated string.")
		return nil
	case p.match(lexer.LBracket):
		var vals []object.Expr

		if !p.check(lexer.RBracket) {
			vals = append(vals, p.expression())
			for p.match(lexer.Comma) {
				vals = append(vals, p.expression())
			}
		}

		if !p.consume(lexer.RBracket, "Expect ']' after list values.") {
			return nil
		}
		return &object.ListExpr{Values: vals}
	case p.match(lexer.LParen):
		exp := p.expression()
		if exp == nil {
			return nil
		}
		if p.consume(lexer.RParen, "Expect ')' after expression.") {
			return &object.GroupingExpr{Expression: exp}
		}
	}

	p.addError(p.curTok, "Expect expression.")
	return nil
}

func (p *Parser) parseNumber() *object.NumberExpr {
	tok := p.prevTok
	if tok.Type == lexer.NumberF {
		n, err := strconv.ParseFloat(p.prevTok.Lexeme, 64)
		if err != nil {
			p.addError(p.prevTok, "Unable to parse value: "+p.prevTok.Lexeme+".")
			return nil
		}
		return &object.NumberExpr{Token: tok, Float: n}
	}

	n, err := strconv.ParseInt(p.prevTok.Lexeme, 10, 64)
	if err != nil {
		p.addError(p.prevTok, "Unable to parse value: "+p.prevTok.Lexeme+".")
		return nil
	}
	return &object.NumberExpr{Token: p.prevTok, Int: int(n)}
}

func (p *Parser) synchronize() {
	p.nextToken()

	for p.curTok.Type != lexer.EOF {
		if p.prevTok.Type == lexer.Semicolon {
			return
		}

		switch p.curTok.Type {
		case lexer.Class:
		case lexer.Fn:
		case lexer.Var:
		case lexer.For:
		case lexer.If:
		case lexer.While:
		case lexer.Print:
		case lexer.Return:
			return
		}

		p.nextToken()
	}
}
