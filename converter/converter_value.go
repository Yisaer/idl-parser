package converter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/typeref"
)

func parseBytesToString(data []byte, converter *IDLConverter) (value string, remained []byte, err error) {
	if len(data) < 4 {
		return "", nil, fmt.Errorf("expect data len larger than %v got len %v", 4, len(data))
	}
	strLen, remained, err := parseBytesToInt64(data, 4, converter)
	if err != nil {
		return "", nil, fmt.Errorf("parse sequence len error:%v", err.Error())
	}
	if strLen == 0 {
		// 字符串长度为 0 的特殊情况
		return "", remained, nil
	}
	if int64(len(remained)) < strLen {
		return "", nil, errors.New("data truncated, insufficient bytes for string")
	}
	return string(remained[:strLen]), remained[strLen:], nil
}

func parseBytesToInt64(data []byte, expLen int, converter *IDLConverter) (int64, []byte, error) {
	if len(data) < expLen {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", expLen, len(data))
	}
	parseData, remainData := data[:expLen], data[expLen:]
	got, err := bytesToInt64(parseData, converter)
	return got, remainData, err
}

func parseBytesToInt16(data []byte, converter *IDLConverter) (int64, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 2, len(data))
	}
	parseData, remainData := data[:2], data[2:]
	var value int16
	if converter.IsLittleEndian {
		value = int16(binary.LittleEndian.Uint16(parseData))
	} else {
		value = int16(binary.BigEndian.Uint16(parseData))
	}
	return int64(value), remainData, nil
}

func parseBytesToUint16(data []byte, converter *IDLConverter) (int64, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 2, len(data))
	}
	parseData, remainData := data[:2], data[2:]
	var value uint16
	if converter.IsLittleEndian {
		value = binary.LittleEndian.Uint16(parseData)
	} else {
		value = binary.BigEndian.Uint16(parseData)
	}
	return int64(value), remainData, nil
}

func parseBytesToInt32(data []byte, converter *IDLConverter) (int64, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	var value int32
	if converter.IsLittleEndian {
		value = int32(binary.LittleEndian.Uint32(parseData))
	} else {
		value = int32(binary.BigEndian.Uint32(parseData))
	}
	return int64(value), remainData, nil
}

func parseBytesToUint32(data []byte, converter *IDLConverter) (int64, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	var value uint32
	if converter.IsLittleEndian {
		value = binary.LittleEndian.Uint32(parseData)
	} else {
		value = binary.BigEndian.Uint32(parseData)
	}
	return int64(value), remainData, nil
}

func parseBytesToUint64(data []byte, converter *IDLConverter) (int64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 8, len(data))
	}
	parseData, remainData := data[:8], data[8:]
	var value uint64
	if converter.IsLittleEndian {
		value = binary.LittleEndian.Uint64(parseData)
	} else {
		value = binary.BigEndian.Uint64(parseData)
	}
	return int64(value), remainData, nil
}

func bytesToInt64(b []byte, converter *IDLConverter) (int64, error) {
	switch len(b) {
	case 1:
		return int64(b[0]), nil
	case 2:
		if converter.IsLittleEndian {
			return int64(binary.LittleEndian.Uint16(b)), nil
		} else {
			return int64(binary.BigEndian.Uint16(b)), nil
		}
	case 4:
		if converter.IsLittleEndian {
			return int64(binary.LittleEndian.Uint32(b)), nil
		} else {
			return int64(binary.BigEndian.Uint32(b)), nil
		}
	case 8:
		if converter.IsLittleEndian {
			return int64(binary.LittleEndian.Uint64(b)), nil
		} else {
			return int64(binary.BigEndian.Uint64(b)), nil
		}
	default:
		return 0, fmt.Errorf("unexpect data len:%v", len(b))
	}
}

func parseBytesToBoolean(data []byte) (bool, []byte, error) {
	if len(data) < 1 {
		return false, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	return data[0] != 0x00, data[1:], nil
}

func parseBytesToFloat64(data []byte, expLen int, converter *IDLConverter) (float64, []byte, error) {
	if len(data) < expLen {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", expLen, len(data))
	}
	if expLen == 4 {
		parseData := data[:4]
		remainData := data[4:]
		var value float32
		if converter.IsLittleEndian {
			value = math.Float32frombits(binary.LittleEndian.Uint32(parseData))
		} else {
			value = math.Float32frombits(binary.BigEndian.Uint32(parseData))
		}
		return float64(value), remainData, nil
	} else if expLen == 8 {
		parseData := data[:8]
		remainData := data[8:]
		var value float64
		if converter.IsLittleEndian {
			value = math.Float64frombits(binary.LittleEndian.Uint64(parseData))
		} else {
			value = math.Float64frombits(binary.BigEndian.Uint64(parseData))
		}
		return value, remainData, nil
	}
	return 0, nil, fmt.Errorf("expect data len 4/8 got len %v", expLen)
}

func parseBytesToList(data []byte, seqType typeref.Sequence, module ast.Module, converter *IDLConverter) ([]interface{}, []byte, error) {
	if len(data) <= 4 {
		return nil, nil, fmt.Errorf("expect data len larger than %v got len %v", 4, len(data))
	}
	sequenceLen, remained, err := parseBytesToInt64(data, 4, converter)
	if err != nil {
		return nil, nil, fmt.Errorf("parse sequence len error:%v", err.Error())
	}
	result := make([]interface{}, 0, sequenceLen)
	var v interface{}
	for i := 0; i < int(sequenceLen); i++ {
		v, remained, err = parseDataByType(remained, seqType.InnerType, module, converter)
		if err != nil {
			return nil, nil, fmt.Errorf("parse sequence %v error:%v", seqType.InnerType, err.Error())
		}
		result = append(result, v)
	}
	return result, remained, nil
}
