package svg_test

import (
	"errors"
	"testing"

	"github.com/itsubaki/qasm/svg"
)

func Test_unwrap(t *testing.T) {
	cases := []struct {
		result any
		want   any
		hasErr bool
	}{
		{
			result: errors.New("something went wrong"),
			want:   nil,
			hasErr: true,
		},
		{
			result: 42,
			want:   42,
			hasErr: false,
		},
		{
			result: "not an int",
			want:   nil,
			hasErr: true,
		},
	}

	for _, c := range cases {
		got, err := svg.UnwrapInt(c.result)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("got error = %v", err)
			continue
		}

		if got != c.want {
			t.Errorf("got = %v, want = %v", got, c.want)
		}
	}
}
