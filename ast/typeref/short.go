package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type ShortType struct{}

func (t ShortType) TypeName() string {
	return "short"
}

func (ShortType) isTypeRef() {}

func ParseShort(code string) gomme.Result[ShortType, string] {
	return gomme.Map(
		gomme.Token[string]("short"),
		func(_ string) (ShortType, error) { return ShortType{}, nil },
	)(code)
}

func (ShortType) TypeRefType() typ.FieldRefType {
	return typ.ShortType
}

type UnsignedShortType struct{}

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
			return UnsignedShortType{}, nil
		},
	)(code)
}

func (UnsignedShortType) TypeRefType() typ.FieldRefType {
	return typ.UnsignedShortType
}
