package utils

import (
	"unicode"
	"unicode/utf8"
	
	"github.com/oleiade/gomme"
)

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

func InLeftEmpty[Output any](parser gomme.Parser[string, Output]) gomme.Parser[string, Output] {
	return gomme.Preceded(
		ParseEmpty0,
		parser,
	)
}

func InEmpty[Output any](parser gomme.Parser[string, Output]) gomme.Parser[string, Output] {
	return gomme.Delimited(
		ParseEmpty0,
		parser,
		ParseEmpty0,
	)
}

func Identifier(code string) gomme.Result[string, string] {
	return gomme.Recognize(
		gomme.Pair(
			gomme.Alpha1[string](),
			gomme.Alphanumeric0[string](),
		))(code)
}

func ParseDigit(code string) gomme.Result[int, string] {
	return gomme.Map(
		gomme.Digit1[string](),
		func(digits string) (int, error) {
			result := 0
			for _, digit := range digits {
				result = result*10 + int(digit-'0')
			}
			return result, nil
		},
	)(code)
}

func ParseValidChar(code string) gomme.Result[string, string] {
	var matched string
	remaining := code
	for len(remaining) > 0 {
		r, size := utf8.DecodeRuneInString(remaining)
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '/' || r == '_') {
			break
		}
		matched += string(r)
		remaining = remaining[size:]
	}
	if len(matched) == 0 && len(code) > 0 {
		return gomme.Failure[string, string](gomme.NewError[string](code, "expect rune"), code)
	}
	return gomme.Success[string, string](matched, remaining)
}
