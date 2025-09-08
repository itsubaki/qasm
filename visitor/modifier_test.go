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
	}

	for _, c := range cases {
		if got := visitor.Pow(c.in, c.p); !got.Equal(c.want) {
			t.Fail()
		}
	}
}
