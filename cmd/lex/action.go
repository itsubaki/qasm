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
		token, str := lex.Tokenize()
		if token == lexer.EOF {
			break
		}

		fmt.Printf("%v %v ", str, lexer.Tokens[token])

		if token == lexer.SEMICOLON {
			fmt.Println()
		}
	}

	return nil
}
