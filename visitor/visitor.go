package visitor

import (
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
)

func New(qsim *q.Q) *Visitor {
	return &Visitor{
		qsim,
		NewEnviron(),
		&parser.Baseqasm3ParserVisitor{},
	}
}

type Visitor struct {
	qsim    *q.Q
	Environ *Environ
	*parser.Baseqasm3ParserVisitor
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
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitPragma(ctx *parser.PragmaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitStatementOrScope(ctx *parser.StatementOrScopeContext) interface{} {
	return v.VisitChildren(ctx)
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
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitMeasureArrowAssignmentStatement(ctx *parser.MeasureArrowAssignmentStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitResetStatement(ctx *parser.ResetStatementContext) interface{} {
	return v.VisitChildren(ctx)
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
	v.Environ.Qubit[ctx.Identifier().GetText()] = v.qsim.Zero()
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDefStatement(ctx *parser.DefStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitExternStatement(ctx *parser.ExternStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGateStatement(ctx *parser.GateStatementContext) interface{} {
	return v.VisitChildren(ctx)
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
	return v.VisitChildren(ctx)
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
	return v.VisitChildren(ctx)
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
	return v.VisitChildren(ctx)
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
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIndexedIdentifier(ctx *parser.IndexedIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitReturnSignature(ctx *parser.ReturnSignatureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGateModifier(ctx *parser.GateModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitScalarType(ctx *parser.ScalarTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitQubitType(ctx *parser.QubitTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArrayType(ctx *parser.ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArrayReferenceType(ctx *parser.ArrayReferenceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDesignator(ctx *parser.DesignatorContext) interface{} {
	return v.VisitChildren(ctx)
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
	return v.VisitChildren(ctx)
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
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIdentifierList(ctx *parser.IdentifierListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGateOperandList(ctx *parser.GateOperandListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitExternArgumentList(ctx *parser.ExternArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}
