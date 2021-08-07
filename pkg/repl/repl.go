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

func New(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)
	e := evaluator.Default()

	fmt.Println(">> OPENQASM 3.0;")
	fmt.Println(">> include \"itsubaki/q.qasm\";")
	for {
		fmt.Printf(">> ")
		if ok := s.Scan(); !ok {
			return
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
				fmt.Printf("[ERROR] %v\n", err)
			}

			continue
		}

		if err := e.Eval(ast); err != nil {
			msg := fmt.Sprintf("[ERROR] eval: %v\n", err)
			io.WriteString(out, msg)
		}

		if len(e.Qubit) == 0 {
			continue
		}

		for _, s := range e.Q.State() {
			fmt.Println(s)
		}
	}
}
