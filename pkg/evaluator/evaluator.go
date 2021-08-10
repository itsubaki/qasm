package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Evaluator struct {
	Const map[string]int
	Bit   *Bit
	Qubit *Qubit
	Q     *q.Q
}

func New(qsim *q.Q) *Evaluator {
	return &Evaluator{
		Const: make(map[string]int),
		Bit: &Bit{
			Name:  make([]string, 0),
			Value: make(map[string][]int),
		},
		Qubit: &Qubit{
			Name:  make([]string, 0),
			Value: make(map[string][]q.Qubit),
		},
		Q: qsim,
	}
}

func Default() *Evaluator {
	return New(q.New())
}

func (e *Evaluator) Eval(p *ast.OpenQASM) error {
	for _, stmt := range p.Statements {
		switch s := stmt.(type) {
		case *ast.DeclConstStmt:
			if err := e.evalDeclConstStmt(s); err != nil {
				return fmt.Errorf("let: %v", err)
			}
		case *ast.DeclStmt:
			if err := e.evalDeclStmt(s); err != nil {
				return fmt.Errorf("let: %v", err)
			}
		case *ast.ResetStmt:
			if err := e.evalResetStmt(s); err != nil {
				return fmt.Errorf("reset: %v", err)
			}
		case *ast.ApplyStmt:
			if err := e.evalApplyStmt(s); err != nil {
				return fmt.Errorf("apply: %v", err)
			}
		case *ast.MeasureStmt:
			if _, err := e.evalMeasureStmt(s); err != nil {
				return fmt.Errorf("measure: %v", err)
			}
		case *ast.ArrowStmt:
			if err := e.evalArrowStmt(s); err != nil {
				return fmt.Errorf("arrow: %v", err)
			}
		case *ast.AssignStmt:
			if err := e.evalAssignStmt(s); err != nil {
				return fmt.Errorf("assign: %v", err)
			}
		case *ast.PrintStmt:
			if err := e.evalPrintStmt(s); err != nil {
				return fmt.Errorf("print: %v", err)
			}
		default:
			return fmt.Errorf("invalid stmt=%v", stmt)
		}
	}

	return nil
}

func (e *Evaluator) evalDeclConstStmt(s *ast.DeclConstStmt) error {
	ident := s.Name.Value
	if _, ok := e.Const[ident]; ok {
		return fmt.Errorf("already exists=%v", ident)
	}

	e.Const[ident] = s.Int()
	return nil
}

func (e *Evaluator) evalDeclStmt(s *ast.DeclStmt) error {
	n := 1
	if s.Index != nil {
		n = s.Index.Int()
	}
	ident := s.Name.Value

	if s.Kind == lexer.QUBIT {
		if ok := e.Qubit.Exists(ident); ok {
			return fmt.Errorf("already exists=%v", ident)
		}

		qb := e.Q.ZeroWith(n)
		e.Qubit.Add(ident, qb)
		return nil
	}

	if s.Kind == lexer.BIT {
		if ok := e.Bit.Exists(ident); ok {
			return fmt.Errorf("already exists=%v", ident)
		}

		e.Bit.Add(ident, make([]int, n))
		return nil
	}

	return fmt.Errorf("invalid token=%v", s.Kind)
}

func (e *Evaluator) evalResetStmt(s *ast.ResetStmt) error {
	for _, t := range s.Target {
		qb, err := e.Qubit.Get(t.Value, t.Index)
		if err != nil {
			return fmt.Errorf("get qubit: %v", err)
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) evalApplyCModExp2(s *ast.ApplyStmt) error {
	if len(s.Target) != 4 {
		return fmt.Errorf("invalid target length %v", len(s.Target))
	}

	a, ok := e.Const[s.Target[0].Value]
	if !ok {
		return fmt.Errorf("IDENT=%v not found", s.Target[0].Value)
	}

	N, ok := e.Const[s.Target[1].Value]
	if !ok {
		return fmt.Errorf("IDENT=%v not found", s.Target[1].Value)
	}

	r0, err := e.Qubit.Get(s.Target[2].Value)
	if err != nil {
		return fmt.Errorf("get qubit: %v", err)
	}

	r1, err := e.Qubit.Get(s.Target[3].Value)
	if err != nil {
		return fmt.Errorf("get qubit: %v", err)
	}

	e.Q.CModExp2(a, N, r0, r1)
	return nil
}

func (e *Evaluator) evalApplyStmt(s *ast.ApplyStmt) error {
	if s.Kind == lexer.CMODEXP2 {
		if err := e.evalApplyCModExp2(s); err != nil {
			return fmt.Errorf("apply: %v", err)
		}

		return nil
	}

	in := make([]q.Qubit, 0)
	for _, t := range s.Target {
		qb, err := e.Qubit.Get(t.Value, t.Index)
		if err != nil {
			return fmt.Errorf("get qubit: %v", err)
		}

		in = append(in, qb...)
	}

	switch s.Kind {
	case lexer.X:
		e.Q.X(in...)
	case lexer.Y:
		e.Q.Y(in...)
	case lexer.Z:
		e.Q.Z(in...)
	case lexer.H:
		e.Q.H(in...)
	case lexer.T:
		e.Q.T(in...)
	case lexer.S:
		e.Q.S(in...)
	case lexer.CX:
		e.Q.CNOT(in[0], in[1])
	case lexer.CZ:
		e.Q.CZ(in[0], in[1])
	case lexer.CCX:
		e.Q.CCNOT(in[0], in[1], in[2])
	case lexer.SWAP:
		e.Q.Swap(in...)
	case lexer.QFT:
		e.Q.QFT(in...)
	case lexer.IQFT:
		e.Q.InvQFT(in...)
	default:
		return fmt.Errorf("gate=%v not found", s.Kind)
	}

	return nil
}

func (e *Evaluator) evalAssignStmt(s *ast.AssignStmt) error {
	c, err := e.Bit.Get(s.Left.Value)
	if err != nil {
		return fmt.Errorf("get bit: %v", err)
	}

	switch s := s.Right.(type) {
	case *ast.MeasureStmt:
		qb, err := e.evalMeasureStmt(s)
		if err != nil {
			return fmt.Errorf("measure: %v", err)
		}
		for _, q := range qb {
			c[q] = e.Q.State(q)[0].Int[0]
		}
	default:
		return fmt.Errorf("invalid stmt=%v", s)
	}

	return nil
}

func (e *Evaluator) evalArrowStmt(s *ast.ArrowStmt) error {
	return nil
}

func (e *Evaluator) evalMeasureStmt(s *ast.MeasureStmt) ([]q.Qubit, error) {
	out := make([]q.Qubit, 0)
	for _, t := range s.Target {
		qb, err := e.Qubit.Get(t.Value, t.Index)
		if err != nil {
			return nil, fmt.Errorf("get qubit: %v", err)
		}

		e.Q.Measure(qb...)
		out = append(out, qb...)
	}

	return out, nil
}

func (e *Evaluator) evalPrintStmt(s *ast.PrintStmt) error {
	return e.Println()
}

func (e *Evaluator) Println() error {
	if len(e.Qubit.Name) == 0 {
		return nil
	}

	index := make([][]int, 0)
	for _, ident := range e.Qubit.Name {
		qb, err := e.Qubit.Get(ident)
		if err != nil {
			return fmt.Errorf("get qubit: %v", err)
		}

		index = append(index, q.Index(qb...))
	}

	for _, s := range e.Q.Raw().State(index...) {
		fmt.Println(s)
	}

	return nil
}
