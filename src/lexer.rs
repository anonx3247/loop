enum Token {
    Comment(str),
    Identifier(str),
    Keyword(Keyword),
    Compare(Compare),
    Operator(Operator),
    Value(Value),
    Symbol(Symbol),
    Braket(Braket),
    Type(str),
}

enum Compare {
    GreaterEqual,
    LessEqual,
    Greater,
    Less,
}

enum Operator {
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

enum Keyword {
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
    Struct,
    Impl,
    Interface,
    Debug,
    And,
    Or,
    Not,
    As,
}

enum Value {
    Int(i64),
    Float(f64),
    String(str),
    Hex(str),
    Oct(str),
    Bin(str),
    True,
    False,
    None,
}

enum Symbol {
    Colon,
    Error,
    Optional,
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

enum Braket {
    OpenBrace,
    CloseBrace,
    OpenBraket,
    CloseBraket,
    OpenParen,
    CloseParen,
}
