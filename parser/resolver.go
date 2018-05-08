package parser

import (
	"github.com/butlermatt/glpc/object"
	"github.com/butlermatt/glpc/lexer"
)

type Resolver struct {
	stack []map[string]bool
	dist  map[object.Expr]int
}

func NewResolver() *Resolver {
	return &Resolver{dist: make(map[object.Expr]int)}
}

func (r *Resolver) Begin() {
	r.stack = append(r.stack, make(map[string]bool))
}

func (r *Resolver) End() {
	if len(r.stack) == 0 {
		return
	}

	r.stack = r.stack[:len(r.stack) - 1]
}

func (r *Resolver) Peek() map[string]bool {
	if len(r.stack) == 0 {
		return nil
	}

	return r.stack[len(r.stack) - 1]
}

func (r *Resolver) Declare(name *lexer.Token) error {
	scope := r.Peek()
	if scope == nil {
		return nil
	}

	if _, ok := scope[name.Lexeme]; ok {
		return ParseError{Line: name.Line, Where: name.Lexeme, Msg: "Variable with this name already declared in this scope."}
	}

	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) Define(name *lexer.Token) {
	scope := r.Peek()
	if scope == nil {
		return
	}

	scope[name.Lexeme] = true
}

func (r *Resolver) Local(expr object.Expr, name *lexer.Token) {
	for i := len(r.stack) - 1; i >= 0; i-- {
		if _, ok := r.stack[i][name.Lexeme]; ok {
			r.dist[expr] = len(r.stack) - 1 - i
			return
		}
	}

	// Not found, assume it's global.
}