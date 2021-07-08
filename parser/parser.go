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
	TERM
	FACTOR
)

// precedence map
var mapPrecedence = map[token.TokenType]int{
	token.PLUS:  TERM,
	token.MINUS: TERM,
	token.MUL:   FACTOR,
	token.DIV:   FACTOR,
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
	p.registerPrefixFn(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpr)
	// register INFIX semantic code
	p.registerInfixFn(token.PLUS, p.parseInfixExpr)
	p.registerInfixFn(token.MINUS, p.parseInfixExpr)
	p.registerInfixFn(token.MUL, p.parseInfixExpr)
	p.registerInfixFn(token.DIV, p.parseInfixExpr)
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

func (p *Parser) Parse() ast.Expr {
	return p.expression(LOWEST)
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

func (p *Parser) parseNumberLiteral() ast.Expr {
	expr := &ast.Literal{}
	if p.match(token.NUMBER) {
		expr.Value = p.prevToken.Lexeme
	} else {
		p.newError(fmt.Sprintf("%v expect expression", p.curToken))
	}
	return expr
}

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.nextToken()
	exp := p.expression(LOWEST)
	p.expect(token.RPAREN, "expect ')' after expression.")

	return exp
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