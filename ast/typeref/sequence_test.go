package typeref

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeq(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"sequence<octet>", Sequence{InnerType: OctetType{}}},
		{"sequence<long long>", Sequence{InnerType: LongLongType{}}},
		{"sequence<long>", Sequence{InnerType: LongType{}}},
		{"sequence<unsigned long>", Sequence{InnerType: UnsignedLongType{}}},
		{"sequence<unsigned long long>", Sequence{InnerType: UnsignedLongLongType{}}},
		{"sequence<idbits>", Sequence{InnerType: TypeName{Name: "idbits"}}},
	}

	for _, test := range tests {
		result := ParseSequence(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}
