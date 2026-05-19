package svg_test

import (
	"fmt"

	"github.com/itsubaki/qasm/svg"
)

func ExampleGate_Wires() {
	g := &svg.Gate{
		Name:     "CNOT",
		Controls: []int{0},
		Targets:  []int{1},
	}

	fmt.Println(g.Wires())

	// Output:
	// [0 1]
}

func ExampleMeasurement_Wires() {
	m := &svg.Measurement{
		Wire: []int{0},
	}
	fmt.Println(m.Wires())

	m01 := &svg.Measurement{
		Wire:   []int{0},
		Target: []int{1},
	}
	fmt.Println(m01.Wires())

	// Output:
	// [0]
	// [0 1]
}

func ExampleBarrier_Wires() {
	b := &svg.Barrier{}
	fmt.Println(b.Wires())

	// Output:
	// []
}
