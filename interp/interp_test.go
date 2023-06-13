package interp

import (
	"reflect"
	"strata/expr"
	"strata/parser"
	"strata/value"
	"testing"
)

func TestNumC(t *testing.T) {

	ast := expr.NumC{Value: 1}
	got, _ := Interp(ast, value.New())
	want := value.NumV{Value: 1}

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestBinop(t *testing.T) {

	ast := expr.Binop{
		Op:    "+",
		Left:  expr.NumC{Value: 1},
		Right: expr.NumC{Value: 2},
	}
	got, _ := Interp(ast, value.New())
	want := value.NumV{Value: 3}

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestLamC(t *testing.T) {

	ast := expr.LamC{
		Params: []expr.Expr{
			expr.IdC{Value: "a"},
		},
		Body: expr.Binop{
			Op:    "+",
			Left:  expr.IdC{Value: "a"},
			Right: expr.NumC{Value: 1},
		},
	}
	got, _ := Interp(ast, value.New())
	want := value.LamV{
		Params: []expr.Expr{
			expr.IdC{Value: "a"},
		},
		Body: expr.Binop{
			Op:    "+",
			Left:  expr.IdC{Value: "a"},
			Right: expr.NumC{Value: 1},
		},
		Closure: value.New(),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestTest(t *testing.T) {
	ast := parser.Parse("1 + 1")
	got, _ := Interp(ast, value.New())
	want := value.NumV{Value: 2}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestTests(t *testing.T) {
	ast := parser.Parse("1 < 1")
	got, _ := Interp(ast, value.New())
	want := value.BoolV{Value: false}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestTestss(t *testing.T) {
	ast := parser.Parse("if 1 < 10 then 1 else 100")
	got, _ := Interp(ast, value.New())
	want := value.NumV{Value: 1}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestE(t *testing.T) {
	ast := parser.ParseAll(`
	let sum = fn n =>
		if 0 < n then n + sum(n - 1) else 0
	let result = sum(5)`)
	env := TopInterp(ast)
	got := env.Get(expr.IdC{Value: "result"})
	want := value.NumV{Value: 15}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output value.Value
	}{
		{
			"recursing",
			`let summer = fn n, acc =>
				if 0 < n then
					summer(n - 1, acc + n)
				else
					acc
			let sum = fn n => summer(n, 0)
			let output = sum(5)
			`,
			value.NumV{Value: 15},
		},
		{
			"church numerals",
			`let iter = fn x => x + 1
			let init = 0
			let one = fn f => fn a => f(a)
			let add = fn c1 => fn c2 => fn f => fn a =>
				c1(f)(c2(f)(a))
			let two = add(one)(one)
			let output = two(iter)(init)`,
			value.NumV{Value: 2},
		},
		{
			"curry",
			`let add = fn a => fn b => fn c => fn d =>
				a + b + c + d
			let output = add(1)(2)(3)(4)`,
			value.NumV{Value: 10},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ast := parser.ParseAll(test.input)
			env := TopInterp(ast)
			want := test.output
			got := env.Get(expr.IdC{Value: "output"})
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, wanted %v", got, want)
			}
		})
	}
}
