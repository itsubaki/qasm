// Code generated from qasm3Parser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // qasm3Parser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by qasm3Parser.
type qasm3ParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by qasm3Parser#program.
	VisitProgram(ctx *ProgramContext) interface{}

	// Visit a parse tree produced by qasm3Parser#version.
	VisitVersion(ctx *VersionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#annotation.
	VisitAnnotation(ctx *AnnotationContext) interface{}

	// Visit a parse tree produced by qasm3Parser#scope.
	VisitScope(ctx *ScopeContext) interface{}

	// Visit a parse tree produced by qasm3Parser#pragma.
	VisitPragma(ctx *PragmaContext) interface{}

	// Visit a parse tree produced by qasm3Parser#statementOrScope.
	VisitStatementOrScope(ctx *StatementOrScopeContext) interface{}

	// Visit a parse tree produced by qasm3Parser#calibrationGrammarStatement.
	VisitCalibrationGrammarStatement(ctx *CalibrationGrammarStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#includeStatement.
	VisitIncludeStatement(ctx *IncludeStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#breakStatement.
	VisitBreakStatement(ctx *BreakStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#continueStatement.
	VisitContinueStatement(ctx *ContinueStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#endStatement.
	VisitEndStatement(ctx *EndStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#forStatement.
	VisitForStatement(ctx *ForStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#ifStatement.
	VisitIfStatement(ctx *IfStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#returnStatement.
	VisitReturnStatement(ctx *ReturnStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#whileStatement.
	VisitWhileStatement(ctx *WhileStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#switchStatement.
	VisitSwitchStatement(ctx *SwitchStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#switchCaseItem.
	VisitSwitchCaseItem(ctx *SwitchCaseItemContext) interface{}

	// Visit a parse tree produced by qasm3Parser#barrierStatement.
	VisitBarrierStatement(ctx *BarrierStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#boxStatement.
	VisitBoxStatement(ctx *BoxStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#delayStatement.
	VisitDelayStatement(ctx *DelayStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#gateCallStatement.
	VisitGateCallStatement(ctx *GateCallStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#measureArrowAssignmentStatement.
	VisitMeasureArrowAssignmentStatement(ctx *MeasureArrowAssignmentStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#resetStatement.
	VisitResetStatement(ctx *ResetStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#aliasDeclarationStatement.
	VisitAliasDeclarationStatement(ctx *AliasDeclarationStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#classicalDeclarationStatement.
	VisitClassicalDeclarationStatement(ctx *ClassicalDeclarationStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#constDeclarationStatement.
	VisitConstDeclarationStatement(ctx *ConstDeclarationStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#ioDeclarationStatement.
	VisitIoDeclarationStatement(ctx *IoDeclarationStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#oldStyleDeclarationStatement.
	VisitOldStyleDeclarationStatement(ctx *OldStyleDeclarationStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#quantumDeclarationStatement.
	VisitQuantumDeclarationStatement(ctx *QuantumDeclarationStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#defStatement.
	VisitDefStatement(ctx *DefStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#externStatement.
	VisitExternStatement(ctx *ExternStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#gateStatement.
	VisitGateStatement(ctx *GateStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#assignmentStatement.
	VisitAssignmentStatement(ctx *AssignmentStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#expressionStatement.
	VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#calStatement.
	VisitCalStatement(ctx *CalStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#defcalStatement.
	VisitDefcalStatement(ctx *DefcalStatementContext) interface{}

	// Visit a parse tree produced by qasm3Parser#bitwiseXorExpression.
	VisitBitwiseXorExpression(ctx *BitwiseXorExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#additiveExpression.
	VisitAdditiveExpression(ctx *AdditiveExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#durationofExpression.
	VisitDurationofExpression(ctx *DurationofExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#parenthesisExpression.
	VisitParenthesisExpression(ctx *ParenthesisExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#comparisonExpression.
	VisitComparisonExpression(ctx *ComparisonExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#multiplicativeExpression.
	VisitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#logicalOrExpression.
	VisitLogicalOrExpression(ctx *LogicalOrExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#castExpression.
	VisitCastExpression(ctx *CastExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#powerExpression.
	VisitPowerExpression(ctx *PowerExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#bitwiseOrExpression.
	VisitBitwiseOrExpression(ctx *BitwiseOrExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#callExpression.
	VisitCallExpression(ctx *CallExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#bitshiftExpression.
	VisitBitshiftExpression(ctx *BitshiftExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#bitwiseAndExpression.
	VisitBitwiseAndExpression(ctx *BitwiseAndExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#equalityExpression.
	VisitEqualityExpression(ctx *EqualityExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#logicalAndExpression.
	VisitLogicalAndExpression(ctx *LogicalAndExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#indexExpression.
	VisitIndexExpression(ctx *IndexExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#unaryExpression.
	VisitUnaryExpression(ctx *UnaryExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#literalExpression.
	VisitLiteralExpression(ctx *LiteralExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#aliasExpression.
	VisitAliasExpression(ctx *AliasExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#declarationExpression.
	VisitDeclarationExpression(ctx *DeclarationExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#measureExpression.
	VisitMeasureExpression(ctx *MeasureExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#rangeExpression.
	VisitRangeExpression(ctx *RangeExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#setExpression.
	VisitSetExpression(ctx *SetExpressionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#arrayLiteral.
	VisitArrayLiteral(ctx *ArrayLiteralContext) interface{}

	// Visit a parse tree produced by qasm3Parser#indexOperator.
	VisitIndexOperator(ctx *IndexOperatorContext) interface{}

	// Visit a parse tree produced by qasm3Parser#indexedIdentifier.
	VisitIndexedIdentifier(ctx *IndexedIdentifierContext) interface{}

	// Visit a parse tree produced by qasm3Parser#returnSignature.
	VisitReturnSignature(ctx *ReturnSignatureContext) interface{}

	// Visit a parse tree produced by qasm3Parser#gateModifier.
	VisitGateModifier(ctx *GateModifierContext) interface{}

	// Visit a parse tree produced by qasm3Parser#scalarType.
	VisitScalarType(ctx *ScalarTypeContext) interface{}

	// Visit a parse tree produced by qasm3Parser#qubitType.
	VisitQubitType(ctx *QubitTypeContext) interface{}

	// Visit a parse tree produced by qasm3Parser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by qasm3Parser#arrayReferenceType.
	VisitArrayReferenceType(ctx *ArrayReferenceTypeContext) interface{}

	// Visit a parse tree produced by qasm3Parser#designator.
	VisitDesignator(ctx *DesignatorContext) interface{}

	// Visit a parse tree produced by qasm3Parser#defcalTarget.
	VisitDefcalTarget(ctx *DefcalTargetContext) interface{}

	// Visit a parse tree produced by qasm3Parser#defcalArgumentDefinition.
	VisitDefcalArgumentDefinition(ctx *DefcalArgumentDefinitionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#defcalOperand.
	VisitDefcalOperand(ctx *DefcalOperandContext) interface{}

	// Visit a parse tree produced by qasm3Parser#gateOperand.
	VisitGateOperand(ctx *GateOperandContext) interface{}

	// Visit a parse tree produced by qasm3Parser#externArgument.
	VisitExternArgument(ctx *ExternArgumentContext) interface{}

	// Visit a parse tree produced by qasm3Parser#argumentDefinition.
	VisitArgumentDefinition(ctx *ArgumentDefinitionContext) interface{}

	// Visit a parse tree produced by qasm3Parser#argumentDefinitionList.
	VisitArgumentDefinitionList(ctx *ArgumentDefinitionListContext) interface{}

	// Visit a parse tree produced by qasm3Parser#defcalArgumentDefinitionList.
	VisitDefcalArgumentDefinitionList(ctx *DefcalArgumentDefinitionListContext) interface{}

	// Visit a parse tree produced by qasm3Parser#defcalOperandList.
	VisitDefcalOperandList(ctx *DefcalOperandListContext) interface{}

	// Visit a parse tree produced by qasm3Parser#expressionList.
	VisitExpressionList(ctx *ExpressionListContext) interface{}

	// Visit a parse tree produced by qasm3Parser#identifierList.
	VisitIdentifierList(ctx *IdentifierListContext) interface{}

	// Visit a parse tree produced by qasm3Parser#gateOperandList.
	VisitGateOperandList(ctx *GateOperandListContext) interface{}

	// Visit a parse tree produced by qasm3Parser#externArgumentList.
	VisitExternArgumentList(ctx *ExternArgumentListContext) interface{}
}
