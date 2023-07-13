mod lexer;

use std::fs;

fn main() {
    let program = fs::read_to_string("main.lp").expect("File wasn't found");

    let tokens = lexer::tokenize(program);

    for token in tokens {
        println!("{}", token);
    }
}
