package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type StringType struct {
	SelfType string `json:"self_type"`
}

func NewStringType() StringType {
	return StringType{
		SelfType: "string",
	}
}

func (t StringType) TypeName() string {
	return "string"
}

func ParseString(code string) gomme.Result[StringType, string] {
	return gomme.Map(
		gomme.Token[string]("string"),
		func(_ string) (StringType, error) { return NewStringType(), nil },
	)(code)
}

func (StringType) TypeRefType() typ.FieldRefType {
	return typ.StringType
}
