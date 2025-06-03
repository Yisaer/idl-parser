package ast

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/bitset"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typ"
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
