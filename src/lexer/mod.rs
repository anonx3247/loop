pub mod token;

use token::Token;

pub fn tokenize(program: String) -> Result<Vec<Token>, LexError> {
    let mut tokens: Vec<Token> = Vec::new();

    let program_end = program.len();

    let mut start = match get_next_non_space_char(program, 0) {
        Ok(pos) => pos,
        Err(e) => return Ok(tokens),
    };
    

    let mut i = start;
    let mut non_space_str = "";
    let mut new_word = false;

    while i < program_end {
        match is_empty_at(program, i) {
            Ok(true) => {
                match &mut lex_tokens(program, start, i){
                    Ok(v) => tokens.append(v),
                    Err(e) => return Err(*e)
                };
                start = match get_next_non_space_char(program, i) {
                    Ok(pos) => pos,
                    Err(e) => return Ok(tokens),
                };
                i = start;
            },
            Ok(false) => {i += 1; break},
            Err(e) => return Err(e)
        }

        
    }
    
    return Ok(tokens)
}


fn get_next_non_space_char(s: String, loc: usize) -> Result<usize, LexError> {
        let mut start = loc;
        let start_char = s.get(start..start+1);
        let mut reached_beginning = false;

        while !reached_beginning && start < s.len() {
            match is_empty_at(s, start) {
                Ok(true) => start += 1,
                Ok(false) => reached_beginning = true,
                Err(e) => return Err(e)
            };
        }

        if start >= s.len() {
            return Err(LexError::EOF)
        }

        return Ok(start)
}

fn is_empty_at(s: String, loc: usize) -> Result<bool, LexError> {
    let char = s.get(loc..loc+1);
    match char {
        Some(" ") | Some("\t") | Some("\n") => Ok(true),
        None => Err(LexError::EOF),
        _ => Ok(false),
    }
}

fn lex_tokens(program: String, start: usize, end: usize) -> Result<Vec<Token>, LexError> {
    return Err(LexError::EOF)
}

enum LexError {
    EOF,
}
