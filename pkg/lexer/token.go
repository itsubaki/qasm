package lexer

type Token int

const (
	// Specials
	ILLEGAL Token = iota
	EOF
	WHITESPACE

	// Literals
	literal_begin
	IDENT
	STRING
	INT
	FLOAT
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
	operator_end

	// Keywords
	keyword_begin
	OPENQASM // OPENQASM
	INCLUDE  // include
	CONST    // const
	QUBIT    // qubit
	BIT      // bit
	RESET    // reset
	X        // x
	Y        // y
	Z        // z
	H        // h
	S        // S
	T        // T
	CX       // cx
	CZ       // cz
	CCX      // ccx
	SWAP     // swap
	QFT      // qft
	IQFT     // iqft
	CMODEXP2 // cmodexp2
	MEASURE  // measure
	PRINT    // print
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

	// Keywords
	OPENQASM: "OPENQASM",
	INCLUDE:  "include",
	CONST:    "const",
	QUBIT:    "qubit",
	BIT:      "bit",
	RESET:    "reset",
	X:        "x",
	Y:        "y",
	Z:        "z",
	H:        "h",
	S:        "s",
	T:        "t",
	CX:       "cx",
	CZ:       "cz",
	CCX:      "ccx",
	SWAP:     "swap",
	QFT:      "qft",
	IQFT:     "iqft",
	CMODEXP2: "cmodexp2",
	MEASURE:  "measure",
	PRINT:    "print",
}
