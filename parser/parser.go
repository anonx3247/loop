package parser

import (
	"com.loop.anonx3247/ast"
	"com.loop.anonx3247/lexer"
	"com.loop.anonx3247/utils"
)

type Parser struct {
	tokens lexer.TokenList
	pos    int
}

func (p *Parser) error(message string) error {
	if p.pos >= len(p.tokens) {
		return utils.Error{Source: p.tokens[len(p.tokens)-1].Value, Message: message}
	}
	return utils.Error{Source: p.tokens[p.pos].Value, Message: message}
}

func NewParser(source string) (*Parser, error) {
	lexer := lexer.NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}
	return &Parser{tokens: tokens, pos: 0}, nil
}

func ParserFromTokens(tokens lexer.TokenList) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) PeekTokens(tokens int) []lexer.Token {
	return p.tokens[p.pos : p.pos+tokens]
}

func (p *Parser) Peek() (lexer.Token, error) {
	if p.pos >= len(p.tokens) {
		return lexer.Token{}, p.error("unexpected EOF")
	}
	return p.tokens[p.pos], nil
}

func (p *Parser) ConsumeTokens(tokens int) ([]lexer.Token, error) {
	consumedTokens := p.PeekTokens(tokens)
	p.pos += tokens
	return consumedTokens, nil
}

func (p *Parser) TryConsume(token lexer.TokenType) (lexer.Token, error) {
	peekedToken, err := p.Peek()
	if err != nil {
		return lexer.Token{}, err
	}
	if peekedToken.Type == token {
		return p.Consume()
	}
	return lexer.Token{}, p.error("expected another token")
}

func (p *Parser) TryConsumeKind(kind lexer.ShapeType) (lexer.Token, error) {
	peekedToken, err := p.Peek()
	if err != nil {
		return lexer.Token{}, err
	}
	if kind.Matches(peekedToken.Type) {
		return p.Consume()
	}
	return lexer.Token{}, p.error("expected another token")
}

func (p *Parser) Consume() (lexer.Token, error) {
	consumedToken, err := p.Peek()
	if err != nil {
		panic(err)
	}
	p.pos++
	return consumedToken, nil
}

func (p *Parser) Next() (lexer.Token, error) {
	if p.pos >= len(p.tokens) {
		return lexer.Token{}, p.error("unexpected EOF")
	}

	return p.Consume()
}

func (p *Parser) Parse() (ast.Program, error) {
	program := ast.Program{}
	for p.pos < len(p.tokens) {
		expr, err := p.ParseExpr()
		if err != nil {
			return program, err
		}
		program.Exprs = append(program.Exprs, expr)
	}
	return program, nil
}
