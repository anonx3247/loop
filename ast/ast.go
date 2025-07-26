package ast

import (
	"com.loop.anonx3247/utils"
)

type Program struct {
	Exprs []Expr
}

type Expr interface {
	Source() utils.String
	Eval() (Value, error)
	CheckDepth(startDepth int) (int, error)
}
type Value interface {
	Type() Type
	Source() utils.String
	IsBase() bool
	String() string
}

type Type interface {
	BaseType() BaseType
	// TupleType, FunctionType, etc.
}

type ParenExpr struct {
	Expr Expr
}

func (p ParenExpr) Source() utils.String {
	return p.Expr.Source()
}

func (p ParenExpr) Eval() (Value, error) {
	return p.Expr.Eval()
}

func (p ParenExpr) CheckDepth(startDepth int) (int, error) {
	if startDepth > 100 {
		return -1, utils.Error{Source: p.Source(), Message: "expression too deep"}
	}
	return p.Expr.CheckDepth(startDepth + 1)
}
