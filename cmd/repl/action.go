package repl

import (
	"fmt"
	"os"

	"github.com/itsubaki/qasm/pkg/repl"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if err := repl.Run(os.Stdin, os.Stdout); err != nil {
		return fmt.Errorf("repl run: %v", err)
	}

	return nil
}
