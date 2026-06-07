package svg

var (
	_ Op = (*Gate)(nil)
	_ Op = (*Measurement)(nil)
)

type Circuit struct {
	Wires []Wire `json:"wires"`
	Ops   []Op   `json:"ops"`
}

type Wire struct {
	Name string `json:"name"`
}

type Op interface {
	Wires() []int
}

type Gate struct {
	Name    string `json:"name"`
	Control []int  `json:"control,omitempty"`
	Target  []int  `json:"target,omitempty"`
}

func (g *Gate) Wires() []int {
	var wires []int
	wires = append(wires, g.Control...)
	wires = append(wires, g.Target...)
	return wires
}

type Subroutine struct {
	Name string `json:"name"`
	Wire []int  `json:"wire,omitempty"`
}

func (s *Subroutine) Wires() []int {
	return s.Wire
}

type Measurement struct {
	Wire   []int `json:"wire"`
	Target []int `json:"target,omitempty"`
}

func (m *Measurement) Wires() []int {
	var wires []int
	wires = append(wires, m.Wire...)
	wires = append(wires, m.Target...)
	return wires
}

type Barrier struct {
	Wire []int `json:"wire"`
}

func (b *Barrier) Wires() []int {
	return b.Wire
}
