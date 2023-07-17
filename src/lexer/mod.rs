pub mod token;

use token::*;
use regex::Regex;

#[derive(Debug)]
pub enum LexingError {
    UnevenRanges,
}

pub type SkipRange = Vec<Range>;

pub fn tokenize(program: String) -> Result<Vec<Token>, LexingError> {
    let mut tokens: Vec<Token> = Vec::new();
    let mut skip: SkipRange = Vec::new();

    fn lex(
        program: &String,
        tokens: &mut Vec<Token>, 
        skip: &mut SkipRange, 
        lexer: fn (&String, &SkipRange) -> Vec<Token>
    ) -> Option<LexingError> {
        let mut new_tokens = lexer(program, skip);
        let errs = update_skip(&new_tokens, skip);
        tokens.append(&mut new_tokens);
        return errs
    }

    let lexers = vec![
        comment_tokens,
        string_tokens,
        char_tokens,
        number_tokens,
        keyword_tokens,
        value_tokens,
        type_tokens,
        identifier_tokens,
        symbol_tokens,
        equal_and_assign_tokens,
        multiplication_and_deref_tokens,
        colon_and_double_tokens,
    ];

    for lx in lexers {
        match lex(&program, &mut tokens, &mut skip, lx) {
            Some(err) => return Err(err),
            None => continue,
        };
    }

    tokens.sort_unstable_by(|tokx, toky|
        if tokx.range[1] > toky.range[1] {
            std::cmp::Ordering::Greater
        } else if tokx.range[1] == toky.range[1] {
            std::cmp::Ordering::Equal
        } else {
            std::cmp::Ordering::Less
        }
    );
    return Ok(tokens)
}

fn keyword_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut keywords: Vec<Token> = Vec::new();    

    let re = Regex::new(r"(\bmut\b|\bfor\b|\bloop\b|\bwhile\b|\bif\b|\belse\b|\bswitch\b|\bin\b|\breturn\b|\bcatch\b|\bfn\b|\bis\b|\benum\b|\bimpl\b|\binterface\b|\bdebug\b|\band\b|\bor\b|\bnot\b|\bas\b|\bmod\b|\bimport\b)").unwrap();
    for capture in re.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            keywords.push(Token{
                token: match capture.as_str() {
                    "mut" => TokenType::Keyword(Keyword::Mut),
                    "for" => TokenType::Keyword(Keyword::For),
                    "loop" => TokenType::Keyword(Keyword::Loop),
                    "while" => TokenType::Keyword(Keyword::While),
                    "if" => TokenType::Keyword(Keyword::If),
                    "else" => TokenType::Keyword(Keyword::Else),
                    "switch" => TokenType::Keyword(Keyword::Switch),
                    "in" => TokenType::Keyword(Keyword::In),
                    "return" => TokenType::Keyword(Keyword::Return),
                    "catch" => TokenType::Keyword(Keyword::Catch),
                    "fn" => TokenType::Keyword(Keyword::Fn),
                    "is" => TokenType::Keyword(Keyword::Is),
                    "enum" => TokenType::Keyword(Keyword::Enum),
                    "impl" => TokenType::Keyword(Keyword::Impl),
                    "interface" => TokenType::Keyword(Keyword::Interface),
                    "debug" => TokenType::Keyword(Keyword::Debug),
                    "and" => TokenType::Keyword(Keyword::And),
                    "or" => TokenType::Keyword(Keyword::Or),
                    "not" => TokenType::Keyword(Keyword::Not),
                    "as" => TokenType::Keyword(Keyword::As),
                    "mod" => TokenType::Keyword(Keyword::Mod),
                    "import" => TokenType::Keyword(Keyword::Import),
                    _ => panic!("matched an unknown keyword: {}", capture.as_str()),
                },
                range: range_from_match(capture),
            });
        }
    }

    return keywords
}

fn value_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut values: Vec<Token> = Vec::new();    

    let re = Regex::new(r"(\bnone\b|\btrue\b|\bfalse\b|\berror\b)").unwrap();
    for capture in re.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            values.push(Token{
                token: match capture.as_str() {
                    "none" => TokenType::Value(Value::None),
                    "true" => TokenType::Value(Value::True),
                    "false" => TokenType::Value(Value::False),
                    "error" => TokenType::Value(Value::Error),
                    _ => panic!("matched an unknown keyword: {}", capture.as_str()),
                },
                range: range_from_match(capture),
            });
        }
    }

    return values
}

fn colon_and_double_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut colons: Vec<Token> = Vec::new();
    let double_colon = Regex::new(r"::").unwrap();

    for capture in double_colon.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            colons.push(Token{
                token: TokenType::Symbol(Symbol::DoubleColon),
                range: range_from_match(capture),
            });
        }
    }
    
    let single_colon = Regex::new(r"[^:]:[^:]").unwrap();

    for capture in single_colon.find_iter(program) {
        let mut rg = range_from_match(capture);
        rg[1] -= 1; // we don't want the matched "[^:]" character
        rg[0] += 1; // we don't want the matched "[^:]" character
        if !range_intersects_skip(rg, skip) {
            colons.push(Token{
                token: TokenType::Symbol(Symbol::Colon),
                range: rg,
            });
        }
    }
    return colons
}

fn multiplication_and_deref_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut deref_and_mult: Vec<Token> = Vec::new();
    let deref = Regex::new(r"\*(\w|\*)").unwrap();

    for capture in deref.find_iter(program) {
        let mut rg = range_from_match(capture);
        rg[1] -= 1; // we don't want the matched "word" character
        if !range_intersects_skip(rg, skip) {
            deref_and_mult.push(Token{
                token: TokenType::Symbol(Symbol::Dereference),
                range: rg,
            });
        }
    }
    
    let mult = Regex::new(r"\*[^a-zA-Z\*]").unwrap();

    for capture in mult.find_iter(program) {
        let mut rg = range_from_match(capture);
        rg[1] -= 1; // we don't want the matched "word" character
        if !range_intersects_skip(rg, skip) {
            deref_and_mult.push(Token{
                token: TokenType::Operator(Operator::Mult),
                range: rg,
            });
        }
    }
    return deref_and_mult
}

fn equal_and_assign_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut equal_and_assigns: Vec<Token> = Vec::new();
    let equal = Regex::new(r"==").unwrap();

    for capture in equal.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            equal_and_assigns.push(Token{
                token: TokenType::Compare(Compare::Equal),
                range: range_from_match(capture),
            });
        }
    }
    
    let assign = Regex::new(r"[^=]=[^=]").unwrap();

    for capture in assign.find_iter(program) {
        let mut rg = range_from_match(capture);
        rg[1] -= 1; // we don't want the matched "[^=]" character
        rg[0] += 1; // we don't want the matched "[^=]" character
        if !range_intersects_skip(rg, skip) {
            equal_and_assigns.push(Token{
                token: TokenType::Symbol(Symbol::Assign),
                range: rg,
            });
        }
    }
    return equal_and_assigns
}

fn symbol_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut symbols: Vec<Token> = Vec::new();    

    let re = Regex::new(r"(!|\?|\.|&|\||->|\.\.|\.\.\.|\{|\}|\[|\]|\(|\)|\+|-|/|\^|%|\+=|-=|\*=|/=|\^=|>=|<=|>|<)").unwrap();


    for capture in re.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            symbols.push(Token{
                token: match capture.as_str() {
                    "!" =>    TokenType::Symbol(Symbol::Bang), 
                    "?" =>    TokenType::Symbol(Symbol::Optional), 
                    "." =>    TokenType::Symbol(Symbol::Dot),
                    "," =>    TokenType::Symbol(Symbol::Comma),
                    "&" =>    TokenType::Symbol(Symbol::Address), 
                    "|" =>    TokenType::Symbol(Symbol::TypeSum), 
                    "->" =>    TokenType::Symbol(Symbol::Arrow),
                    ".." =>      TokenType::Symbol(Symbol::Range),   
                    "..." =>    TokenType::Symbol(Symbol::Elipsis), 
                    "{" =>    TokenType::Braket(Braket::OpenBrace),
                    "}" =>    TokenType::Braket(Braket::CloseBrace),
                    "[" =>    TokenType::Braket(Braket::OpenBraket),
                    "]" =>    TokenType::Braket(Braket::CloseBraket),
                    "(" =>    TokenType::Braket(Braket::OpenParen),
                    ")" =>    TokenType::Braket(Braket::CloseParen),
                    "+" =>    TokenType::Operator(Operator::Plus),
                    "-" =>    TokenType::Operator(Operator::Minus),
                    "/" =>    TokenType::Operator(Operator::Div),
                    "^" =>    TokenType::Operator(Operator::Power),
                    "%" =>    TokenType::Operator(Operator::Modulus),
                    "+=" =>    TokenType::Operator(Operator::PlusAssign),
                    "-=" =>    TokenType::Operator(Operator::MinusAssign),
                    "*=" =>    TokenType::Operator(Operator::MultAssign),
                    "/=" =>    TokenType::Operator(Operator::DivAssign),
                    "^=" =>    TokenType::Operator(Operator::PowerAssign),
                    ">=" =>    TokenType::Compare(Compare::GreaterEqual),
                    "<=" =>    TokenType::Compare(Compare::LessEqual),
                    ">" =>    TokenType::Compare(Compare::Greater),
                    "<" =>    TokenType::Compare(Compare::Less),
                    _ => panic!("matched an unknown symbol: {}", capture.as_str()),
                },
                range: range_from_match(capture),
            });
        }
    }

    return symbols
}
fn number_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut numbers: Vec<Token> = Vec::new();    
    let int = Regex::new(r"\b[0-9]+\b").unwrap();
    let float = Regex::new(r"\b[0-9]+\.[0-9]+(e(\+|-)[0-9]+)?").unwrap();
    let hex = Regex::new(r"\b0x[0-9a-fA-F]+\b").unwrap();
    let oct = Regex::new(r"\b0o[0-7]+\b").unwrap();
    let bin = Regex::new(r"\b0b[0-1]+\b").unwrap();

    for capture in int.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            numbers.push(Token{
                token: TokenType::Value(Value::Int),
                range: range_from_match(capture),
            });
        }
    }
    for capture in float.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            numbers.push(Token{
                token: TokenType::Value(Value::Float),
                range: range_from_match(capture),
            });
        }
    }
    for capture in hex.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            numbers.push(Token{
                token: TokenType::Value(Value::Hex),
                range: range_from_match(capture),
            });
        }
    }
    for capture in oct.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            numbers.push(Token{
                token: TokenType::Value(Value::Oct),
                range: range_from_match(capture),
            });
        }
    }
    for capture in bin.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            numbers.push(Token{
                token: TokenType::Value(Value::Bin),
                range: range_from_match(capture),
            });
        }
    }
    
    return numbers
}

fn identifier_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut identifiers: Vec<Token> = Vec::new();    
    let re = Regex::new(r"[a-z][a-zA-Z0-9_]*").unwrap();

    for capture in re.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            identifiers.push(Token{
                token: TokenType::Identifier,
                range: range_from_match(capture),
            });
        }
    }
    
    return identifiers
}

fn type_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut types: Vec<Token> = Vec::new();    
    let user_type = Regex::new(r"(\[\])*[A-Z][a-zA-Z]*").unwrap();
    let builtin_type = Regex::new(r"(\[\])*(u8|u16|u32|u64|i32|i64|f32|f64|str|usize|isize|bool)").unwrap();

    for capture in user_type.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            types.push(Token{
                token: TokenType::Type,
                range: range_from_match(capture),
            });
        }
    }
    for capture in builtin_type.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            types.push(Token{
                token: TokenType::Type,
                range: range_from_match(capture),
            });
        }
    }
    
    return types
}


fn string_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut strings: Vec<Token> = Vec::new();    
    let re = Regex::new("\"([^\"]|\\\\\")*\"").unwrap();

    for capture in re.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            strings.push(Token{
                token: TokenType::Value(Value::String),
                range: range_from_match(capture),
            });
        }
    }
    
    return strings
}

fn char_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut chars: Vec<Token> = Vec::new();    
    let re = Regex::new("'[^']+'").unwrap();

    for capture in re.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            chars.push(Token{
                token: TokenType::Value(Value::Char),
                range: range_from_match(capture),
            });
        }
    }
    
    return chars
}

fn comment_tokens(program: &String, skip: &SkipRange) -> Vec<Token> {
    let mut comments: Vec<Token> = Vec::new();
    let single_line = Regex::new(r"--.*").unwrap();

    for capture in single_line.find_iter(program) {
        if !range_intersects_skip(range_from_match(capture), skip) {
            comments.push(Token{
                token: TokenType::Comment,
                range: range_from_match(capture)
            });
        }
    }

    return comments
}


/*
fn inc_skip(num: usize, skip: &SkipRange) -> usize {
    let mut cursor = num;
    for r in skip {
        if cursor >= r[0] && cursor <= r[1] {
            cursor = r[1]
        }
    }
    return cursor+1
}
*/

// creates a sequence of ranges for lexers to skip
fn update_skip(new_tokens: &Vec<Token>, skip: &mut SkipRange) -> Option<LexingError> {
    let mut indices: Vec<usize> = Vec::new();

    if new_tokens.is_empty() {
        return None;
    }

    for token in new_tokens {
        indices.push(token.range[0]);
        indices.push(token.range[1]);
    }

    indices.sort_unstable();
    let mut cursor = 0;
    while cursor < indices.len()-1 {
        if indices[cursor] == indices[cursor+1] {
            indices.remove(cursor);
            indices.remove(cursor); // we want to remove the i+1 but now its at i given we removed its predecessor
        }
        cursor += 1;
    }
    cursor = 0;
    while cursor < indices.len()-1 {
        skip.push([indices[cursor], indices[cursor+1]]);
        cursor += 2;
    }

    if indices.len() % 2 != 0 {
        println!("Indices: {:?}", indices);
        return Some(LexingError::UnevenRanges);
    }
    return None
    
}

/*
fn sort_by_length_decreasing(list: &mut Vec<&str>) -> () {
    list.sort_unstable_by(|x, y| 
        if x.len() > y.len() { std::cmp::Ordering::Greater } 
        else if x.len() == y.len() { std::cmp::Ordering::Equal } 
        else {std::cmp::Ordering::Less});
}
*/

fn range_intersects_skip(r: Range, s: &SkipRange) -> bool {
    fn is_in_range(loc: usize, r: Range) -> bool {
        if loc >= r[0] && loc < r[1] {
            return true
        } else {
            return false
        }
    }

    for &rng in s {
        if is_in_range(r[0], rng) || is_in_range(r[1]-1, rng) {
            return true
        }
    }
    return false
}

fn range_from_match(cap: regex::Match) -> Range {
    return [cap.start(), cap.end()]
}
