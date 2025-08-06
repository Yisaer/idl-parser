package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type FloatType struct {
	SelfType string `json:"self_type"`
}

func NewFloatType() FloatType {
	return FloatType{
		SelfType: "float",
	}
}

func (t FloatType) TypeName() string {
	return "float"
}

func ParseFloat(code string) gomme.Result[FloatType, string] {
	return gomme.Map(
		gomme.Token[string]("float"),
		func(_ string) (FloatType, error) { return NewFloatType(), nil },
	)(code)
}

func (FloatType) TypeRefType() typ.FieldRefType {
	return typ.FloatType
}
