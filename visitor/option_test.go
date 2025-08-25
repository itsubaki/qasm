package visitor_test

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func ExampleWithMaxQubits() {
	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env,
		visitor.WithMaxQubits(5),
	)

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(`qubit[10] q;`))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	if err := v.Run(p.Program()); err != nil {
		fmt.Println(err)
	}

	// Output:
	// need=10, max=5: too many qubits
}

func ExampleWithMaxQubits_unlimited() {
	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env,
		visitor.WithMaxQubits(0),
	)

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(`qubit[10] q;`))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	if err := v.Run(p.Program()); err != nil {
		fmt.Println(err)
	}

	fmt.Println(env.Qubit)

	// Output:
	// map[q:[0 1 2 3 4 5 6 7 8 9]]
}
