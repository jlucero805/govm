package value

import (
	"reflect"
	"strata/expr"
)

type Value interface {
	value()
	Equals(Value) bool
}

type MapV struct {
	Binds map[Value]Value
}

func (str MapV) value() {}

func (str MapV) Equals(v Value) bool {
	switch node := v.(type) {
	case MapV:
		return reflect.DeepEqual(node.Binds, str.Binds)
	default:
		return false
	}
}

type ListV struct {
	Values []Value
}

func (str ListV) value() {}

func (str ListV) Equals(v Value) bool {
	switch node := v.(type) {
	case ListV:
		return reflect.DeepEqual(node.Values, str.Values)
	default:
		return false
	}
}

type StrV struct {
	Value string
}

func (str StrV) value() {}

func (str StrV) Equals(v Value) bool {
	switch node := v.(type) {
	case StrV:
		return str.Value == node.Value
	default:
		return false
	}
}

type NumV struct {
	Value float64
}

func (num NumV) value() {}

func (n NumV) Equals(v Value) bool {
	switch node := v.(type) {
	case NumV:
		return n.Value == node.Value
	default:
		return false
	}
}

type NilV struct{}

func (nil NilV) value() {}

func (n NilV) Equals(v Value) bool {
	switch v.(type) {
	case NilV:
		return true
	default:
		return false
	}
}

type LamV struct {
	Params  []expr.Expr
	Body    expr.Expr
	Closure *Env
}

func (lam LamV) value() {}

func (l LamV) Equals(v Value) bool {
	return false
}

type BoolV struct {
	Value bool
}

func (b BoolV) Equals(v Value) bool {
	switch node := v.(type) {
	case BoolV:
		return node.Value == b.Value
	default:
		return false
	}
}

func (b BoolV) value() {}

type PrimV struct {
	Name string
}

func (b PrimV) Equals(v Value) bool {
	switch v.(type) {
	default:
		return false
	}
}

func (b PrimV) value() {}

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
