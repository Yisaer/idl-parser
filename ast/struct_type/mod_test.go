package struct_type

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/annotation"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

func TestParseStruct(t *testing.T) {
	code := `struct AB {
	  octet header;
      long h2;
      unsigned long h3;
      unsigned long long h4;
	}
	`
	result := Parse(code)
	require.Nil(t, result.Err)
	require.Equal(t, result.Output.Name, "AB")
	require.Equal(t, result.Output.Fields[0].Name, "header")
	require.Equal(t, result.Output.Fields[0].Type.TypeRefType(), typ.OctetType)
	require.Equal(t, result.Output.Fields[1].Name, "h2")
	require.Equal(t, result.Output.Fields[1].Type.TypeRefType(), typ.LongType)
	require.Equal(t, result.Output.Fields[2].Name, "h3")
	require.Equal(t, result.Output.Fields[2].Type.TypeRefType(), typ.UnsignedLongType)
	require.Equal(t, result.Output.Fields[3].Name, "h4")
	require.Equal(t, result.Output.Fields[3].Type.TypeRefType(), typ.UnsignedLongLongType)
}

func TestParseStructAnnotation(t *testing.T) {
	code := `struct AB {
	  @format octet header;
	}`
	result := Parse(code)
	require.Nil(t, result.Err)
	require.Equal(t, result.Output.Name, "AB")
	require.Equal(t, result.Output.Fields[0].Name, "header")
	require.Equal(t, result.Output.Fields[0].Type.TypeRefType(), typ.OctetType)
	require.Len(t, result.Output.Fields[0].Annotations, 1)
	require.Equal(t, result.Output.Fields[0].Annotations[0].Name, "format")
}

func TestStructFieldAnnotations(t *testing.T) {
	tests := []struct {
		input    string
		expected Struct
	}{
		{
			input: `struct AB {
						@format(a=b) octet header;
						long h2;
					}`,
			expected: Struct{
				Name: "AB",
				Fields: []Field{
					{
						Name:        "header",
						Type:        typeref.OctetType{SelfType: "octet"},
						Annotations: []annotation.Annotation{{Name: "format", Values: map[string]string{"a": "b"}}},
					},
					{
						Name: "h2",
						Type: typeref.LongType{SelfType: "long"},
					},
				},
			},
		},
		{
			input: `struct AB {
						@format octet header;
					}`,
			expected: Struct{
				Name: "AB",
				Fields: []Field{
					{
						Name:        "header",
						Type:        typeref.OctetType{SelfType: "octet"},
						Annotations: []annotation.Annotation{{Name: "format"}},
					},
				},
			},
		},
	}

	for _, test := range tests {
		result := Parse(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestParseStructSequence(t *testing.T) {
	code := `struct AB {
	  sequence<octet> payload;
	}`
	result := Parse(code)
	require.Nil(t, result.Err)
}
