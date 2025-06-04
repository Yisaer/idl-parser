package ast

import (
	"encoding/json"
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
								Type: typeref.BitFieldType{Width: uint8(4), SelfType: "bitfield"},
							},
						},
					},
					struct_type.Struct{
						Name: "CANFrame",
						Fields: []struct_type.Field{
							{
								Name:        "header",
								Annotations: []annotation.Annotation{{Name: "format"}},
								Type:        typeref.OctetType{SelfType: "octet"},
							},
							{
								Name:        "id",
								Annotations: []annotation.Annotation{{Name: "format", Values: map[string]string{"a": "b"}}},
								Type:        typeref.TypeName{Name: "idbits", SelfType: "idbits"},
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

func TestModuleJson(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: `module spi {
						bitset IdBits {
							bitfield<4> bid; 
							bitfield<12> cid;
						};

						bitset LBits {
							bitfield<1> isUpdate;      
							bitfield<7> plen;
						};

						struct CANFrame {
							octet header;
							IdBits id;
							LBits l;
							sequence<octet> payload;
						};
						struct SPI {
							unsigned short header;           
							unsigned short plen;           
							octet counter;
							octet crc;
							@format (type=canpack,dbc=ab) @merge sequence<CANFrame> messages; 
						};
						struct parquet {
							unsigned long long timestamp; 
							@format (type=binpack) @merge sequence<SPI> packs;
						};
          			}`,
			expected: `{"name":"spi","content":[{"name":"IdBits","fields":[{"type":{"width":4,"self_type":"bitfield"},"name":"bid"},{"type":{"width":12,"self_type":"bitfield"},"name":"cid"}]},{"name":"LBits","fields":[{"type":{"width":1,"self_type":"bitfield"},"name":"isUpdate"},{"type":{"width":7,"self_type":"bitfield"},"name":"plen"}]},{"name":"CANFrame","fields":[{"type":{"self_type":"octet"},"name":"header"},{"type":{"self_type":"IdBits","name":"IdBits"},"name":"id"},{"type":{"self_type":"LBits","name":"LBits"},"name":"l"},{"type":{"self_type":"sequence","inner_type":{"self_type":"octet"}},"name":"payload"}]},{"name":"SPI","fields":[{"type":{"self_type":"unsigned short"},"name":"header"},{"type":{"self_type":"unsigned short"},"name":"plen"},{"type":{"self_type":"octet"},"name":"counter"},{"type":{"self_type":"octet"},"name":"crc"},{"annotations":[{"name":"format","values":{"dbc":"ab","type":"canpack"}},{"name":"merge"}],"type":{"self_type":"sequence","inner_type":{"self_type":"CANFrame","name":"CANFrame"}},"name":"messages"}]},{"name":"parquet","fields":[{"type":{"self_type":"unsigned long long"},"name":"timestamp"},{"annotations":[{"name":"format","values":{"type":"binpack"}},{"name":"merge"}],"type":{"self_type":"sequence","inner_type":{"self_type":"SPI","name":"SPI"}},"name":"packs"}]}]}`,
		},
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
			expected: `{"name":"spi","content":[{"name":"idbits","fields":[{"type":{"width":4,"self_type":"bitfield"},"name":"bid"}]},{"name":"CANFrame","fields":[{"annotations":[{"name":"format"}],"type":{"self_type":"octet"},"name":"header"},{"annotations":[{"name":"format","values":{"a":"b"}}],"type":{"self_type":"idbits","name":"idbits"},"name":"id"}]}]}`,
		},
	}
	for _, test := range tests {
		result := Parse(test.input)
		require.Nil(t, result.Err)
		v, err := json.Marshal(result.Output)
		require.NoError(t, err)
		require.Equal(t, test.expected, string(v))
	}
}
