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
				Name:    "H",
				Targets: []int{1},
			},
			&svg.Gate{
				Name:    "X",
				Targets: []int{0},
			},
			&svg.Gate{
				Name:     "X",
				Controls: []int{1},
				Targets:  []int{2},
			},
			&svg.Gate{
				Name:     "X",
				Controls: []int{0},
				Targets:  []int{1},
			},
			&svg.Gate{
				Name:    "H",
				Targets: []int{0},
			},
			&svg.Gate{
				Name:     "U",
				Controls: []int{2, 0},
				Targets:  []int{1},
			},
			&svg.Gate{
				Name:     "QFT",
				Controls: []int{0},
				Targets:  []int{1, 2, 3},
			},
			&svg.Measurement{
				Wire: []int{
					0,
				},
				Target: []int{
					4,
				},
			},
			&svg.Measurement{
				Wire: []int{
					1,
				},
				Target: []int{
					5,
				},
			},
			&svg.Measurement{
				Wire: []int{
					0,
				},
			},
			&svg.Measurement{
				Wire: []int{
					1, 2,
				},
			},
			&svg.Measurement{
				Wire: []int{
					0,
				},
			},
			&svg.Gate{
				Name:    "H",
				Targets: []int{1},
			},
			&svg.Gate{
				Name:    "H",
				Targets: []int{0},
			},
			&svg.Gate{
				Name:     "QFT",
				Controls: []int{3},
				Targets:  []int{0, 1, 2},
			},
			&svg.Gate{
				Name:    "H",
				Targets: []int{0},
			},
			&svg.Measurement{
				Wire: []int{
					1,
				},
			},
			&svg.Measurement{
				Wire: []int{
					0,
				},
			},
			&svg.Measurement{
				Wire: []int{
					0,
				},
				Target: []int{
					5,
				},
			},
			&svg.Gate{
				Name:    "X",
				Targets: []int{2},
			},
			&svg.Gate{
				Name:    "H",
				Targets: []int{3},
			},
			&svg.Barrier{},
			&svg.Gate{
				Name:    "X",
				Targets: []int{0},
			},
			&svg.Gate{
				Name:    "H",
				Targets: []int{1},
			},
		},
	}

	layout := svg.NewLayout(circuit)
	fmt.Println(layout.Wires)

	// Output:
	// [{q0} {q1} {q2} {q3} {c0} {c1}]
}

func ExampleLayer_CompatibleWith() {
	fmt.Println((&svg.Layer{}).CompatibleWith(&svg.Barrier{}))

	// Output:
	// false
}
