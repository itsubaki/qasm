package lex

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	fmt.Println("Hello World")
	return nil
}
