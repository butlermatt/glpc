package parser

import (
	"testing"

	"github.com/butlermatt/glpc/lexer"
)

func TestResolver_BeginEnd(t *testing.T) {
	r := NewResolver()
	r.Begin()

	if len(r.stack) != 1 {
		t.Errorf("wrong number of scopes on stack. expected=%d, got=%d", 1, len(r.stack))
	}

	r.Begin()
	r.Begin()
	if len(r.stack) != 3 {
		t.Errorf("wrong number of scopes on stack. expected=%d, got=%d", 3, len(r.stack))
	}

	r.End()
	if len(r.stack) != 2 {
		t.Errorf("wrong number of scopes on stack. expected=%d, got=%d", 2, len(r.stack))
	}
	r.End()
	if len(r.stack) != 1 {
		t.Errorf("wrong number of scopes on stack. expected=%d, got=%d", 1, len(r.stack))
	}
	r.End()
	if len(r.stack) != 0 {
		t.Errorf("wrong number of scopes on stack. expected=%d, got=%d", 0, len(r.stack))
	}

	r.End()
	if len(r.stack) != 0 {
		t.Errorf("wrong number of scopes on stack. expected=%d, got=%d", 0, len(r.stack))
	}
}

func TestResolver_Declare(t *testing.T) {
	r := NewResolver()
	r.Begin()

	if len(r.stack) != 1 {
		t.Fatalf("wrong number of scopes on stack. expected=%d, got=%d", 1, len(r.stack))
	}

	name := lexer.NewToken(lexer.Ident, "x", "testfile.gpc", 1)
	err := r.Declare(name)
	if err != nil {
		t.Fatalf("error declaring variable. %v", err)
	}

	p := r.Peek()
	if p == nil {
		t.Fatalf("current scope should not be nil.")
	}

	v, ok := p[name.Lexeme]
	if !ok {
		t.Errorf("unable to locate value in scope.")
	}

	if v {
		t.Errorf("value is incorrect. expected=%t, got=%t", false, v)
	}

	name2 := lexer.NewToken(lexer.Ident, "y", "testfile.gpc", 2)
	v, ok = p[name2.Lexeme]
	if ok {
		t.Errorf("located value in scope that has not been added.")
	}
	if v {
		t.Errorf("value is incorrect. expected=%t, got=%t", false, v)
	}

	r.Begin()
	err = r.Declare(name2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p2 := r.Peek()
	_, ok = p2[name.Lexeme]
	if ok {
		t.Errorf("unexpected value in scope. Should be in prior scope")
	}

	_, ok = p2[name2.Lexeme]
	if !ok {
		t.Errorf("value not located in scope when expected.")
	}

	r.End()
	r.End()
}

func TestResolver_Define(t *testing.T) {
	r := NewResolver()
	r.Begin()

	name := lexer.NewToken(lexer.Ident, "x", "testfile.gpc", 1)
	err := r.Declare(name)

	if err != nil {
		t.Errorf("unexpected error declaring variable: %v", err)
	}

	p := r.Peek()
	v, ok := p[name.Lexeme]
	if !ok {
		t.Errorf("unable to load expected value from scope")
	}

	if v {
		t.Errorf("value is showing defined as well as declared.")
	}

	r.Define(name)

	v = p[name.Lexeme]
	if !v {
		t.Errorf("name is not defined in scope when it should be.")
	}

	r.Begin()
	p = r.Peek()
	v, ok = p[name.Lexeme]
	if ok {
		t.Errorf("name is declared in scope when it should not exist.")
	}

	if v {
		t.Errorf("name is defined in scope when it should not exist.")
	}

	r.End()
	r.End()
}
