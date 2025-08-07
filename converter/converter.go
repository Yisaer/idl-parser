package converter

import (
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

type IDLConverter struct {
	SchemaID   string
	SchemaPath string
	Module     ast.Module
	list       *list.List
	tarStruct  struct_type.Struct
}

func (c *IDLConverter) Init() error {
	v, err := os.ReadFile(c.SchemaPath)
	if err != nil {
		return err
	}
	res := ast.Parse(string(v))
	if res.Err != nil {
		return res.Err
	}
	c.Module = res.Output
	c.list = list.New()
	if err := c.travelModule(); err != nil {
		return err
	}
	if err := c.verifyModule(c.Module); err != nil {
		return err
	}
	if err := c.verifyStruct(c.Module); err != nil {
		return err
	}
	return nil
}

func (c *IDLConverter) verifyModuleSeq() error {
	for _, con := range c.Module.Content {
		st, ok := con.(struct_type.Struct)
		if !ok {
			continue
		}
		for index, field := range st.Fields {
			if field.Type.TypeRefType() == typ.SequenceType {
				if index == 0 {
					return fmt.Errorf("len should defined before sequence filed:%v in struct:%v", field.Name, st.Name)
				}
				if st.Fields[index-1].Name != "len" {
					return fmt.Errorf("len should defined before sequence filed:%v in struct:%v", field.Name, st.Name)
				}
			}
		}
	}
	return nil
}

func (c *IDLConverter) travelModule() error {
	nodes := strings.Split(c.SchemaPath, ".")
	for _, node := range nodes {
		c.list.PushFront(node)
	}
	return c.travel(c.list.Front(), c.Module)
}

func (c *IDLConverter) travel(curr *list.Element, currModule ast.Module) error {
	node := curr.Value.(string)
	if curr.Next() == nil {
		for _, con := range currModule.Content {
			if con.GetName() == node {
				st, ok := con.(struct_type.Struct)
				if !ok {
					return fmt.Errorf("travel node %v not struct", node)
				}
				c.tarStruct = st
				return nil
			}
		}
		return fmt.Errorf("travel node %v not found", node)
	}
	for _, con := range currModule.Content {
		if con.GetName() == node {
			module, ok := con.(ast.Module)
			if !ok {
				return fmt.Errorf("travel node %v not module", node)
			}
			return c.travel(curr.Next(), module)
		}
	}
	return fmt.Errorf("travel node %v not found", node)
}

func (c *IDLConverter) verifyModule(module ast.Module) error {
	moduleCount := 0
	for _, con := range module.Content {
		subModule, ok := con.(ast.Module)
		if ok {
			moduleCount++
			err := c.verifyModule(subModule)
			if err != nil {
				return err
			}
			continue
		}
	}
	if moduleCount > 0 && moduleCount < len(module.Content) {
		return fmt.Errorf("module %v has both sub module and others", module.Name)
	}
	return nil
}

func (c *IDLConverter) verifyStruct(module ast.Module) error {
	for _, con := range module.Content {
		subModule, ok := con.(ast.Module)
		if ok {
			err := c.verifyStruct(subModule)
			if err != nil {
				return err
			}
			continue
		}
		subSt, ok := con.(struct_type.Struct)
		if ok {
			if err := c.verifyStructField(subSt); err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func (c *IDLConverter) verifyStructField(st struct_type.Struct) error {
	for _, field := range st.Fields {
		if !isSupportedTyp(field.Type.TypeRefType()) {
			return fmt.Errorf("st %v has unsupported field %v", st.Name, field.Name)
		}
	}
	return nil
}

var (
	supportedFieldType = []typ.FieldRefType{
		typ.OctetType,
		typ.ShortType,
		typ.UnsignedShortType,
		typ.LongType,
		typ.UnsignedLongType,
		typ.LongLongType,
		typ.UnsignedLongLongType,
		typ.BooleanType,
		typ.FloatType,
		typ.SequenceType,
		typ.StringType,
	}
)

func isSupportedTyp(tar typ.FieldRefType) bool {
	for _, v := range supportedFieldType {
		if v == tar {
			return true
		}
	}
	return false
}

func (c *IDLConverter) Decode(data []byte) (map[string]interface{}, error) {
	m := make(map[string]any, len(c.tarStruct.Fields))
	var v interface{}
	var err error
	var remained []byte
	remained = data
	for _, field := range c.tarStruct.Fields {
		v, remained, err = parseDataByType(remained, field.Type)
		if err != nil {
			return nil, fmt.Errorf("struct %v parse field %v error:%v", c.tarStruct.Name, field.Name, err.Error())
		}
		m[field.Name] = v
	}
	return m, nil
}

func parseDataByType(data []byte, t typeref.TypeRef) (interface{}, []byte, error) {
	switch t.TypeRefType() {
	case typ.OctetType:
		return parseBytesToInt64(data, 1)
	case typ.ShortType:
		return parseBytesToInt16(data)
	case typ.UnsignedShortType:
		return parseBytesToUint16(data)
	case typ.LongType:
		return parseBytesToInt32(data)
	case typ.UnsignedLongType:
		return parseBytesToUint32(data)
	case typ.LongLongType:
		return parseBytesToInt64(data, 8)
	case typ.UnsignedLongLongType:
		return parseBytesToUint64(data)
	case typ.BooleanType:
		return parseBytesToBoolean(data)
	case typ.FloatType:
		return parseBytesToFloat64(data, 4)
	case typ.SequenceType:
		seq := t.(typeref.Sequence)
		return parseBytesToList(data, seq)
	case typ.StringType:
		return parseBytesToString(data)

	}
	return nil, nil, fmt.Errorf("unsupported type:%v", t.TypeName())
}

func parseBytesToString(data []byte) (value string, remained []byte, err error) {
	if len(data) <= 4 {
		return "", nil, fmt.Errorf("expect data len larger than %v got len %v", 4, len(data))
	}
	strLen, remained, err := parseBytesToInt64(data, 4)
	if err != nil {
		return "", nil, fmt.Errorf("parse sequence len error:%v", err.Error())
	}
	if int64(len(data)) < 4+strLen {
		return "", nil, errors.New("data truncated, insufficient bytes for string")
	}
	return string(remained[:strLen]), remained[strLen:], nil
}

func parseBytesToInt64(data []byte, expLen int) (int64, []byte, error) {
	if len(data) < expLen {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", expLen, len(data))
	}
	parseData, remainData := data[:expLen], data[expLen:]
	got, err := bytesToInt64(parseData)
	return got, remainData, err
}

func parseBytesToInt16(data []byte) (int64, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 2, len(data))
	}
	parseData, remainData := data[:2], data[2:]
	value := int16(binary.BigEndian.Uint16(parseData))
	return int64(value), remainData, nil
}

func parseBytesToUint16(data []byte) (int64, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 2, len(data))
	}
	parseData, remainData := data[:2], data[2:]
	value := binary.BigEndian.Uint16(parseData)
	return int64(value), remainData, nil
}

func parseBytesToInt32(data []byte) (int64, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	value := int32(binary.BigEndian.Uint32(parseData))
	return int64(value), remainData, nil
}

func parseBytesToUint32(data []byte) (int64, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 4, len(data))
	}
	parseData, remainData := data[:4], data[4:]
	value := binary.BigEndian.Uint32(parseData)
	return int64(value), remainData, nil
}

func parseBytesToUint64(data []byte) (int64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", 8, len(data))
	}
	parseData, remainData := data[:8], data[8:]
	value := binary.BigEndian.Uint64(parseData)
	return int64(value), remainData, nil
}

func bytesToInt64(b []byte) (int64, error) {
	switch len(b) {
	case 1:
		return int64(b[0]), nil
	case 2:
		return int64(binary.BigEndian.Uint16(b)), nil
	case 4:
		return int64(binary.BigEndian.Uint32(b)), nil
	case 8:
		return int64(binary.BigEndian.Uint64(b)), nil
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

func parseBytesToFloat64(data []byte, expLen int) (float64, []byte, error) {
	if len(data) < expLen {
		return 0, nil, fmt.Errorf("expect data len %v got len %v", expLen, len(data))
	}
	if expLen == 4 {
		parseData := data[:4]
		remainData := data[4:]
		value := math.Float32frombits(binary.BigEndian.Uint32(parseData))
		return float64(value), remainData, nil
	} else if expLen == 8 {
		parseData := data[:8]
		remainData := data[8:]
		value := math.Float64frombits(binary.BigEndian.Uint64(parseData))
		return value, remainData, nil
	}
	return 0, nil, fmt.Errorf("expect data len 4/8 got len %v", len(data))
}

func parseBytesToList(data []byte, seqType typeref.Sequence) ([]interface{}, []byte, error) {
	if len(data) <= 4 {
		return nil, nil, fmt.Errorf("expect data len larger than %v got len %v", 4, len(data))
	}
	sequenceLen, remained, err := parseBytesToInt64(data, 4)
	if err != nil {
		return nil, nil, fmt.Errorf("parse sequence len error:%v", err.Error())
	}
	result := make([]interface{}, 0, sequenceLen)
	var v interface{}
	for i := 0; i < int(sequenceLen); i++ {
		v, remained, err = parseDataByType(remained, seqType.InnerType)
		if err != nil {
			return nil, nil, fmt.Errorf("parse sequence %v error:%v", seqType.InnerType, err.Error())
		}
		result = append(result, v)
	}
	return result, remained, nil
}
