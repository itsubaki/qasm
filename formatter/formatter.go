package formatter

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
)

type Formatter struct {
	*parser.Baseqasm3ParserListener
	rewriter *antlr.TokenStreamRewriter
	indent   int
}

func New(tokens *antlr.CommonTokenStream) *Formatter {
	return &Formatter{
		Baseqasm3ParserListener: &parser.Baseqasm3ParserListener{},
		rewriter:                antlr.NewTokenStreamRewriter(tokens),
		indent:                  0,
	}
}
