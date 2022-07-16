package parser

type Parser struct {
	tokens  []*Token
	current int
}

func New(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() Expr {
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.equality()
}
func (p *Parser) equality() Expr {
	return p.buildBinaryExpr(p.comparison, EQUAL_EQUAL, BANG_EQUAL)
}
func (p *Parser) comparison() Expr {
	return p.buildBinaryExpr(p.term, LESS_EQUAL, LESS, GREATER_EQUAL, GREATER)
}
func (p *Parser) term() Expr {
	return p.buildBinaryExpr(p.factor, MINUS, PLUS)
}
func (p *Parser) factor() Expr {
	return p.buildBinaryExpr(p.unary, STAR, SLASH)
}
func (p *Parser) unary() Expr {
	if op := p.matchAny(BANG, MINUS); op != nil {
		return &UnaryExpr{op, p.unary()}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if tok := p.matchAny(NUMBER, STRING, NIL, TRUE, FALSE); tok != nil {
		return &LiteralExpr{tok}
	}
	if lp := p.matchAny(LEFT_PAREN); lp != nil {
		expr := p.expression()
		if rp := p.matchAny(RIGHT_PAREN); rp == nil {
			panic("missing right paren") // TODO:
		}
		return &GroupingExpr{expr}
	}

	//
	panic("unexpected token")
}

func (p *Parser) buildBinaryExpr(nextLevel func() Expr, tokenTypes ...TokenType) Expr {
	expr := nextLevel()

	for {
		op := p.matchAny(tokenTypes...)
		if op == nil {
			break
		}

		right := nextLevel()
		expr = &BinaryExpr{expr, op, right}
	}

	return expr
}

func (p *Parser) matchAny(tokenTypes ...TokenType) *Token {
	for _, tt := range tokenTypes {
		if p.check(tt) {
			return p.advance()
		}
	}
	return nil
}

func (p *Parser) eof() bool {
	return p.check(EOF)
}

func (p *Parser) check(tt TokenType) bool {
	return p.tokens[p.current].Typ == tt
}

func (p *Parser) advance() *Token {
	t := p.tokens[p.current]
	p.current++
	return t
}
