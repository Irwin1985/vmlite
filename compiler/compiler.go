package compiler

import (
	"fmt"
	"vmlite/ast"
	"vmlite/code"
	"vmlite/token"
)

type Compiler struct {
	co_code   []code.Opcode
	co_names  []string
	co_consts []interface{}
	co_values []interface{}
	map_num   map[float64]int
	errors    []string
}

func NewCompiler(co_names []string, co_consts []interface{}, map_num map[float64]int) *Compiler {
	c := &Compiler{
		co_code:   []code.Opcode{},
		co_names:  co_names,
		co_consts: co_consts,
		co_values: []interface{}{},
		errors:    []string{},
		map_num:   map_num,
	}
	return c
}

func (c *Compiler) GetCodes() []code.Opcode {
	return c.co_code
}

func (c *Compiler) GetConstants() []interface{} {
	return c.co_consts
}

func (c *Compiler) GetNames() []string {
	return c.co_names
}

func (c *Compiler) Errors() []string {
	return c.errors
}

func (c *Compiler) Compile(program []ast.Stmt) {
	for _, stmt := range program {
		c.executeStmt(stmt)
	}
}

// Statements Visitor and Executor
func (c *Compiler) executeStmt(stmt ast.Stmt) interface{} {
	return stmt.Accept(c)
}

func (c *Compiler) VisitVarStmt(stmt *ast.VarStmt) interface{} {
	c.evaluateExpr(stmt.Value)
	i := c.addName(stmt.Name.Lexeme.(string))

	c.emit(code.STORE, i)
	return nil
}

func (c *Compiler) VisitExprStmt(stmt *ast.ExprStmt) interface{} {
	return c.evaluateExpr(stmt.Expression)
}

func (c *Compiler) VisitPrintStmt(stmt *ast.PrintStmt) interface{} {
	c.evaluateExpr(stmt.Value)
	c.emit(code.PRINT)
	return nil
}

// Expressions Visitor and Evaluator
func (c *Compiler) evaluateExpr(expr ast.Expr) interface{} {
	return expr.Accept(c)
}

func (c *Compiler) VisitUnaryExpr(expr *ast.Unary) interface{} {
	c.evaluateExpr(expr.Right)
	c.emit(code.UNARY)
	switch expr.Operator.Type {
	case token.MINUS:
		c.emit(code.UNEG)
	case token.NOT:
		c.emit(code.NOT)
	}
	return nil
}

func (c *Compiler) VisitBinaryExpr(expr *ast.Binary) interface{} {
	c.evaluateExpr(expr.Left)
	c.evaluateExpr(expr.Right)
	c.emit(code.CMP) // COMPARE
	switch expr.Operator.Type {
	case token.PLUS:
		c.emit(code.ADD)
	case token.MINUS:
		c.emit(code.SUB)
	case token.MUL:
		c.emit(code.MUL)
	case token.DIV:
		c.emit(code.DIV)
	case token.LT:
		c.emit(code.LT)
	case token.GT:
		c.emit(code.GT)
	case token.LEQ:
		c.emit(code.LEQ)
	case token.GEQ:
		c.emit(code.GEQ)
	case token.EQ:
		c.emit(code.EQ)
	case token.NEQ:
		c.emit(code.NEQ)
	case token.AND:
		c.emit(code.AND)
	case token.OR:
		c.emit(code.OR)
	}
	return nil
}

func (c *Compiler) VisitLiteralExpr(expr *ast.Literal) interface{} {
	t := expr.Token.Type
	switch t {
	case token.TRUE, token.FALSE:
		c.emit(code.BOOL)
		v := expr.Token.Lexeme.(bool)
		if v {
			c.emit(code.TRUE)
		} else {
			c.emit(code.FALSE)
		}
	case token.STRING:
		i := c.addConstant(expr.Token.Lexeme.(string))
		c.emit(code.PUSH, i)

	case token.IDENT:
		name := expr.Token.Lexeme.(string)
		i := -1
		for j := 0; j < len(c.co_names); j++ {
			if c.co_names[j] == name {
				i = j
				break
			}
		}
		if i < 0 {
			c.addError(fmt.Sprintf("Variable not defined: %s\n", name))
			return nil
		}
		c.emit(code.LOAD, i)
	case token.NUMBER:
		v := expr.Token.Lexeme.(float64)
		// lookup in map_num first
		if i, ok := c.map_num[v]; ok {
			c.emit(code.PUSH, i)
		} else { // feed the map_num with new consts
			i := c.addConstant(v)
			c.map_num[v] = i
			c.emit(code.PUSH, i) // save reference
		}
	}
	return nil
}

/****************************
* COMPILER HELPER FUNCTIONS
*****************************/

// emits a byte instruction
func (c *Compiler) emit(op code.Opcode, args ...int) int {
	ins := code.Make(op, args...)
	i := len(c.co_code)
	c.co_code = append(c.co_code, ins...)
	return i
}

// add a constant in co_consts array
func (c *Compiler) addConstant(cons interface{}) int {
	var i int = len(c.co_consts)
	c.co_consts = append(c.co_consts, cons)
	return i
}

// add a name (symbol) in co_names (global for now)
func (c *Compiler) addName(name string) int {
	// find the index
	i := len(c.co_names)
	for j := 0; j < len(c.co_names); j++ {
		if c.co_names[j] == name {
			i = j
			break
		}
	}
	// add name(symbol)
	c.co_names = append(c.co_names, name)

	return i
}

// add error into array
func (c *Compiler) addError(msg string) {
	c.errors = append(c.errors, msg)
}
