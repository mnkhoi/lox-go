package lox

import (
	"fmt"
	"strconv"
	"strings"
)

type AstPrinter struct {
}

func (a AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) VisitExprBinary(binary Binary) any {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a AstPrinter) VisitExprGrouping(grouping Grouping) any {
	return a.parenthesize("group", grouping.Expression)
}

func (a AstPrinter) VisitExprLiteral(literal Literal) any {
	if literal.Value == nil {
		return "nil"
	}

	switch l := literal.Value.(type) {
	case float64:
		return strconv.FormatFloat(l, 'f', -1, 64)
	default:
		return fmt.Sprint(l)
	}
}

func (a AstPrinter) VisitExprUnary(unary Unary) any {
	return a.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("(%s", name))
	for _, expr := range exprs {
		builder.WriteRune(' ')
		builder.WriteString(expr.Accept(a).(string))
	}
	builder.WriteRune(')')

	return builder.String()
}

func NewAstPrinter() *AstPrinter{
	return &AstPrinter{}
}
