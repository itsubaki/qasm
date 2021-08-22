package main

import (
	"os"

	"github.com/itsubaki/qasm/cmd"
	"github.com/itsubaki/qasm/cmd/lex"
	"github.com/itsubaki/qasm/cmd/parse"
	"github.com/itsubaki/qasm/cmd/repl"
	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	app := cli.NewApp()

	app.Name = "qasm"
	app.HelpName = "qasm"
	app.Usage = "Run Quantum Computation Simulator with OpenQASM 3.0"
	app.Version = "0.1.0"
	app.Action = cmd.Action
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
		},
	}

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
		Usage:  "convert to an ast (abstract syntax tree)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
			},
		},
	}

	repl := cli.Command{
		Name:   "repl",
		Action: repl.Action,
		Usage:  "read-eval-print loop",
	}

	app.Commands = []*cli.Command{
		&lexer,
		&parser,
		&repl,
	}

	return app
}

func main() {
	if err := New().Run(os.Args); err != nil {
		panic(err)
	}
}
