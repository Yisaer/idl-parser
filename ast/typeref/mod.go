package typeref

import (
	"github.com/oleiade/gomme"

	"github.com/yisaer/idl-parser/ast/typ"
)

type TypeRef interface {
	TypeRefType() typ.FieldRefType
	TypeName() string
}

func ParseTypeRef(code string) gomme.Result[TypeRef, string] {
	if len(code) > 0 && containsArraySuffix(code) {
		arrayResult := ParseArray(code)
		if arrayResult.Err == nil {
			return gomme.Result[TypeRef, string]{
				Output:    arrayResult.Output,
				Remaining: arrayResult.Remaining,
				Err:       nil,
			}
		}
	}
	return ParseBaseType(code)
}

// containsArraySuffix
func containsArraySuffix(code string) bool {
	if len(code) < 3 {
		return false
	}

	lastBracket := -1
	for i := len(code) - 1; i >= 0; i-- {
		if code[i] == ']' {
			lastBracket = i
			break
		}
	}

	if lastBracket == -1 {
		return false
	}

	for i := lastBracket - 1; i >= 0; i-- {
		if code[i] == '[' {
			if i+1 < lastBracket {
				for j := i + 1; j < lastBracket; j++ {
					if code[j] < '0' || code[j] > '9' {
						return false
					}
				}
				return true
			}
			return false
		}
	}

	return false
}
