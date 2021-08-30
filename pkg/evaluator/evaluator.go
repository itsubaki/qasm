package evaluator

import (
	"fmt"
	"strconv"

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

func (e *Evaluator) eval(s ast.Stmt) error {
	switch s := s.(type) {
	case *ast.DeclStmt:
		if err := e.evalDeclStmt(s); err != nil {
			return fmt.Errorf("eval decl stmt: %v", err)
		}
	case *ast.ExprStmt:
		if err := e.evalExprStmt(s); err != nil {
			return fmt.Errorf("eval expr stmt: %v", err)
		}
	case *ast.ApplyStmt:
		if err := e.evalApplyStmt(s); err != nil {
			return fmt.Errorf("eval apply stmt: %v", err)
		}
	case *ast.AssignStmt:
		if err := e.evalAssignStmt(s); err != nil {
			return fmt.Errorf("eval assign stmt: %v", err)
		}
	case *ast.ArrowStmt:
		if err := e.evalArrowStmt(s); err != nil {
			return fmt.Errorf("eval arrow stmt: %v", err)
		}
	case *ast.ResetStmt:
		if err := e.evalResetStmt(s); err != nil {
			return fmt.Errorf("eval reset stmt: %v", err)
		}
	case *ast.PrintStmt:
		if err := e.evalPrintStmt(s); err != nil {
			return fmt.Errorf("eval print stmt: %v", err)
		}
	default:
		return fmt.Errorf("unsupported stmt=%#v", s)
	}

	return nil
}

func (e *Evaluator) evalDeclStmt(s *ast.DeclStmt) error {
	switch decl := s.Decl.(type) {
	case *ast.GenConst:
		if _, ok := e.R.Const[ast.Ident(decl)]; ok {
			return fmt.Errorf("already exists=%v", decl.Name.Value)
		}

		switch x := decl.Value.(type) {
		case *ast.BasicLit:
			v, err := strconv.ParseFloat(x.Value, 64)
			if err != nil {
				return fmt.Errorf("parse float=%v: %v", x.Value, err)
			}

			e.R.Const[ast.Ident(decl)] = v
		default:
			return fmt.Errorf("unsupported expr=%#v", x)
		}
	case *ast.GenDecl:
		switch decl.Kind {
		case lexer.BIT:
			if _, ok := e.R.Bit.Get(&ast.IdentExpr{Value: decl.Name.Value}); ok {
				return fmt.Errorf("already exists=%v", decl.Name)
			}
			e.R.Bit.Add(ast.Ident(decl), make([]int, decl.Size()))
		case lexer.QUBIT:
			if _, ok := e.R.Qubit.Get(&ast.IdentExpr{Value: decl.Name.Value}); ok {
				return fmt.Errorf("already exists=%v", decl.Name)
			}
			qb := e.Q.ZeroWith(decl.Size())
			e.R.Qubit.Add(ast.Ident(decl), qb)
		default:
			return fmt.Errorf("unsupported kind=%v", decl.Kind)
		}
	case *ast.GateDecl, *ast.FuncDecl:
		e.R.Func[ast.Ident(decl)] = decl
	default:
		return fmt.Errorf("unsupported decl=%v", decl)
	}

	return nil
}

func (e *Evaluator) evalExprStmt(s *ast.ExprStmt) error {
	if _, err := e.evalExpr(s.X); err != nil {
		return fmt.Errorf("eval expr=%v", s.X)
	}

	return nil
}

func (e *Evaluator) evalExpr(x ast.Expr) ([]int, error) {
	switch x := x.(type) {
	case *ast.MeasureExpr:
		out, err := e.measure(x.QArgs.List...)
		if err != nil {
			return nil, fmt.Errorf("measure: %v", err)
		}

		return out, err
	case *ast.CallExpr:
		out, err := e.call(x)
		if err != nil {
			return nil, fmt.Errorf("call :%v", err)
		}

		return out, err
	default:
		return nil, fmt.Errorf("unsupported expr=%#v", x)
	}
}

func (e *Evaluator) evalApplyStmt(s *ast.ApplyStmt) error {
	if s.Kind == lexer.IDENT {
		x := &ast.CallExpr{
			Name:   s.Name,
			Params: s.Params,
			QArgs:  s.QArgs,
		}
		if _, err := e.call(x); err != nil {
			return fmt.Errorf("call=%s: %v", x, err)
		}

		return nil
	}

	params := make([]float64, 0)
	for _, p := range s.Params.List.List {
		if a, ok := e.R.Const[ast.Ident(p)]; ok {
			params = append(params, a)
			continue
		}

		switch p := p.(type) {
		case *ast.BasicLit:
			params = append(params, p.Float64())
		default:
			return fmt.Errorf("unsupported expr=%#v", p)
		}
	}

	qargs := make([][]q.Qubit, 0)
	for _, a := range s.QArgs.List {
		qb, ok := e.R.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit=%#v not found", a)
		}

		qargs = append(qargs, qb)
	}

	// TODO
	// modifier
	return e.apply(s.Kind, params, qargs)
}

func (e *Evaluator) evalAssignStmt(s *ast.AssignStmt) error {
	// left
	c, ok := e.R.Bit.Get(s.Left)
	if !ok {
		return fmt.Errorf("bit=%#v not found", s.Left)
	}

	// right
	m, err := e.evalExpr(s.Right)
	if err != nil {
		return fmt.Errorf("eval expr: %v", err)
	}

	// assign
	for i := range m {
		c[i] = m[i]
	}

	return nil
}

func (e *Evaluator) evalArrowStmt(s *ast.ArrowStmt) error {
	return e.evalAssignStmt(&ast.AssignStmt{
		Left:  s.Right,
		Right: s.Left,
	})
}

func (e *Evaluator) evalResetStmt(s *ast.ResetStmt) error {
	for _, a := range s.QArgs.List {
		qb, ok := e.R.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit=%#v not found", a)
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) evalPrintStmt(s *ast.PrintStmt) error {
	if len(e.R.Qubit.Name) == 0 {
		return nil
	}

	qargs := s.QArgs.List
	if len(qargs) == 0 {
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

func (e *Evaluator) call(x *ast.CallExpr) ([]int, error) {
	decl, ok := e.R.Func[x.Name]
	if !ok {
		return nil, fmt.Errorf("%v not found", x.Name)
	}

	switch decl := decl.(type) {
	case *ast.GateDecl:
		return []int{}, e.callGate(x, decl)
	case *ast.FuncDecl:
		return e.callFunc(x, decl)
	}

	return nil, fmt.Errorf("unsupported decl=%#v", decl)
}

func (e *Evaluator) callGate(x *ast.CallExpr, decl *ast.GateDecl) error {
	params := make(map[string]ast.Expr)
	for i, p := range decl.Params.List.List {
		params[ast.Ident(p)] = x.Params.List.List[i]
	}

	// gate bell q, p { h q; cx q, p; }
	// bell q0, q1;
	// q -> q0, p -> q1
	qargs := make(map[string]ast.Expr)
	for i, a := range decl.QArgs.List {
		qargs[ast.Ident(a)] = x.QArgs.List[i]
	}

	for _, b := range decl.Body.List {
		switch s := b.(type) {
		case *ast.ApplyStmt:
			if err := e.eval(&ast.ApplyStmt{
				Kind:     s.Kind,
				Name:     s.Name,
				Modifier: s.Modifier,
				Params: ast.ParenExpr{
					List: assign(s.Params.List, params),
				},
				// q -> q0, p -> q1
				QArgs: assign(s.QArgs, qargs),
			}); err != nil {
				return fmt.Errorf("eval: %#v", err)
			}
		default:
			return fmt.Errorf("unsupported stmt=%#v", s)
		}
	}

	return nil
}

func (e *Evaluator) callFunc(x *ast.CallExpr, decl *ast.FuncDecl) ([]int, error) {
	params := make(map[string]ast.Expr)
	for i, p := range decl.Params.List.List {
		params[ast.Ident(p)] = x.Params.List.List[i]
	}

	qargs := make(map[string]ast.Expr)
	for i, a := range decl.QArgs.List {
		qargs[ast.Ident(a)] = x.QArgs.List[i]
	}

	for _, b := range decl.Body.List {
		switch s := b.(type) {
		case *ast.ApplyStmt:
			if err := e.eval(&ast.ApplyStmt{
				Kind:     s.Kind,
				Name:     s.Name,
				Modifier: s.Modifier,
				Params: ast.ParenExpr{
					List: assign(s.Params.List, params),
				},
				QArgs: assign(s.QArgs, qargs),
			}); err != nil {
				return nil, fmt.Errorf("eval: %#v", err)
			}
		case *ast.ReturnStmt:
			switch r := s.Result.(type) {
			case *ast.MeasureExpr:
				out, err := e.evalExpr(&ast.MeasureExpr{
					QArgs: assign(r.QArgs, qargs),
				})
				if err != nil {
					return nil, fmt.Errorf("eval expr: %#v", err)
				}

				return out, nil
			case nil:
				// no return value
				return []int{}, nil
			default:
				return nil, fmt.Errorf("unsupported expr=%#v", x)
			}
		default:
			return nil, fmt.Errorf("unsupported stmt=%#v", s)
		}
	}

	return nil, nil
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

func (e *Evaluator) apply(g lexer.Token, p []float64, qargs [][]q.Qubit) error {
	in := make([]q.Qubit, 0)
	for _, q := range qargs {
		in = append(in, q...)
	}

	switch g {
	case lexer.U:
		e.Q.U(p[0], p[1], p[2], in...)
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
		e.Q.CModExp2(int(p[0]), int(p[1]), qargs[0], qargs[1])
	default:
		return fmt.Errorf("gate=%v(%v) not found", g, lexer.Tokens[g])
	}

	return nil
}

func assign(c ast.ExprList, args map[string]ast.Expr) ast.ExprList {
	out := ast.ExprList{}
	for _, a := range c.List {
		switch x := a.(type) {
		case *ast.BasicLit:
			out.Append(x)
		case *ast.IndexExpr:
			out.Append(&ast.IndexExpr{
				Name: ast.IdentExpr{
					Value: ast.Ident(args[ast.Ident(a)]),
				},
				Value: x.Value,
			})
		default:
			out.Append(args[ast.Ident(a)])
		}
	}

	return out
}

func (e *Evaluator) Println() error {
	if err := e.eval(&ast.PrintStmt{}); err != nil {
		return fmt.Errorf("print qubit: %v", err)
	}

	if err := e.R.Bit.Println(); err != nil {
		return fmt.Errorf("print bit: %v", err)
	}

	return nil
}
