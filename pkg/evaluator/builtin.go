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

func flatten(qargs [][]q.Qubit) []q.Qubit {
	var out []q.Qubit
	for _, q := range qargs {
		out = append(out, q...)
	}

	return out
}
