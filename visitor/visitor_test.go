package visitor_test

import (
	"fmt"
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

func ExampleVisitor_VisitQuantumDeclarationStatement() {
	text := `
	qubit q;
	qubit[2] q2;
	`

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	tree := p.Program()

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := v.Visit(tree); err != nil {
		fmt.Println(err)
	}

	fmt.Println(env.Qubit)

	// Output:
	// map[q:[0] q2:[1 2]]
}

func ExampleVisitor_VisitQuantumDeclarationStatement_errAlreadyDeclared() {
	text := `
	qubit q;
	qubit q;
	`

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	tree := p.Program()

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	fmt.Println(tree.Accept(v))

	// Output:
	// identifier=q: already declared
}

func ExampleVisitor_VisitClassicalDeclarationStatement_bit() {
	text := `
	bit c;
	bit[8] a = "10001111";
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

	fmt.Println(env.ClassicalBit)

	// Output:
	// (program (statementOrScope (statement (classicalDeclarationStatement (scalarType bit) c ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit (designator [ (expression 8) ])) a = (declarationExpression (expression "10001111")) ;))) <EOF>)
	// map[a:[1 0 0 0 1 1 1 1] c:[0]]
}

func ExampleVisitor_VisitDeclarationExpression_measure() {
	text := `
	qubit q;
	U(pi, 0, pi) q;
	bit c = measure q;
	`

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	switch ret := v.Visit(tree).(type) {
	case error:
		fmt.Println(ret)
	default:
		fmt.Println(qsim.State(q.Qubit(0)))
		fmt.Println(env.ClassicalBit["c"])
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (classicalDeclarationStatement (scalarType bit) c = (declarationExpression (measureExpression measure (gateOperand (indexedIdentifier q)))) ;))) <EOF>)
	// [[1][  1]( 1.0000 0.0000i): 1.0000]
	// [1]
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

func TestVisitor_VisitAssignmentStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want [][]string
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
			want: [][]string{
				{
					"map[c:[0 0]]",
					"map[c:[1 1]]",
				},
				{
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
			want: [][]string{
				{
					"map[c:[0 0]]",
					"map[c:[1 1]]",
				},
				{
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

		{
			var found bool
			for _, w := range c.want[0] {
				if fmt.Sprintf("%v", env.ClassicalBit) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", env.ClassicalBit, c.want[0])
			}
		}

		{
			var found bool
			for _, w := range c.want[1] {
				if fmt.Sprintf("%v", qsim.State()) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", qsim.State(), c.want[1])
			}
		}

	}
}

func TestVisitor_VisitMeasureArrowAssignmentStatement(t *testing.T) {
	cases := []struct {
		text string
		tree string
		want [][]string
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
			want: [][]string{
				{
					"map[c:[0 0]]",
					"map[c:[1 1]]",
				},
				{
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
			want: [][]string{
				{
					"map[c:[0 0]]",
					"map[c:[1 1]]",
				},
				{
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

		{
			var found bool
			for _, w := range c.want[0] {
				if fmt.Sprintf("%v", env.ClassicalBit) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", env.ClassicalBit, c.want[0])
			}
		}

		{
			var found bool
			for _, w := range c.want[1] {
				if fmt.Sprintf("%v", qsim.State()) == w {
					found = true
				}
			}

			if !found {
				t.Errorf("got=%v, want=%v", qsim.State(), c.want[1])
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
			// sx **2 = x
			text: `
				qubit q;
				pow(2) @ U(pi/2, -pi/2, pi/2) q;
			`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement (gateModifier pow ( (expression 2) ) @) U ( (expressionList (expression (expression pi) / (expression 2)) , (expression (expression - (expression pi)) / (expression 2)) , (expression (expression pi) / (expression 2))) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
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
