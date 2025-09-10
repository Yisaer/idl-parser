package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

// ArrayType 表示一个数组类型，包装了基础类型和数组大小
type ArrayType struct {
	InnerType TypeRef `json:"inner_type"`
	ArraySize int     `json:"array_size"`
}

func NewArrayType(innerType TypeRef, arraySize int) ArrayType {
	return ArrayType{
		InnerType: innerType,
		ArraySize: arraySize,
	}
}

func (t ArrayType) TypeName() string {
	return t.InnerType.TypeName()
}

func (t ArrayType) TypeRefType() typ.FieldRefType {
	return typ.ArrayType
}

// ParseArray 解析数组格式的类型，如 "short[8]", "long[16]" 等
func ParseArray(code string) gomme.Result[ArrayType, string] {
	// 首先尝试解析基础类型
	baseTypeResult := ParseBaseType(code)
	if baseTypeResult.Err != nil {
		return gomme.Result[ArrayType, string]{Err: baseTypeResult.Err}
	}

	baseType := baseTypeResult.Output
	remaining := baseTypeResult.Remaining

	// 检查是否还有数组后缀
	if len(remaining) > 0 && remaining[0] == '[' {
		// 解析数组大小
		arraySizeResult := gomme.Delimited(
			gomme.Token[string]("["),
			utils.ParseDigit,
			gomme.Token[string]("]"),
		)(remaining)

		if arraySizeResult.Err != nil {
			return gomme.Result[ArrayType, string]{Err: arraySizeResult.Err}
		}

		return gomme.Result[ArrayType, string]{
			Output:    NewArrayType(baseType, arraySizeResult.Output),
			Remaining: arraySizeResult.Remaining,
			Err:       nil,
		}
	}

	// 如果没有数组后缀，返回错误（因为这是 ParseArray 函数）
	return gomme.Result[ArrayType, string]{
		Err: gomme.NewError("expected array format"),
	}
}

// ParseBaseType 解析基础类型（不包含数组）
func ParseBaseType(code string) gomme.Result[TypeRef, string] {
	return gomme.Alternative(
		gomme.Map(ParseSequence, func(seq Sequence) (TypeRef, error) { return seq, nil }),
		gomme.Map(ParsePlainOctet, func(octet OctetType) (TypeRef, error) { return octet, nil }),
		gomme.Map(ParseShort, func(short ShortType) (TypeRef, error) { return short, nil }),
		gomme.Map(ParseUnsignedShort, func(us UnsignedShortType) (TypeRef, error) { return us, nil }),
		gomme.Map(ParseLongLong, func(longlong LongLongType) (TypeRef, error) { return longlong, nil }),
		gomme.Map(ParseLong, func(long LongType) (TypeRef, error) { return long, nil }),
		gomme.Map(ParseUnsignedLongLong, func(ull UnsignedLongLongType) (TypeRef, error) { return ull, nil }),
		gomme.Map(ParseUnsignedLong, func(ul UnsignedLongType) (TypeRef, error) { return ul, nil }),
		gomme.Map(ParseBoolean, func(b BooleanType) (TypeRef, error) { return b, nil }),
		gomme.Map(ParseFloat, func(f FloatType) (TypeRef, error) { return f, nil }),
		gomme.Map(ParseString, func(s StringType) (TypeRef, error) { return s, nil }),
		gomme.Map(ParseBitField, func(bitfield BitFieldType) (TypeRef, error) { return bitfield, nil }),
		gomme.Map(ParsePlainTypeName, func(name TypeName) (TypeRef, error) { return name, nil }),
	)(code)
}