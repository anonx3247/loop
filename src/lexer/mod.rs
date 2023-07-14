pub mod token;

use token::Token;
use token::TokenType;
use token::Range;

pub enum LexingError {
    EOF,
    UnevenRanges,
    UnclosedComment,
}

pub type SkipRange = Vec<Range>;

pub fn tokenize(program: String) -> Result<Vec<Token>, LexingError> {
    let mut tokens: Vec<Token> = Vec::new();

    let mut skip = match skip(tokens) {
        Ok(sk) => sk,
        Err(e) => return Err(e),
    };

    let mut start = match next_not_of(vec![" ", "\t", "\n"], program, 0, skip) {
        Ok(pos) => pos,
        Err(e) => return Ok(tokens),
    };
    

    let mut i = start;
    let mut non_space_str = "";
    let mut new_word = false;

    let mut comments = match extract_single_line_comments(program, start, skip) {
        Ok(comm) => comm,
        Err(e) => return Err(e)
    };

    tokens.append(&mut comments);
    comments = match extract_multi_line_comments(program, start, skip) {
        Ok(comm) => comm,
        Err(e) => return Err(e)
    };

    tokens.append(&mut comments);
    return Ok(tokens)
}

fn extract_multi_line_comments(program: String, start: usize, skip: SkipRange) -> Result<Vec<Token>, LexingError> {
    let mut i  = start;
    let mut comments: Vec<Token> = Vec::new();

    while i < program.len() {
        let next_comment_beg = match next("---", program, i, skip) {
            Ok(val) => val,
            Err(e) => return Err(e),
        };
        let next_comment_end = match next("---", program, next_comment_beg, skip) {
            Ok(val) => val,
            Err(e) => return Err(LexingError::UnclosedComment),
        };

        let comment = Token{
            token: match program.get(next_comment_beg..next_comment_end+1) {
                Some(str) => TokenType::Comment(String::from(str)),
                None => return Err(LexingError::EOF),
            },
            range: [next_comment_beg, next_comment_end],
        };

        comments.push(comment);

        i = inc_skip(i, skip);
        
    }
    return Ok(comments)
}

fn extract_single_line_comments(program: String, start: usize, skip: SkipRange) -> Result<Vec<Token>, LexingError> {
    let mut i  = start;
    let mut comments: Vec<Token> = Vec::new();

    while i < program.len() {
        let next_comment = match next("--", program, i, skip) {
            Ok(val) => val,
            Err(e) => return Err(e),
        };
        let next_newline = match next("\n", program, next_comment, skip) {
            Ok(val) => val,
            Err(e) => return Err(e),
        };

        let comment = Token{
            token: match program.get(next_comment..next_newline+1) {
                Some(str) => TokenType::Comment(String::from(str)),
                None => return Err(LexingError::EOF),
            },
            range: [next_comment, next_newline],
        };

        comments.push(comment);

        i = inc_skip(i, skip);
        
    }
    return Ok(comments)
}

// Finds the next `seq` in `program` starting at `start`
fn next(seq: &str, program: String, start: usize, skip: Vec<Range>) -> Result<usize, LexingError> {
    let mut i = start;
    let mut found_start = false;
    while i < program.len() {
        if !found_start {
            match program.get(i..i+1) {
                Some(char) => if &seq[0..1] == char {found_start = true},
                None => return Err(LexingError::EOF),
            }
        } else {
            match program.get(i..i+seq.len()) {
                Some(s) => if seq == s {return Ok(i)} else {found_start = false},
                None => return Err(LexingError::EOF),
            }
        }
        i = inc_skip(i, skip);
    }

    if i >= program.len() {
        return Err(LexingError::EOF);
    }

    return Ok(i)
}

fn next_of(sequences: Vec<&str>, program: String, start: usize, skip: SkipRange) -> Result<usize, LexingError> {
    let mut i = start;
    let mut found_start = false;

    let lengths = Vec::from_iter(sequences.iter().map(|seq| seq.len()));
    let max_size = max(lengths);
    
    while i < program.len() {
        if !found_start {
            match program.get(i..i+1) {
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
            match program.get(i..i+max_size) {
                Some(s) => {
                    for seq in sequences {
                        if seq == &s[0..seq.len()] {
                            return Ok(i);
                        }
                    }
                    found_start = false;
                },
                None => return Err(LexingError::EOF),
            }
        }
        i = inc_skip(i, skip);
    }

    if i >= program.len() {
        return Err(LexingError::EOF);
    }

    return Ok(i)
}

fn next_not_of(sequences: Vec<&str>, program: String, start: usize, skip: SkipRange) -> Result<usize, LexingError> {
    let mut i = start;
    let mut found_start = false;

    let lengths = Vec::from_iter(sequences.iter().map(|seq| seq.len()));
    let max_size = max(lengths);
    
    while i < program.len() {
        if !found_start {
            match program.get(i..i+1) {
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
            match program.get(i..i+max_size) {
                Some(s) => {
                    let mut found_seq = false;
                    for seq in sequences {
                        if seq == &s[0..seq.len()] {
                            i = i +seq.len();
                            found_seq = true;
                            break
                        }
                    }
                    if !found_seq {
                        return Ok(i)
                    }
                    found_start = false;
                },
                None => return Err(LexingError::EOF),
            }
        }
        i = inc_skip(i, skip);
    }

    if i >= program.len() {
        return Err(LexingError::EOF);
    }

    return Ok(i)
}

// inc_skip expects a simplified skip sequence
fn inc_skip(num: usize, skip: SkipRange) -> usize {
    let mut i = num;
    for r in skip {
        if i >= r[0] && i <= r[1] {
            return r[1]+1
        }
    }
    return i
}

// creates a sequence of ranges for lexers to skip
fn skip(tokens: Vec<Token>) -> Result<SkipRange, LexingError> {
    let mut indices: Vec<usize> = Vec::new();
    let mut skip: Vec<Range> = Vec::new();

    for token in tokens {
        indices.push(token.range[0]);
        indices.push(token.range[1]);
    }

    indices.sort_unstable();
    let mut i = 0;
    while i < indices.len()-1 {
        if indices[i] == indices[i+1] {
            indices.remove(i);
            indices.remove(i); // we want to remove the i+1 but now its at i given we removed its predecessor
        }
        i += 1;
    }
    i = 0;
    while i < indices.len()-1 {
        skip.push([indices[i], indices[i+1]]);
    }

    if indices.len() % 2 == 0 {
        return Err(LexingError::UnevenRanges);
    }
    return Ok(skip)
    
}

fn max(list: Vec<usize>) -> usize {
    let mut maximum = list[0];
    for elem in list {
        if elem > maximum {
            maximum = elem
        }
    }
    return maximum
}
