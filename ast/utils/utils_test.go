package utils

import (
	"testing"

	"github.com/oleiade/gomme"
	"github.com/stretchr/testify/require"
)

func TestInEmpty(t *testing.T) {
	code := `// ;`
	result := InEmpty(gomme.Token[string](";"))(code)
	require.NotNil(t, result.Err)

	code = ` ; // xxx`
	result = InEmpty(gomme.Token[string](";"))(code)
	require.Equal(t, result.Output, ";")
}
