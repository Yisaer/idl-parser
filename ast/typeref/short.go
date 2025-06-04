package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type ShortType struct {
	SelfType string `json:"self_type"`
}

func NewShortType() ShortType {
	return ShortType{SelfType: "short"}
}

func (t ShortType) TypeName() string {
	return "short"
}

func (ShortType) isTypeRef() {}

func ParseShort(code string) gomme.Result[ShortType, string] {
	return gomme.Map(
		gomme.Token[string]("short"),
		func(_ string) (ShortType, error) { return NewShortType(), nil },
	)(code)
}

func (ShortType) TypeRefType() typ.FieldRefType {
	return typ.ShortType
}

type UnsignedShortType struct {
	SelfType string `json:"self_type"`
}

func NewUnsignedShortType() UnsignedShortType {
	return UnsignedShortType{SelfType: "unsigned short"}
}

func (t UnsignedShortType) TypeName() string {
	return "unsigned short"
}

func ParseUnsignedShort(code string) gomme.Result[UnsignedShortType, string] {
	return gomme.Map(
		gomme.SeparatedPair(
			gomme.Token[string]("unsigned"),
			gomme.Whitespace1[string](),
			gomme.Token[string]("short"),
		),
		func(_ gomme.PairContainer[string, string]) (UnsignedShortType, error) {
			return NewUnsignedShortType(), nil
		},
	)(code)
}

func (UnsignedShortType) TypeRefType() typ.FieldRefType {
	return typ.UnsignedShortType
}
