package evaluator_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/itsubaki/qasm/pkg/ast"
	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/evaluator/object"
	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
)

func eval(qasm string, verbose ...bool) error {
	l := lexer.New(strings.NewReader(qasm))
	p := parser.New(l)

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		return fmt.Errorf("parse: %v\n", errs)
	}

	opts := evaluator.Opts{
		Verbose: false,
	}
	if verbose != nil {
		opts.Verbose = verbose[0]
	}

	e := evaluator.Default(opts)
	if err := e.Eval(a); err != nil {
		return fmt.Errorf("eval: %v\n", err)
	}

	if err := e.Println(); err != nil {
		return fmt.Errorf("print: %v\n", err)
	}

	return nil
}

func Example_verbose() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

qubit[2] q;
reset q;

h q[0];
cx q[0], q[1];
`

	if err := eval(qasm, true); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// *ast.OpenQASM
	// .  *ast.DeclStmt(OPENQASM 3.0;)
	// .  []ast.Stmt
	// .  .  *ast.DeclStmt(gate h q { U(pi / 2.0, 0, pi) q; })
	// .  .  .  *ast.GateDecl(gate h q { U(pi / 2.0, 0, pi) q; })
	// .  .  *ast.DeclStmt(gate x q { U(pi, 0, pi) q; })
	// .  .  .  *ast.GateDecl(gate x q { U(pi, 0, pi) q; })
	// .  .  *ast.DeclStmt(gate cx c, t { ctrl @ x c, t; })
	// .  .  .  *ast.GateDecl(gate cx c, t { ctrl @ x c, t; })
	// .  .  *ast.DeclStmt(qubit[2] q;)
	// .  .  .  *ast.GenDecl(qubit[2] q)
	// .  .  *ast.ResetStmt(reset q;)
	// .  .  *ast.ExprStmt(h q[0];)
	// .  .  .  *ast.CallExpr(h q[0])
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
	// .  .  *ast.ExprStmt(cx q[0], q[1];)
	// .  .  .  *ast.CallExpr(cx q[0], q[1])
	// .  .  .  *ast.GateDecl(gate cx c, t { ctrl @ x c, t; })
	// .  .  .  .  *ast.BlockStmt({ ctrl @ x c, t; })
	// .  .  .  .  .  *ast.ExprStmt(ctrl @ x c, t;)
	// .  .  .  .  .  .  *ast.CallExpr(ctrl @ x c, t)
	// .  .  .  .  .  .  *ast.GateDecl(gate x q { U(pi, 0, pi) q; })
	// .  .  .  .  .  .  .  *ast.BlockStmt({ ctrl @ U(pi, 0, pi) c, t; })
	// .  .  .  .  .  .  .  .  *ast.ApplyStmt(ctrl @ U(pi, 0, pi) c, t;)
	// .  .  .  .  .  .  .  .  .  *ast.BasicLit(pi)
	// .  .  .  .  .  .  .  .  .  .  return *object.Float(3.141592653589793)
	// .  .  .  .  .  .  .  .  .  *ast.BasicLit(0)
	// .  .  .  .  .  .  .  .  .  .  return *object.Int(0)
	// .  .  .  .  .  .  .  .  .  *ast.BasicLit(pi)
	// .  .  .  .  .  .  .  .  .  .  return *object.Float(3.141592653589793)
	// .  .  *ast.PrintStmt(print;)
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func Example_hermite() {
	qasm := `
OPENQASM 3.0;

gate i q { U(0, 0, 0) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate y q { U(pi, pi/2.0, pi/2.0) q; }
gate z q { Z q; }

qubit[2] q;
reset q;

i q;
x q; x q;
y q; y q;
z q; z q;
h q; h q;

X q; X q;
Y q; Y q;
Z q; Z q;
H q; H q;
T q; T q;
S q; S q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_measure() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
bit[2] c;
reset q;

x q;
c = measure q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11][  3]( 1.0000 0.0000i): 1.0000
	// c: 11
}

func Example_arrow() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
bit[2] c;
reset q;

x q;
measure q -> c;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11][  3]( 1.0000 0.0000i): 1.0000
	// c: 11
}

func Example_gate() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

qubit[2] q;
reset q;

h q[0];
cx q[0], q[1];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func Example_print() {
	qasm := `
OPENQASM 3.0;

print;

qubit[2] q;
reset q;

print q;
print q[0], q[1];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
	// [0 0][  0   0]( 1.0000 0.0000i): 1.0000
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_inv() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

gate bell q, p { h q; cx q, p; }

qubit[2] q;
reset q;

bell q[0], q[1];
inv @ bell q[0], q[1];

inv @ inv @ bell q[0], q[1];
inv @ inv @ inv @bell q[0], q[1];

QFT q;
IQFT q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_pow0() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

gate bell q, p { h q; cx q, p; }

qubit[2] q;
reset q;

pow(0) @ bell q[0], q[1];
pow(0) @ U(pi/2.0, 0, pi) q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_pow1() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }

qubit q;
reset q;

pow(1) @ U(pi/2.0, 0, pi) q;
pow(1) @ h q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func Example_pow2() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

gate bell q, p { h q; cx q, p; }

qubit[2] q;
reset q;

pow(2)  @ bell q[0], q[1];
pow(-2) @ bell q[0], q[1];
pow(2)  @ U(pi/2.0, 0, pi) q;
pow(-2) @ U(pi/2.0, 0, pi) q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_ctrl() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

gate bell q, p { h q; cx q, p; }

qubit q;
qubit[2] p;
reset q, p;

x q;
ctrl @ bell q, p[0], p[1];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1 00][  1   0]( 0.7071 0.0000i): 0.5000
	// [1 11][  1   3]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl2() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

gate bell q, p { U(pi/2.0, 0, pi) q; cx q, p; }

qubit q;
qubit[2] p;
reset q, p;

x q;
ctrl @ bell q, p[0], p[1];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1 00][  1   0]( 0.7071 0.0000i): 0.5000
	// [1 11][  1   3]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl3() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

gate bell q, p { U(pi/2.0, 0, pi) q; cx q, p; }

qubit q;
qubit[2] p;
reset q, p;

x q;
ctrl(0) @ bell q, p[0], p[1];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1 00][  1   0]( 0.7071 0.0000i): 0.5000
	// [1 11][  1   3]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl4() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate xc a, b { ctrl @ x b, a; }

qubit q;
qubit[2] p;
reset q, p;

x q;
x p[1];
ctrl @ xc q, p[0], p[1];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1 11][  1   3]( 1.0000 0.0000i): 1.0000
}

func Example_negctrl() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate ncx  q0, q1 { negctrl @ x q0, q1; }

qubit[3] q;
reset q;

x q[0];
ctrl @ ncx q[0], q[1], q[2];
`

	// 000 -> 100 -> 101
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [101][  5]( 1.0000 0.0000i): 1.0000
}

func Example_negctrl2() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[3] q;
reset q;

x q[0];
ctrl(0) @ negctrl(1) @ x q[0], q[1], q[2];
`

	// 000 -> 100 -> 101
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [101][  5]( 1.0000 0.0000i): 1.0000
}

func Example_def() {
	qasm := `
OPENQASM 3.0;

def x qubit[n] q -> bit[n] {
	X q;
	return measure q;
}

qubit[2] q;
bit[2] c;
reset q;

c = x q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11][  3]( 1.0000 0.0000i): 1.0000
	// c: 11
}

func Example_bell() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }
gate bell q, p { h q; cx q, p; }

qubit[2] q;
reset q;

bell q[0], q[1];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func Example_shor() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }

const N = 15;
const a = 7;

qubit[3] r0;
qubit[4] r1;
reset r0, r1;

x r1[-1];
h r0;
CMODEXP2(a, N) r0, r1;
IQFT r0;
`

	if err := eval(qasm); err != nil {
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
func TestEvalExpr(t *testing.T) {
	var cases = []struct {
		in   ast.Expr
		want object.Object
	}{
		{
			in: &ast.BasicLit{
				Kind:  lexer.INT,
				Value: "3",
			},
			want: &object.Int{
				Value: 3,
			},
		},
		{
			in: &ast.BasicLit{
				Kind:  lexer.PI,
				Value: "pi",
			},
			want: &object.Float{
				Value: 3.141592653589793,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.PLUS,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "7",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "11",
				},
			},
			want: &object.Int{
				Value: 18,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.PLUS,
				Left: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "7",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "11",
				},
			},
			want: &object.Float{
				Value: 18,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.MINUS,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "7",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "11",
				},
			},
			want: &object.Int{
				Value: -4,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.MINUS,
				Left: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "7",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "11",
				},
			},
			want: &object.Float{
				Value: -4,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.MUL,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "7",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "11",
				},
			},
			want: &object.Int{
				Value: 77,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.MUL,
				Left: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "7",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "11",
				},
			},
			want: &object.Float{
				Value: 77,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.DIV,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "15",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "3",
				},
			},
			want: &object.Int{
				Value: 5,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.DIV,
				Left: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "15",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "3",
				},
			},
			want: &object.Float{
				Value: 5,
			},
		},
		{
			in: &ast.InfixExpr{
				Kind: lexer.MOD,
				Left: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "15",
				},
				Right: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "3",
				},
			},
			want: &object.Int{
				Value: 0,
			},
		},
		{
			in: &ast.UnaryExpr{
				Kind: lexer.MINUS,
				Value: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "3",
				},
			},
			want: &object.Int{
				Value: -3,
			},
		},
		{
			in: &ast.UnaryExpr{
				Kind: lexer.PLUS,
				Value: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "3.0",
				},
			},
			want: &object.Float{
				Value: 3.0,
			},
		},
	}

	for _, c := range cases {
		got, err := evaluator.Eval(c.in)
		if err != nil {
			t.Fatalf("in(%v): %v", c.in, err)
		}

		if got.Type() != c.want.Type() {
			t.Errorf("got=%T, want=%T", got, c.want)
		}

		if got.String() != c.want.String() {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
