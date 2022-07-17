package parser

import "fmt"

type Parser struct {
	tokens  []*Token
	current int
}

type SyntaxError struct {
	Msg   string
	Token *Token
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("syntax error: %s, line %d: '%s'", e.Msg, e.Token.Line, e.Token.Lexeme)
}

func New(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}
func (p *Parser) equality() (Expr, error) {
	return p.buildBinaryExpr(p.comparison, EQUAL_EQUAL, BANG_EQUAL)
}
func (p *Parser) comparison() (Expr, error) {
	return p.buildBinaryExpr(p.term, LESS_EQUAL, LESS, GREATER_EQUAL, GREATER)
}
func (p *Parser) term() (Expr, error) {
	return p.buildBinaryExpr(p.factor, MINUS, PLUS)
}
func (p *Parser) factor() (Expr, error) {
	return p.buildBinaryExpr(p.unary, STAR, SLASH)
}
func (p *Parser) unary() (Expr, error) {
	if op := p.matchAny(BANG, MINUS); op != nil {
		u, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{op, u}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if tok := p.matchAny(NUMBER, STRING, NIL, TRUE, FALSE); tok != nil {
		return &LiteralExpr{tok}, nil
	}
	if lp := p.matchAny(LEFT_PAREN); lp != nil {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if rp := p.matchAny(RIGHT_PAREN); rp == nil {
			return nil, p.genSyntaxError("missing closing parenthesis")
		}
		return &GroupingExpr{expr}, nil
	}

	return nil, p.genSyntaxError("unexpected token")
}

func (p *Parser) buildBinaryExpr(nextLevel func() (Expr, error), tokenTypes ...TokenType) (Expr, error) {
	expr, err := nextLevel()
	if err != nil {
		return nil, err
	}

	for {
		op := p.matchAny(tokenTypes...)
		if op == nil {
			break
		}

		right, err := nextLevel()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{expr, op, right}
	}

	return expr, nil
}

func (p *Parser) matchAny(tokenTypes ...TokenType) *Token {
	for _, tt := range tokenTypes {
		if p.check(tt) {
			return p.advance()
		}
	}
	return nil
}
func (p *Parser) genSyntaxError(format string, v ...interface{}) *SyntaxError {
	return &SyntaxError{
		Msg:   fmt.Sprintf(format, v...),
		Token: p.tokens[p.current],
	}
}

func (p *Parser) check(tt TokenType) bool {
	return p.tokens[p.current].Typ == tt
}

func (p *Parser) advance() *Token {
	t := p.tokens[p.current]
	p.current++
	return t
}
