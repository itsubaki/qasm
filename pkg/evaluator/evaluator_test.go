package evaluator_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/pkg/quantum/qubit"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/evaluator/object"
	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
)

type state struct {
	Name  string
	Value []int64
}

func (s state) Equals(v state) bool {
	if s.Name != v.Name {
		return false
	}

	if len(s.Value) != len(v.Value) {
		return false
	}

	for i := range s.Value {
		if s.Value[i] != v.Value[i] {
			return false
		}
	}

	return true
}

func equals(cs, cv []state) bool {
	if len(cs) != len(cv) {
		return false
	}
	for i := range cs {
		if !cs[i].Equals(cv[i]) {
			return false
		}
	}

	return true
}

func eval(qasm string, verbose ...bool) ([]qubit.State, []state, error) {
	l := lexer.New(strings.NewReader(qasm))
	p := parser.New(l)

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		return nil, nil, fmt.Errorf("parse: %v\n", errs)
	}

	opts := evaluator.Opts{
		Verbose: false,
	}
	if verbose != nil {
		opts.Verbose = verbose[0]
	}

	e := evaluator.Default(opts)
	if err := e.Eval(a); err != nil {
		return nil, nil, fmt.Errorf("eval: %v\n", err)
	}

	s := make([]state, 0)
	for _, n := range e.Env.Bit.Name {
		s = append(s, state{
			Name:  n,
			Value: e.Env.Bit.Value[n],
		})
	}

	if len(e.Env.Qubit.Name) == 0 {
		return []qubit.State{}, s, nil
	}

	var index [][]int
	for _, n := range e.Env.Qubit.Name {
		qb, _ := e.Env.Qubit.Get(&ast.IdentExpr{Name: n})
		index = append(index, q.Index(qb...))
	}

	return e.Q.Raw().State(index...), s, nil
}

func Example_print() {
	qasm := `
OPENQASM 3.0;

print;

qubit[2] q;
U(pi, 0, pi) q[0];

print q[0], q[1];
print q;

bit[2] c;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
	// [10][  2]( 1.0000 0.0000i): 1.0000
	// [10][  2]( 1.0000 0.0000i): 1.0000
	// c: [0 0]
}

func Example_qFT() {
	qasm := `
OPENQASM 3.0;

qubit[3] q;

U(pi, 0, pi) q[1];
QFT q;

print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [000][  0]( 0.3536 0.0000i): 0.1250
	// [001][  1]( 0.3536 0.0000i): 0.1250
	// [010][  2](-0.3536 0.0000i): 0.1250
	// [011][  3](-0.3536 0.0000i): 0.1250
	// [100][  4]( 0.0000 0.3536i): 0.1250
	// [101][  5]( 0.0000 0.3536i): 0.1250
	// [110][  6]( 0.0000-0.3536i): 0.1250
	// [111][  7]( 0.0000-0.3536i): 0.1250
}

func Example_shor() {
	qasm := `
OPENQASM 3.0;

const N = 3 * 5;
const a = 7;

qubit[3] r0;
qubit[4] r1;

X r1[-1];
H r0;
CMODEXP2(a, N) r0, r1;
IQFT r0;

print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [000 0001][  0   1]( 0.2500 0.0000i): 0.0625
	// [000 0100][  0   4]( 0.2500 0.0000i): 0.0625
	// [000 0111][  0   7]( 0.2500 0.0000i): 0.0625
	// [000 1101][  0  13]( 0.2500 0.0000i): 0.0625
	// [010 0001][  2   1]( 0.2500 0.0000i): 0.0625
	// [010 0100][  2   4](-0.2500 0.0000i): 0.0625
	// [010 0111][  2   7]( 0.0000-0.2500i): 0.0625
	// [010 1101][  2  13]( 0.0000 0.2500i): 0.0625
	// [100 0001][  4   1]( 0.2500 0.0000i): 0.0625
	// [100 0100][  4   4]( 0.2500 0.0000i): 0.0625
	// [100 0111][  4   7](-0.2500 0.0000i): 0.0625
	// [100 1101][  4  13](-0.2500 0.0000i): 0.0625
	// [110 0001][  6   1]( 0.2500 0.0000i): 0.0625
	// [110 0100][  6   4](-0.2500 0.0000i): 0.0625
	// [110 0111][  6   7]( 0.0000 0.2500i): 0.0625
	// [110 1101][  6  13]( 0.0000-0.2500i): 0.0625
}

func Example_verbose() {
	qasm := `
	OPENQASM 3.0;
	include "../../testdata/stdgates.qasm";

	qubit[2] c;
	qubit t;
	
	h c[0];
	x c[1];
	ctrl @ cx c[0], c[1], t;
	print;
	`

	if _, _, err := eval(qasm, true); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// *ast.OpenQASM
	// .  *ast.DeclStmt(OPENQASM 3.0;)
	// .  []ast.Stmt
	// .  .  *ast.InclStmt(include "../../testdata/stdgates.qasm";)
	// .  .  .  *ast.DeclStmt(gate i q { U(0, 0, 0) q; })
	// .  .  .  .  *ast.GateDecl(gate i q { U(0, 0, 0) q; })
	// .  .  .  *ast.DeclStmt(gate h q { U(pi / 2.0, 0, pi) q; })
	// .  .  .  .  *ast.GateDecl(gate h q { U(pi / 2.0, 0, pi) q; })
	// .  .  .  *ast.DeclStmt(gate x q { U(pi, 0, pi) q; })
	// .  .  .  .  *ast.GateDecl(gate x q { U(pi, 0, pi) q; })
	// .  .  .  *ast.DeclStmt(gate y q { U(pi, pi / 2.0, pi / 2.0) q; })
	// .  .  .  .  *ast.GateDecl(gate y q { U(pi, pi / 2.0, pi / 2.0) q; })
	// .  .  .  *ast.DeclStmt(gate z q { Z q; })
	// .  .  .  .  *ast.GateDecl(gate z q { Z q; })
	// .  .  .  *ast.DeclStmt(gate cx c, t { ctrl @ x c, t; })
	// .  .  .  .  *ast.GateDecl(gate cx c, t { ctrl @ x c, t; })
	// .  .  .  *ast.DeclStmt(gate cy c, t { ctrl @ y c, t; })
	// .  .  .  .  *ast.GateDecl(gate cy c, t { ctrl @ y c, t; })
	// .  .  .  *ast.DeclStmt(gate cz c, t { ctrl @ z c, t; })
	// .  .  .  .  *ast.GateDecl(gate cz c, t { ctrl @ z c, t; })
	// .  .  .  *ast.DeclStmt(gate ch c, t { ctrl @ h c, t; })
	// .  .  .  .  *ast.GateDecl(gate ch c, t { ctrl @ h c, t; })
	// .  .  .  *ast.DeclStmt(gate ccx c0, c1, t { ctrl @ ctrl @ x c0, c1, t; })
	// .  .  .  .  *ast.GateDecl(gate ccx c0, c1, t { ctrl @ ctrl @ x c0, c1, t; })
	// .  .  .  *ast.DeclStmt(gate swap q, p { cx q, p; cx p, q; cx q, p; })
	// .  .  .  .  *ast.GateDecl(gate swap q, p { cx q, p; cx p, q; cx q, p; })
	// .  .  *ast.DeclStmt(qubit[2] c;)
	// .  .  .  *ast.GenDecl(qubit[2] c)
	// .  .  *ast.DeclStmt(qubit t;)
	// .  .  .  *ast.GenDecl(qubit t)
	// .  .  *ast.ExprStmt(h c[0];)
	// .  .  .  *ast.CallExpr(h c[0])
	// .  .  .  *ast.GateDecl(gate h q { U(pi / 2.0, 0, pi) q; })
	// .  .  .  .  *ast.BlockStmt({ U(pi / 2.0, 0, pi) q; })
	// .  .  .  .  .  *ast.ApplyStmt(U(pi / 2.0, 0, pi) q;)
	// .  .  .  .  .  .  *ast.InfixExpr(pi / 2.0)
	// .  .  .  .  .  .  .  *ast.BasicLit(pi)
	// .  .  .  .  .  .  .  .  return *object.Float(3.141592653589793)
	// .  .  .  .  .  .  .  *ast.BasicLit(2.0)
	// .  .  .  .  .  .  .  .  return *object.Float(2)
	// .  .  .  .  .  .  .  return *object.Float(1.5707963267948966)
	// .  .  .  .  .  .  *ast.BasicLit(0)
	// .  .  .  .  .  .  .  return *object.Int(0)
	// .  .  .  .  .  .  *ast.BasicLit(pi)
	// .  .  .  .  .  .  .  return *object.Float(3.141592653589793)
	// .  .  *ast.ExprStmt(x c[1];)
	// .  .  .  *ast.CallExpr(x c[1])
	// .  .  .  *ast.GateDecl(gate x q { U(pi, 0, pi) q; })
	// .  .  .  .  *ast.BlockStmt({ U(pi, 0, pi) q; })
	// .  .  .  .  .  *ast.ApplyStmt(U(pi, 0, pi) q;)
	// .  .  .  .  .  .  *ast.BasicLit(pi)
	// .  .  .  .  .  .  .  return *object.Float(3.141592653589793)
	// .  .  .  .  .  .  *ast.BasicLit(0)
	// .  .  .  .  .  .  .  return *object.Int(0)
	// .  .  .  .  .  .  *ast.BasicLit(pi)
	// .  .  .  .  .  .  .  return *object.Float(3.141592653589793)
	// .  .  *ast.ExprStmt(ctrl @ cx c[0], c[1], t;)
	// .  .  .  *ast.CallExpr(ctrl @ cx c[0], c[1], t)
	// .  .  .  *ast.GateDecl(gate cx c, t { ctrl @ x c, t; })
	// .  .  .  .  *ast.BlockStmt({ ctrl @ x c, t; })
	// .  .  .  .  .  *ast.ExprStmt(ctrl @ ctrl @ x _v0[0], c, t;)
	// .  .  .  .  .  .  *ast.CallExpr(ctrl @ ctrl @ x _v0[0], c, t)
	// .  .  .  .  .  .  *ast.GateDecl(gate x q { U(pi, 0, pi) q; })
	// .  .  .  .  .  .  .  *ast.BlockStmt({ U(pi, 0, pi) q; })
	// .  .  .  .  .  .  .  .  *ast.ApplyStmt(ctrl @ ctrl @ U(pi, 0, pi) _v1, _v0[0], q;)
	// .  .  .  .  .  .  .  .  .  *ast.BasicLit(pi)
	// .  .  .  .  .  .  .  .  .  .  return *object.Float(3.141592653589793)
	// .  .  .  .  .  .  .  .  .  *ast.BasicLit(0)
	// .  .  .  .  .  .  .  .  .  .  return *object.Int(0)
	// .  .  .  .  .  .  .  .  .  *ast.BasicLit(pi)
	// .  .  .  .  .  .  .  .  .  .  return *object.Float(3.141592653589793)
	// .  .  *ast.PrintStmt(print;)
	// [01 0][  1   0]( 0.7071 0.0000i): 0.5000
	// [11 1][  3   1]( 0.7071 0.0000i): 0.5000
}

func TestEvaluator_Eval(t *testing.T) {
	var cases = []struct {
		in     string
		qstate []qubit.State
		cstate []state
		hasErr bool
	}{
		{
			`
			qubit q;
			U(pi/2.0, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{1},
					BinaryString: []string{"1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			U(pi/2.0, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.5, 0),
					Int:          []int64{0},
					BinaryString: []string{"00"},
				},
				{
					Amplitude:    complex(0.5, 0),
					Int:          []int64{1},
					BinaryString: []string{"01"},
				},
				{
					Amplitude:    complex(0.5, 0),
					Int:          []int64{2},
					BinaryString: []string{"10"},
				},
				{
					Amplitude:    complex(0.5, 0),
					Int:          []int64{3},
					BinaryString: []string{"11"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			U(pi/2.0, 0, pi) q[0];
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{0},
					BinaryString: []string{"00"},
				},
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{2},
					BinaryString: []string{"10"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			bit[2] c;
			
			U(pi, 0, pi) q;
			c = measure q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3},
					BinaryString: []string{"11"},
				},
			},
			[]state{
				{
					Name:  "c",
					Value: []int64{1, 1},
				},
			},
			false,
		},
		{
			`
			qubit[2] q;
			bit[2] c;
			
			U(pi, 0, pi) q;
			measure q -> c;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3},
					BinaryString: []string{"11"},
				},
			},
			[]state{
				{
					Name:  "c",
					Value: []int64{1, 1},
				},
			},
			false,
		},
		{
			`
			qubit[2] q;			
			U(pi/2.0, 0, pi) q;
			reset q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"00"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit q;			
			U(pi, tau, euler) q;
			inv @ U(pi, tau, euler) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit q;			
			inv @ U(1.0, 2.0, 3.0) q;
			inv @ inv @ U(1.0, 2.0, 3.0) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit q;			
			pow(2) @ U(pi/2.0, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit q;			
			pow(3) @ pow(2) @ pow(2) @ U(pi/2.0, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit q;			
			pow(3) @ U(pi/2.0, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{1},
					BinaryString: []string{"1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit q;
			qubit t;
			
			U(pi/2.0, 0, pi) q;
			ctrl @ U(pi, 0, pi) q, t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{0, 0},
					BinaryString: []string{"0", "0"},
				},
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{1, 1},
					BinaryString: []string{"1", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			qubit t;
			
			U(pi, 0, pi) q;
			ctrl(2) @ U(pi, 0, pi) q, t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3, 1},
					BinaryString: []string{"11", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			qubit t;
			
			U(pi, 0, pi) q[0];
			ctrl @ ctrl @ U(pi, 0, pi) q, t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{2, 0},
					BinaryString: []string{"10", "0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			qubit t;
			
			U(pi, 0, pi) q;
			ctrl @ ctrl @ U(pi, 0, pi) q, t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3, 1},
					BinaryString: []string{"11", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			qubit[2] p;
			qubit t;
			
			U(pi, 0, pi) q;
			U(pi, 0, pi) p;

			ctrl(3) @ ctrl(1) @ U(pi, 0, pi) q, p, t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3, 3, 1},
					BinaryString: []string{"11", "11", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			qubit[2] p;
			qubit t;
			
			U(pi, 0, pi) q;
			U(pi, 0, pi) p;

			ctrl(1) @ ctrl(3) @ U(pi, 0, pi) q, p, t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3, 3, 1},
					BinaryString: []string{"11", "11", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit[2] q;
			qubit t;
			
			U(pi, 0, pi) q[0];
			ctrl @ negctrl @ U(pi, 0, pi) q, t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{2, 1},
					BinaryString: []string{"10", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit q;
			qubit p;
			
			U(pi, 0, pi) q;
			ctrl @ U(1.0, 2.0, 3.0) q, p;
			ctrl @ inv @ U(1.0, 2.0, 3.0) q, p;
			U(pi, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0, 0},
					BinaryString: []string{"0", "0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			qubit q;
			qubit p;
			
			U(pi, 0, pi) q;
			ctrl @ U(1.0, 2.0, 3.0) q, p;
			inv @ ctrl @ U(1.0, 2.0, 3.0) q, p;
			U(pi, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0, 0},
					BinaryString: []string{"0", "0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate u(a, b, c) q { U(a, b, c) q; }

			qubit[2] q;
			u(pi, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3},
					BinaryString: []string{"11"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate  u(a, b, c) q {       U(a, b, c) q; }
			gate iu(a, b, c) q { inv @ U(a, b, c) q; }


			qubit[2] q;
			u(pi/2.0, 0, pi) q;
			iu(pi/2.0, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"00"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate pu(p, a, b, c) q { pow(p) @ U(a, b, c) q; }

			qubit q;
			pu(2, pi, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate pu(p, a, b, c) q { pow(p) @ U(a, b, c) q; }

			qubit q;
			pu(3, pi, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{1},
					BinaryString: []string{"1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate u(a, b, c) q {
				U(a, b, c) q;
				U(c, b, a) q;
			}
			
			gate iu(a, b, c) q {
				inv @ U(c, b, a) q;
				inv @ U(a, b, c) q;
			}
			
			qubit q;
			u(1.0, 2.0, 3.0) q;
			iu(1.0, 2.0, 3.0) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			include "../../testdata/stdgates.qasm";

			qubit q;
			qubit p;
			
			x q;
			cx q, p;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{1, 1},
					BinaryString: []string{"1", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			include "../../testdata/stdgates.qasm";

			qubit[2] q;
			qubit[2] p;
			
			x q;
			cx q, p;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3, 3},
					BinaryString: []string{"11", "11"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			include "../../testdata/stdgates.qasm";

			qubit[2] q;
			qubit p;
			
			x q;
			cx q, p;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{3, 0},
					BinaryString: []string{"11", "0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			include "../../testdata/stdgates.qasm";

			qubit q;
			qubit[2] p;
			
			x q;
			cx q, p;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{1, 3},
					BinaryString: []string{"1", "11"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate u(a, b, c) q {
				U(pi, 0, pi) q;
				U(a, b, c) q;
			}
			
			qubit q;
			
			u(pi/2.0, 0, pi) q;
			inv @ u(pi/2.0, 0, pi) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate u(a, b, c) q {
				U(pi, 0, pi) q;
				U(a, b, c) q;
			}
			
			qubit q;
			
			pow(3) @ u(1.0, 2.0, 3.0) q;
			inv @ pow(3) @ u(1.0, 2.0, 3.0) q;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{0},
					BinaryString: []string{"0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate icu(x, y, z) c, t { inv @ ctrl @ U(x, y, z) c, t; }
			
			qubit q;
			qubit p;
			
			U(pi, 0, pi) q;
			icu(1, 2, 3) q, p;
			inv @ icu(1, 2, 3) q, p;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{1, 0},
					BinaryString: []string{"1", "0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate cu(x, y, z) c, t { ctrl @ U(x, y, z) c, t; }
			
			qubit q;
			qubit p;
			
			U(pi, 0, pi) q;
			pow(2) @ cu(pi/2.0, 0, pi) q, p;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(1, 0),
					Int:          []int64{1, 0},
					BinaryString: []string{"1", "0"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			gate cu(x, y, z) c, t { ctrl @ U(x, y, z) c, t; }
			
			qubit q;
			qubit p;
			
			U(pi, 0, pi) q;
			pow(3) @ cu(pi/2.0, 0, pi) q, p;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{1, 0},
					BinaryString: []string{"1", "0"},
				},
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{1, 1},
					BinaryString: []string{"1", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			include "../../testdata/stdgates.qasm";

			qubit c;
			qubit t;
			
			h c;
			ctrl @ x c, t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{0, 0},
					BinaryString: []string{"0", "0"},
				},
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{1, 1},
					BinaryString: []string{"1", "1"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			include "../../testdata/stdgates.qasm";

			qubit[2] q;

			h q[0];
			ctrl @ x q[0], q[1];
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{0},
					BinaryString: []string{"00"},
				},
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{3},
					BinaryString: []string{"11"},
				},
			},
			[]state{},
			false,
		},
		{
			`
			include "../../testdata/stdgates.qasm";

			qubit[2] c;
			qubit t;
			
			h c[0];
			x c[1];
			ctrl @ cx c[0], c[1], t;
			`,
			[]qubit.State{
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{1, 0},
					BinaryString: []string{"01", "0"},
				},
				{
					Amplitude:    complex(0.7071067811865475, 0),
					Int:          []int64{3, 1},
					BinaryString: []string{"11", "1"},
				},
			},
			[]state{},
			false,
		},
	}

	for _, c := range cases {
		qs, cs, err := eval(c.in)
		if (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}

		if !equals(cs, c.cstate) {
			t.Errorf("got=%v, want=%v", cs, c.cstate)
			continue
		}

		if !qubit.Equals(qs, c.qstate) {
			t.Errorf("got=%v, want=%v", qs, c.qstate)
			continue
		}
	}
}

func TestEvaluator_Assign(t *testing.T) {
	var cases = []struct {
		p      *ast.OpenQASM
		s      *ast.AssignStmt
		hasErr bool
	}{
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{},
			},
			&ast.AssignStmt{
				Right: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{Name: "q"},
						},
					},
				},
				Left: &ast.IdentExpr{
					Name: "c",
				},
			},
			true,
		},
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{
					&ast.DeclStmt{
						Decl: &ast.GenDecl{
							Kind: lexer.BIT,
							Name: "c",
							Type: &ast.IdentExpr{
								Name: "bit",
							},
						},
					},
				},
			},
			&ast.AssignStmt{
				Right: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{Name: "q"},
						},
					},
				},
				Left: &ast.IdentExpr{
					Name: "c",
				},
			},
			true,
		},
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{
					&ast.DeclStmt{
						Decl: &ast.GenDecl{
							Kind: lexer.QUBIT,
							Name: "q",
							Type: &ast.IdentExpr{
								Name: "qubit",
							},
						},
					},
				},
			},
			&ast.AssignStmt{
				Right: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{Name: "q"},
						},
					},
				},
				Left: &ast.IdentExpr{
					Name: "c",
				},
			},
			true,
		},
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{
					&ast.DeclStmt{
						Decl: &ast.GenDecl{
							Kind: lexer.QUBIT,
							Name: "q",
							Type: &ast.IdentExpr{
								Name: "qubit",
							},
						},
					},
					&ast.DeclStmt{
						Decl: &ast.GenDecl{
							Kind: lexer.BIT,
							Name: "c",
							Type: &ast.IdentExpr{
								Name: "bit",
							},
						},
					},
				},
			},
			&ast.AssignStmt{
				Right: &ast.MeasureExpr{
					QArgs: ast.ExprList{
						List: []ast.Expr{
							&ast.IdentExpr{Name: "q"},
						},
					},
				},
				Left: &ast.IdentExpr{
					Name: "c",
				},
			},
			false,
		},
	}

	for _, c := range cases {
		e := evaluator.Default()
		if err := e.Eval(c.p); err != nil {
			t.Errorf("err: %v", err)
			continue
		}

		if err := e.Assign(c.s, e.Env); (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}
	}
}

func TestEvaluator_Reset(t *testing.T) {
	var cases = []struct {
		p      *ast.OpenQASM
		s      *ast.ResetStmt
		hasErr bool
	}{
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{},
			},
			&ast.ResetStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Name: "q",
						},
					},
				},
			},
			true,
		},
	}

	for _, c := range cases {
		e := evaluator.Default()
		if err := e.Eval(c.p); err != nil {
			t.Errorf("err: %v", err)
			continue
		}

		if err := e.Reset(c.s, e.Env); (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}
	}
}

func TestEvaluator_Print(t *testing.T) {
	var cases = []struct {
		p      *ast.OpenQASM
		s      *ast.PrintStmt
		hasErr bool
	}{
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{
					&ast.DeclStmt{
						Decl: &ast.GenDecl{
							Kind: lexer.QUBIT,
							Name: "q",
							Type: &ast.IdentExpr{
								Name: "qubit",
							},
						},
					},
				},
			},
			&ast.PrintStmt{
				QArgs: ast.ExprList{
					List: []ast.Expr{
						&ast.IdentExpr{
							Name: "p",
						},
					},
				},
			},
			true,
		},
	}

	for _, c := range cases {
		e := evaluator.Default()
		if err := e.Eval(c.p); err != nil {
			t.Errorf("err: %v", err)
			continue
		}

		if err := e.Print(c.s, e.Env); (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}
	}
}

func TestEvaluator_Println(t *testing.T) {
	var cases = []struct {
		p      *ast.OpenQASM
		hasErr bool
	}{
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{},
			},
			false,
		},
	}

	for _, c := range cases {
		e := evaluator.Default()
		if err := e.Eval(c.p); err != nil {
			t.Errorf("err: %v", err)
			continue
		}

		if err := e.Println(); (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}
	}
}

func TestEvaluator_Block(t *testing.T) {
	var cases = []struct {
		p      *ast.OpenQASM
		s      *ast.BlockStmt
		want   object.Object
		hasErr bool
	}{
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{},
			},
			&ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Result: &ast.BasicLit{
							Kind:  lexer.INT,
							Value: "123",
						},
					},
				},
			},
			&object.Int{Value: 123},
			false,
		},
	}

	for _, c := range cases {
		e := evaluator.Default()
		if err := e.Eval(c.p); err != nil {
			t.Errorf("err: %v", err)
			continue
		}

		ret, err := e.Block(c.s, e.Env)
		if (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}

		got := ret.(*object.ReturnValue).Value
		if got.Type() != c.want.Type() {
			t.Errorf("got=%T, want=%T", got, c.want)
			continue
		}

		if got.String() != c.want.String() {
			t.Errorf("got=%v, want=%v", got, c.want)
			continue
		}
	}
}

func TestEvaluator_GenConst(t *testing.T) {
	var cases = []struct {
		p      *ast.OpenQASM
		s      *ast.GenConst
		hasErr bool
	}{
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{},
			},
			&ast.GenConst{
				Name: "foo",
				Value: &ast.BasicLit{
					Kind:  lexer.EOF,
					Value: "1.2",
				},
			},
			true,
		},
	}

	for _, c := range cases {
		e := evaluator.Default()
		if err := e.Eval(c.p); err != nil {
			t.Errorf("err: %v", err)
			continue
		}

		if err := e.GenConst(c.s, e.Env); (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}
	}
}

func TestEvaluator_Return(t *testing.T) {
	var cases = []struct {
		p      *ast.OpenQASM
		s      *ast.ReturnStmt
		hasErr bool
	}{
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{},
			},
			&ast.ReturnStmt{
				Result: &ast.BasicLit{
					Kind: lexer.PI,
				},
			},
			false,
		},
		{
			&ast.OpenQASM{
				Stmts: []ast.Stmt{},
			},
			&ast.ReturnStmt{
				Result: &ast.IdentExpr{
					Name: "foo",
				},
			},
			true,
		},
	}

	for _, c := range cases {
		e := evaluator.Default()
		if err := e.Eval(c.p); err != nil {
			t.Errorf("err: %v", err)
			continue
		}

		if _, err := e.Return(c.s, e.Env); (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}
	}
}
