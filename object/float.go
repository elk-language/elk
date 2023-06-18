package object

var FloatClass *Class // ::Std::Float

// Elk's Float value
type Float float64

func (Float) Class() *Class {
	return FloatClass
}

func (Float) IsFrozen() bool {
	return true
}

func (Float) SetFrozen() {}

func initFloat() {
	FloatClass = NewClass(ClassWithParent(NumericClass))
	StdModule.AddConstant("Float", FloatClass)
}
