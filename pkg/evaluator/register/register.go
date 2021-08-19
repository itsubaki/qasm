package register

import (
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
)

type Register struct {
	Gate  Gate
	Const Const
	Bit   *Bit
	Qubit *Qubit
}

func New() *Register {
	return &Register{
		Gate:  make(map[string]ast.GateStmt),
		Const: make(map[string]int),
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
