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
	QUBIT    // qubit
	BIT      // bit
	RESET    // reset
	H        // h
	X        // x
	Y        // y
	Z        // z
	CX       // cx
	CCX      // ccx
	SWAP     // swap
	MEASURE  // measure
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
	INCLUDE:  "INCLUDE",
	QUBIT:    "QUBIT",
	BIT:      "BIT",
	RESET:    "RESET",
	H:        "H",
	X:        "X",
	Y:        "Y",
	Z:        "Z",
	CX:       "CX",
	CCX:      "CCX",
	SWAP:     "SWAP",
	MEASURE:  "MEASURE",
}
