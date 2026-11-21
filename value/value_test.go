package value_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/qasm/value"
)

func ExampleValue_Add() {
	a := value.New(1)
	b := value.New(2)

	sum, err := a.Add(b)
	if err != nil {
		panic(err)
	}

	fmt.Println(sum.Value())

	// Output:
	// 3
}

func Equal(a, b *value.Value) bool {
	ok, err := a.Eq(b)
	if err != nil {
		return false
	}

	return ok.Value().(bool)
}

func TestPromote(t *testing.T) {
	cases := []struct {
		a any
		b any
		v any
		w any
	}{
		{
			a: int(1),
			b: int64(2),
			v: int64(1),
			w: int64(2),
		},
		{
			a: int(1),
			b: float64(2.0),
			v: float64(1.0),
			w: float64(2.0),
		},
		{
			a: int64(1),
			b: int(2),
			v: int64(1),
			w: int64(2),
		},
		{
			a: int64(1),
			b: int64(2),
			v: int64(1),
			w: int64(2),
		},
		{
			a: int64(1),
			b: float64(2),
			v: float64(1),
			w: float64(2),
		},
		{
			a: float64(1),
			b: int(2),
			v: float64(1),
			w: float64(2),
		},
		{
			a: float64(1),
			b: int64(2),
			v: float64(1),
			w: float64(2),
		},
		{
			a: float64(1),
			b: float64(2),
			v: float64(1),
			w: float64(2),
		},
	}

	for _, c := range cases {
		v, w, err := value.Promote(value.New(c.a), value.New(c.b))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.v)); !ok {
			t.Error("unexpected error")
		}

		if ok := Equal(w, value.New(c.w)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestPromote_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
	}

	for _, c := range cases {
		if _, _, err := value.Promote(value.New(c.a), value.New(c.b)); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueAdd(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want any
	}{
		{
			a:    int(1),
			b:    int(2),
			want: int(3),
		},
		{
			a:    int(1),
			b:    int64(2),
			want: int64(3),
		},
		{
			a:    int(1),
			b:    float64(2),
			want: float64(3),
		},
		{
			a:    uint(1),
			b:    uint(2),
			want: uint(3),
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.Add(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueAdd_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.Add(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueSub(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want any
	}{
		{
			a:    int(1),
			b:    int(2),
			want: int(-1),
		},
		{
			a:    int(1),
			b:    int64(2),
			want: int64(-1),
		},
		{
			a:    int(1),
			b:    float64(2),
			want: float64(-1),
		},
		{
			a:    uint(3),
			b:    uint(2),
			want: uint(1),
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.Sub(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueSub_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.Sub(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueMul(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want any
	}{
		{
			a:    int(1),
			b:    int(2),
			want: int(2),
		},
		{
			a:    int(1),
			b:    int64(2),
			want: int64(2),
		},
		{
			a:    int(1),
			b:    float64(2),
			want: float64(2),
		},
		{
			a:    uint(3),
			b:    uint(2),
			want: uint(6),
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.Mul(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueMul_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.Mul(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueDiv(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want any
	}{
		{
			a:    int(4),
			b:    int(2),
			want: int(2),
		},
		{
			a:    int(4),
			b:    int64(2),
			want: int64(2),
		},
		{
			a:    int(4),
			b:    float64(2),
			want: float64(2),
		},
		{
			a:    uint(4),
			b:    uint(2),
			want: uint(2),
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.Div(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueDiv_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.Div(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueMod(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want any
	}{
		{
			a:    int(4),
			b:    int(2),
			want: int(0),
		},
		{
			a:    int(4),
			b:    int64(2),
			want: int64(0),
		},
		{
			a:    uint(13),
			b:    uint(3),
			want: uint(1),
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.Mod(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueMod_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.Mod(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValuePow(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want any
	}{
		{
			a:    int(3),
			b:    int(2),
			want: int(9),
		},
		{
			a:    int(3),
			b:    int64(2),
			want: int64(9),
		},
		{
			a:    int(3),
			b:    float64(2),
			want: float64(9),
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.Pow(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValuePow_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.Pow(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueEq(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want bool
	}{
		{
			a:    true,
			b:    true,
			want: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		got, err := a.Eq(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if got.Value() != c.want {
			t.Errorf("got=%v, want=%v", got.Value(), c.want)
		}
	}
}

func TestValueEq_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.Eq(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueNotEq(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want bool
	}{
		{
			a:    true,
			b:    true,
			want: false,
		},
		{
			a:    int(4),
			b:    int(3),
			want: true,
		},
		{
			a:    int(4),
			b:    int64(3),
			want: true,
		},
		{
			a:    int(4),
			b:    float64(3),
			want: true,
		},
		{
			a:    uint(4),
			b:    uint(3),
			want: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.NotEq(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueNotEq_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.NotEq(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueLessThan(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want bool
	}{
		{
			a:    1,
			b:    2,
			want: true,
		},
		{
			a:    1.5,
			b:    2.5,
			want: true,
		},
		{
			a:    int(4),
			b:    int(8),
			want: true,
		},
		{
			a:    int(4),
			b:    int64(8),
			want: true,
		},
		{
			a:    int(4),
			b:    float64(8),
			want: true,
		},
		{
			a:    uint(4),
			b:    uint(8),
			want: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.LessThan(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueLessThan_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.LessThan(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueLessThanOrEqual(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want bool
	}{
		{
			a:    1,
			b:    2,
			want: true,
		},
		{
			a:    1.5,
			b:    2.5,
			want: true,
		},
		{
			a:    int(4),
			b:    int(8),
			want: true,
		},
		{
			a:    int(4),
			b:    int64(8),
			want: true,
		},
		{
			a:    int(4),
			b:    float64(8),
			want: true,
		},
		{
			a:    uint(4),
			b:    uint(8),
			want: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.LessThanOrEqual(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueLessThanOrEqual_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.LessThanOrEqual(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueGreaterThan(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want bool
	}{
		{
			a: 1,
			b: 2,
		},
		{
			a: 1.5,
			b: 2.5,
		},
		{
			a:    int(4),
			b:    int(2),
			want: true,
		},
		{
			a:    int(4),
			b:    int64(2),
			want: true,
		},
		{
			a:    int(4),
			b:    float64(2),
			want: true,
		},
		{
			a:    uint(4),
			b:    uint(2),
			want: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.GreaterThan(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueGreaterThan_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.GreaterThan(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueGreaterThanOrEqual(t *testing.T) {
	cases := []struct {
		a    any
		b    any
		want bool
	}{
		{
			a:    4,
			b:    2,
			want: true,
		},
		{
			a:    4.5,
			b:    2.5,
			want: true,
		},
		{
			a:    int(4),
			b:    int(2),
			want: true,
		},
		{
			a:    int(4),
			b:    int64(2),
			want: true,
		},
		{
			a:    int(4),
			b:    float64(2),
			want: true,
		},
		{
			a:    uint(4),
			b:    uint(2),
			want: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		v, err := a.GreaterThanOrEqual(b)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueGreaterThanOrEqual_invalid(t *testing.T) {
	cases := []struct {
		a any
		b any
	}{
		{
			a: complex(1, 2),
			b: 1,
		},
		{
			a: true,
			b: true,
		},
	}

	for _, c := range cases {
		a, b := value.New(c.a), value.New(c.b)

		if _, err := a.GreaterThanOrEqual(b); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueInt(t *testing.T) {
	cases := []struct {
		v    any
		want int
	}{
		{
			v:    int(4),
			want: 4,
		},
		{
			v:    int64(4),
			want: 4,
		},
		{
			v:    float64(4.5),
			want: 4,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).Int()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueInt_invalid(t *testing.T) {
	cases := []struct {
		v any
	}{
		{
			v: complex(1, 2),
		},
		{
			v: true,
		},
	}

	for _, c := range cases {
		if _, err := value.New(c.v).Int(); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueInt64(t *testing.T) {
	cases := []struct {
		v    any
		want int64
	}{
		{
			v:    int(4),
			want: int64(4),
		},
		{
			v:    int64(4),
			want: int64(4),
		},
		{
			v:    float64(4.5),
			want: int64(4),
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).Int64()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Error("unexpected error")
		}
	}
}

func TestValueInt64_invalid(t *testing.T) {
	cases := []struct {
		v any
	}{
		{
			v: complex(1, 2),
		},
		{
			v: true,
		},
	}

	for _, c := range cases {
		if _, err := value.New(c.v).Int64(); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueUInt(t *testing.T) {
	cases := []struct {
		v    any
		want uint
	}{
		{
			v:    int(4),
			want: uint(4),
		},
		{
			v:    int64(4),
			want: uint(4),
		},
		{
			v:    float64(4.5),
			want: uint(4),
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).UInt()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValueUInt_invalid(t *testing.T) {
	cases := []struct {
		v any
	}{
		{
			v: complex(1, 2),
		},
		{
			v: true,
		},
	}

	for _, c := range cases {
		if _, err := value.New(c.v).UInt(); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueFloat64(t *testing.T) {
	cases := []struct {
		v    any
		want float64
	}{
		{
			v:    int(4),
			want: 4.0,
		},
		{
			v:    int64(4),
			want: 4.0,
		},
		{
			v:    float64(4.5),
			want: 4.5,
		},
		{
			v:    uint(4),
			want: 4.0,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).Float64()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValueFloat64_invalid(t *testing.T) {
	cases := []struct {
		v any
	}{
		{
			v: complex(1, 2),
		},
		{
			v: true,
		},
	}

	for _, c := range cases {
		if _, err := value.New(c.v).Float64(); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueNegative(t *testing.T) {
	cases := []struct {
		v    any
		want any
	}{
		{
			v:    int(4),
			want: int(-4),
		},
		{
			v:    int64(4),
			want: int64(-4),
		},
		{
			v:    float64(4.5),
			want: float64(-4.5),
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).Negative()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValueNegative_invalid(t *testing.T) {
	cases := []struct {
		v any
	}{
		{
			v: complex(1, 2),
		},
		{
			v: true,
		},
	}

	for _, c := range cases {
		if _, err := value.New(c.v).Negative(); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueBitNot(t *testing.T) {
	cases := []struct {
		v    any
		want any
	}{
		{
			v:    true,
			want: false,
		},
		{
			v:    false,
			want: true,
		},
		{
			v:    uint(4),
			want: ^uint(4),
		},
		{
			v:    int(4),
			want: int(^4),
		},
		{
			v:    int64(4),
			want: int64(^4),
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).BitNot()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValueBitNot_invalid(t *testing.T) {
	cases := []struct {
		v any
	}{
		{
			v: complex(1, 2),
		},
	}

	for _, c := range cases {
		if _, err := value.New(c.v).BitNot(); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}

func TestValueBoolNot(t *testing.T) {
	cases := []struct {
		v    any
		want any
	}{
		{
			v:    true,
			want: false,
		},
	}

	for _, c := range cases {
		v, err := value.New(c.v).BoolNot()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if ok := Equal(v, value.New(c.want)); !ok {
			t.Errorf("got=%v, want=%v", v.Value(), c.want)
		}
	}
}

func TestValueBoolNot_invalid(t *testing.T) {
	cases := []struct {
		v any
	}{
		{
			v: complex(1, 2),
		},
		{
			v: int(4),
		},
		{
			v: int64(4),
		},
		{
			v: float64(4.5),
		},
	}

	for _, c := range cases {
		if _, err := value.New(c.v).BoolNot(); err == nil {
			t.Error("expected an error")
			continue
		}
	}
}
