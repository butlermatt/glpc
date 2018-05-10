package object

import "github.com/butlermatt/glpc/lexer"

const __GlobalEnv__ = "__global__"

var fileScopes map[string]*Environment = make(map[string]*Environment)

func GetGlobal() *Environment {
	if env, ok := fileScopes[__GlobalEnv__]; ok {
		return env
	}

	env := NewEnvironment(__GlobalEnv__)
	// TODO: Populate builtin functions.
	return env
}

type Environment struct {
	parent *Environment
	m      map[string]Object
}

func NewEnvironment(filename string) *Environment {
	env := &Environment{m: make(map[string]Object)}
	fileScopes[filename] = env
	return env
}

func NewEnclosedEnvironment(enclosing *Environment) *Environment {
	return &Environment{parent: enclosing, m: make(map[string]Object)}
}

func (e *Environment) Define(name *lexer.Token, value Object) error {
	if _, ok := e.m[name.Lexeme]; ok {
		return NewRuntimeError(name, "Variable has already been declared.")
	}

	e.m[name.Lexeme] = value
	return nil
}

func (e *Environment) Get(name *lexer.Token) (Object, error) {
	v, ok := e.m[name.Lexeme]
	if ok {
		return v, nil
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	return nil, NewRuntimeError(name, "Undefined variable.")
}

func (e *Environment) GetAt(distance int, name *lexer.Token) (Object, error) {
	env := e
	for i := 0; i < distance; i++ {
		env = e.parent
	}

	return env.Get(name)
}

func (e *Environment) Assign(name *lexer.Token, value Object) error {
	if _, ok := e.m[name.Lexeme]; ok {
		e.m[name.Lexeme] = value
	}

	if e.parent != nil {
		return e.parent.Assign(name, value)
	}

	return NewRuntimeError(name, "Undefined variable.")
}

func (e *Environment) AssignAt(distance int, name *lexer.Token, value Object) error {
	env := e
	for i := 0; i < distance; i++ {
		env = env.parent
	}

	return env.Assign(name, value)
}
