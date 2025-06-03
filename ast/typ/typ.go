package typ

type ModuleContentType int

const (
	BitSetType ModuleContentType = iota
	StructType
)

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
)
