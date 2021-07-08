package ast

import "fmt"

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(e Expr) string {
	return fmt.Sprintf("%v", a.evaluate(e))
}

func (a *AstPrinter) evaluate(e Expr) interface{} {
	return e.Accept(a)
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return fmt.Sprintf("(%v %v)", expr.Operator.Lexeme, a.evaluate(expr.Right))
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return fmt.Sprintf("(%v %v %v)", a.evaluate(expr.Left), expr.Operator.Lexeme, a.evaluate(expr.Right))
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}
