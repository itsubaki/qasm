package visitor_test

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/environ"
	"github.com/itsubaki/qasm/parser"
	"github.com/itsubaki/qasm/visitor"
)

func ExampleWithMaxQubits() {
	v := visitor.New(
		q.New(),
		environ.New(),
		visitor.WithMaxQubits(5),
	)

	program, err := parser.Parse(`qubit[10] q;`)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := v.Run(program); err != nil {
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

	program, err := parser.Parse(`qreg q[10];`)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := v.Run(program); err != nil {
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

	program, err := parser.Parse(`qubit[10] q;`)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := v.Run(program); err != nil {
		fmt.Println(err)
	}

	fmt.Println(env.Qubit)

	// Output:
	// map[q:[0 1 2 3 4 5 6 7 8 9]]
}
