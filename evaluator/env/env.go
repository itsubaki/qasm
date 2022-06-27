package env

import (
	"fmt"

	"github.com/itsubaki/qasm/ast"
	"github.com/itsubaki/qasm/evaluator/object"
)

type Func map[string]ast.Decl

type Const map[string]object.Object

type Environ struct {
	Const     Const
	Func      Func
	Bit       *Bit
	Qubit     *Qubit
	Modifier  []ast.Modifier
	Decl      []ast.Decl
	CtrlQArgs []ast.Expr
	Outer     *Environ
}

func New() *Environ {
	return &Environ{
		Func:      make(map[string]ast.Decl),
		Const:     make(map[string]object.Object),
		Bit:       NewBit(),
		Qubit:     NewQubit(),
		Modifier:  make([]ast.Modifier, 0),
		Decl:      make([]ast.Decl, 0),
		CtrlQArgs: nil,
		Outer:     nil,
	}
}

func (e *Environ) NewEnclosed(decl ast.Decl, mod []ast.Modifier) *Environ {
	return &Environ{
		Func:     e.Func,
		Const:    e.Const,
		Bit:      NewBit(),
		Qubit:    NewQubit(),
		Modifier: mod,
		Decl:     append(e.Decl, decl),
		Outer:    e,
	}
}

func (e *Environ) String() string {
	return fmt.Sprintf("func: %v, const: %v, bit: %v, qubit: %v, modifier: %v, decl: %v", e.Func, e.Const, e.Bit, e.Qubit, e.Modifier, e.Decl)
}
