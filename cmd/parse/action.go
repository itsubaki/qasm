package parse

import (
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/qasm/ast"
	"github.com/itsubaki/qasm/lexer"
	"github.com/itsubaki/qasm/parser"
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

	a := p.Parse()
	if errs := p.Errors(); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}

		return fmt.Errorf("parse: %v", errs)
	}

	ast.Println(a)
	return nil
}
