package svg_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/svg"
)

func TestSVG(t *testing.T) {
	cases := []struct {
		text   string
		hasErr bool
	}{
		{
			text:   `OPENQASM 3.0; qubit[2] q; h q[0]; ctrl @ x q[0], q[1]; measure q;`,
			hasErr: false,
		},
		{
			text:   `qubit[2] q; {x a;}`,
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
			text:   `qubit[ q;`,
			hasErr: true,
		},
		{
			text:   `qubit q; x a;`,
			hasErr: true,
		},
	}

	for _, c := range cases {
		diagram, err := svg.SVG(c.text, svg.DefaultConfig)
		if c.hasErr && err != nil {
			continue
		}

		if err != nil {
			t.Errorf("got error = %v", err)
			continue
		}

		got := fmt.Sprintf("%s ... %s", diagram[:4], diagram[len(diagram)-6:])
		if got != "<svg ... </svg>" {
			t.Errorf("got=%q, want=%q", got, "<svg ... </svg>")
		}
	}
}
