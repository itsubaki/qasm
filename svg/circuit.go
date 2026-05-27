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
	Name    string `json:"name"`
	Targets []int  `json:"targets,omitempty"`
}

func (s *Subroutine) Wires() []int {
	return s.Targets
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

type Barrier struct{}

func (b *Barrier) Wires() []int {
	return nil
}
