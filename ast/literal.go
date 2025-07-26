package ast

import (
	"strings"

	"com.loop.anonx3247/lexer"
	"com.loop.anonx3247/utils"
)

type Literal struct {
	Value Value
}

func (l Literal) Source() utils.String {
	return l.Value.Source()
}

func (l Literal) Eval() (Value, error) {
	return l.Value, nil
}

func (l Literal) CheckDepth(startDepth int) (int, error) {
	if startDepth > 100 {
		return -1, utils.Error{Source: l.Source(), Message: "expression too deep"}
	}
	return startDepth, nil
}

func LiteralFromToken(tok lexer.Token) (Literal, error) {
	switch tok.Type {
	case lexer.NUMBER_LITERAL:
		if strings.Contains(tok.Value.String(), ".") || strings.Contains(tok.Value.String(), "e") {
			val, err := TryFloatFrom[float32](tok)
			if err != nil {
				return Literal{}, err
			}
			return Literal{Value: val}, nil
		} else {
			val, err := TryIntFrom[int32](tok)
			if err != nil {
				return Literal{}, err
			}
			return Literal{Value: val}, nil
		}
	case lexer.STRING_LITERAL:
		val, err := TryStrFrom(tok)
		if err != nil {
			return Literal{}, err
		}
		return Literal{Value: val}, nil
	case lexer.TRUE, lexer.FALSE:
		val, err := TryBoolFrom(tok)
		if err != nil {
			return Literal{}, err
		}
		return Literal{Value: val}, nil
	}
	return Literal{}, tok.Error("invalid literal")
}

func NewLiteral(value Value) Literal {
	return Literal{Value: value}
}
