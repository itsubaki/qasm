package visitor_test

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func ExampleErrorListener() {
	text := `// line 1
OPENQASM 3.0;
qubit q0 q1;
`

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	errListener := &visitor.ErrorListener{}
	lexer.AddErrorListener(errListener)
	p.AddErrorListener(errListener)

	_ = p.Program()
	for _, e := range errListener.Errors {
		fmt.Println(e)
	}

	// Output:
	// syntax error at line:3:9: extraneous input 'q1' expecting ';'
}
