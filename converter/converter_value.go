package converter

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typeref"
)

func (c *IDLConverter) parseBytesToStruct(data []byte, st struct_type.Struct, module ast.Module) (map[string]interface{}, []byte, error) {
	var v interface{}
	var err error
	var remained []byte
	remained = data
	m := make(map[string]any, len(st.Fields))
	for _, field := range st.Fields {
		v, remained, err = c.parseDataByType(remained, field.Type, module)
		if err != nil {
			return nil, nil, fmt.Errorf("struct %v parse field %v error:%v", st, field.Name, err.Error())
		}
		m[field.Name] = v
	}
	return m, remained, nil
}

func (c *IDLConverter) parseBytesToInt64(data []byte) (int64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 8, len(data))
	}
	parseData, remainData := data[:8], data[8:]
	var value int64
	if c.IsLittleEndian {
		value = int64(binary.LittleEndian.Uint64(parseData))
	} else {
		value = int64(binary.BigEndian.Uint64(parseData))
	}
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToInt8(data []byte) (int8, []byte, error) {
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	parseData, remainData := data[:1], data[1:]
	var value int8
	value = int8(parseData[0])
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToInt16(data []byte) (int16, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 2, len(data))
	}
	parseData, remainData := data[:2], data[2:]
	var value int16
	if c.IsLittleEndian {
		value = int16(binary.LittleEndian.Uint16(parseData))
	} else {
		value = int16(binary.BigEndian.Uint16(parseData))
	}
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToUint8(data []byte) (uint8, []byte, error) {
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	parseData, remainData := data[:1], data[1:]
	var value uint8
	value = parseData[0]
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToUint16(data []byte) (uint16, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 2, len(data))
	}
	parseData, remainData := data[:2], data[2:]
	var value uint16
	if c.IsLittleEndian {
		value = binary.LittleEndian.Uint16(parseData)
	} else {
		value = binary.BigEndian.Uint16(parseData)
	}
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToInt32(data []byte) (int32, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	var value int32
	if c.IsLittleEndian {
		value = int32(binary.LittleEndian.Uint32(parseData))
	} else {
		value = int32(binary.BigEndian.Uint32(parseData))
	}
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToUint32(data []byte) (uint32, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	var value uint32
	if c.IsLittleEndian {
		value = binary.LittleEndian.Uint32(parseData)
	} else {
		value = binary.BigEndian.Uint32(parseData)
	}
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToUint64(data []byte) (uint64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 8, len(data))
	}
	parseData, remainData := data[:8], data[8:]
	var value uint64
	if c.IsLittleEndian {
		value = binary.LittleEndian.Uint64(parseData)
	} else {
		value = binary.BigEndian.Uint64(parseData)
	}
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToBoolean(data []byte) (bool, []byte, error) {
	if len(data) < 1 {
		return false, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	return data[0] != 0x00, data[1:], nil
}

func (c *IDLConverter) parseBytesToFloat32(data []byte) (float32, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	var value float32
	if c.IsLittleEndian {
		value = math.Float32frombits(binary.LittleEndian.Uint32(parseData))
	} else {
		value = math.Float32frombits(binary.BigEndian.Uint32(parseData))
	}
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToFloat64(data []byte) (float64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 8, len(data))
	}
	parseData, remainData := data[:8], data[8:]
	var value float64
	if c.IsLittleEndian {
		value = math.Float64frombits(binary.LittleEndian.Uint64(parseData))
	} else {
		value = math.Float64frombits(binary.BigEndian.Uint64(parseData))
	}
	return value, remainData, nil
}

func (c *IDLConverter) parseBytesToList(data []byte, seqType typeref.Sequence, module ast.Module) ([]interface{}, []byte, error) {
	sequenceLen, remained, err := c.parseBytesToLength(data)
	if err != nil {
		return nil, nil, err
	}
	result := make([]interface{}, 0, sequenceLen)
	var v interface{}
	for i := 0; i < int(sequenceLen); i++ {
		v, remained, err = c.parseDataByType(remained, seqType.InnerType, module)
		if err != nil {
			return nil, nil, fmt.Errorf("parse sequence %v error:%v", seqType.InnerType, err.Error())
		}
		result = append(result, v)
	}
	return result, remained, nil
}

func (c *IDLConverter) parseBytesToLength(data []byte) (int, []byte, error) {
	switch c.LengthFieldLength {
	case 1:
		v, remained, err := c.parseBytesToUint8(data)
		if err != nil {
			return 0, nil, err
		}
		return int(v), remained, nil
	case 2:
		v, remained, err := c.parseBytesToUint16(data)
		if err != nil {
			return 0, nil, err
		}
		return int(v), remained, nil
	case 4:
		v, remained, err := c.parseBytesToUint32(data)
		if err != nil {
			return 0, nil, err
		}
		return int(v), remained, nil
	default:
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 1, c.LengthFieldLength)
	}
}
