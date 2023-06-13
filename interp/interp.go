package interp

import (
	"fmt"
	"strata/expr"
	"strata/value"
)

func Interp(node expr.Expr, env value.Env) (value.Value, value.Env) {
	switch node.(type) {
	case expr.IdC:
		id := node.(expr.IdC)
		return env.Get(id), env
	case expr.NumC:
		return value.NumV{Value: node.(expr.NumC).Value}, env
	case expr.Binop:
		switch node.(expr.Binop).Op {
		case "+":
			binop := node.(expr.Binop)
			left, env1 := Interp(binop.Left, env)
			right, env2 := Interp(binop.Right, env1)
			fmt.Println(left)
			fmt.Println(right)
			return value.NumV{Value: left.(value.NumV).Value + right.(value.NumV).Value}, env2
		case "<":
			binop := node.(expr.Binop)
			left, env1 := Interp(binop.Left, env)
			right, env2 := Interp(binop.Right, env1)
			return value.BoolV{Value: left.(value.NumV).Value < right.(value.NumV).Value}, env2
		default:
			return value.NumV{Value: 0}, env
		}
	case expr.LamC:
		lam := node.(expr.LamC)
		return value.LamV{
			Params:  lam.Params,
			Body:    lam.Body,
			Closure: env,
		}, env
	case expr.Let:
		let := node.(expr.Let)
		val, e := Interp(let.Value, env)
		newEnv := value.Set(e, let.Id.(expr.IdC), val)
		return value.NilV{}, newEnv
	case expr.Call:
		call := node.(expr.Call)
		proc, newEnv := Interp(call.Proc, env)
		args := []value.Value{}
		cEnv := newEnv
		var v value.Value
		for _, val := range call.Args {
			v, cEnv = Interp(val, env)
			args = append(args, v)
		}
		binds := []*value.Bind{}
		for i := 0; i < len(call.Args); i += 1 {
			binds = append(binds, &value.Bind{
				Id:    proc.(value.LamV).Params[i].(expr.IdC),
				Value: args[i],
			})
		}
		extendedEnv := value.Extend(proc.(value.LamV).Closure, value.Env{Binds: binds})
		fmt.Println(extendedEnv)
		body, _ := Interp(proc.(value.LamV).Body, extendedEnv)
		return body, cEnv
	case expr.If:
		iff := node.(expr.If)
		i, eenv := Interp(iff.Cond, env)
		cond := i.(value.BoolV).Value
		if cond {
			return Interp(iff.Then, eenv)
		}
		return Interp(iff.Else, eenv)
	default:
		return value.NumV{Value: 0}, env
	}
}

func TopInterp(nodes []expr.Expr) value.Env {
	env := value.New()
	for _, node := range nodes {
		_, newEnv := Interp(node, env)
		env = newEnv
	}
	return env
}
