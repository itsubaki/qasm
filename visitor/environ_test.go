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
	env.Gate["x"] = &visitor.Gate{
		Name: "x",
	}

	g, ok := env.GetGate("x")
	fmt.Println(g.Name, ok)

	enclosed := env.NewEnclosed()
	encg, ok := enclosed.GetGate("x")
	fmt.Println(encg.Name, ok)

	// Output:
	// x true
	// x true
}

func ExampleEnviron_GetSubroutine() {
	env := visitor.NewEnviron()
	env.Subroutine["qft"] = &visitor.Subroutine{
		Name: "qft",
	}

	enclosed := env.NewEnclosed()
	sub, ok := enclosed.GetSubroutine("qft")
	fmt.Println(sub.Name, ok)

	// Output:
	// qft true
}

func ExampleEnviron_Index() {
	env := visitor.NewEnviron()
	env.SetQubit("q0", []q.Qubit{0, 1})
	env.SetQubit("q1", []q.Qubit{2, 3, 4})

	index := env.Index()
	fmt.Println(index)

	// Output:
	// [[0 1] [2 3 4]]
}
