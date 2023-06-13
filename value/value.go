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
	Closure Env
}

func (lam LamV) value() {}

type BoolV struct {
	Value bool
}

func (b BoolV) value() {}

type Env struct {
	Binds []*Bind
}

type Bind struct {
	Id    expr.IdC
	Value Value
}

func (b *Bind) SetVal(value Value) {
	b.Value = value
}

func New() Env {
	return Env{Binds: []*Bind{}}
}

func (e *Env) set(id expr.IdC, value Value) {
	e.Binds = append(e.Binds, &Bind{Id: id, Value: value})
}

func Set(env Env, id expr.IdC, value Value) Env {
	return Env{
		Binds: append(env.Binds, &Bind{Id: id, Value: value}),
	}
}

func (e *Env) Get(id expr.IdC) Value {
	for i := len(e.Binds) - 1; i >= 0; i = i - 1 {
		if e.Binds[i].Id == id {
			return e.Binds[i].Value
		}
	}
	return nil
}

func (e *Env) GetBind(id expr.IdC) *Bind {
	for i := len(e.Binds) - 1; i >= 0; i = i - 1 {
		if e.Binds[i].Id == id {
			return e.Binds[i]
		}
	}
	return nil
}

func Extend(e1 Env, e2 Env) Env {
	return Env{Binds: append(e1.Binds, e2.Binds...)}
}
