package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/cmd"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func main() {
	text := cmd.MustScan(os.Stdin)
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))

	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	qsim := q.New()
	env := visitor.NewEnviron()
	v := visitor.New(qsim, env)

	if err := tree.Accept(v); err != nil {
		fmt.Println(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}
}
