package value_test

import (
	"reflect"
	"testing"

	"github.com/itsubaki/qasm/value"
)

func TestValue_BitAnd(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{a: int(10), b: int(12), want: int(8)},
		{a: int64(10), b: int64(12), want: int64(8)},
		{a: int(10), b: int64(12), want: int64(8)},
		{a: int64(10), b: int(12), want: int64(8)},
		{a: uint(10), b: uint(12), want: uint(8)},
		{a: false, b: false, want: false},
		{a: false, b: true, want: false},
		{a: true, b: false, want: false},
		{a: true, b: true, want: true},
		{a: []bool{true, false, true, false}, b: []bool{true, true, false, false}, want: []bool{true, false, false, false}},
		{a: []bool{}, b: []bool{}, want: []bool{}},
		{a: 1.5, b: 2.5, hasErr: true},
		{a: int64(1), b: 2.5, hasErr: true},
		{a: true, b: int64(1), hasErr: true},
		{a: []bool{true}, b: true, hasErr: true},
		{a: true, b: []bool{true}, hasErr: true},
		{a: []bool{true}, b: []bool{true, false}, hasErr: true},
		{a: nil, b: int64(1), hasErr: true},
	}

	for _, c := range cases {
		got, err := value.New(c.a).BitAnd(value.New(c.b))
		if c.hasErr {
			if err == nil {
				t.Errorf("%v & %v: expected an error", c.a, c.b)
			}
			continue
		}

		if err != nil {
			t.Errorf("%v & %v: unexpected error: %v", c.a, c.b, err)
			continue
		}

		if !reflect.DeepEqual(got.Value(), c.want) {
			t.Errorf("%v & %v: got=%v (%T), want=%v (%T)", c.a, c.b, got.Value(), got.Value(), c.want, c.want)
		}
	}
}

func TestValue_BitOr(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{a: int(10), b: int(12), want: int(14)},
		{a: int64(10), b: int64(12), want: int64(14)},
		{a: int(10), b: int64(12), want: int64(14)},
		{a: int64(10), b: int(12), want: int64(14)},
		{a: uint(10), b: uint(12), want: uint(14)},
		{a: false, b: false, want: false},
		{a: false, b: true, want: true},
		{a: true, b: false, want: true},
		{a: true, b: true, want: true},
		{a: []bool{true, false, true, false}, b: []bool{true, true, false, false}, want: []bool{true, true, true, false}},
		{a: []bool{}, b: []bool{}, want: []bool{}},
		{a: 1.5, b: 2.5, hasErr: true},
		{a: int64(1), b: 2.5, hasErr: true},
		{a: true, b: int64(1), hasErr: true},
		{a: []bool{true}, b: true, hasErr: true},
		{a: true, b: []bool{true}, hasErr: true},
		{a: []bool{true}, b: []bool{true, false}, hasErr: true},
		{a: nil, b: int64(1), hasErr: true},
	}

	for _, c := range cases {
		got, err := value.New(c.a).BitOr(value.New(c.b))
		if c.hasErr {
			if err == nil {
				t.Errorf("%v | %v: expected an error", c.a, c.b)
			}
			continue
		}

		if err != nil {
			t.Errorf("%v | %v: unexpected error: %v", c.a, c.b, err)
			continue
		}

		if !reflect.DeepEqual(got.Value(), c.want) {
			t.Errorf("%v | %v: got=%v (%T), want=%v (%T)", c.a, c.b, got.Value(), got.Value(), c.want, c.want)
		}
	}
}

func TestValue_BitXor(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{a: int(10), b: int(12), want: int(6)},
		{a: int64(10), b: int64(12), want: int64(6)},
		{a: int(10), b: int64(12), want: int64(6)},
		{a: int64(10), b: int(12), want: int64(6)},
		{a: uint(10), b: uint(12), want: uint(6)},
		{a: false, b: false, want: false},
		{a: false, b: true, want: true},
		{a: true, b: false, want: true},
		{a: true, b: true, want: false},
		{a: []bool{true, false, true, false}, b: []bool{true, true, false, false}, want: []bool{false, true, true, false}},
		{a: []bool{}, b: []bool{}, want: []bool{}},
		{a: 1.5, b: 2.5, hasErr: true},
		{a: int64(1), b: 2.5, hasErr: true},
		{a: true, b: int64(1), hasErr: true},
		{a: []bool{true}, b: true, hasErr: true},
		{a: true, b: []bool{true}, hasErr: true},
		{a: []bool{true}, b: []bool{true, false}, hasErr: true},
		{a: nil, b: int64(1), hasErr: true},
	}

	for _, c := range cases {
		got, err := value.New(c.a).BitXor(value.New(c.b))
		if c.hasErr {
			if err == nil {
				t.Errorf("%v ^ %v: expected an error", c.a, c.b)
			}
			continue
		}

		if err != nil {
			t.Errorf("%v ^ %v: unexpected error: %v", c.a, c.b, err)
			continue
		}

		if !reflect.DeepEqual(got.Value(), c.want) {
			t.Errorf("%v ^ %v: got=%v (%T), want=%v (%T)", c.a, c.b, got.Value(), got.Value(), c.want, c.want)
		}
	}
}

func TestValue_ShiftLeft(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{a: int(11), b: int(2), want: int(44)},
		{a: int64(11), b: int64(2), want: int64(44)},
		{a: uint(11), b: uint(2), want: uint(44)},
		{a: int64(11), b: uint(2), want: int64(44)},
		{a: int64(11), b: int64(0), want: int64(11)},
		{a: int64(11), b: int64(64), want: int64(0)},
		{a: true, b: int64(0), want: true},
		{a: false, b: int64(0), want: false},
		{a: true, b: int64(1), want: false},
		{a: []bool{true, false, true, true}, b: int64(1), want: []bool{false, true, true, false}},
		{a: []bool{true, false, true, true}, b: int64(0), want: []bool{true, false, true, true}},
		{a: []bool{true, false, true, true}, b: int64(4), want: []bool{false, false, false, false}},
		{a: []bool{true, false, true, true}, b: ^uint(0), want: []bool{false, false, false, false}},
		{a: []bool{}, b: int64(0), want: []bool{}},
		{a: int64(1), b: int(-1), hasErr: true},
		{a: int64(1), b: int64(-1), hasErr: true},
		{a: []bool{true}, b: int64(-1), hasErr: true},
		{a: int64(1), b: 1.5, hasErr: true},
		{a: int64(1), b: true, hasErr: true},
		{a: 1.5, b: int64(1), hasErr: true},
	}

	for _, c := range cases {
		got, err := value.New(c.a).ShiftLeft(value.New(c.b))
		if c.hasErr {
			if err == nil {
				t.Errorf("%v << %v: expected an error", c.a, c.b)
			}
			continue
		}

		if err != nil {
			t.Errorf("%v << %v: unexpected error: %v", c.a, c.b, err)
			continue
		}

		if !reflect.DeepEqual(got.Value(), c.want) {
			t.Errorf("%v << %v: got=%v (%T), want=%v (%T)", c.a, c.b, got.Value(), got.Value(), c.want, c.want)
		}
	}
}

func TestValue_ShiftRight(t *testing.T) {
	cases := []struct {
		a, b   any
		want   any
		hasErr bool
	}{
		{a: int(11), b: int(1), want: int(5)},
		{a: int64(11), b: int64(1), want: int64(5)},
		{a: uint(11), b: uint(1), want: uint(5)},
		{a: int64(11), b: uint(1), want: int64(5)},
		{a: int64(11), b: int64(0), want: int64(11)},
		{a: int64(-8), b: int64(1), want: int64(-4)},
		{a: int64(11), b: int64(64), want: int64(0)},
		{a: true, b: int64(0), want: true},
		{a: false, b: int64(0), want: false},
		{a: true, b: int64(1), want: false},
		{a: []bool{true, false, true, true}, b: int64(1), want: []bool{false, true, false, true}},
		{a: []bool{true, false, true, true}, b: int64(0), want: []bool{true, false, true, true}},
		{a: []bool{true, false, true, true}, b: int64(4), want: []bool{false, false, false, false}},
		{a: []bool{true, false, true, true}, b: ^uint(0), want: []bool{false, false, false, false}},
		{a: []bool{}, b: int64(0), want: []bool{}},
		{a: int64(1), b: int(-1), hasErr: true},
		{a: int64(1), b: int64(-1), hasErr: true},
		{a: []bool{true}, b: int64(-1), hasErr: true},
		{a: int64(1), b: 1.5, hasErr: true},
		{a: int64(1), b: true, hasErr: true},
		{a: 1.5, b: int64(1), hasErr: true},
	}

	for _, c := range cases {
		got, err := value.New(c.a).ShiftRight(value.New(c.b))
		if c.hasErr {
			if err == nil {
				t.Errorf("%v >> %v: expected an error", c.a, c.b)
			}
			continue
		}

		if err != nil {
			t.Errorf("%v >> %v: unexpected error: %v", c.a, c.b, err)
			continue
		}

		if !reflect.DeepEqual(got.Value(), c.want) {
			t.Errorf("%v >> %v: got=%v (%T), want=%v (%T)", c.a, c.b, got.Value(), got.Value(), c.want, c.want)
		}
	}
}

func TestValue_BitNotBits(t *testing.T) {
	cases := []struct {
		v, want []bool
	}{
		{
			v:    []bool{true, false, true, false},
			want: []bool{false, true, false, true},
		},
		{
			v:    []bool{},
			want: []bool{},
		},
	}

	for _, c := range cases {
		got, err := value.New(c.v).BitNot()
		if err != nil {
			t.Errorf("~%v: unexpected error: %v", c.v, err)
			continue
		}

		if !reflect.DeepEqual(got.Value(), c.want) {
			t.Errorf("~%v: got=%v, want=%v", c.v, got.Value(), c.want)
		}
	}
}
