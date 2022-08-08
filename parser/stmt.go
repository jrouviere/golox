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

// we define it as an error so it can bubble up like an exception
// could use panic instead but that seemed overkill
type ReturnValue struct {
	val interface{}
}

func (r *ReturnValue) Error() string {
	return ""
}

type ReturnStmt struct {
	value Expr
}

func (e *ReturnStmt) String() string {
	return "(return " + e.value.String() + ")"
}

func (e *ReturnStmt) Evaluate(env *Env) error {
	if e.value == nil {
		return &ReturnValue{nil}
	}

	v, err := e.value.Evaluate(env)
	if err != nil {
		return err
	}
	return &ReturnValue{v}
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

type FunStmt struct {
	name   *Token
	params []*Token
	body   Stmt
}

func (e *FunStmt) String() string {
	return e.name.String()
}

func (e *FunStmt) Evaluate(env *Env) error {
	env.Define(e.name.Lexeme, &LoxFunction{
		Declaration: e,
	})
	return nil
}

type VarDecl struct {
	name *Token
	init Expr
}

func (e *VarDecl) String() string {
	if e.init == nil {
		return "(var " + e.name.String() + " )"
	}
	return "(var " + e.name.String() + " = " + e.init.String() + " )"
}

func (e *VarDecl) Evaluate(env *Env) error {
	var init interface{}
	if e.init != nil {
		v, err := e.init.Evaluate(env)
		if err != nil {
			return err
		}
		init = v
	}
	env.Define(e.name.Lexeme, init)
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

type IfStmt struct {
	expr     Expr
	thenBrch Stmt
	elseBrch Stmt
}

func (e *IfStmt) String() string {
	var b strings.Builder
	b.WriteString("(if " + e.expr.String() + "\n")
	b.WriteString(e.thenBrch.String() + "\n")
	if e.elseBrch != nil {
		b.WriteString(") else (\n")
		b.WriteString(e.elseBrch.String() + "\n")
	}
	b.WriteString(")")
	return b.String()
}

func (e *IfStmt) Evaluate(env *Env) error {

	val, err := e.expr.Evaluate(env)
	if err != nil {
		return err
	}

	if isTruthy(val) {
		return e.thenBrch.Evaluate(env)
	} else {
		if e.elseBrch != nil {
			return e.elseBrch.Evaluate(env)
		}
	}
	return nil
}

type WhileStmt struct {
	expr Expr
	body Stmt
}

func (e *WhileStmt) String() string {
	var b strings.Builder
	b.WriteString("(while " + e.expr.String() + "\n")
	b.WriteString(e.body.String() + "\n")
	b.WriteString(")")
	return b.String()
}

func (e *WhileStmt) Evaluate(env *Env) error {
	for {
		cond, err := e.expr.Evaluate(env)
		if err != nil {
			return err
		}
		if !isTruthy(cond) {
			return nil
		}

		if err := e.body.Evaluate(env); err != nil {
			return err
		}
	}
}
