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
			return nil, nil, fmt.Errorf("type %v is not octet type", t.TypeRefType())
		}
		return parseShort(data, short, module, converter)
	case typ.UnsignedShortType:
		return parseBytesToUint16(data, converter)
	case typ.LongType:
		return parseBytesToInt32(data, converter)
	case typ.UnsignedLongType:
		return parseBytesToUint32(data, converter)
	case typ.LongLongType:
		return parseBytesToInt64(data, 8, converter)
	case typ.UnsignedLongLongType:
		return parseBytesToUint64(data, converter)
	case typ.BooleanType:
		return parseBytesToBoolean(data)
	case typ.FloatType:
		return parseBytesToFloat64(data, 4, converter)
	case typ.SequenceType:
		seq := t.(typeref.Sequence)
		return parseBytesToList(data, seq, module, converter)
	case typ.StringType:
		return parseBytesToString(data, converter)
	case typ.SelfDefinedTypeType:
		definedType := t.(typeref.TypeName)
		for _, c := range module.Content {
			if c.GetName() == definedType.Name && c.ModuleContentType() == typ.StructType {
				st := c.(struct_type.Struct)
				return parseBytesToStruct(data, st, module, converter)
			}
		}
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
