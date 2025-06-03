package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type OctetType struct{}

func (t OctetType) TypeName() string {
	return "octet"
}

func ParseOctet(code string) gomme.Result[OctetType, string] {
	return gomme.Map(
		gomme.Token[string]("octet"),
		func(token string) (OctetType, error) {
			return OctetType{}, nil
		},
	)(code)
}

func (OctetType) TypeRefType() typ.FieldRefType {
	return typ.OctetType
}
