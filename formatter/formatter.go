package formatter

import (
	"github.com/antlr4-go/antlr/v4"
)

type Formatter struct {
	rewriter *antlr.TokenStreamRewriter
	indent   int
}

func New(tokens *antlr.CommonTokenStream) *Formatter {
	return &Formatter{
		rewriter: antlr.NewTokenStreamRewriter(tokens),
		indent:   0,
	}
}
