package ast

import (
	"com.loop.anonx3247/env"
	"com.loop.anonx3247/utils"
)

type Expr interface {
	Source() utils.String
	Eval(env *env.Env) (env.Value, error)
}

type ParenExpr struct {
	Expr Expr
}

func (p ParenExpr) Source() utils.String {
	return p.Expr.Source()
}

func (p ParenExpr) Eval(env *env.Env) (env.Value, error) {
	return p.Expr.Eval(env)
}

type Scope struct {
	Exprs []Expr
}

func (s *Scope) Eval(env *env.Env) (env.Value, error) {
	for i, expr := range s.Exprs {
		if i == len(s.Exprs)-1 {
			return expr.Eval(env)
		}
		_, err := expr.Eval(env)
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
