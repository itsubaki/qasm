package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator/env"
	"github.com/itsubaki/qasm/pkg/evaluator/object"
	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
)

const indent = ".  "

type Evaluator struct {
	Q      *q.Q
	Env    *env.Environment
	Opts   Opts
	indent int
}

type Opts struct {
	Verbose bool
}

func New(qsim *q.Q, env *env.Environment, opts ...Opts) *Evaluator {
	e := &Evaluator{
		Q:   qsim,
		Env: env,
	}

	if opts != nil {
		e.Opts = opts[0]
	}

	return e
}

func Default(opts ...Opts) *Evaluator {
	return New(q.New(), env.New(), opts...)
}

func Eval(n ast.Node) (object.Object, error) {
	return Default().eval(n, env.New())
}

func (e *Evaluator) Eval(p *ast.OpenQASM) error {
	if e.Opts.Verbose {
		fmt.Printf("%T\n", p)

		e.indent++
		if p.Version != nil {
			fmt.Printf("%v", strings.Repeat(indent, e.indent))
			fmt.Printf("%T(%v)\n", p.Version, p.Version)
		}
		fmt.Printf("%v", strings.Repeat(indent, e.indent))
		fmt.Printf("%T\n", p.Stmts)
	}

	for _, s := range p.Stmts {
		if _, err := e.eval(s, e.Env); err != nil {
			return fmt.Errorf("eval(%v): %v", s, err)
		}
	}

	return nil
}

func (e *Evaluator) eval(n ast.Node, env *env.Environment) (obj object.Object, err error) {
	if e.Opts.Verbose {
		defer func() {
			if obj != nil && obj.Type() != object.NIL {
				fmt.Printf("%v", strings.Repeat(indent, e.indent+1))
				fmt.Printf("return %T(%v)\n", obj, obj)
			}

			e.indent--
		}()

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

	case *ast.ResetStmt:
		if err := e.evalReset(n, env); err != nil {
			return nil, fmt.Errorf("apply(%v): %v", n, err)
		}

	case *ast.ApplyStmt:
		if err := e.evalApply(n, env); err != nil {
			return nil, fmt.Errorf("apply(%v): %v", n, err)
		}

	case *ast.PrintStmt:
		if err := e.evalPrint(n, env); err != nil {
			return nil, fmt.Errorf("print(%v): %v", n, err)
		}

	case *ast.InclStmt:
		path := strings.Trim(n.Path.Value, "\"")
		f, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read file=%s: %v", path, err)
		}

		l := lexer.New(strings.NewReader(string(f)))
		p := parser.New(l)

		a := p.Parse()
		if errs := p.Errors(); len(errs) != 0 {
			return nil, fmt.Errorf("parse: %v", errs)
		}

		for _, s := range a.Stmts {
			if _, err := e.eval(s, e.Env); err != nil {
				return nil, fmt.Errorf("eval(%v): %v", s, err)
			}
		}

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

			if v != nil && v.Type() == object.RETURN_VALUE {
				return v, nil
			}
		}

	case *ast.ReturnStmt:
		v, err := e.eval(n.Result, env)
		if err != nil {
			return nil, fmt.Errorf("return: %v", err)
		}

		return &object.ReturnValue{Value: v}, nil

	case *ast.GenConst:
		v, err := e.eval(n.Value, env)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", n, err)
		}
		env.Const[ast.Must(ast.Ident(n))] = v

	case *ast.GenDecl:
		switch n.Kind {
		case lexer.BIT:
			env.Bit.Add(n, make([]int64, n.Size()))
		case lexer.QUBIT:
			env.Qubit.Add(n, e.Q.ZeroWith(n.Size()))
		}

	case *ast.GateDecl:
		env.Func[ast.Must(ast.Ident(n))] = n

	case *ast.FuncDecl:
		env.Func[ast.Must(ast.Ident(n))] = n

	case *ast.CallExpr:
		return e.call(n, env)

	case *ast.MeasureExpr:
		return e.measure(n, env)

	case *ast.UnaryExpr:
		v, err := e.eval(n.Value, env)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", n.Value, err)
		}

		switch n.Kind {
		case lexer.PLUS:
			return v, nil
		case lexer.MINUS:
			switch v := v.(type) {
			case *object.Int:
				return &object.Int{Value: -1 * v.Value}, nil
			case *object.Float:
				return &object.Float{Value: -1 * v.Value}, nil
			}
		}

	case *ast.InfixExpr:
		lhs, err := e.eval(n.Left, env)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", n.Left, err)
		}

		rhs, err := e.eval(n.Right, env)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", n.Right, err)
		}

		return e.evalInfix(n.Kind, lhs, rhs)

	case *ast.IdentExpr:
		if v, ok := env.Const[ast.Must(ast.Ident(n))]; ok {
			return v, nil
		}

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
	}

	return &object.Nil{}, nil
}

func (e *Evaluator) evalInfix(kind lexer.Token, lhs, rhs object.Object) (object.Object, error) {
	switch kind {
	case lexer.PLUS:
		switch lhs := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: lhs.Value + rhs.(*object.Int).Value}, nil
		case *object.Float:
			return &object.Float{Value: lhs.Value + rhs.(*object.Float).Value}, nil
		}

	case lexer.MINUS:
		switch lhs := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: lhs.Value - rhs.(*object.Int).Value}, nil
		case *object.Float:
			return &object.Float{Value: lhs.Value - rhs.(*object.Float).Value}, nil
		}

	case lexer.MUL:
		switch lhs := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: lhs.Value * rhs.(*object.Int).Value}, nil
		case *object.Float:
			return &object.Float{Value: lhs.Value * rhs.(*object.Float).Value}, nil
		}

	case lexer.DIV:
		switch lhs := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: lhs.Value / rhs.(*object.Int).Value}, nil
		case *object.Float:
			return &object.Float{Value: lhs.Value / rhs.(*object.Float).Value}, nil
		}

	case lexer.MOD:
		switch lhs := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: lhs.Value % rhs.(*object.Int).Value}, nil
		}

	}

	return nil, fmt.Errorf("unsupported(%v)", kind)
}

func (e *Evaluator) evalPrint(s *ast.PrintStmt, env *env.Environment) error {
	if len(env.Qubit.Name) == 0 {
		return nil
	}

	qargs := s.QArgs.List
	if len(qargs) == 0 {
		for _, n := range env.Qubit.Name {
			qargs = append(qargs, &ast.IdentExpr{Name: n})
		}
	}

	var index [][]int
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

func (e *Evaluator) evalReset(s *ast.ResetStmt, env *env.Environment) error {
	for _, a := range s.QArgs.List {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit(%v) not found", a)
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) evalApply(s *ast.ApplyStmt, env *env.Environment) error {
	var params []float64
	for _, p := range s.Params.List.List {
		v, err := e.eval(p, env)
		if err != nil {
			return fmt.Errorf("eval(%v): %v", p, err)
		}

		switch o := v.(type) {
		case *object.Float:
			params = append(params, o.Value)
		case *object.Int:
			params = append(params, float64(o.Value))
		default:
			return fmt.Errorf("unsupported(%v)", p)
		}
	}

	var qargs [][]q.Qubit
	for _, a := range s.QArgs.List {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit(%v) not found", a)
		}

		qargs = append(qargs, qb)
	}

	return e.apply(s.Modifier, s.Kind, params, qargs, env)
}

func (e *Evaluator) measure(x *ast.MeasureExpr, env *env.Environment) (object.Object, error) {
	qargs := x.QArgs.List
	if len(qargs) == 0 {
		return nil, fmt.Errorf("qargs is empty")
	}

	var m []q.Qubit
	for _, a := range qargs {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return nil, fmt.Errorf("qubit(%v) not found", a)
		}

		m = append(m, qb...)
	}

	e.Q.Measure(m...)

	var bit []object.Object
	for _, q := range m {
		bit = append(bit, &object.Int{Value: int64(e.Q.State(q)[0].Int[0])})
	}

	return &object.Array{Elm: bit}, nil
}

func (e *Evaluator) apply(mod []ast.Modifier, g lexer.Token, params []float64, qargs [][]q.Qubit, env *env.Environment) error {
	// QFT, IQFT, CMODEXP2
	if e.builtinApply(g, params, qargs) {
		return nil
	}

	// U, X, Y, Z, H, S, T
	u, ok := builtin(g, params)
	if !ok {
		return fmt.Errorf("gate=%v(%v) not found", lexer.Tokens[g], g)
	}

	// inv @ U
	if len(ast.ModInv(mod))%2 == 1 {
		u = u.Dagger()
	}

	// pow(2) @ U
	u, err := e.pow(mod, u, env)
	if err != nil {
		return fmt.Errorf("pow(%v): %v", mod, err)
	}

	// ctrl/negc @ U
	ctrl, negc := e.ctrl(mod, qargs, env)
	if len(ctrl)+len(negc) > 0 {
		e.ctrlApply(mod, ctrl, negc, u, qargs)
		return nil
	}

	// no ctrl/negc
	e.Q.Apply(u, flatten(qargs)...)
	return nil
}

func (e *Evaluator) builtinApply(g lexer.Token, p []float64, qargs [][]q.Qubit) bool {
	switch g {
	case lexer.QFT:
		e.Q.QFT(flatten(qargs)...)
		return true
	case lexer.IQFT:
		e.Q.InvQFT(flatten(qargs)...)
		return true
	case lexer.CMODEXP2:
		e.Q.CModExp2(int(p[0]), int(p[1]), qargs[0], qargs[1])
		return true
	}

	return false
}

func (e *Evaluator) ctrl(mod []ast.Modifier, qargs [][]q.Qubit, env *env.Environment) ([]q.Qubit, []q.Qubit) {
	var ctrl, negc []q.Qubit
	for i, m := range ast.ModCtrl(mod) {
		switch m.Kind {
		case lexer.CTRL:
			ctrl = append(ctrl, qargs[i]...)
		case lexer.NEGCTRL:
			negc = append(negc, qargs[i]...)
		}
	}

	return ctrl, negc
}

func (e *Evaluator) ctrlApply(mod []ast.Modifier, ctrl, negc []q.Qubit, u matrix.Matrix, qargs [][]q.Qubit) {
	c := append(ctrl, negc...)
	t := qargs[len(qargs)-1]
	m := ast.ModCtrl(mod)

	if len(c) == len(t) {
		// qubit[2] q;
		// qubit[2] r;
		// ctrl @ x q, r;
		//
		// equals to
		// ctrl @ x q[0], r[0];
		// ctrl @ x q[1], r[1];
		for i := range c {
			e.negc(negc, func() { e.Q.C(u, c[i], t[i]) })
		}

		return
	}

	if len(c) == 1 {
		// qubit[2] q;
		// qubit[2] r;
		// ctrl @ x q[0], r;
		//
		// equals to
		// ctrl @ x q[0], r[0];
		// ctrl @ x q[0], r[1];
		//
		// or
		// ctrl @ x q[0], r[0];
		for i := range t {
			e.negc(negc, func() { e.Q.C(u, c[0], t[i]) })
		}

		return
	}

	if len(t) == 1 && len(m) == 1 {
		// qubit[2] q;
		// qubit[2] r;
		// ctrl @ x q, r[0];
		//
		// equals to
		// ctrl @ q[0], r[0];
		// ctrl @ q[1], r[0];
		for i := range c {
			e.negc(negc, func() { e.Q.C(u, c[i], t[0]) })
		}

		return
	}

	// ctrl @ ctrl @ x q[0], q[1], r[0];
	for i := range t {
		e.negc(negc, func() { e.Q.Controlled(u, c, t[i]) })
	}

	// len(c) > 1 && len(t) > 1 && len(c) != len(t)
	// ?
}

func (e *Evaluator) negc(negc []q.Qubit, f func()) {
	if len(negc) > 0 {
		e.Q.X(negc...)
	}

	f()

	if len(negc) > 0 {
		e.Q.X(negc...)
	}
}

func (e *Evaluator) pow(mod []ast.Modifier, u matrix.Matrix, env *env.Environment) (matrix.Matrix, error) {
	// U
	pow := ast.ModPow(mod)
	if len(pow) == 0 {
		return u, nil
	}

	// pow(2) @ pow(-2) @ U equals to pow(0) @ U
	var p int
	for _, m := range pow {
		n := m.Index.List.List[0]
		v, err := e.eval(n, env)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", n, err)
		}

		p = p + int(v.(*object.Int).Value)
	}

	// pow(0) equals to identity
	if p == 0 {
		return gate.I(), nil
	}

	// pow(-1) equals to inv
	if p < 0 {
		p = -1 * p
		u = u.Dagger()
	}

	// apply pow
	out := u
	for i := 1; i < p; i++ {
		out = out.Apply(u)
	}

	return out, nil
}

func (e *Evaluator) call(x *ast.CallExpr, env *env.Environment) (object.Object, error) {
	f, ok := env.Func[x.Name]
	if !ok {
		return nil, fmt.Errorf("decl(%v) not found", x.Name)
	}

	if e.Opts.Verbose {
		fmt.Printf("%v", strings.Repeat(indent, e.indent))
		fmt.Printf("%T(%v)\n", f, f)
	}

	switch f := f.(type) {
	case *ast.GateDecl:
		return &object.Nil{}, e.applyg(x, f, env)
	case *ast.FuncDecl:
		return e.applyf(x, f, env)
	}

	return nil, fmt.Errorf("unsupported. call=%v", f)
}

func (e *Evaluator) ctrlCall(x *ast.CallExpr, body ast.BlockStmt, env *env.Environment) error {
	ctrl := ast.ModCtrl(x.Modifier)
	if len(ctrl) > 0 {
		body = override(addmod(body, ctrl), env.Qubit.Name)
	}

	// bell q0, q1;
	if _, err := e.eval(&body, env); err != nil {
		return fmt.Errorf("eval(%v): %v", &body, err)
	}

	return nil
}

func (e *Evaluator) applyg(x *ast.CallExpr, g *ast.GateDecl, outer *env.Environment) error {
	// extend
	env, err := e.extg(x, g, outer)
	if err != nil {
		return fmt.Errorf("extend: %v", err)
	}

	// inv
	body := g.Body
	if len(ast.ModInv(x.Modifier))%2 == 1 {
		body = inverse(body)
	}

	// U
	pow := ast.ModPow(x.Modifier)
	if len(pow) == 0 {
		if err := e.ctrlCall(x, body, env); err != nil {
			return fmt.Errorf("ctrl(%v): %v", x, err)
		}

		return nil
	}

	// pow(2) @ pow(-2) @ U equals to pow(0) @ U
	var p int
	for _, m := range pow {
		v, err := e.eval(m.Index.List.List[0], env)
		if err != nil {
			return fmt.Errorf("eval(%v): %v", m.Index.List.List[0], err)
		}

		p = p + int(v.(*object.Int).Value)
	}

	// pow(0) equals to Identity
	if p == 0 {
		return nil
	}

	// pow(-1) equals to Inv
	if p < 0 {
		p = -1 * p
		body = inverse(body)
	}

	// apply pow(2) @ U
	for i := 0; i < p; i++ {
		// U
		if err := e.ctrlCall(x, body, env); err != nil {
			return fmt.Errorf("ctrl(%v): %v", x, err)
		}
	}

	return nil
}

func (e *Evaluator) extg(x *ast.CallExpr, g *ast.GateDecl, outer *env.Environment) (*env.Environment, error) {
	env := env.NewEnclosed(outer)

	for i, p := range g.Params.List.List {
		n := ast.Must(ast.Ident(p))
		if v, ok := outer.Const[n]; ok {
			env.Const[n] = v
			continue
		}

		v, err := e.eval(x.Params.List.List[i], outer)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", x.Params.List.List[i], err)
		}

		env.Const[n] = v
	}

	// count ctrl/negctrl equals to target
	// ctrl @ negctrl @ x q0, q1, q2
	// q2 is target
	c := len(ast.ModCtrl(x.Modifier))

	// ctrl
	// loop with x.QArgs.List
	for i := 0; i < c; i++ {
		if v, ok := outer.Qubit.Get(x.QArgs.List[i]); ok {
			env.Qubit.Add(&ast.IdentExpr{Name: fmt.Sprintf("c_%v", x.QArgs.List[i])}, v)
			continue
		}

		return nil, fmt.Errorf("qubit(%v) not found", x.QArgs.List[i])
	}

	// target
	for i := range g.QArgs.List {
		if v, ok := outer.Qubit.Get(x.QArgs.List[c+i]); ok {
			env.Qubit.Add(g.QArgs.List[i], v)
			continue
		}

		return nil, fmt.Errorf("qubit(%v) not found", x.QArgs.List[c+i])
	}

	return env, nil
}

func (e *Evaluator) applyf(x *ast.CallExpr, g *ast.FuncDecl, env *env.Environment) (object.Object, error) {
	// extned
	ext, err := e.extf(x, g, env)
	if err != nil {
		return nil, fmt.Errorf("extend: %v", err)
	}

	// eval
	v, err := e.eval(&g.Body, ext)
	if err != nil {
		return nil, fmt.Errorf("eval(%v): %v", &g.Body, err)
	}

	return v.(*object.ReturnValue).Value, nil
}

func (e *Evaluator) extf(x *ast.CallExpr, g *ast.FuncDecl, outer *env.Environment) (*env.Environment, error) {
	env := env.NewEnclosed(outer)

	for i := range g.QArgs.List {
		if v, ok := outer.Qubit.Get(x.QArgs.List[i]); ok {
			env.Qubit.Add(g.QArgs.List[i], v)
			continue
		}

		return nil, fmt.Errorf("qubit(%v) not found", x.QArgs.List[i])
	}

	return env, nil
}

func (e *Evaluator) Println() error {
	if _, err := e.eval(&ast.PrintStmt{}, e.Env); err != nil {
		return fmt.Errorf("print qubit: %v", err)
	}

	for _, n := range e.Env.Bit.Name {
		fmt.Printf("%v: ", n)

		c, ok := e.Env.Bit.Get(&ast.IdentExpr{Name: n})
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

func override(block ast.BlockStmt, name []string) ast.BlockStmt {
	var qargs ast.ExprList
	for _, n := range name {
		qargs.Append(&ast.IdentExpr{Name: n})
	}

	var out ast.BlockStmt
	for _, b := range block.List {
		switch s := b.(type) {
		case *ast.ApplyStmt:
			out.Append(&ast.ApplyStmt{
				Modifier: s.Modifier,
				Kind:     s.Kind,
				Name:     s.Name,
				Params:   s.Params,
				QArgs:    qargs,
			})

		case *ast.ExprStmt:
			switch X := s.X.(type) {
			case *ast.CallExpr:
				out.Append(&ast.ExprStmt{
					X: &ast.CallExpr{
						Modifier: X.Modifier,
						Name:     X.Name,
						Params:   X.Params,
						QArgs:    qargs,
					},
				})
			}
		}
	}

	return out
}

func inverse(block ast.BlockStmt) ast.BlockStmt {
	// inv @ and reverse
	out := addmod(block, []ast.Modifier{{Kind: lexer.INV}})
	return out.Reverse()
}

func addmod(block ast.BlockStmt, mod []ast.Modifier) ast.BlockStmt {
	var out ast.BlockStmt
	for _, b := range block.List {
		switch s := b.(type) {
		case *ast.ApplyStmt:
			out.Append(&ast.ApplyStmt{
				Modifier: append(s.Modifier, mod...),
				Kind:     s.Kind,
				Name:     s.Name,
				Params:   s.Params,
				QArgs:    s.QArgs,
			})

		case *ast.ExprStmt:
			switch X := s.X.(type) {
			case *ast.CallExpr:
				out.Append(&ast.ExprStmt{
					X: &ast.CallExpr{
						Modifier: append(X.Modifier, mod...),
						Name:     X.Name,
						Params:   X.Params,
						QArgs:    X.QArgs,
					},
				})
			}
		}
	}

	return out
}

func builtin(g lexer.Token, p []float64) (matrix.Matrix, bool) {
	switch g {
	case lexer.U:
		return gate.U(p[0], p[1], p[2]), true
	case lexer.X:
		return gate.X(), true
	case lexer.Y:
		return gate.Y(), true
	case lexer.Z:
		return gate.Z(), true
	case lexer.H:
		return gate.H(), true
	case lexer.T:
		return gate.T(), true
	case lexer.S:
		return gate.S(), true
	}

	return nil, false
}

func flatten(qargs [][]q.Qubit) []q.Qubit {
	var out []q.Qubit
	for _, q := range qargs {
		out = append(out, q...)
	}

	return out
}
