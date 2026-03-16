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
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
