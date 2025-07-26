package ast

import (
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

func (u UnaryExpr) Eval() (Value, error) {
	val, err := u.Value.Eval()
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

func AddressOf(val Value, source utils.String) (Value, error) {
	panic("address_of not implemented yet")
}

func BitwiseNot(val Value, source utils.String) (Value, error) {
	switch val.Type().BaseType() {
	case I32:
		if intVal, ok := val.(BaseValue[int32]); ok {
			return NewI32Value(^intVal.GetValue(), source), nil
		}
	case I64:
		if intVal, ok := val.(BaseValue[int64]); ok {
			return NewI64Value(^intVal.GetValue(), source), nil
		}
	case I16:
		if intVal, ok := val.(BaseValue[int16]); ok {
			return NewI16Value(^intVal.GetValue(), source), nil
		}
	case I8:
		if intVal, ok := val.(BaseValue[int8]); ok {
			return NewI8Value(^intVal.GetValue(), source), nil
		}
	case U32:
		if intVal, ok := val.(BaseValue[uint32]); ok {
			return NewU32Value(^intVal.GetValue(), source), nil
		}
	case U64:
		if intVal, ok := val.(BaseValue[uint64]); ok {
			return NewU64Value(^intVal.GetValue(), source), nil
		}
	case U16:
		if intVal, ok := val.(BaseValue[uint16]); ok {
			return NewU16Value(^intVal.GetValue(), source), nil
		}
	case U8:
		if intVal, ok := val.(BaseValue[uint8]); ok {
			return NewU8Value(^intVal.GetValue(), source), nil
		}
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for bitwise not"}
}

func Not(val Value, source utils.String) (Value, error) {
	if val.Type().BaseType() == Bool {
		if boolVal, ok := val.(BaseValue[bool]); ok {
			return NewBoolValue(!boolVal.GetValue(), source), nil
		}
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for not"}
}

func Minus(val Value, source utils.String) (Value, error) {
	switch val.Type().BaseType() {
	case I8:
		if intVal, ok := val.(BaseValue[int8]); ok {
			return NewI8Value(-intVal.GetValue(), source), nil
		}
	case I16:
		if intVal, ok := val.(BaseValue[int16]); ok {
			return NewI16Value(-intVal.GetValue(), source), nil
		}
	case I32:
		if intVal, ok := val.(BaseValue[int32]); ok {
			return NewI32Value(-intVal.GetValue(), source), nil
		}
	case I64:
		if intVal, ok := val.(BaseValue[int64]); ok {
			return NewI64Value(-intVal.GetValue(), source), nil
		}
	case F32:
		if floatVal, ok := val.(BaseValue[float32]); ok {
			return NewF32Value(-floatVal.GetValue(), source), nil
		}
	case F64:
		if floatVal, ok := val.(BaseValue[float64]); ok {
			return NewF64Value(-floatVal.GetValue(), source), nil
		}
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for minus"}
}
