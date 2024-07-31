package compile

import (
	"reflect"
	"testing"
)

func TestAssemble(t *testing.T) {
	tests := []struct {
		input  string
		output []int64
	}{
		{"1", []int64{PUSH, 1}},
		{"1 + 1", []int64{PUSH, 1, PUSH, 1, ADD}},
		{"let x = 1", []int64{
			PUSH, 1,
			STORE, 1,
		}},
		{"fn x => x", []int64{
			FUNDEF,
			STORE, 1,
			LOAD, 1,
			RET,
			FUNEND,
		}},
	}
	for _, test := range tests {
		t.Run("order of ops", func(t *testing.T) {
			output := Assemble(test.input)

			if !reflect.DeepEqual(output, test.output) {
				t.Errorf("got %v, wanted %v", output, test.output)
			}
		})
	}
}
