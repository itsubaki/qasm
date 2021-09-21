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

U(pi, 0, pi) q;
U(pi, 0, pi) q;

U(pi, pi/2.0, pi/2.0) q;
U(pi, pi/2.0, pi/2.0) q;

U(pi, 0, pi) q[0];
U(pi, 0, pi) q[0];

U(pi, 0, pi) q[0], q[1];
U(pi, 0, pi) q[0], q[1];
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

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
reset q0, q1, q2;

U(pi, 0, pi) q0, q1, q2;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11 11][  3   3   3]( 1.0000 0.0000i): 1.0000
}

func Example_inv() {
	qasm := `
OPENQASM 3.0;

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
reset q0, q1, q2;

U(1.0, 2.0, 3.0) q0, q1, q2;
inv @ U(1.0, 2.0, 3.0) q0, q1, q2;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00 00 00][  0   0   0]( 1.0000 0.0000i): 1.0000
}

func Example_pow0() {
	qasm := `
OPENQASM 3.0;

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
reset q0, q1, q2;

pow(0) @ U(1.0, 2.0, 3.0) q0, q1, q2;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00 00 00][  0   0   0]( 1.0000 0.0000i): 1.0000
}

func Example_pow1() {
	qasm := `
OPENQASM 3.0;

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
reset q0, q1, q2;

pow(1) @ U(1.0, 2.0, 3.0) q0, q1, q2;
inv    @ U(1.0, 2.0, 3.0) q0, q1, q2;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00 00 00][  0   0   0]( 1.0000 0.0000i): 1.0000
}

func Example_pow2() {
	qasm := `
OPENQASM 3.0;

qubit[2] q0;
qubit[2] q1;
qubit[2] q2;
reset q0, q1, q2;

pow(2)  @ U(1.0, 2.0, 3.0) q0, q1, q2;
pow(-2) @ U(1.0, 2.0, 3.0) q0, q1, q2;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00 00 00][  0   0   0]( 1.0000 0.0000i): 1.0000
}

func Example_ctrlu() {
	qasm := `
OPENQASM 3.0;

qubit[3] q0;
qubit[3] q1;
	
U(pi, 0, pi) q0[0];
ctrl @ U(pi, 0, pi) q0, q1;	
`

	// [000 000] -> [100 000] -> [100 100]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	//
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
