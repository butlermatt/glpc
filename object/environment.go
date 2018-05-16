package object

import (
	"github.com/butlermatt/glpc/lexer"
)

const __GlobalEnv__ = "__global__"

var fileScopes = make(map[string]*FileEnvironment)

func GetGlobal() *FileEnvironment {
	if env, ok := fileScopes[__GlobalEnv__]; ok {
		return env
	}

	env := NewFileEnvironment(__GlobalEnv__)
	return env
}

func GetFileEnvironment(filename string) *FileEnvironment {
	return fileScopes[filename]
}

type FileEnvironment struct {
	filename string
	depth    map[Expr]int
	env      *Environment
}

func (fe *FileEnvironment) Env() *Environment {
	return fe.env
}

func (fe *FileEnvironment) Name() string {
	return fe.filename
}

type Environment struct {
	parent *Environment
	m      map[string]Object
}

func NewFileEnvironment(filename string) *FileEnvironment {
	env := &Environment{m: make(map[string]Object)}
	fe := &FileEnvironment{filename: filename, env: env}
	fileScopes[filename] = fe
	return fe
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

func (e *Environment) DefineString(name string, value Object) {
	e.m[name] = value
}

func (e *Environment) GetString(name string) Object {
	return e.m[name]
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

func (e *Environment) Copy(other *Environment) {
	for lex, obj := range other.m {
		e.m[lex] = obj
	}
}
