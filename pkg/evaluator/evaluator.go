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

func (e *Evaluator) Eval(p *ast.OpenQASM) error {
	for _, stmt := range p.Statements {
		switch s := stmt.(type) {
		case *ast.LetStmt:
			if err := e.evalLetStmt(s); err != nil {
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

func (e *Evaluator) evalLetStmt(s *ast.LetStmt) error {
	n := 1
	if s.Index != nil {
		n = s.Index.Int()
	}

	if s.Kind == lexer.QUBIT {
		if _, ok := e.Qubit[s.Name.Value]; ok {
			return fmt.Errorf("already exists=%v", s.Name.Value)
		}

		qb := e.Q.ZeroWith(n)
		e.Qubit[s.Name.Value] = qb
		return nil
	}

	if s.Kind == lexer.BIT {
		if _, ok := e.Bit[s.Name.Value]; ok {
			return fmt.Errorf("already exists=%v", s.Name.Value)
		}

		e.Bit[s.Name.Value] = make([]int, n)
		return nil
	}

	return fmt.Errorf("invalid token=%v", s.Kind)
}

func (e *Evaluator) evalResetStmt(s *ast.ResetStmt) error {
	for _, t := range s.Target {
		qb, ok := e.Qubit[t.Value]
		if !ok {
			return fmt.Errorf("IDENT=%v not found", t.Value)
		}

		if t.Index != nil {
			e.Q.Reset(qb[t.Index.Int()])
			continue
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) evalApplyStmt(s *ast.ApplyStmt) error {
	for _, t := range s.Target {
		qb, ok := e.Qubit[t.Value]
		if !ok {
			return fmt.Errorf("IDENT=%v not found", t.Value)
		}

		if t.Index != nil {
			index := t.Index.Int()
			if index > len(qb)-1 {
				return fmt.Errorf("index out of range[%v] with length %v", index, len(qb))
			}

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
			return fmt.Errorf("gate=%v not found", s.Kind)
		}
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
		qb, ok := e.Qubit[t.Value]
		if !ok {
			return nil, fmt.Errorf("IDENT=%v not found", t.Value)
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
	if len(e.Qubit) == 0 {
		return nil
	}

	for _, s := range e.Q.State() {
		fmt.Println(s)
	}

	return nil
}
