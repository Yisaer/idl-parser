package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

type Sequence struct {
	InnerType TypeRef `json:"inner_type"`
}

func (s Sequence) TypeRefType() typ.FieldRefType {
	return typ.SequenceType
}

func (s Sequence) TypeName() string {
	return s.InnerType.TypeName()
}

func ParseSequence(code string) gomme.Result[Sequence, string] {
	result := gomme.Map(
		gomme.Preceded(
			gomme.Token[string]("sequence"),
			utils.InLeftEmpty(gomme.Delimited(
				gomme.Token[string]("<"),
				utils.InEmpty(ParseTypeRef),
				gomme.Token[string](">"),
			))),
		func(innerType TypeRef) (Sequence, error) {
			return Sequence{InnerType: innerType}, nil
		},
	)(code)
	return result
}
