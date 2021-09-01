package evaluator

import (
	"fmt"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator/object"
	"github.com/itsubaki/qasm/pkg/lexer"
)

type Evaluator struct {
	Q   *q.Q
	Env *object.Environment
}

func New(qsim *q.Q) *Evaluator {
	return &Evaluator{
		Q:   qsim,
		Env: object.NewEnvironment(),
	}
}

func Default() *Evaluator {
	return New(q.New())
}

func (e *Evaluator) Eval(p *ast.OpenQASM) error {
	for _, s := range p.Stmts {
		if _, err := e.eval(s, e.Env); err != nil {
			return fmt.Errorf("eval(%v): %v", s, err)
		}
	}

	return nil
}

func (e *Evaluator) eval(n ast.Node, env *object.Environment) (object.Object, error) {
	switch n := n.(type) {
	case *ast.ExprStmt:
		return e.eval(n.X, env)

	case *ast.DeclStmt:
		return e.eval(n.Decl, env)

	case *ast.ArrowStmt:
		return e.eval(&ast.AssignStmt{Left: n.Right, Right: n.Left}, env)

	case *ast.ReturnStmt:
		v, err := e.eval(n.Result, env)
		if err != nil {
			return nil, fmt.Errorf("return: %v", err)
		}

		return &object.ReturnValue{Value: v}, nil

	case *ast.AssignStmt:
		rhs, err := e.eval(n.Right, env)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", n.Right, err)
		}

		c, ok := env.Bit.Get(n.Left)
		if !ok {
			return nil, fmt.Errorf("bit(%v) not found", n.Left)
		}

		e := rhs.(*object.Array).Elm
		for i := range e {
			c[i] = e[i].(*object.Int).Value
		}

	case *ast.BlockStmt:
		for _, b := range n.List {
			v, err := e.eval(b, env)
			if err != nil {
				return nil, fmt.Errorf("eval(%v): %v", b, err)
			}

			if v.Type() == object.ReturnValueType {
				return v, nil
			}
		}

	case *ast.ResetStmt:
		if err := e.evalReset(n, env); err != nil {
			return nil, fmt.Errorf("apply(%v): %v", n, err)
		}

	case *ast.PrintStmt:
		if err := e.evalPrint(n, env); err != nil {
			return nil, fmt.Errorf("print(%v): %v", n, err)
		}

	case *ast.ApplyStmt:
		if err := e.evalApply(n, env); err != nil {
			return nil, fmt.Errorf("apply(%v): %v", n, err)
		}

	case *ast.GenConst:
		// TODO check already exists
		v, err := e.eval(n.Value, env)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", n, err)
		}
		env.Const[ast.Ident(n)] = v

	case *ast.GenDecl:
		// TODO check already exists
		switch n.Kind {
		case lexer.BIT:
			env.Bit.Add(n, make([]int64, n.Size()))
		case lexer.QUBIT:
			env.Qubit.Add(n, e.Q.ZeroWith(n.Size()))
		}

	case *ast.GateDecl:
		// TODO check already exists
		env.Func[ast.Ident(n)] = n

	case *ast.FuncDecl:
		// TODO check already exists
		env.Func[ast.Ident(n)] = n

	case *ast.CallExpr:
		return e.call(n, env)

	case *ast.MeasureExpr:
		return e.measure(n, env)

	case *ast.BasicLit:
		switch n.Kind {
		case lexer.INT:
			return &object.Int{Value: n.Int64()}, nil
		case lexer.FLOAT:
			return &object.Float{Value: n.Float64()}, nil
		case lexer.STRING:
			return &object.String{Value: n.Value}, nil
		}
	}

	return &object.Nil{}, nil
}

func (e *Evaluator) evalPrint(s *ast.PrintStmt, env *object.Environment) error {
	if len(env.Qubit.Name) == 0 {
		return nil
	}

	qargs := s.QArgs.List
	if len(qargs) == 0 {
		for _, n := range env.Qubit.Name {
			qargs = append(qargs, &ast.IdentExpr{Value: n})
		}
	}

	index := make([][]int, 0)
	for _, a := range qargs {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit(%v) not found", a)
		}

		index = append(index, q.Index(qb...))
	}

	for _, s := range e.Q.Raw().State(index...) {
		fmt.Println(s)
	}

	return nil
}

func (e *Evaluator) evalReset(s *ast.ResetStmt, env *object.Environment) error {
	for _, a := range s.QArgs.List {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit(%v) not found", a)
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) evalApply(s *ast.ApplyStmt, env *object.Environment) error {
	obj := make([]object.Object, 0)
	for _, p := range s.Params.List.List {
		if a, ok := env.Const[ast.Ident(p)]; ok {
			obj = append(obj, a)
			continue
		}

		v, err := e.eval(p, env)
		if err != nil {
			return fmt.Errorf("eval(%v): %v", p, err)
		}

		obj = append(obj, v)
	}

	p := make([]float64, 0)
	for _, o := range obj {
		switch o := o.(type) {
		case *object.Float:
			p = append(p, o.Value)
		case *object.Int:
			p = append(p, float64(o.Value))
		default:
			return fmt.Errorf("unsupported(%v)", o)
		}
	}

	q := make([][]q.Qubit, 0)
	for _, a := range s.QArgs.List {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit(%v) not found", a)
		}

		q = append(q, qb)
	}

	return e.apply(s.Modifier, s.Kind, p, q)
}

func (e *Evaluator) apply(mod []ast.Modifiler, g lexer.Token, p []float64, qargs [][]q.Qubit) error {
	in := make([]q.Qubit, 0)
	for _, q := range qargs {
		in = append(in, q...)
	}

	var u matrix.Matrix
	switch g {
	case lexer.U:
		u = gate.U(p[0], p[1], p[2])
	case lexer.X:
		u = gate.X()
	case lexer.Y:
		u = gate.Y()
	case lexer.Z:
		u = gate.Z()
	case lexer.H:
		u = gate.H()
	case lexer.T:
		u = gate.T()
	case lexer.S:
		u = gate.S()
	case lexer.CX:
		for i := range qargs[0] {
			e.Q.CNOT(qargs[0][i], qargs[1][i])
		}
		return nil
	case lexer.CZ:
		for i := range qargs[0] {
			e.Q.CZ(qargs[0][i], qargs[1][i])
		}
		return nil
	case lexer.CCX:
		for i := range qargs[0] {
			e.Q.CCNOT(qargs[0][i], qargs[1][i], qargs[2][i])
		}
		return nil
	case lexer.SWAP:
		e.Q.Swap(in...)
		return nil
	case lexer.QFT:
		e.Q.QFT(in...)
		return nil
	case lexer.IQFT:
		e.Q.InvQFT(in...)
		return nil
	case lexer.CMODEXP2:
		e.Q.CModExp2(int(p[0]), int(p[1]), qargs[0], qargs[1])
		return nil
	default:
		return fmt.Errorf("gate=%v(%v) not found", lexer.Tokens[g], g)
	}

	// Inverse U
	for _, m := range mod {
		if m.Kind == lexer.INV {
			u = u.Dagger()
		}
	}

	// Controlled-U
	for _, m := range mod {
		if m.Kind == lexer.INV || m.Kind == lexer.POW {
			continue
		}

		var c int
		if len(m.Index.List.List) > 0 {
			switch x := m.Index.List.List[0].(type) {
			case *ast.BasicLit:
				c = int(x.Float64())
			}
		}

		switch m.Kind {
		case lexer.CTRL:
			for i := range qargs[c] {
				e.Q.C(u, qargs[c][i], qargs[len(qargs)-1][i])
			}

			return nil
		case lexer.NEGCTRL:
			for i := range qargs[c] {
				e.Q.X(qargs[c][i])
			}
			for i := range qargs[c] {
				e.Q.C(u, qargs[c][i], qargs[len(qargs)-1][i])
			}
			for i := range qargs[c] {
				e.Q.X(qargs[c][i])
			}

			return nil
		}
	}

	// U
	e.Q.Apply(u, in...)
	return nil
}

func (e *Evaluator) call(x *ast.CallExpr, env *object.Environment) (object.Object, error) {
	decl, ok := env.Func[x.Name]
	if !ok {
		return nil, fmt.Errorf("decl(%v) not found", x.Name)
	}

	switch decl := decl.(type) {
	case *ast.GateDecl:
		return e.callGate(x, decl, env)
	case *ast.FuncDecl:
		return e.callFunc(x, decl, env)
	}

	return nil, fmt.Errorf("unsupported(%v)", decl)
}

func (e *Evaluator) callGate(x *ast.CallExpr, decl *ast.GateDecl, env *object.Environment) (object.Object, error) {
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
			a := &ast.ApplyStmt{
				Kind:     s.Kind,
				Name:     s.Name,
				Modifier: s.Modifier,
				Params: ast.ParenExpr{
					List: assign(s.Params.List, params),
				},
				QArgs: assign(s.QArgs, qargs), // q -> q0, p -> q1
			}

			for a.Kind == lexer.IDENT {
				decl := env.Func[a.Name].(*ast.GateDecl)
				a.Kind = decl.Body.List[0].(*ast.ApplyStmt).Kind
				a.Name = decl.Body.List[0].(*ast.ApplyStmt).Name
				a.Params = decl.Body.List[0].(*ast.ApplyStmt).Params
			}

			if _, err := e.eval(a, env); err != nil {
				return nil, fmt.Errorf("apply(%v): %v", a, err)
			}
		default:
			return nil, fmt.Errorf("unsupported(%v)", s)
		}
	}

	return nil, nil
}

func (e *Evaluator) callFunc(x *ast.CallExpr, decl *ast.FuncDecl, env *object.Environment) (object.Object, error) {
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
			a := &ast.ApplyStmt{
				Kind:     s.Kind,
				Name:     s.Name,
				Modifier: s.Modifier,
				Params: ast.ParenExpr{
					List: assign(s.Params.List, params),
				},
				QArgs: assign(s.QArgs, qargs),
			}
			if _, err := e.eval(a, env); err != nil {
				return nil, fmt.Errorf("apply(%v): %#v", a, err)
			}
		case *ast.ReturnStmt:
			switch r := s.Result.(type) {
			case *ast.MeasureExpr:
				m := &ast.MeasureExpr{
					QArgs: assign(r.QArgs, qargs),
				}
				out, err := e.eval(m, env)
				if err != nil {
					return nil, fmt.Errorf("measure(%v): %#v", m, err)
				}

				return out, nil
			case nil:
				// no return value
				return &object.Nil{}, nil
			default:
				return nil, fmt.Errorf("unsupported(%v)", x)
			}
		default:
			return nil, fmt.Errorf("unsupported(%v)", s)
		}
	}

	return nil, nil
}

func (e *Evaluator) measure(x *ast.MeasureExpr, env *object.Environment) (object.Object, error) {
	qargs := x.QArgs.List
	if len(qargs) == 0 {
		return nil, fmt.Errorf("qargs is empty")
	}

	m := make([]q.Qubit, 0)
	for _, a := range qargs {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return nil, fmt.Errorf("qubit(%v) not found", a)
		}

		e.Q.Measure(qb...)
		m = append(m, qb...)
	}

	bit := make([]object.Object, 0)
	for _, q := range m {
		bit = append(bit, &object.Int{Value: int64(e.Q.State(q)[0].Int[0])})
	}

	return &object.Array{Elm: bit}, nil
}

func (e *Evaluator) Println() error {
	if _, err := e.eval(&ast.PrintStmt{}, e.Env); err != nil {
		return fmt.Errorf("print qubit: %v", err)
	}

	for _, n := range e.Env.Bit.Name {
		fmt.Printf("%v: ", n)

		c, ok := e.Env.Bit.Get(&ast.IdentExpr{Value: n})
		if !ok {
			return fmt.Errorf("name=%v not found", n)
		}

		for _, v := range c {
			fmt.Printf("%v", v)
		}

		fmt.Println()
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
