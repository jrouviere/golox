package parser

import (
	"fmt"
	"strings"
)

type Expr interface {
	Evaluate(env *Env) (interface{}, error)
	String() string
}

type RuntimeError struct {
	Msg string
}

func (e *RuntimeError) Error() string {
	return "runtime error: " + e.Msg
}

type BinaryExpr struct {
	left  Expr
	op    *Token
	right Expr
}

func (e *BinaryExpr) String() string {
	return "(" + e.op.Lexeme + " " + e.left.String() + " " + e.right.String() + ")"
}

func (e *BinaryExpr) Evaluate(env *Env) (interface{}, error) {
	l, err := e.left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	r, err := e.right.Evaluate(env)
	if err != nil {
		return nil, err
	}

	switch e.op.Typ {
	case PLUS:
		if allNumbers(l, r) {
			return l.(float64) + r.(float64), nil
		}
		if allStrings(l, r) {
			return l.(string) + r.(string), nil
		}
	case MINUS:
		if allNumbers(l, r) {
			return l.(float64) - r.(float64), nil
		}
	case STAR:
		if allNumbers(l, r) {
			return l.(float64) * r.(float64), nil
		}
	case SLASH:
		if allNumbers(l, r) {
			return l.(float64) / r.(float64), nil
		}
	case EQUAL_EQUAL:
		return isEqual(l, r)
	case BANG_EQUAL:
		eq, err := isEqual(l, r)
		return !eq, err
	case LESS_EQUAL:
		if allNumbers(l, r) {
			return l.(float64) <= r.(float64), nil
		}
	case LESS:
		if allNumbers(l, r) {
			return l.(float64) < r.(float64), nil
		}
		if allStrings(l, r) {
			return l.(string) < r.(string), nil
		}
	case GREATER_EQUAL:
		if allNumbers(l, r) {
			return l.(float64) >= r.(float64), nil
		}
	case GREATER:
		if allNumbers(l, r) {
			return l.(float64) > r.(float64), nil
		}
		if allStrings(l, r) {
			return l.(string) > r.(string), nil
		}
	}

	return nil, &RuntimeError{
		Msg: fmt.Sprintf("unimplemented operation %T %v %T", l, e.op.Lexeme, r),
	}
}

type UnaryExpr struct {
	op    *Token
	right Expr
}

func (e *UnaryExpr) String() string {
	return "(" + e.op.Lexeme + " " + e.right.String() + ")"
}

func (e *UnaryExpr) Evaluate(env *Env) (interface{}, error) {
	r, err := e.right.Evaluate(env)
	if err != nil {
		return nil, err
	}

	switch e.op.Typ {
	case MINUS:
		return -r.(float64), nil
	}
	return nil, &RuntimeError{"unimplemented"}
}

type LiteralExpr struct {
	op *Token
}

func (e *LiteralExpr) String() string {
	return e.op.Lexeme
}

func (e *LiteralExpr) Evaluate(env *Env) (interface{}, error) {
	switch e.op.Typ {
	case NIL:
		return nil, nil
	case FALSE:
		return false, nil
	case TRUE:
		return true, nil
	}
	return e.op.Literal, nil
}

type GroupingExpr struct {
	expr Expr
}

func (e *GroupingExpr) String() string {
	return "(group " + e.expr.String() + ")"
}

func (e *GroupingExpr) Evaluate(env *Env) (interface{}, error) {
	return e.expr.Evaluate(env)
}

type Variable struct {
	name *Token
}

func (e *Variable) String() string {
	return "(value " + e.name.Lexeme + ")"
}

func (e *Variable) Evaluate(env *Env) (interface{}, error) {
	return env.Get(e.name.Lexeme)
}

type Assign struct {
	name  *Token
	value Expr
}

func (e *Assign) String() string {
	return "(assign " + e.name.Lexeme + " " + e.value.String() + ")"
}

func (e *Assign) Evaluate(env *Env) (interface{}, error) {
	v, err := e.value.Evaluate(env)
	if err != nil {
		return nil, err
	}
	return v, env.Set(e.name.Lexeme, v)
}

type Logical struct {
	left     Expr
	operator *Token
	right    Expr
}

func (e *Logical) String() string {
	return "(" + e.operator.Lexeme + " " + e.left.String() + ", " + e.right.String() + ")"
}

func (e *Logical) Evaluate(env *Env) (interface{}, error) {
	l, err := e.left.Evaluate(env)
	if err != nil {
		return nil, err
	}

	if e.operator.Typ == OR {
		if isTruthy(l) {
			return l, nil
		}
	} else {
		if !isTruthy(l) {
			return l, nil
		}
	}

	return e.right.Evaluate(env)
}

type Call struct {
	callee Expr
	paren  *Token
	args   []Expr
}

func (e *Call) String() string {
	var args []string
	for _, arg := range e.args {
		args = append(args, arg.String())
	}
	return "(call " + e.callee.String() + "(" + strings.Join(args, ",") + ")"
}

func (e *Call) Evaluate(env *Env) (interface{}, error) {
	callee, err := e.callee.Evaluate(env)
	if err != nil {
		return nil, err
	}

	var args []interface{}
	for _, a := range e.args {
		arg, err := a.Evaluate(env)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	callable, ok := callee.(Callable)
	if !ok {
		return nil, &RuntimeError{
			Msg: "can only call functions and classes",
		}
	}

	if callable.Arity() != len(args) {
		return nil, &RuntimeError{
			Msg: fmt.Sprintf("expected %d arguments but got %d", callable.Arity(), len(args)),
		}
	}

	return callable.Call(env, args)
}

type Callable interface {
	Arity() int
	Call(env *Env, args []interface{}) (interface{}, error)
}

// ---

func allNumbers(vals ...interface{}) bool {
	for _, v := range vals {
		if _, ok := v.(float64); !ok {
			return false
		}
	}
	return true
}

func allStrings(vals ...interface{}) bool {
	for _, v := range vals {
		if _, ok := v.(string); !ok {
			return false
		}
	}
	return true
}
func allBools(vals ...interface{}) bool {
	for _, v := range vals {
		if _, ok := v.(bool); !ok {
			return false
		}
	}
	return true
}
func isEqual(l, r interface{}) (bool, error) {
	if allNumbers(l, r) {
		return l.(float64) == r.(float64), nil
	}
	if allStrings(l, r) {
		return l.(string) == r.(string), nil
	}
	if allBools(l, r) {
		return l.(bool) == r.(bool), nil
	}
	return false, &RuntimeError{
		Msg: fmt.Sprintf("cannot compare %T and %T", l, r),
	}
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v := v.(type) {
	case bool:
		return v
	default:
		return true
	}
}
