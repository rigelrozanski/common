package parse

import "go/token"

type ParseContext struct {
	fset *token.FileSet
}

// NewParseContext creates a new ParseContext object
func NewParseContext(fset *token.FileSet) ParseContext {
	return ParseContext{
		fset: fset,
	}
}
