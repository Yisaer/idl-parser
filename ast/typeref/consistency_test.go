package typeref

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/typ"
)

func TestParseTypeRefConsistency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected typ.FieldRefType
	}{
		// 基础类型测试
		{
			name:     "short",
			input:    "short",
			expected: typ.ShortType,
		},
		{
			name:     "long",
			input:    "long",
			expected: typ.LongType,
		},
		{
			name:     "boolean",
			input:    "boolean",
			expected: typ.BooleanType,
		},
		{
			name:     "float",
			input:    "float",
			expected: typ.FloatType,
		},
		{
			name:     "string",
			input:    "string",
			expected: typ.StringType,
		},
		{
			name:     "octet",
			input:    "octet",
			expected: typ.OctetType,
		},
		{
			name:     "unsigned short",
			input:    "unsigned short",
			expected: typ.UnsignedShortType,
		},
		{
			name:     "unsigned long",
			input:    "unsigned long",
			expected: typ.UnsignedLongType,
		},
		{
			name:     "long long",
			input:    "long long",
			expected: typ.LongLongType,
		},
		{
			name:     "unsigned long long",
			input:    "unsigned long long",
			expected: typ.UnsignedLongLongType,
		},
		// 数组类型测试
		{
			name:     "short array",
			input:    "short[8]",
			expected: typ.ArrayType,
		},
		{
			name:     "long array",
			input:    "long[16]",
			expected: typ.ArrayType,
		},
		{
			name:     "boolean array",
			input:    "boolean[3]",
			expected: typ.ArrayType,
		},
		{
			name:     "float array",
			input:    "float[32]",
			expected: typ.ArrayType,
		},
		{
			name:     "octet array",
			input:    "octet[8]",
			expected: typ.ArrayType,
		},
		{
			name:     "unsigned short array",
			input:    "unsigned short[4]",
			expected: typ.ArrayType,
		},
		// 类型名测试
		{
			name:     "type name",
			input:    "MyType",
			expected: typ.SelfDefinedTypeType,
		},
		{
			name:     "type name array",
			input:    "MyType[8]",
			expected: typ.ArrayType,
		},
		// 序列测试
		{
			name:     "sequence",
			input:    "sequence<long>",
			expected: typ.SequenceType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseTypeRef(test.input)
			if result.Err != nil {
				t.Errorf("Unexpected error: %v", result.Err)
				return
			}
			require.Equal(t, test.expected, result.Output.TypeRefType())
		})
	}
}

func TestParseBaseTypeVsParseTypeRefConsistency(t *testing.T) {
	// 测试基础类型在 ParseBaseType 和 ParseTypeRef 中的一致性
	baseTypes := []string{
		"short",
		"long", 
		"boolean",
		"float",
		"string",
		"octet",
		"unsigned short",
		"unsigned long",
		"long long",
		"unsigned long long",
		"MyType",
		"sequence<long>",
	}

	for _, input := range baseTypes {
		t.Run(input, func(t *testing.T) {
			baseResult := ParseBaseType(input)
			typeRefResult := ParseTypeRef(input)
			
			if baseResult.Err != nil {
				t.Errorf("ParseBaseType unexpected error: %v", baseResult.Err)
				return
			}
			if typeRefResult.Err != nil {
				t.Errorf("ParseTypeRef unexpected error: %v", typeRefResult.Err)
				return
			}
			
			require.Equal(t, baseResult.Output.TypeRefType(), typeRefResult.Output.TypeRefType())
			require.Equal(t, baseResult.Remaining, typeRefResult.Remaining)
		})
	}
}