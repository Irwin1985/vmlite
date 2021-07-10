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
	errors    []string
	lop       byte
}

func NewCompiler(co_names []string, co_consts []interface{}) *Compiler {
	c := &Compiler{
		co_code:   []code.Opcode{},
		co_names:  co_names,
		co_consts: co_consts,
		co_values: []interface{}{},
		errors:    []string{},
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

	c.emit(code.STORE, float32(i))
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
	op := c.lop
	c.emit(code.UNARY)
	switch expr.Operator.Type {
	case token.MINUS:
		if typeOf(op) == 'f' {
			c.emit(code.UNEG)
		} else {
			c.addError("the [minus] operator only works with numeric types.")
		}
	case token.NOT:
		if typeOf(op) == 'l' {
			c.emit(code.NOT)
		} else {
			c.addError("the [not] operator only works with boolean types.")
		}
	}
	return nil
}

func (c *Compiler) VisitBinaryExpr(expr *ast.Binary) interface{} {
	c.evaluateExpr(expr.Left)
	op1 := c.lop
	c.evaluateExpr(expr.Right)
	op2 := c.lop

	// DEBUG
	fmt.Printf("leftOp: %v, rightOp: %v typeLeft: %v, typeRight: %v\n", code.CodeMap[op1], code.CodeMap[op2], string(typeOf(op1)), string(typeOf(op2)))
	// DEBUG

	c.emit(code.CMP) // COMPARE
	t := expr.Operator.Type
	if typeOf(op1) == 's' && typeOf(op2) == 's' {
		switch t {
		case token.PLUS:
			c.emit(code.ADDS)
		case token.MINUS:
			c.emit(code.SUBS)
		default:
			c.addError("unsupported operator for string type.")
		}
	} else if typeOf(op1) == 'f' && typeOf(op2) == 'f' {
		switch t {
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
		default:
			c.addError("unsupported operator for string type.")
		}
	} else if typeOf(op1) == 'l' && typeOf(op2) == 'l' {
		switch t {
		case token.AND:
			c.emit(code.AND)
		case token.OR:
			c.emit(code.OR)
		default:
			c.addError("unsupported operator for boolean type.")
		}
	} else {
		c.addError("invalid operands.")
	}
	return nil
}

func (c *Compiler) VisitLiteralExpr(expr *ast.Literal) interface{} {
	t := expr.Token.Type
	switch t {
	case token.TRUE, token.FALSE:
		c.emit(code.BOOL)
		if t == token.TRUE {
			c.emit(code.TRUE)
		} else {
			c.emit(code.FALSE)
		}
	case token.STRING:
		i := c.addConstant(expr.Token.Lexeme.(string))
		c.emit(code.PUSHS, float32(i))

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
		c.emit(code.LOAD, float32(i))
	case token.NUMBER:
		c.emit(code.PUSHF, expr.Token.Lexeme.(float32))
	}
	return nil
}

/****************************
* COMPILER HELPER FUNCTIONS
*****************************/

// emits a byte instruction
func (c *Compiler) emit(op code.Opcode, args ...float32) int {
	ins := code.Make(op, args...)
	i := len(c.co_code)
	c.co_code = append(c.co_code, ins...)

	c.lop = op // save last op

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

func typeOf(op code.Opcode) byte {
	if op == code.PUSHS || op == code.ADDS || op == code.SUBS {
		return 's'
	} else if op == code.PUSHF || op == code.ADD || op == code.SUB || op == code.MUL || op == code.DIV || op == code.LT ||
		op == code.LEQ || op == code.GT || op == code.GEQ || op == code.EQ || op == code.NEQ || op == code.UNEG {
		return 'f'
	} else if op == code.TRUE || op == code.FALSE || op == code.AND || op == code.OR || op == code.NOT {
		return 'l'
	}
	return 'u'
}
