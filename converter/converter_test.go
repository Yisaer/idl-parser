package converter

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

func TestConverterDecode(t *testing.T) {
	c, err := NewIDLConverter("a.b.c", "./testdata/test.idl")
	require.NoError(t, err)
	testdata := []byte{41, 42}
	m, err := c.Decode(testdata)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{"id1": int64(41), "id2": int64(42)}, m)
}

func TestParseDataByType_Octet(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       int64
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse octet value 42 successfully",
			data:           []byte{42, 1, 2, 3},
			expected:       42,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse octet value 255 successfully",
			data:           []byte{255, 10, 20},
			expected:       255,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "parse octet value 0 successfully",
			data:           []byte{0, 100, 200},
			expected:       0,
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient",
			data:           []byte{},
			expected:       0,
			expectedRemain: nil,
			expectError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			octetType := typeref.NewOctetType()
			result, remain, err := parseDataByType(tt.data, octetType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_SequenceOfOctet(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       []interface{}
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse sequence with 3 octet elements",
			data:           []byte{0, 0, 0, 3, 10, 20, 30, 100, 200},
			expected:       []interface{}{int64(10), int64(20), int64(30)},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "parse sequence with 0 octet elements",
			data:           []byte{0, 0, 0, 0, 100, 200},
			expected:       []interface{}{},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "parse sequence with 1 octet element",
			data:           []byte{0, 0, 0, 1, 255, 100, 200},
			expected:       []interface{}{int64(255)},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "should return error when insufficient data for sequence length",
			data:           []byte{0, 0, 0},
			expected:       nil,
			expectedRemain: nil,
			expectError:    true,
		},
		{
			name:           "should return error when sequence length correct but data insufficient",
			data:           []byte{0, 0, 0, 3, 10, 20},
			expected:       nil,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			octetType := typeref.NewOctetType()
			sequenceType := typeref.NewSequence(octetType)

			result, remain, err := parseDataByType(tt.data, sequenceType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_UnsupportedType(t *testing.T) {
	mockType := &mockUnsupportedType{}

	result, remain, err := parseDataByType([]byte{1, 2, 3}, mockType, ast.Module{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported type")
	require.Nil(t, result)
	require.Nil(t, remain)
}

type mockUnsupportedType struct{}

func (m *mockUnsupportedType) TypeRefType() typ.FieldRefType {
	return typ.FieldRefType(999)
}

func (m *mockUnsupportedType) TypeName() string {
	return "unsupported_type"
}

func createSequenceData(length int, elements ...byte) []byte {
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(length))
	return append(lengthBytes, elements...)
}

func TestParseDataByType_Short(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       int64
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse short value 12345 successfully",
			data:           []byte{0x30, 0x39, 1, 2, 3},
			expected:       12345,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse short value -12345 successfully",
			data:           []byte{0xCF, 0xC7, 10, 20},
			expected:       -12345,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient for short",
			data:           []byte{0x30},
			expected:       0,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortType := typeref.NewShortType()
			result, remain, err := parseDataByType(tt.data, shortType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_UnsignedShort(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       int64
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse unsigned short value 65535 successfully",
			data:           []byte{0xFF, 0xFF, 1, 2, 3},
			expected:       65535,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse unsigned short value 0 successfully",
			data:           []byte{0x00, 0x00, 10, 20},
			expected:       0,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient for unsigned short",
			data:           []byte{0xFF},
			expected:       0,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unsignedShortType := typeref.NewUnsignedShortType()
			result, remain, err := parseDataByType(tt.data, unsignedShortType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_Long(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       int64
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse long value 1234567890 successfully",
			data:           []byte{0x49, 0x96, 0x02, 0xD2, 1, 2, 3},
			expected:       1234567890,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse long value -1234567890 successfully",
			data:           []byte{0xB6, 0x69, 0xFD, 0x2E, 10, 20},
			expected:       -1234567890,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient for long",
			data:           []byte{0x49, 0x96, 0x02},
			expected:       0,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			longType := typeref.NewLongType()
			result, remain, err := parseDataByType(tt.data, longType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_UnsignedLong(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       int64
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse unsigned long value 4294967295 successfully",
			data:           []byte{0xFF, 0xFF, 0xFF, 0xFF, 1, 2, 3},
			expected:       4294967295,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse unsigned long value 0 successfully",
			data:           []byte{0x00, 0x00, 0x00, 0x00, 10, 20},
			expected:       0,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient for unsigned long",
			data:           []byte{0xFF, 0xFF, 0xFF},
			expected:       0,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unsignedLongType := typeref.NewUnsignedLong()
			result, remain, err := parseDataByType(tt.data, unsignedLongType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_LongLong(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       int64
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse long long value 9223372036854775807 successfully",
			data:           []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 1, 2, 3},
			expected:       9223372036854775807,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse long long value 0 successfully",
			data:           []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 10, 20},
			expected:       0,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient for long long",
			data:           []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			expected:       0,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			longLongType := typeref.NewLongLongType()
			result, remain, err := parseDataByType(tt.data, longLongType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_UnsignedLongLong(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       int64
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse unsigned long long value 9223372036854775807 successfully",
			data:           []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 1, 2, 3},
			expected:       9223372036854775807,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse unsigned long long value 0 successfully",
			data:           []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 10, 20},
			expected:       0,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient for unsigned long long",
			data:           []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			expected:       0,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unsignedLongLongType := typeref.NewUnsignedLongLong()
			result, remain, err := parseDataByType(tt.data, unsignedLongLongType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_Boolean(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       bool
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse boolean true successfully",
			data:           []byte{0x01, 1, 2, 3},
			expected:       true,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse boolean false successfully",
			data:           []byte{0x00, 10, 20},
			expected:       false,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "parse boolean true with non-zero value",
			data:           []byte{0xFF, 100, 200},
			expected:       true,
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient for boolean",
			data:           []byte{},
			expected:       false,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booleanType := typeref.NewBooleanType()
			result, remain, err := parseDataByType(tt.data, booleanType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_Float(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       float64
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse float value 3.14 successfully",
			data:           []byte{0x40, 0x48, 0xF5, 0xC3, 1, 2, 3},
			expected:       3.14,
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse float value 0.0 successfully",
			data:           []byte{0x00, 0x00, 0x00, 0x00, 10, 20},
			expected:       0.0,
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "parse float value -2.5 successfully",
			data:           []byte{0xC0, 0x20, 0x00, 0x00, 100, 200},
			expected:       -2.5,
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "should return error when data insufficient for float",
			data:           []byte{0x40, 0x48, 0xF5},
			expected:       0.0,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			floatType := typeref.NewFloatType()
			result, remain, err := parseDataByType(tt.data, floatType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.InDelta(t, tt.expected, result, 0.001) // 使用 InDelta 比较浮点数
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_SequenceOfShort(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       []interface{}
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse sequence with 2 short elements",
			data:           []byte{0, 0, 0, 2, 0x12, 0x34, 0x56, 0x78, 100, 200},
			expected:       []interface{}{int64(0x1234), int64(0x5678)},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "parse sequence with 0 short elements",
			data:           []byte{0, 0, 0, 0, 100, 200},
			expected:       []interface{}{},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "should return error when insufficient data for sequence of shorts",
			data:           []byte{0, 0, 0, 2, 0x12},
			expected:       nil,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortType := typeref.NewShortType()
			sequenceType := typeref.NewSequence(shortType)

			result, remain, err := parseDataByType(tt.data, sequenceType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_SequenceOfBoolean(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       []interface{}
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse sequence with 3 boolean elements",
			data:           []byte{0, 0, 0, 3, 0x01, 0x00, 0xFF, 100, 200},
			expected:       []interface{}{true, false, true},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "parse sequence with 1 boolean element",
			data:           []byte{0, 0, 0, 1, 0x00, 100, 200},
			expected:       []interface{}{false},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "should return error when insufficient data for sequence of booleans",
			data:           []byte{0, 0, 0, 2, 0x01},
			expected:       nil,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			booleanType := typeref.NewBooleanType()
			sequenceType := typeref.NewSequence(booleanType)

			result, remain, err := parseDataByType(tt.data, sequenceType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_String(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       string
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse string 'hello' successfully",
			data:           []byte{0, 0, 0, 5, 'h', 'e', 'l', 'l', 'o', 1, 2, 3},
			expected:       "hello",
			expectedRemain: []byte{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "parse empty string successfully",
			data:           []byte{0, 0, 0, 0, 100, 200},
			expected:       "",
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "parse string with special characters successfully",
			data:           []byte{0, 0, 0, 3, 'a', 'b', 'c', 10, 20},
			expected:       "abc",
			expectedRemain: []byte{10, 20},
			expectError:    false,
		},
		{
			name:           "parse string with unicode characters successfully",
			data:           []byte{0, 0, 0, 6, 0xE4, 0xB8, 0xAD, 0xE6, 0x96, 0x87, 100, 200},
			expected:       "中文",
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "should return error when insufficient data for string length",
			data:           []byte{0, 0, 0},
			expected:       "",
			expectedRemain: nil,
			expectError:    true,
		},
		{
			name:           "should return error when string length correct but data insufficient",
			data:           []byte{0, 0, 0, 5, 'h', 'e', 'l'},
			expected:       "",
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stringType := typeref.NewStringType()
			result, remain, err := parseDataByType(tt.data, stringType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}

func TestParseDataByType_SequenceOfString(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expected       []interface{}
		expectedRemain []byte
		expectError    bool
	}{
		{
			name:           "parse sequence with 2 string elements",
			data:           []byte{0, 0, 0, 2, 0, 0, 0, 3, 'a', 'b', 'c', 0, 0, 0, 2, 'x', 'y', 100, 200},
			expected:       []interface{}{"abc", "xy"},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "parse sequence with 0 string elements",
			data:           []byte{0, 0, 0, 0, 100, 200},
			expected:       []interface{}{},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "parse sequence with 1 string element",
			data:           []byte{0, 0, 0, 1, 0, 0, 0, 4, 't', 'e', 's', 't', 100, 200},
			expected:       []interface{}{"test"},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "parse sequence with empty strings",
			data:           []byte{0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 100, 200},
			expected:       []interface{}{"", ""},
			expectedRemain: []byte{100, 200},
			expectError:    false,
		},
		{
			name:           "should return error when insufficient data for sequence length",
			data:           []byte{0, 0, 0},
			expected:       nil,
			expectedRemain: nil,
			expectError:    true,
		},
		{
			name:           "should return error when sequence length correct but first string data insufficient",
			data:           []byte{0, 0, 0, 1, 0, 0, 0, 5, 'h', 'e'},
			expected:       nil,
			expectedRemain: nil,
			expectError:    true,
		},
		{
			name:           "should return error when first string length correct but data insufficient",
			data:           []byte{0, 0, 0, 1, 0, 0, 0, 3},
			expected:       nil,
			expectedRemain: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stringType := typeref.NewStringType()
			sequenceType := typeref.NewSequence(stringType)

			result, remain, err := parseDataByType(tt.data, sequenceType, ast.Module{})

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, remain)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.expectedRemain, remain)
			}
		})
	}
}
