package visitor_test

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func ExampleVisitor_VisitQuantumDeclarationStatement() {
	text := "OPENQASM 3.0;qubit q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	v := visitor.New(qsim)
	tree.Accept(v)

	fmt.Println(qsim.NumberOfBit())
	fmt.Println(qsim.M().IsZero())
	fmt.Println(v.Environ.Qubit)

	// Output:
	// (program (version OPENQASM 3.0 ;) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) <EOF>)
	// 1
	// true
	// map[q:0]
}
