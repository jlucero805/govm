package parser

import (
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

func (p *parser) mmap() expr.Expr {
	p.eat() // {
	m := make(map[expr.Expr]expr.Expr)
	for !p.end() && p.cur().Type != "}" {
		id := p.atom() // id can only be literals
		p.eat()        // :
		val := p.expr()
		if p.cur().Type == "," {
			p.eat()
		}
		m[id] = val
	}
	p.eat() // }
	return expr.MapC{Binds: m}
}

func (p *parser) list() expr.Expr {
	p.eat() // [
	var vals = []expr.Expr{}
	for !p.end() && p.cur().Type != "]" {
		val := p.expr()
		vals = append(vals, val)
		if p.cur().Type == "," {
			p.eat()
		}
	}
	p.eat() // ]
	return expr.ListC{Values: vals}
}

func (p *parser) atom() expr.Expr {
	switch p.cur().Type {
	case "string":
		string := p.cur().Lexeme
		p.eat()
		return expr.StrC{Value: string}
	case "num":
		float, _ := strconv.ParseFloat(p.cur().Lexeme, 64)
		strconv.ParseFloat(p.cur().Lexeme, 64)
		p.eat()
		return expr.NumC{Value: float}
	case "id":
		id := p.cur().Lexeme
		p.eat()
		return expr.IdC{Value: id}
	case "{":
		return p.mmap()
	case "[":
		return p.list()
	case "LPAR":
		p.eat()
		expression := p.fn()
		p.eat()
		return expr.Group{Body: expression}
	default:
		return nil
	}
}

func (p *parser) dot() expr.Expr {
	left := p.atom()

	if p.end() {
		return left
	}

	isDot := p.cur().Type == "."

	if isDot {
		exprs := []expr.Expr{}
		for !p.end() && p.cur().Type == "." {
			p.eat()
			atom := p.atom()
			exprs = append(exprs, atom)
		}
		dotExpr := expr.Binop{
			Op: ".",
			Left: left,
			Right: exprs[0],
		}
		exprs = exprs[1:]
		for len(exprs) > 0 {
			temp := exprs[0]
			exprs = exprs[1:]
			dotExpr = expr.Binop{
				Op: ".",
				Left: dotExpr,
				Right: temp,
			}
		}
		return dotExpr
	}

	return left
}

func (p *parser) call() expr.Expr {
	first := p.dot()
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

	if p.end() {
		return left
	}

	isStar := p.cur().Type == "*"
	isSlash := p.cur().Type == "/"

	if isStar || isSlash {
		op := p.cur().Type
		p.eat()
		right := p.factor()
		binop := expr.Binop{
			Op:    op,
			Left:  left,
			Right: right,
		}
		return binop
	}
	return left
}

func (p *parser) term() expr.Expr {
	left := p.factor()

	if p.end() {
		return left
	}

	isPlus := p.cur().Type == "+"
	isMinus := p.cur().Type == "-"

	if isPlus || isMinus {
		op := p.cur().Type
		p.eat()
		right := p.term()
		binop := expr.Binop{Op: op, Left: left, Right: right}
		return binop
	}
	return left
}

func (p *parser) comparison() expr.Expr {
	left := p.term()

	if p.end() {
		return left
	}

	isGreater := p.cur().Type == ">"
	isLessThan := p.cur().Type == "<"
	isGreaterEqual := p.cur().Type == ">="
	isLessEqual := p.cur().Type == "<="

	if isGreater || isLessThan || isGreaterEqual || isLessEqual {
		op := p.cur().Type
		p.eat() // operator
		right := p.comparison()
		binop := expr.Binop{Op: op, Left: left, Right: right}
		return binop
	}

	return left
}

func (p *parser) equality() expr.Expr {
	left := p.comparison()

	if p.end() {
		return left
	}

	isEqualEqual := p.cur().Type == "=="

	if isEqualEqual {
		op := p.cur().Type
		p.eat()
		right := p.equality()
		return expr.Binop{Op: op, Left: left, Right: right}
	}

	return left
}

func (p *parser) fn() expr.Expr {
	if p.cur().Type == "fn" {
		p.eat()
		args := p.fnargs()
		p.eat() // =>
		body := p.expr()
		return expr.LamC{
			Params: args,
			Body:   body,
		}
	}
	return p.equality()
}

func (p *parser) fnargs() []expr.Expr {
	args := []expr.Expr{}
	for p.cur().Type != "=>" {
		id := p.atom()
		args = append(args, id)
		if p.cur().Type == "," {
			p.eat()
		}
	}
	return args
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

func (p *parser) expr() expr.Expr {
	if p.cur().Type == "let" {
		p.eat()
		let := p.let()
		return let
	} else if p.cur().Type == "if" {
		p.eat()
		iff := p.iff()
		return iff
	} else if p.cur().Type == "do" {
		p.eat()
		es := []expr.Expr{}
		for p.cur().Type != "end" {
			e := p.expr()
			es = append(es, e)
		}
		p.eat()
		return expr.DoC{Exprs: es}
	}
	value := p.fn()
	return value
}

func (p *parser) parse() expr.Expr {
	value := p.expr()
	p.add(value)
	return value
}

func (p *parser) end() bool {
	return p.position >= len(p.tokens)
}
