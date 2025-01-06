package io_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/itsubaki/qasm/io"
)

func TestMustScan(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "const int a = 42;",
			want:  "const int a = 42;\n",
		},
	}

	for _, c := range cases {
		got := io.MustScan(strings.NewReader(c.input))
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestMustPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec == nil {
			t.Fail()
		}
	}()

	io.Must("", errors.New("something went wrong"))
	t.Fail()
}
