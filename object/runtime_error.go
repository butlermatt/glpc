package object

import (
	"fmt"
	"github.com/butlermatt/glpc/lexer"
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
