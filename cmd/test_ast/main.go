package main

import (
	"fmt"
	"lox-go/lox"
)

func main() {
	expression := lox.Binary{
		Left: lox.Unary{
			Operator: lox.Token{Type: lox.MINUS, Lexeme: "-", Literal: nil, Line: 1},
			Right:    lox.Literal{Value: 123},
		},
		Operator: lox.Token{Type: lox.STAR, Lexeme: "*", Literal: nil, Line: 1},
		Right: lox.Grouping{Expression: lox.Literal{Value: 45.67}},
	}

	astPrinter := lox.NewAstPrinter()
	fmt.Println(astPrinter.Print(expression))
}
