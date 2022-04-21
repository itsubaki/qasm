package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
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
				fmt.Printf("[ERROR] parse: %v\n", err)
			}
			continue
		}

		if err := e.Eval(ast); err != nil {
			fmt.Printf("[ERROR] eval: %v\n", err)
			continue
		}

		if strings.HasPrefix(txt, "print") {
			continue
		}

		s, err := e.State()
		if err != nil {
			fmt.Printf("[ERROR] state: %v\n", err)
			continue
		}

		e.Println(s)
	}

	return nil
}
