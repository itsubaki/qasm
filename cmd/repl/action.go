package repl

import (
	"fmt"
	"os"

	"github.com/itsubaki/qasm/repl"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	if err := repl.Run(os.Stdin); err != nil {
		return fmt.Errorf("repl run: %v", err)
	}

	return nil
}
