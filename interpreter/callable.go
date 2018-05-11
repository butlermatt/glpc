package interpreter

import "github.com/butlermatt/glpc/object"

type Callable interface {
	// Arity is the number of expected arguments
	Arity() int
	Call(interpreter *Interpreter, args []object.Object) (object.Object, error)
}

type Function struct {
	declaration *object.FunctionStmt
	closure     *object.Environment
	isInit      bool
}

func NewFunction(declaration *object.FunctionStmt, env *object.Environment, isInit bool) *Function {
	return &Function{declaration: declaration, closure: env, isInit: isInit}
}

func (f *Function) Type() object.Type { return object.Function }
func (f *Function) String() string    { return "<fn " + f.declaration.Name.Lexeme + ">" }
func (f *Function) Arity() int        { return len(f.declaration.Parameters) }
func (f *Function) Call(interpreter *Interpreter, args []object.Object) (object.Object, error) {
	if len(args) != f.Arity() {
		return nil, object.NewRuntimeError(f.declaration.Name, "Incorrect number of arguments passed.")
	}

	env := object.NewEnclosedEnvironment(f.closure)
	for i, p := range f.declaration.Parameters {
		env.Define(p, args[i])
	}

	err := interpreter.executeBlock(f.declaration.Body, env)
	if err != nil {
		if e, ok := err.(*ReturnError); ok {
			return e.Value, nil
		}
		return nil, err
	}

	if f.isInit {
		return f.closure.GetString("this"), nil
	}

	return NullOb, nil
}

func (f *Function) Bind(inst *Instance) *Function {
	env := object.NewEnclosedEnvironment(f.closure)
	env.DefineString("this", inst)
	return NewFunction(f.declaration, env, f.isInit)
}
