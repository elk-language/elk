package object

var Float32Class *Class // ::Std::Float64

// Elk's Float32 value
type Float32 float32

func (Float32) Class() *Class {
	return Float32Class
}

func (Float32) IsFrozen() bool {
	return true
}

func (Float32) SetFrozen() {}

func initFloat32() {
	Float32Class = NewClass(ClassWithParent(NumericClass))
	StdModule.AddConstant("Float32", Float32Class)
}
