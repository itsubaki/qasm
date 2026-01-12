package visitor_test

import (
	"errors"
	"math"
	"testing"

	"github.com/itsubaki/qasm/visitor"
)

func TestAngle(t *testing.T) {
	type want struct {
		bitString string
		k         uint
		radian    float64
	}

	cases := []struct {
		bits   uint
		radian float64
		want   want
	}{
		{
			bits:   4,
			radian: math.Pi,
			want: want{
				bitString: "1000",
				k:         8,
				radian:    math.Pi,
			},
		},
		{
			bits:   6,
			radian: math.Pi / 2,
			want: want{
				bitString: "010000",
				k:         16,
				radian:    math.Pi / 2,
			},
		},
		{
			bits:   8,
			radian: 7 * math.Pi / 8,
			want: want{
				bitString: "01110000",
				k:         112,
				radian:    7 * math.Pi / 8,
			},
		},
		{
			bits:   4,
			radian: -math.Pi / 2,
			want: want{
				bitString: "1100",
				k:         12,
				radian:    3 * math.Pi / 2, // normalized to [0, 2pi)
			},
		},
	}

	for _, c := range cases {
		angle := visitor.NewAngle(c.bits, c.radian)
		if angle.BitString != c.want.bitString {
			t.Errorf("got=%v, want=%v", angle.BitString, c.want.bitString)
		}

		if angle.K != c.want.k {
			t.Errorf("got=%v, want=%v", angle.K, c.want.k)
		}

		if angle.Radian() != c.want.radian {
			t.Errorf("got=%v, want=%v", angle.Radian(), c.want.radian)
		}
	}
}

func TestObject_Inspect(t *testing.T) {
	cases := []struct {
		obj     visitor.Object
		objType visitor.Type
		want    string
	}{
		{
			obj: &visitor.String{
				Value: "Hello, World!",
			},
			objType: visitor.StringType,
			want:    "Hello, World!",
		},
		{
			obj: &visitor.Int{
				Value: 42,
			},
			objType: visitor.IntType,
			want:    "42",
		},
		{
			obj: &visitor.Float{
				Value: 42.123456,
			},
			objType: visitor.FloatType,
			want:    "42.123456",
		},
		{
			obj: &visitor.Bool{
				Value: true,
			},
			objType: visitor.BoolType,
			want:    "true",
		},
		{
			obj: &visitor.Array{Elements: []visitor.Object{
				&visitor.String{Value: "Hello, World!"},
				&visitor.Int{Value: 42},
				&visitor.Float{Value: 42.123456},
			}},
			objType: visitor.ArrayType,
			want:    "[Hello, World!, 42, 42.123456]",
		},
		{
			obj:     &visitor.Nil{},
			objType: visitor.NilType,
			want:    "nil",
		},
		{
			obj: &visitor.ReturnValue{
				Value: &visitor.String{
					Value: "Hello, World!",
				},
			},
			objType: visitor.ReturnValueType,
			want:    "Hello, World!",
		},
		{
			obj: &visitor.Error{
				Err: errors.New("something went wrong"),
			},
			objType: visitor.ErrorType,
			want:    "something went wrong",
		},
		{
			obj: &visitor.Annotation{
				Keyword:              "TODO",
				RemainingLineContent: "Fix this issue",
			},
			objType: visitor.AnnotationType,
			want:    "TODO Fix this issue",
		},
		{
			obj: &visitor.Pragma{
				RemainingLineContent: "This is a pragma",
			},
			objType: visitor.PragmaType,
			want:    "This is a pragma",
		},
		{
			obj: &visitor.Angle{
				Bits:      4,
				BitString: "1000",
				K:         8,
			},
			objType: visitor.AngleType,
			want:    "8(1000)",
		},
	}

	for _, c := range cases {
		gotT := c.obj.Type()
		if gotT != c.objType {
			t.Errorf("got=%q, want %q", gotT, c.objType)
		}

		got := c.obj.Inspect()
		if got != c.want {
			t.Errorf("got=%q, want %q", got, c.want)
		}
	}
}
