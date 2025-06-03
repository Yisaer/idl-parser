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
		gomme.Map(ParseOctet, func(octet OctetType) (TypeRef, error) { return octet, nil }),
		gomme.Map(ParseShort, func(short ShortType) (TypeRef, error) { return short, nil }),
		gomme.Map(ParseUnsignedShort, func(us UnsignedShortType) (TypeRef, error) { return us, nil }),
		gomme.Map(ParseLong, func(long LongType) (TypeRef, error) { return long, nil }),
		gomme.Map(ParseUnsignedLongLong, func(ull UnsignedLongLongType) (TypeRef, error) { return ull, nil }),
		gomme.Map(ParseUnsignedLong, func(ul UnsignedLongType) (TypeRef, error) { return ul, nil }),
		gomme.Map(ParseLongLong, func(longlong LongLongType) (TypeRef, error) { return longlong, nil }),
		gomme.Map(ParseBitField, func(name BitFieldType) (TypeRef, error) { return name, nil }),
		gomme.Map(ParseTypeName, func(name TypeName) (TypeRef, error) { return name, nil }),
	)(code)
}
