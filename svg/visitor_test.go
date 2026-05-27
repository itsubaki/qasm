package svg_test

import (
	"errors"
	"fmt"
	"testing"

	xparser "github.com/itsubaki/qasm/parser"
	"github.com/itsubaki/qasm/svg"
)

func ExampleVisitor_Add() {
	v := svg.NewVisitor()
	if err := v.Add("q"); err != nil {
		panic(err)
	}

	if err := v.Add("q"); err != nil {
		fmt.Println(err)
	}

	// Output:
	// "q" redeclared
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
			text:   `qubit[2] q; oracle(q);`,
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
			text:   `qubit q; x a;`,
			hasErr: true,
			errMsg: `undefined "a"`,
		},
		{
			text:   `qubit[2] q; ctrl(a) @ x q[0], q[1];`,
			hasErr: true,
			errMsg: `unexpected type string`,
		},
		{
			text:   `qubit[2] q; h a[0];`,
			hasErr: true,
			errMsg: `undefined "a[0]"`,
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
			errMsg: `unexpected type string`,
		},
		{
			text:   `bit[x] c;`,
			hasErr: true,
			errMsg: `unexpected type string`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(c.text)
		if err != nil {
			panic(err)
		}

		_, err = svg.NewVisitor().Build(program)
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
