pub mod ast;
use crate::lexer::token::{self, TokenType};

use self::ast::{Construction, Node};

#[derive(Debug)]
pub enum ParsingError {
    NonMatchable(token::TokenType),
    MatchNotFound((usize, token::Braket)),
}

/*
pub fn parse(tokens: &Vec<token::Token>) -> Result<Box<dyn ast::Node>, ParsingError> {
}
*/

pub fn init_tree(tokens: &Vec<token::Token>) -> Result<Vec<Box<ast::Construction>>, ParsingError> {
    let mut constructions = Vec::new();
    let mut cursor = 0;

    fn braket_construct(
        tokens: &Vec<token::Token>,
        cursor: &mut usize,
    ) -> Result<ast::Construction, ParsingError> {
        let origin = *cursor;
        let matching = match find_matching(&tokens, *cursor) {
            Ok(m) => m,
            Err(e) => return Err(e),
        };
        let content: Vec<Box<Construction>>;
        if matching > *cursor + 1 {
            content = match init_tree(&tokens[*cursor + 1..matching].to_vec()) {
                Ok(c) => c,
                Err(e) => return Err(e),
            };
        } else {
            content = Vec::new();
        }
        *cursor = matching + 1; // move cursor to after match
        match &tokens[origin].token {
            token::TokenType::Braket(token::Braket::OpenBrace) => Ok(Construction::Brace(content)),
            token::TokenType::Braket(token::Braket::OpenBraket) => {
                Ok(Construction::Braket(content))
            }
            token::TokenType::Braket(token::Braket::OpenParen) => Ok(Construction::Paren(content)),
            _s => Err(ParsingError::NonMatchable(*_s)),
        }
    }

    while cursor < tokens.len() {
        match tokens[cursor].token {
            token::TokenType::Braket(token::Braket::OpenBrace)
            | token::TokenType::Braket(token::Braket::OpenBraket)
            | token::TokenType::Braket(token::Braket::OpenParen) => {
                let con = match braket_construct(tokens, &mut cursor) {
                    // the cursor is already moved by the function
                    Ok(c) => c,
                    Err(e) => return Err(e),
                };

                constructions.push(Box::new(con));
            }
            _ => {
                let con = Construction::Token(tokens[cursor]);
                constructions.push(Box::new(con));
                cursor += 1;
            }
        }
    }

    return Ok(constructions);
}

fn find_matching(tokens: &Vec<token::Token>, idx: usize) -> Result<usize, ParsingError> {
    let open: token::Braket; // set one up to get the type

    let mut matching = idx + 1;

    match tokens[idx].token {
        token::TokenType::Braket(k) => open = k,
        _s => return Err(ParsingError::NonMatchable(_s)),
    };

    let close = match open {
        token::Braket::OpenBrace => token::Braket::CloseBrace,
        token::Braket::OpenBraket => token::Braket::CloseBraket,
        token::Braket::OpenParen => token::Braket::CloseParen,
        _s => return Err(ParsingError::NonMatchable(token::TokenType::Braket(_s))),
    };

    let mut counter = 1;

    while counter > 0 && matching < tokens.len() {
        match tokens[matching].token {
            token::TokenType::Braket(k) if k == close => counter -= 1,
            token::TokenType::Braket(k) if k == open => counter += 1,
            _ => (),
        };

        matching += 1;
    }

    if matching >= tokens.len() && counter != 0 {
        return Err(ParsingError::MatchNotFound((idx, open)));
    }

    return Ok(matching - 1);
}

pub fn print_constructions(cons: Vec<Box<ast::Construction>>, offset: usize) {
    for con in cons {
        println!("{}", (*con).display(offset));
        match (*con) {
            ast::Construction::Brace(k)
            | ast::Construction::Paren(k)
            | ast::Construction::Braket(k) => print_constructions(k, offset + 1),
            ast::Construction::Token(t) => {
                let mut spaces: Vec<String> = Vec::new();
                for i in 0..offset + 1 {
                    spaces.push(String::from("    "));
                }

                let off = spaces.join("");
                println!("{}{}", off, t);
            }
        };
    }
}
