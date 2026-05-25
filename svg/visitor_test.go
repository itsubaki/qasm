package svg_test

import (
	"errors"
	"testing"

	"github.com/itsubaki/qasm/svg"
)

func Test_isError(t *testing.T) {
	cases := []struct {
		result any
		hasErr bool
	}{
		{
			result: errors.New("something went wrong"),
			hasErr: true,
		},
		{
			result: "not an error",
			hasErr: false,
		},
	}

	for _, c := range cases {
		if err := svg.IsError(c.result); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("got error = %v", err)
		}
	}
}

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
