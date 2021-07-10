package ast

import "vmlite/token"

type VisitorExpr interface {
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitBinaryExpr(expr *Binary) interface{}
	// VisitIdentifierExpr(expr *Identifier) interface{}
}

type Expr interface {
	Accept(v VisitorExpr) interface{}
}

type Literal struct {
	Token token.Token
}

func (expr *Literal) Accept(v VisitorExpr) interface{} {
	return v.VisitLiteralExpr(expr)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (expr *Unary) Accept(v VisitorExpr) interface{} {
	return v.VisitUnaryExpr(expr)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (expr *Binary) Accept(v VisitorExpr) interface{} {
	return v.VisitBinaryExpr(expr)
}

// type Identifier struct {
// 	Value token.Token
// }

// func (expr *Identifier) Accept(v VisitorExpr) interface{} {
// 	return v.VisitIdentifierExpr(expr)
// }
