package main

import "vmlite/repl"

func main() {
	mode := "repl"
	input := `5 + 5 * 3 * 3 * 4`
	repl.Start(mode, input)
}
