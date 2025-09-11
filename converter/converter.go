package converter

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/yisaer/idl-parser/ast"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/typeref"
)

type IDLConverter struct {
	IsLittleEndian bool
	SchemaPath     string
	Module         ast.Module
}

func NewIDLConverter(schemaPath string, isLittleEndian bool) (*IDLConverter, error) {
	c := &IDLConverter{
		SchemaPath:     schemaPath,
		IsLittleEndian: isLittleEndian,
	}
	if err := c.init(); err != nil {
		return nil, err
	}
	return c, nil
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
	v, _, err := c.parseBytesToStruct(data, *targetSts, *targetModule)
	return v, err
}

func (c *IDLConverter) parseDataByType(data []byte, t typeref.TypeRef, module ast.Module) (interface{}, []byte, error) {
	switch t.TypeRefType() {
	case typ.ArrayType:
		arrayType, ok := t.(typeref.ArrayType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not array type", t.TypeRefType())
		}
		return c.parseArray(data, arrayType, module)
	case typ.OctetType:
		ot, ok := t.(typeref.OctetType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not octet type", t.TypeRefType())
		}
		return c.parseOctet(data, ot, module)
	case typ.ShortType:
		short, ok := t.(typeref.ShortType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not short type", t.TypeRefType())
		}
		return c.parseShort(data, short, module)
	case typ.UnsignedShortType:
		unsignedShort, ok := t.(typeref.UnsignedShortType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned short type", t.TypeRefType())
		}
		return c.parseUnsignedShort(data, unsignedShort, module)
	case typ.LongType:
		long, ok := t.(typeref.LongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not long type", t.TypeRefType())
		}
		return c.parseLong(data, long, module)
	case typ.UnsignedLongType:
		unsignedLong, ok := t.(typeref.UnsignedLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned long type", t.TypeRefType())
		}
		return c.parseUnsignedLong(data, unsignedLong, module)
	case typ.LongLongType:
		longLong, ok := t.(typeref.LongLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not long long type", t.TypeRefType())
		}
		return c.parseLongLong(data, longLong, module)
	case typ.UnsignedLongLongType:
		unsignedLongLong, ok := t.(typeref.UnsignedLongLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned long long type", t.TypeRefType())
		}
		return c.parseUnsignedLongLong(data, unsignedLongLong, module)
	case typ.BooleanType:
		boolean, ok := t.(typeref.BooleanType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not boolean type", t.TypeRefType())
		}
		return c.parseBoolean(data, boolean, module)
	case typ.FloatType:
		float, ok := t.(typeref.FloatType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not float type", t.TypeRefType())
		}
		return c.parseFloat(data, float, module)
	case typ.SequenceType:
		seq := t.(typeref.Sequence)
		return c.parseBytesToList(data, seq, module)
	case typ.StringType:
		st, ok := t.(typeref.StringType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not string type", t.TypeRefType())
		}
		return c.parseString(data, st, module)
	case typ.SelfDefinedTypeType:
		definedType := t.(typeref.TypeName)
		return c.parseSelfDefineType(data, definedType, module)
	}
	return nil, nil, fmt.Errorf("unsupported type:%v", t.TypeName())
}

func (c *IDLConverter) parseArray(data []byte, arrayType typeref.ArrayType, module ast.Module) (interface{}, []byte, error) {
	remained := data
	var err error
	var v interface{}
	result := make([]interface{}, 0, arrayType.ArraySize)

	for i := 0; i < arrayType.ArraySize; i++ {
		v, remained, err = c.parseDataByType(remained, arrayType.InnerType, module)
		if err != nil {
			return nil, nil, fmt.Errorf("parse array element %v error:%v", i, err.Error())
		}
		result = append(result, v)
	}
	return result, remained, nil
}

func (c *IDLConverter) parseOctet(data []byte, ot typeref.OctetType, _ ast.Module) (interface{}, []byte, error) {
	if ot.ArraySize > 0 {
		remained := data
		var err error
		var v int64
		result := make([]int64, 0)
		for i := 0; i < ot.ArraySize; i++ {
			v, remained, err = c.parseBytesToInt64(remained, 1)
			if err != nil {
				return nil, nil, err
			}
			result = append(result, v)
		}
		return result, remained, nil
	}
	return c.parseBytesToInt64(data, 1)
}

func (c *IDLConverter) parseShort(data []byte, short typeref.ShortType, _ ast.Module) (interface{}, []byte, error) {
	return c.parseBytesToInt16(data)
}

func (c *IDLConverter) parseUnsignedShort(data []byte, unsignedShort typeref.UnsignedShortType, _ ast.Module) (interface{}, []byte, error) {
	return c.parseBytesToUint16(data)
}

func (c *IDLConverter) parseLong(data []byte, long typeref.LongType, _ ast.Module) (interface{}, []byte, error) {
	return c.parseBytesToInt32(data)
}

func (c *IDLConverter) parseUnsignedLong(data []byte, unsignedLong typeref.UnsignedLongType, _ ast.Module) (interface{}, []byte, error) {
	return c.parseBytesToUint32(data)
}

func (c *IDLConverter) parseLongLong(data []byte, longLong typeref.LongLongType, _ ast.Module) (interface{}, []byte, error) {
	return c.parseBytesToInt64(data, 8)
}

func (c *IDLConverter) parseUnsignedLongLong(data []byte, unsignedLongLong typeref.UnsignedLongLongType, _ ast.Module) (interface{}, []byte, error) {
	return c.parseBytesToUint64(data)
}

func (c *IDLConverter) parseBoolean(data []byte, boolean typeref.BooleanType, _ ast.Module) (interface{}, []byte, error) {
	return c.parseBytesToBoolean(data)
}

func (c *IDLConverter) parseFloat(data []byte, float typeref.FloatType, _ ast.Module) (interface{}, []byte, error) {
	return c.parseBytesToFloat64(data, 4)
}

func (c *IDLConverter) parseString(data []byte, st typeref.StringType, _ ast.Module) (interface{}, []byte, error) {
	if st.Length != nil {
		return ParseBytesToFixedLengthString(data, *st.Length)
	}
	return ParseBytesToDynamicStringWithEndian(data, c.IsLittleEndian)
}

func (c *IDLConverter) parseSelfDefineType(data []byte, definedType typeref.TypeName, module ast.Module) (interface{}, []byte, error) {
	for _, content := range module.Content {
		if content.GetName() == definedType.Name && content.ModuleContentType() == typ.StructType {
			st := content.(struct_type.Struct)
			return c.parseBytesToStruct(data, st, module)
		}
	}
	return nil, nil, fmt.Errorf("unsupported type:%v", definedType.Name)
}
