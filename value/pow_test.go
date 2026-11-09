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
		{a: 0, r: 4, want: 0},
		{a: 2, r: 0, want: 1},
		{a: 2, r: 3, want: 8},
	}

	for _, c := range cases {
		got := value.Pow(c.a, c.r)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
