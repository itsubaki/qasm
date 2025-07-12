package visitor

import "github.com/itsubaki/q/math/matrix"

func Pow(u *matrix.Matrix, p float64) *matrix.Matrix {
	if p < 0 {
		p, u = -p, u.Dagger()
	}

	// TODO: support float type
	return matrix.ApplyN(u, int(p))
}
