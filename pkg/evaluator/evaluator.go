package evaluator

import (
	"fmt"
	"math"
	"strings"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator/object"
	"github.com/itsubaki/qasm/pkg/lexer"
)

const indent = ".  "

type Evaluator struct {
	Q      *q.Q
	Env    *object.Environment
	Opts   Opts
	indent int
}

type Opts struct {
	Verbose bool
}

func New(qsim *q.Q, opts ...Opts) *Evaluator {
	e := &Evaluator{
		Q:   qsim,
		Env: object.NewEnvironment(),
	}

	if opts != nil {
		e.Opts = opts[0]
	}

	return e
}

func Default(opts ...Opts) *Evaluator {
	return New(q.New(), opts...)
}

func (e *Evaluator) Eval(p *ast.OpenQASM) error {
	for _, s := range p.Stmts {
		if _, err := e.eval(s, e.Env); err != nil {
			return fmt.Errorf("eval(%v): %v", s, err)
		}
	}

	return nil
}

func (e *Evaluator) eval(n ast.Node, env *object.Environment) (obj object.Object, err error) {
	defer func() {
		if obj == nil || obj.Type() == object.NIL {
			e.indent--
			return
		}

		if e.Opts.Verbose {
			fmt.Printf("%v", strings.Repeat(indent, e.indent+1))
			fmt.Printf("return %T(%v)\n", obj, obj)
			e.indent--
		}
	}()

	if e.Opts.Verbose {
		e.indent++
		fmt.Printf("%v", strings.Repeat(indent, e.indent))
		fmt.Printf("%T(%v)\n", n, n)
	}

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

			if v.Type() == object.RETURN_VALUE {
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
		case lexer.PI:
			return &object.Float{Value: math.Pi}, nil
		case lexer.TAU:
			return &object.Float{Value: math.Pi * 2}, nil
		case lexer.EULER:
			return &object.Float{Value: math.E}, nil
		}
	case *ast.IdentExpr:
		if v, ok := env.Const[ast.Ident(n)]; ok {
			return v, nil
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
	if s.Kind == lexer.IDENT {
		x := &ast.CallExpr{
			Name:     s.Name,
			Modifier: s.Modifier,
			Params:   s.Params,
			QArgs:    s.QArgs,
		}

		if _, err := e.eval(x, env); err != nil {
			return fmt.Errorf("eval(%v): %v", x, err)
		}

		return nil
	}

	obj := make([]object.Object, 0)
	for _, p := range s.Params.List.List {
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

func (e *Evaluator) apply(mod []ast.Modifier, g lexer.Token, p []float64, qargs [][]q.Qubit) error {
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
	var c int
	for _, m := range mod {
		if m.Kind == lexer.INV {
			c++
		}
	}

	if c%2 == 0 {
		u = u.Dagger()
	}

	// Controlled-U
	for _, m := range mod {
		if m.Kind == lexer.INV || m.Kind == lexer.POW {
			continue
		}

		var c int
		if len(m.Index.List.List) > 0 {
			c = int(m.Index.List.List[0].(*ast.BasicLit).Float64())
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

func (e *Evaluator) callGate(x *ast.CallExpr, d *ast.GateDecl, outer *object.Environment) (object.Object, error) {
	env := e.extend(x, d, outer)

	for _, b := range d.Body.List {
		switch a := b.(type) {
		case *ast.ApplyStmt:
			if a.Kind != lexer.IDENT && len(x.QArgs.List) == len(d.QArgs.List) {
				if _, err := e.eval(a, env); err != nil {
					return nil, fmt.Errorf("eval(%v): %v", &d.Body, err)
				}
				continue
			}

			if e.Opts.Verbose {
				fmt.Printf("%v", strings.Repeat(indent, e.indent))
				fmt.Printf("%T(%v)\n", d, a)
			}

			if a.Kind != lexer.IDENT {
				a.QArgs = x.QArgs
				a.Modifier = append(x.Modifier, a.Modifier...)
				if _, err := e.eval(a, outer); err != nil {
					return nil, fmt.Errorf("eval(%v): %v", &d.Body, err)
				}
				continue
			}

			//	call declared gate
			decl := env.Func[a.Name].(*ast.GateDecl)
			for j := range decl.Body.List {
				s := &ast.ApplyStmt{
					Kind:     decl.Body.List[j].(*ast.ApplyStmt).Kind,
					Name:     decl.Body.List[j].(*ast.ApplyStmt).Name,
					Params:   decl.Body.List[j].(*ast.ApplyStmt).Params,
					QArgs:    x.QArgs,
					Modifier: append(x.Modifier, append(a.Modifier, decl.Body.List[j].(*ast.ApplyStmt).Modifier...)...),
				}

				if _, err := e.eval(s, outer); err != nil {
					return nil, fmt.Errorf("eval(%v): %v", &d.Body, err)
				}
			}

		default:
			if _, err := e.eval(b, env); err != nil {
				return nil, fmt.Errorf("eval(%v): %v", &d.Body, err)
			}
		}
	}

	return nil, nil
}

func (e *Evaluator) extend(x *ast.CallExpr, d *ast.GateDecl, outer *object.Environment) *object.Environment {
	env := object.NewEnclosedEnvironment(outer)
	env.Func = outer.Func
	env.Const = outer.Const

	for i := range d.QArgs.List {
		v, ok := outer.Qubit.Get(x.QArgs.List[i])
		if !ok {
			panic(fmt.Sprintf("qubit(%v) not found", x.QArgs.List[i]))
		}

		env.Qubit.Add(d.QArgs.List[i], v)
	}

	return env
}

func (e *Evaluator) callFunc(x *ast.CallExpr, d *ast.FuncDecl, outer *object.Environment) (object.Object, error) {
	env := e.extendFunc(x, d, outer)

	v, err := e.eval(&d.Body, env)
	if err != nil {
		return nil, fmt.Errorf("eval(%v): %v", &d.Body, err)
	}

	return v.(*object.ReturnValue).Value, nil
}

func (e *Evaluator) extendFunc(x *ast.CallExpr, d *ast.FuncDecl, outer *object.Environment) *object.Environment {
	env := object.NewEnclosedEnvironment(outer)
	env.Func = outer.Func
	env.Const = outer.Const

	for i := range d.QArgs.List {
		v, ok := outer.Qubit.Get(x.QArgs.List[i])
		if !ok {
			panic(fmt.Sprintf("qubit(%v) not found", x.QArgs.List[i]))
		}

		env.Qubit.Add(d.QArgs.List[i], v)
	}

	return env
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
			return fmt.Errorf("bit(%v) not found", n)
		}

		for _, v := range c {
			fmt.Printf("%v", v)
		}

		fmt.Println()
	}

	return nil
}
