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

qubit q;
reset q;

h q;
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
	// .  .  *ast.DeclStmt(qubit q;)
	// .  .  .  *ast.GenDecl(qubit q)
	// .  .  *ast.ResetStmt(reset q;)
	// .  .  *ast.ExprStmt(h q;)
	// .  .  .  *ast.CallExpr(h q)
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
	// .  .  *ast.PrintStmt(print;)
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
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

func Example_gateCtrlXq0r0() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;
reset q, r;

x q[0];
ctrl @ x q[0], r[0];
`

	// [00 00] -> [10 00] -> [10 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 10][  2   2]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrlXq0r() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;
reset q, r;
	
x q[0];
ctrl @ x q[0], r;
`

	// [00 00] -> [10 00] -> [10 11]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 11][  2   3]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrlXqr() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;
reset q, r;
	
x q[0];
ctrl @ x q, r;	
`

	// [00 00] -> [10 00] -> [10 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 10][  2   2]( 1.0000 0.0000i): 1.0000
}

func Example_gateCtrl2Xqr() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;

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

func Example_gateCtrl1Ctrl1Xq0q1r() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[2] q;
qubit[2] r;

x q[0];
ctrl(1) @ ctrl(1) @ x q[0], q[1], r;	
`

	// [00 00] -> [10 00] -> [10 10]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 10][  2   2]( 1.0000 0.0000i): 1.0000
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
gate cx a, b { ctrl @ x a, b; }

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
gate cx a, b { ctrl @ x a, b; }

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
gate cx a, b { ctrl @ x a, b; }

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

func Example_gateCtrlCX() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx a, b { ctrl @ x a, b; }

qubit q0;
qubit q1;
qubit q2;
reset q0, q1, q2;

x q0;
x q1;
ctrl @ cx q0, q1, q2;
`

	// [0 0 0] -> [1 1 0] -> [1 1 1]
	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1 1 1][  1   1   1]( 1.0000 0.0000i): 1.0000
}

func Example_gateCCX() {
	qasm := `
OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate ccx a, b, c { ctrl @ ctrl @ x a, b, c; }

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
QFT q; IQFT q;

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

U(tau, pi, euler) q;
inv @ U(tau, pi, euler) q;
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

pow(+2) @ U(1.0, 2.0, 3.0) q;
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

func Example_def() {
	qasm := `
OPENQASM 3.0;

def hoge(int[32] a, int[32] N) qubit[n] r0 -> bit[n] {
    X r0[-1];
    return measure r0;
}

const N = 3 * 5;
const a = 7;

qubit[3] r0;
bit[3] c;
reset r0;

c = hoge(a, N) r0;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [001][  1]( 1.0000 0.0000i): 1.0000
	// c: 001
}

func Example_shor() {
	qasm := `
OPENQASM 3.0;

const N = 3 * 5;
const a = 7;

qubit[3] r0;
qubit[4] r1;
reset r0, r1;

X r1[-1];
H r0;
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

func Example_shorgate() {
	qasm := `
OPENQASM 3.0;

gate shor(a ,N) r0, r1 {
	CMODEXP2(a, N) r0, r1;
	IQFT r0;
}

const N = 3 * 5;
const a = 4;

qubit[3] r0;
qubit[4] r1;
reset r0, r1;

X r1[-1];
H r0;
shor(a, N) r0, r1;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [000 0001][  0   1]( 0.5000 0.0000i): 0.2500
	// [000 0100][  0   4]( 0.5000 0.0000i): 0.2500
	// [100 0001][  4   1]( 0.5000 0.0000i): 0.2500
	// [100 0100][  4   4](-0.5000 0.0000i): 0.2500
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
			in: &ast.BasicLit{
				Kind:  lexer.STRING,
				Value: "hoge",
			},
			want: &object.String{
				Value: "hoge",
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
				Kind: lexer.MINUS,
				Value: &ast.BasicLit{
					Kind:  lexer.FLOAT,
					Value: "3.0",
				},
			},
			want: &object.Float{
				Value: -3.0,
			},
		},
		{
			in: &ast.UnaryExpr{
				Kind: lexer.PLUS,
				Value: &ast.BasicLit{
					Kind:  lexer.INT,
					Value: "3",
				},
			},
			want: &object.Int{
				Value: 3,
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
