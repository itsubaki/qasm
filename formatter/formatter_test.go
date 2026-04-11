package formatter_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/formatter"
)

func ExampleFormatter_Format() {
	text := `
	OPENQASM 3.0;include "../testdata/stdgates.qasm";
	qubit q0;qubit[2] q1;qubit[2] q2;
	x q1;h q2;
	`

	formatted, err := formatter.Format(text)
	if err != nil {
		panic(err)
	}

	fmt.Println(formatted)

	// Output:
	// OPENQASM 3.0;
	// include "../testdata/stdgates.qasm";
	// qubit q0;
	// qubit[2] q1;
	// qubit[2] q2;
	// x q1;
	// h q2;
}

func TestFormat(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
	}{
		{
			text: `OPENQASM 3.0; qubit q;`,
			want: `OPENQASM 3.0;
qubit q;
`,
		},
		{
			text:   `qubit[ q;`,
			errMsg: `1:8: mismatched input ';' expecting ']'`,
		},
	}

	for _, c := range cases {
		formatted, err := formatter.Format(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%q, want=%q", err.Error(), c.errMsg)
			}

			continue
		}

		if formatted != c.want {
			t.Errorf("got=%q, want=%q", formatted, c.want)
		}
	}
}
