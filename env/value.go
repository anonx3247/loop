package env

import "com.loop.anonx3247/utils"

type Value interface {
	Type() Type
	Source() utils.String
	IsBase() bool
	String() string
}

type Type interface {
	BaseType() BaseType
	// TupleType, FunctionType, etc.
}
