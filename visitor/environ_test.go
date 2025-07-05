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

func ExampleEnviron_SetVariable() {
	env := visitor.NewEnviron()

	env.SetVariable("a", 42)
	fmt.Println("env:", env.Variable)

	enclosed := env.NewEnclosed()
	enclosed.SetVariable("a", 43)
	fmt.Println("env:", env.Variable)
	fmt.Println("enclosed:", enclosed.Variable)

	enclosed.SetVariable("b", 100)
	enclosed.SetVariable("b", 101)
	fmt.Println("env:", env.Variable)
	fmt.Println("enclosed:", enclosed.Variable)

	// Output:
	// env: map[a:42]
	// env: map[a:43]
	// enclosed: map[]
	// env: map[a:43]
	// enclosed: map[b:101]
}

func ExampleEnviron_GetGate() {
	env := visitor.NewEnviron()
	_, ok := env.GetGate("x")
	fmt.Println(ok)

	// Output:
	// false
}

func ExampleEnviron_GetSubroutine() {
	env := visitor.NewEnviron()
	env.Subroutine["x"] = &visitor.Subroutine{
		Name: "x",
	}

	enclosed := env.NewEnclosed()
	sub, ok := enclosed.GetSubroutine("x")
	fmt.Println(sub.Name, ok)

	// Output:
	// x true
}
