package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

type OctetType struct {
	SelfType  string `json:"self_type"`
	ArraySize int    `json:"array_size,omitempty"` // 0 means not an array
}

func NewOctetType() OctetType {
	return OctetType{SelfType: "octet"}
}

func NewOctetArrayType(arraySize int) OctetType {
	return OctetType{SelfType: "octet", ArraySize: arraySize}
}

func (t OctetType) TypeName() string {
	return t.SelfType
}

// ParseOctet 解析 octet 类型，支持数组格式
func ParseOctet(code string) gomme.Result[OctetType, string] {
	// Check if input starts with "octet["
	if len(code) >= 6 && code[:6] == "octet[" {
		// Try to parse as array format
		arrayResult := gomme.Map(
			gomme.Pair(
				gomme.Token[string]("octet"),
				gomme.Delimited(
					gomme.Token[string]("["),
					utils.ParseDigit,
					gomme.Token[string]("]"),
				),
			),
			func(output gomme.PairContainer[string, int]) (OctetType, error) {
				return NewOctetArrayType(output.Right), nil
			},
		)(code)

		// If array parsing succeeded and consumed all input, return it
		if arrayResult.Err == nil && arrayResult.Remaining == "" {
			return arrayResult
		}

		// If we started with "octet[" but parsing failed, return the error
		return arrayResult
	}

	// Otherwise, try plain octet format
	return ParsePlainOctet(code)
}

// ParsePlainOctet 解析纯 octet 类型（不包含数组）
func ParsePlainOctet(code string) gomme.Result[OctetType, string] {
	return gomme.Map(
		gomme.Token[string]("octet"),
		func(token string) (OctetType, error) {
			return NewOctetType(), nil
		},
	)(code)
}

func (t OctetType) TypeRefType() typ.FieldRefType {
	if t.ArraySize > 0 {
		return typ.OctetArrayType
	}
	return typ.OctetType
}
