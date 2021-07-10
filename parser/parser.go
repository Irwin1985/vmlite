package parser

import (
	"fmt"
	"vmlite/ast"
	"vmlite/lexer"
	"vmlite/token"
)

// precedence order
const (
	LOWEST int = iota
	LOGIC_OR
	LOGIC_AND
	EQUALITY
	COMPARISON
	TERM
	FACTOR
	PREFIX
	CALL
	INDEX
)

// precedence map
var mapPrecedence = map[token.TokenType]int{
	// logical operators
	token.OR:  LOGIC_OR,
	token.AND: LOGIC_AND,
	// equality
	token.EQ:  EQUALITY,
	token.NEQ: EQUALITY,
	// comparison
	token.LT:  COMPARISON,
	token.GT:  COMPARISON,
	token.LEQ: COMPARISON,
	token.GEQ: COMPARISON,
	// term
	token.PLUS:  TERM,
	token.MINUS: TERM,
	// factor
	token.MUL: FACTOR,
	token.DIV: FACTOR,
}

// semantic function types
type PrefixFnType func() ast.Expr
type InfixFnType func(ast.Expr) ast.Expr

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	prevToken token.Token
	peekToken token.Token
	errors    []string
	// semantic map
	mapPrefixFn map[token.TokenType]PrefixFnType
	mapInfixFn  map[token.TokenType]InfixFnType
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:           l,
		mapPrefixFn: make(map[token.TokenType]PrefixFnType),
		mapInfixFn:  make(map[token.TokenType]InfixFnType),
	}

	// register PREFIX semantic code
	p.registerPrefixFn(token.NUMBER, p.parseLiteral)
	p.registerPrefixFn(token.STRING, p.parseLiteral)
	p.registerPrefixFn(token.IDENT, p.parseLiteral)
	p.registerPrefixFn(token.TRUE, p.parseLiteral)
	p.registerPrefixFn(token.FALSE, p.parseLiteral)

	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpr)
	p.registerPrefixFn(token.MINUS, p.parseUnaryExpr)
	p.registerPrefixFn(token.NOT, p.parseUnaryExpr)

	// register INFIX semantic code
	p.registerInfixFn(token.PLUS, p.parseInfixExpr)
	p.registerInfixFn(token.MINUS, p.parseInfixExpr)
	p.registerInfixFn(token.MUL, p.parseInfixExpr)
	p.registerInfixFn(token.DIV, p.parseInfixExpr)
	p.registerInfixFn(token.OR, p.parseInfixExpr)
	p.registerInfixFn(token.AND, p.parseInfixExpr)
	p.registerInfixFn(token.LT, p.parseInfixExpr)
	p.registerInfixFn(token.GT, p.parseInfixExpr)
	p.registerInfixFn(token.LEQ, p.parseInfixExpr)
	p.registerInfixFn(token.GEQ, p.parseInfixExpr)
	p.registerInfixFn(token.EQ, p.parseInfixExpr)
	p.registerInfixFn(token.NEQ, p.parseInfixExpr)

	// move tokens
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefixFn(t token.TokenType, fn PrefixFnType) {
	p.mapPrefixFn[t] = fn
}

func (p *Parser) registerInfixFn(t token.TokenType, fn InfixFnType) {
	p.mapInfixFn[t] = fn
}

func (p *Parser) curPrecedence() int {
	if pre, ok := mapPrecedence[p.curToken.Type]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Program() []ast.Stmt {
	stmt := []ast.Stmt{}
	for !p.match(token.EOF) {
		stmt = append(stmt, p.statement())
	}
	return stmt
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.VAR) {
		return p.varStatement()
	} else if p.match(token.PRINT) {
		return p.printStmt()
	} else {
		return p.exprStmt()
	}
}

func (p *Parser) varStatement() ast.Stmt {
	stmt := &ast.VarStmt{}
	p.expect(token.IDENT, "expect IDENTIFIER after 'var' declaration.")
	stmt.Name = p.prevToken

	p.expect(token.ASSIGN, "expect '=' before expression.")
	stmt.Value = p.expression(LOWEST)

	return stmt
}

func (p *Parser) printStmt() ast.Stmt {
	stmt := &ast.PrintStmt{}
	stmt.Value = p.expression(LOWEST)

	return stmt
}

func (p *Parser) exprStmt() ast.Stmt {
	stmt := &ast.ExprStmt{}
	stmt.Expression = p.expression(LOWEST)

	return stmt
}

func (p *Parser) expression(precedence int) ast.Expr {
	prefixFn := p.mapPrefixFn[p.curToken.Type]
	if prefixFn == nil {
		p.newError(fmt.Sprintf("%v no prefix parsing function for this token.", p.curToken))
		return nil
	}
	leftExpr := prefixFn()
	for precedence < p.curPrecedence() {
		infixFn := p.mapInfixFn[p.curToken.Type]
		if infixFn == nil {
			return leftExpr
		}
		leftExpr = infixFn(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parseLiteral() ast.Expr {
	expr := &ast.Literal{Token: p.curToken}
	p.nextToken()
	return expr
}

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.nextToken()
	exp := p.expression(LOWEST)
	p.expect(token.RPAREN, "expect ')' after expression.")

	return exp
}

func (p *Parser) parseUnaryExpr() ast.Expr {
	expr := &ast.Unary{Operator: p.curToken}
	p.nextToken()
	expr.Right = p.expression(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.Binary{
		Left:     left,
		Operator: p.curToken,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.expression(precedence)

	return expr
}

func (p *Parser) expect(t token.TokenType, msg string) {
	if p.match(t) {
		return
	}
	p.newError(msg)
}

func (p *Parser) match(t token.TokenType) bool {
	if p.curToken.Type == t {
		p.prevToken = p.curToken
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) newError(msg string) {
	p.errors = append(p.errors, msg)
}
