package value

import (
	"fmt"
	"math"

	"github.com/google/go-cmp/cmp"
)

var Float64Class *Class // ::Std::Float64

// Elk's Float64 value
type Float64 float64

// Positive infinity
func Float64Inf() Float64 {
	return Float64(math.Inf(1))
}

// Negative infinity
func Float64NegInf() Float64 {
	return Float64(math.Inf(-1))
}

// Not a number
func Float64NaN() Float64 {
	return Float64(math.NaN())
}

func (Float64) Class() *Class {
	return Float64Class
}

func (Float64) DirectClass() *Class {
	return Float64Class
}

func (Float64) SingletonClass() *Class {
	return nil
}

func (Float64) IsFrozen() bool {
	return true
}

func (Float64) SetFrozen() {}

func (f Float64) Inspect() string {
	if f.IsNaN() {
		return fmt.Sprintf("%s::NAN", f.Class().PrintableName())
	}
	if f.IsInf(1) {
		return fmt.Sprintf("%s::INF", f.Class().PrintableName())
	}
	if f.IsInf(-1) {
		return fmt.Sprintf("%s::NEG_INF", f.Class().PrintableName())
	}
	return fmt.Sprintf("%gf64", f)
}

func (f Float64) InstanceVariables() SimpleSymbolMap {
	return nil
}

// IsNaN reports whether f is a “not-a-number” value.
func (f Float64) IsNaN() bool {
	return math.IsNaN(float64(f))
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func (f Float64) IsInf(sign int) bool {
	return math.IsInf(float64(f), sign)
}

var Float64Comparer = cmp.Comparer(func(x, y Float64) bool {
	if x.IsNaN() || y.IsNaN() {
		return x.IsNaN() && y.IsNaN()
	}
	return x == y
})

func initFloat64() {
	Float64Class = NewClassWithOptions(
		ClassWithParent(NumericClass),
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Float64", Float64Class)
	Float64Class.AddConstantString("NAN", Float64NaN())
	Float64Class.AddConstantString("INF", Float64Inf())
	Float64Class.AddConstantString("NEG_INF", Float64NegInf())
}
