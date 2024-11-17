package parser

import (
	"errors"

	"com.loop.anonx3247/src/ast"
	"com.loop.anonx3247/src/lexer"
)

type Parser struct {
	Tokens []lexer.Token
	Cursor int
	Ast    *ast.Program
}

func (p *Parser) Parse() *ast.Program {
	return p.Ast
}

// Finds matches for tokens like ( ) [ ] { } < >
func findMatchingToken(tokens []lexer.Token, cursor int) (int, error) {
	if matchingToken, ok := matchingTokens[tokens[cursor].Type]; ok {
		depth := 1
		for i := cursor + 1; i < len(tokens); i++ {
			if tokens[i].Type == matchingToken {
				depth--
				if depth == 0 {
					return i, nil
				}
			} else if tokens[i].Type == tokens[cursor].Type {
				depth++
			}
		}
	}
	return 0, errors.New("no matching token found for " + tokens[cursor].Value)
}

func splitTokens(tokens lexer.TokenPhrase, splitType lexer.TokenType) (splits []lexer.TokenPhrase) {
	splits = []lexer.TokenPhrase{}
	currentSplit := lexer.TokenPhrase{}
	for _, token := range tokens {
		if token.Type == splitType {
			splits = append(splits, currentSplit)
			currentSplit = lexer.TokenPhrase{}
		} else {
			currentSplit = append(currentSplit, token)
		}
	}
	splits = append(splits, currentSplit)
	return splits
}

// takes into account matching tokens and skips them correctly
func smartSplitTokens(tokens lexer.TokenPhrase, splitType lexer.TokenType) (splits []lexer.TokenPhrase, err error) {
	splits = []lexer.TokenPhrase{}
	currentSplit := lexer.TokenPhrase{}
	index := 0
	for index < len(tokens) {
		if tokens[index].Type == splitType {
			splits = append(splits, currentSplit)
			currentSplit = lexer.TokenPhrase{}
		} else if isOpeningToken(tokens[index]) {
			matchIdx, err := findMatchingToken(tokens, index)
			if err != nil {
				return nil, err
			}
			currentSplit = append(currentSplit, tokens[index:matchIdx+1]...)
			index = matchIdx
		} else if isClosingToken(tokens[index]) {
			return nil, errors.New("unmatched closing token")
		} else {
			currentSplit = append(currentSplit, tokens[index])
		}
		index++
	}
	splits = append(splits, currentSplit)
	return splits, nil
}

var matchingTokens = map[lexer.TokenType]lexer.TokenType{
	lexer.OPEN_PARENTHESIS: lexer.CLOSE_PARENTHESIS,
	lexer.OPEN_BRACKET:     lexer.CLOSE_BRACKET,
	lexer.OPEN_BRACE:       lexer.CLOSE_BRACE,
	lexer.LESS:             lexer.GREATER,
}

func isOpeningToken(token lexer.Token) bool {
	_, ok := matchingTokens[token.Type]
	return ok
}

func isClosingToken(token lexer.Token) bool {
	for _, v := range matchingTokens {
		if v == token.Type {
			return true
		}
	}
	return false
}
