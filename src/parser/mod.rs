pub mod ast;
use std::collections::HashMap;

use crate::lexer::token;

use ast::{Brace, Construction, Node};

type Tree = Vec<Box<Node>>;

#[derive(Debug)]
pub enum ParsingError {
    NonMatchable(token::TokenType),
    MatchNotFound((usize, token::Braket)),
    Unparsable(String),
}

pub fn parse(tokens: &Vec<token::Token>, program: &String) -> Result<Tree, ParsingError> {
    let mut tree: Tree = match generate_braket_heirarchy(tokens) {
        Ok(c) => c.iter().map(|k| k.clone() as Box<Node>).collect(),
        Err(e) => return Err(e),
    };

    tree = match parse_literals(&tree, program) {
        Ok(t) => t,
        Err(e) => return Err(e),
    };

    tree = match generate_scopes(&tree, program) {
        Ok((tr, id)) => vec![Box::new(Node::Construction(Construction::Brace(Brace {
            identifiers: id,
            content: tr,
        })))],
        Err(e) => return Err(e),
    };

    /*
    tree = match parse_function_calls(&tree, program) {
        Ok(t) => t,
        Err(e) => return Err(e),
    };
    */

    return Ok(tree);
}

/*
fn parse_function_calls(tree: &Tree, program: &String) -> Result<Tree, ParsingError> {
    let mut i = 0;
    let mut new_tree: Tree = Vec::new();

    // to be defined
    fn extract_args() {}

    while i < tree.len() - 1 {
        if let (
            Node::Construction(Construction::Token(token::Token {
                token: token::TokenType::Identifier,
                range: r,
            })),
            Node::Construction(Construction::Paren(p)),
        ) = (tree[i].as_ref(), tree[i + 1].as_ref())
        {
            let fun = Node::Undetermined(ast::Undetermined::FnCall(ast::FnCall {
                identifier: program[r[0]..r[1]].to_string(),
                arguments: extract_args(p, program),
            }));
            i += 1
        } else {
            new_tree.push(match tree[i].as_ref() {
                Node::Construction(Construction::Brace(Brace {
                    identifiers: id,
                    content: c,
                })) => match parse_function_calls(c, program) {
                    Ok(v) => Box::new(Node::Construction(Construction::Brace(Brace {
                        identifiers: *id,
                        content: v,
                    }))),
                    Err(e) => return Err(e),
                },
                Node::Construction(Construction::Paren(k)) => {
                    match parse_function_calls(k, program) {
                        Ok(v) => Box::new(Node::Construction(Construction::Paren(v))),
                        Err(e) => return Err(e),
                    }
                }
                Node::Construction(Construction::SquareBraket(k)) => {
                    match parse_function_calls(k, program) {
                        Ok(v) => Box::new(Node::Construction(Construction::SquareBraket(v))),
                        Err(e) => return Err(e),
                    }
                }
                _s => Box::new(*_s),
            });
        }
        i += 1
    }

    return Ok(new_tree);
}
    */

fn generate_scopes(
    tree: &Tree,
    program: &String,
) -> Result<(Tree, Option<ast::IdentMap>), ParsingError> {
    let mut new_tree: Tree = Vec::new();
    let mut idents: ast::IdentMap = HashMap::new();
    for node in tree {
        if let Node::Construction(Construction::Brace(t)) = node.as_ref() {
            new_tree.push(match generate_scopes(&t.content, &program) {
                Ok((tr, id)) => Box::new(Node::Construction(Construction::Brace(Brace {
                    identifiers: id,
                    content: tr,
                }))),
                Err(e) => return Err(e),
            })
        } else if let Node::Construction(Construction::Token(token::Token { token: t, range: r })) =
            node.as_ref()
        {
            if let token::TokenType::Identifier = t {
                idents.insert(program[r[0]..r[1]].to_string(), Some(node.clone()));
            }
            new_tree.push(node.clone());
        } else {
            new_tree.push(node.clone());
        }
    }

    return Ok((new_tree, Some(idents)));
}

fn generate_braket_heirarchy(tokens: &Vec<token::Token>) -> Result<Tree, ParsingError> {
    let mut constructions: Tree = Vec::new();
    let mut cursor = 0;

    fn braket_construct(
        tokens: &Vec<token::Token>,
        cursor: &mut usize,
    ) -> Result<Node, ParsingError> {
        let origin = *cursor;
        let matching = match find_matching(&tokens, *cursor) {
            Ok(m) => m,
            Err(e) => return Err(e),
        };
        let content: Tree;
        if matching > *cursor + 1 {
            content = match generate_braket_heirarchy(&tokens[*cursor + 1..matching].to_vec()) {
                Ok(c) => c,
                Err(e) => return Err(e),
            };
        } else {
            content = Vec::new();
        }
        *cursor = matching + 1; // move cursor to after match
        match &tokens[origin].token {
            token::TokenType::Braket(token::Braket::OpenBrace) => {
                Ok(Node::Construction(Construction::Brace(Brace {
                    identifiers: None,
                    content,
                })))
            }
            token::TokenType::Braket(token::Braket::OpenSquareBraket) => {
                Ok(Node::Construction(Construction::SquareBraket(content)))
            }
            token::TokenType::Braket(token::Braket::OpenParen) => {
                Ok(Node::Construction(Construction::Paren(content)))
            }
            _s => Err(ParsingError::NonMatchable(*_s)),
        }
    }

    while cursor < tokens.len() {
        match tokens[cursor].token {
            token::TokenType::Braket(token::Braket::OpenBrace)
            | token::TokenType::Braket(token::Braket::OpenSquareBraket)
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
                constructions.push(Box::new(Node::Construction(con)));
                cursor += 1;
            }
        }
    }

    return Ok(constructions);
}

fn parse_literals(tree: &Tree, program: &String) -> Result<Tree, ParsingError> {
    let mut parsed_tree: Tree = Vec::new();

    fn literal(t: token::Token, v: token::Value, program: &String) -> Result<Node, ParsingError> {
        let lit: ast::Literal;
        lit = match v {
            token::Value::Int => ast::Literal {
                ttype: ast::Type::U64,
                value: ast::RealValue::Number(ast::Number::U64(
                    match program[t.range[0]..t.range[1]].parse() {
                        Ok(val) => val,
                        Err(_) => {
                            return Err(ParsingError::Unparsable(
                                program[t.range[0]..t.range[1]].to_string(),
                            ))
                        }
                    },
                )),
            },
            token::Value::Hex => ast::Literal {
                ttype: ast::Type::U64,
                value: ast::RealValue::Number(ast::Number::U64({
                    u64::from_str_radix(
                        program[t.range[0]..t.range[1]].trim_start_matches("0x"),
                        16,
                    )
                    .unwrap()
                })),
            },
            token::Value::Oct => ast::Literal {
                ttype: ast::Type::U64,
                value: ast::RealValue::Number(ast::Number::U64({
                    u64::from_str_radix(program[t.range[0]..t.range[1]].trim_start_matches("0o"), 8)
                        .unwrap()
                })),
            },
            token::Value::Bin => ast::Literal {
                ttype: ast::Type::U64,
                value: ast::RealValue::Number(ast::Number::U64({
                    u64::from_str_radix(program[t.range[0]..t.range[1]].trim_start_matches("0b"), 2)
                        .unwrap()
                })),
            },
            token::Value::Float => ast::Literal {
                ttype: ast::Type::F64,
                value: ast::RealValue::Number(ast::Number::F64(
                    program[t.range[0]..t.range[1]].parse().unwrap(),
                )),
            },
            token::Value::String => ast::Literal {
                ttype: ast::Type::String,
                value: ast::RealValue::String(program[t.range[0] + 1..t.range[1] - 1].to_string()),
            },
            token::Value::True | token::Value::False => ast::Literal {
                ttype: ast::Type::Bool,
                value: ast::RealValue::Bool(program[t.range[0]..t.range[1]].parse().unwrap()),
            },
            token::Value::Byte => ast::Literal {
                ttype: ast::Type::Byte,
                value: ast::RealValue::Byte(
                    program[t.range[0] + 1..t.range[1] - 1].chars().collect(),
                ),
            },
            token::Value::None => ast::Literal {
                ttype: ast::Type::None,
                value: ast::RealValue::None,
            },
            token::Value::Error => ast::Literal {
                ttype: ast::Type::Error,
                value: ast::RealValue::Error,
            },
        };
        return Ok(Node::Expression(ast::Expression::Value(
            ast::Value::Literal(lit),
        )));
    }

    for node in tree {
        parsed_tree.push(match node.as_ref() {
            Node::Construction(c) => match c {
                Construction::Brace(Brace {
                    identifiers: id,
                    content: k,
                }) => match parse_literals(&k, program) {
                    Ok(a) => Box::new(Node::Construction(Construction::Brace(Brace {
                        identifiers: id.clone(),
                        content: a,
                    }))),
                    Err(e) => return Err(e),
                },
                Construction::SquareBraket(k) => match parse_literals(&k, program) {
                    Ok(a) => Box::new(Node::Construction(Construction::SquareBraket(a))),
                    Err(e) => return Err(e),
                },
                Construction::Paren(k) => match parse_literals(&k, program) {
                    Ok(a) => Box::new(Node::Construction(Construction::Paren(a))),
                    Err(e) => return Err(e),
                },
                Construction::Token(t) => match t.token {
                    token::TokenType::Value(v) => match literal(*t, v, program) {
                        Ok(a) => Box::new(a),
                        Err(e) => return Err(e),
                    },
                    _s => Box::new(Node::Construction(Construction::Token(*t))),
                },
            },
            _s => Box::new(_s.clone()),
        });
    }

    return Ok(parsed_tree);
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
        token::Braket::OpenSquareBraket => token::Braket::CloseSquareBraket,
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

pub fn print_tree(tree: Tree, offset: usize) {
    for node in tree {
        println!("{}", (*node).display(offset));
        match *node {
            Node::Construction(Construction::Paren(k))
            | Node::Construction(Construction::SquareBraket(k)) => print_tree(k, offset + 1),
            Node::Construction(Construction::Token(t)) => {
                let mut spaces: Vec<String> = Vec::new();
                for _ in 0..offset + 1 {
                    spaces.push(String::from("    "));
                }

                let off = spaces.join("");
                println!("{}{}", off, t);
            }
            Node::Construction(Construction::Brace(ast::Brace {
                identifiers: _,
                content: c,
            })) => print_tree(c, offset + 1),
            _ => (),
        };
    }
}
