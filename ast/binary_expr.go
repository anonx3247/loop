package ast

import (
	"com.loop.anonx3247/lexer"
	"com.loop.anonx3247/utils"
)

type BinaryExpr struct {
	source utils.String
	Op     lexer.TokenType
	Left   *Expr
	Right  *Expr
}

func (b BinaryExpr) Source() utils.String {
	return b.source
}

func (b BinaryExpr) Eval() (Value, error) {
	left, err := (*b.Left).Eval()
	if err != nil {
		return nil, err
	}
	right, err := (*b.Right).Eval()
	if err != nil {
		return nil, err
	}

	switch b.Op {
	case lexer.PLUS:
		return AddValues(left, right, b.Source())
	case lexer.MINUS:
		return SubtractValues(left, right, b.Source())
	case lexer.MULTIPLY:
		return MultiplyValues(left, right, b.Source())
	case lexer.DIVIDE:
		return DivideValues(left, right, b.Source())
	case lexer.MODULO:
		return ModuloValues(left, right, b.Source())
	case lexer.EQUAL:
		return EqualsValues(left, right, b.Source())
	case lexer.NOT_EQUAL:
		return NotEqualsValues(left, right, b.Source())
	case lexer.GREATER_THAN:
		return GreaterThanValues(left, right, b.Source())
	case lexer.GREATER_THAN_OR_EQUAL:
		return GreaterThanOrEqualValues(left, right, b.Source())
	case lexer.LESS_THAN:
		return LessThanValues(left, right, b.Source())
	case lexer.LESS_THAN_OR_EQUAL:
		return LessThanOrEqualValues(left, right, b.Source())
	case lexer.AND:
		return AndValues(left, right, b.Source())
	case lexer.OR:
		return OrValues(left, right, b.Source())
	case lexer.BITWISE_XOR:
		return BitwiseXorValues(left, right, b.Source())
	case lexer.BITWISE_AND:
		return BitwiseAndValues(left, right, b.Source())
	case lexer.BITWISE_OR:
		return BitwiseOrValues(left, right, b.Source())
	case lexer.BITWISE_LEFT_SHIFT:
		return BitwiseLeftShiftValues(left, right, b.Source())
	case lexer.BITWISE_RIGHT_SHIFT:
		return BitwiseRightShiftValues(left, right, b.Source())
	default:
		return nil, utils.Error{Source: b.Source(), Message: "unsupported binary operator"}
	}
}

func addBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left + right,
		source: source,
	}, nil
}

func subtractBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left - right,
		source: source,
	}, nil
}

func multiplyBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left * right,
		source: source,
	}, nil
}

func divideBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left / right,
		source: source,
	}, nil
}

func moduloBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left % right,
		source: source,
	}, nil
}

func equalsBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string](left, right T, source utils.String) (BaseValue[bool], error) {
	return BaseValue[bool]{
		value:  left == right,
		source: source,
	}, nil
}

func notEqualsBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string](left, right T, source utils.String) (BaseValue[bool], error) {
	return BaseValue[bool]{
		value:  left != right,
		source: source,
	}, nil
}

func greaterThanBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](left, right T, source utils.String) (BaseValue[bool], error) {
	return BaseValue[bool]{
		value:  left > right,
		source: source,
	}, nil
}

func greaterThanOrEqualBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](left, right T, source utils.String) (BaseValue[bool], error) {
	return BaseValue[bool]{
		value:  left >= right,
		source: source,
	}, nil
}

func lessThanBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](left, right T, source utils.String) (BaseValue[bool], error) {
	return BaseValue[bool]{
		value:  left < right,
		source: source,
	}, nil
}

func lessThanOrEqualBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](left, right T, source utils.String) (BaseValue[bool], error) {
	return BaseValue[bool]{
		value:  left <= right,
		source: source,
	}, nil
}

func bitwiseXorBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left ^ right,
		source: source,
	}, nil
}

func bitwiseAndBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left & right,
		source: source,
	}, nil
}

func bitwiseOrBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left | right,
		source: source,
	}, nil
}

func bitwiseLeftShiftBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left << right,
		source: source,
	}, nil
}

func bitwiseRightShiftBaseValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](left, right T, source utils.String) (BaseValue[T], error) {
	return BaseValue[T]{
		value:  left >> right,
		source: source,
	}, nil
}

func andBaseValues(left, right bool, source utils.String) (BaseValue[bool], error) {
	return BaseValue[bool]{
		value:  left && right,
		source: source,
	}, nil
}

func orBaseValues(left, right bool, source utils.String) (BaseValue[bool], error) {
	return BaseValue[bool]{
		value:  left || right,
		source: source,
	}, nil
}

// Helper function to add two values
func AddValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in addition"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	case Str:
		underlyingValues := GetBaseTypeValues[string](left, right)
		return addBaseValues(underlyingValues[0], underlyingValues[1], source)
	}

	return nil, utils.Error{Source: source, Message: "unsupported types for addition"}
}

// Helper function to multiply two values
func MultiplyValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in multiplication"}
	}

	switch left.Type().BaseType() {
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return multiplyBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for multiplication"}
}

func SubtractValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in subtraction"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return subtractBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for subtraction"}
}

func DivideValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in division"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return divideBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for division"}
}

func ModuloValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in modulo"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return moduloBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return moduloBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return moduloBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return moduloBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return moduloBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return moduloBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return moduloBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return moduloBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for modulo"}
}

func EqualsValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in equality"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case Str:
		underlyingValues := GetBaseTypeValues[string](left, right)
		return equalsBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for equality"}
}

func NotEqualsValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in inequality"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	case Str:
		underlyingValues := GetBaseTypeValues[string](left, right)
		return notEqualsBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for inequality"}
}

func GreaterThanValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in greater than"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return greaterThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for greater than"}
}

func GreaterThanOrEqualValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in greater than or equal"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return greaterThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	default:
		return nil, utils.Error{Source: source, Message: "unsupported types for greater than or equal"}
	}
}

func LessThanValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in less than"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return lessThanBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for less than"}
}

func LessThanOrEqualValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in less than or equal"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F32:
		underlyingValues := GetBaseTypeValues[float32](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	case F64:
		underlyingValues := GetBaseTypeValues[float64](left, right)
		return lessThanOrEqualBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for less than or equal"}
}

func AndValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in greater than or equal"}
	}

	switch left.Type().BaseType() {
	case Bool:
		underlyingValues := GetBaseTypeValues[bool](left, right)
		return andBaseValues(underlyingValues[0], underlyingValues[1], source)
	default:
		return nil, utils.Error{Source: source, Message: "unsupported types for greater than or equal"}
	}
}

func OrValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in greater than or equal"}
	}

	switch left.Type().BaseType() {
	case Bool:
		underlyingValues := GetBaseTypeValues[bool](left, right)
		return orBaseValues(underlyingValues[0], underlyingValues[1], source)
	default:
		return nil, utils.Error{Source: source, Message: "unsupported types for greater than or equal"}
	}
}

func BitwiseXorValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in division"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for division"}
}

func BitwiseAndValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in division"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for division"}
}

func BitwiseOrValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in division"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return bitwiseOrBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for division"}
}

func BitwiseLeftShiftValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in division"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return bitwiseLeftShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return bitwiseLeftShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return bitwiseLeftShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return bitwiseLeftShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return bitwiseLeftShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return bitwiseLeftShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return bitwiseLeftShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return bitwiseLeftShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for division"}
}

func BitwiseRightShiftValues(left, right Value, source utils.String) (Value, error) {
	if (!left.IsBase() || !right.IsBase()) || (left.Type().BaseType() != right.Type().BaseType()) {
		return nil, utils.Error{Source: source, Message: "type mismatch in division"}
	}

	switch left.Type().BaseType() {
	case I8:
		underlyingValues := GetBaseTypeValues[int8](left, right)
		return bitwiseRightShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I16:
		underlyingValues := GetBaseTypeValues[int16](left, right)
		return bitwiseRightShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I32:
		underlyingValues := GetBaseTypeValues[int32](left, right)
		return bitwiseRightShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case I64:
		underlyingValues := GetBaseTypeValues[int64](left, right)
		return bitwiseRightShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U8:
		underlyingValues := GetBaseTypeValues[uint8](left, right)
		return bitwiseRightShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U16:
		underlyingValues := GetBaseTypeValues[uint16](left, right)
		return bitwiseRightShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U32:
		underlyingValues := GetBaseTypeValues[uint32](left, right)
		return bitwiseRightShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	case U64:
		underlyingValues := GetBaseTypeValues[uint64](left, right)
		return bitwiseRightShiftBaseValues(underlyingValues[0], underlyingValues[1], source)
	}
	return nil, utils.Error{Source: source, Message: "unsupported types for division"}
}
