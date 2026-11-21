package svg_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/svg"
)

func ExampleSVG() {
	diagram, err := svg.SVG("qubit q; h q;", svg.DefaultConfig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s ... %s", diagram[:4], diagram[len(diagram)-6:])

	// Output:
	// <svg ... </svg>
}

func TestSVG(t *testing.T) {
	cases := []struct {
		text string
	}{
		{
			text: `OPENQASM 3.0; qubit[2] q; h q[0]; ctrl @ x q[0], q[1]; measure q;`,
		},
	}

	for _, c := range cases {
		diagram, err := svg.SVG(c.text, svg.DefaultConfig)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if got := fmt.Sprintf("%s ... %s", diagram[:4], diagram[len(diagram)-6:]); got != "<svg ... </svg>" {
			t.Errorf("got=%q, want=%q", got, "<svg ... </svg>")
		}
	}
}

func TestSVG_invalid(t *testing.T) {
	cases := []struct {
		text string
	}{
		{
			text: `qubit[ q;`,
		},
		{
			text: `qubit q; x a;`,
		},
	}

	for _, c := range cases {
		if _, err := svg.SVG(c.text, svg.DefaultConfig); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}
