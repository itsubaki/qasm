package repl

import (
	"os"

	"github.com/itsubaki/qasm/pkg/repl"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	repl.New(os.Stdin, os.Stdout)
	return nil
}
