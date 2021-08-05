package lex

import (
	"fmt"
	"os"
	"strings"

	"github.com/itsubaki/qasm/pkg/lexer"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	path := c.String("file")

	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file=%s: %v", path, err)
	}

	lex := lexer.New(strings.NewReader(string(f)))
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

	if errs := lex.Errors(); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}

		return fmt.Errorf("tokenize: %v", errs)
	}

	return nil
}
