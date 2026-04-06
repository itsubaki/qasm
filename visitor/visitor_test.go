package visitor_test

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/environ"
	"github.com/itsubaki/qasm/gen/parser"
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

func ExampleVisitor_Run_error() {
	text := `
	const int a = 42;
	const int a = 43;
	`

	if _, _, err := visitor.Run(text); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// declare const "a": already declared
}

func ExampleVisit() {
	text := `
	const int a = 42;
	const int a = 43;
	`

	_, err := visitor.Visit(text)
	fmt.Println(err)

	// Output:
	// declare const "a": already declared
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

func ExampleVisitor_VisitIncludeStatement_invalid() {
	text := `
	include "../testdata/invalid.qasm";
	qubit q;
	h q;
	`

	if _, _, err := visitor.Run(text); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// include: literal "invalid": undefined
}

func ExampleVisitor_VisitIncludeStatement_fileNotFound() {
	text := `include "file_not_found.qasm";`

	if _, _, err := visitor.Run(text); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// read file=file_not_found.qasm: open file_not_found.qasm: no such file or directory
}

func ExampleVisitor_VisitErrorNode() {
	token := &antlr.BaseToken{}
	token.SetText("something went wrong")

	v := visitor.New(q.New(), environ.New())
	fmt.Println(v.VisitErrorNode(antlr.NewErrorNodeImpl(token)))

	// Output:
	// something went wrong
}

func TestVisitor_VisitChildren(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
	}{
		{
			text: "OPENQASM 3.0; const int a = 42;",
			want: "map[a:42]",
		},
		{
			text:   "const int a = 42; const int a = 43;",
			errMsg: `declare const "a": already declared`,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		env := environ.New()
		result := visitor.New(q.New(), env).VisitChildren(p.Program())
		if err, ok := result.(error); ok && err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", env.Const) != c.want {
			t.Errorf("got=%v, want=%v", env.Const, c.want)
		}
	}
}

func TestVisitor_VisitConstDeclarationStatement(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
	}{
		{
			text: "const int a = 42;",
			want: "map[a:42]",
		},
		{
			text: "const uint N = 3 * 5;",
			want: "map[N:15]",
		},
		{
			text:   "const int a = 42; const int a = 43;",
			errMsg: `declare const "a": already declared`,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", env.Const) != c.want {
			t.Errorf("got=%v, want=%v", env.Const, c.want)
		}
	}
}

func TestVisitor_VisitClassicalDeclarationStatement(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
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
			text:   `qubit[3] q; bit c = measure q;`,
			errMsg: "assign 3 bits to a single bit",
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
		{
			text:   "float a = 1; float a = 0;",
			errMsg: `declare float "a": already declared`,
		},
		{
			text:   "int a = 1; int a = 0;",
			errMsg: `declare int "a": already declared`,
		},
		{
			text:   "uint a = 1; uint a = 0;",
			errMsg: `declare uint "a": already declared`,
		},
		{
			text:   "bool a = true; bool a = false;",
			errMsg: `declare bool "a": already declared`,
		},
		{
			text:   "bit a; bool a;",
			errMsg: `declare bool "a": already declared`,
		},
		{
			text:   "bit[1] a; bool a;",
			errMsg: `declare bool "a": already declared`,
		},
		{
			text:   "bit a; bit a;",
			errMsg: `declare bit "a": already declared`,
		},
		{
			text:   "int a = 1; bit a;",
			errMsg: `declare bit "a": already declared`,
		},
		{
			text:   "bit[1] a; bit a;",
			errMsg: `declare bit "a": already declared`,
		},
		{
			text:   "angle[4] a; angle[4] a;",
			errMsg: `declare angle "a": already declared`,
		},
		{
			text:   "bit a; angle[4] a;",
			errMsg: `declare angle "a": already declared`,
		},
		{
			text:   "bit[1] a; angle[4] a;",
			errMsg: `declare angle "a": already declared`,
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
		_, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if len(env.Bit) > 0 && fmt.Sprintf("%v", env.Bit) != c.want {
			t.Errorf("got=%v, want=%v", env.Bit, c.want)
		}

		if len(env.BitArray) > 0 && fmt.Sprintf("%v", env.BitArray) != c.want {
			t.Errorf("got=%v, want=%v", env.BitArray, c.want)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitQuantumDeclarationStatement(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
	}{
		{
			text: "qubit q;",
			want: "map[q:[0]]",
		},
		{
			text: "qubit[2] q;",
			want: "map[q:[0 1]]",
		},
		{
			text: "qubit q0; qubit[2] q1; qubit[3] q3; qubit q4;",
			want: "map[q0:[0] q1:[1 2] q3:[3 4 5] q4:[6]]",
		},
		{
			text:   "qubit[1.1] q;",
			errMsg: `size must be an integer "qubit[1.1]"`,
		},
		{
			text:   "qubit q; qubit q;",
			errMsg: `declare qubit "q": already declared`,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", env.Qubit) != c.want {
			t.Errorf("got=%v, want=%v", env.Qubit, c.want)
		}
	}
}

func TestVisitor_VisitAliasDeclarationStatement(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
	}{
		{
			text: "qubit[5] q; let myreg = q[1:4];",
			want: "map[myreg:[1 2 3] q:[0 1 2 3 4]]",
		},
		{
			text: "qubit[5] q; let myreg = q[1];",
			want: "map[myreg:[1] q:[0 1 2 3 4]]",
		},
		{
			text: "qubit[2] one; qubit[10] two; let concatenated = one ++ two;",
			want: "map[concatenated:[0 1 2 3 4 5 6 7 8 9 10 11] one:[0 1] two:[2 3 4 5 6 7 8 9 10 11]]",
		},
		{
			text:   "qubit[5] q; let q = q[1:4];",
			errMsg: `declare alias "q": already declared`,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", env.Qubit) != c.want {
			t.Errorf("got=%v, want=%v", env.Qubit, c.want)
		}
	}
}

func TestVisitor_VisitOldStyleDeclarationStatement(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
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
			text:   "qreg q; qreg q;",
			errMsg: `declare qreg "q": already declared`,
		},
		{
			text: "creg c;",
			want: "map[c:[false]]",
		},
		{
			text: "creg c[2];",
			want: "map[c:[false false]]",
		},
		{
			text:   "creg c; creg c;",
			errMsg: `declare creg "c": already declared`,
		},
		{
			text:   "bit c; creg c;",
			errMsg: `declare creg "c": already declared`,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if len(env.Qubit) > 0 && fmt.Sprintf("%v", env.Qubit) != c.want {
			t.Errorf("got=%v, want=%v", env.Qubit, c.want)
		}

		if len(env.Bit) > 0 && fmt.Sprintf("%v", env.Bit) != c.want {
			t.Errorf("got=%v, want=%v", env.Bit, c.want)
		}

		if len(env.BitArray) > 0 && fmt.Sprintf("%v", env.BitArray) != c.want {
			t.Errorf("got=%v, want=%v", env.BitArray, c.want)
		}
	}
}

func TestVisitor_VisitAssignmentStatement(t *testing.T) {
	type Want struct {
		bit      []string
		qubit    []string
		variable []string
	}

	cases := []struct {
		text   string
		want   Want
		errMsg string
	}{
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
				c = measure q;
			`,
			want: Want{
				bit: []string{
					"map[c:[false false]]",
					"map[c:[true true]]",
				},
				qubit: []string{
					"[[00][  0]( 1.0000 0.0000i): 1.0000]",
					"[[11][  3]( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
				c[0] = measure q[0];
				c[1] = measure q[1];
			`,
			want: Want{
				bit: []string{
					"map[c:[false false]]",
					"map[c:[true true]]",
				},
				qubit: []string{
					"[[00][  0]( 1.0000 0.0000i): 1.0000]",
					"[[11][  3]( 1.0000 0.0000i): 1.0000]",
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
			text:   `bit c; c = "10";`,
			errMsg: `assign 2 bits to a single bit`,
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
		{
			text: `
				qubit q;
				U(pi/2.0, 0, pi) q;
				c = measure q;
			`,
			errMsg: `operand "c": undefined`,
		},
	}

	for _, c := range cases {
		qsim, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

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
				if fmt.Sprintf("%v", env.Bit) == w {
					found = true
				}

				if fmt.Sprintf("%v", env.BitArray) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v/%v, want=%v", env.Bit, env.BitArray, c.want.bit)
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

func TestVisitor_VisitMeasureArrowAssignmentStatement(t *testing.T) {
	type Want struct {
		bit   []string
		qubit []string
	}

	cases := []struct {
		text   string
		want   Want
		errMsg string
	}{
		{
			text: `
				qubit q;
				bit c;
				U(pi, 0, pi) q;
				measure q -> c;
			`,
			want: Want{
				bit:   []string{"map[c:true]"},
				qubit: []string{"[[1][  1]( 1.0000 0.0000i): 1.0000]"},
			},
		},
		{
			text: `
				qubit[3] q;
				bit c;
				measure q -> c;
			`,
			errMsg: "assign 3 bits to a single bit",
		},
		{
			text: `
				qubit[3] q;
				bit c;
				U(pi, 0, pi) q[1];
				measure q[1] -> c;
			`,
			want: Want{
				bit:   []string{"map[c:true]"},
				qubit: []string{"[[010][  2]( 1.0000 0.0000i): 1.0000]"},
			},
		},
		{
			text: `
				qubit q;
				bit[1] c;
				U(pi, 0, pi) q;
				measure q -> c;
			`,
			want: Want{
				bit:   []string{"map[c:[true]]"},
				qubit: []string{"[[1][  1]( 1.0000 0.0000i): 1.0000]"},
			},
		},
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
				measure q -> c;
			`,
			want: Want{
				bit: []string{
					"map[c:[false false]]",
					"map[c:[true true]]",
				},
				qubit: []string{
					"[[00][  0]( 1.0000 0.0000i): 1.0000]",
					"[[11][  3]( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
				measure q[0] -> c[0];
				measure q[1] -> c[1];
			`,
			want: Want{
				bit: []string{
					"map[c:[false false]]",
					"map[c:[true true]]",
				},
				qubit: []string{
					"[[00][  0]( 1.0000 0.0000i): 1.0000]",
					"[[11][  3]( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2.0, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
				measure q;
			`,
			want: Want{
				qubit: []string{
					"[[00][  0]( 1.0000 0.0000i): 1.0000]",
					"[[11][  3]( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
	}

	for _, c := range cases {
		qsim, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

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
				if fmt.Sprintf("%v", env.Bit) == w {
					found = true
				}

				if fmt.Sprintf("%v", env.BitArray) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v/%v, want=%v", env.Bit, env.BitArray, c.want.bit)
			}
		}
	}
}

func TestVisitor_VisitResetStatement(t *testing.T) {
	cases := []struct {
		text   string
		want   []string
		errMsg string
	}{
		{
			text: `
				qubit q;
				U(pi/2, 0, pi) q;
				reset q;
			`,
			want: []string{
				"[0][  0]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text:   "int a = 1; reset a;",
			errMsg: `invalid operand "a"`,
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Fatalf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		for i, s := range qsim.State() {
			if s.String() != c.want[i] {
				t.Fatalf("got=%v, want=%v", s.String(), c.want[i])
			}
		}
	}
}

func TestVisitor_VisitMultiplicativeExpression(t *testing.T) {
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
			t.Fail()
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitAdditiveExpression(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
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
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitParenthesisExpression(t *testing.T) {
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

func TestVisitor_VisitCallExpression(t *testing.T) {
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

func TestVisitor_VisitPowerExpression(t *testing.T) {
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
			t.Fail()
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitLogicalOrExpression(t *testing.T) {
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

func TestVisitor_VisitLogicalAndExpression(t *testing.T) {
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

func TestVisitor_VisitBitwiseAndExpression(t *testing.T) {
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

func TestVisitor_VisitBitwiseOrExpression(t *testing.T) {
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

func TestVisitor_VisitBitwiseXorExpression(t *testing.T) {
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

func TestVisitor_VisitBitshiftExpression(t *testing.T) {
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

func TestVisitor_VisitEqualityExpression(t *testing.T) {
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
			t.Fail()
		}

		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitUnaryExpression(t *testing.T) {
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
			t.Fail()
		}

		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitComparisonExpression(t *testing.T) {
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
			t.Fail()
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitGateCallStatement(t *testing.T) {
	cases := []struct {
		text   string
		want   []string
		errMsg string
	}{
		{
			text: "qubit q; U(pi, 0, pi) q;",
			want: []string{
				"[1][  1]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `qubit q; U(pi/2, 0, pi) q;`,
			want: []string{
				"[0][  0]( 0.7071 0.0000i): 0.5000",
				"[1][  1]( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `qubit q; gphase(-π/2);`,
			want: []string{
				"[0][  0]( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
			`,
			want: []string{
				"[00][  0]( 0.7071 0.0000i): 0.5000",
				"[11][  3]( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				inv @ U(pi/2, 0, pi) q[0];
			`,
			want: []string{
				"[00][  0]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit[2] q;
				pow(2) @ U(pi/2, 0, pi) q[0];
			`,
			want: []string{
				"[00][  0]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `qubit[2] q; pow(3) @ U(pi/2, 0, pi) q[0];`,
			want: []string{
				"[00][  0]( 0.7071 0.0000i): 0.5000",
				"[10][  2]( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text:   `qubit q; gphase(true);`,
			errMsg: "invalid param 'bool(true)'",
		},
		{
			text:   `qubit q; gphase(a);`,
			errMsg: `literal "a": undefined`,
		},
		{
			text:   `qubit[2] q; U(true, 0, pi) q;`,
			errMsg: "invalid param 'bool(true)'",
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				ctrl @ pow(2) @ U(pi, 0, pi) q[0], q[1];
			`,
			errMsg: "pow with control modifier is not implemented: not implemented",
		},
		{
			text: `
				int a = 1;
				qubit t;
				ctrl @ U(pi, 0, pi) a, t;
			`,
			errMsg: `invalid operand "a,t"`,
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
			errMsg: `invalid operand "q"`,
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

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

func TestVisitor_VisitGateCallStatement_userdefined(t *testing.T) {
	cases := []struct {
		text   string
		want   []string
		errMsg string
	}{
		{
			text: `
				gate u q { }
				gate u q { }
			`,
			errMsg: `declare gate "u": already declared`,
		},
		{
			text:   `qubit[2] q; myg(pi, 0, pi) q;`,
			errMsg: `gate "myg": undefined`,
		},
		{
			text: `
				gate h q0 { U(pi/2, 0, pi) q0; }
				gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }
				qubit[2] q;
				h q[0];
				cx q[0], q[1];
			`,
			want: []string{
				"[00][  0]( 0.7071 0.0000i): 0.5000",
				"[11][  3]( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				gate x q0 { U(pi, 0, pi) q0; }
				gate y q0 { U(pi, pi/2, pi/2) q0; }
				qubit q;
				x q;
				y q;
			`,
			want: []string{
				"[0][  0]( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				gate x q { U(pi, 0, pi) q; }
				gate y q { U(pi, pi/2, pi/2) q; }
				gate xy q { x q; y q; }
				qubit q;
				xy q;
			`,
			want: []string{
				"[0][  0]( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				qubit q;
				u(pi, 0, pi) q;
			`,
			want: []string{
				"[1][  1]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				qubit[2] q;
				u(pi, 0, pi) q;
			`,
			want: []string{
				"[11][  3]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				qubit q;
				u(a, 0, pi) q;
			`,
			errMsg: `literal "a": undefined`,
		},
		{
			text: `
				gate myg q { }
				int a = 1;
				myg a;
			`,
			errMsg: `invalid operand "a"`,
		},
		{
			text: `
				gate myg q { myx q; }
				qubit q;
				myg q;
			`,
			errMsg: `gate call[0]: gate "myx": undefined`,
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				qubit q;
				inv @ u(pi, 0, pi) q;
			`,
			errMsg: "modifier is not implemented in user-defined: not implemented",
		},
		// not implemented.
		// {
		// 	text: `
		// 		gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
		// 		const int n = 3;
		// 		qubit q;
		// 		pow(n) @ u(pi, 0, pi) q;
		// 	`,
		// 	want: []string{
		// 		"[1][  1]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
		// 		qubit q;
		// 		u(pi, 0, pi) q;
		// 		inv @ u(pi, 0, pi) q;
		// 	`,
		// 	want: []string{
		// 		"[0][  0]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
		// 		gate invu(p0, p1, p2) q { inv @ u(p0, p1, p2) q; }
		// 		qubit q;
		// 		u(1, 2, 3) q;
		// 		invu(1, 2, 3) q;
		// 	`,
		// 	want: []string{
		// 		"[0][  0]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate x q { U(pi, 0, pi) q; }
		// 		qubit[2] q;
		// 		x q[0];
		// 		ctrl @ x q[0], q[1];
		// 	`,
		// 	want: []string{
		// 		"[11][  3]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate x q { U(pi, 0, pi) q; }
		// 		gate cx q0, q1 { ctrl @ x q0, q1; }
		// 		qubit[2] q;
		// 		qubit t;
		// 		x q;
		// 		ctrl @ cx q[0], q[1], t;
		// 	`,
		// 	want: []string{
		// 		"[111][  7]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate x q { U(pi, 0, pi) q; }
		// 		qubit[2] q;
		// 		negctrl @ x q[0], q[1];
		// 	`,
		// 	want: []string{
		// 		"[01][  1]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate x q { U(pi, 0, pi) q; }
		// 		gate negcx q0, q1 { negctrl @ x q0, q1; }
		// 		qubit[3] q;
		// 		negctrl @ negcx q[0], q[1], q[2];
		// 	`,
		// 	want: []string{
		// 		"[001][  1]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate x q { U(pi, 0, pi) q; }
		// 		qubit[2] q;
		// 		x q[1];
		// 		ctrl @ x q[1], q[0];
		// 	`,
		// 	want: []string{
		// 		"[11][  3]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate x q { U(pi, 0, pi) q; }
		// 		gate negcx q0, q1 { negctrl @ x q0, q1; }
		// 		qubit[3] q;
		// 		negctrl @ negcx q[1], q[2], q[0];
		// 	`,
		// 	want: []string{
		// 		"[100][  4]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate x q { U(pi, 0, pi) q; }
		// 		gate negcx q0, q1 { negctrl @ x q0, q1; }
		// 		qubit[3] q;
		// 		negctrl @ negcx q[2], q[0], q[1];
		// 	`,
		// 	want: []string{
		// 		"[010][  2]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
		// {
		// 	text: `
		// 		gate x q { U(pi, 0, pi) q; }
		// 		gate cx q0, q1 { ctrl @ x q0, q1; }
		// 		qubit[3] q;
		// 		x q[1];
		// 		x q[2];
		// 		ctrl @ cx q[1], q[2], q[0];
		// 	`,
		// 	want: []string{
		// 		"[111][  7]( 1.0000 0.0000i): 1.0000",
		// 	},
		// },
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

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

func TestVisitor_VisitGateModifier(t *testing.T) {
	cases := []struct {
		text   string
		want   []string
		errMsg string
	}{
		{
			text: `
				qubit q0;
				qubit q1;
				U(pi/2.0, 0, pi) q0;
				ctrl @ U(pi, 0, pi) q0, q1;
			`,
			want: []string{
				"[00][  0]( 0.7071 0.0000i): 0.5000",
				"[11][  3]( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				qubit[3] q;
				U(pi/2.0, 0, pi) q[0], q[1];
				ctrl @ ctrl @ U(pi, 0, pi) q[0], q[1], q[2];
			`,
			want: []string{
				"[000][  0]( 0.5000 0.0000i): 0.2500",
				"[010][  2]( 0.5000 0.0000i): 0.2500",
				"[100][  4]( 0.5000 0.0000i): 0.2500",
				"[111][  7]( 0.5000 0.0000i): 0.2500",
			},
		},
		{
			text: `
				qubit[3] q;
				U(pi/2.0, 0, pi) q[0], q[1];
				ctrl(2) @ U(pi, 0, pi) q[0], q[1], q[2];
			`,
			want: []string{
				"[000][  0]( 0.5000 0.0000i): 0.2500",
				"[010][  2]( 0.5000 0.0000i): 0.2500",
				"[100][  4]( 0.5000 0.0000i): 0.2500",
				"[111][  7]( 0.5000 0.0000i): 0.2500",
			},
		},
		{
			text: `
				qubit q0;
				qubit q1;
				U(pi/2.0, 0, pi) q0;
				negctrl @ U(pi, 0, pi) q0, q1;
			`,
			want: []string{
				"[01][  1]( 0.7071 0.0000i): 0.5000",
				"[10][  2]( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: `
				qubit q;
				U(pi, tau, euler) q;
				inv @ U(pi, tau, euler) q;
			`,
			want: []string{
				"[0][  0]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				const float half = pi / 2;
				qubit q;
				pow(2) @ U(half, -half, half) q;
			`,
			want: []string{
				"[1][  1]( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit q;
				pow(true) @ U(pi, 0, pi) q;
			`,
			errMsg: `apply "pow(true)@": unexpected type: bool`,
		},
		{
			text: `
				qubit q;
				inv @ pow(2) @ U(pi, 0, pi) q;
			`,
			want: []string{
				"[0][  0]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit[3] q;
				U(pi/2.0, 0, pi) q[0], q[1];
				ctrl @ negctrl @ U(pi, 0, pi) q[0], q[1], q[2];
			`,
			want: []string{
				"[000][  0]( 0.5000 0.0000i): 0.2500",
				"[010][  2]( 0.5000 0.0000i): 0.2500",
				"[101][  5]( 0.5000 0.0000i): 0.2500",
				"[110][  6]( 0.5000 0.0000i): 0.2500",
			},
		},
		{
			// sx**2 = x
			text: `
				qubit q;
				pow(2.0) @ U(pi/2, -pi/2, pi/2) q;
			`,
			want: []string{
				"[1][  1]( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				qubit q;
				U(pi/2, -pi/2, pi/2) q;
			`,
			want: []string{
				"[0][  0]( 0.7071 0.0000i): 0.5000",
				"[1][  1]( 0.0000-0.7071i): 0.5000",
			},
		},
		{
			text: `
				qubit[2] q;
				U(pi/2, 0, pi) q[0];
				inv @ ctrl @ U(pi, 0, pi) q[0], q[1];
			`,
			want: []string{
				"[00][  0]( 0.7071 0.0000i): 0.5000",
				"[11][  3]( 0.7071 0.0000i): 0.5000",
			},
		},
	}

	for _, c := range cases {
		qsim, _, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

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

func TestVisitor_VisitDefStatement(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
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
		{
			text: `
				qubit q;
				xm(q);
			`,
			errMsg: `subroutine "xm": undefined`,
		},
		{
			text:   "def f(qubit q) -> bit { return 1; } def f(qubit q) -> bit { return 0; }",
			errMsg: `declare subroutine "f": already declared`,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if len(env.Bit) > 0 && fmt.Sprintf("%v", env.Bit) != c.want {
			t.Errorf("got=%v, want=%v", env.Bit, c.want)
		}

		if len(env.BitArray) > 0 && fmt.Sprintf("%v", env.BitArray) != c.want {
			t.Errorf("got=%v, want=%v", env.BitArray, c.want)
		}
	}
}

func TestVisitor_VisitIfStatement(t *testing.T) {
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

func TestVisitor_VisitForStatement(t *testing.T) {
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

func TestVisitor_VisitBreakStatement(t *testing.T) {
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

func TestVisitor_VisitContinueStatement(t *testing.T) {
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

func TestVisitor_VisitWhileStatement(t *testing.T) {
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

func TestVisitor_VisitSwitchStatement(t *testing.T) {
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

func TestVisitor_VisitCastExpression(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
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
			text: "bool a = true; int b = int(a);",
			want: "map[a:true b:int(true): unexpected type: bool]",
		},
		{
			text: "bool a = true; uint b = uint(a);",
			want: "map[a:true b:uint(true): unexpected type: bool]",
		},
		{
			text: "bool a = true; float b = float(a);",
			want: "map[a:true b:float64(true): unexpected type: bool]",
		},
		{
			text: "bool a = bool(1);",
			want: `map[a:unsupported scalar type "bool"]`,
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
			t.Fail()
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitIndexedIdentifier(t *testing.T) {
	cases := []struct {
		text   string
		want   []int64
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
		if err, ok := result.(error); ok {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", c.want) {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitIndexExpression(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
	}{
		{
			text:   "int a = 1; a[0];",
			errMsg: `invalid operand "a"`,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		x := p.Program().
			StatementOrScope(1).
			Statement().
			ExpressionStatement().
			Expression()

		ctx, ok := x.(*parser.IndexExpressionContext)
		if !ok {
			t.Fatalf("index expression type=%T", x)
		}

		result := visitor.New(q.New(), environ.New()).VisitIndexExpression(ctx)
		if err, ok := result.(error); ok {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitScalarType(t *testing.T) {
	cases := []struct {
		text   string
		want   int64
		errMsg string
	}{
		{
			text: "angle a = 0;",
			want: 1,
		},
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
		if err, ok := result.(error); ok {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

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

func TestVisitor_VisitQubitType(t *testing.T) {
	cases := []struct {
		text   string
		want   int64
		errMsg string
	}{
		{
			text: "qubit q;",
			want: 1,
		},
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
		if err, ok := result.(error); ok {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

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

func TestVisitor_VisitRangeExpression(t *testing.T) {
	cases := []struct {
		text   string
		want   []int64
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
		if err, ok := result.(error); ok {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", c.want) {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitArrayType(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
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
			text: "array[int[11], 5] aa;",
			want: "map[aa:unsupported bit size 11]",
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
			text: "array[uint[11], 5] aa;",
			want: "map[aa:unsupported bit size 11]",
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
			text: "array[float[11], 5] aa;",
			want: "map[aa:unsupported bit size 11]",
		},
		{
			text: "array[angle[4], 3] aa;",
			want: `map[aa:unsupported scalar type "angle[4]"]`,
		},
		{
			text: "array[bool, 3] aa;",
			want: "map[aa:[false false false]]",
		},
		{
			text:   "array[int[64], 5] aa;array[float[64], 5] aa;",
			errMsg: `declare array "aa": already declared`,
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitArrayLiteral(t *testing.T) {
	cases := []struct {
		text   string
		want   string
		errMsg string
	}{
		{
			text: "array[int[8], 3] aa = {0, 1, 2};",
			want: "map[aa:[0 1 2]]",
		},
	}

	for _, c := range cases {
		_, env, err := visitor.Run(c.text)
		if err != nil {
			if err.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", err, c.errMsg)
			}

			continue
		}

		if fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitMeasureExpression(t *testing.T) {
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

func TestVisitor_VisitMeasureExpression_assign(t *testing.T) {
	cases := []struct {
		text   string
		errMsg string
	}{
		{
			text:   "int a = 1; bit c = measure a;",
			errMsg: `invalid operand "a"`,
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

func TestVisitor_unimplemented(t *testing.T) {
	v := visitor.New(q.New(), environ.New())

	cases := []struct {
		name string
		err  error
	}{
		{name: "VisitPragma", err: v.VisitPragma(&parser.PragmaContext{}).(error)},
		{name: "VisitAnnotation", err: v.VisitAnnotation(&parser.AnnotationContext{}).(error)},
		{name: "VisitDurationofExpression", err: v.VisitDurationofExpression(&parser.DurationofExpressionContext{}).(error)},
		{name: "VisitSetExpression", err: v.VisitSetExpression(&parser.SetExpressionContext{}).(error)},
		{name: "VisitArrayReferenceType", err: v.VisitArrayReferenceType(&parser.ArrayReferenceTypeContext{}).(error)},
		{name: "VisitDefcalArgumentDefinitionList", err: v.VisitDefcalArgumentDefinitionList(&parser.DefcalArgumentDefinitionListContext{}).(error)},
		{name: "VisitDefcalArgumentDefinition", err: v.VisitDefcalArgumentDefinition(&parser.DefcalArgumentDefinitionContext{}).(error)},
		{name: "VisitDefcalTarget", err: v.VisitDefcalTarget(&parser.DefcalTargetContext{}).(error)},
		{name: "VisitDefcalOperandList", err: v.VisitDefcalOperandList(&parser.DefcalOperandListContext{}).(error)},
		{name: "VisitDefcalOperand", err: v.VisitDefcalOperand(&parser.DefcalOperandContext{}).(error)},
		{name: "VisitIoDeclarationStatement", err: v.VisitIoDeclarationStatement(&parser.IoDeclarationStatementContext{}).(error)},
		{name: "VisitExternStatement", err: v.VisitExternStatement(&parser.ExternStatementContext{}).(error)},
		{name: "VisitExternArgumentList", err: v.VisitExternArgumentList(&parser.ExternArgumentListContext{}).(error)},
		{name: "VisitExternArgument", err: v.VisitExternArgument(&parser.ExternArgumentContext{}).(error)},
		{name: "VisitCalStatement", err: v.VisitCalStatement(&parser.CalStatementContext{}).(error)},
		{name: "VisitDefcalStatement", err: v.VisitDefcalStatement(&parser.DefcalStatementContext{}).(error)},
		{name: "VisitCalibrationGrammarStatement", err: v.VisitCalibrationGrammarStatement(&parser.CalibrationGrammarStatementContext{}).(error)},
		{name: "VisitBarrierStatement", err: v.VisitBarrierStatement(&parser.BarrierStatementContext{}).(error)},
		{name: "VisitBoxStatement", err: v.VisitBoxStatement(&parser.BoxStatementContext{}).(error)},
		{name: "VisitDelayStatement", err: v.VisitDelayStatement(&parser.DelayStatementContext{}).(error)},
	}

	for _, c := range cases {
		if !errors.Is(c.err, visitor.ErrNotImplemented) {
			t.Fatalf("got=%v", c.err)
		}

		if !strings.Contains(c.err.Error(), c.name) {
			t.Fatalf("got=%v", c.err)
		}
	}
}
