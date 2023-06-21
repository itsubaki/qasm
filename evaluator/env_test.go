package evaluator_test

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/ast"
	"github.com/itsubaki/qasm/evaluator"
)

func ExampleQubit_All() {
	qb := &evaluator.Qubit{
		Name:  make([]string, 0),
		Value: make(map[string][]q.Qubit),
	}

	qb.Add(&ast.GenDecl{Name: "q0"}, []q.Qubit{1, 2, 3})
	qb.Add(&ast.GenDecl{Name: "q1"}, []q.Qubit{4, 5, 6})
	qb.Add(&ast.GenDecl{Name: "q2"}, []q.Qubit{7, 8, 9})

	fmt.Println(qb.All())

	// Output:
	// [1 2 3 4 5 6 7 8 9] <nil>
}
