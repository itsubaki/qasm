package svg_test

import (
	"fmt"

	"github.com/itsubaki/qasm/svg"
)

func ExampleNewLayout() {
	circuit := &svg.Circuit{
		Wires: []svg.Wire{
			{Name: "q0"},
			{Name: "q1"},
			{Name: "q2"},
			{Name: "q3"},
			{Name: "c0"},
			{Name: "c1"},
		},
		Ops: []svg.Op{
			&svg.Gate{
				Name:   "H",
				Target: []svg.WireID{1},
			},
			&svg.Gate{
				Name:   "X",
				Target: []svg.WireID{0},
			},
			&svg.Gate{
				Name:    "X",
				Control: []svg.WireID{1},
				Target:  []svg.WireID{2},
			},
			&svg.Gate{
				Name:    "X",
				Control: []svg.WireID{0},
				Target:  []svg.WireID{1},
			},
			&svg.Gate{
				Name:   "H",
				Target: []svg.WireID{0},
			},
			&svg.Gate{
				Name:    "U",
				Control: []svg.WireID{2, 0},
				Target:  []svg.WireID{1},
			},
			&svg.Gate{
				Name:    "QFT",
				Control: []svg.WireID{0},
				Target:  []svg.WireID{1, 2, 3},
			},
			&svg.Measurement{
				Wire: []svg.WireID{
					0,
				},
				Target: []svg.WireID{
					4,
				},
			},
			&svg.Measurement{
				Wire: []svg.WireID{
					1,
				},
				Target: []svg.WireID{
					5,
				},
			},
			&svg.Measurement{
				Wire: []svg.WireID{
					0,
				},
			},
			&svg.Measurement{
				Wire: []svg.WireID{
					1, 2,
				},
			},
			&svg.Measurement{
				Wire: []svg.WireID{
					0,
				},
			},
			&svg.Gate{
				Name:   "H",
				Target: []svg.WireID{1},
			},
			&svg.Gate{
				Name:   "H",
				Target: []svg.WireID{0},
			},
			&svg.Gate{
				Name:    "QFT",
				Control: []svg.WireID{3},
				Target:  []svg.WireID{0, 1, 2},
			},
			&svg.Gate{
				Name:   "H",
				Target: []svg.WireID{0},
			},
			&svg.Measurement{
				Wire: []svg.WireID{
					1,
				},
			},
			&svg.Measurement{
				Wire: []svg.WireID{
					0,
				},
			},
			&svg.Measurement{
				Wire: []svg.WireID{
					0,
				},
				Target: []svg.WireID{
					5,
				},
			},
			&svg.Gate{
				Name:   "X",
				Target: []svg.WireID{2},
			},
			&svg.Gate{
				Name:   "H",
				Target: []svg.WireID{3},
			},
			&svg.Barrier{
				Wire: []svg.WireID{0, 1, 2, 3},
			},
			&svg.Gate{
				Name:   "X",
				Target: []svg.WireID{0},
			},
			&svg.Gate{
				Name:   "H",
				Target: []svg.WireID{1},
			},
		},
	}

	layout := svg.NewLayout(circuit)
	for _, wire := range layout.Wires {
		fmt.Println(wire.Label())
	}

	// Output:
	// q0
	// q1
	// q2
	// q3
	// c0
	// c1
}
