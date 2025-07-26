package ast

import (
	"com.loop.anonx3247/utils"
)

type ConditionalExpr struct {
	Condition Expr
	Content   Scope
	Next      *ConditionalExpr
}

func (c ConditionalExpr) Eval() (Value, error) {
	condition, err := c.Condition.Eval()
	if err != nil {
		return nil, err
	}

	if !condition.IsBase() {
		return nil, utils.Error{Source: condition.Source(), Message: "condition is not a boolean"}
	}

	conditionValue := condition.(BaseValue[bool])

	if conditionValue.GetValue() {
		return c.Content.Eval()
	} else if c.Next != nil {
		return c.Next.Eval()
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
		Condition: NewLiteral(NewBoolValue(true, source)),
		Content:   content,
		Next:      nil,
	}
}
