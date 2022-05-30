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

	case *ast.CallExpr:
		return e.Call(n, env)

	case *ast.BlockStmt:
		return e.Block(e.ModifyStmt(n, env), env)

	case *ast.ReturnStmt:
		return e.Return(n, env)

	case *ast.MeasureExpr:
		return e.Measure(n, env)

	case *ast.UnaryExpr:
		return e.Unary(n, env)

	case *ast.InfixExpr:
		return e.Infix(n, env)

	case *ast.InclStmt:
		return &object.Nil{}, e.Include(n, env)

	case *ast.AssignStmt:
		return &object.Nil{}, e.Assign(n, env)

	case *ast.ResetStmt:
		return &object.Nil{}, e.Reset(n, env)

	case *ast.PrintStmt:
		return &object.Nil{}, e.Print(n, env)

	case *ast.ApplyStmt:
		return &object.Nil{}, e.Apply(n, env)

	case *ast.GenConst:
		return &object.Nil{}, e.GenConst(n, env)

	case *ast.GenDecl:
		return &object.Nil{}, e.GenDecl(n, env)

	case *ast.GateDecl:
		env.Func[ast.Must(ast.Ident(n))] = n
		return &object.Nil{}, nil

	case *ast.FuncDecl:
		env.Func[ast.Must(ast.Ident(n))] = n
		return &object.Nil{}, nil

	case *ast.IdentExpr:
		if v, ok := env.Const[ast.Must(ast.Ident(n))]; ok {
			return v, nil
		}

		return nil, fmt.Errorf("const=%v not found", n)

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

		return nil, fmt.Errorf("unsupported type=%v", n)
	}

	return nil, fmt.Errorf("unsupported type=%v", n)
}

func (e *Evaluator) Include(s *ast.InclStmt, env *env.Environ) error {
	path := strings.Trim(s.Path.Value, "\"")
	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file=%s: %v", path, err)
	}

	l := lexer.New(strings.NewReader(string(f)))
	p := parser.New(l)

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		return fmt.Errorf("parse: %v", errs)
	}

	for _, s := range a.Stmts {
		if _, err := e.eval(s, e.Env); err != nil {
			return fmt.Errorf("eval(%v): %v", s, err)
		}
	}

	return nil
}

func (e *Evaluator) Assign(s *ast.AssignStmt, env *env.Environ) error {
	rhs, err := e.eval(s.Right, env)
	if err != nil {
		return fmt.Errorf("eval(%v): %v", s.Right, err)
	}

	c, ok := env.Bit.Get(s.Left)
	if !ok {
		return fmt.Errorf("bit=%v not found", s.Left)
	}

	elm := rhs.(*object.Array).Elm
	for i := range elm {
		c[i] = elm[i].(*object.Int).Value
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
		c, _ := e.Env.Bit.Get(&ast.IdentExpr{Name: n})
		fmt.Println(c)
	}

	return nil
}

func (e *Evaluator) Println() error {
	if _, err := e.eval(&ast.PrintStmt{}, e.Env); err != nil {
		return fmt.Errorf("eval(%v): %v", &ast.PrintStmt{}, err)
	}

	return nil
}

func (e *Evaluator) GenConst(s *ast.GenConst, env *env.Environ) error {
	v, err := e.eval(s.Value, env)
	if err != nil {
		return fmt.Errorf("eval(%v): %v", s.Value, err)
	}

	env.Const[ast.Must(ast.Ident(s))] = v
	return nil
}

func (e *Evaluator) GenDecl(s *ast.GenDecl, env *env.Environ) error {
	switch s.Kind {
	case lexer.BIT:
		env.Bit.Add(s, make([]int64, s.Size()))
	case lexer.QUBIT:
		env.Qubit.Add(s, e.Q.ZeroWith(s.Size()))
	}

	return nil
}

func (e *Evaluator) Block(s *ast.BlockStmt, env *env.Environ) (object.Object, error) {
	for _, b := range s.List {
		v, err := e.eval(b, env)
		if err != nil {
			return nil, fmt.Errorf("eval(%v): %v", b, err)
		}

		if v != nil && v.Type() == object.RETURN_VALUE {
			return v, nil
		}
	}

	return nil, nil
}

func (e *Evaluator) Return(s *ast.ReturnStmt, env *env.Environ) (object.Object, error) {
	v, err := e.eval(s.Result, env)
	if err != nil {
		return nil, fmt.Errorf("eval(%v): %v", s.Result, err)
	}

	return &object.ReturnValue{Value: v}, nil
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

func (e *Evaluator) Unary(s *ast.UnaryExpr, env *env.Environ) (object.Object, error) {
	v, err := e.eval(s.Value, env)
	if err != nil {
		return nil, fmt.Errorf("eval(%v): %v", s.Value, err)
	}

	switch s.Kind {
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

	return nil, fmt.Errorf("unsupported type=%v", s.Kind)
}

func (e *Evaluator) Infix(s *ast.InfixExpr, env *env.Environ) (object.Object, error) {
	lhs, err := e.eval(s.Left, env)
	if err != nil {
		return nil, fmt.Errorf("eval(%v): %v", s.Left, err)
	}

	rhs, err := e.eval(s.Right, env)
	if err != nil {
		return nil, fmt.Errorf("eval(%v): %v", s.Right, err)
	}

	switch s.Kind {
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

	return nil, fmt.Errorf("unsupported type=%v", s.Kind)
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

	// QFT, IQFT, CMODEXP2
	if e.ApplyBuiltin(s.Kind, params, qargs) {
		return nil
	}

	// U, X, Y, Z, H, T, S
	u, ok := Builtin(s.Kind, params)
	if !ok {
		return fmt.Errorf("gate=%v not found", lexer.Tokens[s.Kind])
	}

	if len(env.Decl) > 0 && len(qargs) > 1 {
		// for j ← 0, 1 do
		//   g qr0[0],qr1[j],qr2[0],qr3[j];
		// https://qiskit.github.io/openqasm/language/gates.html#hierarchically-defined-unitary-gates
		e.ApplyUParallel(s.Modifier, u, qargs, env)
		return nil
	}

	e.ApplyU(s.Modifier, u, qargs, env)
	return nil
}

func (e *Evaluator) ApplyBuiltin(g lexer.Token, p []float64, qargs [][]q.Qubit) bool {
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

func (e *Evaluator) ApplyUParallel(mod []ast.Modifier, u matrix.Matrix, qargs [][]q.Qubit, env *env.Environ) error {
	size := 0
	for i := range qargs {
		if len(qargs[i]) > size {
			size = len(qargs[i])
		}
	}

	// validation
	for i := range qargs {
		if len(qargs[i]) == 1 || len(qargs[i]) == size {
			continue
		}

		return fmt.Errorf("invalid qargs size=%v", qargs)
	}

	// gates can be applied in parallel
	// REVIEW: Should not need to be reversed when inv @ ctrl @ u qr0, qr1;
	for i := 0; i < size; i++ {
		e.ApplyUAt(i, mod, u, qargs, env)
	}

	return nil
}

func (e *Evaluator) ApplyUAt(i int, mod []ast.Modifier, u matrix.Matrix, qargs [][]q.Qubit, env *env.Environ) {
	cqargs := make([][]q.Qubit, 0)
	for j := range qargs {
		if len(qargs[j]) == 1 {
			cqargs = append(cqargs, []q.Qubit{qargs[j][0]})
			continue
		}

		cqargs = append(cqargs, []q.Qubit{qargs[j][i]})
	}

	e.ApplyU(mod, u, cqargs, env)
}

func (e *Evaluator) ApplyU(mod []ast.Modifier, u matrix.Matrix, qargs [][]q.Qubit, env *env.Environ) error {
	// Modify inv, pow
	u = e.ModifyU(mod, u, env)

	// Modify ctrl
	ctrl := ast.ModCtrl(mod)
	if len(ctrl) > 0 {
		// NOTE: That is, inv @ ctrl @ U = ctrl @ inv @ U.
		// https://qiskit.github.io/openqasm/language/gates.html#quantum-gate-modifiers

		u, _, negc := e.Ctrl(ctrl, u, qargs, env)
		e.X(negc, func() { e.Q.Apply(u) })
		return nil
	}

	// no ctrl
	e.Q.Apply(u, flatten(qargs)...)
	return nil
}

// Ctrl returns ctrl @ U
func (e *Evaluator) Ctrl(modctrl []ast.Modifier, u matrix.Matrix, qargs [][]q.Qubit, env *env.Environ) (matrix.Matrix, []q.Qubit, []q.Qubit) {
	var ctrl, negc []q.Qubit
	if len(modctrl) == 0 {
		return u, ctrl, negc
	}

	fqargs, begin := flatten(qargs), 0
	for _, m := range modctrl {
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

func (e *Evaluator) Params(s *ast.ApplyStmt, env *env.Environ) ([]float64, error) {
	var params []float64
	for _, p := range s.Params.List.List {
		v, err := e.eval(p, env)
		if err != nil {
			return nil, fmt.Errorf("params=%v: %v", p, err)
		}

		switch o := v.(type) {
		case *object.Float:
			params = append(params, o.Value)
		case *object.Int:
			params = append(params, float64(o.Value))
		default:
			return nil, fmt.Errorf("unsupported type=%v", o)
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

	return nil, fmt.Errorf("unsupported type=%v", f)
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
	if len(ast.ModCtrl(env.Modifier)) == 0 {
		for i := range decl {
			if qb, ok := outer.Qubit.Get(args[i]); ok {
				env.Qubit.Add(decl[i], qb)
			}
		}

		return
	}

	ctrl := make([]ast.Expr, 0)
	for i := range ast.ModCtrl(env.Modifier) {
		switch x := args[i].(type) {
		case *ast.IdentExpr:
			ctrl = append(ctrl, &ast.IdentExpr{Name: fmt.Sprintf("_v%d", i)})
		case *ast.IndexExpr:
			ctrl = append(ctrl, &ast.IndexExpr{Name: fmt.Sprintf("_v%d", i), Value: x.Value})
		}
	}
	env.CtrlQArgs = ctrl

	cdecl := append(make([]ast.Expr, 0), ctrl...)
	for i := range decl {
		cdecl = append(cdecl, decl[i])
	}

	for i := range cdecl {
		if qb, ok := outer.Qubit.Get(args[i]); ok {
			env.Qubit.Add(cdecl[i], qb)
			continue
		}
	}
}

// ModifyStmt returns modified BlockStmt
func (e *Evaluator) ModifyStmt(s *ast.BlockStmt, env *env.Environ) *ast.BlockStmt {
	// NOTE: pow(2) @ inv @ u is not equal to inv @ pow(2) @ u
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

	var i int
	for _, m := range env.Modifier {
		switch m.Kind {
		case lexer.CTRL, lexer.NEGCTRL:
			s = s.Add(m, []ast.Expr{env.CtrlQArgs[i]}...)
			i++
		case lexer.INV:
			s = s.Inv()
		case lexer.POW:
			s = s.Pow(int(ast.Must(e.eval(m.Index.List.List[0], env)).(*object.Int).Value))
		}
	}

	return s
}

// ModifyU returns modified(inv, pow) U
func (e *Evaluator) ModifyU(mod []ast.Modifier, u matrix.Matrix, env *env.Environ) matrix.Matrix {
	// NOTE: pow(2) @ inv @ u is not equal to inv @ pow(2) @ u
	for _, m := range mod {
		switch m.Kind {
		case lexer.INV:
			u = u.Dagger()
		case lexer.POW:
			u = matrix.ApplyN(u, int(ast.Must(e.eval(m.Index.List.List[0], env)).(*object.Int).Value))
		}
	}

	return u
}
