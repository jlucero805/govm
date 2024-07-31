package compile

import (
	"reflect"
	"strata/parser"
	"strings"
	"testing"
)

func TestOrderOfOps(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{`x`, `load x`},
		{`1`, `push 1`},
		{`1 + 2`, "push 2\npush 1\nadd"},
		{`1 + 2 * 4`, "push 4\npush 2\nmul\npush 1\nadd"},
		{`call(1, 2, x)`, "load x\npush 2\npush 1\nload call\ncall"},
		{
			`let sum = fn n =>
				if 0 < n
					then n + sum(n - 1)
					else 0
			sum(5)`,
			strings.Join(
				[]string{
					"fundef",
					"store n",
					"load n",
					"push 0",
					"lt",
					"jmprif 4",
					"push 0",
					"jmpr 11",
					"push 1",
					"load n",
					"sub",
					"load sum",
					"call",
					"load n",
					"add",
					"ret",
					"funend",
					"store sum",
				},
				"\n",
			),
		},
		{
			`fn x => x + 1`,
			strings.Join(
				[]string{
					"fundef",
					"store x",
					"push 1",
					"load x",
					"add",
					"ret",
					"funend",
				},
				"\n",
			),
		},
		{
			`let lt1 =
				fn x => if x < 100 then 1 else 0`,
			strings.Join(
				[]string{
					"fundef",
					"store x",
					"push 100",
					"load x",
					"lt",
					"jmprif 4",
					"push 0",
					"jmpr 2",
					"push 1",
					"ret",
					"funend",
					"store lt1",
				},
				"\n",
			),
		},
		{
			`let add =
				fn x, y => x + 1`,
			strings.Join(
				[]string{
					"fundef",
					"store x",
					"store y",
					"push 1",
					"load x",
					"add",
					"ret",
					"funend",
					"store add",
				},
				"\n",
			),
		},
		{
			`let add =
				fn x, y => do
					let add1 = fn x => x + 1
					add1(x) + add1(y)
				end`,
			strings.Join(
				[]string{
					"fundef",
					"store x",
					"store y",
					"envnew",
					"fundef",
					"store x",
					"push 1",
					"load x",
					"add",
					"ret",
					"funend",
					"store add1",
					"load y",
					"load add1",
					"call",
					"load x",
					"load add1",
					"call",
					"add",
					"envend",
					"ret",
					"funend",
					"store add",
				},
				"\n",
			),
		},
		{
			`let x = do
				let y = 2
				y
			end`,
			strings.Join(
				[]string{
					"envnew",
					"push 2",
					"store y",
					"load y",
					"envend",
					"store x",
				},
				"\n",
			),
		},
	}
	for _, test := range tests {
		t.Run("order of ops", func(t *testing.T) {
			ast := parser.Parse(test.input)
			got := Compile(ast)

			if !reflect.DeepEqual(got, test.output) {
				t.Errorf("got %v, wanted %v", got, test.output)
			}
		})
	}
}
