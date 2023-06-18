package object

var Float64Class *Class // ::Std::Float64

// Elk's Float64 value
type Float64 float64

func (Float64) Class() *Class {
	return Float64Class
}

func (Float64) IsFrozen() bool {
	return true
}

func (Float64) SetFrozen() {}

func initFloat64() {
	Float64Class = NewClass(ClassWithParent(NumericClass))
	StdModule.AddConstant("Float64", Float64Class)
}
