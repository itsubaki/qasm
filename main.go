package main

import (
	"os"

	"github.com/itsubaki/qasm/cmd/lex"
	"github.com/itsubaki/qasm/cmd/parse"
	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	app := cli.NewApp()

	app.Name = "qasm"
	app.Usage = "Run Quantum Computation Simulator with OpenQASM 3.0"
	app.Version = "0.1.0"

	lexer := cli.Command{
		Name:   "lex",
		Action: lex.Action,
		Usage:  "convert to a sequence of tokens",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
			},
		},
	}

	parser := cli.Command{
		Name:   "parse",
		Action: parse.Action,
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
		&parser,
	}

	return app
}

func main() {
	if err := New().Run(os.Args); err != nil {
		panic(err)
	}
}
