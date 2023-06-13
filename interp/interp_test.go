package interp

import (
	"fmt"
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
	ast := parser.Parse("if 1 < 0 then 1 else 0")
	got, _ := Interp(ast, value.New())
	want := value.NumV{Value: 0}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestTestTest(t *testing.T) {
	ast := parser.ParseAll("let x = 1 + 1")
	got := TopInterp(ast)

	fmt.Println(got)
	t.Logf("%v", got)
}
