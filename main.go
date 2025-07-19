package main

import (
	"bufio"
	"flag"
	"fmt"
	"strings"

	"maps"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/scan"
	"github.com/itsubaki/qasm/visitor"
)

func main() {
	var filepath string
	var repl, lex, parse bool
	flag.StringVar(&filepath, "f", "", "filepath")
	flag.BoolVar(&repl, "repl", false, "REPL(read-eval-print loop) mode")
	flag.BoolVar(&lex, "lex", false, "the input into a sequence of tokens")
	flag.BoolVar(&parse, "parse", false, "the input into an convert to an AST (abstract syntax tree)")
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
		REPL()
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
	Print(tree.ToStringTree(nil, p))
}

func REPL() {
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
				for _, s := range qsim.State() {
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

func Print(input string) {
	var indent int
	var token string

	flush := func() {
		if token == "" {
			return
		}

		fmt.Print(token)
		token = ""
	}

	first := true
	for i := range len(input) {
		c := input[i]
		switch c {
		case '(':
			flush()
			if !first {
				fmt.Println()
			}

			fmt.Print(strings.Repeat("  ", indent), "(")
			indent++
			first = false
		case ')':
			flush()
			indent--
			fmt.Println()
			fmt.Print(strings.Repeat("  ", indent), ")")
		case ' ', '\n', '\t':
			flush()
			if i+1 < len(input) && input[i+1] != '(' {
				fmt.Print(" ")
			}
		default:
			token += string(c)
		}
	}

	flush()
	fmt.Println()
}
