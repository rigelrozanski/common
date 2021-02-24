package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func main() {
	FullParseAndPrintFile("testdir/testfile.go")
	ParseFile("testdir/testfile.go")
}

func FullParseAndPrintFile(file string) ParsedGoFile {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(fset, f)
}

type ParsedInterfaceFunc struct {
	StartLine       int
	EndLine         int
	FuncMaterialStr string // SomeFuncName(input string) string
}

type ParsedInterface struct {
	Name      string
	Functions []ParsedInterfaceFunc
	StartLine int
	EndLine   int
}

type ParsedStructField struct {
	Name      string // field name
	Type      string
	StartLine int
	EndLine   int
}

type ParsedStruct struct {
	Name      string
	Fields    []ParsedStructField
	StartLine int
	EndLine   int
}

type ParsedGoFile struct {
	parsedInterfaces []ParsedInterface
	parsedStruct     []ParsedStruct
}

func ParseFile(file string) ParsedGoFile {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	/////////////////////////
	// Interfaces
	for _, decl := range f.Decls {
		declGen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if declGen.Tok != token.TYPE {
			continue
		}
		if len(declGen.Specs) == 0 {
			continue
		}
		spec, ok := declGen.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		// is interface?
		it, ok := spec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		interfaceName := spec.Name.Name
		interfaceStartLine := fset.PositionFor(it.Methods.Opening, false).Line
		interfaceEndLine := fset.PositionFor(it.Methods.Closing, false).Line

		fmt.Println("found interface")
		fmt.Printf("name: %v\n", interfaceName)
		fmt.Printf("start Line: %v\n", interfaceStartLine)
		fmt.Printf("end Line: %v\n", interfaceEndLine)

		fns := it.Methods.List
		for _, fn := range fns {
			if len(fn.Names) == 0 {
				continue
			}
			if fn.Names[0].Obj.Kind != ast.Fun { // skip inherited interfaces (only consider functions of this interface)
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

			fmt.Printf("function on line: %v\n", fnLine)
			fmt.Printf("reproduced function specifics: %v(%v) %v\n", fnName, paramsConcat, resultsConcat)
		}

	}
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
