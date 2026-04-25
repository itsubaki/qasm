package visitor_test

import (
	"math"
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/qasm/visitor"
)

func TestPow2x2(t *testing.T) {
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
		got := visitor.Pow2x2(c.in, c.p)
		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestPow2x2_compose(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		a, b float64
	}{
		{
			in: gate.X(),
			a:  0.5,
			b:  0.5,
		},
		{
			in: gate.H(),
			a:  0.3,
			b:  0.7,
		},
		{
			in: gate.T(),
			a:  0.2,
			b:  0.4,
		},
		{
			in: gate.U(0.7, 0.3, 1.1),
			a:  0.25,
			b:  0.75,
		},
		{
			in: gate.U(1.2, 2.0, 0.5),
			a:  0.6,
			b:  0.2,
		},
		{
			in: gate.U(1e-10, 0, 0),
			a:  0.5,
			b:  0.5,
		},
		{
			in: gate.U(math.Pi-1e-10, 0, 0),
			a:  0.3,
			b:  0.7,
		},
		{
			in: gate.U(0.9, 0.2, 1.3),
			a:  0.5,
			b:  -0.5,
		},
	}

	for _, c := range cases {
		// U^a * U^b = U^(a+b)
		a := visitor.Pow2x2(c.in, c.a)
		b := visitor.Pow2x2(c.in, c.b)
		ab := visitor.Pow2x2(c.in, c.a+c.b)

		if !a.MatMul(b).Equal(ab) {
			t.Fail()
		}
	}
}

func FuzzPow2x2(f *testing.F) {
	f.Add(0.5, 0.1, 0.2, 0.3)
	f.Add(1.0, 0.0, 0.0, 0.0)
	f.Add(0.3, math.Pi, 0.0, 0.0)

	isValid := func(v ...float64) bool {
		for _, f := range v {
			if math.IsNaN(f) || math.IsInf(f, 0) {
				return false
			}
		}

		return true
	}

	f.Fuzz(func(t *testing.T, p, theta, phi, lambda float64) {
		if !isValid(p, theta, phi, lambda) {
			return
		}

		theta = math.Mod(theta, math.Pi)
		phi = math.Mod(phi, 2*math.Pi)
		lambda = math.Mod(lambda, 2*math.Pi)

		u := gate.U(theta, phi, lambda)
		up := visitor.Pow2x2(u, p)
		back := visitor.Pow2x2(up, 1/p)

		if !back.Equal(u) {
			t.Errorf("p=%v, theta=%v, phi=%v, lambda=%v", p, theta, phi, lambda)
		}
	})
}
