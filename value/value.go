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
	switch left := a.v.(type) {
	case int:
		switch b.v.(type) {
		case int:
			return a, b, nil
		case int64:
			return New(int64(left)), b, nil
		case float64:
			return New(float64(left)), b, nil
		}
	case int64:
		switch right := b.v.(type) {
		case int:
			return a, New(int64(right)), nil
		case int64:
			return a, b, nil
		case float64:
			return New(float64(left)), b, nil
		}
	case float64:
		switch right := b.v.(type) {
		case int:
			return a, New(float64(right)), nil
		case int64:
			return a, New(float64(right)), nil
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

	switch left := a.v.(type) {
	case int:
		return New(left + b.v.(int)), nil
	case int64:
		return New(left + b.v.(int64)), nil
	case float64:
		return New(left + b.v.(float64)), nil
	case uint:
		return New(left + b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T + %T", a.v, b.v)
}

func (v *Value) Sub(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left - b.v.(int)), nil
	case int64:
		return New(left - b.v.(int64)), nil
	case float64:
		return New(left - b.v.(float64)), nil
	case uint:
		return New(left - b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T - %T", a.v, b.v)
}

func (v *Value) Mul(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left * b.v.(int)), nil
	case int64:
		return New(left * b.v.(int64)), nil
	case float64:
		return New(left * b.v.(float64)), nil
	case uint:
		return New(left * b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T * %T", a.v, b.v)
}

func (v *Value) Div(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left / b.v.(int)), nil
	case int64:
		return New(left / b.v.(int64)), nil
	case float64:
		return New(left / b.v.(float64)), nil
	case uint:
		return New(left / b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T / %T", a.v, b.v)
}

func (v *Value) Mod(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left % b.v.(int)), nil
	case int64:
		return New(left % b.v.(int64)), nil
	case uint:
		return New(left % b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T %% %T", a.v, b.v)
}

func (v *Value) Pow(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(Pow(left, b.v.(int))), nil
	case int64:
		return New(Pow(left, b.v.(int64))), nil
	case float64:
		return New(math.Pow(left, b.v.(float64))), nil
	}

	return nil, fmt.Errorf("unexpected %T ** %T", a.v, b.v)
}

func (v *Value) Eq(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left == b.v.(int)), nil
	case int64:
		return New(left == b.v.(int64)), nil
	case float64:
		return New(isClose(left, b.v.(float64))), nil
	case bool:
		return New(left == b.v.(bool)), nil
	case uint:
		return New(left == b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T == %T", a.v, b.v)
}

func (v *Value) NotEq(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left != b.v.(int)), nil
	case int64:
		return New(left != b.v.(int64)), nil
	case float64:
		return New(!isClose(left, b.v.(float64))), nil
	case bool:
		return New(left != b.v.(bool)), nil
	case uint:
		return New(left != b.v.(uint)), nil
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

	switch left := a.v.(type) {
	case int:
		return New(left < b.v.(int)), nil
	case int64:
		return New(left < b.v.(int64)), nil
	case float64:
		return New(left < b.v.(float64)), nil
	case uint:
		return New(left < b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T < %T", a.v, b.v)
}

func (v *Value) LessThanOrEqual(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left <= b.v.(int)), nil
	case int64:
		return New(left <= b.v.(int64)), nil
	case float64:
		return New(left < b.v.(float64) || isClose(left, b.v.(float64))), nil
	case uint:
		return New(left <= b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T <= %T", a.v, b.v)
}

func (v *Value) GreaterThan(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left > b.v.(int)), nil
	case int64:
		return New(left > b.v.(int64)), nil
	case float64:
		return New(left > b.v.(float64)), nil
	case uint:
		return New(left > b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T > %T", a.v, b.v)
}

func (v *Value) GreaterThanOrEqual(w *Value) (*Value, error) {
	a, b, err := Promote(v, w)
	if err != nil {
		return nil, err
	}

	switch left := a.v.(type) {
	case int:
		return New(left >= b.v.(int)), nil
	case int64:
		return New(left >= b.v.(int64)), nil
	case float64:
		return New(left > b.v.(float64) || isClose(left, b.v.(float64))), nil
	case uint:
		return New(left >= b.v.(uint)), nil
	}

	return nil, fmt.Errorf("unexpected %T >= %T", a.v, b.v)
}

func (v *Value) Int() (*Value, error) {
	switch val := v.v.(type) {
	case int:
		return New(val), nil
	case int64:
		return New(int(val)), nil
	case float64:
		return New(int(val)), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) Int64() (*Value, error) {
	switch val := v.v.(type) {
	case int:
		return New(int64(val)), nil
	case int64:
		return New(int64(val)), nil
	case float64:
		return New(int64(val)), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) UInt() (*Value, error) {
	switch val := v.v.(type) {
	case int:
		return New(uint(val)), nil
	case int64:
		return New(uint(val)), nil
	case float64:
		return New(uint(val)), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) Float64() (*Value, error) {
	switch val := v.v.(type) {
	case int:
		return New(float64(val)), nil
	case int64:
		return New(float64(val)), nil
	case float64:
		return New(val), nil
	case uint:
		return New(float64(val)), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) Negative() (*Value, error) {
	switch val := v.v.(type) {
	case int:
		return New(-val), nil
	case int64:
		return New(-val), nil
	case float64:
		return New(-val), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) BitNot() (*Value, error) {
	switch val := v.v.(type) {
	case int:
		return New(^val), nil
	case int64:
		return New(^val), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}

func (v *Value) BoolNot() (*Value, error) {
	switch val := v.v.(type) {
	case bool:
		return New(!val), nil
	}

	return nil, fmt.Errorf("unexpected type: %T", v.v)
}
