package compile

import (
	"fmt"
	"strata/expr"
	"strings"
)

func Compile(node expr.Expr) string {
	switch node := node.(type) {
	case expr.IdC:
		return fmt.Sprintf("load %v", node.Value)
	case expr.NumC:
		return fmt.Sprintf("push %v", node.Value)
	case expr.Binop:
		var op string
		switch node.Op {
		case "+":
			op = "add"
		case "-":
			op = "sub"
		case "*":
			op = "mul"
		case "/":
			op = "div"
		case ">":
			op = "gt"
		case "<":
			op = "lt"
		case "<=":
			op = "lteq"
		case ">=":
			op = "gteq"
		case "==":
			op = "eq"
		}
		emit := []string{
			Compile(node.Right),
			Compile(node.Left),
			op,
		}
		return strings.Join(emit, "\n")
	case expr.Call:
		args := []string{}
		for i := len(node.Args) - 1; i >= 0; i-- {
			args = append(args, Compile(node.Args[i]))
		}
		args = append(args, Compile(node.Proc))
		args = append(args, "call")
		return strings.Join(args, "\n")
	case expr.LamC:
		emit := []string{
			"fundef",
		}
		for _, param := range node.Params {
			instr := fmt.Sprintf("store %v", param.(expr.IdC).Value)
			emit = append(emit, instr)
		}
		emit = append(emit, Compile(node.Body))
		emit = append(emit, "ret")
		emit = append(emit, "funend")
		return strings.Join(emit, "\n")
	case expr.Let:
		emit := []string{
			Compile(node.Value),
			fmt.Sprintf("store %v", node.Id.(expr.IdC).Value),
		}
		return strings.Join(emit, "\n")
	case expr.DoC:
		emit := []string{}
		emit = append(emit, "envnew")
		for _, stmt := range node.Exprs {
			instr := Compile(stmt)
			emit = append(emit, instr)
		}
		emit = append(emit, "envend")
		return strings.Join(emit, "\n")
	case expr.If:
		cond := Compile(node.Cond)
		then := Compile(node.Then)
		el := Compile(node.Else)
		elLength := len(strings.Fields(el)) + 2
		jmprif := fmt.Sprintf("jmprif %v", elLength)
		thenLength := len(strings.Fields(then))
		jmpr := fmt.Sprintf("jmpr %v", thenLength)
		emit := []string{
			cond,
			jmprif,
			el,
			jmpr,
			then,
		}
		return strings.Join(emit, "\n")
	default:
		panic(fmt.Sprintf("Value: %#v doesn't exist yet.", node))
	}
}
