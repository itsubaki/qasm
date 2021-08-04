package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Evaluator struct {
	qubit map[string][]q.Qubit
	qsim  *q.Q
}

func New() *Evaluator {
	return &Evaluator{
		qubit: make(map[string][]q.Qubit),
		qsim:  q.New(),
	}
}

func (e *Evaluator) QSim() *q.Q {
	return e.qsim
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
		default:
			return fmt.Errorf("invalid stmt=%v", stmt)
		}
	}

	return nil
}

func (e *Evaluator) evalLetStmt(s *ast.LetStmt) error {
	if s.Kind == lexer.QUBIT {
		n := 1
		if s.Name.Index != nil {
			n = s.Name.IndexValue()
		}

		q := e.qsim.ZeroWith(n)
		e.qubit[s.Name.Value] = q
	}

	return nil
}

func (e *Evaluator) evalResetStmt(s *ast.ResetStmt) error {
	for _, n := range s.Name {
		q, ok := e.qubit[n.String()]
		if !ok {
			return fmt.Errorf("invalid ident=%v", n.String())
		}

		e.qsim.Reset(q...)
	}

	return nil
}

func (e *Evaluator) evalApplyStmt(s *ast.ApplyStmt) error {
	qb, ok := e.qubit[s.Name.Value]
	if !ok {
		return fmt.Errorf("invalid ident=%v", s.Name.String())
	}

	if s.Name.Index != nil {
		index := s.Name.IndexValue()
		qb = append(make([]q.Qubit, 0), qb[index])
	}

	switch s.Kind {
	case lexer.X:
		e.qsim.X(qb...)
	case lexer.Y:
		e.qsim.Y(qb...)
	case lexer.Z:
		e.qsim.Z(qb...)
	case lexer.H:
		e.qsim.H(qb...)
	default:
		return fmt.Errorf("invalid token=%v", s.Kind)
	}

	return nil
}

func (e *Evaluator) evalMeasureStmt(s *ast.MeasureStmt) error {
	q, ok := e.qubit[s.Name.String()]
	if !ok {
		return fmt.Errorf("invalid ident=%v", s.Name.String())
	}

	e.qsim.Measure(q...)
	return nil
}
