package svg

import "github.com/antlr4-go/antlr/v4"

type Visitor struct {
	wire map[string]int
}

func NewVisitor() *Visitor {
	return &Visitor{
		wire: make(map[string]int),
	}
}

func (v *Visitor) Run(tree antlr.ParseTree) (*Circuit, error) {
	// bytea, err := json.MarshalIndent(TODO, "", "  ")
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(string(bytea))

	return TODO, nil
}

var TODO = &Circuit{
	Wires: []Wire{
		{Name: "q0"},
		{Name: "q1"},
		{Name: "q2"},
		{Name: "q3"},
		{Name: "c0"},
		{Name: "c1"},
	},
	Ops: []Op{
		&Gate{
			Name:    "H",
			Targets: []int{1},
		},
		&Gate{
			Name:    "X",
			Targets: []int{0},
		},
		&Gate{
			Name:     "X",
			Controls: []int{1},
			Targets:  []int{2},
		},
		&Gate{
			Name:     "X",
			Controls: []int{0},
			Targets:  []int{1},
		},
		&Gate{
			Name:    "H",
			Targets: []int{0},
		},
		&Gate{
			Name:     "U",
			Controls: []int{2, 0},
			Targets:  []int{1},
		},
		&Gate{
			Name:     "QFT",
			Controls: []int{0},
			Targets:  []int{1, 2, 3},
		},
		&Measurement{
			Wire: []int{
				0,
			},
			Target: []int{
				4,
			},
		},
		&Measurement{
			Wire: []int{
				1,
			},
			Target: []int{
				5,
			},
		},
		&Measurement{
			Wire: []int{
				0,
			},
		},
		&Measurement{
			Wire: []int{
				1, 2,
			},
		},
		&Measurement{
			Wire: []int{
				0,
			},
		},
		&Gate{
			Name:    "H",
			Targets: []int{1},
		},
		&Gate{
			Name:    "H",
			Targets: []int{0},
		},
		&Gate{
			Name:     "QFT",
			Controls: []int{3},
			Targets:  []int{0, 1, 2},
		},
		&Gate{
			Name:    "H",
			Targets: []int{0},
		},
		&Measurement{
			Wire: []int{
				1,
			},
		},
		&Measurement{
			Wire: []int{
				0,
			},
		},
		&Measurement{
			Wire: []int{
				0,
			},
			Target: []int{
				5,
			},
		},
		&Gate{
			Name:    "X",
			Targets: []int{2},
		},
		&Gate{
			Name:    "H",
			Targets: []int{3},
		},
		&Barrier{},
		&Gate{
			Name:    "X",
			Targets: []int{0},
		},
		&Gate{
			Name:    "H",
			Targets: []int{1},
		},
	},
}
