package parser

import (
	"fmt"
	"strings"
)

type Stmt interface {
	Evaluate(env *Env) error
	String() string
}

type PrintStmt struct {
	value Expr
}

func (e *PrintStmt) String() string {
	return "(print " + e.value.String() + ")"
}

func (e *PrintStmt) Evaluate(env *Env) error {
	v, err := e.value.Evaluate(env)
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}

type ExprStmt struct {
	value Expr
}

func (e *ExprStmt) String() string {
	return e.value.String()
}

func (e *ExprStmt) Evaluate(env *Env) error {
	_, err := e.value.Evaluate(env)
	return err
}

type VarDecl struct {
	name *Token
	init Expr
}

func (e *VarDecl) String() string {
	return "(var " + e.name.String() + " = " + e.init.String() + " )"
}

func (e *VarDecl) Evaluate(env *Env) error {
	v, err := e.init.Evaluate(env)
	if err != nil {
		return err
	}
	env.Define(e.name.Lexeme, v)
	return nil
}

type Block struct {
	statements []Stmt
}

func (e *Block) String() string {
	var b strings.Builder
	b.WriteString("(block \n")
	for _, s := range e.statements {
		b.WriteString(s.String() + "\n")
	}
	b.WriteString(")")
	return b.String()
}

func (e *Block) Evaluate(env *Env) error {
	scope := NewEnv(env)

	for _, s := range e.statements {
		if err := s.Evaluate(scope); err != nil {
			return err
		}
	}

	return nil
}
