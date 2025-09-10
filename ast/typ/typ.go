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
	OctetArrayType // 保留以向后兼容
	ShortType
	UnsignedShortType
	LongType
	UnsignedLongType
	LongLongType
	UnsignedLongLongType
	SelfDefinedTypeType
	SequenceType
	BooleanType
	FloatType
	StringType
	ArrayType // 新的通用数组类型
)
