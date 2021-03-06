package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/itsubaki/qasm/evaluator"
	"github.com/itsubaki/qasm/lexer"
	"github.com/itsubaki/qasm/parser"
)

func Run(in io.Reader) error {
	s := bufio.NewScanner(in)
	e := evaluator.Default()

	fmt.Println(">> OPENQASM 3.0;")
	for {
		fmt.Printf(">> ")
		if ok := s.Scan(); !ok {
			return fmt.Errorf("scanner.Scan")
		}

		txt := s.Text()
		if len(txt) < 1 {
			continue
		}

		if txt == "quit" || txt == "exit" {
			break
		}

		if txt == "clear" {
			e = evaluator.Default()
			continue
		}

		l := lexer.New(strings.NewReader(txt))
		p := parser.New(l)

		ast := p.Parse()
		if errs := p.Errors(); len(errs) != 0 {
			for _, err := range errs {
				fmt.Printf("parse: %v\n", err)
			}
			continue
		}

		if err := e.Eval(ast); err != nil {
			fmt.Printf("eval: %v\n", err)
			continue
		}

		if strings.HasPrefix(txt, "print") {
			continue
		}

		if err := e.Println(); err != nil {
			fmt.Printf("println: %v\n", err)
		}
	}

	return nil
}
