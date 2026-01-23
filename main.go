package main

import (
	"bufio"
	"flag"
	"fmt"
	"sort"

	"maps"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/q/quantum/qubit"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/scan"
	"github.com/itsubaki/qasm/visitor"
)

func main() {
	var filepath string
	var repl, lex, parse bool
	var top int
	flag.StringVar(&filepath, "f", "", "filepath")
	flag.IntVar(&top, "top", -1, "")
	flag.BoolVar(&repl, "repl", false, "REPL(read-eval-print loop) mode")
	flag.BoolVar(&lex, "lex", false, "Lex the input into a sequence of tokens")
	flag.BoolVar(&parse, "parse", false, "Parse the input and convert it into an AST (abstract syntax tree)")
	flag.Parse()

	switch {
	case lex:
		text, err := Read(filepath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		Lex(text)
	case parse:
		text, err := Read(filepath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		Parse(text)
	case repl:
		REPL(top)
	default:
		text, err := Read(filepath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
		p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

		qsim := q.New()
		env := visitor.NewEnviron()
		v := visitor.New(qsim, env)

		if err := v.Run(p.Program()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, s := range Top(qsim.State(), top) {
			fmt.Println(s)
		}

		fmt.Printf("%-10s: %v\n", "const", env.Const)
		fmt.Printf("%-10s: %v\n", "variable", env.Variable)
		fmt.Printf("%-10s: %v\n", "bit", env.ClassicalBit)
		fmt.Printf("%-10s: %v\n", "qubit", env.Qubit)
		fmt.Printf("%-10s: %v\n", "gate", slices.Sorted(maps.Keys(env.Gate)))
		fmt.Printf("%-10s: %v\n", "subroutine", slices.Sorted(maps.Keys(env.Subroutine)))
	}
}

func Read(filepath string) (string, error) {
	if filepath != "" {
		read, err := os.ReadFile(filepath)
		if err != nil {
			return "", fmt.Errorf("read file %s: %w", filepath, err)
		}

		return string(read), nil
	}

	text, err := scan.Text(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}

	return text, nil
}

func Lex(text string) {
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	for _, token := range lexer.GetAllTokens() {
		fmt.Println(token)
	}
}

func Parse(text string) {
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))
}

func REPL(top int) {
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

			if text == "exit;" || text == "quit;" {
				return
			}

			if text == "print;" {
				for _, s := range Top(qsim.State(), top) {
					fmt.Println(s)
				}

				fmt.Printf("%-10s: %v\n", "const", env.Const)
				fmt.Printf("%-10s: %v\n", "variable", env.Variable)
				fmt.Printf("%-10s: %v\n", "bit", env.ClassicalBit)
				fmt.Printf("%-10s: %v\n", "qubit", env.Qubit)
				fmt.Printf("%-10s: %v\n", "gate", slices.Sorted(maps.Keys(env.Gate)))
				fmt.Printf("%-10s: %v\n", "subroutine", slices.Sorted(maps.Keys(env.Subroutine)))

				continue
			}

			if text == "clear;" {
				qsim = q.New()
				env = visitor.NewEnviron()
				v = visitor.New(qsim, env)
				continue
			}

			lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
			p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

			if err := v.Run(p.Program()); err != nil {
				fmt.Println(err)
				continue
			}

			for _, s := range qsim.State() {
				fmt.Println(s)
			}
		}
	}
}

func Top(s []qubit.State, n int) []qubit.State {
	if n < 0 {
		return s
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].Probability() > s[j].Probability()
	})

	return s[:min(n, len(s))]
}
