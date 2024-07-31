package value

import "reflect"

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
