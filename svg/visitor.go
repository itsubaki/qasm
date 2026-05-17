package svg

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/environ"
)

type Visitor struct {
	circuit *Circuit
	env     *environ.Environ
}

func New(circuit *Circuit, env *environ.Environ) *Visitor {
	return &Visitor{
		circuit: circuit,
		env:     env,
	}
}

func (v *Visitor) Run(tree antlr.ParseTree) error {
	return nil
}
