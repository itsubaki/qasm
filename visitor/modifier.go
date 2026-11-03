package visitor

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/qasm/gen/parser"
)

// NotImplementedOrder returns an error if the order of modifiers is not implemented.
func NotImplementedOrder(ctx *parser.GateCallStatementContext) error {
	var pow bool
	for _, mod := range ctx.AllGateModifier() {
		switch {
		case mod.POW() != nil:
			pow = true
		case mod.CTRL() != nil:
			if !pow {
				continue
			}

			return fmt.Errorf("pow applied after ctrl: %w", ErrNotImplemented)
		case mod.NEGCTRL() != nil:
			if !pow {
				continue
			}

			return fmt.Errorf("pow applied after negctrl: %w", ErrNotImplemented)
		}
	}

	return nil
}

// HasControlModifier returns true if the gate call statement has control modifiers.
func HasControlModifier(ctx parser.IGateCallStatementContext) bool {
	for _, mod := range ctx.AllGateModifier() {
		if mod.CTRL() != nil || mod.NEGCTRL() != nil {
			return true
		}
	}

	return false
}

// Pow2x2 returns u^p for 2x2 matrix u and float p.
func Pow2x2(u *matrix.Matrix, p float64, tol ...float64) *matrix.Matrix {
	// SU
	det := u.At(0, 0)*u.At(1, 1) - u.At(0, 1)*u.At(1, 0)
	phase := cmplx.Sqrt(det)
	su := u.Mul(1 / phase)

	// theta
	tr := real(su.Trace())
	cosTheta := max(min(tr/2, 1), -1)
	theta := math.Acos(cosTheta)
	sinTheta := math.Sin(theta)

	// phase^p
	phaseP := cmplx.Pow(phase, complex(p, 0))

	// if sin(theta) is close to zero, su is close to I or -I.
	if epsilon.IsZeroF64(sinTheta, tol...) {
		idp := matrix.Identity(2).Mul(phaseP)
		if cosTheta > 0 {
			return idp
		}

		return idp.Mul(cmplx.Exp(complex(0, p*math.Pi)))
	}

	id := matrix.Identity(2)
	a := su.Sub(su.Dagger()).Mul(complex(0, -0.5/sinTheta))

	// p*theta
	cos := complex(math.Cos(p*theta), 0)
	sin := complex(0, math.Sin(p*theta))
	return id.Mul(cos).Add(a.Mul(sin)).Mul(phaseP)
}

type QubitOperand struct {
	Qubits  []q.Qubit
	IsArray bool
}

// broadcast expands operands into rows, repeating scalar operands in each row.
func broadcast(args []QubitOperand) ([][]q.Qubit, error) {
	size := 0
	for _, arg := range args {
		if len(arg.Qubits) == 0 {
			return nil, fmt.Errorf("gate operand has size 0")
		}

		if !arg.IsArray {
			continue
		}

		if size != 0 && len(arg.Qubits) != size {
			return nil, fmt.Errorf("gate operand size mismatch: %d and %d", size, len(arg.Qubits))
		}

		size = len(arg.Qubits)
	}

	if size == 0 {
		size = 1
	}

	rows := make([][]q.Qubit, size)
	for i := range rows {
		rows[i] = make([]q.Qubit, len(args))
		for j, arg := range args {
			idx := 0
			if arg.IsArray {
				idx = i
			}

			rows[i][j] = arg.Qubits[idx]
		}
	}

	return rows, nil
}
