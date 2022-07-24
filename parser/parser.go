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

func (p *Parser) Parse() ([]Stmt, error) {
	var stmts []Stmt
	for !p.check(EOF) {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.matchAny(VAR) != nil {
		return p.varDeclaration()
	}
	//TODO: SYNCHRONISE
	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name := p.matchAny(IDENTIFIER)
	if name == nil {
		return nil, p.genSyntaxError("missing variable name")
	}

	var init Expr
	if p.matchAny(EQUAL) != nil {
		in, err := p.expression()
		if err != nil {
			return nil, err
		}
		init = in
	}

	if p.matchAny(SEMICOLON) == nil {
		return nil, p.genSyntaxError("missing semicolon after value")
	}

	return &VarDecl{name: name, init: init}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.matchAny(PRINT) != nil {
		return p.printStmt()
	}
	if p.matchAny(LEFT_BRACE) != nil {
		return p.blockStmt()
	}
	return p.exprStmt()
}

func (p *Parser) blockStmt() (Stmt, error) {
	var lst []Stmt

	for !p.check(RIGHT_BRACE) && !p.check(EOF) {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		lst = append(lst, stmt)
	}
	if p.matchAny(RIGHT_BRACE) == nil {
		return nil, p.genSyntaxError("missing closing } after block")
	}

	return &Block{statements: lst}, nil
}

func (p *Parser) printStmt() (Stmt, error) {
	exp, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.matchAny(SEMICOLON) == nil {
		return nil, p.genSyntaxError("missing semicolon after value")
	}
	return &PrintStmt{exp}, nil
}

func (p *Parser) exprStmt() (Stmt, error) {
	exp, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.matchAny(SEMICOLON) == nil {
		return nil, p.genSyntaxError("missing semicolon after expression")
	}
	return &ExprStmt{exp}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.matchAny(EQUAL) != nil {
		val, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*Variable); ok {
			return &Assign{
				name:  v.name,
				value: val,
			}, nil
		}
		return nil, p.genSyntaxError("invalid assignment target")
	}
	return expr, nil
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
	if name := p.matchAny(IDENTIFIER); name != nil {
		return &Variable{name}, nil
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
