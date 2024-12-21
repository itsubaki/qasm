package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/cmd"
	"github.com/itsubaki/qasm/parser"
)

func main() {
	text := cmd.MustScan(os.Stdin)
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	fmt.Println(lexer.GetAllTokens())
}
