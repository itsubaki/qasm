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
	for _, s := range p.Stmts {
		if err := e.eval(s); err != nil {
			return fmt.Errorf("eval: %v", err)
		}
	}

	return nil
}

func (e *Evaluator) eval(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.DeclStmt:
		if err := e.evalDeclStmt(s); err != nil {
			return fmt.Errorf("decl: %v", err)
		}
	case *ast.ExprStmt:
		if err := e.evalExprStmt(s); err != nil {
			return fmt.Errorf("expr: %v", err)
		}
	case *ast.AssignStmt:
		if err := e.evalAssignStmt(s); err != nil {
			return fmt.Errorf("assign: %v", err)
		}
	case *ast.ArrowStmt:
		if err := e.evalArrowStmt(s); err != nil {
			return fmt.Errorf("arrow: %v", err)
		}
	default:
		return fmt.Errorf("invalid stmt=%#v", stmt)
	}

	return nil
}

func (e *Evaluator) evalDeclStmt(s *ast.DeclStmt) error {
	switch d := s.Decl.(type) {
	case *ast.GenConst:
		if _, ok := e.R.Const[d.Name.Value]; ok {
			return fmt.Errorf("already exists=%v", d.Name.Value)
		}

		e.R.Const[d.Name.Value] = d.Int()
	case *ast.GenDecl:
		switch d.Kind {
		case lexer.BIT:
			if _, ok := e.R.Bit.Get(&ast.IdentExpr{Value: d.Name.Value}); ok {
				return fmt.Errorf("already exists=%v", d.Name)
			}
			e.R.Bit.Add(d.Name.Value, make([]int, d.Size()))

		case lexer.QUBIT:
			if _, ok := e.R.Qubit.Get(&ast.IdentExpr{Value: d.Name.Value}); ok {
				return fmt.Errorf("already exists=%v", d.Name)
			}
			qb := e.Q.ZeroWith(d.Size())
			e.R.Qubit.Add(d.Name.Value, qb)

		default:
			return fmt.Errorf("invalid kind=%v", d.Kind)
		}
	case *ast.GateDecl:
		e.R.Func[d.Name] = d
	case *ast.FuncDecl:
		e.R.Func[d.Name] = d
	default:
		return fmt.Errorf("invalid decl=%v", s.Decl)
	}

	return nil
}

func (e *Evaluator) evalExprStmt(s *ast.ExprStmt) error {
	return e.evalExpr(s.X)
}

func (e *Evaluator) evalExpr(x ast.Expr) error {
	switch x := x.(type) {
	case *ast.ResetExpr:
		for _, a := range x.QArgs.List {
			qb, ok := e.R.Qubit.Get(a)
			if !ok {
				return fmt.Errorf("qubit=%#v not found", a)
			}

			e.Q.Reset(qb...)
		}
	case *ast.PrintExpr:
		if err := e.Println(x.QArgs.List...); err != nil {
			return fmt.Errorf("println: %v", err)
		}
	case *ast.MeasureExpr:
		if _, err := e.measure(x.QArgs.List...); err != nil {
			return fmt.Errorf("measure: %v", err)
		}
	case *ast.ApplyExpr:
		params := make([]int, 0)
		for _, p := range x.Params.List.List {
			if a, ok := e.R.Const[p.String()]; ok {
				params = append(params, a)
			}
		}

		qargs := make([][]q.Qubit, 0)
		for _, a := range x.QArgs.List {
			qb, ok := e.R.Qubit.Get(a)
			if !ok {
				return fmt.Errorf("qubit=%#v not found", a)
			}

			qargs = append(qargs, qb)
		}

		return e.apply(x.Kind, params, qargs)
	case *ast.CallExpr:
		return e.call(x)
	default:
		return fmt.Errorf("invalid expr=%#v", x)
	}

	return nil
}

func (e *Evaluator) evalAssignStmt(s *ast.AssignStmt) error {
	switch x := s.Right.(type) {
	case *ast.MeasureExpr:
		// left
		c, ok := e.R.Bit.Get(s.Left)
		if !ok {
			return fmt.Errorf("bit=%#v not found", s.Left)
		}

		// right
		m, err := e.measure(x.QArgs.List...)
		if err != nil {
			return fmt.Errorf("measure: %v", err)
		}

		// assign
		for i := range m {
			c[i] = m[i]
		}

		return nil
	case *ast.CallExpr:

		return nil
	default:
		return fmt.Errorf("invalid stmt=%#v", s)
	}
}

func (e *Evaluator) evalArrowStmt(s *ast.ArrowStmt) error {
	return e.evalAssignStmt(&ast.AssignStmt{
		Left:  s.Right,
		Right: s.Left,
	})
}

func (e *Evaluator) call(x *ast.CallExpr) error {
	decl, ok := e.R.Func[x.Name]
	if !ok {
		return fmt.Errorf("%v not found", x.Name)
	}

	switch f := decl.(type) {
	case *ast.GateDecl:
		return e.callGate(x, f)
	case *ast.FuncDecl:
		return e.callFunc(x, f)
	}

	return fmt.Errorf("invalid func=%#v", decl)
}

func (e *Evaluator) callGate(x *ast.CallExpr, g *ast.GateDecl) error {
	prms := make(map[string]ast.Expr)
	for i, p := range g.Params.List.List {
		prms[ast.Ident(p)] = x.Params.List.List[i]
	}

	args := make(map[string]ast.Expr)
	for i, a := range g.QArgs.List {
		args[ast.Ident(a)] = x.QArgs.List[i]
	}

	for _, b := range g.Body.List {
		switch s := b.(type) {
		case *ast.ExprStmt:
			switch X := s.X.(type) {
			case *ast.ApplyExpr:
				x := &ast.ApplyExpr{
					Kind: X.Kind,
				}

				for _, p := range X.Params.List.List {
					x.Params.List.Append(prms[ast.Ident(p)])
				}

				for _, a := range X.QArgs.List {
					x.QArgs.Append(args[ast.Ident(a)])
				}

				if err := e.evalExpr(x); err != nil {
					return fmt.Errorf("eval expr: %#v", err)
				}
			default:
				return fmt.Errorf("invalid expr=%#v", X)
			}
		default:
			return fmt.Errorf("invalid stmt=%#v", s)
		}
	}

	return nil
}

func (e *Evaluator) callFunc(x *ast.CallExpr, f *ast.FuncDecl) error {
	// TODO
	return fmt.Errorf("%v is not implemented", f)
}

func (e *Evaluator) measure(qargs ...ast.Expr) ([]int, error) {
	if len(qargs) == 0 {
		return nil, fmt.Errorf("qargs is empty")
	}

	m := make([]q.Qubit, 0)
	for _, a := range qargs {
		qb, ok := e.R.Qubit.Get(a)
		if !ok {
			return nil, fmt.Errorf("%#v not found", a)
		}

		e.Q.Measure(qb...)
		m = append(m, qb...)
	}

	bit := make([]int, 0)
	for _, q := range m {
		bit = append(bit, e.Q.State(q)[0].Int[0])
	}

	return bit, nil
}

func (e *Evaluator) apply(gate lexer.Token, params []int, qargs [][]q.Qubit) error {
	in := make([]q.Qubit, 0)
	for _, q := range qargs {
		in = append(in, q...)
	}

	switch gate {
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
		for i := range qargs[0] {
			e.Q.CNOT(qargs[0][i], qargs[1][i])
		}
	case lexer.CZ:
		for i := range qargs[0] {
			e.Q.CZ(qargs[0][i], qargs[1][i])
		}
	case lexer.CCX:
		for i := range qargs[0] {
			e.Q.CCNOT(qargs[0][i], qargs[1][i], qargs[2][i])
		}
	// itsubaki/q
	case lexer.SWAP:
		e.Q.Swap(in...)
	case lexer.QFT:
		e.Q.QFT(in...)
	case lexer.IQFT:
		e.Q.InvQFT(in...)
	case lexer.CMODEXP2:
		e.Q.CModExp2(params[0], params[1], qargs[0], qargs[1])
	default:
		return fmt.Errorf("gate=%v(%v) not found", gate, lexer.Tokens[gate])
	}

	return nil
}

func (e *Evaluator) Println(qargs ...ast.Expr) error {
	if len(e.R.Qubit.Name) == 0 {
		return nil
	}

	if len(qargs) == 0 {
		qargs = make([]ast.Expr, 0)
		for _, n := range e.R.Qubit.Name {
			qargs = append(qargs, &ast.IdentExpr{
				Value: n,
			})
		}
	}

	index := make([][]int, 0)
	for _, a := range qargs {
		qb, ok := e.R.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit=%#v not found", a)
		}

		index = append(index, q.Index(qb...))
	}

	for _, s := range e.Q.Raw().State(index...) {
		fmt.Println(s)
	}

	return nil
}
