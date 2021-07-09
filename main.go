package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"vmlite/ast"
	"vmlite/compiler"
	"vmlite/lexer"
	"vmlite/parser"
	"vmlite/vm"
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
	//debugCompiler()
	//debugVM()
}

func repl() {
	co_consts := []interface{}{}
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
		if len(input) == 0 {
			continue
		}
		if input == "quit" {
			break
		}
		l := lexer.NewLexer(input)
		p := parser.NewParser(l)
		expr := p.Parse()
		if len(p.Errors()) > 0 {
			printParseErrors(p.Errors())
		}
		c := compiler.NewCompiler(co_consts)
		c.Compile(expr)
		co_codes := c.GetCodes()
		co_consts = c.GetConstants()
		vm := vm.NewVM(co_codes, co_consts)
		err := vm.Run()
		if err != nil {
			panic(err)
		}
		tos := vm.TOS()
		fmt.Printf("%v\n", tos)
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

func debugCompiler() {
	co_consts := []interface{}{}
	input := `(1 + 2) * 3`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	expr := p.Parse()
	if len(p.Errors()) > 0 {
		printParseErrors(p.Errors())
		return
	}
	c := compiler.NewCompiler(co_consts)
	c.Compile(expr)
	co_codes := c.GetCodes()
	co_consts = c.GetConstants()
	vm := vm.NewVM(co_codes, co_consts)
	err := vm.Run()
	if err != nil {
		panic(err)
	}
	tos := vm.TOS()
	fmt.Printf("%v\n", tos)
}

func debugVM() {
	co_consts := []interface{}{}
	input := `2 * -3`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	expr := p.Parse()
	if len(p.Errors()) > 0 {
		printParseErrors(p.Errors())
		return
	}
	c := compiler.NewCompiler(co_consts)
	c.Compile(expr)
	co_codes := c.GetCodes()
	co_consts = c.GetConstants()
	vm := vm.NewVM(co_codes, co_consts)
	err := vm.Run()
	if err != nil {
		panic(err)
	}
	tos := vm.TOS()
	fmt.Printf("%v\n", tos)
}

func printParseErrors(errors []string) {
	fmt.Print("Ups! something went wrong in parsing phase!\n")
	fmt.Printf("%s\n", BUG_ERROR)
	for _, err := range errors {
		fmt.Println(err)
	}
}
