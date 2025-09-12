package converter

import (
	"testing"
	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/typeref"
)

// createTestConverterForParseTests 创建用于解析测试的 IDLConverter 实例
func createTestConverterForParseTests(isLittleEndian bool) *IDLConverter {
	return &IDLConverter{
		IDlConverterConfig: IDlConverterConfig{
			IsLittleEndian:    isLittleEndian,
			LengthFieldLength: 4,
		},
		TotalBytes:  0,
		ParsedBytes: 0,
	}
}

func TestParseOctet(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：octet(1字节) + 剩余数据(2字节)
	testData := []byte{0x42, 0x11, 0x22}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	ot := typeref.OctetType{}
	result, remained, err := converter.parseOctet(testData, ot, ast.Module{})
	
	if err != nil {
		t.Errorf("parseOctet failed: %v", err)
		return
	}
	
	// 验证解析结果
	if result != uint8(0x42) {
		t.Errorf("Expected result 0x42, got %v", result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x11, 0x22}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 1
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseShort(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：short(2字节) + 剩余数据(2字节)
	testData := []byte{0x12, 0x34, 0x55, 0x66}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	short := typeref.ShortType{}
	result, remained, err := converter.parseShort(testData, short, ast.Module{})
	
	if err != nil {
		t.Errorf("parseShort failed: %v", err)
		return
	}
	
	// 验证解析结果 (big endian: 0x1234 = 4660)
	if result != int16(4660) {
		t.Errorf("Expected result 4660, got %v", result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x55, 0x66}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 2
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseShortLittleEndian(t *testing.T) {
	converter := createTestConverterForParseTests(true)
	
	// 测试数据：short(2字节) + 剩余数据(2字节)
	testData := []byte{0x34, 0x12, 0x55, 0x66}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	short := typeref.ShortType{}
	result, remained, err := converter.parseShort(testData, short, ast.Module{})
	
	if err != nil {
		t.Errorf("parseShort failed: %v", err)
		return
	}
	
	// 验证解析结果 (little endian: 0x1234 = 4660)
	if result != int16(4660) {
		t.Errorf("Expected result 4660, got %v", result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x55, 0x66}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 2
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseUnsignedShort(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：unsigned short(2字节) + 剩余数据(2字节)
	testData := []byte{0xFF, 0xFF, 0x55, 0x66}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	unsignedShort := typeref.UnsignedShortType{}
	result, remained, err := converter.parseUnsignedShort(testData, unsignedShort, ast.Module{})
	
	if err != nil {
		t.Errorf("parseUnsignedShort failed: %v", err)
		return
	}
	
	// 验证解析结果 (0xFFFF = 65535)
	if result != uint16(65535) {
		t.Errorf("Expected result 65535, got %v", result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x55, 0x66}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 2
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseLong(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：long(4字节) + 剩余数据(2字节)
	testData := []byte{0x12, 0x34, 0x56, 0x78, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	long := typeref.LongType{}
	result, remained, err := converter.parseLong(testData, long, ast.Module{})
	
	if err != nil {
		t.Errorf("parseLong failed: %v", err)
		return
	}
	
	// 验证解析结果 (big endian: 0x12345678 = 305419896)
	if result != int32(305419896) {
		t.Errorf("Expected result 305419896, got %v", result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 4
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseUnsignedLong(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：unsigned long(4字节) + 剩余数据(2字节)
	testData := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	unsignedLong := typeref.UnsignedLongType{}
	result, remained, err := converter.parseUnsignedLong(testData, unsignedLong, ast.Module{})
	
	if err != nil {
		t.Errorf("parseUnsignedLong failed: %v", err)
		return
	}
	
	// 验证解析结果 (0xFFFFFFFF = 4294967295)
	if result != uint32(4294967295) {
		t.Errorf("Expected result 4294967295, got %v", result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 4
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseLongLong(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：long long(8字节) + 剩余数据(2字节)
	testData := []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	longLong := typeref.LongLongType{}
	result, remained, err := converter.parseLongLong(testData, longLong, ast.Module{})
	
	if err != nil {
		t.Errorf("parseLongLong failed: %v", err)
		return
	}
	
	// 验证解析结果 (big endian: 0x123456789ABCDEF0 = 1311768467294899696)
	expectedResult := int64(0x123456789ABCDEF0)
	if result != expectedResult {
		t.Errorf("Expected result %d, got %v", expectedResult, result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 8
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseUnsignedLongLong(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：unsigned long long(8字节) + 剩余数据(2字节)
	testData := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	unsignedLongLong := typeref.UnsignedLongLongType{}
	result, remained, err := converter.parseUnsignedLongLong(testData, unsignedLongLong, ast.Module{})
	
	if err != nil {
		t.Errorf("parseUnsignedLongLong failed: %v", err)
		return
	}
	
	// 验证解析结果 (0xFFFFFFFFFFFFFFFF = 18446744073709551615 in uint64)
	expectedResult := uint64(0xFFFFFFFFFFFFFFFF)
	if result != expectedResult {
		t.Errorf("Expected result %d, got %v", expectedResult, result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 8
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseBoolean(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：boolean(1字节) + 剩余数据(2字节)
	testData := []byte{0x01, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	boolean := typeref.BooleanType{}
	result, remained, err := converter.parseBoolean(testData, boolean, ast.Module{})
	
	if err != nil {
		t.Errorf("parseBoolean failed: %v", err)
		return
	}
	
	// 验证解析结果
	if result != true {
		t.Errorf("Expected result true, got %v", result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 1
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseFloat(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：float(4字节) + 剩余数据(2字节)
	// 使用 3.14 的 IEEE 754 表示 (big endian)
	testData := []byte{0x40, 0x48, 0xF5, 0xC3, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	float := typeref.FloatType{}
	result, remained, err := converter.parseFloat(testData, float, ast.Module{})
	
	if err != nil {
		t.Errorf("parseFloat failed: %v", err)
		return
	}
	
	// 验证解析结果 (3.14)
	expectedResult := float32(3.14)
	if result != expectedResult {
		t.Errorf("Expected result %v, got %v", expectedResult, result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 4
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseDouble(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：double(8字节) + 剩余数据(2字节)
	// 使用 3.141592653589793 的 IEEE 754 表示 (big endian)
	testData := []byte{0x40, 0x09, 0x21, 0xFB, 0x54, 0x44, 0x2D, 0x18, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	double := typeref.DoubleType{}
	result, remained, err := converter.parseDouble(testData, double, ast.Module{})
	
	if err != nil {
		t.Errorf("parseDouble failed: %v", err)
		return
	}
	
	// 验证解析结果 (3.141592653589793)
	expectedResult := 3.141592653589793
	if result != expectedResult {
		t.Errorf("Expected result %v, got %v", expectedResult, result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes
	expectedParsedBytes := 8
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseString(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：长度字段(4字节) + UTF-8 BOM(3字节) + 字符串内容(5字节) + 终止符(1字节) + 剩余数据(2字节)
	// 长度字段: 0x00000009 (9，包含BOM、字符串内容和终止符)
	// UTF-8 BOM: 0xEF 0xBB 0xBF
	// 字符串内容: "hello"
	// 终止符: 0x00
	testData := []byte{0x00, 0x00, 0x00, 0x09, 0xEF, 0xBB, 0xBF, 'h', 'e', 'l', 'l', 'o', 0x00, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	stringType := typeref.StringType{}
	result, remained, err := converter.parseString(testData, stringType, ast.Module{})
	
	if err != nil {
		t.Errorf("parseString failed: %v", err)
		return
	}
	
	// 验证解析结果
	expectedResult := "hello"
	if result != expectedResult {
		t.Errorf("Expected result %q, got %q", expectedResult, result)
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes (长度字段4字节 + BOM、字符串内容和终止符9字节 = 13字节)
	expectedParsedBytes := 13
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseArray(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：3个short元素(6字节) + 剩余数据(2字节)
	testData := []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	arrayType := typeref.ArrayType{
		ArraySize:  3,
		InnerType:  typeref.ShortType{},
	}
	result, remained, err := converter.parseArray(testData, arrayType, ast.Module{})
	
	if err != nil {
		t.Errorf("parseArray failed: %v", err)
		return
	}
	
	// 验证解析结果
	if len(result) != 3 {
		t.Errorf("Expected array length 3, got %d", len(result))
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes (parseArray通过parseDataByType调用，每个short元素2字节，3个元素=6字节)
	expectedParsedBytes := 6
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}

func TestParseBytesToList(t *testing.T) {
	converter := createTestConverterForParseTests(false)
	
	// 测试数据：长度字段(4字节) + 2个short元素(4字节) + 剩余数据(2字节)
	// 长度字段: 0x00000002 (2)
	testData := []byte{0x00, 0x00, 0x00, 0x02, 0x12, 0x34, 0x56, 0x78, 0x99, 0xAA}
	converter.TotalBytes = len(testData)
	converter.ParsedBytes = 0
	
	seqType := typeref.Sequence{
		InnerType: typeref.ShortType{},
	}
	result, remained, err := converter.parseBytesToList(testData, seqType, ast.Module{})
	
	if err != nil {
		t.Errorf("parseBytesToList failed: %v", err)
		return
	}
	
	// 验证解析结果
	if len(result) != 2 {
		t.Errorf("Expected sequence length 2, got %d", len(result))
	}
	
	// 验证剩余数据
	expectedRemained := []byte{0x99, 0xAA}
	if len(remained) != len(expectedRemained) {
		t.Errorf("Expected remained length %d, got %d", len(expectedRemained), len(remained))
	}
	for i, b := range expectedRemained {
		if remained[i] != b {
			t.Errorf("Expected remained[%d] = 0x%02x, got 0x%02x", i, b, remained[i])
		}
	}
	
	// 验证 ParsedBytes (长度字段4字节 + 2个short元素4字节 = 8字节)
	expectedParsedBytes := 8
	if converter.ParsedBytes != expectedParsedBytes {
		t.Errorf("Expected ParsedBytes %d, got %d", expectedParsedBytes, converter.ParsedBytes)
	}
}