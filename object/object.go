package object

import (
	"fmt"
	"github.com/butlermatt/glpc/lexer"
)

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
	String
	Printer
)

type RuntimeError struct {
	Token   *lexer.Token
	Message string
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("[Runtime Error] - line %d at %q - %s", re.Token.Line, re.Token.Lexeme, re.Message)
}

func NewRuntimeError(token *lexer.Token, msg string) *RuntimeError {
	return &RuntimeError{Token: token, Message: msg}
}
