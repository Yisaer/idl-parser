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

func (c *IDLConverter) parseBytesToInt64(data []byte, expLen int) (int64, []byte, error) {
	if len(data) < expLen {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", expLen, len(data))
	}
	parseData, remainData := data[:expLen], data[expLen:]
	got, err := c.bytesToInt64(parseData)
	return got, remainData, err
}

func (c *IDLConverter) parseBytesToInt16(data []byte) (int64, []byte, error) {
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
	return int64(value), remainData, nil
}

func (c *IDLConverter) parseBytesToUint16(data []byte) (int64, []byte, error) {
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
	return int64(value), remainData, nil
}

func (c *IDLConverter) parseBytesToInt32(data []byte) (int64, []byte, error) {
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
	return int64(value), remainData, nil
}

func (c *IDLConverter) parseBytesToUint32(data []byte) (int64, []byte, error) {
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
	return int64(value), remainData, nil
}

func (c *IDLConverter) parseBytesToUint64(data []byte) (int64, []byte, error) {
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
	return int64(value), remainData, nil
}

func (c *IDLConverter) bytesToInt64(b []byte) (int64, error) {
	switch len(b) {
	case 1:
		return int64(b[0]), nil
	case 2:
		if c.IsLittleEndian {
			return int64(binary.LittleEndian.Uint16(b)), nil
		} else {
			return int64(binary.BigEndian.Uint16(b)), nil
		}
	case 4:
		if c.IsLittleEndian {
			return int64(binary.LittleEndian.Uint32(b)), nil
		} else {
			return int64(binary.BigEndian.Uint32(b)), nil
		}
	case 8:
		if c.IsLittleEndian {
			return int64(binary.LittleEndian.Uint64(b)), nil
		} else {
			return int64(binary.BigEndian.Uint64(b)), nil
		}
	default:
		return 0, fmt.Errorf("unexpect data len:%v", len(b))
	}
}

func (c *IDLConverter) parseBytesToBoolean(data []byte) (bool, []byte, error) {
	if len(data) < 1 {
		return false, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	return data[0] != 0x00, data[1:], nil
}

func (c *IDLConverter) parseBytesToFloat64(data []byte, expLen int) (float64, []byte, error) {
	if len(data) < expLen {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", expLen, len(data))
	}
	if expLen == 4 {
		parseData := data[:4]
		remainData := data[4:]
		var value float32
		if c.IsLittleEndian {
			value = math.Float32frombits(binary.LittleEndian.Uint32(parseData))
		} else {
			value = math.Float32frombits(binary.BigEndian.Uint32(parseData))
		}
		return float64(value), remainData, nil
	} else if expLen == 8 {
		parseData := data[:8]
		remainData := data[8:]
		var value float64
		if c.IsLittleEndian {
			value = math.Float64frombits(binary.LittleEndian.Uint64(parseData))
		} else {
			value = math.Float64frombits(binary.BigEndian.Uint64(parseData))
		}
		return value, remainData, nil
	}
	return 0, nil, fmt.Errorf("expect data len 4/8 got len %v", expLen)
}

func (c *IDLConverter) parseBytesToList(data []byte, seqType typeref.Sequence, module ast.Module) ([]interface{}, []byte, error) {
	if len(data) <= 4 {
		return nil, nil, fmt.Errorf("expect data len larger than %v got len %v", 4, len(data))
	}
	sequenceLen, remained, err := c.parseBytesToInt64(data, 4)
	if err != nil {
		return nil, nil, fmt.Errorf("parse sequence len error:%v", err.Error())
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
