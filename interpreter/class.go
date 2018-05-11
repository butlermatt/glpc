package interpreter

import (
	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/object"
)

type Class struct {
	Name       string
	superclass *Class
	methods    map[string]*Function
}

func (c *Class) Type() object.Type { return object.Class }
func (c *Class) String() string    { return c.Name }

func (c *Class) Arity() int {
	init := c.methods["init"]
	if init == nil {
		return 0
	}
	return init.Arity()
}

func (c *Class) Call(interp *Interpreter, args []object.Object) (object.Object, error) {
	inst := &Instance{klass: c, fields: make(map[string]object.Object)}
	init := c.methods["init"]
	if init != nil {
		_, err := init.Bind(inst).Call(interp, args)
		if err != nil {
			return nil, err
		}
	}

	return inst, nil
}

func (c *Class) findMethod(inst *Instance, name string) *Function {
	method := c.methods[name]
	if method != nil {
		return method.Bind(inst)
	}

	if c.superclass != nil {
		return c.superclass.findMethod(inst, name)
	}

	return nil
}

type Instance struct {
	klass  *Class
	fields map[string]object.Object
}

func (in *Instance) Type() object.Type { return object.Instance }
func (in *Instance) String() string    { return in.klass.Name + " instance" }

func (in *Instance) Get(name *lexer.Token) (object.Object, error) {
	if v, ok := in.fields[name.Lexeme]; ok {
		return v, nil
	}

	if m := in.klass.findMethod(in, name.Lexeme); m != nil {
		return m, nil
	}

	return nil, object.NewRuntimeError(name, "Undefined property.")
}

func (in *Instance) Set(name *lexer.Token, value object.Object) {
	in.fields[name.Lexeme] = value
}
