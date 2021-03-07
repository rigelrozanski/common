package parse

import (
	"go/ast"
)

type ParsedStruct struct {
	Name      string
	Fields    []ParsedField
	Comment   []string
	StartLine int
	EndLine   int
}

// GetParsedStruct
func GetCurrentParsedStruct(file string, lineNo int) (strct ParsedStruct, found bool) {
	pfile := ParseFile(file)
	for _, st := range pfile.parsedStructs {
		if st.StartLine <= lineNo && lineNo <= st.EndLine {
			return st, true
		}
	}
	return strct, false
}

// NewParsedStruct creates a new ParsedStruct object
func NewParsedStruct(name string, fields []ParsedField, comment []string, startLine int, endLine int) ParsedStruct {
	return ParsedStruct{
		Name:      name,
		Fields:    fields,
		Comment:   comment,
		StartLine: startLine,
		EndLine:   endLine,
	}
}

func (pc ParseContext) ParseStruct(decl ast.Decl) (out ParsedStruct, found bool) {

	spec, comment, found := GetSpecAndComment(decl)
	if !found {
		return out, false
	}
	name := spec.Name.Name

	strct, ok := spec.Type.(*ast.StructType)
	if !ok {
		return out, false
	}
	flds, _ := pc.FieldList(strct.Fields)

	startLine := pc.fset.PositionFor(strct.Fields.Opening, false).Line
	endLine := pc.fset.PositionFor(strct.Fields.Closing, false).Line

	out = NewParsedStruct(name, flds, comment, startLine, endLine)
	return out, true
}
