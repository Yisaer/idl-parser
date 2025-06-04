package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type OctetType struct {
	SelfType string `json:"self_type"`
}

func NewOctetType() OctetType {
	return OctetType{SelfType: "octet"}
}

func (t OctetType) TypeName() string {
	return t.SelfType
}

func ParseOctet(code string) gomme.Result[OctetType, string] {
	return gomme.Map(
		gomme.Token[string]("octet"),
		func(token string) (OctetType, error) {
			return NewOctetType(), nil
		},
	)(code)
}

func (OctetType) TypeRefType() typ.FieldRefType {
	return typ.OctetType
}
