package main

import (
	"fmt"
	"strata/interp"
	"strata/parser"
	"strata/value"
)

func printEnv(env value.Env) {
	for _, b := range env.Binds {
		fmt.Printf("%v :: %v \n", b.Id.Value, b.Value)
	}
}

func main() {
	ast := parser.ParseAll(`
	let add = fn a => fn b => a + b
	let result = add(69)(2)
	`)

	//got := interp.TopInterp(ast)
	interp.TopInterp(ast)

	//printEnv(got)
}
