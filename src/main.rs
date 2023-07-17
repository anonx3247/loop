mod lexer;
mod parser;

use std::fs;

fn main() {
    let program = fs::read_to_string("main.lp").expect("File wasn't found");

    let tokens = match lexer::tokenize(program) {
        Ok(tok) => tok,
        Err(e) => panic!("Got lexer error: {:?}", e),
    };

    for token in tokens {
        println!("{}", token);
    }
}
