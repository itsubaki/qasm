package listener_test

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/listener"
)

func Example() {
	text := "OPENQASM 3.0;"
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	p := parser.Newqasm3Parser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	p.AddParseListener(listener.New())

	tree := p.Program()
	fmt.Println(tree.ToStringTree(nil, p))

	// Output:
	// [DEBUG] EnterProgram
	// [DEBUG] ExitProgram
	// (program (version OPENQASM 3.0 ;) <EOF>)
}
