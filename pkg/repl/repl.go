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
		}

		l := lexer.New(strings.NewReader(txt))
		p := parser.New(l)

		ast := p.Parse()
		if errs := p.Errors(); len(errs) != 0 {
			for _, err := range errs {
				fmt.Println(err)
			}

			continue
		}

		if err := e.Eval(ast); err != nil {
			msg := fmt.Sprintf("eval: %v\n", err)
			io.WriteString(out, msg)
		}
	}
}
