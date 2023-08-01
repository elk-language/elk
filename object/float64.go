package object

import "fmt"

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

func (f Float64) Inspect() string {
	return fmt.Sprintf("%gf64", f)
}

func (f Float64) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initFloat64() {
	Float64Class = NewClass(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("Float64", Float64Class)
}
