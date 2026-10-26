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
				Name:   "H",
				Target: []svg.WireID{0},
			},
			&svg.Gate{
				Name:    "CNOT",
				Control: []svg.WireID{0},
				Target:  []svg.WireID{1},
			},
			&svg.Barrier{
				Wire: []svg.WireID{0, 1},
			},
			&svg.Barrier{
				Wire: []svg.WireID{0, 2},
			},
			&svg.Subroutine{
				Name: "oracle",
				Wire: []svg.WireID{0, 1},
			},
			&svg.Measurement{
				Wire:   []svg.WireID{0},
				Target: []svg.WireID{2},
			},
		},
	}

	diagram := svg.Render(svg.NewLayout(circuit), svg.DefaultConfig)
	fmt.Printf("%s ... %s", diagram[:4], diagram[len(diagram)-6:])

	// Output:
	// <svg ... </svg>
}
