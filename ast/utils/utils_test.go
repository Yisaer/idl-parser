package utils

import (
	"testing"

	"github.com/oleiade/gomme"
	"github.com/stretchr/testify/assert"
)

func TestInEmpty(t *testing.T) {
	code := `// ;`
	result := InEmpty(gomme.Token[string](";"))(code)
	assert.NotNil(t, result.Err)

	code = ` ; // xxx`
	result = InEmpty(gomme.Token[string](";"))(code)
	assert.Equal(t, result.Output, ";")
}
