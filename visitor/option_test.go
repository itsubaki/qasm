package visitor_test

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/environ"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func ExampleWithMaxQubits() {
	v := visitor.New(
		q.New(),
		environ.New(),
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

func ExampleWithMaxQubits_oldstyle() {
	v := visitor.New(
		q.New(),
		environ.New(),
		visitor.WithMaxQubits(5),
	)

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(`qreg q[10];`))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	if err := v.Run(p.Program()); err != nil {
		fmt.Println(err)
	}

	// Output:
	// need=10, max=5: too many qubits
}

func ExampleWithMaxQubits_unlimited() {
	env := environ.New()
	v := visitor.New(
		q.New(),
		env,
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
