package token

import "fmt"

type TokenType uint

const (
	IDENT TokenType = iota
	NUMBER
	STRING

	// single characters
	PLUS
	MINUS
	MUL
	DIV
	LPAREN
	RPAREN
	ASSIGN

	// comparison
	LT
	GT
	LEQ
	GEQ
	EQ
	NOT
	NEQ

	// keywords
	VAR
	PRINT
	TRUE
	FALSE
	AND
	OR
	EOF
)

var tokenNames = []string{
	"IDENT",
	"NUMBER",
	"STRING",
	"PLUS",
	"MINUS",
	"MUL",
	"DIV",
	"LPAREN",
	"RPAREN",
	"ASSIGN",
	"LT",
	"GT",
	"LEQ",
	"GEQ",
	"EQ",
	"NOT",
	"NEQ",
	"VAR",
	"PRINT",
	"TRUE",
	"FALSE",
	"AND",
	"OR",
	"EOF",
}

var symbolMap = map[string]TokenType{
	"+":  PLUS,
	"-":  MINUS,
	"*":  MUL,
	"/":  DIV,
	"(":  LPAREN,
	")":  RPAREN,
	"=":  ASSIGN,
	"<":  LT,
	">":  GT,
	"<=": LEQ,
	">=": GEQ,
	"==": EQ,
	"!":  NOT,
	"!=": NEQ,
}

var keywords = map[string]TokenType{
	"var":   VAR,
	"print": PRINT,
	"true":  TRUE,
	"false": FALSE,
	"and":   AND,
	"or":    OR,
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
