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
		for _, b := range n.List {
			v, err := e.eval(b, env)
			if err != nil {
				return nil, fmt.Errorf("block: eval(%v): %v", b, err)
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

func (e *Evaluator) Print(s *ast.PrintStmt, env *env.Environment) error {
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

	return nil
}

func (e *Evaluator) Println() error {
	if _, err := e.eval(&ast.PrintStmt{}, e.Env); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	for _, n := range e.Env.Bit.Name {
		fmt.Printf("%v: ", n)

		c, ok := e.Env.Bit.Get(&ast.IdentExpr{Name: n})
		if !ok {
			return fmt.Errorf("bit=%v not found", n)
		}

		for _, v := range c {
			fmt.Printf("%v", v)
		}

		fmt.Println()
	}
	return nil
}

func (e *Evaluator) Reset(s *ast.ResetStmt, env *env.Environment) error {
	for _, a := range s.QArgs.List {
		qb, ok := env.Qubit.Get(a)
		if !ok {
			return fmt.Errorf("qubit=%v not found", a)
		}

		e.Q.Reset(qb...)
	}

	return nil
}

func (e *Evaluator) Measure(x *ast.MeasureExpr, env *env.Environment) (object.Object, error) {
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

func (e *Evaluator) Apply(s *ast.ApplyStmt, env *env.Environment) error {
	params, err := e.Params(s, env)
	if err != nil {
		return fmt.Errorf("params: %v", err)
	}

	qargs, err := e.QArgs(s, env)
	if err != nil {
		return fmt.Errorf("qargs: %v", err)
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

	// modifier
	u, ctrl, negc := e.Mod(s.Modifier, u, flatten(qargs), env)
	if len(ctrl)+len(negc) > 0 {
		// ctrl @ negctrl @ u
		e.Negc(negc, func() { e.Q.Apply(u) })
		return nil
	}

	e.Q.Apply(u, flatten(qargs)...)
	return nil
}

func (e *Evaluator) Params(s *ast.ApplyStmt, env *env.Environment) ([]float64, error) {
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

func (e *Evaluator) QArgs(s *ast.ApplyStmt, env *env.Environment) ([][]q.Qubit, error) {
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

func (e *Evaluator) Mod(mod []ast.Modifier, u matrix.Matrix, qargs []q.Qubit, env *env.Environment) (matrix.Matrix, []q.Qubit, []q.Qubit) {
	u = e.Inv(mod, u)
	u = e.Pow(mod, u, env)
	return e.Ctrl(mod, u, qargs, env)
}

func (e *Evaluator) Pow(mod []ast.Modifier, u matrix.Matrix, env *env.Environment) matrix.Matrix {
	pow := ast.ModPow(mod)
	if len(pow) == 0 {
		return u
	}

	// pow(2) @ pow(-2) is equal to pow(0)
	var p int
	for _, m := range pow {
		v := Must(e.eval(m.Index.List.List[0], env))
		p = p + int(v.(*object.Int).Value)
	}

	// pow(0) is equal to identity
	if p == 0 {
		return gate.I()
	}

	// pow(-1) is equals to inv
	if p < 0 {
		p = -1 * p
		u = u.Dagger()
	}

	// pow(p) @ g
	out := u
	for i := 1; i < p; i++ {
		out = out.Apply(u)
	}

	return out
}

func (e *Evaluator) Inv(mod []ast.Modifier, u matrix.Matrix) matrix.Matrix {
	// inv @ U
	if len(ast.ModInv(mod))%2 == 1 {
		u = u.Dagger()
	}

	return u
}

func (e *Evaluator) Ctrl(mod []ast.Modifier, u matrix.Matrix, qargs []q.Qubit, env *env.Environment) (matrix.Matrix, []q.Qubit, []q.Qubit) {
	var ctrl, negc []q.Qubit
	if len(ast.ModCtrl(mod)) == 0 {
		return u, ctrl, negc
	}

	begin := 0
	for _, m := range ast.ModCtrl(mod) {
		p := 1
		if len(m.Index.List.List) > 0 {
			v := Must(e.eval(m.Index.List.List[0], env))
			p = int(v.(*object.Int).Value)
		}

		switch m.Kind {
		case lexer.CTRL:
			ctrl = append(ctrl, qargs[begin:begin+p]...)
		case lexer.NEGCTRL:
			negc = append(negc, qargs[begin:begin+p]...)
		}

		begin = begin + p
	}

	return gate.Controlled(u, len(qargs), q.Index(append(ctrl, negc...)...), qargs[len(qargs)-1].Index()), ctrl, negc
}

func (e *Evaluator) Negc(negc []q.Qubit, f func()) {
	if len(negc) > 0 {
		e.Q.X(negc...)
	}

	f()

	if len(negc) > 0 {
		e.Q.X(negc...)
	}
}

func (e *Evaluator) Call(s *ast.CallExpr, env *env.Environment) (object.Object, error) {
	fmt.Println(s)
	return nil, nil
}
