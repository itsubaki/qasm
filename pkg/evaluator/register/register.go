package register

import (
	"math"

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
	c := make(map[string]float64)
	c["pi"] = math.Pi
	c["tau"] = math.Pi * 2
	c["euler"] = math.E

	return &Register{
		Const: c,
		Bit: &Bit{
			Name:  make([]string, 0),
			Value: make(map[string][]int),
		},
		Qubit: &Qubit{
			Name:  make([]string, 0),
			Value: make(map[string][]q.Qubit),
		},
		Func: make(map[string]ast.Decl),
	}
}
