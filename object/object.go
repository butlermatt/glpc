package object

// Type indicates the type of object the base object represents.
type Type int

// Object is a runtime representation of a GLPC Object
type Object interface {
	Type() Type
	String() string
}

const (
	Null Type = iota
	Boolean
	BuiltIn
	Class
	Function
	Instance
	List
	Number
	Return
	String
	Printer
)
