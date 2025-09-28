package typeref

import (
	"fmt"

	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

type StringType struct {
	SelfType string `json:"self_type"`
	// If Length is nil, this is a dynamic string, otherwise this is a fixed length string
	Length *int `json:"length,omitempty"`
}

func NewStringType() StringType {
	return StringType{
		SelfType: "string",
		Length:   nil,
	}
}

func NewFixedLengthStringType(length int) StringType {
	return StringType{
		SelfType: "string",
		Length:   &length,
	}
}

func (t StringType) TypeName() string {
	if t.Length != nil {
		return fmt.Sprintf("string<%d>", *t.Length)
	}
	return "string"
}

func ParseString(code string) gomme.Result[StringType, string] {
	stringTokenResult := gomme.Token[string]("string")(code)
	if stringTokenResult.Err != nil {
		return gomme.Result[StringType, string]{Err: stringTokenResult.Err}
	}

	remaining := stringTokenResult.Remaining

	if len(remaining) > 0 && remaining[0] == '<' {
		lengthResult := gomme.Delimited(
			gomme.Token[string]("<"),
			utils.ParseDigit,
			gomme.Token[string](">"),
		)(remaining)

		if lengthResult.Err != nil {
			return gomme.Result[StringType, string]{Err: lengthResult.Err}
		}

		return gomme.Result[StringType, string]{
			Output:    NewFixedLengthStringType(lengthResult.Output),
			Remaining: lengthResult.Remaining,
			Err:       nil,
		}
	}
	return gomme.Result[StringType, string]{
		Output:    NewStringType(),
		Remaining: remaining,
		Err:       nil,
	}
}

func (StringType) TypeRefType() typ.FieldRefType {
	return typ.StringType
}
