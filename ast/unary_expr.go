package ast

import (
	"com.loop.anonx3247/env"
	"com.loop.anonx3247/lexer"
	"com.loop.anonx3247/utils"
)

type UnaryExpr struct {
	source utils.String
	Op     lexer.TokenType
	Value  Expr
}

func (u UnaryExpr) Source() utils.String {
	return u.source
}

func (u UnaryExpr) Eval(env *env.Env) (env.Value, error) {
	val, err := u.Value.Eval(env)
	if err != nil {
		return nil, err
	}
	switch u.Op {
	case lexer.ADDRESS_OF:
		// TODO: implement address_of operation
		panic("address_of not implemented yet")
	case lexer.BITWISE_NOT:
		return BitwiseNot(val, u.Source())
	case lexer.NOT:
		return Not(val, u.Source())
	case lexer.MINUS:
		return Minus(val, u.Source())
	}
	return nil, utils.Error{Source: u.Source(), Message: "unsupported unary operator BITWISE_NOT on non-integer type"}

}

func AddressOf(val env.Value, source utils.String) (env.Value, error) {
	panic("address_of not implemented yet")
}

func BitwiseNot(val env.Value, source utils.String) (env.Value, error) {
	switch val.Type().BaseType() {
	case env.I32:
		if intVal, ok := val.(env.BaseValue[int32]); ok {
			return env.NewI32Value(^intVal.GetValue(), source), nil
		}
	case env.I64:
		if intVal, ok := val.(env.BaseValue[int64]); ok {
			return env.NewI64Value(^intVal.GetValue(), source), nil
		}
	case env.I16:
		if intVal, ok := val.(env.BaseValue[int16]); ok {
			return env.NewI16Value(^intVal.GetValue(), source), nil
		}
	case env.I8:
		if intVal, ok := val.(env.BaseValue[int8]); ok {
			return env.NewI8Value(^intVal.GetValue(), source), nil
		}
	case env.U32:
		if intVal, ok := val.(env.BaseValue[uint32]); ok {
			return env.NewU32Value(^intVal.GetValue(), source), nil
		}
	case env.U64:
		if intVal, ok := val.(env.BaseValue[uint64]); ok {
			return env.NewU64Value(^intVal.GetValue(), source), nil
		}
	case env.U16:
		if intVal, ok := val.(env.BaseValue[uint16]); ok {
			return env.NewU16Value(^intVal.GetValue(), source), nil
		}
	case env.U8:
		if intVal, ok := val.(env.BaseValue[uint8]); ok {
			return env.NewU8Value(^intVal.GetValue(), source), nil
		}
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for bitwise not"}
}

func Not(val env.Value, source utils.String) (env.Value, error) {
	if val.Type().BaseType() == env.Bool {
		if boolVal, ok := val.(env.BaseValue[bool]); ok {
			return env.NewBoolValue(!boolVal.GetValue(), source), nil
		}
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for not"}
}

func Minus(val env.Value, source utils.String) (env.Value, error) {
	switch val.Type().BaseType() {
	case env.I8:
		if intVal, ok := val.(env.BaseValue[int8]); ok {
			return env.NewI8Value(-intVal.GetValue(), source), nil
		}
	case env.I16:
		if intVal, ok := val.(env.BaseValue[int16]); ok {
			return env.NewI16Value(-intVal.GetValue(), source), nil
		}
	case env.I32:
		if intVal, ok := val.(env.BaseValue[int32]); ok {
			return env.NewI32Value(-intVal.GetValue(), source), nil
		}
	case env.I64:
		if intVal, ok := val.(env.BaseValue[int64]); ok {
			return env.NewI64Value(-intVal.GetValue(), source), nil
		}
	case env.F32:
		if floatVal, ok := val.(env.BaseValue[float32]); ok {
			return env.NewF32Value(-floatVal.GetValue(), source), nil
		}
	case env.F64:
		if floatVal, ok := val.(env.BaseValue[float64]); ok {
			return env.NewF64Value(-floatVal.GetValue(), source), nil
		}
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for minus"}
}
