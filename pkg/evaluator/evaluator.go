package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Evaluator struct {
	Bit   map[string][]int
	Qubit map[string][]q.Qubit
	Q     *q.Q
}

func New(qsim *q.Q) *Evaluator {
	return &Evaluator{
		Bit:   make(map[string][]int),
		Qubit: make(map[string][]q.Qubit),
		Q:     qsim,
	}
}

func Default() *Evaluator {
	return New(q.New())
}

func (e *Evaluator) Eval(p *ast.Program) error {
	for _, stmt := range p.Statements {
		switch s := stmt.(type) {
		case *ast.LetStmt:
			if err := e.evalLetStmt(s); err != nil {
				return fmt.Errorf("eval let: %v", err)
			}
		case *ast.ResetStmt:
			if err := e.evalResetStmt(s); err != nil {
				return fmt.Errorf("eval reset: %v", err)
			}
		case *ast.ApplyStmt:
			if err := e.evalApplyStmt(s); err != nil {
				return fmt.Errorf("eval apply: %v", err)
			}
		case *ast.MeasureStmt:
			if err := e.evalMeasureStmt(s); err != nil {
				return fmt.Errorf("eval measure: %v", err)
			}
		case *ast.AssignStmt:
			if err := e.evalAssignStmt(s); err != nil {
				return fmt.Errorf("eval assign: %v", err)
			}
		default:
			return fmt.Errorf("invalid stmt=%v", stmt)
		}
	}

	return nil
}

func (e *Evaluator) evalLetStmt(s *ast.LetStmt) error {
	num := func(i *ast.IdentExpr) int {
		if i.Index == nil {
			return 1
		}

		return i.IndexValue()
	}

	if s.Kind == lexer.QUBIT {
		n := num(s.Name)
		qb := e.Q.ZeroWith(n)
		e.Qubit[s.Name.Value] = qb

		return nil
	}

	if s.Kind == lexer.BIT {
		n := num(s.Name)
		e.Bit[s.Name.Value] = make([]int, n)

		return nil
	}

	return fmt.Errorf("invalid token=%v", s.Kind)
}

func (e *Evaluator) evalResetStmt(s *ast.ResetStmt) error {
	for _, n := range s.Target {
		qb, ok := e.Qubit[n.Value]
		if !ok {
			return fmt.Errorf("invalid ident=%v", n.Value)
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) evalApplyStmt(s *ast.ApplyStmt) error {
	qb, ok := e.Qubit[s.Target.Value]
	if !ok {
		return fmt.Errorf("invalid ident=%v", s.Target.Value)
	}

	if s.Target.Index != nil {
		index := s.Target.IndexValue()
		qb = append(make([]q.Qubit, 0), qb[index])
	}

	switch s.Kind {
	case lexer.X:
		e.Q.X(qb...)
	case lexer.Y:
		e.Q.Y(qb...)
	case lexer.Z:
		e.Q.Z(qb...)
	case lexer.H:
		e.Q.H(qb...)
	default:
		return fmt.Errorf("invalid token=%v", s.Kind)
	}

	return nil
}
func (e *Evaluator) evalAssignStmt(s *ast.AssignStmt) error {
	c := e.Bit[s.Left.Value]

	switch s := s.Right.(type) {
	case *ast.MeasureStmt:
		if err := e.evalMeasureStmt(s); err != nil {
			return fmt.Errorf("eval measure: %v", err)
		}

		qb, ok := e.Qubit[s.Target.Value]
		if !ok {
			return fmt.Errorf("invalid ident=%v", s.Target.Value)
		}

		for _, q := range qb {
			c[q] = e.Q.State(q)[0].Int[0]
		}
	default:
		return fmt.Errorf("invalid stmt=%v", s)
	}

	return nil
}

func (e *Evaluator) evalMeasureStmt(s *ast.MeasureStmt) error {
	qb, ok := e.Qubit[s.Target.Value]
	if !ok {
		return fmt.Errorf("invalid ident=%v", s.Target.Value)
	}

	e.Q.Measure(qb...)
	return nil
}
