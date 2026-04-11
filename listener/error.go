package listener

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type SyntaxError struct {
	Line    int
	Column  int
	Message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Column, e.Message)
}

type ErrorListener struct {
	antlr.DefaultErrorListener
	Errors []*SyntaxError
}

func (l *ErrorListener) SyntaxError(
	_ antlr.Recognizer,
	_ any, // offending symbol
	line int,
	column int,
	msg string,
	_ antlr.RecognitionException,
) {
	l.Errors = append(l.Errors, &SyntaxError{
		Line:    line,
		Column:  column,
		Message: msg,
	})
}

func NewErrorListener(r ...antlr.Recognizer) *ErrorListener {
	listener := &ErrorListener{}
	for _, v := range r {
		v.AddErrorListener(listener)
	}

	return listener
}
