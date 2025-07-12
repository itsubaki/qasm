package main

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/scan"
	"github.com/itsubaki/qasm/visitor"
)

func main() {
	var filepath string
	flag.StringVar(&filepath, "f", "", "filepath")
	flag.Parse()

	var text string
	if filepath != "" {
		read, err := os.ReadFile(filepath)
		if err != nil {
			panic(err)
		}

		text = string(read)
	} else {
		text = scan.MustText(os.Stdin)
	}

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(string(text)))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := v.Run(p.Program()); err != nil {
		panic(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	fmt.Printf("%-10s: %v\n", "const", env.Const)
	fmt.Printf("%-10s: %v\n", "variable", env.Variable)
	fmt.Printf("%-10s: %v\n", "bit", env.ClassicalBit)
	fmt.Printf("%-10s: %v\n", "qubit", env.Qubit)
	fmt.Printf("%-10s: %v\n", "gate", slices.Sorted(maps.Keys(env.Gate)))
	fmt.Printf("%-10s: %v\n", "subroutine", slices.Sorted(maps.Keys(env.Subroutine)))
}
