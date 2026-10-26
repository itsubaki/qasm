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
	if err := v.AddWire(svg.Wire{Name: "q"}); err != nil {
		panic(err)
	}

	if err := v.AddWire(svg.Wire{Name: "q"}); err != nil {
		fmt.Println(err)
	}

	// Output:
	// "q" redeclared
}

func TestVisitor_WireIDs(t *testing.T) {
	program, err := xparser.Parse(`
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
	`)
	if err != nil {
		t.Fatal(err)
	}

	v := svg.NewVisitor(environ.New())
	circuit, err := v.Build(program)
	if err != nil {
		t.Fatal(err)
	}

	zero, one := 0, 1
	wantWires := []svg.Wire{
		{Name: "ancilla"},
		{Name: "q", Index: &zero},
		{Name: "q", Index: &one},
		{Name: "r", Index: &zero},
		{Name: "flag"},
		{Name: "c", Index: &zero},
		{Name: "c", Index: &one},
		{Name: "d", Index: &zero},
	}
	if !reflect.DeepEqual(circuit.Wires, wantWires) {
		t.Errorf("wires = %v, want %v", circuit.Wires, wantWires)
	}

	wantOps := []svg.Op{
		&svg.Gate{Name: "H", Target: []svg.WireID{1}},
		&svg.Gate{Name: "X", Control: []svg.WireID{0}, Target: []svg.WireID{2}},
		&svg.Subroutine{Name: "ORACLE", Wire: []svg.WireID{1, 2}},
		&svg.Barrier{Wire: []svg.WireID{1, 2}},
		&svg.Measurement{Wire: []svg.WireID{1}, Target: []svg.WireID{5}},
		&svg.Measurement{Wire: []svg.WireID{2}, Target: []svg.WireID{6}},
		&svg.Measurement{Wire: []svg.WireID{3}, Target: []svg.WireID{7}},
		&svg.Measurement{Wire: []svg.WireID{0}, Target: []svg.WireID{4}},
	}
	if !reflect.DeepEqual(circuit.Ops, wantOps) {
		t.Errorf("ops = %#v, want %#v", circuit.Ops, wantOps)
	}

	cases := []struct {
		name  string
		index []int64
		want  []svg.WireID
	}{
		{name: "ancilla", want: []svg.WireID{0}},
		{name: "q", want: []svg.WireID{1, 2}},
		{name: "q", index: []int64{0}, want: []svg.WireID{1}},
		{name: "r", want: []svg.WireID{3}},
		{name: "r", index: []int64{0}, want: []svg.WireID{3}},
		{name: "d", want: []svg.WireID{7}},
		{name: "q", index: []int64{2}},
		{name: "q", index: []int64{-1}},
		{name: "ancilla", index: []int64{0}},
		{name: "missing"},
	}
	for _, c := range cases {
		got, ok := v.GetWireIDs(c.name, c.index...)
		if ok != (c.want != nil) || !reflect.DeepEqual(got, c.want) {
			t.Errorf("GetWireIDs(%q, %v) = %v, %v; want %v, %v", c.name, c.index, got, ok, c.want, c.want != nil)
		}
	}

	if err := v.AddWire(svg.Wire{Name: "q", Index: new(int)}); err == nil {
		t.Error("expected duplicate q[0] to be rejected")
	}
}

func TestVisitor_Build(t *testing.T) {
	cases := []struct {
		text   string
		hasErr bool
		errMsg string
	}{
		{
			text:   `qubit[3] q; ctrl(2) @ x q[0], q[1], q[2];`,
			hasErr: false,
		},
		{
			text:   `qubit[2] q; cx q[0], q[1]; oracle(q);`,
			hasErr: false,
		},
		{
			text:   `qubit q; oracle(q);`,
			hasErr: false,
		},
		{
			text:   `qubit q; barrier q;`,
			hasErr: false,
		},
		{
			text:   `qubit q; {x a;}`,
			hasErr: false,
		},
		{
			text:   `qubit q; bit c; measure q -> c;`,
			hasErr: false,
		},
		{
			text:   `qubit[2] q; bit[2] c; measure q; measure q -> c; measure q[0] -> c[0];`,
			hasErr: false,
		},
		{
			text:   `qubit[2] q; bit[2] c; bit b;`,
			hasErr: false,
		},
		{
			text:   `qubit q; h q;`,
			hasErr: false,
		},
		{
			text:   `const int n = 3; qubit[n] q;`,
			hasErr: false,
		},
		{
			text:   `int x = 1; int x = 2;`,
			hasErr: true,
			errMsg: `"x" redeclared`,
		},
		{
			text:   `const int x = 1; const int x = 2;`,
			hasErr: true,
			errMsg: `"x" redeclared`,
		},
		{
			text:   `qubit q; x a;`,
			hasErr: true,
			errMsg: `undefined "a"`,
		},
		{
			text:   `qubit[2] q; ctrl(a) @ x q[0], q[1];`,
			hasErr: true,
			errMsg: `undefined "a"`,
		},
		{
			text:   `qubit[2] q; h a[0];`,
			hasErr: true,
			errMsg: `undefined "a"`,
		},
		{
			text:   `qubit[2] q; qubit[2] q;`,
			hasErr: true,
			errMsg: `"q[0]" redeclared`,
		},
		{
			text:   `bit[2] c; bit[2] c;`,
			hasErr: true,
			errMsg: `"c[0]" redeclared`,
		},
		{
			text:   `bit c; bit c;`,
			hasErr: true,
			errMsg: `"c" redeclared`,
		},
		{
			text:   `measure q;`,
			hasErr: true,
			errMsg: `undefined "q"`,
		},
		{
			text:   `qubit[x] q;`,
			hasErr: true,
			errMsg: `undefined "x"`,
		},
		{
			text:   `bit[x] c;`,
			hasErr: true,
			errMsg: `undefined "x"`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			panic(err)
		}

		_, err = svg.Build(program)
		if c.hasErr && err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got = %v, want = %v", err, c.errMsg)
			}

			continue
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func Test_cast(t *testing.T) {
	cases := []struct {
		result any
		want   any
		hasErr bool
	}{
		{
			result: errors.New("something went wrong"),
			want:   nil,
			hasErr: true,
		},
		{
			result: 42,
			want:   42,
			hasErr: false,
		},
		{
			result: "not an int",
			want:   nil,
			hasErr: true,
		},
	}

	for _, c := range cases {
		got, err := svg.CastInt(c.result)
		if c.hasErr && err != nil {
			continue
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if got != c.want {
			t.Errorf("got = %v, want = %v", got, c.want)
		}
	}
}
