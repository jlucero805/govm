package vm

// func TestVmAll(t *testing.T) {
// 	tests := []struct {
// 		input  string
// 		output int64
// 	}{
// 		{"1 + 1", 2},
// 		{"10 < 100", 1},
// 		{"1 + 1 + 2", 4},
// 		{"1 + 3 * 2", 7},
// 		{"1 + 3 * 2 1 + 1", 2},
// 		{"1 + 2 / 2", 2},
// 		{"4 - 2", 2},
// 		{"fn x => x + 1", 1},
// 		{"fn x => x + 1 2 * 3", 6},
// 		{"let x = 2 x + x", 4},
// 		{"let x = fn x => x + 1 1 + 1", 2},
// 		{"let add = fn x => x + 1 add(1)", 2},
// 		{"if 1 < 2 then 69 else 420", 69},
// 		{"if 2 < 2 then 69 else 420", 420},
// 		{
// 			`let add = fn x, y => x + y
// 			let mul = fn x, y => x * y
// 			add(1, 2) + mul(3, 4)`,
// 			15,
// 		},
// 		{
// 			`let lt10 = fn x => x < 10
// 			lt10(5)`,
// 			1,
// 		},
// 		{
// 			`let lt10 = fn x => x < 10
// 			if lt10(10) then 69 else 420`,
// 			420,
// 		},
// 		{
// 			`let lt10 = fn x => x < 10
// 			if lt10(0) then 69 else 420`,
// 			69,
// 		},
// 		{
// 			`10 - 9`,
// 			1,
// 		},
// 		{
// 			`let sum = fn n =>
// 				if 0 < n then n + sum(0)
// 					else n
// 			sum(15)`,
// 			15,
// 		},
// 		{
// 			`let sum = 1==2
// 			sum`,
// 			0,
// 		},
// 		{
// 			`let sum = 1==1
// 			sum`,
// 			1,
// 		},
// 		{
// 			`let fib = fn n =>
// 				if n == 1 then 1
// 				else if n == 2 then 1
// 				else fib(n - 1) + fib(n - 2)
// 			fib(10)`,
// 			55,
// 		},
// 		{
// 			`
// 			let a = fn x => x
// 			a(1)
// 			`,
// 			1,
// 		},
// 		{
// 			`
// 			let a = fn x => x
// 			let b = fn x => a(x)
// 			b(1)
// 			`,
// 			1,
// 		},
// 		{
// 			`
// 			let a = fn x => x
// 			let b = fn x => a(x)
// 			let c = fn x => b(x)
// 			c(1)
// 			`,
// 			1,
// 		},
// 		{
// 			`
// 			let a = fn x => x
// 			let b = fn x => a(x)
// 			let c = fn x => b(x)
// 			let d = fn x => c(x)
// 			d(1)
// 			`,
// 			1,
// 		},
// 		{
// 			`
// 			let a = fn x => x
// 			let b = fn x => a(x)
// 			let c = fn x => b(x)
// 			let d = fn x => c(x)
// 			let e = fn x => d(x)
// 			e(1)
// 			`,
// 			1,
// 		},
// 		{
// 			`
// 			let a = fn x => x
// 			let b = fn x => a(x)
// 			let c = fn x => b(x)
// 			let d = fn x => c(x)
// 			let e = fn x => d(x)
// 			let f = fn x => e(x)
// 			f(1)
// 			`,
// 			1,
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run("order of ops", func(t *testing.T) {
// 			bytecode := compile.AssembleAll(test.input)
// 			vm := NewVm(bytecode)
// 			vm.Run()
// 			fmt.Println(vm)
// 			got := vm.C.Pop()
// 			if !reflect.DeepEqual(got, test.output) {
// 				t.Errorf("got %v, wanted %v", got, test.output)
// 			}
// 		})
// 	}
// }

// func TestVmSingle(t *testing.T) {
// 	tests := []struct {
// 		input  string
// 		output int64
// 	}{
// 		{"1 + 1", 2},
// 		{"1 + 1 + 2", 4},
// 		{"1 + 3 * 2", 7},
// 	}
// 	for _, test := range tests {
// 		t.Run("order of ops", func(t *testing.T) {
// 			bytecode := compile.Assemble(test.input)
// 			vm := NewVm(bytecode)
// 			vm.Run()
// 			fmt.Println(vm)
// 			got := vm.C.Pop()
// 			if !reflect.DeepEqual(got, test.output) {
// 				t.Errorf("got %v, wanted %v", got, test.output)
// 			}
// 		})
// 	}
// }
