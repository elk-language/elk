package object

// Elk's Float value
type Float float64

func (Float) Class() *Class {
	return FloatClass
}

func (Float) IsFrozen() bool {
	return true
}

func (Float) SetFrozen() {}

var FloatClass *Class // ::Std::Float

func initFloat() {
	FloatClass = NewClass()
	StdModule.AddConstant("Float", FloatClass)
}
