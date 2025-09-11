package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type DoubleType struct {
	SelfType string `json:"self_type"`
}

func NewDoubleType() DoubleType {
	return DoubleType{
		SelfType: "double",
	}
}

func (t DoubleType) TypeName() string {
	return "double"
}

func ParseDouble(code string) gomme.Result[DoubleType, string] {
	return gomme.Map(
		gomme.Token[string]("double"),
		func(_ string) (DoubleType, error) { return NewDoubleType(), nil },
	)(code)
}

func (DoubleType) TypeRefType() typ.FieldRefType {
	return typ.DoubleType
}
