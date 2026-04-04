package value_test

import (
	"testing"

	"github.com/itsubaki/qasm/value"
)

func TestPow(t *testing.T) {
	cases := []struct {
		a, r int
		want int
	}{
		{0, 4, 0},
		{2, 0, 1},
		{2, 3, 8},
	}

	for _, c := range cases {
		got := value.Pow(c.a, c.r)
		if got != c.want {
			t.Errorf("Pow(%v, %v): got=%v, want=%v", c.a, c.r, got, c.want)
		}
	}
}

func TestPowFloat64(t *testing.T) {
	cases := []struct {
		a, r float64
		want float64
	}{
		{0, 4, 0},
		{2, 0, 1},
		{2, 3, 8},
		{9, 0.5, 3},
		{2, -1, 0.5},
	}

	for _, c := range cases {
		got := value.Pow(c.a, c.r)
		if got != c.want {
			t.Errorf("Pow(%v, %v): got=%v, want=%v", c.a, c.r, got, c.want)
		}
	}
}
