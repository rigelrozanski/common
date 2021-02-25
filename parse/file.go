package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

// AUX debugging function
func FullParseAndPrintFile(file string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(fset, f)
}

type ParsedGoFile struct {
	parsedInterfaces []ParsedInterface
	parsedStruct     []ParsedStruct
}

func ParseFile(file string) (out ParsedGoFile) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		infc, found := ParseInterface(fset, decl)
		if found {
			out.parsedInterfaces = append(out.parsedInterfaces, infc)
		}
	}
	return out
}
