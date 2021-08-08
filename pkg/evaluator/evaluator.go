package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Evaluator struct {
	Const map[string]int
	Bit   map[string][]int
	Qubit map[string][]q.Qubit
	Order []string
	Q     *q.Q
}

func New(qsim *q.Q) *Evaluator {
	return &Evaluator{
		Const: make(map[string]int),
		Bit:   make(map[string][]int),
		Qubit: make(map[string][]q.Qubit),
		Order: make([]string, 0),
		Q:     qsim,
	}
}

func Default() *Evaluator {
	return New(q.New())
}

func (e *Evaluator) Clear() {
	e.Bit = make(map[string][]int)
	e.Qubit = make(map[string][]q.Qubit)
	e.Order = make([]string, 0)
	e.Q = q.New()
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
			if s.Kind == lexer.CMODEXP2 {
				if err := e.evalApplyCModExp2(s); err != nil {
					return fmt.Errorf("apply: %v", err)
				}

				continue
			}

			if err := e.evalApplyStmt(s); err != nil {
				return fmt.Errorf("apply: %v", err)
			}
		case *ast.MeasureStmt:
			if _, err := e.evalMeasureStmt(s); err != nil {
				return fmt.Errorf("measure: %v", err)
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
		if _, ok := e.Qubit[ident]; ok {
			return fmt.Errorf("already exists=%v", ident)
		}

		qb := e.Q.ZeroWith(n)
		e.Qubit[ident] = qb
		e.Order = append(e.Order, ident)
		return nil
	}

	if s.Kind == lexer.BIT {
		if _, ok := e.Bit[ident]; ok {
			return fmt.Errorf("already exists=%v", ident)
		}

		e.Bit[ident] = make([]int, n)
		return nil
	}

	return fmt.Errorf("invalid token=%v", s.Kind)
}

func (e *Evaluator) evalResetStmt(s *ast.ResetStmt) error {
	for _, t := range s.Target {
		ident := t.Value
		qb, ok := e.Qubit[ident]
		if !ok {
			return fmt.Errorf("IDENT=%v not found", ident)
		}

		if t.Index != nil {
			e.Q.Reset(qb[t.Index.Int()])
			continue
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

	r0, ok := e.Qubit[s.Target[2].Value]
	if !ok {
		return fmt.Errorf("IDENT=%v not found", s.Target[2].Value)
	}

	r1, ok := e.Qubit[s.Target[3].Value]
	if !ok {
		return fmt.Errorf("IDENT=%v not found", s.Target[3].Value)
	}

	e.Q.CModExp2(a, N, r0, r1)
	return nil
}

func (e *Evaluator) evalApplyStmt(s *ast.ApplyStmt) error {
	in := make([]q.Qubit, 0)
	for _, t := range s.Target {
		ident := t.Value
		qb, ok := e.Qubit[ident]
		if !ok {
			return fmt.Errorf("IDENT=%v not found", ident)
		}

		if t.Index != nil {
			index := t.Index.Int()
			if index > len(qb)-1 {
				return fmt.Errorf("index out of range[%v] with length %v", index, len(qb))
			}

			qb = append(make([]q.Qubit, 0), qb[index])
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
	c := e.Bit[s.Left.Value]
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

func (e *Evaluator) evalMeasureStmt(s *ast.MeasureStmt) ([]q.Qubit, error) {
	for _, t := range s.Target {
		ident := t.Value
		qb, ok := e.Qubit[ident]
		if !ok {
			return nil, fmt.Errorf("IDENT=%v not found", ident)
		}

		if t.Index != nil {
			index := t.Index.Int()
			qb = append(make([]q.Qubit, 0), qb[index])
		}

		e.Q.Measure(qb...)
	}

	return e.Qubit[s.Target[0].Value], nil
}

func (e *Evaluator) evalPrintStmt(s *ast.PrintStmt) error {
	return e.Println()
}

func (e *Evaluator) Println() error {
	if len(e.Qubit) == 0 {
		return nil
	}

	index := make([][]int, 0)
	for _, ident := range e.Order {
		qb, ok := e.Qubit[ident]
		if !ok {
			return fmt.Errorf("IDENT=%v not found", ident)
		}

		index = append(index, q.Index(qb...))
	}

	for _, s := range e.Q.Raw().State(index...) {
		fmt.Println(s)
	}

	return nil
}
