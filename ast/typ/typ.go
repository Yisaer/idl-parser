package typ

type ModuleContentType int

const (
	BitSetType ModuleContentType = iota
	StructType
	ModuleType
)

func ModuleContentTypeToString(ct ModuleContentType) string {
	switch ct {
	case BitSetType:
		return "BitSet"
	case StructType:
		return "Struct"
	case ModuleType:
		return "Module"
	}
	return ""
}

type FieldRefType int

const (
	BitFieldType FieldRefType = iota
	OctetType
	ShortType
	UnsignedShortType
	LongType
	UnsignedLongType
	LongLongType
	UnsignedLongLongType
	SelfDefinedTypeType
	SequenceType
)
