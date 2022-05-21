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
	Decl     []ast.Decl
	Modifier []ast.Modifier
	Outer    *Environ
}

func New() *Environ {
	return &Environ{
		Func:     make(map[string]ast.Decl),
		Const:    make(map[string]object.Object),
		Bit:      NewBit(),
		Qubit:    NewQubit(),
		Decl:     make([]ast.Decl, 0),
		Modifier: make([]ast.Modifier, 0),
		Outer:    nil,
	}
}

func (e *Environ) NewEnclosed(decl ast.Decl, mod []ast.Modifier) *Environ {
	return &Environ{
		Func:     e.Func,
		Const:    e.Const,
		Bit:      NewBit(),
		Qubit:    NewQubit(),
		Decl:     append(e.Decl, decl),
		Modifier: append(e.Modifier, mod...),
		Outer:    e,
	}
}

func (e *Environ) String() string {
	return fmt.Sprintf("func: %v, const: %v, bit: %v, qubit: %v, modifier: %v, decl: %v", e.Func, e.Const, e.Bit, e.Qubit, e.Modifier, e.Decl)
}
