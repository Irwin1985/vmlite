package lexer

import (
	"fmt"
	"unicode"
	"vmlite/token"
)

var EOF_CHAR = rune(0)

type Lexer struct {
	pos   int
	c     rune
	ln    int
	col   int
	input []rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		pos:   -1,
		ln:    1,
		col:   0,
		input: []rune(input),
	}
	l.consume() // prime first char
	return l
}

func (l *Lexer) consume() {
	l.pos += 1
	if l.pos >= len(l.input) {
		l.c = EOF_CHAR
		return
	}
	l.c = l.input[l.pos]
	if l.c == '\n' {
		l.ln += 1
		l.col = 0
		return
	}
	l.col += 1
}

func (l *Lexer) ws() {
	for !l.isAtEnd() && unicode.IsSpace(l.c) {
		l.consume()
	}
}

func (l *Lexer) getNum() token.Token {
	pos := l.pos
	ln := l.ln
	col := l.col
	for !l.isAtEnd() && unicode.IsNumber(l.c) {
		l.consume()
	}
	return token.Token{Ln: ln, Col: col, Type: token.NUMBER, Lexeme: string(l.input[pos:l.pos])}
}

func (l *Lexer) NextToken() token.Token {
	for !l.isAtEnd() {
		if unicode.IsSpace(l.c) {
			l.ws()
			continue
		}
		if unicode.IsNumber(l.c) {
			return l.getNum()
		}
		if tok, ok := token.IsSymbol(string(l.c)); ok {
			c := string(l.c)
			l.consume()
			return token.Token{Ln: l.ln, Col: l.col, Type: tok, Lexeme: c}
		}
		panic(fmt.Sprintf("unknown character '%c' at Ln: %d, Col: %d", l.c, l.ln, l.col))
	}
	return token.Token{Ln: l.ln, Col: l.col, Type: token.EOF, Lexeme: ""}
}

func (l *Lexer) isAtEnd() bool {
	return l.c == EOF_CHAR
}
