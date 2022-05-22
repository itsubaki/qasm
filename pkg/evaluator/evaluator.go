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
	Env    *env.Environ
	Opts   Opts
	indent int
}

type Opts struct {
	Verbose bool
}

func New(qsim *q.Q, env *env.Environ, opts ...Opts) *Evaluator {
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

func (e *Evaluator) eval(n ast.Node, env *env.Environ) (obj object.Object, err error) {
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
		if err := e.Reset(n, env); err != nil {
			return nil, fmt.Errorf("reset(%v): %v", n, err)
		}

	case *ast.ApplyStmt:
		if err := e.Apply(n, env); err != nil {
			return nil, fmt.Errorf("apply(%v): %v", n, err)
		}

	case *ast.PrintStmt:
		if err := e.Print(n, env); err != nil {
			return nil, fmt.Errorf("print(%v): %v", n, err)
		}

	case *ast.InclStmt:
		path := strings.Trim(n.Path.Value, "\"")
		f, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("include: read file=%s: %v", path, err)
		}

		l := lexer.New(strings.NewReader(string(f)))
		p := parser.New(l)

		a := p.Parse()
		if errs := p.Errors(); len(errs) != 0 {
			return nil, fmt.Errorf("include: parse: %v", errs)
		}

		for _, s := range a.Stmts {
			if _, err := e.eval(s, e.Env); err != nil {
				return nil, fmt.Errorf("include: eval(stmt=%v): %v", s, err)
			}
		}

	case *ast.AssignStmt:
		rhs, err := e.eval(n.Right, env)
		if err != nil {
			return nil, fmt.Errorf("assign: eval(right=%v): %v", n.Right, err)
		}

		c, ok := env.Bit.Get(n.Left)
		if !ok {
			return nil, fmt.Errorf("assign: bit=%v not found", n.Left)
		}

		e := rhs.(*object.Array).Elm
		for i := range e {
			c[i] = e[i].(*object.Int).Value
		}

	case *ast.BlockStmt:
		// NOTE: Modifier
		//
		// gate inv(a, b, c) q { inv @ U(a, b, c) q; inv @ U(c, b, a) q;}
		// pow(2) @ inv(pi/2.0, 0, pi) q;
		//
		// is
		// inv @ U(a, b, c) q;
		// inv @ U(c, b, a) q;
		// inv @ U(a, b, c) q;
		// inv @ U(c, b, a) q;
		//
		// is not
		// inv @ U(a, b, c) q;
		// inv @ U(a, b, c) q;
		// inv @ U(c, b, a) q;
		// inv @ U(c, b, a) q;
		//
		//
		// gate pow23(a, b, c) q { pow(2) @ U(a, b, c) q; pow(3) @ U(c, b, a) q;}
		// inv @ pow23(pi/2.0, 0, pi) q;
		//
		// is
		// inv @ pow(3) @ U(c, b, a) q;
		// inv @ pow(2) @ U(a, b, c) q;
		//
		// is not
		// inv @ pow(2) @ U(a, b, c) q;
		// inv @ pow(3) @ U(c, b, a) q;

		for _, m := range env.Modifier {
			switch m.Kind {
			case lexer.INV:
				n = Inverse(n, env)
			case lexer.POW:
				n.Pow(int(ast.Must(e.eval(m.Index.List.List[0], env)).(*object.Int).Value))
			}
		}

		return e.ApplyBolck(n, env)

	case *ast.ReturnStmt:
		v, err := e.eval(n.Result, env)
		if err != nil {
			return nil, fmt.Errorf("return: %v", err)
		}

		return &object.ReturnValue{Value: v}, nil

	case *ast.GenConst:
		v, err := e.eval(n.Value, env)
		if err != nil {
			return nil, fmt.Errorf("gen const: eval(%v): %v", n, err)
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
		return e.Call(n, env)

	case *ast.MeasureExpr:
		return e.Measure(n, env)

	case *ast.UnaryExpr:
		v, err := e.eval(n.Value, env)
		if err != nil {
			return nil, fmt.Errorf("unary: eval(%v): %v", n.Value, err)
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
			return nil, fmt.Errorf("infix: eval(left=%v): %v", n.Left, err)
		}

		rhs, err := e.eval(n.Right, env)
		if err != nil {
			return nil, fmt.Errorf("infix: eval(left=%v): %v", n.Right, err)
		}

		return e.Infix(n.Kind, lhs, rhs)

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

func (e *Evaluator) Infix(kind lexer.Token, lhs, rhs object.Object) (object.Object, error) {
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

func (e *Evaluator) Print(s *ast.PrintStmt, env *env.Environ) error {
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
			return fmt.Errorf("qubit=%v not found", a)
		}

		index = append(index, q.Index(qb...))
	}

	for _, s := range e.Q.Raw().State(index...) {
		fmt.Println(s)
	}

	for _, n := range e.Env.Bit.Name {
		fmt.Printf("%v: ", n)

		c, ok := e.Env.Bit.Get(&ast.IdentExpr{Name: n})
		if !ok {
			return fmt.Errorf("bit=%v not found", n)
		}

		fmt.Println(c)
	}

	return nil
}

func (e *Evaluator) Println() error {
	if _, err := e.eval(&ast.PrintStmt{}, e.Env); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func (e *Evaluator) Reset(s *ast.ResetStmt, env *env.Environ) error {
	for _, a := range s.QArgs.List {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit=%v not found", a)
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) Measure(x *ast.MeasureExpr, env *env.Environ) (object.Object, error) {
	qargs := x.QArgs.List
	if len(qargs) == 0 {
		return nil, fmt.Errorf("qargs is empty")
	}

	var m []q.Qubit
	for _, a := range qargs {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return nil, fmt.Errorf("qubit=%v not found", a)
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

func (e *Evaluator) ApplyBolck(s *ast.BlockStmt, env *env.Environ) (object.Object, error) {
	for _, b := range s.List {
		v, err := e.eval(b, env)
		if err != nil {
			return nil, fmt.Errorf("block=%v: %v", b, err)
		}

		if v != nil && v.Type() == object.RETURN_VALUE {
			return v, nil
		}
	}

	return nil, nil
}

func (e *Evaluator) Apply(s *ast.ApplyStmt, env *env.Environ) error {
	params, err := e.Params(s, env)
	if err != nil {
		return fmt.Errorf("params: %v", err)
	}

	qargs, err := e.QArgs(s, env)
	if err != nil {
		return fmt.Errorf("qargs: %v", err)
	}

	if len(env.Decl) > 0 && len(qargs) > 1 {
		// for j ← 0, 1 do
		//   g qr0[0],qr1[j],qr2[0],qr3[j];
		return e.ApplyInBlock(params, qargs, s, env)
	}

	// QFT, IQFT, CMODEXP2
	if BuiltinApply(e.Q, s.Kind, params, qargs) {
		return nil
	}

	// U, X, Y, Z, H, T, S
	u, ok := Builtin(s.Kind, params)
	if !ok {
		return fmt.Errorf("gate=%v not found", lexer.Tokens[s.Kind])
	}

	// Modifier
	// FIXME: ctrl and inv, pow cannot be specified at the same time.
	if s.HasModCtrl() && (s.HasModInv() || s.HasModPow()) {
		return fmt.Errorf("invalid modifier=%v", s.Modifier)
	}

	u, ctrl, negc := e.Ctrl(s.Modifier, u, qargs, env)
	if len(ctrl)+len(negc) > 0 {
		// ctrl @ negctrl @ u
		e.X(negc, func() { e.Q.Apply(u) })
		return nil
	}

	// inv, pow or nothing.
	u = e.Mod(s.Modifier, u, env)
	e.Q.Apply(u, flatten(qargs)...)
	return nil
}

func (e *Evaluator) ApplyInBlock(params []float64, qargs [][]q.Qubit, s *ast.ApplyStmt, env *env.Environ) error {
	fmt.Printf("here is in decl=%v, qargs=%v\n", env.Decl, qargs)
	return nil
}

func (e *Evaluator) Params(s *ast.ApplyStmt, env *env.Environ) ([]float64, error) {
	var params []float64
	for _, p := range s.Params.List.List {
		v, err := e.eval(p, env)
		if err != nil {
			return nil, fmt.Errorf("params: eval(%v): %v", p, err)
		}

		switch o := v.(type) {
		case *object.Float:
			params = append(params, o.Value)
		case *object.Int:
			params = append(params, float64(o.Value))
		default:
			return nil, fmt.Errorf("unsupported(%v)", o)
		}
	}

	return params, nil
}

func (e *Evaluator) QArgs(s *ast.ApplyStmt, env *env.Environ) ([][]q.Qubit, error) {
	var qargs [][]q.Qubit
	for _, a := range s.QArgs.List {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return nil, fmt.Errorf("qubit=%v not found", a)
		}

		qargs = append(qargs, qb)
	}

	return qargs, nil
}

func (e *Evaluator) Mod(mod []ast.Modifier, u matrix.Matrix, env *env.Environ) matrix.Matrix {
	// NOTE: pow(2) @ inv @ u is not equal to inv @ pow(2) @ u
	for _, m := range mod {
		switch m.Kind {
		case lexer.INV:
			u = u.Dagger()
		case lexer.POW:
			p := ast.Must(e.eval(m.Index.List.List[0], env)).(*object.Int).Value
			u = matrix.ApplyN(u, int(p))
		}
	}

	return u
}

func (e *Evaluator) Ctrl(mod []ast.Modifier, u matrix.Matrix, qargs [][]q.Qubit, env *env.Environ) (matrix.Matrix, []q.Qubit, []q.Qubit) {
	var ctrl, negc []q.Qubit
	if len(ast.ModCtrl(mod)) == 0 {
		return u, ctrl, negc
	}

	fqargs, begin := flatten(qargs), 0
	for _, m := range ast.ModCtrl(mod) {
		p := 1
		if len(m.Index.List.List) > 0 {
			v := ast.Must(e.eval(m.Index.List.List[0], env))
			p = int(v.(*object.Int).Value)
		}

		switch m.Kind {
		case lexer.CTRL:
			ctrl = append(ctrl, fqargs[begin:begin+p]...)
		case lexer.NEGCTRL:
			negc = append(negc, fqargs[begin:begin+p]...)
		}

		begin = begin + p
	}

	n := e.Q.NumberOfBit()
	c := q.Index(append(ctrl, negc...)...)
	t := q.Index(qargs[len(qargs)-1]...)

	// FIXME: fixed target.
	return gate.Controlled(u, n, c, t[0]), ctrl, negc
}

func (e *Evaluator) X(target []q.Qubit, f func()) {
	if len(target) > 0 {
		e.Q.X(target...)
	}

	f()

	if len(target) > 0 {
		e.Q.X(target...)
	}
}

func (e *Evaluator) Call(x *ast.CallExpr, outer *env.Environ) (object.Object, error) {
	f, ok := outer.Func[x.Name]
	if !ok {
		return nil, fmt.Errorf("decl=%v not found", x.Name)
	}

	if e.Opts.Verbose {
		fmt.Printf("%v", strings.Repeat(indent, e.indent))
		fmt.Printf("%T(%v)\n", f, f)
	}

	switch decl := f.(type) {
	case *ast.GateDecl:
		return e.eval(&decl.Body, e.Enclosed(x, decl, outer))
	}

	return nil, fmt.Errorf("unsupported decl=%v", f)
}

func (e *Evaluator) Enclosed(x *ast.CallExpr, decl *ast.GateDecl, outer *env.Environ) *env.Environ {
	env := outer.NewEnclosed(decl, x.Modifier)
	e.SetConst(env, outer, decl.Params.List.List, x.Params.List.List)
	e.SetQArgs(env, outer, decl.QArgs.List, x.QArgs.List)
	return env
}

func (e *Evaluator) SetConst(env, outer *env.Environ, decl, args []ast.Expr) {
	for i, d := range decl {
		n := ast.Must(ast.Ident(d))
		v := ast.Must(e.eval(args[i], outer))
		env.Const[n] = v
	}
}

func (e *Evaluator) SetQArgs(env, outer *env.Environ, decl, args []ast.Expr) {
	for i := range decl {
		if qb, ok := outer.Qubit.Get(args[i]); ok {
			env.Qubit.Add(decl[i], qb)
		}
	}
}

func Inverse(n *ast.BlockStmt, env *env.Environ) *ast.BlockStmt {
	env.Modifier = append(env.Modifier, ast.Modifier{Kind: lexer.INV})
	return n.Reverse()
}
