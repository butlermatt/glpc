package parser

import (
	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
	"strconv"
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

func (p *Parser) statement() object.Stmt {
	switch {
	// TODO: Cases for break, continue, If, Print, Return, While, For, and LBrace
	}

	return p.expressionStatement()
}

func (p *Parser) expressionStatement() object.Stmt {
	expr := p.expression()
	if !p.consume(lexer.Semicolon, "Expect ';' after value.") {
		return nil
	}
	return &object.ExpressionStmt{Expression: expr}
}

func (p *Parser) expression() object.Expr {
	return p.unary()
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

	return p.primary()
}

func (p *Parser) primary() object.Expr {
	switch {
	case p.match(lexer.False, lexer.True):
		return &object.BooleanExpr{Token: p.prevTok, Value: p.prevTok.Type == lexer.True}
	case p.match(lexer.Null):
		return &object.NullExpr{Token: p.prevTok, Value: nil}
	case p.match(lexer.NumberF):
		n, err := strconv.ParseFloat(p.prevTok.Lexeme, 64)
		if err != nil {
			p.addError(p.prevTok, "Unable to parse value: "+p.prevTok.Lexeme+".")
			return nil
		}
		return &object.NumberExpr{Token: p.prevTok, Float: n}
	case p.match(lexer.NumberI):
		n, err := strconv.ParseInt(p.prevTok.Lexeme, 10, 64)
		if err != nil {
			p.addError(p.prevTok, "Unable to parse value: "+p.prevTok.Lexeme+".")
			return nil
		}
		return &object.NumberExpr{Token: p.prevTok, Int: int(n)}
	case p.match(lexer.String, lexer.RawString):
		return &object.StringExpr{Token: p.prevTok, Value: p.prevTok.Lexeme}
	case p.match(lexer.UTString):
		p.addError(p.prevTok, "unterminated string")
		return nil
	case p.match(lexer.LBracket):
		var vals []object.Expr

		if !p.check(lexer.RBracket) {
			vals = append(vals, p.expression())
			for p.match(lexer.Comma) {
				vals = append(vals, p.expression())
			}
		}

		if !p.consume(lexer.RBracket, "expect ']' after list values.") {
			return nil
		}
		return &object.ListExpr{Values: vals}
	case p.match(lexer.LParen):
		exp := p.expression()
		if exp == nil {
			return nil
		}
		if p.consume(lexer.RParen, "expect ')' after expression.") {
			return &object.GroupingExpr{Expression: exp}
		}
	}

	p.addError(p.curTok, "Expect expression.")
	return nil
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
