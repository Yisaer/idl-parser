package converter

import (
	"encoding/binary"
	"errors"
	"unicode/utf16"
)

// StringEncoding 字符串编码类型
type StringEncoding int

const (
	EncodingUnknown StringEncoding = iota
	EncodingUTF8
	EncodingUTF16BE
	EncodingUTF16LE
)

// detectEncodingFromBOM 通过BOM探测编码
func detectEncodingFromBOM(data []byte) (StringEncoding, int) {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return EncodingUTF8, 3
	}
	if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
		return EncodingUTF16BE, 2
	}
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		return EncodingUTF16LE, 2
	}
	return EncodingUnknown, 0
}

// getTerminatorLength 获取终止符长度
func getTerminatorLength(encoding StringEncoding) int {
	switch encoding {
	case EncodingUTF8:
		return 1
	case EncodingUTF16BE, EncodingUTF16LE:
		return 2
	default:
		return 0
	}
}

// findTerminatorPosition 查找终止符的位置
func findTerminatorPosition(data []byte, encoding StringEncoding) int {
	terminatorLen := getTerminatorLength(encoding)
	if terminatorLen == 0 {
		return -1
	}

	switch encoding {
	case EncodingUTF8:
		for i := 0; i <= len(data)-terminatorLen; i++ {
			if data[i] == 0x00 {
				return i
			}
		}
	case EncodingUTF16BE, EncodingUTF16LE:
		for i := 0; i <= len(data)-terminatorLen; i += 2 {
			if data[i] == 0x00 && data[i+1] == 0x00 {
				return i
			}
		}
	}
	return -1
}

// convertBytesToString 将字节数据按编码转换为字符串
func convertBytesToString(data []byte, encoding StringEncoding) (string, error) {
	switch encoding {
	case EncodingUTF8:
		return string(data), nil
	case EncodingUTF16BE:
		return convertUTF16BytesToString(data, binary.BigEndian)
	case EncodingUTF16LE:
		return convertUTF16BytesToString(data, binary.LittleEndian)
	default:
		return "", errors.New("unsupported encoding")
	}
}

// convertUTF16BytesToString 转换UTF-16字节到字符串
func convertUTF16BytesToString(data []byte, byteOrder binary.ByteOrder) (string, error) {
	if len(data)%2 != 0 {
		return "", errors.New("UTF-16 data length must be even")
	}
	u16s := make([]uint16, len(data)/2)
	for i := range u16s {
		u16s[i] = byteOrder.Uint16(data[2*i : 2*i+2])
	}
	runes := utf16.Decode(u16s)
	return string(runes), nil
}

// parseBytesToFixedLengthString 解析定长字符串
// 根据SOME/IP协议，定长字符串必须以BOM开头，以终止符结束，长度固定
func parseBytesToFixedLengthString(data []byte, fixedLength int) (value string, remained []byte, err error) {
	if len(data) < fixedLength {
		return "", nil, errors.New("insufficient data for fixed-length string")
	}

	// 1. 检测BOM确定编码方式
	encoding, bomLen := detectEncodingFromBOM(data)
	if encoding == EncodingUnknown {
		return "", nil, errors.New("fixed-length string must start with BOM")
	}

	// 2. 验证固定长度
	if len(data) < fixedLength {
		return "", nil, errors.New("data length less than fixed length")
	}

	// 3. 提取字符串内容（去除BOM和终止符）
	stringData := data[:fixedLength]
	content, err := extractFixedStringContent(stringData, encoding, bomLen)
	if err != nil {
		return "", nil, err
	}

	// 4. 返回解析结果和剩余数据
	return content, data[fixedLength:], nil
}

// extractFixedStringContent 从定长字符串数据中提取内容
func extractFixedStringContent(data []byte, encoding StringEncoding, bomLen int) (string, error) {
	// 跳过BOM
	contentStart := bomLen
	terminatorLen := getTerminatorLength(encoding)

	// 检查是否有足够的空间包含终止符
	if len(data) < contentStart+terminatorLen {
		return "", errors.New("insufficient space for terminator")
	}

	// 查找终止符位置
	terminatorPos := findTerminatorPosition(data[contentStart:], encoding)
	if terminatorPos == -1 {
		return "", errors.New("terminator not found in fixed-length string")
	}

	// 提取内容（从BOM后到终止符前）
	contentEnd := contentStart + terminatorPos
	contentBytes := data[contentStart:contentEnd]

	// 移除填充的0x00字节
	contentBytes = removePadding(contentBytes, encoding)

	// 转换为字符串
	return convertBytesToString(contentBytes, encoding)
}

// removePadding 移除字符串末尾的填充字节
func removePadding(data []byte, encoding StringEncoding) []byte {
	if len(data) == 0 {
		return data
	}

	switch encoding {
	case EncodingUTF8:
		// UTF-8: 移除末尾的0x00字节
		for len(data) > 0 && data[len(data)-1] == 0x00 {
			data = data[:len(data)-1]
		}
	case EncodingUTF16BE, EncodingUTF16LE:
		// UTF-16: 移除末尾的0x00 0x00字节对
		for len(data) >= 2 && data[len(data)-2] == 0x00 && data[len(data)-1] == 0x00 {
			data = data[:len(data)-2]
		}
	}
	return data
}

func ParseBytesToDynamicString(data []byte) (value string, remained []byte, err error) {
	if len(data) < 4 {
		return "", nil, errors.New("insufficient data for dynamic string length field")
	}

	// 1. 读取4字节大端序长度字段
	lengthField := binary.BigEndian.Uint32(data[:4])
	stringDataStart := 4

	// 2. 检查是否有足够的数据
	if len(data) < int(stringDataStart)+int(lengthField) {
		return "", nil, errors.New("insufficient data for dynamic string content")
	}

	// 3. 提取字符串数据
	stringData := data[stringDataStart : stringDataStart+int(lengthField)]
	totalConsumed := stringDataStart + int(lengthField)

	// 4. 解析字符串内容
	content, err := parseDynamicStringContent(stringData)
	if err != nil {
		return "", nil, err
	}

	// 5. 返回解析结果和剩余数据
	return content, data[totalConsumed:], nil
}

// parseDynamicStringContent 解析动态字符串的内容部分
func parseDynamicStringContent(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty string data")
	}

	// 1. 检测BOM确定编码方式
	encoding, bomLen := detectEncodingFromBOM(data)
	if encoding == EncodingUnknown {
		// 无BOM，尝试推断编码
		encoding = inferEncodingFromContent(data)
		bomLen = 0
	}

	// 2. 提取字符串内容（去除BOM和终止符）
	contentStart := bomLen
	terminatorLen := getTerminatorLength(encoding)

	// 检查是否有足够的空间包含终止符
	if len(data) < contentStart+terminatorLen {
		return "", errors.New("insufficient space for terminator")
	}

	// 查找终止符位置
	terminatorPos := findTerminatorPosition(data[contentStart:], encoding)
	if terminatorPos == -1 {
		return "", errors.New("terminator not found in dynamic string")
	}

	// 提取内容（从BOM后到终止符前）
	contentEnd := contentStart + terminatorPos
	contentBytes := data[contentStart:contentEnd]

	// 移除填充的0x00字节
	contentBytes = removePadding(contentBytes, encoding)

	// 转换为字符串
	return convertBytesToString(contentBytes, encoding)
}

// inferEncodingFromContent 从内容推断编码方式
func inferEncodingFromContent(data []byte) StringEncoding {
	if len(data) < 2 {
		return EncodingUTF8
	}

	// 检查是否为UTF-16模式
	// UTF-16BE: 高字节在前，低字节在后，且低字节通常为0x00（ASCII字符）
	// UTF-16LE: 低字节在前，高字节在后，且高字节通常为0x00（ASCII字符）

	// 检查前几个字符的模式
	utf16BECount := 0
	utf16LECount := 0

	// 检查前几个字符（最多检查前20个字节）
	maxCheck := len(data)
	if maxCheck > 20 {
		maxCheck = 20
	}

	for i := 0; i < maxCheck-1; i += 2 {
		if i+1 >= len(data) {
			break
		}

		// UTF-16BE模式：高字节非零，低字节为零（ASCII字符）
		if data[i] != 0x00 && data[i+1] == 0x00 {
			utf16BECount++
		}
		// UTF-16LE模式：低字节非零，高字节为零（ASCII字符）
		if data[i+1] != 0x00 && data[i] == 0x00 {
			utf16LECount++
		}
	}

	// 根据模式数量判断编码
	if utf16BECount > utf16LECount {
		return EncodingUTF16BE
	}
	if utf16LECount > utf16BECount {
		return EncodingUTF16LE
	}

	// 默认返回UTF-8
	return EncodingUTF8
}
