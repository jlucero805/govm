package value

import "strata/expr"

type LamV struct {
	Params  []expr.Expr
	Body    expr.Expr
	Closure *Env
}

func (lam LamV) value() {}

func (l LamV) Equals(v Value) bool {
	return false
}
