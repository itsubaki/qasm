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
	errListener := listener.NewErrorListener(l, p)

	_ = p.Program()
	for _, e := range errListener.Errors {
		fmt.Println(e)
	}

	// Output:
	// syntax error at line:3:9: extraneous input 'q1' expecting ';'
}
