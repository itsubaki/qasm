package svg_test

import (
	"fmt"

	"github.com/itsubaki/qasm/svg"
)

func ExampleRender() {
	circuit := &svg.Circuit{
		Wires: []svg.Wire{
			{Name: "q0"},
			{Name: "q1"},
			{Name: "c0"},
		},
		Ops: []svg.Op{
			&svg.Gate{
				Name:    "H",
				Targets: []int{0},
			},
			&svg.Gate{
				Name:     "CNOT",
				Controls: []int{0},
				Targets:  []int{1},
			},
			&svg.Barrier{},
			&svg.Subroutine{
				Name: "oracle",
				Wire: []int{0, 1},
			},
			&svg.Measurement{
				Wire:    []int{0},
				Targets: []int{2},
			},
		},
	}

	diagram := svg.Render(svg.NewLayout(circuit), svg.DefaultConfig)
	fmt.Printf("%s ... %s", diagram[:4], diagram[len(diagram)-6:])

	// Output:
	// <svg ... </svg>
}
