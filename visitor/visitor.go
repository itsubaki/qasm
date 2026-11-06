package visitor

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/qasm/angle"
	"github.com/itsubaki/qasm/environ"
	"github.com/itsubaki/qasm/gen/parser"
	xparser "github.com/itsubaki/qasm/parser"
	"github.com/itsubaki/qasm/value"
)

var (
	ErrTooManyQubits  = errors.New("too many qubits")
	ErrNotImplemented = errors.New("not implemented")
)

type Visitor struct {
	*parser.Baseqasm3ParserVisitor
	qsim      *q.Q
	env       *environ.Environ
	maxQubits int
}

type returnValue struct {
	value any
}

type endSignal struct{}

func New(qsim *q.Q, env *environ.Environ, opt ...Option) *Visitor {
	v := &Visitor{
		Baseqasm3ParserVisitor: &parser.Baseqasm3ParserVisitor{},
		qsim:                   qsim,
		env:                    env,
	}

	for _, f := range opt {
		f(v)
	}

	return v
}

func Run(text string) (*q.Q, *environ.Environ, error) {
	program, err := xparser.Parse(text)
	if err != nil {
		return nil, nil, err
	}

	qsim := q.New()
	env := environ.New()
	v := New(qsim, env)
	if err := v.Run(program); err != nil {
		return nil, nil, err
	}

	return qsim, env, nil
}

func Visit(text string) (any, error) {
	program, err := xparser.Parse(text)
	if err != nil {
		return nil, err
	}

	result := New(q.New(), environ.New()).Visit(program)
	if err, ok := result.(error); ok && err != nil {
		return nil, err
	}

	if _, ok := result.(endSignal); ok {
		return End, nil
	}

	return result, nil
}

func (v *Visitor) Enclosed() *Visitor {
	return New(v.qsim, v.env.NewEnclosed(), WithMaxQubits(v.maxQubits))
}

func (v *Visitor) Run(tree antlr.ParseTree) error {
	if err, ok := v.Visit(tree).(error); ok && err != nil {
		return err
	}

	return nil
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

func (v *Visitor) VisitProgram(ctx *parser.ProgramContext) any {
	if ctx.Version() != nil {
		v.env.Version = v.Visit(ctx.Version()).(string)
	}

	var result any
	for _, s := range ctx.AllStatementOrScope() {
		result = v.Visit(s)
		switch result.(type) {
		case error, returnValue, endSignal:
			return result
		}

		if contains(result, Break, Continue) {
			return result
		}
	}

	return result
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

	return fmt.Errorf("unsupported statement %q", ctx.GetText())
}

func (v *Visitor) VisitScope(ctx *parser.ScopeContext) any {
	enclosed := v.Enclosed()

	for _, s := range ctx.AllStatementOrScope() {
		result := enclosed.Visit(s)
		switch result.(type) {
		case error, returnValue, endSignal:
			return result
		}

		if contains(result, Break, Continue) {
			return result
		}
	}

	return nil
}

func (v *Visitor) VisitIncludeStatement(ctx *parser.IncludeStatementContext) any {
	path := strings.Trim(v.Visit(ctx.StringLiteral()).(string), "\"")
	text, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %v", path, err)
	}

	program, err := xparser.Parse(string(text))
	if err != nil {
		return fmt.Errorf("include %s: %w", path, err)
	}

	result := v.Visit(program)
	if err, ok := result.(error); ok && err != nil {
		return fmt.Errorf("include %s: %w", path, err)
	}

	if _, ok := result.(endSignal); ok {
		return result
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
	return endSignal{}
}

func (v *Visitor) VisitIfStatement(ctx *parser.IfStatementContext) any {
	condition := v.Visit(ctx.Expression())
	switch condition.(type) {
	case error, endSignal:
		return condition
	}

	enclosed := v.Enclosed()
	if condition.(bool) {
		return enclosed.Visit(ctx.GetIf_body())
	}

	if ctx.GetElse_body() != nil {
		return enclosed.Visit(ctx.GetElse_body())
	}

	return nil
}

func (v *Visitor) VisitForStatement(ctx *parser.ForStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	result := v.Visit(ctx.RangeExpression())
	switch result.(type) {
	case error, endSignal:
		return result
	}

	rx := result.([]int64)
	enclosed := v.Enclosed()

	for i := rx[0]; i <= rx[1]; i++ {
		enclosed.env.Variable[id] = i
		result := enclosed.Visit(ctx.StatementOrScope())
		switch result.(type) {
		case error, returnValue, endSignal:
			return result
		}

		if contains(result, Break) {
			return nil
		}
	}

	return nil
}

func (v *Visitor) VisitWhileStatement(ctx *parser.WhileStatementContext) any {
	enclosed := v.Enclosed()
	for {
		condition := v.Visit(ctx.Expression())
		switch condition.(type) {
		case error, endSignal:
			return condition
		}

		if !condition.(bool) {
			return nil
		}

		result := enclosed.Visit(ctx.GetBody())
		switch result.(type) {
		case error, returnValue, endSignal:
			return result
		}

		if contains(result, Break) {
			return nil
		}
	}
}

func (v *Visitor) VisitSwitchStatement(ctx *parser.SwitchStatementContext) any {
	enclosed := v.Enclosed()

	x := v.Visit(ctx.Expression())
	switch x.(type) {
	case error, endSignal:
		return x
	}

	for _, item := range ctx.AllSwitchCaseItem() {
		if item.DEFAULT() != nil {
			return enclosed.Visit(item)
		}

		result := v.Visit(item.ExpressionList())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		for _, r := range result.([]any) {
			if r != x {
				continue
			}

			return enclosed.Visit(item)
		}
	}

	return nil
}

func (v *Visitor) VisitSwitchCaseItem(ctx *parser.SwitchCaseItemContext) any {
	return v.Visit(ctx.Scope())
}

func (v *Visitor) VisitReturnStatement(ctx *parser.ReturnStatementContext) any {
	var result any
	if ctx.MeasureExpression() != nil {
		result = v.Visit(ctx.MeasureExpression())
	} else if ctx.Expression() != nil {
		result = v.Visit(ctx.Expression())
	}

	switch result.(type) {
	case error, endSignal:
		return result
	}

	return returnValue{value: result}
}

func (v *Visitor) VisitGateStatement(ctx *parser.GateStatementContext) any {
	name := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetGate(name); ok {
		return fmt.Errorf("%q redeclared", name)
	}

	var params, qargs []string
	switch len(ctx.AllIdentifierList()) {
	case 1:
		qargs = v.Visit(ctx.IdentifierList(0)).([]string)
	case 2:
		params = v.Visit(ctx.IdentifierList(0)).([]string)
		qargs = v.Visit(ctx.IdentifierList(1)).([]string)
	}

	v.env.Gate[name] = &environ.Gate{
		Name:   name,
		Params: params,
		QArgs:  qargs,
		Body:   ctx.Scope(),
	}

	return nil
}

// Params returns []float64, an error, or an endSignal.
func (v *Visitor) Params(xlist parser.IExpressionListContext) any {
	result := v.Visit(xlist)
	switch result.(type) {
	case error, endSignal:
		return result
	}

	var params []float64
	for _, e := range result.([]any) {
		switch e.(type) {
		case error, endSignal:
			return e
		}

		switch param := e.(type) {
		case float64:
			params = append(params, param)
		case int64:
			params = append(params, float64(param))
		default:
			return fmt.Errorf("invalid param '%T(%v)'", param, param)
		}
	}

	return params
}

// Builtin returns *matrix.Matrix, an error, or an endSignal.
// It returns nil if the gate is not built-in.
func (v *Visitor) Builtin(ctx *parser.GateCallStatementContext) any {
	if ctx.GPHASE() != nil {
		result := v.Params(ctx.ExpressionList())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		params := result.([]float64)
		// u 2x2
		u := gate.I().Mul(cmplx.Exp(complex(0, params[0])))
		return u
	}

	id := v.Visit(ctx.Identifier()).(string)
	switch id {
	case U:
		// params
		result := v.Params(ctx.ExpressionList())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		params := result.([]float64)
		// u 2x2
		u := gate.U(params[0], params[1], params[2])
		return u
	default:
		return nil
	}
}

func (v *Visitor) UserDefinedGateCall(ctx *parser.GateCallStatementContext) any {
	if len(ctx.AllGateModifier()) > 0 {
		// NOTE: Modifiers on user-defined gate calls are not implemented yet.
		return fmt.Errorf("user-defined call with modifier: %w", ErrNotImplemented)
	}

	id := v.Visit(ctx.Identifier()).(string)
	g, ok := v.env.GetGate(id)
	if !ok {
		return fmt.Errorf("undefined %q", id)
	}

	// params
	enclosed := v.Enclosed()
	if ctx.ExpressionList() != nil {
		result := v.Params(ctx.ExpressionList())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		params := result.([]float64)
		for i, p := range g.Params {
			enclosed.env.Variable[p] = params[i]
		}
	}

	// Broadcast at the call site so gate parameters are bound to scalar qubits.
	var qargs []QubitOperand
	if ctx.GateOperandList() != nil {
		result := v.Visit(ctx.GateOperandList())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		qargs = result.([]QubitOperand)
	}

	if len(qargs) != len(g.QArgs) {
		return fmt.Errorf("expected %d gate operands, got %d", len(g.QArgs), len(qargs))
	}

	rows, err := broadcast(qargs)
	if err != nil {
		return err
	}

	for _, row := range rows {
		for i, id := range g.QArgs {
			enclosed.env.Qubit[id] = row[i]
		}

		for i, s := range g.Body.AllStatementOrScope() {
			call := s.Statement().GateCallStatement().(*parser.GateCallStatementContext)
			result := enclosed.VisitGateCallStatement(call)
			if err, ok := result.(error); ok && err != nil {
				return fmt.Errorf("gate call[%d]: %w", i, err)
			}
			if _, ok := result.(endSignal); ok {
				return result
			}
		}
	}

	return nil
}

func (v *Visitor) VisitGateCallStatement(ctx *parser.GateCallStatementContext) any {
	// NOTE: The "pow@ ctrl@" and "pow@ negctrl@" modifier orders are not supported.
	if err := NotImplementedOrder(ctx); err != nil {
		return err
	}

	// builtin gate
	result := v.Builtin(ctx)
	switch result.(type) {
	case error, endSignal:
		return result
	}

	// If no built-in gate was found, execute a user-defined gate.
	if result == nil {
		return v.UserDefinedGateCall(ctx)
	}

	u := result.(*matrix.Matrix)

	// inv and pow modifiers
	for _, mod := range ctx.AllGateModifier() {
		switch {
		case mod.INV() != nil:
			u = u.Dagger()
		case mod.POW() != nil:
			result := v.Visit(mod)
			switch result.(type) {
			case error, endSignal:
				return result
			}

			p, err := value.New(result).Float64()
			if err != nil {
				return fmt.Errorf("apply %q: %w", mod.GetText(), err)
			}

			// u is builtin gate, so Pow2x2 is used.
			u = Pow2x2(u, p.Value().(float64))
		}
	}

	// control modifiers
	if HasControlModifier(ctx) {
		// qubit[2] c;
		// qubit t;
		// U(pi/2, 0, pi) c;
		// U(pi/2, 0, pi) c[0], c[1];
		// ctrl @ U(pi, 0, pi) c, t;
		// ctrl @ U(pi, 0, pi) c[0], t;
		result := v.Visit(ctx.GateOperandList())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		qargs := result.([]QubitOperand)
		rows, err := broadcast(qargs)
		if err != nil {
			return err
		}

		var cursor int
		var negative []int
		for _, mod := range ctx.AllGateModifier() {
			if mod.CTRL() == nil && mod.NEGCTRL() == nil {
				continue
			}

			result := v.Visit(mod)
			switch result.(type) {
			case error, endSignal:
				return result
			}

			n, ok := result.(int64)
			if !ok || n < 1 {
				return fmt.Errorf("apply %q", mod.GetText())
			}

			if n > int64(len(qargs)-cursor-1) {
				return fmt.Errorf("not enough operands for %q", mod.GetText())
			}

			for range n {
				if mod.NEGCTRL() != nil {
					negative = append(negative, cursor)
				}

				cursor++
			}
		}

		if len(qargs) != cursor+1 {
			return fmt.Errorf("expected %d gate operands, got %d", cursor+1, len(qargs))
		}

		for _, row := range rows {
			neg := make([]q.Qubit, len(negative))
			for i, idx := range negative {
				neg[i] = row[idx]
			}

			v.qsim.X(neg...)
			v.qsim.Controlled(u, row[:cursor], row[cursor:])
			v.qsim.X(neg...)
		}

		return nil
	}

	// qargs
	var qargs []q.Qubit
	if ctx.GateOperandList() != nil {
		// qubit q0; qubit q1; U q0, q1;
		// qubit[2] q; U q;
		result := v.Visit(ctx.GateOperandList())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		operand := result.([]QubitOperand)
		for _, o := range operand {
			qargs = append(qargs, o.Qubits...)
		}
	} else {
		// all qubits for gphase
		for i := range v.qsim.NumQubits() {
			qargs = append(qargs, q.Qubit(i))
		}
	}

	v.qsim.G(u, qargs...)
	return nil
}

func (v *Visitor) MeasureAssignment(id parser.IIndexedIdentifierContext, measure parser.IMeasureExpressionContext) any {
	measured := v.Visit(measure)
	switch measured.(type) {
	case error, endSignal:
		return measured
	}

	if id == nil {
		return nil
	}

	operand := v.Visit(id.Identifier()).(string)
	result := v.Visit(id)
	switch result.(type) {
	case error, endSignal:
		return result
	}

	index := result.([]int64)
	if bits, ok := v.env.GetClBitArray(operand); ok {
		var val []bool
		switch m := measured.(type) {
		case bool:
			val = []bool{m}
		case []bool:
			val = m
		}

		if len(index) == 0 {
			v.env.SetClBitArray(operand, val)
			return nil
		}

		for i, bit := range val {
			bits[index[i]] = bit
		}

		return nil
	}

	if _, ok := v.env.GetClBit(operand); ok {
		switch bit := measured.(type) {
		case bool:
			v.env.SetClBit(operand, bit)
			return nil
		case []bool:
			return fmt.Errorf("assign %d bits to a single bit", len(bit))
		}
	}

	return fmt.Errorf("undefined %q", operand)
}

func (v *Visitor) VisitMeasureArrowAssignmentStatement(ctx *parser.MeasureArrowAssignmentStatementContext) any {
	return v.MeasureAssignment(ctx.IndexedIdentifier(), ctx.MeasureExpression())
}

func (v *Visitor) VisitAssignmentStatement(ctx *parser.AssignmentStatementContext) any {
	if ctx.MeasureExpression() != nil {
		return v.MeasureAssignment(ctx.IndexedIdentifier(), ctx.MeasureExpression())
	}

	operand := v.Visit(ctx.IndexedIdentifier().Identifier()).(string)
	result := v.Visit(ctx.IndexedIdentifier())
	switch result.(type) {
	case error, endSignal:
		return result
	}

	index := result.([]int64)
	x := v.Visit(ctx.Expression())
	switch x.(type) {
	case error, endSignal:
		return x
	}

	if bits, ok := v.env.GetClBitArray(operand); ok {
		var val []bool
		switch m := x.(type) {
		case bool:
			val = []bool{m}
		case []bool:
			val = m
		}

		if len(index) == 0 {
			v.env.SetClBitArray(operand, val)
			return nil
		}

		for i, bit := range val {
			bits[index[i]] = bit
		}

		return nil
	}

	if _, ok := v.env.GetClBit(operand); ok {
		switch bit := x.(type) {
		case bool:
			v.env.SetClBit(operand, bit)
			return nil
		case []bool:
			if len(bit) != 1 {
				return fmt.Errorf("assign %d bits to a single bit", len(bit))
			}

			v.env.SetClBit(operand, bit[0])
			return nil
		}
	}

	v.env.SetVariable(operand, x)
	return nil
}

func (v *Visitor) VisitResetStatement(ctx *parser.ResetStatementContext) any {
	result := v.Visit(ctx.GateOperand())
	switch result.(type) {
	case error, endSignal:
		return result
	}

	v.qsim.Reset(result.(QubitOperand).Qubits...)
	return nil
}

func (v *Visitor) VisitConstDeclarationStatement(ctx *parser.ConstDeclarationStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetConst(id); ok {
		return fmt.Errorf("%q redeclared", id)
	}

	result := v.Visit(ctx.DeclarationExpression())
	switch result.(type) {
	case error, endSignal:
		return result
	}

	v.env.Const[id] = result
	return nil
}

func (v *Visitor) VisitQuantumDeclarationStatement(ctx *parser.QuantumDeclarationStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetQubit(id); ok {
		return fmt.Errorf("%q redeclared", id)
	}

	if _, ok := v.env.GetQubitArray(id); ok {
		return fmt.Errorf("%q redeclared", id)
	}

	result := v.Visit(ctx.QubitType())
	switch result.(type) {
	case error, endSignal:
		return result
	}

	size, ok := result.(int64)
	if !ok {
		return fmt.Errorf("size must be an integer %q", ctx.QubitType().GetText())
	}

	need := v.qsim.NumQubits() + int(size)
	if v.maxQubits > 0 && need > v.maxQubits {
		return fmt.Errorf("need=%d, max=%d: %w", need, v.maxQubits, ErrTooManyQubits)
	}

	qubits := v.qsim.Zeros(int(size))
	if ctx.QubitType().Designator() != nil {
		v.env.SetQubitArray(id, qubits)
		return nil
	}

	v.env.SetQubit(id, qubits[0])
	return nil
}

func (v *Visitor) VisitClassicalDeclarationStatement(ctx *parser.ClassicalDeclarationStatementContext) any {
	setVariable := func(id string, tree antlr.ParseTree) any {
		result := v.Visit(tree)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		v.env.SetVariable(id, result)
		return nil
	}

	switch {
	case ctx.ArrayType() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if ctx.DeclarationExpression() != nil {
			return setVariable(id, ctx.DeclarationExpression())
		}

		return setVariable(id, ctx.ArrayType())
	case ctx.ScalarType().INT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if ctx.DeclarationExpression() != nil {
			return setVariable(id, ctx.DeclarationExpression())
		}

		v.env.SetVariable(id, int(0))
		return nil
	case ctx.ScalarType().UINT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if ctx.DeclarationExpression() != nil {
			return setVariable(id, ctx.DeclarationExpression())
		}

		v.env.SetVariable(id, uint(0))
		return nil
	case ctx.ScalarType().FLOAT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if ctx.DeclarationExpression() != nil {
			return setVariable(id, ctx.DeclarationExpression())
		}

		v.env.SetVariable(id, float32(0))
		return nil
	case ctx.ScalarType().BOOL() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetClBit(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if _, ok := v.env.GetClBitArray(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if ctx.DeclarationExpression() != nil {
			return setVariable(id, ctx.DeclarationExpression())
		}

		v.env.SetVariable(id, false)
		return nil
	case ctx.ScalarType().ANGLE() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetClBit(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if _, ok := v.env.GetClBitArray(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if ctx.DeclarationExpression() != nil {
			result := v.Visit(ctx.ScalarType())
			switch result.(type) {
			case error, endSignal:
				return result
			}

			bits := result.(int64)
			result = v.Visit(ctx.DeclarationExpression())
			switch result.(type) {
			case error, endSignal:
				return result
			}

			radian := result.(float64)
			v.env.SetVariable(id, angle.New(uint(bits), radian))
			return nil
		}

		result := v.Visit(ctx.ScalarType())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		bits := result.(int64)
		v.env.SetVariable(id, angle.New(uint(bits), 0))
		return nil
	case ctx.ScalarType().BIT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if _, ok := v.env.GetClBit(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if _, ok := v.env.GetClBitArray(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if ctx.DeclarationExpression() != nil {
			x := v.Visit(ctx.DeclarationExpression())
			switch x.(type) {
			case error, endSignal:
				return x
			}

			switch {
			case ctx.ScalarType().Designator() != nil:
				switch bits := x.(type) {
				case bool:
					v.env.SetClBitArray(id, []bool{bits})
					return nil
				case []bool:
					v.env.SetClBitArray(id, bits)
					return nil
				default:
					return fmt.Errorf("assign %v(%T) to %q", x, x, id)
				}
			default:
				switch bit := x.(type) {
				case bool:
					v.env.SetClBit(id, bit)
					return nil
				case []bool:
					if len(bit) != 1 {
						return fmt.Errorf("assign %d bits to a single bit", len(bit))
					}

					v.env.SetClBit(id, bit[0])
					return nil
				default:
					return fmt.Errorf("assign %v(%T) to %q", x, x, id)
				}
			}
		}

		// default value
		switch {
		case ctx.ScalarType().Designator() != nil:
			result := v.Visit(ctx.ScalarType())
			switch result.(type) {
			case error, endSignal:
				return result
			}

			size := result.(int64)
			v.env.SetClBitArray(id, make([]bool, int(size)))
			return nil
		default:
			v.env.SetClBit(id, false)
			return nil
		}
	default:
		return fmt.Errorf("unsupported scalar type %q", ctx.ScalarType().GetText())
	}
}

func (v *Visitor) VisitDefStatement(ctx *parser.DefStatementContext) any {
	name := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetSubroutine(name); ok {
		return fmt.Errorf("%q redeclared", name)
	}

	var qargs []string
	if ctx.ArgumentDefinitionList() != nil {
		args := v.Visit(ctx.ArgumentDefinitionList()).([]any)
		for _, a := range args {
			qargs = append(qargs, a.(string))
		}
	}

	var retType any
	if ctx.ReturnSignature() != nil {
		retType = v.Visit(ctx.ReturnSignature())
	}

	v.env.Subroutine[name] = &environ.Subroutine{
		Name:       name,
		QArgs:      qargs,
		Body:       ctx.Scope(),
		ReturnType: retType,
	}

	return nil
}

func (v *Visitor) VisitAliasDeclarationStatement(ctx *parser.AliasDeclarationStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	if _, ok := v.env.GetQubit(id); ok {
		return fmt.Errorf("%q redeclared", id)
	}

	if _, ok := v.env.GetQubitArray(id); ok {
		return fmt.Errorf("%q redeclared", id)
	}

	result := v.Visit(ctx.AliasExpression())
	switch result.(type) {
	case error, endSignal:
		return result
	}

	switch alias := result.(type) {
	case q.Qubit:
		v.env.Qubit[id] = alias
	case []q.Qubit:
		v.env.QubitArray[id] = alias
	default:
		return fmt.Errorf("invalid qubit alias %q", ctx.AliasExpression().GetText())
	}
	return nil
}

func (v *Visitor) VisitOldStyleDeclarationStatement(ctx *parser.OldStyleDeclarationStatementContext) any {
	switch {
	case ctx.QREG() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetQubit(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if _, ok := v.env.GetQubitArray(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		var size int64 = 1
		if ctx.Designator() != nil {
			result := v.Visit(ctx.Designator())
			switch result.(type) {
			case error, endSignal:
				return result
			}

			size = result.(int64)
		}

		need := v.qsim.NumQubits() + int(size)
		if v.maxQubits > 0 && need > v.maxQubits {
			return fmt.Errorf("need=%d, max=%d: %w", need, v.maxQubits, ErrTooManyQubits)
		}

		v.env.SetQubitArray(id, v.qsim.Zeros(int(size)))
		return nil
	case ctx.CREG() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetClBit(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		if _, ok := v.env.GetClBitArray(id); ok {
			return fmt.Errorf("%q redeclared", id)
		}

		var size int64 = 1
		if ctx.Designator() != nil {
			result := v.Visit(ctx.Designator())
			switch result.(type) {
			case error, endSignal:
				return result
			}

			size = result.(int64)
		}

		v.env.SetClBitArray(id, make([]bool, int(size)))
		return nil
	default:
		return fmt.Errorf("unsupported old style declaration %q", ctx.GetText())
	}
}

func (v *Visitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) any {
	switch {
	case ctx.Identifier() != nil:
		s := v.Visit(ctx.Identifier()).(string)
		if lit, ok := Const[s]; ok {
			return lit
		}

		if lit, ok := v.env.GetVariable(s); ok {
			return lit
		}

		if lit, ok := v.env.GetQubit(s); ok {
			return lit
		}

		if lit, ok := v.env.GetQubitArray(s); ok {
			return lit
		}

		if lit, ok := v.env.GetClBit(s); ok {
			return lit
		}

		if lit, ok := v.env.GetClBitArray(s); ok {
			return lit
		}

		if lit, ok := v.env.GetConst(s); ok {
			return lit
		}

		return fmt.Errorf("undefined %q", s)
	case ctx.DecimalIntegerLiteral() != nil:
		s := v.Visit(ctx.DecimalIntegerLiteral()).(string)
		lit, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int %q: %w", s, err)
		}

		return lit
	case ctx.HexIntegerLiteral() != nil:
		s := v.Visit(ctx.HexIntegerLiteral()).(string)
		lit, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return fmt.Errorf("parse int %q: %w", s, err)
		}

		return lit
	case ctx.FloatLiteral() != nil:
		s := v.Visit(ctx.FloatLiteral()).(string)
		lit, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("parse float %q: %w", s, err)
		}

		return lit
	case ctx.BooleanLiteral() != nil:
		s := v.Visit(ctx.BooleanLiteral()).(string)
		lit, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("parse bool %q: %w", s, err)
		}

		return lit
	case ctx.BitstringLiteral() != nil:
		s := v.Visit(ctx.BitstringLiteral()).(string)
		bitstring := strings.Trim(s, "\"")

		lit := make([]bool, len(bitstring))
		for i, b := range bitstring {
			lit[i] = b == '1'
		}

		return lit
	default:
		return fmt.Errorf("unsupported literal %q", ctx.GetText())
	}
}

func (v *Visitor) VisitCastExpression(ctx *parser.CastExpressionContext) any {
	x := v.Visit(ctx.Expression())
	switch x.(type) {
	case error, endSignal:
		return x
	}

	switch {
	case ctx.ScalarType().INT() != nil:
		val, err := value.New(x).Int()
		if err != nil {
			return fmt.Errorf("int(%v): %w", x, err)
		}

		return val.Value()
	case ctx.ScalarType().UINT() != nil:
		val, err := value.New(x).UInt()
		if err != nil {
			return fmt.Errorf("uint(%v): %w", x, err)
		}

		return val.Value()
	case ctx.ScalarType().FLOAT() != nil:
		val, err := value.New(x).Float64()
		if err != nil {
			return fmt.Errorf("float64(%v): %w", x, err)
		}

		return val.Value()
	default:
		return fmt.Errorf("unsupported scalar type %q", ctx.ScalarType().GetText())
	}
}

func (v *Visitor) VisitAdditiveExpression(ctx *parser.AdditiveExpressionContext) any {
	left := v.Visit(ctx.Expression(0))
	switch left.(type) {
	case error, endSignal:
		return left
	}

	right := v.Visit(ctx.Expression(1))
	switch right.(type) {
	case error, endSignal:
		return right
	}

	a, b := value.New(left), value.New(right)
	if ctx.PLUS() != nil {
		val, err := a.Add(b)
		if err != nil {
			return fmt.Errorf("%v+%v: %w", left, right, err)
		}

		return val.Value()
	}

	if ctx.MINUS() != nil {
		val, err := a.Sub(b)
		if err != nil {
			return fmt.Errorf("%v-%v: %w", left, right, err)
		}

		return val.Value()
	}

	return fmt.Errorf("unsupported operator %q", ctx.GetText())
}

func (v *Visitor) VisitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) any {
	left := v.Visit(ctx.Expression(0))
	switch left.(type) {
	case error, endSignal:
		return left
	}

	right := v.Visit(ctx.Expression(1))
	switch right.(type) {
	case error, endSignal:
		return right
	}

	a, b := value.New(left), value.New(right)
	if ctx.ASTERISK() != nil {
		val, err := a.Mul(b)
		if err != nil {
			return fmt.Errorf("%v*%v: %w", left, right, err)
		}

		return val.Value()
	}

	if ctx.SLASH() != nil {
		val, err := a.Div(b)
		if err != nil {
			return fmt.Errorf("%v/%v: %w", left, right, err)
		}

		return val.Value()
	}

	if ctx.PERCENT() != nil {
		val, err := a.Mod(b)
		if err != nil {
			return fmt.Errorf("%v%%%v: %w", left, right, err)
		}

		return val.Value()
	}

	return fmt.Errorf("unsupported operator %q", ctx.GetText())
}

func (v *Visitor) VisitEqualityExpression(ctx *parser.EqualityExpressionContext) any {
	left := v.Visit(ctx.Expression(0))
	switch left.(type) {
	case error, endSignal:
		return left
	}

	right := v.Visit(ctx.Expression(1))
	switch right.(type) {
	case error, endSignal:
		return right
	}

	a, b := value.New(left), value.New(right)
	op := v.Visit(ctx.EqualityOperator()).(string)
	switch op {
	case "==":
		val, err := a.Eq(b)
		if err != nil {
			return fmt.Errorf("%v==%v: %w", left, right, err)
		}

		return val.Value()
	case "!=":
		val, err := a.NotEq(b)
		if err != nil {
			return fmt.Errorf("%v!=%v: %w", left, right, err)
		}

		return val.Value()
	}

	return fmt.Errorf("unsupported operator %q", op)
}

func (v *Visitor) VisitComparisonExpression(ctx *parser.ComparisonExpressionContext) any {
	left := v.Visit(ctx.Expression(0))
	switch left.(type) {
	case error, endSignal:
		return left
	}

	right := v.Visit(ctx.Expression(1))
	switch right.(type) {
	case error, endSignal:
		return right
	}

	a, b := value.New(left), value.New(right)
	op := v.Visit(ctx.ComparisonOperator()).(string)

	switch op {
	case "<":
		val, err := a.LessThan(b)
		if err != nil {
			return fmt.Errorf("%v<%v: %w", left, right, err)
		}

		return val.Value()
	case "<=":
		val, err := a.LessThanOrEqual(b)
		if err != nil {
			return fmt.Errorf("%v<=%v: %w", left, right, err)
		}

		return val.Value()
	case ">":
		val, err := a.GreaterThan(b)
		if err != nil {
			return fmt.Errorf("%v>%v: %w", left, right, err)
		}

		return val.Value()
	case ">=":
		val, err := a.GreaterThanOrEqual(b)
		if err != nil {
			return fmt.Errorf("%v>=%v: %w", left, right, err)
		}

		return val.Value()
	}

	return fmt.Errorf("unsupported operator %q", op)
}

func (v *Visitor) VisitUnaryExpression(ctx *parser.UnaryExpressionContext) any {
	x := v.Visit(ctx.Expression())
	switch x.(type) {
	case error, endSignal:
		return x
	}

	switch {
	case ctx.MINUS() != nil:
		val, err := value.New(x).Negative()
		if err != nil {
			return fmt.Errorf("negate(%v): %w", x, err)
		}

		return val.Value()
	case ctx.EXCLAMATION_POINT() != nil:
		val, err := value.New(x).BoolNot()
		if err != nil {
			return fmt.Errorf("!%v: %w", x, err)
		}

		return val.Value()
	case ctx.TILDE() != nil:
		val, err := value.New(x).BitNot()
		if err != nil {
			return fmt.Errorf("~%v: %w", x, err)
		}

		return val.Value()
	default:
		return fmt.Errorf("unsupported operator %q", ctx.GetOp().GetText())
	}
}

func (v *Visitor) VisitPowerExpression(ctx *parser.PowerExpressionContext) any {
	base := v.Visit(ctx.Expression(0))
	switch base.(type) {
	case error, endSignal:
		return base
	}

	exp := v.Visit(ctx.Expression(1))
	switch exp.(type) {
	case error, endSignal:
		return exp
	}

	a, b := value.New(base), value.New(exp)
	val, err := a.Pow(b)
	if err != nil {
		return fmt.Errorf("%v^%v: %w", base, exp, err)
	}

	return val.Value()
}

func (v *Visitor) VisitLogicalAndExpression(ctx *parser.LogicalAndExpressionContext) any {
	for _, x := range ctx.AllExpression() {
		result := v.Visit(x)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		if !result.(bool) {
			return false
		}
	}

	return true
}

func (v *Visitor) VisitLogicalOrExpression(ctx *parser.LogicalOrExpressionContext) any {
	for _, x := range ctx.AllExpression() {
		result := v.Visit(x)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		if result.(bool) {
			return true
		}
	}

	return false
}

func (v *Visitor) VisitBitwiseAndExpression(ctx *parser.BitwiseAndExpressionContext) any {
	left := v.Visit(ctx.AllExpression()[0])
	switch left.(type) {
	case error, endSignal:
		return left
	}

	result := left.(int64)
	for _, x := range ctx.AllExpression()[1:] {
		right := v.Visit(x)
		switch right.(type) {
		case error, endSignal:
			return right
		}

		result = result & right.(int64)
	}

	return result
}

func (v *Visitor) VisitBitwiseOrExpression(ctx *parser.BitwiseOrExpressionContext) any {
	left := v.Visit(ctx.AllExpression()[0])
	switch left.(type) {
	case error, endSignal:
		return left
	}

	result := left.(int64)
	for _, x := range ctx.AllExpression()[1:] {
		right := v.Visit(x)
		switch right.(type) {
		case error, endSignal:
			return right
		}

		result = result | right.(int64)
	}

	return result
}

func (v *Visitor) VisitBitwiseXorExpression(ctx *parser.BitwiseXorExpressionContext) any {
	left := v.Visit(ctx.AllExpression()[0])
	switch left.(type) {
	case error, endSignal:
		return left
	}

	result := left.(int64)
	for _, x := range ctx.AllExpression()[1:] {
		right := v.Visit(x)
		switch right.(type) {
		case error, endSignal:
			return right
		}

		result = result ^ right.(int64)
	}

	return result
}

func (v *Visitor) VisitBitshiftExpression(ctx *parser.BitshiftExpressionContext) any {
	op := v.Visit(ctx.BitshiftOperator()).(string)
	left := v.Visit(ctx.AllExpression()[0])
	switch left.(type) {
	case error, endSignal:
		return left
	}
	result := left.(int64)
	for _, x := range ctx.AllExpression()[1:] {
		right := v.Visit(x)
		switch right.(type) {
		case error, endSignal:
			return right
		}
		switch op {
		case "<<":
			result = result << right.(int64)
		case ">>":
			result = result >> right.(int64)
		}
	}

	return result
}

func (v *Visitor) VisitDeclarationExpression(ctx *parser.DeclarationExpressionContext) any {
	if ctx.MeasureExpression() != nil {
		return v.Visit(ctx.MeasureExpression())
	}

	if ctx.ArrayLiteral() != nil {
		return v.Visit(ctx.ArrayLiteral())
	}

	return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitCallExpression(ctx *parser.CallExpressionContext) any {
	var args []any
	if ctx.ExpressionList() != nil {
		result := v.Visit(ctx.ExpressionList())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		args = result.([]any)
	}

	id := v.Visit(ctx.Identifier()).(string)
	switch id {
	case "sin":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Sin(a.Value().(float64))
	case "cos":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Cos(a.Value().(float64))
	case "tan":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Tan(a.Value().(float64))
	case "arcsin":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Asin(a.Value().(float64))
	case "arccos":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Acos(a.Value().(float64))
	case "arctan":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Atan(a.Value().(float64))
	case "ceiling":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Ceil(a.Value().(float64))
	case "floor":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Floor(a.Value().(float64))
	case "sqrt":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Sqrt(a.Value().(float64))
	case "exp":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Exp(a.Value().(float64))
	case "log":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		return math.Log(a.Value().(float64))
	case "mod":
		a, err := value.New(args[0]).Float64()
		if err != nil {
			return err
		}

		b, err := value.New(args[1]).Float64()
		if err != nil {
			return err
		}

		return math.Mod(a.Value().(float64), b.Value().(float64))
	default:
		routine, ok := v.env.GetSubroutine(id)
		if !ok {
			return fmt.Errorf("undefined %q", id)
		}

		enclosed := v.Enclosed()
		for i, p := range routine.QArgs {
			switch arg := args[i].(type) {
			case q.Qubit:
				enclosed.env.Qubit[p] = arg
			case []q.Qubit:
				enclosed.env.QubitArray[p] = arg
			default:
				return fmt.Errorf("invalid qubit argument %q", p)
			}
		}

		result := enclosed.Visit(routine.Body)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		if ret, ok := result.(returnValue); ok {
			return ret.value
		}

		return result
	}
}

func (v *Visitor) VisitRangeExpression(ctx *parser.RangeExpressionContext) any {
	var list []int64
	for _, x := range ctx.AllExpression() {
		result := v.Visit(x)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		val, err := value.New(result).Int64()
		if err != nil {
			return fmt.Errorf("int64(%v): %w", x.GetText(), err)
		}

		list = append(list, val.Value().(int64))
	}

	return list
}

func (v *Visitor) VisitAliasExpression(ctx *parser.AliasExpressionContext) any {
	var result []q.Qubit
	for _, x := range ctx.AllExpression() {
		operand := v.Visit(x)
		switch operand.(type) {
		case error, endSignal:
			return operand
		}

		if len(ctx.AllExpression()) == 1 {
			return operand
		}

		switch arg := operand.(type) {
		case q.Qubit:
			result = append(result, arg)
		case []q.Qubit:
			result = append(result, arg...)
		default:
			return fmt.Errorf("invalid qubit alias %q", x.GetText())
		}
	}

	return result
}

func (v *Visitor) VisitIndexExpression(ctx *parser.IndexExpressionContext) any {
	operand := v.Visit(ctx.Expression())
	switch operand.(type) {
	case error, endSignal:
		return operand
	}

	switch operand.(type) {
	case q.Qubit, []q.Qubit:
	default:
		return fmt.Errorf("invalid operand %q", ctx.Expression().GetText())
	}

	return v.IndexQubits(operand, ctx.IndexOperator())
}

func (v *Visitor) IndexQubits(operand any, ctx parser.IIndexOperatorContext) any {
	var qubit []q.Qubit
	switch arg := operand.(type) {
	case q.Qubit:
		qubit = []q.Qubit{arg}
	case []q.Qubit:
		qubit = arg
	default:
		return fmt.Errorf("invalid qubit array operand")
	}

	var result []q.Qubit
	for _, x := range ctx.AllExpression() {
		value := v.Visit(x)
		switch value.(type) {
		case error, endSignal:
			return value
		}

		index, ok := value.(int64)
		if !ok {
			return fmt.Errorf("index must be an integer '%v'", value)
		}

		if index < 0 || index >= int64(len(qubit)) {
			return fmt.Errorf("qubit index %d out of range for size %d", index, len(qubit))
		}

		result = append(result, qubit[index])
	}

	for _, x := range ctx.AllRangeExpression() {
		value := v.Visit(x)
		switch value.(type) {
		case error, endSignal:
			return value
		}

		index := value.([]int64)
		if len(index) != 2 {
			return fmt.Errorf("qubit range %q: %w", x.GetText(), ErrNotImplemented)
		}

		if index[0] < 0 || index[1] < index[0] || index[1] > int64(len(qubit)) {
			return fmt.Errorf("qubit range %q out of range for size %d", x.GetText(), len(qubit))
		}

		result = append(result, qubit[index[0]:index[1]]...)
	}

	if len(ctx.AllExpression()) == 1 && len(ctx.AllRangeExpression()) == 0 {
		return result[0]
	}

	return result
}

func (v *Visitor) VisitMeasureExpression(ctx *parser.MeasureExpressionContext) any {
	result := v.Visit(ctx.GateOperand())
	switch result.(type) {
	case error, endSignal:
		return result
	}

	qargs := result.(QubitOperand).Qubits
	v.qsim.Measure(qargs...)

	var bits []bool
	for _, q := range qargs {
		binary := v.qsim.State(q)[0].BinaryString()[0]
		intv := number.MustParseInt(binary)
		bits = append(bits, intv == 1)
	}

	if len(bits) == 1 {
		return bits[0]
	}

	return bits
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

func (v *Visitor) VisitIndexOperator(ctx *parser.IndexOperatorContext) any {
	var list []any
	for _, x := range ctx.AllExpression() {
		result := v.Visit(x)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		list = append(list, result)
	}

	return list
}

func (v *Visitor) VisitIndexedIdentifier(ctx *parser.IndexedIdentifierContext) any {
	var index []int64
	for _, operator := range ctx.AllIndexOperator() {
		result := v.Visit(operator)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		for _, o := range result.([]any) {
			idx, ok := o.(int64)
			if !ok {
				return fmt.Errorf("index must be an integer '%v'", o)
			}

			index = append(index, idx)
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
		result := v.Visit(ctx.Designator())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		val, ok := result.(int64)
		if !ok {
			return fmt.Errorf("size must be an integer %q", ctx.Designator().GetText())
		}

		return val
	}

	return int64(1)
}

func (v *Visitor) VisitQubitType(ctx *parser.QubitTypeContext) any {
	if ctx.Designator() != nil {
		result := v.Visit(ctx.Designator())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		val, ok := result.(int64)
		if !ok {
			return fmt.Errorf("size must be an integer %q", ctx.Designator().GetText())
		}

		return val
	}

	return int64(1)
}

func (v *Visitor) VisitArrayType(ctx *parser.ArrayTypeContext) any {
	result := v.Visit(ctx.ExpressionList())
	switch result.(type) {
	case error, endSignal:
		return result
	}

	size := result.([]any)[0].(int64)
	scalar := ctx.ScalarType()

	switch {
	case scalar.INT() != nil:
		result := v.Visit(scalar.Designator())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		bits := result.(int64)
		switch bits {
		case 8:
			return make([]int8, size)
		case 16:
			return make([]int16, size)
		case 32:
			return make([]int32, size)
		case 64:
			return make([]int64, size)
		default:
			return fmt.Errorf("unsupported bit size %d", bits)
		}

	case scalar.UINT() != nil:
		result := v.Visit(scalar.Designator())
		switch result.(type) {
		case error, endSignal:
			return result
		}

		bits := result.(int64)
		switch bits {
		case 8:
			return make([]uint8, size)
		case 16:
			return make([]uint16, size)
		case 32:
			return make([]uint32, size)
		case 64:
			return make([]uint64, size)
		default:
			return fmt.Errorf("unsupported bit size %d", bits)
		}

	case scalar.FLOAT() != nil:
		result := v.Visit(scalar.Designator())
		switch result.(type) {
		case error, endSignal:
			return result
		}
		bits := result.(int64)
		switch bits {
		case 32:
			return make([]float32, size)
		case 64:
			return make([]float64, size)
		default:
			return fmt.Errorf("unsupported bit size %d", bits)
		}
	case scalar.BOOL() != nil:
		return make([]bool, size)
	default:
		return fmt.Errorf("unsupported scalar type %q", scalar.GetText())
	}
}

func (v *Visitor) VisitArrayLiteral(ctx *parser.ArrayLiteralContext) any {
	var list []any
	for _, x := range ctx.AllExpression() {
		result := v.Visit(x)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		list = append(list, result)
	}

	return list
}

func (v *Visitor) VisitReturnSignature(ctx *parser.ReturnSignatureContext) any {
	return ctx.ScalarType()
}

func (v *Visitor) VisitArgumentDefinition(ctx *parser.ArgumentDefinitionContext) any {
	return v.Visit(ctx.Identifier())
}

func (v *Visitor) VisitGateOperand(ctx *parser.GateOperandContext) any {
	operand := v.Visit(ctx.IndexedIdentifier().Identifier()).(string)
	var result any
	if qb, ok := v.env.GetQubit(operand); ok {
		result = qb
	} else if qb, ok := v.env.GetQubitArray(operand); ok {
		result = qb
	} else {
		return fmt.Errorf("undefined %q", operand)
	}

	for _, index := range ctx.IndexedIdentifier().AllIndexOperator() {
		result = v.IndexQubits(result, index)
		switch result.(type) {
		case error, endSignal:
			return result
		}
	}

	switch arg := result.(type) {
	case q.Qubit:
		return QubitOperand{Qubits: []q.Qubit{arg}}
	case []q.Qubit:
		return QubitOperand{Qubits: arg, IsArray: true}
	}

	return fmt.Errorf("invalid gate operand %q", ctx.GetText())
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
		result := v.Visit(x)
		switch result.(type) {
		case error, endSignal:
			return result
		}

		list = append(list, result)
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
	var list []QubitOperand
	for _, operand := range ctx.AllGateOperand() {
		op := v.Visit(operand)
		switch op.(type) {
		case error, endSignal:
			return op
		}

		list = append(list, op.(QubitOperand))
	}

	return list
}

const (
	Break    string = "break;"
	Continue string = "continue;"
	End      string = "end;"
)

// contains returns true if the result contains substrings.
func contains(result any, substrings ...string) bool {
	switch v := result.(type) {
	case string:
		for _, s := range substrings {
			if strings.Contains(v, s) {
				return true
			}
		}
	case []any:
		for _, r := range v {
			if contains(r, substrings...) {
				return true
			}
		}
	}

	return false
}

func (v *Visitor) VisitPragma(ctx *parser.PragmaContext) any {
	return fmt.Errorf("VisitPragma: %w", ErrNotImplemented)
}

func (v *Visitor) VisitAnnotation(ctx *parser.AnnotationContext) any {
	return fmt.Errorf("VisitAnnotation: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDurationofExpression(ctx *parser.DurationofExpressionContext) any {
	return fmt.Errorf("VisitDurationofExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitSetExpression(ctx *parser.SetExpressionContext) any {
	return fmt.Errorf("VisitSetExpression: %w", ErrNotImplemented)
}

func (v *Visitor) VisitArrayReferenceType(ctx *parser.ArrayReferenceTypeContext) any {
	return fmt.Errorf("VisitArrayReferenceType: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalArgumentDefinitionList(ctx *parser.DefcalArgumentDefinitionListContext) any {
	return fmt.Errorf("VisitDefcalArgumentDefinitionList: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalArgumentDefinition(ctx *parser.DefcalArgumentDefinitionContext) any {
	return fmt.Errorf("VisitDefcalArgumentDefinition: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalTarget(ctx *parser.DefcalTargetContext) any {
	return fmt.Errorf("VisitDefcalTarget: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalOperandList(ctx *parser.DefcalOperandListContext) any {
	return fmt.Errorf("VisitDefcalOperandList: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDefcalOperand(ctx *parser.DefcalOperandContext) any {
	return fmt.Errorf("VisitDefcalOperand: %w", ErrNotImplemented)
}

func (v *Visitor) VisitIoDeclarationStatement(ctx *parser.IoDeclarationStatementContext) any {
	return fmt.Errorf("VisitIoDeclarationStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitExternStatement(ctx *parser.ExternStatementContext) any {
	return fmt.Errorf("VisitExternStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitExternArgumentList(ctx *parser.ExternArgumentListContext) any {
	return fmt.Errorf("VisitExternArgumentList: %w", ErrNotImplemented)
}

func (v *Visitor) VisitExternArgument(ctx *parser.ExternArgumentContext) any {
	return fmt.Errorf("VisitExternArgument: %w", ErrNotImplemented)
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

func (v *Visitor) VisitBoxStatement(ctx *parser.BoxStatementContext) any {
	return fmt.Errorf("VisitBoxStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDelayStatement(ctx *parser.DelayStatementContext) any {
	return fmt.Errorf("VisitDelayStatement: %w", ErrNotImplemented)
}
