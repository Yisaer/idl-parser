package bitset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBitSetField(t *testing.T) {
	code := `bitfield<1> a; // 1bit
	`
	result := parseField(code)
	assert.Equal(t, result.Output.Name, "a")
	assert.Equal(t, result.Output.Type.Width, 1)
}

func TestParseBitSet(t *testing.T) {
	code := `bitset S {
	bitfield<1> a; // 1bit
	bitfield<4> b; // 4bit
	}
	`
	result := Parse(code)
	assert.Equal(t, result.Output.Name, "S")
	assert.Equal(t, len(result.Output.Fields), 2)
}
