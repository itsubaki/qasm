package visitor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/qasm/gen/parser"
)

var (
	ErrAlreadyDeclared = errors.New("already declared")
	ErrQubitNotFound   = errors.New("qubit not found")
	ErrGateNotFound    = errors.New("gate not found")
	ErrUnexpectedType  = errors.New("unexpected type")
)

func New(qsim *q.Q, env *Environ) *Visitor {
	return &Visitor{
		qsim,
		env,
	}
}

type Visitor struct {
	qsim    *q.Q
	Environ *Environ
}

func (v *Visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *Visitor) VisitTerminal(_ antlr.TerminalNode) interface{} {
	return nil
}

func (v *Visitor) VisitErrorNode(_ antlr.ErrorNode) interface{} {
	return nil
}

func (v *Visitor) VisitChildren(node antlr.RuleNode) interface{} {
	var result interface{}
	for _, n := range node.GetChildren() {
		tree, ok := n.(antlr.ParseTree)
		if !ok {
			continue
		}

		if res := v.Visit(tree); res != nil {
			result = res
			break
		}
	}

	return result
}

func (v *Visitor) VisitProgram(ctx *parser.ProgramContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitVersion(ctx *parser.VersionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitStatement(ctx *parser.StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitAnnotation(ctx *parser.AnnotationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitScope(ctx *parser.ScopeContext) interface{} {
	var list []interface{}
	for _, s := range ctx.AllStatementOrScope() {
		list = append(list, v.Visit(s))
	}

	return list
}

func (v *Visitor) VisitPragma(ctx *parser.PragmaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitStatementOrScope(ctx *parser.StatementOrScopeContext) interface{} {
	if ctx.Statement() != nil {
		return v.Visit(ctx.Statement())
	}

	return v.Visit(ctx.Scope())
}

func (v *Visitor) VisitCalibrationGrammarStatement(ctx *parser.CalibrationGrammarStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIncludeStatement(ctx *parser.IncludeStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBreakStatement(ctx *parser.BreakStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitContinueStatement(ctx *parser.ContinueStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitEndStatement(ctx *parser.EndStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitForStatement(ctx *parser.ForStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIfStatement(ctx *parser.IfStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitReturnStatement(ctx *parser.ReturnStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitWhileStatement(ctx *parser.WhileStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSwitchStatement(ctx *parser.SwitchStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSwitchCaseItem(ctx *parser.SwitchCaseItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBarrierStatement(ctx *parser.BarrierStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBoxStatement(ctx *parser.BoxStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDelayStatement(ctx *parser.DelayStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGateCallStatement(ctx *parser.GateCallStatementContext) interface{} {
	var params []float64
	for _, e := range v.Visit(ctx.ExpressionList()).([]interface{}) {
		switch val := e.(type) {
		case float64:
			params = append(params, val)
		case int64:
			params = append(params, float64(val))
		default:
			return fmt.Errorf("param=%v: %w", e, ErrUnexpectedType)
		}
	}

	var qb []q.Qubit
	for _, o := range v.Visit(ctx.GateOperandList()).([]interface{}) {
		if q, ok := v.Environ.Qubit[o.(string)]; ok {
			qb = append(qb, q...)
			continue
		}

		return fmt.Errorf("operand=%s: %w", o, ErrQubitNotFound)
	}

	id := ctx.Identifier().GetText()
	switch id {
	case "U":
		u := gate.U(params[0], params[1], params[2])

		var ctrl bool
		for _, mod := range ctx.AllGateModifier() {
			m := v.Visit(mod).(string)
			switch m {
			case "ctrl@":
				n := v.qsim.NumberOfBit()
				index := q.Index(qb...)
				u = gate.C(u, n, index[0], index[1])
				ctrl = true
			case "negctrl@":
				// TODO: implement
				ctrl = true
			case "inv@":
				u = u.Inverse()
			default:
				if !strings.HasPrefix(m, "pow") {
					return fmt.Errorf("modifier=%s: %w", m, ErrUnexpectedType)
				}

				x := mod.Expression()
				p, ok := v.Visit(x).(int64)
				if !ok {
					return fmt.Errorf("pow=%v: %w", x.GetText(), ErrUnexpectedType)
				}

				u = matrix.ApplyN(u, int(p))
			}
		}

		if ctrl {
			v.qsim.Apply(u)
			return nil
		}

		v.qsim.Apply(u, qb...)
		return nil
	default:
		return fmt.Errorf("idenfitier=%s: %w", id, ErrGateNotFound)
	}
}

func (v *Visitor) VisitMeasureArrowAssignmentStatement(ctx *parser.MeasureArrowAssignmentStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitResetStatement(ctx *parser.ResetStatementContext) interface{} {
	operand := v.Visit(ctx.GateOperand()).(string)
	q, ok := v.Environ.Qubit[operand]
	if !ok {
		return fmt.Errorf("operand=%s: %w", operand, ErrQubitNotFound)
	}

	v.qsim.Reset(q...)
	return nil
}

func (v *Visitor) VisitAliasDeclarationStatement(ctx *parser.AliasDeclarationStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitClassicalDeclarationStatement(ctx *parser.ClassicalDeclarationStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitConstDeclarationStatement(ctx *parser.ConstDeclarationStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIoDeclarationStatement(ctx *parser.IoDeclarationStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitOldStyleDeclarationStatement(ctx *parser.OldStyleDeclarationStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitQuantumDeclarationStatement(ctx *parser.QuantumDeclarationStatementContext) interface{} {
	id := ctx.Identifier().GetText()
	if _, ok := v.Environ.Qubit[id]; ok {
		return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
	}

	designator := ctx.QubitType().Designator()
	if designator == nil {
		v.Environ.Qubit[id] = []q.Qubit{v.qsim.Zero()}
		return nil
	}

	size := v.Visit(designator).(int64)
	v.Environ.Qubit[id] = v.qsim.ZeroWith(int(size))
	return nil
}

func (v *Visitor) VisitDefStatement(ctx *parser.DefStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitExternStatement(ctx *parser.ExternStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGateStatement(ctx *parser.GateStatementContext) interface{} {
	fmt.Println(ctx.GetText())
	return nil
}

func (v *Visitor) VisitAssignmentStatement(ctx *parser.AssignmentStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitExpressionStatement(ctx *parser.ExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitCalStatement(ctx *parser.CalStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDefcalStatement(ctx *parser.DefcalStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBitwiseXorExpression(ctx *parser.BitwiseXorExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitAdditiveExpression(ctx *parser.AdditiveExpressionContext) interface{} {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		default:
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpectedType)
		}
	}

	op := ctx.GetOp().GetText()
	switch op {
	case "+":
		return operand[0] + operand[1]
	case "-":
		return operand[0] - operand[1]
	default:
		return fmt.Errorf("operator=%s: %w", op, ErrUnexpectedType)
	}
}

func (v *Visitor) VisitDurationofExpression(ctx *parser.DurationofExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitParenthesisExpression(ctx *parser.ParenthesisExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitComparisonExpression(ctx *parser.ComparisonExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) interface{} {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		default:
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpectedType)
		}
	}

	op := ctx.GetOp().GetText()
	switch op {
	case "*":
		return operand[0] * operand[1]
	case "/":
		return operand[0] / operand[1]
	default:
		return fmt.Errorf("operator=%s: %w", op, ErrUnexpectedType)
	}
}

func (v *Visitor) VisitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitCastExpression(ctx *parser.CastExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitPowerExpression(ctx *parser.PowerExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBitwiseOrExpression(ctx *parser.BitwiseOrExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitCallExpression(ctx *parser.CallExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBitshiftExpression(ctx *parser.BitshiftExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBitwiseAndExpression(ctx *parser.BitwiseAndExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitEqualityExpression(ctx *parser.EqualityExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIndexExpression(ctx *parser.IndexExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitUnaryExpression(ctx *parser.UnaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) interface{} {
	x := ctx.GetText()
	if lit, ok := BultinConst[x]; ok {
		return lit
	}

	if ctx.FloatLiteral() != nil {
		lit, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return fmt.Errorf("parse float: x=%s: %w", x, err)
		}

		return lit
	}

	if ctx.DecimalIntegerLiteral() != nil {
		lit, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int: x=%s: %w", x, err)
		}

		return lit
	}

	return fmt.Errorf("x=%s: %w", x, ErrUnexpectedType)
}

func (v *Visitor) VisitAliasExpression(ctx *parser.AliasExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDeclarationExpression(ctx *parser.DeclarationExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitMeasureExpression(ctx *parser.MeasureExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitRangeExpression(ctx *parser.RangeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSetExpression(ctx *parser.SetExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArrayLiteral(ctx *parser.ArrayLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIndexOperator(ctx *parser.IndexOperatorContext) interface{} {
	var list []interface{}
	for _, x := range ctx.AllExpression() {
		list = append(list, v.Visit(x))
	}

	return list
}

func (v *Visitor) VisitIndexedIdentifier(ctx *parser.IndexedIdentifierContext) interface{} {
	for _, op := range ctx.AllIndexOperator() {
		fmt.Println(v.Visit(op))
	}

	return ctx.Identifier().GetText()
}

func (v *Visitor) VisitReturnSignature(ctx *parser.ReturnSignatureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGateModifier(ctx *parser.GateModifierContext) interface{} {
	return ctx.GetText()
}

func (v *Visitor) VisitScalarType(ctx *parser.ScalarTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitQubitType(ctx *parser.QubitTypeContext) interface{} {
	return ctx.QUBIT()
}

func (v *Visitor) VisitArrayType(ctx *parser.ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArrayReferenceType(ctx *parser.ArrayReferenceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDesignator(ctx *parser.DesignatorContext) interface{} {
	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitDefcalTarget(ctx *parser.DefcalTargetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDefcalArgumentDefinition(ctx *parser.DefcalArgumentDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDefcalOperand(ctx *parser.DefcalOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGateOperand(ctx *parser.GateOperandContext) interface{} {
	return v.Visit(ctx.IndexedIdentifier())
}

func (v *Visitor) VisitExternArgument(ctx *parser.ExternArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArgumentDefinition(ctx *parser.ArgumentDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArgumentDefinitionList(ctx *parser.ArgumentDefinitionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDefcalArgumentDefinitionList(ctx *parser.DefcalArgumentDefinitionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDefcalOperandList(ctx *parser.DefcalOperandListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitExpressionList(ctx *parser.ExpressionListContext) interface{} {
	var list []interface{}
	for _, x := range ctx.AllExpression() {
		list = append(list, v.Visit(x))
	}

	return list
}

func (v *Visitor) VisitIdentifierList(ctx *parser.IdentifierListContext) interface{} {
	var list []interface{}
	for _, id := range ctx.AllIdentifier() {
		list = append(list, id.GetText())
	}

	return list
}

func (v *Visitor) VisitGateOperandList(ctx *parser.GateOperandListContext) interface{} {
	var list []interface{}
	for _, o := range ctx.AllGateOperand() {
		list = append(list, v.Visit(o))
	}

	return list
}

func (v *Visitor) VisitExternArgumentList(ctx *parser.ExternArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}
