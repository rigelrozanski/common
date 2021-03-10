package parse

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"path"
	"strings"
)

func GetAllFuncsOfStruct(packageDir, structName string) (pf []ParsedFunc) {
	fileInfo, err := ioutil.ReadDir(packageDir)
	if err != nil {
		return pf
	}
	for _, file := range fileInfo {
		fname := file.Name()
		if file.IsDir() || !strings.HasSuffix(fname, ".go") {
			continue
		}
		fpath := path.Join(packageDir, fname)
		pfile := ParseFile(fpath)

		for _, fnc := range pfile.parsedFuncs {
			if fnc.FunctionOf == structName {
				pf = append(pf, fnc)
			}
		}
	}
	return pf
}

type ParsedFunc struct {
	Name                  string
	FunctionOf            string
	RecreatedCode         string
	RecreatedCodeForInter string
	Comment               []string
	StartLine             int
	EndLine               int
}

// NewParsedFunc creates a new ParsedFunc object
func NewParsedFunc(name, functionOf, recreatedCode, recreatedCodeForInter string,
	comment []string, startLine, endLine int) ParsedFunc {
	return ParsedFunc{
		Name:                  name,
		FunctionOf:            functionOf,
		RecreatedCode:         recreatedCode,
		RecreatedCodeForInter: recreatedCodeForInter,
		Comment:               comment,
		StartLine:             startLine,
		EndLine:               endLine,
	}
}

func (pc ParseContext) ParseFunc(decl ast.Decl) (out ParsedFunc, found bool) {

	fd, ok := decl.(*ast.FuncDecl)
	if !ok {
		return out, false
	}

	r := fd.Recv
	fnOf := ""
	if r != nil {
		if len(r.List) > 0 {
			_, fnOf, _, found = pc.ExprTypeString(r.List[0].Type)
			if !found {
				return out, false
			}
		}
	}
	name := fd.Name.Name
	comment := strings.Split(fd.Doc.Text(), `\n`)
	sl := pc.fset.PositionFor(fd.Body.Lbrace, false).Line
	el := pc.fset.PositionFor(fd.Body.Rbrace, false).Line
	rc, _ := pc.FuncTypeString(fd.Type)
	rcfi := name + string(rc[4:len(rc)])

	out = NewParsedFunc(name, fnOf, rc, rcfi, comment, sl, el)
	return out, true
}

func (pc ParseContext) FuncTypeString(ft *ast.FuncType) (typeString, dummyFulfill string) {
	typeString = "func("
	_, outP := pc.FieldList(ft.Params)
	typeString += outP
	typeString += ")"

	flds, res := pc.FieldList(ft.Results)
	if len(res) > 0 {
		typeString += " " + res
	}
	dummyFulfill = typeString

	if len(res) > 0 {
		dummyFnReturn := ""
		for _, fld := range flds {
			if len(dummyFnReturn) > 0 {
				dummyFnReturn += ", "
			}
			dummyFnReturn += fld.DummyFulfill
		}
		dummyFulfill += fmt.Sprintf("{ return %v }", dummyFnReturn)
	} else {
		dummyFulfill += "{}"
	}
	return typeString, dummyFulfill
}
