package visitor

import (
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
)

// AddControlled returns a controlled-u gate with control bit.
// u is a (2**n x 2**n) unitary matrix and returns a (2**n x 2**n) matrix.
func AddControlled(u matrix.Matrix, c []int) matrix.Matrix {
	d, _ := u.Dimension()
	n := number.Log2(d)
	g := gate.I(n)

	var mask int
	for _, bit := range c {
		mask |= (1 << (n - 1 - bit))
	}

	for i := 0; i < d; i++ {
		if (i & mask) != mask {
			continue
		}

		for j := 0; j < d; j++ {
			g[i][j] = u[i][j]
		}
	}

	return g
}

func Pow(u matrix.Matrix, p float64) matrix.Matrix {
	// TODO: support float type
	return matrix.ApplyN(u, int(p))
}
