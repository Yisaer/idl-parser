package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

type Sequence struct {
	SelfType  string  `json:"self_type"`
	InnerType TypeRef `json:"inner_type"`
}

func NewSequence(innerType TypeRef) Sequence {
	return Sequence{SelfType: "sequence", InnerType: innerType}
}

func (s Sequence) TypeRefType() typ.FieldRefType {
	return typ.SequenceType
}

func (s Sequence) TypeName() string {
	return "sequence"
}

func ParseSequence(code string) gomme.Result[Sequence, string] {
	return gomme.Map(
		gomme.Preceded(
			gomme.Token[string]("sequence"),
			utils.InEmpty(gomme.Delimited(
				gomme.Token[string]("<"),
				utils.InEmpty(ParseTypeRef),
				gomme.Token[string](">"),
			))),
		func(innerType TypeRef) (Sequence, error) {
			return NewSequence(innerType), nil
		},
	)(code)
}
