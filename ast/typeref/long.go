package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type LongType struct {
	SelfType string `json:"self_type"`
}

func NewLongType() LongType {
	return LongType{
		SelfType: "long",
	}
}

func (t LongType) TypeName() string {
	return "long"
}

func ParseLong(code string) gomme.Result[LongType, string] {
	return gomme.Map(
		gomme.Token[string]("long"),
		func(_ string) (LongType, error) { return NewLongType(), nil },
	)(code)
}

func (LongType) TypeRefType() typ.FieldRefType {
	return typ.LongType
}

type UnsignedLongType struct {
	SelfType string `json:"self_type"`
}

func NewUnsignedLong() UnsignedLongType {
	return UnsignedLongType{SelfType: "unsigned long"}
}

func (t UnsignedLongType) TypeName() string {
	return "unsigned long"
}

func ParseUnsignedLong(code string) gomme.Result[UnsignedLongType, string] {
	return gomme.Map(
		gomme.SeparatedPair(
			gomme.Token[string]("unsigned"),
			gomme.Whitespace1[string](),
			gomme.Token[string]("long"),
		),
		func(_ gomme.PairContainer[string, string]) (UnsignedLongType, error) { return NewUnsignedLong(), nil },
	)(code)
}

func (UnsignedLongType) TypeRefType() typ.FieldRefType {
	return typ.UnsignedLongType
}

type LongLongType struct {
	SelfType string `json:"self_type"`
}

func NewLongLongType() LongLongType {
	return LongLongType{SelfType: "long long"}
}

func (t LongLongType) TypeName() string {
	return "long long"
}

func (LongLongType) isTypeRef() {}

func ParseLongLong(code string) gomme.Result[LongLongType, string] {
	return gomme.Map(
		gomme.SeparatedPair(
			gomme.Token[string]("long"),
			gomme.Whitespace1[string](),
			gomme.Token[string]("long"),
		),
		func(_ gomme.PairContainer[string, string]) (LongLongType, error) { return NewLongLongType(), nil },
	)(code)
}

func (LongLongType) TypeRefType() typ.FieldRefType {
	return typ.LongLongType
}

type UnsignedLongLongType struct {
	SelfType string `json:"self_type"`
}

func NewUnsignedLongLong() UnsignedLongLongType {
	return UnsignedLongLongType{SelfType: "unsigned long long"}
}

func (t UnsignedLongLongType) TypeName() string {
	return "unsigned long long"
}

func ParseUnsignedLongLong(code string) gomme.Result[UnsignedLongLongType, string] {
	return gomme.Map(
		gomme.SeparatedPair(
			gomme.Token[string]("unsigned"),
			gomme.Whitespace1[string](),
			gomme.SeparatedPair(
				gomme.Token[string]("long"),
				gomme.Whitespace1[string](),
				gomme.Token[string]("long"),
			)),
		func(pair gomme.PairContainer[string, gomme.PairContainer[string, string]]) (UnsignedLongLongType, error) {
			return NewUnsignedLongLong(), nil
		},
	)(code)
}

func (UnsignedLongLongType) TypeRefType() typ.FieldRefType {
	return typ.UnsignedLongLongType
}
