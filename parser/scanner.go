package parser

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	input   string
	line    int
	start   int
	current int
}

type ScanningError struct {
	Line int
	Msg  string
}

func (e ScanningError) Error() string {
	return fmt.Sprintf("Line %d, %v", e.Line, e.Msg)
}

func NewScanner(input string) *Scanner {
	return &Scanner{
		input:   input,
		line:    1,
		start:   0,
		current: 0,
	}
}

func (s *Scanner) Scan() ([]*Token, error) {
	var tokens []*Token

	for !s.eof() {
		s.start = s.current
		tok, err := s.scanToken()
		if err != nil {
			return nil, err
		}
		if tok != nil {
			tokens = append(tokens, tok)
		}
	}

	tokens = append(tokens, &Token{EOF, "", nil, s.line})
	return tokens, nil
}

func (s *Scanner) scanToken() (*Token, error) {
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
			return nil, nil
		}
		return s.genToken(SLASH, nil)

	case ' ', '\t', '\r':
		//skip whitespaces
		return nil, nil
	case '\n':
		s.line++
		return nil, nil

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
	case '<':
		if s.match('=') {
			return s.genToken(LESS_EQUAL, nil)
		}
		return s.genToken(LESS, nil)
	case '>':
		if s.match('=') {
			return s.genToken(GREATER_EQUAL, nil)
		}
		return s.genToken(GREATER, nil)

	case '"':
		return s.genString()

	default:
		if isDigit(c) {
			return s.genNumber()
		}
		if isAlpha(c) {
			return s.genIdent()
		}
		return s.genError("unexpected token: '%s'", string(c))
	}
}

func (s *Scanner) genString() (*Token, error) {
	for s.peek() != '"' && !s.eof() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.eof() {
		return s.genError("unexpected end of string")
	}
	s.advance() // closing "

	return s.genToken(STRING, s.input[s.start+1:s.current-1])
}

func (s *Scanner) genNumber() (*Token, error) {
	for isDigit(s.peek()) && !s.eof() {
		s.advance()
	}
	// parse decimal part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) && !s.eof() {
			s.advance()
		}
	}
	txt := s.input[s.start:s.current]

	val, err := strconv.ParseFloat(txt, 64)
	if err != nil {
		return s.genError("invalid number: %v, %v", txt, err)
	}
	return s.genToken(NUMBER, val)
}

func (s *Scanner) genIdent() (*Token, error) {
	for isAlphaNum(s.peek()) && !s.eof() {
		s.advance()
	}

	txt := s.input[s.start:s.current]
	if tt, f := keywords[txt]; f {
		return s.genToken(tt, nil)
	}
	return s.genToken(IDENTIFIER, txt)
}

func (s *Scanner) genToken(typ TokenType, val interface{}) (*Token, error) {
	lex := s.input[s.start:s.current]
	return &Token{typ, lex, val, s.line}, nil
}

func (s *Scanner) genError(format string, v ...interface{}) (*Token, error) {
	return nil, &ScanningError{
		Line: s.line,
		Msg:  fmt.Sprintf(format, v...),
	}
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
func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.input) {
		return '\x00'
	}
	return rune(s.input[s.current+1])
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

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}
func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
func isAlphaNum(c rune) bool {
	return isDigit(c) || isAlpha(c)
}
