package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

type TypeName struct {
	SelfType string `json:"self_type"`
	Name     string `json:"name"`
}

func NewTypeName(name string) TypeName {
	return TypeName{Name: name, SelfType: name}
}

func (TypeName) TypeRefType() typ.FieldRefType {
	return typ.SelfDefinedTypeType
}

func (t TypeName) TypeName() string {
	return t.Name
}

// ParseTypeName 解析类型名（现在只支持基础格式，数组格式由通用解析器处理）
func ParseTypeName(code string) gomme.Result[TypeName, string] {
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
