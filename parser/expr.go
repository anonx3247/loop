package parser

import (
	"com.loop.anonx3247/ast"
	"com.loop.anonx3247/lexer"
)

var operatorPrecedence = map[lexer.TokenType]int{
	lexer.EQUAL:                 1,
	lexer.NOT_EQUAL:             1,
	lexer.GREATER_THAN:          1,
	lexer.GREATER_THAN_OR_EQUAL: 1,
	lexer.LESS_THAN:             1,
	lexer.LESS_THAN_OR_EQUAL:    1,
	lexer.AND:                   1,
	lexer.OR:                    1,
	lexer.NOT:                   1,
	lexer.IF:                    2,
	// lexer.MATCH:                 2,
	lexer.PLUS:                3,
	lexer.MINUS:               3,
	lexer.ADDRESS_OF:          3,
	lexer.MULTIPLY:            4,
	lexer.DIVIDE:              4,
	lexer.MODULO:              4,
	lexer.BITWISE_AND:         5,
	lexer.BITWISE_OR:          5,
	lexer.BITWISE_XOR:         5,
	lexer.BITWISE_NOT:         5,
	lexer.BITWISE_LEFT_SHIFT:  6,
	lexer.BITWISE_RIGHT_SHIFT: 6,
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
		lit, err := ast.LiteralFromToken(leftToken)
		if err != nil {
			return nil, err
		}
		return lit, nil
	} else if leftToken.Type == lexer.IF {
		return p.parseIfExpr()
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

// assumes that the if token has already been consumed
func (p *Parser) parseIfExpr() (ast.Expr, error) {
	condition, err := p.ParseExpr()

	if err != nil {
		return ast.ConditionalExpr{}, err
	}
	_, err = p.TryConsume(lexer.L_BRACE)
	if err != nil {
		return ast.ConditionalExpr{}, err
	}
	thenExpr, err := p.Parse()
	if err != nil {
		return ast.ConditionalExpr{}, err
	}
	_, err = p.TryConsume(lexer.R_BRACE)

	next, err := p.Peek()
	if err != nil {
		return ast.ConditionalExpr{
			Condition: condition,
			Content:   thenExpr,
			Next:      nil,
		}, nil
	}

	if next.Type == lexer.ELIF {
		p.Consume()
		next, err := p.parseIfExpr()
		nextCond := next.(ast.ConditionalExpr)
		if err != nil {
			return ast.ConditionalExpr{}, err
		}
		return ast.ConditionalExpr{
			Condition: condition,
			Content:   thenExpr,
			Next:      &nextCond,
		}, nil
	}

	if next.Type == lexer.ELSE {
		p.ConsumeTokens(2) // consume the else and the brace
		nextScope, err := p.Parse()
		if err != nil {
			return ast.ConditionalExpr{}, err
		}
		_, err = p.TryConsume(lexer.R_BRACE)
		if err != nil {
			return ast.ConditionalExpr{}, err
		}
		elseCond := ast.NewElseExpr(nextScope, next.Value)
		return ast.ConditionalExpr{
			Condition: condition,
			Content:   thenExpr,
			Next:      &elseCond,
		}, nil
	}

	return ast.ConditionalExpr{
		Condition: condition,
		Content:   thenExpr,
		Next:      nil,
	}, nil
}
