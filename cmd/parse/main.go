package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/io"
)

func main() {
	text := io.MustScan(os.Stdin)
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))
}
