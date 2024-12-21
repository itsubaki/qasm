package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/cmd"
	"github.com/itsubaki/qasm/parser"
	"github.com/itsubaki/qasm/visitor"
)

func main() {
	text := cmd.MustScan(os.Stdin)
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))

	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	tree.Accept(visitor.New())
}
