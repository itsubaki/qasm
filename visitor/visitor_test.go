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

func ExampleVisitor_VisitGateStatement() {
	text := `
	gate u(p0, p1, p2) q { U(p0, p1, p2) q; }
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
		panic(ret)
	}

	for _, g := range env.Gate {
		fmt.Print(g.Name, g.Params, g.QArgs, " > ")
		for _, s := range g.Body {
			fmt.Print(s.GetText())
		}
		fmt.Println()
	}

	// Output:
	// (program (statementOrScope (statement (gateStatement gate u ( (identifierList p0 , p1 , p2) ) (identifierList q) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression p0) , (expression p1) , (expression p2)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) })))) <EOF>)
	// u[p0 p1 p2] [q] > U(p0,p1,p2)q;
}

func ExampleVisitor_VisitGateCallStatement() {
	text := `
	gate h q0 { U(pi/2, 0, pi) q0; }
	gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }

	qubit q0;
	qubit q1;
	h q0;
	cx q0, q1;
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
		panic(ret)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (gateStatement gate h (identifierList q0) (scope { (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) })))) (statementOrScope (statement (gateStatement gate cx (identifierList q0 , q1) (scope { (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0)) , (gateOperand (indexedIdentifier q1))) ;))) })))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q0 ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q1 ;))) (statementOrScope (statement (gateCallStatement h (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) (statementOrScope (statement (gateCallStatement cx (gateOperandList (gateOperand (indexedIdentifier q0)) , (gateOperand (indexedIdentifier q1))) ;))) <EOF>)
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
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
			text: "float zero;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType float) zero ;))) <EOF>)",
			want: "map[zero:0]",
		},
		{
			text: "int zero;",
			tree: "(program (statementOrScope (statement (classicalDeclarationStatement (scalarType int) zero ;))) <EOF>)",
			want: "map[zero:0]",
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
		{
			text: `
				qubit[2] q;
				U(pi/2.0, 0, pi) q[0];
				ctrl @ U(pi, 0, pi) q[0], q[1];
				measure q;
			`,
			tree: "(program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ]))) , (gateOperand (indexedIdentifier q (indexOperator [ (expression 1) ])))) ;))) (statementOrScope (statement (measureArrowAssignmentStatement (measureExpression measure (gateOperand (indexedIdentifier q))) ;))) <EOF>)",
			want: [][]string{
				{
					"map[]",
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
