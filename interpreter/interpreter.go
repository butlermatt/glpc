package interpreter

import (
	"github.com/butlermatt/glpc/object"
)

type Interpreter struct {
	stmts   []object.Stmt
	local   map[object.Expr]int
	env     *object.Environment
	globals *object.Environment
}

func New(statements []object.Stmt, depthMap map[object.Expr]int, file string) *Interpreter {
	env := object.NewEnvironment(file)
	glob := object.GetGlobal()
	return &Interpreter{stmts: statements, local: depthMap, env: env, globals: glob}
}
