package visitor

import (
	"math"
	"math/cmplx"
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

	if p == math.Trunc(p) && p <= math.MaxInt {
		return matrix.ApplyN(u, int(p))
	}

	return matPowFloat(u, p)
}

// matPowFloat computes u^p for a unitary matrix u and non-integer real exponent p >= 0.
// For 2x2 matrices it uses the analytic Sylvester formula via eigenvalues derived
// from the characteristic polynomial (trace and determinant).
// For matrices larger than 2x2 it falls back to the nearest integer power.
func matPowFloat(u *matrix.Matrix, p float64) *matrix.Matrix {
	rows, _ := u.Dim()
	if rows != 2 {
		// For matrices larger than 2x2, fall back to the nearest integer power.
		n := math.Round(p)
		if n > math.MaxInt {
			n = math.MaxInt
		}
		return matrix.ApplyN(u, int(n))
	}

	// Characteristic polynomial for 2x2: λ² - tr(u)λ + det(u) = 0
	t := u.At(0, 0) + u.At(1, 1)                              // trace
	d := u.At(0, 0)*u.At(1, 1) - u.At(0, 1)*u.At(1, 0)       // determinant
	disc := t*t - 4*d
	sqrtDisc := cmplx.Sqrt(disc)
	λ1 := (t + sqrtDisc) / 2
	λ2 := (t - sqrtDisc) / 2

	cp := complex(p, 0)

	if cmplx.Abs(λ1-λ2) < 1e-10 {
		// Equal eigenvalues: u is a scalar multiple of identity, u^p = λ^p·I.
		λp := cmplx.Pow(λ1, cp)
		return matrix.New(
			[]complex128{λp, 0},
			[]complex128{0, λp},
		)
	}

	// Sylvester's formula: u^p = a·u + b·I
	// where a = (λ1^p - λ2^p)/(λ1-λ2) and b = (λ2^p·λ1 - λ1^p·λ2)/(λ1-λ2).
	λ1p := cmplx.Pow(λ1, cp)
	λ2p := cmplx.Pow(λ2, cp)
	diff := λ1 - λ2
	a := (λ1p - λ2p) / diff
	b := (λ2p*λ1 - λ1p*λ2) / diff

	return matrix.New(
		[]complex128{a*u.At(0, 0) + b, a*u.At(0, 1)},
		[]complex128{a*u.At(1, 0), a*u.At(1, 1) + b},
	)
}
