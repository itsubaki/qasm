package lexer

type Token int

const (
	// Specials
	ILLEGAL Token = iota
	EOF
	WHITESPACE

	// Literals
	literal_begin
	IDENT  // q
	STRING // "stdgates.qasm"
	INT    // 42
	FLOAT  // 1.23
	literal_end

	// Operators
	operator_begin
	LBRACKET  // '['
	RBRACKET  // ']'
	LBRACE    // '{'
	RBRACE    // '}'
	LPAREN    // '('
	RPAREN    // ')'
	COLON     // ':'
	SEMICOLON // ';'
	DOT       // '.'
	COMMA     // ','
	EQUALS    // '='
	PLUS      // '+'
	MINUS     // '-'
	MUL       // '*'
	DIV       // '/'
	MOD       // '%'
	ARROW     // "->"
	AT        // '@'
	operator_end

	// Keywords
	keyword_begin
	OPENQASM // OPENQASM
	INCLUDE  // include
	CONST    // const
	QUBIT    // qubit
	BIT      // bit
	RESET    // reset
	X        // X
	Y        // Y
	Z        // Z
	H        // H
	S        // S
	T        // T
	U        // U
	QFT      // QFT
	IQFT     // IQFT
	CMODEXP2 // CMODEXP2
	MEASURE  // measure
	GATE     // gate
	PRINT    // print
	DEF      // def
	RETURN   // return
	CTRL     // ctrl
	NEGCTRL  // negctrl
	INV      // inv
	POW      // pow
	PI       // pi
	TAU      // tau
	EULER    // euler
	keyword_end
)

var Tokens = [...]string{
	// Specials
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	WHITESPACE: "WHITESPACE",

	// Literals
	IDENT:  "IDENT",
	STRING: "STRING",
	INT:    "INT",
	FLOAT:  "FLOAT",

	// Operators
	LBRACKET:  "[",
	RBRACKET:  "]",
	LBRACE:    "{",
	RBRACE:    "}",
	LPAREN:    "(",
	RPAREN:    ")",
	COLON:     ":",
	SEMICOLON: ";",
	DOT:       ".",
	COMMA:     ",",
	EQUALS:    "=",
	PLUS:      "+",
	MINUS:     "-",
	MUL:       "*",
	DIV:       "/",
	MOD:       "%",
	ARROW:     "->",
	AT:        "@",

	// Keywords
	OPENQASM: "OPENQASM",
	INCLUDE:  "include",
	CONST:    "const",
	QUBIT:    "qubit",
	BIT:      "bit",
	RESET:    "reset",
	X:        "X",
	Y:        "Y",
	Z:        "Z",
	H:        "H",
	S:        "S",
	T:        "T",
	U:        "U",
	QFT:      "QFT",
	IQFT:     "IQFT",
	CMODEXP2: "CMODEXP2",
	MEASURE:  "measure",
	GATE:     "gate",
	PRINT:    "print",
	DEF:      "def",
	RETURN:   "return",
	CTRL:     "ctrl",
	NEGCTRL:  "negctrl",
	INV:      "inv",
	POW:      "pow",
	PI:       "pi",
	TAU:      "tau",
	EULER:    "euler",
}

func IsModifiler(token Token) bool {
	if token == CTRL || token == NEGCTRL || token == INV || token == POW {
		return true
	}

	return false
}

func IsBinaryOperator(token Token) bool {
	if token == PLUS || token == MINUS || token == MUL || token == DIV || token == MOD {
		return true
	}

	return false
}

func IsBasicLit(token Token) bool {
	if token == IDENT || token == STRING || token == INT || token == FLOAT {
		return true
	}

	return false
}

func IsConst(token Token) bool {
	if token == PI || token == TAU || token == EULER {
		return true
	}

	return false
}
