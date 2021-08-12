package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/qasm/pkg/evaluator"
	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	path := c.String("file")
	if len(path) == 0 {
		cli.ShowAppHelp(c)
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
		return fmt.Errorf("parse: %v", errs)
	}

	fmt.Println(ast)

	e := evaluator.Default()
	if err := e.Eval(ast); err != nil {
		return fmt.Errorf("eval: %v\n", err)
	}

	if err := e.Print(); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}
