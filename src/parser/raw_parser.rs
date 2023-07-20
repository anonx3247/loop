use std::collections::HashMap;

use crate::lexer::token::{self, Keyword, Token, TokenType};

use construction_ast::{self, Tree};

#[derive(Debug)]
pub enum ParsingError {
    NonMatchable(token::TokenType),
    MatchNotFound((usize, token::Braket)),
    Unparsable(String),
    MissingColonBeforeTypeInAssignment,
    CannotAssignTo(Node),
    TooManyCommas,
}

pub fn parse_raw(tokens: &Vec<token::Token>, program: &String) -> Result<Tree, ParsingError> {
    let mut tree: Tree = match generate_braket_heirarchy(tokens) {
        Ok(c) => c.iter().map(|k| k.clone() as Box<Node>).collect(),
        Err(e) => return Err(e),
    };

    tree = match parse_literals(&tree, program) {
        Ok(t) => t,
        Err(e) => return Err(e),
    };

    tree = match generate_scopes_from_braces(&tree, program) {
        Ok((tr, id)) => vec![Box::new(Node::Construction(Construction::Brace(Brace {
            identifiers: id,
            content: tr,
        })))],
        Err(e) => return Err(e),
    };

    tree = match parse_construction_lists(&tree) {
        Ok(t) => t,
        Err(e) => return Err(e),
    };

    return Ok(tree);
}

fn parse_construction_lists(tree: &Tree) -> Result<Tree, ParsingError> {
    let mut new_tree = Vec::new();

    // will return an empty list if there are no commas
    fn parse_list_construction(tree: &Tree) -> Result<ConstructionList, ParsingError> {
        let mut list = Vec::new();
        let mut commas = Vec::new();
        for i in 0..tree.len() {
            match tree[i].as_ref() {
                Node::Construction(Construction::Token(Token {
                    token: TokenType::Symbol(token::Symbol::Comma),
                    range: _,
                })) => {
                    if i != tree.len() - 1 {
                        commas.push(i)
                    } else {
                        return Err(ParsingError::TooManyCommas);
                    }
                }
                _ => continue,
            };
        }

        if commas.len() == 0 {
            return Ok(list);
        }
        commas.push(tree.len() - 1);

        for i in 0..commas.len() {
            list.push(tree[commas[i]..commas[i + 1]].to_vec());
        }

        return Ok(list);
    }

    for elem in tree {
        match elem.as_ref() {
            Node::Construction(Construction::Brace(Brace {
                identifiers: id,
                content: c,
            })) => {
                let list = match parse_list_construction(c) {
                    Ok(v) => v,
                    Err(e) => return Err(e),
                };
                new_tree.push(Box::new(Node::Construction(Construction::Struct(
                    ConstructionStruct {
                        identifiers: *id,
                        content: list,
                    },
                ))));
            }
            Node::Construction(Construction::Paren(p)) => {
                let list = match parse_list_construction(p) {
                    Ok(v) => v,
                    Err(e) => return Err(e),
                };
                new_tree.push(Box::new(Node::Construction(Construction::Tuple(list))));
            }
            _s => new_tree.push(*elem),
        };
    }

    Ok(new_tree)
}

fn parse_value(tree: &Tree, program: &String) -> Result<ast::Value, ParsingError> {
    Err(ParsingError::Unparsable(String::from("Value")))
}

// separates along commas into a list of tokens

/*
fn parse_assignments(tree: &Tree, program: &String) -> Result<Tree, ParsingError> {
    let mut new_tree: Tree = Vec::new();
    let mut i = 1; // cannot begin program with an '='

    fn parse_val(node: &Node, program: &String) -> Result<ast::Tuple, ParsingError> {
        match node {
            Node::Construction(Construction::Tuple(list)) => {
                let vals = Vec::new();
                for elem in list {
                    vals.push(match parse_value(&elem, program) {
                        Ok(v) => v,
                        Err(e) => return Err(e),
                    });
                }
                Ok(vals)
            }
            Node::Construction(Construction::Token(Token {
                token: TokenType::Value(v),
                range: r,
            })) => match parse_value(vec![node], program) {
                Ok(v) => v,
                Err(e) => return Err(e),
            },
            _ => return Err(ParsingError::CannotAssignTo(node.clone())),
        }
    }

    fn parse_var(node: &Node, program: &String) -> Result<ast::Tuple, ParsingError> {
        match node {
            Node::Construction(Construction::Tuple(list)) => {
                let mut vars = Vec::new();
                for item in list {
                    for elem in item {
                        match elem.as_ref() {
                            Node::Construction(Construction::Token(Token {
                                token: TokenType::Identifier,
                                range: r,
                            })) => vars.push(Value::Object(ast::Object {
                                root: ast::ObjectMember::Identifier(
                                    program[r[0]..r[1]].to_string(),
                                ),
                                child: None,
                            })),
                            _ => return Err(ParsingError::CannotAssignTo(*elem.clone())),
                        }
                    }
                }
                Ok(vars)
            }
            Node::Construction(Construction::Token(Token {
                token: TokenType::Identifier,
                range: r,
            })) => Ok(vec![Value::Object(ast::Object {
                root: ast::ObjectMember::Identifier(program[r[0]..r[1]].to_string()),
                child: None,
            })]),
            _ => return Err(ParsingError::CannotAssignTo(node.clone())),
        }
    }

    while i < tree.len() {
        if let Node::Construction(Construction::Token(Token {
            token: TokenType::Symbol(token::Symbol::Assign),
            range: _,
        })) = *tree[i]
        {
            let val = match parse_val(tree[i + 1].as_ref()) {
                Ok(v) => v,
                Err(e) => return Err(e),
            }; // ... = val
            let mut mutable = true;
            let mut ttype: Option<ast::Type> = None;
            let mut variables: ast::Tuple = Vec::new();
            let mut values: ast::Tuple = Vec::new();

            if i > 2 {
                // space for 'ident : Type'
                if let Node::Construction(Construction::Token(Token {
                    token: TokenType::Type,
                    range: r,
                })) = *tree[i - 1]
                {
                    ttype = Some(parse_type(r, program));
                    match *tree[i - 2] {
                        // check there is a ':' before the type
                        Node::Construction(Construction::Token(Token {
                            token: TokenType::Symbol(token::Symbol::Colon),
                            range: _,
                        })) => {}
                        _ => return Err(ParsingError::MissingColonBeforeTypeInAssignment),
                    }

                    variables = match parse_var(tree[i - 3].as_ref(), program) {
                        Ok(v) => v,
                        Err(e) => return Err(e),
                    };

                    if i > 3 {
                        if let Node::Construction(Construction::Token(Token {
                            token: TokenType::Keyword(Keyword::Mut),
                            range: _,
                        })) = *tree[i - 4]
                        {
                            mutable = true;
                            new_tree.remove(i - 4); // remove the 'mut' keyword
                        }
                    }
                    new_tree.remove(i - 3); // remove the colon
                    new_tree.remove(i - 2); // remove the type
                }
            }
            new_tree.remove(i - 1); // remove the vars
            new_tree.push(Box::new(Node::Statement(Statement::Assign(ast::Assign {
                mutable,
                variables,
                values,
                ttype,
            }))));

            i += 1; // skip the value now
        } else {
            new_tree.push(tree[i]);
        }
        i += 1;
    }

    return Ok(new_tree);
}
*/

fn parse_function_calls(tree: &Tree, program: &String) -> Result<Tree, ParsingError> {
    let mut i = 0;
    let mut new_tree: Tree = Vec::new();

    // to be defined
    fn extract_args(
        t: ConstructionTuple,
        program: &String,
    ) -> Result<Vec<ast::ArgValue>, ParsingError> {
        let mut args = Vec::new();
        for arg in t {
            let name: Option<String>;
            if let Node::Construction(Construction::Token(Token {
                token: TokenType::Identifier,
                range: r,
            })) = arg[0]
            {
                name = Some(program[r[0]..r[1]].to_string());
            }
        }
        return Ok(args);
    }

    while i < tree.len() - 1 {
        if let (
            Node::Construction(Construction::Token(token::Token {
                token: token::TokenType::Identifier,
                range: r,
            })),
            Node::Construction(Construction::Tuple(p)),
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

fn generate_scopes_from_braces(
    tree: &Tree,
    program: &String,
) -> Result<(Tree, ast::IdentMap), ParsingError> {
    let mut new_tree: Tree = Vec::new();
    let mut idents: ast::IdentMap = HashMap::new();
    for node in tree {
        if let Node::Construction(Construction::Brace(t)) = node.as_ref() {
            new_tree.push(match generate_scopes_from_braces(&t.content, &program) {
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

    return Ok((new_tree, idents));
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
                    identifiers: HashMap::new(),
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
        return Ok(Node::Expression(ast::Expression::Value(Value::Literal(
            lit,
        ))));
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
fn parse_type(r: token::Range, program: &String) -> ast::Type {
    match program[r[0]..r[1]].as_ref() {
        "u8" => ast::Type::U8,
        "u16" => ast::Type::U16,
        "u32" => ast::Type::U32,
        "u64" => ast::Type::U64,
        "i32" => ast::Type::I32,
        "i64" => ast::Type::I64,
        "f32" => ast::Type::F32,
        "f64" => ast::Type::F64,
        "str" => ast::Type::String,
        _s => ast::Type::UserType(String::from(_s)),
    }
}
