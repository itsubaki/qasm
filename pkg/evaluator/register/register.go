package register

import (
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
)

type Register struct {
	Const Const
	Bit   *Bit
	Qubit *Qubit
	Func  Func
}

func New() *Register {
	return &Register{
		Func:  make(map[string]ast.Decl),
		Const: make(map[string]float64),
		Bit: &Bit{
			Name:  make([]string, 0),
			Value: make(map[string][]int),
		},
		Qubit: &Qubit{
			Name:  make([]string, 0),
			Value: make(map[string][]q.Qubit),
		},
	}
}