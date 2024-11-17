package parser

import (
	"testing"

	"com.loop.anonx3247/src/ast"
	"com.loop.anonx3247/src/lexer"
)

func TestParsePrimitiveType(t *testing.T) {
	exampleTokens := lexer.TokenPhrase{
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.NEWLINE, Value: "\n"},
		lexer.Token{Type: lexer.F32_TYPE, Value: "f32"},
		lexer.Token{Type: lexer.STR_TYPE, Value: "str"},
	}
	results := []*ast.Type{
		{Kind: ast.U8},
		nil,
		{Kind: ast.F32},
		{Kind: ast.STR},
	}
	p := Parser{Tokens: exampleTokens}
	for i := range exampleTokens {
		p.Cursor = i
		result, err := p.parsePrimitiveType()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil && results[i] != nil {
			t.Errorf("Expected a nil result, got %v", result)
		} else if result != nil && result.Kind != results[i].Kind {
			t.Errorf("Expected %v, got %v", results[i].Kind, result.Kind)
		}
	}
}

func TestParseUserDefinedType(t *testing.T) {
	exampleTokens := lexer.TokenPhrase{
		lexer.Token{Type: lexer.TYPE, Value: "MyType"},
		lexer.Token{Type: lexer.NEWLINE, Value: "\n"},
	}
	p := Parser{Tokens: exampleTokens}
	p.Cursor = 0
	for i := range exampleTokens {
		p.Cursor = i
		result, err := p.parseUserDefinedType()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil && i != 1 {
			t.Errorf("Expected a nil result, got %v", result)
		} else if result != nil && result.Kind != ast.USER_DEFINED_TYPE {
			t.Errorf("Expected a user defined type, got %v", result.Kind)
		}
		if result != nil && result.Name != "MyType" {
			t.Errorf("Expected name to be MyType, got %v", result.Name)
		}
	}
}

func TestSplitTokens(t *testing.T) {
	exampleTokens := lexer.TokenPhrase{
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.F32_TYPE, Value: "f32"},
		lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
	}
	result := splitTokens(exampleTokens, lexer.COMMA)
	if len(result) != 2 {
		t.Errorf("Expected 2 splits, got %v", len(result))
	}
	if result[0][0].Type != lexer.OPEN_PARENTHESIS {
		t.Errorf("Expected first split to be open parenthesis, got %v", result[0][0].Type)
	}
	if result[0][1].Type != lexer.U8_TYPE {
		t.Errorf("Expected first split to be u8, got %v", result[0][1].Type)
	}
	if result[1][0].Type != lexer.F32_TYPE {
		t.Errorf("Expected second split to be f32, got %v", result[1][0].Type)
	}
	if result[1][1].Type != lexer.CLOSE_PARENTHESIS {
		t.Errorf("Expected second split to be close parenthesis, got %v", result[1][1].Type)
	}
}

func TestSmartSplitTokens(t *testing.T) {
	// u8,u8,u8
	exampleTokens := lexer.TokenPhrase{
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
	}
	result, err := smartSplitTokens(exampleTokens, lexer.COMMA)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 splits, got %v", len(result))
	}
	// u8,(u8,u8)
	exampleTokens = lexer.TokenPhrase{
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
	}
	result, err = smartSplitTokens(exampleTokens, lexer.COMMA)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 splits, got %v", len(result))
	}
	// (u8,u8),u8
	exampleTokens = lexer.TokenPhrase{
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
	}
	result, err = smartSplitTokens(exampleTokens, lexer.COMMA)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 splits, got %v", len(result))
	}
}

func TestFindMatchingToken(t *testing.T) {
	example1Tokens := lexer.TokenPhrase{
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
	}
	result, err := findMatchingToken(example1Tokens, 0)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}

	example2Tokens := lexer.TokenPhrase{
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.F32_TYPE, Value: "f32"},
		lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
	}
	result, err = findMatchingToken(example2Tokens, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 4 {
		t.Errorf("Expected result to be 4, got %v", result)
	}
	// multiple open parentheses
	example3Tokens := lexer.TokenPhrase{
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
		lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
	}
	result, err = findMatchingToken(example3Tokens, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 4 {
		t.Errorf("Expected result to be 4, got %v", result)
	}

	// test the case (int, (str, f32))
	example4Tokens := lexer.TokenPhrase{
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.STR_TYPE, Value: "str"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.F32_TYPE, Value: "f32"},
		lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
		lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
	}
	result, err = findMatchingToken(example4Tokens, 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 8 {
		t.Errorf("Expected result to be 8, got %v", result)
	}
}

func TestParseTupleType(t *testing.T) {
	examplePhrases := []lexer.TokenPhrase{
		{
			lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
			lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
			lexer.Token{Type: lexer.COMMA, Value: ","},
			lexer.Token{Type: lexer.F32_TYPE, Value: "f32"},
			lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
		},
		{
			lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
			lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
		},
		{
			lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
			lexer.Token{Type: lexer.STR_TYPE, Value: "str"},
			lexer.Token{Type: lexer.COMMA, Value: ","},
			lexer.Token{Type: lexer.F32_TYPE, Value: "f32"},
			lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
		},
		{
			lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
			lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
			lexer.Token{Type: lexer.COMMA, Value: ","},
			lexer.Token{Type: lexer.TYPE, Value: "MyType"},
			lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
		},
		// (u8,MyType,(f32,str))
		{
			lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
			lexer.Token{Type: lexer.U8_TYPE, Value: "u8"},
			lexer.Token{Type: lexer.COMMA, Value: ","},
			lexer.Token{Type: lexer.TYPE, Value: "MyType"},
			lexer.Token{Type: lexer.COMMA, Value: ","},
			lexer.Token{Type: lexer.OPEN_PARENTHESIS, Value: "("},
			lexer.Token{Type: lexer.F32_TYPE, Value: "f32"},
			lexer.Token{Type: lexer.COMMA, Value: ","},
			lexer.Token{Type: lexer.STR_TYPE, Value: "str"},
			lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
			lexer.Token{Type: lexer.CLOSE_PARENTHESIS, Value: ")"},
		},
	}
	results := []*ast.Type{
		{Kind: ast.TUPLE_TYPE, Params: []ast.Type{{Kind: ast.U8}, {Kind: ast.F32}}},
		nil,
		{Kind: ast.TUPLE_TYPE, Params: []ast.Type{{Kind: ast.STR}, {Kind: ast.F32}}},
		{Kind: ast.TUPLE_TYPE, Params: []ast.Type{{Kind: ast.U8}, {Kind: ast.USER_DEFINED_TYPE, Name: "MyType"}}},
		{Kind: ast.TUPLE_TYPE, Params: []ast.Type{{Kind: ast.U8}, {Kind: ast.USER_DEFINED_TYPE, Name: "MyType"}, {Kind: ast.TUPLE_TYPE, Params: []ast.Type{{Kind: ast.F32}, {Kind: ast.STR}}}}},
	}
	p := Parser{}
	for i := range examplePhrases {
		p.Tokens = examplePhrases[i]
		p.Cursor = 0
		result, err := p.parseTupleType()
		if err != nil && i != 1 {
			t.Errorf("Expected no error, got: %v for phrase %v", err, i)
		} else if result == nil && results[i] != nil {
			t.Errorf("Expected a non-nil result for phrase %v", i)
		} else if result != nil && result.Kind != results[i].Kind {
			t.Errorf("Expected %v, got %v", results[i].Kind, result.Kind)
		}
	}
}
