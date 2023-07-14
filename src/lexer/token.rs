#[derive(Debug)]
pub struct Token {
    pub token: TokenType,
    pub range: Range,
}

pub type Range = [usize; 2];


#[derive(Debug)]
pub enum TokenType {
    Comment, // done
    Identifier,
    Keyword(Keyword),
    Compare(Compare),
    Operator(Operator),
    Value(Value),
    Symbol(Symbol),
    Braket(Braket),
    Type,
}

#[derive(Debug)]
pub enum Compare {
    GreaterEqual,
    LessEqual,
    Greater,
    Less,
    Equal,
}

#[derive(Debug)]
pub enum Operator {
    Plus,
    Minus,
    Mult,
    Div,
    Power,
    Modulus,
    PlusAssign,
    MinusAssign,
    MultAssign,
    DivAssign,
    PowerAssign,
}

#[derive(Debug)]
pub enum Keyword {
    Mut,
    For,
    Loop,
    While,
    If,
    Else,
    Switch,
    In,
    Return,
    Catch,
    Fn,
    Is,
    Enum,
    Stringuct,
    Impl,
    Interface,
    Debug,
    And,
    Or,
    Not,
    As,
}

#[derive(Debug)]
pub enum Value {
    Int,
    Float,
    String,
    Char,
    Hex,
    Oct,
    Bin,
    True,
    False,
    None,
    Error,
}

#[derive(Debug)]
pub enum Symbol {
    Colon,
    Bang, // !
    Optional, // ?
    Dot,
    Comma,
    Dereference, // "*Ident"
    Address,     // "&Ident"
    Assign,
    Equal,
    TypeSum, // "|"
    Arrow,
    Range,   // ".."
    Elipsis, // "..."
}

#[derive(Debug)]
pub enum Braket {
    OpenBrace,
    CloseBrace,
    OpenBraket,
    CloseBraket,
    OpenParen,
    CloseParen,
}
