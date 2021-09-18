package object

import (
	"bytes"
	"strconv"
	"strings"
)

type Type string

const (
	STRING       = "String"
	INT          = "Int"
	FLOAT        = "Float"
	ARRAY        = "Array"
	NIL          = "nil"
	RETURN_VALUE = "ReturnValue"
)

type Object interface {
	Type() Type
	String() string
}

type Int struct {
	Value int64
}

func (i *Int) Type() Type {
	return INT
}

func (i *Int) String() string {
	return strconv.FormatInt(i.Value, 10)
}

type Float struct {
	Value float64
}

func (f *Float) Type() Type {
	return FLOAT
}

func (f *Float) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

type Nil struct{}

func (n *Nil) Type() Type {
	return NIL
}

func (n *Nil) String() string {
	return "nil"
}

type Array struct {
	Elm []Object
}

func (a *Array) Type() Type {
	return ARRAY
}

func (a *Array) String() string {
	var out bytes.Buffer

	var elm []string
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
	return RETURN_VALUE
}

func (v *ReturnValue) String() string {
	return v.Value.String()
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING
}

func (s *String) String() string {
	return s.Value
}
