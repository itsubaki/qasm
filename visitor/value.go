package visitor

import (
	"fmt"
	"math"
)

type Value struct {
	v any
}

func (v Value) Value() any {
	return v.v
}

func NewValue(v any) Value {
	return Value{v: v}
}

func Promote(a, b Value) (Value, Value, error) {
	switch av := a.v.(type) {
	case int:
		switch b.v.(type) {
		case int:
			return a, b, nil
		case int64:
			return NewValue(int64(av)), b, nil
		case float64:
			return NewValue(float64(av)), b, nil
		}
	case int64:
		switch bv := b.v.(type) {
		case int:
			return a, NewValue(int64(bv)), nil
		case int64:
			return a, b, nil
		case float64:
			return NewValue(float64(av)), b, nil
		}
	case float64:
		switch bv := b.v.(type) {
		case int:
			return a, NewValue(float64(bv)), nil
		case int64:
			return a, NewValue(float64(bv)), nil
		case float64:
			return a, b, nil
		}
	case bool:
		switch b.v.(type) {
		case bool:
			return a, b, nil
		}
	}

	return Value{}, Value{}, fmt.Errorf("unexpected %T and %T", a.v, b.v)
}

func (v Value) Add(w Value) (Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return Value{}, err
	}

	switch v := a.v.(type) {
	case int:
		return NewValue(v + b.v.(int)), nil
	case int64:
		return NewValue(v + b.v.(int64)), nil
	case float64:
		return NewValue(v + b.v.(float64)), nil
	}

	return Value{}, fmt.Errorf("unexpected %T + %T", a.v, b.v)
}

func (v Value) Sub(w Value) (Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return Value{}, err
	}

	switch v := a.v.(type) {
	case int:
		return NewValue(v - b.v.(int)), nil
	case int64:
		return NewValue(v - b.v.(int64)), nil
	case float64:
		return NewValue(v - b.v.(float64)), nil
	}

	return Value{}, fmt.Errorf("unexpected %T - %T", a.v, b.v)
}

func (v Value) Mul(w Value) (Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return Value{}, err
	}

	switch v := a.v.(type) {
	case int:
		return NewValue(v * b.v.(int)), nil
	case int64:
		return NewValue(v * b.v.(int64)), nil
	case float64:
		return NewValue(v * b.v.(float64)), nil
	}

	return Value{}, fmt.Errorf("unexpected %T * %T", a.v, b.v)
}

func (v Value) Div(w Value) (Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return Value{}, err
	}

	switch v := a.v.(type) {
	case int:
		return NewValue(v / b.v.(int)), nil
	case int64:
		return NewValue(v / b.v.(int64)), nil
	case float64:
		return NewValue(v / b.v.(float64)), nil
	}

	return Value{}, fmt.Errorf("unexpected %T / %T", a.v, b.v)
}

func (v Value) Mod(w Value) (Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return Value{}, err
	}

	switch v := a.v.(type) {
	case int:
		return NewValue(v % b.v.(int)), nil
	case int64:
		return NewValue(v % b.v.(int64)), nil
	}

	return Value{}, fmt.Errorf("unexpected %T %% %T", a.v, b.v)
}

func (v Value) Eq(w Value) (Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return Value{}, err
	}

	switch v := a.v.(type) {
	case int:
		return NewValue(v == b.v.(int)), nil
	case int64:
		return NewValue(v == b.v.(int64)), nil
	case float64:
		return NewValue(isClose(v, b.v.(float64))), nil
	case bool:
		return NewValue(v == b.v.(bool)), nil
	}

	return Value{}, fmt.Errorf("unexpected %T == %T", a.v, b.v)
}

func (v Value) NotEq(w Value) (Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return Value{}, err
	}

	switch v := a.v.(type) {
	case int:
		return NewValue(v != b.v.(int)), nil
	case int64:
		return NewValue(v != b.v.(int64)), nil
	case float64:
		return NewValue(!isClose(v, b.v.(float64))), nil
	case bool:
		return NewValue(v != b.v.(bool)), nil
	}

	return Value{}, fmt.Errorf("unexpected %T != %T", a.v, b.v)
}

func isClose(a, b float64) bool {
	return math.Abs(a-b) <= 1e-8+1e-5*math.Max(math.Abs(a), math.Abs(b))
}
