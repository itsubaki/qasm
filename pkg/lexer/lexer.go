package lexer

import (
	"bufio"
	"bytes"
	"io"
)

var (
	operator map[string]Token = make(map[string]Token)
	keyword  map[string]Token = make(map[string]Token)
	cnst     map[string]Token = make(map[string]Token)
)

func init() {
	for i := operator_begin + 1; i < operator_end; i++ {
		operator[Tokens[i]] = i
	}

	for i := keyword_begin + 1; i < keyword_end; i++ {
		keyword[Tokens[i]] = i
	}

	for i := const_begin + 1; i < const_end; i++ {
		cnst[Tokens[i]] = i
	}
}

type Lexer struct {
	eof    rune
	r      *bufio.Reader
	errors []error
}

func New(r io.Reader) *Lexer {
	return &Lexer{
		eof:    rune(-1),
		r:      bufio.NewReader(r),
		errors: make([]error, 0),
	}
}

func (l *Lexer) Errors() []error {
	return l.errors
}

func (l *Lexer) Tokenize() (Token, string) {
	return l.TokenizeIgnore(WHITESPACE)
}

func (l *Lexer) TokenizeIgnore(t ...Token) (Token, string) {
	ignore := make(map[Token]bool)
	for _, tt := range t {
		ignore[tt] = true
	}

	for {
		token, literal := l.Scan()
		if _, ok := ignore[token]; ok {
			continue
		}

		return token, literal
	}
}

func (l *Lexer) Scan() (Token, string) {
	ch := l.read()
	if ch == l.eof {
		return EOF, ""
	}

	if isWhitespace(ch) {
		l.unread()
		return l.whitespace()
	}

	if isLetter(ch) {
		l.unread()
		str := l.scan()

		if v, ok := keyword[str]; ok {
			return v, str
		}

		if v, ok := cnst[str]; ok {
			return v, str
		}

		return IDENT, str
	}

	if isDigit(ch) {
		l.unread()
		return l.scanNumber()
	}

	if isString(ch) {
		l.unread()
		return STRING, l.scanString()
	}

	if ch == '-' {
		if l.read() == '>' {
			return ARROW, "->"
		}

		l.unread()
	}

	if v, ok := operator[string(ch)]; ok {
		return v, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (l *Lexer) scan() string {
	var buf bytes.Buffer
	if _, err := buf.WriteRune(l.read()); err != nil {
		l.errors = append(l.errors, err)
	}

	for {
		ch := l.read()
		if ch == l.eof {
			break
		}

		if isLetter(ch) || isDigit(ch) {
			if _, err := buf.WriteRune(ch); err != nil {
				l.errors = append(l.errors, err)
			}

			continue
		}

		l.unread()
		break
	}

	return buf.String()
}

func (l *Lexer) scanString() string {
	var buf bytes.Buffer
	if _, err := buf.WriteRune(l.read()); err != nil {
		l.errors = append(l.errors, err)
	}

	for {
		ch := l.read()
		if ch == l.eof {
			break
		}

		if _, err := buf.WriteRune(ch); err != nil {
			l.errors = append(l.errors, err)
		}

		if isString(ch) {
			break
		}
	}

	return buf.String()
}

func (l *Lexer) scanNumber() (Token, string) {
	var buf bytes.Buffer
	if _, err := buf.WriteRune(l.read()); err != nil {
		l.errors = append(l.errors, err)
	}

	token := INT
	for {
		ch := l.read()
		if ch == l.eof {
			break
		}

		if ch == '.' {
			if _, err := buf.WriteRune(ch); err != nil {
				l.errors = append(l.errors, err)
			}

			token = FLOAT
			continue
		}

		if isDigit(ch) {
			if _, err := buf.WriteRune(ch); err != nil {
				l.errors = append(l.errors, err)
			}

			continue
		}

		l.unread()
		break
	}

	return token, buf.String()
}

func (l *Lexer) whitespace() (Token, string) {
	var buf bytes.Buffer
	if _, err := buf.WriteRune(l.read()); err != nil {
		l.errors = append(l.errors, err)
	}

	for {
		ch := l.read()
		if ch == l.eof {
			break
		}

		if isWhitespace(ch) {
			if _, err := buf.WriteRune(ch); err != nil {
				l.errors = append(l.errors, err)
			}

			continue
		}

		l.unread()
		break
	}

	return WHITESPACE, buf.String()
}

func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return l.eof
	}

	return ch
}

func (l *Lexer) unread() {
	if err := l.r.UnreadRune(); err != nil {
		l.errors = append(l.errors, err)
	}
}

func isWhitespace(ch rune) bool {
	if ch == ' ' {
		return true
	}
	if ch == '\t' {
		return true
	}
	if ch == '\n' {
		return true
	}

	return false
}

func isLetter(ch rune) bool {
	if ch >= 'a' && ch <= 'z' {
		return true
	}
	if ch >= 'A' && ch <= 'Z' {
		return true
	}

	return false
}

func isDigit(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}

	return false
}

func isString(ch rune) bool {
	if ch == '"' {
		return true
	}

	if ch == '\'' {
		return true
	}

	return false
}
