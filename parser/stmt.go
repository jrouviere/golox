package parser

import "fmt"

type Stmt interface {
	Evaluate() error
	String() string
}

type PrintStmt struct {
	value Expr
}

func (e *PrintStmt) String() string {
	return "(print " + e.value.String() + ")"
}

func (e *PrintStmt) Evaluate() error {
	v, err := e.value.Evaluate()
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

func (e *ExprStmt) Evaluate() error {
	_, err := e.value.Evaluate()
	return err
}
