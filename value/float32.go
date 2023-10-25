package value

import (
	"fmt"
	"math"

	"github.com/google/go-cmp/cmp"
)

var Float32Class *Class // ::Std::Float64

// Elk's Float32 value
type Float32 float32

// Positive infinity
func Float32Inf() Float32 {
	return Float32(math.Inf(1))
}

// Negative infinity
func Float32NegInf() Float32 {
	return Float32(math.Inf(-1))
}

// Not a number
func Float32NaN() Float32 {
	return Float32(math.NaN())
}

func (Float32) Class() *Class {
	return Float32Class
}

func (Float32) IsFrozen() bool {
	return true
}

func (Float32) SetFrozen() {}

func (f Float32) Inspect() string {
	if f.IsNaN() {
		return fmt.Sprintf("%s::NAN", f.Class().PrintableName())
	}
	if f.IsInf(1) {
		return fmt.Sprintf("%s::INF", f.Class().PrintableName())
	}
	if f.IsInf(-1) {
		return fmt.Sprintf("%s::NEG_INF", f.Class().PrintableName())
	}
	return fmt.Sprintf("%ff32", f)
}

func (f Float32) InstanceVariables() SimpleSymbolMap {
	return nil
}

// IsNaN reports whether f is a “not-a-number” value.
func (f Float32) IsNaN() bool {
	return math.IsNaN(float64(f))
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func (f Float32) IsInf(sign int) bool {
	return math.IsInf(float64(f), sign)
}

var Float32Comparer = cmp.Comparer(func(x, y Float32) bool {
	if x.IsNaN() || y.IsNaN() {
		return x.IsNaN() && y.IsNaN()
	}
	return x == y
})

func initFloat32() {
	Float32Class = NewClass(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("Float32", Float32Class)
	Float32Class.AddConstant("NAN", Float32NaN())
	Float32Class.AddConstant("INF", Float32Inf())
	Float32Class.AddConstant("NEG_INF", Float32NegInf())
}
