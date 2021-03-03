package parse

import (
	"fmt"
	"go/ast"
	"go/token"
)

type ParsedField struct {
	Names     []string // field name(s)
	Type      string
	StartLine int
}

// NewParsedField creates a new ParsedField object
func NewParsedField(names []string, Type string, startLine int) ParsedField {
	return ParsedField{
		Names:     names,
		Type:      Type,
		StartLine: startLine,
	}
}

func (pc ParseContext) FieldList(fl *ast.FieldList) (flds []ParsedField, recreatedCode string) {

	if fl == nil {
		return flds, ""
	}

	var fieldsConcat string

OUTER:
	for _, field := range fl.List {
		if field == nil {
			continue
		}

		fldNames := []string{}
		for _, fieldName := range field.Names {
			fieldNameStr := fieldName.Name
			fieldsConcat += fmt.Sprintf("%v, ", fieldNameStr)
			fldNames = append(fldNames, fieldNameStr)
		}
		if len(fieldsConcat) > 2 { // remove final comma+space from the names
			fieldsConcat = string(fieldsConcat[:len(fieldsConcat)-2])
		}

		var fldTypeStr string
		var pos token.Pos
		switch ft := field.Type.(type) {
		case *ast.Ident:
			if len(fieldsConcat) > 0 {
				fieldsConcat += " "
			}
			fldTypeStr = ft.Name
			pos = ft.NamePos
		case *ast.FuncType:
			if len(fieldsConcat) > 0 {
				fieldsConcat += " "
			}
			fldTypeStr = pc.FuncTypeString(ft)
			pos = ft.Func
		default:
			// undefined type, move along
			continue OUTER
		}
		fieldsConcat += fmt.Sprintf("%v, ", fldTypeStr)

		posInt := pc.fset.PositionFor(pos, false).Line
		fld := NewParsedField(fldNames, fldTypeStr, posInt)
		flds = append(flds, fld)
	}
	if len(fieldsConcat) > 2 { // remove final comma+space from all the fields
		fieldsConcat = string(fieldsConcat[:len(fieldsConcat)-2])
	}

	return flds, fieldsConcat
}

func (pc ParseContext) FuncTypeString(ft *ast.FuncType) string {
	out := "func("
	_, outP := pc.FieldList(ft.Params)
	out += outP
	out += ")"
	_, res := pc.FieldList(ft.Results)
	if len(res) > 0 {
		out += " " + res
	}
	return out
}
