package lexer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// we are going to load the basic_tokenizer_test.lp file and tokenize it

func TestSplit(t *testing.T) {
	code := Code("hello\nworld")
	split := code.Split()
	if len(split) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(split))
	}
	if split[0] != "hello" {
		t.Fatalf("expected first line to be 'hello', got %s", split[0])
	}
	if split[1] != "world" {
		t.Fatalf("expected second line to be 'world', got %s", split[1])
	}
}

func TestSortKeys(t *testing.T) {
	keys := []string{"a", "aa", "b", "bb", "c", "cc"}
	sortKeys(keys)
	// make sure each key is bigger in length than the previous key
	for i := 1; i < len(keys); i++ {
		if len(keys[i]) > len(keys[i-1]) {
			t.Fatalf("expected %s to be bigger than %s", keys[i], keys[i-1])
		}
	}
}
func TestGetKeys(t *testing.T) {
	m := map[string]TokenType{"a": PERIOD, "aa": PERIOD, "b": PERIOD, "bb": PERIOD, "c": PERIOD, "cc": PERIOD}
	expected := []string{"a", "aa", "b", "bb", "c", "cc"}
	sort.Strings(expected)
	keys := getKeys(m)
	sort.Strings(keys)
	for i, key := range keys {
		if key != expected[i] {
			t.Fatalf("expected %s to be %s", key, expected[i])
		}
	}
}

func TestTokenizeFromMap(t *testing.T) {
	line := "hello"
	tokens := []Token{}
	idx := 0
	newTokens, newIdx, err := tokenizeFromMap(line, tokens, idx, map[string]TokenType{"hello": PERIOD})
	if err != nil {
		t.Fatalf("failed to tokenize from map: %v", err)
	}
	tokenCheck(t, newTokens, newIdx, line, PERIOD, "hello")
}

func TestTokenizeKeyword(t *testing.T) {
	// we want to check each keyword is matched correctly
	for k, v := range keywords {
		line := k
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeKeyword(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize keyword: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, v, k)
	}
}

func TestTokenizeSymbol(t *testing.T) {
	// we want to check each symbol is matched correctly
	for k, v := range symbols {
		line := k
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeSymbol(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize symbol: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, v, k)
	}
}

func TestTokenizeTypeMap(t *testing.T) {
	// we want to check each type is matched correctly
	for k, v := range types {
		line := k
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeType(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize type: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, v, k)
	}
}

func TestTokenizeType(t *testing.T) {
	examples := []string{"A", "B0", "C3", "MyType"}
	for _, e := range examples {
		line := e
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeType(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize type: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, TYPE, e)
	}
}

func tokenCheck(t *testing.T, newTokens []Token, newIdx int, line string, expectedType TokenType, expectedValue string) {
	if newIdx != len(line) {
		t.Fatalf("expected newIdx to be %d, got %d\nline was: %s", len(line), newIdx, line)
	}
	if len(newTokens) != 1 && strings.TrimSpace(line) != "" {
		t.Fatalf("expected 1 token, got %d\nline was: %s", len(newTokens), line)
	} else if len(newTokens) == 1 {
		if newTokens[0].Type != expectedType {
			t.Fatalf("expected %s token, got %s\nline was: %s", expectedType, newTokens[0].Type, line)
		}
		if newTokens[0].Value != expectedValue {
			t.Fatalf("expected value to be %s, got %s\nline was: %s", expectedValue, newTokens[0].Value, line)
		}
	}

}

func TestTokenizeIdentifier(t *testing.T) {
	examples := []string{"hello", "world", "hello_world", "a123"}
	for _, e := range examples {
		line := e
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeIdentifier(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize identifier: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, IDENTIFIER, e)
	}
}

func TestTokenizeInteger(t *testing.T) {
	examples := []string{"123", "456", "789"}
	for _, e := range examples {
		line := e
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeInteger(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize integer: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, INTEGER, e)
	}
}

func TestTokenizeFloat(t *testing.T) {
	examples := []string{"1.23", "0.32e12", "42.8e-10"}
	for _, e := range examples {
		line := e
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeFloat(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize float: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, FLOAT, e)
	}
}

func TestTokenizeString(t *testing.T) {
	examples := []string{"\"hello\"", "'world'"}
	for _, e := range examples {
		line := e
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeString(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize string: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, STRING, e)
	}
}

func TestTokenizeDecorator(t *testing.T) {
	examples := []string{"@hello", "@world"}
	for _, e := range examples {
		line := e
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeDecorator(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize decorator: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, DECORATOR, e)
	}
}

func TestTokenizeComment(t *testing.T) {
	examples := []string{"-- hello", "--- hello ---"}
	for _, e := range examples {
		line := e
		tokens := []Token{}
		idx := 0
		newTokens, newIdx, err := tokenizeComment(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize comment: %v", err)
		}
		tokenCheck(t, newTokens, newIdx, line, LINE_COMMENT, e)
	}
}

func TestTokenizeNext(t *testing.T) {
	examples := []string{"for x := 0 { x }", "for x := 0 -- hello world"}
	expected := [][]TokenType{
		{FOR, IDENTIFIER, ASSIGN, INTEGER, OPEN_BRACE, IDENTIFIER, CLOSE_BRACE},
		{FOR, IDENTIFIER, ASSIGN, INTEGER, LINE_COMMENT},
	}
	expectedValue := [][]string{
		{"for", "x", ":=", "0", "{", "x", "}"},
		{"for", "x", ":=", "0", "-- hello world"},
	}
	for i, e := range examples {
		line := e
		tokens := []Token{}
		idx := 0
		for j := 0; j < len(expected[i]); j++ {
			newTokens, newIdx, err := tokenizeNext(line, tokens, idx)
			if err != nil {
				t.Fatalf("failed to tokenize next: %v", err)
			}
			tokenCheck(t, newTokens, newIdx-idx, line[idx:newIdx], expected[i][j], expectedValue[i][j])
			idx = newIdx
		}
		newTokens, _, err := tokenizeNext(line, tokens, idx)
		if err != nil {
			t.Fatalf("failed to tokenize next on empty line: %v", err)
		}
		// make sure we have no tokens left
		if len(newTokens) != 0 {
			t.Fatalf("expected no tokens left, got %d", len(newTokens))
		}
	}
}

func TestTokenizeLine(t *testing.T) {
	line := "for x:= if x > 0 { x } else { -x } -- make x positive"

	expected := []TokenType{FOR, IDENTIFIER, ASSIGN, IF, IDENTIFIER, GREATER, INTEGER, OPEN_BRACE, IDENTIFIER, CLOSE_BRACE, ELSE, OPEN_BRACE, MINUS, IDENTIFIER, CLOSE_BRACE, LINE_COMMENT}
	expectedValue := []string{"for", "x", ":=", "if", "x", ">", "0", "{", "x", "}", "else", "{", "-", "x", "}", "-- make x positive"}
	tokens, err := tokenizeLine(line)
	if err != nil {
		t.Fatalf("failed to tokenize line: %v", err)
	}
	fmt.Println("tokens:")
	fmt.Println(tokens)

	for i, token := range tokens {
		if token.Type != expected[i] {
			t.Fatalf("expected %s token, got %s", expected[i], token.Type)
		}
		if token.Value != expectedValue[i] {
			t.Fatalf("expected value to be %s, got %s", expectedValue[i], token.Value)
		}
	}
}

func TestTokenizer(t *testing.T) {
	filePath := filepath.Join("basic_tokenizer_test.lp")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	fmt.Println("TestTokenizer")
	fmt.Println("Testing on the following file:")
	fmt.Println(filePath)
	fmt.Println(string(content))
	code := Code(string(content))
	tokens, err := code.Tokenize()
	if err != nil {
		t.Fatalf("failed to tokenize file: %v", err)
	}
	fmt.Println(tokens)
}
