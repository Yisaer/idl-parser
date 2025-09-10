package struct_type

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/annotation"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
	"github.com/yisaer/idl-parser/ast/utils"
)

type Field struct {
	Annotations annotation.Annotations `json:"annotations,omitempty"`
	Type        typeref.TypeRef        `json:"type"`
	Name        string                 `json:"name"`
}

type Struct struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
	Type   string  `json:"type"`
}

func (s Struct) GetName() string {
	return s.Name
}

func (Struct) ModuleContentType() typ.ModuleContentType {
	return typ.StructType
}

func parseField(code string) gomme.Result[Field, string] {
	var typeRefParser gomme.Parser[string, typeref.TypeRef] = typeref.ParseTypeRef
	var annotationsParser gomme.Parser[string, annotation.Annotations] = annotation.ParseAnnotations
	var optionalWhitespace gomme.Parser[string, string] = gomme.Whitespace0[string]()
	return gomme.Map(
		gomme.Pair(
			gomme.SeparatedPair(
				gomme.Optional(annotationsParser),
				optionalWhitespace,
				gomme.SeparatedPair(
					typeRefParser,
					gomme.Whitespace1[string](),
					gomme.Recognize(
						utils.ParseValidChar,
					),
				),
			),
			gomme.Pair(
				gomme.Whitespace0[string](),
				gomme.Token[string](";"),
			),
		),
		func(output gomme.PairContainer[gomme.PairContainer[annotation.Annotations, gomme.PairContainer[typeref.TypeRef, string]], gomme.PairContainer[string, string]]) (Field, error) {
			fieldData := output.Left
			if len(fieldData.Left) < 1 {
				fieldData.Left = nil
			}
			return Field{
				Annotations: fieldData.Left,
				Type:        fieldData.Right.Left,
				Name:        fieldData.Right.Right,
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
			gomme.Recognize(utils.ParseValidChar),
		)(structTokenResult.Remaining)
	if nameResult.Err != nil {
		return gomme.Failure[string, Struct](nameResult.Err, code)
	}
	fieldsResult := utils.InEmpty(
		gomme.Delimited(
			utils.InEmpty(gomme.Token[string]("{")),
			gomme.SeparatedList0(parseField, gomme.Whitespace0[string]()),
			utils.InEmpty(gomme.Token[string]("}")),
		))(nameResult.Remaining)
	if fieldsResult.Err != nil {
		return gomme.Failure[string, Struct](fieldsResult.Err, code)
	}
	return gomme.Success(
		Struct{
			Name:   nameResult.Output,
			Fields: fieldsResult.Output,
			Type:   typ.ModuleContentTypeToString(typ.StructType),
		},
		fieldsResult.Remaining,
	)
}
