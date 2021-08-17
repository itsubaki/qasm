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
	if len(path) == 0 {
		cli.ShowAppHelp(c)
		return nil
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file=%s: %v", path, err)
	}

	l := lexer.New(strings.NewReader(string(f)))
	newline := true
	for {
		token, _ := l.Tokenize()
		if token == lexer.EOF {
			break
		}

		fmt.Printf("%v ", lexer.Tokens[token])

		if token == lexer.LBRACE {
			newline = false
		}

		if token == lexer.RBRACE {
			newline = true
		}

		if !newline {
			continue
		}

		if token == lexer.SEMICOLON || token == lexer.RBRACE {
			fmt.Println()
		}
	}

	if errs := l.Errors(); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}

		return fmt.Errorf("tokenize: %v", errs)
	}

	return nil
}
