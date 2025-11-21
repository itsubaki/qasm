package formatter_test

import (
	"fmt"

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
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.Newqasm3Parser(stream)

	tree := p.Program()
	f := formatter.New(stream)
	antlr.ParseTreeWalkerDefault.Walk(f, tree)

	result := f.Format()
	fmt.Println(result)

	// Output:
	// OPENQASM 3.0;
	// include "../testdata/stdgates.qasm";
	// qubit q0;
	// qubit[2] q1;
	// qubit[2] q2;
	// x q1;
	// h q2;
}
