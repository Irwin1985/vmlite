package token

import "fmt"

type TokenType uint

const (
	NUMBER TokenType = iota
	PLUS
	MINUS
	MUL
	DIV
	LPAREN
	RPAREN
	EOF
)

var tokenNames = []string{
	"NUMBER",
	"PLUS",
	"MINUS",
	"MUL",
	"DIV",
	"LPAREN",
	"RPAREN",
	"EOF",
}

var symbolMap = map[string]TokenType{
	"+": PLUS,
	"-": MINUS,
	"*": MUL,
	"/": DIV,
	"(": LPAREN,
	")": RPAREN,
}

type Token struct {
	Type   TokenType
	Lexeme interface{}
	Ln     int
	Col    int
}

func (t Token) ToString() string {
	return fmt.Sprintf("Ln %d, Col: %d -> <%s, '%v'>", t.Ln, t.Col, tokenNames[t.Type], t.Lexeme)
}

func IsSymbol(k string) (TokenType, bool) {
	if v, ok := symbolMap[k]; ok {
		return v, ok
	}
	return EOF, false
}

func NewToken(ln int, col int, t TokenType, v interface{}) Token {
	tok := Token{
		Ln:     ln,
		Col:    col,
		Type:   t,
		Lexeme: v,
	}
	return tok
}
