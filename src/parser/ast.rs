use crate::lexer::token::Token;

enum Node {
    Expression(Expression),
    Statement(Statement),
    Construction(Construction)
}

struct ProgramFile {
    imports: Vec<ImportStatement>,
    mods: Vec<ModStatement>,
    body: Vec<Statement>
}

struct BinaryOperator {
    op: BinaryOperatorKind,
    left: Box<Expression>,
    right: Box<Expression>,
}

struct UnaryOperator {
    op: BinaryOperatorKind,
    right: Box<Expression>,
}

struct Debug {
    child: Box<Statement>,
}

struct IfElseStatement {
    if_branch: IfStatement,
    elseif_branches: Vec<IfStatement>,
    else_branch: Option<Block>,
}

struct IfStatement {
    condition: Expression,
    body: Block,
}

struct Block {
    contents: Vec<Node>,
}

struct IfElse {
    if_branch: If,
    elseif_branches: Vec<If>,
    else_branch: Option<Box<Node>>,
}

struct If {
    condition: Box<Expression>,
    body: Block,
}

enum Expression {
    BinaryOperator(BinaryOperator),
    UnaryOperator(UnaryOperator),
    Value(Value),
    IfElse(IfElse),
    Switch(Switch),
    ParenExpression(ParenExpression),
}

type ParenExpression = Box<Expression>;

enum Statement {
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

struct ImportStatement {
    root: Object,
    children: Vec<Object>,
}

type ModStatement = String;

enum BinaryOperatorKind {
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

enum UnaryOperatorKind {
    Negative,
    Not,
}


struct Switch {
    item: Value,
    branches: Vec<SwitchBranch>,
    default: Option<Block>,
}
enum Value {
    FnCall(FnCall),
    Object(Object),
    Literal(Literal),
    List(List),
}

type List = Vec<Object>;

type SwitchBranch = If;

type Identifier = String;

struct Object {
    root: ObjectMember,
    child: Option<Box<Object>>,
}

enum ObjectMember {
    Identifier(Identifier),
    FnCall(FnCall),
    Index(Index),
    Namespace(String),
}

enum Index {
    Int(u64),
    Range,
}

struct Literal {
    ttype: Type,
    value: RealValue,
}

enum Type {
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

enum RealValue {
    Number(Number),
    String(String),
    Byte(u8),
}
struct Assign {
    mutable: bool,
    variables: Tuple,
    values: Tuple,
    ttype: Option<Type>,
}

struct Mutate {
    kind: MutationKind,
    variables: Tuple,
    values: Tuple,
}

enum MutationKind {
    AssignAdd,
    AssignSubtract,
    AssignMultiply,
    AssignDivide,
    AssignExponent,
    Assign,
}

struct FnCall {
    identifier: String,
    arguments: Vec<ArgValue>,
}

type StructInit = FnCall;

struct Fn {
    identifier: String,
    arguments: Vec<Arg>,
    return_types: Vec<Type>,
    body: Block,
}

struct Arg {
    name: String,
    ttype: Type,
}

struct ArgValue {
    name: Option<String>,
    value: Expression
}

struct For {
    elem: Tuple,
    iterator: Iterable,
    body: Block,
}

type Tuple = Vec<Value>;

enum Iterable {
    Range(Range),
    Object(Object),
}

struct Range {
    start: Number,
    end: Number,
    step: Number,
}

enum Number {
    U8(u8),
    U16(u16),
    U32(u32),
    U64(u64),
    I32(i32),
    I64(i64),
    F32(f32),
    F64(f64),
}

struct While {
    condition: Expression,
    body: Block,
}

struct Loop {
    body: Block,
}

struct Struct {
    name: String,
    members: Vec<Arg>,
}

struct Enum {
    name: String,
    members: Vec<EnumMember>,
}

struct EnumMember {
    name: String,
    ttype: Option<Type>,
}

enum Construction {
    Brace(Brace),
    Paren(Paren),
    Braket(Braket),
    Token(Token),
}

type Brace = Vec<Node>;
type Paren = Brace;
type Braket = Paren;
