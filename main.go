package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/qasm/cmd/lex"
	"github.com/urfave/cli/v2"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "qasm"
	app.Usage = "Run Quantum Computation Simulator with OpenQASM 3.0"
	app.Version = version

	lexer := cli.Command{
		Name:   "lex",
		Action: lex.Action,
		Usage:  "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
			},
		},
	}

	app.Commands = []*cli.Command{
		&lexer,
	}

	return app
}

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
