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
