package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

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

// ParseArray
func ParseArray(code string) gomme.Result[ArrayType, string] {
	baseTypeResult := ParseBaseType(code)
	if baseTypeResult.Err != nil {
		return gomme.Result[ArrayType, string]{Err: baseTypeResult.Err}
	}

	baseType := baseTypeResult.Output
	remaining := baseTypeResult.Remaining

	if len(remaining) > 0 && remaining[0] == '[' {
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

	return gomme.Result[ArrayType, string]{
		Err: gomme.NewError("expected array format"),
	}
}

// ParseBaseType
func ParseBaseType(code string) gomme.Result[TypeRef, string] {
	return gomme.Alternative(
		gomme.Map(ParseSequence, func(seq Sequence) (TypeRef, error) { return seq, nil }),
		gomme.Map(ParseOctet, func(octet OctetType) (TypeRef, error) { return octet, nil }),
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
		gomme.Map(ParseTypeName, func(name TypeName) (TypeRef, error) { return name, nil }),
	)(code)
}
