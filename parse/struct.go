package parse

import (
	"go/ast"
)

type ParsedStruct struct {
	Name             string
	Fields           []ParsedField
	Comment          []string
	CommentStartLine int
	StartLine        int
	EndLine          int
}

// get the struct (if it exists) containing the lineNo within the given file.
// If the line number is over a field within the given struct, also retrieve
// that particular field
func GetCurrentParsedStruct(file string, lineNo int) (
	strct ParsedStruct, field ParsedField, found bool) {

	pfile := ParseFile(file)
	for _, st := range pfile.parsedStructs {
		if st.StartLine <= lineNo && lineNo <= st.EndLine {
			for _, fld := range st.Fields {
				if fld.StartLine == lineNo {
					return st, fld, true
				}
			}
			return st, ParsedField{}, true
		}
	}
	return ParsedStruct{}, ParsedField{}, false
}

// NewParsedStruct creates a new ParsedStruct object
func NewParsedStruct(name string, fields []ParsedField,
	comment []string, commentStartLine, startLine, endLine int) ParsedStruct {

	return ParsedStruct{
		Name:             name,
		Fields:           fields,
		Comment:          comment,
		CommentStartLine: startLine,
		StartLine:        startLine,
		EndLine:          endLine,
	}
}

func (pc ParseContext) ParseStruct(decl ast.Decl) (out ParsedStruct, found bool) {

	spec, comment, commentStartLine, found := pc.GetSpecAndComment(decl)
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

	if len(comment) == 0 {
		commentStartLine = startLine
	}
	out = NewParsedStruct(name, flds, comment, commentStartLine, startLine, endLine)
	return out, true
}
