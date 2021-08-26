package register

import "math"

type Const map[string]float64

func NewConst() Const {
	c := make(map[string]float64)
	c["pi"] = math.Pi
	c["tau"] = math.Pi * 2
	c["euler"] = math.E

	return c
}
