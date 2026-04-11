package listener_test

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/listener"
)

func ExampleErrorListener() {
	text := `// line 1
OPENQASM 3.0;
qubit q0 q1;
`

	l := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(l, antlr.TokenDefaultChannel))
	listener := &listener.ErrorListener{}
	p.RemoveErrorListeners()     // remove default error listeners
	p.AddErrorListener(listener) // add custom error listener

	_ = p.Program()
	for _, err := range listener.Errors {
		fmt.Println(err)
	}

	// Output:
	// 3:9: extraneous input 'q1' expecting ';'
}
