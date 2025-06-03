package ast

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/annotation"
	"github.com/yisaer/idl-parser/ast/bitset"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

func TestParsing(t *testing.T) {
	code := `module spi {
    bitset idbits {
        bitfield<4> bid; // 4 bits for bus_id
        bitfield<12> cid;  // 12 bits for can_id
    };

	struct CANFrame {
		octet header;
		idbits id;
	};
}`
	result := Parse(code)
	require.Nil(t, result.Err)
	require.Equal(t, result.Output.Name, "spi")
	require.Equal(t, len(result.Output.Content), 2)
	first, ok := result.Output.Content[0].(bitset.BitSet)
	require.True(t, ok)
	require.Equal(t, first.Name, "idbits")
	require.Equal(t, len(first.Fields), 2)
	require.Equal(t, first.Fields[0].Name, "bid")
	require.Equal(t, first.Fields[0].Type.Width, uint8(4))
	require.Equal(t, first.Fields[1].Name, "cid")
	require.Equal(t, first.Fields[1].Type.Width, uint8(12))

	second, ok := result.Output.Content[1].(struct_type.Struct)
	require.True(t, ok)
	require.Equal(t, second.Name, "CANFrame")
	require.Equal(t, len(second.Fields), 2)
	require.Equal(t, second.Fields[0].Name, "header")
	require.Equal(t, second.Fields[0].Type.TypeRefType(), typ.OctetType)
	require.Equal(t, second.Fields[1].Name, "id")
	require.Equal(t, second.Fields[1].Type.TypeRefType(), typ.SelfDefinedTypeType)
	require.Equal(t, second.Fields[1].Type.TypeName(), "idbits")
}

func TestParseModule(t *testing.T) {
	tests := []struct {
		input    string
		expected Module
	}{
		{
			input: `module spi {
						bitset idbits {
							bitfield<4> bid; // 4 bits for bus_id
						};

						struct CANFrame {
							@format octet header;
							@format(a=b) idbits id;
						};
					}`,
			expected: Module{
				Name: "spi",
				Content: []ModuleContent{
					bitset.BitSet{
						Name: "idbits",
						Fields: []bitset.Field{
							{
								Name: "bid",
								Type: typeref.BitFieldType{Width: uint8(4)},
							},
						},
					},
					struct_type.Struct{
						Name: "CANFrame",
						Fields: []struct_type.Field{
							{
								Name:        "header",
								Annotations: []annotation.Annotation{{Name: "format"}},
								Type:        typeref.OctetType{},
							},
							{
								Name:        "id",
								Annotations: []annotation.Annotation{{Name: "format", Values: map[string]string{"a": "b"}}},
								Type:        typeref.TypeName{Name: "idbits"},
							},
						},
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
