package evaluator

import "github.com/itsubaki/qasm/pkg/evaluator/object"

func Must(obj object.Object, err error) object.Object {
	if err != nil {
		panic(err)
	}

	return obj
}
