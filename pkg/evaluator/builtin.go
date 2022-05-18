package evaluator

import (
	"github.com/itsubaki/q"
	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/qasm/pkg/lexer"
)

func Builtin(g lexer.Token, p []float64) (matrix.Matrix, bool) {
	switch g {
	case lexer.U:
		return gate.U(p[0], p[1], p[2]), true
	case lexer.X:
		return gate.X(), true
	case lexer.Y:
		return gate.Y(), true
	case lexer.Z:
		return gate.Z(), true
	case lexer.H:
		return gate.H(), true
	case lexer.T:
		return gate.T(), true
	case lexer.S:
		return gate.S(), true
	}

	return nil, false
}

func BuiltinApply(qsim *q.Q, g lexer.Token, p []float64, qargs [][]q.Qubit) bool {
	switch g {
	case lexer.QFT:
		qsim.QFT(flatten(qargs)...)
		return true
	case lexer.IQFT:
		qsim.InvQFT(flatten(qargs)...)
		return true
	case lexer.CMODEXP2:
		qsim.CModExp2(int(p[0]), int(p[1]), qargs[0], qargs[1])
		return true
	}

	return false
}

func flatten(qargs [][]q.Qubit) []q.Qubit {
	var out []q.Qubit
	for _, q := range qargs {
		out = append(out, q...)
	}

	return out
}
