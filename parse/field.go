package parse

import (
	"fmt"
	"go/ast"
	"go/token"
)

// NOTE cannot be used for interfaces of functions
// empty dummy being "" for a string for instance
// or SomeStruct{} for SomeStruct
func BasicFieldTypeToEmptyDummy(Type string) (emptyDummy string) {
	if len(Type) > 0 && string(Type[0]) == "*" {
		return "nil"
	}
	switch Type {
	case "error":
		return "nil"
	case "string":
		return `""`
	case "int", "int8", "int16", "int32", "int64",
		"uint8", "uint16", "uint32", "uint64",
		"float64", "float32", "complex64", "complex128":
		return "0"
	case "byte":
		return "0x00"
	case "rune":
		return `' '`
	case "bool":
		return "false"
	default: // good for arrays, maps, structs
		return fmt.Sprintf("%v{}", Type)
	}

}

type ParsedField struct {
	Names        []string // field name(s)
	Type         string
	DummyFulfill string // for instance "" if Type was string
	StartLine    int
}

// NewParsedField creates a new ParsedField object
func NewParsedField(names []string, Type, dummyFulfill string, startLine int) ParsedField {
	return ParsedField{
		Names:        names,
		Type:         Type,
		DummyFulfill: dummyFulfill,
		StartLine:    startLine,
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
		pos, fldTypeStr, dummyFulfill, found := pc.ExprTypeString(field.Type)

		if !found {
			// undefined type, move along
			continue OUTER
		}

		if len(fieldsConcat) > 0 {
			fieldsConcat += " "
		}
		fieldsConcat += fmt.Sprintf("%v, ", fldTypeStr)

		posInt := pc.fset.PositionFor(pos, false).Line

		fld := NewParsedField(fldNames, fldTypeStr, dummyFulfill, posInt)
		flds = append(flds, fld)
	}
	if len(fieldsConcat) > 2 { // remove final comma+space from all the fields
		fieldsConcat = string(fieldsConcat[:len(fieldsConcat)-2])
	}

	return flds, fieldsConcat
}

func (pc ParseContext) ExprTypeString(expr ast.Expr) (
	pos token.Pos, fldTypeStr, dummyFulfill string, found bool) {

	switch et := expr.(type) {
	case *ast.Ident:
		fldTypeStr = et.Name
		dummyFulfill = BasicFieldTypeToEmptyDummy(fldTypeStr)
		pos = et.NamePos
	case *ast.ArrayType:
		_, arrType, _, found := pc.ExprTypeString(et.Elt)
		if !found {
			return 0, "", "", false
		}
		fldTypeStr = fmt.Sprintf("[]%v", arrType)
		dummyFulfill = BasicFieldTypeToEmptyDummy(fldTypeStr)
		pos = et.Lbrack
	case *ast.MapType:
		pos = et.Map
		_, keyStr, _, found := pc.ExprTypeString(et.Key)

		if !found {
			return 0, "", "", false
		}
		_, valStr, _, found := pc.ExprTypeString(et.Value)
		if !found {
			return 0, "", "", false
		}
		fldTypeStr = fmt.Sprintf("map[%v]%v", keyStr, valStr)
		dummyFulfill = BasicFieldTypeToEmptyDummy(fldTypeStr)
	case *ast.FuncType:
		fldTypeStr, dummyFulfill = pc.FuncTypeString(et)
		pos = et.Func
	case *ast.InterfaceType:
		pos = et.Interface
		fldTypeStr = "interface{}"
		dummyFulfill = "nil"

	default:
		return 0, "", "", false
	}
	return pos, fldTypeStr, dummyFulfill, true
}
