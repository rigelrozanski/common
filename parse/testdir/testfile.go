package main

import "fmt"

func somefunc() {
	fmt.Println("vim-go")
}

// First Comment
// Second Comment
type TestInter interface { // HEY
	FirstTestFn(hello string) string
	SecondTestFn() (goodbye int)
	ThirdTestFn(mixed string, input, types bool) (a, bunch string, of, outputs bool)
	FourthTestFn(sup, is, up bool)
	FifthTestFn(sup, is, up bool) func()
	SixthTestFn(fnInput func(func(bloop string)) string) func()
}

/*
SOME BIG COMMENT
IN STARS
*/
type TestStrct struct { // HELLO
	firstVar  string
	secondVar []string
	thirdVar  map[string]int
	fourthVar map[string][]int
	fifthVar  interface{}
	sixthVar  map[string]interface{}
}
