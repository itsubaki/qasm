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

func Example_gateHermite() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }

qubit q;
reset q;

h q; h q;
x q; x q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gateQargs() {
	qasm := `
OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }

qubit[2] q;
reset q;

h q; h q;
x q; x q;

h q[0]; h q[0];
x q[1]; x q[1];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gateInv() {
	qasm := `
OPENQASM 3.0;

gate u q { U(1.0, 2.0, 3.0) q; }

qubit[2] q;
reset q;

u q;
inv @ u q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gatePow0() {
	qasm := `
OPENQASM 3.0;

gate u q { U(1.0, 2.0, 3.0) q; }

qubit[2] q;
reset q;

pow(0) @ u(1.0, 2.0, 3.0) q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gatePow1() {
	qasm := `
OPENQASM 3.0;

gate u q { U(1.0, 2.0, 3.0) q; }

qubit[2] q;
reset q;

pow(1) @ u(1.0, 2.0, 3.0) q;
inv    @ u(1.0, 2.0, 3.0) q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gatePow2() {
	qasm := `
OPENQASM 3.0;

gate u q { U(1.0, 2.0, 3.0) q; }

qubit[2] q;
reset q;

pow(2)  @ u(1.0, 2.0, 3.0) q;
pow(-2) @ u(1.0, 2.0, 3.0) q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrlq0q0r0() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;
reset q, r;

x q[0];
ctrl(1) @ x q[0], r[0];
`

	// [00 00] -> [10 00] -> [10 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 10][  2   2]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrlCX() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx a, b { ctrl(1) @ x b, a; }

qubit q0;
qubit q1;
qubit q2;
reset q0, q1, q2;

x q0;
x q1;
ctrl(1) @ cx q0, q1, q2;
`

	// [0 0 0] -> [1 1 0] -> [1 1 1]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1 1 1][  1   1   1]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrlCCX() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate ccx a, b, c { ctrl(1) @ ctrl(1) @ x a, b, c; }

qubit q0;
qubit q1;
qubit q2;
reset q0, q1, q2;

x q0;
x q1;
ccx q0, q1, q2;
`

	// [0 0 0] -> [1 1 0] -> [1 1 1]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1 1 1][  1   1   1]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrlq0qr() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;
reset q, r;
	
x q[0];
ctrl(2) @ x q, r;	
`

	// [00 00] -> [10 00] -> [10 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 10][  2   2]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrlq0q0r() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;
reset q, r;
	
x q[0];
ctrl(1) @ x q[0], r;
`

	// [00 00] -> [10 00] -> [10 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 11][  2   3]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrlqqr() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;

x q;
ctrl(1) @ x q, r;	
`

	// [00 00] -> [11 00] -> [11 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11][  3   3]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrl2ctrl2() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
	
x q0;
x q1;
ctrl(2) @ ctrl(2) @ x q0, q1, q2;	
`

	// [00 00 00] -> [11 11 00] -> [11 11 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11 11][  3   3   3]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrl2negc2() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
	
x q0;
ctrl(2) @ negctrl(2) @ x q0, q1, q2;	
`

	// [00 00 00] -> [11 00 00] -> [11 00 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 00 11][  3   0   3]( 1.0000 0.0000i): 1.0000
}

func Example_gateCXqr() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx a, b { ctrl(1) @ x a, b; }

qubit[2] q;
qubit[2] r;
	
x q;
cx q, r;
`

	// [00 00] -> [11 00] -> [11 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11][  3   3]( 1.0000 0.0000i): 1.0000
}

func Example_gateCXq0r0() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx a, b { ctrl(1) @ x a, b; }

qubit[2] q;
qubit[2] r;
	
x q[0];
cx q[0], r[0];
`

	// [00 00] -> [10 00] -> [10 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 10][  2   2]( 1.0000 0.0000i): 1.0000
}

func Example_gateCXq0r() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx a, b { ctrl(1) @ x a, b; }

qubit[2] q;
qubit[2] r;
	
x q[0];
cx q[0], r;
`

	// [00 00] -> [10 00] -> [10 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 11][  2   3]( 1.0000 0.0000i): 1.0000
}

func Example_gateCXqr0() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx a, b { ctrl(1) @ x a, b; }

qubit[2] q;
qubit[2] r;
	
x q;
cx q, r[0];
`

	// [00 00] -> [11 00] -> [11 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 10][  3   2]( 1.0000 0.0000i): 1.0000
}

func Example_gateCXba() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx a, b { ctrl(1) @ x b, a; }

qubit[2] q0;
qubit[2] q1;
reset q0, q1;
	
x q1;
cx q0, q1;
`

	// [00 00] -> [00 11] -> [11 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11][  3   3]( 1.0000 0.0000i): 1.0000
}

func Example_hermite() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
reset q;

X q; X q;
Y q; Y q;
Z q; Z q;
H q; H q;
T q; T q;
S q; S q;

U(0, 0, 0) q;

U(pi/2.0, 0, pi) q;
U(pi/2.0, 0, pi) q;

U(pi, 0, pi) q[0];
U(pi, 0, pi) q[0];
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_qargs() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
reset q;

U(pi, 0, pi) q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11][  3]( 1.0000 0.0000i): 1.0000
}

func Example_inv() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
reset q;

U(1.0, 2.0, 3.0) q;
inv @ U(1.0, 2.0, 3.0) q;
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

qubit[2] q;
reset q;

pow(0) @ U(1.0, 2.0, 3.0) q;
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

qubit[2] q;
reset q;

pow(1) @ U(1.0, 2.0, 3.0) q;
inv    @ U(1.0, 2.0, 3.0) q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_pow2() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
reset q;

pow(2)  @ U(1.0, 2.0, 3.0) q;
pow(-2) @ U(1.0, 2.0, 3.0) q;
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

qubit[2] q;
qubit[2] r;
reset q, r;

U(pi, 0, pi) q[0];
ctrl(1) @ U(pi, 0, pi) q[0], r[0];
`

	// [00 00] -> [10 00] -> [10 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 10][  2   2]( 1.0000 0.0000i): 1.0000
}

func Example_ctrlqr() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
qubit[2] r;
reset q, r;
	
U(pi, 0, pi) q[0];
ctrl(2) @ U(pi, 0, pi) q, r;	
`

	// [00 00] -> [10 00] -> [10 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 10][  2   2]( 1.0000 0.0000i): 1.0000
}

func Example_ctrlq0r() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
qubit[2] r;
reset q, r;
	
U(pi, 0, pi) q[0];
ctrl(1) @ U(pi, 0, pi) q[0], r;
`

	// [00 00] -> [10 00] -> [10 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 11][  2   3]( 1.0000 0.0000i): 1.0000
}

func Example_ctrlq() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
qubit[2] r;
reset q, r;
	
U(pi, 0, pi) q;
ctrl(1) @ U(pi, 0, pi) q, r;	
`

	// [00 00] -> [11 00] -> [11 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11][  3   3]( 1.0000 0.0000i): 1.0000
}

func Example_ctrl2ctrl2() {
	qasm := `
OPENQASM 3.0;

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
reset q0, q1, q2;
	
U(pi, 0, pi) q0;
U(pi, 0, pi) q1;
ctrl(1) @ ctrl(1) @ U(pi, 0, pi) q0, q1, q2;	
`

	// [00 00 00] -> [11 11 00] -> [11 11 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11 11][  3   3   3]( 1.0000 0.0000i): 1.0000
}

func Example_ctrl2negc2() {
	qasm := `
OPENQASM 3.0;

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
reset q0, q1, q2;

U(pi, 0, pi) q0;
ctrl(2) @ negctrl(2) @ U(pi, 0, pi) q0, q1, q2;	
`

	// [00 00 00] -> [11 00 00] -> [11 00 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 00 11][  3   0   3]( 1.0000 0.0000i): 1.0000
}

func Example_measure() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
bit[2] c;
reset q;

U(pi, 0, pi) q;
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

qubit[2] q;
bit[2] c;
reset q;

U(pi, 0, pi) q;
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
