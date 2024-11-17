package parser

import (
	"errors"
	"fmt"

	"com.loop.anonx3247/src/ast"
	"com.loop.anonx3247/src/lexer"
)

// A type can be of multiple kinds
// - primitive
// - Generic
// - User defined
// - Tuple
// - Function type
// - Type with Type parameters

func (p *Parser) parsePrimitiveType() (*ast.Type, error) {
	switch p.Tokens[p.Cursor].Type {
	case lexer.U8_TYPE:
		return &ast.Type{Kind: ast.U8}, nil
	case lexer.U16_TYPE:
		return &ast.Type{Kind: ast.U16}, nil
	case lexer.U32_TYPE:
		return &ast.Type{Kind: ast.U32}, nil
	case lexer.U64_TYPE:
		return &ast.Type{Kind: ast.U64}, nil
	case lexer.I8_TYPE:
		return &ast.Type{Kind: ast.I8}, nil
	case lexer.I16_TYPE:
		return &ast.Type{Kind: ast.I16}, nil
	case lexer.I32_TYPE:
		return &ast.Type{Kind: ast.I32}, nil
	case lexer.I64_TYPE:
		return &ast.Type{Kind: ast.I64}, nil
	case lexer.F32_TYPE:
		return &ast.Type{Kind: ast.F32}, nil
	case lexer.F64_TYPE:
		return &ast.Type{Kind: ast.F64}, nil
	case lexer.BOOL_TYPE:
		return &ast.Type{Kind: ast.BOOL}, nil
	case lexer.CHAR_TYPE:
		return &ast.Type{Kind: ast.CHAR}, nil
	case lexer.STR_TYPE:
		return &ast.Type{Kind: ast.STR}, nil
	}
	return nil, errors.New("expected primitive type")
}

func (p *Parser) parseGenericType() (*ast.Type, error) {
	if p.Tokens[p.Cursor].Type == lexer.TYPE && len(p.Tokens[p.Cursor].Value) == 1 {
		return &ast.Type{Kind: ast.GENERIC_TYPE, Name: p.Tokens[p.Cursor].Value}, nil
	}
	return nil, errors.New("expected generic type")
}

func (p *Parser) parseUserDefinedType() (*ast.Type, error) {
	if p.Tokens[p.Cursor].Type == lexer.TYPE && len(p.Tokens[p.Cursor].Value) > 1 {
		return &ast.Type{Kind: ast.USER_DEFINED_TYPE, Name: p.Tokens[p.Cursor].Value}, nil
	}
	return nil, errors.New("expected user defined type")
}

func (p *Parser) parseTupleType() (*ast.Type, error) {
	// first make sure we are starting a tuple with an open parenthesis
	if p.Cursor >= len(p.Tokens) {
		fmt.Println("Cursor out of bounds")
		return nil, errors.New("cursor out of bounds")
	}
	if p.Tokens[p.Cursor].Type != lexer.OPEN_PARENTHESIS {
		return nil, errors.New("expected open parenthesis")
	}
	tuple := &ast.Type{Kind: ast.TUPLE_TYPE}
	matchIdx, err := findMatchingToken(p.Tokens, p.Cursor)
	if err != nil {
		return nil, err
	}
	defer func() { p.Cursor = matchIdx }()
	p.Cursor++
	splitTokens, err := smartSplitTokens(p.Tokens[p.Cursor:matchIdx], lexer.COMMA)
	if err != nil {
		return nil, err
	}
	if len(splitTokens) == 0 {
		return nil, errors.New("expected at least one type")
	}
	for _, tokenPhrase := range splitTokens {
		parser := Parser{Tokens: tokenPhrase, Cursor: 0}

		parsedType, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		tuple.Params = append(tuple.Params, *parsedType)
	}
	return tuple, nil
}
func (p *Parser) parseType() (*ast.Type, error) {
	if p.Cursor >= len(p.Tokens) {
		return nil, errors.New("cursor out of bounds")
	}

	switch p.Tokens[p.Cursor].Type {
	case lexer.OPEN_PARENTHESIS:
		fmt.Println("Parsing tuple type")
		return p.parseTupleType()
	case lexer.TYPE:
		if len(p.Tokens[p.Cursor].Value) == 1 {
			fmt.Println("Parsing generic type")
			return p.parseGenericType()
		} else {
			fmt.Println("Parsing user defined type")
			return p.parseUserDefinedType()
		}
	case lexer.U8_TYPE, lexer.U16_TYPE, lexer.U32_TYPE, lexer.U64_TYPE, lexer.I8_TYPE, lexer.I16_TYPE, lexer.I32_TYPE, lexer.I64_TYPE, lexer.F32_TYPE, lexer.F64_TYPE, lexer.BOOL_TYPE, lexer.CHAR_TYPE, lexer.STR_TYPE:
		fmt.Println("Parsing primitive type")
		return p.parsePrimitiveType()
	default:
		fmt.Println("Parsing unknown type")
		return nil, errors.New("expected type")
	}
}
