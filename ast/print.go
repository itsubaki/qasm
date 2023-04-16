package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/itsubaki/qasm/lexer"
)

const indent = ".  "

func Println(x interface{}) {
	Print(x)
	fmt.Println()
}

func Print(x interface{}) {
	p := printer{
		indent: 0,
		last:   0,
	}

	p.print(reflect.ValueOf(x))
}

type printer struct {
	indent int
	last   byte
}

func (p *printer) Write(data []byte) (int, error) {
	indent := strings.Repeat(indent, p.indent)

	for _, b := range data {
		if p.last == '\n' {
			fmt.Print(indent)
		}

		fmt.Print(string(b))
		p.last = b
	}

	return len(data), nil
}

func (p *printer) printf(format string, a ...interface{}) error {
	if _, err := fmt.Fprintf(p, format, a...); err != nil {
		return fmt.Errorf("fmt.Fprintf: %v", err)
	}

	return nil
}

func (p *printer) slice(x reflect.Value) {
	p.printf("%s (len = %d) {", x.Type(), x.Len())
	if x.Len() > 0 {
		p.indent++
		p.printf("\n")
		for i, n := 0, x.Len(); i < n; i++ {
			p.printf("%d: ", i)
			p.print(x.Index(i))
			p.printf("\n")
		}
		p.indent--
	}
	p.printf("}")
}

func (p *printer) field(i int, t reflect.Type, v reflect.Value) bool {
	name := t.Field(i).Name
	p.printf("%s: ", name)
	if name == "Kind" {
		lit := lexer.Tokens[v.Int()]
		p.print(reflect.ValueOf(lit))
		p.printf("\n")
		return true
	}

	p.print(v)
	return false
}

func (p *printer) strct(x reflect.Value) {
	t := x.Type()
	p.printf("%s {", t)
	p.indent++

	first := true
	for i, n := 0, t.NumField(); i < n; i++ {
		value := x.Field(i)
		if isNil(value) {
			continue
		}

		switch v := value.Interface().(type) {
		case string:
			if len(v) == 0 {
				continue
			}
		}

		if first {
			p.printf("\n")
			first = false
		}

		if p.field(i, t, value) {
			continue
		}

		p.printf("\n")
	}

	p.indent--
	p.printf("}")
}

func (p *printer) print(x reflect.Value) {
	if isNil(x) {
		p.printf("nil")
		return
	}

	switch x.Kind() {
	case reflect.Ptr:
		p.printf("*")
		p.print(x.Elem())

	case reflect.Interface:
		p.print(x.Elem())

	case reflect.Slice:
		p.slice(x)

	case reflect.Struct:
		p.strct(x)

	default:
		p.printf("%v", x)
	}
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}
