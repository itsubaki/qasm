package ast

import (
	"fmt"
	"reflect"
)

var indent = []byte(". ")

func Print(x interface{}) {
	p := printer{}
	p.print(reflect.ValueOf(x))
}

type printer struct {
	indent int
}

func (p *printer) Write(data []byte) (int, error) {
	return fmt.Print(string(data))
}

func (p *printer) printf(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(p, format, args...); err != nil {
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
			p.print(value)
			p.printf("\n")
		}
		p.printf("}")
	case reflect.Array:
		p.printf("%s {", x.Type())
		if x.Len() > 0 {
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf("\n")
			}
		}
		p.printf("}")
	case reflect.Slice:
		p.printf("%s {", x.Type())
		if x.Len() > 0 {
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf("\n")
			}
		}
		p.printf("}")
	default:
		v := x.Interface()
		switch v := v.(type) {
		case string:
			p.printf("%q", v)
			return
		default:
			p.printf("%v", x)
		}
	}
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}
