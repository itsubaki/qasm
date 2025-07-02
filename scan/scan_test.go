package scan_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/itsubaki/qasm/scan"
)

var ErrSomtingWentWrong = errors.New("something went wrong")

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, ErrSomtingWentWrong
}

func TestText(t *testing.T) {
	if _, err := scan.Text(&errorReader{}); err != nil {
		return
	}

	t.Fatal("unexpected")
}

func TestMustText(t *testing.T) {
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
		got := scan.MustText(strings.NewReader(c.input))
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

	scan.Must("", errors.New("something went wrong"))
	t.Fail()
}
