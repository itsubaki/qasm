package lex

import (
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	f := c.String("file")
	bin, err := os.ReadFile(f)
	if err != nil {
		return fmt.Errorf("read file=%s: %v", f, err)
	}

	lex := lexer.New(strings.NewReader(string(bin)))
	for {
		token, _ := lex.Tokenize()
		if token == lexer.EOF {
			break
		}

		fmt.Printf("%v ", lexer.Tokens[token])

		if token == lexer.SEMICOLON {
			fmt.Println()
		}
	}

	errs := lex.Errors()
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	return nil
}
