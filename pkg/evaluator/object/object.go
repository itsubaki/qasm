package object

import (
	"bytes"
	"strconv"
	"strings"
)

type Type string

const (
	StringType      = "String"
	IntType         = "Int"
	FloatType       = "Float"
	ArrayType       = "Array"
	NilType         = "nil"
	ReturnValueType = "ReturnValue"
)

type Object interface {
	Type() Type
	String() string
}

type Int struct {
	Value int64
}

func (i *Int) Type() Type {
	return IntType
}

func (i *Int) String() string {
	return strconv.FormatInt(i.Value, 10)
}

type Float struct {
	Value float64
}

func (f *Float) Type() Type {
	return FloatType
}

func (f *Float) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

type Nil struct{}

func (n *Nil) Type() Type {
	return NilType
}

func (n *Nil) String() string {
	return "nil"
}

type Array struct {
	Elm []Object
}

func (a *Array) Type() Type {
	return ArrayType
}

func (a *Array) String() string {
	var out bytes.Buffer

	elm := make([]string, 0)
	for _, e := range a.Elm {
		elm = append(elm, e.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elm, ", "))
	out.WriteString("]")

	return out.String()
}

type ReturnValue struct {
	Value Object
}

func (v *ReturnValue) Type() Type {
	return ReturnValueType
}

func (v *ReturnValue) String() string {
	return v.Value.String()
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return StringType
}

func (s *String) String() string {
	return s.Value
}
