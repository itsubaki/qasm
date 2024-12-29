package visitor

import "math"

var BuiltinConst = map[string]float64{
	"pi":    math.Pi,
	"tau":   2 * math.Pi,
	"euler": math.E,
	"π":     math.Pi,
	"τ":     2 * math.Pi,
	"ℇ":     math.E,
}
