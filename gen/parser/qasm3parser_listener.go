// Code generated from qasm3Parser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // qasm3Parser
import "github.com/antlr4-go/antlr/v4"

// qasm3ParserListener is a complete listener for a parse tree produced by qasm3Parser.
type qasm3ParserListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterVersion is called when entering the version production.
	EnterVersion(c *VersionContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterAnnotation is called when entering the annotation production.
	EnterAnnotation(c *AnnotationContext)

	// EnterScope is called when entering the scope production.
	EnterScope(c *ScopeContext)

	// EnterPragma is called when entering the pragma production.
	EnterPragma(c *PragmaContext)

	// EnterStatementOrScope is called when entering the statementOrScope production.
	EnterStatementOrScope(c *StatementOrScopeContext)

	// EnterCalibrationGrammarStatement is called when entering the calibrationGrammarStatement production.
	EnterCalibrationGrammarStatement(c *CalibrationGrammarStatementContext)

	// EnterIncludeStatement is called when entering the includeStatement production.
	EnterIncludeStatement(c *IncludeStatementContext)

	// EnterBreakStatement is called when entering the breakStatement production.
	EnterBreakStatement(c *BreakStatementContext)

	// EnterContinueStatement is called when entering the continueStatement production.
	EnterContinueStatement(c *ContinueStatementContext)

	// EnterEndStatement is called when entering the endStatement production.
	EnterEndStatement(c *EndStatementContext)

	// EnterForStatement is called when entering the forStatement production.
	EnterForStatement(c *ForStatementContext)

	// EnterIfStatement is called when entering the ifStatement production.
	EnterIfStatement(c *IfStatementContext)

	// EnterReturnStatement is called when entering the returnStatement production.
	EnterReturnStatement(c *ReturnStatementContext)

	// EnterWhileStatement is called when entering the whileStatement production.
	EnterWhileStatement(c *WhileStatementContext)

	// EnterSwitchStatement is called when entering the switchStatement production.
	EnterSwitchStatement(c *SwitchStatementContext)

	// EnterSwitchCaseItem is called when entering the switchCaseItem production.
	EnterSwitchCaseItem(c *SwitchCaseItemContext)

	// EnterBarrierStatement is called when entering the barrierStatement production.
	EnterBarrierStatement(c *BarrierStatementContext)

	// EnterBoxStatement is called when entering the boxStatement production.
	EnterBoxStatement(c *BoxStatementContext)

	// EnterDelayStatement is called when entering the delayStatement production.
	EnterDelayStatement(c *DelayStatementContext)

	// EnterGateCallStatement is called when entering the gateCallStatement production.
	EnterGateCallStatement(c *GateCallStatementContext)

	// EnterMeasureArrowAssignmentStatement is called when entering the measureArrowAssignmentStatement production.
	EnterMeasureArrowAssignmentStatement(c *MeasureArrowAssignmentStatementContext)

	// EnterResetStatement is called when entering the resetStatement production.
	EnterResetStatement(c *ResetStatementContext)

	// EnterAliasDeclarationStatement is called when entering the aliasDeclarationStatement production.
	EnterAliasDeclarationStatement(c *AliasDeclarationStatementContext)

	// EnterClassicalDeclarationStatement is called when entering the classicalDeclarationStatement production.
	EnterClassicalDeclarationStatement(c *ClassicalDeclarationStatementContext)

	// EnterConstDeclarationStatement is called when entering the constDeclarationStatement production.
	EnterConstDeclarationStatement(c *ConstDeclarationStatementContext)

	// EnterIoDeclarationStatement is called when entering the ioDeclarationStatement production.
	EnterIoDeclarationStatement(c *IoDeclarationStatementContext)

	// EnterOldStyleDeclarationStatement is called when entering the oldStyleDeclarationStatement production.
	EnterOldStyleDeclarationStatement(c *OldStyleDeclarationStatementContext)

	// EnterQuantumDeclarationStatement is called when entering the quantumDeclarationStatement production.
	EnterQuantumDeclarationStatement(c *QuantumDeclarationStatementContext)

	// EnterDefStatement is called when entering the defStatement production.
	EnterDefStatement(c *DefStatementContext)

	// EnterExternStatement is called when entering the externStatement production.
	EnterExternStatement(c *ExternStatementContext)

	// EnterGateStatement is called when entering the gateStatement production.
	EnterGateStatement(c *GateStatementContext)

	// EnterAssignmentStatement is called when entering the assignmentStatement production.
	EnterAssignmentStatement(c *AssignmentStatementContext)

	// EnterExpressionStatement is called when entering the expressionStatement production.
	EnterExpressionStatement(c *ExpressionStatementContext)

	// EnterCalStatement is called when entering the calStatement production.
	EnterCalStatement(c *CalStatementContext)

	// EnterDefcalStatement is called when entering the defcalStatement production.
	EnterDefcalStatement(c *DefcalStatementContext)

	// EnterBitwiseXorExpression is called when entering the bitwiseXorExpression production.
	EnterBitwiseXorExpression(c *BitwiseXorExpressionContext)

	// EnterAdditiveExpression is called when entering the additiveExpression production.
	EnterAdditiveExpression(c *AdditiveExpressionContext)

	// EnterDurationofExpression is called when entering the durationofExpression production.
	EnterDurationofExpression(c *DurationofExpressionContext)

	// EnterParenthesisExpression is called when entering the parenthesisExpression production.
	EnterParenthesisExpression(c *ParenthesisExpressionContext)

	// EnterComparisonExpression is called when entering the comparisonExpression production.
	EnterComparisonExpression(c *ComparisonExpressionContext)

	// EnterMultiplicativeExpression is called when entering the multiplicativeExpression production.
	EnterMultiplicativeExpression(c *MultiplicativeExpressionContext)

	// EnterLogicalOrExpression is called when entering the logicalOrExpression production.
	EnterLogicalOrExpression(c *LogicalOrExpressionContext)

	// EnterCastExpression is called when entering the castExpression production.
	EnterCastExpression(c *CastExpressionContext)

	// EnterPowerExpression is called when entering the powerExpression production.
	EnterPowerExpression(c *PowerExpressionContext)

	// EnterBitwiseOrExpression is called when entering the bitwiseOrExpression production.
	EnterBitwiseOrExpression(c *BitwiseOrExpressionContext)

	// EnterCallExpression is called when entering the callExpression production.
	EnterCallExpression(c *CallExpressionContext)

	// EnterBitshiftExpression is called when entering the bitshiftExpression production.
	EnterBitshiftExpression(c *BitshiftExpressionContext)

	// EnterBitwiseAndExpression is called when entering the bitwiseAndExpression production.
	EnterBitwiseAndExpression(c *BitwiseAndExpressionContext)

	// EnterEqualityExpression is called when entering the equalityExpression production.
	EnterEqualityExpression(c *EqualityExpressionContext)

	// EnterLogicalAndExpression is called when entering the logicalAndExpression production.
	EnterLogicalAndExpression(c *LogicalAndExpressionContext)

	// EnterIndexExpression is called when entering the indexExpression production.
	EnterIndexExpression(c *IndexExpressionContext)

	// EnterUnaryExpression is called when entering the unaryExpression production.
	EnterUnaryExpression(c *UnaryExpressionContext)

	// EnterLiteralExpression is called when entering the literalExpression production.
	EnterLiteralExpression(c *LiteralExpressionContext)

	// EnterAliasExpression is called when entering the aliasExpression production.
	EnterAliasExpression(c *AliasExpressionContext)

	// EnterDeclarationExpression is called when entering the declarationExpression production.
	EnterDeclarationExpression(c *DeclarationExpressionContext)

	// EnterMeasureExpression is called when entering the measureExpression production.
	EnterMeasureExpression(c *MeasureExpressionContext)

	// EnterRangeExpression is called when entering the rangeExpression production.
	EnterRangeExpression(c *RangeExpressionContext)

	// EnterSetExpression is called when entering the setExpression production.
	EnterSetExpression(c *SetExpressionContext)

	// EnterArrayLiteral is called when entering the arrayLiteral production.
	EnterArrayLiteral(c *ArrayLiteralContext)

	// EnterIndexOperator is called when entering the indexOperator production.
	EnterIndexOperator(c *IndexOperatorContext)

	// EnterIndexedIdentifier is called when entering the indexedIdentifier production.
	EnterIndexedIdentifier(c *IndexedIdentifierContext)

	// EnterReturnSignature is called when entering the returnSignature production.
	EnterReturnSignature(c *ReturnSignatureContext)

	// EnterGateModifier is called when entering the gateModifier production.
	EnterGateModifier(c *GateModifierContext)

	// EnterScalarType is called when entering the scalarType production.
	EnterScalarType(c *ScalarTypeContext)

	// EnterQubitType is called when entering the qubitType production.
	EnterQubitType(c *QubitTypeContext)

	// EnterArrayType is called when entering the arrayType production.
	EnterArrayType(c *ArrayTypeContext)

	// EnterArrayReferenceType is called when entering the arrayReferenceType production.
	EnterArrayReferenceType(c *ArrayReferenceTypeContext)

	// EnterDesignator is called when entering the designator production.
	EnterDesignator(c *DesignatorContext)

	// EnterDefcalTarget is called when entering the defcalTarget production.
	EnterDefcalTarget(c *DefcalTargetContext)

	// EnterDefcalArgumentDefinition is called when entering the defcalArgumentDefinition production.
	EnterDefcalArgumentDefinition(c *DefcalArgumentDefinitionContext)

	// EnterDefcalOperand is called when entering the defcalOperand production.
	EnterDefcalOperand(c *DefcalOperandContext)

	// EnterGateOperand is called when entering the gateOperand production.
	EnterGateOperand(c *GateOperandContext)

	// EnterExternArgument is called when entering the externArgument production.
	EnterExternArgument(c *ExternArgumentContext)

	// EnterArgumentDefinition is called when entering the argumentDefinition production.
	EnterArgumentDefinition(c *ArgumentDefinitionContext)

	// EnterArgumentDefinitionList is called when entering the argumentDefinitionList production.
	EnterArgumentDefinitionList(c *ArgumentDefinitionListContext)

	// EnterDefcalArgumentDefinitionList is called when entering the defcalArgumentDefinitionList production.
	EnterDefcalArgumentDefinitionList(c *DefcalArgumentDefinitionListContext)

	// EnterDefcalOperandList is called when entering the defcalOperandList production.
	EnterDefcalOperandList(c *DefcalOperandListContext)

	// EnterExpressionList is called when entering the expressionList production.
	EnterExpressionList(c *ExpressionListContext)

	// EnterIdentifierList is called when entering the identifierList production.
	EnterIdentifierList(c *IdentifierListContext)

	// EnterGateOperandList is called when entering the gateOperandList production.
	EnterGateOperandList(c *GateOperandListContext)

	// EnterExternArgumentList is called when entering the externArgumentList production.
	EnterExternArgumentList(c *ExternArgumentListContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitVersion is called when exiting the version production.
	ExitVersion(c *VersionContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitAnnotation is called when exiting the annotation production.
	ExitAnnotation(c *AnnotationContext)

	// ExitScope is called when exiting the scope production.
	ExitScope(c *ScopeContext)

	// ExitPragma is called when exiting the pragma production.
	ExitPragma(c *PragmaContext)

	// ExitStatementOrScope is called when exiting the statementOrScope production.
	ExitStatementOrScope(c *StatementOrScopeContext)

	// ExitCalibrationGrammarStatement is called when exiting the calibrationGrammarStatement production.
	ExitCalibrationGrammarStatement(c *CalibrationGrammarStatementContext)

	// ExitIncludeStatement is called when exiting the includeStatement production.
	ExitIncludeStatement(c *IncludeStatementContext)

	// ExitBreakStatement is called when exiting the breakStatement production.
	ExitBreakStatement(c *BreakStatementContext)

	// ExitContinueStatement is called when exiting the continueStatement production.
	ExitContinueStatement(c *ContinueStatementContext)

	// ExitEndStatement is called when exiting the endStatement production.
	ExitEndStatement(c *EndStatementContext)

	// ExitForStatement is called when exiting the forStatement production.
	ExitForStatement(c *ForStatementContext)

	// ExitIfStatement is called when exiting the ifStatement production.
	ExitIfStatement(c *IfStatementContext)

	// ExitReturnStatement is called when exiting the returnStatement production.
	ExitReturnStatement(c *ReturnStatementContext)

	// ExitWhileStatement is called when exiting the whileStatement production.
	ExitWhileStatement(c *WhileStatementContext)

	// ExitSwitchStatement is called when exiting the switchStatement production.
	ExitSwitchStatement(c *SwitchStatementContext)

	// ExitSwitchCaseItem is called when exiting the switchCaseItem production.
	ExitSwitchCaseItem(c *SwitchCaseItemContext)

	// ExitBarrierStatement is called when exiting the barrierStatement production.
	ExitBarrierStatement(c *BarrierStatementContext)

	// ExitBoxStatement is called when exiting the boxStatement production.
	ExitBoxStatement(c *BoxStatementContext)

	// ExitDelayStatement is called when exiting the delayStatement production.
	ExitDelayStatement(c *DelayStatementContext)

	// ExitGateCallStatement is called when exiting the gateCallStatement production.
	ExitGateCallStatement(c *GateCallStatementContext)

	// ExitMeasureArrowAssignmentStatement is called when exiting the measureArrowAssignmentStatement production.
	ExitMeasureArrowAssignmentStatement(c *MeasureArrowAssignmentStatementContext)

	// ExitResetStatement is called when exiting the resetStatement production.
	ExitResetStatement(c *ResetStatementContext)

	// ExitAliasDeclarationStatement is called when exiting the aliasDeclarationStatement production.
	ExitAliasDeclarationStatement(c *AliasDeclarationStatementContext)

	// ExitClassicalDeclarationStatement is called when exiting the classicalDeclarationStatement production.
	ExitClassicalDeclarationStatement(c *ClassicalDeclarationStatementContext)

	// ExitConstDeclarationStatement is called when exiting the constDeclarationStatement production.
	ExitConstDeclarationStatement(c *ConstDeclarationStatementContext)

	// ExitIoDeclarationStatement is called when exiting the ioDeclarationStatement production.
	ExitIoDeclarationStatement(c *IoDeclarationStatementContext)

	// ExitOldStyleDeclarationStatement is called when exiting the oldStyleDeclarationStatement production.
	ExitOldStyleDeclarationStatement(c *OldStyleDeclarationStatementContext)

	// ExitQuantumDeclarationStatement is called when exiting the quantumDeclarationStatement production.
	ExitQuantumDeclarationStatement(c *QuantumDeclarationStatementContext)

	// ExitDefStatement is called when exiting the defStatement production.
	ExitDefStatement(c *DefStatementContext)

	// ExitExternStatement is called when exiting the externStatement production.
	ExitExternStatement(c *ExternStatementContext)

	// ExitGateStatement is called when exiting the gateStatement production.
	ExitGateStatement(c *GateStatementContext)

	// ExitAssignmentStatement is called when exiting the assignmentStatement production.
	ExitAssignmentStatement(c *AssignmentStatementContext)

	// ExitExpressionStatement is called when exiting the expressionStatement production.
	ExitExpressionStatement(c *ExpressionStatementContext)

	// ExitCalStatement is called when exiting the calStatement production.
	ExitCalStatement(c *CalStatementContext)

	// ExitDefcalStatement is called when exiting the defcalStatement production.
	ExitDefcalStatement(c *DefcalStatementContext)

	// ExitBitwiseXorExpression is called when exiting the bitwiseXorExpression production.
	ExitBitwiseXorExpression(c *BitwiseXorExpressionContext)

	// ExitAdditiveExpression is called when exiting the additiveExpression production.
	ExitAdditiveExpression(c *AdditiveExpressionContext)

	// ExitDurationofExpression is called when exiting the durationofExpression production.
	ExitDurationofExpression(c *DurationofExpressionContext)

	// ExitParenthesisExpression is called when exiting the parenthesisExpression production.
	ExitParenthesisExpression(c *ParenthesisExpressionContext)

	// ExitComparisonExpression is called when exiting the comparisonExpression production.
	ExitComparisonExpression(c *ComparisonExpressionContext)

	// ExitMultiplicativeExpression is called when exiting the multiplicativeExpression production.
	ExitMultiplicativeExpression(c *MultiplicativeExpressionContext)

	// ExitLogicalOrExpression is called when exiting the logicalOrExpression production.
	ExitLogicalOrExpression(c *LogicalOrExpressionContext)

	// ExitCastExpression is called when exiting the castExpression production.
	ExitCastExpression(c *CastExpressionContext)

	// ExitPowerExpression is called when exiting the powerExpression production.
	ExitPowerExpression(c *PowerExpressionContext)

	// ExitBitwiseOrExpression is called when exiting the bitwiseOrExpression production.
	ExitBitwiseOrExpression(c *BitwiseOrExpressionContext)

	// ExitCallExpression is called when exiting the callExpression production.
	ExitCallExpression(c *CallExpressionContext)

	// ExitBitshiftExpression is called when exiting the bitshiftExpression production.
	ExitBitshiftExpression(c *BitshiftExpressionContext)

	// ExitBitwiseAndExpression is called when exiting the bitwiseAndExpression production.
	ExitBitwiseAndExpression(c *BitwiseAndExpressionContext)

	// ExitEqualityExpression is called when exiting the equalityExpression production.
	ExitEqualityExpression(c *EqualityExpressionContext)

	// ExitLogicalAndExpression is called when exiting the logicalAndExpression production.
	ExitLogicalAndExpression(c *LogicalAndExpressionContext)

	// ExitIndexExpression is called when exiting the indexExpression production.
	ExitIndexExpression(c *IndexExpressionContext)

	// ExitUnaryExpression is called when exiting the unaryExpression production.
	ExitUnaryExpression(c *UnaryExpressionContext)

	// ExitLiteralExpression is called when exiting the literalExpression production.
	ExitLiteralExpression(c *LiteralExpressionContext)

	// ExitAliasExpression is called when exiting the aliasExpression production.
	ExitAliasExpression(c *AliasExpressionContext)

	// ExitDeclarationExpression is called when exiting the declarationExpression production.
	ExitDeclarationExpression(c *DeclarationExpressionContext)

	// ExitMeasureExpression is called when exiting the measureExpression production.
	ExitMeasureExpression(c *MeasureExpressionContext)

	// ExitRangeExpression is called when exiting the rangeExpression production.
	ExitRangeExpression(c *RangeExpressionContext)

	// ExitSetExpression is called when exiting the setExpression production.
	ExitSetExpression(c *SetExpressionContext)

	// ExitArrayLiteral is called when exiting the arrayLiteral production.
	ExitArrayLiteral(c *ArrayLiteralContext)

	// ExitIndexOperator is called when exiting the indexOperator production.
	ExitIndexOperator(c *IndexOperatorContext)

	// ExitIndexedIdentifier is called when exiting the indexedIdentifier production.
	ExitIndexedIdentifier(c *IndexedIdentifierContext)

	// ExitReturnSignature is called when exiting the returnSignature production.
	ExitReturnSignature(c *ReturnSignatureContext)

	// ExitGateModifier is called when exiting the gateModifier production.
	ExitGateModifier(c *GateModifierContext)

	// ExitScalarType is called when exiting the scalarType production.
	ExitScalarType(c *ScalarTypeContext)

	// ExitQubitType is called when exiting the qubitType production.
	ExitQubitType(c *QubitTypeContext)

	// ExitArrayType is called when exiting the arrayType production.
	ExitArrayType(c *ArrayTypeContext)

	// ExitArrayReferenceType is called when exiting the arrayReferenceType production.
	ExitArrayReferenceType(c *ArrayReferenceTypeContext)

	// ExitDesignator is called when exiting the designator production.
	ExitDesignator(c *DesignatorContext)

	// ExitDefcalTarget is called when exiting the defcalTarget production.
	ExitDefcalTarget(c *DefcalTargetContext)

	// ExitDefcalArgumentDefinition is called when exiting the defcalArgumentDefinition production.
	ExitDefcalArgumentDefinition(c *DefcalArgumentDefinitionContext)

	// ExitDefcalOperand is called when exiting the defcalOperand production.
	ExitDefcalOperand(c *DefcalOperandContext)

	// ExitGateOperand is called when exiting the gateOperand production.
	ExitGateOperand(c *GateOperandContext)

	// ExitExternArgument is called when exiting the externArgument production.
	ExitExternArgument(c *ExternArgumentContext)

	// ExitArgumentDefinition is called when exiting the argumentDefinition production.
	ExitArgumentDefinition(c *ArgumentDefinitionContext)

	// ExitArgumentDefinitionList is called when exiting the argumentDefinitionList production.
	ExitArgumentDefinitionList(c *ArgumentDefinitionListContext)

	// ExitDefcalArgumentDefinitionList is called when exiting the defcalArgumentDefinitionList production.
	ExitDefcalArgumentDefinitionList(c *DefcalArgumentDefinitionListContext)

	// ExitDefcalOperandList is called when exiting the defcalOperandList production.
	ExitDefcalOperandList(c *DefcalOperandListContext)

	// ExitExpressionList is called when exiting the expressionList production.
	ExitExpressionList(c *ExpressionListContext)

	// ExitIdentifierList is called when exiting the identifierList production.
	ExitIdentifierList(c *IdentifierListContext)

	// ExitGateOperandList is called when exiting the gateOperandList production.
	ExitGateOperandList(c *GateOperandListContext)

	// ExitExternArgumentList is called when exiting the externArgumentList production.
	ExitExternArgumentList(c *ExternArgumentListContext)
}
