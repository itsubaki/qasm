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
