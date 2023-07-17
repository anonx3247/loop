use crate::lexer::token::Token;

#[derive(Debug)]
pub enum NodeKind {
    Expression,
    Statement,
    Construction(ConstructionKind),
}

#[derive(Debug)]
pub enum ConstructionKind {
    Brace,
    Braket,
    Paren,
    Token,
}

#[derive(Debug)]
pub struct ProgramFile {
    pub imports: Vec<ImportStatement>,
    pub mods: Vec<ModStatement>,
    pub body: Vec<Statement>,
}

#[derive(Debug)]
pub struct BinaryOperator {
    pub op: BinaryOperatorKind,
    pub left: Box<Expression>,
    pub right: Box<Expression>,
}

#[derive(Debug)]
pub struct UnaryOperator {
    pub op: BinaryOperatorKind,
    pub right: Box<Expression>,
}

#[derive(Debug)]
pub struct Debug {
    pub child: Box<Statement>,
}

#[derive(Debug)]
pub struct IfElseStatement {
    pub if_branch: IfStatement,
    pub elseif_branches: Vec<IfStatement>,
    pub else_branch: Option<Block>,
}

#[derive(Debug)]
pub struct IfStatement {
    pub condition: Expression,
    pub body: Block,
}

#[derive(Debug)]
pub struct Block {
    pub contents: Vec<Box<dyn Node>>,
}

#[derive(Debug)]
pub struct IfElse {
    pub if_branch: If,
    pub elseif_branches: Vec<If>,
    pub else_branch: Option<Box<dyn Node>>,
}

#[derive(Debug)]
pub struct If {
    pub condition: Box<Expression>,
    pub body: Block,
}

#[derive(Debug)]
pub enum Expression {
    BinaryOperator(BinaryOperator),
    UnaryOperator(UnaryOperator),
    Value(Value),
    IfElse(IfElse),
    Switch(Switch),
    ParenExpression(ParenExpression),
}

type ParenExpression = Box<Expression>;

#[derive(Debug)]
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

#[derive(Debug)]
pub struct ImportStatement {
    pub root: Object,
    pub children: Vec<Object>,
}

type ModStatement = String;

#[derive(Debug)]
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

#[derive(Debug)]
pub enum UnaryOperatorKind {
    Negative,
    Not,
}

#[derive(Debug)]
pub struct Switch {
    pub item: Value,
    pub branches: Vec<SwitchBranch>,
    pub default: Option<Block>,
}
#[derive(Debug)]
pub enum Value {
    FnCall(FnCall),
    Object(Object),
    Literal(Literal),
    List(List),
}

type List = Vec<Object>;

type SwitchBranch = If;

type Identifier = String;

#[derive(Debug)]
pub struct Object {
    pub root: ObjectMember,
    pub child: Option<Box<Object>>,
}

#[derive(Debug)]
pub enum ObjectMember {
    Identifier(Identifier),
    FnCall(FnCall),
    Index(Index),
    Namespace(String),
}

#[derive(Debug)]
pub enum Index {
    Int(u64),
    Range,
}

#[derive(Debug)]
pub struct Literal {
    pub ttype: Type,
    pub value: RealValue,
}

#[derive(Debug)]
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
}

#[derive(Debug)]
pub enum RealValue {
    Number(Number),
    String(String),
    Byte(u8),
}
#[derive(Debug)]
pub struct Assign {
    pub mutable: bool,
    pub variables: Tuple,
    pub values: Tuple,
    pub ttype: Option<Type>,
}

#[derive(Debug)]
pub struct Mutate {
    pub kind: MutationKind,
    pub variables: Tuple,
    pub values: Tuple,
}

#[derive(Debug)]
pub enum MutationKind {
    AssignAdd,
    AssignSubtract,
    AssignMultiply,
    AssignDivide,
    AssignExponent,
    Assign,
}

#[derive(Debug)]
pub struct FnCall {
    pub identifier: String,
    pub arguments: Vec<ArgValue>,
}

type StructInit = FnCall;

#[derive(Debug)]
pub struct Fn {
    pub identifier: String,
    pub arguments: Vec<Arg>,
    pub return_types: Vec<Type>,
    pub body: Block,
}

#[derive(Debug)]
pub struct Arg {
    pub name: String,
    pub ttype: Type,
}

#[derive(Debug)]
pub struct ArgValue {
    pub name: Option<String>,
    pub value: Expression,
}

#[derive(Debug)]
pub struct For {
    pub elem: Tuple,
    pub iterator: Iterable,
    pub body: Block,
}

type Tuple = Vec<Value>;

#[derive(Debug)]
pub enum Iterable {
    Range(Range),
    Object(Object),
}

#[derive(Debug)]
pub struct Range {
    pub start: Number,
    pub end: Number,
    pub step: Number,
}

#[derive(Debug)]
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

#[derive(Debug)]
pub struct While {
    pub condition: Expression,
    pub body: Block,
}

#[derive(Debug)]
pub struct Loop {
    pub body: Block,
}

#[derive(Debug)]
pub struct Struct {
    pub name: String,
    pub members: Vec<Arg>,
}

#[derive(Debug)]
pub struct Enum {
    pub name: String,
    pub members: Vec<EnumMember>,
}

#[derive(Debug)]
pub struct EnumMember {
    pub name: String,
    pub ttype: Option<Type>,
}

#[derive(Debug)]
pub enum Construction {
    Brace(Brace),
    Paren(Paren),
    Braket(Braket),
    Token(Token),
}

type Brace = Vec<Box<Construction>>;
type Paren = Brace;
type Braket = Paren;

pub trait Node {
    fn display(&self, offset: usize) -> String {
        let mut spaces: Vec<String> = Vec::new();
        for i in 0..offset {
            spaces.push(String::from("    "));
        }

        let off = spaces.join("");
        return String::from(format!("{} {:?}:", off, self.kind()));
    }
    fn kind(&self) -> NodeKind;
}

impl Node for Statement {
    fn kind(&self) -> NodeKind {
        NodeKind::Statement
    }
}

impl std::fmt::Debug for dyn Node {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "{}", &self.display(0))
    }
}

impl Node for Expression {
    fn kind(&self) -> NodeKind {
        NodeKind::Expression
    }
}

impl Node for Construction {
    fn kind(&self) -> NodeKind {
        match self {
            Self::Brace(_) => NodeKind::Construction(ConstructionKind::Brace),
            Self::Paren(_) => NodeKind::Construction(ConstructionKind::Paren),
            Self::Braket(_) => NodeKind::Construction(ConstructionKind::Braket),
            Self::Token(_) => NodeKind::Construction(ConstructionKind::Token),
        }
    }
}
