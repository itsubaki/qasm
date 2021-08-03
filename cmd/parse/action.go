package parse

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	fmt.Println("hello world")

	return nil
}
