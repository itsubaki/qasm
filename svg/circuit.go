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
	Name     string `json:"name"`
	Controls []int  `json:"controls,omitempty"`
	Targets  []int  `json:"targets,omitempty"`
}

func (g *Gate) Wires() []int {
	var wires []int
	wires = append(wires, g.Controls...)
	wires = append(wires, g.Targets...)
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
	Wire    []int `json:"wire"`
	Targets []int `json:"targets,omitempty"`
}

func (m *Measurement) Wires() []int {
	var wires []int
	wires = append(wires, m.Wire...)
	wires = append(wires, m.Targets...)
	return wires
}

type Barrier struct {
	Wire []int `json:"wires"`
}

func (b *Barrier) Wires() []int {
	return b.Wire
}
