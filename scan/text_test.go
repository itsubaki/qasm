package scan_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/itsubaki/qasm/scan"
)

var ErrSomtingWentWrong = errors.New("something went wrong")

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, ErrSomtingWentWrong
}

func ExampleText() {
	text, err := scan.Text(strings.NewReader("const int a = 42;"))
	if err != nil {
		panic(err)
	}

	fmt.Print(text)

	// Output:
	// const int a = 42;
}

func TestText(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "const int a = 42;",
			want:  "const int a = 42;\n",
		},
		{
			input: "",
			want:  "",
		},
	}

	for _, c := range cases {
		got, err := scan.Text(strings.NewReader(c.input))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if got != c.want {
			t.Errorf("got=%q, want=%q", got, c.want)
		}
	}
}

func TestText_invalid(t *testing.T) {
	cases := []struct {
		reader io.Reader
		want   error
	}{
		{
			reader: &errorReader{},
			want:   ErrSomtingWentWrong,
		},
	}

	for _, c := range cases {
		if _, err := scan.Text(c.reader); !errors.Is(err, c.want) {
			t.Errorf("got=%v, want=%v", err, c.want)
		}
	}
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
		if got := scan.MustText(strings.NewReader(c.input)); got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestMust_invalid(t *testing.T) {
	cases := []struct {
		text string
		err  error
	}{
		{
			text: "",
			err:  ErrSomtingWentWrong,
		},
	}

	for _, c := range cases {
		var got any
		func() {
			defer func() {
				got = recover()
			}()

			scan.Must(c.text, c.err)
		}()

		if got != c.err {
			t.Errorf("panic=%v, want=%v", got, c.err)
		}
	}
}
