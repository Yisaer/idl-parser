package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestStructPaddingIntegration 测试结构体字段间的 padding 行为
func TestStructPaddingIntegration(t *testing.T) {
	t.Run("简单结构体_两个动态字符串字段", func(t *testing.T) {
		// 测试场景：module m { struct s { string A; string B; }}
		// string A 不是最后一个字段，需要 padding
		// string B 是最后一个字段，不需要 padding

		converter, err := NewIDLConverter(IDlConverterConfig{
			IsLittleEndian:    false,
			LengthFieldLength: 4,
			PaddingLength:     4,
			SchemaPath:        "testdata/test_string_padding.idl",
		})
		require.NoError(t, err)

		// 构造测试数据：
		// string A: 长度字段(4字节) + UTF-8字符串内容(10字节) + padding(2字节)
		// string B: 长度字段(4字节) + UTF-8字符串内容(10字节) + 无padding
		testData := []byte{
			// string A: "Hello"
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00, // "Hello" + 终止符

			// string B: "World"
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x57, 0x6F, 0x72, 0x6C, 0x64, 0x00, // "World" + 终止符
		}

		result, err := converter.Decode("m.s", testData)
		require.NoError(t, err)

		// 验证解析结果
		require.Equal(t, "Hello", result["A"])
		require.Equal(t, "World", result["B"])
	})

	t.Run("嵌套结构体_padding行为", func(t *testing.T) {
		// 测试场景：嵌套结构体中的 padding 行为
		// module m { struct outer { string fieldA; inner innerStruct; string fieldB; }}
		// inner struct: { string field1; string field2; }

		converter, err := NewIDLConverter(IDlConverterConfig{
			IsLittleEndian:    false,
			LengthFieldLength: 4,
			PaddingLength:     4,
			SchemaPath:        "testdata/test_string_padding.idl",
		})
		require.NoError(t, err)

		// 构造测试数据：
		// fieldA: 长度字段(4字节) + UTF-8字符串内容(10字节) + padding(2字节)
		// innerStruct.field1: 长度字段(4字节) + UTF-8字符串内容(10字节) + padding(2字节)
		// innerStruct.field2: 长度字段(4字节) + UTF-8字符串内容(10字节) + 无padding
		// fieldB: 长度字段(4字节) + UTF-8字符串内容(10字节) + 无padding
		testData := []byte{
			// fieldA: "OuterA"
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x4F, 0x75, 0x74, 0x65, 0x72, 0x41, 0x00, // "OuterA" + 终止符

			// innerStruct.field1: "Inner1"
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x49, 0x6E, 0x6E, 0x65, 0x72, 0x31, 0x00, // "Inner1" + 终止符

			// innerStruct.field2: "Inner2"
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x49, 0x6E, 0x6E, 0x65, 0x72, 0x32, 0x00, // "Inner2" + 终止符

			// fieldB: "OuterB"
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x4F, 0x75, 0x74, 0x65, 0x72, 0x42, 0x00, // "OuterB" + 终止符
		}

		result, err := converter.Decode("m.outer", testData)
		require.NoError(t, err)

		// 验证解析结果
		require.Equal(t, "OuterA", result["fieldA"])
		require.Equal(t, "OuterB", result["fieldB"])

		// 验证嵌套结构体
		innerStruct, ok := result["innerStruct"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, "Inner1", innerStruct["field1"])
		require.Equal(t, "Inner2", innerStruct["field2"])
	})
}
