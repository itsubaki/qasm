package visitor

import (
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
)

// AddControlled returns a controlled-u gate with control bit.
// u is a (2**n x 2**n) unitary matrix and returns a (2**n x 2**n) matrix.
func AddControlled(u matrix.Matrix, c int) matrix.Matrix {
	d, _ := u.Dimension()
	n := number.Log2(d)
	g := gate.I(n)

	for i := 0; i < d; i++ {
		if (i>>(n-1-c))&1 == 0 {
			continue
		}

		for j := 0; j < d; j++ {
			g[i][j] = u[i][j]
		}
	}

	return g
}
