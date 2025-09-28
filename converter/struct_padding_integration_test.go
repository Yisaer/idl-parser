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
			0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x00, 0x00, // "Hello" + 终止符 + 1字节填充
			0x00, 0x00, // padding 字节 (2字节)

			// string B: "World"
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x57, 0x6F, 0x72, 0x6C, 0x64, 0x00, 0x00, // "World" + 终止符 + 1字节填充
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
		// innerStruct.field2: 长度字段(4字节) + UTF-8字符串内容(10字节) + padding(2字节)
		// fieldB: 长度字段(4字节) + UTF-8字符串内容(10字节) + 无padding
		testData := []byte{
			// fieldA: "OuterA" (16字节)
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x4F, 0x75, 0x74, 0x65, 0x72, 0x41, 0x00, // "OuterA" + 终止符
			0x00, 0x00, // padding 字节 (2字节)

			// innerStruct.field1: "Inner1" (16字节)
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x49, 0x6E, 0x6E, 0x65, 0x72, 0x31, 0x00, // "Inner1" + 终止符
			0x00, 0x00, // padding 字节 (2字节)

			// innerStruct.field2: "Inner2" (14字节)
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x49, 0x6E, 0x6E, 0x65, 0x72, 0x32, 0x00, // "Inner2" + 终止符
			0x00, 0x00, // innerStruct 的 padding 字节 (2字节)

			// fieldB: "OuterB" (14字节)
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

	t.Run("简单结构体_两个动态数组字段", func(t *testing.T) {
		// 测试场景：module m { struct arrayStruct { sequence<short> arrayA; sequence<short> arrayB; }}
		// arrayA 不是最后一个字段，需要 padding
		// arrayB 是最后一个字段，不需要 padding

		converter, err := NewIDLConverter(IDlConverterConfig{
			IsLittleEndian:    false,
			LengthFieldLength: 4,
			PaddingLength:     4,
			SchemaPath:        "testdata/test_array_padding.idl",
		})
		require.NoError(t, err)

		// 构造测试数据：
		// arrayA: 长度字段(4字节) + 3个short元素(6字节) = 10字节，10 % 4 = 2，需要2字节padding
		// arrayB: 长度字段(4字节) + 2个short元素(4字节) = 8字节，是最后一个字段，不需要padding
		testData := []byte{
			// arrayA: [1, 2, 3] (12字节)
			0x00, 0x00, 0x00, 0x03, // 长度字段 (3个元素)
			0x00, 0x01, 0x00, 0x02, 0x00, 0x03, // 三个short元素: 1, 2, 3
			0x00, 0x00, // padding 字节 (2字节)

			// arrayB: [4, 5]
			0x00, 0x00, 0x00, 0x02, // 长度字段 (2个元素) 4字节
			0x00, 0x04, 0x00, 0x05, // 两个short元素: 4, 5 4字节
		}

		result, err := converter.Decode("m.arrayStruct", testData)
		require.NoError(t, err)

		// 验证解析结果
		arrayA, ok := result["arrayA"].([]interface{})
		require.True(t, ok)
		require.Equal(t, 3, len(arrayA))
		require.Equal(t, int16(1), arrayA[0])
		require.Equal(t, int16(2), arrayA[1])
		require.Equal(t, int16(3), arrayA[2])

		arrayB, ok := result["arrayB"].([]interface{})
		require.True(t, ok)
		require.Equal(t, 2, len(arrayB))
		require.Equal(t, int16(4), arrayB[0])
		require.Equal(t, int16(5), arrayB[1])
	})

	t.Run("嵌套结构体_动态数组padding行为", func(t *testing.T) {
		// 测试场景：嵌套结构体中的动态数组 padding 行为
		// module m { struct outerArray { sequence<short> fieldA; innerArray innerStruct; sequence<short> fieldB; }}
		// innerArray: { sequence<short> field1; sequence<short> field2; }

		converter, err := NewIDLConverter(IDlConverterConfig{
			IsLittleEndian:    false,
			LengthFieldLength: 4,
			PaddingLength:     4,
			SchemaPath:        "testdata/test_array_padding.idl",
		})
		require.NoError(t, err)

		// 构造测试数据：
		// fieldA: 长度字段(4字节) + 3个short元素(6字节) = 10字节，10 % 4 = 2，需要2字节padding
		// innerStruct.field1: 长度字段(4字节) + 3个short元素(6字节) = 10字节，10 % 4 = 2，需要2字节padding
		// innerStruct.field2: 长度字段(4字节) + 2个short元素(4字节) = 8字节，是最后一个字段，不需要padding
		// fieldB: 长度字段(4字节) + 2个short元素(4字节) = 8字节，是最后一个字段，不需要padding
		testData := []byte{
			// fieldA: [10, 20, 30] (12字节)
			0x00, 0x00, 0x00, 0x03, // 长度字段 (3个元素)
			0x00, 0x0A, 0x00, 0x14, 0x00, 0x1E, // 三个short元素: 10, 20, 30
			0x00, 0x00, // padding 字节 (2字节)

			// innerStruct.field1: [40, 50, 60] (12字节)
			0x00, 0x00, 0x00, 0x03, // 长度字段 (3个元素)
			0x00, 0x28, 0x00, 0x32, 0x00, 0x3C, // 三个short元素: 40, 50, 60
			0x00, 0x00, // padding 字节 (2字节)

			// innerStruct.field2: [70, 80] (8字节)
			0x00, 0x00, 0x00, 0x02, // 长度字段 (2个元素)
			0x00, 0x46, 0x00, 0x50, // 两个short元素: 70, 80

			// fieldB: [90, 100] (8字节)
			0x00, 0x00, 0x00, 0x02, // 长度字段 (2个元素)
			0x00, 0x5A, 0x00, 0x64, // 两个short元素: 90, 100
		}

		result, err := converter.Decode("m.outerArray", testData)
		require.NoError(t, err)

		// 验证解析结果
		fieldA, ok := result["fieldA"].([]interface{})
		require.True(t, ok)
		require.Equal(t, 3, len(fieldA))
		require.Equal(t, int16(10), fieldA[0])
		require.Equal(t, int16(20), fieldA[1])
		require.Equal(t, int16(30), fieldA[2])

		fieldB, ok := result["fieldB"].([]interface{})
		require.True(t, ok)
		require.Equal(t, 2, len(fieldB))
		require.Equal(t, int16(90), fieldB[0])
		require.Equal(t, int16(100), fieldB[1])

		// 验证嵌套结构体
		innerStruct, ok := result["innerStruct"].(map[string]interface{})
		require.True(t, ok)

		field1, ok := innerStruct["field1"].([]interface{})
		require.True(t, ok)
		require.Equal(t, 3, len(field1))
		require.Equal(t, int16(40), field1[0])
		require.Equal(t, int16(50), field1[1])
		require.Equal(t, int16(60), field1[2])

		field2, ok := innerStruct["field2"].([]interface{})
		require.True(t, ok)
		require.Equal(t, 2, len(field2))
		require.Equal(t, int16(70), field2[0])
		require.Equal(t, int16(80), field2[1])
	})
}
