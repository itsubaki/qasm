package evaluator_test

import (
	"fmt"
	"strings"

	"github.com/itsubaki/qasm/pkg/evaluator"
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

	return nil
}

func Example_include() {
	qasm := `
OPENQASM 3.0;
include "../../testdata/stdgates.qasm";
`

	if err := eval(qasm, true); err != nil {
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

func Example_measure() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
qubit[2] p;

U(pi, 0, pi) q, p;
measure q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11][  3   3]( 1.0000 0.0000i): 1.0000
}

func Example_measure_assign() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
bit[2] c;
bit[2] b;

U(pi, 0, pi) q;
c = measure q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11][  3]( 1.0000 0.0000i): 1.0000
	// c: [1 1]
	// b: [0 0]
}

func Example_measure_arrow() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
bit[2] c;

U(pi, 0, pi) q;

measure q -> c;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11][  3]( 1.0000 0.0000i): 1.0000
	// c: [1 1]
}

func Example_reset() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;

U(pi/2.0, 0, pi) q;
reset q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func Example_print() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;

U(pi, 0, pi) q[0];

print q[0];
print q[0], q[1];
print q;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [1][  1]( 1.0000 0.0000i): 1.0000
	// [1 0][  1   0]( 1.0000 0.0000i): 1.0000
	// [10][  2]( 1.0000 0.0000i): 1.0000
}

func Example_u() {
	qasm := `
OPENQASM 3.0;

qubit q;

U(pi/2.0, 0, pi) q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func Example_u_arr() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;

U(pi/2.0, 0, pi) q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func Example_u_index() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;

U(pi/2.0, 0, pi) q[0];
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [10][  2]( 0.7071 0.0000i): 0.5000
}

func Example_qFT() {
	qasm := `
OPENQASM 3.0;

qubit[3] q;

U(pi, 0, pi) q[1];
QFT q;

print;
`

	if err := eval(qasm); err != nil {
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

	if err := eval(qasm); err != nil {
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

func Example_inv() {
	qasm := `
OPENQASM 3.0;

qubit q;

U(1.0, 2.0, 3.0) q;
print;

inv @ U(1.0, 2.0, 3.0) q;
print;

inv @ inv @ U(1.0, 2.0, 3.0) q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 0.8776 0.0000i): 0.7702
	// [1][  1](-0.1995 0.4359i): 0.2298
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 0.8776 0.0000i): 0.7702
	// [1][  1](-0.1995 0.4359i): 0.2298
}

func Example_pow() {
	qasm := `
OPENQASM 3.0;

qubit q;

U(pi/2.0, 0, pi) q;
U(pi/2.0, 0, pi) q;
print;

pow(2) @ U(pi/2.0, 0, pi) q;
print;
reset q;

pow(3) @ pow(2) @ pow(2) @ U(pi/2.0, 0, pi) q;
print;
reset q;

pow(3) @ U(pi/2.0, 0, pi) q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl_U() {
	qasm := `
OPENQASM 3.0;

qubit q;
qubit r;

U(pi/2.0, 0, pi) q;
ctrl @ U(pi, 0, pi) q, r;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0 0][  0   0]( 0.7071 0.0000i): 0.5000
	// [1 1][  1   1]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl2_U() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
qubit r;

U(pi, 0, pi) q;
print;

ctrl(2) @ U(pi, 0, pi) q, r;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 0][  3   0]( 1.0000 0.0000i): 1.0000
	// [11 1][  3   1]( 1.0000 0.0000i): 1.0000
}

func Example_ctrlctrl_U() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
qubit r;

U(pi, 0, pi) q[0];
ctrl @ ctrl @ U(pi, 0, pi) q, r;
print;

U(pi, 0, pi) q[1];
ctrl @ ctrl @ U(pi, 0, pi) q, r;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 0][  2   0]( 1.0000 0.0000i): 1.0000
	// [11 1][  3   1]( 1.0000 0.0000i): 1.0000
}

func Example_ctrl3ctrl1_U() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
qubit[2] p;
qubit r;

U(pi, 0, pi) q;
U(pi, 0, pi) p;
print;

ctrl(3) @ ctrl(1) @ U(pi, 0, pi) q, p, r;
print;

ctrl(1) @ ctrl(3) @ U(pi, 0, pi) q, p, r;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [11 11 0][  3   3   0]( 1.0000 0.0000i): 1.0000
	// [11 11 1][  3   3   1]( 1.0000 0.0000i): 1.0000
	// [11 11 0][  3   3   0]( 1.0000 0.0000i): 1.0000
}

func Example_negctrl_U() {
	qasm := `
OPENQASM 3.0;

qubit[2] q;
qubit r;

U(pi, 0, pi) q[0];
ctrl @ negctrl @ U(pi, 0, pi) q, r;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [10 1][  2   1]( 1.0000 0.0000i): 1.0000
}

func Example_gate_u() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q { U(a, b, c) q; }

qubit[2] q;
u(pi, 0, pi) q;
print;
`

	if err := eval(qasm); err != nil {
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

	if err := eval(qasm); err != nil {
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

pow3(pi, 0, pi) q;
print;
reset q;

powu(pi, 0, pi, 2) q;
print;
reset q;

powu(pi, 0, pi, 3) q;
print;
reset q;

powu(pi, 0, pi, 4) q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [1][  1]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [1][  1]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func Example_inv_gate() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q { U(a, b, c) q; }

qubit q;

U(1, 2, 3) q;
print;

inv @ u(1, 2, 3) q;
print;

inv @ u(1, 2, 3) q;
print;

U(1, 2, 3) q;
print;

inv @ inv @ u(1, 2, 3) q;
print;

inv @ u(1, 2, 3) q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 0.8776 0.0000i): 0.7702
	// [1][  1](-0.1995 0.4359i): 0.2298
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 0.8776 0.0000i): 0.7702
	// [1][  1]( 0.4746 0.0677i): 0.2298
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 0.8776 0.0000i): 0.7702
	// [1][  1](-0.1995 0.4359i): 0.2298
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func Example_pow_gate() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q { U(a, b, c) q; }

qubit q;

pow(2) @ u(pi/2.0, 0, pi) q;
print;
reset q;

pow(3) @ pow(2) @ u(pi/2.0, 0, pi) q;
print;
reset q;

pow(3) @ u(pi/2.0, 0, pi) q;
print;
reset q;

pow(3) @ pow(3) @ u(pi/2.0, 0, pi) q;
print;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func Example_ctrl_gate() {
	qasm := `
OPENQASM 3.0;

gate u(a, b, c) q { U(a, b, c) q; }

qubit q;
qubit r;

ctrl @ u(pi/2.0, 0, pi) q, r;
`

	if err := eval(qasm); err != nil {
		fmt.Printf("eval: %v\n", err)
		return
	}

	// Not Implemented.
}
