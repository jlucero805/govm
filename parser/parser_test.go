package parser

import (
	"fmt"
	"reflect"
	"strata/expr"
	"strata/lexer"
	"testing"
)

func Test(t *testing.T) {
	input := "12 + 21 * 2"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	fmt.Println(parser.stmts)
	expect := []expr.Expr{
		expr.Binop{
			Op:   "+",
			Left: expr.NumC{Value: 12},
			Right: expr.Binop{
				Op:    "*",
				Left:  expr.NumC{Value: 21},
				Right: expr.NumC{Value: 2},
			},
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
	fmt.Println(parser.stmts)
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

func TestCall(t *testing.T) {
	input := "call(1 , 2)"
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	fmt.Println(parser.stmts)
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
