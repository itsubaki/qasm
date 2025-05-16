package visitor_test

import (
	"math/rand/v2"
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/qasm/visitor"
)

func TestControlled(t *testing.T) {
	u := gate.U(rand.Float64(), rand.Float64(), rand.Float64())

	cases := []struct {
		in   *matrix.Matrix
		want *matrix.Matrix
		bit  int
	}{
		{
			in:   gate.TensorProduct(gate.X(), 2, []int{1}),
			want: gate.ControlledNot(2, []int{0}, 1),
			bit:  0,
		},
		{
			in:   gate.TensorProduct(gate.X(), 2, []int{0}),
			want: gate.ControlledNot(2, []int{1}, 0),
			bit:  1,
		},
		{
			in:   gate.ControlledNot(3, []int{0}, 2),
			want: gate.ControlledNot(3, []int{0, 1}, 2),
			bit:  1,
		},
		{
			in:   gate.ControlledNot(3, []int{0}, 2),
			want: gate.ControlledNot(3, []int{0}, 2),
			bit:  0,
		},
		{
			in:   gate.Controlled(u, 3, []int{0}, 1),
			want: gate.Controlled(u, 3, []int{0, 2}, 1),
			bit:  2,
		},
		{
			in:   gate.Controlled(u, 3, []int{0}, 2),
			want: gate.Controlled(u, 3, []int{0, 1}, 2),
			bit:  1,
		},
		{
			in:   gate.Controlled(u, 3, []int{1}, 2),
			want: gate.Controlled(u, 3, []int{0, 1}, 2),
			bit:  0,
		},
		{
			in:   gate.Controlled(u, 3, []int{1}, 0),
			want: gate.Controlled(u, 3, []int{1, 2}, 0),
			bit:  2,
		},
		{
			in:   gate.Controlled(u, 3, []int{2}, 0),
			want: gate.Controlled(u, 3, []int{2, 1}, 0),
			bit:  1,
		},
		{
			in:   gate.Controlled(u, 3, []int{2}, 1),
			want: gate.Controlled(u, 3, []int{0, 2}, 1),
			bit:  0,
		},
	}

	for _, c := range cases {
		got := visitor.Controlled(c.in, []int{c.bit})
		if !got.Equals(c.want) {
			t.Fail()
		}
	}
}

func TestNegControlled(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		want *matrix.Matrix
		bit  int
	}{
		{
			in: gate.TensorProduct(gate.X(), 2, []int{1}),
			want: matrix.Apply(
				gate.TensorProduct(gate.X(), 2, []int{0}),
				gate.ControlledNot(2, []int{0}, 1),
				gate.TensorProduct(gate.X(), 2, []int{0}),
			),
			bit: 0,
		},
		{
			in: gate.TensorProduct(gate.X(), 2, []int{0}),
			want: matrix.Apply(
				gate.TensorProduct(gate.X(), 2, []int{1}),
				gate.ControlledNot(2, []int{1}, 0),
				gate.TensorProduct(gate.X(), 2, []int{1}),
			),
			bit: 1,
		},
		{
			in: gate.ControlledNot(3, []int{0}, 2),
			want: matrix.Apply(
				gate.TensorProduct(gate.X(), 3, []int{1}),
				gate.ControlledNot(3, []int{0, 1}, 2),
				gate.TensorProduct(gate.X(), 3, []int{1}),
			),
			bit: 1,
		},
		{
			in: gate.ControlledNot(3, []int{0}, 1),
			want: matrix.Apply(
				gate.TensorProduct(gate.X(), 3, []int{2}),
				gate.ControlledNot(3, []int{0, 2}, 1),
				gate.TensorProduct(gate.X(), 3, []int{2}),
			),
			bit: 2,
		},
	}

	for _, c := range cases {
		got := visitor.NegControlled(c.in, []int{c.bit})
		if !got.Equals(c.want) {
			t.Fail()
		}
	}
}
