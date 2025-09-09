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
		{
			name:     "parse octet array with size 8",
			input:    "octet[8]",
			expected: NewOctetArrayType(8),
		},
		{
			name:     "parse octet array with size 1",
			input:    "octet[1]",
			expected: NewOctetArrayType(1),
		},
		{
			name:     "parse octet array with size 16",
			input:    "octet[16]",
			expected: NewOctetArrayType(16),
		},
		{
			name:     "parse octet array with size 255",
			input:    "octet[255]",
			expected: NewOctetArrayType(255),
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
			name:  "invalid format",
			input: "octet[",
		},
		{
			name:  "invalid format missing closing bracket",
			input: "octet[8",
		},
		{
			name:  "invalid format with non-digit",
			input: "octet[abc]",
		},
		{
			name:  "invalid format with empty brackets",
			input: "octet[]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseOctet(test.input)
			require.NotNil(t, result.Err)
		})
	}
}
