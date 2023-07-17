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
    Impl,
    Interface,
    Debug,
    And,
    Or,
    Not,
    As,
    Mod,
    Import,
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
    DoubleColon,
    Bang, // !
    Optional, // ?
    Dot,
    Comma,
    Dereference, // "*Ident"
    Address,     // "&Ident"
    Assign,
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

impl std::fmt::Display for Token {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        // Write strictly the first element into the supplied output
        // stream: `f`. Returns `fmt::Result` which indicates whether the
        // operation succeeded or failed. Note that `write!` uses syntax which
        // is very similar to `println!`.
        write!(f, "{:?} : {:?}", self.token, self.range)
    }
}
