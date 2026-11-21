package svg_test

import (
	"encoding/xml"
	"fmt"
	"strings"
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

func TestSVG_gphase(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `gphase(pi/2);`,
			want: "Global phase: gphase(pi/2)",
		},
		{
			text: `gphase(1 & 3); qubit q; h q;`,
			want: "Global phase: gphase(1&amp;3)",
		},
		{
			text: `qubit q; ctrl @ gphase(pi/2) q;`,
			want: ">P</text>",
		},
		{
			text: `qubit[2] q; gphase(pi/4); h q[0]; ctrl(2) @ gphase(pi/2) q[0], q[1]; measure q;`,
			want: "Global phase: gphase(pi/4)",
		},
	}

	for _, c := range cases {
		diagram, err := svg.SVG(c.text, svg.DefaultConfig)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.text, err)
			continue
		}

		if !strings.Contains(diagram, c.want) {
			t.Errorf("%s: missing %q", c.text, c.want)
		}

		var root struct {
			XMLName xml.Name
		}
		if err := xml.Unmarshal([]byte(diagram), &root); err != nil {
			t.Errorf("%s: invalid SVG XML: %v", c.text, err)
		}
	}
}

func TestSVG_measurementAssignment(t *testing.T) {
	cases := []struct {
		text string
		want int
	}{
		{
			text: `qubit q; bit c; c = measure q;`,
			want: 1,
		},
		{
			text: `qubit q; bit c = measure q;`,
			want: 1,
		},
		{
			text: `qubit[2] q; bit[2] c; c = measure q;`,
			want: 2,
		},
		{
			text: `qubit[2] q; bit[2] c = measure q;`,
			want: 2,
		},
	}

	for _, c := range cases {
		diagram, err := svg.SVG(c.text, svg.DefaultConfig)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.text, err)
			continue
		}

		if got := strings.Count(diagram, "<polygon"); got != c.want {
			t.Errorf("%s: measurement arrows=%d, want=%d", c.text, got, c.want)
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

func TestSVG_negativeControl(t *testing.T) {
	cases := []struct {
		name  string
		theme svg.Theme
	}{
		{
			name:  "paper",
			theme: svg.Paper,
		},
		{
			name:  "cyber",
			theme: svg.Cyber,
		},
	}

	for _, c := range cases {
		config := svg.DefaultConfig
		config.Theme = c.theme

		diagram, err := svg.SVG(`qubit[3] q; ctrl @ negctrl @ x q[0], q[1], q[2];`, config)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}

		var root struct {
			Circles []struct {
				Y      int    `xml:"cy,attr"`
				Fill   string `xml:"fill,attr"`
				Stroke string `xml:"stroke,attr"`
			} `xml:"circle"`
		}
		if err := xml.Unmarshal([]byte(diagram), &root); err != nil {
			t.Errorf("%s: invalid SVG XML: %v", c.name, err)
			continue
		}

		if len(root.Circles) != 2 {
			t.Errorf("%s: control circles=%d, want=2", c.name, len(root.Circles))
			continue
		}

		if got := root.Circles[0]; got.Y != config.WireStartY || got.Fill != c.theme.GateStroke {
			t.Errorf("%s: incorrect positive control: %+v", c.name, got)
		}

		if got := root.Circles[1]; got.Y != config.WireStartY+config.WireGap || got.Fill != c.theme.Background || got.Stroke != c.theme.GateStroke {
			t.Errorf("%s: incorrect negative control: %+v", c.name, got)
		}
	}
}
