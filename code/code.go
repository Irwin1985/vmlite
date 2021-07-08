package code

type code = byte

const (
	push code = iota
	add
	sub
	mul
	div
	halt
)

var codeMap = map[string]code{
	"push": push,
	"add":  add,
	"sub":  sub,
	"mul":  mul,
	"div":  div,
	"halt": halt,
}
