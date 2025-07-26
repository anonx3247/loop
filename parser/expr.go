package parser

import (
	"com.loop.anonx3247/ast"
	"com.loop.anonx3247/lexer"
)

var operatorPrecedence = map[lexer.TokenType]int{
	lexer.PLUS:                  2,
	lexer.MINUS:                 2,
	lexer.MULTIPLY:              3,
	lexer.DIVIDE:                3,
	lexer.MODULO:                3,
	lexer.BITWISE_AND:           4,
	lexer.BITWISE_OR:            4,
	lexer.BITWISE_XOR:           4,
	lexer.BITWISE_NOT:           4,
	lexer.BITWISE_LEFT_SHIFT:    5,
	lexer.BITWISE_RIGHT_SHIFT:   5,
	lexer.EQUAL:                 1,
	lexer.NOT_EQUAL:             1,
	lexer.GREATER_THAN:          1,
	lexer.GREATER_THAN_OR_EQUAL: 1,
	lexer.LESS_THAN:             1,
	lexer.LESS_THAN_OR_EQUAL:    1,
	lexer.ADDRESS_OF:            2,
}

func (p *Parser) ParseExpr() (ast.Expr, error) {
	return p.parseExprWithPrecedence(0)
}

// note that unary expressions are considered atoms, as well as parenthesis
func (p *Parser) parseAtom() (ast.Expr, error) {
	leftToken, err := p.Consume()
	if err != nil {
		return nil, err
	}

	if lexer.S_UNARY_OPERATOR.Matches(leftToken.Type) {
		expr, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}
		return ast.UnaryExpr{Op: leftToken.Type, Value: expr}, nil
	} else if leftToken.Type == lexer.L_PAREN {
		expr, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}
		_, err = p.TryConsume(lexer.R_PAREN)
		if err != nil {
			return nil, err
		}
		return ast.ParenExpr{Expr: expr}, nil
	} else if lexer.S_VALUE.Matches(leftToken.Type) {
		lit, err := ast.NewLiteral(leftToken)
		if err != nil {
			return nil, err
		}
		return lit, nil
	}
	return nil, p.error("expected atom")
}

func (p *Parser) parseExprWithPrecedence(minPrecedence int) (ast.Expr, error) {
	left, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	var currentToken lexer.Token
	var currentPrecedence int
	for p.pos < len(p.tokens) {

		currentToken, err = p.TryConsumeKind(lexer.S_BINARY_OPERATOR)
		if err != nil {
			return left, nil
		}
		currentPrecedence = operatorPrecedence[currentToken.Type]
		if currentPrecedence < minPrecedence {
			return left, nil
		}
		nextPrecedence := currentPrecedence + 1
		right, err := p.parseExprWithPrecedence(nextPrecedence)
		if err != nil {
			return nil, err
		}
		leftCopy := new(ast.Expr)
		*leftCopy = left
		bin := ast.BinaryExpr{
			Op:    currentToken.Type,
			Left:  leftCopy,
			Right: &right,
		}

		left = &bin
	}

	return left, nil
}
