package parse

import (
	"go/ast"
	"go/token"
	"strings"
)

func (pc ParseContext) GetSpecAndComment(decl ast.Decl) (spec *ast.TypeSpec, comment []string, commentStartLine int, found bool) {

	declGen, ok := decl.(*ast.GenDecl)
	if !ok {
		return spec, comment, commentStartLine, false
	}
	commentStartLine = pc.fset.PositionFor(declGen.Doc.Pos(), false).Line
	comment = strings.Split(declGen.Doc.Text(), `\n`)
	if declGen.Tok != token.TYPE {
		return spec, comment, commentStartLine, false
	}
	if len(declGen.Specs) == 0 {
		return spec, comment, commentStartLine, false
	}
	spec, ok = declGen.Specs[0].(*ast.TypeSpec)
	if !ok {
		return spec, comment, commentStartLine, false
	}
	return spec, comment, commentStartLine, true
}
