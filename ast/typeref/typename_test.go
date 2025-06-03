package typeref

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTypeName(t *testing.T) {
	code := "idbits id;"
	result := ParseTypeName(code)
	require.Equal(t, "idbits", result.Output.Name)
}
