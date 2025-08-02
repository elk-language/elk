package value

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/cespare/xxhash/v2"
)

var Float32Class *Class                    // ::Std::Float32
var Float32ConvertibleInterface *Interface // ::Std::Float32::Convertible

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

func (f Float32) ToValue() Value {
	return Value{
		flag: FLOAT32_FLAG,
		data: uintptr(math.Float32bits(float32(f))),
	}
}

func (Float32) Class() *Class {
	return Float32Class
}

func (Float32) DirectClass() *Class {
	return Float32Class
}

func (Float32) SingletonClass() *Class {
	return nil
}

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
	return fmt.Sprintf("%gf32", f)
}

func (f Float32) Error() string {
	return f.Inspect()
}

func (f Float32) InstanceVariables() *InstanceVariables {
	return nil
}

func (f Float32) ToString() String {
	return String(fmt.Sprintf("%g", f))
}

func (f Float32) Hash() UInt64 {
	d := xxhash.New()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(float32(f)))
	d.Write(b)
	return UInt64(d.Sum64())
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

func (f Float32) Add(other Value) (Float32, Value) {
	if !other.IsFloat32() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return f + o, Undefined
}

// ExponentiateVal by the right value.
func (f Float32) ExponentiateVal(other Value) (Float32, Value) {
	if !other.IsFloat32() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return Float32(math.Pow(float64(f), float64(o))), Undefined
}

func (f Float32) Subtract(other Value) (Float32, Value) {
	if !other.IsFloat32() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return f - o, Undefined
}

func (f Float32) Multiply(other Value) (Float32, Value) {
	if !other.IsFloat32() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return f * o, Undefined
}

func (f Float32) ModuloVal(other Value) (Float32, Value) {
	if !other.IsFloat32() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return Float32(math.Mod(float64(f), float64(o))), Undefined
}

func (f Float32) Divide(other Value) (Float32, Value) {
	if !other.IsFloat32() {
		return 0, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return f / o, Undefined
}

func (f Float32) CompareVal(other Value) (Value, Value) {
	if !other.IsFloat32() {
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
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

func (f Float32) GreaterThanVal(other Value) (Value, Value) {
	result, err := f.GreaterThan(other)
	return ToElkBool(result), err
}

func (f Float32) GreaterThan(other Value) (bool, Value) {
	if !other.IsFloat32() {
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return f > o, Undefined
}

func (f Float32) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := f.GreaterThanEqual(other)
	return ToElkBool(result), err
}

func (f Float32) GreaterThanEqual(other Value) (bool, Value) {
	if !other.IsFloat32() {
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return f >= o, Undefined
}

func (f Float32) LessThanVal(other Value) (Value, Value) {
	result, err := f.LessThan(other)
	return ToElkBool(result), err
}

func (f Float32) LessThan(other Value) (bool, Value) {
	if !other.IsFloat32() {
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return f < o, Undefined
}

func (f Float32) LessThanEqualVal(other Value) (Value, Value) {
	result, err := f.LessThanEqual(other)
	return ToElkBool(result), err
}

func (f Float32) LessThanEqual(other Value) (bool, Value) {
	if !other.IsFloat32() {
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}

	o := other.AsFloat32()
	return f <= o, Undefined
}

func (f Float32) EqualVal(other Value) Value {
	return ToElkBool(f.Equal(other))
}

func (f Float32) Equal(other Value) bool {
	if !other.IsFloat32() {
		return false
	}

	o := other.AsFloat32()
	return f == o
}

func (f Float32) StrictEqualVal(other Value) Value {
	return f.EqualVal(other)
}

func initFloat32() {
	Float32Class = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Float32", Ref(Float32Class))
	Float32Class.AddConstantString("NAN", Float32NaN().ToValue())
	Float32Class.AddConstantString("INF", Float32Inf().ToValue())
	Float32Class.AddConstantString("NEG_INF", Float32NegInf().ToValue())

	Float32ConvertibleInterface = NewInterface()
	Float32Class.AddConstantString("Convertible", Ref(Float32ConvertibleInterface))
}
