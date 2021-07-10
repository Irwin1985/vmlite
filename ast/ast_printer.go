package ast

import (
	"bytes"
	"fmt"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(program []Stmt) string {
	var out bytes.Buffer
	for _, stmt := range program {
		out.WriteString(fmt.Sprintf("%v\n", a.executeStmt(stmt)))
	}
	return out.String()
}

func (a *AstPrinter) executeStmt(stmt Stmt) interface{} {
	return stmt.Accept(a)
}

func (a *AstPrinter) VisitVarStmt(stmt *VarStmt) interface{} {
	return fmt.Sprintf("var %v = %v", stmt.Name.Lexeme, a.evaluateExpr(stmt.Value))
}

func (a *AstPrinter) VisitExprStmt(stmt *ExprStmt) interface{} {
	return a.evaluateExpr(stmt.Expression)
}

func (a *AstPrinter) VisitPrintStmt(stmt *PrintStmt) interface{} {
	return fmt.Sprintf("print(%v)", a.evaluateExpr(stmt.Value))
}

// Expression evaluator
func (a *AstPrinter) evaluateExpr(e Expr) interface{} {
	return e.Accept(a)
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return fmt.Sprintf("(%v %v)", expr.Operator.Lexeme, a.evaluateExpr(expr.Right))
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return fmt.Sprintf("(%v %v %v)", a.evaluateExpr(expr.Left), expr.Operator.Lexeme, a.evaluateExpr(expr.Right))
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if v, ok := expr.Token.Lexeme.(string); ok {
		return fmt.Sprintf("'%s'", v)
	}
	return expr.Token.Lexeme
}

// func (a *AstPrinter) VisitIdentifierExpr(expr *Identifier) interface{} {
// 	return expr.Value.Lexeme
// }
