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
		{
			"basic equality",
			`let name = "name"
			let output = name == "name"`,
			value.BoolV{Value: true},
		},
		{
			"rectangle",
			`let Rect = fn x, y =>
				fn msg =>
					if msg == "area" then
						x * y
					else if msg == "perimeter" then
						2 * x + 2 * y
					else
						"fail"
			let rect1 = Rect(2, 4)
			let perm = rect1("perimeter")
			let area = rect1("area")
			let output = 20`,
			value.NumV{Value: 20},
		},
		{
			"advanced oop",
			`let Rect = fn x, y =>
				fn msg =>
					if get(msg, "name") == "area" then
						x * y
					else if get(msg, "name") == "z" then
						x * y * get(msg, "data")
					else
						"fail"
			let rect = Rect(2, 3)
			let output = rect({"name" : "z", "data" : 6})`,
			value.NumV{Value: 36},
		},
		{
			"basic ops",
			`let square = fn x, y => x * x
			let output = square(2)`,
			value.NumV{Value: 4},
		},
		{
			"list sum",
			`let lst = [1,2,3,4,5]
			let sum = fn l, i =>
				if i < len(l) then
					nth(lst, i) + sum(l, i + 1)
				else 
					0
			let output = 15`,
			value.NumV{Value: 15},
		},
		{
			"list sum with tail",
			`let lst = [1, 2, 3, 4, 5]
			let sum = fn l =>
				if i < len(l) then
					head(l) + sum(tail(l))
				else
					0
			let output = 15`,
			value.NumV{Value: 15},
		},
		{
			"fibonacci",
			`let fib = fn n =>
				if n == 1 then
					1
				else if n == 2 then
					1
				else
					fib(n - 1) + fib(n - 2)
			let output = fib(20)`,
			value.NumV{Value: 6765},
		},
		{
			"btree",
			`let tree = {
				"value": 1,
				"left": {
					"value": 2,
					"left": "nil",
					"right": {
						"value": 4,
						"left": {
							"value": 5,
							"left": "nil",
							"right": "nil"
						},
						"right": "nil"
					}
				},
				"right": {
					"value": 3,
					"left": "nil",
					"right": "nil"
				},
			}
			let addTree = fn tree =>
				if tree == "nil" then 0
				else get(tree, "value")
					+ addTree(get(tree, "left"))
					+ addTree(get(tree, "right"))
			let output = addTree(tree)`,
			value.NumV{Value: 15},
		},
		{
			"nested funcs",
			`let add = fn a, b => do
				let add1 = fn x => x + 1
				add1(a) + add1(b)
			end
			let output = add(1, 2)`,
			value.NumV{Value: 5},
		},
		{
			"generate List",
			`let genList = fn n => do
				let gen = fn n, acc =>
					if n > 0 then gen(n - 1, push(acc, n))
					else acc
				gen(n, [])
			end
			let output = genList(5)`,
			value.ListV{
				Values: []value.Value{
					value.NumV{Value: 5},
					value.NumV{Value: 4},
					value.NumV{Value: 3},
					value.NumV{Value: 2},
					value.NumV{Value: 1},
				},
			},
		},
		{
			"is even",
			`let isEven = fn n =>
				if n >= 2 then isEven(n - 2)
				else n == 0
			let output = [
				isEven(2),
				isEven(11),
				isEven(52),
				isEven(103)
			]`,
			value.ListV{
				Values: []value.Value{
					value.BoolV{Value: true},
					value.BoolV{Value: false},
					value.BoolV{Value: true},
					value.BoolV{Value: false},
				},
			},
		},
		{
			"lexical scoping",
			`let x = 1
			let num = fn => do
				let x = 2
				let y = do
					let x = 3
					x
				end
				[x, y]
			end
			let output = num()`,
			value.ListV{
				Values: []value.Value{
					value.NumV{Value: 2},
					value.NumV{Value: 3},
				},
			},
		},
		{
			"matrix",
			`let matrix = [
				[1, 2, 3, 4, 5],
				[1, 2, 3, 4, 5],
				[1, 2, 3, 4, 5],
				[1, 2, 3, 4, 5],
				[1, 2, 3, 4, 5],
			]
			let sum = fn lst =>
				if len(lst) > 0 then
					head(lst) + sum(tail(lst))
				else
					0
			let output = Enum.map(matrix, sum)`,
			value.ListV{
				Values: []value.Value{
					value.NumV{Value: 15},
					value.NumV{Value: 15},
					value.NumV{Value: 15},
					value.NumV{Value: 15},
					value.NumV{Value: 15},
				},
			},
		},
		{
			"stdlib",
			`let output = Enum.map([1, 2, 3], (fn x => x + 1))`,
			value.ListV{
				Values: []value.Value{
					value.NumV{Value: 2},
					value.NumV{Value: 3},
					value.NumV{Value: 4},
				},
			},
		},
		{
			"filter",
			"let output = Enum.filter([1, 2, 3, 4], (fn x => x == 2))",
			value.ListV{
				Values: []value.Value{
					value.NumV{Value: 2},
				},
			},
		},
		{
			"foldl",
			`let output = Enum.foldl([1, 2, 3, 4, 5],
				(fn x, acc => x + acc),
				0)`,
			value.NumV{Value: 15},
		},
		{
			". operator",
			`let m = {hey: 69}
			let output = m.hey`,
			value.NumV{Value: 69},
		},
		{
			"nested object dot operator",
			`let m = {one: {two: {three: 4}, rand: 3}}
			let output = m.one.two.three`,
			value.NumV{Value: 4},
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

func TestOrderOfOps(t *testing.T) {
	tests := []struct {
		input1 string
		input2 string
	}{
		{
			`let x = 2 * 2 + 4 * 4`,
			`let x = (2 * 2) + (4 * 4)`,
		},
		{
			`let x = 2 * 2 + 4 * 4`,
			`let x = ((2 * 2) + (4 * 4))`,
		},
		{
			`let x = 2 * 2 / 4 + 3 - 1 * 3 + 4 > 2 / 2 + 1 * 3`,
			`let x = ((((2 * 2) / 4) + 3 - (1 * 3) + 4) > ((2 / 2) + (1 * 3)))`,
		},
	}
	for _, test := range tests {
		t.Run("order of ops", func(t *testing.T) {
			ast1 := parser.ParseAll(test.input1)
			env1 := TopInterp(ast1)

			ast2 := parser.ParseAll(test.input2)
			env2 := TopInterp(ast2)
			if !reflect.DeepEqual(env2, env1) {
				t.Errorf("got %v, wanted %v", env2, env1)
			}
		})
	}
}
