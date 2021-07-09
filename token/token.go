package token

import "fmt"

type TokenType uint

const (
	IDENT TokenType = iota
	NUMBER
	PLUS
	MINUS
	MUL
	DIV
	LPAREN
	RPAREN
	VAR
	ASSIGN
	PRINT
	EOF
)

var tokenNames = []string{
	"IDENT",
	"NUMBER",
	"PLUS",
	"MINUS",
	"MUL",
	"DIV",
	"LPAREN",
	"RPAREN",
	"VAR",
	"ASSIGN",
	"PRINT",
	"EOF",
}

var symbolMap = map[string]TokenType{
	"+": PLUS,
	"-": MINUS,
	"*": MUL,
	"/": DIV,
	"(": LPAREN,
	")": RPAREN,
	"=": ASSIGN,
}

var keywords = map[string]TokenType{
	"var":   VAR,
	"print": PRINT,
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

func GetKeywordOrIdent(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return IDENT
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
