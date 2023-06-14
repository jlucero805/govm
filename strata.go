package main

import (
	"fmt"
	"strata/interp"
	"strata/parser"
	"strata/value"
	"time"
)

func strata() {
	ast := parser.Parse(`
	fn lst, f => do
				let go = fn lst, f, acc =>
					if len(lst) > 0 then
						go(tail(lst), f, push(acc, f(head(lst))))
					else
						acc
				go(lst, f, [])
			end`)
	val, _ := interp.Interp(ast, value.New())
	fmt.Printf("%#v", val)
}

func strataFib() {
	ast := parser.ParseAll(`
	let fib = fn n =>
		if n == 1 then
			1
		else if n == 2 then
			1
		else
			fib(n - 1) + fib(n - 2)
	fib(40)`)
	interp.TopInterp(ast)
	//fmt.Println(got.Binds)
}

func fib(n int) int {
	if n == 1 {
		return 1
	} else if n == 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func timeStrata(times int) float64 {
	var experiments int64 = 0
	for i := 0; i < times; i += 1 {
		start := time.Now()
		strataFib()
		end := time.Since(start).Microseconds()
		experiments += end
	}
	return float64(experiments) / float64(times)
}

func timeGolang(times int) float64 {
	var experiments int64 = 0
	for i := 0; i < times; i += 1 {
		start := time.Now()
		fib(40)
		end := time.Since(start).Microseconds()
		experiments += end
	}
	return float64(experiments) / float64(times)
}

func main() {
	// t := timeStrata(10)
	// fmt.Println(t)
	// tgo := timeGolang(10)
	// fmt.Println(tgo)
	// fmt.Println(tgo / t)
	strata()
}
