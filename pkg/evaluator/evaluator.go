package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator/register"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Evaluator struct {
	Q *q.Q
	R *register.Register
}

func New(qsim *q.Q) *Evaluator {
	return &Evaluator{
		Q: qsim,
		R: register.New(),
	}
}

func Default() *Evaluator {
	return New(q.New())
}

func (e *Evaluator) Eval(p *ast.OpenQASM) error {
	for _, expr := range p.Gates {
		switch g := expr.(type) {
		case *ast.GateExpr:
			e.R.Gate[expr.Literal()] = *g
		default:
			return fmt.Errorf("invalid expr=%v", g)
		}
	}

	for _, stmt := range p.Statements {
		switch s := stmt.(type) {
		case *ast.ConstStmt:
			if err := e.evalConstStmt(s); err != nil {
				return fmt.Errorf("const: %v", err)
			}
		case *ast.DeclStmt:
			if err := e.evalDeclStmt(s); err != nil {
				return fmt.Errorf("decl: %v", err)
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

func (e *Evaluator) evalConstStmt(s *ast.ConstStmt) error {
	ident := s.Name.Value
	if _, ok := e.R.Const[ident]; ok {
		return fmt.Errorf("already exists=%v", ident)
	}

	e.R.Const[ident] = s.Int()
	return nil
}

func (e *Evaluator) evalDeclStmt(s *ast.DeclStmt) error {
	ident := s.Name.Value

	n := 1
	if s.Name.Index != nil {
		n = s.Name.Index.Int()
	}

	if s.Kind == lexer.BIT {
		if ok := e.R.Bit.Exists(ident); ok {
			return fmt.Errorf("already exists=%v", ident)
		}

		e.R.Bit.Add(ident, make([]int, n))
		return nil
	}

	if s.Kind == lexer.QUBIT {
		if ok := e.R.Qubit.Exists(ident); ok {
			return fmt.Errorf("already exists=%v", ident)
		}

		qb := e.Q.ZeroWith(n)
		e.R.Qubit.Add(ident, qb)
		return nil
	}

	return fmt.Errorf("invalid token=%v", s.Kind)
}

func (e *Evaluator) evalResetStmt(s *ast.ResetStmt) error {
	for _, t := range s.Target {
		qb, err := e.R.Qubit.Get(t.Value, t.Index)
		if err != nil {
			return fmt.Errorf("get qubit=%v: %v", t.Value, err)
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) evalApplyCModExp2(s *ast.ApplyStmt) error {
	if len(s.Target) != 4 {
		return fmt.Errorf("invalid target length %v", len(s.Target))
	}

	a, ok := e.R.Const[s.Target[0].Value]
	if !ok {
		return fmt.Errorf("IDENT=%v not found", s.Target[0].Value)
	}

	N, ok := e.R.Const[s.Target[1].Value]
	if !ok {
		return fmt.Errorf("IDENT=%v not found", s.Target[1].Value)
	}

	r0, err := e.R.Qubit.Get(s.Target[2].Value)
	if err != nil {
		return fmt.Errorf("get qubit=%v: %v", s.Target[2].Value, err)
	}

	r1, err := e.R.Qubit.Get(s.Target[3].Value)
	if err != nil {
		return fmt.Errorf("get qubit=%v: %v", s.Target[3].Value, err)
	}

	e.Q.CModExp2(a, N, r0, r1)
	return nil
}

func (e *Evaluator) evalApplyStmt(s *ast.ApplyStmt) error {
	if s.Kind == lexer.CMODEXP2 {
		if err := e.evalApplyCModExp2(s); err != nil {
			return fmt.Errorf("cmodexp2: %v", err)
		}

		return nil
	}

	in := make([]q.Qubit, 0)
	for _, t := range s.Target {
		qb, err := e.R.Qubit.Get(t.Value, t.Index)
		if err != nil {
			return fmt.Errorf("get qubit=%v: %v", t.Value, err)
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
		return fmt.Errorf("gate=%v(%v) not found", s.Kind, s.Literal())
	}

	return nil
}

func (e *Evaluator) evalAssignStmt(s *ast.AssignStmt) error {
	c, err := e.R.Bit.Get(s.Left.Value)
	if err != nil {
		return fmt.Errorf("get bit=%v: %v", s.Left.Value, err)
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
	return e.evalAssignStmt(&ast.AssignStmt{
		Kind:  lexer.EQUALS,
		Left:  s.Right,
		Right: s.Left,
	})
}

func (e *Evaluator) evalMeasureStmt(s *ast.MeasureStmt) ([]q.Qubit, error) {
	out := make([]q.Qubit, 0)
	for _, t := range s.Target {
		qb, err := e.R.Qubit.Get(t.Value, t.Index)
		if err != nil {
			return nil, fmt.Errorf("get qubit=%v: %v", t.Value, err)
		}

		e.Q.Measure(qb...)
		out = append(out, qb...)
	}

	return out, nil
}

func (e *Evaluator) evalPrintStmt(s *ast.PrintStmt) error {
	if s.Target == nil || len(s.Target) == 0 {
		return e.Println()
	}

	name := make([]string, 0)
	for _, t := range s.Target {
		name = append(name, t.Value)
	}

	return e.Println(name...)
}

func (e *Evaluator) Println(name ...string) error {
	if len(e.R.Qubit.Name) == 0 {
		return nil
	}

	if len(name) == 0 {
		name = e.R.Qubit.Name
	}

	index := make([][]int, 0)
	for _, n := range name {
		qb, err := e.R.Qubit.Get(n)
		if err != nil {
			return fmt.Errorf("get qubit=%v: %v", n, err)
		}

		index = append(index, q.Index(qb...))
	}

	for _, s := range e.Q.Raw().State(index...) {
		fmt.Println(s)
	}

	return nil
}
