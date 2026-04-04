package visitor_test

import (
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/qasm/visitor"
)

func TestPow(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		p    float64
		want *matrix.Matrix
	}{
		{
			in:   gate.X(),
			p:    1,
			want: gate.X(),
		},
		{
			in:   gate.X(),
			p:    2,
			want: gate.I(),
		},
		{
			in:   gate.T(),
			p:    -1,
			want: gate.T().Dagger(),
		},
		{
			// X^0.5 is the SX (√X) gate: [[(1+i)/2, (1-i)/2], [(1-i)/2, (1+i)/2]]
			in: gate.X(),
			p:  0.5,
			want: matrix.New(
				[]complex128{(1 + 1i) / 2, (1 - 1i) / 2},
				[]complex128{(1 - 1i) / 2, (1 + 1i) / 2},
			),
		},
		{
			// X^0.5 composed with itself must equal X: (X^0.5)^2 = X
			in:   visitor.Pow(gate.X(), 0.5),
			p:    2,
			want: gate.X(),
		},
		{
			// Z^0.5 is the S gate: [[1, 0], [0, i]]
			in: gate.Z(),
			p:  0.5,
			want: matrix.New(
				[]complex128{1, 0},
				[]complex128{0, 1i},
			),
		},
	}

	for _, c := range cases {
		if got := visitor.Pow(c.in, c.p); !got.Equal(c.want) {
			t.Errorf("Pow(%v, %v)\n got: %v\nwant: %v", c.in, c.p, got, c.want)
		}
	}
}
