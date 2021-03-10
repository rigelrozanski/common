package parse

import (
	"fmt"
	"go/ast"
	"strings"
)

type ParsedInterface struct {
	Name             string
	Functions        []ParsedInterfaceFunc
	Comment          []string
	CommentStartLine int
	StartLine        int
	EndLine          int
}

type ParsedInterfaceFunc struct {
	StartLine     int
	FuncName      string        // SomeFuncName
	RecreatedCode string        // SomeFuncName(input string) string
	OutputFields  []ParsedField // all the output fields, ordered
}

func GetCurrentParsedInterface(file string, lineNo int) (
	inter ParsedInterface, found bool) {

	pfile := ParseFile(file)
	for _, in := range pfile.parsedInterfaces {
		if in.StartLine <= lineNo && lineNo <= in.EndLine {
			return in, true
		}
	}
	return ParsedInterface{}, false
}

func (pc ParseContext) ParseInterface(decl ast.Decl) (out ParsedInterface, found bool) {

	spec, comment, commentStartLine, found := pc.GetSpecAndComment(decl)
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
	if len(comment) > 0 {
		out.CommentStartLine = commentStartLine
	} else {
		out.CommentStartLine = out.StartLine
	}

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
		outFlds, resultsConcat := pc.FieldList(fnType.Results)
		if strings.Contains(resultsConcat, " ") {
			resultsConcat = fmt.Sprintf("(%v)", resultsConcat)
		}

		fnName := fn.Names[0].Name
		fnLine := pc.fset.PositionFor(fn.Names[0].NamePos, false).Line
		fnStr := fmt.Sprintf("%v(%v) %v", fnName, paramsConcat, resultsConcat)

		out.Functions = append(out.Functions,
			ParsedInterfaceFunc{
				StartLine:     fnLine,
				FuncName:      fnName,
				RecreatedCode: fnStr,
				OutputFields:  outFlds,
			})

	}
	return out, true
}

// -----------------------------------------------------

// enterChar could be \n or \r
func InterfaceCodeFromFuncs(interName string, pf []ParsedFunc, enterChar string) string {
	out := fmt.Sprintf("type %v interface {%v", interName, enterChar)
	for _, fn := range pf {
		out += fmt.Sprintf("    %v%v", fn.RecreatedCode, enterChar)
	}
	out += fmt.Sprintf("}%v", enterChar)
	return out
}
