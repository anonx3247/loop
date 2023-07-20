mod lexer;
mod parser;

use std::fs;

fn main() {
    let program = fs::read_to_string("main.lp").expect("File wasn't found");

    let tokens = match lexer::tokenize(&program) {
        Ok(tok) => tok,
        Err(e) => panic!("Got lexer error: {:?}", e),
    };

    for i in 0..tokens.len() {
        //println!("{}: {}", i, tokens[i]);
    }

    let tree = match parser::raw_parser::parse_raw(&tokens, &program) {
        Ok(tr) => tr,
        Err(e) => panic!("{:?}", e),
    };

    parser::print_tree(tree, 0);
}
