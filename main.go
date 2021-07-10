package main

import "vmlite/repl"

func main() {
	mode := "repl"
	input := `print 1 + 2`
	repl.Start(mode, input)
}
