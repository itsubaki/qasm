package visitor

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/qasm/gen/parser"
)

var (
	ErrAlreadyDeclared      = errors.New("already declared")
	ErrQubitNotFound        = errors.New("qubit not found")
	ErrClassicalBitNotFound = errors.New("classical bit not found")
	ErrGateNotFound         = errors.New("gate not found")
	ErrConstNotFound        = errors.New("const not found")
	ErrUnexpected           = errors.New("unexpected")
	ErrNotImplemented       = errors.New("not implemented")
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

func (v *Visitor) VisitTerminal(node antlr.TerminalNode) interface{} {
	return node.GetText()
}

func (v *Visitor) VisitErrorNode(node antlr.ErrorNode) interface{} {
	return node.GetText()
}

func (v *Visitor) VisitChildren(node antlr.RuleNode) interface{} {
	return fmt.Errorf("VisitChildren: %w", ErrNotImplemented)
}

func (v *Visitor) VisitProgram(ctx *parser.ProgramContext) interface{} {
	if ctx.Version() != nil {
		v.Environ.Version = v.Visit(ctx.Version()).(string)
	}

	var result interface{}
	for _, s := range ctx.AllStatementOrScope() {
		if res := v.Visit(s); res != nil {
			result = res
			break
		}
	}

	return result
}

func (v *Visitor) VisitVersion(ctx *parser.VersionContext) interface{} {
	return v.Visit(ctx.VersionSpecifier())
}

func (v *Visitor) VisitPragma(ctx *parser.PragmaContext) interface{} {
	return v.Visit(ctx.RemainingLineContent())
}

func (v *Visitor) VisitAnnotation(ctx *parser.AnnotationContext) interface{} {
	return fmt.Errorf("VisitAnnotation: %w", ErrNotImplemented)
}

func (v *Visitor) VisitStatement(ctx *parser.StatementContext) interface{} {
	statements := []antlr.ParseTree{
		ctx.Pragma(),
		ctx.AliasDeclarationStatement(),
		ctx.AssignmentStatement(),
		ctx.BarrierStatement(),
		ctx.BoxStatement(),
		ctx.BreakStatement(),
		ctx.CalStatement(),
		ctx.CalibrationGrammarStatement(),
		ctx.ClassicalDeclarationStatement(),
		ctx.ConstDeclarationStatement(),
		ctx.ContinueStatement(),
		ctx.DefStatement(),
		ctx.DefcalStatement(),
		ctx.DelayStatement(),
		ctx.EndStatement(),
		ctx.ExpressionStatement(),
		ctx.ExternStatement(),
		ctx.ForStatement(),
		ctx.GateCallStatement(),
		ctx.GateStatement(),
		ctx.IfStatement(),
		ctx.IncludeStatement(),
		ctx.IoDeclarationStatement(),
		ctx.MeasureArrowAssignmentStatement(),
		ctx.OldStyleDeclarationStatement(),
		ctx.QuantumDeclarationStatement(),
		ctx.ResetStatement(),
		ctx.ReturnStatement(),
		ctx.SwitchStatement(),
		ctx.WhileStatement(),
	}

	for _, s := range statements {
		if s == nil {
			continue
		}

		return v.Visit(s)
	}

	return fmt.Errorf("statement=%s: %w", ctx.GetText(), ErrUnexpected)
}

func (v *Visitor) VisitScope(ctx *parser.ScopeContext) interface{} {
	var list []interface{}
	for _, s := range ctx.AllStatementOrScope() {
		list = append(list, v.Visit(s))
	}

	return list
}

func (v *Visitor) VisitStatementOrScope(ctx *parser.StatementOrScopeContext) interface{} {
	if ctx.Statement() != nil {
		return v.Visit(ctx.Statement())
	}

	return v.Visit(ctx.Scope())
}

func (v *Visitor) VisitIncludeStatement(ctx *parser.IncludeStatementContext) interface{} {
	return fmt.Errorf("VisitIncludeStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitIfStatement(ctx *parser.IfStatementContext) interface{} {
	return fmt.Errorf("VisitIfStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitForStatement(ctx *parser.ForStatementContext) interface{} {
	return fmt.Errorf("VisitForStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBreakStatement(ctx *parser.BreakStatementContext) interface{} {
	return fmt.Errorf("VisitBreakStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitContinueStatement(ctx *parser.ContinueStatementContext) interface{} {
	return fmt.Errorf("VisitContinueStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitWhileStatement(ctx *parser.WhileStatementContext) interface{} {
	return fmt.Errorf("VisitWhileStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitSwitchStatement(ctx *parser.SwitchStatementContext) interface{} {
	return fmt.Errorf("VisitSwitchStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitSwitchCaseItem(ctx *parser.SwitchCaseItemContext) interface{} {
	return fmt.Errorf("VisitSwitchCaseItem: %w", ErrNotImplemented)
}

func (v *Visitor) VisitEndStatement(ctx *parser.EndStatementContext) interface{} {
	return fmt.Errorf("VisitEndStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitCalibrationGrammarStatement(ctx *parser.CalibrationGrammarStatementContext) interface{} {
	return fmt.Errorf("VisitCalibrationGrammarStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBarrierStatement(ctx *parser.BarrierStatementContext) interface{} {
	return fmt.Errorf("VisitBarrierStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBoxStatement(ctx *parser.BoxStatementContext) interface{} {
	return fmt.Errorf("VisitBoxStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDelayStatement(ctx *parser.DelayStatementContext) interface{} {
	return fmt.Errorf("VisitDelayStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitReturnStatement(ctx *parser.ReturnStatementContext) interface{} {
	if ctx.MeasureExpression() != nil {
		return v.Visit(ctx.MeasureExpression())
	}

	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitGateStatement(ctx *parser.GateStatementContext) interface{} {
	fmt.Println(ctx.GetText())
	return nil
}

func (v *Visitor) Params(xlist parser.IExpressionListContext) ([]float64, error) {
	var params []float64
	for _, e := range v.Visit(xlist).([]interface{}) {
		switch val := e.(type) {
		case float64:
			params = append(params, val)
		case int64:
			params = append(params, float64(val))
		default:
			return nil, fmt.Errorf("param=%v: %w", val, ErrUnexpected)
		}
	}

	return params, nil
}

func (v *Visitor) Modify(u matrix.Matrix, qargs []q.Qubit, modifier []parser.IGateModifierContext) (matrix.Matrix, error) {
	var ctrl, negctrl []q.Qubit
	for i, mod := range modifier {
		// https://openqasm.com/language/gates.html#inverse-modifier
		// The inverse of a controlled operation is defined by inverting the control unitary. That is, inv @ ctrl @ U = ctrl @ inv @ U.
		switch {
		case mod.CTRL() != nil:
			ctrl = append(ctrl, qargs[i])
		case mod.NEGCTRL() != nil:
			ctrl, negctrl = append(ctrl, qargs[i]), append(negctrl, qargs[i])
		case mod.INV() != nil:
			u = u.Dagger()
		case mod.POW() != nil:
			x := v.Visit(mod).(int64)
			u = matrix.ApplyN(u, int(x))
		default:
			return nil, fmt.Errorf("modifier=%s: %w", mod.GetText(), ErrUnexpected)
		}
	}

	switch len(ctrl) {
	case 0:
		n := v.qsim.NumberOfBit()
		u = gate.TensorProduct(u, n, q.Index(qargs...))
	default:
		n := v.qsim.NumberOfBit()
		u = gate.Controlled(u, n, q.Index(ctrl...), qargs[len(qargs)-1].Index())

		if len(negctrl) > 0 {
			x := gate.TensorProduct(gate.X(), n, q.Index(negctrl...))
			u = matrix.Apply(x, u, x)
		}
	}

	return u, nil
}

func (v *Visitor) VisitGateCallStatement(ctx *parser.GateCallStatementContext) interface{} {
	id := v.Visit(ctx.Identifier())
	switch id {
	case "U":
		params, err := v.Params(ctx.ExpressionList())
		if err != nil {
			return fmt.Errorf("params: %w", err)
		}
		u := gate.U(params[0], params[1], params[2])

		qargs := v.Visit(ctx.GateOperandList()).([]q.Qubit)
		u, err = v.Modify(u, qargs, ctx.AllGateModifier())
		if err != nil {
			return fmt.Errorf("modify: %w", err)
		}

		v.qsim.Apply(u)
		return nil
	default:
		return fmt.Errorf("idenfitier=%s: %w", id, ErrGateNotFound)
	}
}

func (v *Visitor) MeasureAssignment(identifier parser.IIndexedIdentifierContext, measure parser.IMeasureExpressionContext) error {
	measured := v.Visit(measure)
	if identifier == nil {
		return nil
	}

	operand := v.Visit(identifier.Identifier()).(string)
	bits, ok := v.Environ.GetClassicalBit(operand)
	if !ok {
		return fmt.Errorf("operand=%s: %w", operand, ErrClassicalBitNotFound)
	}

	index := v.Visit(identifier).([]int64)
	if len(index) == 0 {
		copy(bits, measured.([]int64))
		return nil
	}

	for i, m := range measured.([]int64) {
		bits[index[i]] = m
	}

	return nil
}

func (v *Visitor) VisitMeasureArrowAssignmentStatement(ctx *parser.MeasureArrowAssignmentStatementContext) interface{} {
	return v.MeasureAssignment(ctx.IndexedIdentifier(), ctx.MeasureExpression())
}

func (v *Visitor) VisitAssignmentStatement(ctx *parser.AssignmentStatementContext) interface{} {
	if ctx.MeasureExpression() != nil {
		return v.MeasureAssignment(ctx.IndexedIdentifier(), ctx.MeasureExpression())
	}

	// TODO: v.Visit(ctx.Expression())
	return fmt.Errorf("statement=%s: %w", ctx.GetText(), ErrNotImplemented)
}

func (v *Visitor) VisitResetStatement(ctx *parser.ResetStatementContext) interface{} {
	qargs := v.Visit(ctx.GateOperand()).([]q.Qubit)
	v.qsim.Reset(qargs...)
	return nil
}

func (v *Visitor) VisitQuantumDeclarationStatement(ctx *parser.QuantumDeclarationStatementContext) interface{} {
	id := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.Environ.GetQubit(id); ok {
		return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
	}

	size := v.Visit(ctx.QubitType()).(int64)
	v.Environ.Qubit[id] = v.qsim.ZeroWith(int(size))
	return nil
}

func (v *Visitor) VisitClassicalDeclarationStatement(ctx *parser.ClassicalDeclarationStatementContext) interface{} {
	switch {
	case ctx.ScalarType().BIT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.Environ.GetClassicalBit(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			bits := v.Visit(ctx.DeclarationExpression()).([]int64)
			v.Environ.ClassicalBit[id] = bits
			return nil
		}

		size := v.Visit(ctx.ScalarType()).(int64)
		v.Environ.ClassicalBit[id] = make([]int64, int(size))
		return nil
	default:
		return fmt.Errorf("scalar type=%s: %w", ctx.ScalarType().GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitConstDeclarationStatement(ctx *parser.ConstDeclarationStatementContext) interface{} {
	return fmt.Errorf("VisitConstDeclarationStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitAliasDeclarationStatement(ctx *parser.AliasDeclarationStatementContext) interface{} {
	return fmt.Errorf("VisitAliasDeclarationStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitIoDeclarationStatement(ctx *parser.IoDeclarationStatementContext) interface{} {
	return fmt.Errorf("VisitIoDeclarationStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitOldStyleDeclarationStatement(ctx *parser.OldStyleDeclarationStatementContext) interface{} {
	return fmt.Errorf("VisitOldStyleDeclarationStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefStatement(ctx *parser.DefStatementContext) interface{} {
	return fmt.Errorf("VisitDefStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitExternStatement(ctx *parser.ExternStatementContext) interface{} {
	return fmt.Errorf("VisitExternStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitCalStatement(ctx *parser.CalStatementContext) interface{} {
	return fmt.Errorf("VisitCalStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalStatement(ctx *parser.DefcalStatementContext) interface{} {
	return fmt.Errorf("VisitDefcalStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitParenthesisExpression(ctx *parser.ParenthesisExpressionContext) interface{} {
	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitExpressionStatement(ctx *parser.ExpressionStatementContext) interface{} {
	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) interface{} {
	switch {
	case ctx.Identifier() != nil:
		x := v.Visit(ctx.Identifier()).(string)
		lit, ok := BuiltinConst[x]
		if !ok {
			return fmt.Errorf("identifier=%s: %w", x, ErrConstNotFound)
		}

		return lit
	case ctx.DecimalIntegerLiteral() != nil:
		x := v.Visit(ctx.DecimalIntegerLiteral()).(string)
		lit, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int: x=%s: %w", x, err)
		}

		return lit
	case ctx.FloatLiteral() != nil:
		x := v.Visit(ctx.FloatLiteral()).(string)
		lit, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return fmt.Errorf("parse float: x=%s: %w", x, err)
		}

		return lit
	case ctx.BooleanLiteral() != nil:
		x := v.Visit(ctx.BooleanLiteral()).(string)
		lit, err := strconv.ParseBool(x)
		if err != nil {
			return fmt.Errorf("parse bool: x=%s: %w", x, err)
		}

		return lit
	case ctx.BitstringLiteral() != nil:
		x := v.Visit(ctx.BitstringLiteral()).(string)
		bistring := strings.Trim(x, "\"")

		lit := make([]int64, len(bistring))
		for i, b := range bistring {
			lit[i] = int64(b - '0')
		}

		return lit
	default:
		return fmt.Errorf("x=%s: %w", ctx.GetText(), ErrUnexpected)
	}
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
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpected)
		}
	}

	switch {
	case ctx.PLUS() != nil:
		return operand[0] + operand[1]
	case ctx.MINUS() != nil:
		return operand[0] - operand[1]
	default:
		return fmt.Errorf("operator=%s: %w", ctx.GetOp().GetText(), ErrUnexpected)
	}
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
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpected)
		}
	}

	switch {
	case ctx.ASTERISK() != nil:
		return operand[0] * operand[1]
	case ctx.SLASH() != nil:
		return operand[0] / operand[1]
	case ctx.PERCENT() != nil:
		return float64(int64(operand[0]) % int64(operand[1]))
	default:
		return fmt.Errorf("operator=%s: %w", ctx.GetOp().GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitEqualityExpression(ctx *parser.EqualityExpressionContext) interface{} {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		case bool:
			var o float64
			if val {
				o = 1
			}

			operand = append(operand, o)
		default:
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpected)
		}
	}

	op := v.Visit(ctx.EqualityOperator()).(string)
	switch op {
	case "==":
		return operand[0] == operand[1]
	case "!=":
		return operand[0] != operand[1]
	default:
		return fmt.Errorf("operator=%s: %w", op, ErrUnexpected)
	}
}

func (v *Visitor) VisitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) interface{} {
	result := true
	for _, x := range ctx.AllExpression() {
		if result = result && v.Visit(x).(bool); !result {
			return false
		}
	}

	return true
}

func (v *Visitor) VisitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) interface{} {
	var result bool
	for _, x := range ctx.AllExpression() {
		if result = result || v.Visit(x).(bool); result {
			return true
		}
	}

	return false
}

func (v *Visitor) VisitBitwiseAndExpression(ctx *parser.BitwiseAndExpressionContext) interface{} {
	return fmt.Errorf("VisitBitwiseAndExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBitwiseOrExpression(ctx *parser.BitwiseOrExpressionContext) interface{} {
	return fmt.Errorf("VisitBitwiseOrExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBitwiseXorExpression(ctx *parser.BitwiseXorExpressionContext) interface{} {
	return fmt.Errorf("VisitBitwiseXorExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBitshiftExpression(ctx *parser.BitshiftExpressionContext) interface{} {
	return fmt.Errorf("VisitBitshiftExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDurationofExpression(ctx *parser.DurationofExpressionContext) interface{} {
	return fmt.Errorf("VisitDurationofExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitComparisonExpression(ctx *parser.ComparisonExpressionContext) interface{} {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		default:
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpected)
		}
	}

	op := v.Visit(ctx.ComparisonOperator()).(string)
	switch op {
	case "<":
		return operand[0] < operand[1]
	case "<=":
		return operand[0] <= operand[1]
	case ">":
		return operand[0] > operand[1]
	case ">=":
		return operand[0] >= operand[1]
	default:
		return fmt.Errorf("operator=%s: %w", op, ErrUnexpected)
	}
}

func (v *Visitor) VisitCastExpression(ctx *parser.CastExpressionContext) interface{} {
	return fmt.Errorf("VisitCastExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitPowerExpression(ctx *parser.PowerExpressionContext) interface{} {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		default:
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpected)
		}
	}

	return math.Pow(operand[0], operand[1])
}

func (v *Visitor) VisitCallExpression(ctx *parser.CallExpressionContext) interface{} {
	return fmt.Errorf("VisitCallExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitUnaryExpression(ctx *parser.UnaryExpressionContext) interface{} {
	operand := v.Visit(ctx.Expression())
	switch {
	case ctx.MINUS() != nil:
		switch val := operand.(type) {
		case float64:
			return -1 * val
		case int64:
			return -1 * val
		default:
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpected)
		}
	case ctx.EXCLAMATION_POINT() != nil:
		switch val := operand.(type) {
		case bool:
			return !val
		default:
			return fmt.Errorf("operand=%v: %w", val, ErrUnexpected)
		}
	case ctx.TILDE() != nil:
		return fmt.Errorf("VisitUnaryExpression: %w", ErrNotImplemented)
	default:
		return fmt.Errorf("operator=%s: %w", ctx.GetOp().GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitAliasExpression(ctx *parser.AliasExpressionContext) interface{} {
	return fmt.Errorf("VisitAliasExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDeclarationExpression(ctx *parser.DeclarationExpressionContext) interface{} {
	if ctx.MeasureExpression() != nil {
		return v.Visit(ctx.MeasureExpression())
	}

	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitRangeExpression(ctx *parser.RangeExpressionContext) interface{} {
	return fmt.Errorf("VisitRangeExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitSetExpression(ctx *parser.SetExpressionContext) interface{} {
	return fmt.Errorf("VisitSetExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitIndexExpression(ctx *parser.IndexExpressionContext) interface{} {
	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitMeasureExpression(ctx *parser.MeasureExpressionContext) interface{} {
	qargs := v.Visit(ctx.GateOperand()).([]q.Qubit)
	v.qsim.Measure(qargs...)

	var bits []int64
	for _, q := range qargs {
		bits = append(bits, v.qsim.State(q)[0].Int(0))
	}

	return bits
}

func (v *Visitor) VisitArrayLiteral(ctx *parser.ArrayLiteralContext) interface{} {
	return fmt.Errorf("VisitArrayLiteral: %w", ErrNotImplemented)
}

func (v *Visitor) VisitIndexOperator(ctx *parser.IndexOperatorContext) interface{} {
	var list []interface{}
	for _, x := range ctx.AllExpression() {
		list = append(list, v.Visit(x))
	}

	return list
}

func (v *Visitor) VisitIndexedIdentifier(ctx *parser.IndexedIdentifierContext) interface{} {
	var index []int64
	for _, op := range ctx.AllIndexOperator() {
		for _, v := range v.Visit(op).([]interface{}) {
			index = append(index, v.(int64))
		}
	}

	return index
}

func (v *Visitor) VisitGateModifier(ctx *parser.GateModifierContext) interface{} {
	if ctx.Expression() != nil {
		return v.Visit(ctx.Expression()).(int64)
	}

	return int64(1)
}

func (v *Visitor) VisitReturnSignature(ctx *parser.ReturnSignatureContext) interface{} {
	return fmt.Errorf("VisitReturnSignature: %w", ErrNotImplemented)
}

func (v *Visitor) VisitScalarType(ctx *parser.ScalarTypeContext) interface{} {
	if ctx.Designator() != nil {
		return v.Visit(ctx.Designator()).(int64)
	}

	return int64(1)
}

func (v *Visitor) VisitQubitType(ctx *parser.QubitTypeContext) interface{} {
	if ctx.Designator() != nil {
		return v.Visit(ctx.Designator()).(int64)
	}

	return int64(1)
}

func (v *Visitor) VisitArrayType(ctx *parser.ArrayTypeContext) interface{} {
	return fmt.Errorf("VisitArrayType: %w", ErrNotImplemented)
}

func (v *Visitor) VisitArrayReferenceType(ctx *parser.ArrayReferenceTypeContext) interface{} {
	return fmt.Errorf("VisitArrayReferenceType: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDesignator(ctx *parser.DesignatorContext) interface{} {
	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitDefcalTarget(ctx *parser.DefcalTargetContext) interface{} {
	return fmt.Errorf("VisitDefcalTarget: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalArgumentDefinition(ctx *parser.DefcalArgumentDefinitionContext) interface{} {
	return fmt.Errorf("VisitDefcalArgumentDefinition: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalOperand(ctx *parser.DefcalOperandContext) interface{} {
	return fmt.Errorf("VisitDefcalOperand: %w", ErrNotImplemented)
}

func (v *Visitor) VisitGateOperand(ctx *parser.GateOperandContext) interface{} {
	indexID := ctx.IndexedIdentifier()

	operand := v.Visit(indexID.Identifier()).(string)
	qb, ok := v.Environ.GetQubit(operand)
	if !ok {
		return fmt.Errorf("operand=%s: %w", operand, ErrQubitNotFound)
	}

	index := v.Visit(indexID).([]int64)
	if len(index) == 0 {
		return qb
	}

	var list []q.Qubit
	for _, idx := range index {
		list = append(list, qb[idx])
	}

	return list
}

func (v *Visitor) VisitExternArgument(ctx *parser.ExternArgumentContext) interface{} {
	return fmt.Errorf("VisitExternArgument: %w", ErrNotImplemented)
}

func (v *Visitor) VisitArgumentDefinition(ctx *parser.ArgumentDefinitionContext) interface{} {
	return fmt.Errorf("VisitArgumentDefinition: %w", ErrNotImplemented)
}

func (v *Visitor) VisitArgumentDefinitionList(ctx *parser.ArgumentDefinitionListContext) interface{} {
	var list []interface{}
	for _, def := range ctx.AllArgumentDefinition() {
		list = append(list, v.Visit(def))
	}

	return list
}

func (v *Visitor) VisitDefcalArgumentDefinitionList(ctx *parser.DefcalArgumentDefinitionListContext) interface{} {
	var list []interface{}
	for _, def := range ctx.AllDefcalArgumentDefinition() {
		list = append(list, v.Visit(def))
	}

	return list
}

func (v *Visitor) VisitDefcalOperandList(ctx *parser.DefcalOperandListContext) interface{} {
	var list []interface{}
	for _, o := range ctx.AllDefcalOperand() {
		list = append(list, v.Visit(o))
	}

	return list
}

func (v *Visitor) VisitExpressionList(ctx *parser.ExpressionListContext) interface{} {
	var list []interface{}
	for _, x := range ctx.AllExpression() {
		list = append(list, v.Visit(x))
	}

	return list
}

func (v *Visitor) VisitIdentifierList(ctx *parser.IdentifierListContext) interface{} {
	var list []string
	for _, id := range ctx.AllIdentifier() {
		list = append(list, v.Visit(id).(string))
	}

	return list
}

func (v *Visitor) VisitGateOperandList(ctx *parser.GateOperandListContext) interface{} {
	var list []q.Qubit
	for _, o := range ctx.AllGateOperand() {
		list = append(list, v.Visit(o).([]q.Qubit)...)
	}

	return list
}

func (v *Visitor) VisitExternArgumentList(ctx *parser.ExternArgumentListContext) interface{} {
	var list []interface{}
	for _, arg := range ctx.AllExternArgument() {
		list = append(list, v.Visit(arg))
	}

	return list
}
