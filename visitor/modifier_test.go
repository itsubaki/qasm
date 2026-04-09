package visitor_test

import (
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/qasm/visitor"
)

func TestPow2x2(t *testing.T) {
	cases := []struct {
		in     *matrix.Matrix
		p      float64
		want   *matrix.Matrix
		errMsg string
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
			in:     gate.New([]complex128{1, 2, 3}),
			p:      2,
			errMsg: "unsupported matrix size 1x3",
		},
	}

	for _, c := range cases {
		got, err := visitor.Pow2x2(c.in, c.p)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
