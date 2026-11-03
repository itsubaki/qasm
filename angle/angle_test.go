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

func TestAngle(t *testing.T) {
	cases := []struct {
		bits       uint
		radian     float64
		wantK      uint
		wantBits   string
		wantRadian float64
	}{
		{bits: 3, radian: 2*math.Pi - 0.01, wantK: 0, wantBits: "000", wantRadian: 0},
		{bits: 3, radian: -0.01, wantK: 0, wantBits: "000", wantRadian: 0},
		{bits: 3, radian: -math.SmallestNonzeroFloat64, wantK: 0, wantBits: "000", wantRadian: 0},
	}

	for _, c := range cases {
		a := angle.New(c.bits, c.radian)
		if a.K != c.wantK || a.BitString != c.wantBits || a.Radian() != c.wantRadian {
			t.Errorf("bits=%d radian=%v: got=%v, want K=%d bits=%s radian=%v", c.bits, c.radian, a, c.wantK, c.wantBits, c.wantRadian)
		}
	}
}
