package value

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

func New(v any) *Value {
	return &Value{v: v}
}

func Promote(a, b *Value) (*Value, *Value, error) {
	switch av := a.v.(type) {
	case int:
		switch b.v.(type) {
		case int:
			return a, b, nil
		case int64:
			return New(int64(av)), b, nil
		case float64:
			return New(float64(av)), b, nil
		}
	case int64:
		switch bv := b.v.(type) {
		case int:
			return a, New(int64(bv)), nil
		case int64:
			return a, b, nil
		case float64:
			return New(float64(av)), b, nil
		}
	case float64:
		switch bv := b.v.(type) {
		case int:
			return a, New(float64(bv)), nil
		case int64:
			return a, New(float64(bv)), nil
		case float64:
			return a, b, nil
		}
	case uint:
		switch b.v.(type) {
		case uint:
			return a, b, nil
		}
	case bool:
		switch b.v.(type) {
		case bool:
			return a, b, nil
		}
	}

	return nil, nil, fmt.Errorf("unexpected %T and %T", a.v, b.v)
}

func (v *Value) Add(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch v := a.v.(type) {
	case int:
		return New(v + b.v.(int)), nil
	case int64:
		return New(v + b.v.(int64)), nil
	case float64:
		return New(v + b.v.(float64)), nil
	}

	return nil, fmt.Errorf("unexpected %T + %T", a.v, b.v)
}

func (v *Value) Sub(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch v := a.v.(type) {
	case int:
		return New(v - b.v.(int)), nil
	case int64:
		return New(v - b.v.(int64)), nil
	case float64:
		return New(v - b.v.(float64)), nil
	}

	return nil, fmt.Errorf("unexpected %T - %T", a.v, b.v)
}

func (v *Value) Mul(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch v := a.v.(type) {
	case int:
		return New(v * b.v.(int)), nil
	case int64:
		return New(v * b.v.(int64)), nil
	case float64:
		return New(v * b.v.(float64)), nil
	}

	return nil, fmt.Errorf("unexpected %T * %T", a.v, b.v)
}

func (v *Value) Div(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch v := a.v.(type) {
	case int:
		return New(v / b.v.(int)), nil
	case int64:
		return New(v / b.v.(int64)), nil
	case float64:
		return New(v / b.v.(float64)), nil
	}

	return nil, fmt.Errorf("unexpected %T / %T", a.v, b.v)
}

func (v *Value) Mod(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch v := a.v.(type) {
	case int:
		return New(v % b.v.(int)), nil
	case int64:
		return New(v % b.v.(int64)), nil
	}

	return nil, fmt.Errorf("unexpected %T %% %T", a.v, b.v)
}

func (v *Value) Eq(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch v := a.v.(type) {
	case int:
		return New(v == b.v.(int)), nil
	case int64:
		return New(v == b.v.(int64)), nil
	case float64:
		return New(isClose(v, b.v.(float64))), nil
	case bool:
		return New(v == b.v.(bool)), nil
	case uint:
		return New(v == b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T == %T", a.v, b.v)
}

func (v *Value) NotEq(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch v := a.v.(type) {
	case int:
		return New(v != b.v.(int)), nil
	case int64:
		return New(v != b.v.(int64)), nil
	case float64:
		return New(!isClose(v, b.v.(float64))), nil
	case bool:
		return New(v != b.v.(bool)), nil
	}

	return nil, fmt.Errorf("unexpected %T != %T", a.v, b.v)
}

func isClose(a, b float64) bool {
	return math.Abs(a-b) <= 1e-8+1e-5*math.Max(math.Abs(a), math.Abs(b))
}

func (v *Value) LessThan(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch av := a.v.(type) {
	case int:
		return New(av < b.v.(int)), nil
	case int64:
		return New(av < b.v.(int64)), nil
	case float64:
		return New(av < b.v.(float64)), nil
	}

	return nil, fmt.Errorf("unexpected %T < %T", a.v, b.v)
}

func (v *Value) LessThanOrEqual(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch av := a.v.(type) {
	case int:
		return New(av <= b.v.(int)), nil
	case int64:
		return New(av <= b.v.(int64)), nil
	case float64:
		bv := b.v.(float64)
		return New(av < bv || isClose(av, bv)), nil
	}

	return nil, fmt.Errorf("unexpected %T <= %T", a.v, b.v)
}

func (v *Value) GreaterThan(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch av := a.v.(type) {
	case int:
		return New(av > b.v.(int)), nil
	case int64:
		return New(av > b.v.(int64)), nil
	case float64:
		return New(av > b.v.(float64)), nil
	}

	return nil, fmt.Errorf("unexpected %T > %T", a.v, b.v)
}

func (v *Value) GreaterThanOrEqual(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch av := a.v.(type) {
	case int:
		return New(av >= b.v.(int)), nil
	case int64:
		return New(av >= b.v.(int64)), nil
	case float64:
		bv := b.v.(float64)
		return New(av > bv || isClose(av, bv)), nil
	}

	return nil, fmt.Errorf("unexpected %T >= %T", a.v, b.v)
}

func (v *Value) Int() (*Value, error) {
	switch v := v.v.(type) {
	case int:
		return New(v), nil
	case int64:
		return New(int(v)), nil
	case float64:
		return New(int(v)), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) Int64() (*Value, error) {
	switch v := v.v.(type) {
	case int:
		return New(int64(v)), nil
	case int64:
		return New(int64(v)), nil
	case float64:
		return New(int64(v)), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) UInt() (*Value, error) {
	switch v := v.v.(type) {
	case int:
		return New(uint(v)), nil
	case int64:
		return New(uint(v)), nil
	case float64:
		return New(uint(v)), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) Float64() (*Value, error) {
	switch v := v.v.(type) {
	case int:
		return New(float64(v)), nil
	case int64:
		return New(float64(v)), nil
	case float64:
		return New(v), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) Negative() (*Value, error) {
	switch v := v.v.(type) {
	case int:
		return New(-v), nil
	case int64:
		return New(-v), nil
	case float64:
		return New(-v), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) BitNot() (*Value, error) {
	switch v := v.v.(type) {
	case int:
		return New(^v), nil
	case int64:
		return New(^v), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) BoolNot() (*Value, error) {
	switch v := v.v.(type) {
	case bool:
		return New(!v), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}
