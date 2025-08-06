package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type BooleanType struct {
	SelfType string `json:"self_type"`
}

func NewBooleanType() BooleanType {
	return BooleanType{
		SelfType: "boolean",
	}
}

func (t BooleanType) TypeName() string {
	return "boolean"
}

func ParseBoolean(code string) gomme.Result[BooleanType, string] {
	return gomme.Map(
		gomme.Token[string]("boolean"),
		func(_ string) (BooleanType, error) { return NewBooleanType(), nil },
	)(code)
}

func (BooleanType) TypeRefType() typ.FieldRefType {
	return typ.BooleanType
}
