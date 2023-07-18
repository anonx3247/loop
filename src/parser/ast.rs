use crate::lexer::token::Token;

#[derive(Debug, Clone)]
pub enum Node {
    Expression(Expression),
    Statement(Statement),
    Construction(Construction),
    Undetermined(Undetermined),
    Scope(Scope),
}

#[derive(Debug, Clone)]
pub struct ProgramFile {
    pub imports: Vec<ImportStatement>,
    pub mods: Vec<ModStatement>,
    pub body: Vec<Statement>,
}

#[derive(Debug, Clone)]
pub struct BinaryOperator {
    pub op: BinaryOperatorKind,
    pub left: Box<Expression>,
    pub right: Box<Expression>,
}

#[derive(Debug, Clone)]
pub struct UnaryOperator {
    pub op: BinaryOperatorKind,
    pub right: Box<Expression>,
}

#[derive(Debug, Clone)]
pub struct Debug {
    pub child: Box<Statement>,
}

#[derive(Debug, Clone)]
pub struct IfElseStatement {
    pub if_branch: IfStatement,
    pub elseif_branches: Vec<IfStatement>,
    pub else_branch: Option<Block>,
}

#[derive(Debug, Clone)]
pub struct IfStatement {
    pub condition: Expression,
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Block {
    pub contents: Vec<Box<Node>>,
}

#[derive(Debug, Clone)]
pub struct IfElse {
    pub if_branch: If,
    pub elseif_branches: Vec<If>,
    pub else_branch: Option<Box<Node>>,
}

#[derive(Debug, Clone)]
pub struct If {
    pub condition: Box<Expression>,
    pub body: Block,
}

#[derive(Debug, Clone)]
pub enum Expression {
    BinaryOperator(BinaryOperator),
    UnaryOperator(UnaryOperator),
    Value(Value),
    IfElse(IfElse),
    Switch(Switch),
    ParenExpression(ParenExpression),
}

type ParenExpression = Box<Expression>;

#[derive(Debug, Clone)]
pub enum Statement {
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
    IfElse(IfElse),
    Switch(Switch),
    ImportStatement(ImportStatement),
    ModStatement(ModStatement),
    Break,
    Continue,
    Scope,
}

#[derive(Debug, Clone)]
pub enum Undetermined {
    FnCall(FnCall),
    Fn(Fn),
    IfElse(IfElse),
    Switch(Switch),
}

#[derive(Debug, Clone)]
pub struct Scope {
    pub identifiers: IdentMap,
    pub content: Vec<Box<Node>>,
}

pub type IdentMap = std::collections::HashMap<String, Option<Box<Node>>>;

#[derive(Debug, Clone)]
pub struct ImportStatement {
    pub root: Object,
    pub children: Vec<Object>,
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
    pub item: Value,
    pub branches: Vec<SwitchBranch>,
    pub default: Option<Block>,
}
#[derive(Debug, Clone)]
pub enum Value {
    FnCall(FnCall),
    Object(Object),
    Literal(Literal),
    List(List),
}

type List = Vec<Object>;

type SwitchBranch = If;

type Identifier = String;

#[derive(Debug, Clone)]
pub struct Object {
    pub root: ObjectMember,
    pub child: Option<Box<Object>>,
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
    pub value: RealValue,
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
pub enum RealValue {
    Number(Number),
    String(String),
    Bool(bool),
    Byte(Vec<char>),
    Error,
    None,
}
#[derive(Debug, Clone)]
pub struct Assign {
    pub mutable: bool,
    pub variables: Tuple,
    pub values: Tuple,
    pub ttype: Option<Type>,
}

#[derive(Debug, Clone)]
pub struct Mutate {
    pub kind: MutationKind,
    pub variables: Tuple,
    pub values: Tuple,
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
    pub arguments: Vec<ArgValue>,
}

type StructInit = FnCall;

#[derive(Debug, Clone)]
pub struct Fn {
    pub identifier: String,
    pub arguments: Vec<Arg>,
    pub return_types: Vec<Type>,
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Arg {
    pub name: String,
    pub ttype: Type,
}

#[derive(Debug, Clone)]
pub struct ArgValue {
    pub name: Option<String>,
    pub value: Expression,
}

#[derive(Debug, Clone)]
pub struct For {
    pub elem: Tuple,
    pub iterator: Iterable,
    pub body: Block,
}

type Tuple = Vec<Value>;

#[derive(Debug, Clone)]
pub enum Iterable {
    Range(Range),
    Object(Object),
}

#[derive(Debug, Clone)]
pub struct Range {
    pub start: Number,
    pub end: Number,
    pub step: Number,
}

#[derive(Debug, Clone)]
pub enum Number {
    U8(u8),
    U16(u16),
    U32(u32),
    U64(u64),
    I32(i32),
    I64(i64),
    F32(f32),
    F64(f64),
}

#[derive(Debug, Clone)]
pub struct While {
    pub condition: Expression,
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Loop {
    pub body: Block,
}

#[derive(Debug, Clone)]
pub struct Struct {
    pub name: String,
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

#[derive(Debug, Clone)]
pub enum Construction {
    Brace(Brace),
    Paren(Paren),
    SquareBraket(SquareBraket),
    Token(Token),
}

type Brace = Vec<Box<Node>>;
type Paren = Brace;
type SquareBraket = Paren;

pub trait NodeDisplay {
    fn display(&self, offset: usize) -> String;

    fn offset(&self, offset: usize) -> String {
        let mut spaces: Vec<String> = Vec::new();
        for i in 0..offset {
            spaces.push(String::from("    "));
        }

        spaces.join("")
    }
}

impl std::fmt::Debug for dyn NodeDisplay {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "{}", &self.display(0))
    }
}

impl NodeDisplay for Node {
    fn display(&self, offset: usize) -> String {
        return String::from(format!("{} {:?}:", self.offset(offset), self));
    }
}
