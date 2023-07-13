pub struct Token {
    token: TokenType,
    start: [u64; 2],
    end: [u64; 2],
}


pub enum TokenType {
    Comment(String),
    Identifier(String),
    Keyword(Keyword),
    Compare(Compare),
    Operator(Operator),
    Value(Value),
    Symbol(Symbol),
    Braket(Braket),
    Type(String),
}

pub enum Compare {
    GreaterEqual,
    LessEqual,
    Greater,
    Less,
    Equal,
}

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

pub enum Value {
    Int(i64),
    Float(f64),
    Stringing(String),
    Hex(String),
    Oct(String),
    Bin(String),
    True,
    False,
    None,
    Error,
}

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

pub enum Braket {
    OpenBrace,
    CloseBrace,
    OpenBraket,
    CloseBraket,
    OpenParen,
    CloseParen,
}
