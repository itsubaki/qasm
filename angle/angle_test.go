package angle_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/qasm/angle"
)

func ExampleAngle() {
	a := angle.New(4, math.Pi)
	fmt.Println(a.Radian())

	fmt.Println("bits:", a.Bits)
	fmt.Println("bit string:", a.BitString)
	fmt.Println("k:", a.K)
	fmt.Println(a)

	// Output:
	// 3.141592653589793
	// bits: 4
	// bit string: 1000
	// k: 8
	// 3.141592653589793(8,1000)
}

func TestMod2Pi(t *testing.T) {
	cases := []struct {
		radian float64
		want   float64
	}{
		{
			radian: 3 * math.Pi,
			want:   math.Pi,
		},
		{
			radian: -math.Pi / 2,
			want:   3 * math.Pi / 2,
		},
	}

	for _, c := range cases {
		got := angle.Mod2Pi(c.radian)
		if got != c.want {
			t.Errorf("Mod2Pi(%f) = %f, want %f", c.radian, got, c.want)
		}
	}
}
