package visitor_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/environ"
	"github.com/itsubaki/qasm/parser"
	"github.com/itsubaki/qasm/visitor"
)

func TestWithMaxQubits_enclosed(t *testing.T) {
	cases := []struct {
		name      string
		text      string
		maxQubits int
		depth     int
		want      int
		wantErr   bool
	}{
		{name: "quantum declaration", text: "qubit[2] b;", maxQubits: 2, depth: 1, want: 1, wantErr: true},
		{name: "old style declaration", text: "qreg b[2];", maxQubits: 2, depth: 1, want: 1, wantErr: true},
		{name: "nested scope", text: "qubit[2] b;", maxQubits: 2, depth: 2, want: 1, wantErr: true},
		{name: "nested old style declaration", text: "qreg b[2];", maxQubits: 2, depth: 2, want: 1, wantErr: true},
		{name: "at limit", text: "qubit b;", maxQubits: 2, depth: 1, want: 2},
		{name: "at limit old style declaration", text: "qreg b[1];", maxQubits: 2, depth: 1, want: 2},
		{name: "unlimited", text: "qubit[2] b;", maxQubits: 0, depth: 2, want: 3},
		{name: "unlimited old style declaration", text: "qreg b[2];", maxQubits: 0, depth: 2, want: 3},
	}

	for _, c := range cases {
		qsim := q.New()
		v := visitor.New(qsim, environ.New(), visitor.WithMaxQubits(c.maxQubits))
		parent, err := parser.Parse("qubit a;")
		if err != nil {
			t.Fatal(err)
		}
		if err := v.Run(parent); err != nil {
			t.Fatal(err)
		}

		for range c.depth {
			v = v.Enclosed()
		}
		program, err := parser.Parse(c.text)
		if err != nil {
			t.Fatal(err)
		}
		err = v.Run(program)
		if c.wantErr {
			if !errors.Is(err, visitor.ErrTooManyQubits) {
				t.Errorf("got=%v, want=%v", err, visitor.ErrTooManyQubits)
			}
		} else if err != nil {
			t.Fatal(err)
		}

		if got := qsim.NumQubits(); got != c.want {
			t.Errorf("got=%d qubits, want=%d", got, c.want)
		}
	}
}

func ExampleWithMaxQubits() {
	v := visitor.New(
		q.New(),
		environ.New(),
		visitor.WithMaxQubits(5),
	)

	program, err := parser.Parse(`qubit[10] q;`)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := v.Run(program); err != nil {
		fmt.Println(err)
	}

	// Output:
	// need=10, max=5: too many qubits
}

func ExampleWithMaxQubits_oldstyle() {
	v := visitor.New(
		q.New(),
		environ.New(),
		visitor.WithMaxQubits(5),
	)

	program, err := parser.Parse(`qreg q[10];`)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := v.Run(program); err != nil {
		fmt.Println(err)
	}

	// Output:
	// need=10, max=5: too many qubits
}

func ExampleWithMaxQubits_unlimited() {
	env := environ.New()
	v := visitor.New(
		q.New(),
		env,
		visitor.WithMaxQubits(0),
	)

	program, err := parser.Parse(`qubit[10] q;`)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := v.Run(program); err != nil {
		fmt.Println(err)
	}

	fmt.Println(env.Qubit)

	// Output:
	// map[q:[0 1 2 3 4 5 6 7 8 9]]
}
