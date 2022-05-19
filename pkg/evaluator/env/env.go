package env

import (
	"fmt"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator/object"
)

type Func map[string]ast.Decl

type Const map[string]object.Object

type Environ struct {
	Const    Const
	Func     Func
	Bit      *Bit
	Qubit    *Qubit
	Outer    *Environ
	Decl     ast.Decl
	Modifier []ast.Modifier
}

func New() *Environ {
	return &Environ{
		Func:     make(map[string]ast.Decl),
		Const:    make(map[string]object.Object),
		Bit:      NewBit(),
		Qubit:    NewQubit(),
		Decl:     nil,
		Outer:    nil,
		Modifier: make([]ast.Modifier, 0),
	}
}

func (e *Environ) NewEnclosed(decl ast.Decl) *Environ {
	return &Environ{
		Func:     e.Func,
		Const:    e.Const,
		Bit:      NewBit(),
		Qubit:    NewQubit(),
		Decl:     decl,
		Outer:    e,
		Modifier: e.Modifier,
	}
}

func (e *Environ) String() string {
	return fmt.Sprintf("func: %v, const: %v, bit: %v, qubit: %v, modifier: %v", e.Func, e.Const, e.Bit, e.Qubit, e.Modifier)
}
