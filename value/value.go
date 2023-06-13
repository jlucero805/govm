package value

import (
	"strata/expr"
)

type Value interface {
	value()
}

type NumV struct {
	Value float64
}

func (num NumV) value() {}

type NilV struct{}

func (nil NilV) value() {}

type LamV struct {
	Params  []expr.Expr
	Body    expr.Expr
	Closure *Env
}

func (lam LamV) value() {}

type BoolV struct {
	Value bool
}

func (b BoolV) value() {}

type Env struct {
	Binds    map[expr.IdC]Value
	Previous *Env
	Parent   bool
}

type Bind struct {
	Id     expr.IdC
	Value  Value
	Parent bool
}

func (b *Bind) SetVal(value Value) {
	b.Value = value
}

func New() *Env {
	return &Env{Binds: make(map[expr.IdC]Value), Parent: true}
}

func (e *Env) Set(id expr.IdC, value Value) {
	e.Binds[id] = value
}

func (e *Env) Get(id expr.IdC) Value {
	if val, ok := e.Binds[id]; ok {
		return val
	}
	if e.Parent {
		panic(id.Value)
	}
	return e.Previous.Get(id)
}

func Extend(e1 *Env, e2 *Env) *Env {
	e2.Previous = e1
	return e2
}
