package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type TypeName struct {
	Name string `json:"name"`
}

func (TypeName) TypeRefType() typ.FieldRefType {
	return typ.SelfDefinedTypeType
}

func (t TypeName) TypeName() string {
	return t.Name
}

func ParseTypeName(code string) gomme.Result[TypeName, string] {
	return gomme.Map(
		gomme.Recognize(
			gomme.Pair(
				gomme.Alpha1[string](),
				gomme.Alphanumeric0[string](),
			),
		),
		func(name string) (TypeName, error) {
			return TypeName{Name: name}, nil
		},
	)(code)
}
