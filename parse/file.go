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
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(fset, f)
}

type ParsedGoFile struct {
	parsedInterfaces []ParsedInterface
	parsedStructs    []ParsedStruct
	parsedFuncs      []ParsedFunc
}

func ParseFile(file string) (out ParsedGoFile) {
	fset := token.NewFileSet()
	pc := NewParseContext(fset)
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		infc, found := pc.ParseInterface(decl)
		if found {
			out.parsedInterfaces = append(out.parsedInterfaces, infc)
			continue
		}
		strct, found := pc.ParseStruct(decl)
		if found {
			out.parsedStructs = append(out.parsedStructs, strct)
			continue
		}
		fn, found := pc.ParseFunc(decl)
		if found {
			out.parsedFuncs = append(out.parsedFuncs, fn)
			continue
		}
	}
	return out
}
