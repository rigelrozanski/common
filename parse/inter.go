package parse

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type ParsedInterface struct {
	Name      string
	Functions []ParsedInterfaceFunc
	StartLine int
	EndLine   int
}

type ParsedInterfaceFunc struct {
	StartLine       int
	FuncMaterialStr string // SomeFuncName(input string) string
}

func ParseInterface(fset *token.FileSet, decl ast.Decl) (out ParsedInterface, found bool) {

	declGen, ok := decl.(*ast.GenDecl)
	if !ok {
		return out, false
	}
	if declGen.Tok != token.TYPE {
		return out, false
	}
	if len(declGen.Specs) == 0 {
		return out, false
	}
	spec, ok := declGen.Specs[0].(*ast.TypeSpec)
	if !ok {
		return out, false
	}

	// is interface?
	it, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return out, false
	}

	out.Name = spec.Name.Name
	out.StartLine = fset.PositionFor(it.Methods.Opening, false).Line
	out.EndLine = fset.PositionFor(it.Methods.Closing, false).Line

	fns := it.Methods.List
	for _, fn := range fns {
		if len(fn.Names) == 0 {
			continue
		}

		// skip inherited interfaces (only consider functions of this interface)
		if fn.Names[0].Obj.Kind != ast.Fun {
			continue
		}

		fnType, ok := fn.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		// Get all the params
		paramsConcat := FieldListString(fnType.Params)

		// Get all the function results
		resultsConcat := FieldListString(fnType.Results)
		if strings.Contains(resultsConcat, " ") {
			resultsConcat = fmt.Sprintf("(%v)", resultsConcat)
		}

		fnName := fn.Names[0].Name
		fnLine := fset.PositionFor(fn.Names[0].NamePos, false).Line
		fnStr := fmt.Sprintf("%v(%v) %v", fnName, paramsConcat, resultsConcat)
		out.Functions = append(out.Functions,
			ParsedInterfaceFunc{
				StartLine:       fnLine,
				FuncMaterialStr: fnStr,
			})

	}
	return out, true
}

func FieldListString(fl *ast.FieldList) string {

	if fl == nil {
		return ""
	}

	fieldsConcat := ""
	for _, field := range fl.List {
		if field == nil {
			continue
		}

		for _, fieldName := range field.Names {
			fieldNameStr := fieldName.Name
			fieldsConcat += fmt.Sprintf("%v, ", fieldNameStr)
		}
		if len(fieldsConcat) > 2 { // remove final comma+space from the names
			fieldsConcat = string(fieldsConcat[:len(fieldsConcat)-2])
		}

		switch ft := field.Type.(type) {
		case *ast.Ident:
			if len(fieldsConcat) > 0 {
				fieldsConcat += " "
			}
			fieldsConcat += fmt.Sprintf("%v, ", ft.Name)
		case *ast.FuncType:
			if len(fieldsConcat) > 0 {
				fieldsConcat += " "
			}
			fieldsConcat += fmt.Sprintf("%v, ", FuncTypeString(ft))
		}
	}
	if len(fieldsConcat) > 2 { // remove final comma+space from all the fields
		fieldsConcat = string(fieldsConcat[:len(fieldsConcat)-2])
	}

	return fieldsConcat
}

func FuncTypeString(ft *ast.FuncType) string {
	out := "func("
	out += FieldListString(ft.Params)
	out += ")"
	res := FieldListString(ft.Results)
	if len(res) > 0 {
		out += " " + res
	}
	return out
}
