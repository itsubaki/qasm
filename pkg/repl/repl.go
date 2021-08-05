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
	scanner := bufio.NewScanner(in)
	eval := evaluator.Default()

	for {
		fmt.Printf(">> ")
		if ok := scanner.Scan(); !ok {
			return
		}

		txt := scanner.Text()
		if len(txt) < 1 {
			continue
		}

		if txt == "quit" || txt == "exit" {
			break
		}

		lex := lexer.New(strings.NewReader(txt))
		p := parser.New(lex)

		ast := p.Parse()
		if errs := p.Errors(); len(errs) != 0 {
			for _, err := range errs {
				fmt.Println(err)
			}

			continue
		}

		if err := eval.Eval(ast); err != nil {
			msg := fmt.Sprintf("eval: %v\n", err)
			io.WriteString(out, msg)
		}
	}
}
