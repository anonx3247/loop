package ast

import (
	"com.loop.anonx3247/env"
	envv "com.loop.anonx3247/env"
	"com.loop.anonx3247/utils"
)

type ConditionalExpr struct {
	Condition Expr
	Content   Scope
	Next      *ConditionalExpr
}

func (c ConditionalExpr) Eval(env *envv.Env) (envv.Value, error) {
	condition, err := c.Condition.Eval(env)
	if err != nil {
		return nil, err
	}

	if !condition.IsBase() {
		return nil, utils.Error{Source: condition.Source(), Message: "condition is not a boolean"}
	}

	conditionValue := condition.(envv.BaseValue[bool])

	if conditionValue.GetValue() {
		return c.Content.Eval(env)
	} else if c.Next != nil {
		return c.Next.Eval(env)
	}
	return nil, nil
}

func (c ConditionalExpr) Source() utils.String {
	if c.Next != nil {
		return utils.Encompass(c.Condition.Source(), c.Content.Source(), c.Next.Source())
	}
	return utils.Encompass(c.Condition.Source(), c.Content.Source())
}

func NewElseExpr(content Scope, source utils.String) ConditionalExpr {
	return ConditionalExpr{
		Condition: NewLiteral(env.NewBoolValue(true, source)),
		Content:   content,
		Next:      nil,
	}
}
