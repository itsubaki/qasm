package value

import "math"

// Pow returns a**r, the base-a exponential of r.
func Pow[T int | int64 | float64](a, r T) T {
	if left, ok := any(a).(int); ok {
		return T(powInt(left, any(r).(int)))
	}

	if left, ok := any(a).(int64); ok {
		return T(powInt(left, any(r).(int64)))
	}

	return T(math.Pow(float64(a), float64(r)))
}

// powInt returns a**r for integer values using exponentiation by squaring.
func powInt[T int | int64](a, r T) T {
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
