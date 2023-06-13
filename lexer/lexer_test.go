package lexer

import (
	"fmt"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	lexer := New("123 + (if) bro ")
	lexer.lexTokens()
	tokens := lexer.tokens
	fmt.Println(lexer)
	if !reflect.DeepEqual(tokens,
		[]Token{
			newToken("num", "123"),
			newToken("+", "+"),
			newToken("LPAR", "("),
			newToken("if", "if"),
			newToken("RPAR", ")"),
			newToken("id", "bro")}) {
		t.Errorf("")
	}
}

func Test2(t *testing.T) {
	lexer := New("12 + 21 * 2")
	lexer.lexTokens()
	tokens := lexer.tokens
	fmt.Println(lexer)
	if !reflect.DeepEqual(tokens,
		[]Token{
			newToken("num", "12"),
			newToken("+", "+"),
			newToken("num", "21"),
			newToken("*", "*"),
			newToken("num", "2")}) {
		t.Errorf("")
	}
}

func Test3(t *testing.T) {
	lexer := New("fn a, b => a + b")
	lexer.lexTokens()
	tokens := lexer.tokens
	fmt.Println(lexer)
	if !reflect.DeepEqual(tokens,
		[]Token{
			newToken("fn", "fn"),
			newToken("id", "a"),
			newToken(",", ","),
			newToken("id", "b"),
			newToken("=>", "=>"),
			newToken("id", "a"),
			newToken("+", "+"),
			newToken("id", "b"),
		}) {
		t.Errorf("")
	}
}

func Test4(t *testing.T) {
	lexer := New("call(1, 2)")
	lexer.lexTokens()
	tokens := lexer.tokens
	fmt.Println(lexer)
	if !reflect.DeepEqual(tokens,
		[]Token{
			newToken("id", "call"),
			newToken("LPAR", "("),
			newToken("num", "1"),
			newToken(",", ","),
			newToken("num", "2"),
			newToken("RPAR", ")"),
		}) {
		t.Errorf("")
	}
}
