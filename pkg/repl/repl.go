package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
)

func New(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(">> ")
		if ok := scanner.Scan(); !ok {
			return
		}

		txt := scanner.Text()
		if len(txt) < 1 {
			continue
		}

		if txt == "quit" {
			break
		}

		lex := lexer.New(strings.NewReader(txt))
		for {
			token, _ := lex.Tokenize()
			if token == lexer.EOF {
				break
			}

			msg := fmt.Sprintf("%v ", lexer.Tokens[token])
			io.WriteString(out, msg)

			if errs := lex.Errors(); len(errs) != 0 {
				for _, err := range errs {
					msg := fmt.Sprintf("%v\n", err)
					io.WriteString(out, msg)
				}
			}
		}
		io.WriteString(out, "\n")
	}
}
