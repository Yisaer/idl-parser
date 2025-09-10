package typeref

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/typ"
)

func TestParseOctet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected OctetType
	}{
		{
			name:     "parse plain octet",
			input:    "octet",
			expected: NewOctetType(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseOctet(test.input)
			require.Nil(t, result.Err)
			require.Equal(t, test.expected, result.Output)
		})
	}
}

func TestParseOctetArrayViaTypeRef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ArrayType
	}{
		{
			name: "parse octet array with size 8",
			input: "octet[8]",
			expected: ArrayType{
				InnerType: NewOctetType(),
				ArraySize: 8,
			},
		},
		{
			name: "parse octet array with size 1",
			input: "octet[1]",
			expected: ArrayType{
				InnerType: NewOctetType(),
				ArraySize: 1,
			},
		},
		{
			name: "parse octet array with size 16",
			input: "octet[16]",
			expected: ArrayType{
				InnerType: NewOctetType(),
				ArraySize: 16,
			},
		},
		{
			name: "parse octet array with size 255",
			input: "octet[255]",
			expected: ArrayType{
				InnerType: NewOctetType(),
				ArraySize: 255,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseTypeRef(test.input)
			require.Nil(t, result.Err)
			require.Equal(t, typ.ArrayType, result.Output.TypeRefType())
			
			arrayType, ok := result.Output.(ArrayType)
			require.True(t, ok)
			require.Equal(t, test.expected, arrayType)
		})
	}
}

func TestParseOctetTypeRefType(t *testing.T) {
	tests := []struct {
		name     string
		octet    OctetType
		expected typ.FieldRefType
	}{
		{
			name:     "plain octet returns OctetType",
			octet:    NewOctetType(),
			expected: typ.OctetType,
		},
		{
			name:     "octet array returns OctetArrayType",
			octet:    NewOctetArrayType(8),
			expected: typ.OctetArrayType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.octet.TypeRefType())
		})
	}
}

func TestParseOctetError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "invalid token",
			input: "invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseOctet(test.input)
			require.NotNil(t, result.Err)
		})
	}
}
