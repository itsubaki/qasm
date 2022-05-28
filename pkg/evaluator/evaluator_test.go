package evaluator_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/pkg/quantum/qubit"
	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
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

func Example_include() {
	qasm := `
OPENQASM 3.0;
include "../../testdata/stdgates.qasm";
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
}

func Example_print() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;

U(pi, 0, pi) q[0];

print q;
print q[0], q[1];
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10][  2]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
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

func Example_iQFT() {
	qasm := `
OPENQASM 3.0;

qubit[3] q;

U(pi, 0, pi) q[1];
QFT q;
IQFT q;

print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [010][  2]( 1.0000 0.0000i): 1.0000
}

func Example_cMODEXP2() {
	qasm := `
OPENQASM 3.0;

qubit[3] r0;
qubit[4] r1;

X r1[-1];
H r0;
CMODEXP2(7, 15) r0, r1;
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

func TestEvaluator_Eval(t *testing.T) {
	var cases = []struct {
		in     string
		qstate []qubit.State
		cstate []state
		hasErr bool
	}{
		{
			`
			qubit[2] q;
			qubit[2] p;
			
			U(pi, 0, pi) q, p;
			measure q;
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
			qubit[2] q;
			bit[2] c;
			bit[2] b;
			
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
				{
					Name:  "b",
					Value: []int64{0, 0},
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
			qubit q;			
			U(1.0, 2.0, 3.0) q;
			inv @ U(1.0, 2.0, 3.0) q;
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

func Example_u_ctrl_inv() {
	qasm := `
OPENQASM 3.0;

qubit q;
qubit p;
print;

U(pi, 0, pi) q;
print;

ctrl @ inv @ U(1.0, 2.0, 3.0) q, p;
print;

ctrl @ U(1.0, 2.0, 3.0) q, p;
print;

inv @ ctrl @ U(1.0, 2.0, 3.0) q, p;
print;

ctrl @ U(1.0, 2.0, 3.0) q, p;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0 0][  0   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 0.8776 0.0000i): 0.7702
	// [1 1][  1   1]( 0.4746 0.0677i): 0.2298
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 0.8776 0.0000i): 0.7702
	// [1 1][  1   1]( 0.4746 0.0677i): 0.2298
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
}

func Example_gate() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q { U(a, b, c) q; }

qubit[2] q;
u(pi, 0, pi) q;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11][  3]( 1.0000 0.0000i): 1.0000
}

func Example_gate_inv() {
	qasm := `
OPENQASM 3.0;

gate invu(a, b, c) q { inv @ U(a, b, c) q; }

qubit q;

U(1, 2, 3) q;
print;

invu(1, 2, 3) q;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 0.8776 0.0000i): 0.7702
	// [1][  1](-0.1995 0.4359i): 0.2298
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gate_pow() {
	qasm := `
OPENQASM 3.0;

gate pow2(a, b, c) q    { pow(2) @ U(a, b, c) q; }
gate pow3(a, b, c) q    { pow(3) @ U(a, b, c) q; }
gate powu(a, b, c, p) q { pow(p) @ U(a, b, c) q; }

qubit q;

pow2(pi, 0, pi) q;
print;
reset q;

powu(pi, 0, pi, 2) q;
print;
reset q;

powu(pi, 0, pi, 4) q;
print;
reset q;

pow3(pi, 0, pi) q;
print;
reset q;

powu(pi, 0, pi, 3) q;
print;
reset q;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [1][  1]( 1.0000 0.0000i): 1.0000
	// [1][  1]( 1.0000 0.0000i): 1.0000
}

func Example_gate_block() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q {
	U(a, b, c) q;
	print;
	U(c, b, a) q;
	print;
}

qubit[2] q;
u(pi, 0, pi) q;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11][  3]( 1.0000 0.0000i): 1.0000
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gate_block_inv() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q {
	U(a, b, c) q;
	U(c, b, a) q;
}

gate invu(a, b, c) q {
	inv @ U(c, b, a) q;
	inv @ U(a, b, c) q;
}

qubit q;

u(pi/2.0, 0, pi) q;
print;

invu(pi/2.0, 0, pi) q;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 0.0000-0.7071i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gate_block_pow() {
	qasm := `
OPENQASM 3.0;

gate pow23(a, b, c) q    {
	pow(2) @ U(a, b, c) q;
	pow(3) @ U(a, b, c) q;
}

qubit q;

pow(2) @ U (pi/2.0, 0, pi) q;
pow(3) @ U (pi/2.0, 0, pi) q;
print;
reset q;

pow23(pi/2.0, 0, pi) q;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func Example_gate_ctrl() {
	qasm := `
OPENQASM 3.0;

gate cx c, t { ctrl @ U(pi, 0, pi) c, t; }

qubit q;
qubit p;
print;

U(pi, 0, pi) q;
print;

cx q, p;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0 0][  0   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
	// [1 1][  1   1]( 1.0000 0.0000i): 1.0000
}

func Example_gate_ctrl_qreg() {
	qasm := `
OPENQASM 3.0;

gate cx c, t { ctrl @ U(pi, 0, pi) c, t; }

qubit[2] q;
qubit[2] p;
print;

U(pi, 0, pi) q;
print;

cx q, p;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00 00][  0   0]( 1.0000 0.0000i): 1.0000
	// [11 00][  3   0]( 1.0000 0.0000i): 1.0000
	// [11 11][  3   3]( 1.0000 0.0000i): 1.0000
}

func Example_gate_ctrl_qreg2() {
	qasm := `
OPENQASM 3.0;

gate cx c, t { ctrl @ U(pi, 0, pi) c, t; }

qubit[2] q;
qubit p;
print;

U(pi, 0, pi) q;
print;

cx q, p;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00 0][  0   0]( 1.0000 0.0000i): 1.0000
	// [11 0][  3   0]( 1.0000 0.0000i): 1.0000
	// [11 0][  3   0]( 1.0000 0.0000i): 1.0000
}

func Example_gate_ctrl_qreg3() {
	qasm := `
OPENQASM 3.0;

gate cx c, t { ctrl @ U(pi, 0, pi) c, t; }

qubit q;
qubit[2] p;
print;

U(pi, 0, pi) q;
print;

cx q, p;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0 00][  0   0]( 1.0000 0.0000i): 1.0000
	// [1 00][  1   0]( 1.0000 0.0000i): 1.0000
	// [1 11][  1   3]( 1.0000 0.0000i): 1.0000
}

func Example_inv_gate() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q {
	U(pi, 0, pi) q;
	U(a, b, c) q;
}

qubit q;

u(pi/2.0, 0, pi) q;
print;

inv @ u(pi/2.0, 0, pi) q;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1](-0.7071 0.0000i): 0.5000
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func Example_pow_gate() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q {
	U(pi, 0, pi) q;
	U(a, b, c) q;
}

qubit q;
print;

U(pi, 0, pi) q;
U(pi/2.0, pi, pi) q;

U(pi, 0, pi) q;
U(pi/2.0, pi, pi) q;

U(pi, 0, pi) q;
U(pi/2.0, pi, pi) q;
print;

reset q;
print;

pow(3) @ u(pi/2.0, pi, pi) q;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func Example_inv_gate_ctrl() {
	qasm := `
OPENQASM 3.0;

gate cu(x, y, z) c, t { inv @ ctrl @ U(x, y, z) c, t; }

qubit q;
qubit p;
print;

U(pi, 0, pi) q;
print;

cu(1, 2, 3) q, p;
print;

inv @ cu(1, 2, 3) q, p;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0 0][  0   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 0.8776 0.0000i): 0.7702
	// [1 1][  1   1]( 0.4746 0.0677i): 0.2298
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
}

func Example_pow_gate_ctrl() {
	qasm := `
OPENQASM 3.0;

gate cu(x, y, z) c, t { ctrl @ U(x, y, z) c, t; }

qubit q;
qubit p;
print;

U(pi, 0, pi) q;
print;

cu(pi/2.0, 0, pi) q, p;
cu(pi/2.0, 0, pi) q, p;
cu(pi/2.0, 0, pi) q, p;
print;

reset q, p;
print;

U(pi, 0, pi) q;
print;

pow(3) @ cu(pi/2.0, 0, pi) q, p;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0 0][  0   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 0.7071 0.0000i): 0.5000
	// [1 1][  1   1]( 0.7071 0.0000i): 0.5000
	// [0 0][  0   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 0.7071 0.0000i): 0.5000
	// [1 1][  1   1]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl_x() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }

qubit c;
qubit t;

h c;
print;

ctrl @ x c, t;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0 0][  0   0]( 0.7071 0.0000i): 0.5000
	// [1 0][  1   0]( 0.7071 0.0000i): 0.5000
	// [0 0][  0   0]( 0.7071 0.0000i): 0.5000
	// [1 1][  1   1]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl_ctrl_x() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx p, q { ctrl @ x p, q; }

qubit[2] c;
qubit t;

h c[0];
x c[1];
print;

ctrl @ cx c[0], c[1], t;
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [01 0][  1   0]( 0.7071 0.0000i): 0.5000
	// [11 0][  3   0]( 0.7071 0.0000i): 0.5000
	// [01 0][  1   0]( 0.7071 0.0000i): 0.5000
	// [11 1][  3   1]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl_x_index() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }

qubit[2] q;

h q[0];
print;

ctrl @ x q[0], q[1];
print;
`

	if _, _, err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [10][  2]( 0.7071 0.0000i): 0.5000
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}
