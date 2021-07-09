package repl

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"vmlite/ast"
	"vmlite/compiler"
	"vmlite/lexer"
	"vmlite/parser"
	"vmlite/token"
	"vmlite/vm"
)

const VERSION = "1.0"
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

const VALUES_SIZE = 65536

var co_consts = []interface{}{}
var co_names = []string{}
var co_values = make([]interface{}, VALUES_SIZE)

func Start(mode string, input string) {
	if mode == "repl" {
		repl()
	} else if mode == "lexer" {
		debugLexer(input)
	} else if mode == "parser" {
		debugParser(input)
	} else if mode == "compiler" {
		debugCompiler(input)
	} else if mode == "vm" {
		debugVM(input)
	}
}

func repl() {
	displayWelcome()

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
		run(input)
	}
}

func run(input string) {

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Program()
	if len(p.Errors()) > 0 {
		printErrors(p.Errors())
	}

	c := compiler.NewCompiler(co_names, co_consts)
	c.Compile(program)
	errors := c.Errors()
	if len(errors) > 0 {
		printErrors(errors)
		return
	}

	co_codes := c.GetCodes()
	co_names = c.GetNames()
	co_consts = c.GetConstants()

	// debug
	//fmt.Printf("co_names[%v]\nco_consts[%v]\n", co_names, co_consts)
	// debug

	vm := vm.NewVM(co_codes, co_consts, co_names, co_values)
	err := vm.Run()
	if err != nil {
		fmt.Print(err)
	}
}

func debugLexer(input string) {
	l := lexer.NewLexer(input)
	tok := l.NextToken()
	for tok.Type != token.EOF {
		fmt.Println(tok.ToString())
		tok = l.NextToken()
	}
}

func debugParser(input string) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Program()
	if len(p.Errors()) > 0 {
		printErrors(p.Errors())
	}
	o := ast.NewAstPrinter()
	fmt.Printf("%s\n", o.Print(program))
}

func debugCompiler(input string) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Program()
	if len(p.Errors()) > 0 {
		printErrors(p.Errors())
		return
	}
	c := compiler.NewCompiler(co_names, co_consts)
	c.Compile(program)
	errors := c.Errors()
	if len(errors) > 0 {
		printErrors(errors)
		return
	}

	co_codes := c.GetCodes()
	co_names = c.GetNames()
	co_consts = c.GetConstants()

	output := compiler.PrintByteCode(co_codes, co_consts)
	fmt.Printf("%v\n", output)
}

func debugVM(input string) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Program()
	if len(p.Errors()) > 0 {
		printErrors(p.Errors())
		return
	}
	c := compiler.NewCompiler(co_names, co_consts)
	c.Compile(program)
	errors := c.Errors()
	if len(errors) > 0 {
		printErrors(errors)
		return
	}

	co_codes := c.GetCodes()
	co_names = c.GetNames()
	co_consts = c.GetConstants()

	vm := vm.NewVM(co_codes, co_consts, co_names, co_values)
	err := vm.Run()
	if err != nil {
		panic(err)
	}
	tos := vm.TOS()
	if tos != nil {
		fmt.Printf("%v\n", tos)
	}
}

func printErrors(errors []string) {
	fmt.Print("Ups! something went wrong!\n")
	fmt.Printf("%s\n", BUG_ERROR)
	for _, err := range errors {
		fmt.Println(err)
	}
}

func displayWelcome() {
	fmt.Printf("%s\n", PROMPT)
	fmt.Printf("Welcome to LVM (Little virtual machine) Version: %s \n", VERSION)
	fmt.Printf("Date and time %v\n", time.Now().Format(time.ANSIC))
	fmt.Print("Type 'quit' to exit.\n")
}
