package struct_type

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
	"github.com/yisaer/idl-parser/ast/utils"
)

type Field struct {
	Type typeref.TypeRef `json:"type"`
	Name string          `json:"name"`
}

type Struct struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

func (Struct) ModuleContentType() typ.ModuleContentType {
	return typ.StructType
}

func parseField(code string) gomme.Result[Field, string] {
	var typeRefParser gomme.Parser[string, typeref.TypeRef] = typeref.ParseTypeRef
	return gomme.Map(
		gomme.SeparatedPair(
			typeRefParser,
			gomme.Whitespace1[string](),
			gomme.Recognize(
				gomme.Pair(gomme.Alpha1[string](), gomme.Alphanumeric0[string]()),
			),
		),
		func(output gomme.PairContainer[typeref.TypeRef, string]) (Field, error) {
			return Field{
				Type: output.Left,
				Name: output.Right,
			}, nil
		},
	)(code)
}

func Parse(code string) gomme.Result[Struct, string] {
	bitsetTokenResult := gomme.Token[string]("struct")(code)
	if bitsetTokenResult.Err != nil {
		return gomme.Failure[string, Struct](bitsetTokenResult.Err, code)
	}
	nameResult :=
		utils.InEmpty(
			gomme.Recognize(gomme.Pair(gomme.Alpha1[string](), gomme.Alphanumeric0[string]())),
		)(bitsetTokenResult.Remaining)
	if nameResult.Err != nil {
		return gomme.Failure[string, Struct](nameResult.Err, code)
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
		return gomme.Failure[string, Struct](fieldsResult.Err, code)
	}
	return gomme.Success(
		Struct{
			Name:   nameResult.Output,
			Fields: fieldsResult.Output,
		},
		fieldsResult.Remaining,
	)
}
