package parse

import (
	"fmt"
	"go/ast"
	"strings"
)

type ParsedInterface struct {
	Name      string
	Functions []ParsedInterfaceFunc
	Comment   []string
	StartLine int
	EndLine   int
}

type ParsedInterfaceFunc struct {
	StartLine     int
	RecreatedCode string // SomeFuncName(input string) string
}

func (pc ParseContext) ParseInterface(decl ast.Decl) (out ParsedInterface, found bool) {

	spec, comment, found := GetSpecAndComment(decl)
	if !found {
		return out, false
	}

	it, ok := spec.Type.(*ast.InterfaceType) // is interface?
	if !ok {
		return out, false
	}

	out.Name = spec.Name.Name
	out.Comment = comment
	out.StartLine = pc.fset.PositionFor(it.Methods.Opening, false).Line
	out.EndLine = pc.fset.PositionFor(it.Methods.Closing, false).Line

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
		_, paramsConcat := pc.FieldList(fnType.Params)

		// Get all the function results
		_, resultsConcat := pc.FieldList(fnType.Results)
		if strings.Contains(resultsConcat, " ") {
			resultsConcat = fmt.Sprintf("(%v)", resultsConcat)
		}

		fnName := fn.Names[0].Name
		fnLine := pc.fset.PositionFor(fn.Names[0].NamePos, false).Line
		fnStr := fmt.Sprintf("%v(%v) %v", fnName, paramsConcat, resultsConcat)
		out.Functions = append(out.Functions,
			ParsedInterfaceFunc{
				StartLine:     fnLine,
				RecreatedCode: fnStr,
			})

	}
	return out, true
}