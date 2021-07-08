package main

import (
	"fmt"
	"vmlite/lexer"
	"vmlite/token"
)

func main() {
	input := `1 + 2 * 3
	(2 + 3) * 5
	7 * 2 / 4
	`
	l := lexer.NewLexer(input)
	tok := l.NextToken()
	for tok.Type != token.EOF {
		fmt.Println(tok.ToString())
		tok = l.NextToken()
	}
	fmt.Println(tok.ToString())
}
