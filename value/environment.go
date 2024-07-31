package value

import "strata/expr"

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
	m := make(map[expr.IdC]Value)
	m[expr.IdC{Value: "get"}] = PrimV{Name: "MapGet"}
	m[expr.IdC{Value: "nth"}] = PrimV{Name: "nth"}
	m[expr.IdC{Value: "len"}] = PrimV{Name: "len"}
	m[expr.IdC{Value: "tail"}] = PrimV{Name: "tail"}
	m[expr.IdC{Value: "head"}] = PrimV{Name: "head"}
	m[expr.IdC{Value: "push"}] = PrimV{Name: "push"}
	m[expr.IdC{Value: "print"}] = PrimV{Name: "print"}
	env := &Env{Binds: m, Parent: true}
	return env
}

func NewChild() *Env {
	m := make(map[expr.IdC]Value)
	return &Env{Binds: m}
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
