package parse

import (
	"go/ast"
	"go/token"
	"strings"
)

func GetSpecAndComment(decl ast.Decl) (spec *ast.TypeSpec, comment []string, found bool) {

	declGen, ok := decl.(*ast.GenDecl)
	if !ok {
		return spec, comment, false
	}
	comment = strings.Split(declGen.Doc.Text(), `\n`)
	if declGen.Tok != token.TYPE {
		return spec, comment, false
	}
	if len(declGen.Specs) == 0 {
		return spec, comment, false
	}
	spec, ok = declGen.Specs[0].(*ast.TypeSpec)
	if !ok {
		return spec, comment, false
	}
	return spec, comment, true
}
