package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
)

var indent = []byte(".  ")

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
	indent := strings.Repeat(string(indent), p.indent)

	for _, b := range data {
		if p.last == '\n' {
			fmt.Printf(indent)
		}

		fmt.Printf(string(b))
		p.last = b
	}

	return len(data), nil
}

func (p *printer) printf(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(p, format, a...); err != nil {
		panic(err)
	}
}

func (p *printer) print(x reflect.Value) {
	if isNil(x) {
		p.printf("nil")
		return
	}

	switch x.Kind() {
	case reflect.Interface:
		p.print(x.Elem())

	case reflect.Ptr:
		p.printf("*")
		p.print(x.Elem())

	case reflect.Struct:
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

			name := t.Field(i).Name
			p.printf("%s: ", name)
			if name == "Kind" {
				lit := lexer.Tokens[value.Int()]
				p.print(reflect.ValueOf(lit))
				p.printf("\n")
				continue
			}

			p.print(value)
			p.printf("\n")
		}
		p.indent--
		p.printf("}")

	case reflect.Array:
		p.printf("%s {", x.Type())
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

	case reflect.Slice:
		p.printf("%s {", x.Type())
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
