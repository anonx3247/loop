package ast

import (
	"fmt"

	"com.loop.anonx3247/env"
	"com.loop.anonx3247/lexer"
	"com.loop.anonx3247/utils"
)

type AssignmentKind int

const (
	ASSIGNMENT AssignmentKind = iota
	DECLARATION
	PLUS_ASSIGNMENT
	MINUS_ASSIGNMENT
	MULTIPLY_ASSIGNMENT
	DIVIDE_ASSIGNMENT
	MODULO_ASSIGNMENT
	BITWISE_AND_ASSIGNMENT
	BITWISE_OR_ASSIGNMENT
	BITWISE_XOR_ASSIGNMENT
	BITWISE_LEFT_SHIFT_ASSIGNMENT
	BITWISE_RIGHT_SHIFT_ASSIGNMENT
)

type AssignmentExpr struct {
	Kind   AssignmentKind
	source utils.String
	Const  bool
	Name   string
	Value  Expr
}

func NewAssignmentExpr(identifier lexer.Token, assignmentToken lexer.Token, value Expr) AssignmentExpr {
	fmt.Println("Assigning!", identifier.Value.String(), assignmentToken.Value.String(), value.Source())
	kind := ASSIGNMENT
	switch assignmentToken.Type {
	case lexer.COLON_ASSIGN:
		kind = DECLARATION
	case lexer.ASSIGN:
		kind = ASSIGNMENT
	case lexer.PLUS_ASSIGN:
		kind = PLUS_ASSIGNMENT
	case lexer.MINUS_ASSIGN:
		kind = MINUS_ASSIGNMENT
	case lexer.MULTIPLY_ASSIGN:
		kind = MULTIPLY_ASSIGNMENT
	case lexer.DIVIDE_ASSIGN:
		kind = DIVIDE_ASSIGNMENT
	case lexer.MODULO_ASSIGN:
		kind = MODULO_ASSIGNMENT
	case lexer.BITWISE_AND_ASSIGN:
		kind = BITWISE_AND_ASSIGNMENT
	case lexer.BITWISE_OR_ASSIGN:
		kind = BITWISE_OR_ASSIGNMENT
	case lexer.BITWISE_XOR_ASSIGN:
		kind = BITWISE_XOR_ASSIGNMENT
	case lexer.BITWISE_LEFT_SHIFT_ASSIGN:
		kind = BITWISE_LEFT_SHIFT_ASSIGNMENT
	case lexer.BITWISE_RIGHT_SHIFT_ASSIGN:
		kind = BITWISE_RIGHT_SHIFT_ASSIGNMENT
	}
	return AssignmentExpr{
		Const:  false, // TODO: should be true by default but set to false for testing
		Kind:   kind,
		source: identifier.Value,
		Name:   identifier.Value.String(),
		Value:  value,
	}
}

func (a AssignmentExpr) Eval(env *env.Env) (env.Value, error) {
	value, err := a.Value.Eval(env)
	if err != nil {
		return nil, err
	}
	oldValue, ok := env.Get(a.Name)
	switch a.Kind {
	case ASSIGNMENT:
		if ok {
			if !a.Const {
				env.Set(a.Name, value, a.Const)
			} else {
				return nil, utils.Error{Source: a.Source(), Message: "cannot assign to constant"}
			}
		} else {
			env.Set(a.Name, value, a.Const)
		}
	case DECLARATION:
		env.Set(a.Name, value, a.Const)
	case PLUS_ASSIGNMENT:
		newValue, err := AddValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case MINUS_ASSIGNMENT:
		newValue, err := SubtractValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case MULTIPLY_ASSIGNMENT:
		newValue, err := MultiplyValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case DIVIDE_ASSIGNMENT:
		newValue, err := DivideValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case MODULO_ASSIGNMENT:
		newValue, err := ModuloValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case BITWISE_AND_ASSIGNMENT:
		newValue, err := BitwiseAndValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case BITWISE_OR_ASSIGNMENT:
		newValue, err := BitwiseOrValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case BITWISE_XOR_ASSIGNMENT:
		newValue, err := BitwiseXorValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case BITWISE_LEFT_SHIFT_ASSIGNMENT:
		newValue, err := BitwiseLeftShiftValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	case BITWISE_RIGHT_SHIFT_ASSIGNMENT:
		newValue, err := BitwiseRightShiftValues(value, oldValue, a.Source())
		if err != nil {
			return nil, err
		}
		env.Set(a.Name, newValue, a.Const)
	}
	return value, nil
}

func (a AssignmentExpr) Source() utils.String {
	return a.source
}
