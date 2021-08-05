package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
	"github.com/itsubaki/qasm/pkg/repl"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	path := c.String("file")

	if len(path) == 0 {
		repl.New(os.Stdin, os.Stdout)
		return nil
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file=%s: %v", path, err)
	}

	l := lexer.New(strings.NewReader(string(f)))
	p := parser.New(l)

	ast := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}

		return fmt.Errorf("parse: %v", errs)
	}

	fmt.Println(ast)

	if err := evaluator.Default().Eval(ast); err != nil {
		return fmt.Errorf("eval: %v\n", err)
	}

	return nil
}
