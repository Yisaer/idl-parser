package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestParseBytesToFixedLengthString 测试 IDLConverter.ParseBytesToFixedLengthString 方法的定长字符串解析
func TestParseBytesToFixedLengthString(t *testing.T) {
	// 创建测试用的 IDLConverter
	converter := &IDLConverter{
		IDlConverterConfig: IDlConverterConfig{
			IsLittleEndian: false,
			LengthFieldLength: 4,
		},
	}
	t.Run("UTF-8定长字符串_无填充", func(t *testing.T) {
		// UTF-8 BOM: EF BB BF
		// 内容: "Hello" (UTF-8)
		// 终止符: 00
		// 固定长度: 9字节 (3 BOM + 5 内容 + 1 终止符)
		testData := []byte{
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x65, 0x6C, 0x6C, 0x6F, // "Hello"
			0x00, // UTF-8 终止符
		}

		result, remained, err := converter.ParseBytesToFixedLengthString(testData, 9)
		require.NoError(t, err)
		require.Equal(t, "Hello", result)
		require.Empty(t, remained)
	})

	t.Run("UTF-8定长字符串_有填充", func(t *testing.T) {
		// UTF-8 BOM: EF BB BF
		// 内容: "Hi" (UTF-8)
		// 终止符: 00
		// 填充: 00 00 00 00
		// 固定长度: 10字节 (3 BOM + 2 内容 + 1 终止符 + 4 填充)
		testData := []byte{
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x69, // "Hi"
			0x00,                   // UTF-8 终止符
			0x00, 0x00, 0x00, 0x00, // 填充字节
		}

		result, remained, err := converter.ParseBytesToFixedLengthString(testData, 10)
		require.NoError(t, err)
		require.Equal(t, "Hi", result)
		require.Empty(t, remained)
	})

	t.Run("UTF-16BE定长字符串_无填充", func(t *testing.T) {
		// UTF-16BE BOM: FE FF
		// 内容: "Go" (UTF-16BE: 00 47 00 6F)
		// 终止符: 00 00
		// 固定长度: 8字节 (2 BOM + 4 内容 + 2 终止符)
		testData := []byte{
			0xFE, 0xFF, // UTF-16BE BOM
			0x00, 0x47, // 'G'
			0x00, 0x6F, // 'o'
			0x00, 0x00, // UTF-16 终止符
		}

		result, remained, err := converter.ParseBytesToFixedLengthString(testData, 8)
		require.NoError(t, err)
		require.Equal(t, "Go", result)
		require.Empty(t, remained)
	})

	t.Run("UTF-16BE定长字符串_有填充", func(t *testing.T) {
		// UTF-16BE BOM: FE FF
		// 内容: "Hi" (UTF-16BE: 00 48 00 69)
		// 终止符: 00 00
		// 填充: 00 00 00 00
		// 固定长度: 10字节 (2 BOM + 4 内容 + 2 终止符 + 2 填充)
		testData := []byte{
			0xFE, 0xFF, // UTF-16BE BOM
			0x00, 0x48, // 'H'
			0x00, 0x69, // 'i'
			0x00, 0x00, // UTF-16 终止符
			0x00, 0x00, // 填充字节
		}

		result, remained, err := converter.ParseBytesToFixedLengthString(testData, 10)
		require.NoError(t, err)
		require.Equal(t, "Hi", result)
		require.Empty(t, remained)
	})

	t.Run("UTF-16LE定长字符串_无填充", func(t *testing.T) {
		// UTF-16LE BOM: FF FE
		// 内容: "Go" (UTF-16LE: 47 00 6F 00)
		// 终止符: 00 00
		// 固定长度: 8字节 (2 BOM + 4 内容 + 2 终止符)
		testData := []byte{
			0xFF, 0xFE, // UTF-16LE BOM
			0x47, 0x00, // 'G'
			0x6F, 0x00, // 'o'
			0x00, 0x00, // UTF-16 终止符
		}

		result, remained, err := converter.ParseBytesToFixedLengthString(testData, 8)
		require.NoError(t, err)
		require.Equal(t, "Go", result)
		require.Empty(t, remained)
	})

	t.Run("UTF-16LE定长字符串_有填充", func(t *testing.T) {
		// UTF-16LE BOM: FF FE
		// 内容: "Hi" (UTF-16LE: 48 00 69 00)
		// 终止符: 00 00
		// 填充: 00 00 00 00
		// 固定长度: 10字节 (2 BOM + 4 内容 + 2 终止符 + 2 填充)
		testData := []byte{
			0xFF, 0xFE, // UTF-16LE BOM
			0x48, 0x00, // 'H'
			0x69, 0x00, // 'i'
			0x00, 0x00, // UTF-16 终止符
			0x00, 0x00, // 填充字节
		}

		result, remained, err := converter.ParseBytesToFixedLengthString(testData, 10)
		require.NoError(t, err)
		require.Equal(t, "Hi", result)
		require.Empty(t, remained)
	})

	t.Run("带剩余数据的测试", func(t *testing.T) {
		// UTF-8 定长字符串 + 额外数据
		testData := []byte{
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x65, 0x6C, 0x6C, 0x6F, // "Hello"
			0x00,             // UTF-8 终止符
			0xFF, 0xFF, 0xFF, // 剩余数据
		}

		result, remained, err := converter.ParseBytesToFixedLengthString(testData, 9)
		require.NoError(t, err)
		require.Equal(t, "Hello", result)
		require.Equal(t, []byte{0xFF, 0xFF, 0xFF}, remained)
	})

	t.Run("中文字符测试", func(t *testing.T) {
		// UTF-8 中文字符: "你好"
		// UTF-8 BOM: EF BB BF
		// 内容: "你好" (UTF-8: E4 BD A0 E5 A5 BD)
		// 终止符: 00
		// 固定长度: 10字节 (3 BOM + 6 内容 + 1 终止符)
		testData := []byte{
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0xE4, 0xBD, 0xA0, // "你" (UTF-8)
			0xE5, 0xA5, 0xBD, // "好" (UTF-8)
			0x00, // UTF-8 终止符
		}

		result, remained, err := converter.ParseBytesToFixedLengthString(testData, 10)
		require.NoError(t, err)
		require.Equal(t, "你好", result)
		require.Empty(t, remained)
	})

	t.Run("错误情况_无BOM", func(t *testing.T) {
		// 没有BOM的定长字符串应该失败
		testData := []byte{
			0x48, 0x65, 0x6C, 0x6C, 0x6F, // "Hello"
			0x00, // UTF-8 终止符
		}

		_, _, err := converter.ParseBytesToFixedLengthString(testData, 6)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown encoding: no valid BOM found")
	})

	t.Run("错误情况_数据不足", func(t *testing.T) {
		// 数据长度小于固定长度
		testData := []byte{
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x65, // "He"
		}

		_, _, err := converter.ParseBytesToFixedLengthString(testData, 10)
		require.Error(t, err)
		require.Contains(t, err.Error(), "insufficient data")
	})

	t.Run("错误情况_无终止符", func(t *testing.T) {
		// 没有终止符的定长字符串应该失败
		testData := []byte{
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x65, 0x6C, 0x6C, 0x6F, // "Hello"
			// 缺少终止符
		}

		_, _, err := converter.ParseBytesToFixedLengthString(testData, 8)
		require.Error(t, err)
		require.Contains(t, err.Error(), "terminator not found")
	})
}

// TestParseBytesToDynamicString 测试 IDLConverter.ParseBytesToDynamicStringWithConfigurableLength 方法的动态字符串解析
func TestParseBytesToDynamicString(t *testing.T) {
	// 创建测试用的 IDLConverter
	converter := &IDLConverter{
		IDlConverterConfig: IDlConverterConfig{
			IsLittleEndian: false,
			LengthFieldLength: 4,
		},
	}
	t.Run("UTF-8动态字符串_带BOM", func(t *testing.T) {
		// 长度字段: 00 00 00 08 (8字节，包含BOM和终止符)
		// UTF-8 BOM: EF BB BF
		// 内容: "Test" (UTF-8)
		// 终止符: 00
		testData := []byte{
			0x00, 0x00, 0x00, 0x08, // 长度字段 (8字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x54, 0x65, 0x73, 0x74, // "Test"
			0x00, // UTF-8 终止符
		}

		result, remained, err := converter.ParseBytesToDynamicString(testData)
		require.NoError(t, err)
		require.Equal(t, "Test", result)
		require.Empty(t, remained)
	})

	t.Run("UTF-16BE动态字符串_带BOM", func(t *testing.T) {
		// 长度字段: 00 00 00 08 (8字节，包含BOM和终止符)
		// UTF-16BE BOM: FE FF
		// 内容: "Go" (UTF-16BE: 00 47 00 6F)
		// 终止符: 00 00
		testData := []byte{
			0x00, 0x00, 0x00, 0x08, // 长度字段 (8字节)
			0xFE, 0xFF, // UTF-16BE BOM
			0x00, 0x47, // 'G'
			0x00, 0x6F, // 'o'
			0x00, 0x00, // UTF-16 终止符
		}

		result, remained, err := converter.ParseBytesToDynamicString(testData)
		require.NoError(t, err)
		require.Equal(t, "Go", result)
		require.Empty(t, remained)
	})

	t.Run("UTF-16LE动态字符串_带BOM", func(t *testing.T) {
		// 长度字段: 00 00 00 08 (8字节，包含BOM和终止符)
		// UTF-16LE BOM: FF FE
		// 内容: "Go" (UTF-16LE: 47 00 6F 00)
		// 终止符: 00 00
		testData := []byte{
			0x00, 0x00, 0x00, 0x08, // 长度字段 (8字节)
			0xFF, 0xFE, // UTF-16LE BOM
			0x47, 0x00, // 'G'
			0x6F, 0x00, // 'o'
			0x00, 0x00, // UTF-16 终止符
		}

		result, remained, err := converter.ParseBytesToDynamicString(testData)
		require.NoError(t, err)
		require.Equal(t, "Go", result)
		require.Empty(t, remained)
	})

	t.Run("中文字符测试", func(t *testing.T) {
		// UTF-8 中文字符: "你好"
		// 长度字段: 00 00 00 0A (10字节，包含BOM和终止符)
		// UTF-8 BOM: EF BB BF
		// 内容: "你好" (UTF-8: E4 BD A0 E5 A5 BD)
		// 终止符: 00
		testData := []byte{
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0xE4, 0xBD, 0xA0, // "你" (UTF-8)
			0xE5, 0xA5, 0xBD, // "好" (UTF-8)
			0x00, // UTF-8 终止符
		}

		result, remained, err := converter.ParseBytesToDynamicString(testData)
		require.NoError(t, err)
		require.Equal(t, "你好", result)
		require.Empty(t, remained)
	})

	t.Run("带填充的动态字符串", func(t *testing.T) {
		// UTF-8 动态字符串，带填充
		// 长度字段: 00 00 00 0A (10字节，包含BOM、内容和填充)
		// UTF-8 BOM: EF BB BF
		// 内容: "Hi" (UTF-8)
		// 终止符: 00
		// 填充: 00 00 00 00
		testData := []byte{
			0x00, 0x00, 0x00, 0x0A, // 长度字段 (10字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x48, 0x69, // "Hi"
			0x00,                   // UTF-8 终止符
			0x00, 0x00, 0x00, 0x00, // 填充字节
		}

		result, remained, err := converter.ParseBytesToDynamicString(testData)
		require.NoError(t, err)
		require.Equal(t, "Hi", result)
		require.Empty(t, remained)
	})

	t.Run("空字符串测试_应该失败", func(t *testing.T) {
		// 空字符串（只有终止符，没有BOM）
		testData := []byte{
			0x00, 0x00, 0x00, 0x01, // 长度字段 (1字节)
			0x00, // UTF-8 终止符
		}

		_, _, err := converter.ParseBytesToDynamicString(testData)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown encoding: no valid BOM found")
	})

	t.Run("带BOM的空字符串测试", func(t *testing.T) {
		// 带BOM的空字符串（BOM + 终止符）
		testData := []byte{
			0x00, 0x00, 0x00, 0x04, // 长度字段 (4字节)
			0xEF, 0xBB, 0xBF, // UTF-8 BOM
			0x00, // UTF-8 终止符
		}

		result, remained, err := converter.ParseBytesToDynamicString(testData)
		require.NoError(t, err)
		require.Equal(t, "", result)
		require.Empty(t, remained)
	})

	t.Run("错误情况_无BOM的动态字符串", func(t *testing.T) {
		// 没有BOM的动态字符串应该失败
		testData := []byte{
			0x00, 0x00, 0x00, 0x06, // 长度字段 (6字节)
			0x48, 0x65, 0x6C, 0x6C, 0x6F, // "Hello"
			0x00, // UTF-8 终止符
		}

		_, _, err := converter.ParseBytesToDynamicString(testData)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown encoding: no valid BOM found")
	})

	t.Run("错误情况_无效BOM", func(t *testing.T) {
		// 无效的BOM应该失败
		testData := []byte{
			0x00, 0x00, 0x00, 0x08, // 长度字段 (8字节)
			0xAA, 0xBB, 0xCC, // 无效的BOM
			0x48, 0x65, 0x6C, 0x6C, 0x6F, // "Hello"
			0x00, // UTF-8 终止符
		}

		_, _, err := converter.ParseBytesToDynamicString(testData)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown encoding: no valid BOM found")
	})
}
