package converter

import (
	"testing"
)

func TestParsingProgress(t *testing.T) {
	// 创建测试用的 IDLConverter
	converter := &IDLConverter{
		IDlConverterConfig: IDlConverterConfig{
			IsLittleEndian:    true,
			LengthFieldLength: 4,
		},
	}

	// 测试数据
	testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	// 初始化进度跟踪
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0

	// 测试初始状态
	parsedBytes, totalBytes, progressPercent := converter.GetParsingProgress()
	if parsedBytes != 0 {
		t.Errorf("Expected parsedBytes to be 0, got %d", parsedBytes)
	}
	if totalBytes != len(testData) {
		t.Errorf("Expected totalBytes to be %d, got %d", len(testData), totalBytes)
	}
	if progressPercent != 0.0 {
		t.Errorf("Expected progressPercent to be 0.0, got %f", progressPercent)
	}

	// 模拟解析过程
	converter.ParsedBytes = 4 // 模拟解析了4个字节

	// 测试解析后的状态
	parsedBytes, totalBytes, progressPercent = converter.GetParsingProgress()
	if parsedBytes != 4 {
		t.Errorf("Expected parsedBytes to be 4, got %d", parsedBytes)
	}
	if totalBytes != len(testData) {
		t.Errorf("Expected totalBytes to be %d, got %d", len(testData), totalBytes)
	}
	expectedProgress := float64(4) / float64(len(testData)) * 100.0
	if progressPercent != expectedProgress {
		t.Errorf("Expected progressPercent to be %f, got %f", expectedProgress, progressPercent)
	}

	// 测试完成解析
	converter.ParsedBytes = len(testData)
	parsedBytes, totalBytes, progressPercent = converter.GetParsingProgress()
	if parsedBytes != len(testData) {
		t.Errorf("Expected parsedBytes to be %d, got %d", len(testData), parsedBytes)
	}
	if progressPercent != 100.0 {
		t.Errorf("Expected progressPercent to be 100.0, got %f", progressPercent)
	}
}

func TestParsingProgressWithZeroTotalBytes(t *testing.T) {
	// 创建测试用的 IDLConverter
	converter := &IDLConverter{
		IDlConverterConfig: IDlConverterConfig{
			IsLittleEndian:    true,
			LengthFieldLength: 4,
		},
	}

	// 测试零字节情况
	converter.TotalBytes = 0
	converter.ParsedBytes = 0

	parsedBytes, totalBytes, progressPercent := converter.GetParsingProgress()
	if parsedBytes != 0 {
		t.Errorf("Expected parsedBytes to be 0, got %d", parsedBytes)
	}
	if totalBytes != 0 {
		t.Errorf("Expected totalBytes to be 0, got %d", totalBytes)
	}
	if progressPercent != 0.0 {
		t.Errorf("Expected progressPercent to be 0.0, got %f", progressPercent)
	}
}

func TestGetParsedBytes(t *testing.T) {
	converter := &IDLConverter{}
	converter.ParsedBytes = 100

	if converter.GetParsedBytes() != 100 {
		t.Errorf("Expected GetParsedBytes() to return 100, got %d", converter.GetParsedBytes())
	}
}

func TestGetTotalBytes(t *testing.T) {
	converter := &IDLConverter{}
	converter.TotalBytes = 200

	if converter.GetTotalBytes() != 200 {
		t.Errorf("Expected GetTotalBytes() to return 200, got %d", converter.GetTotalBytes())
	}
}