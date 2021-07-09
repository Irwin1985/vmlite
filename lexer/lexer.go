package lexer

import (
	"fmt"
	"strconv"
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
	lex := string(l.input[pos:l.pos])
	v, ok := strconv.ParseFloat(lex, 64)
	if ok != nil {
		panic(ok)
	}
	return token.NewToken(ln, col, token.NUMBER, v)
}

func (l *Lexer) getIdent() token.Token {
	pos := l.pos
	ln := l.ln
	col := l.col
	for !l.isAtEnd() && l.isIdent(l.c) {
		l.consume()
	}
	v := string(l.input[pos:l.pos])
	return token.NewToken(ln, col, token.GetKeywordOrIdent(v), v)
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
		if l.isIdent(l.c) {
			return l.getIdent()
		}
		if tok, ok := token.IsSymbol(string(l.c)); ok {
			c := string(l.c)
			l.consume()
			return token.NewToken(l.ln, l.col, tok, c)
		}
		panic(fmt.Sprintf("unknown character '%c' at Ln: %d, Col: %d", l.c, l.ln, l.col))
	}
	return token.NewToken(l.ln, l.col, token.EOF, "")
}

func (l *Lexer) isAtEnd() bool {
	return l.c == EOF_CHAR
}

func (l *Lexer) isIdent(c rune) bool {
	return unicode.IsLetter(c) || c == '_'
}
