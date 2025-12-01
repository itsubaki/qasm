package visitor

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type ErrorListener struct {
	antlr.ErrorListener
	Errors []error
}

func (l *ErrorListener) SyntaxError(
	recognizer antlr.Recognizer,
	offendingSymbol any,
	line int,
	column int,
	msg string,
	e antlr.RecognitionException,
) {
	l.Errors = append(l.Errors, fmt.Errorf("syntax error at line:%d:%d: %s", line, column, msg))
}
