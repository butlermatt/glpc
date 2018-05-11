package interpreter

import "bytes"
import (
	"fmt"
	"github.com/butlermatt/glpc/object"
)

type Null struct{}

func (n *Null) Type() object.Type { return object.Null }
func (n *Null) String() string    { return "null" }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() object.Type { return object.Boolean }
func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type List struct {
	Elements []object.Object
}

func (l *List) Type() object.Type { return object.List }
func (l *List) String() string {
	var out bytes.Buffer

	if len(l.Elements) == 0 {
		out.WriteString("[]")
	} else if len(l.Elements) == 1 {
		out.WriteByte('[')
		out.WriteString(l.Elements[0].String())
		out.WriteByte(']')
	} else {
		out.WriteByte('[')
		out.WriteString(l.Elements[0].String())
		for i := 1; i < len(l.Elements); i++ {
			out.WriteString(", ")
			out.WriteString(l.Elements[i].String())
		}

		out.WriteByte(']')
	}

	return out.String()
}

type Number struct {
	IsInt bool
	Int   int
	Float float64
}

func (n *Number) Type() object.Type { return object.Number }
func (n *Number) String() string {
	if n.IsInt {
		return fmt.Sprintf("%d", n.Int)
	}
	return fmt.Sprintf("%.2f", n.Float)
}

type String struct {
	Value string
}

func (s *String) Type() object.Type { return object.String }
func (s *String) String() string    { return s.Value }
