package struct_type

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/annotation"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
	"github.com/yisaer/idl-parser/ast/utils"
)

type Field struct {
	Annotations annotation.Annotations `json:"annotations"`
	Type        typeref.TypeRef        `json:"type"`
	Name        string                 `json:"name"`
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
	var annotationsParser gomme.Parser[string, annotation.Annotations] = annotation.ParseAnnotations
	var optionalWhitespace gomme.Parser[string, string] = gomme.Whitespace0[string]()
	return gomme.Map(
		gomme.SeparatedPair(
			gomme.Optional(annotationsParser),
			optionalWhitespace,
			gomme.SeparatedPair(
				typeRefParser,
				gomme.Whitespace1[string](),
				gomme.Recognize(
					gomme.Pair(gomme.Alpha1[string](), gomme.Alphanumeric0[string]()),
				),
			),
		),
		func(output gomme.PairContainer[annotation.Annotations, gomme.PairContainer[typeref.TypeRef, string]]) (Field, error) {
			if len(output.Left) < 1 {
				output.Left = nil
			}
			return Field{
				Annotations: output.Left,
				Type:        output.Right.Left,
				Name:        output.Right.Right,
			}, nil
		},
	)(code)
}

func Parse(code string) gomme.Result[Struct, string] {
	structTokenResult := gomme.Token[string]("struct")(code)
	if structTokenResult.Err != nil {
		return gomme.Failure[string, Struct](structTokenResult.Err, code)
	}
	nameResult :=
		utils.InEmpty(
			gomme.Recognize(gomme.Pair(gomme.Alpha1[string](), gomme.Alphanumeric0[string]())),
		)(structTokenResult.Remaining)
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
