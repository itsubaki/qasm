package formatter

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsubaki/qasm/gen/parser"
	"github.com/itsubaki/qasm/listener"
)

type Formatter struct {
	*parser.Baseqasm3ParserListener
	rewriter    *antlr.TokenStreamRewriter
	programName string
}

func New(tokens *antlr.CommonTokenStream) *Formatter {
	return &Formatter{
		Baseqasm3ParserListener: &parser.Baseqasm3ParserListener{},
		rewriter:                antlr.NewTokenStreamRewriter(tokens),
		programName:             "default",
	}
}

func Format(text string) (string, error) {
	lexer := parser.Newqasm3Lexer(antlr.NewInputStream(text))
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.Newqasm3Parser(stream)
	listener := listener.NewErrorListener(p)

	tree := p.Program()
	if len(listener.Errors) > 0 {
		return "", listener.Errors[0]
	}

	f := New(stream)
	antlr.ParseTreeWalkerDefault.Walk(f, tree)
	return f.Format(), nil
}

func (f *Formatter) ExitVersion(ctx *parser.VersionContext) {
	f.rewriter.InsertBefore(f.programName, ctx.GetStart().GetTokenIndex()+1, " ")
	f.rewriter.InsertAfter(f.programName, ctx.GetStop().GetTokenIndex(), "\n")
}

func (f *Formatter) ExitStatement(ctx *parser.StatementContext) {
	f.rewriter.InsertAfter(f.programName, ctx.GetStop().GetTokenIndex(), "\n")
}

func (f *Formatter) ExitIncludeStatement(ctx *parser.IncludeStatementContext) {
	f.rewriter.InsertBefore(f.programName, ctx.GetStart().GetTokenIndex()+1, " ")
}

func (f *Formatter) ExitGateCallStatement(ctx *parser.GateCallStatementContext) {
	f.rewriter.InsertBefore(f.programName, ctx.GetStart().GetTokenIndex()+1, " ")
}

func (f *Formatter) ExitQubitType(ctx *parser.QubitTypeContext) {
	f.rewriter.InsertAfter(f.programName, ctx.GetStop().GetTokenIndex(), " ")
}

func (f *Formatter) Format() string {
	return f.rewriter.GetTextDefault()
}
