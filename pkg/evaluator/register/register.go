package register

import (
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
		Const: NewConst(),
		Bit:   NewBit(),
		Qubit: NewQubit(),
		Func:  make(map[string]ast.Decl),
	}
}
