package svg

import "fmt"

var (
	_ Op = (*Gate)(nil)
	_ Op = (*Subroutine)(nil)
	_ Op = (*Measurement)(nil)
	_ Op = (*Barrier)(nil)
)

type Op interface {
	WireIDs() []WireID
}

type Circuit struct {
	Wires []Wire `json:"wires"`
	Ops   []Op   `json:"ops"`
}

// WireID is the position of a wire in Circuit.Wires.
type WireID int

// Wire identifies a source variable or an element of a source register.
type Wire struct {
	Name  string `json:"name"`
	Index *int   `json:"index,omitempty"`
}

func (w Wire) Label() string {
	if w.Index != nil {
		return fmt.Sprintf("%s[%d]", w.Name, *w.Index)
	}

	return w.Name
}

type Gate struct {
	Name    string   `json:"name"`
	Control []WireID `json:"control,omitempty"`
	Target  []WireID `json:"target,omitempty"`
}

func (g *Gate) WireIDs() []WireID {
	var wires []WireID
	wires = append(wires, g.Control...)
	wires = append(wires, g.Target...)
	return wires
}

type Subroutine struct {
	Name string   `json:"name"`
	Wire []WireID `json:"wire,omitempty"`
}

func (s *Subroutine) WireIDs() []WireID {
	return s.Wire
}

type Measurement struct {
	Wire   []WireID `json:"wire"`
	Target []WireID `json:"target,omitempty"`
}

func (m *Measurement) WireIDs() []WireID {
	var wires []WireID
	wires = append(wires, m.Wire...)
	wires = append(wires, m.Target...)
	return wires
}

type Barrier struct {
	Wire []WireID `json:"wire"`
}

func (b *Barrier) WireIDs() []WireID {
	return b.Wire
}
