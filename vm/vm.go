package vm

import (
	"fmt"
	compile "strata/compiler"
)

const LARGEST int64 = 9_223_372_036_854_775_807

type Stack struct {
	Values []int64
}

func NewStack() *Stack {
	return &Stack{Values: []int64{}}
}

func (s *Stack) Push(v int64) {
	s.Values = append(s.Values, v)
}

func (s *Stack) Pop() int64 {
	l := len(s.Values)
	value := s.Values[l-1]
	s.Values = s.Values[:l-1]
	return value
}

type Env struct {
	Envs []map[int64]int64
}

func NewEnv() *Env {
	envs := []map[int64]int64{}
	envs = append(envs, make(map[int64]int64))
	return &Env{Envs: envs}
}

func NewFrame() map[int64]int64 {
	return make(map[int64]int64)
}

type Vm struct {
	Ip int64
	Ra int64
	Is []int64
	C  *Stack
	E  *Env
}

func NewVm(instrs []int64) *Vm {
	return &Vm{
		Ip: 0,
		Ra: 0,
		Is: instrs,
		C:  NewStack(),
		E:  NewEnv(),
	}
}

func (vm *Vm) Next() int64 {
	instr := vm.Is[vm.Ip]
	vm.Ip += 1
	return instr
}

func (vm *Vm) Run() {
	for vm.Ip < int64(len(vm.Is)) {
		vm.RunInstr()
	}
}

func (vm *Vm) GetEnv(input int64) int64 {
	for i := len(vm.E.Envs) - 1; i >= 0; i-- {
		if val, ok := vm.E.Envs[i][input]; ok {
			return val
		}
	}
	panic(fmt.Sprintf("%#v", vm.E.Envs))
	panic(fmt.Sprintf("No Variable %v", input))
}

func (vm *Vm) CurrentEnv() map[int64]int64 {
	return vm.E.Envs[len(vm.E.Envs)-1]
}

func (vm *Vm) RunInstr() {
	instr := vm.Next()
	switch instr {
	case compile.PUSH:
		value := vm.Next()
		// fmt.Printf("push %v\n", value)
		vm.C.Push(value)
	case compile.EQ:
		// fmt.Printf("eq\n")
		left := vm.C.Pop()
		right := vm.C.Pop()
		cond := left == right
		var condInt int64
		if cond {
			condInt = 1
		} else {
			condInt = 0
		}
		vm.C.Push(condInt)
	case compile.LT:
		// fmt.Printf("lt\n")
		left := vm.C.Pop()
		right := vm.C.Pop()
		cond := left < right
		var condInt int64
		if cond {
			condInt = 1
		} else {
			condInt = 0
		}
		vm.C.Push(condInt)
	case compile.ADD:
		// fmt.Printf("add\n")
		left := vm.C.Pop()
		right := vm.C.Pop()
		vm.C.Push(left + right)
	case compile.SUB:
		// fmt.Printf("sub\n")
		left := vm.C.Pop()
		right := vm.C.Pop()
		vm.C.Push(left - right)
	case compile.MUL:
		// fmt.Printf("mul\n")
		left := vm.C.Pop()
		right := vm.C.Pop()
		vm.C.Push(left * right)
	case compile.DIV:
		// fmt.Printf("mul\n")
		left := vm.C.Pop()
		right := vm.C.Pop()
		vm.C.Push(left / right)
	case compile.FUNDEF:
		// fmt.Printf("fundef\n")
		vm.C.Push(vm.Ip)
		for vm.Is[vm.Ip] != compile.FUNEND {
			vm.Next()
		}
	case compile.STORE:
		location := vm.Next()
		// fmt.Printf("store %v\n", location)
		value := vm.C.Pop()
		vm.CurrentEnv()[location] = value
	case compile.LOAD:
		location := vm.Next()
		// fmt.Printf("load %v\n", location)
		vm.C.Push(vm.GetEnv(location))
	case compile.RET:
		// vm.Ip = vm.Ra
		// fmt.Printf("ret\n")
		vm.E.Envs = vm.E.Envs[:len(vm.E.Envs)-1]
		vm.Ip = vm.GetEnv(LARGEST)
	case compile.CALL:
		// fmt.Printf("call\n")
		vm.Ra = vm.Ip
		vm.CurrentEnv()[LARGEST] = vm.Ra
		location := vm.C.Pop()
		vm.Ip = location
		vm.E.Envs = append(vm.E.Envs, NewFrame())
	case compile.JMPR:
		location := vm.Next()
		// fmt.Printf("jmpr %v\n", location)
		vm.Ip += location
	case compile.JMPRIF:
		location := vm.Next()
		// fmt.Printf("jmprif %v\n", location)
		cond := vm.C.Pop()
		if cond > 0 {
			vm.Ip += location
		}
	}
}
