package visitor

import (
	"slices"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/qasm/gen/parser"
)

func HasControlModifier(ctx parser.IGateCallStatementContext) bool {
	for _, mod := range ctx.AllGateModifier() {
		if mod.CTRL() != nil || mod.NEGCTRL() != nil {
			return true
		}
	}

	return false
}

func ReversedModifier(ctx parser.IGateCallStatementContext) []parser.IGateModifierContext {
	modifier := make([]parser.IGateModifierContext, len(ctx.AllGateModifier()))
	copy(modifier, ctx.AllGateModifier())
	slices.Reverse(modifier)

	return modifier
}

func Pow(u *matrix.Matrix, p float64) *matrix.Matrix {
	if p < 0 {
		p, u = -p, u.Dagger()
	}

	// TODO: support float type
	return matrix.ApplyN(u, int(p))
}
