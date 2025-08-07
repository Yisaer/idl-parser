package converter

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

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
			result, remain, err := parseDataByType(tt.data, octetType)

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

			result, remain, err := parseDataByType(tt.data, sequenceType)

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

	result, remain, err := parseDataByType([]byte{1, 2, 3}, mockType)

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
