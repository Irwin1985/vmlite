package main

import "vmlite/repl"

func main() {
	mode := "repl"
	input := `print 12 + 15 + 12 + 15 + 10`
	repl.Start(mode, input)
}
