package visitor_test

import (
	"errors"
	"testing"

	"github.com/itsubaki/qasm/visitor"
)

func TestObject_Inspect(t *testing.T) {
	cases := []struct {
		obj     visitor.Object
		objType visitor.Type
		want    string
	}{
		{
			obj:     &visitor.String{Value: "Hello, World!"},
			objType: visitor.StringType,
			want:    "Hello, World!",
		},
		{
			obj:     &visitor.Int{Value: 42},
			objType: visitor.IntType,
			want:    "42",
		},
		{
			obj:     &visitor.Float{Value: 42.123456},
			objType: visitor.FloatType,
			want:    "42.123456",
		},
		{
			obj:     &visitor.Bool{Value: true},
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
			obj:     &visitor.ReturnValue{Value: &visitor.String{Value: "Hello, World!"}},
			objType: visitor.ReturnValueType,
			want:    "Hello, World!",
		},
		{
			obj:     &visitor.Error{Err: errors.New("something went wrong")},
			objType: visitor.ErrorType,
			want:    "something went wrong",
		},
		{
			obj:     &visitor.Annotation{Keyword: "TODO", RemainingLineContent: "Fix this issue"},
			objType: visitor.AnnotationType,
			want:    "TODO Fix this issue",
		},
		{
			obj:     &visitor.Pragma{RemainingLineContent: "This is a pragma"},
			objType: visitor.PragmaType,
			want:    "This is a pragma",
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
