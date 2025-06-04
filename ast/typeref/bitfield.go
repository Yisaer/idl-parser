package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type BitFieldType struct {
	Width    uint8  `json:"width"`
	SelfType string `json:"self_type"`
}

func NewBitField(width uint8) BitFieldType {
	return BitFieldType{
		Width:    width,
		SelfType: "bitfield",
	}
}

func (bt BitFieldType) TypeRefType() typ.FieldRefType {
	return typ.BitFieldType
}

func (bt BitFieldType) TypeName() string { return "bitset" }

func ParseBitField(code string) gomme.Result[BitFieldType, string] {
	return gomme.Map(
		gomme.Pair(
			gomme.Token[string]("bitfield"),
			gomme.Delimited(
				gomme.Token[string]("<"),
				gomme.UInt8[string](),
				gomme.Token[string](">"),
			),
		),
		func(pair gomme.PairContainer[string, uint8]) (BitFieldType, error) {
			return NewBitField(pair.Right), nil
		},
	)(code)
}
