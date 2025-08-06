package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type TypeRef interface {
	TypeRefType() typ.FieldRefType
	TypeName() string
}

func ParseTypeRef(code string) gomme.Result[TypeRef, string] {
	return gomme.Alternative(
		gomme.Map(ParseSequence, func(seq Sequence) (TypeRef, error) { return seq, nil }),
		gomme.Map(ParseOctet, func(octet OctetType) (TypeRef, error) { return octet, nil }),
		gomme.Map(ParseShort, func(short ShortType) (TypeRef, error) { return short, nil }),
		gomme.Map(ParseUnsignedShort, func(us UnsignedShortType) (TypeRef, error) { return us, nil }),
		gomme.Map(ParseLongLong, func(longlong LongLongType) (TypeRef, error) { return longlong, nil }),
		gomme.Map(ParseLong, func(long LongType) (TypeRef, error) { return long, nil }),
		gomme.Map(ParseUnsignedLongLong, func(ull UnsignedLongLongType) (TypeRef, error) { return ull, nil }),
		gomme.Map(ParseUnsignedLong, func(ul UnsignedLongType) (TypeRef, error) { return ul, nil }),
		gomme.Map(ParseBoolean, func(b BooleanType) (TypeRef, error) { return b, nil }),
		gomme.Map(ParseBitField, func(bitfield BitFieldType) (TypeRef, error) { return bitfield, nil }),
		gomme.Map(ParseTypeName, func(name TypeName) (TypeRef, error) { return name, nil }),
	)(code)
}
