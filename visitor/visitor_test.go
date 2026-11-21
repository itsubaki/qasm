package visitor_test

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"math/cmplx"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/environ"
	"github.com/itsubaki/qasm/gen/parser"
	xparser "github.com/itsubaki/qasm/parser"
	"github.com/itsubaki/qasm/visitor"
)

func Example_comment() {
	text := `
	// this is a comment
	/* this is a comment block */
	end;
	`

	fmt.Println(visitor.Visit(text))

	// Output:
	// end; <nil>
}

func ExampleVisitor_Run() {
	text := "OPENQASM 3.0;"

	_, env, err := visitor.Run(text)
	if err != nil {
		panic(err)
	}

	fmt.Println(env.Version)

	// Output:
	// 3.0
}

func ExampleVisitor_VisitVersion() {
	text := "OPENQASM 3.0;"

	_, env, err := visitor.Run(text)
	if err != nil {
		panic(err)
	}

	fmt.Println(env.Version)

	// Output:
	// 3.0
}

func ExampleVisitor_VisitIncludeStatement() {
	text := `include "../testdata/stdgates.qasm";`

	_, env, err := visitor.Run(text)
	if err != nil {
		panic(err)
	}

	fmt.Println(slices.Sorted(maps.Keys(env.Gate)))

	// Output:
	// [cx h i x y z]
}

func ExampleVisitor_VisitErrorNode() {
	token := &antlr.BaseToken{}
	token.SetText("something went wrong")

	v := visitor.New(q.New(), environ.New())
	fmt.Println(v.VisitErrorNode(antlr.NewErrorNodeImpl(token)))

	// Output:
	// something went wrong
}

func TestRun_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "const int a = 42; const int a = 42;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "qubit[ q;",
			errMsg: "1:8: mismatched input ';' expecting ']'",
		},
		{
			text:   "qubit[2] q; U(pi/2.0, 0, pi) q[1.1];",
			errMsg: `index must be an integer '1.1'`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisit_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "const int a = 42; const int a = 42;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "qubit[ q;",
			errMsg: "1:8: mismatched input ';' expecting ']'",
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitIncludeStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   `include "../testdata/invalid.qasm"; qubit q; h q;`,
			errMsg: `include ../testdata/invalid.qasm: undefined "invalid"`,
		},
		{
			text:   `include "../testdata/invalid_syntax.qasm";`,
			errMsg: `include ../testdata/invalid_syntax.qasm: 1:8: mismatched input ';' expecting ']'`,
		},
		{
			text:   `include "file_not_found.qasm";`,
			errMsg: `read file file_not_found.qasm: open file_not_found.qasm: no such file or directory`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitConstDeclarationStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "const int a = 42;",
			want: "map[a:42]",
		},
		{
			text: "const uint N = 3 * 5;",
			want: "map[N:15]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", env.Const) != c.want {
			t.Errorf("got=%v, want=%v", env.Const, c.want)
		}
	}
}

func TestVisitorVisitConstDeclarationStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "const int a = 42; const int a = 43;",
			errMsg: `"a" redeclared`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitClassicalDeclarationStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "bool b;",
			want: "map[b:false]",
		},
		{
			text: "bool b = ((1 + 3) * 4 == 16);",
			want: "map[b:true]",
		},
		{
			text: "bit c;",
			want: "map[c:false]",
		},
		{
			text: "bit[1] c;",
			want: "map[c:[false]]",
		},
		{
			text: "bit[4] c;",
			want: "map[c:[false false false false]]",
		},
		{
			text: `bit[8] a = "10001111";`,
			want: "map[a:[true false false false true true true true]]",
		},
		{
			text: `bit[2] c = "10"; c;`,
			want: "map[c:[true false]]",
		},
		{
			text: "bit[1] a = true;",
			want: "map[a:[true]]",
		},
		{
			text: `bit a = "1";`,
			want: "map[a:true]",
		},
		{
			text: `qubit q; U(pi, 0, pi) q; bit c = measure q;`,
			want: "map[c:true]",
		},
		{
			text: `qubit[3] q; U(pi, 0, pi) q[1]; bit c = measure q[1];`,
			want: "map[c:true]",
		},
		{
			text: "int ans = 42;",
			want: "map[ans:42]",
		},
		{
			text: "float ratio = 22 / 7;",
			want: "map[ratio:3]",
		},
		{
			text: "float ratio = 22 / 7.0;",
			want: "map[ratio:3.142857142857143]",
		},
		{
			text: "float f1 = 1.2;",
			want: "map[f1:1.2]",
		},
		{
			text: "float f2 = .1;",
			want: "map[f2:0.1]",
		},
		{
			text: "float f3 = 0.;",
			want: "map[f3:0]",
		},
		{
			text: "float f4 = 3.14e10;",
			want: "map[f4:3.14e+10]",
		},
		{
			text: "float f5 = 2e+1;",
			want: "map[f5:20]",
		},
		{
			text: "float f6 = 2.0E-1;",
			want: "map[f6:0.2]",
		},
		{
			text: "uint ans = 42;",
			want: "map[ans:42]",
		},
		{
			text: "float z;",
			want: "map[z:0]",
		},
		{
			text: "int z;",
			want: "map[z:0]",
		},
		{
			text: "uint z;",
			want: "map[z:0]",
		},
		{
			text: "int hex = 0xffff;",
			want: "map[hex:65535]",
		},
		{
			text: "int hex = 0XBEEF;",
			want: "map[hex:48879]",
		},
		{
			text: "angle[4] a;",
			want: "map[a:0(0,0000)]",
		},
		{
			text: "angle[4] a = pi;",
			want: "map[a:3.141592653589793(8,1000)]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if len(env.ClBit) > 0 && fmt.Sprintf("%v", env.ClBit) != c.want {
			t.Errorf("got=%v, want=%v", env.ClBit, c.want)
		}

		if len(env.ClBitArray) > 0 && fmt.Sprintf("%v", env.ClBitArray) != c.want {
			t.Errorf("got=%v, want=%v", env.ClBitArray, c.want)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitClassicalDeclarationStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   `qubit[3] q; bit c = measure q;`,
			errMsg: "assign 3 bits to a single bit",
		},
		{
			text:   "float a = 1; float a = 0;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "int a = 1; int a = 0;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "uint a = 1; uint a = 0;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "bool a = true; bool a = false;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "bit a; bool a;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "bit[1] a; bool a;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "bit a; bit a;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "int a = 1; bit a;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "bit[1] a; bit a;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "angle[4] a; angle[4] a;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "bit a; angle[4] a;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   "bit[1] a; angle[4] a;",
			errMsg: `"a" redeclared`,
		},
		{
			text:   `bit[8] a = 3.14;`,
			errMsg: `assign 3.14(float64) to "a"`,
		},
		{
			text:   `bit a = 3.14;`,
			errMsg: `assign 3.14(float64) to "a"`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitQuantumDeclarationStatement(t *testing.T) {
	cases := []struct {
		text      string
		want      string
		wantArray string
	}{
		{
			text:      "qubit q;",
			want:      "map[q:0]",
			wantArray: "map[]",
		},
		{
			text:      "qubit[1] q;",
			want:      "map[]",
			wantArray: "map[q:[0]]",
		},
		{
			text:      "qubit[2] q;",
			want:      "map[]",
			wantArray: "map[q:[0 1]]",
		},
		{
			text:      "qubit q0; qubit[2] q1; qubit[3] q3; qubit q4;",
			want:      "map[q0:0 q4:6]",
			wantArray: "map[q1:[1 2] q3:[3 4 5]]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", env.Qubit) != c.want {
			t.Errorf("got=%v, want=%v", env.Qubit, c.want)
		}

		if fmt.Sprintf("%v", env.QubitArray) != c.wantArray {
			t.Errorf("got=%v, want=%v", env.QubitArray, c.wantArray)
		}
	}
}

func TestVisitorVisitQuantumDeclarationStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "qubit[1.1] q;",
			errMsg: `size must be an integer "[1.1]"`,
		},
		{
			text:   "qubit q; qubit q;",
			errMsg: `"q" redeclared`,
		},
		{
			text:   "qubit q; qubit[1] q;",
			errMsg: `"q" redeclared`,
		},
		{
			text:   "qubit[1] q; qubit q;",
			errMsg: `"q" redeclared`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitAliasDeclarationStatement(t *testing.T) {
	cases := []struct {
		text      string
		want      string
		wantArray string
	}{
		{
			text:      "qubit[5] q; let myreg = q[1:4];",
			want:      "map[]",
			wantArray: "map[myreg:[1 2 3] q:[0 1 2 3 4]]",
		},
		{
			text:      "qubit[5] q; let myreg = q[1];",
			want:      "map[myreg:1]",
			wantArray: "map[q:[0 1 2 3 4]]",
		},
		{
			text:      "qubit[2] one; qubit[10] two; let concatenated = one ++ two;",
			want:      "map[]",
			wantArray: "map[concatenated:[0 1 2 3 4 5 6 7 8 9 10 11] one:[0 1] two:[2 3 4 5 6 7 8 9 10 11]]",
		},
		{
			text:      "qubit q; let alias = q;",
			want:      "map[alias:0 q:0]",
			wantArray: "map[]",
		},
		{
			text:      "qubit[1] q; let alias = q;",
			want:      "map[]",
			wantArray: "map[alias:[0] q:[0]]",
		},
		{
			text:      "qubit a; qubit b; let pair = a ++ b;",
			want:      "map[a:0 b:1]",
			wantArray: "map[pair:[0 1]]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", env.Qubit) != c.want {
			t.Errorf("got=%v, want=%v", env.Qubit, c.want)
		}

		if fmt.Sprintf("%v", env.QubitArray) != c.wantArray {
			t.Errorf("got=%v, want=%v", env.QubitArray, c.wantArray)
		}
	}
}

func TestVisitorVisitAliasDeclarationStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "qubit[5] q; let q = q[1:4];",
			errMsg: `"q" redeclared`,
		},
		{
			text:   "qubit q; let q = q;",
			errMsg: `"q" redeclared`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitOldStyleDeclarationStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "qreg q;",
			want: "map[q:[0]]",
		},
		{
			text: "qreg q[2];",
			want: "map[q:[0 1]]",
		},
		{
			text: "creg c;",
			want: "map[c:[false]]",
		},
		{
			text: "creg c[2];",
			want: "map[c:[false false]]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if len(env.QubitArray) > 0 && fmt.Sprintf("%v", env.QubitArray) != c.want {
			t.Errorf("got=%v, want=%v", env.QubitArray, c.want)
		}

		if len(env.ClBit) > 0 && fmt.Sprintf("%v", env.ClBit) != c.want {
			t.Errorf("got=%v, want=%v", env.ClBit, c.want)
		}

		if len(env.ClBitArray) > 0 && fmt.Sprintf("%v", env.ClBitArray) != c.want {
			t.Errorf("got=%v, want=%v", env.ClBitArray, c.want)
		}
	}
}

func TestVisitorVisitOldStyleDeclarationStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "qreg q; qreg q;",
			errMsg: `"q" redeclared`,
		},
		{
			text:   "qubit q; qreg q[1];",
			errMsg: `"q" redeclared`,
		},
		{
			text:   "creg c; creg c;",
			errMsg: `"c" redeclared`,
		},
		{
			text:   "bit c; creg c;",
			errMsg: `"c" redeclared`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitAssignmentStatement(t *testing.T) {
	type Want struct {
		bit      []string
		qubit    []string
		variable []string
	}

	cases := []struct {
		text string
		want Want
	}{
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				gphase(-pi/4);
				ctrl @ U(pi, 0, pi) q[0], q[1];
				ctrl @ gphase(-pi/2) q[0];
				c = measure q;
			`,
			want: Want{
				bit: []string{
					"map[c:[false false]]",
					"map[c:[true true]]",
				},
				qubit: []string{
					"[[00] ( 1.0000 0.0000i): 1.0000]",
					"[[11] ( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				gphase(-pi/4);
				ctrl @ U(pi, 0, pi) q[0], q[1];
				ctrl @ gphase(-pi/2) q[0];
				c[0] = measure q[0];
				c[1] = measure q[1];
			`,
			want: Want{
				bit: []string{
					"map[c:[false false]]",
					"map[c:[true true]]",
				},
				qubit: []string{
					"[[00] ( 1.0000 0.0000i): 1.0000]",
					"[[11] ( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
		{
			text: "bit c; c = true;",
			want: Want{
				bit: []string{"map[c:true]"},
			},
		},
		{
			text: `bit c; c = "1";`,
			want: Want{
				bit: []string{"map[c:true]"},
			},
		},
		{
			text: `qubit[3] q; bit c; U(pi, 0, pi) q[1]; c = measure q[1];`,
			want: Want{
				bit: []string{"map[c:true]"},
			},
		},
		{
			text: "bit[1] c; c = true;",
			want: Want{
				bit: []string{"map[c:[true]]"},
			},
		},
		{
			text: `bit[2] c; c = "10";`,
			want: Want{
				bit: []string{"map[c:[true false]]"},
			},
		},
		{
			text: `bit[2] c; c[0] = "1"; c[1] = "0";`,
			want: Want{
				bit: []string{"map[c:[true false]]"},
			},
		},
		{
			text: `int ans = 42; ans = ans * 2;`,
			want: Want{
				variable: []string{
					"map[ans:84]",
				},
			},
		},
	}

	for _, c := range cases {
		qsim, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if len(c.want.qubit) > 0 {
			var found bool

			for _, w := range c.want.qubit {
				if fmt.Sprintf("%v", qsim.State()) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", qsim.State(), c.want.qubit)
			}
		}

		if len(c.want.bit) > 0 {
			var found bool

			for _, w := range c.want.bit {
				if fmt.Sprintf("%v", env.ClBit) == w {
					found = true
				}

				if fmt.Sprintf("%v", env.ClBitArray) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v/%v, want=%v", env.ClBit, env.ClBitArray, c.want.bit)
			}
		}

		if len(c.want.variable) > 0 {
			var found bool

			for _, w := range c.want.variable {
				if fmt.Sprintf("%v", env.Variable) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", env.Variable, c.want.variable)
			}
		}
	}
}

func TestVisitorVisitAssignmentStatement_invalid(t *testing.T) {
	type Want struct {
		bit      []string
		qubit    []string
		variable []string
	}

	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   `bit c; c = "10";`,
			errMsg: `assign 2 bits to a single bit`,
		},
		{
			text: `
				qubit q;
				U(pi/2.0, 0, pi) q;
				c = measure q;
			`,
			errMsg: `undefined "c"`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitMeasureArrowAssignmentStatement(t *testing.T) {
	type Want struct {
		bit   []string
		qubit []string
	}

	cases := []struct {
		text string
		want Want
	}{
		{
			text: `
				qubit q;
				bit c;
				U(pi, 0, pi) q;
				gphase(-pi/2);
				measure q -> c;
			`,
			want: Want{
				bit:   []string{"map[c:true]"},
				qubit: []string{"[[1] ( 1.0000 0.0000i): 1.0000]"},
			},
		},
		{
			text: `
				qubit[3] q;
				bit c;
				U(pi, 0, pi) q[1];
				gphase(-pi/2);
				measure q[1] -> c;
			`,
			want: Want{
				bit:   []string{"map[c:true]"},
				qubit: []string{"[[010] ( 1.0000 0.0000i): 1.0000]"},
			},
		},
		{
			text: `
				qubit q;
				bit[1] c;
				U(pi, 0, pi) q;
				gphase(-pi/2);
				measure q -> c;
			`,
			want: Want{
				bit:   []string{"map[c:[true]]"},
				qubit: []string{"[[1] ( 1.0000 0.0000i): 1.0000]"},
			},
		},
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				gphase(-pi/4);
				ctrl @ U(pi, 0, pi) q[0], q[1];
				ctrl @ gphase(-pi/2) q[0];
				measure q -> c;
			`,
			want: Want{
				bit: []string{
					"map[c:[false false]]",
					"map[c:[true true]]",
				},
				qubit: []string{
					"[[00] ( 1.0000 0.0000i): 1.0000]",
					"[[11] ( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				gphase(-pi/4);
				ctrl @ U(pi, 0, pi) q[0], q[1];
				ctrl @ gphase(-pi/2) q[0];
				measure q[0] -> c[0];
				measure q[1] -> c[1];
			`,
			want: Want{
				bit: []string{
					"map[c:[false false]]",
					"map[c:[true true]]",
				},
				qubit: []string{
					"[[00] ( 1.0000 0.0000i): 1.0000]",
					"[[11] ( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2.0, 0, pi) q[0];
				gphase(-pi/4);
				ctrl @ U(pi, 0, pi) q[0], q[1];
				ctrl @ gphase(-pi/2) q[0];
				measure q;
			`,
			want: Want{
				qubit: []string{
					"[[00] ( 1.0000 0.0000i): 1.0000]",
					"[[11] ( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
	}

	for _, c := range cases {
		qsim, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if len(c.want.qubit) > 0 {
			var found bool

			for _, w := range c.want.qubit {
				if fmt.Sprintf("%v", qsim.State()) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", qsim.State(), c.want.qubit)
			}
		}

		if len(c.want.bit) > 0 {
			var found bool

			for _, w := range c.want.bit {
				if fmt.Sprintf("%v", env.ClBit) == w {
					found = true
				}

				if fmt.Sprintf("%v", env.ClBitArray) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v/%v, want=%v", env.ClBit, env.ClBitArray, c.want.bit)
			}
		}
	}
}

func TestVisitorVisitMeasureArrowAssignmentStatement_invalid(t *testing.T) {
	type Want struct {
		bit   []string
		qubit []string
	}

	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text: `
				qubit[3] q;
				bit c;
				measure q -> c;
			`,
			errMsg: "assign 3 bits to a single bit",
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitResetStatement(t *testing.T) {
	cases := []struct {
		text string
		want []string
	}{
		{
			text: `
				qubit q;
				U(pi/2, 0, pi) q;
				gphase(-pi/4);
				reset q;
			`,
			want: []string{
				"[0] ( 1.0000 0.0000i): 1.0000",
			},
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		for i, s := range qsim.State() {
			if s.String() != c.want[i] {
				t.Fatalf("got=%v, want=%v", s.String(), c.want[i])
			}
		}
	}
}

func TestVisitorVisitResetStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "int a = 1; reset a;",
			errMsg: `undefined "a"`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitAdditiveExpression(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "1 + 3;",
			want: "4",
		},
		{
			text: "1.1 + 3;",
			want: "4.1",
		},
		{
			text: "1.0 - 3;",
			want: "-2",
		},
		{
			text: "1 - 3.1;",
			want: "-2.1",
		},
		{
			text: "1.3 - 3.1;",
			want: "-1.8",
		},
		{
			text: "int(1.5)+int(2.5);",
			want: "3",
		},
		{
			text: "int(1.5)-int(2.5);",
			want: "-1",
		},
		{
			text: "int(2.5)+1;",
			want: "3",
		},
		{
			text: "int(2.5)+1.5;",
			want: "3.5",
		},
		{
			text: "1+int(2.5);",
			want: "3",
		},
		{
			text: "1.5+int(2.5);",
			want: "3.5",
		},
		{
			text: "int(2.5)-1;",
			want: "1",
		},
		{
			text: "int(2.5)-1.5;",
			want: "0.5",
		},
		{
			text: "1-int(2.5);",
			want: "-1",
		},
		{
			text: "1.5-int(2.5);",
			want: "-0.5",
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitAdditiveExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "9.0 + foobar(1);",
			errMsg: `undefined "foobar"`,
		},
		{
			text:   "foobar(1) + 9.0;",
			errMsg: `undefined "foobar"`,
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitMultiplicativeExpression(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "1 * 3;",
			want: "3",
		},
		{
			text: "1.0 * 3;",
			want: "3",
		},
		{
			text: "2 * 3.1;",
			want: "6.2",
		},
		{
			text: "2.5 * 3.1;",
			want: "7.75",
		},
		{
			text: "4 / 2;",
			want: "2",
		},
		{
			text: "10 % 3;",
			want: "1",
		},
		{
			text: "int(1) * int(3);",
			want: "3",
		},
		{
			text: "int(9) / int(3);",
			want: "3",
		},
		{
			text: "int(9) % int(3);",
			want: "0",
		},
		{
			text: "int(1) * 3;",
			want: "3",
		},
		{
			text: "int(9) / 3;",
			want: "3",
		},
		{
			text: "int(9) % 3;",
			want: "0",
		},
		{
			text: "int(1) * 3.0;",
			want: "3",
		},
		{
			text: "int(9) / 3.0;",
			want: "3",
		},
		{
			text: "9 % int(3);",
			want: "0",
		},
		{
			text: "1 * int(3);",
			want: "3",
		},
		{
			text: "9 / int(3);",
			want: "3",
		},
		{
			text: "1.0 * int(3);",
			want: "3",
		},
		{
			text: "9.0 / int(3);",
			want: "3",
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitMultiplicativeExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "9.0 / foobar(1);",
			errMsg: `undefined "foobar"`,
		},
		{
			text:   "foobar(1) / 9.0;",
			errMsg: `undefined "foobar"`,
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitParenthesisExpression(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "(1 + 3) * 4;",
			want: "16",
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Fail()
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitCallExpression(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "sin(1.0);",
			want: "0.8414709848078965",
		},
		{
			text: "cos(1.0);",
			want: "0.5403023058681398",
		},
		{
			text: "tan(1.0);",
			want: "1.557407724654902",
		},
		{
			text: "arcsin(0.5);",
			want: "0.5235987755982989",
		},
		{
			text: "arccos(0.5);",
			want: "1.0471975511965976",
		},
		{
			text: "arctan(1.0);",
			want: "0.7853981633974483",
		},
		{
			text: "ceiling(1.1);",
			want: "2",
		},
		{
			text: "floor(1.1);",
			want: "1",
		},
		{
			text: "sqrt(2.0);",
			want: "1.4142135623730951",
		},
		{
			text: "exp(1.0);",
			want: "2.718281828459045",
		},
		{
			text: "log(2.0);",
			want: "0.6931471805599453",
		},
		{
			text: "mod(10.0, 3.0);",
			want: "1",
		},
		{
			text: "sin(0);",
			want: "0",
		},
		{
			text: "cos(0);",
			want: "1",
		},
		{
			text: "tan(0);",
			want: "0",
		},
		{
			text: "arcsin(0);",
			want: "0",
		},
		{
			text: "arccos(1);",
			want: "0",
		},
		{
			text: "arctan(0);",
			want: "0",
		},
		{
			text: "ceiling(1);",
			want: "1",
		},
		{
			text: "floor(1);",
			want: "1",
		},
		{
			text: "sqrt(4);",
			want: "2",
		},
		{
			text: "exp(0);",
			want: "1",
		},
		{
			text: "log(1);",
			want: "0",
		},
		{
			text: "mod(5,2);",
			want: "1",
		},
		{
			text: "sqrt(uint(4));",
			want: "2",
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("%s: got=%v, want=%v", c.text, result, c.want)
		}
	}
}

func TestVisitorVisitCallExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "sqrt(true);",
			errMsg: "unexpected type: bool",
		},
		{
			text:   "mod(true,2);",
			errMsg: "unexpected type: bool",
		},
		{
			text:   "mod(5,true);",
			errMsg: "unexpected type: bool",
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitPowerExpression(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "2**3;",
			want: "8",
		},
		{
			text: "2**0.5;",
			want: "1.4142135623730951",
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitPowerExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "foobar**0.5;",
			errMsg: `undefined "foobar"`,
		},
		{
			text:   "2**foobar;",
			errMsg: `undefined "foobar"`,
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitLogicalOrExpression(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{
			text: "false || false;",
			want: false,
		},
		{
			text: "true || false;",
			want: true,
		},
		{
			text: "false || true;",
			want: true,
		},
		{
			text: "true || true;",
			want: true,
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Fail()
		}

		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitLogicalAndExpression(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{
			text: "false && false;",
			want: false,
		},
		{
			text: "true && false;",
			want: false,
		},
		{
			text: "false && true;",
			want: false,
		},
		{
			text: "true && true;",
			want: true,
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Fail()
		}

		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitBitwiseAndExpression(t *testing.T) {
	cases := []struct {
		text string
		want int64
	}{
		{
			text: "10 & 12;",
			want: 8,
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Fail()
		}

		if result.(int64) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitBitwiseOrExpression(t *testing.T) {
	cases := []struct {
		text string
		want int64
	}{
		{
			text: "10 | 12;",
			want: 14,
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Fail()
		}

		if result.(int64) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitBitwiseXorExpression(t *testing.T) {
	cases := []struct {
		text string
		want int64
	}{
		{
			text: "10 ^ 12;",
			want: 6,
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Fail()
		}

		if result.(int64) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitBitshiftExpression(t *testing.T) {
	cases := []struct {
		text string
		want int64
	}{
		{
			text: "11 << 2;",
			want: 44,
		},
		{
			text: "11 >> 1;",
			want: 5,
		},
		{
			text: "11 >> 1 << 1;",
			want: 10,
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Fail()
		}

		if result.(int64) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitEqualityExpression(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{
			text: "10.0 == 2.0 * 5;",
			want: true,
		},
		{
			text: "10.0 != 2.0 * 5;",
			want: false,
		},
		{
			text: "5 != 3;",
			want: true,
		},
		{
			text: "5 == 3.0;",
			want: false,
		},
		{
			text: "5 != 3.0;",
			want: true,
		},
		{
			text: "5.0 == 3;",
			want: false,
		},
		{
			text: "5.0 != 3;",
			want: true,
		},
		{
			text: "(1 == 1) == true;",
			want: true,
		},
		{
			text: "(1 == 1) != true;",
			want: false,
		},
		{
			text: "int(5) == int(3);",
			want: false,
		},
		{
			text: "int(5) != int(3);",
			want: true,
		},
		{
			text: "int(5) == 3;",
			want: false,
		},
		{
			text: "int(5) != 3;",
			want: true,
		},
		{
			text: "int(5) == 3.5;",
			want: false,
		},
		{
			text: "int(5) != 3.5;",
			want: true,
		},
		{
			text: "5 == int(3);",
			want: false,
		},
		{
			text: "5 != int(3);",
			want: true,
		},
		{
			text: "float(5.5) == int(3);",
			want: false,
		},
		{
			text: "float(5.5) != int(3);",
			want: true,
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitEqualityExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "foobar == 3;",
			errMsg: `undefined "foobar"`,
		},
		{
			text:   "3 == foobar;",
			errMsg: `undefined "foobar"`,
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitUnaryExpression(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{
			text: "!false;",
			want: true,
		},
		{
			text: "!true;",
			want: false,
		},
		{
			text: "-1 == 1;",
			want: false,
		},
		{
			text: "-1.0 != 1.0;",
			want: true,
		},
		{
			text: "!(-1.0 != 1.0);",
			want: false,
		},
		{
			text: "~5 == -6;",
			want: true,
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitUnaryExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "-foobar;",
			errMsg: `undefined "foobar"`,
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitComparisonExpression(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "2.0 < 3;",
			want: "true",
		},
		{
			text: "2 <= 2;",
			want: "true",
		},
		{
			text: "2.0 > 3;",
			want: "false",
		},
		{
			text: "2 >= 3;",
			want: "false",
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitComparisonExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "foobar > 3;",
			errMsg: `undefined "foobar"`,
		},
		{
			text:   "3 < foobar;",
			errMsg: `undefined "foobar"`,
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitGateCallStatement(t *testing.T) {
	cases := []struct {
		text string
		want []string
	}{
		{
			text: "qubit q; U(pi, 0, pi) q; gphase(-pi/2);",
			want: []string{
				"[1] ( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `qubit q; U(pi/2, 0, pi) q; gphase(-pi/4);`,
			want: []string{
				"[0] ( 0.7071 0.0000i): 0.5000",
				"[1] ( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `qubit q; gphase(-π/2);`,
			want: []string{
				"[0] ( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				gphase(-pi/4);
				ctrl @ U(pi, 0, pi) q[0], q[1];
				ctrl @ gphase(-pi/2) q[0];
			`,
			want: []string{
				"[00] ( 0.7071 0.0000i): 0.5000",
				"[11] ( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				inv @ U(pi/2, 0, pi) q[0];
			`,
			want: []string{
				"[00] ( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit[2] q;
				pow(2) @ U(pi/2, 0, pi) q[0];
				gphase(-pi/2);
			`,
			want: []string{
				"[00] ( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit[2] q;
				pow(0.5) @ U(pi, 0, pi) q[0];
				gphase(-pi/4);
				pow(0.5) @ U(pi, 0, pi) q[0];
				gphase(-pi/4);
			`,
			want: []string{
				"[10] ( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `qubit[2] q; pow(3) @ U(pi/2, 0, pi) q[0]; gphase(-3*pi/4);`,
			want: []string{
				"[00] ( 0.7071 0.0000i): 0.5000",
				"[10] ( 0.7071 0.0000i): 0.5000",
			},
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		for i, s := range qsim.State() {
			if s.String() == c.want[i] {
				continue
			}

			t.Errorf("got=%v, want=%v", s.String(), c.want[i])
		}
	}
}

func TestVisitorVisitGateCallStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   `qubit q; gphase(true);`,
			errMsg: "invalid param 'bool(true)'",
		},
		{
			text:   `qubit q; gphase(a);`,
			errMsg: `undefined "a"`,
		},
		{
			text:   `qubit[2] q; U(true, 0, pi) q;`,
			errMsg: "invalid param 'bool(true)'",
		},
		{
			text: `
				int a = 1;
				qubit t;
				ctrl @ U(pi, 0, pi) a, t;
			`,
			errMsg: `undefined "a"`,
		},
		{
			text: `
				qubit c;
				qubit t;
				ctrl(true) @ U(pi, 0, pi) c, t;
			`,
			errMsg: `apply "ctrl(true)@"`,
		},
		{
			text: `
				qubit c;
				qubit t;
				negctrl(true) @ U(pi, 0, pi) c, t;
			`,
			errMsg: `apply "negctrl(true)@"`,
		},
		{
			text:   `U(pi, 0, pi) q;`,
			errMsg: `undefined "q"`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitGateCallStatement_userDefined(t *testing.T) {
	cases := []struct {
		text string
		want []string
	}{
		{
			text: `
				gate h q0 { U(pi/2, 0, pi) q0; gphase(-pi/4); }
				gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; ctrl @ gphase(-pi/2) q0; }
				qubit[2] q;
				h q[0];
				cx q[0], q[1];
			`,
			want: []string{
				"[00] ( 0.7071 0.0000i): 0.5000",
				"[11] ( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				gate x q0 { U(pi, 0, pi) q0; gphase(-pi/2); }
				gate y q0 { U(pi, pi/2, pi/2) q0; gphase(-pi/2); }
				qubit q;
				x q;
				y q;
			`,
			want: []string{
				"[0] ( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				gate x q { U(pi, 0, pi) q; gphase(-pi/2); }
				gate y q { U(pi, pi/2, pi/2) q; gphase(-pi/2); }
				gate xy q { x q; y q; }
				qubit q;
				xy q;
			`,
			want: []string{
				"[0] ( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; gphase(-p0/2); }
				qubit q;
				u(pi, 0, pi) q;
			`,
			want: []string{
				"[1] ( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; gphase(-p0/2); }
				qubit[2] q;
				u(pi, 0, pi) q;
			`,
			want: []string{
				"[11] ( 1.0000 0.0000i): 1.0000",
			},
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		for i, s := range qsim.State() {
			if s.String() == c.want[i] {
				continue
			}

			t.Errorf("got=%v, want=%v", s.String(), c.want[i])
		}
	}
}

func TestVisitorVisitGateCallStatement_invalidUserDefined(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text: `
				gate u q { }
				gate u q { }
			`,
			errMsg: `"u" redeclared`,
		},
		{
			text:   `qubit[2] q; myg(pi, 0, pi) q;`,
			errMsg: `undefined "myg"`,
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				qubit q;
				u(a, 0, pi) q;
			`,
			errMsg: `undefined "a"`,
		},
		{
			text: `
				gate myg q { }
				int a = 1;
				myg a;
			`,
			errMsg: `undefined "a"`,
		},
		{
			text: `
				gate myg q { myx q; }
				qubit q;
				myg q;
			`,
			errMsg: `gate call[0]: undefined "myx"`,
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				qubit q;
				inv @ u(pi, 0, pi) q;
			`,
			errMsg: "user-defined call with modifier: not implemented",
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitGateModifier(t *testing.T) {
	cases := []struct {
		text string
		want []string
	}{
		{
			text: `
				qubit q0;
				qubit q1;
				U(pi/2.0, 0, pi) q0;
				gphase(-pi/4);
				ctrl @ U(pi, 0, pi) q0, q1;
				ctrl @ gphase(-pi/2) q0;
			`,
			want: []string{
				"[00] ( 0.7071 0.0000i): 0.5000",
				"[11] ( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				qubit[3] q;
				U(pi/2.0, 0, pi) q[0], q[1];
				gphase(-pi/2);
				ctrl @ ctrl @ U(pi, 0, pi) q[0], q[1], q[2];
				ctrl @ ctrl @ gphase(-pi/2) q[0], q[1];
			`,
			want: []string{
				"[000] ( 0.5000 0.0000i): 0.2500",
				"[010] ( 0.5000 0.0000i): 0.2500",
				"[100] ( 0.5000 0.0000i): 0.2500",
				"[111] ( 0.5000 0.0000i): 0.2500",
			},
		},
		{
			text: `
				qubit[3] q;
				U(pi/2.0, 0, pi) q[0], q[1];
				gphase(-pi/2);
				ctrl(2) @ U(pi, 0, pi) q[0], q[1], q[2];
				ctrl(2) @ gphase(-pi/2) q[0], q[1];
			`,
			want: []string{
				"[000] ( 0.5000 0.0000i): 0.2500",
				"[010] ( 0.5000 0.0000i): 0.2500",
				"[100] ( 0.5000 0.0000i): 0.2500",
				"[111] ( 0.5000 0.0000i): 0.2500",
			},
		},
		{
			text: `
				qubit q0;
				qubit q1;
				U(pi/2.0, 0, pi) q0;
				gphase(-pi/4);
				negctrl @ U(pi, 0, pi) q0, q1;
				negctrl @ gphase(-pi/2) q0;
			`,
			want: []string{
				"[01] ( 0.7071 0.0000i): 0.5000",
				"[10] ( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				qubit q;
				U(pi, tau, euler) q;
				inv @ U(pi, tau, euler) q;
			`,
			want: []string{
				"[0] ( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				const float half = pi / 2;
				qubit q;
				pow(2) @ U(half, -half, half) q;
				gphase(-half);
			`,
			want: []string{
				"[1] ( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit q;
				inv @ pow(2) @ U(pi, 0, pi) q;
				gphase(pi);
			`,
			want: []string{
				"[0] ( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `qubit q; inv @ pow(0.5) @ U(0.3, 1.1, 2.7) q; gphase(0.3/4);`,
			want: []string{
				"[0] ( 0.9958 0.0040i): 0.9915",
				"[1] ( 0.0905-0.0164i): 0.0085",
			},
		},
		{
			text: `qubit q; pow(0.5) @ inv @ U(0.3, 1.1, 2.7) q; gphase(0.3/4);`,
			want: []string{
				"[0] ( 0.9958 0.0040i): 0.9915",
				"[1] ( 0.0905-0.0164i): 0.0085",
			},
		},
		{
			text: `
				qubit[3] q;
				U(pi/2.0, 0, pi) q[0], q[1];
				gphase(-pi/2);
				ctrl @ negctrl @ U(pi, 0, pi) q[0], q[1], q[2];
				ctrl @ negctrl @ gphase(-pi/2) q[0], q[1];
			`,
			want: []string{
				"[000] ( 0.5000 0.0000i): 0.2500",
				"[010] ( 0.5000 0.0000i): 0.2500",
				"[101] ( 0.5000 0.0000i): 0.2500",
				"[110] ( 0.5000 0.0000i): 0.2500",
			},
		},
		{
			// RX(pi/2)^2 = -iX
			text: `
				qubit q;
				pow(2.0) @ U(pi/2, -pi/2, pi/2) q;
				gphase(-pi/2);
			`,
			want: []string{
				"[1] ( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit q;
				U(pi/2, -pi/2, pi/2) q;
				gphase(-pi/4);
			`,
			want: []string{
				"[0] ( 0.7071 0.0000i): 0.5000",
				"[1] ( 0.0000-0.7071i): 0.5000",
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				gphase(-pi/4);
				inv @ ctrl @ U(pi, 0, pi) q[0], q[1];
				ctrl @ gphase(pi/2) q[0];
			`,
			want: []string{
				"[00] ( 0.7071 0.0000i): 0.5000",
				"[11] ( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				gphase(-pi/4);
				ctrl @ inv @ U(pi, 0, pi) q[0], q[1];
				ctrl @ gphase(pi/2) q[0];
			`,
			want: []string{
				"[00] ( 0.7071 0.0000i): 0.5000",
				"[11] ( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				gphase(-pi/4);
				ctrl @ pow(2) @ U(pi/2, -pi/2, pi/2) q[0], q[1];
				ctrl @ gphase(-pi/2) q[0];
			`,
			want: []string{
				"[00] ( 0.7071 0.0000i): 0.5000",
				"[11] ( 0.0000-0.7071i): 0.5000",
			},
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		for i, s := range qsim.State() {
			if s.String() == c.want[i] {
				continue
			}

			t.Errorf("got=%v, want=%v", s.String(), c.want[i])
		}
	}
}

func TestVisitorVisitGateModifier_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text: `
				qubit q;
				pow(true) @ U(pi, 0, pi) q;
			`,
			errMsg: `apply "pow(true)@": unexpected type: bool`,
		},
		{
			// not implemented
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				pow(2) @ ctrl @ U(pi/2, -pi/2, pi/2) q[0], q[1];
			`,
			errMsg: "pow applied after ctrl: not implemented",
		},
		{
			// not implemented
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				pow(2) @ negctrl @ U(pi/2, -pi/2, pi/2) q[0], q[1];
			`,
			errMsg: "pow applied after negctrl: not implemented",
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitDefStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `
				def xm(qubit q1) { U(pi, 0, pi) q1; measure q1; }
				qubit q;
				xm(q);
			`,
			want: "map[]",
		},
		{
			text: `
				def xm(qubit q1) -> bit { U(pi, 0, pi) q1; return measure q1; }
				qubit q;
				bit c = xm(q);
			`,
			want: "map[c:true]",
		},
		{
			text: `
				def xm(qubit q1) -> bit { U(pi, 0, pi) q1; bit m = measure q1; return m;}
				qubit q;
				bit c = xm(q);
			`,
			want: "map[c:true]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if len(env.ClBit) > 0 && fmt.Sprintf("%v", env.ClBit) != c.want {
			t.Errorf("got=%v, want=%v", env.ClBit, c.want)
		}

		if len(env.ClBitArray) > 0 && fmt.Sprintf("%v", env.ClBitArray) != c.want {
			t.Errorf("got=%v, want=%v", env.ClBitArray, c.want)
		}
	}
}

func TestVisitorVisitDefStatement_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text: `
				qubit q;
				xm(q);
			`,
			errMsg: `undefined "xm"`,
		},
		{
			text:   "def f(qubit q) -> bit { return 1; } def f(qubit q) -> bit { return 0; }",
			errMsg: `"f" redeclared`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitIfStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `
				int a = 10;
				if (a == 10) { a = 20; }
				int b = 30;
			`,
			want: "map[a:20 b:30]",
		},
		{
			text: `
				int a = 10;
				if (a == 20) { a = 100; } else { a = 30; }
			`,
			want: "map[a:30]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Fail()
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitForStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `
				int a = 0;
				for int i in [0:9] {
					a = a + 1;
				}
			`,
			want: "map[a:10]",
		},
		{
			text: `
				int a = 0;
				for int i in [0:9] {
					a = a + i;
				}
			`,
			want: "map[a:45]",
		},
		{
			text: `
				int a = 0;
				for int i in [(1-1):9] {
					a = a + i;
				}
			`,
			want: "map[a:45]",
		},
		{
			text: `
				int R = 10;
				int a = 0;
				for int i in [0:R-1] {
					a = a + i;
				}
			`,
			want: "map[R:10 a:45]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Fail()
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitBreakStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `
				int a = 0;
				for int i in [0:10] {
					a = a + 1;
					if ( a > 10 ) {
						break;
					}
					a = a + 1;
				}
			`,
			want: "map[a:11]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Fail()
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitContinueStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `
				int a = 0;
				for int i in [0:9] {
					a = a + 1;
					if ( a > 10 ) {
						continue;
					}
					a = a + 1;
				}
			`,
			want: "map[a:15]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Fail()
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitWhileStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `
				int a = 0;
				while (a < 10) {
					a = a + 1;
				}
			`,
			want: "map[a:10]",
		},
		{
			text: `
				int a = 0;
				while (a < 100) {
					a = a + 1;
					if ( a > 10 ) {
						break;
					}
					a = a + 10;
				}
			`,
			want: "map[a:12]",
		},
		{
			text: `
				int a = 0;
				while (a < 100) {
					a = a + 1;
					if ( a < 10 ) {
						continue;
					}
					a = a + 10;
				}
			`,
			want: "map[a:108]",
		},
		{
			text: `
				int a = 0;
				for int i in [0:10] {
					{
						if (a > 2) {
							break;
						}
						a = a + 1;
					}
					a = a + 100;
				}
			`,
			want: "map[a:101]",
		},
		{
			text: `
				int a = 0;
				for int i in [0:10] {
					{
						if (a > 2) {
							continue;
						}
						a = a + 1;
					}
					a = a + 100;
				}
			`,
			want: "map[a:101]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Fail()
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitSwitchStatement(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `
				int a = 15;
				int b = 0;
				switch (a) {
					case 15 {
						b = 15;
					}
					default {
						b = -1;
					}
				}
			`,
			want: "map[a:15 b:15]",
		},
		{
			text: `
				int a = 20;
				int b = 0;
				switch (a) {
					case 15 {
						b = 15;
					}
					default {
						b = -1;
					}
				}
			`,
			want: "map[a:20 b:-1]",
		},
		{
			text: `
				int a = 20;
				int b = 0;
				switch (a) {
					case 1, 2, 3{
						b = 15;
					}
					case 20 {
						b = -1;
					}
				}
			`,
			want: "map[a:20 b:-1]",
		},
		{
			text: `
				int a = 20;
				int b = 0;
				switch (a) { }
			`,
			want: "map[a:20 b:0]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Fail()
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitCastExpression(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "int a = 42; float b = float(a); float c = b + 0.1;",
			want: "map[a:42 b:42 c:42.1]",
		},
		{
			text: "int a = 42; float b = float(a+0.1);",
			want: "map[a:42 b:42.1]",
		},
		{
			text: "int a = int(42);",
			want: "map[a:42]",
		},
		{
			text: "int a = int(42.123);",
			want: "map[a:42]",
		},
		{
			text: "uint a = uint(42);",
			want: "map[a:42]",
		},
		{
			text: "uint a = uint(42.123);",
			want: "map[a:42]",
		},
		{
			text: "int[32] a = 42;",
			want: "map[a:42]",
		},
		{
			text: "int[64] a = 42;",
			want: "map[a:42]",
		},
		{
			text: "int[64] a = 42;",
			want: "map[a:42]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitCastExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "bool a = true; int b = int(a);",
			errMsg: "int(true): unexpected type: bool",
		},
		{
			text:   "bool a = true; uint b = uint(a);",
			errMsg: "uint(true): unexpected type: bool",
		},
		{
			text:   "bool a = true; float b = float(a);",
			errMsg: "float64(true): unexpected type: bool",
		},
		{
			text:   "bool a = bool(1);",
			errMsg: `unsupported scalar type "bool"`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitIndexedIdentifier_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text: `
				qubit[2] q;
				U(pi, 0, pi) q[1.1];
			`,
			errMsg: `index must be an integer '1.1'`,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))

		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		id := p.Program().
			StatementOrScope(1).
			Statement().
			GateCallStatement().
			GateOperandList().
			GateOperand(0).
			IndexedIdentifier()

		ctx, ok := id.(*parser.IndexedIdentifierContext)
		if !ok {
			t.Fatalf("indexed identifier type=%T", id)
		}

		result := visitor.New(q.New(), environ.New()).VisitIndexedIdentifier(ctx)
		err, ok := result.(error)
		if !ok {
			t.Errorf("got=%v (%T), want error", result, result)
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitIndexExpression(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: `bit[2] m = "10"; m[0];`,
			want: "true",
		},
		{
			text: `bit[2] m = "10"; m[1];`,
			want: "false",
		},
		{
			text: `bit[2] m = "10"; int i = 0; m[i + 1];`,
			want: "false",
		},
		{
			text: `bit[2] m = "10"; m[0] != m[1];`,
			want: "true",
		},
		{
			text: `bit[2] m = "11"; m[0] != m[1];`,
			want: "false",
		},
		{
			text: `bit[2] m = "10"; bit s = m[0]; s;`,
			want: "true",
		},
		{
			text: `OPENQASM 3.0; qubit[2] q; bit[2] m; bit s;
				m[0] = measure q[0]; m[1] = measure q[1];
				s = m[0] != m[1]; s;`,
			want: "false",
		},
		{
			text: `OPENQASM 3.0; qubit[2] q; bit[2] m; bit s;
				U(pi, 0, pi) q[0];
				m[0] = measure q[0]; m[1] = measure q[1];
				s = m[0] != m[1]; s;`,
			want: "true",
		},
		{
			text: `qubit[7] q; bit[7] m; bit s1;
				U(pi, 0, pi) q[0], q[2], q[4]; m = measure q;
				s1 = ((m[0] != m[2]) != m[4]) != m[6]; s1;`,
			want: "true",
		},
		{
			text: `bit[3] m = "101"; m[2, 0];`,
			want: "[true true]",
		},
		{
			text: `bit[3] m = "101"; m[0:2];`,
			want: "[true false]",
		},
		{
			text: `bit[3] m = "101"; m[0:3][1];`,
			want: "false",
		},
		{
			text: `bit[3] m = "101"; m[1:1];`,
			want: "[]",
		},
	}

	for _, c := range cases {
		result, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitorVisitIndexExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "bit[2] m; m[-1];",
			errMsg: "bit index -1 out of range for size 2",
		},
		{
			text:   "bit[2] m; m[2];",
			errMsg: "bit index 2 out of range for size 2",
		},
		{
			text:   "bit[2] m; m[1.1];",
			errMsg: "index must be an integer '1.1'",
		},
		{
			text:   "bit[2] m; m[i];",
			errMsg: `undefined "i"`,
		},
		{
			text:   "bit[0] m; m[0];",
			errMsg: "bit index 0 out of range for size 0",
		},
		{
			text:   "bit[2] m; m[-1:1];",
			errMsg: `bit range "-1:1" out of range for size 2`,
		},
		{
			text:   "bit[2] m; m[1:0];",
			errMsg: `bit range "1:0" out of range for size 2`,
		},
		{
			text:   "bit[2] m; m[0:3];",
			errMsg: `bit range "0:3" out of range for size 2`,
		},
		{
			text:   "bit[2] m; m[0:1:2];",
			errMsg: `bit range "0:1:2": not implemented`,
		},
		{
			text:   "bit[2] m; m[0:i];",
			errMsg: `undefined "i"`,
		},
		{
			text:   "int a = 1; a[0];",
			errMsg: `invalid operand "a"`,
		},
		{
			text:   "a[0];",
			errMsg: `undefined "a"`,
		},
	}

	for _, c := range cases {
		_, err := visitor.Visit(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitScalarType(t *testing.T) {
	cases := []struct {
		text string
		want int64
	}{
		{
			text: "angle a = 0;",
			want: 1,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))

		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		scalar := p.Program().
			StatementOrScope(0).
			Statement().
			ClassicalDeclarationStatement().
			ScalarType()

		ctx, ok := scalar.(*parser.ScalarTypeContext)
		if !ok {
			t.Fatalf("scalar type=%T", scalar)
		}

		result := visitor.New(q.New(), environ.New()).VisitScalarType(ctx)
		if err, ok := result.(error); ok {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		got, ok := result.(int64)
		if !ok {
			t.Errorf("got=%T, want int64", result)
		}

		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestVisitorVisitScalarType_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "angle[1.1] a;",
			errMsg: `size must be an integer "[1.1]"`,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))

		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		scalar := p.Program().
			StatementOrScope(0).
			Statement().
			ClassicalDeclarationStatement().
			ScalarType()

		ctx, ok := scalar.(*parser.ScalarTypeContext)
		if !ok {
			t.Fatalf("scalar type=%T", scalar)
		}

		result := visitor.New(q.New(), environ.New()).VisitScalarType(ctx)
		err, ok := result.(error)
		if !ok {
			t.Errorf("got=%v (%T), want error", result, result)
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitQubitType(t *testing.T) {
	cases := []struct {
		text string
		want int64
	}{
		{
			text: "qubit q;",
			want: 1,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))

		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		qb := p.Program().
			StatementOrScope(0).
			Statement().
			QuantumDeclarationStatement().
			QubitType()

		ctx, ok := qb.(*parser.QubitTypeContext)
		if !ok {
			t.Fatalf("qubit type=%T", qb)
		}

		result := visitor.New(q.New(), environ.New()).VisitQubitType(ctx)
		if err, ok := result.(error); ok {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		got, ok := result.(int64)
		if !ok {
			t.Errorf("got=%T, want int64", result)
		}

		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestVisitorVisitQubitType_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "qubit[1.1] q;",
			errMsg: `size must be an integer "[1.1]"`,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))

		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		qb := p.Program().
			StatementOrScope(0).
			Statement().
			QuantumDeclarationStatement().
			QubitType()

		ctx, ok := qb.(*parser.QubitTypeContext)
		if !ok {
			t.Fatalf("qubit type=%T", qb)
		}

		result := visitor.New(q.New(), environ.New()).VisitQubitType(ctx)
		err, ok := result.(error)
		if !ok {
			t.Errorf("got=%v (%T), want error", result, result)
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitRangeExpression_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "for int i in [true:9] { }",
			errMsg: "int64(true): unexpected type: bool",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))

		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		x := p.Program().
			StatementOrScope(0).
			Statement().
			ForStatement().
			RangeExpression()

		ctx, ok := x.(*parser.RangeExpressionContext)
		if !ok {
			t.Fatalf("range expression type=%T", x)
		}

		result := visitor.New(q.New(), environ.New()).VisitRangeExpression(ctx)
		err, ok := result.(error)
		if !ok {
			t.Errorf("got=%v (%T), want error", result, result)
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitArrayType(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "array[int[8], 5] aa;",
			want: "map[aa:[0 0 0 0 0]]",
		},
		{
			text: "array[int[16], 5] aa;",
			want: "map[aa:[0 0 0 0 0]]",
		},
		{
			text: "array[int[32], 5] aa;",
			want: "map[aa:[0 0 0 0 0]]",
		},
		{
			text: "array[int[64], 5] aa;",
			want: "map[aa:[0 0 0 0 0]]",
		},
		{
			text: "array[uint[8], 5] aa;",
			want: "map[aa:[0 0 0 0 0]]",
		},
		{
			text: "array[uint[16], 5] aa;",
			want: "map[aa:[0 0 0 0 0]]",
		},
		{
			text: "array[uint[32], 5] aa;",
			want: "map[aa:[0 0 0 0 0]]",
		},
		{
			text: "array[uint[64], 3] aa;",
			want: "map[aa:[0 0 0]]",
		},
		{
			text: "array[float[32], 3] aa;",
			want: "map[aa:[0 0 0]]",
		},
		{
			text: "array[float[64], 3] aa;",
			want: "map[aa:[0 0 0]]",
		},
		{
			text: "array[bool, 3] aa;",
			want: "map[aa:[false false false]]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitArrayType_invalid(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "array[int[11], 5] aa;",
			errMsg: "unsupported bit size 11",
		},
		{
			text:   "array[uint[11], 5] aa;",
			errMsg: "unsupported bit size 11",
		},
		{
			text:   "array[float[11], 5] aa;",
			errMsg: "unsupported bit size 11",
		},
		{
			text:   "array[angle[4], 3] aa;",
			errMsg: `unsupported scalar type "angle[4]"`,
		},
		{
			text:   "array[int[64], 5] aa;array[float[64], 5] aa;",
			errMsg: `"aa" redeclared`,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Run(c.text)
		if err == nil {
			t.Error("expected an error")
			continue
		}

		if err.Error() != c.errMsg {
			t.Errorf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitArrayLiteral(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{
			text: "array[int[8], 3] aa = {0, 1, 2};",
			want: "map[aa:[0 1 2]]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("unexpected error: %v", err)

			continue
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitorVisitMeasureExpression(t *testing.T) {
	cases := []struct {
		text string
		want any
	}{
		{
			text: `
				qubit q;
				U(pi, 0, pi) q;
				measure q;
			`,
			want: true,
		},
		{
			text: `
				qubit[2] q;
				U(pi, 0, pi) q[0];
				U(pi, 0, pi) q[1];
				measure q;
			`,
			want: []bool{true, true},
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
		statements := p.Program().AllStatementOrScope()

		v := visitor.New(q.New(), environ.New())

		for _, s := range statements[:len(statements)-1] {
			if err := v.Run(s); err != nil {
				t.Fatalf("got=%v, want no error", err)
			}
		}

		x := v.Visit(statements[len(statements)-1].
			Statement().
			MeasureArrowAssignmentStatement().
			MeasureExpression(),
		)

		switch want := c.want.(type) {
		case bool:
			bit, ok := x.(bool)
			if !ok {
				t.Fatalf("got=%T, want bool", x)
			}

			if bit != want {
				t.Fatalf("got=%v, want=%v", bit, want)
			}
		case []bool:
			bits, ok := x.([]bool)
			if !ok {
				t.Fatalf("got=%T, want []bool", x)
			}

			if fmt.Sprintf("%v", bits) != fmt.Sprintf("%v", want) {
				t.Fatalf("got=%v, want=%v", bits, want)
			}
		}
	}
}

func TestVisitorVisitMeasureExpression_invalidAssignment(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "int a = 1; bit c = measure a;",
			errMsg: `undefined "a"`,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
		statements := p.Program().AllStatementOrScope()

		v := visitor.New(q.New(), environ.New())

		if err := v.Run(statements[0]); err != nil {
			t.Fatalf("got=%v, want no error", err)
		}

		x := v.Visit(statements[1].
			Statement().
			ClassicalDeclarationStatement().
			DeclarationExpression().
			MeasureExpression(),
		)
		err, ok := x.(error)
		if !ok {
			t.Fatalf("got=%T, want error", x)
		}

		if err.Error() != c.errMsg {
			t.Fatalf("got=%v, want=%v", err, c.errMsg)
		}
	}
}

func TestVisitorVisitEndStatement(t *testing.T) {
	cases := []struct {
		name string
		text string
		want string
	}{
		{
			name: "direct",
			text: `stop();`,
			want: "map[touched:0]",
		},
		{
			name: "initializer",
			text: `int result=stop();`,
			want: "map[touched:0]",
		},
		{
			name: "assignment",
			text: `int result=0; result=stop();`,
			want: "map[result:0 touched:0]",
		},
		{
			name: "return",
			text: `def f() -> int { return stop(); } int result=f();`,
			want: "map[touched:0]",
		},
		{
			name: "argument",
			text: `mod(stop(),touch(q));`,
			want: "map[touched:0]",
		},
		{
			name: "gate argument",
			text: `U(stop(),0,0) q;`,
			want: "map[touched:0]",
		},
		{
			name: "loop condition",
			text: `while (stop()==0) { break; }`,
			want: "map[touched:0]",
		},
	}

	for _, c := range cases {
		qsim, env, err := visitor.Run(`
			int touched=0; qubit q;
			def stop() -> int { end; }
			def touch(qubit arg) -> int { U(pi,0,pi) arg; return 1; }
		` + c.text + ` touched=1; U(pi,0,pi) q;`)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}

		if got := fmt.Sprint(env.Variable); got != c.want {
			t.Errorf("%s: variables=%s, want=%s", c.name, got, c.want)
		}

		if !qsim.Qubit().IsZero() {
			t.Errorf("%s: executed gate after end: %v", c.name, qsim.Qubit())
		}
	}
}

func TestVisitor_invalidUnimplemented(t *testing.T) {
	cases := []struct {
		name  string
		visit func(*visitor.Visitor) any
	}{
		{
			name: "VisitPragma",
			visit: func(v *visitor.Visitor) any {
				return v.VisitPragma(&parser.PragmaContext{})
			},
		},
		{
			name: "VisitAnnotation",
			visit: func(v *visitor.Visitor) any {
				return v.VisitAnnotation(&parser.AnnotationContext{})
			},
		},
		{
			name: "VisitDurationofExpression",
			visit: func(v *visitor.Visitor) any {
				return v.VisitDurationofExpression(&parser.DurationofExpressionContext{})
			},
		},
		{
			name: "VisitSetExpression",
			visit: func(v *visitor.Visitor) any {
				return v.VisitSetExpression(&parser.SetExpressionContext{})
			},
		},
		{
			name: "VisitArrayReferenceType",
			visit: func(v *visitor.Visitor) any {
				return v.VisitArrayReferenceType(&parser.ArrayReferenceTypeContext{})
			},
		},
		{
			name: "VisitDefcalArgumentDefinitionList",
			visit: func(v *visitor.Visitor) any {
				return v.VisitDefcalArgumentDefinitionList(&parser.DefcalArgumentDefinitionListContext{})
			},
		},
		{
			name: "VisitDefcalArgumentDefinition",
			visit: func(v *visitor.Visitor) any {
				return v.VisitDefcalArgumentDefinition(&parser.DefcalArgumentDefinitionContext{})
			},
		},
		{
			name: "VisitDefcalTarget",
			visit: func(v *visitor.Visitor) any {
				return v.VisitDefcalTarget(&parser.DefcalTargetContext{})
			},
		},
		{
			name: "VisitDefcalOperandList",
			visit: func(v *visitor.Visitor) any {
				return v.VisitDefcalOperandList(&parser.DefcalOperandListContext{})
			},
		},
		{
			name: "VisitDefcalOperand",
			visit: func(v *visitor.Visitor) any {
				return v.VisitDefcalOperand(&parser.DefcalOperandContext{})
			},
		},
		{
			name: "VisitIoDeclarationStatement",
			visit: func(v *visitor.Visitor) any {
				return v.VisitIoDeclarationStatement(&parser.IoDeclarationStatementContext{})
			},
		},
		{
			name: "VisitExternStatement",
			visit: func(v *visitor.Visitor) any {
				return v.VisitExternStatement(&parser.ExternStatementContext{})
			},
		},
		{
			name: "VisitExternArgumentList",
			visit: func(v *visitor.Visitor) any {
				return v.VisitExternArgumentList(&parser.ExternArgumentListContext{})
			},
		},
		{
			name: "VisitExternArgument",
			visit: func(v *visitor.Visitor) any {
				return v.VisitExternArgument(&parser.ExternArgumentContext{})
			},
		},
		{
			name: "VisitCalStatement",
			visit: func(v *visitor.Visitor) any {
				return v.VisitCalStatement(&parser.CalStatementContext{})
			},
		},
		{
			name: "VisitDefcalStatement",
			visit: func(v *visitor.Visitor) any {
				return v.VisitDefcalStatement(&parser.DefcalStatementContext{})
			},
		},
		{
			name: "VisitCalibrationGrammarStatement",
			visit: func(v *visitor.Visitor) any {
				return v.VisitCalibrationGrammarStatement(&parser.CalibrationGrammarStatementContext{})
			},
		},
		{
			name: "VisitBoxStatement",
			visit: func(v *visitor.Visitor) any {
				return v.VisitBoxStatement(&parser.BoxStatementContext{})
			},
		},
		{
			name: "VisitDelayStatement",
			visit: func(v *visitor.Visitor) any {
				return v.VisitDelayStatement(&parser.DelayStatementContext{})
			},
		},
	}

	for _, c := range cases {
		v := visitor.New(q.New(), environ.New())

		result := c.visit(v)
		err, ok := result.(error)
		if !ok {
			t.Errorf("%s: got=%v (%T), want error", c.name, result, result)
			continue
		}

		if !errors.Is(err, visitor.ErrNotImplemented) {
			t.Fatalf("got=%v", err)
		}

		if !strings.Contains(err.Error(), c.name) {
			t.Fatalf("got=%v", err)
		}
	}
}

func TestVisitor_invalidNested(t *testing.T) {
	scopes := []struct {
		name string
		text string
	}{
		{
			name: "scope",
			text: `{ %s }`,
		},
		{
			name: "nested scope",
			text: `{ { %s } touched = 1; }`,
		},
		{
			name: "if",
			text: `if (true) { %s }`,
		},
		{
			name: "else",
			text: `if (false) {} else { %s }`,
		},
		{
			name: "for",
			text: `for int i in [0:1] { %s }`,
		},
		{
			name: "while",
			text: `while (touched < 1) { %s }`,
		},
		{
			name: "switch case",
			text: `switch (1) { case 1 { %s } }`,
		},
		{
			name: "switch default",
			text: `switch (1) { default { %s } }`,
		},
		{
			name: "function",
			text: `def f(qubit arg) { %s } f(q);`,
		},
		{
			name: "function without arguments",
			text: `def f() { %s } f();`,
		},
		{
			name: "function initializer",
			text: `def f(qubit arg) -> int { %s return 1; } int result = f(q);`,
		},
		{
			name: "function initializer without arguments",
			text: `def f() -> int { %s return 1; } int result = f();`,
		},
		{
			name: "nested function",
			text: `def f(qubit arg) { %s } def g(qubit arg) { f(arg); touched = 1; } g(q);`,
		},
	}

	cases := []struct {
		name      string
		statement string
		wantErr   string
		wantLimit bool
	}{
		{
			name:      "undefined operand",
			statement: `reset missing;`,
			wantErr:   `undefined "missing"`,
		},
		{
			name:      "measurement",
			statement: `measure missing;`,
			wantErr:   `undefined "missing"`,
		},
		{
			name:      "measurement assignment",
			statement: `bit result; measure missing -> result;`,
			wantErr:   `undefined "missing"`,
		},
		{
			name:      "redeclared variable",
			statement: `int a = 1; int a = 2;`,
			wantErr:   `"a" redeclared`,
		},
		{
			name:      "qubit limit",
			statement: `qubit extra;`,
			wantLimit: true,
		},
		{
			name:      "old style qubit limit",
			statement: `qreg extra[1];`,
			wantLimit: true,
		},
	}

	for _, scope := range scopes {
		for _, c := range cases {
			name := scope.name + "/" + c.name
			text := `int touched = 0; qubit q; ` + fmt.Sprintf(scope.text, c.statement+` touched = 1;`) + ` touched = 2;`

			program, err := xparser.Parse(text)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}

			qsim, env := q.New(), environ.New()
			v := visitor.New(
				qsim,
				env,
				visitor.WithMaxQubits(1),
			)

			err = v.Run(program)

			switch {
			case c.wantLimit:
				if !errors.Is(err, visitor.ErrTooManyQubits) {
					t.Errorf("%s: got=%v, want=%v", name, err, visitor.ErrTooManyQubits)
				}
			case err == nil || err.Error() != c.wantErr:
				t.Errorf("%s: got=%v, want=%s", name, err, c.wantErr)
			}

			if got := env.Variable["touched"]; got != int64(0) {
				t.Errorf("%s: executed statements after error: touched=%v", name, got)
			}

			if got := qsim.NumQubits(); got != 1 {
				t.Errorf("%s: got=%d qubits, want=1", name, got)
			}
		}
	}
}

func TestVisitor_invalidFunction(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{
			name: "assignment",
			text: `int result = 0; result = fail(q);`,
		},
		{
			name: "constant",
			text: `const int result = fail(q);`,
		},
		{
			name: "integer initializer",
			text: `int result = fail(q);`,
		},
		{
			name: "unsigned initializer",
			text: `uint result = fail(q);`,
		},
		{
			name: "float initializer",
			text: `float result = fail(q);`,
		},
		{
			name: "bool initializer",
			text: `bool result = fail(q);`,
		},
		{
			name: "angle initializer",
			text: `angle result = fail(q);`,
		},
		{
			name: "bit initializer",
			text: `bit result = fail(q);`,
		},
		{
			name: "bit array initializer",
			text: `bit[1] result = fail(q);`,
		},
		{
			name: "array initializer",
			text: `array[int[64], 2] result = {fail(q), touch(q)};`,
		},
		{
			name: "if condition",
			text: `if (fail(q)) { touched = 1; }`,
		},
		{
			name: "while condition",
			text: `while (fail(q)) { break; }`,
		},
		{
			name: "for range",
			text: `for int i in [0:fail(q)] { break; }`,
		},
		{
			name: "switch expression",
			text: `switch (fail(q)) { default { touched = 1; } }`,
		},
		{
			name: "switch case expression",
			text: `switch (1) { case fail(q), touch(q) { touched = 1; } }`,
		},
		{
			name: "builtin argument",
			text: `mod(fail(q), touch(q));`,
		},
		{
			name: "function argument",
			text: `def g(qubit arg) { touched = 1; } g(fail(q));`,
		},
		{
			name: "gate argument",
			text: `U(fail(q), touch(q), 0) q;`,
		},
		{
			name: "gate power",
			text: `pow(fail(q)) @ U(0, 0, 0) q;`,
		},
		{
			name: "gate control",
			text: `ctrl(fail(q)) @ U(0, 0, 0) q;`,
		},
		{
			name: "gate negative control",
			text: `negctrl(fail(q)) @ U(0, 0, 0) q;`,
		},
		{
			name: "cast",
			text: `int(fail(q));`,
		},
		{
			name: "logical and",
			text: `fail(q) && touch(q);`,
		},
		{
			name: "logical and right operand",
			text: `true && fail(q);`,
		},
		{
			name: "logical or",
			text: `fail(q) || touch(q);`,
		},
		{
			name: "logical or right operand",
			text: `false || fail(q);`,
		},
		{
			name: "bitwise and",
			text: `fail(q) & touch(q);`,
		},
		{
			name: "bitwise and right operand",
			text: `1 & fail(q);`,
		},
		{
			name: "bitwise or",
			text: `fail(q) | touch(q);`,
		},
		{
			name: "bitwise or right operand",
			text: `1 | fail(q);`,
		},
		{
			name: "bitwise xor",
			text: `fail(q) ^ touch(q);`,
		},
		{
			name: "bitwise xor right operand",
			text: `1 ^ fail(q);`,
		},
		{
			name: "left shift",
			text: `fail(q) << touch(q);`,
		},
		{
			name: "left shift right operand",
			text: `1 << fail(q);`,
		},
		{
			name: "right shift",
			text: `fail(q) >> touch(q);`,
		},
		{
			name: "right shift right operand",
			text: `1 >> fail(q);`,
		},
		{
			name: "index expression",
			text: `q[fail(q)];`,
		},
		{
			name: "index operand",
			text: `fail(q)[0];`,
		},
		{
			name: "index range",
			text: `q[0:fail(q)];`,
		},
		{
			name: "gate operand index",
			text: `reset q[fail(q)];`,
		},
		{
			name: "assignment index",
			text: `bit[1] result; result[fail(q)] = touch(q);`,
		},
		{
			name: "measurement index",
			text: `bit[1] result; measure q -> result[fail(q)];`,
		},
		{
			name: "qubit size",
			text: `qubit[fail(q)] result;`,
		},
		{
			name: "old style qubit size",
			text: `qreg result[fail(q)];`,
		},
		{
			name: "old style bit size",
			text: `creg result[fail(q)];`,
		},
		{
			name: "bit size",
			text: `bit[fail(q)] result;`,
		},
		{
			name: "angle size",
			text: `angle[fail(q)] result;`,
		},
		{
			name: "array size",
			text: `array[int[64], fail(q)] result;`,
		},
		{
			name: "array element size",
			text: `array[int[fail(q)], 1] result;`,
		},
		{
			name: "unsigned array element size",
			text: `array[uint[fail(q)], 1] result;`,
		},
		{
			name: "float array element size",
			text: `array[float[fail(q)], 1] result;`,
		},
		{
			name: "alias",
			text: `let result = fail(q) ++ touch(q);`,
		},
	}

	for _, c := range cases {
		text := `
			int touched = 0;
			qubit q;
			def fail(qubit arg) -> int { qubit extra; touched = 1; return 1; }
			def touch(qubit arg) -> int { touched = 2; return 1; }
		` + c.text + ` touched = 3;`
		program, err := xparser.Parse(text)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}

		qsim, env := q.New(), environ.New()
		v := visitor.New(
			qsim,
			env,
			visitor.WithMaxQubits(1),
		)

		if err := v.Run(program); !errors.Is(err, visitor.ErrTooManyQubits) {
			t.Errorf("%s: got=%v, want=%v", c.name, err, visitor.ErrTooManyQubits)
		}

		if got := env.Variable["touched"]; got != int64(0) {
			t.Errorf("%s: executed statements after error: touched=%v", c.name, got)
		}

		for id, val := range env.Variable {
			if err, ok := val.(error); ok {
				t.Errorf("%s: stored error in variable %q: %v", c.name, id, err)
			}
		}

		if got := qsim.NumQubits(); got != 1 {
			t.Errorf("%s: got=%d qubits, want=1", c.name, got)
		}
	}
}

func TestVisitor_errorShortCircuit(t *testing.T) {
	cases := []struct {
		text string
	}{
		{
			text: `if (false) { fail(q); }`,
		},
		{
			text: `if (true) { touched = 0; } else { fail(q); }`,
		},
		{
			text: `while (false) { fail(q); }`,
		},
		{
			text: `for int i in [1:0] { fail(q); }`,
		},
		{
			text: `switch (1) { case 0 { fail(q); } default { touched = 0; } }`,
		},
		{
			text: `bool result = false && fail(q);`,
		},
		{
			text: `bool result = true || fail(q);`,
		},
		{
			text: `def empty(qubit arg) {} empty(q);`,
		},
		{
			text: `def empty() {} empty();`,
		},
	}

	for _, c := range cases {
		program, err := xparser.Parse(`
			int touched = 0;
			qubit q;
			def fail(qubit arg) -> bit { qubit extra; touched = 1; return true; }
		` + c.text)
		if err != nil {
			t.Fatal(err)
		}

		qsim, env := q.New(), environ.New()
		v := visitor.New(
			qsim,
			env,
			visitor.WithMaxQubits(1),
		)

		if err := v.Run(program); err != nil {
			t.Errorf("%s: unexpected error: %v", c.text, err)
		}

		if got := env.Variable["touched"]; got != int64(0) {
			t.Errorf("%s: evaluated skipped branch: touched=%v", c.text, got)
		}

		if got := qsim.NumQubits(); got != 1 {
			t.Errorf("%s: got=%d qubits, want=1", c.text, got)
		}
	}
}

func TestVisitor_programContinuesAfterScope(t *testing.T) {
	cases := []struct {
		name string
		body string
		want int64
	}{
		{
			name: "if",
			body: `if (true) { n = 1; n = 2; }`,
			want: 12,
		},
		{
			name: "else",
			body: `if (false) {} else { n = 1; n = 2; }`,
			want: 12,
		},
		{
			name: "empty if",
			body: `if (true) {}`,
			want: 10,
		},
		{
			name: "scope",
			body: `{ n = 1; n = 2; }`,
			want: 12,
		},
		{
			name: "nested scope",
			body: `{ if (true) { n = 1; n = 2; } n = 3; }`,
			want: 13,
		},
		{
			name: "expression",
			body: `1 + 2;`,
			want: 10,
		},
		{
			name: "function call",
			body: `def f() -> int { return 1; } f();`,
			want: 10,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(`int n = 0; ` + c.body + ` n = n + 10;`)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}

		if got := env.Variable["n"]; got != c.want {
			t.Errorf("%s: n=%v, want %v", c.name, got, c.want)
		}
	}
}

func TestVisitor_returnStopsExecution(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{
			name: "direct",
			body: `return 1;`,
		},
		{
			name: "if",
			body: `if (true) { return 1; touched = 1; }`,
		},
		{
			name: "else",
			body: `if (false) {} else { return 1; touched = 1; }`,
		},
		{
			name: "nested scope",
			body: `{ { return 1; } touched = 1; }`,
		},
		{
			name: "for",
			body: `for int i in [0:2] { return 1; touched = 1; }`,
		},
		{
			name: "while",
			body: `while (true) { return 1; touched = 1; break; }`,
		},
		{
			name: "switch case",
			body: `switch (1) { case 1 { return 1; touched = 1; } }`,
		},
		{
			name: "switch default",
			body: `switch (1) { default { return 1; touched = 1; } }`,
		},
		{
			name: "nested function",
			body: `def g() -> int { return 1; touched = 1; } return g();`,
		},
	}

	for _, c := range cases {
		qsim, env, err := visitor.Run(`
				int touched = 0;
				qubit q;
				def f(qubit arg) -> int { ` + c.body + ` U(pi, 0, pi) arg; return 2; }
				int result = f(q);
				int after = 3;
		`)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}

		if got := env.Variable["result"]; got != int64(1) {
			t.Errorf("%s: result=%v, want 1", c.name, got)
		}

		if got := env.Variable["touched"]; got != int64(0) {
			t.Errorf("%s: executed statements after return: touched=%v", c.name, got)
		}

		if got := env.Variable["after"]; got != int64(3) {
			t.Errorf("%s: caller did not continue: after=%v", c.name, got)
		}

		if got := qsim.Probability()[0]; math.Abs(got-1) > 1e-8 {
			t.Errorf("%s: executed gate after return: P(0)=%v", c.name, got)
		}
	}
}

func TestVisitor_emptyReturn(t *testing.T) {
	cases := []struct {
		text        string
		wantTouched int64
		wantAfter   int64
	}{
		{
			text: `
		int touched = 0;
		def f() { if (true) { return; } touched = 1; }
		f();
		int after = 2;
	`,
			wantTouched: 0,
			wantAfter:   2,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			t.Fatal(err)
		}

		if got := env.Variable["touched"]; got != c.wantTouched {
			t.Errorf("touched=%v, want=%v", got, c.wantTouched)
		}

		if got := env.Variable["after"]; got != c.wantAfter {
			t.Errorf("after=%v, want=%v", got, c.wantAfter)
		}
	}
}

func TestVisitor_invalidReturn(t *testing.T) {
	cases := []struct {
		code string
	}{
		{
			code: `def f() -> int { return missing; } int result = f();`,
		},
	}

	for _, c := range cases {
		if _, _, err := visitor.Run(c.code); err != nil {
			if err.Error() == `undefined "missing"` {
				continue
			}

			t.Errorf("got=%v, want undefined identifier error", err)
		}

		t.Fatal("got=nil, want undefined identifier error")
	}
}

func TestVisitor_broadcast(t *testing.T) {
	cases := []struct {
		name string
		text string
		want int
	}{
		{
			name: "register control and scalar target",
			text: `qubit[2] c; qubit t; U(pi,0,pi) c[0]; ctrl @ U(pi,0,pi) c,t;`,
			want: 0b101,
		},
		{
			name: "register control and register target",
			text: `qubit[2] c; qubit[2] t; U(pi,0,pi) c[0]; ctrl @ U(pi,0,pi) c,t;`,
			want: 0b1010,
		},
		{
			name: "scalar control and register target",
			text: `qubit c; qubit[2] t; U(pi,0,pi) c; ctrl @ U(pi,0,pi) c,t;`,
			want: 0b111,
		},
		{
			name: "negative register control",
			text: `qubit[2] c; qubit[2] t; U(pi,0,pi) c[0]; negctrl @ U(pi,0,pi) c,t;`,
			want: 0b1001,
		},
		{
			name: "negative scalar control",
			text: `qubit c; qubit[2] t; negctrl @ U(pi,0,pi) c,t;`,
			want: 0b011,
		},
		{
			name: "multiple register controls",
			text: `qubit[2] a; qubit[2] b; qubit[2] t; U(pi,0,pi) a[0],b[0]; ctrl(2) @ U(pi,0,pi) a,b,t;`,
			want: 0b101010,
		},
		{
			name: "mixed controls",
			text: `qubit[2] a; qubit b; qubit[2] t; U(pi,0,pi) a[0]; ctrl @ negctrl @ U(pi,0,pi) a,b,t;`,
			want: 0b10010,
		},
		{
			name: "user-defined gate",
			text: `gate cx c,t { ctrl @ U(pi,0,pi) c,t; } qubit[2] c; qubit[2] t; U(pi,0,pi) c[0]; cx c,t;`,
			want: 0b1010,
		},
		{
			name: "single-element arrays",
			text: `qubit[1] c; qubit[1] t; U(pi,0,pi) c; ctrl @ U(pi,0,pi) c,t;`,
			want: 0b11,
		},
		{
			name: "indexed single-element control array",
			text: `qubit[1] c; qubit[2] t; U(pi,0,pi) c; ctrl @ U(pi,0,pi) c[0],t;`,
			want: 0b111,
		},
		{
			name: "indexed single-element target array",
			text: `qubit[2] c; qubit[1] t; U(pi,0,pi) c[0]; ctrl @ U(pi,0,pi) c,t[0];`,
			want: 0b101,
		},
		{
			name: "scalar alias",
			text: `qubit[1] a; let c=a[0]; qubit[2] t; U(pi,0,pi) c; ctrl @ U(pi,0,pi) c,t;`,
			want: 0b111,
		},
		{
			name: "user-defined gate with scalar control",
			text: `gate cx c,t { ctrl @ U(pi,0,pi) c,t; } qubit c; qubit[2] t; U(pi,0,pi) c; cx c,t;`,
			want: 0b111,
		},
		{
			name: "nested user-defined gate with single-element arrays",
			text: `gate cx c,t { ctrl @ U(pi,0,pi) c,t; } gate g c,t { cx c,t; } qubit[1] c; qubit[1] t; U(pi,0,pi) c; g c,t;`,
			want: 0b11,
		},
		{
			name: "subroutine scalar parameter",
			text: `def f(qubit c, qubit[2] t) { ctrl @ U(pi,0,pi) c,t; } qubit[1] c; qubit[2] t; U(pi,0,pi) c; f(c[0],t);`,
			want: 0b111,
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}

		if got := qsim.Probability()[c.want]; math.Abs(got-1) > 1e-8 {
			t.Errorf("%s: P(%b)=%v, want 1; probabilities=%v", c.name, c.want, got, qsim.Probability())
		}
	}
}

func TestVisitor_invalidBroadcast(t *testing.T) {
	cases := []struct {
		name string
		code string
	}{
		{
			name: "different array lengths",
			code: `qubit[2] c; qubit[3] t; ctrl @ U(pi,0,pi) c,t;`,
		},
		{
			name: "single-element control array",
			code: `qubit[1] c; qubit[2] t; ctrl @ U(pi,0,pi) c,t;`,
		},
		{
			name: "single-element target array",
			code: `qubit[2] c; qubit[1] t; ctrl @ U(pi,0,pi) c,t;`,
		},
		{
			name: "negative control",
			code: `qubit[1] c; qubit[2] t; negctrl @ U(pi,0,pi) c,t;`,
		},
		{
			name: "multiple controls",
			code: `qubit[1] a; qubit[2] b; qubit t; ctrl(2) @ U(pi,0,pi) a,b,t;`,
		},
		{
			name: "old-style register",
			code: `qreg c[1]; qreg t[2]; ctrl @ U(pi,0,pi) c,t;`,
		},
		{
			name: "array alias",
			code: `qubit[1] c; qubit[2] t; let a=c; ctrl @ U(pi,0,pi) a,t;`,
		},
		{
			name: "concatenated alias",
			code: `qubit a; qubit b; let c=a++b; qubit[1] t; ctrl @ U(pi,0,pi) c,t;`,
		},
		{
			name: "user-defined gate",
			code: `gate cx c,t { ctrl @ U(pi,0,pi) c,t; } qubit[1] c; qubit[2] t; cx c,t;`,
		},
		{
			name: "user-defined gate without a controlled operation",
			code: `gate g a,b { U(pi,0,pi) a; U(pi,0,pi) b; } qubit[1] a; qubit[2] b; g a,b;`,
		},
		{
			name: "empty user-defined gate",
			code: `gate g a,b {} qubit[1] a; qubit[2] b; g a,b;`,
		},
		{
			name: "subroutine array parameters",
			code: `def f(qubit[1] c, qubit[2] t) { ctrl @ U(pi,0,pi) c,t; } qubit[1] c; qubit[2] t; f(c,t);`,
		},
	}

	for _, c := range cases {
		if _, _, err := visitor.Run(c.code); err != nil {
			if strings.Contains(err.Error(), "gate operand size mismatch") {
				continue
			}

			t.Errorf("got=%v, want register size error", err)
		}

		t.Errorf("register size error")
	}
}

func TestVisitor_localBindings(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		want     string
		wantProb []float64
	}{
		{
			name: "loop variable",
			text: `int i=100; int sum=0; for int i in [0:2] { sum=sum+i; }`,
			want: "map[i:100 sum:3]",
		},
		{
			name: "nested loop variables",
			text: `int i=100; int sum=0; for int i in [0:1] { for int i in [0:1] { sum=sum+i; } sum=sum+i; }`,
			want: "map[i:100 sum:3]",
		},
		{
			name:     "gate parameter",
			text:     `float theta=0.0; gate g(theta) a { U(theta,0,0) a; } qubit q; g(pi) q;`,
			want:     "map[theta:0]",
			wantProb: []float64{0, 1},
		},
		{
			name: "nested gate parameters",
			text: `float theta=0.0;
				gate inner(theta) a { U(theta,0,0) a; }
				gate outer(theta) a { inner(2*theta) a; U(theta,0,0) a; }
				qubit q; outer(pi/2) q;`,
			want:     "map[theta:0]",
			wantProb: []float64{0.5, 0.5},
		},
	}

	for _, c := range cases {
		qsim, env, err := visitor.Run(c.text)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}

		if got := fmt.Sprint(env.Variable); got != c.want {
			t.Errorf("%s: got=%s, want=%s", c.name, got, c.want)
		}

		if got := qsim.Probability(); !slices.EqualFunc(got, c.wantProb, func(a, b float64) bool {
			return math.Abs(a-b) <= 1e-8
		}) {
			t.Errorf("%s: probabilities=%v, want=%v", c.name, got, c.wantProb)
		}
	}
}

func TestVisitor_endInInclude(t *testing.T) {
	cases := []struct {
		code string
		want int64
	}{
		{
			code: "end;",
			want: 0,
		},
		{
			code: "int inside=1;",
			want: 1,
		},
		{
			code: "def stop() -> int { end; } int inside=stop();",
			want: 0,
		},
	}

	for _, c := range cases {
		path := filepath.Join(t.TempDir(), "included.qasm")
		if err := os.WriteFile(path, []byte(c.code), 0600); err != nil {
			t.Fatal(err)
		}

		_, env, err := visitor.Run(fmt.Sprintf(`int touched=0; include %q; touched=1;`, path))
		if err != nil {
			t.Errorf("%s: %v", c.code, err)
			continue
		}

		if got := env.Variable["touched"]; got != c.want {
			t.Errorf("%s: touched=%v, want=%v", c.code, got, c.want)
		}
	}
}

func TestVisitor_openQASM3(t *testing.T) {
	cases := []struct {
		theta  float64
		phi    float64
		lambda float64
	}{
		{
			theta:  0,
			phi:    0,
			lambda: 0,
		},
		{
			theta:  math.Pi,
			phi:    0,
			lambda: math.Pi,
		},
		{
			theta:  math.Pi / 2,
			phi:    0,
			lambda: math.Pi,
		},
		{
			theta:  2 * math.Pi,
			phi:    0,
			lambda: 0,
		},
		{
			theta:  0.7,
			phi:    0.3,
			lambda: 1.1,
		},
		{
			theta:  0.7 + 2*math.Pi,
			phi:    0.3,
			lambda: 1.1,
		},
		{
			theta:  0.7,
			phi:    0.3 + 2*math.Pi,
			lambda: 1.1,
		},
		{
			theta:  0.7,
			phi:    0.3,
			lambda: 1.1 + 2*math.Pi,
		},
	}

	for _, c := range cases {
		theta, phi, lambda := c.theta, c.phi, c.lambda
		// Directly use the specification's matrix, independently of gate.U.
		// https://openqasm.com/versions/3.0/language/gates.html#built-in-single-qubit-gate-u
		e := cmplx.Exp(complex(0, theta))
		want := [2][2]complex128{
			{(1 + e) / 2, -1i * cmplx.Exp(complex(0, lambda)) * (1 - e) / 2},
			{1i * cmplx.Exp(complex(0, phi)) * (1 - e) / 2, cmplx.Exp(complex(0, phi+lambda)) * (1 + e) / 2},
		}

		tree, err := xparser.Parse(fmt.Sprintf("U(%.17g, %.17g, %.17g) q;", theta, phi, lambda))
		if err != nil {
			t.Fatal(err)
		}

		for col := range 2 {
			var basis [2]complex128
			basis[col] = 1

			qsim := q.New()
			env := environ.New()
			env.Qubit["q"] = qsim.New(basis[:]...)

			if err := visitor.New(qsim, env).Run(tree); err != nil {
				t.Fatal(err)
			}

			for row, got := range qsim.Amplitude() {
				if cmplx.Abs(got-want[row][col]) > 1e-12 {
					t.Errorf("U[%d,%d]=%v, want %v", row, col, got, want[row][col])
				}
			}
		}
	}
}

func TestVisitor_controlledPhase(t *testing.T) {
	cases := []struct {
		name string
		text string
		want []complex128
	}{
		{
			name: "controlled 2*pi is identity",
			text: `qubit[2] q; h q[0]; ctrl @ U(2*pi, 0, 0) q[0], q[1]; h q[0];`,
			want: []complex128{1, 0, 0, 0},
		},
		{
			name: "negative controlled 2*pi is identity",
			text: `qubit[2] q; h q[0]; negctrl @ U(2*pi, 0, 0) q[0], q[1]; h q[0];`,
			want: []complex128{1, 0, 0, 0},
		},
		{
			name: "controlled iX produces a relative phase",
			text: `qubit[2] q; h q[0]; ctrl @ U(pi, 0, pi) q[0], q[1]; cx q[0], q[1]; h q[0];`,
			want: []complex128{0.5 + 0.5i, 0, 0.5 - 0.5i, 0},
		},
		{
			name: "gphase applies once with spectators",
			text: `qubit[3] q; gphase(pi/2);`,
			want: []complex128{1i, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "broadcast H includes one correction per call",
			text: `qubit[2] q; qubit spectator; h q;`,
			want: []complex128{0.5, 0, 0.5, 0, 0.5, 0, 0.5, 0},
		},
		{
			name: "broadcast CX preserves phase",
			text: `qubit[2] c; qubit[2] t; x c; cx c, t;`,
			want: []complex128{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(`include "../testdata/stdgates.qasm";` + c.text)
		if err != nil {
			t.Fatal(err)
		}

		got := qsim.Amplitude()
		if len(got) != len(c.want) {
			t.Fatalf("got %d amplitudes, want %d", len(got), len(c.want))
		}

		for i := range got {
			if cmplx.Abs(got[i]-c.want[i]) > 1e-12 {
				t.Errorf("amplitude[%d]=%v, want %v", i, got[i], c.want[i])
			}
		}
	}
}

func TestVisitor_controlledGphase(t *testing.T) {
	cases := []struct {
		text   string
		active func(int) bool
	}{
		{
			text: `ctrl @ gphase(pi/2) q[0];`,
			active: func(i int) bool {
				return i&0b100 != 0
			},
		},
		{
			text: `negctrl @ gphase(pi/2) q[0];`,
			active: func(i int) bool {
				return i&0b100 == 0
			},
		},
		{
			text: `ctrl(2) @ gphase(pi/2) q[0], q[1];`,
			active: func(i int) bool {
				return i&0b110 == 0b110
			},
		},
		{
			text: `ctrl @ negctrl @ gphase(pi/2) q[0], q[1];`,
			active: func(i int) bool {
				return i&0b110 == 0b100
			},
		},
		{
			text: `inv @ ctrl @ gphase(-pi/2) q[0];`,
			active: func(i int) bool {
				return i&0b100 != 0
			},
		},
		{
			text: `ctrl @ pow(2) @ gphase(pi/4) q[0];`,
			active: func(i int) bool {
				return i&0b100 != 0
			},
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(`include "../testdata/stdgates.qasm"; qubit[3] q; h q;` + c.text)
		if err != nil {
			t.Fatal(err)
		}

		for i, got := range qsim.Amplitude() {
			want := complex(1/math.Sqrt(8), 0)

			if c.active(i) {
				want *= 1i
			}

			if cmplx.Abs(got-want) > 1e-12 {
				t.Errorf("amplitude[%d]=%v, want %v", i, got, want)
			}
		}
	}
}

func TestVisitor_classicalBitOperators(t *testing.T) {
	cases := []struct {
		op   string
		a    bool
		b    bool
		want bool
	}{
		{
			op:   "&",
			a:    false,
			b:    false,
			want: false,
		},
		{
			op:   "&",
			a:    false,
			b:    true,
			want: false,
		},
		{
			op:   "&",
			a:    true,
			b:    false,
			want: false,
		},
		{
			op:   "&",
			a:    true,
			b:    true,
			want: true,
		},
		{
			op:   "|",
			a:    false,
			b:    false,
			want: false,
		},
		{
			op:   "|",
			a:    false,
			b:    true,
			want: true,
		},
		{
			op:   "|",
			a:    true,
			b:    false,
			want: true,
		},
		{
			op:   "|",
			a:    true,
			b:    true,
			want: true,
		},
		{
			op:   "^",
			a:    false,
			b:    false,
			want: false,
		},
		{
			op:   "^",
			a:    false,
			b:    true,
			want: true,
		},
		{
			op:   "^",
			a:    true,
			b:    false,
			want: true,
		},
		{
			op:   "^",
			a:    true,
			b:    true,
			want: false,
		},
	}

	for _, c := range cases {
		texts := []string{
			fmt.Sprintf(`bit a = %t; bit b = %t; bit s = a %s b; s;`, c.a, c.b, c.op),
			fmt.Sprintf(`bit[2] m; m[0] = %t; m[1] = %t; bit s = m[0] %s m[1]; s;`, c.a, c.b, c.op),
		}

		for _, text := range texts {
			got, err := visitor.Visit(text)
			if err != nil {
				t.Errorf("%s: unexpected error: %v", text, err)
				continue
			}

			if got != c.want {
				t.Errorf("%s: got=%v (%T), want=%v", text, got, got, c.want)
			}
		}
	}
}

func TestVisitor_bitwiseExpressions(t *testing.T) {
	cases := []struct {
		name string
		text string
		want any
	}{
		{
			name: "issue 52",
			text: `OPENQASM 3.0; bit[2] m; bit s; m[0] = 1; m[1] = 0; s = m[0] ^ m[1]; s;`,
			want: true,
		},
		{
			name: "assign one to bit",
			text: `bit s; s = 1; s;`,
			want: true,
		},
		{
			name: "assign zero to bit",
			text: `bit s = true; s = 0; s;`,
			want: false,
		},
		{
			name: "clear indexed bit",
			text: `bit[2] m = "11"; m[0] = 0; m[0];`,
			want: false,
		},
		{
			name: "syndrome",
			text: `bit[7] m = "1010101"; bit s = m[0] ^ m[2] ^ m[4] ^ m[6]; s;`,
			want: false,
		},
		{
			name: "syndrome odd parity",
			text: `bit[7] m = "1010100"; bit s = m[0] ^ m[2] ^ m[4] ^ m[6]; s;`,
			want: true,
		},
		{
			name: "bit not true",
			text: `bit a = true; bit s = ~a; s;`,
			want: false,
		},
		{
			name: "bit not false",
			text: `bit a; bit s = ~a; s;`,
			want: true,
		},
		{
			name: "indexed not",
			text: `bit[2] m = "10"; ~m[1];`,
			want: true,
		},
		{
			name: "bit left shift zero",
			text: `bit a = true; a << 0;`,
			want: true,
		},
		{
			name: "bit left shift",
			text: `bit a = true; a << 1;`,
			want: false,
		},
		{
			name: "bit right shift zero",
			text: `bit a = true; a >> 0;`,
			want: true,
		},
		{
			name: "bit right shift",
			text: `bit a = true; a >> 1;`,
			want: false,
		},
		{
			name: "register and",
			text: `bit[4] a = "1010"; bit[4] b = "1100"; bit[4] c = a & b; c;`,
			want: []bool{true, false, false, false},
		},
		{
			name: "register or",
			text: `"1010" | "1100";`,
			want: []bool{true, true, true, false},
		},
		{
			name: "register xor",
			text: `"1010" ^ "1100";`,
			want: []bool{false, true, true, false},
		},
		{
			name: "register not",
			text: `~"1010";`,
			want: []bool{false, true, false, true},
		},
		{
			name: "register left shift",
			text: `"10001111" << 1;`,
			want: []bool{false, false, false, true, true, true, true, false},
		},
		{
			name: "register right shift",
			text: `"10001111" >> 2;`,
			want: []bool{false, false, true, false, false, false, true, true},
		},
		{
			name: "register shift zero",
			text: `"1010" << 0;`,
			want: []bool{true, false, true, false},
		},
		{
			name: "register shift width",
			text: `"1010" << 4;`,
			want: []bool{false, false, false, false},
		},
		{
			name: "register shift beyond width",
			text: `"1010" >> 9223372036854775807;`,
			want: []bool{false, false, false, false},
		},
		{
			name: "register chained shift",
			text: `"1001" << 1 >> 1;`,
			want: []bool{false, false, false, true},
		},
		{
			name: "precedence",
			text: `bit[3] m = "101"; m[0] | m[1] ^ m[2] & ~m[0];`,
			want: true,
		},
		{
			name: "cast integer and",
			text: `int(10) & 12;`,
			want: int64(8),
		},
		{
			name: "cast integer or",
			text: `10 | int(12);`,
			want: int64(14),
		},
		{
			name: "cast integer xor",
			text: `int(10) ^ int(12);`,
			want: int(6),
		},
		{
			name: "unsigned and",
			text: `uint(10) & uint(12);`,
			want: uint(8),
		},
		{
			name: "unsigned or",
			text: `uint(10) | uint(12);`,
			want: uint(14),
		},
		{
			name: "unsigned xor",
			text: `uint(10) ^ uint(12);`,
			want: uint(6),
		},
		{
			name: "unsigned not",
			text: `~uint(5);`,
			want: ^uint(5),
		},
		{
			name: "cast integer shift",
			text: `int(11) << 2;`,
			want: int(44),
		},
		{
			name: "unsigned shift",
			text: `uint(11) >> uint(1);`,
			want: uint(5),
		},
		{
			name: "bit logical and",
			text: `bit[2] m = "10"; m[0] && m[1];`,
			want: false,
		},
		{
			name: "bit logical or",
			text: `bit[2] m = "10"; m[0] || m[1];`,
			want: true,
		},
	}

	for _, c := range cases {
		got, err := visitor.Visit(c.text)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.name, err)
			continue
		}

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got=%v (%T), want=%v (%T)", c.name, got, got, c.want, c.want)
		}
	}
}

func TestVisitor_invalidBitwiseOperand(t *testing.T) {
	cases := []struct {
		text string
	}{
		{
			text: `true & 1;`,
		},
		{
			text: `1 | true;`,
		},
		{
			text: `1 ^ false;`,
		},
		{
			text: `1.5 & 2;`,
		},
		{
			text: `1 | 2.5;`,
		},
		{
			text: `1.5 ^ 2.5;`,
		},
		{
			text: `"10" & "1";`,
		},
		{
			text: `"1" | "10";`,
		},
		{
			text: `"10" ^ "1";`,
		},
		{
			text: `true & "1";`,
		},
		{
			text: `"1" | true;`,
		},
		{
			text: `"1" ^ 1;`,
		},
		{
			text: `~1.5;`,
		},
		{
			text: `1.5 << 1;`,
		},
		{
			text: `1 >> 1.5;`,
		},
		{
			text: `1 << true;`,
		},
		{
			text: `1 << -1;`,
		},
		{
			text: `1 >> -1;`,
		},
		{
			text: `"10" << -1;`,
		},
		{
			text: `true >> -1;`,
		},
		{
			text: `1 << int(-1);`,
		},
		{
			text: `1 && true;`,
		},
		{
			text: `true && 1;`,
		},
		{
			text: `false || 1;`,
		},
		{
			text: `1 || false;`,
		},
	}

	for _, c := range cases {
		if _, err := visitor.Visit(c.text); err == nil {
			t.Errorf("%s: expected an operand error", c.text)
		}
	}
}

func TestVisitor_invalidUndefinedBitwiseOperand(t *testing.T) {
	cases := []struct {
		text string
	}{
		{
			text: `missing & true;`,
		},
		{
			text: `true | missing;`,
		},
		{
			text: `true ^ missing;`,
		},
		{
			text: `true << missing;`,
		},
	}

	for _, c := range cases {
		if _, err := visitor.Visit(c.text); err == nil || !strings.Contains(err.Error(), `undefined "missing"`) {
			t.Errorf("%s: got=%v, want undefined operand error", c.text, err)
		}
	}
}

func TestVisitor_invalidBitAssignments(t *testing.T) {
	cases := []struct {
		text string
	}{
		{
			text: `bit s; s = 2;`,
		},
		{
			text: `bit s; s = 1.5;`,
		},
		{
			text: `bit[2] m; m[0] = -1;`,
		},
		{
			text: `bit[2] m; m[0] = 1.5;`,
		},
	}

	for _, c := range cases {
		if _, err := visitor.Visit(c.text); err == nil {
			t.Errorf("%s: expected an assignment error", c.text)
		}
	}
}

func TestVisitor_bitwiseDoesNotMutateOperands(t *testing.T) {
	cases := []struct {
		expr string
	}{
		{
			expr: `a & b`,
		},
		{
			expr: `a | b`,
		},
		{
			expr: `a ^ b`,
		},
		{
			expr: `~a`,
		},
		{
			expr: `a << 0`,
		},
		{
			expr: `a << 1`,
		},
		{
			expr: `a >> 1`,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(`bit[4] a = "1010"; bit[4] b = "1100"; bit[4] c = ` + c.expr + `; c[0] = false;`)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", c.expr, err)
			continue
		}

		if !reflect.DeepEqual(env.ClBitArray["a"], []bool{true, false, true, false}) {
			t.Errorf("%s: operands changed: %v", c.expr, env.ClBitArray)
		}

		if !reflect.DeepEqual(env.ClBitArray["b"], []bool{true, true, false, false}) {
			t.Errorf("%s: operands changed: %v", c.expr, env.ClBitArray)
		}
	}
}
