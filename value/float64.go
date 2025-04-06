package value

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/cespare/xxhash/v2"
)

var Float64Class *Class                    // ::Std::Float64
var Float64ConvertibleInterface *Interface // ::Std::Float64::Convertible

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

func (f Float64) Copy() Reference {
	return f
}

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

func (f Float64) Error() string {
	return f.Inspect()
}

func (f Float64) InstanceVariables() SymbolMap {
	return nil
}

func (f Float64) ToString() String {
	return String(fmt.Sprintf("%g", f))
}

func (f Float64) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(float64(f)))
	d.Write(b)
	return UInt64(d.Sum64())
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

func (f Float64) Add(other Value) (Float64, Value) {
	if !other.IsFloat64() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return f + o, Undefined
}

// ExponentiateVal by the right value.
func (f Float64) ExponentiateVal(other Value) (Float64, Value) {
	if !other.IsFloat64() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return Float64(math.Pow(float64(f), float64(o))), Undefined
}

func (f Float64) Subtract(other Value) (Float64, Value) {
	if !other.IsFloat64() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return f - o, Undefined
}

func (f Float64) Multiply(other Value) (Float64, Value) {
	if !other.IsFloat64() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return f * o, Undefined
}

func (f Float64) ModuloVal(other Value) (Float64, Value) {
	if !other.IsFloat64() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return Float64(math.Mod(float64(f), float64(o))), Undefined
}

func (f Float64) Divide(other Value) (Float64, Value) {
	if !other.IsFloat64() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return f / o, Undefined
}

func (f Float64) CompareVal(other Value) (Value, Value) {
	if !other.IsFloat64() {
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	if math.IsNaN(float64(f)) || math.IsNaN(float64(o)) {
		return Nil, Undefined
	}

	if f > o {
		return SmallInt(1).ToValue(), Undefined
	}
	if f < o {
		return SmallInt(-1).ToValue(), Undefined
	}
	return SmallInt(0).ToValue(), Undefined
}

func (f Float64) GreaterThanVal(other Value) (Value, Value) {
	result, err := f.GreaterThan(other)
	return ToElkBool(result), err
}

func (f Float64) GreaterThan(other Value) (bool, Value) {
	if !other.IsFloat64() {
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return f > o, Undefined
}

func (f Float64) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := f.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (f Float64) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsFloat64() {
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return f >= o, Undefined
}

func (f Float64) LessThanVal(other Value) (Value, Value) {
	result, err := f.LessThan(other)
	return ToElkBool(result), err
}

func (f Float64) LessThan(other Value) (bool, Value) {
	if !other.IsFloat64() {
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return f < o, Undefined
}

func (f Float64) LessThanEqualVal(other Value) (Value, Value) {
	result, err := f.LessThanEqual(other)
	return ToElkBool(result), err
}

func (f Float64) LessThanEqual(other Value) (bool, Value) {
	if !other.IsFloat64() {
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat64()
	return f <= o, Undefined
}

func (f Float64) EqualVal(other Value) Value {
	return ToElkBool(f.Equal(other))
}

func (f Float64) Equal(other Value) bool {
	if !other.IsFloat64() {
		return false
	}

	o := other.AsFloat64()
	return f == o
}

func (f Float64) StrictEqualVal(other Value) Value {
	return f.EqualVal(other)
}

func initFloat64() {
	Float64Class = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Float64", Ref(Float64Class))
	Float64Class.AddConstantString("NAN", Float64NaN().ToValue())
	Float64Class.AddConstantString("INF", Float64Inf().ToValue())
	Float64Class.AddConstantString("NEG_INF", Float64NegInf().ToValue())

	Float64ConvertibleInterface = NewInterface()
	Float64Class.AddConstantString("Convertible", Ref(Float64ConvertibleInterface))
}
