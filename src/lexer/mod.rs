pub mod token;

use token::*;

#[derive(Debug)]
pub enum LexingError {
    EOF,
    UnevenRanges,
    Unclosed(TokenType),
    Unknown,
}

pub type SkipRange = Vec<Range>;

pub fn tokenize(program: String) -> Result<Vec<Token>, LexingError> {
    let mut tokens: Vec<Token> = Vec::new();

    let mut skip = match skip_list(&tokens) {
        Ok(sk) => sk,
        Err(e) => return Err(e),
    };

    let mut single_line_comments = match extract_single_line_comments(&program, &skip) {
        Ok(comm) => comm,
        Err(e) => return Err(e)
    };

    tokens.append(&mut single_line_comments);
    skip = match skip_list(&tokens) {
        Ok(sk) => sk,
        Err(e) => return Err(e),
    };
    let mut multi_line_comments = match extract_multi_line_comments(&program, &skip) {
        Ok(comm) => comm,
        Err(e) => return Err(e)
    };

    tokens.append(&mut multi_line_comments);
    skip = match skip_list(&tokens) {
        Ok(sk) => sk,
        Err(e) => return Err(e),
    };

    let mut strings = match extract_strings(&program, &skip) {
        Ok(s) => s,
        Err(e) => return Err(e)
    };

    tokens.append(&mut strings);
    skip = match skip_list(&tokens) {
        Ok(sk) => sk,
        Err(e) => return Err(e),
    };

    let mut chars = match extract_chars(&program, &skip) {
        Ok(s) => s,
        Err(e) => return Err(e)
    };

    tokens.append(&mut chars);
    skip = match skip_list(&tokens) {
        Ok(sk) => sk,
        Err(e) => return Err(e),
    };

    let mut keywords = match extract_keywords(&program, &skip) {
        Ok(s) => s,
        Err(e) => return Err(e)
    };

    tokens.append(&mut keywords);
    skip = match skip_list(&tokens) {
        Ok(sk) => sk,
        Err(e) => return Err(e),
    };

    let mut symbols = match extract_symbols(&program, &skip) {
        Ok(s) => s,
        Err(e) => return Err(e)
    };

    tokens.append(&mut symbols);
    skip = match skip_list(&tokens) {
        Ok(sk) => sk,
        Err(e) => return Err(e),
    };
    return Ok(tokens)
}
fn extract_keywords(program: &String, skip: &SkipRange) -> Result<Vec<Token>, LexingError> {
    let mut keywords: Vec<Token> = Vec::new();    
    let mut cursor = 0;

    let words = vec![
        "mut",
        "for",
        "loop",
        "while",
        "if",
        "else",
        "switch",
        "in",
        "return",
        "catch",
        "fn",
        "is",
        "enum",
        "impl",
        "interface",
        "debug",
        "and",
        "or",
        "not",
        "as",
    ];

    while cursor < program.len() {
        let (keyword_start, keyword) = match next_word_of(&words, program, cursor, &skip) {
            Ok(Some(i)) => i,
            Err(e) => return Err(e),
            Ok(None) => break,
        };

        keywords.push(Token{
            token: match keyword.as_str() {
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
                _ => return Err(LexingError::Unknown),
            },
            range: [keyword_start, keyword_start+keyword.len()],
        });

        cursor = inc_skip(keyword_start+keyword.len(), &&skip)
       
    }

    return Ok(keywords)
}

fn extract_symbols(program: &String, skip: &SkipRange) -> Result<Vec<Token>, LexingError> {
    let mut symbols: Vec<Token> = Vec::new();    
    let mut cursor = 0;

    let symbol_str = vec![
        ":",
        "!",
        "?",
        ".",
        ",",
        "&",
        "*",
        ":=",
        "=",
        "|",
        "->",
        "..",  
        "...",
        "{",
        "}",
        "[",
        "]",
        "(",
        ")",
        "+",
        "-",
        "*",
        "/",
        "^",
        "%",
        "+=",
        "-=",
        "*=",
        "/=",
        "^=",
        ">=",
        "<=",
        ">",
        "<",
        "=="
    ];

    while cursor < program.len() {
        let (symbol_start, symbol) = match next_of(&symbol_str, program, cursor, &skip) {
            Ok(Some(i)) => i,
            Err(e) => return Err(e),
            Ok(None) => break,
        };

        symbols.push(Token{
            token: match symbol.as_str() {
                ":" =>    TokenType::Symbol(Symbol::Colon),
                "!" =>    TokenType::Symbol(Symbol::Bang), 
                "?" =>    TokenType::Symbol(Symbol::Optional), 
                "." =>    TokenType::Symbol(Symbol::Dot),
                "," =>    TokenType::Symbol(Symbol::Comma),
                "&" =>    TokenType::Symbol(Symbol::Dereference), 
                //"*" =>    TokenType::Symbol(Symbol::Address),     
                ":=" =>    TokenType::Symbol(Symbol::Assign),
                "=" =>    TokenType::Symbol(Symbol::Equal),
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
                //"*" =>    TokenType::Operator(Operator::Mult),
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
                "==" =>   TokenType::Compare(Compare::Equal),
                _ => return Err(LexingError::Unknown),
            },
            range: [symbol_start, symbol_start+symbol.len()],
        });

        cursor = inc_skip(symbol_start+symbol.len(), &&skip)
       
    }

    return Ok(symbols)
}

fn extract_strings(program: &String, skip: &SkipRange) -> Result<Vec<Token>, LexingError> {
    let mut strings: Vec<Token> = Vec::new();    
    let mut cursor = 0;

    while cursor < program.len() {
        let str_start = match next("\"", program, cursor, &skip) {
            Ok(Some(i)) => i,
            Err(e) => return Err(e),
            Ok(None) => break,
        };

        // a string " can be preceeded by a \
        let mut unmatched = true;
        let mut int_cursor = str_start;
        let mut str_end = 0;
        while unmatched {
            str_end = match next("\"", program, int_cursor, &skip) {
                Ok(Some(i)) => i,
                Err(e) => return Err(e),
                Ok(None) => return Err(LexingError::Unclosed(TokenType::Value(Value::String)))
            };

            match program.get(str_end-1..str_end) {
                Some("\\") => int_cursor = inc_skip(int_cursor, &&skip),
                None => return Err(LexingError::EOF),
                _ => unmatched = false,
            }
        }

        let string = Token{
            token: match program.get(str_start..str_end+1) {
                Some(_s) => TokenType::Value(token::Value::String),
                None => return Err(LexingError::EOF),
            },
            range: [str_start, str_end],
        };

        strings.push(string);
        cursor = inc_skip(str_end, &&skip)
       
    }

    return Ok(strings)
}

fn extract_chars(program: &String, skip: &SkipRange) -> Result<Vec<Token>, LexingError> {
    let mut chars: Vec<Token> = Vec::new();    
    let mut cursor = 0;

    while cursor < program.len() {
        let char_start = match next("'", program, cursor, &skip) {
            Ok(Some(i)) => i,
            Err(e) => return Err(e),
            Ok(None) => break,
        };

        // a string " can be preceeded by a \
        let mut unmatched = true;
        let mut int_cursor = char_start;
        let mut char_end = 0;
        while unmatched {
            char_end = match next("'", program, int_cursor, &skip) {
                Ok(Some(i)) => i,
                Err(e) => return Err(e),
                Ok(None) => return Err(LexingError::Unclosed(TokenType::Value(Value::Char))),
            };

            match program.get(char_end-1..char_end) {
                Some("\\") => int_cursor = inc_skip(int_cursor, &&skip),
                None => return Err(LexingError::EOF),
                _ => unmatched = false,
            }
        }

        let char = Token{
            token: match program.get(char_start..char_end+1) {
                Some(_s) => TokenType::Value(Value::String),
                None => return Err(LexingError::EOF),
            },
            range: [char_start, char_end],
        };

        chars.push(char);
        cursor = inc_skip(char_end, &&skip)
       
    }

    return Ok(chars)
}
fn extract_multi_line_comments(program: &String, skip: &SkipRange) -> Result<Vec<Token>, LexingError> {
    let mut cursor  = 0;
    let mut comments: Vec<Token> = Vec::new();

    while cursor < program.len() {
        let next_comment_beg = match next("---", program, cursor, &skip) {
            Ok(Some(val)) => val,
            Err(e) => return Err(e),
            Ok(None) => break,
        };
        let next_comment_end = match next("---", program, next_comment_beg, &skip) {
            Ok(Some(val)) => val,
            Err(e) => return Err(e),
            Ok(None) => return Err(LexingError::Unclosed(TokenType::Comment)),
        };

        let comment = Token{
            token: match program.get(next_comment_beg..next_comment_end+1) {
                Some(_s) => TokenType::Comment,
                None => return Err(LexingError::EOF),
            },
            range: [next_comment_beg, next_comment_end],
        };

        comments.push(comment);
        cursor = inc_skip(next_comment_end, &&skip);
        
    }
    return Ok(comments)
}

fn extract_single_line_comments(program: &String, skip: &SkipRange) -> Result<Vec<Token>, LexingError> {
    let mut cursor  = 0;
    let mut comments: Vec<Token> = Vec::new();

    while cursor < program.len() {
        let comment_pos = match next("--", program, cursor, &skip) {
            Ok(Some(val)) => val,
            Err(e) => return Err(e),
            Ok(None) => break,
        };
        let newline = match next("\n", program, comment_pos, &skip) {
            Ok(Some(val)) => val,
            Err(e) => return Err(e),
            Ok(None) => program.len()-1, // means we are at end of program
        };

        let comment = Token{
            token: match program.get(comment_pos..newline+1) {
                Some(_s) => TokenType::Comment,
                None => return Err(LexingError::EOF),
            },
            range: [comment_pos, newline],
        };

        comments.push(comment);
        cursor = inc_skip(newline, &&skip);
        
    }
    return Ok(comments)
}

// Finds the next `seq` in `program` starting at `start`
fn next(seq: &str, program: &String, start: usize, skip: &SkipRange) -> Result<Option<usize>, LexingError> {
    let mut cursor = start;
    let mut found_start = false;
    while cursor < program.len() {
        if !found_start {
            match program.get(cursor..cursor+1) {
                Some(char) => if &seq[0..1] == char {found_start = true},
                None => return Err(LexingError::EOF),
            }
        } else {
            match program.get(cursor..cursor+seq.len()) {
                Some(s) => if seq == s {return Ok(Some(cursor))} else {found_start = false},
                None => return Err(LexingError::EOF),
            }
        }
        cursor = inc_skip(cursor, &&skip);
    }

    return Err(LexingError::EOF)
}

fn next_word_of(words: &Vec<&str>, program: &String, start: usize, skip: &SkipRange) -> Result<Option<(usize, String)>, LexingError> {
    let mut cursor = start;
    let length = program.len();

    while cursor < length {
        let (next, word) = match next_of(words, program, cursor, &skip) {
            Ok(Some(a)) => a,
            Err(e) => return Err(e),
            Ok(None) => return Ok(None),
        };

        if is_whitespace_at(next-1, &program)  && is_whitespace_at(next+word.len(), &program) {
            return Ok(Some((next, word)))
        } else {
            cursor = inc_skip(next+word.len(), &&skip);
        }
    }
    return Ok(None)
}


fn next_of(sequences: &Vec<&str>, program: &String, start: usize, skip: &SkipRange) -> Result<Option<(usize, String)>, LexingError> {
    let mut cursor = start;
    let mut found_start = false;

    let lengths = Vec::from_iter(sequences.iter().map(|seq| seq.len()));
    let max_size = max(&lengths);
    
    while cursor < program.len() {
        if !found_start {
            match program.get(cursor..cursor+1) {
                Some(char) => {
                    for seq in sequences {
                        if &seq[0..1] == char {
                            found_start = true;
                            break
                        }
                    }
                },
                None => return Err(LexingError::EOF),
            }
        } else {
            match program.get(cursor..cursor+max_size) {
                Some(s) => {
                    for seq in sequences {
                        if *seq == &s[0..seq.len()] {
                            return Ok(Some((cursor, String::from(*seq))));
                        }
                    }
                    found_start = false;
                },
                None => return Err(LexingError::EOF),
            }
        }
        cursor = inc_skip(cursor, &&skip);
    }

    return Ok(None)
}

fn next_not_of(sequences: &Vec<&str>, program: &String, start: usize, skip: &SkipRange) -> Result<Option<usize>, LexingError> {
    let mut cursor = start;
    let mut found_start = false;

    let lengths = Vec::from_iter(sequences.iter().map(|seq| seq.len()));
    let max_size = max(&lengths);
    
    while cursor < program.len() {
        if !found_start {
            match program.get(cursor..cursor+1) {
                Some(char) => {
                    for seq in sequences {
                        if &seq[0..1] == char {
                            found_start = true;
                            break
                        }
                    }
                },
                None => return Err(LexingError::EOF),
            }
        } else {
            match program.get(cursor..cursor+max_size) {
                Some(s) => {
                    let mut found_seq = false;
                    for seq in sequences {
                        if *seq == &s[0..seq.len()] {
                            cursor = cursor +seq.len();
                            found_seq = true;
                            break
                        }
                    }
                    if !found_seq {
                        return Ok(Some(cursor))
                    }
                    found_start = false;
                },
                None => return Err(LexingError::EOF),
            }
        }
        cursor = inc_skip(cursor, &&skip);
    }

    return Ok(None)
}

// inc_skip expects a simplified skip sequence
fn inc_skip(num: usize, skip: &SkipRange) -> usize {
    let mut cursor = num;
    for r in skip {
        if cursor >= r[0] && cursor <= r[1] {
            cursor = r[1]+1
        }
    }
    return cursor
}

// creates a sequence of ranges for lexers to skip
fn skip_list(tokens: &Vec<Token>) -> Result<SkipRange, LexingError> {
    let mut indices: Vec<usize> = Vec::new();
    let mut skip: Vec<Range> = Vec::new();

    for token in tokens {
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
    }

    if indices.len() % 2 == 0 {
        return Err(LexingError::UnevenRanges);
    }
    return Ok(skip)
    
}

fn max(list: &Vec<usize>) -> usize {
    let mut maximum = list[0];
    for elem in list {
        if *elem > maximum {
            maximum = *elem
        }
    }
    return maximum
}

fn is_whitespace_at(loc: usize, program: &String) -> bool {
    if loc == 0 || loc == program.len() {
        true
    } else {
       match program.get(loc..loc+1) {
            Some(" ") | Some("\t") | Some("\n") | Some("\r") => true,
            None => true, // only possible if out of range
            _ => false,
        }
    }
}
