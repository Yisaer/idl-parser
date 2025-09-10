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
	if len(code) > 0 && containsArraySuffix(code) {
		arrayResult := ParseArray(code)
		if arrayResult.Err == nil {
			return gomme.Result[TypeRef, string]{
				Output:    arrayResult.Output,
				Remaining: arrayResult.Remaining,
				Err:       nil,
			}
		}
	}

	// 否则使用原来的解析逻辑
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
		gomme.Map(ParseFloat, func(f FloatType) (TypeRef, error) { return f, nil }),
		gomme.Map(ParseString, func(s StringType) (TypeRef, error) { return s, nil }),
		gomme.Map(ParseBitField, func(bitfield BitFieldType) (TypeRef, error) { return bitfield, nil }),
		gomme.Map(ParseTypeName, func(name TypeName) (TypeRef, error) { return name, nil }),
	)(code)
}

// containsArraySuffix 检查字符串是否包含数组后缀 [数字]
func containsArraySuffix(code string) bool {
	if len(code) < 3 { // 至少需要 "a[1]" 的长度
		return false
	}

	// 从后往前查找 ']'
	lastBracket := -1
	for i := len(code) - 1; i >= 0; i-- {
		if code[i] == ']' {
			lastBracket = i
			break
		}
	}

	if lastBracket == -1 {
		return false
	}

	// 向前查找匹配的 '['
	for i := lastBracket - 1; i >= 0; i-- {
		if code[i] == '[' {
			// 检查 [ 和 ] 之间是否都是数字
			if i+1 < lastBracket {
				for j := i + 1; j < lastBracket; j++ {
					if code[j] < '0' || code[j] > '9' {
						return false
					}
				}
				return true
			}
			return false
		}
	}

	return false
}
