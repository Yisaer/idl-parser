package ast

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/bitset"
	"github.com/yisaer/idl-parser/ast/struct_type"
	"github.com/yisaer/idl-parser/ast/typ"
	"github.com/yisaer/idl-parser/ast/utils"
)

type ModuleContent interface {
	GetName() string
	ModuleContentType() typ.ModuleContentType
}

type Module struct {
	Name    string          `json:"name"`
	Content []ModuleContent `json:"content"`
	Type    string          `json:"type"`
}

func (m Module) GetName() string {
	return m.Name
}

func (Module) ModuleContentType() typ.ModuleContentType {
	return typ.ModuleType
}

func Parse(code string) gomme.Result[Module, string] {
	moduleTokenResult := utils.InLeftEmpty(gomme.Token[string]("module"))(code)
	if moduleTokenResult.Err != nil {
		return gomme.Failure[string, Module](moduleTokenResult.Err, code)
	}
	nameResult :=
		utils.InEmpty(
			gomme.Recognize(gomme.Pair(gomme.Alpha1[string](), gomme.Alphanumeric0[string]())),
		)(moduleTokenResult.Remaining)
	if nameResult.Err != nil {
		return gomme.Failure[string, Module](nameResult.Err, code)
	}
	contentResult := gomme.Delimited(
		utils.InEmpty(gomme.Token[string]("{")),
		gomme.Many0(utils.InEmpty(
			gomme.Terminated(gomme.Alternative(
				gomme.Map(bitset.Parse, func(output bitset.BitSet) (ModuleContent, error) { return output, nil }),
				gomme.Map(struct_type.Parse, func(output struct_type.Struct) (ModuleContent, error) { return output, nil }),
				gomme.Map(Parse, func(output Module) (ModuleContent, error) { return output, nil }),
			),
				gomme.Optional(utils.InEmpty(gomme.Token[string](";"))),
			),
		)),
		utils.InEmpty(gomme.Token[string]("}")),
	)(nameResult.Remaining)
	if contentResult.Err != nil {
		return gomme.Failure[string, Module](contentResult.Err, code)
	}
	return gomme.Success(Module{
		Name:    nameResult.Output,
		Content: contentResult.Output,
		Type:    typ.ModuleContentTypeToString(typ.ModuleType),
	}, contentResult.Remaining)
}
