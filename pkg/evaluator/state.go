package evaluator

import (
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Classical struct {
	Name  string
	Value []int64
}

type State struct {
	Quantum   []qubit.State
	Classical []Classical
}
