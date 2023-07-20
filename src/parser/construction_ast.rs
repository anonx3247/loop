use crate::lexer::token::Token;

pub type Tree = Vec<Box<Node>>;

#[derive(Debug, Clone)]
pub enum Node {
    FnCall(FnCall),
    For(For),
    While(While),
    Loop(Loop),
    Debug(Debug),
    Assign(Assign),
    Mutate(Mutate),
    Fn(Fn),
    Struct(Struct),
    Enum(Enum),
    ImportStatement(ImportStatement),
    ModStatement(ModStatement),
    Break,
    Continue,
    BinaryOperator(BinaryOperator),
    UnaryOperator(UnaryOperator),
    IfElse(IfElse),
    Switch(Switch),
    ParenExpression(ParenExpression),
    Object(Object),
    Literal(Literal),
    Many(Tree),
}

#[derive(Debug, Clone)]
pub struct BinaryOperator {
    pub op: BinaryOperatorKind,
    pub left: Box<Node>,
    pub right: Box<Node>,
}

#[derive(Debug, Clone)]
pub struct UnaryOperator {
    pub op: BinaryOperatorKind,
    pub right: Box<Node>,
}

#[derive(Debug, Clone)]
pub struct Debug {
    pub child: Box<Node>,
}

#[derive(Debug, Clone)]
pub struct IfElseStatement {
    pub if_branch: IfStatement,
    pub elseif_branches: Vec<IfStatement>,
    pub else_branch: Option<Block>,
}

#[derive(Debug, Clone)]
pub struct IfStatement {
    pub condition: Box<Node>,
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Block {
    pub identifiers: IdentMap,
    pub contents: Box<Node>,
}

#[derive(Debug, Clone)]
pub struct IfElse {
    pub if_branch: If,
    pub elseif_branches: Vec<If>,
    pub else_branch: Option<Box<Node>>,
}

#[derive(Debug, Clone)]
pub struct If {
    pub condition: Box<Node>,
    pub body: Block,
}

type ParenExpression = Box<Node>;

pub type IdentMap = std::collections::HashMap<String, Option<Box<Node>>>;

#[derive(Debug, Clone)]
pub struct ImportStatement {
    pub root: Object,
    pub children: Box<Node>,
}

type ModStatement = String;

#[derive(Debug, Clone)]
pub enum BinaryOperatorKind {
    Sum,
    Multiplication,
    Subtraction,
    Division,
    Exponent,
    And,
    Or,
    Greater,
    Less,
    Equal,
    GreaterEqual,
    LessEqual,
}

#[derive(Debug, Clone)]
pub enum UnaryOperatorKind {
    Negative,
    Not,
}

#[derive(Debug, Clone)]
pub struct Switch {
    pub item: Box<Node>,
    pub branches: Vec<SwitchBranch>,
    pub default: Option<Block>,
}

type SwitchBranch = If;

type Identifier = String;

#[derive(Debug, Clone)]
pub struct Object {
    pub root: ObjectMember,
    pub child: Option<Box<Node>>,
}

#[derive(Debug, Clone)]
pub enum ObjectMember {
    Identifier(Identifier),
    FnCall(FnCall),
    Index(Index),
    Namespace(String),
}

#[derive(Debug, Clone)]
pub enum Index {
    Int(u64),
    Range,
}

#[derive(Debug, Clone)]
pub struct Literal {
    pub ttype: Type,
    pub value: String,
}

#[derive(Debug, Clone)]
pub enum Type {
    U8,
    U16,
    U32,
    U64,
    I32,
    I64,
    F32,
    F64,
    String,
    Byte,
    UserType(String),
    Bool,
    Error,
    None,
}

#[derive(Debug, Clone)]
pub struct Assign {
    pub mutable: bool,
    pub variables: Box<Node>,
    pub values: Box<Node>,
    pub ttype: Option<Type>,
}

#[derive(Debug, Clone)]
pub struct Mutate {
    pub kind: MutationKind,
    pub variables: Box<Node>,
    pub values: Box<Node>,
}

#[derive(Debug, Clone)]
pub enum MutationKind {
    AssignAdd,
    AssignSubtract,
    AssignMultiply,
    AssignDivide,
    AssignExponent,
    Assign,
}

#[derive(Debug, Clone)]
pub struct FnCall {
    pub identifier: String,
    pub arguments: Vec<Arg>,
}

type StructInit = FnCall;

#[derive(Debug, Clone)]
pub struct Fn {
    pub identifier: Identifier,
    pub arguments: Vec<Arg>,
    pub return_types: Vec<Type>,
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Arg {
    pub name: Identifier,
    pub ttype: Option<Type>,
    pub value: Option<Box<Node>>,
}

#[derive(Debug, Clone)]
pub struct For {
    pub elem: Box<Node>,
    pub iterator: Box<Node>,
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Range {
    pub start: Box<Node>,
    pub end: Box<Node>,
    pub step: Box<Node>,
}

#[derive(Debug, Clone)]
pub struct While {
    pub condition: Box<Node>,
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Loop {
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Struct {
    pub name: Identifier,
    pub members: Vec<Arg>,
}

#[derive(Debug, Clone)]
pub struct Enum {
    pub name: String,
    pub members: Vec<EnumMember>,
}

#[derive(Debug, Clone)]
pub struct EnumMember {
    pub name: String,
    pub ttype: Option<Type>,
}

#[derive(Clone)]
pub struct Brace {
    pub content: Vec<Box<Node>>,
}

type Paren = Vec<Box<Node>>;
type SquareBraket = Paren;

pub type List = Vec<Box<Node>>;

impl Node {
    pub fn display(&self, offset: usize) -> String {
        return String::from(format!("{} {:?}:", self.offset(offset), self));
    }

    fn offset(&self, offset: usize) -> String {
        let mut spaces: Vec<String> = Vec::new();
        for i in 0..offset {
            spaces.push(String::from("    "));
        }

        spaces.join("")
    }
}

impl std::fmt::Debug for Brace {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        // Write strictly the first element into the supplied output
        // stream: `f`. Returns `fmt::Result` which indicates whether the
        // operation succeeded or failed. Note that `write!` uses syntax which
        // is very similar to `println!`.
        write!(f, "Brace: id=*, content:")
    }
}
