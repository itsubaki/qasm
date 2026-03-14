package visitor

import (
	"fmt"
	"strings"

	"github.com/itsubaki/qasm/angle"
)

type Type string

const (
	StringType      Type = "String"
	IntType         Type = "Int"
	FloatType       Type = "Float"
	BoolType        Type = "Bool"
	AngleType       Type = "Angle"
	ArrayType       Type = "Array"
	ReturnValueType Type = "ReturnValue"
	ErrorType       Type = "Error"
	NilType         Type = "Nil"
	AnnotationType  Type = "Annotation"
	PragmaType      Type = "Pragma"
)

type Object interface {
	Type() Type
	Inspect() string
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return StringType
}

func (s *String) Inspect() string {
	return s.Value
}

type Int struct {
	Value int64
}

func (v *Int) Type() Type {
	return IntType
}

func (v *Int) Inspect() string {
	return fmt.Sprintf("%d", v.Value)
}

type Float struct {
	Value float64
}

func (v *Float) Type() Type {
	return FloatType
}

func (v *Float) Inspect() string {
	return fmt.Sprintf("%f", v.Value)
}

type Bool struct {
	Value bool
}

func (v *Bool) Type() Type {
	return BoolType
}

func (v *Bool) Inspect() string {
	return fmt.Sprintf("%t", v.Value)
}

type Array struct {
	Elements []Object
}

func (v *Array) Type() Type {
	return ArrayType
}

func (v *Array) Inspect() string {
	var sb strings.Builder
	for i, e := range v.Elements {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(e.Inspect())
	}

	return fmt.Sprintf("[%s]", sb.String())
}

type ReturnValue struct {
	Value Object
}

func (v *ReturnValue) Type() Type {
	return ReturnValueType
}

func (v *ReturnValue) Inspect() string {
	return v.Value.Inspect()
}

type Error struct {
	Err error
}

func (v *Error) Type() Type {
	return ErrorType
}

func (v *Error) Inspect() string {
	return v.Err.Error()
}

type Nil struct{}

func (v *Nil) Type() Type {
	return NilType
}

func (v *Nil) Inspect() string {
	return "<nil>"
}

type Annotation struct {
	Keyword              string
	RemainingLineContent string
}

func (v *Annotation) Type() Type {
	return AnnotationType
}

func (v *Annotation) Inspect() string {
	return fmt.Sprintf("%s %s", v.Keyword, v.RemainingLineContent)
}

type Pragma struct {
	RemainingLineContent string
}

func (v *Pragma) Type() Type {
	return PragmaType
}

func (v *Pragma) Inspect() string {
	return v.RemainingLineContent
}

type Angle struct {
	*angle.Angle
}

func (v *Angle) Type() Type {
	return AngleType
}

func (v *Angle) Inspect() string {
	return fmt.Sprintf("%d(%s)", v.K, v.BitString)
}

func (v *Angle) String() string {
	return v.Inspect()
}
