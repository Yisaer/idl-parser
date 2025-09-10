package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
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

// ParseOctet 解析 octet 类型（现在只支持基础格式，数组格式由通用解析器处理）
func ParseOctet(code string) gomme.Result[OctetType, string] {
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
