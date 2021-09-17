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
	"github.com/itsubaki/qasm/pkg/evaluator/object"
	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
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

func New(qsim *q.Q, env *object.Environment, opts ...Opts) *Evaluator {
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
	return New(q.New(), object.NewEnvironment(), opts...)
}

func Eval(n ast.Node) (object.Object, error) {
	return Default().eval(n, object.NewEnvironment())
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

func (e *Evaluator) eval(n ast.Node, env *object.Environment) (obj object.Object, err error) {
	defer func() {
		if !e.Opts.Verbose {
			return
		}

		if obj != nil && obj.Type() != object.NIL {
			fmt.Printf("%v", strings.Repeat(indent, e.indent+1))
			fmt.Printf("return %T(%v)\n", obj, obj)
		}

		e.indent--
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

	case *ast.CallExpr:
		return e.call(n, env)

	case *ast.MeasureExpr:
		return e.measure(n, env)

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
		if v, ok := env.Const[ast.Ident(n)]; ok {
			return v, nil
		}
	}

	return &object.Nil{}, nil
}

func (e *Evaluator) evalInfix(kind lexer.Token, lhs, rhs object.Object) (object.Object, error) {
	switch kind {
	case lexer.PLUS:
		switch t := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: t.Value + rhs.(*object.Int).Value}, nil
		case *object.Float:
			return &object.Float{Value: t.Value + rhs.(*object.Float).Value}, nil
		}

	case lexer.MINUS:
		switch t := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: t.Value - rhs.(*object.Int).Value}, nil
		case *object.Float:
			return &object.Float{Value: t.Value - rhs.(*object.Float).Value}, nil
		}

	case lexer.MUL:
		switch t := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: t.Value * rhs.(*object.Int).Value}, nil
		case *object.Float:
			return &object.Float{Value: t.Value * rhs.(*object.Float).Value}, nil
		}

	case lexer.DIV:
		switch t := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: t.Value / rhs.(*object.Int).Value}, nil
		case *object.Float:
			return &object.Float{Value: t.Value / rhs.(*object.Float).Value}, nil
		}

	case lexer.MOD:
		switch t := lhs.(type) {
		case *object.Int:
			return &object.Int{Value: t.Value % rhs.(*object.Int).Value}, nil
		}

	}

	return nil, fmt.Errorf("unsupported(%v)", kind)
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
	params := make([]float64, 0)
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
			return fmt.Errorf("unsupported(%v)", o)
		}
	}

	qargs := make([][]q.Qubit, 0)
	for _, a := range s.QArgs.List {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit(%v) not found", a)
		}

		qargs = append(qargs, qb)
	}

	return e.apply(s.Modifier, s.Kind, params, qargs)
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

func inv(mod []ast.Modifier, u matrix.Matrix) matrix.Matrix {
	var c int
	for _, m := range mod {
		if m.Kind == lexer.INV {
			c++
		}
	}

	if c%2 == 1 {
		u = u.Dagger()
	}

	return u
}

func pow(mod []ast.Modifier, u matrix.Matrix) matrix.Matrix {
	for _, m := range mod {
		if m.Kind != lexer.POW {
			continue
		}

		var c int
		if len(m.Index.List.List) > 0 {
			c = int(m.Index.List.List[0].(*ast.BasicLit).Int64())
		}

		// pow(-1) U equals to inv @ U
		if c < 0 {
			u = u.Dagger()
			c = -1 * c
		}

		tmp := u
		for i := 1; i < c; i++ {
			u = u.Apply(tmp)
		}
	}

	return u
}

func flatten(qargs [][]q.Qubit) []q.Qubit {
	out := make([]q.Qubit, 0)
	for _, q := range qargs {
		out = append(out, q...)
	}

	return out
}

func (e *Evaluator) tryBuiltinApply(g lexer.Token, p []float64, qargs [][]q.Qubit) bool {
	switch g {
	case lexer.SWAP:
		e.Q.Swap(flatten(qargs)...)
		return true
	case lexer.QFT:
		e.Q.QFT(flatten(qargs)...)
		return true
	case lexer.IQFT:
		e.Q.InvQFT(flatten(qargs)...)
		return true
	case lexer.CMODEXP2:
		e.Q.CModExp2(int(p[0]), int(p[1]), qargs[0], qargs[1])
		return true
	case lexer.CX:
		for i := range qargs[0] {
			e.Q.CNOT(qargs[0][i], qargs[1][i])
		}
		return true
	case lexer.CZ:
		for i := range qargs[0] {
			e.Q.CZ(qargs[0][i], qargs[1][i])
		}
		return true
	case lexer.CCX:
		for i := range qargs[0] {
			e.Q.CCNOT(qargs[0][i], qargs[1][i], qargs[2][i])
		}
		return true
	}

	return false
}

func (e *Evaluator) tryCtrlApply(mod []ast.Modifier, u matrix.Matrix, qargs [][]q.Qubit) bool {
	ctrl := make([][]q.Qubit, 0)
	negc := make([][]q.Qubit, 0)

	// ctrl @ ctrl @ X equals to ctrl(0) @ ctrl(1) @ X
	defaultIndex := 0
	for _, m := range mod {
		c := defaultIndex
		if len(m.Index.List.List) > 0 {
			c = int(m.Index.List.List[0].(*ast.BasicLit).Int64())
		}

		switch m.Kind {
		case lexer.CTRL:
			ctrl = append(ctrl, qargs[c])
		case lexer.NEGCTRL:
			negc = append(negc, qargs[c])
		default:
			continue
		}

		defaultIndex++
	}

	if len(ctrl) == 0 && len(negc) == 0 {
		return false
	}

	if len(negc) > 0 {
		e.Q.X(flatten(negc)...)
	}

	c := append(flatten(ctrl), flatten(negc)...)
	e.Q.Controlled(u, c, qargs[len(qargs)-1][0])

	if len(negc) > 0 {
		e.Q.X(flatten(negc)...)
	}

	return true
}

func (e *Evaluator) apply(mod []ast.Modifier, gate lexer.Token, params []float64, qargs [][]q.Qubit) error {
	// cx, swap, qft, cmodexp2, ...
	if e.tryBuiltinApply(gate, params, qargs) {
		return nil
	}

	// U, x, y, z, h, ...
	u, ok := builtin(gate, params)
	if !ok {
		return fmt.Errorf("gate=%v(%v) not found", lexer.Tokens[gate], gate)
	}

	// Inverse U
	u = inv(mod, u)

	// Pow(2) U
	u = pow(mod, u)

	// Controlled-U
	if e.tryCtrlApply(mod, u, qargs) {
		return nil
	}

	// U
	e.Q.Apply(u, flatten(qargs)...)
	return nil
}

func (e *Evaluator) call(x *ast.CallExpr, env *object.Environment) (object.Object, error) {
	decl, ok := env.Func[x.Name]
	if !ok {
		return nil, fmt.Errorf("decl(%v) not found", x.Name)
	}

	if e.Opts.Verbose {
		fmt.Printf("%v", strings.Repeat(indent, e.indent))
		fmt.Printf("%T(%v)\n", decl, decl)
	}

	switch decl := decl.(type) {
	case *ast.GateDecl:
		return nil, e.callGate(x, decl, env)
	case *ast.FuncDecl:
		return e.callFunc(x, decl, env)
	}

	return nil, fmt.Errorf("unsupported(%v)", decl)
}

func isCtrl(mod []ast.Modifier) bool {
	for _, m := range mod {
		if m.Kind == lexer.CTRL || m.Kind == lexer.NEGCTRL {
			return true
		}
	}

	return false
}

func (e *Evaluator) callGate(x *ast.CallExpr, d *ast.GateDecl, outer *object.Environment) error {
	if !isCtrl(x.Modifier) {
		// U(pi, 0, pi) q;
		if _, err := e.eval(&d.Body, e.extend(x, d, outer)); err != nil {
			return fmt.Errorf("eval(%v): %v", d.Body, err)
		}

		return nil
	}

	// ctrl @
	for _, b := range d.Body.List {
		var c ast.Node

		switch s := b.(type) {
		case *ast.ApplyStmt:
			// ctrl @ U(pi, 0, pi) q, p;
			c = &ast.ApplyStmt{
				Modifier: append(x.Modifier, s.Modifier...),
				Kind:     s.Kind,
				Name:     s.Name,
				Params:   s.Params,
				QArgs:    x.QArgs,
			}

		case *ast.ExprStmt:
			switch X := s.X.(type) {
			case *ast.CallExpr:
				// ctrl @ h q, p;
				c = &ast.CallExpr{
					Modifier: append(x.Modifier, X.Modifier...),
					Name:     X.Name,
					Params:   X.Params,
					QArgs:    x.QArgs,
				}
			default:
				return fmt.Errorf("unsupported(%v)", X)
			}

		default:
			return fmt.Errorf("unsupported(%v)", s)
		}

		if _, err := e.eval(c, outer); err != nil {
			return fmt.Errorf("eval(%v): %v", c, err)
		}
	}

	return nil
}

func (e *Evaluator) extend(x *ast.CallExpr, d *ast.GateDecl, outer *object.Environment) *object.Environment {
	env := object.NewEnclosedEnvironment(outer)

	for i := range d.Params.List.List {
		if v, ok := outer.Const[ast.Ident(x.Params.List.List[i])]; ok {
			env.Const[ast.Ident(d.Params.List.List[i])] = v
			continue
		}

		v, err := e.eval(x.Params.List.List[i], outer)
		if err != nil {
			panic(fmt.Sprintf("eval(%v): %v", x.Params.List.List[i], err))
		}

		env.Const[ast.Ident(d.Params.List.List[i])] = v
	}

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
