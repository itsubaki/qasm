package parse

import (
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/itsubaki/qasm/pkg/parser"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	path := c.String("file")

	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file=%s: %v", path, err)
	}

	p := parser.New(lexer.New(strings.NewReader(string(f))))
	ast := p.Parse()
	fmt.Print(ast)

	if errs := p.Errors(); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	return nil
}
