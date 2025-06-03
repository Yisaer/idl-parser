package utils

import "github.com/oleiade/gomme"

func ParseComment(code string) gomme.Result[string, string] {
	return gomme.Recognize(
		gomme.Pair(gomme.Pair(
			gomme.Token[string]("//"),
			gomme.TakeUntil(
				gomme.Alternative(
					gomme.CRLF[string](),
					gomme.Map(gomme.CR[string](), func(_ rune) (string, error) { return "", nil }),
					gomme.Map(gomme.LF[string](), func(_ rune) (string, error) { return "", nil }),
				),
			),
		),
			gomme.Alternative(
				gomme.CRLF[string](),
				gomme.Map(gomme.CR[string](), func(_ rune) (string, error) { return "", nil }),
				gomme.Map(gomme.LF[string](), func(_ rune) (string, error) { return "", nil }),
			),
		))(code)
}

func ParseEmpty0(code string) gomme.Result[string, string] {
	return gomme.Recognize(
		gomme.Many0(
			gomme.Alternative(
				ParseComment,
				gomme.Whitespace1[string](),
			)))(code)
}

func ParseEmpty1(code string) gomme.Result[string, string] {
	return gomme.Recognize(
		gomme.Many1(
			gomme.Alternative(
				ParseComment,
				gomme.Whitespace1[string](),
			)))(code)
}

func InEmpty[Output any](parser gomme.Parser[string, Output]) gomme.Parser[string, Output] {
	return gomme.Delimited(
		ParseEmpty0,
		parser,
		ParseEmpty0,
	)
}
