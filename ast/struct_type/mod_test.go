package struct_type

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast/typ"
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
