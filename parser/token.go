//go:generate stringer -type=TokenType

package parser

import "fmt"

type Token struct {
	Typ     TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Typ, t.Lexeme, t.Literal)
}

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	IDENTIFIER
	STRING
	NUMBER

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)
