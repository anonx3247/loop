package ast

import (
	"com.loop.anonx3247/utils"
)

type Expr interface {
	Source() utils.String
	Eval() (Value, error)
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

type Scope struct {
	Exprs []Expr
}

func (s *Scope) Eval() (Value, error) {
	for i, expr := range s.Exprs {
		if i == len(s.Exprs)-1 {
			return expr.Eval()
		}
		_, err := expr.Eval()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *Scope) Source() utils.String {
	sources := make([]utils.String, len(s.Exprs))
	for i, expr := range s.Exprs {
		sources[i] = expr.Source()
	}
	return utils.Encompass(sources...)
}
