package main

import (
	"fmt"
	"strata/interp"
	"strata/parser"
)

func main() {
	ast := parser.ParseAll(`
	let sum = fn n =>
		if 0 < n then
			n + sum(n - 1)
		else
			0
	let result = sum(5)
	`)

	got := interp.TopInterp(ast)
	//interp.TopInterp(ast)
	fmt.Println(got.Binds)
	//printEnv(got)
}
