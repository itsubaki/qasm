package svg

import (
	"fmt"
	"strconv"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/environ"
	"github.com/itsubaki/qasm/gen/parser"
)

type Visitor struct {
	*parser.Baseqasm3ParserVisitor
	circuit *Circuit
	env     *environ.Environ
	wire    map[string]int
}

func NewVisitor() *Visitor {
	return &Visitor{
		Baseqasm3ParserVisitor: &parser.Baseqasm3ParserVisitor{},
		circuit:                &Circuit{},
		env:                    environ.New(),
		wire:                   make(map[string]int),
	}
}

func (v *Visitor) Build(tree antlr.ParseTree) (*Circuit, error) {
	if err, ok := v.Visit(tree).(error); ok && err != nil {
		return nil, err
	}

	return v.circuit, nil
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

func (v *Visitor) VisitGateCallStatement(ctx *parser.GateCallStatementContext) any {
	// qargs
	qargs, err := unwrap[[]int](v.Visit(ctx.GateOperandList()))
	if err != nil {
		return err
	}

	// gate
	g, err := unwrap[string](v.Visit(ctx.Identifier()))
	if err != nil {
		return err
	}

	// TODO: modifier
	// TODO: support control modifier
	v.circuit.Ops = append(v.circuit.Ops, &Gate{
		Name:    g,
		Targets: qargs,
	})

	return nil
}

func (v *Visitor) VisitGateOperandList(ctx *parser.GateOperandListContext) any {
	var list []int
	for _, operand := range ctx.AllGateOperand() {
		op, err := unwrap[int](v.Visit(operand))
		if err != nil {
			return err
		}

		list = append(list, op)
	}

	return list
}

func (v *Visitor) VisitGateOperand(ctx *parser.GateOperandContext) any {
	qargs, err := unwrap[string](v.Visit(ctx.IndexedIdentifier().Identifier()))
	if err != nil {
		return err
	}

	// q or q[0]
	index, err := unwrap[[]int64](v.Visit(ctx.IndexedIdentifier()))
	if err != nil {
		return err
	}

	if len(index) > 0 {
		// h q[0];
		qargs = fmt.Sprintf("%s[%d]", qargs, index[0])
	}

	return v.wire[qargs]
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
		op, err := unwrap[[]any](v.Visit(operator))
		if err != nil {
			return err
		}

		for _, o := range op {
			idx, err := unwrap[int64](o)
			if err != nil {
				return err
			}

			index = append(index, idx)
		}
	}

	return index
}

func (v *Visitor) VisitQuantumDeclarationStatement(ctx *parser.QuantumDeclarationStatementContext) any {
	id, err := unwrap[string](v.Visit(ctx.Identifier()))
	if err != nil {
		return err
	}

	if _, ok := v.wire[id]; ok {
		return fmt.Errorf("wire %q already exists", id)
	}

	size, err := unwrap[int64](v.Visit(ctx.QubitType()))
	if err != nil {
		return err
	}

	ids := []string{id}
	if size > 1 {
		ids = make([]string, 0, size)
		for i := range size {
			ids = append(ids, fmt.Sprintf("%s[%d]", id, i))
		}
	}

	for _, id := range ids {
		v.wire[id] = len(v.circuit.Wires)
		v.circuit.Wires = append(v.circuit.Wires, Wire{
			Name: id,
		})
	}

	return nil
}

func (v *Visitor) VisitQubitType(ctx *parser.QubitTypeContext) any {
	if ctx.Designator() != nil {
		val, err := unwrap[int64](v.Visit(ctx.Designator()))
		if err != nil {
			return err
		}

		return val
	}

	return int64(1)
}

func (v *Visitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) any {
	switch {
	case ctx.DecimalIntegerLiteral() != nil:
		s, err := unwrap[string](v.Visit(ctx.DecimalIntegerLiteral()))
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

func (v *Visitor) VisitDesignator(ctx *parser.DesignatorContext) any {
	return v.Visit(ctx.Expression())
}

func unwrap[T any](result any) (T, error) {
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
