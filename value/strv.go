package value

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
