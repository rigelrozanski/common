package common

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// ParseStruct - parse types retrieves the line numbers of a golang file which contains struct declarations
func ParseStruct(lines []string) (funcs, structs, interfaces [][]string) {

	// Parse the file string
	fset := token.NewFileSet()                     // positions are relative to fset
	f, err := parser.ParseFile(fset, "", lines, 0) // parser.ImportsOnly)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(f.Decls); i++ {

		switch f.Decls[i].(type) {
		case *ast.GenDecl:
			typeDecl := f.Decls[i].(*ast.GenDecl)
			for k := 0; k < len(typeDecl.Specs); k++ {
				switch typeDecl.Specs[k].(type) {
				case *ast.TypeSpec:
					writeTypeOutput := false
					typeSpec := typeDecl.Specs[k].(*ast.TypeSpec)
					typeType := typeSpec.Type

					var start, end token.Pos

					switch typeType.(type) {
					case *ast.StructType:
						structDecl := typeType.(*ast.StructType)
						fields := structDecl.Fields.List

						//define start position
						start = structDecl.Struct - 1

						// Define the end position of final field
						for _, field := range fields {
							end = field.Type.End() - 1
						}
						writeTypeOutput = true

					case *ast.InterfaceType:
						structDecl := typeType.(*ast.InterfaceType)
						fields := structDecl.Methods.List

						//define start position
						start = structDecl.Interface - 1

						// Define the end position of final field
						for _, field := range fields {
							end = field.Type.End() - 1

						}
						writeTypeOutput = true

					}
					if writeTypeOutput {
						toAdd := typeSpec.Name.Name + " " +
							lines[start:end] +
							"}"
						toAdd = strings.Replace(toAdd, "\r\n", ", ", -1)
						toAdd = strings.Replace(toAdd, "\n", ", ", -1)
						toAdd = strings.Replace(toAdd, "\t", "", -1)
						toAdd = strings.Replace(toAdd, "{,", "{ ", -1)
						for j := 0; j < 10; j++ {
							toAdd = strings.Replace(toAdd, "  ", " ", -1)
						}

						switch typeType.(type) {
						case *ast.StructType:
							structs = append(structs, Str2Multiline(toAdd))
						case *ast.InterfaceType:
							types = append(types, Str2Multiline(toAdd))
						}

					}

				}
			}
		case *ast.FuncDecl:
			typeDecl := f.Decls[i].(*ast.FuncDecl)
			finalElement = "\tfunc " + typeDecl.Name.Name + "\n"
			funcs := append(funcs, Str2Multi(finalElement)
		}
	}
}
