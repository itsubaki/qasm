package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/scan"
)

func main() {
	var filepath string
	flag.StringVar(&filepath, "f", "", "filepath")
	flag.Parse()

	var text string
	if filepath != "" {
		read, err := os.ReadFile(filepath)
		if err != nil {
			panic(err)
		}

		text = string(read)
	} else {
		text = scan.MustText(os.Stdin)
	}

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))

	for _, token := range lexer.GetAllTokens() {
		fmt.Println(token)
	}
}
