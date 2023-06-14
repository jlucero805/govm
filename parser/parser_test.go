package parser

import (
	"reflect"
	"strata/expr"
	"strata/lexer"
	"testing"
)

func TestComparison(t *testing.T) {
	input := "21 > 2"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.Binop{
			Op:    ">",
			Left:  expr.NumC{Value: 21},
			Right: expr.NumC{Value: 2},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("")
	}
}

func Test(t *testing.T) {
	input := "12 + 21 < 2"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.Binop{
			Op: "<",
			Left: expr.Binop{
				Op:    "+",
				Left:  expr.NumC{Value: 12},
				Right: expr.NumC{Value: 21},
			},
			Right: expr.NumC{Value: 2},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("")
	}
}

func Test2(t *testing.T) {
	input := "fn a, b => a + b"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.LamC{
			Params: []expr.Expr{
				expr.IdC{Value: "a"},
				expr.IdC{Value: "b"},
			},
			Body: expr.Binop{
				Op:    "+",
				Left:  expr.IdC{Value: "a"},
				Right: expr.IdC{Value: "b"},
			},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("")
	}
}

func Test22(t *testing.T) {
	input := "fn a, b => a + b"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.LamC{
			Params: []expr.Expr{
				expr.IdC{Value: "a"},
				expr.IdC{Value: "b"},
			},
			Body: expr.Binop{
				Op:    "+",
				Left:  expr.IdC{Value: "a"},
				Right: expr.IdC{Value: "b"},
			},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("")
	}
}

func TestCalla(t *testing.T) {
	input := "(1 + 2)"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.Group{
			Body: expr.Binop{
				Op:    "+",
				Left:  expr.NumC{Value: 1},
				Right: expr.NumC{Value: 2},
			},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("%#v %#v", expect, parser.stmts)
	}
}

func TestCall(t *testing.T) {
	input := "call(1 , 2)"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.Call{
			Proc: expr.IdC{Value: "call"},
			Args: []expr.Expr{
				expr.NumC{Value: 1},
				expr.NumC{Value: 2},
			},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("%#v %#v", expect, parser.stmts)
	}
}

func TestTsdf(t *testing.T) {
	input := "if x < 0 then 0 else 1"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.If{
			Cond: expr.Binop{
				Op:    "<",
				Left:  expr.IdC{Value: "x"},
				Right: expr.NumC{Value: 0},
			},
			Then: expr.NumC{Value: 0},
			Else: expr.NumC{Value: 1},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("%#v %#v", expect, parser.stmts)
	}
}

func TestTsdfasf(t *testing.T) {
	input := "let x = if x < 0 then 0 else 1"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.Let{
			Id: expr.IdC{Value: "x"},
			Value: expr.If{
				Cond: expr.Binop{
					Op:    "<",
					Left:  expr.IdC{Value: "x"},
					Right: expr.NumC{Value: 0},
				},
				Then: expr.NumC{Value: 0},
				Else: expr.NumC{Value: 1},
			},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("%#v %#v", expect, parser.stmts)
	}
}
func TestTskjdfasf(t *testing.T) {
	input := "let x = x.y.z"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.Let{
			Id: expr.IdC{Value: "x"},
			Value: expr.Binop{
				Op: ".",
				Left: expr.Binop{
					Op:    ".",
					Left:  expr.IdC{Value: "x"},
					Right: expr.IdC{Value: "y"},
				},
				Right: expr.IdC{Value: "z"},
			},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("\n%#v\n \n%#v\n", expect, parser.stmts)
	}
}

func TestTsdfaasdfsf(t *testing.T) {
	input := "let x = fn n => if n < 0 then 0 else 1"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	expect := []expr.Expr{
		expr.Let{
			Id: expr.IdC{Value: "x"},
			Value: expr.LamC{
				Params: []expr.Expr{
					expr.IdC{Value: "n"},
				},
				Body: expr.If{
					Cond: expr.Binop{
						Op:    "<",
						Left:  expr.IdC{Value: "n"},
						Right: expr.NumC{Value: 0},
					},
					Then: expr.NumC{Value: 0},
					Else: expr.NumC{Value: 1},
				},
			},
		},
	}
	if !reflect.DeepEqual(expect, parser.stmts) {
		t.Errorf("%#v", parser.stmts)
	}
}
