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

func ExampleEnviron_SetClBit() {
	env := environ.New()

	env.SetClBit("a", true)
	fmt.Println("env:", env.ClBit)

	enclosed := env.NewEnclosed()
	enclosed.SetClBit("a", false)
	fmt.Println("env:", env.ClBit)
	fmt.Println("enclosed:", enclosed.ClBit)

	enclosed.SetClBit("b", true)
	enclosed.SetClBit("b", false)
	fmt.Println("env:", env.ClBit)
	fmt.Println("enclosed:", enclosed.ClBit)

	// Output:
	// env: map[a:true]
	// env: map[a:false]
	// enclosed: map[]
	// env: map[a:false]
	// enclosed: map[b:false]
}

func ExampleEnviron_SetClBitArray() {
	env := environ.New()

	env.SetClBitArray("a", []bool{true, false})
	fmt.Println("env:", env.ClBitArray)

	enclosed := env.NewEnclosed()
	enclosed.SetClBitArray("a", []bool{false, true})
	fmt.Println("env:", env.ClBitArray)
	fmt.Println("enclosed:", enclosed.ClBitArray)

	enclosed.SetClBitArray("b", []bool{true, true})
	enclosed.SetClBitArray("b", []bool{false, false})
	fmt.Println("env:", env.ClBitArray)
	fmt.Println("enclosed:", enclosed.ClBitArray)

	// Output:
	// env: map[a:[true false]]
	// env: map[a:[false true]]
	// enclosed: map[]
	// env: map[a:[false true]]
	// enclosed: map[b:[false false]]
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

func ExampleEnviron_GetClBit() {
	env := environ.New()
	env.ClBit["c"] = true

	enclosed := env.NewEnclosed()
	enclosed.ClBit["d"] = false

	fmt.Println(enclosed.GetClBit("not found"))
	fmt.Println(enclosed.GetClBit("c"))
	fmt.Println(enclosed.GetClBit("d"))

	// Output:
	// false false
	// true true
	// false true
}

func ExampleEnviron_GetClBitArray() {
	env := environ.New()
	env.ClBitArray["c"] = []bool{true, false}

	enclosed := env.NewEnclosed()
	enclosed.ClBitArray["d"] = []bool{false, true}

	fmt.Println(enclosed.GetClBitArray("not found"))
	fmt.Println(enclosed.GetClBitArray("c"))
	fmt.Println(enclosed.GetClBitArray("d"))

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
