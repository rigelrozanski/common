package parse

type ParsedStruct struct {
	Name      string
	Fields    []ParsedStructField
	StartLine int
	EndLine   int
}

type ParsedStructField struct {
	Name      string // field name
	Type      string
	StartLine int
}
