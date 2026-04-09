package visitor

import (
	"fmt"
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

// Pow2x2 returns u**p for 2x2 matrix u and float p.
// If p is negative, returns (u-dagger)**p.
func Pow2x2(u *matrix.Matrix, p float64) (*matrix.Matrix, error) {
	if u.Rows != 2 || u.Cols != 2 {
		return nil, fmt.Errorf("unsupported matrix size %dx%d", u.Rows, u.Cols)
	}

	if p < 0 {
		p, u = -p, u.Dagger()
	}

	// TODO: support float type
	return matrix.ApplyN(u, int(p)), nil
}
