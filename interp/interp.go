package interp

import (
	"strata/expr"
	"strata/parser"
	"strata/value"
)

func Interp(node expr.Expr, env *value.Env) (value.Value, *value.Env) {
	switch node.(type) {
	case expr.DoC:
		do := node.(expr.DoC)
		childEnv := value.Extend(env, value.NewChild())
		es := []value.Value{}
		for _, e := range do.Exprs {
			ee, en := Interp(e, childEnv)
			childEnv = en
			es = append(es, ee)
		}
		return es[len(es)-1], env
	case expr.MapC:
		em := node.(expr.MapC)
		m := make(map[value.Value]value.Value)
		env = env
		for key, val := range em.Binds {
			var id value.Value
			switch key := key.(type) {
			case expr.IdC:
				id = value.StrV{Value: key.Value}
			default:
				id, env = Interp(key, env)
			}
			mval, env := Interp(val, env)
			env = env
			m[id] = mval
		}
		return value.MapV{Binds: m}, env
	case expr.ListC:
		list := node.(expr.ListC)
		listV := []value.Value{}
		newEnv := env
		for _, val := range list.Values {
			interpedVal, e := Interp(val, env)
			newEnv = e
			listV = append(listV, interpedVal)
		}
		return value.ListV{Values: listV}, newEnv
	case expr.StrC:
		return value.StrV{Value: node.(expr.StrC).Value}, env
	case expr.IdC:
		id := node.(expr.IdC)
		return env.Get(id), env
	case expr.NumC:
		return value.NumV{Value: node.(expr.NumC).Value}, env
	case expr.Binop:
		switch node.(expr.Binop).Op {
		case ".":
			binop := node.(expr.Binop)
			var right expr.Expr
			switch r := binop.Right.(type) {
			case expr.IdC:
				right = expr.StrC{Value: r.Value}
			default:
				right = binop.Right
			}
			val, env := Interp(expr.Call{
				Proc: expr.IdC{Value: "get"},
				Args: []expr.Expr{
					binop.Left,
					right,
				},
			}, env)
			return val, env
		case "+":
			binop := node.(expr.Binop)
			left, env1 := Interp(binop.Left, env)
			right, env2 := Interp(binop.Right, env1)
			return value.NumV{Value: left.(value.NumV).Value + right.(value.NumV).Value}, env2
		case "*":
			binop := node.(expr.Binop)
			left, env1 := Interp(binop.Left, env)
			right, env2 := Interp(binop.Right, env1)
			return value.NumV{Value: left.(value.NumV).Value * right.(value.NumV).Value}, env2
		case "/":
			binop := node.(expr.Binop)
			left, env1 := Interp(binop.Left, env)
			right, env2 := Interp(binop.Right, env1)
			return value.NumV{Value: left.(value.NumV).Value / right.(value.NumV).Value}, env2
		case "-":
			binop := node.(expr.Binop)
			left, env1 := Interp(binop.Left, env)
			right, env2 := Interp(binop.Right, env1)
			return value.NumV{Value: left.(value.NumV).Value - right.(value.NumV).Value}, env2
		case "<":
			binop := node.(expr.Binop)
			left, env1 := Interp(binop.Left, env)
			right, env2 := Interp(binop.Right, env1)
			return value.BoolV{Value: left.(value.NumV).Value < right.(value.NumV).Value}, env2
		case ">":
			binop := node.(expr.Binop)
			left, env1 := Interp(binop.Left, env)
			right, env2 := Interp(binop.Right, env1)
			return value.BoolV{Value: left.(value.NumV).Value > right.(value.NumV).Value}, env2
		case "<=":
			binop := node.(expr.Binop)
			left, env := Interp(binop.Left, env)
			right, env := Interp(binop.Right, env)
			return value.BoolV{Value: left.(value.NumV).Value <= right.(value.NumV).Value}, env
		case ">=":
			binop := node.(expr.Binop)
			left, env := Interp(binop.Left, env)
			right, env := Interp(binop.Right, env)
			return value.BoolV{Value: left.(value.NumV).Value >= right.(value.NumV).Value}, env
		case "==":
			binop := node.(expr.Binop)
			left, env := Interp(binop.Left, env)
			right, env := Interp(binop.Right, env)
			return value.BoolV{Value: left.Equals(right)}, env
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
		env.Set(let.Id.(expr.IdC), value.NilV{})
		val, e := Interp(let.Value, env)
		e.Set(let.Id.(expr.IdC), val)
		return value.NilV{}, e
	case expr.Call:
		call := node.(expr.Call)
		proc, newEnv := Interp(call.Proc, env)
		switch proc := proc.(type) {
		case value.LamV:
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
					Id:    proc.Params[i].(expr.IdC),
					Value: args[i],
				})
			}
			callEnvMap := make(map[expr.IdC]value.Value)
			for _, b := range binds {
				callEnvMap[b.Id] = b.Value
			}
			callEnv := &value.Env{Binds: callEnvMap}
			extendedEnv := value.Extend(proc.Closure, callEnv)
			body, _ := Interp(proc.Body, extendedEnv)
			return body, cEnv
		case value.PrimV:
			args := []value.Value{}
			cEnv := newEnv
			var v value.Value
			for _, val := range call.Args {
				v, cEnv = Interp(val, env)
				args = append(args, v)
			}
			switch proc.Name {
			case "MapGet":
				return args[0].(value.MapV).Binds[args[1]], cEnv
			case "nth":
				return args[0].(value.ListV).Values[int(args[1].(value.NumV).Value)], cEnv
			case "len":
				return value.NumV{Value: float64(len(args[0].(value.ListV).Values))}, cEnv
			case "head":
				list := args[0].(value.ListV).Values
				if len(list) <= 0 {
					panic("list is length 0")
				} else {
					return list[0], cEnv
				}
			case "tail":
				list := args[0].(value.ListV).Values
				if len(list) <= 0 {
					panic("list is length 0")
				}
				return value.ListV{Values: list[1:]}, cEnv
			case "push":
				list := args[0].(value.ListV).Values
				return value.ListV{Values: append(list, args[1])}, cEnv
			default:
				panic("prim doesn't exist")
			}

		default:
			panic("")
		}
	case expr.If:
		iff := node.(expr.If)
		i, e := Interp(iff.Cond, env)
		cond := i.(value.BoolV).Value
		if cond {
			val, env := Interp(iff.Then, e)
			return val, env
		}
		return Interp(iff.Else, e)
	case expr.Group:
		var group = node.(expr.Group)
		return Interp(group.Body, env)
	default:
		return value.NumV{Value: 0}, env
	}
}

func TopInterp(nodes []expr.Expr) *value.Env {
	env := value.New()
	_, env = Interp(parser.Parse(`
	let Enum = {
		map: (fn lst, f => do
			let go = fn lst, f, acc =>
				if len(lst) > 0 then
					go(tail(lst), f, push(acc, f(head(lst))))
				else
					acc
			go(lst, f, [])
		end),
	 filter: (fn lst, f => do
		let go = fn lst, f, acc =>
			if len(lst) > 0 then do
				let accVals = if f(head(lst)) then
					push(acc, head(lst))
				else
					acc
				go(tail(lst), f, accVals)
			end else
				acc
		go(lst, f, [])
	end),
	foldl: (fn col, f, acc => do
		let iter = fn col, acc => 
			if len(col) > 0 then
				iter(tail(col), f(head(col), acc))
			else
				acc
		iter(col, acc)
	end)
	}`), env)
	for _, node := range nodes {
		_, newEnv := Interp(node, env)
		env = newEnv
	}
	return env
}
