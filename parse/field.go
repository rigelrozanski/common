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
		pos, fldTypeStr, found := pc.ExprTypeString(field.Type)

		if !found {
			// undefined type, move along
			continue OUTER
		}

		if len(fieldsConcat) > 0 {
			fieldsConcat += " "
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

func (pc ParseContext) ExprTypeString(expr ast.Expr) (pos token.Pos, fldTypeStr string, found bool) {
	switch et := expr.(type) {
	case *ast.Ident:
		fldTypeStr = et.Name
		pos = et.NamePos
	case *ast.ArrayType:
		_, arrType, found := pc.ExprTypeString(et.Elt)
		if !found {
			return 0, "", false
		}
		fldTypeStr = fmt.Sprintf("[]%v", arrType)
		pos = et.Lbrack
	case *ast.MapType:
		pos = et.Map
		_, keyStr, found := pc.ExprTypeString(et.Key)

		if !found {
			return 0, "", false
		}
		_, valStr, found := pc.ExprTypeString(et.Value)
		if !found {
			return 0, "", false
		}
		fldTypeStr = fmt.Sprintf("map[%v]%v", keyStr, valStr)
	case *ast.FuncType:
		fldTypeStr = pc.FuncTypeString(et)
		pos = et.Func
	default:
		return 0, "", false
	}
	return pos, fldTypeStr, true
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
