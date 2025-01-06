// Code generated from qasm3Parser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // qasm3Parser
import "github.com/antlr4-go/antlr/v4"

// Baseqasm3ParserListener is a complete listener for a parse tree produced by qasm3Parser.
type Baseqasm3ParserListener struct{}

var _ qasm3ParserListener = &Baseqasm3ParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *Baseqasm3ParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *Baseqasm3ParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *Baseqasm3ParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *Baseqasm3ParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *Baseqasm3ParserListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *Baseqasm3ParserListener) ExitProgram(ctx *ProgramContext) {}

// EnterVersion is called when production version is entered.
func (s *Baseqasm3ParserListener) EnterVersion(ctx *VersionContext) {}

// ExitVersion is called when production version is exited.
func (s *Baseqasm3ParserListener) ExitVersion(ctx *VersionContext) {}

// EnterStatement is called when production statement is entered.
func (s *Baseqasm3ParserListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *Baseqasm3ParserListener) ExitStatement(ctx *StatementContext) {}

// EnterAnnotation is called when production annotation is entered.
func (s *Baseqasm3ParserListener) EnterAnnotation(ctx *AnnotationContext) {}

// ExitAnnotation is called when production annotation is exited.
func (s *Baseqasm3ParserListener) ExitAnnotation(ctx *AnnotationContext) {}

// EnterScope is called when production scope is entered.
func (s *Baseqasm3ParserListener) EnterScope(ctx *ScopeContext) {}

// ExitScope is called when production scope is exited.
func (s *Baseqasm3ParserListener) ExitScope(ctx *ScopeContext) {}

// EnterPragma is called when production pragma is entered.
func (s *Baseqasm3ParserListener) EnterPragma(ctx *PragmaContext) {}

// ExitPragma is called when production pragma is exited.
func (s *Baseqasm3ParserListener) ExitPragma(ctx *PragmaContext) {}

// EnterStatementOrScope is called when production statementOrScope is entered.
func (s *Baseqasm3ParserListener) EnterStatementOrScope(ctx *StatementOrScopeContext) {}

// ExitStatementOrScope is called when production statementOrScope is exited.
func (s *Baseqasm3ParserListener) ExitStatementOrScope(ctx *StatementOrScopeContext) {}

// EnterCalibrationGrammarStatement is called when production calibrationGrammarStatement is entered.
func (s *Baseqasm3ParserListener) EnterCalibrationGrammarStatement(ctx *CalibrationGrammarStatementContext) {
}

// ExitCalibrationGrammarStatement is called when production calibrationGrammarStatement is exited.
func (s *Baseqasm3ParserListener) ExitCalibrationGrammarStatement(ctx *CalibrationGrammarStatementContext) {
}

// EnterIncludeStatement is called when production includeStatement is entered.
func (s *Baseqasm3ParserListener) EnterIncludeStatement(ctx *IncludeStatementContext) {}

// ExitIncludeStatement is called when production includeStatement is exited.
func (s *Baseqasm3ParserListener) ExitIncludeStatement(ctx *IncludeStatementContext) {}

// EnterBreakStatement is called when production breakStatement is entered.
func (s *Baseqasm3ParserListener) EnterBreakStatement(ctx *BreakStatementContext) {}

// ExitBreakStatement is called when production breakStatement is exited.
func (s *Baseqasm3ParserListener) ExitBreakStatement(ctx *BreakStatementContext) {}

// EnterContinueStatement is called when production continueStatement is entered.
func (s *Baseqasm3ParserListener) EnterContinueStatement(ctx *ContinueStatementContext) {}

// ExitContinueStatement is called when production continueStatement is exited.
func (s *Baseqasm3ParserListener) ExitContinueStatement(ctx *ContinueStatementContext) {}

// EnterEndStatement is called when production endStatement is entered.
func (s *Baseqasm3ParserListener) EnterEndStatement(ctx *EndStatementContext) {}

// ExitEndStatement is called when production endStatement is exited.
func (s *Baseqasm3ParserListener) ExitEndStatement(ctx *EndStatementContext) {}

// EnterForStatement is called when production forStatement is entered.
func (s *Baseqasm3ParserListener) EnterForStatement(ctx *ForStatementContext) {}

// ExitForStatement is called when production forStatement is exited.
func (s *Baseqasm3ParserListener) ExitForStatement(ctx *ForStatementContext) {}

// EnterIfStatement is called when production ifStatement is entered.
func (s *Baseqasm3ParserListener) EnterIfStatement(ctx *IfStatementContext) {}

// ExitIfStatement is called when production ifStatement is exited.
func (s *Baseqasm3ParserListener) ExitIfStatement(ctx *IfStatementContext) {}

// EnterReturnStatement is called when production returnStatement is entered.
func (s *Baseqasm3ParserListener) EnterReturnStatement(ctx *ReturnStatementContext) {}

// ExitReturnStatement is called when production returnStatement is exited.
func (s *Baseqasm3ParserListener) ExitReturnStatement(ctx *ReturnStatementContext) {}

// EnterWhileStatement is called when production whileStatement is entered.
func (s *Baseqasm3ParserListener) EnterWhileStatement(ctx *WhileStatementContext) {}

// ExitWhileStatement is called when production whileStatement is exited.
func (s *Baseqasm3ParserListener) ExitWhileStatement(ctx *WhileStatementContext) {}

// EnterSwitchStatement is called when production switchStatement is entered.
func (s *Baseqasm3ParserListener) EnterSwitchStatement(ctx *SwitchStatementContext) {}

// ExitSwitchStatement is called when production switchStatement is exited.
func (s *Baseqasm3ParserListener) ExitSwitchStatement(ctx *SwitchStatementContext) {}

// EnterSwitchCaseItem is called when production switchCaseItem is entered.
func (s *Baseqasm3ParserListener) EnterSwitchCaseItem(ctx *SwitchCaseItemContext) {}

// ExitSwitchCaseItem is called when production switchCaseItem is exited.
func (s *Baseqasm3ParserListener) ExitSwitchCaseItem(ctx *SwitchCaseItemContext) {}

// EnterBarrierStatement is called when production barrierStatement is entered.
func (s *Baseqasm3ParserListener) EnterBarrierStatement(ctx *BarrierStatementContext) {}

// ExitBarrierStatement is called when production barrierStatement is exited.
func (s *Baseqasm3ParserListener) ExitBarrierStatement(ctx *BarrierStatementContext) {}

// EnterBoxStatement is called when production boxStatement is entered.
func (s *Baseqasm3ParserListener) EnterBoxStatement(ctx *BoxStatementContext) {}

// ExitBoxStatement is called when production boxStatement is exited.
func (s *Baseqasm3ParserListener) ExitBoxStatement(ctx *BoxStatementContext) {}

// EnterDelayStatement is called when production delayStatement is entered.
func (s *Baseqasm3ParserListener) EnterDelayStatement(ctx *DelayStatementContext) {}

// ExitDelayStatement is called when production delayStatement is exited.
func (s *Baseqasm3ParserListener) ExitDelayStatement(ctx *DelayStatementContext) {}

// EnterGateCallStatement is called when production gateCallStatement is entered.
func (s *Baseqasm3ParserListener) EnterGateCallStatement(ctx *GateCallStatementContext) {}

// ExitGateCallStatement is called when production gateCallStatement is exited.
func (s *Baseqasm3ParserListener) ExitGateCallStatement(ctx *GateCallStatementContext) {}

// EnterMeasureArrowAssignmentStatement is called when production measureArrowAssignmentStatement is entered.
func (s *Baseqasm3ParserListener) EnterMeasureArrowAssignmentStatement(ctx *MeasureArrowAssignmentStatementContext) {
}

// ExitMeasureArrowAssignmentStatement is called when production measureArrowAssignmentStatement is exited.
func (s *Baseqasm3ParserListener) ExitMeasureArrowAssignmentStatement(ctx *MeasureArrowAssignmentStatementContext) {
}

// EnterResetStatement is called when production resetStatement is entered.
func (s *Baseqasm3ParserListener) EnterResetStatement(ctx *ResetStatementContext) {}

// ExitResetStatement is called when production resetStatement is exited.
func (s *Baseqasm3ParserListener) ExitResetStatement(ctx *ResetStatementContext) {}

// EnterAliasDeclarationStatement is called when production aliasDeclarationStatement is entered.
func (s *Baseqasm3ParserListener) EnterAliasDeclarationStatement(ctx *AliasDeclarationStatementContext) {
}

// ExitAliasDeclarationStatement is called when production aliasDeclarationStatement is exited.
func (s *Baseqasm3ParserListener) ExitAliasDeclarationStatement(ctx *AliasDeclarationStatementContext) {
}

// EnterClassicalDeclarationStatement is called when production classicalDeclarationStatement is entered.
func (s *Baseqasm3ParserListener) EnterClassicalDeclarationStatement(ctx *ClassicalDeclarationStatementContext) {
}

// ExitClassicalDeclarationStatement is called when production classicalDeclarationStatement is exited.
func (s *Baseqasm3ParserListener) ExitClassicalDeclarationStatement(ctx *ClassicalDeclarationStatementContext) {
}

// EnterConstDeclarationStatement is called when production constDeclarationStatement is entered.
func (s *Baseqasm3ParserListener) EnterConstDeclarationStatement(ctx *ConstDeclarationStatementContext) {
}

// ExitConstDeclarationStatement is called when production constDeclarationStatement is exited.
func (s *Baseqasm3ParserListener) ExitConstDeclarationStatement(ctx *ConstDeclarationStatementContext) {
}

// EnterIoDeclarationStatement is called when production ioDeclarationStatement is entered.
func (s *Baseqasm3ParserListener) EnterIoDeclarationStatement(ctx *IoDeclarationStatementContext) {}

// ExitIoDeclarationStatement is called when production ioDeclarationStatement is exited.
func (s *Baseqasm3ParserListener) ExitIoDeclarationStatement(ctx *IoDeclarationStatementContext) {}

// EnterOldStyleDeclarationStatement is called when production oldStyleDeclarationStatement is entered.
func (s *Baseqasm3ParserListener) EnterOldStyleDeclarationStatement(ctx *OldStyleDeclarationStatementContext) {
}

// ExitOldStyleDeclarationStatement is called when production oldStyleDeclarationStatement is exited.
func (s *Baseqasm3ParserListener) ExitOldStyleDeclarationStatement(ctx *OldStyleDeclarationStatementContext) {
}

// EnterQuantumDeclarationStatement is called when production quantumDeclarationStatement is entered.
func (s *Baseqasm3ParserListener) EnterQuantumDeclarationStatement(ctx *QuantumDeclarationStatementContext) {
}

// ExitQuantumDeclarationStatement is called when production quantumDeclarationStatement is exited.
func (s *Baseqasm3ParserListener) ExitQuantumDeclarationStatement(ctx *QuantumDeclarationStatementContext) {
}

// EnterDefStatement is called when production defStatement is entered.
func (s *Baseqasm3ParserListener) EnterDefStatement(ctx *DefStatementContext) {}

// ExitDefStatement is called when production defStatement is exited.
func (s *Baseqasm3ParserListener) ExitDefStatement(ctx *DefStatementContext) {}

// EnterExternStatement is called when production externStatement is entered.
func (s *Baseqasm3ParserListener) EnterExternStatement(ctx *ExternStatementContext) {}

// ExitExternStatement is called when production externStatement is exited.
func (s *Baseqasm3ParserListener) ExitExternStatement(ctx *ExternStatementContext) {}

// EnterGateStatement is called when production gateStatement is entered.
func (s *Baseqasm3ParserListener) EnterGateStatement(ctx *GateStatementContext) {}

// ExitGateStatement is called when production gateStatement is exited.
func (s *Baseqasm3ParserListener) ExitGateStatement(ctx *GateStatementContext) {}

// EnterAssignmentStatement is called when production assignmentStatement is entered.
func (s *Baseqasm3ParserListener) EnterAssignmentStatement(ctx *AssignmentStatementContext) {}

// ExitAssignmentStatement is called when production assignmentStatement is exited.
func (s *Baseqasm3ParserListener) ExitAssignmentStatement(ctx *AssignmentStatementContext) {}

// EnterExpressionStatement is called when production expressionStatement is entered.
func (s *Baseqasm3ParserListener) EnterExpressionStatement(ctx *ExpressionStatementContext) {}

// ExitExpressionStatement is called when production expressionStatement is exited.
func (s *Baseqasm3ParserListener) ExitExpressionStatement(ctx *ExpressionStatementContext) {}

// EnterCalStatement is called when production calStatement is entered.
func (s *Baseqasm3ParserListener) EnterCalStatement(ctx *CalStatementContext) {}

// ExitCalStatement is called when production calStatement is exited.
func (s *Baseqasm3ParserListener) ExitCalStatement(ctx *CalStatementContext) {}

// EnterDefcalStatement is called when production defcalStatement is entered.
func (s *Baseqasm3ParserListener) EnterDefcalStatement(ctx *DefcalStatementContext) {}

// ExitDefcalStatement is called when production defcalStatement is exited.
func (s *Baseqasm3ParserListener) ExitDefcalStatement(ctx *DefcalStatementContext) {}

// EnterBitwiseXorExpression is called when production bitwiseXorExpression is entered.
func (s *Baseqasm3ParserListener) EnterBitwiseXorExpression(ctx *BitwiseXorExpressionContext) {}

// ExitBitwiseXorExpression is called when production bitwiseXorExpression is exited.
func (s *Baseqasm3ParserListener) ExitBitwiseXorExpression(ctx *BitwiseXorExpressionContext) {}

// EnterAdditiveExpression is called when production additiveExpression is entered.
func (s *Baseqasm3ParserListener) EnterAdditiveExpression(ctx *AdditiveExpressionContext) {}

// ExitAdditiveExpression is called when production additiveExpression is exited.
func (s *Baseqasm3ParserListener) ExitAdditiveExpression(ctx *AdditiveExpressionContext) {}

// EnterDurationofExpression is called when production durationofExpression is entered.
func (s *Baseqasm3ParserListener) EnterDurationofExpression(ctx *DurationofExpressionContext) {}

// ExitDurationofExpression is called when production durationofExpression is exited.
func (s *Baseqasm3ParserListener) ExitDurationofExpression(ctx *DurationofExpressionContext) {}

// EnterParenthesisExpression is called when production parenthesisExpression is entered.
func (s *Baseqasm3ParserListener) EnterParenthesisExpression(ctx *ParenthesisExpressionContext) {}

// ExitParenthesisExpression is called when production parenthesisExpression is exited.
func (s *Baseqasm3ParserListener) ExitParenthesisExpression(ctx *ParenthesisExpressionContext) {}

// EnterComparisonExpression is called when production comparisonExpression is entered.
func (s *Baseqasm3ParserListener) EnterComparisonExpression(ctx *ComparisonExpressionContext) {}

// ExitComparisonExpression is called when production comparisonExpression is exited.
func (s *Baseqasm3ParserListener) ExitComparisonExpression(ctx *ComparisonExpressionContext) {}

// EnterMultiplicativeExpression is called when production multiplicativeExpression is entered.
func (s *Baseqasm3ParserListener) EnterMultiplicativeExpression(ctx *MultiplicativeExpressionContext) {
}

// ExitMultiplicativeExpression is called when production multiplicativeExpression is exited.
func (s *Baseqasm3ParserListener) ExitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) {
}

// EnterLogicalOrExpression is called when production logicalOrExpression is entered.
func (s *Baseqasm3ParserListener) EnterLogicalOrExpression(ctx *LogicalOrExpressionContext) {}

// ExitLogicalOrExpression is called when production logicalOrExpression is exited.
func (s *Baseqasm3ParserListener) ExitLogicalOrExpression(ctx *LogicalOrExpressionContext) {}

// EnterCastExpression is called when production castExpression is entered.
func (s *Baseqasm3ParserListener) EnterCastExpression(ctx *CastExpressionContext) {}

// ExitCastExpression is called when production castExpression is exited.
func (s *Baseqasm3ParserListener) ExitCastExpression(ctx *CastExpressionContext) {}

// EnterPowerExpression is called when production powerExpression is entered.
func (s *Baseqasm3ParserListener) EnterPowerExpression(ctx *PowerExpressionContext) {}

// ExitPowerExpression is called when production powerExpression is exited.
func (s *Baseqasm3ParserListener) ExitPowerExpression(ctx *PowerExpressionContext) {}

// EnterBitwiseOrExpression is called when production bitwiseOrExpression is entered.
func (s *Baseqasm3ParserListener) EnterBitwiseOrExpression(ctx *BitwiseOrExpressionContext) {}

// ExitBitwiseOrExpression is called when production bitwiseOrExpression is exited.
func (s *Baseqasm3ParserListener) ExitBitwiseOrExpression(ctx *BitwiseOrExpressionContext) {}

// EnterCallExpression is called when production callExpression is entered.
func (s *Baseqasm3ParserListener) EnterCallExpression(ctx *CallExpressionContext) {}

// ExitCallExpression is called when production callExpression is exited.
func (s *Baseqasm3ParserListener) ExitCallExpression(ctx *CallExpressionContext) {}

// EnterBitshiftExpression is called when production bitshiftExpression is entered.
func (s *Baseqasm3ParserListener) EnterBitshiftExpression(ctx *BitshiftExpressionContext) {}

// ExitBitshiftExpression is called when production bitshiftExpression is exited.
func (s *Baseqasm3ParserListener) ExitBitshiftExpression(ctx *BitshiftExpressionContext) {}

// EnterBitwiseAndExpression is called when production bitwiseAndExpression is entered.
func (s *Baseqasm3ParserListener) EnterBitwiseAndExpression(ctx *BitwiseAndExpressionContext) {}

// ExitBitwiseAndExpression is called when production bitwiseAndExpression is exited.
func (s *Baseqasm3ParserListener) ExitBitwiseAndExpression(ctx *BitwiseAndExpressionContext) {}

// EnterEqualityExpression is called when production equalityExpression is entered.
func (s *Baseqasm3ParserListener) EnterEqualityExpression(ctx *EqualityExpressionContext) {}

// ExitEqualityExpression is called when production equalityExpression is exited.
func (s *Baseqasm3ParserListener) ExitEqualityExpression(ctx *EqualityExpressionContext) {}

// EnterLogicalAndExpression is called when production logicalAndExpression is entered.
func (s *Baseqasm3ParserListener) EnterLogicalAndExpression(ctx *LogicalAndExpressionContext) {}

// ExitLogicalAndExpression is called when production logicalAndExpression is exited.
func (s *Baseqasm3ParserListener) ExitLogicalAndExpression(ctx *LogicalAndExpressionContext) {}

// EnterIndexExpression is called when production indexExpression is entered.
func (s *Baseqasm3ParserListener) EnterIndexExpression(ctx *IndexExpressionContext) {}

// ExitIndexExpression is called when production indexExpression is exited.
func (s *Baseqasm3ParserListener) ExitIndexExpression(ctx *IndexExpressionContext) {}

// EnterUnaryExpression is called when production unaryExpression is entered.
func (s *Baseqasm3ParserListener) EnterUnaryExpression(ctx *UnaryExpressionContext) {}

// ExitUnaryExpression is called when production unaryExpression is exited.
func (s *Baseqasm3ParserListener) ExitUnaryExpression(ctx *UnaryExpressionContext) {}

// EnterLiteralExpression is called when production literalExpression is entered.
func (s *Baseqasm3ParserListener) EnterLiteralExpression(ctx *LiteralExpressionContext) {}

// ExitLiteralExpression is called when production literalExpression is exited.
func (s *Baseqasm3ParserListener) ExitLiteralExpression(ctx *LiteralExpressionContext) {}

// EnterAliasExpression is called when production aliasExpression is entered.
func (s *Baseqasm3ParserListener) EnterAliasExpression(ctx *AliasExpressionContext) {}

// ExitAliasExpression is called when production aliasExpression is exited.
func (s *Baseqasm3ParserListener) ExitAliasExpression(ctx *AliasExpressionContext) {}

// EnterDeclarationExpression is called when production declarationExpression is entered.
func (s *Baseqasm3ParserListener) EnterDeclarationExpression(ctx *DeclarationExpressionContext) {}

// ExitDeclarationExpression is called when production declarationExpression is exited.
func (s *Baseqasm3ParserListener) ExitDeclarationExpression(ctx *DeclarationExpressionContext) {}

// EnterMeasureExpression is called when production measureExpression is entered.
func (s *Baseqasm3ParserListener) EnterMeasureExpression(ctx *MeasureExpressionContext) {}

// ExitMeasureExpression is called when production measureExpression is exited.
func (s *Baseqasm3ParserListener) ExitMeasureExpression(ctx *MeasureExpressionContext) {}

// EnterRangeExpression is called when production rangeExpression is entered.
func (s *Baseqasm3ParserListener) EnterRangeExpression(ctx *RangeExpressionContext) {}

// ExitRangeExpression is called when production rangeExpression is exited.
func (s *Baseqasm3ParserListener) ExitRangeExpression(ctx *RangeExpressionContext) {}

// EnterSetExpression is called when production setExpression is entered.
func (s *Baseqasm3ParserListener) EnterSetExpression(ctx *SetExpressionContext) {}

// ExitSetExpression is called when production setExpression is exited.
func (s *Baseqasm3ParserListener) ExitSetExpression(ctx *SetExpressionContext) {}

// EnterArrayLiteral is called when production arrayLiteral is entered.
func (s *Baseqasm3ParserListener) EnterArrayLiteral(ctx *ArrayLiteralContext) {}

// ExitArrayLiteral is called when production arrayLiteral is exited.
func (s *Baseqasm3ParserListener) ExitArrayLiteral(ctx *ArrayLiteralContext) {}

// EnterIndexOperator is called when production indexOperator is entered.
func (s *Baseqasm3ParserListener) EnterIndexOperator(ctx *IndexOperatorContext) {}

// ExitIndexOperator is called when production indexOperator is exited.
func (s *Baseqasm3ParserListener) ExitIndexOperator(ctx *IndexOperatorContext) {}

// EnterIndexedIdentifier is called when production indexedIdentifier is entered.
func (s *Baseqasm3ParserListener) EnterIndexedIdentifier(ctx *IndexedIdentifierContext) {}

// ExitIndexedIdentifier is called when production indexedIdentifier is exited.
func (s *Baseqasm3ParserListener) ExitIndexedIdentifier(ctx *IndexedIdentifierContext) {}

// EnterReturnSignature is called when production returnSignature is entered.
func (s *Baseqasm3ParserListener) EnterReturnSignature(ctx *ReturnSignatureContext) {}

// ExitReturnSignature is called when production returnSignature is exited.
func (s *Baseqasm3ParserListener) ExitReturnSignature(ctx *ReturnSignatureContext) {}

// EnterGateModifier is called when production gateModifier is entered.
func (s *Baseqasm3ParserListener) EnterGateModifier(ctx *GateModifierContext) {}

// ExitGateModifier is called when production gateModifier is exited.
func (s *Baseqasm3ParserListener) ExitGateModifier(ctx *GateModifierContext) {}

// EnterScalarType is called when production scalarType is entered.
func (s *Baseqasm3ParserListener) EnterScalarType(ctx *ScalarTypeContext) {}

// ExitScalarType is called when production scalarType is exited.
func (s *Baseqasm3ParserListener) ExitScalarType(ctx *ScalarTypeContext) {}

// EnterQubitType is called when production qubitType is entered.
func (s *Baseqasm3ParserListener) EnterQubitType(ctx *QubitTypeContext) {}

// ExitQubitType is called when production qubitType is exited.
func (s *Baseqasm3ParserListener) ExitQubitType(ctx *QubitTypeContext) {}

// EnterArrayType is called when production arrayType is entered.
func (s *Baseqasm3ParserListener) EnterArrayType(ctx *ArrayTypeContext) {}

// ExitArrayType is called when production arrayType is exited.
func (s *Baseqasm3ParserListener) ExitArrayType(ctx *ArrayTypeContext) {}

// EnterArrayReferenceType is called when production arrayReferenceType is entered.
func (s *Baseqasm3ParserListener) EnterArrayReferenceType(ctx *ArrayReferenceTypeContext) {}

// ExitArrayReferenceType is called when production arrayReferenceType is exited.
func (s *Baseqasm3ParserListener) ExitArrayReferenceType(ctx *ArrayReferenceTypeContext) {}

// EnterDesignator is called when production designator is entered.
func (s *Baseqasm3ParserListener) EnterDesignator(ctx *DesignatorContext) {}

// ExitDesignator is called when production designator is exited.
func (s *Baseqasm3ParserListener) ExitDesignator(ctx *DesignatorContext) {}

// EnterDefcalTarget is called when production defcalTarget is entered.
func (s *Baseqasm3ParserListener) EnterDefcalTarget(ctx *DefcalTargetContext) {}

// ExitDefcalTarget is called when production defcalTarget is exited.
func (s *Baseqasm3ParserListener) ExitDefcalTarget(ctx *DefcalTargetContext) {}

// EnterDefcalArgumentDefinition is called when production defcalArgumentDefinition is entered.
func (s *Baseqasm3ParserListener) EnterDefcalArgumentDefinition(ctx *DefcalArgumentDefinitionContext) {
}

// ExitDefcalArgumentDefinition is called when production defcalArgumentDefinition is exited.
func (s *Baseqasm3ParserListener) ExitDefcalArgumentDefinition(ctx *DefcalArgumentDefinitionContext) {
}

// EnterDefcalOperand is called when production defcalOperand is entered.
func (s *Baseqasm3ParserListener) EnterDefcalOperand(ctx *DefcalOperandContext) {}

// ExitDefcalOperand is called when production defcalOperand is exited.
func (s *Baseqasm3ParserListener) ExitDefcalOperand(ctx *DefcalOperandContext) {}

// EnterGateOperand is called when production gateOperand is entered.
func (s *Baseqasm3ParserListener) EnterGateOperand(ctx *GateOperandContext) {}

// ExitGateOperand is called when production gateOperand is exited.
func (s *Baseqasm3ParserListener) ExitGateOperand(ctx *GateOperandContext) {}

// EnterExternArgument is called when production externArgument is entered.
func (s *Baseqasm3ParserListener) EnterExternArgument(ctx *ExternArgumentContext) {}

// ExitExternArgument is called when production externArgument is exited.
func (s *Baseqasm3ParserListener) ExitExternArgument(ctx *ExternArgumentContext) {}

// EnterArgumentDefinition is called when production argumentDefinition is entered.
func (s *Baseqasm3ParserListener) EnterArgumentDefinition(ctx *ArgumentDefinitionContext) {}

// ExitArgumentDefinition is called when production argumentDefinition is exited.
func (s *Baseqasm3ParserListener) ExitArgumentDefinition(ctx *ArgumentDefinitionContext) {}

// EnterArgumentDefinitionList is called when production argumentDefinitionList is entered.
func (s *Baseqasm3ParserListener) EnterArgumentDefinitionList(ctx *ArgumentDefinitionListContext) {}

// ExitArgumentDefinitionList is called when production argumentDefinitionList is exited.
func (s *Baseqasm3ParserListener) ExitArgumentDefinitionList(ctx *ArgumentDefinitionListContext) {}

// EnterDefcalArgumentDefinitionList is called when production defcalArgumentDefinitionList is entered.
func (s *Baseqasm3ParserListener) EnterDefcalArgumentDefinitionList(ctx *DefcalArgumentDefinitionListContext) {
}

// ExitDefcalArgumentDefinitionList is called when production defcalArgumentDefinitionList is exited.
func (s *Baseqasm3ParserListener) ExitDefcalArgumentDefinitionList(ctx *DefcalArgumentDefinitionListContext) {
}

// EnterDefcalOperandList is called when production defcalOperandList is entered.
func (s *Baseqasm3ParserListener) EnterDefcalOperandList(ctx *DefcalOperandListContext) {}

// ExitDefcalOperandList is called when production defcalOperandList is exited.
func (s *Baseqasm3ParserListener) ExitDefcalOperandList(ctx *DefcalOperandListContext) {}

// EnterExpressionList is called when production expressionList is entered.
func (s *Baseqasm3ParserListener) EnterExpressionList(ctx *ExpressionListContext) {}

// ExitExpressionList is called when production expressionList is exited.
func (s *Baseqasm3ParserListener) ExitExpressionList(ctx *ExpressionListContext) {}

// EnterIdentifierList is called when production identifierList is entered.
func (s *Baseqasm3ParserListener) EnterIdentifierList(ctx *IdentifierListContext) {}

// ExitIdentifierList is called when production identifierList is exited.
func (s *Baseqasm3ParserListener) ExitIdentifierList(ctx *IdentifierListContext) {}

// EnterGateOperandList is called when production gateOperandList is entered.
func (s *Baseqasm3ParserListener) EnterGateOperandList(ctx *GateOperandListContext) {}

// ExitGateOperandList is called when production gateOperandList is exited.
func (s *Baseqasm3ParserListener) ExitGateOperandList(ctx *GateOperandListContext) {}

// EnterExternArgumentList is called when production externArgumentList is entered.
func (s *Baseqasm3ParserListener) EnterExternArgumentList(ctx *ExternArgumentListContext) {}

// ExitExternArgumentList is called when production externArgumentList is exited.
func (s *Baseqasm3ParserListener) ExitExternArgumentList(ctx *ExternArgumentListContext) {}
