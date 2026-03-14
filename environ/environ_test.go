package environ_test

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/environ"
)

func ExampleEnviron_NewEnclosed() {
	env := environ.New()
	env.Qubit["q0"] = []q.Qubit{0, 1}

	enclosed := env.NewEnclosed()
	enclosed.Qubit["q0"] = []q.Qubit{2, 3}
	enclosed.Qubit["q1"] = []q.Qubit{4, 5}

	fmt.Println(enclosed.GetQubit("q0"))
	fmt.Println(enclosed.GetQubit("q1"))

	// Output:
	// [2 3] true
	// [4 5] true
}

func ExampleEnviron_SetVariable() {
	env := environ.New()

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

func ExampleEnviron_GetConst() {
	env := environ.New()
	env.Const["c0"] = 42

	enclosed := env.NewEnclosed()
	enclosed.Const["c1"] = 43

	fmt.Println(enclosed.GetConst("not found"))
	fmt.Println(enclosed.GetConst("c0"))
	fmt.Println(enclosed.GetConst("c1"))

	// Output:
	// <nil> false
	// 42 true
	// 43 true
}

func ExampleEnviron_GetVariable() {
	env := environ.New()
	env.SetVariable("v0", 42)

	enclosed := env.NewEnclosed()
	enclosed.SetVariable("v1", 43)

	fmt.Println(enclosed.GetVariable("not found"))
	fmt.Println(enclosed.GetVariable("v0"))
	fmt.Println(enclosed.GetVariable("v1"))

	// Output:
	// <nil> false
	// 42 true
	// 43 true
}

func ExampleEnviron_GetQubit() {
	env := environ.New()
	env.Qubit["q"] = []q.Qubit{0, 1, 2}

	enclosed := env.NewEnclosed()
	enclosed.Qubit["r"] = []q.Qubit{3, 4}

	fmt.Println(enclosed.GetQubit("not found"))
	fmt.Println(enclosed.GetQubit("q"))
	fmt.Println(enclosed.GetQubit("r"))

	// Output:
	// [] false
	// [0 1 2] true
	// [3 4] true
}

func ExampleEnviron_GetClassicalBit() {
	env := environ.New()
	env.ClassicalBit["c"] = []bool{true, false}

	enclosed := env.NewEnclosed()
	enclosed.ClassicalBit["d"] = []bool{false, true}

	fmt.Println(enclosed.GetClassicalBit("not found"))
	fmt.Println(enclosed.GetClassicalBit("c"))
	fmt.Println(enclosed.GetClassicalBit("d"))

	// Output:
	// [] false
	// [true false] true
	// [false true] true
}

func ExampleEnviron_GetGate() {
	env := environ.New()
	env.Gate["x"] = &environ.Gate{
		Name: "x",
	}

	enclosed := env.NewEnclosed()
	enclosed.Gate["y"] = &environ.Gate{
		Name: "y",
	}

	x, xok := enclosed.GetGate("x")
	y, yok := enclosed.GetGate("y")

	fmt.Println(enclosed.GetGate("not found"))
	fmt.Println(x.Name, xok)
	fmt.Println(y.Name, yok)

	// Output:
	// <nil> false
	// x true
	// y true
}

func ExampleEnviron_GetSubroutine() {
	env := environ.New()
	env.Subroutine["qft"] = &environ.Subroutine{
		Name: "qft",
	}

	enclosed := env.NewEnclosed()
	enclosed.Subroutine["swap"] = &environ.Subroutine{
		Name: "swap",
	}

	qft, qftok := enclosed.GetSubroutine("qft")
	swap, swapok := enclosed.GetSubroutine("swap")

	fmt.Println(enclosed.GetSubroutine("not found"))
	fmt.Println(qft.Name, qftok)
	fmt.Println(swap.Name, swapok)

	// Output:
	// <nil> false
	// qft true
	// swap true
}

func ExampleEnviron_Index() {
	env := environ.New()
	env.SetQubit("q0", []q.Qubit{0, 1})
	env.SetQubit("q1", []q.Qubit{2, 3, 4})

	index := env.Index()
	fmt.Println(index)

	// Output:
	// [[0 1] [2 3 4]]
}
