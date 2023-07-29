package evaluator_test

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/ast"
	"github.com/itsubaki/qasm/evaluator"
)

func ExampleBit() {
	b := &evaluator.Bit{
		Name:  make([]string, 0),
		Value: make(map[string][]int64),
	}
	b.Add(&ast.IdentExpr{Name: "a"}, []int64{1, 2, 3})
	b.Add(&ast.IdentExpr{Name: "b"}, []int64{4, 5, 6})
	fmt.Println(b)

	fmt.Println(b.Get(&ast.IndexExpr{Name: "c"}))
	fmt.Println(b.Get(&ast.IndexExpr{Name: "c", Value: "1"}))
	fmt.Println(b.Get(&ast.IndexExpr{Name: "a", Value: "0"}))
	fmt.Println(b.Get(&ast.IndexExpr{Name: "b", Value: "-1"}))
	fmt.Println(b.Get(&ast.IndexExpr{Name: "b", Value: "10"}))
	fmt.Println(b.Get(&ast.IndexExpr{Name: "b", Value: "-10"}))
	fmt.Println(b.Get(&ast.BadExpr{}))

	// Output:
	// [a b], map[a:[1 2 3] b:[4 5 6]]
	// [] false
	// [] false
	// [1] true
	// [6] true
	// [] false
	// [] false
	// [] false
}

func ExampleQubit() {
	qb := &evaluator.Qubit{
		Name:  make([]string, 0),
		Value: make(map[string][]q.Qubit),
	}

	qb.Add(&ast.GenDecl{Name: "q0"}, []q.Qubit{1, 2, 3})
	qb.Add(&ast.GenDecl{Name: "q1"}, []q.Qubit{4, 5, 6})
	qb.Add(&ast.GenDecl{Name: "q2"}, []q.Qubit{7, 8, 9})

	fmt.Println(qb)
	fmt.Println(qb.All())
	fmt.Println(qb.Get(&ast.IndexExpr{Name: "c"}))
	fmt.Println(qb.Get(&ast.IndexExpr{Name: "q0", Value: "10"}))
	fmt.Println(qb.Get(&ast.IndexExpr{Name: "q0", Value: "-10"}))
	fmt.Println(qb.Get(&ast.BadExpr{}))

	// Output:
	// [q0 q1 q2], map[q0:[1 2 3] q1:[4 5 6] q2:[7 8 9]]
	// [1 2 3 4 5 6 7 8 9] <nil>
	// [] false
	// [] false
	// [] false
	// [] false
}

func ExampleEnviron() {
	env := evaluator.NewEnviron()
	fmt.Println(env)

	// Output:
	// gate: map[], const: map[], bit: [], map[], qubit: [], map[], modifier: [], decl: []
}
