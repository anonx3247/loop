package ast

type TypeKind int

const (
	U8 TypeKind = iota
	U16
	U32
	U64
	I8
	I16
	I32
	I64
	F32
	F64
	BOOL
	CHAR
	STR
	GENERIC_TYPE
	USER_DEFINED_TYPE
	TUPLE_TYPE
)

type Type struct {
	Kind   TypeKind
	Name   string
	Params []Type
}

type Program struct {
	Expressions []Expression
}

type Expression interface {
}
