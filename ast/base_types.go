package ast

import (
	"fmt"
	"reflect"
	"strconv"

	"com.loop.anonx3247/lexer"
	"com.loop.anonx3247/utils"
)

type BaseType int

const (
	I8 BaseType = iota
	I16
	I32
	I64
	U8
	U16
	U32
	U64
	F32
	F64
	Bool
	Str
)

func ToBaseType[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string | bool](t T) BaseType {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Int8:
		return I8
	case reflect.Int16:
		return I16
	case reflect.Int32:
		return I32
	case reflect.Int64:
		return I64
	case reflect.Uint8:
		return U8
	case reflect.Uint16:
		return U16
	case reflect.Uint32:
		return U32
	case reflect.Uint64:
		return U64
	case reflect.Float32:
		return F32
	case reflect.Float64:
		return F64
	case reflect.Bool:
		return Bool
	case reflect.String:
		return Str
	case reflect.Int:
		return I32
	case reflect.Uint:
		return U32
	}
	return -1
}

// Generic base value type that implements both Value and Type interfaces
type BaseValue[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string | bool] struct {
	source utils.String
	value  T
}

func (bv BaseValue[T]) BaseType() BaseType {
	return ToBaseType(bv.value)
}

func HaveSameType[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string | bool](values ...BaseValue[T]) bool {
	for _, value := range values {
		if value.BaseType() != values[0].BaseType() {
			return false
		}
	}
	return true
}

func GetBaseTypeValues[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string | bool](values ...Value) (out []T) {
	underlyingValues := make([]BaseValue[T], 0)
	for _, value := range values {
		val, ok := value.(BaseValue[T])
		if !ok {
			return nil
		}
		underlyingValues = append(underlyingValues, val)
	}
	for _, value := range underlyingValues {
		out = append(out, value.value)
	}
	return out
}

func (bv BaseValue[T]) IsBase() bool {
	return true
}

func (bv BaseValue[T]) Value() Value {
	return bv
}

func (bv BaseValue[T]) Source() utils.String {
	return bv.source
}

func (bv BaseValue[T]) Type() Type {
	return bv
}

func (bv BaseValue[T]) String() string {
	return fmt.Sprintf("%v", bv.value)
}

// GetValue returns the underlying value
func (bv BaseValue[T]) GetValue() T {
	return bv.value
}

// Convenient constructor functions for each type
func NewI8Value(v int8, source utils.String) BaseValue[int8] {
	return BaseValue[int8]{value: v, source: source}
}

func NewI16Value(v int16, source utils.String) BaseValue[int16] {
	return BaseValue[int16]{value: v, source: source}
}

func NewI32Value(v int32, source utils.String) BaseValue[int32] {
	return BaseValue[int32]{value: v, source: source}
}

func NewI64Value(v int64, source utils.String) BaseValue[int64] {
	return BaseValue[int64]{value: v, source: source}
}

func NewU8Value(v uint8, source utils.String) BaseValue[uint8] {
	return BaseValue[uint8]{value: v, source: source}
}

func NewU16Value(v uint16, source utils.String) BaseValue[uint16] {
	return BaseValue[uint16]{value: v, source: source}
}

func NewU32Value(v uint32, source utils.String) BaseValue[uint32] {
	return BaseValue[uint32]{value: v, source: source}
}

func NewU64Value(v uint64, source utils.String) BaseValue[uint64] {
	return BaseValue[uint64]{value: v, source: source}
}

func NewF32Value(v float32, source utils.String) BaseValue[float32] {
	return BaseValue[float32]{value: v, source: source}
}

func NewF64Value(v float64, source utils.String) BaseValue[float64] {
	return BaseValue[float64]{value: v, source: source}
}

func NewBoolValue(v bool, source utils.String) BaseValue[bool] {
	return BaseValue[bool]{value: v, source: source}
}

func NewStrValue(v string, source utils.String) BaseValue[string] {
	return BaseValue[string]{value: v, source: source}
}

func TryIntFrom[T int8 | int16 | int32 | int64 | uint16 | uint32 | uint64](tok lexer.Token) (BaseValue[T], error) {
	if tok.Type == lexer.NUMBER_LITERAL {
		val, ok := strconv.ParseInt(tok.Value.String(), 10, 64)
		if ok != nil {
			return BaseValue[T]{}, tok.Error("cannot convert token value to target type")
		}
		return BaseValue[T]{value: T(val), source: tok.Value}, nil
	}
	return BaseValue[T]{}, tok.Error("cannot convert token value to target type")
}

func TryFloatFrom[T float32 | float64](tok lexer.Token) (BaseValue[T], error) {
	if tok.Type == lexer.NUMBER_LITERAL {
		val, ok := strconv.ParseFloat(tok.Value.String(), 64)
		if ok != nil {
			return BaseValue[T]{}, tok.Error("cannot convert token value to target type")
		}
		return BaseValue[T]{value: T(val), source: tok.Value}, nil
	}
	return BaseValue[T]{}, tok.Error("cannot convert token value to target type")
}

func TryBoolFrom(tok lexer.Token) (BaseValue[bool], error) {
	if tok.Type == lexer.TRUE {
		return BaseValue[bool]{value: true, source: tok.Value}, nil
	} else if tok.Type == lexer.FALSE {
		return BaseValue[bool]{value: false, source: tok.Value}, nil
	}
	return BaseValue[bool]{}, tok.Error("cannot convert token value to target type")
}
func TryStrFrom(tok lexer.Token) (BaseValue[string], error) {
	if tok.Type == lexer.STRING_LITERAL {
		return BaseValue[string]{value: tok.Value.String()}, nil
	}
	return BaseValue[string]{}, utils.Error{Source: tok.Value, Message: "cannot convert token value to target type"}
}
