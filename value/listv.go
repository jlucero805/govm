package value

import "reflect"

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
