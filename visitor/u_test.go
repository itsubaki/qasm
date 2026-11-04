package visitor_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/environ"
	xparser "github.com/itsubaki/qasm/parser"
	"github.com/itsubaki/qasm/visitor"
)

func TestVisitor_U_OpenQASM3(t *testing.T) {
	for _, params := range [][3]float64{
		{0, 0, 0},
		{math.Pi, 0, math.Pi},
		{math.Pi / 2, 0, math.Pi},
		{2 * math.Pi, 0, 0},
		{0.7, 0.3, 1.1},
		{0.7 + 2*math.Pi, 0.3, 1.1},
		{0.7, 0.3 + 2*math.Pi, 1.1},
		{0.7, 0.3, 1.1 + 2*math.Pi},
	} {
		theta, phi, lambda := params[0], params[1], params[2]
		// Directly use the specification's matrix, independently of gate.U.
		// https://openqasm.com/versions/3.0/language/gates.html#built-in-single-qubit-gate-u
		e := cmplx.Exp(complex(0, theta))
		want := [2][2]complex128{
			{(1 + e) / 2, -1i * cmplx.Exp(complex(0, lambda)) * (1 - e) / 2},
			{1i * cmplx.Exp(complex(0, phi)) * (1 - e) / 2, cmplx.Exp(complex(0, phi+lambda)) * (1 + e) / 2},
		}

		tree, err := xparser.Parse(fmt.Sprintf("U(%.17g, %.17g, %.17g) q;", theta, phi, lambda))
		if err != nil {
			t.Fatal(err)
		}

		for col := range 2 {
			var basis [2]complex128
			basis[col] = 1
			qsim := q.New()
			env := environ.New()
			env.Qubit["q"] = qsim.New(basis[:]...)
			if err := visitor.New(qsim, env).Run(tree); err != nil {
				t.Fatal(err)
			}

			for row, got := range qsim.Amplitude() {
				if cmplx.Abs(got-want[row][col]) > 1e-12 {
					t.Errorf("U[%d,%d]=%v, want %v", row, col, got, want[row][col])
				}
			}
		}
	}
}

func TestVisitor_U_controlledPhase(t *testing.T) {
	cases := []struct {
		name string
		text string
		want []complex128
	}{
		{
			name: "controlled 2*pi is identity",
			text: `qubit[2] q; h q[0]; ctrl @ U(2*pi, 0, 0) q[0], q[1]; h q[0];`,
			want: []complex128{1, 0, 0, 0},
		},
		{
			name: "negative controlled 2*pi is identity",
			text: `qubit[2] q; h q[0]; negctrl @ U(2*pi, 0, 0) q[0], q[1]; h q[0];`,
			want: []complex128{1, 0, 0, 0},
		},
		{
			name: "controlled iX produces a relative phase",
			text: `qubit[2] q; h q[0]; ctrl @ U(pi, 0, pi) q[0], q[1]; cx q[0], q[1]; h q[0];`,
			want: []complex128{0.5 + 0.5i, 0, 0.5 - 0.5i, 0},
		},
		{
			name: "gphase applies once with spectators",
			text: `qubit[3] q; gphase(pi/2);`,
			want: []complex128{1i, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "broadcast H includes one correction per call",
			text: `qubit[2] q; qubit spectator; h q;`,
			want: []complex128{0.5, 0, 0.5, 0, 0.5, 0, 0.5, 0},
		},
		{
			name: "broadcast CX preserves phase",
			text: `qubit[2] c; qubit[2] t; x c; cx c, t;`,
			want: []complex128{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(`include "../testdata/stdgates.qasm";` + c.text)
		if err != nil {
			t.Fatal(err)
		}

		got := qsim.Amplitude()
		if len(got) != len(c.want) {
			t.Fatalf("got %d amplitudes, want %d", len(got), len(c.want))
		}

		for i := range got {
			if cmplx.Abs(got[i]-c.want[i]) > 1e-12 {
				t.Errorf("amplitude[%d]=%v, want %v", i, got[i], c.want[i])
			}
		}
	}
}

func TestVisitor_controlledGphase(t *testing.T) {
	cases := []struct {
		text   string
		active func(int) bool
	}{
		{`ctrl @ gphase(pi/2) q[0];`, func(i int) bool { return i&0b100 != 0 }},
		{`negctrl @ gphase(pi/2) q[0];`, func(i int) bool { return i&0b100 == 0 }},
		{`ctrl(2) @ gphase(pi/2) q[0], q[1];`, func(i int) bool { return i&0b110 == 0b110 }},
		{`ctrl @ negctrl @ gphase(pi/2) q[0], q[1];`, func(i int) bool { return i&0b110 == 0b100 }},
		{`inv @ ctrl @ gphase(-pi/2) q[0];`, func(i int) bool { return i&0b100 != 0 }},
		{`ctrl @ pow(2) @ gphase(pi/4) q[0];`, func(i int) bool { return i&0b100 != 0 }},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(`include "../testdata/stdgates.qasm"; qubit[3] q; h q;` + c.text)
		if err != nil {
			t.Fatal(err)
		}

		for i, got := range qsim.Amplitude() {
			want := complex(1/math.Sqrt(8), 0)
			if c.active(i) {
				want *= 1i
			}

			if cmplx.Abs(got-want) > 1e-12 {
				t.Errorf("amplitude[%d]=%v, want %v", i, got, want)
			}
		}
	}
}
