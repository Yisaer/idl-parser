package converter

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

type IDlConverterConfig struct {
	IsLittleEndian    bool
	LengthFieldLength int // should be 1,2,4 bytes
	SchemaPath        string
	PaddingLength     int // should be 1,2,4,8,16,32,64
}

type IDLConverter struct {
	IDlConverterConfig
	Module     ast.Module
	TotalBytes int
}

func NewIDLConverterWithModule(Config IDlConverterConfig, module ast.Module) (*IDLConverter, error) {
	c := &IDLConverter{
		IDlConverterConfig: Config,
	}
	if err := c.verifyConfig(); err != nil {
		return nil, err
	}
	c.Module = module
	if err := c.verifyModule(c.Module); err != nil {
		return nil, err
	}
	if err := c.verifyStruct(c.Module); err != nil {
		return nil, err
	}
	return c, nil
}

func NewIDLConverter(Config IDlConverterConfig) (*IDLConverter, error) {
	c := &IDLConverter{
		IDlConverterConfig: Config,
	}
	if err := c.verifyConfig(); err != nil {
		return nil, err
	}
	if err := c.init(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *IDLConverter) verifyConfig() error {
	switch c.LengthFieldLength {
	case 1, 2, 4:
		break
	default:
		return errors.New("invalid length field, should be 1, 2 or 4")
	}
	switch c.PaddingLength {
	case 1, 2, 4, 8, 32, 64:
		break
	default:
		return errors.New("invalid padding length field, should be 1, 2, 4, 16, 32, 64")
	}
	return nil
}

func (c *IDLConverter) init() error {
	v, err := os.ReadFile(c.SchemaPath)
	if err != nil {
		return err
	}
	res := ast.Parse(string(v))
	if res.Err != nil {
		return fmt.Errorf("remaining: %v, err:%v", res.Remaining, res.Err)
	}
	c.Module = res.Output
	if err := c.verifyModule(c.Module); err != nil {
		return err
	}
	if err := c.verifyStruct(c.Module); err != nil {
		return err
	}
	return nil
}

func (c *IDLConverter) travelModule(targetSchemaID string) (*struct_type.Struct, *ast.Module, error) {
	linkList := list.New()
	nodes := strings.Split(targetSchemaID, ".")
	for _, node := range nodes {
		linkList.PushBack(node)
	}
	mockRootModule := ast.Module{
		Name: "header",
		Content: []ast.ModuleContent{
			c.Module,
		},
		Type: "module",
	}
	return c.travel(linkList.Front(), mockRootModule)
}

func (c *IDLConverter) travel(curr *list.Element, currModule ast.Module) (*struct_type.Struct, *ast.Module, error) {
	node := curr.Value.(string)
	if curr.Next() == nil {
		for _, con := range currModule.Content {
			if con.GetName() == node {
				st, ok := con.(struct_type.Struct)
				if !ok {
					return nil, nil, fmt.Errorf("travel node %v not struct", node)
				}
				return &st, &currModule, nil
			}
		}
		return nil, nil, fmt.Errorf("travel node %v not found", node)
	}
	for _, con := range currModule.Content {
		if con.GetName() == node {
			module, ok := con.(ast.Module)
			if !ok {
				return nil, nil, fmt.Errorf("travel node %v not module", node)
			}
			return c.travel(curr.Next(), module)
		}
	}
	return nil, nil, fmt.Errorf("travel node %v not found", node)
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
	if len(st.Fields) < 1 {
		return fmt.Errorf("struct %v has no fields", st)
	}
	for _, field := range st.Fields {
		if !isSupportedTyp(field.Type.TypeRefType()) {
			return fmt.Errorf("st %v has unsupported field %v", st.Name, field.Name)
		}
		strTyp, ok := field.Type.(typeref.StringType)
		if ok && strTyp.Length == nil {
			return fmt.Errorf("st:%v not support field %v, dynamic string", st.Name, field.Name)
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
		// typ.SequenceType, no support dynamic slice padding
		typ.StringType,
		typ.SelfDefinedTypeType,
		typ.ArrayType,
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

func (c *IDLConverter) IsSchemaIDExist(schemaID string) bool {
	_, _, err := c.travelModule(schemaID)
	if err != nil {
		return false
	}
	return true
}

func (c *IDLConverter) Decode(schemaID string, data []byte) (map[string]interface{}, error) {
	targetSts, targetModule, err := c.travelModule(schemaID)
	if err != nil {
		return nil, err
	}
	c.TotalBytes = len(data)
	v, _, err := c.parseBytesToStruct(data, *targetSts, *targetModule, true)
	return v, err
}

func (c *IDLConverter) parseDataByType(data []byte, t typeref.TypeRef, module ast.Module, isLastField bool) (interface{}, []byte, error) {

	var result interface{}
	var remained []byte
	var err error

	switch t.TypeRefType() {
	case typ.ArrayType:
		arrayType, ok := t.(typeref.ArrayType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not array type", t.TypeRefType())
		}
		result, remained, err = c.parseArray(data, arrayType, module, isLastField)
	case typ.OctetType:
		ot, ok := t.(typeref.OctetType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not octet type", t.TypeRefType())
		}
		otResult, newRemained, newErr := c.parseOctet(data, ot, module)
		result, remained, err = otResult, newRemained, newErr
	case typ.ShortType:
		short, ok := t.(typeref.ShortType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not short type", t.TypeRefType())
		}
		shortResult, newRemained, newErr := c.parseShort(data, short, module)
		result, remained, err = shortResult, newRemained, newErr
	case typ.UnsignedShortType:
		unsignedShort, ok := t.(typeref.UnsignedShortType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned short type", t.TypeRefType())
		}
		ushortResult, newRemained, newErr := c.parseUnsignedShort(data, unsignedShort, module)
		result, remained, err = ushortResult, newRemained, newErr
	case typ.LongType:
		long, ok := t.(typeref.LongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not long type", t.TypeRefType())
		}
		longResult, newRemained, newErr := c.parseLong(data, long, module)
		result, remained, err = longResult, newRemained, newErr
	case typ.UnsignedLongType:
		unsignedLong, ok := t.(typeref.UnsignedLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned long type", t.TypeRefType())
		}
		ulongResult, newRemained, newErr := c.parseUnsignedLong(data, unsignedLong, module)
		result, remained, err = ulongResult, newRemained, newErr
	case typ.LongLongType:
		longLong, ok := t.(typeref.LongLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not long long type", t.TypeRefType())
		}
		llongResult, newRemained, newErr := c.parseLongLong(data, longLong, module)
		result, remained, err = llongResult, newRemained, newErr
	case typ.UnsignedLongLongType:
		unsignedLongLong, ok := t.(typeref.UnsignedLongLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned long long type", t.TypeRefType())
		}
		ullongResult, newRemained, newErr := c.parseUnsignedLongLong(data, unsignedLongLong, module)
		result, remained, err = ullongResult, newRemained, newErr
	case typ.BooleanType:
		boolean, ok := t.(typeref.BooleanType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not boolean type", t.TypeRefType())
		}
		boolResult, newRemained, newErr := c.parseBoolean(data, boolean, module)
		result, remained, err = boolResult, newRemained, newErr
	case typ.FloatType:
		float, ok := t.(typeref.FloatType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not float type", t.TypeRefType())
		}
		floatResult, newRemained, newErr := c.parseFloat(data, float, module)
		result, remained, err = floatResult, newRemained, newErr
	case typ.DoubleType:
		double, ok := t.(typeref.DoubleType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not double type", t.TypeRefType())
		}
		doubleResult, newRemained, newErr := c.parseDouble(data, double, module)
		result, remained, err = doubleResult, newRemained, newErr
	case typ.SequenceType:
		seq := t.(typeref.Sequence)
		result, remained, err = c.parseBytesToList(data, seq, module, isLastField)
	case typ.StringType:
		st, ok := t.(typeref.StringType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not string type", t.TypeRefType())
		}
		result, remained, err = c.parseString(data, st, module, isLastField)
	case typ.SelfDefinedTypeType:
		definedType := t.(typeref.TypeName)
		result, remained, err = c.parseSelfDefineType(data, definedType, module, isLastField)
	default:
		return nil, nil, fmt.Errorf("unsupported type:%v", t.TypeName())
	}
	if err != nil {
		return nil, nil, err
	}
	return result, remained, nil
}

func (c *IDLConverter) parseArray(data []byte, arrayType typeref.ArrayType, module ast.Module, isParLast bool) ([]interface{}, []byte, error) {
	remained := data
	var err error
	var v interface{}
	result := make([]interface{}, 0, arrayType.ArraySize)
	for i := 0; i < arrayType.ArraySize; i++ {
		v, remained, err = c.parseDataByType(remained, arrayType.InnerType, module, isParLast && i == arrayType.ArraySize-1)
		if err != nil {
			return nil, nil, fmt.Errorf("parse array element %v error:%v", i, err.Error())
		}
		result = append(result, v)
	}
	return result, remained, nil
}

func (c *IDLConverter) parseOctet(data []byte, ot typeref.OctetType, _ ast.Module) (uint8, []byte, error) {
	result, remained, err := c.parseBytesToUint8(data)
	return uint8(result), remained, err
}

func (c *IDLConverter) parseShort(data []byte, short typeref.ShortType, _ ast.Module) (int16, []byte, error) {
	result, remained, err := c.parseBytesToInt16(data)
	return int16(result), remained, err
}

func (c *IDLConverter) parseUnsignedShort(data []byte, unsignedShort typeref.UnsignedShortType, _ ast.Module) (uint16, []byte, error) {
	result, remained, err := c.parseBytesToUint16(data)
	return uint16(result), remained, err
}

func (c *IDLConverter) parseLong(data []byte, long typeref.LongType, _ ast.Module) (int32, []byte, error) {
	result, remained, err := c.parseBytesToInt32(data)
	return int32(result), remained, err
}

func (c *IDLConverter) parseUnsignedLong(data []byte, unsignedLong typeref.UnsignedLongType, _ ast.Module) (uint32, []byte, error) {
	result, remained, err := c.parseBytesToUint32(data)
	return uint32(result), remained, err
}

func (c *IDLConverter) parseLongLong(data []byte, longLong typeref.LongLongType, _ ast.Module) (int64, []byte, error) {
	result, remained, err := c.parseBytesToInt64(data)
	return result, remained, err
}

func (c *IDLConverter) parseUnsignedLongLong(data []byte, unsignedLongLong typeref.UnsignedLongLongType, _ ast.Module) (uint64, []byte, error) {
	result, remained, err := c.parseBytesToUint64(data)
	return uint64(result), remained, err
}

func (c *IDLConverter) parseBoolean(data []byte, boolean typeref.BooleanType, _ ast.Module) (bool, []byte, error) {
	result, remained, err := c.parseBytesToBoolean(data)
	return result, remained, err
}

func (c *IDLConverter) parseFloat(data []byte, float typeref.FloatType, _ ast.Module) (float32, []byte, error) {
	result, remained, err := c.parseBytesToFloat32(data)
	return result, remained, err
}

func (c *IDLConverter) parseDouble(data []byte, float typeref.DoubleType, _ ast.Module) (float64, []byte, error) {
	result, remained, err := c.parseBytesToFloat64(data)
	return result, remained, err
}

func (c *IDLConverter) parseString(data []byte, st typeref.StringType, _ ast.Module, isLast bool) (interface{}, []byte, error) {
	var result interface{}
	var remained []byte
	var err error

	if st.Length != nil {
		result, remained, err = c.ParseBytesToFixedLengthString(data, *st.Length)
	} else {
		result, remained, err = c.ParseBytesToDynamicString(data)
		if err != nil {
			return nil, nil, err
		}
		remained, err = c.consumePadding(remained, isLast)
		if err != nil {
			return nil, nil, fmt.Errorf("consuming dynamic string field meet error:%v", err.Error())
		}
	}
	return result, remained, nil
}

func (c *IDLConverter) parseSelfDefineType(data []byte, definedType typeref.TypeName, module ast.Module, isParLast bool) (interface{}, []byte, error) {
	for _, content := range module.Content {
		if content.GetName() == definedType.Name && content.ModuleContentType() == typ.StructType {
			st := content.(struct_type.Struct)
			return c.parseBytesToStruct(data, st, module, isParLast)
		}
	}
	return nil, nil, fmt.Errorf("unsupported type:%v", definedType.Name)
}

func (c *IDLConverter) consumePadding(remained []byte, isLast bool) ([]byte, error) {
	if len(remained) > 0 && !isLast {
		consumedBytes := c.TotalBytes - len(remained)
		extraBytesLength := consumedBytes % c.PaddingLength
		paddingBytesLength := c.PaddingLength - extraBytesLength
		if paddingBytesLength > len(remained) {
			return remained, fmt.Errorf("padding length exceeds remained length")
		}
		if extraBytesLength > 0 {
			return remained[paddingBytesLength:], nil
		}
	}
	return remained, nil
}
