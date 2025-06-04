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
		{"sequence<octet>", Sequence{SelfType: "sequence", InnerType: OctetType{SelfType: "octet"}}},
		{"sequence<long long>", Sequence{SelfType: "sequence", InnerType: LongLongType{SelfType: "long long"}}},
		{"sequence<long>", Sequence{SelfType: "sequence", InnerType: LongType{SelfType: "long"}}},
		{"sequence<unsigned long>", Sequence{SelfType: "sequence", InnerType: UnsignedLongType{SelfType: "unsigned long"}}},
		{"sequence<unsigned long long>", Sequence{SelfType: "sequence", InnerType: UnsignedLongLongType{SelfType: "unsigned long long"}}},
		{"sequence<idbits>", Sequence{SelfType: "sequence", InnerType: TypeName{Name: "idbits", SelfType: "idbits"}}},
	}

	for _, test := range tests {
		result := ParseSequence(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}
