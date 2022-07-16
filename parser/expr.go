package parser

type Expr interface {
	String() string
}

type BinaryExpr struct {
	left  Expr
	op    *Token
	right Expr
}

func (e *BinaryExpr) String() string {
	return "(" + e.op.Lexeme + " " + e.left.String() + " " + e.right.String() + ")"
}

type UnaryExpr struct {
	op    *Token
	right Expr
}

func (e *UnaryExpr) String() string {
	return "(" + e.op.Lexeme + " " + e.right.String() + ")"
}

type LiteralExpr struct {
	op *Token
}

func (e *LiteralExpr) String() string {
	return e.op.Lexeme
}

type GroupingExpr struct {
	expr Expr
}

func (e *GroupingExpr) String() string {
	return "(group " + e.expr.String() + ")"
}
