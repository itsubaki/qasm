package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/visitor"
)

func main() {
	var filepath string
	flag.StringVar(&filepath, "f", "", "filepath")
	flag.Parse()

	if filepath == "" {
		fmt.Printf("Usage: %s -f filepath\n", os.Args[0])
		return
	}

	text, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(string(text)))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))

	qsim := q.New()
	v := visitor.New(qsim, visitor.NewEnviron())
	switch err := v.Visit(p.Program()).(type) {
	case error:
		panic(err)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}
}
