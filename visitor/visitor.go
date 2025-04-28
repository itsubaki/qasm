package visitor

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"slices"
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
	ErrIdentifierNotFound   = errors.New("identifier not found")
	ErrQubitNotFound        = errors.New("qubit not found")
	ErrClassicalBitNotFound = errors.New("classical bit not found")
	ErrVariableNotFound     = errors.New("variable not found")
	ErrGateNotFound         = errors.New("gate not found")
	ErrFunctionNotFound     = errors.New("function not found")
	ErrUnexpected           = errors.New("unexpected")
	ErrNotImplemented       = errors.New("not implemented")
)

type Visitor struct {
	qsim *q.Q
	env  *Environ
}

func New(qsim *q.Q, env *Environ) *Visitor {
	return &Visitor{
		qsim,
		env,
	}
}

func (v *Visitor) Enclosed() *Visitor {
	return New(v.qsim, v.env.NewEnclosed())
}

func (v *Visitor) Visit(tree antlr.ParseTree) any {
	return tree.Accept(v)
}

func (v *Visitor) VisitTerminal(node antlr.TerminalNode) any {
	return node.GetText()
}

func (v *Visitor) VisitErrorNode(node antlr.ErrorNode) any {
	return node.GetText()
}

func (v *Visitor) VisitPragma(ctx *parser.PragmaContext) any {
	return v.Visit(ctx.RemainingLineContent())
}

func (v *Visitor) VisitAnnotation(ctx *parser.AnnotationContext) any {
	return fmt.Errorf("VisitAnnotation: %w", ErrNotImplemented)
}

func (v *Visitor) VisitChildren(node antlr.RuleNode) any {
	return fmt.Errorf("VisitChildren: %w", ErrNotImplemented)
}

func (v *Visitor) VisitProgram(ctx *parser.ProgramContext) any {
	if ctx.Version() != nil {
		v.env.Version = v.Visit(ctx.Version()).(string)
	}

	for _, s := range ctx.AllStatementOrScope() {
		if res := v.Visit(s); res != nil {
			return res
		}
	}

	return nil
}

func (v *Visitor) VisitVersion(ctx *parser.VersionContext) any {
	return v.Visit(ctx.VersionSpecifier())
}

func (v *Visitor) VisitStatementOrScope(ctx *parser.StatementOrScopeContext) any {
	if ctx.Statement() != nil {
		return v.Visit(ctx.Statement())
	}

	return v.Visit(ctx.Scope())
}

func (v *Visitor) VisitStatement(ctx *parser.StatementContext) any {
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

func (v *Visitor) VisitScope(ctx *parser.ScopeContext) any {
	var list []any
	for _, s := range ctx.AllStatementOrScope() {
		result := v.Visit(s)
		list = append(list, result)

		if contains(result, Break, Continue) {
			return list
		}
	}

	return list
}

func (v *Visitor) VisitIncludeStatement(ctx *parser.IncludeStatementContext) any {
	path := strings.Trim(v.Visit(ctx.StringLiteral()).(string), "\"")
	text, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file=%s: %v", path, err)
	}

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(string(text)))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	switch ret := v.Visit(p.Program()).(type) {
	case error:
		return fmt.Errorf("include: %w", ret)
	}

	return nil
}

func (v *Visitor) VisitBreakStatement(ctx *parser.BreakStatementContext) any {
	return ctx.GetText()
}

func (v *Visitor) VisitContinueStatement(ctx *parser.ContinueStatementContext) any {
	return ctx.GetText()
}

func (v *Visitor) VisitEndStatement(ctx *parser.EndStatementContext) any {
	return ctx.GetText()
}

func (v *Visitor) VisitIfStatement(ctx *parser.IfStatementContext) any {
	enclosed := v.Enclosed()
	if v.Visit(ctx.Expression()).(bool) {
		return enclosed.Visit(ctx.GetIf_body())
	}

	if ctx.GetElse_body() != nil {
		return enclosed.Visit(ctx.GetElse_body())
	}

	return nil
}

func (v *Visitor) VisitForStatement(ctx *parser.ForStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	rx := v.Visit(ctx.RangeExpression()).([]int64)

	enclosed := v.Enclosed()
	for i := rx[0]; i < rx[1]; i++ {
		enclosed.env.Variable[id] = i
		result := enclosed.Visit(ctx.StatementOrScope())

		if contains(result, Break) {
			return nil
		}
	}

	return nil
}

func (v *Visitor) VisitWhileStatement(ctx *parser.WhileStatementContext) any {
	enclosed := v.Enclosed()
	for {
		if !v.Visit(ctx.Expression()).(bool) {
			return nil
		}

		result := enclosed.Visit(ctx.GetBody())
		if contains(result, Break) {
			return nil
		}
	}
}

func (v *Visitor) VisitSwitchStatement(ctx *parser.SwitchStatementContext) any {
	enclosed := v.Enclosed()

	x := v.Visit(ctx.Expression())
	for _, item := range ctx.AllSwitchCaseItem() {
		if item.DEFAULT() != nil {
			enclosed.Visit(item)
			return nil
		}

		result := v.Visit(item.ExpressionList()).([]any)
		for _, r := range result {
			if r != x {
				continue
			}

			enclosed.Visit(item)
			return nil
		}
	}

	return nil
}

func (v *Visitor) VisitSwitchCaseItem(ctx *parser.SwitchCaseItemContext) any {
	return v.Visit(ctx.Scope())
}

func (v *Visitor) VisitReturnStatement(ctx *parser.ReturnStatementContext) any {
	if ctx.MeasureExpression() != nil {
		return v.Visit(ctx.MeasureExpression())
	}

	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitGateStatement(ctx *parser.GateStatementContext) any {
	name := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetGate(name); ok {
		return fmt.Errorf("identifier=%s: %w", name, ErrAlreadyDeclared)
	}

	var params, qargs []string
	switch len(ctx.AllIdentifierList()) {
	case 1:
		qargs = v.Visit(ctx.IdentifierList(0)).([]string)
	case 2:
		params = v.Visit(ctx.IdentifierList(0)).([]string)
		qargs = v.Visit(ctx.IdentifierList(1)).([]string)
	default:
		return fmt.Errorf("len(identifier list)=%d: %w", len(ctx.AllIdentifierList()), ErrUnexpected)
	}

	var body []*parser.GateCallStatementContext
	for _, s := range ctx.Scope().AllStatementOrScope() {
		body = append(body, s.Statement().GateCallStatement().(*parser.GateCallStatementContext))
	}

	v.env.Gate[name] = Gate{
		Name:   name,
		Params: params,
		QArgs:  qargs,
		Body:   body,
	}

	return nil
}

func (v *Visitor) Params(xlist parser.IExpressionListContext) ([]float64, error) {
	var params []float64
	for _, e := range v.Visit(xlist).([]any) {
		switch val := e.(type) {
		case float64:
			params = append(params, val)
		case int64:
			params = append(params, float64(val))
		default:
			return nil, fmt.Errorf("param=%v(%T): %w", val, e, ErrUnexpected)
		}
	}

	return params, nil
}

func (v *Visitor) Builtin(ctx *parser.GateCallStatementContext) (matrix.Matrix, bool, error) {
	if ctx.GPHASE() != nil {
		params, err := v.Params(ctx.ExpressionList())
		if err != nil {
			return nil, false, fmt.Errorf("params: %w", err)
		}

		n := v.qsim.NumberOfBit()
		u := gate.I(n).Mul(cmplx.Exp(complex(0, params[0])))
		return u, true, nil
	}

	id := v.Visit(ctx.Identifier()).(string)
	switch id {
	case U:
		// params, qargs
		params, err := v.Params(ctx.ExpressionList())
		if err != nil {
			return nil, false, fmt.Errorf("params: %w", err)
		}
		qargs := v.Visit(ctx.GateOperandList()).([][]q.Qubit)

		// u
		n := v.qsim.NumberOfBit()
		u := gate.U(params[0], params[1], params[2])
		u = gate.TensorProduct(u, n, q.Index(flatten(qargs)...))

		// modify
		u, err = v.Modify(u, qargs, ctx.AllGateModifier())
		if err != nil {
			return nil, false, fmt.Errorf("modify: %w", err)
		}

		return u, true, nil
	default:
		return nil, false, nil
	}
}

func (v *Visitor) Modify(u matrix.Matrix, qargs [][]q.Qubit, modifier []parser.IGateModifierContext) (matrix.Matrix, error) {
	rev := make([]parser.IGateModifierContext, len(modifier))
	copy(rev, modifier)
	slices.Reverse(rev)
	for _, mod := range rev {
		switch {
		case mod.INV() != nil:
			u = u.Dagger()
		case mod.POW() != nil:
			var p float64
			switch n := v.Visit(mod).(type) {
			case float64:
				p = n
			case int64:
				p = float64(n)
			default:
				return nil, fmt.Errorf("pow=%v(%T): %w", n, n, ErrUnexpected)
			}

			u = Pow(u, p)
		}
	}

	var i int
	for _, mod := range modifier {
		switch {
		case mod.CTRL() != nil:
			u, i = Controlled(u, q.Index(qargs[i]...)), i+1
		case mod.NEGCTRL() != nil:
			u, i = NegControlled(u, q.Index(qargs[i]...)), i+1
		}
	}

	return u, nil
}

func (v *Visitor) Defined(ctx *parser.GateCallStatementContext) (matrix.Matrix, error) {
	id := v.Visit(ctx.Identifier()).(string)
	g, ok := v.env.GetGate(id)
	if !ok {
		return nil, fmt.Errorf("idenfitier=%s: %w", id, ErrGateNotFound)
	}

	// params, qargs
	enclosed := v.Enclosed()
	if ctx.ExpressionList() != nil {
		params, err := v.Params(ctx.ExpressionList())
		if err != nil {
			return nil, fmt.Errorf("params: %w", err)
		}

		for i, p := range g.Params {
			enclosed.env.Variable[p] = params[i]
		}
	}

	var qargs [][]q.Qubit
	if ctx.GateOperandList() != nil {
		var shift int
		for _, mod := range ctx.AllGateModifier() {
			switch {
			case
				mod.CTRL() != nil,
				mod.NEGCTRL() != nil:
				shift++
			}
		}

		qargs = v.Visit(ctx.GateOperandList()).([][]q.Qubit)
		for i, id := range g.QArgs {
			enclosed.env.Qubit[id] = qargs[i+shift]
		}
	}

	// u
	var list []matrix.Matrix
	for _, c := range g.Body {
		u, ok, err := enclosed.Builtin(c)
		if err != nil {
			return nil, fmt.Errorf("builtin: %w", err)
		}

		if !ok {
			u, err = enclosed.Defined(c)
			if err != nil {
				return nil, fmt.Errorf("defined: %w", err)
			}
		}

		list = append(list, u)
	}

	// matrix.Apply(A, B, C, ...) is ...CBA
	u := matrix.Apply(list...)

	// modify
	u, err := v.Modify(u, qargs, ctx.AllGateModifier())
	if err != nil {
		return nil, fmt.Errorf("modify: %w", err)
	}

	return u, nil
}

func (v *Visitor) VisitGateCallStatement(ctx *parser.GateCallStatementContext) any {
	u, ok, err := v.Builtin(ctx)
	if err != nil {
		return fmt.Errorf("builtin: %w", err)
	}

	if !ok {
		u, err = v.Defined(ctx)
		if err != nil {
			return fmt.Errorf("defined: %w", err)
		}
	}

	v.qsim.Apply(u)
	return nil
}

func (v *Visitor) MeasureAssignment(identifier parser.IIndexedIdentifierContext, measure parser.IMeasureExpressionContext) error {
	measured := v.Visit(measure)
	if identifier == nil {
		return nil
	}

	operand := v.Visit(identifier.Identifier()).(string)
	bits, ok := v.env.GetClassicalBit(operand)
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

func (v *Visitor) VisitMeasureArrowAssignmentStatement(ctx *parser.MeasureArrowAssignmentStatementContext) any {
	return v.MeasureAssignment(ctx.IndexedIdentifier(), ctx.MeasureExpression())
}

func (v *Visitor) VisitAssignmentStatement(ctx *parser.AssignmentStatementContext) any {
	if ctx.MeasureExpression() != nil {
		return v.MeasureAssignment(ctx.IndexedIdentifier(), ctx.MeasureExpression())
	}

	id := v.Visit(ctx.IndexedIdentifier().Identifier()).(string)
	v.env.SetVariable(id, v.Visit(ctx.Expression()))
	return nil
}

func (v *Visitor) VisitResetStatement(ctx *parser.ResetStatementContext) any {
	qargs := v.Visit(ctx.GateOperand()).([]q.Qubit)
	v.qsim.Reset(qargs...)
	return nil
}

func (v *Visitor) VisitConstDeclarationStatement(ctx *parser.ConstDeclarationStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetConst(id); ok {
		return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
	}

	v.env.Const[id] = v.Visit(ctx.DeclarationExpression())
	return nil
}

func (v *Visitor) VisitQuantumDeclarationStatement(ctx *parser.QuantumDeclarationStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetQubit(id); ok {
		return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
	}

	size := v.Visit(ctx.QubitType()).(int64)
	v.env.Qubit[id] = v.qsim.Zeros(int(size))
	return nil
}

func (v *Visitor) VisitClassicalDeclarationStatement(ctx *parser.ClassicalDeclarationStatementContext) any {
	switch {
	case ctx.ArrayType() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		size := v.Visit(ctx.ArrayType().ExpressionList()).([]any)[0].(int64)
		scalar := ctx.ArrayType().ScalarType()

		// make array
		switch {
		case scalar.INT() != nil:
			bits := v.Visit(scalar.Designator()).(int64)
			switch bits {
			case 8:
				v.env.Variable[id] = make([]int8, size)
			case 16:
				v.env.Variable[id] = make([]int16, size)
			case 32:
				v.env.Variable[id] = make([]int32, size)
			case 64:
				v.env.Variable[id] = make([]int64, size)
			default:
				return fmt.Errorf("invalid bit size=%d: %w", bits, ErrUnexpected)
			}

			return nil
		case scalar.UINT() != nil:
			bits := v.Visit(scalar.Designator()).(int64)
			switch bits {
			case 8:
				v.env.Variable[id] = make([]uint8, size)
			case 16:
				v.env.Variable[id] = make([]uint16, size)
			case 32:
				v.env.Variable[id] = make([]uint32, size)
			case 64:
				v.env.Variable[id] = make([]uint64, size)
			default:
				return fmt.Errorf("invalid bit size=%d: %w", bits, ErrUnexpected)
			}

			return nil
		case scalar.FLOAT() != nil:
			bits := v.Visit(scalar.Designator()).(int64)
			switch bits {
			case 32:
				v.env.Variable[id] = make([]float32, size)
			case 64:
				v.env.Variable[id] = make([]float64, size)
			default:
				return fmt.Errorf("invalid bit size=%d: %w", bits, ErrUnexpected)
			}

			return nil
		case scalar.BOOL() != nil:
			v.env.Variable[id] = make([]bool, size)
			return nil
		default:
			return fmt.Errorf("scalar type=%s: %w", scalar.GetText(), ErrUnexpected)
		}
	case ctx.ScalarType().INT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.Variable[id] = v.Visit(ctx.DeclarationExpression())
			return nil
		}

		v.env.Variable[id] = int(0)
		return nil
	case ctx.ScalarType().UINT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.Variable[id] = v.Visit(ctx.DeclarationExpression())
			return nil
		}

		v.env.Variable[id] = uint(0)
		return nil
	case ctx.ScalarType().FLOAT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.Variable[id] = v.Visit(ctx.DeclarationExpression())
			return nil
		}

		v.env.Variable[id] = float64(0)
		return nil
	case ctx.ScalarType().BOOL() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.Variable[id] = v.Visit(ctx.DeclarationExpression())
			return nil
		}

		v.env.Variable[id] = false
		return nil
	case ctx.ScalarType().BIT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetClassicalBit(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			bits := v.Visit(ctx.DeclarationExpression()).([]int64)
			v.env.ClassicalBit[id] = bits
			return nil
		}

		size := v.Visit(ctx.ScalarType()).(int64)
		v.env.ClassicalBit[id] = make([]int64, int(size))
		return nil
	default:
		return fmt.Errorf("scalar type=%s: %w", ctx.ScalarType().GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitDefStatement(ctx *parser.DefStatementContext) any {
	name := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetSubroutine(name); ok {
		return fmt.Errorf("identifier=%s: %w", name, ErrAlreadyDeclared)
	}

	args := v.Visit(ctx.ArgumentDefinitionList()).([]any)
	var qargs []string
	for _, a := range args {
		qargs = append(qargs, a.(string))
	}

	v.env.Subroutine[name] = Subroutine{
		Name:            name,
		QArgs:           qargs,
		Body:            ctx.Scope().(*parser.ScopeContext),
		ReturnSignature: ctx.ReturnSignature().(*parser.ReturnSignatureContext),
	}

	return nil
}

func (v *Visitor) VisitAliasDeclarationStatement(ctx *parser.AliasDeclarationStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetQubit(id); ok {
		return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
	}

	alias := v.Visit(ctx.AliasExpression()).([]q.Qubit)
	v.env.Qubit[id] = alias
	return nil
}

func (v *Visitor) VisitOldStyleDeclarationStatement(ctx *parser.OldStyleDeclarationStatementContext) any {
	switch {
	case ctx.QREG() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetQubit(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		var size int64 = 1
		if ctx.Designator() != nil {
			size = v.Visit(ctx.Designator()).(int64)
		}

		v.env.Qubit[id] = v.qsim.Zeros(int(size))
		return nil
	case ctx.CREG() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetClassicalBit(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		var size int64 = 1
		if ctx.Designator() != nil {
			size = v.Visit(ctx.Designator()).(int64)
		}

		v.env.ClassicalBit[id] = make([]int64, int(size))
		return nil
	default:
		return fmt.Errorf("x=%s: %w", ctx.GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitIoDeclarationStatement(ctx *parser.IoDeclarationStatementContext) any {
	return fmt.Errorf("VisitIoDeclarationStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitExternStatement(ctx *parser.ExternStatementContext) any {
	return fmt.Errorf("VisitExternStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitCalStatement(ctx *parser.CalStatementContext) any {
	return fmt.Errorf("VisitCalStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalStatement(ctx *parser.DefcalStatementContext) any {
	return fmt.Errorf("VisitDefcalStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitCalibrationGrammarStatement(ctx *parser.CalibrationGrammarStatementContext) any {
	return fmt.Errorf("VisitCalibrationGrammarStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBarrierStatement(ctx *parser.BarrierStatementContext) any {
	return fmt.Errorf("VisitBarrierStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBoxStatement(ctx *parser.BoxStatementContext) any {
	return fmt.Errorf("VisitBoxStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDelayStatement(ctx *parser.DelayStatementContext) any {
	return fmt.Errorf("VisitDelayStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) any {
	switch {
	case ctx.Identifier() != nil:
		s := v.Visit(ctx.Identifier()).(string)
		if lit, ok := Const[s]; ok {
			return lit
		}

		if lit, ok := v.env.GetConst(s); ok {
			return lit
		}

		if lit, ok := v.env.GetVariable(s); ok {
			return lit
		}

		if lit, ok := v.env.GetQubit(s); ok {
			return lit
		}

		if lit, ok := v.env.GetClassicalBit(s); ok {
			return lit
		}

		return fmt.Errorf("identifier=%s: %w", s, ErrIdentifierNotFound)
	case ctx.DecimalIntegerLiteral() != nil:
		s := v.Visit(ctx.DecimalIntegerLiteral()).(string)
		lit, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int: s=%s: %w", s, err)
		}

		return lit
	case ctx.HexIntegerLiteral() != nil:
		s := v.Visit(ctx.HexIntegerLiteral()).(string)
		lit, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return fmt.Errorf("parse int: s=%s: %w", s, err)
		}

		return lit
	case ctx.FloatLiteral() != nil:
		s := v.Visit(ctx.FloatLiteral()).(string)
		lit, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("parse float: s=%s: %w", s, err)
		}

		return lit
	case ctx.BooleanLiteral() != nil:
		s := v.Visit(ctx.BooleanLiteral()).(string)
		lit, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("parse bool: s=%s: %w", s, err)
		}

		return lit
	case ctx.BitstringLiteral() != nil:
		s := v.Visit(ctx.BitstringLiteral()).(string)
		bistring := strings.Trim(s, "\"")

		lit := make([]int64, len(bistring))
		for i, b := range bistring {
			lit[i] = int64(b - '0')
		}

		return lit
	default:
		return fmt.Errorf("x=%s: %w", ctx.GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitAdditiveExpression(ctx *parser.AdditiveExpressionContext) any {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		default:
			return fmt.Errorf("operand=%v(%T): %w", val, val, ErrUnexpected)
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

func (v *Visitor) VisitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) any {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		default:
			return fmt.Errorf("operand=%v(%T): %w", val, val, ErrUnexpected)
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

func (v *Visitor) VisitEqualityExpression(ctx *parser.EqualityExpressionContext) any {
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
			return fmt.Errorf("operand=%v(%T): %w", val, val, ErrUnexpected)
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

func (v *Visitor) VisitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) any {
	result := true
	for _, x := range ctx.AllExpression() {
		if result = result && v.Visit(x).(bool); !result {
			return false
		}
	}

	return true
}

func (v *Visitor) VisitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) any {
	var result bool
	for _, x := range ctx.AllExpression() {
		if result = result || v.Visit(x).(bool); result {
			return true
		}
	}

	return false
}

func (v *Visitor) VisitBitwiseAndExpression(ctx *parser.BitwiseAndExpressionContext) any {
	result := v.Visit(ctx.AllExpression()[0]).(int64)
	for _, x := range ctx.AllExpression()[1:] {
		result = result & v.Visit(x).(int64)
	}

	return result
}

func (v *Visitor) VisitBitwiseOrExpression(ctx *parser.BitwiseOrExpressionContext) any {
	result := v.Visit(ctx.AllExpression()[0]).(int64)
	for _, x := range ctx.AllExpression()[1:] {
		result = result | v.Visit(x).(int64)
	}

	return result
}

func (v *Visitor) VisitBitwiseXorExpression(ctx *parser.BitwiseXorExpressionContext) any {
	result := v.Visit(ctx.AllExpression()[0]).(int64)
	for _, x := range ctx.AllExpression()[1:] {
		result = result ^ v.Visit(x).(int64)
	}

	return result
}

func (v *Visitor) VisitBitshiftExpression(ctx *parser.BitshiftExpressionContext) any {
	op := v.Visit(ctx.BitshiftOperator()).(string)
	result := v.Visit(ctx.AllExpression()[0]).(int64)
	for _, x := range ctx.AllExpression()[1:] {
		switch op {
		case "<<":
			result = result << v.Visit(x).(int64)
		case ">>":
			result = result >> v.Visit(x).(int64)
		}
	}

	return result
}

func (v *Visitor) VisitComparisonExpression(ctx *parser.ComparisonExpressionContext) any {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		default:
			return fmt.Errorf("operand=%v(%T): %w", val, val, ErrUnexpected)
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

func (v *Visitor) VisitPowerExpression(ctx *parser.PowerExpressionContext) any {
	var operand []float64
	for _, x := range ctx.AllExpression() {
		switch val := v.Visit(x).(type) {
		case float64:
			operand = append(operand, val)
		case int64:
			operand = append(operand, float64(val))
		default:
			return fmt.Errorf("operand=%v(%T): %w", val, val, ErrUnexpected)
		}
	}

	return math.Pow(operand[0], operand[1])
}

func (v *Visitor) VisitUnaryExpression(ctx *parser.UnaryExpressionContext) any {
	operand := v.Visit(ctx.Expression())
	switch {
	case ctx.MINUS() != nil:
		switch val := operand.(type) {
		case float64:
			return -1 * val
		case int64:
			return -1 * val
		default:
			return fmt.Errorf("operand=%v(%T): %w", val, val, ErrUnexpected)
		}
	case ctx.EXCLAMATION_POINT() != nil:
		switch val := operand.(type) {
		case bool:
			return !val
		default:
			return fmt.Errorf("operand=%v(%T): %w", val, val, ErrUnexpected)
		}
	case ctx.TILDE() != nil:
		switch val := operand.(type) {
		case int64:
			return ^val
		default:
			return fmt.Errorf("operand=%v(%T): %w", val, val, ErrUnexpected)
		}
	default:
		return fmt.Errorf("operator=%s: %w", ctx.GetOp().GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitDeclarationExpression(ctx *parser.DeclarationExpressionContext) any {
	if ctx.MeasureExpression() != nil {
		return v.Visit(ctx.MeasureExpression())
	}

	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitCallExpression(ctx *parser.CallExpressionContext) any {
	args := v.Visit(ctx.ExpressionList()).([]any)

	id := v.Visit(ctx.Identifier()).(string)
	switch id {
	case "sin":
		return math.Sin(args[0].(float64))
	case "cos":
		return math.Cos(args[0].(float64))
	case "tan":
		return math.Tan(args[0].(float64))
	case "arcsin":
		return math.Asin(args[0].(float64))
	case "arccos":
		return math.Acos(args[0].(float64))
	case "arctan":
		return math.Atan(args[0].(float64))
	case "ceiling":
		return math.Ceil(args[0].(float64))
	case "floor":
		return math.Floor(args[0].(float64))
	case "sqrt":
		return math.Sqrt(args[0].(float64))
	case "exp":
		return math.Exp(args[0].(float64))
	case "log":
		return math.Log(args[0].(float64))
	case "mod":
		return math.Mod(args[0].(float64), args[1].(float64))
	default:
		sub, ok := v.env.GetSubroutine(id)
		if !ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrFunctionNotFound)
		}

		enclosed := v.Enclosed()
		for i, p := range sub.QArgs {
			enclosed.env.Qubit[p] = args[i].([]q.Qubit)
		}

		result := enclosed.Visit(sub.Body).([]any)
		return result[len(result)-1]
	}
}

func (v *Visitor) VisitRangeExpression(ctx *parser.RangeExpressionContext) any {
	var list []int64
	for _, x := range ctx.AllExpression() {
		list = append(list, v.Visit(x).(int64))
	}

	return list
}

func (v *Visitor) VisitAliasExpression(ctx *parser.AliasExpressionContext) any {
	var result []q.Qubit
	for _, x := range ctx.AllExpression() {
		result = append(result, v.Visit(x).([]q.Qubit)...)
	}

	return result
}

func (v *Visitor) VisitIndexExpression(ctx *parser.IndexExpressionContext) any {
	qubit := v.Visit(ctx.Expression()).([]q.Qubit)

	var result []q.Qubit
	for _, x := range ctx.IndexOperator().AllRangeExpression() {
		index := v.Visit(x).([]int64)
		result = append(result, qubit[index[0]:index[1]]...)
	}

	return result
}

func (v *Visitor) VisitMeasureExpression(ctx *parser.MeasureExpressionContext) any {
	qargs := v.Visit(ctx.GateOperand()).([]q.Qubit)
	v.qsim.Measure(qargs...)

	var bits []int64
	for _, q := range qargs {
		bits = append(bits, v.qsim.State(q)[0].Int(0))
	}

	return bits
}

func (v *Visitor) VisitCastExpression(ctx *parser.CastExpressionContext) any {
	val := v.Visit(ctx.Expression())
	switch {
	case ctx.ScalarType().INT() != nil:
		switch v := val.(type) {
		case float64:
			return int(v)
		case int64:
			return int(v)
		}
	case ctx.ScalarType().UINT() != nil:
		switch v := val.(type) {
		case float64:
			return uint(v)
		case int64:
			return uint(v)
		}
	}

	return fmt.Errorf("x=%s: %w", ctx.GetText(), ErrUnexpected)
}

func (v *Visitor) VisitParenthesisExpression(ctx *parser.ParenthesisExpressionContext) any {
	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitExpressionStatement(ctx *parser.ExpressionStatementContext) any {
	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitDesignator(ctx *parser.DesignatorContext) any {
	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitDurationofExpression(ctx *parser.DurationofExpressionContext) any {
	return fmt.Errorf("VisitDurationofExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitSetExpression(ctx *parser.SetExpressionContext) any {
	return fmt.Errorf("VisitSetExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitIndexOperator(ctx *parser.IndexOperatorContext) any {
	var list []any
	for _, x := range ctx.AllExpression() {
		list = append(list, v.Visit(x))
	}

	return list
}

func (v *Visitor) VisitIndexedIdentifier(ctx *parser.IndexedIdentifierContext) any {
	var index []int64
	for _, op := range ctx.AllIndexOperator() {
		for _, v := range v.Visit(op).([]any) {
			index = append(index, v.(int64))
		}
	}

	return index
}

func (v *Visitor) VisitGateModifier(ctx *parser.GateModifierContext) any {
	if ctx.Expression() != nil {
		return v.Visit(ctx.Expression())
	}

	return int64(1)
}

func (v *Visitor) VisitScalarType(ctx *parser.ScalarTypeContext) any {
	if ctx.Designator() != nil {
		return v.Visit(ctx.Designator()).(int64)
	}

	return int64(1)
}

func (v *Visitor) VisitQubitType(ctx *parser.QubitTypeContext) any {
	if ctx.Designator() != nil {
		return v.Visit(ctx.Designator()).(int64)
	}

	return int64(1)
}

func (v *Visitor) VisitArrayType(ctx *parser.ArrayTypeContext) any {
	return fmt.Errorf("VisitArrayType: %w", ErrNotImplemented)
}

func (v *Visitor) VisitArrayReferenceType(ctx *parser.ArrayReferenceTypeContext) any {
	return fmt.Errorf("VisitArrayReferenceType: %w", ErrNotImplemented)
}

func (v *Visitor) VisitArrayLiteral(ctx *parser.ArrayLiteralContext) any {
	return fmt.Errorf("VisitArrayLiteral: %w", ErrNotImplemented)
}

func (v *Visitor) VisitReturnSignature(ctx *parser.ReturnSignatureContext) any {
	return fmt.Errorf("VisitReturnSignature: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalArgumentDefinition(ctx *parser.DefcalArgumentDefinitionContext) any {
	return fmt.Errorf("VisitDefcalArgumentDefinition: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalTarget(ctx *parser.DefcalTargetContext) any {
	return fmt.Errorf("VisitDefcalTarget: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalOperand(ctx *parser.DefcalOperandContext) any {
	return fmt.Errorf("VisitDefcalOperand: %w", ErrNotImplemented)
}

func (v *Visitor) VisitExternArgument(ctx *parser.ExternArgumentContext) any {
	return fmt.Errorf("VisitExternArgument: %w", ErrNotImplemented)
}

func (v *Visitor) VisitArgumentDefinition(ctx *parser.ArgumentDefinitionContext) any {
	return v.Visit(ctx.Identifier())
}

func (v *Visitor) VisitGateOperand(ctx *parser.GateOperandContext) any {
	indexID := ctx.IndexedIdentifier()

	operand := v.Visit(indexID.Identifier()).(string)
	qb, ok := v.env.GetQubit(operand)
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

func (v *Visitor) VisitArgumentDefinitionList(ctx *parser.ArgumentDefinitionListContext) any {
	var list []any
	for _, def := range ctx.AllArgumentDefinition() {
		list = append(list, v.Visit(def))
	}

	return list
}

func (v *Visitor) VisitExpressionList(ctx *parser.ExpressionListContext) any {
	var list []any
	for _, x := range ctx.AllExpression() {
		list = append(list, v.Visit(x))
	}

	return list
}

func (v *Visitor) VisitIdentifierList(ctx *parser.IdentifierListContext) any {
	var list []string
	for _, id := range ctx.AllIdentifier() {
		list = append(list, v.Visit(id).(string))
	}

	return list
}

func (v *Visitor) VisitGateOperandList(ctx *parser.GateOperandListContext) any {
	var list [][]q.Qubit
	for _, o := range ctx.AllGateOperand() {
		list = append(list, v.Visit(o).([]q.Qubit))
	}

	return list
}

func (v *Visitor) VisitDefcalArgumentDefinitionList(ctx *parser.DefcalArgumentDefinitionListContext) any {
	var list []any
	for _, def := range ctx.AllDefcalArgumentDefinition() {
		list = append(list, v.Visit(def))
	}

	return list
}

func (v *Visitor) VisitDefcalOperandList(ctx *parser.DefcalOperandListContext) any {
	var list []any
	for _, o := range ctx.AllDefcalOperand() {
		list = append(list, v.Visit(o))
	}

	return list
}

func (v *Visitor) VisitExternArgumentList(ctx *parser.ExternArgumentListContext) any {
	var list []any
	for _, arg := range ctx.AllExternArgument() {
		list = append(list, v.Visit(arg))
	}

	return list
}
