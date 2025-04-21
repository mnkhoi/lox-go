package lox

type Parser struct {
	Tokens  []Token
	Current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (p *Parser) Expression() Expr {
	return p.Equality()
}

func (p *Parser) Equality() Expr {
	var expr Expr
	expr = p.Comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.Comparison()
		expr = Binary{
			Left:     expr,
			Operator: *operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) Comparison() Expr {
	var expr Expr
	expr = p.Term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.Term()
		expr = Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) Term() Expr {
	var expr Expr
	expr = p.Factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.Factor()
		expr = Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) Factor() Expr {
	var expr Expr
	expr = p.Unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.Unary()
		expr = Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) Unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.Unary()
		return Unary{
			Operator: operator,
			Right:    right,
		}
	}
	return p.Primary()
}

func (p *Parser) Primary() Expr {
	if p.match(FALSE) {
		return Literal{
			Value: false,
		}
	} else if p.match(TRUE) {
		return Literal{
			Value: true,
		}
	} else if p.match(NIL) {
		return Literal{
			Value: nil,
		}
	} else if p.match(NUMBER, STRING) {
		return Literal{
			Value: p.previous().Literal,
		}
	} else if p.match(LEFT_PAREN) {
		expr := p.Expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{
			Expression: expr,
		}
	}

	// TODO: Change this to more appropriate error handling
	return Literal{
		Value: nil,
	}
}

func (p *Parser) consume(ttype TokenType, message string) *Token {
	if p.check(ttype) {
		return & (p.advance())
	}

	
	return nil
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == ttype
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() *Token {
	return &p.Tokens[p.Current-1]
}
