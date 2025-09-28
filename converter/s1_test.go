package converter

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestS1(t *testing.T) {
	c, err := NewIDLConverter(IDlConverterConfig{
		IsLittleEndian:    false,
		LengthFieldLength: 4,
		PaddingLength:     4,
		SchemaPath:        "./testdata/s1_test.idl",
	})
	require.NoError(t, err)
	hexStr := "0000000200000090efbbbfe4b8ade69687205749464900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000022efbbbf456e676c697368205749464900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000380000004e"
	data, err := hex.DecodeString(hexStr)
	require.NoError(t, err)
	v, err := c.Decode("m.WifiApList", data)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"wiFiApNum": int32(2),
		"wiFiApArray": []interface{}{
			map[string]interface{}{
				"wiFiApName":       "中文 WIFI",
				"wiFiApStrength":   int32(12),
				"wiFiApEncryption": int32(34),
			},
			map[string]interface{}{
				"wiFiApName":       "English WIFI",
				"wiFiApStrength":   int32(56),
				"wiFiApEncryption": int32(78),
			},
		},
	}, v)
}
