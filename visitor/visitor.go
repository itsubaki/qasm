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
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/qasm/angle"
	"github.com/itsubaki/qasm/environ"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/value"
)

var (
	ErrAlreadyDeclared    = errors.New("already declared")
	ErrIdentifierNotFound = errors.New("identifier not found")
	ErrQubitNotFound      = errors.New("qubit not found")
	ErrBitNotFound        = errors.New("bit not found")
	ErrVariableNotFound   = errors.New("variable not found")
	ErrGateNotFound       = errors.New("gate not found")
	ErrFunctionNotFound   = errors.New("function not found")
	ErrTooManyQubits      = errors.New("too many qubits")
	ErrBitWidthMismatch   = errors.New("bit width mismatch")
	ErrInvalidOperand     = errors.New("invalid operand")
	ErrUnexpected         = errors.New("unexpected")
	ErrNotImplemented     = errors.New("not implemented")
)

type Visitor struct {
	qsim      *q.Q
	env       *environ.Environ
	maxQubits int
}

func New(qsim *q.Q, env *environ.Environ, opt ...Option) *Visitor {
	v := &Visitor{
		qsim: qsim,
		env:  env,
	}

	for _, f := range opt {
		f(v)
	}

	return v
}

func (v *Visitor) Enclosed() *Visitor {
	return New(v.qsim, v.env.NewEnclosed())
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
	return fmt.Errorf("%v: %w", node.GetText(), ErrUnexpected)
}

func (v *Visitor) VisitChildren(node antlr.RuleNode) any {
	for _, c := range node.GetChildren() {
		tree, ok := c.(antlr.ParseTree)
		if !ok {
			return fmt.Errorf("unexpected child(%T): %w", c, ErrUnexpected)
		}

		if err := v.Run(tree); err != nil {
			return err
		}
	}

	return nil
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
	enclosed := v.Enclosed()

	var list []any
	for _, s := range ctx.AllStatementOrScope() {
		result := enclosed.Visit(s)
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

	if err := v.Run(p.Program()); err != nil {
		return fmt.Errorf("include: %w", err)
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
	unwrap := func(v any) any {
		if s, ok := v.([]any); ok {
			if len(s) == 1 && s[0] == nil {
				return nil
			}
		}

		return v
	}

	enclosed := v.Enclosed()
	if v.Visit(ctx.Expression()).(bool) {
		return unwrap(enclosed.Visit(ctx.GetIf_body()))
	}

	if ctx.GetElse_body() != nil {
		return unwrap(enclosed.Visit(ctx.GetElse_body()))
	}

	return nil
}

func (v *Visitor) VisitForStatement(ctx *parser.ForStatementContext) any {
	id := v.Visit(ctx.Identifier()).(string)
	rx := v.Visit(ctx.RangeExpression()).([]int64)

	enclosed := v.Enclosed()
	for i := rx[0]; i <= rx[1]; i++ {
		enclosed.env.SetVariable(id, i)
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

	v.env.Gate[name] = &environ.Gate{
		Name:   name,
		Params: params,
		QArgs:  qargs,
		Body:   ctx.Scope(),
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

func (v *Visitor) Builtin(ctx *parser.GateCallStatementContext) (*matrix.Matrix, bool, error) {
	if ctx.GPHASE() != nil {
		params, err := v.Params(ctx.ExpressionList())
		if err != nil {
			return nil, false, err
		}

		// u 2x2
		u := gate.I().Mul(cmplx.Exp(complex(0, params[0])))
		return u, true, nil
	}

	id := v.Visit(ctx.Identifier()).(string)
	switch id {
	case U:
		// params
		params, err := v.Params(ctx.ExpressionList())
		if err != nil {
			return nil, false, err
		}

		// u 2x2
		u := gate.U(params[0], params[1], params[2])
		return u, true, nil
	default:
		return nil, false, nil
	}
}

func (v *Visitor) UserDefinedGateCall(ctx *parser.GateCallStatementContext) error {
	if len(ctx.AllGateModifier()) > 0 {
		// NOTE: modifier is not implemented in user-defined gate call
		return fmt.Errorf("modifier is not implemented in user-defined gate call: %w", ErrNotImplemented)
	}

	id := v.Visit(ctx.Identifier()).(string)
	g, ok := v.env.GetGate(id)
	if !ok {
		return fmt.Errorf("identifier=%s: %w", id, ErrGateNotFound)
	}

	// params
	enclosed := v.Enclosed()
	if ctx.ExpressionList() != nil {
		params, err := v.Params(ctx.ExpressionList())
		if err != nil {
			return err
		}

		for i, p := range g.Params {
			enclosed.env.SetVariable(p, params[i])
		}
	}

	// qargs
	if ctx.GateOperandList() != nil {
		qargs, ok := v.Visit(ctx.GateOperandList()).([][]q.Qubit)
		if !ok {
			return fmt.Errorf("operand=%v: %w", ctx.GateOperandList().GetText(), ErrInvalidOperand)
		}

		for i, id := range g.QArgs {
			enclosed.env.Qubit[id] = qargs[i]
		}
	}

	// call body
	for i, s := range g.Body.AllStatementOrScope() {
		call := s.Statement().GateCallStatement().(*parser.GateCallStatementContext)
		result := enclosed.VisitGateCallStatement(call)
		if err, ok := result.(error); ok && err != nil {
			return fmt.Errorf("gate call[%d]: %w", i, err)
		}
	}

	return nil
}

func (v *Visitor) VisitGateCallStatement(ctx *parser.GateCallStatementContext) any {
	u, ok, err := v.Builtin(ctx)
	if err != nil {
		return fmt.Errorf("builtin: %w", err)
	}

	if !ok {
		if err := v.UserDefinedGateCall(ctx); err != nil {
			return fmt.Errorf("user-defined gate call: %w", err)
		}

		return nil
	}

	if HasControlModifier(ctx) {
		// qubit[2] c;
		// qubit t;
		// U(pi/2, 0, pi) c;
		// U(pi/2, 0, pi) c[0], c[1];
		// ctrl @ U(pi, 0, pi) c, t;
		// ctrl @ U(pi, 0, pi) c[0], t;
		qargs, ok := v.Visit(ctx.GateOperandList()).([][]q.Qubit)
		if !ok {
			return fmt.Errorf("operand=%v: %w", ctx.GateOperandList().GetText(), ErrInvalidOperand)
		}

		var ctrl, neg []q.Qubit
		var cursor int
		for _, mod := range ctx.AllGateModifier() {
			switch {
			case mod.INV() != nil:
				u = u.Dagger()
			case mod.POW() != nil:
				// NOTE: pow is not implemented with control modifier
				return fmt.Errorf("pow with control modifier is not implemented: %w", ErrNotImplemented)
			case mod.CTRL() != nil:
				n, ok := v.Visit(mod).(int64)
				if !ok {
					return fmt.Errorf("%v: %w", mod.GetText(), ErrUnexpected)
				}

				for range n {
					ctrl = append(ctrl, qargs[cursor]...)
					cursor++
				}
			case mod.NEGCTRL() != nil:
				n, ok := v.Visit(mod).(int64)
				if !ok {
					return fmt.Errorf("%v: %w", mod.GetText(), ErrUnexpected)
				}

				for range n {
					ctrl = append(ctrl, qargs[cursor]...)
					neg = append(neg, qargs[cursor]...)
					cursor++
				}
			}
		}

		v.qsim.X(neg...)
		defer v.qsim.X(neg...)

		target := qargs[len(qargs)-1]
		v.qsim.Controlled(u, ctrl, target)
		return nil
	}

	// no control modifier
	for _, mod := range ReversedModifier(ctx) {
		switch {
		case mod.INV() != nil:
			u = u.Dagger()
		case mod.POW() != nil:
			p, err := value.New(v.Visit(mod)).Float64()
			if err != nil {
				return fmt.Errorf("%v: %w", mod.GetText(), err)
			}

			u = Pow(u, p.Value().(float64))
		}
	}

	// qargs
	var qargs []q.Qubit
	if ctx.GateOperandList() != nil {
		// qubit q0; qubit q1; U q0, q1;
		// qubit[2] q; U q;
		operand, ok := v.Visit(ctx.GateOperandList()).([][]q.Qubit)
		if !ok {
			return fmt.Errorf("operand=%v: %w", ctx.GateOperandList().GetText(), ErrInvalidOperand)
		}

		for _, o := range operand {
			qargs = append(qargs, o...)
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

func (v *Visitor) MeasureAssignment(id parser.IIndexedIdentifierContext, measure parser.IMeasureExpressionContext) error {
	measured := v.Visit(measure)
	if id == nil {
		return nil
	}

	operand := v.Visit(id.Identifier()).(string)
	index := v.Visit(id).([]int64)

	if bits, ok := v.env.GetBitArray(operand); ok {
		var val []bool
		switch m := measured.(type) {
		case bool:
			val = []bool{m}
		case []bool:
			val = m
		}

		if len(index) == 0 {
			v.env.SetBitArray(operand, val)
			return nil
		}

		for i, bit := range val {
			bits[index[i]] = bit
		}

		return nil
	}

	if _, ok := v.env.GetBit(operand); ok {
		switch m := measured.(type) {
		case bool:
			v.env.SetBit(operand, m)
			return nil
		case []bool:
			return fmt.Errorf("bit length=%d: %w", len(m), ErrBitWidthMismatch)
		}
	}

	return fmt.Errorf("operand=%s: %w", operand, ErrBitNotFound)
}

func (v *Visitor) VisitMeasureArrowAssignmentStatement(ctx *parser.MeasureArrowAssignmentStatementContext) any {
	return v.MeasureAssignment(ctx.IndexedIdentifier(), ctx.MeasureExpression())
}

func (v *Visitor) VisitAssignmentStatement(ctx *parser.AssignmentStatementContext) any {
	if ctx.MeasureExpression() != nil {
		return v.MeasureAssignment(ctx.IndexedIdentifier(), ctx.MeasureExpression())
	}

	id := ctx.IndexedIdentifier()
	operand := v.Visit(id.Identifier()).(string)
	index := v.Visit(id).([]int64)
	x := v.Visit(ctx.Expression())

	if bits, ok := v.env.GetBitArray(operand); ok {
		var val []bool
		switch m := x.(type) {
		case bool:
			val = []bool{m}
		case []bool:
			val = m
		}

		if len(index) == 0 {
			v.env.SetBitArray(operand, val)
			return nil
		}

		for i, bit := range val {
			bits[index[i]] = bit
		}

		return nil
	}

	if _, ok := v.env.GetBit(operand); ok {
		switch bit := x.(type) {
		case bool:
			v.env.SetBit(operand, bit)
			return nil
		case []bool:
			if len(bit) != 1 {
				return fmt.Errorf("bit length=%d: %w", len(bit), ErrBitWidthMismatch)
			}

			v.env.SetBit(operand, bit[0])
			return nil
		}
	}

	v.env.SetVariable(operand, x)
	return nil
}

func (v *Visitor) VisitResetStatement(ctx *parser.ResetStatementContext) any {
	qargs, ok := v.Visit(ctx.GateOperand()).([]q.Qubit)
	if !ok {
		return fmt.Errorf("operand=%v: %w", ctx.GateOperand().GetText(), ErrInvalidOperand)
	}

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
	need := v.qsim.NumQubits() + int(size)
	if v.maxQubits > 0 && need > v.maxQubits {
		return fmt.Errorf("need=%d, max=%d: %w", need, v.maxQubits, ErrTooManyQubits)
	}

	v.env.SetQubit(id, v.qsim.Zeros(int(size)))
	return nil
}

func (v *Visitor) VisitClassicalDeclarationStatement(ctx *parser.ClassicalDeclarationStatementContext) any {
	switch {
	case ctx.ArrayType() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.SetVariable(id, v.Visit(ctx.DeclarationExpression()))
			return nil
		}

		v.env.SetVariable(id, v.Visit(ctx.ArrayType()))
		return nil
	case ctx.ScalarType().INT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.SetVariable(id, v.Visit(ctx.DeclarationExpression()))
			return nil
		}

		v.env.SetVariable(id, int(0))
		return nil
	case ctx.ScalarType().UINT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.SetVariable(id, v.Visit(ctx.DeclarationExpression()))
			return nil
		}

		v.env.SetVariable(id, uint(0))
		return nil
	case ctx.ScalarType().FLOAT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.SetVariable(id, v.Visit(ctx.DeclarationExpression()))
			return nil
		}

		v.env.SetVariable(id, float32(0))
		return nil
	case ctx.ScalarType().BOOL() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetBit(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if _, ok := v.env.GetBitArray(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			v.env.SetVariable(id, v.Visit(ctx.DeclarationExpression()))
			return nil
		}

		v.env.SetVariable(id, false)
		return nil
	case ctx.ScalarType().ANGLE() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetBit(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if _, ok := v.env.GetBitArray(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			bits := v.Visit(ctx.ScalarType()).(int64)
			radian := v.Visit(ctx.DeclarationExpression()).(float64)
			v.env.SetVariable(id, angle.New(uint(bits), radian))

			return nil
		}

		bits := v.Visit(ctx.ScalarType()).(int64)
		v.env.SetVariable(id, angle.New(uint(bits), 0))
		return nil
	case ctx.ScalarType().BIT() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetVariable(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if _, ok := v.env.GetBit(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if _, ok := v.env.GetBitArray(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if ctx.DeclarationExpression() != nil {
			x := v.Visit(ctx.DeclarationExpression())
			switch {
			case ctx.ScalarType().Designator() != nil:
				switch bits := x.(type) {
				case bool:
					v.env.SetBitArray(id, []bool{bits})
					return nil
				case []bool:
					v.env.SetBitArray(id, bits)
					return nil
				default:
					return fmt.Errorf("declaration=%v(%T): %w", x, x, ErrUnexpected)
				}
			default:
				switch bit := x.(type) {
				case bool:
					v.env.SetBit(id, bit)
					return nil
				case []bool:
					if len(bit) != 1 {
						return fmt.Errorf("declaration length=%d: %w", len(bit), ErrBitWidthMismatch)
					}

					v.env.SetBit(id, bit[0])
					return nil
				default:
					return fmt.Errorf("declaration=%v(%T): %w", x, x, ErrUnexpected)
				}
			}
		}

		// default value
		switch {
		case ctx.ScalarType().Designator() != nil:
			size := v.Visit(ctx.ScalarType()).(int64)
			v.env.SetBitArray(id, make([]bool, int(size)))
			return nil
		default:
			v.env.SetBit(id, false)
			return nil
		}
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

		need := v.qsim.NumQubits() + int(size)
		if v.maxQubits > 0 && need > v.maxQubits {
			return fmt.Errorf("need=%d, max=%d: %w", need, v.maxQubits, ErrTooManyQubits)
		}

		v.env.SetQubit(id, v.qsim.Zeros(int(size)))
		return nil
	case ctx.CREG() != nil:
		id := v.Visit(ctx.Identifier()).(string)
		if _, ok := v.env.GetBit(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		if _, ok := v.env.GetBitArray(id); ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrAlreadyDeclared)
		}

		var size int64 = 1
		if ctx.Designator() != nil {
			size = v.Visit(ctx.Designator()).(int64)
		}

		v.env.SetBitArray(id, make([]bool, int(size)))
		return nil
	default:
		return fmt.Errorf("x=%s: %w", ctx.GetText(), ErrUnexpected)
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

		if lit, ok := v.env.GetBit(s); ok {
			return lit
		}

		if lit, ok := v.env.GetBitArray(s); ok {
			return lit
		}

		if lit, ok := v.env.GetConst(s); ok {
			return lit
		}

		return fmt.Errorf("identifier=%s: %w", s, ErrIdentifierNotFound)
	case ctx.DecimalIntegerLiteral() != nil:
		s := v.Visit(ctx.DecimalIntegerLiteral()).(string)
		lit, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int=%s: %w", s, err)
		}

		return lit
	case ctx.HexIntegerLiteral() != nil:
		s := v.Visit(ctx.HexIntegerLiteral()).(string)
		lit, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return fmt.Errorf("parse int=%s: %w", s, err)
		}

		return lit
	case ctx.FloatLiteral() != nil:
		s := v.Visit(ctx.FloatLiteral()).(string)
		lit, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("parse float=%s: %w", s, err)
		}

		return lit
	case ctx.BooleanLiteral() != nil:
		s := v.Visit(ctx.BooleanLiteral()).(string)
		lit, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("parse bool=%s: %w", s, err)
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
		return fmt.Errorf("x=%s: %w", ctx.GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitCastExpression(ctx *parser.CastExpressionContext) any {
	val := v.Visit(ctx.Expression())
	switch {
	case ctx.ScalarType().INT() != nil:
		v, err := value.New(val).Int()
		if err != nil {
			return fmt.Errorf("int(%v): %w", val, err)
		}

		return v.Value()
	case ctx.ScalarType().UINT() != nil:
		v, err := value.New(val).UInt()
		if err != nil {
			return fmt.Errorf("uint(%v): %w", val, err)
		}

		return v.Value()
	case ctx.ScalarType().FLOAT() != nil:
		v, err := value.New(val).Float64()
		if err != nil {
			return fmt.Errorf("float64(%v): %w", val, err)
		}

		return v.Value()
	default:
		return fmt.Errorf("scalar type=%s: %w", ctx.ScalarType().GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitAdditiveExpression(ctx *parser.AdditiveExpressionContext) any {
	left := v.Visit(ctx.Expression(0))
	right := v.Visit(ctx.Expression(1))
	a, b := value.New(left), value.New(right)

	if ctx.PLUS() != nil {
		v, err := a.Add(b)
		if err != nil {
			return fmt.Errorf("%v+%v: %w", left, right, err)
		}

		return v.Value()
	}

	if ctx.MINUS() != nil {
		v, err := a.Sub(b)
		if err != nil {
			return fmt.Errorf("%v-%v: %w", left, right, err)
		}

		return v.Value()
	}

	return fmt.Errorf("unexpected operator=%q: %w", ctx.GetText(), ErrUnexpected)
}

func (v *Visitor) VisitMultiplicativeExpression(ctx *parser.MultiplicativeExpressionContext) any {
	left := v.Visit(ctx.Expression(0))
	right := v.Visit(ctx.Expression(1))
	a, b := value.New(left), value.New(right)

	if ctx.ASTERISK() != nil {
		w, err := a.Mul(b)
		if err != nil {
			return fmt.Errorf("%v*%v: %w", left, right, err)
		}

		return w.Value()
	}

	if ctx.SLASH() != nil {
		w, err := a.Div(b)
		if err != nil {
			return fmt.Errorf("%v/%v: %w", left, right, err)
		}

		return w.Value()
	}

	if ctx.PERCENT() != nil {
		w, err := a.Mod(b)
		if err != nil {
			return fmt.Errorf("%v%%%v: %w", left, right, err)
		}

		return w.Value()
	}

	return fmt.Errorf("unexpected operator=%q: %w", ctx.GetText(), ErrUnexpected)
}

func (v *Visitor) VisitEqualityExpression(ctx *parser.EqualityExpressionContext) any {
	left := v.Visit(ctx.Expression(0))
	right := v.Visit(ctx.Expression(1))
	a, b := value.New(left), value.New(right)

	op := v.Visit(ctx.EqualityOperator()).(string)
	switch op {
	case "==":
		w, err := a.Eq(b)
		if err != nil {
			return fmt.Errorf("%v==%v: %w", left, right, err)
		}

		return w.Value()
	case "!=":
		w, err := a.NotEq(b)
		if err != nil {
			return fmt.Errorf("%v!=%v: %w", left, right, err)
		}

		return w.Value()
	}

	return fmt.Errorf("unexpected operator=%q: %w", op, ErrUnexpected)
}

func (v *Visitor) VisitComparisonExpression(ctx *parser.ComparisonExpressionContext) any {
	left := v.Visit(ctx.Expression(0))
	right := v.Visit(ctx.Expression(1))
	a, b := value.New(left), value.New(right)

	op := v.Visit(ctx.ComparisonOperator()).(string)
	switch op {
	case "<":
		w, err := a.LessThan(b)
		if err != nil {
			return fmt.Errorf("%v<%v: %w", left, right, err)
		}

		return w.Value()
	case "<=":
		w, err := a.LessThanOrEqual(b)
		if err != nil {
			return fmt.Errorf("%v<=%v: %w", left, right, err)
		}

		return w.Value()
	case ">":
		w, err := a.GreaterThan(b)
		if err != nil {
			return fmt.Errorf("%v>%v: %w", left, right, err)
		}

		return w.Value()
	case ">=":
		w, err := a.GreaterThanOrEqual(b)
		if err != nil {
			return fmt.Errorf("%v>=%v: %w", left, right, err)
		}

		return w.Value()
	}

	return fmt.Errorf("unexpected operator=%q: %w", op, ErrUnexpected)
}

func (v *Visitor) VisitUnaryExpression(ctx *parser.UnaryExpressionContext) any {
	x := v.Visit(ctx.Expression())
	switch {
	case ctx.MINUS() != nil:
		w, err := value.New(x).Negative()
		if err != nil {
			return fmt.Errorf("negate(%v): %w", x, err)
		}

		return w.Value()
	case ctx.EXCLAMATION_POINT() != nil:
		w, err := value.New(x).BoolNot()
		if err != nil {
			return fmt.Errorf("!%v: %w", x, err)
		}

		return w.Value()
	case ctx.TILDE() != nil:
		w, err := value.New(x).BitNot()
		if err != nil {
			return fmt.Errorf("~%v: %w", x, err)
		}

		return w.Value()
	default:
		return fmt.Errorf("operator=%s: %w", ctx.GetOp().GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitPowerExpression(ctx *parser.PowerExpressionContext) any {
	base := v.Visit(ctx.Expression(0))
	exp := v.Visit(ctx.Expression(1))
	a, b := value.New(base), value.New(exp)

	w, err := a.Pow(b)
	if err != nil {
		return fmt.Errorf("%v**%v: %w", base, exp, err)
	}

	return w.Value()
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
		routine, ok := v.env.GetSubroutine(id)
		if !ok {
			return fmt.Errorf("identifier=%s: %w", id, ErrFunctionNotFound)
		}

		enclosed := v.Enclosed()
		for i, p := range routine.QArgs {
			enclosed.env.Qubit[p] = args[i].([]q.Qubit)
		}

		result := enclosed.Visit(routine.Body).([]any)
		return result[len(result)-1]
	}
}

func (v *Visitor) VisitRangeExpression(ctx *parser.RangeExpressionContext) any {
	var list []int64
	for _, x := range ctx.AllExpression() {
		visited := v.Visit(x)
		val, err := value.New(visited).Int64()
		if err != nil {
			return fmt.Errorf("int64(%v): %w", visited, err)
		}

		list = append(list, val.Value().(int64))
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
	qubit, ok := v.Visit(ctx.Expression()).([]q.Qubit)
	if !ok {
		return fmt.Errorf("%v: %w", ctx.Expression().GetText(), ErrUnexpected)
	}

	var result []q.Qubit
	for _, x := range ctx.IndexOperator().AllExpression() {
		index := v.Visit(x).(int64)
		result = append(result, qubit[index])
	}

	for _, x := range ctx.IndexOperator().AllRangeExpression() {
		index := v.Visit(x).([]int64)
		result = append(result, qubit[index[0]:index[1]]...)
	}

	return result
}

func (v *Visitor) VisitMeasureExpression(ctx *parser.MeasureExpressionContext) any {
	qargs, ok := v.Visit(ctx.GateOperand()).([]q.Qubit)
	if !ok {
		return fmt.Errorf("operand=%v: %w", ctx.GateOperand().GetText(), ErrInvalidOperand)
	}

	v.qsim.Measure(qargs...)

	var bits []bool
	for _, q := range qargs {
		bits = append(bits, v.qsim.State(q)[0].Int()[0] == 1)
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
		list = append(list, v.Visit(x))
	}

	return list
}

func (v *Visitor) VisitIndexedIdentifier(ctx *parser.IndexedIdentifierContext) any {
	var index []int64
	for _, op := range ctx.AllIndexOperator() {
		for _, val := range v.Visit(op).([]any) {
			idx, ok := val.(int64)
			if !ok {
				return fmt.Errorf("cast %v to int64: %w", val, ErrUnexpected)
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
		val, ok := v.Visit(ctx.Designator()).(int64)
		if !ok {
			return fmt.Errorf("%v: %w", ctx.Designator().GetText(), ErrUnexpected)
		}

		return val
	}

	return int64(1)
}

func (v *Visitor) VisitQubitType(ctx *parser.QubitTypeContext) any {
	if ctx.Designator() != nil {
		val, ok := v.Visit(ctx.Designator()).(int64)
		if !ok {
			return fmt.Errorf("%v: %w", ctx.Designator().GetText(), ErrUnexpected)
		}

		return val
	}

	return int64(1)
}

func (v *Visitor) VisitArrayType(ctx *parser.ArrayTypeContext) any {
	size := v.Visit(ctx.ExpressionList()).([]any)[0].(int64)

	scalar := ctx.ScalarType()
	switch {
	case scalar.INT() != nil:
		bits := v.Visit(scalar.Designator()).(int64)
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
			return fmt.Errorf("bit size=%d: %w", bits, ErrUnexpected)
		}

	case scalar.UINT() != nil:
		bits := v.Visit(scalar.Designator()).(int64)
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
			return fmt.Errorf("bit size=%d: %w", bits, ErrUnexpected)
		}

	case scalar.FLOAT() != nil:
		bits := v.Visit(scalar.Designator()).(int64)
		switch bits {
		case 32:
			return make([]float32, size)
		case 64:
			return make([]float64, size)
		default:
			return fmt.Errorf("bit size=%d: %w", bits, ErrUnexpected)
		}
	case scalar.BOOL() != nil:
		return make([]bool, size)
	default:
		return fmt.Errorf("scalar type=%s: %w", scalar.GetText(), ErrUnexpected)
	}
}

func (v *Visitor) VisitArrayLiteral(ctx *parser.ArrayLiteralContext) any {
	var list []any
	for _, x := range ctx.AllExpression() {
		list = append(list, v.Visit(x))
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
	id := ctx.IndexedIdentifier()
	operand := v.Visit(id.Identifier()).(string)

	qb, ok := v.env.GetQubit(operand)
	if !ok {
		return fmt.Errorf("operand=%s: %w", operand, ErrQubitNotFound)
	}

	index := v.Visit(id).([]int64)
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
		operand, ok := v.Visit(o).([]q.Qubit)
		if !ok {
			return fmt.Errorf("operand=%v: %w", o.GetText(), ErrInvalidOperand)
		}

		list = append(list, operand)
	}

	return list
}

const (
	Break    string = "break;"
	Continue string = "continue;"
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

func (v *Visitor) VisitBarrierStatement(ctx *parser.BarrierStatementContext) any {
	return fmt.Errorf("VisitBarrierStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitBoxStatement(ctx *parser.BoxStatementContext) any {
	return fmt.Errorf("VisitBoxStatement: %w", ErrNotImplemented)
}

func (v *Visitor) VisitDelayStatement(ctx *parser.DelayStatementContext) any {
	return fmt.Errorf("VisitDelayStatement: %w", ErrNotImplemented)
}
