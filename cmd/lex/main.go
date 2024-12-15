package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/parser"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var text string
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	fmt.Println(lexer.GetAllTokens())
}
