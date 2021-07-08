package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"vmlite/ast"
	"vmlite/lexer"
	"vmlite/parser"
)

const BUG_ERROR = `
   !__!
  (@)(-)
 \.'||'./
-:  ::  :-
/'..''..'\
`
const PROMPT = `
__                     
[  |                    
 | |_   __ _ .--..--.   
 | [ \ [  |  .-. .-. |  
 | |\ \/ / | | | | | |  
[___]\__/ [___||__||__] 
`

func main() {
	repl()
	//debugParser()
}

func repl() {
	fmt.Printf("%s\n", PROMPT)
	fmt.Print("Welcome to lvm (Little virtual machine)\n")
	fmt.Printf("Date and time %v\n", time.Now().Format(time.ANSIC))
	fmt.Print("Type 'quit' to exit.\n")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		input := scanner.Text()
		l := lexer.NewLexer(input)
		p := parser.NewParser(l)
		expr := p.Parse()
		if len(p.Errors()) > 0 {
			printParseErrors(p.Errors())
		}
		o := ast.NewAstPrinter()
		fmt.Printf("%s\n", o.Print(expr))
	}
}

func debugParser() {
	input := `(1 + 2) * 3`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	expr := p.Parse()
	if len(p.Errors()) > 0 {
		printParseErrors(p.Errors())
	}
	o := ast.NewAstPrinter()
	fmt.Printf("%s\n", o.Print(expr))
}

func printParseErrors(errors []string) {
	fmt.Print("Ups! something went wrong in parsing phase!\n")
	fmt.Printf("%s\n", BUG_ERROR)
	for _, err := range errors {
		fmt.Println(err)
	}
}
