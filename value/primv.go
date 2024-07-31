package value

type PrimV struct {
	Name string
}

func (b PrimV) Equals(v Value) bool {
	switch v.(type) {
	default:
		return false
	}
}

func (b PrimV) value() {}
