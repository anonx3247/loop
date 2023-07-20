use crate::lexer::token::Token;
use crate::parser::construction_ast;

#[derive(Debug, Clone)]
pub enum Node {
    Expression(Expression),
    Statement(Statement),
    Construction(construction_ast::Node),
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
    pub scope: IdentMap,
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

type List = Vec<Value>;

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
    U8(u8),
    U16(u16),
    U32(u32),
    U64(u64),
    I32(i32),
    I64(i64),
    F32(f32),
    F64(f64),
    String(String),
    Byte(Vec<char>),
    Error,
    Bool(bool),
    None,
}
#[derive(Debug, Clone)]
pub struct Assign {
    pub mutable: bool,
    pub variables: Vec<Identifier>,
    pub values: Tuple,
    pub ttype: Option<Type>,
}

#[derive(Debug, Clone)]
pub struct Mutate {
    pub kind: MutationKind,
    pub variables: Vec<Identifier>,
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
    pub identifier: Identifier,
    pub arguments: Vec<ArgValue>,
}

pub type StructValue = FnCall;

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
    pub ttype: Type,
    pub default_value: Option<Literal>,
}

#[derive(Debug, Clone)]
pub struct ArgValue {
    pub name: Option<Identifier>,
    pub value: Expression,
}

#[derive(Debug, Clone)]
pub struct For {
    pub elem: Vec<Identifier>,
    pub iterator: Iterable,
    pub body: Block,
}

pub type Tuple = Vec<Value>;

#[derive(Debug, Clone)]
pub enum Iterable {
    Range(Range),
    Object(Object),
}

#[derive(Debug, Clone)]
pub struct Range {
    pub start: Value,
    pub end: Value,
    pub step: Option<Value>,
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
