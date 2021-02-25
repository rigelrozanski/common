package main

import "fmt"

func somefunc() {
	fmt.Println("vim-go")
}

// First Comment
// Second Comment
type TestInter interface {
	FirstTestFn(hello string) string
	SecondTestFn() (goodbye int)
	ThirdTestFn(mixed string, input, types bool) (a, bunch string, of, outputs bool)
	FourthTestFn(sup, is, up bool)
	FifthTestFn(sup, is, up bool) func()
	SixthTestFn(fnInput func(func(bloop string)) string) func()
}

type TestStrct struct {
	firstVar  string
	secondVar string
	thirdVar  string
}
