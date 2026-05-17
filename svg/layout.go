package svg

type Layout struct {
	Wires  []Wire  `json:"wires"`
	Layers []Layer `json:"layers"`
}

func (l *Layout) NewLayer(ops []Op, separated ...bool) {
	wires := make(map[int]bool)
	for _, op := range ops {
		for _, w := range op.Wires() {
			wires[w] = true
		}
	}

	sep := true
	if len(separated) > 0 {
		sep = separated[0]
	}

	l.Layers = append(l.Layers, Layer{
		Wires:     wires,
		Ops:       ops,
		separated: sep,
	})
}

type Layer struct {
	Wires     map[int]bool `json:"wires"`
	Ops       []Op         `json:"ops"`
	separated bool
}

func (l *Layer) Add(op Op) {
	for _, w := range op.Wires() {
		l.Wires[w] = true
	}

	l.Ops = append(l.Ops, op)
}

func (l *Layer) Conflicts(wire []int) bool {
	for _, w := range wire {
		if l.Wires[w] {
			return true
		}
	}

	return false
}

func (l *Layer) CompatibleWith(cur Op) bool {
	if l.separated {
		return false
	}

	if _, ok := cur.(*Barrier); ok {
		return false
	}

	// measurements must be in a separate layer from non-measurements
	if _, ok := cur.(*Measurement); ok {
		for _, m := range l.Ops {
			if _, ok := m.(*Measurement); !ok {
				return false
			}
		}
	}

	if _, ok := cur.(*Measurement); !ok {
		for _, m := range l.Ops {
			if _, ok := m.(*Measurement); ok {
				return false
			}
		}
	}

	return true
}

func NewLayout(circuit *Circuit) *Layout {
	layout := &Layout{
		Wires:  circuit.Wires,
		Layers: []Layer{},
	}

	for _, cur := range circuit.Ops {
		// controlled gates must be in their own layer
		if g, ok := cur.(*Gate); ok && len(g.Controls) > 0 {
			layout.NewLayer([]Op{cur})
			continue
		}

		// arrow measurements must be in their own layer
		if m, ok := cur.(*Measurement); ok && len(m.Target) > 0 {
			layout.NewLayer([]Op{cur})
			continue
		}

		// barriers must be in their own layer
		if _, ok := cur.(*Barrier); ok {
			layout.NewLayer([]Op{cur})
			continue
		}

		var placed bool
		for i := len(layout.Layers) - 1; i >= 0; i-- {
			if !layout.Layers[i].CompatibleWith(cur) {
				break
			}

			if layout.Layers[i].Conflicts(cur.Wires()) {
				continue
			}

			layout.Layers[i].Add(cur)
			placed = true
			break
		}

		if placed {
			continue
		}

		layout.NewLayer([]Op{cur}, false)
	}

	return layout
}
