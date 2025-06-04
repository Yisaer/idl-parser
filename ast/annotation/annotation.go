package annotation

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/utils"
)

type Annotations []Annotation
type Annotation struct {
	Name   string            `json:"name"`
	Values map[string]string `json:"values,omitempty"`
}

func parseKVPairs(code string) gomme.Result[map[string]string, string] {
	return gomme.Map(gomme.SeparatedList0(
		gomme.SeparatedPair(
			gomme.Recognize(
				gomme.Pair(
					gomme.Alpha1[string](),
					gomme.Alphanumeric0[string](),
				)),
			utils.InEmpty(gomme.Token[string]("=")),
			gomme.Recognize(
				gomme.Pair(
					gomme.Alpha1[string](),
					gomme.Alphanumeric0[string](),
				)),
		),
		utils.InEmpty(gomme.Token[string](",")),
	),
		func(pairs []gomme.PairContainer[string, string]) (map[string]string, error) {
			values := make(map[string]string)
			for _, pair := range pairs {
				values[pair.Left] = pair.Right
			}
			return values, nil
		})(code)
}

func ParseAnnotation(code string) gomme.Result[Annotation, string] {
	return gomme.Map(
		gomme.SeparatedPair(
			gomme.Preceded(
				gomme.Token[string]("@"),
				utils.Identifier,
			),
			gomme.Whitespace0[string](),
			gomme.Optional(
				gomme.Delimited(
					gomme.Token[string]("("),
					parseKVPairs,
					gomme.Token[string](")"),
				),
			),
		),
		func(output gomme.PairContainer[string, map[string]string]) (Annotation, error) {
			if len(output.Right) < 1 {
				return Annotation{
					Name: output.Left,
				}, nil
			}
			return Annotation{
				Name:   output.Left,
				Values: output.Right,
			}, nil
		},
	)(code)
}

func ParseAnnotations(code string) gomme.Result[Annotations, string] {
	return gomme.Map(
		gomme.Many0(
			gomme.Preceded(
				gomme.Whitespace0[string](),
				ParseAnnotation,
			),
		),
		func(annos []Annotation) (Annotations, error) {
			return Annotations(annos), nil
		},
	)(code)
}
