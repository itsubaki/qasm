package listener

import (
	"fmt"

	"github.com/itsubaki/qasm/gen/parser"
)

type Listener struct {
	*parser.Baseqasm3ParserListener
}

func New() *Listener {
	return &Listener{
		&parser.Baseqasm3ParserListener{},
	}
}

func (l *Listener) EnterProgram(ctx *parser.ProgramContext) {
	fmt.Println("[DEBUG] EnterProgram")
}

func (l *Listener) ExitProgram(ctx *parser.ProgramContext) {
	fmt.Println("[DEBUG] ExitProgram")
}
