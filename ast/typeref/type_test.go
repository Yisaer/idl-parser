package typeref

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/typ"
)

func TestParseBitField(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"bitfield<8>", BitFieldType{Width: 8}},
		{"bitfield<16>", BitFieldType{Width: 16}},
		{"bitfield<32>", BitFieldType{Width: 32}},
	}

	for _, test := range tests {
		result := ParseBitField(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestShort(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"short", ShortType{}},
	}

	for _, test := range tests {
		result := ParseShort(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestUnsignedShort(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"unsigned short", UnsignedShortType{}},
		{"unsigned  short", UnsignedShortType{}},
	}

	for _, test := range tests {
		result := ParseUnsignedShort(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestLong(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"long", LongType{}},
	}

	for _, test := range tests {
		result := ParseLong(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestUnsignedLong(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"unsigned long", UnsignedLongType{}},
		{"unsigned  long", UnsignedLongType{}},
	}

	for _, test := range tests {
		result := ParseUnsignedLong(test.input)
		require.Equal(t, typ.UnsignedLongType, result.Output.TypeRefType())
	}
}

func TestLongLong(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"long long", LongLongType{}},
		{"long  long", LongLongType{}},
	}

	for _, test := range tests {
		result := ParseLongLong(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestUnsignedLongLong(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"unsigned long long", UnsignedLongLongType{}},
		{"unsigned  long  long", UnsignedLongLongType{}},
	}

	for _, test := range tests {
		result := ParseUnsignedLongLong(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}
