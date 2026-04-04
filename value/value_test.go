package value_test

import (
	"testing"

	"github.com/itsubaki/qasm/value"
)

func Equal(a, b *value.Value) bool {
	ok, err := a.Eq(b)
	if err != nil {
		return false
	}

	return ok.Value().(bool)
}

func TestPromote(t *testing.T) {
	cases := []struct {
		a, b   any
		v, w   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      int(1),
			b:      int64(2),
			v:      int64(1),
			w:      int64(2),
			hasErr: false,
		},
		{
			a:      int(1),
			b:      float64(2.0),
			v:      float64(1.0),
			w:      float64(2.0),
			hasErr: false,
		},
		{
			a:      int64(1),
			b:      int(2),
			v:      int64(1),
			w:      int64(2),
			hasErr: false,
		},
		{
			a:      int64(1),
			b:      int64(2),
			v:      int64(1),
			w:      int64(2),
			hasErr: false,
		},
		{
			a:      int64(1),
			b:      float64(2),
			v:      float64(1),
			w:      float64(2),
			hasErr: false,
		},
		{
			a:      float64(1),
			b:      int(2),
			v:      float64(1),
			w:      float64(2),
			hasErr: false,
		},
		{
			a:      float64(1),
			b:      int64(2),
			v:      float64(1),
			w:      float64(2),
			hasErr: false,
		},
		{
			a:      float64(1),
			b:      float64(2),
			v:      float64(1),
			w:      float64(2),
			hasErr: false,
		},
	}

	for _, c := range cases {
		v, w, err := value.Promote(value.New(c.a), value.New(c.b))
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.v)); !ok {
			t.Error("unexpected error")
		}

		if ok := Equal(w, value.New(c.w)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Add(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      int(1),
			b:      int(2),
			want:   int(3),
			hasErr: false,
		},
		{
			a:      int(1),
			b:      int64(2),
			want:   int64(3),
			hasErr: false,
		},
		{
			a:      int(1),
			b:      float64(2),
			want:   float64(3),
			hasErr: false,
		},
		{
			a:      uint(1),
			b:      uint(2),
			want:   uint(3),
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.Add(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Sub(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      int(1),
			b:      int(2),
			want:   int(-1),
			hasErr: false,
		},
		{
			a:      int(1),
			b:      int64(2),
			want:   int64(-1),
			hasErr: false,
		},
		{
			a:      int(1),
			b:      float64(2),
			want:   float64(-1),
			hasErr: false,
		},
		{
			a:      uint(3),
			b:      uint(2),
			want:   uint(1),
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.Sub(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Mul(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      int(1),
			b:      int(2),
			want:   int(2),
			hasErr: false,
		},
		{
			a:      int(1),
			b:      int64(2),
			want:   int64(2),
			hasErr: false,
		},
		{
			a:      int(1),
			b:      float64(2),
			want:   float64(2),
			hasErr: false,
		},
		{
			a:      uint(3),
			b:      uint(2),
			want:   uint(6),
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.Mul(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Div(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      int(4),
			b:      int(2),
			want:   int(2),
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int64(2),
			want:   int64(2),
			hasErr: false,
		},
		{
			a:      int(4),
			b:      float64(2),
			want:   float64(2),
			hasErr: false,
		},
		{
			a:      uint(4),
			b:      uint(2),
			want:   uint(2),
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.Div(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Mod(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      int(4),
			b:      int(2),
			want:   int(0),
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int64(2),
			want:   int64(0),
			hasErr: false,
		},
		{
			a:      uint(13),
			b:      uint(3),
			want:   uint(1),
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.Mod(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Pow(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      int(3),
			b:      int(2),
			want:   int(9),
			hasErr: false,
		},
		{
			a:      int(3),
			b:      int64(2),
			want:   int64(9),
			hasErr: false,
		},
		{
			a:      int(3),
			b:      float64(2),
			want:   float64(9),
			hasErr: false,
		},
		{
			a:      float64(9),
			b:      int(2),
			want:   float64(81),
			hasErr: false,
		},
		{
			a:      float64(9),
			b:      float64(0.5),
			want:   float64(3),
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.Pow(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Eq(t *testing.T) {
	cases := []struct {
		a, b   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		if _, err := a.Eq(b); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestValue_NotEq(t *testing.T) {
	cases := []struct {
		a, b   any
		want   bool
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      int(4),
			b:      int(3),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int64(3),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      float64(3),
			want:   true,
			hasErr: false,
		},
		{
			a:      uint(4),
			b:      uint(3),
			want:   true,
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.NotEq(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_LessThan(t *testing.T) {
	cases := []struct {
		a, b   any
		want   bool
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      1,
			b:      2,
			want:   true,
			hasErr: false,
		},
		{
			a:      1.5,
			b:      2.5,
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int(8),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int64(8),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      float64(8),
			want:   true,
			hasErr: false,
		},
		{
			a:      uint(4),
			b:      uint(8),
			want:   true,
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.LessThan(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_LessThanOrEqual(t *testing.T) {
	cases := []struct {
		a, b   any
		want   bool
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      1,
			b:      2,
			want:   true,
			hasErr: false,
		},
		{
			a:      1.5,
			b:      2.5,
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int(8),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int64(8),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      float64(8),
			want:   true,
			hasErr: false,
		},
		{
			a:      uint(4),
			b:      uint(8),
			want:   true,
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.LessThanOrEqual(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_GreaterThan(t *testing.T) {
	cases := []struct {
		a, b   any
		want   bool
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      1,
			b:      2,
			hasErr: false,
		},
		{
			a:      1.5,
			b:      2.5,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int(2),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int64(2),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      float64(2),
			want:   true,
			hasErr: false,
		},
		{
			a:      uint(4),
			b:      uint(2),
			want:   true,
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.GreaterThan(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_GreaterThanOrEqual(t *testing.T) {
	cases := []struct {
		a, b   any
		want   bool
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
		{
			a:      true,
			b:      true,
			hasErr: true,
		},
		{
			a:      4,
			b:      2,
			want:   true,
			hasErr: false,
		},
		{
			a:      4.5,
			b:      2.5,
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int(2),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      int64(2),
			want:   true,
			hasErr: false,
		},
		{
			a:      int(4),
			b:      float64(2),
			want:   true,
			hasErr: false,
		},
		{
			a:      uint(4),
			b:      uint(2),
			want:   true,
			hasErr: false,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)
		v, err := a.GreaterThanOrEqual(b)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Int(t *testing.T) {
	cases := []struct {
		v      any
		want   int
		hasErr bool
	}{
		{
			v:      complex(1, 2),
			hasErr: true,
		},
		{
			v:      true,
			hasErr: true,
		},
		{
			v:      int(4),
			want:   4,
			hasErr: false,
		},
		{
			v:      int64(4),
			want:   4,
			hasErr: false,
		},
		{
			v:      float64(4.5),
			want:   4,
			hasErr: false,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).Int()
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_Int64(t *testing.T) {
	cases := []struct {
		v      any
		want   int64
		hasErr bool
	}{
		{
			v:      complex(1, 2),
			hasErr: true,
		},
		{
			v:      true,
			hasErr: true,
		},
		{
			v:      int(4),
			want:   int64(4),
			hasErr: false,
		},
		{
			v:      int64(4),
			want:   int64(4),
			hasErr: false,
		},
		{
			v:      float64(4.5),
			want:   int64(4),
			hasErr: false,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).Int64()
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValue_UInt(t *testing.T) {
	cases := []struct {
		v      any
		want   uint
		hasErr bool
	}{
		{
			v:      complex(1, 2),
			hasErr: true,
		},
		{
			v:      true,
			hasErr: true,
		},
		{
			v:      int(4),
			want:   uint(4),
			hasErr: false,
		},
		{
			v:      int64(4),
			want:   uint(4),
			hasErr: false,
		},
		{
			v:      float64(4.5),
			want:   uint(4),
			hasErr: false,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).UInt()
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValue_Float64(t *testing.T) {
	cases := []struct {
		v      any
		want   float64
		hasErr bool
	}{
		{
			v:      complex(1, 2),
			hasErr: true,
		},
		{
			v:      true,
			hasErr: true,
		},
		{
			v:      int(4),
			want:   4.0,
			hasErr: false,
		},
		{
			v:      int64(4),
			want:   4.0,
			hasErr: false,
		},
		{
			v:      float64(4.5),
			want:   4.5,
			hasErr: false,
		},
		{
			v:      uint(4),
			want:   4.0,
			hasErr: false,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).Float64()
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValue_Negative(t *testing.T) {
	cases := []struct {
		v      any
		want   any
		hasErr bool
	}{
		{
			v:      complex(1, 2),
			hasErr: true,
		},
		{
			v:      true,
			hasErr: true,
		},
		{
			v:      int(4),
			want:   int(-4),
			hasErr: false,
		},
		{
			v:      int64(4),
			want:   int64(-4),
			hasErr: false,
		},
		{
			v:      float64(4.5),
			want:   float64(-4.5),
			hasErr: false,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).Negative()
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValue_BitNot(t *testing.T) {
	cases := []struct {
		v      any
		want   any
		hasErr bool
	}{
		{
			v:      complex(1, 2),
			hasErr: true,
		},
		{
			v:      true,
			hasErr: true,
		},
		{
			v:      int(4),
			want:   int(^4),
			hasErr: false,
		},
		{
			v:      int64(4),
			want:   int64(^4),
			hasErr: false,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).BitNot()
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValue_BoolNot(t *testing.T) {
	cases := []struct {
		v      any
		want   any
		hasErr bool
	}{
		{
			v:      complex(1, 2),
			hasErr: true,
		},
		{
			v:      int(4),
			hasErr: true,
		},
		{
			v:      int64(4),
			hasErr: true,
		},
		{
			v:      float64(4.5),
			hasErr: true,
		},
		{
			v:      true,
			want:   false,
			hasErr: false,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).BoolNot()
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}
