package compiler

import (
	"vmlite/ast"
	"vmlite/code"
	"vmlite/token"
)

type Compiler struct {
	co_code   []code.Opcode
	co_consts []interface{}
}

func NewCompiler(co_consts []interface{}) *Compiler {
	c := &Compiler{
		co_code:   []code.Opcode{},
		co_consts: co_consts,
	}
	return c
}

func (c *Compiler) GetCodes() []code.Opcode {
	return c.co_code
}

func (c *Compiler) GetConstants() []interface{} {
	return c.co_consts
}

func (c *Compiler) Compile(expr ast.Expr) {
	c.evaluate(expr)
}

func (c *Compiler) evaluate(expr ast.Expr) interface{} {
	return expr.Accept(c)
}

func (c *Compiler) VisitUnaryExpr(expr *ast.Unary) interface{} {
	c.evaluate(expr.Right)
	switch expr.Operator.Type {
	case token.MINUS:
		c.emit(code.UNEG)
	}
	return nil
}

func (c *Compiler) VisitBinaryExpr(expr *ast.Binary) interface{} {
	c.evaluate(expr.Left)
	c.evaluate(expr.Right)
	switch expr.Operator.Type {
	case token.PLUS:
		c.emit(code.ADD)
	case token.MINUS:
		c.emit(code.SUB)
	case token.MUL:
		c.emit(code.MUL)
	case token.DIV:
		c.emit(code.DIV)
	}
	return nil
}

func (c *Compiler) VisitLiteralExpr(expr *ast.Literal) interface{} {
	switch v := expr.Value.(type) {
	case float64:
		// load constant and get the index
		i := c.addConstant(v)
		// emit the instruction
		c.emit(code.PUSH, i)
	}
	return nil
}

func (c *Compiler) emit(op code.Opcode, args ...int) int {
	ins := code.Make(op, args...)
	i := len(c.co_code)
	c.co_code = append(c.co_code, ins...)
	return i
}

func (c *Compiler) addConstant(cons interface{}) int {
	i := len(c.co_consts)
	c.co_consts = append(c.co_consts, cons)
	return i
}
