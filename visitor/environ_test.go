package visitor_test

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/visitor"
)

func ExampleEnviron_NewEnclosed() {
	env := visitor.NewEnviron()
	env.Qubit["q0"] = []q.Qubit{0, 1}
	env.ClassicalBit["c0"] = []int64{0, 1}

	enclosed := env.NewEnclosed()
	enclosed.Qubit["q1"] = []q.Qubit{2, 3}
	enclosed.ClassicalBit["c1"] = []int64{2, 3}

	fmt.Println(enclosed.GetQubit("q0"))
	fmt.Println(enclosed.GetQubit("q1"))
	fmt.Println(enclosed.GetClassicalBit("c0"))
	fmt.Println(enclosed.GetClassicalBit("c1"))

	// Output:
	// [0 1] true
	// [2 3] true
	// [0 1] true
	// [2 3] true
}
