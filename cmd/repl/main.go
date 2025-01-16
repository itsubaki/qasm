package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func main() {
	sigint := make(chan os.Signal, 2)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	input := make(chan string)
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			input <- s.Text()
		}
	}()

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	fmt.Println(">> OPENQASM 3.0;")
	for {
		fmt.Printf(">> ")

		select {
		case <-sigint:
			return
		case text := <-input:
			if len(text) < 1 {
				continue
			}

			if text == "print" {
				fmt.Printf("%-10s: %v\n", "const", env.Const)
				fmt.Printf("%-10s: %v\n", "variable", env.Variable)
				fmt.Printf("%-10s: %v\n", "bit", env.ClassicalBit)
				fmt.Printf("%-10s: %v\n", "qubit", env.Qubit)
				fmt.Printf("%-10s: %v\n", "gate", slices.Sorted(maps.Keys(env.Gate)))
				fmt.Printf("%-10s: %v\n", "subroutine", slices.Sorted(maps.Keys(env.Subroutine)))

				for _, s := range qsim.State() {
					fmt.Println(s)
				}

				continue
			}

			lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
			p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
			tree := p.Program()

			switch err := v.Visit(tree).(type) {
			case error:
				fmt.Println(err)
				continue
			}

			for _, s := range qsim.State() {
				fmt.Println(s)
			}
		}
	}
}
