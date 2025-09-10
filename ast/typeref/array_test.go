package typeref

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/typ"
)

func TestParseArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ArrayType
	}{
		{
			name: "parse short array with size 8",
			input: "short[8]",
			expected: ArrayType{
				InnerType: NewShortType(),
				ArraySize: 8,
			},
		},
		{
			name: "parse long array with size 16",
			input: "long[16]",
			expected: ArrayType{
				InnerType: NewLongType(),
				ArraySize: 16,
			},
		},
		{
			name: "parse unsigned short array with size 4",
			input: "unsigned short[4]",
			expected: ArrayType{
				InnerType: NewUnsignedShortType(),
				ArraySize: 4,
			},
		},
		{
			name: "parse boolean array with size 1",
			input: "boolean[1]",
			expected: ArrayType{
				InnerType: NewBooleanType(),
				ArraySize: 1,
			},
		},
		{
			name: "parse float array with size 32",
			input: "float[32]",
			expected: ArrayType{
				InnerType: NewFloatType(),
				ArraySize: 32,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseArray(test.input)
			require.Nil(t, result.Err)
			require.Equal(t, test.expected, result.Output)
		})
	}
}

func TestParseArrayError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "plain short without array",
			input: "short",
		},
		{
			name:  "invalid array format",
			input: "short[abc]",
		},
		{
			name:  "missing closing bracket",
			input: "short[8",
		},
		{
			name:  "missing opening bracket",
			input: "short8]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseArray(test.input)
			require.Error(t, result.Err)
		})
	}
}

func TestParseTypeRefWithArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected typ.FieldRefType
	}{
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
			name:     "unsigned short array",
			input:    "unsigned short[4]",
			expected: typ.ArrayType,
		},
		{
			name:     "plain short",
			input:    "short",
			expected: typ.ShortType,
		},
		{
			name:     "plain long",
			input:    "long",
			expected: typ.LongType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseTypeRef(test.input)
			require.Nil(t, result.Err)
			require.Equal(t, test.expected, result.Output.TypeRefType())
		})
	}
}

func TestArrayTypeMethods(t *testing.T) {
	arrayType := ArrayType{
		InnerType: NewShortType(),
		ArraySize: 8,
	}

	require.Equal(t, "short", arrayType.TypeName())
	require.Equal(t, typ.ArrayType, arrayType.TypeRefType())
}