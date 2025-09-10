package typeref

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/typ"
)

func TestParseTypeName(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedName   string
		expectedSize   int
		expectedRemain string
		expectError    bool
	}{
		{
			name:           "parse plain type name",
			input:          "idbits",
			expectedName:   "idbits",
			expectedSize:   0,
			expectedRemain: "",
			expectError:    false,
		},
		{
			name:           "parse type name with remaining text",
			input:          "idbits id;",
			expectedName:   "idbits",
			expectedSize:   0,
			expectedRemain: " id;",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTypeName(tt.input)

			if tt.expectError {
				require.Error(t, result.Err)
			} else {
				if result.Err != nil {
					t.Errorf("Unexpected error: %v", result.Err)
					return
				}
				require.Equal(t, tt.expectedName, result.Output.Name)
				require.Equal(t, tt.expectedName, result.Output.SelfType)
				require.Equal(t, tt.expectedSize, result.Output.ArraySize)
				require.Equal(t, tt.expectedRemain, result.Remaining)
			}
		})
	}
}

func TestParseTypeNameArrayViaTypeRef(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedName   string
		expectedSize   int
		expectedRemain string
		expectError    bool
	}{
		{
			name:           "parse type name with array size",
			input:          "idbits[8]",
			expectedName:   "idbits",
			expectedSize:   8,
			expectedRemain: "",
			expectError:    false,
		},
		{
			name:           "parse type name with large array size",
			input:          "MyType[255]",
			expectedName:   "MyType",
			expectedSize:   255,
			expectedRemain: "",
			expectError:    false,
		},
		{
			name:           "parse type name with single digit array size",
			input:          "Test[1]",
			expectedName:   "Test",
			expectedSize:   1,
			expectedRemain: "",
			expectError:    false,
		},
		{
			name:           "parse type name array with remaining text",
			input:          "idbits[8] some other text",
			expectedName:   "idbits",
			expectedSize:   8,
			expectedRemain: " some other text",
			expectError:    false,
		},
		{
			name:           "parse invalid array format - missing closing bracket",
			input:          "idbits[8",
			expectedName:   "",
			expectedSize:   0,
			expectedRemain: "",
			expectError:    true,
		},
		{
			name:           "parse invalid array format - non-digit in brackets",
			input:          "idbits[abc]",
			expectedName:   "",
			expectedSize:   0,
			expectedRemain: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTypeRef(tt.input)

			if tt.expectError {
				require.Error(t, result.Err)
				return
			}

			if result.Err != nil {
				t.Errorf("Unexpected error: %v", result.Err)
				return
			}
			require.Equal(t, typ.ArrayType, result.Output.TypeRefType())

			arrayType, ok := result.Output.(ArrayType)
			require.True(t, ok)

			typeName, ok := arrayType.InnerType.(TypeName)
			require.True(t, ok)

			require.Equal(t, tt.expectedName, typeName.Name)
			require.Equal(t, tt.expectedSize, arrayType.ArraySize)
			require.Equal(t, tt.expectedRemain, result.Remaining)
		})
	}
}

func TestNewTypeName(t *testing.T) {
	typeName := NewTypeName("MyType")
	require.Equal(t, "MyType", typeName.Name)
	require.Equal(t, "MyType", typeName.SelfType)
	require.Equal(t, 0, typeName.ArraySize)
}

func TestNewTypeNameArrayType(t *testing.T) {
	typeName := NewTypeNameArrayType("MyType", 8)
	require.Equal(t, "MyType", typeName.Name)
	require.Equal(t, "MyType", typeName.SelfType)
	require.Equal(t, 8, typeName.ArraySize)
}
