package angle

import (
	"fmt"
	"math"
)

type Angle struct {
	Bits      uint
	BitString string
	K         uint
}

func New(bits uint, radian float64) *Angle {
	k := math.Ldexp(Mod2Pi(radian)/(2*math.Pi), int(bits))
	kr := uint(math.Round(k))
	return &Angle{
		Bits:      bits,
		BitString: fmt.Sprintf("%0*b", bits, kr),
		K:         kr,
	}
}

func (v *Angle) Radian() float64 {
	return 2 * math.Pi * math.Ldexp(float64(v.K), -int(v.Bits))
}

func (v *Angle) String() string {
	return fmt.Sprintf("%v(%d,%s)", v.Radian(), v.K, v.BitString)
}

func Mod2Pi(radian float64) float64 {
	mod := math.Mod(radian, 2*math.Pi)
	if mod < 0 {
		mod += 2 * math.Pi
	}

	return mod
}
