package formatter_test

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/formatter"
	"github.com/itsubaki/qasm/gen/parser"
)

func ExampleNew() {
	text := `
	OPENQASM 3.0;include "../testdata/stdgates.qasm";
	qubit q0;qubit[2] q1;qubit[2] q2;
	x q1;h q2;
	`

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	formatter.New(antlr.NewCommonTokenStream(lexer, antlr.TokenHiddenChannel))

	// Output:
}
