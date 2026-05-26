package svg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
)

type Visitor struct {
	*parser.Baseqasm3ParserVisitor
	circuit *Circuit
	wire    map[string]int
}

func NewVisitor() *Visitor {
	return &Visitor{
		Baseqasm3ParserVisitor: &parser.Baseqasm3ParserVisitor{},
		circuit:                &Circuit{},
		wire:                   make(map[string]int),
	}
}

func (v *Visitor) Build(tree antlr.ParseTree) (*Circuit, error) {
	if err, ok := v.Visit(tree).(error); ok && err != nil {
		return nil, err
	}

	return v.circuit, nil
}

func (v *Visitor) Add(wireID string) error {
	if _, ok := v.wire[wireID]; ok {
		return fmt.Errorf("%q redeclared", wireID)
	}

	v.wire[wireID] = len(v.circuit.Wires)
	v.circuit.Wires = append(v.circuit.Wires, Wire{
		Name: wireID,
	})

	return nil
}

func (v *Visitor) Visit(tree antlr.ParseTree) any {
	return tree.Accept(v)
}

func (v *Visitor) VisitTerminal(node antlr.TerminalNode) any {
	return node.GetText()
}

func (v *Visitor) VisitProgram(ctx *parser.ProgramContext) any {
	for _, s := range ctx.AllStatementOrScope() {
		if res := v.Visit(s); res != nil {
			return res
		}
	}

	return nil
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

func (v *Visitor) VisitMeasureExpression(ctx *parser.MeasureExpressionContext) any {
	wire, err := cast[[]int](v.Visit(ctx.GateOperand()))
	if err != nil {
		return err
	}

	v.circuit.Ops = append(v.circuit.Ops, &Measurement{
		Wire: wire,
	})

	return nil
}

func (v *Visitor) VisitMeasureArrowAssignmentStatement(ctx *parser.MeasureArrowAssignmentStatementContext) any {
	// TODO: support assignment
	return v.Visit(ctx.MeasureExpression())
}

func (v *Visitor) VisitGateCallStatement(ctx *parser.GateCallStatementContext) any {
	g, err := cast[string](v.Visit(ctx.Identifier()))
	if err != nil {
		return err
	}

	qargs, err := cast[[]int](v.Visit(ctx.GateOperandList()))
	if err != nil {
		return err
	}

	var cursor int
	var ctrls []int
	for _, mod := range ctx.AllGateModifier() {
		switch {
		case mod.CTRL() != nil:
			n, err := cast[int64](v.Visit(mod))
			if err != nil {
				return err
			}

			for range n {
				ctrls = append(ctrls, cursor)
				cursor++
			}
		}
	}

	ctrlSet := make(map[int]struct{})
	for _, c := range ctrls {
		ctrlSet[c] = struct{}{}
	}

	var targets []int
	for i := range qargs {
		if _, ok := ctrlSet[qargs[i]]; ok {
			continue
		}

		targets = append(targets, qargs[i])
	}

	v.circuit.Ops = append(v.circuit.Ops, &Gate{
		Name:     strings.ToUpper(g),
		Controls: ctrls,
		Targets:  targets,
	})

	return nil
}

func (v *Visitor) VisitGateOperandList(ctx *parser.GateOperandListContext) any {
	var wireIDs []int
	for _, operand := range ctx.AllGateOperand() {
		op, err := cast[[]int](v.Visit(operand))
		if err != nil {
			return err
		}

		wireIDs = append(wireIDs, op...)
	}

	return wireIDs
}

func (v *Visitor) VisitGateOperand(ctx *parser.GateOperandContext) any {
	qargs, err := cast[string](v.Visit(ctx.IndexedIdentifier().Identifier()))
	if err != nil {
		return err
	}

	index, err := cast[[]int64](v.Visit(ctx.IndexedIdentifier()))
	if err != nil {
		return err
	}

	if len(index) == 1 {
		// h q[0]
		wireID := fmt.Sprintf("%s[%d]", qargs, index[0])
		w, ok := v.wire[wireID]
		if !ok {
			return fmt.Errorf("undefined %q", wireID)
		}

		return []int{w}
	}

	// qubit[2] q; h q;
	var wireIDs []int
	for i := 0; ; i++ {
		wireID := fmt.Sprintf("%s[%d]", qargs, i)
		w, ok := v.wire[wireID]
		if !ok {
			break
		}

		wireIDs = append(wireIDs, w)
	}

	if len(wireIDs) == 0 {
		// qubit q; h q;
		wireID, ok := v.wire[qargs]
		if !ok {
			return fmt.Errorf("undefined %q", qargs)
		}

		return []int{wireID}
	}

	return wireIDs
}

func (v *Visitor) VisitGateModifier(ctx *parser.GateModifierContext) any {
	if ctx.Expression() != nil {
		return v.Visit(ctx.Expression())
	}

	return int64(1)
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
	for _, operator := range ctx.AllIndexOperator() {
		op, err := cast[[]any](v.Visit(operator))
		if err != nil {
			return err
		}

		for _, o := range op {
			idx, err := cast[int64](o)
			if err != nil {
				return err
			}

			index = append(index, idx)
		}
	}

	return index
}

func (v *Visitor) VisitQuantumDeclarationStatement(ctx *parser.QuantumDeclarationStatementContext) any {
	wireID, err := cast[string](v.Visit(ctx.Identifier()))
	if err != nil {
		return err
	}

	size, err := cast[int64](v.Visit(ctx.QubitType()))
	if err != nil {
		return err
	}

	wireIDs := []string{wireID}
	if size > 1 {
		wireIDs = make([]string, 0, size)
		for i := range size {
			wireIDs = append(wireIDs, fmt.Sprintf("%s[%d]", wireID, i))
		}
	}

	for _, wireID := range wireIDs {
		if err := v.Add(wireID); err != nil {
			return err
		}
	}

	return nil
}

func (v *Visitor) VisitQubitType(ctx *parser.QubitTypeContext) any {
	if ctx.Designator() != nil {
		val, err := cast[int64](v.Visit(ctx.Designator()))
		if err != nil {
			return err
		}

		return val
	}

	return int64(1)
}

func (v *Visitor) VisitClassicalDeclarationStatement(ctx *parser.ClassicalDeclarationStatementContext) any {
	switch {
	case ctx.ScalarType().BIT() != nil:
		wireID, err := cast[string](v.Visit(ctx.Identifier()))
		if err != nil {
			return err
		}

		switch {
		case ctx.ScalarType().Designator() != nil:
			size, err := cast[int64](v.Visit(ctx.ScalarType()))
			if err != nil {
				return err
			}

			for i := range size {
				if err := v.Add(fmt.Sprintf("%s[%d]", wireID, i)); err != nil {
					return err
				}
			}

			return nil
		default:
			if err := v.Add(wireID); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func (v *Visitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) any {
	switch {
	case ctx.DecimalIntegerLiteral() != nil:
		s, err := cast[string](v.Visit(ctx.DecimalIntegerLiteral()))
		if err != nil {
			return err
		}

		lit, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int %q: %w", s, err)
		}

		return lit
	default:
		return fmt.Errorf("unsupported literal %q", ctx.GetText())
	}
}

func (v *Visitor) VisitScalarType(ctx *parser.ScalarTypeContext) any {
	if ctx.Designator() != nil {
		val, err := cast[int64](v.Visit(ctx.Designator()))
		if err != nil {
			return err
		}

		return val
	}

	return int64(1)
}

func (v *Visitor) VisitDesignator(ctx *parser.DesignatorContext) any {
	return v.Visit(ctx.Expression())
}

func cast[T any](result any) (T, error) {
	if err, ok := result.(error); ok && err != nil {
		var zero T
		return zero, err
	}

	resultT, ok := result.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("unexpected type %T", result)
	}

	return resultT, nil
}
