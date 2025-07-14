package visitor

import (
	"fmt"
	"strings"
)

type Type string

const (
	StringType      Type = "String"
	IntType         Type = "Int"
	FloatType       Type = "Float"
	BoolType        Type = "Bool"
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

func (i *Int) Type() Type {
	return IntType
}

func (i *Int) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Float struct {
	Value float64
}

func (f *Float) Type() Type {
	return FloatType
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

type Bool struct {
	Value bool
}

func (b *Bool) Type() Type {
	return BoolType
}

func (b *Bool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() Type {
	return ArrayType
}

func (a *Array) Inspect() string {
	var sb strings.Builder
	for i, e := range a.Elements {
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

func (rv *ReturnValue) Type() Type {
	return ReturnValueType
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Err error
}

func (e *Error) Type() Type {
	return ErrorType
}

func (e *Error) Inspect() string {
	return e.Err.Error()
}

type Nil struct{}

func (n *Nil) Type() Type {
	return NilType
}

func (n *Nil) Inspect() string {
	return "nil"
}

type Annotation struct {
	Keyword              string
	RemainingLineContent string
}

func (a *Annotation) Type() Type {
	return AnnotationType
}

func (a *Annotation) Inspect() string {
	return fmt.Sprintf("%s %s", a.Keyword, a.RemainingLineContent)
}

type Pragma struct {
	RemainingLineContent string
}

func (p *Pragma) Type() Type {
	return PragmaType
}

func (p *Pragma) Inspect() string {
	return p.RemainingLineContent
}
