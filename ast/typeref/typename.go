package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

type TypeName struct {
	SelfType  string `json:"self_type"`
	Name      string `json:"name"`
	ArraySize int    `json:"array_size,omitempty"` // 0 means not an array
}

func NewTypeName(name string) TypeName {
	return TypeName{Name: name, SelfType: name}
}

func NewTypeNameArrayType(name string, arraySize int) TypeName {
	return TypeName{Name: name, SelfType: name, ArraySize: arraySize}
}

func (TypeName) TypeRefType() typ.FieldRefType {
	return typ.SelfDefinedTypeType
}

func (t TypeName) TypeName() string {
	return t.Name
}

// ParseTypeName 解析类型名，支持数组格式
func ParseTypeName(code string) gomme.Result[TypeName, string] {
	// First, parse the identifier (type name)
	identifierResult := gomme.Recognize(
		utils.ParseValidChar,
	)(code)

	if identifierResult.Err != nil {
		return gomme.Result[TypeName, string]{Err: identifierResult.Err}
	}

	typeName := identifierResult.Output
	remaining := identifierResult.Remaining

	// Check if the remaining code starts with "["
	if len(remaining) > 0 && remaining[0] == '[' {
		// Try to parse as array format
		arrayResult := gomme.Map(
			gomme.Pair(
				gomme.Token[string](typeName),
				gomme.Delimited(
					gomme.Token[string]("["),
					utils.ParseDigit,
					gomme.Token[string]("]"),
				),
			),
			func(output gomme.PairContainer[string, int]) (TypeName, error) {
				return NewTypeNameArrayType(output.Left, output.Right), nil
			},
		)(code)

		// If array parsing succeeded, return it
		if arrayResult.Err == nil {
			return arrayResult
		}

		// If we started with "[" but parsing failed, return the error
		// This ensures that idbits[abc] will return an error instead of falling back
		return arrayResult
	}

	// Otherwise, return plain type name
	return ParsePlainTypeName(code)
}

// ParsePlainTypeName 解析纯类型名（不包含数组）
func ParsePlainTypeName(code string) gomme.Result[TypeName, string] {
	// Parse the identifier (type name)
	identifierResult := gomme.Recognize(
		utils.ParseValidChar,
	)(code)

	if identifierResult.Err != nil {
		return gomme.Result[TypeName, string]{Err: identifierResult.Err}
	}

	typeName := identifierResult.Output
	remaining := identifierResult.Remaining

	return gomme.Result[TypeName, string]{
		Output:    NewTypeName(typeName),
		Remaining: remaining,
		Err:       nil,
	}
}
