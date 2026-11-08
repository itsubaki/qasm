package value

import "fmt"

func (v *Value) BitAnd(w *Value) (*Value, error) {
	if _, ok := v.v.([]bool); ok {
		return v.bitwiseBits(w, func(a, b bool) bool { return a && b })
	}

	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left & b.v.(int)), nil
	case int64:
		return New(left & b.v.(int64)), nil
	case uint:
		return New(left & b.v.(uint)), nil
	case bool:
		return New(left && b.v.(bool)), nil
	}

	return nil, fmt.Errorf("unexpected %T & %T", a.v, b.v)
}

func (v *Value) BitOr(w *Value) (*Value, error) {
	if _, ok := v.v.([]bool); ok {
		return v.bitwiseBits(w, func(a, b bool) bool { return a || b })
	}

	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left | b.v.(int)), nil
	case int64:
		return New(left | b.v.(int64)), nil
	case uint:
		return New(left | b.v.(uint)), nil
	case bool:
		return New(left || b.v.(bool)), nil
	}

	return nil, fmt.Errorf("unexpected %T | %T", a.v, b.v)
}

func (v *Value) BitXor(w *Value) (*Value, error) {
	if _, ok := v.v.([]bool); ok {
		return v.bitwiseBits(w, func(a, b bool) bool { return a != b })
	}

	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left ^ b.v.(int)), nil
	case int64:
		return New(left ^ b.v.(int64)), nil
	case uint:
		return New(left ^ b.v.(uint)), nil
	case bool:
		return New(left != b.v.(bool)), nil
	}

	return nil, fmt.Errorf("unexpected %T ^ %T", a.v, b.v)
}

func (v *Value) bitwiseBits(w *Value, op func(bool, bool) bool) (*Value, error) {
	left := v.v.([]bool)
	right, ok := w.v.([]bool)
	if !ok {
		return nil, fmt.Errorf("unexpected %T and %T", v.v, w.v)
	}

	if len(left) != len(right) {
		return nil, fmt.Errorf("bitwise operands have different sizes: %d and %d", len(left), len(right))
	}

	result := make([]bool, len(left))
	for i := range result {
		result[i] = op(left[i], right[i])
	}

	return New(result), nil
}

func (v *Value) ShiftLeft(w *Value) (*Value, error) {
	return v.bitshift(w, true)
}

func (v *Value) ShiftRight(w *Value) (*Value, error) {
	return v.bitshift(w, false)
}

func (v *Value) bitshift(w *Value, left bool) (*Value, error) {
	var count uint64
	switch n := w.v.(type) {
	case int:
		if n < 0 {
			return nil, fmt.Errorf("negative shift count: %d", n)
		}
		count = uint64(n)
	case int64:
		if n < 0 {
			return nil, fmt.Errorf("negative shift count: %d", n)
		}
		count = uint64(n)
	case uint:
		count = uint64(n)
	default:
		return nil, fmt.Errorf("unexpected shift count type: %T", w.v)
	}

	switch val := v.v.(type) {
	case int:
		if left {
			return New(val << count), nil
		}
		return New(val >> count), nil
	case int64:
		if left {
			return New(val << count), nil
		}
		return New(val >> count), nil
	case uint:
		if left {
			return New(val << count), nil
		}
		return New(val >> count), nil
	case bool:
		return New(val && count == 0), nil
	case []bool:
		result := make([]bool, len(val))
		if count < uint64(len(val)) {
			// Bit strings are stored from the most significant bit to the least.
			if left {
				copy(result, val[int(count):])
			} else {
				copy(result[int(count):], val)
			}
		}

		return New(result), nil
	}

	return nil, fmt.Errorf("unexpected shift operand type: %T", v.v)
}
