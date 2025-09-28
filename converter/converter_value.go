package converter

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

func isNeedPaddingType(t typeref.TypeRef) bool {
	switch t.TypeRefType() {
	case typ.StringType:
		st := t.(typeref.StringType)
		return st.Length == nil
	case typ.SequenceType:
		return true
	}
	return false
}

func (c *IDLConverter) ParseBytesToStruct(data []byte, st struct_type.Struct, module ast.Module) (map[string]interface{}, []byte, error) {
	var v interface{}
	var err error
	var remained []byte
	remained = data
	m := make(map[string]any, len(st.Fields))
	for _, field := range st.Fields {
		v, remained, err = c.ParseDataByType(remained, field.Type, module)
		if err != nil {
			return nil, nil, fmt.Errorf("struct %v parse field %v error:%v", st, field.Name, err.Error())
		}
		m[field.Name] = v
		if len(remained) > 0 && isNeedPaddingType(field.Type) {
			remained, err = c.consumePadding(remained)
			if err != nil {
				return nil, nil, fmt.Errorf("struct %v parse field %v error:%v", st, field.Name, err.Error())
			}
		}
	}
	return m, remained, nil
}

func (c *IDLConverter) ParseBytesToInt8(data []byte) (int8, []byte, error) {
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	parseData, remainData := data[:1], data[1:]
	var value int8
	value = int8(parseData[0])
	return value, remainData, nil
}

func (c *IDLConverter) ParseBytesToInt16(data []byte) (int16, []byte, error) {
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

func (c *IDLConverter) ParseBytesToInt32(data []byte) (int32, []byte, error) {
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

func (c *IDLConverter) ParseBytesToInt64(data []byte) (int64, []byte, error) {
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

func (c *IDLConverter) ParseBytesToUint8(data []byte) (uint8, []byte, error) {
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	parseData, remainData := data[:1], data[1:]
	var value uint8
	value = parseData[0]
	return value, remainData, nil
}

func (c *IDLConverter) ParseBytesToUint16(data []byte) (uint16, []byte, error) {
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

func (c *IDLConverter) ParseBytesToUint32(data []byte) (uint32, []byte, error) {
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

func (c *IDLConverter) ParseBytesToUint64(data []byte) (uint64, []byte, error) {
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

func (c *IDLConverter) ParseBytesToBoolean(data []byte) (bool, []byte, error) {
	if len(data) < 1 {
		return false, nil, fmt.Errorf("expect data len %v got len %v", 1, len(data))
	}
	return data[0] != 0x00, data[1:], nil
}

func (c *IDLConverter) ParseBytesToFloat32(data []byte) (float32, []byte, error) {
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

func (c *IDLConverter) ParseBytesToFloat64(data []byte) (float64, []byte, error) {
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

func (c *IDLConverter) ParseBytesToList(data []byte, seqType typeref.Sequence, module ast.Module) ([]interface{}, []byte, error) {
	sequenceBytesLen, remained, err := c.ParseBytesToLength(data)
	if err != nil {
		return nil, nil, err
	}

	// 检查是否有足够的数据来解析sequence
	if len(remained) < sequenceBytesLen {
		return nil, nil, fmt.Errorf("insufficient data for sequence: expected %d bytes, got %d bytes", sequenceBytesLen, len(remained))
	}

	// sequenceBytesLen 是后续元素的总字节数，不是元素个数
	// remained 始终表示当前总共还剩多少字节未解析
	result := make([]interface{}, 0)
	sequenceData := remained[:sequenceBytesLen] // 提取sequence的数据部分
	sequenceRemained := sequenceData            // 用于解析sequence内部元素

	for len(sequenceRemained) > 0 {
		var v interface{}
		v, sequenceRemained, err = c.ParseDataByType(sequenceRemained, seqType.InnerType, module)
		if err != nil {
			return nil, nil, fmt.Errorf("parse sequence %v error:%v", seqType.InnerType, err.Error())
		}
		result = append(result, v)

		// 如果还有剩余数据且当前元素类型需要padding，则消费padding
		if len(sequenceRemained) > 0 && isNeedPaddingType(seqType.InnerType) {
			sequenceRemained, err = c.consumePadding(sequenceRemained)
			if err != nil {
				return nil, nil, fmt.Errorf("parse sequence %v error:%v", seqType.InnerType, err.Error())
			}
		}
	}

	// 检查 sequenceRemained 是否恰好为 0，如果不是 0，说明解析有问题
	if len(sequenceRemained) != 0 {
		return nil, nil, fmt.Errorf("sequence parsing incomplete: %d bytes remaining after parsing all elements", len(sequenceRemained))
	}

	// remained 现在指向sequence数据之后的位置，即整个数据流中剩余未解析的字节
	remained = remained[sequenceBytesLen:]

	return result, remained, nil
}

func (c *IDLConverter) ParseBytesToLength(data []byte) (int, []byte, error) {
	switch c.LengthFieldLength {
	case 1:
		v, remained, err := c.ParseBytesToUint8(data)
		if err != nil {
			return 0, nil, err
		}
		return int(v), remained, nil
	case 2:
		v, remained, err := c.ParseBytesToUint16(data)
		if err != nil {
			return 0, nil, err
		}
		return int(v), remained, nil
	case 4:
		v, remained, err := c.ParseBytesToUint32(data)
		if err != nil {
			return 0, nil, err
		}
		return int(v), remained, nil
	default:
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 1, c.LengthFieldLength)
	}
}
