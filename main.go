package main

import "vmlite/repl"

func main() {
	mode := "repl"
	input := `!true and !false`
	repl.Start(mode, input)
}
