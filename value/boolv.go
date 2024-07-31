package value

type BoolV struct {
	Value bool
}

func (b BoolV) Equals(v Value) bool {
	switch node := v.(type) {
	case BoolV:
		return node.Value == b.Value
	default:
		return false
	}
}

func (b BoolV) value() {}
