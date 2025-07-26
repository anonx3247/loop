package lexer

import (
	"errors"

	"com.loop.anonx3247/src/utils"
)

type Token struct {
	Type  TokenType
	Value utils.String
}

type TokenType int

const (
	// Primitive types
	U8 TokenType = iota
	U16
	U32
	U64
	U128
	I8
	I16
	I32
	I64
	I128
	F32
	F64
	BOOL
	CHAR
	STRING
	GENERIC
	USER_DEFINED

	// Keywords
	IF
	ELIF
	ELSE
	WHILE
	FOR
	LOOP
	RET
	BREAK
	CONTINUE
	MATCH
	COMP
	TYPE
	ABS
	IMPL
	MOD
	USE // local import
	IMPORT
	AS
	FROM
	FN
	LET
	MUT
	IN
	IS
	AND
	OR
	NOT
	TRUE
	FALSE
	NONE
	SELF
	SUPER
	EXCEPT
	NEW
	DEL
	EXIT

	// Punctuation
	L_PAREN                    // (
	R_PAREN                    // )
	L_BRACE                    // {
	R_BRACE                    // }
	L_BRACKET                  // [
	R_BRACKET                  // ]
	COLON                      // :
	COLON_ASSIGN               // :=
	RANGE                      // ..
	PLUS                       // +
	PLUS_ASSIGN                // +=
	MINUS                      // -
	MINUS_ASSIGN               // -=
	MULTIPLY                   // *
	MULTIPLY_ASSIGN            // *=
	DIVIDE                     // /
	DIVIDE_ASSIGN              // /=
	MODULO                     // %
	MODULO_ASSIGN              // %=
	OPTIONAL                   // ?
	OPTIONAL_ASSIGN            // ?=
	ERROR_MARK                 // !
	BITWISE_AND                // &
	BITWISE_AND_ASSIGN         // &=
	BITWISE_OR                 // |
	BITWISE_OR_ASSIGN          // |=
	BITWISE_XOR                // ^
	BITWISE_XOR_ASSIGN         // ^=
	BITWISE_NOT                // ~
	BITWISE_NOT_ASSIGN         // ~=
	BITWISE_LEFT_SHIFT         // <<
	BITWISE_LEFT_SHIFT_ASSIGN  // <<=
	BITWISE_RIGHT_SHIFT        // >>
	BITWISE_RIGHT_SHIFT_ASSIGN // >>=
	ADDRESS_OF                 // #
	PERIOD                     // .
	COMMA                      // ,
	MATCH_ARROW                // =>
	MAP_ARROW                  // ->
	EQUAL                      // ==
	NOT_EQUAL                  // !=
	GREATER_THAN               // >
	GREATER_THAN_OR_EQUAL      // >=
	LESS_THAN                  // <
	LESS_THAN_OR_EQUAL         // <=

	// Literals
	NUMBER_LITERAL // 123, 123.456, 0x123, 0b10101010
	STRING_LITERAL // "Hello, world!", 'a', 'bye bye', r'\my string'
	IDENTIFIER

	NEWLINE
	EOF

	_ANY_TOKEN // used internally, should not match any real token
)

type ShapeType int

const (
	S_OPEN_BRACKET ShapeType = iota
	S_CLOSE_BRACKET
	S_TYPE
	S_VALUE
	S_OPERATOR
	S_ASSIGN_OPERATOR
	S_KEYWORD
	S_TOKEN
)

type ShapeElement struct {
	Type      ShapeType
	TokenType TokenType
}

func (s ShapeElement) Matches(token TokenType) bool {
	var check bool

	switch s.Type {
	case S_TYPE:
		check = token == U8 || token == U16 || token == U32 || token == U64 || token == U128 || token == I8 || token == I16 || token == I32 || token == I64 || token == I128 || token == F32 || token == F64 || token == BOOL || token == CHAR || token == STRING || token == GENERIC || token == USER_DEFINED
	case S_VALUE:
		check = token == NUMBER_LITERAL || token == STRING_LITERAL || token == IDENTIFIER || token == TRUE || token == FALSE || token == NONE || token == SELF || token == SUPER
	case S_OPEN_BRACKET:
		check = token == L_BRACKET || token == L_BRACE || token == L_PAREN
	case S_CLOSE_BRACKET:
		check = token == R_BRACKET || token == R_BRACE || token == R_PAREN
	case S_OPERATOR:
		check = token == PLUS || token == MINUS || token == MULTIPLY || token == DIVIDE || token == MODULO || token == BITWISE_AND || token == BITWISE_OR || token == BITWISE_XOR || token == BITWISE_NOT || token == BITWISE_LEFT_SHIFT || token == BITWISE_RIGHT_SHIFT || token == EQUAL || token == NOT_EQUAL || token == GREATER_THAN || token == GREATER_THAN_OR_EQUAL || token == LESS_THAN || token == LESS_THAN_OR_EQUAL || token == ADDRESS_OF || token == OPTIONAL
	case S_ASSIGN_OPERATOR:
		check = token == COLON_ASSIGN || token == PLUS_ASSIGN || token == MINUS_ASSIGN || token == MULTIPLY_ASSIGN || token == DIVIDE_ASSIGN || token == MODULO_ASSIGN || token == BITWISE_AND_ASSIGN || token == BITWISE_OR_ASSIGN || token == BITWISE_XOR_ASSIGN || token == BITWISE_NOT_ASSIGN || token == BITWISE_LEFT_SHIFT_ASSIGN || token == BITWISE_RIGHT_SHIFT_ASSIGN
	case S_KEYWORD:
		check = token == IF || token == ELIF || token == ELSE || token == WHILE || token == FOR || token == LOOP || token == RET || token == BREAK || token == CONTINUE || token == MATCH || token == COMP || token == TYPE || token == ABS || token == IMPL || token == MOD || token == USE || token == IMPORT || token == AS || token == FN || token == LET || token == MUT || token == IN || token == IS || token == AND || token == OR || token == NOT || token == EXCEPT || token == NEW || token == DEL || token == EXIT
	default:
		check = true
	}

	if s.TokenType != _ANY_TOKEN {
		check = check && s.TokenType == token
	}

	return check
}

type Shape []ShapeElement

func Equal(a, b Token) bool {
	switch a.Type {
	case NUMBER_LITERAL:
		return a.Value.Equal(b.Value) && a.Type == b.Type
	case STRING_LITERAL:
		return a.Value.Equal(b.Value) && a.Type == b.Type
	case IDENTIFIER:
		return a.Value.Equal(b.Value) && a.Type == b.Type
	case GENERIC:
		return a.Value.Equal(b.Value) && a.Type == b.Type
	case USER_DEFINED:
		return a.Value.Equal(b.Value) && a.Type == b.Type
	default:
		return a.Type == b.Type
	}
}

type TokenList []Token

func NewTokenList(tokens ...Token) TokenList {
	return TokenList(tokens)
}

func (tokens TokenList) MatchExact(index int, tokensToMatch TokenList) bool {
	if index >= len(tokens) || index+len(tokensToMatch) > len(tokens) {
		return false
	}

	for i := 0; i < len(tokensToMatch); i++ {
		if !Equal(tokens[index+i], tokensToMatch[i]) {
			return false
		}
	}

	return true
}

func (tokens TokenList) MatchShape(index int, shape Shape) bool {
	if index >= len(tokens) || index+len(shape) > len(tokens) {
		return false
	}

	for i := 0; i < len(shape); i++ {
		if !shape[i].Matches(tokens[index+i].Type) {
			return false
		}
	}

	return true
}

func (tokens TokenList) MatchAnyShape(index int, shapes []Shape) bool {
	for _, shape := range shapes {
		if tokens.MatchShape(index, shape) {
			return true
		}
	}
	return false
}

func (tokens TokenList) FindMatchingBracket(index int) (int, error) {
	if index >= len(tokens) {
		return -1, errors.New("index out of bounds")
	}

	bracket := tokens[index].Type
	var matchingBracket TokenType
	switch bracket {
	case L_BRACE:
		matchingBracket = R_BRACE
	case L_BRACKET:
		matchingBracket = R_BRACKET
	case L_PAREN:
		matchingBracket = R_PAREN
	default:
		return -1, errors.New("invalid bracket")
	}

	depth := 1
	for i := index + 1; i < len(tokens); i++ {
		if tokens[i].Type == bracket {
			depth++
		} else if tokens[i].Type == matchingBracket {
			depth--
			if depth == 0 {
				return i, nil
			}
		}
	}
	return -1, errors.New("unmatched bracket")
}
