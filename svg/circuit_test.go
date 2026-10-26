package svg_test

import (
	"fmt"

	"github.com/itsubaki/qasm/svg"
)

func ExampleWire_Label() {
	index := 0
	fmt.Println((svg.Wire{Name: "q"}).Label())
	fmt.Println((svg.Wire{Name: "q", Index: &index}).Label())

	// Output:
	// q
	// q[0]
}

func ExampleGate_WireIDs() {
	g := &svg.Gate{
		Name:    "CNOT",
		Control: []svg.WireID{0},
		Target:  []svg.WireID{1},
	}

	fmt.Println(g.WireIDs())

	// Output:
	// [0 1]
}

func ExampleSubroutine_WireIDs() {
	s := &svg.Subroutine{
		Name: "qft",
		Wire: []svg.WireID{0, 1},
	}

	fmt.Println(s.WireIDs())

	// Output:
	// [0 1]
}

func ExampleMeasurement_WireIDs() {
	m := &svg.Measurement{
		Wire: []svg.WireID{0},
	}
	fmt.Println(m.WireIDs())

	m01 := &svg.Measurement{
		Wire:   []svg.WireID{0},
		Target: []svg.WireID{1},
	}
	fmt.Println(m01.WireIDs())

	// Output:
	// [0]
	// [0 1]
}

func ExampleBarrier_WireIDs() {
	b := &svg.Barrier{}
	fmt.Println(b.WireIDs())

	// Output:
	// []
}
