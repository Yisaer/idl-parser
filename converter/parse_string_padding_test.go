package converter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/typeref"
)

// TestParseStringPadding 测试 ParseString 方法中动态字符串的 padding 功能
func TestParseStringPadding(t *testing.T) {
	t.Run("非最后一个字段_需要padding", func(t *testing.T) {
		// 测试场景：动态字符串不是最后一个字段，需要 padding
		// PaddingLength = 4, 字符串数据长度为 8 字节（4字节长度字段 + 4字节字符串内容）
		// 总消费字节数 = 8, 8 % 4 = 0, 所以不需要额外 padding
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    false,
				LengthFieldLength: 4,
				PaddingLength:     4,
			},
			TotalBytes: 20, // 设置总字节数用于 padding 计算
		}

		// 测试数据：长度字段(4字节) + UTF-8字符串内容(4字节) + 后续数据
		testData := []byte{
			0x00, 0x00, 0x00, 0x04, // 长度字段 (4字节)
			0xEF, 0xBB, 0xBF, 0x00, // UTF-8 BOM + 终止符 (4字节)
			0xFF, 0xFF, 0xFF, 0xFF, // 后续数据
		}

		result, remained, err := converter.ParseString(testData, typeref.StringType{}, ast.Module{}, false)
		require.NoError(t, err)
		require.Equal(t, "", result) // 空字符串
		require.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF}, remained)
	})

	t.Run("非最后一个字段_需要padding_奇数长度", func(t *testing.T) {
		// 测试场景：动态字符串不是最后一个字段，字符串长度为奇数，需要 padding
		// PaddingLength = 4, 字符串数据长度为 5 字节（4字节长度字段 + 1字节字符串内容）
		// TotalBytes = 20, consumedBytes = 20 - 15 = 5, 5 % 4 = 1, 需要 3 字节 padding
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    false,
				LengthFieldLength: 4,
				PaddingLength:     4,
			},
			TotalBytes: 20,
		}

		// 测试数据：长度字段(4字节) + UTF-8字符串内容(5字节) + padding(3字节) + 后续数据
		testData := []byte{
			0x00, 0x00, 0x00, 0x05, // 长度字段 (5字节)
			0xEF, 0xBB, 0xBF, 0x48, 0x00, // UTF-8 BOM + "H" + 终止符 (5字节)
			0x00, 0x00, 0x00, // padding 字节
			0xFF, 0xFF, 0xFF, // 后续数据
		}

		result, remained, err := converter.ParseString(testData, typeref.StringType{}, ast.Module{}, false)
		require.NoError(t, err)
		require.Equal(t, "H", result)
		require.Equal(t, []byte{0x00, 0xFF, 0xFF, 0xFF}, remained)
	})

	t.Run("最后一个字段_不需要padding", func(t *testing.T) {
		// 测试场景：动态字符串是最后一个字段，不需要 padding
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    false,
				LengthFieldLength: 4,
				PaddingLength:     4,
			},
			TotalBytes: 20,
		}

		// 测试数据：长度字段(4字节) + UTF-8字符串内容(4字节)
		testData := []byte{
			0x00, 0x00, 0x00, 0x04, // 长度字段 (4字节)
			0xEF, 0xBB, 0xBF, 0x00, // UTF-8 BOM + 终止符 (4字节)
		}

		result, remained, err := converter.ParseString(testData, typeref.StringType{}, ast.Module{}, true)
		require.NoError(t, err)
		require.Equal(t, "", result) // 空字符串
		require.Empty(t, remained)
	})

	t.Run("PaddingLength为8字节", func(t *testing.T) {
		// 测试场景：PaddingLength = 8, 字符串数据长度为 6 字节
		// TotalBytes = 20, consumedBytes = 4 + 6 = 10, 10 % 8 = 2, 需要 6 字节 padding
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    false,
				LengthFieldLength: 4,
				PaddingLength:     8,
			},
			TotalBytes: 20,
		}

		// 测试数据：长度字段(4字节) + UTF-8字符串内容(6字节) + padding(6字节) + 后续数据
		testData := []byte{
			0x00, 0x00, 0x00, 0x06, // 长度字段 (6字节)
			0xEF, 0xBB, 0xBF, 0x48, 0x69, 0x00, // UTF-8 BOM + "Hi" + 终止符 (6字节)
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // padding 字节 (6字节)
			0xFF, 0xFF, // 后续数据
		}

		result, remained, err := converter.ParseString(testData, typeref.StringType{}, ast.Module{}, false)
		require.NoError(t, err)
		require.Equal(t, "Hi", result)
		require.Equal(t, []byte{0x00, 0x00, 0xFF, 0xFF}, remained)
	})

	t.Run("PaddingLength为2字节", func(t *testing.T) {
		// 测试场景：PaddingLength = 2, 字符串数据长度为 5 字节
		// TotalBytes = 20, consumedBytes = 20 - 7 = 13, 13 % 2 = 1, 需要 1 字节 padding
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    false,
				LengthFieldLength: 4,
				PaddingLength:     2,
			},
			TotalBytes: 20,
		}

		// 测试数据：长度字段(4字节) + UTF-8字符串内容(5字节) + padding(1字节) + 后续数据
		testData := []byte{
			0x00, 0x00, 0x00, 0x05, // 长度字段 (5字节)
			0xEF, 0xBB, 0xBF, 0x48, 0x00, // UTF-8 BOM + "H" + 终止符 (5字节)
			0x00,       // padding 字节
			0xFF, 0xFF, // 后续数据
		}

		result, remained, err := converter.ParseString(testData, typeref.StringType{}, ast.Module{}, false)
		require.NoError(t, err)
		require.Equal(t, "H", result)
		require.Equal(t, []byte{0xFF, 0xFF}, remained)
	})

	t.Run("有实际内容的动态字符串_padding", func(t *testing.T) {
		// 测试场景：有实际内容的动态字符串，需要 padding
		// TotalBytes = 20, 字符串消费 11 字节（4字节长度字段 + 7字节字符串内容）
		// consumedBytes = 4 + 7 = 11, 11 % 4 = 3, 需要 1 字节 padding
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    false,
				LengthFieldLength: 4,
				PaddingLength:     4,
			},
			TotalBytes: 20,
		}

		// 测试数据：长度字段(4字节) + UTF-8字符串内容(7字节) + padding(1字节) + 后续数据
		testData := []byte{
			0x00, 0x00, 0x00, 0x07, // 长度字段 (7字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x69, 0x00, // "Hi" + 终止符
			0x00,             // padding 字节
			0xFF, 0xFF, 0xFF, // 后续数据
		}

		result, remained, err := converter.ParseString(testData, typeref.StringType{}, ast.Module{}, false)
		require.NoError(t, err)
		require.Equal(t, "Hi", result)
		require.Empty(t, remained)
	})

	t.Run("padding长度超过剩余数据_应该失败", func(t *testing.T) {
		// 测试场景：需要的 padding 长度超过剩余数据长度
		// TotalBytes = 20, 字符串消费 5 字节（4字节长度字段 + 1字节字符串内容）
		// consumedBytes = 5, 5 % 8 = 5, 需要 3 字节 padding，但只有 1 字节剩余数据
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    false,
				LengthFieldLength: 4,
				PaddingLength:     8,
			},
			TotalBytes: 20,
		}

		// 测试数据：长度字段(4字节) + UTF-8字符串内容(4字节) + 只有1字节剩余数据
		// 但需要 3 字节 padding，会失败
		testData := []byte{
			0x00, 0x00, 0x00, 0x04, // 长度字段 (4字节)
			0xEF, 0xBB, 0xBF, 0x00, // UTF-8 BOM + 终止符 (4字节)
			0xFF, // 只有1字节剩余数据，但需要3字节padding
		}

		_, _, err := converter.ParseString(testData, typeref.StringType{}, ast.Module{}, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "padding length exceeds remained length")
	})

	t.Run("定长字符串_不需要padding", func(t *testing.T) {
		// 测试场景：定长字符串不需要 padding
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    false,
				LengthFieldLength: 4,
				PaddingLength:     4,
			},
			TotalBytes: 20,
		}

		// 创建定长字符串类型
		fixedLength := 4
		stringType := typeref.StringType{
			Length: &fixedLength,
		}

		// 测试数据：UTF-8字符串内容(4字节) + 后续数据
		testData := []byte{
			0xEF, 0xBB, 0xBF, 0x00, // UTF-8 BOM + 终止符 (4字节)
			0xFF, 0xFF, 0xFF, // 后续数据
		}

		result, remained, err := converter.ParseString(testData, stringType, ast.Module{}, false)
		require.NoError(t, err)
		require.Equal(t, "", result) // 空字符串
		require.Equal(t, []byte{0xFF, 0xFF, 0xFF}, remained)
	})

	t.Run("小端序_动态字符串_padding", func(t *testing.T) {
		// 测试场景：小端序的动态字符串 padding
		// TotalBytes = 20, 字符串消费 11 字节（4字节长度字段 + 7字节字符串内容）
		// consumedBytes = 4 + 7 = 11, 11 % 4 = 3, 需要 1 字节 padding
		converter := &IDLConverter{
			IDlConverterConfig: IDlConverterConfig{
				IsLittleEndian:    true,
				LengthFieldLength: 4,
				PaddingLength:     4,
			},
			TotalBytes: 20,
		}

		// 测试数据：长度字段(4字节，小端序) + UTF-8字符串内容(7字节) + padding(1字节) + 后续数据
		testData := []byte{
			0x07, 0x00, 0x00, 0x00, // 长度字段 (7字节，小端序)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x69, 0x00, // "Hi" + 终止符
			0x00,             // padding 字节
			0xFF, 0xFF, 0xFF, // 后续数据
		}

		result, remained, err := converter.ParseString(testData, typeref.StringType{}, ast.Module{}, false)
		require.NoError(t, err)
		require.Equal(t, "Hi", result)
		require.Empty(t, remained)
	})
}
