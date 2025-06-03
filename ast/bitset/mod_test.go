package bitset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBitSetField(t *testing.T) {
	code := `bitfield<1> a; // 1bit
	`
	result := parseField(code)
	require.Equal(t, result.Output.Name, "a")
	require.Equal(t, result.Output.Type.Width, uint8(1))
}

func TestParseBitSet(t *testing.T) {
	code := `bitset S {
	bitfield<1> a; // 1bit
	bitfield<4> b; // 4bit
	}
	`
	result := Parse(code)
	require.Equal(t, result.Output.Name, "S")
	require.Equal(t, len(result.Output.Fields), 2)
}
