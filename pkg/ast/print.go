package ast

import (
	"fmt"
	"reflect"
)

func Print(x interface{}) {
	print(reflect.ValueOf(x))
}

func printf(format string, args ...interface{}) {
	if _, err := fmt.Printf(format, args...); err != nil {
		panic(err)
	}
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}

func print(x reflect.Value) {
	if isNil(x) {
		printf("nil")
		return
	}

	switch x.Kind() {
	case reflect.Interface:
		print(x.Elem())
	case reflect.Ptr:
		printf("*")
		print(x.Elem())
	case reflect.Struct:
		t := x.Type()
		printf("%s {", t)

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
				printf("\n")
				first = false
			}

			name := t.Field(i).Name
			printf("%s: ", name)
			print(value)
			printf("\n")
		}
		printf("}")
	case reflect.Array:
		printf("%s {", x.Type())
		if x.Len() > 0 {
			printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				printf("%d: ", i)
				print(x.Index(i))
				printf("\n")
			}
		}
		printf("}")
	case reflect.Slice:
		printf("%s {", x.Type())
		if x.Len() > 0 {
			printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				printf("%d: ", i)
				print(x.Index(i))
				printf("\n")
			}
		}
		printf("}")
	default:
		v := x.Interface()
		switch v := v.(type) {
		case string:
			printf("%q", v)
			return
		default:
			printf("%v", x)
		}
	}
}
