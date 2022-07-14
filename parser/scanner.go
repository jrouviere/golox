package parser

type Scanner struct {
	input   string
	line    int
	start   int
	current int
}

func NewScanner(input string) *Scanner {
	return &Scanner{
		input:   input,
		line:    1,
		start:   0,
		current: 0,
	}
}

func (s *Scanner) Scan() []*Token {
	var tokens []*Token

	for !s.eof() {
		s.start = s.current
		if tok := s.scanToken(); tok != nil {
			tokens = append(tokens, tok)
		}
	}

	return tokens
}

func (s *Scanner) scanToken() *Token {
	c := s.advance()
	switch c {
	case '(':
		return s.genToken(LEFT_PAREN, nil)
	case ')':
		return s.genToken(RIGHT_PAREN, nil)
	case '{':
		return s.genToken(LEFT_BRACE, nil)
	case '}':
		return s.genToken(RIGHT_BRACE, nil)
	case ',':
		return s.genToken(COMMA, nil)
	case '.':
		return s.genToken(DOT, nil)
	case '-':
		return s.genToken(MINUS, nil)
	case '+':
		return s.genToken(PLUS, nil)
	case ';':
		return s.genToken(SEMICOLON, nil)
	case '*':
		return s.genToken(STAR, nil)
	case '/':
		if s.match('/') {
			//comment, consume until the end of line
			for s.peek() != '\n' && !s.eof() {
				s.advance()
			}
			return nil
		}
		return s.genToken(SLASH, nil)

	case ' ', '\t', '\r':
		//skip whitespaces
		return nil
	case '\n':
		s.line++
		return nil

	case '!':
		if s.match('=') {
			return s.genToken(BANG_EQUAL, nil)
		}
		return s.genToken(BANG, nil)

	case '=':
		if s.match('=') {
			return s.genToken(EQUAL_EQUAL, nil)
		}
		return s.genToken(EQUAL, nil)
	}

	return nil
}

func (s *Scanner) genToken(typ TokenType, val interface{}) *Token {
	lex := s.input[s.start:s.current]
	return &Token{typ, lex, val, s.line}
}

func (s *Scanner) eof() bool {
	return s.current >= len(s.input)
}

func (s *Scanner) peek() rune {
	if s.eof() {
		return '\x00'
	}
	return rune(s.input[s.current])
}

func (s *Scanner) advance() rune {
	r := rune(s.input[s.current])
	s.current++
	return r
}

func (s *Scanner) match(expected rune) bool {
	if s.eof() {
		return false
	}
	if rune(s.input[s.current]) != expected {
		return false
	}
	s.current++
	return true
}
