package value

import (
	"fmt"
	"math"
)

// Pow returns a**r, the base-a exponential of r.
func Pow[T int | int64 | float64](a, r T) T {
	switch left := any(a).(type) {
	case int:
		return T(pow(left, any(r).(int)))
	case int64:
		return T(pow(left, any(r).(int64)))
	case float64:
		return T(math.Pow(left, any(r).(float64)))
	}

	panic(fmt.Sprintf("unexpected pow types %T and %T", a, r))
}

func pow[T int | int64](a, r T) T {
	if a == 0 {
		return 0
	}

	if r == 0 {
		return 1
	}

	// exponentiation by squaring
	p := T(1)
	for r > 0 {
		if r&1 == 1 {
			p = p * a
		}

		a = a * a
		r >>= 1
	}

	return p
}
