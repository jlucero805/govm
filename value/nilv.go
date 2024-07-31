package value

type NilV struct{}

func (nil NilV) value() {}

func (n NilV) Equals(v Value) bool {
	switch v.(type) {
	case NilV:
		return true
	default:
		return false
	}
}
