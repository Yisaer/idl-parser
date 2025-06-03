package bitset

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
	"github.com/yisaer/idl-parser/ast/utils"
)

type Field struct {
	Type typeref.BitFieldType `json:"type"`
	Name string               `json:"name"`
}

type BitSet struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

func (BitSet) ModuleContentType() typ.ModuleContentType {
	return typ.BitSetType
}

func parseField(code string) gomme.Result[Field, string] {
	var bitFieldParser gomme.Parser[string, typeref.BitFieldType] = typeref.ParseBitField
	return gomme.Map(
		gomme.SeparatedPair(
			bitFieldParser,
			gomme.Whitespace1[string](),
			gomme.Recognize(
				gomme.Pair(gomme.Alpha1[string](), gomme.Alphanumeric0[string]()),
			),
		),
		func(output gomme.PairContainer[typeref.BitFieldType, string]) (Field, error) {
			return Field{
				Type: output.Left,
				Name: output.Right,
			}, nil
		},
	)(code)
}

func Parse(code string) gomme.Result[BitSet, string] {
	bitsetTokenResult := gomme.Token[string]("bitset")(code)
	if bitsetTokenResult.Err != nil {
		return gomme.Failure[string, BitSet](bitsetTokenResult.Err, code)
	}
	nameResult :=
		utils.InEmpty(
			gomme.Recognize(gomme.Pair(gomme.Alpha1[string](), gomme.Alphanumeric0[string]())),
		)(bitsetTokenResult.Remaining)
	if nameResult.Err != nil {
		return gomme.Failure[string, BitSet](nameResult.Err, code)
	}
	fieldsResult := utils.InEmpty(
		gomme.Delimited(
			utils.InEmpty(gomme.Token[string]("{")),
			gomme.SeparatedList0(parseField, utils.InEmpty(gomme.Token[string](";"))),
			gomme.Pair(
				gomme.Optional(utils.InEmpty(gomme.Token[string](";"))),
				utils.InEmpty(gomme.Token[string]("}")),
			),
		))(nameResult.Remaining)
	if fieldsResult.Err != nil {
		return gomme.Failure[string, BitSet](fieldsResult.Err, code)
	}
	return gomme.Success(
		BitSet{
			Name:   nameResult.Output,
			Fields: fieldsResult.Output,
		},
		fieldsResult.Remaining,
	)
}
