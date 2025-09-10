package converter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

func TestParseArrayTypes(t *testing.T) {
	tests := []struct {
		name     string
		typeStr  string
		data     []byte
		expected interface{}
	}{
		{
			name:     "short array [8]",
			typeStr:  "short[8]",
			data:     []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08},
			expected: []interface{}{int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7), int64(8)},
		},
		{
			name:     "long array [4]",
			typeStr:  "long[4]",
			data:     []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x04},
			expected: []interface{}{int64(1), int64(2), int64(3), int64(4)},
		},
		{
			name:     "boolean array [3]",
			typeStr:  "boolean[3]",
			data:     []byte{0x01, 0x00, 0x01},
			expected: []interface{}{true, false, true},
		},
		{
			name:     "unsigned short array [2]",
			typeStr:  "unsigned short[2]",
			data:     []byte{0x00, 0x01, 0x00, 0x02},
			expected: []interface{}{int64(1), int64(2)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			typeResult := typeref.ParseTypeRef(test.typeStr)
			require.Nil(t, typeResult.Err)
			require.Equal(t, typ.ArrayType, typeResult.Output.TypeRefType())
			converter := createTestConverter(false) // big endian
			result, remain, err := converter.parseDataByType(test.data, typeResult.Output, ast.Module{})
			require.NoError(t, err)
			require.Equal(t, test.expected, result)
			require.Empty(t, remain)
		})
	}
}

func TestParseArrayTypesLittleEndian(t *testing.T) {
	tests := []struct {
		name     string
		typeStr  string
		data     []byte
		expected interface{}
	}{
		{
			name:     "short array [4] little endian",
			typeStr:  "short[4]",
			data:     []byte{0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00},
			expected: []interface{}{int64(1), int64(2), int64(3), int64(4)},
		},
		{
			name:     "long array [2] little endian",
			typeStr:  "long[2]",
			data:     []byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00},
			expected: []interface{}{int64(1), int64(2)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 解析类型
			typeResult := typeref.ParseTypeRef(test.typeStr)
			require.Nil(t, typeResult.Err)
			require.Equal(t, typ.ArrayType, typeResult.Output.TypeRefType())
			converter := createTestConverter(true) // little endian
			result, remain, err := converter.parseDataByType(test.data, typeResult.Output, ast.Module{})
			require.NoError(t, err)
			require.Equal(t, test.expected, result)
			require.Empty(t, remain)
		})
	}
}

func TestParseArrayError(t *testing.T) {
	tests := []struct {
		name    string
		typeStr string
		data    []byte
	}{
		{
			name:    "insufficient data for short array",
			typeStr: "short[4]",
			data:    []byte{0x00, 0x01, 0x00, 0x02}, // 只有2个short，需要4个
		},
		{
			name:    "insufficient data for long array",
			typeStr: "long[2]",
			data:    []byte{0x00, 0x00, 0x00, 0x01}, // 只有1个long，需要2个
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 解析类型
			typeResult := typeref.ParseTypeRef(test.typeStr)
			require.Nil(t, typeResult.Err)

			// 创建 converter
			converter := createTestConverter(false)

			// 解析数据应该失败
			result, remain, err := converter.parseDataByType(test.data, typeResult.Output, ast.Module{})
			require.Error(t, err)
			require.Nil(t, result)
			require.Nil(t, remain)
		})
	}
}
