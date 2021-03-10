package main

import (
	"fmt"

	"github.com/rigelrozanski/common/parse"
)

func main() {
	parse.FullParseAndPrintFile("testfile.go")
	pf := parse.ParseFile("testfile.go")
	fmt.Printf("%v\n", pf)

	allfuncs := parse.GetAllFuncsOfStruct("./", "TestStrct")
	fmt.Printf("\n\nALL FUNCS%v\n", allfuncs)

	//s, err := json.MarshalIndent(pf, "", "\t")
	//if err != nil {
	//log.Fatal(err)
	//}
	//fmt.Println(string(s))
}
