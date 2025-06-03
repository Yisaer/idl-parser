package typeref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTypeName(t *testing.T) {
	code := "idbits id;"
	result := ParseTypeName(code)
	assert.Equal(t, "idbits", result.Output.Name)
}
