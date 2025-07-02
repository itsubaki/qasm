package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/scan"
)

func main() {
	text := scan.MustText(os.Stdin)
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))

	for _, token := range lexer.GetAllTokens() {
		fmt.Println(token)
	}
}
