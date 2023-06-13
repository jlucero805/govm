package parser

import (
	"fmt"
	"strata/expr"
	"strata/lexer"
	"strconv"
)

type parser struct {
	tokens   []lexer.Token
	stmts    []expr.Expr
	position int
}

func New(tokens []lexer.Token) parser {
	p := parser{tokens: tokens}
	return p
}

func Parse(input string) expr.Expr {
	tokens := lexer.Lex(input)
	parser := New(tokens)
	parser.parse()
	return parser.stmts[0]
}

func ParseAll(input string) []expr.Expr {
	tokens := lexer.Lex(input)
	parser := New(tokens)
	for parser.position < len(parser.tokens) {
		parser.parse()
	}
	return parser.stmts
}

func (p *parser) cur() lexer.Token {
	return p.tokens[p.position]
}

func (p *parser) add(expr expr.Expr) {
	p.stmts = append(p.stmts, expr)
}

func (p *parser) eat() {
	p.position += 1
}

func (p *parser) atom() expr.Expr {
	switch p.cur().Type {
	case "num":
		float, _ := strconv.ParseFloat(p.cur().Lexeme, 64)
		strconv.ParseFloat(p.cur().Lexeme, 64)
		p.eat()
		return expr.NumC{Value: float}
	case "id":
		id := p.cur().Lexeme
		p.eat()
		return expr.IdC{Value: id}
	case "LPAR":
		p.eat()
		expression := p.fn()
		return expr.Group{Body: expression}
	default:
		return nil
	}
}

func (p *parser) call() expr.Expr {
	first := p.atom()
	if !p.end() && p.cur().Type == "LPAR" {
		p.eat()
		cargs := p.cargs()
		ecall := expr.Call{
			Proc: first,
			Args: cargs,
		}
		for !p.end() && p.cur().Type == "LPAR" {
			p.eat()
			ecargs := p.cargs()
			ecall = expr.Call{
				Proc: ecall,
				Args: ecargs,
			}
		}
		return ecall
	}
	return first
}

func (p *parser) cargs() []expr.Expr {
	res := []expr.Expr{}
	for p.cur().Type != "RPAR" {
		arg := p.fn()
		res = append(res, arg)
		if p.cur().Type == "," {
			p.eat()
		}
	}
	p.eat()
	return res
}

func (p *parser) factor() expr.Expr {
	left := p.call()
	if !p.end() && (p.cur().Type == "*" || p.cur().Type == "<") {
		p.eat()
		right := p.factor()
		binop := expr.Binop{
			Op:    "<",
			Left:  left,
			Right: right,
		}
		return binop
	}
	return left
}

func (p *parser) term() expr.Expr {
	left := p.factor()
	if !p.end() && p.cur().Type == "+" {
		p.eat()
		right := p.term()
		binop := expr.Binop{Op: "+", Left: left, Right: right}
		return binop
	}
	return left
}

func (p *parser) fn() expr.Expr {
	if p.cur().Type == "fn" {
		p.eat()
		args := p.fnargs()
		p.eat() // =>
		body := p.fn()
		return expr.LamC{
			Params: args,
			Body:   body,
		}
	}
	return p.term()
}

func (p *parser) fnargs() []expr.Expr {
	args := []expr.Expr{}
	for p.cur().Type != "=>" {
		id := p.atom()
		fmt.Println(id)
		args = append(args, id)
		if p.cur().Type == "," {
			p.eat()
		}
	}
	return args
}

func (p *parser) expr() expr.Expr {
	return p.fn()
}

func (p *parser) let() expr.Expr {
	id := p.atom()
	p.eat() // =
	val := p.expr()
	return expr.Let{
		Id:    id,
		Value: val,
	}
}

func (p *parser) iff() expr.Expr {
	cond := p.expr()
	p.eat() // then
	then := p.expr()
	p.eat() // else
	el := p.expr()
	return expr.If{
		Cond: cond,
		Then: then,
		Else: el,
	}
}

func (p *parser) parse() expr.Expr {
	if p.cur().Type == "let" {
		p.eat()
		let := p.let()
		p.add(let)
		return let
	} else if p.cur().Type == "if" {
		p.eat()
		iff := p.iff()
		p.add(iff)
		return iff
	}
	value := p.expr()
	p.add(value)
	return value
}

func (p *parser) end() bool {
	return p.position >= len(p.tokens)
}
