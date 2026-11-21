package svg_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/itsubaki/qasm/environ"
	xparser "github.com/itsubaki/qasm/parser"
	"github.com/itsubaki/qasm/svg"
)

func ExampleVisitor_AddWire() {
	v := svg.NewVisitor(environ.New())

	if err := v.AddWire(svg.Wire{
		Name: "q",
	}); err != nil {
		panic(err)
	}

	if err := v.AddWire(svg.Wire{
		Name: "q",
	}); err != nil {
		fmt.Println(err)
	}

	// Output:
	// "q" redeclared
}

func TestVisitorBuild_wireIDs(t *testing.T) {
	var zero int
	one := 1
	cases := []struct {
		text      string
		wantWires []svg.Wire
		wantOps   []svg.Op
	}{
		{
			text: `
		qubit ancilla;
		qubit[2] q;
		qubit[1] r;
		bit flag;
		bit[2] c;
		bit[1] d;
		h q[0];
		ctrl @ x ancilla, q[1];
		oracle(q);
		barrier q;
		measure q -> c;
		measure r[0] -> d[0];
		measure ancilla -> flag;
	`,
			wantWires: []svg.Wire{
				{
					Name: "ancilla",
				},
				{
					Name:  "q",
					Index: &zero,
				},
				{
					Name:  "q",
					Index: &one,
				},
				{
					Name:  "r",
					Index: &zero,
				},
				{
					Name: "flag",
				},
				{
					Name:  "c",
					Index: &zero,
				},
				{
					Name:  "c",
					Index: &one,
				},
				{
					Name:  "d",
					Index: &zero,
				},
			},
			wantOps: []svg.Op{
				&svg.Gate{
					Name:   "H",
					Target: []svg.WireID{1},
				},
				&svg.Gate{
					Name:    "X",
					Control: []svg.WireID{0},
					Target:  []svg.WireID{2},
				},
				&svg.Subroutine{
					Name: "ORACLE",
					Wire: []svg.WireID{1, 2},
				},
				&svg.Barrier{
					Wire: []svg.WireID{1, 2},
				},
				&svg.Measurement{
					Wire:   []svg.WireID{1},
					Target: []svg.WireID{5},
				},
				&svg.Measurement{
					Wire:   []svg.WireID{2},
					Target: []svg.WireID{6},
				},
				&svg.Measurement{
					Wire:   []svg.WireID{3},
					Target: []svg.WireID{7},
				},
				&svg.Measurement{
					Wire:   []svg.WireID{0},
					Target: []svg.WireID{4},
				},
			},
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			t.Fatal(err)
		}

		circuit, err := svg.Build(program)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if !reflect.DeepEqual(circuit.Wires, c.wantWires) {
			t.Errorf("wires=%v, want=%v", circuit.Wires, c.wantWires)
		}

		if !reflect.DeepEqual(circuit.Ops, c.wantOps) {
			t.Errorf("ops=%#v, want=%#v", circuit.Ops, c.wantOps)
		}
	}
}

func ExampleBuild_gphase() {
	program, err := xparser.Parse(`qubit q; gphase(pi/4); ctrl @ gphase(pi/2) q;`)
	if err != nil {
		panic(err)
	}

	circuit, err := svg.Build(program)
	if err != nil {
		panic(err)
	}

	fmt.Println(circuit.GlobalPhases)
	gate := circuit.Ops[0].(*svg.Gate)
	fmt.Println(gate.Name, gate.Control, gate.Target)

	// Output:
	// [gphase(pi/4)]
	// P [] [0]
}

func TestVisitorBuild_gphase(t *testing.T) {
	cases := []struct {
		text       string
		wantPhases []string
		wantOps    []svg.Op
	}{
		{
			text:       `gphase(pi/2);`,
			wantPhases: []string{"gphase(pi/2)"},
		},
		{
			text:       `qubit q; h q; gphase(pi/4);`,
			wantPhases: []string{"gphase(pi/4)"},
			wantOps: []svg.Op{
				&svg.Gate{
					Name:   "H",
					Target: []svg.WireID{0},
				},
			},
		},
		{
			text:       `qubit q; inv @ gphase(pi/4); pow(2) @ gphase(pi/8);`,
			wantPhases: []string{"inv@gphase(pi/4)", "pow(2)@gphase(pi/8)"},
		},
		{
			text: `qubit q; ctrl @ gphase(pi/2) q;`,
			wantOps: []svg.Op{
				&svg.Gate{
					Name:   "P",
					Target: []svg.WireID{0},
				},
			},
		},
		{
			text: `qubit[2] q; ctrl(2) @ gphase(pi/2) q[0], q[1];`,
			wantOps: []svg.Op{
				&svg.Gate{
					Name:    "P",
					Control: []svg.WireID{0},
					Target:  []svg.WireID{1},
				},
			},
		},
		{
			text: `qubit[2] q; ctrl @ ctrl @ gphase(pi/2) q[0], q[1];`,
			wantOps: []svg.Op{
				&svg.Gate{
					Name:    "P",
					Control: []svg.WireID{0},
					Target:  []svg.WireID{1},
				},
			},
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			t.Errorf("%s: unexpected parse error: %v", c.text, err)
			continue
		}

		circuit, err := svg.Build(program)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.text, err)
			continue
		}

		if !reflect.DeepEqual(circuit.GlobalPhases, c.wantPhases) {
			t.Errorf("%s: phases=%v, want=%v", c.text, circuit.GlobalPhases, c.wantPhases)
		}

		if !reflect.DeepEqual(circuit.Ops, c.wantOps) {
			t.Errorf("%s: ops=%#v, want=%#v", c.text, circuit.Ops, c.wantOps)
		}
	}
}

func TestVisitorBuild_invalidGphase(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `gphase();`,
			want: "gphase expects one parameter",
		},
		{
			text: `gphase(1, 2);`,
			want: "gphase expects one parameter",
		},
		{
			text: `qubit q; gphase(pi) q;`,
			want: "gphase does not take target operands",
		},
		{
			text: `ctrl @ gphase(pi/2);`,
			want: `not enough operands for "ctrl@"`,
		},
		{
			text: `qubit q; ctrl(2) @ gphase(pi/2) q;`,
			want: `not enough operands for "ctrl(2)@"`,
		},
		{
			text: `qubit q; ctrl(0) @ gphase(pi/2) q;`,
			want: "control count must be positive: 0",
		},
		{
			text: `qubit[2] q; ctrl @ gphase(pi/2) q[0], q[1];`,
			want: "gphase control operand mismatch: controls=1, operands=2",
		},
		{
			text: `negctrl @ gphase(pi/2);`,
			want: `not enough operands for "negctrl@"`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			t.Errorf("%s: unexpected parse error: %v", c.text, err)
			continue
		}

		if _, err := svg.Build(program); err == nil || err.Error() != c.want {
			t.Errorf("%s: got=%v, want=%v", c.text, err, c.want)
		}
	}
}

func ExampleBuild_measurementAssignment() {
	program, err := xparser.Parse(`qubit q; bit c = measure q;`)
	if err != nil {
		panic(err)
	}

	circuit, err := svg.Build(program)
	if err != nil {
		panic(err)
	}

	measurement := circuit.Ops[0].(*svg.Measurement)
	fmt.Println(measurement.Wire, measurement.Target)

	// Output:
	// [0] [1]
}

func TestVisitorBuild_measurementAssignment(t *testing.T) {
	cases := []struct {
		text      string
		reference string
	}{
		{
			text:      `qubit q; bit c; c = measure q;`,
			reference: `qubit q; bit c; measure q -> c;`,
		},
		{
			text:      `qubit q; bit c = measure q;`,
			reference: `qubit q; bit c; measure q -> c;`,
		},
		{
			text:      `qubit[2] q; bit[2] c; c = measure q;`,
			reference: `qubit[2] q; bit[2] c; measure q -> c;`,
		},
		{
			text:      `qubit[2] q; bit[2] c = measure q;`,
			reference: `qubit[2] q; bit[2] c; measure q -> c;`,
		},
		{
			text:      `qubit[2] q; bit[2] c; c[1] = measure q[0];`,
			reference: `qubit[2] q; bit[2] c; measure q[0] -> c[1];`,
		},
		{
			text:      `qubit[1] q; bit[1] c = measure q;`,
			reference: `qubit[1] q; bit[1] c; measure q -> c;`,
		},
		{
			text:      `qubit[2] q; bit flag = measure q[0]; bit[2] c; c[1] = measure q[1];`,
			reference: `qubit[2] q; bit flag; measure q[0] -> flag; bit[2] c; measure q[1] -> c[1];`,
		},
		{
			text:      `qubit q; bit c; h q; c = measure q; x q; bit d = measure q;`,
			reference: `qubit q; bit c; h q; measure q -> c; x q; bit d; measure q -> d;`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			t.Errorf("%s: unexpected parse error: %v", c.text, err)
			continue
		}

		circuit, err := svg.Build(program)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.text, err)
			continue
		}

		reference, err := xparser.Parse(c.reference)
		if err != nil {
			t.Errorf("%s: unexpected reference parse error: %v", c.reference, err)
			continue
		}

		want, err := svg.Build(reference)
		if err != nil {
			t.Errorf("%s: unexpected reference error: %v", c.reference, err)
			continue
		}

		if !reflect.DeepEqual(circuit, want) {
			t.Errorf("%s: circuit=%#v, want=%#v", c.text, circuit, want)
		}
	}
}

func TestVisitorBuild_invalidMeasurementAssignment(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `bit c; c = measure missing;`,
			want: `undefined "missing"`,
		},
		{
			text: `bit c = measure missing;`,
			want: `undefined "missing"`,
		},
		{
			text: `qubit q; c = measure q;`,
			want: `undefined "c"`,
		},
		{
			text: `qubit q; bit[2] c; c[2] = measure q;`,
			want: `undefined "c"`,
		},
		{
			text: `qubit q; bit[2] c; c = measure q;`,
			want: "measurement size mismatch: qubits=1, bits=2",
		},
		{
			text: `qubit[3] q; bit[2] c = measure q;`,
			want: "measurement size mismatch: qubits=3, bits=2",
		},
		{
			text: `qubit[3] q; bit[2] c; measure q -> c;`,
			want: "measurement size mismatch: qubits=3, bits=2",
		},
		{
			text: `qubit[2] q; bit c; c = measure q;`,
			want: "measurement size mismatch: qubits=2, bits=1",
		},
		{
			text: `qubit q; bit c; c ^= measure q;`,
			want: `measurement assignment with operator "^=": not implemented`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			t.Errorf("%s: unexpected parse error: %v", c.text, err)
			continue
		}

		if _, err := svg.Build(program); err == nil || err.Error() != c.want {
			t.Errorf("%s: got=%v, want=%v", c.text, err, c.want)
		}
	}
}

func TestVisitorGetWireIDs(t *testing.T) {
	cases := []struct {
		name  string
		index []int64
		want  []svg.WireID
	}{
		{
			name: "ancilla",
			want: []svg.WireID{0},
		},
		{
			name: "q",
			want: []svg.WireID{1, 2},
		},
		{
			name:  "q",
			index: []int64{0},
			want:  []svg.WireID{1},
		},
		{
			name: "r",
			want: []svg.WireID{3},
		},
		{
			name:  "r",
			index: []int64{0},
			want:  []svg.WireID{3},
		},
		{
			name: "d",
			want: []svg.WireID{7},
		},
		{
			name:  "q",
			index: []int64{2},
		},
		{
			name:  "q",
			index: []int64{-1},
		},
		{
			name:  "ancilla",
			index: []int64{0},
		},
		{
			name: "missing",
		},
	}

	program, err := xparser.Parse(`qubit ancilla; qubit[2] q; qubit[1] r; bit flag; bit[2] c; bit[1] d;`)
	if err != nil {
		t.Fatal(err)
	}

	v := svg.NewVisitor(environ.New())

	if _, err := v.Build(program); err != nil {
		t.Fatal(err)
	}

	for _, c := range cases {
		got, ok := v.GetWireIDs(c.name, c.index...)
		if ok != (c.want != nil) {
			t.Errorf("GetWireIDs(%q, %v): ok=%v, want=%v", c.name, c.index, ok, c.want != nil)
		}

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("GetWireIDs(%q, %v)=%v, want=%v", c.name, c.index, got, c.want)
		}
	}
}

func TestVisitorAddWire_invalid(t *testing.T) {
	var index int
	cases := []struct {
		wire svg.Wire
		want string
	}{
		{
			wire: svg.Wire{
				Name: "q",
			},
			want: `"q" redeclared`,
		},
		{
			wire: svg.Wire{
				Name:  "q",
				Index: &index,
			},
			want: `"q[0]" redeclared`,
		},
	}

	for _, c := range cases {
		v := svg.NewVisitor(environ.New())

		if err := v.AddWire(c.wire); err != nil {
			t.Fatal(err)
		}

		err := v.AddWire(c.wire)
		if err == nil {
			t.Error("expected duplicate wire to be rejected")
			continue
		}

		if err.Error() != c.want {
			t.Errorf("got=%v, want=%v", err, c.want)
		}
	}
}

func TestVisitorBuild(t *testing.T) {
	cases := []struct {
		text string
	}{
		{
			text: `qubit[3] q; ctrl(2) @ x q[0], q[1], q[2];`,
		},
		{
			text: `qubit[2] q; cx q[0], q[1]; oracle(q);`,
		},
		{
			text: `qubit q; oracle(q);`,
		},
		{
			text: `qubit q; barrier q;`,
		},
		{
			text: `qubit q; {x a;}`,
		},
		{
			text: `qubit q; bit c; measure q -> c;`,
		},
		{
			text: `qubit[2] q; bit[2] c; measure q; measure q -> c; measure q[0] -> c[0];`,
		},
		{
			text: `qubit[2] q; bit[2] c; bit b;`,
		},
		{
			text: `qubit q; h q;`,
		},
		{
			text: `const int n = 3; qubit[n] q;`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			panic(err)
		}

		if _, err := svg.Build(program); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestVisitorBuild_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   `int x = 1; int x = 2;`,
			errMsg: `"x" redeclared`,
		},
		{
			text:   `const int x = 1; const int x = 2;`,
			errMsg: `"x" redeclared`,
		},
		{
			text:   `qubit q; x a;`,
			errMsg: `undefined "a"`,
		},
		{
			text:   `qubit[2] q; ctrl(a) @ x q[0], q[1];`,
			errMsg: `undefined "a"`,
		},
		{
			text:   `qubit[2] q; h a[0];`,
			errMsg: `undefined "a"`,
		},
		{
			text:   `qubit[2] q; qubit[2] q;`,
			errMsg: `"q[0]" redeclared`,
		},
		{
			text:   `bit[2] c; bit[2] c;`,
			errMsg: `"c[0]" redeclared`,
		},
		{
			text:   `bit c; bit c;`,
			errMsg: `"c" redeclared`,
		},
		{
			text:   `measure q;`,
			errMsg: `undefined "q"`,
		},
		{
			text:   `qubit[x] q;`,
			errMsg: `undefined "x"`,
		},
		{
			text:   `bit[x] c;`,
			errMsg: `undefined "x"`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			panic(err)
		}

		_, err = svg.Build(program)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestCast(t *testing.T) {
	cases := []struct {
		result any
		want   any
	}{
		{
			result: 42,
			want:   42,
		},
	}

	for _, c := range cases {
		got, err := svg.CastInt(c.result)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if got != c.want {
			t.Errorf("got = %v, want = %v", got, c.want)
		}
	}
}

func TestCast_invalid(t *testing.T) {
	cases := []struct {
		result any
	}{
		{
			result: errors.New("something went wrong"),
		},
		{
			result: "not an int",
		},
	}

	for _, c := range cases {
		if _, err := svg.CastInt(c.result); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func ExampleBuild_negativeControl() {
	program, err := xparser.Parse(`qubit[2] q; negctrl @ x q[0], q[1];`)
	if err != nil {
		panic(err)
	}

	circuit, err := svg.Build(program)
	if err != nil {
		panic(err)
	}

	gate := circuit.Ops[0].(*svg.Gate)
	fmt.Println(gate.Name, gate.Control, gate.NegativeControl, gate.Target)

	// Output:
	// X [] [0] [1]
}

func TestVisitorBuild_negativeControl(t *testing.T) {
	cases := []struct {
		name string
		text string
		want []svg.Op
	}{
		{
			name: "single",
			text: `qubit[2] q; negctrl @ x q[0], q[1];`,
			want: []svg.Op{
				&svg.Gate{
					Name:            "X",
					NegativeControl: []svg.WireID{0},
					Target:          []svg.WireID{1},
				},
			},
		},
		{
			name: "multiple",
			text: `qubit[3] q; negctrl(2) @ x q[0], q[1], q[2];`,
			want: []svg.Op{
				&svg.Gate{
					Name:            "X",
					NegativeControl: []svg.WireID{0, 1},
					Target:          []svg.WireID{2},
				},
			},
		},
		{
			name: "positiveThenNegative",
			text: `qubit[3] q; ctrl @ negctrl @ x q[0], q[1], q[2];`,
			want: []svg.Op{
				&svg.Gate{
					Name:            "X",
					Control:         []svg.WireID{0},
					NegativeControl: []svg.WireID{1},
					Target:          []svg.WireID{2},
				},
			},
		},
		{
			name: "negativeThenPositive",
			text: `qubit[3] q; negctrl @ ctrl @ x q[0], q[1], q[2];`,
			want: []svg.Op{
				&svg.Gate{
					Name:            "X",
					Control:         []svg.WireID{1},
					NegativeControl: []svg.WireID{0},
					Target:          []svg.WireID{2},
				},
			},
		},
		{
			name: "phase",
			text: `qubit q; negctrl @ gphase(pi/2) q;`,
			want: []svg.Op{
				&svg.Gate{
					Name:   "X",
					Target: []svg.WireID{0},
				},
				&svg.Gate{
					Name:   "P",
					Target: []svg.WireID{0},
				},
				&svg.Gate{
					Name:   "X",
					Target: []svg.WireID{0},
				},
			},
		},
		{
			name: "negativePhaseControl",
			text: `qubit[2] q; negctrl @ ctrl @ gphase(pi/2) q[0], q[1];`,
			want: []svg.Op{
				&svg.Gate{
					Name:            "P",
					NegativeControl: []svg.WireID{0},
					Target:          []svg.WireID{1},
				},
			},
		},
		{
			name: "positivePhaseControl",
			text: `qubit[2] q; ctrl @ negctrl @ gphase(pi/2) q[0], q[1];`,
			want: []svg.Op{
				&svg.Gate{
					Name:   "X",
					Target: []svg.WireID{1},
				},
				&svg.Gate{
					Name:    "P",
					Control: []svg.WireID{0},
					Target:  []svg.WireID{1},
				},
				&svg.Gate{
					Name:   "X",
					Target: []svg.WireID{1},
				},
			},
		},
		{
			name: "multiplePhaseControls",
			text: `qubit[2] q; negctrl(2) @ gphase(pi/2) q[0], q[1];`,
			want: []svg.Op{
				&svg.Gate{
					Name:   "X",
					Target: []svg.WireID{1},
				},
				&svg.Gate{
					Name:            "P",
					NegativeControl: []svg.WireID{0},
					Target:          []svg.WireID{1},
				},
				&svg.Gate{
					Name:   "X",
					Target: []svg.WireID{1},
				},
			},
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}

		circuit, err := svg.Build(program)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}

		if !reflect.DeepEqual(circuit.Ops, c.want) {
			t.Errorf("%s: got=%#v, want=%#v", c.name, circuit.Ops, c.want)
		}
	}
}

func TestVisitorBuild_invalidNegativeControl(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `qubit q; negctrl(0) @ x q;`,
			want: "control count must be positive: 0",
		},
		{
			text: `qubit[2] q; negctrl(3) @ x q[0], q[1];`,
			want: `not enough operands for "negctrl(3)@"`,
		},
		{
			text: `qubit[2] q; ctrl @ negctrl(2) @ x q[0], q[1];`,
			want: `not enough operands for "negctrl(2)@"`,
		},
		{
			text: `qubit[2] q; negctrl @ gphase(pi/2) q[0], q[1];`,
			want: "gphase control operand mismatch: controls=1, operands=2",
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := svg.Build(program); err == nil || err.Error() != c.want {
			t.Errorf("%s: got=%v, want=%q", c.text, err, c.want)
		}
	}
}

func ExampleBuild_broadcast() {
	program, err := xparser.Parse(`qubit[2] c; qubit[2] t; ctrl @ x c, t;`)
	if err != nil {
		panic(err)
	}

	circuit, err := svg.Build(program)
	if err != nil {
		panic(err)
	}

	for _, op := range circuit.Ops {
		gate := op.(*svg.Gate)
		fmt.Println(gate.Name, gate.Control, gate.Target)
	}

	// Output:
	// X [0] [2]
	// X [1] [3]
}

func TestVisitorBuild_broadcast(t *testing.T) {
	cases := []struct {
		name      string
		decl      string
		text      string
		reference string
	}{
		{
			name:      "arrayControlAndArrayTarget",
			decl:      `qubit[2] c; qubit[2] t;`,
			text:      `ctrl @ x c, t;`,
			reference: `ctrl @ x c[0], t[0]; ctrl @ x c[1], t[1];`,
		},
		{
			name:      "arrayControlAndScalarTarget",
			decl:      `qubit[2] c; qubit t;`,
			text:      `ctrl @ x c, t;`,
			reference: `ctrl @ x c[0], t; ctrl @ x c[1], t;`,
		},
		{
			name:      "scalarControlAndArrayTarget",
			decl:      `qubit c; qubit[2] t;`,
			text:      `ctrl @ x c, t;`,
			reference: `ctrl @ x c, t[0]; ctrl @ x c, t[1];`,
		},
		{
			name:      "negativeControl",
			decl:      `qubit[2] c; qubit[2] t;`,
			text:      `negctrl @ x c, t;`,
			reference: `negctrl @ x c[0], t[0]; negctrl @ x c[1], t[1];`,
		},
		{
			name:      "multipleControls",
			decl:      `qubit[2] a; qubit[2] b; qubit[2] t;`,
			text:      `ctrl(2) @ x a, b, t;`,
			reference: `ctrl(2) @ x a[0], b[0], t[0]; ctrl(2) @ x a[1], b[1], t[1];`,
		},
		{
			name:      "mixedControls",
			decl:      `qubit[2] a; qubit b; qubit[2] t;`,
			text:      `ctrl @ negctrl @ x a, b, t;`,
			reference: `ctrl @ negctrl @ x a[0], b, t[0]; ctrl @ negctrl @ x a[1], b, t[1];`,
		},
		{
			name:      "singleElementArrays",
			decl:      `qubit[1] c; qubit[1] t;`,
			text:      `ctrl @ x c, t;`,
			reference: `ctrl @ x c[0], t[0];`,
		},
		{
			name:      "indexedSingleElementControl",
			decl:      `qubit[1] c; qubit[2] t;`,
			text:      `ctrl @ x c[0], t;`,
			reference: `ctrl @ x c[0], t[0]; ctrl @ x c[0], t[1];`,
		},
		{
			name:      "indexedSingleElementTarget",
			decl:      `qubit[2] c; qubit[1] t;`,
			text:      `ctrl @ x c, t[0];`,
			reference: `ctrl @ x c[0], t[0]; ctrl @ x c[1], t[0];`,
		},
		{
			name:      "singleQubitGate",
			decl:      `qubit[2] q;`,
			text:      `h q;`,
			reference: `h q[0]; h q[1];`,
		},
		{
			name:      "multipleTargets",
			decl:      `qubit c; qubit[2] a; qubit[2] b;`,
			text:      `ctrl @ swap c, a, b;`,
			reference: `ctrl @ swap c, a[0], b[0]; ctrl @ swap c, a[1], b[1];`,
		},
		{
			name:      "controlledPhase",
			decl:      `qubit[2] q;`,
			text:      `ctrl @ gphase(pi/2) q;`,
			reference: `ctrl @ gphase(pi/2) q[0]; ctrl @ gphase(pi/2) q[1];`,
		},
		{
			name:      "negativeControlledPhase",
			decl:      `qubit[2] q;`,
			text:      `negctrl @ gphase(pi/2) q;`,
			reference: `negctrl @ gphase(pi/2) q[0]; negctrl @ gphase(pi/2) q[1];`,
		},
		{
			name:      "mixedControlledPhase",
			decl:      `qubit[2] c; qubit t;`,
			text:      `ctrl @ negctrl @ gphase(pi/2) c, t;`,
			reference: `ctrl @ negctrl @ gphase(pi/2) c[0], t; ctrl @ negctrl @ gphase(pi/2) c[1], t;`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.decl + c.text)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}

		circuit, err := svg.Build(program)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}

		reference, err := xparser.Parse(c.decl + c.reference)
		if err != nil {
			t.Fatalf("%s: reference: %v", c.name, err)
		}

		want, err := svg.Build(reference)
		if err != nil {
			t.Fatalf("%s: reference: %v", c.name, err)
		}

		if !reflect.DeepEqual(circuit, want) {
			t.Errorf("%s: ops=%#v, want=%#v", c.name, circuit.Ops, want.Ops)
		}
	}
}

func TestVisitorBuild_invalidBroadcast(t *testing.T) {
	cases := []struct {
		name string
		text string
		want string
	}{
		{
			name: "differentArrayLengths",
			text: `qubit[2] c; qubit[3] t; ctrl @ x c, t;`,
			want: "gate operand size mismatch: 2 and 3",
		},
		{
			name: "singleElementControlArray",
			text: `qubit[1] c; qubit[2] t; ctrl @ x c, t;`,
			want: "gate operand size mismatch: 1 and 2",
		},
		{
			name: "singleElementTargetArray",
			text: `qubit[2] c; qubit[1] t; ctrl @ x c, t;`,
			want: "gate operand size mismatch: 2 and 1",
		},
		{
			name: "negativeControl",
			text: `qubit[1] c; qubit[2] t; negctrl @ x c, t;`,
			want: "gate operand size mismatch: 1 and 2",
		},
		{
			name: "multipleControls",
			text: `qubit[1] a; qubit[2] b; qubit t; ctrl(2) @ x a, b, t;`,
			want: "gate operand size mismatch: 1 and 2",
		},
		{
			name: "arrayIsOneOperand",
			text: `qubit[3] q; ctrl(2) @ x q;`,
			want: `not enough operands for "ctrl(2)@"`,
		},
		{
			name: "phaseControlCount",
			text: `qubit[2] q; ctrl(2) @ gphase(pi/2) q;`,
			want: `not enough operands for "ctrl(2)@"`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}

		if _, err := svg.Build(program); err == nil || err.Error() != c.want {
			t.Errorf("%s: got=%v, want=%q", c.name, err, c.want)
		}
	}
}
