package visitor_test

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func ExampleVisitor_VisitQuantumDeclarationStatement() {
	text := "qubit q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)
	tree.Accept(v)

	fmt.Println(env.Qubit)

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) <EOF>)
	// map[q:[0]]
}

func ExampleVisitor_VisitQuantumDeclarationStatement_register() {
	text := "qubit[2] q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) <EOF>)
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func ExampleVisitor_VisitQuantumDeclarationStatement_errAlreadyDeclared() {
	text := "qubit q; qubit q;"

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

func ExampleVisitor_VisitGateCallStatement_x() {
	text := "qubit q; U(pi, 0, pi) q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)
	// [1][  1]( 1.0000 0.0000i): 1.0000
}

func ExampleVisitor_VisitGateCallStatement_h() {
	text := "qubit q; U(pi/2, 0, pi) q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func ExampleVisitor_VisitMultiplicativeExpression() {
	text := "qubit q; U(pi/2, 0.0*3/3, pi) q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2)) , (expression (expression (expression 0.0) * (expression 3)) / (expression 3)) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func ExampleVisitor_VisitAdditiveExpression() {
	text := "qubit q; U(pi/2, 0.0+3-3, pi) q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2)) , (expression (expression (expression 0.0) + (expression 3)) - (expression 3)) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func ExampleVisitor_VisitResetStatement() {
	text := "qubit q; reset q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (resetStatement reset (gateOperand (indexedIdentifier q)) ;))) <EOF>)
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func ExampleVisitor_VisitGateModifier_ctrl() {
	text := "qubit q0; qubit q1; U(pi/2.0, 0, pi) q0; ctrl @ U(pi, 0, pi) q0, q1;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q0 ;))) (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q1 ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression (expression pi) / (expression 2.0)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier ctrl @) U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q0)) , (gateOperand (indexedIdentifier q1))) ;))) <EOF>)
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleVisitor_VisitGateModifier_inv() {
	text := "qubit q; U(pi, tau, euler) q; inv @ U(pi, tau, euler) q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression tau) , (expression euler)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) (statementOrScope (statement (gateCallStatement (gateModifier inv @) U ( (expressionList (expression pi) , (expression tau) , (expression euler)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func ExampleVisitor_VisitGateModifier_pow() {
	text := "qubit q; pow(2) @ U(pi/2, 0, pi) q;"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit) q ;))) (statementOrScope (statement (gateCallStatement (gateModifier pow ( (expression 2) ) @) U ( (expressionList (expression (expression pi) / (expression 2)) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q))) ;))) <EOF>)
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func ExampleVisitor_VisitIndexedIdentifier() {
	text := "qubit[2] q; U(pi, 0, pi) q[0];"

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// (program (statementOrScope (statement (quantumDeclarationStatement (qubitType qubit (designator [ (expression 2) ])) q ;))) (statementOrScope (statement (gateCallStatement U ( (expressionList (expression pi) , (expression 0) , (expression pi)) ) (gateOperandList (gateOperand (indexedIdentifier q (indexOperator [ (expression 0) ])))) ;))) <EOF>)
	// [10][  2]( 1.0000 0.0000i): 1.0000
}
