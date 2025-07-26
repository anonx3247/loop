package lexer

import (
	"regexp"
	"sort"

	"com.loop.anonx3247/utils"
)

var (
	NUMBER_RE        = regexp.MustCompile(`^\d+(\.\d+)?([eE][-]?\d+)?`)
	IDENTIFIER_RE    = regexp.MustCompile(`^[a-z_][a-zA-Z0-9_]*`)
	GENERIC_RE       = regexp.MustCompile(`^[A-Z]`)
	USER_DEFINED_RE  = regexp.MustCompile(`^[A-Z][a-zA-Z0-9_]+`)
	STRING_RE_SINGLE = regexp.MustCompile(`^'([^']|\\')+'`)
	STRING_RE_DOUBLE = regexp.MustCompile(`^"([^"]|\\")+"`)
	STRING_RE_RAW    = regexp.MustCompile("^`([^`]|\\`)+`")

	SINGLE_LINE_COMMENT_RE = regexp.MustCompile(`^--.*`)
	MULTI_LINE_COMMENT_RE  = regexp.MustCompile(`^---`)

	KEYWORDS   = regexp.MustCompile(`^(if|elif|else|while|for|loop|ret|break|continue|match|comp|type|abs|impl|mod|use|import|as|from|fn|let|mut|in|is|and|or|not|true|false|none|self|super|except|new|del|exit)`)
	BASE_TYPES = regexp.MustCompile(`^(u8|u16|u32|u64|u128|i8|i16|i32|i64|i128|f32|f64|bool|char|string)`)
	OPERATORS  = regexp.MustCompile(`^(\(|\)|\{|\}|\[|\]|\:=|\:|\.\.|\+|\+=|-|-=|\*|\*=|/|/=|%|%=|~|~=|&|&=|\||\|=|\^|\^=|#|\.|\,|->|=>|==|!=|>|>=|<|<=|=)`)
)

type Lexer struct {
	source string
	pos    int
}

func NewLexer(source string) *Lexer {
	return &Lexer{source: source}
}

func (l *Lexer) slice(length int) utils.String {
	return utils.StringFrom(l.source, l.pos, length)
}

func (l *Lexer) Tokenize() (TokenList, error) {
	tokens := TokenList{}
	for l.pos < len(l.source) {
		token, err := l.Next()
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (l *Lexer) Match(re *regexp.Regexp) (bool, int) {
	match := re.FindStringIndex(l.source[l.pos:])
	if match == nil {
		return false, 0
	}
	return true, match[1]
}

func (l *Lexer) Next() (Token, error) {
	offset := 0
	defer func() {
		l.pos += offset
	}()
	if l.pos >= len(l.source) {
		return Token{Type: EOF, Value: utils.String{}}, nil
	} else {

		match, length := l.Match(SINGLE_LINE_COMMENT_RE)
		if match {
			l.pos += length
			return l.Next()
		}

		match, length = l.Match(MULTI_LINE_COMMENT_RE)
		if match {
			l.pos += length
			end_found := false
			for l.pos < len(l.source) && !end_found {
				match, length = l.Match(MULTI_LINE_COMMENT_RE)
				if match {
					end_found = true
					l.pos += length
					break
				}
				l.pos++
			}
			if !end_found {
				return Token{}, l.error("unterminated multi-line comment")
			}
			return l.Next()
		}

		atom, err := l.tryTokenizeAtom()

		if err == nil {
			offset = len(atom.Value.String())
			return atom, nil
		}

		for re, t := range map[*regexp.Regexp]TokenType{
			IDENTIFIER_RE:    IDENTIFIER,
			USER_DEFINED_RE:  USER_DEFINED,
			GENERIC_RE:       GENERIC,
			NUMBER_RE:        NUMBER_LITERAL,
			STRING_RE_SINGLE: STRING_LITERAL,
			STRING_RE_DOUBLE: STRING_LITERAL,
			STRING_RE_RAW:    STRING_LITERAL,
		} {
			match, offset = l.Match(re)
			if match {
				return Token{Type: t, Value: l.slice(offset)}, nil
			}
		}

		if l.source[l.pos] == ' ' || l.source[l.pos] == '\t' {
			l.pos++
			return l.Next()
		}

		if l.source[l.pos] == '\n' {
			l.pos++
			return Token{Type: NEWLINE, Value: l.slice(1)}, nil
		}

		return Token{}, l.error("unexpected token")
	}
}

func (l *Lexer) error(message string) utils.Error {
	return utils.Error{Source: l.slice(0), Message: message}
}

func (l *Lexer) tryTokenizeAtom() (Token, error) {

	type Pair struct {
		Word string
		Type TokenType
	}
	var words = []Pair{
		{"if", IF},
		{"elif", ELIF},
		{"else", ELSE},
		{"while", WHILE},
		{"for", FOR},
		{"loop", LOOP},
		{"ret", RET},
		{"break", BREAK},
		{"continue", CONTINUE},
		{"match", MATCH},
		{"comp", COMP},
		{"type", TYPE},
		{"abs", ABS},
		{"impl", IMPL},
		{"mod", MOD},
		{"use", USE},
		{"import", IMPORT},
		{"as", AS},
		{"from", FROM},
		{"fn", FN},
		{"let", LET},
		{"mut", MUT},
		{"in", IN},
		{"is", IS},
		{"and", AND},
		{"or", OR},
		{"not", NOT},
		{"true", TRUE},
		{"false", FALSE},
		{"none", NONE},
		{"self", SELF},
		{"super", SUPER},
		{"except", EXCEPT},
		{"new", NEW},
		{"del", DEL},
		{"exit", EXIT},
		{"u8", U8},
		{"u16", U16},
		{"u32", U32},
		{"u64", U64},
		{"i8", I8},
		{"i16", I16},
		{"i32", I32},
		{"i64", I64},
		{"f32", F32},
		{"f64", F64},
		{"bool", BOOL},
		{"char", CHAR},
		{"string", STRING},
		{"(", L_PAREN},
		{")", R_PAREN},
		{"{", L_BRACE},
		{"}", R_BRACE},
		{"[", L_BRACKET},
		{"]", R_BRACKET},
		{"=", ASSIGN},
		{":", COLON},
		{":=", COLON_ASSIGN},
		{"..", RANGE},
		{"+", PLUS},
		{"+=", PLUS_ASSIGN},
		{"-", MINUS},
		{"-=", MINUS_ASSIGN},
		{"*", MULTIPLY},
		{"*=", MULTIPLY_ASSIGN},
		{"/", DIVIDE},
		{"/=", DIVIDE_ASSIGN},
		{"%", MODULO},
		{"%=", MODULO_ASSIGN},
		{"?", OPTIONAL},
		{"?=", OPTIONAL_ASSIGN},
		{"!", ERROR_MARK},
		{"~", BITWISE_NOT},
		{"&", BITWISE_AND},
		{"&=", BITWISE_AND_ASSIGN},
		{"|", BITWISE_OR},
		{"|=", BITWISE_OR_ASSIGN},
		{"^", BITWISE_XOR},
		{"^=", BITWISE_XOR_ASSIGN},
		{"<<", BITWISE_LEFT_SHIFT},
		{"<<=", BITWISE_LEFT_SHIFT_ASSIGN},
		{">>", BITWISE_RIGHT_SHIFT},
		{">>=", BITWISE_RIGHT_SHIFT_ASSIGN},
		{"#", ADDRESS_OF},
		{".", PERIOD},
		{",", COMMA},
		{"->", MAP_ARROW},
		{"=>", MATCH_ARROW},
		{"==", EQUAL},
		{"!=", NOT_EQUAL},
		{">", GREATER_THAN},
		{">=", GREATER_THAN_OR_EQUAL},
		{"<", LESS_THAN},
		{"<=", LESS_THAN_OR_EQUAL},
	}

	sort.Slice(words, func(i, j int) bool {
		return len(words[i].Word) > len(words[j].Word)
	})

	for _, word := range words {
		if len(l.source[l.pos:]) >= len(word.Word) && l.source[l.pos:l.pos+len(word.Word)] == word.Word {
			return Token{word.Type, l.slice(len(word.Word))}, nil
		}
	}

	return Token{_ANY_TOKEN, utils.String{}}, l.error("unknown token")
}
