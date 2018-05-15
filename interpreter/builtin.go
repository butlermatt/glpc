package interpreter

import (
	"fmt"
	"github.com/butlermatt/glpc/object"
)

type CallFn func(interpreter *Interpreter, args []object.Object) (object.Object, error)

type BIError string

func (be BIError) Error() string {
	return string(be)
}

type BuiltIn struct {
	arity  int
	callFn CallFn
}

func (b *BuiltIn) Type() object.Type { return object.BuiltIn }
func (b *BuiltIn) String() string    { return "builtin function" }

func (b *BuiltIn) Arity() int { return b.arity }
func (b *BuiltIn) Call(interpreter *Interpreter, args []object.Object) (object.Object, error) {
	return b.callFn(interpreter, args)
}

func newBuiltin(arity int, fn CallFn) *BuiltIn {
	return &BuiltIn{arity: arity, callFn: fn}
}

func SetupGlobal(env *object.Environment) *object.Environment {
	if env.GetString("len") != nil {
		return env
	}

	env.DefineString("len", newBuiltin(1, bLen))
	env.DefineString("debugPrint", newBuiltin(-1, bDebugPrint))

	return env
}

func bLen(interp *Interpreter, args []object.Object) (object.Object, error) {
	switch obj := args[0]; obj.Type() {
	case object.String:
		s := obj.(*String)
		return &Number{IsInt: true, Int: len(s.Value)}, nil
	case object.List:
		l := obj.(*List)
		return &Number{IsInt: true, Int: len(l.Elements)}, nil
	}

	return NullOb, BIError("'len' argument must be of a type STRING or LIST.")
}

// TODO Remove this when I get something better
func bDebugPrint(inter *Interpreter, args []object.Object) (object.Object, error) {
	if len(args) < 1 {
		fmt.Println("")
		return NullOb, nil
	}

	fmt.Printf("%s", args[0].String())
	for i := 1; i < len(args); i++ {
		fmt.Printf(" %s", args[i].String())
	}
	fmt.Println("")
	return NullOb, nil
}
