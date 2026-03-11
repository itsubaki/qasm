package visitor_test

import (
	"testing"

	"github.com/itsubaki/qasm/visitor"
)

func TestPromote(t *testing.T) {
	cases := []struct {
		a, b   any
		hasErr bool
	}{
		{
			a:      complex(1, 2),
			b:      1,
			hasErr: true,
		},
	}

	for _, c := range cases {
		_, _, err := visitor.Promote(visitor.NewValue(c.a), visitor.NewValue(c.b))
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestAdd(t *testing.T) {
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
		a, b := visitor.NewValue(c.a), visitor.NewValue(c.b)
		if _, err := a.Add(b); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestSub(t *testing.T) {
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
		a, b := visitor.NewValue(c.a), visitor.NewValue(c.b)
		if _, err := a.Sub(b); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestMul(t *testing.T) {
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
		a, b := visitor.NewValue(c.a), visitor.NewValue(c.b)
		if _, err := a.Mul(b); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestDiv(t *testing.T) {
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
		a, b := visitor.NewValue(c.a), visitor.NewValue(c.b)
		if _, err := a.Div(b); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestMod(t *testing.T) {
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
		a, b := visitor.NewValue(c.a), visitor.NewValue(c.b)
		if _, err := a.Mod(b); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestEq(t *testing.T) {
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
		a, b := visitor.NewValue(c.a), visitor.NewValue(c.b)
		if _, err := a.Eq(b); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestNotEq(t *testing.T) {
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
		a, b := visitor.NewValue(c.a), visitor.NewValue(c.b)
		if _, err := a.NotEq(b); err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
		}
	}
}
