package visitor_test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func ExampleVisitor_VisitVersion() {
	text := "OPENQASM 3.0;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := v.Visit(tree); err != nil {
		fmt.Println(err)
	}

	fmt.Println(env.Version)

	// Output:
	// (program (version OPENQASM 3.0 ;) <EOF>)
	// 3.0
}

func ExampleVisitor_VisitPragma() {
	text := `
	pragma qiskit.simulator noise model "qpu1.noise";
	`

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	fmt.Println(v.Visit(tree))

	// Output:
	// (program (statementOrScope (statement (pragma pragma qiskit.simulator noise model "qpu1.noise";))) <EOF>)
	// qiskit.simulator noise model "qpu1.noise";
}

func ExampleVisitor_VisitResetStatement() {
	text := `
	qubit q;
	U(pi/2, 0, pi) q;
	reset q;
	`

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := v.Visit(tree); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (resetStatement reset (gateOperand (indexedIdentifier q)) ;))) <EOF>)
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func ExampleVisitor_VisitIncludeStatement() {
	text := `
	include "../_testdata/stdgates.qasm";
	qubit q;
	h q;
	`

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := v.Visit(tree); err != nil {
		fmt.Println(err)
	}

	fmt.Println(slices.Sorted(maps.Keys(env.Gate)))
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (includeStatement include "../_testdata/stdgates.qasm" ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement h (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)
	// [cx h i x y z]
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func TestVisitor_VisitConstDeclarationStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: "const int a = 42;",
			tree: "(program (statementOrScope (statement (constDeclarationStatement const (scalarType int) a = (declarationExpression (expression 42)) ;))) <EOF>)",
			want: "map[a:42]",
		},
		{
			text: "const uint N = 3 * 5;",
			tree: "(program (statementOrScope (statement (constDeclarationStatement const (scalarType uint) N = (declarationExpression (expression (expression 3) * (expression 5))) ;))) <EOF>)",
			want: "map[N:15]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.Const) > 0 && fmt.Sprintf("%v", env.Const) != c.want {
			t.Errorf("got=%v, want=%v", env.Const, c.want)
		}
	}
}

func TestVisitor_VisitClassicalDeclarationStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: "bit c;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType bit) c ;))) <EOF>)",
			want: "map[c:[0]]",
		},
		{
			text: "bit[4] c;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType bit (designator [ (expression 4) ])) c ;))) <EOF>)",
			want: "map[c:[0 0 0 0]]",
		},
		{
			text: `bit[8] a = "10001111";`,
			tree: `(program (statementOrScope (statement (classicalDeclarationStatement (scalarType bit (designator [ (expression 8) ])) a = (declarationExpression (expression "10001111")) ;))) <EOF>)`,
			want: "map[a:[1 0 0 0 1 1 1 1]]",
		},
		{
			text: `qubit q; U(pi, 0, pi) q; bit c = measure q;`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit) c = (declarationExpression (measureExpression measure (gateOperand (indexedIdentifier q)))) ;))) <EOF>)",
			want: "map[c:[1]]",
		},
		{
			text: "float ratio = 22 / 7;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType float) ratio = (declarationExpression (expression (expression 22) / (expression 7))) ;))) <EOF>)",
			want: "map[ratio:3.142857142857143]",
		},
		{
			text: "int ans = 42;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) ans = (declarationExpression (expression 42)) ;))) <EOF>)",
			want: "map[ans:42]",
		},
		{
			text: "uint ans = 42;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType uint) ans = (declarationExpression (expression 42)) ;))) <EOF>)",
			want: "map[ans:42]",
		},
		{
			text: "float z;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType float) z ;))) <EOF>)",
			want: "map[z:0]",
		},
		{
			text: "int z;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) z ;))) <EOF>)",
			want: "map[z:0]",
		},
		{
			text: "uint z;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType uint) z ;))) <EOF>)",
			want: "map[z:0]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.ClassicalBit) > 0 && fmt.Sprintf("%v", env.ClassicalBit) != c.want {
			t.Errorf("got=%v, want=%v", env.ClassicalBit, c.want)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitQuantumDeclarationStatement(t *testing.T) {
	cases := []struct {
		text   string
		tree   string
		want   string
		errMsg string
	}{
		{
			text: "qubit q;",
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) <EOF>)",
			want: "map[q:[0]]",
		},
		{
			text: "qubit[2] q;",
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) <EOF>)",
			want: "map[q:[0 1]]",
		},
		{
			text: "qubit q0; qubit[2] q1; qubit[3] q3; qubit q4;",
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q0 ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q1 ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 3) ])) q3 ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q4 ;))) <EOF>)",
			want: "map[q0:[0] q1:[1 2] q3:[3 4 5] q4:[6]]",
		},
		{
			text:   "qubit q; qubit q;",
			tree:   "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) <EOF>)",
			errMsg: "identifier=q: already declared",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			if ret.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", ret, c.errMsg)
			}
		default:
			if fmt.Sprintf("%v", env.Qubit) != c.want {
				t.Errorf("got=%v, want=%v", env.Qubit, c.want)
			}
		}
	}
}

func TestVisitor_VisitAliasDeclarationStatement(t *testing.T) {
	cases := []struct {
		text   string
		tree   string
		want   string
		errMsg string
	}{
		{
			text: `
				qubit[5] q;
				let myreg = q[1:4];
			`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 5) ])) q ;))) (statementOrScope (statement (aliasDeclarationStatement let myreg = (aliasExpression (expression (expression q) (indexOperator [ (rangeExpression (expression 1) : (expression 4)) ]))) ;))) <EOF>)",
			want: "map[myreg:[1 2 3] q:[0 1 2 3 4]]",
		},
		{
			text: `
				qubit[2] one;
				qubit[10] two;
				let concatenated = one ++ two;
			`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) one ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 10) ])) two ;))) (statementOrScope (statement (aliasDeclarationStatement let concatenated = (aliasExpression (expression one) ++ (expression two)) ;))) <EOF>)",
			want: "map[concatenated:[0 1 2 3 4 5 6 7 8 9 10 11] one:[0 1] two:[2 3 4 5 6 7 8 9 10 11]]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			if ret.Error() != c.errMsg {
				t.Errorf("got=%v, want=%v", ret, c.errMsg)
			}
		default:
			if fmt.Sprintf("%v", env.Qubit) != c.want {
				t.Errorf("got=%v, want=%v", env.Qubit, c.want)
			}
		}
	}
}

func TestVisitor_VisitAssignmentStatement(t *testing.T) {
	type Want struct {
		classicalBit []string
		qubit        []string
		variable     []string
	}

	cases := []struct {
		text string
		tree string
		want Want
	}{
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
				c = measure q;
			`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit (designator [ (expression 2) ])) c ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) (statementOrScope (statement (assignmentStatement (indexedIdentifier c) = (measureExpression measure (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: Want{
				classicalBit: []string{
					"map[c:[0 0]]",
					"map[c:[1 1]]",
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
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit (designator [ (expression 2) ])) c ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) (statementOrScope (statement (assignmentStatement (indexedIdentifier c (indexOperator [ (expression 0) ])) = (measureExpression measure (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) (statementOrScope (statement (assignmentStatement (indexedIdentifier c (indexOperator [ (expression 1) ])) = (measureExpression measure (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) <EOF>)",
			want: Want{
				classicalBit: []string{
					"map[c:[0 0]]",
					"map[c:[1 1]]",
				},
				qubit: []string{
					"[[00][  0]( 1.0000 0.0000i): 1.0000]",
					"[[11][  3]( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
		{
			text: `
				int ans = 42;
				ans = ans * 2;
			`,
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) ans = (declarationExpression (expression 42)) ;))) (statementOrScope (statement (assignmentStatement (indexedIdentifier ans) = (expression (expression ans) * (expression 2)) ;))) <EOF>)",
			want: Want{
				variable: []string{
					"map[ans:84]",
				},
			},
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(c.want.classicalBit) != 0 {
			var found bool
			for _, w := range c.want.classicalBit {
				if fmt.Sprintf("%v", env.ClassicalBit) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", env.ClassicalBit, c.want.classicalBit)
			}
		}

		if len(c.want.qubit) != 0 {
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

		if len(c.want.variable) != 0 {
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
		classicalBit []string
		qubit        []string
	}

	cases := []struct {
		text string
		tree string
		want Want
	}{
		{
			text: `
				qubit[2] q;
				bit[2] c;
				U(pi/2.0, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
				measure q -> c;
			`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit (designator [ (expression 2) ])) c ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) (statementOrScope (statement (measureArrowAssignmentStatement (measureExpression measure (gateOperand (indexedIdentifier q))) -> (indexedIdentifier c) ;))) <EOF>)",
			want: Want{
				classicalBit: []string{
					"map[c:[0 0]]",
					"map[c:[1 1]]",
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
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit (designator [ (expression 2) ])) c ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) (statementOrScope (statement (measureArrowAssignmentStatement (measureExpression measure (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) -> (indexedIdentifier c (indexOperator [ (expression 0) ])) ;))) (statementOrScope (statement (measureArrowAssignmentStatement (measureExpression measure (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) -> (indexedIdentifier c (indexOperator [ (expression 1) ])) ;))) <EOF>)",
			want: Want{
				classicalBit: []string{
					"map[c:[0 0]]",
					"map[c:[1 1]]",
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
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) (statementOrScope (statement (measureArrowAssignmentStatement (measureExpression measure (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: Want{
				qubit: []string{
					"[[00][  0]( 1.0000 0.0000i): 1.0000]",
					"[[11][  3]( 1.0000 0.0000i): 1.0000]",
				},
			},
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}
		if len(c.want.classicalBit) != 0 {
			var found bool
			for _, w := range c.want.classicalBit {
				if fmt.Sprintf("%v", env.ClassicalBit) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", env.ClassicalBit, c.want.classicalBit)
			}
		}

		if len(c.want.qubit) != 0 {
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
	}
}

func TestVisitor_VisitMultiplicativeExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: "1 * 3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 1) * (expression 3)) ;))) <EOF>)",
			want: "3",
		},
		{
			text: "1.0 * 3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 1.0) * (expression 3)) ;))) <EOF>)",
			want: "3",
		},
		{
			text: "4 / 2;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 4) / (expression 2)) ;))) <EOF>)",
			want: "2",
		},
		{
			text: "10 % 3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 10) % (expression 3)) ;))) <EOF>)",
			want: "1",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitAdditiveExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: "1 + 3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 1) + (expression 3)) ;))) <EOF>)",
			want: "4",
		},
		{
			text: "1.0 - 3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 1.0) - (expression 3)) ;))) <EOF>)",
			want: "-2",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitParenthesisExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: "(1 + 3) * 4;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression ( (expression (expression 1) + (expression 3)) )) * (expression 4)) ;))) <EOF>)",
			want: "16",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitCallExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: "sin(1.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression sin ( (expressionList (expression 1.0)) )) ;))) <EOF>)",
			want: "0.8414709848078965",
		},
		{
			text: "cos(1.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression cos ( (expressionList (expression 1.0)) )) ;))) <EOF>)",
			want: "0.5403023058681398",
		},
		{
			text: "tan(1.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression tan ( (expressionList (expression 1.0)) )) ;))) <EOF>)",
			want: "1.557407724654902",
		},
		{
			text: "arcsin(0.5);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression arcsin ( (expressionList (expression 0.5)) )) ;))) <EOF>)",
			want: "0.5235987755982989",
		},
		{
			text: "arccos(0.5);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression arccos ( (expressionList (expression 0.5)) )) ;))) <EOF>)",
			want: "1.0471975511965976",
		},
		{
			text: "arctan(1.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression arctan ( (expressionList (expression 1.0)) )) ;))) <EOF>)",
			want: "0.7853981633974483",
		},
		{
			text: "ceiling(1.1);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression ceiling ( (expressionList (expression 1.1)) )) ;))) <EOF>)",
			want: "2",
		},
		{
			text: "floor(1.1);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression floor ( (expressionList (expression 1.1)) )) ;))) <EOF>)",
			want: "1",
		},
		{
			text: "sqrt(2.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression sqrt ( (expressionList (expression 2.0)) )) ;))) <EOF>)",
			want: "1.4142135623730951",
		},
		{
			text: "exp(1.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression exp ( (expressionList (expression 1.0)) )) ;))) <EOF>)",
			want: "2.718281828459045",
		},
		{
			text: "log(2.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression log ( (expressionList (expression 2.0)) )) ;))) <EOF>)",
			want: "0.6931471805599453",
		},
		{
			text: "mod(10.0, 3.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression mod ( (expressionList (expression 10.0) , (expression 3.0)) )) ;))) <EOF>)",
			want: "1",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		default:
			if fmt.Sprintf("%v", ret) != c.want {
				t.Errorf("got=%v, want=%v", ret, c.want)
			}
		}
	}
}

func TestVisitor_VisitPowerExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: "2**3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 2) ** (expression 3)) ;))) <EOF>)",
			want: "8",
		},
		{
			text: "2**0.5;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 2) ** (expression 0.5)) ;))) <EOF>)",
			want: "1.4142135623730951",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitLogicalOrExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want bool
	}{
		{
			text: "false || false;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression false) || (expression false)) ;))) <EOF>)",
			want: false,
		},
		{
			text: "true || false;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression true) || (expression false)) ;))) <EOF>)",
			want: true,
		},
		{
			text: "false || true;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression false) || (expression true)) ;))) <EOF>)",
			want: true,
		},
		{
			text: "true || true;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression true) || (expression true)) ;))) <EOF>)",
			want: true,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitLogicalAndExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want bool
	}{
		{
			text: "false && false;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression false) && (expression false)) ;))) <EOF>)",
			want: false,
		},
		{
			text: "true && false;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression true) && (expression false)) ;))) <EOF>)",
			want: false,
		},
		{
			text: "false && true;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression false) && (expression true)) ;))) <EOF>)",
			want: false,
		},
		{
			text: "true && true;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression true) && (expression true)) ;))) <EOF>)",
			want: true,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitBitwiseAndExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want int64
	}{
		{
			text: "10 & 12;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 10) & (expression 12)) ;))) <EOF>)",
			want: 8,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if result.(int64) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitBitwiseOrExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want int64
	}{
		{
			text: "10 | 12;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 10) | (expression 12)) ;))) <EOF>)",
			want: 14,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if result.(int64) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitBitwiseXorExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want int64
	}{
		{
			text: "10 ^ 12;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 10) ^ (expression 12)) ;))) <EOF>)",
			want: 6,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if result.(int64) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitBitshiftExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want int64
	}{
		{
			text: "11 << 2;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 11) << (expression 2)) ;))) <EOF>)",
			want: 44,
		},
		{
			text: "11 >> 1;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 11) >> (expression 1)) ;))) <EOF>)",
			want: 5,
		},
		{
			text: "11 >> 1 << 1;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression (expression 11) >> (expression 1)) << (expression 1)) ;))) <EOF>)",
			want: 10,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if result.(int64) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitEqualityExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want bool
	}{
		{
			text: "10.0 == 2.0 * 5;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 10.0) == (expression (expression 2.0) * (expression 5))) ;))) <EOF>)",
			want: true,
		},
		{
			text: "10.0 != 2.0 * 5;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 10.0) != (expression (expression 2.0) * (expression 5))) ;))) <EOF>)",
			want: false,
		},
		{
			text: "(1 == 1) == true;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression ( (expression (expression 1) == (expression 1)) )) == (expression true)) ;))) <EOF>)",
			want: true,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitUnaryExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want bool
	}{
		{
			text: "!false;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression ! (expression false)) ;))) <EOF>)",
			want: true,
		},
		{
			text: "!true;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression ! (expression true)) ;))) <EOF>)",
			want: false,
		},
		{
			text: "-1 == 1;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression - (expression 1)) == (expression 1)) ;))) <EOF>)",
			want: false,
		},
		{
			text: "-1.0 != 1.0;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression - (expression 1.0)) != (expression 1.0)) ;))) <EOF>)",
			want: true,
		},
		{
			text: "!(-1.0 != 1.0);",
			tree: "(program (statementOrScope (statement (expressionStatement (expression ! (expression ( (expression (expression - (expression 1.0)) != (expression 1.0)) ))) ;))) <EOF>)",
			want: false,
		},
		{
			text: "~5 == -6;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression ~ (expression 5)) == (expression - (expression 6))) ;))) <EOF>)",
			want: true,
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if result.(bool) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitComparisonExpression(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: "2.0 < 3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 2.0) < (expression 3)) ;))) <EOF>)",
			want: "true",
		},
		{
			text: "2 <= 2;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 2) <= (expression 2)) ;))) <EOF>)",
			want: "true",
		},
		{
			text: "2.0 > 3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 2.0) > (expression 3)) ;))) <EOF>)",
			want: "false",
		},
		{
			text: "2 >= 3;",
			tree: "(program (statementOrScope (statement (expressionStatement (expression (expression 2) >= (expression 3)) ;))) <EOF>)",
			want: "false",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		result := v.Visit(tree)
		if fmt.Sprintf("%v", result) != c.want {
			t.Errorf("got=%v, want=%v", result, c.want)
		}
	}
}

func TestVisitor_VisitGateCallStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want []string
	}{
		{
			text: "qubit q; U(pi, 0, pi) q;",
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[1][  1]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: "qubit q; U(pi/2, 0, pi) q;",
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[0][  0]( 0.7071 0.0000i): 0.5000",
				"[1][  1]( 0.7071 0.0000i): 0.5000",
			},
		},
		{
			text: "qubit q; gphase(-π/2);",
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement gphase ( (expressionList (expression (expression - (expression π)) / (expression 2))) ) ;))) <EOF>)",
			want: []string{
				"[0][  0]( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				gate h q0 { U(pi/2, 0, pi) q0; }
				gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }
				qubit[2] q;
				h q[0];
				cx q[0], q[1];
			`,
			tree: "(program (statementOrScope (statement (gateStatement gate h (identifierList q0) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) })))) (statementOrScope (statement (gateStatement gate cx (identifierList q0 , q1) (scope { (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0)) , (gateOperand (indexedIdentifier q1))) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) (statementOrScope (statement (gateCallStatement h (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) (statementOrScope (statement (gateCallStatement cx (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (gateStatement gate x (identifierList q0) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) })))) (statementOrScope (statement (gateStatement gate y (identifierList q0) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression (expression pi) / (expression 2)) , (expression (expression pi) / (expression 2))) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement x (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (gateCallStatement y (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (gateStatement gate x (identifierList q) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) (statementOrScope (statement (gateStatement gate y (identifierList q) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression (expression pi) / (expression 2)) , (expression (expression pi) / (expression 2))) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) (statementOrScope (statement (gateStatement gate xy (identifierList q) (scope { (statementOrScope (statement (gateCallStatement x (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (gateCallStatement y (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement xy (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (gateStatement gate u ( (identifierList p0 , p1 , p2) ) (identifierList q) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression p0) , (expression p1) , (expression p2)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement u ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[1][  1]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				qubit q;
				u(pi, 0, pi) q;
				inv @ u(pi, 0, pi) q;
			`,
			tree: "(program (statementOrScope (statement (gateStatement gate u ( (identifierList p0 , p1 , p2) ) (identifierList q) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression p0) , (expression p1) , (expression p2)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement u ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier inv @) u ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[0][  0]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				gate invu(p0, p1, p2) q { inv @ u(p0, p1, p2) q; }
				qubit q;
				u(1, 2, 3) q;
				invu(1, 2, 3) q;
			`,
			tree: "(program (statementOrScope (statement (gateStatement gate u ( (identifierList p0 , p1 , p2) ) (identifierList q) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression p0) , (expression p1) , (expression p2)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) (statementOrScope (statement (gateStatement gate invu ( (identifierList p0 , p1 , p2) ) (identifierList q) (scope { (statementOrScope (statement (gateCallStatement (gateModifier inv @) u ( (expressionList (expression p0) , (expression p1) , (expression p2)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement u ( (expressionList (expression 1) , (expression 2) , (expression 3)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (gateCallStatement invu ( (expressionList (expression 1) , (expression 2) , (expression 3)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[0][  0]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			text: `
				gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
				const int n = 3;
				qubit q;
				pow(n) @ u(pi, 0, pi) q;
			`,
			tree: "(program (statementOrScope (statement (gateStatement gate u ( (identifierList p0 , p1 , p2) ) (identifierList q) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression p0) , (expression p1) , (expression p2)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) (statementOrScope (statement (constDeclarationStatement const (scalarType int) n = (declarationExpression (expression 3)) ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement (gateModifier pow ( (expression n) ) @) u ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[1][  1]( 1.0000 0.0000i): 1.0000",
			},
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)
		if err := v.Visit(tree); err != nil {
			panic(err)
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
		text string
		tree string
		want []string
	}{
		{
			text: `
				qubit q0;
				qubit q1;
				U(pi/2.0, 0, pi) q0;
				ctrl @ U(pi, 0, pi) q0, q1;
			`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q0 ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q1 ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0)) , (gateOperand (indexedIdentifier q1))) ;))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 3) ])) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 2) ])))) ;))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q0 ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q1 ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier negctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0)) , (gateOperand (indexedIdentifier q1))) ;))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression tau) , (expression euler)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier inv @) U ( (expressionList (expression pi) , (expression tau) , (expression euler)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[0][  0]( 1.0000 0.0000i): 1.0000",
			},
		},
		{
			// sx**2 = x
			text: `
				qubit q;
				pow(2) @ U(pi/2, -pi/2, pi/2) q;
			`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement (gateModifier pow ( (expression 2) ) @) U ( (expressionList (expression (expression pi) / (expression 2)) , (expression (expression - (expression pi)) / (expression 2)) , (expression (expression pi) / (expression 2))) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[1][  1]( 0.0000-1.0000i): 1.0000",
			},
		},
		{
			text: `
				const float half = pi / 2;
				qubit q;
				pow(2) @ U(half, -half, half) q;
			`,
			tree: "(program (statementOrScope (statement (constDeclarationStatement const (scalarType float) half = (declarationExpression (expression (expression pi) / (expression 2))) ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement (gateModifier pow ( (expression 2) ) @) U ( (expressionList (expression half) , (expression - (expression half)) , (expression half)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: []string{
				"[1][  1]( 0.0000-1.0000i): 1.0000",
			},
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)
		if err := v.Visit(tree); err != nil {
			panic(err)
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
		text string
		tree string
		want string
	}{
		{
			text: `
				gate x q0 { U(pi, 0, pi) q0; }
				def xm(qubit q1) -> bit { x q1; return measure q1; }
				qubit q;
				bit c = xm(q);
			`,
			tree: "(program (statementOrScope (statement (gateStatement gate x (identifierList q0) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) })))) (statementOrScope (statement (defStatement def xm ( (argumentDefinitionList (argumentDefinition (qubitType qubit) q1)) ) (returnSignature -> (scalarType bit)) (scope { (statementOrScope (statement (gateCallStatement x (gateOperandList (gateOperand (indexedIdentifier q1))) ;))) (statementOrScope (statement (returnStatement return (measureExpression measure (gateOperand (indexedIdentifier q1))) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit) c = (declarationExpression (expression xm ( (expressionList (expression q)) ))) ;))) <EOF>)",
			want: "map[c:[1]]",
		},
		{
			text: `
				gate x q0 { U(pi, 0, pi) q0; }
				def xm(qubit q1) -> bit { x q1; bit m = measure q1; return m;}
				qubit q;
				bit c = xm(q);
			`,
			tree: "(program (statementOrScope (statement (gateStatement gate x (identifierList q0) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) })))) (statementOrScope (statement (defStatement def xm ( (argumentDefinitionList (argumentDefinition (qubitType qubit) q1)) ) (returnSignature -> (scalarType bit)) (scope { (statementOrScope (statement (gateCallStatement x (gateOperandList (gateOperand (indexedIdentifier q1))) ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit) m = (declarationExpression (measureExpression measure (gateOperand (indexedIdentifier q1)))) ;))) (statementOrScope (statement (returnStatement return (expression m) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit) c = (declarationExpression (expression xm ( (expressionList (expression q)) ))) ;))) <EOF>)",
			want: "map[c:[1]]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.ClassicalBit) > 0 && fmt.Sprintf("%v", env.ClassicalBit) != c.want {
			t.Errorf("got=%v, want=%v", env.ClassicalBit, c.want)
		}
	}
}

func TestVisitor_VisitIfStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: `
				int a = 10;
				if (a == 10) { a = 20; }
			`,
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 10)) ;))) (statementOrScope (statement (ifStatement if ( (expression (expression a) == (expression 10)) ) (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression 20) ;))) }))))) <EOF>)",
			want: "map[a:20]",
		},
		{
			text: `
				int a = 10;
				if (a == 20) { a = 100; } else { a = 30; }
			`,
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 10)) ;))) (statementOrScope (statement (ifStatement if ( (expression (expression a) == (expression 20)) ) (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression 100) ;))) })) else (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression 30) ;))) }))))) <EOF>)",
			want: "map[a:30]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitForStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: `
				int a = 0;
				for int i in [0:10] {
					a = a + 1;
				}
			`,
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (forStatement for (scalarType int) i in [ (rangeExpression (expression 0) : (expression 10)) ] (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 1)) ;))) }))))) <EOF>)",
			want: "map[a:10]",
		},
		{
			text: `
				int a = 0;
				for int i in [0:10] {
					a = a + i;
				}
			`,
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (forStatement for (scalarType int) i in [ (rangeExpression (expression 0) : (expression 10)) ] (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression i)) ;))) }))))) <EOF>)",
			want: "map[a:45]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitBreakStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
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
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (forStatement for (scalarType int) i in [ (rangeExpression (expression 0) : (expression 10)) ] (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 1)) ;))) (statementOrScope (statement (ifStatement if ( (expression (expression a) > (expression 10)) ) (statementOrScope (scope { (statementOrScope (statement (breakStatement break ;))) }))))) (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 1)) ;))) }))))) <EOF>)",
			want: "map[a:11]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitContinueStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: `
				int a = 0;
				for int i in [0:10] {
					a = a + 1;
					if ( a > 10 ) {
						continue;
					}
					a = a + 1;
				}
			`,
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (forStatement for (scalarType int) i in [ (rangeExpression (expression 0) : (expression 10)) ] (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 1)) ;))) (statementOrScope (statement (ifStatement if ( (expression (expression a) > (expression 10)) ) (statementOrScope (scope { (statementOrScope (statement (continueStatement continue ;))) }))))) (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 1)) ;))) }))))) <EOF>)",
			want: "map[a:15]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitWhileStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want string
	}{
		{
			text: `
				int a = 0;
				while (a < 10) {
					a = a + 1;
				}
			`,
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (whileStatement while ( (expression (expression a) < (expression 10)) ) (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 1)) ;))) }))))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (whileStatement while ( (expression (expression a) < (expression 100)) ) (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 1)) ;))) (statementOrScope (statement (ifStatement if ( (expression (expression a) > (expression 10)) ) (statementOrScope (scope { (statementOrScope (statement (breakStatement break ;))) }))))) (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 10)) ;))) }))))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (whileStatement while ( (expression (expression a) < (expression 100)) ) (statementOrScope (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 1)) ;))) (statementOrScope (statement (ifStatement if ( (expression (expression a) < (expression 10)) ) (statementOrScope (scope { (statementOrScope (statement (continueStatement continue ;))) }))))) (statementOrScope (statement (assignmentStatement (indexedIdentifier a) = (expression (expression a) + (expression 10)) ;))) }))))) <EOF>)",
			want: "map[a:108]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}

func TestVisitor_VisitSwitchStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
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
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 15)) ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType int) b = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (switchStatement switch ( (expression a) ) { (switchCaseItem case (expressionList (expression 15)) (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier b) = (expression 15) ;))) })) (switchCaseItem default (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier b) = (expression - (expression 1)) ;))) })) }))) <EOF>)",
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
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) a = (declarationExpression (expression 20)) ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType int) b = (declarationExpression (expression 0)) ;))) (statementOrScope (statement (switchStatement switch ( (expression a) ) { (switchCaseItem case (expressionList (expression 15)) (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier b) = (expression 15) ;))) })) (switchCaseItem default (scope { (statementOrScope (statement (assignmentStatement (indexedIdentifier b) = (expression - (expression 1)) ;))) })) }))) <EOF>)",
			want: "map[a:20 b:-1]",
		},
	}

	for _, c := range cases {
		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(c.text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		tree := p.Program()
		if tree.ToStringTree(nil, p) != c.tree {
			t.Errorf("got=%v, want=%v", tree.ToStringTree(nil, p), c.tree)
		}

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		switch ret := v.Visit(tree).(type) {
		case error:
			panic(ret)
		}

		if len(env.Variable) > 0 && fmt.Sprintf("%v", env.Variable) != c.want {
			t.Errorf("got=%v, want=%v", env.Variable, c.want)
		}
	}
}
