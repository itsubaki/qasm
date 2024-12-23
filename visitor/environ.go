package visitor

import "github.com/itsubaki/q"

type Environ struct {
	Qubit map[string]q.Qubit
}

func NewEnviron() *Environ {
	return &Environ{
		Qubit: make(map[string]q.Qubit),
	}
}
