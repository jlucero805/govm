package value

type NumV struct {
	Value float64
}

func (num NumV) value() {}

func (n NumV) Equals(v Value) bool {
	switch node := v.(type) {
	case NumV:
		return n.Value == node.Value
	default:
		return false
	}
}
