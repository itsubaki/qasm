package listener

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type ErrorListener struct {
	antlr.DefaultErrorListener
	Errors []error
}

func (l *ErrorListener) SyntaxError(
	_ antlr.Recognizer,
	_ any, // offending symbol
	line int,
	column int,
	msg string,
	_ antlr.RecognitionException,
) {
	l.Errors = append(l.Errors, fmt.Errorf("syntax error at line:%d:%d: %s", line, column, msg))
}

func NewErrorListener(r ...antlr.Recognizer) *ErrorListener {
	// remove default error listeners
	for _, v := range r {
		v.RemoveErrorListeners()
	}

	// add custom error listener
	listener := &ErrorListener{}
	for _, v := range r {
		v.AddErrorListener(listener)
	}

	return listener
}
