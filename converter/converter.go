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
	SchemaID       string
	SchemaPath     string
	Module         ast.Module
	list           *list.List
	tarStruct      struct_type.Struct
	belongModule   ast.Module
}

func NewIDLConverter(schemaID string, schemaPath string, isLittleEndian bool) (*IDLConverter, error) {
	c := &IDLConverter{
		SchemaID:       schemaID,
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

func (c *IDLConverter) travelModule() error {
	nodes := strings.Split(c.SchemaID, ".")
	for _, node := range nodes {
		c.list.PushBack(node)
	}
	mockRootModule := ast.Module{
		Name: "header",
		Content: []ast.ModuleContent{
			c.Module,
		},
		Type: "module",
	}
	return c.travel(c.list.Front(), mockRootModule)
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
				c.belongModule = currModule
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
	v, _, err := parseBytesToStruct(data, c.tarStruct, c.belongModule, c)
	return v, err
}

func parseDataByType(data []byte, t typeref.TypeRef, module ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	switch t.TypeRefType() {
	case typ.OctetType:
		ot, ok := t.(typeref.OctetType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not octet type", t.TypeRefType())
		}
		return parseOctet(data, ot, module, converter)
	case typ.ShortType:
		short, ok := t.(typeref.ShortType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not short type", t.TypeRefType())
		}
		return parseShort(data, short, module, converter)
	case typ.UnsignedShortType:
		unsignedShort, ok := t.(typeref.UnsignedShortType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned short type", t.TypeRefType())
		}
		return parseUnsignedShort(data, unsignedShort, module, converter)
	case typ.LongType:
		long, ok := t.(typeref.LongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not long type", t.TypeRefType())
		}
		return parseLong(data, long, module, converter)
	case typ.UnsignedLongType:
		unsignedLong, ok := t.(typeref.UnsignedLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned long type", t.TypeRefType())
		}
		return parseUnsignedLong(data, unsignedLong, module, converter)
	case typ.LongLongType:
		longLong, ok := t.(typeref.LongLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not long long type", t.TypeRefType())
		}
		return parseLongLong(data, longLong, module, converter)
	case typ.UnsignedLongLongType:
		unsignedLongLong, ok := t.(typeref.UnsignedLongLongType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not unsigned long long type", t.TypeRefType())
		}
		return parseUnsignedLongLong(data, unsignedLongLong, module, converter)
	case typ.BooleanType:
		boolean, ok := t.(typeref.BooleanType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not boolean type", t.TypeRefType())
		}
		return parseBoolean(data, boolean, module, converter)
	case typ.FloatType:
		float, ok := t.(typeref.FloatType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not float type", t.TypeRefType())
		}
		return parseFloat(data, float, module, converter)
	case typ.SequenceType:
		seq := t.(typeref.Sequence)
		return parseBytesToList(data, seq, module, converter)
	case typ.StringType:
		st, ok := t.(typeref.StringType)
		if !ok {
			return nil, nil, fmt.Errorf("type %v is not string type", t.TypeRefType())
		}
		return parseString(data, st, module, converter)
	case typ.SelfDefinedTypeType:
		definedType := t.(typeref.TypeName)
		return parseSelfDefineType(data, definedType, module, converter)
	}
	return nil, nil, fmt.Errorf("unsupported type:%v", t.TypeName())
}

func parseOctet(data []byte, ot typeref.OctetType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	if ot.ArraySize > 0 {
		remained := data
		var err error
		var v int64
		result := make([]int64, 0)
		for i := 0; i < ot.ArraySize; i++ {
			v, remained, err = parseBytesToInt64(remained, 1, converter)
			if err != nil {
				return nil, nil, err
			}
			result = append(result, v)
		}
		return result, remained, nil
	}
	return parseBytesToInt64(data, 1, converter)
}

func parseShort(data []byte, short typeref.ShortType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToInt16(data, converter)
}

func parseUnsignedShort(data []byte, unsignedShort typeref.UnsignedShortType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToUint16(data, converter)
}

func parseLong(data []byte, long typeref.LongType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToInt32(data, converter)
}

func parseUnsignedLong(data []byte, unsignedLong typeref.UnsignedLongType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToUint32(data, converter)
}

func parseLongLong(data []byte, longLong typeref.LongLongType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToInt64(data, 8, converter)
}

func parseUnsignedLongLong(data []byte, unsignedLongLong typeref.UnsignedLongLongType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToUint64(data, converter)
}

func parseBoolean(data []byte, boolean typeref.BooleanType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToBoolean(data)
}

func parseFloat(data []byte, float typeref.FloatType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToFloat64(data, 4, converter)
}

func parseString(data []byte, string typeref.StringType, _ ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	return parseBytesToString(data, converter)
}

func parseSelfDefineType(data []byte, definedType typeref.TypeName, module ast.Module, converter *IDLConverter) (interface{}, []byte, error) {
	for _, c := range module.Content {
		if c.GetName() == definedType.Name && c.ModuleContentType() == typ.StructType {
			st := c.(struct_type.Struct)
			return parseBytesToStruct(data, st, module, converter)
		}
	}
	return nil, nil, fmt.Errorf("unsupported type:%v", definedType.Name)
}
