package value

type Value interface {
	value()
	Equals(Value) bool
}
