package ast

import "vmlite/token"

type VisitorStmt interface {
	VisitVarStmt(stmt *VarStmt) interface{}
	VisitExprStmt(stmt *ExprStmt) interface{}
	VisitPrintStmt(stmt *PrintStmt) interface{}
}

type Stmt interface {
	Accept(v VisitorStmt) interface{}
}

type VarStmt struct {
	Name  token.Token
	Value Expr
}

func (stmt *VarStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitVarStmt(stmt)
}

type ExprStmt struct {
	Expression Expr
}

func (stmt *ExprStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitExprStmt(stmt)
}

type PrintStmt struct {
	Value Expr
}

func (stmt *PrintStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitPrintStmt(stmt)
}
