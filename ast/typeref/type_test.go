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
		{"bitfield<8>", BitFieldType{Width: 8, SelfType: "bitfield"}},
		{"bitfield<16>", BitFieldType{Width: 16, SelfType: "bitfield"}},
		{"bitfield<32>", BitFieldType{Width: 32, SelfType: "bitfield"}},
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
		{"short", ShortType{SelfType: "short"}},
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
		{"unsigned short", UnsignedShortType{SelfType: "unsigned short"}},
		{"unsigned  short", UnsignedShortType{SelfType: "unsigned short"}},
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
		{"long", LongType{SelfType: "long"}},
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
		{"unsigned long", UnsignedLongType{SelfType: "unsigned long"}},
		{"unsigned  long", UnsignedLongType{SelfType: "unsigned long"}},
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
		{"long long", LongLongType{SelfType: "long long"}},
		{"long  long", LongLongType{SelfType: "long long"}},
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
		{"unsigned long long", UnsignedLongLongType{SelfType: "unsigned long long"}},
		{"unsigned  long  long", UnsignedLongLongType{SelfType: "unsigned long long"}},
	}

	for _, test := range tests {
		result := ParseUnsignedLongLong(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"boolean", BooleanType{SelfType: "boolean"}},
	}

	for _, test := range tests {
		result := ParseBoolean(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"float", FloatType{SelfType: "float"}},
	}

	for _, test := range tests {
		result := ParseFloat(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		input    string
		expected TypeRef
	}{
		{"string", StringType{SelfType: "string"}},
	}

	for _, test := range tests {
		result := ParseString(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}
