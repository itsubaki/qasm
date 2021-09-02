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
	verbose := c.Bool("verbose")
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

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		return fmt.Errorf("parse: %v", errs)
	}

	e := evaluator.Default(evaluator.Opts{
		Verbose: verbose,
	})
	if err := e.Eval(a); err != nil {
		return fmt.Errorf("eval: %v\n", err)
	}

	if err := e.Println(); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}
