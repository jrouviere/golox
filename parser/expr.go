package parser

type Expr interface {
	Evaluate() (interface{}, error)
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

func (e *BinaryExpr) Evaluate() (interface{}, error) {
	l, err := e.left.Evaluate()
	if err != nil {
		return nil, err
	}
	r, err := e.right.Evaluate()
	if err != nil {
		return nil, err
	}

	// TODO: check types

	switch e.op.Typ {
	case PLUS:
		return l.(float64) + r.(float64), nil
	case MINUS:
		return l.(float64) - r.(float64), nil
	case STAR:
		return l.(float64) * r.(float64), nil
	case SLASH:
		return l.(float64) / r.(float64), nil
	case EQUAL_EQUAL:
		return l.(float64) == r.(float64), nil
	case BANG_EQUAL:
		return l.(float64) != r.(float64), nil
	case LESS_EQUAL:
		return l.(float64) <= r.(float64), nil
	case LESS:
		return l.(float64) < r.(float64), nil
	case GREATER_EQUAL:
		return l.(float64) >= r.(float64), nil
	case GREATER:
		return l.(float64) > r.(float64), nil
	}

	return nil, &RuntimeError{"unimplemented"}
}

type UnaryExpr struct {
	op    *Token
	right Expr
}

func (e *UnaryExpr) String() string {
	return "(" + e.op.Lexeme + " " + e.right.String() + ")"
}

func (e *UnaryExpr) Evaluate() (interface{}, error) {
	r, err := e.right.Evaluate()
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

func (e *LiteralExpr) Evaluate() (interface{}, error) {
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

func (e *GroupingExpr) Evaluate() (interface{}, error) {
	return e.expr.Evaluate()
}
