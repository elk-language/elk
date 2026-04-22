package value

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ALTree/bigfloat"
	"github.com/cespare/xxhash/v2"
)

var BigFloatClass *Class // ::Std::BigFloat

// Elk's BigFloat value
type BigFloat big.Float

const BigFloatNaNMode big.RoundingMode = 0b111

var (
	// positive infinity
	BigFloatInfVal = BigFloatInf()
	// negative infinity
	BigFloatNegInfVal = BigFloatNegInf()
	// not a number value
	BigFloatNaNVal = BigFloatNaN()
)

// Convert Go's big.Float values to Elk's BigFloat values.
func ToElkBigFloat(f *big.Float) *BigFloat {
	return (*BigFloat)(f)
}

// Create a new BigFloat with the specified value.
func NewBigFloat(f float64) *BigFloat {
	if math.IsNaN(f) {
		return BigFloatNaN()
	}
	return ToElkBigFloat(big.NewFloat(f))
}

// Create NaN
func BigFloatNaN() *BigFloat {
	return ToElkBigFloat((&big.Float{}).SetMode(BigFloatNaNMode))
}

// Create +Inf
func BigFloatInf() *BigFloat {
	return ToElkBigFloat((&big.Float{}).SetInf(false))
}

// Create -Inf
func BigFloatNegInf() *BigFloat {
	return ToElkBigFloat((&big.Float{}).SetInf(true))
}

// Sets f to the value of i.
func (f *BigFloat) SetSmallInt(i SmallInt) *BigFloat {
	return f.SetInt64(int64(i))
}

// Sets f to the value of i.
func (f *BigFloat) SetElkInt64(i Int64) *BigFloat {
	f.AsGoBigFloat().SetInt64(int64(i))
	return f
}

// Sets f to the value of i.
func (f *BigFloat) SetInt64(i int64) *BigFloat {
	f.AsGoBigFloat().SetInt64(i)
	return f
}

// Sets f to the value of i.
func (f *BigFloat) SetUInt64(i UInt64) *BigFloat {
	f.AsGoBigFloat().SetUint64(uint64(i))
	return f
}

// Sets f to the value of i.
func (f *BigFloat) SetUint64(i uint64) *BigFloat {
	f.AsGoBigFloat().SetUint64(i)
	return f
}

// Sets f to the value of i.
func (f *BigFloat) SetBigInt(i *BigInt) *BigFloat {
	f.AsGoBigFloat().SetInt(i.ToGoBigInt())
	return f
}

// Sets f to the possibly rounded value of x.
func (f *BigFloat) SetFloat(x Float) *BigFloat {
	return f.SetFloat64(float64(x))
}

// Sets f to the possibly rounded value of x.
func (f *BigFloat) SetFloat64(x float64) *BigFloat {
	if math.IsNaN(x) {
		return f.SetNaN()
	}
	f.AsGoBigFloat().SetFloat64(x)
	return f
}

// Sets f to the possibly rounded value of x.
func (f *BigFloat) SetElkFloat64(x Float64) *BigFloat {
	return f.SetFloat64(float64(x))
}

// Sets f to the possibly rounded value of x.
func (f *BigFloat) SetFloat32(x float32) *BigFloat {
	return f.SetFloat64(float64(x))
}

// Sets f to the possibly rounded value of x.
func (f *BigFloat) SetElkFloat32(x Float32) *BigFloat {
	return f.SetFloat64(float64(x))
}

func (f *BigFloat) Hash() UInt64 {
	d := xxhash.New()
	bytes, err := f.AsGoBigFloat().GobEncode()
	if err != nil {
		panic(fmt.Sprintf("could not create a hash for big float: %s", err))
	}
	d.Write(bytes)
	return UInt64(d.Sum64())
}

// Convert to a Float value.
func (f *BigFloat) ToFloat() Float {
	if f.IsNaN() {
		return FloatNaN()
	}

	f64, _ := f.AsGoBigFloat().Float64()
	return Float(f64)
}

// Convert to a Float64 value.
func (f *BigFloat) ToFloat64() Float64 {
	return Float64(f.Float64())
}

func (f *BigFloat) Float64() float64 {
	if f.IsNaN() {
		return float64(Float64NaN())
	}

	f64, _ := f.AsGoBigFloat().Float64()
	return f64
}

// Convert to a Float32 value.
func (f *BigFloat) ToFloat32() Float32 {
	return Float32(f.Float32())
}

func (f *BigFloat) Float32() float32 {
	if f.IsNaN() {
		return float32(Float32NaN())
	}

	f32, _ := f.AsGoBigFloat().Float32()
	return f32
}

// Convert to a Float value.
func (f *BigFloat) ToBigInt() *BigInt {
	i, _ := f.ToGoBigFloat().Int(&big.Int{})
	return ToElkBigInt(i)
}

// Convert to an Int value.
func (f *BigFloat) ToInt() Value {
	bigInt := f.ToBigInt()
	if bigInt.IsSmallInt() {
		return bigInt.ToSmallInt().ToValue()
	}
	return Ref(bigInt)
}

// Convert to an Int64 value.
func (f *BigFloat) ToInt64() Int64 {
	return Int64(f.ToBigInt().ToGoBigInt().Int64())
}

func (f *BigFloat) Int64() int64 {
	return int64(f.ToBigInt().ToGoBigInt().Int64())
}

// Convert to an Int32 value.
func (f *BigFloat) ToInt32() Int32 {
	return Int32(f.ToBigInt().ToGoBigInt().Int64())
}

func (f *BigFloat) Int32() int32 {
	return int32(f.ToBigInt().ToGoBigInt().Int64())
}

// Convert to an Int16 value.
func (f *BigFloat) ToInt16() Int16 {
	return Int16(f.ToBigInt().ToGoBigInt().Int64())
}

func (f *BigFloat) Int16() int16 {
	return int16(f.ToBigInt().ToGoBigInt().Int64())
}

// Convert to an Int8 value.
func (f *BigFloat) ToInt8() Int8 {
	return Int8(f.ToBigInt().ToGoBigInt().Int64())
}

func (f *BigFloat) int8() int8 {
	return int8(f.ToBigInt().ToGoBigInt().Int64())
}

// Convert to an UInt value.
func (f *BigFloat) ToUInt() UInt {
	return UInt(f.ToBigInt().ToGoBigInt().Uint64())
}

// Convert to an UInt64 value.
func (f *BigFloat) ToUInt64() UInt64 {
	return UInt64(f.ToBigInt().ToGoBigInt().Uint64())
}

func (f *BigFloat) Uint64() uint64 {
	return uint64(f.ToBigInt().ToGoBigInt().Uint64())
}

// Convert to an UInt32 value.
func (f *BigFloat) ToUInt32() UInt32 {
	return UInt32(f.ToBigInt().ToGoBigInt().Uint64())
}

func (f *BigFloat) Uint32() uint32 {
	return uint32(f.ToBigInt().ToGoBigInt().Uint64())
}

// Convert to an UInt16 value.
func (f *BigFloat) ToUInt16() UInt16 {
	return UInt16(f.ToBigInt().ToGoBigInt().Uint64())
}

func (f *BigFloat) Uint16() uint16 {
	return uint16(f.ToBigInt().ToGoBigInt().Uint64())
}

// Convert to an UInt8 value.
func (f *BigFloat) ToUInt8() UInt8 {
	return UInt8(f.ToBigInt().ToGoBigInt().Uint64())
}

func (f *BigFloat) Uint8() uint8 {
	return uint8(f.ToBigInt().ToGoBigInt().Uint64())
}

// Set z = x
func (z *BigFloat) Set(x *BigFloat) *BigFloat {
	z.AsGoBigFloat().Set(x.AsGoBigFloat())
	return z
}

// Set z = +Inf
func (z *BigFloat) SetInf() *BigFloat {
	z.AsGoBigFloat().SetInf(false)
	return z
}

// Set z = -Inf
func (z *BigFloat) SetNegInf() *BigFloat {
	z.AsGoBigFloat().SetInf(true)
	return z
}

// Set z = NaN
func (z *BigFloat) SetNaN() *BigFloat {
	z.AsGoBigFloat().Set(&big.Float{}).SetMode(BigFloatNaNMode)
	return z
}

// Sign returns:
//
// -1 if f <   0
//
//	0 if f is ±0
//
// +1 if f >   0
func (f *BigFloat) Sign() int {
	return f.AsGoBigFloat().Sign()
}

// IsNaN reports whether f is a “not-a-number” value.
func (f *BigFloat) IsNaN() bool {
	return f.AsGoBigFloat().Mode() == BigFloatNaNMode
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func (f *BigFloat) IsInf(sign int) bool {
	return f.AsGoBigFloat().IsInf() && (sign == 0 || f.Sign() == sign)
}

// Parse a big float value from the given string.
func ParseBigFloat(str string) (*BigFloat, Value) {
	f, _, err := big.ParseFloat(
		str,
		10,
		PrecisionForFloatString(str),
		big.ToNearestEven,
	)
	if err != nil {
		return nil, Ref(NewError(FormatErrorClass, err.Error()))
	}

	return ToElkBigFloat(f), Undefined
}

// Same as [ParseBigFloat] but panics on error.
func ParseBigFloatPanic(str string) *BigFloat {
	result, err := ParseBigFloat(str)
	if !err.IsUndefined() {
		panic(err)
	}

	return result
}

// A BigFloatErrNaN panic is raised when converting a NaN value from Elk's BigFloat to Go's big.Float.
// Implements the error interface.
type BigFloatErrNaN struct {
	msg string
}

func (e BigFloatErrNaN) Error() string {
	return e.msg
}

// Convert Elk's BigFloat values to Go's big.Float values.
// Panics with big.ErrNaN if f is a NaN.
func (f *BigFloat) ToGoBigFloat() *big.Float {
	if f.IsNaN() {
		panic(BigFloatErrNaN{msg: "big.Float(NaN)"})
	}
	return (*big.Float)(f)
}

// Convert Elk's BigFloat values to Go's big.Float values.
// Does a cast without any checks.
func (f *BigFloat) AsGoBigFloat() *big.Float {
	return (*big.Float)(f)
}

// Sets the f's precision to prec and possibly
// rounds the value.
func (f *BigFloat) SetPrecision(prec uint) *BigFloat {
	return ToElkBigFloat(f.AsGoBigFloat().SetPrec(prec))
}

// Gets the f's precision.
func (f *BigFloat) Precision() uint {
	return f.AsGoBigFloat().Prec()
}

// Calculates the precision required to represent
// the float in the given string.
func PrecisionForFloatString(str string) uint {
	prec := uint(math.Ceil(float64(CountFloatDigits(str)) * math.Log2(10.0)))
	if prec < FloatPrecision {
		return FloatPrecision
	}

	return prec
}

// Counts the number of decimal digits in the string.
func CountFloatDigits(str string) int {
	var count int
charLoop:
	for _, char := range str {
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			count++
		case 'e', 'p':
			break charLoop
		}
	}

	return count
}

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
//
// Panics if x or y are NaN.
func (x *BigFloat) Cmp(y *BigFloat) int {
	return x.ToGoBigFloat().Cmp(y.ToGoBigFloat())
}

// Negate the number and return the result.
func (f *BigFloat) Negate() *BigFloat {
	return ToElkBigFloat(
		(&big.Float{}).Neg(f.AsGoBigFloat()),
	)
}

func (*BigFloat) Class() *Class {
	return BigFloatClass
}

func (*BigFloat) DirectClass() *Class {
	return BigFloatClass
}

func (*BigFloat) SingletonClass() *Class {
	return nil
}

func (f *BigFloat) Copy() Reference {
	return f
}

func (f *BigFloat) ToValue() Value {
	return Ref(f)
}

func (f *BigFloat) Error() string {
	return f.Inspect()
}

func (f *BigFloat) Inspect() string {
	if f.IsNaN() {
		return fmt.Sprintf("%s::NAN", f.Class().PrintableName())
	}
	if f.IsInf(1) {
		return fmt.Sprintf("%s::INF", f.Class().PrintableName())
	}
	if f.IsInf(-1) {
		return fmt.Sprintf("%s::NEG_INF", f.Class().PrintableName())
	}
	return fmt.Sprintf("%sbf", f.AsGoBigFloat().Text('g', -1))
}

func (f *BigFloat) ToString() String {
	if f.IsNaN() {
		return "NaN"
	}
	if f.IsInf(1) {
		return "Inf"
	}
	if f.IsInf(-1) {
		return "-Inf"
	}
	return String(fmt.Sprintf("%s", f.AsGoBigFloat().Text('g', -1)))
}

func (f *BigFloat) InstanceVariables() *InstanceVariables {
	return nil
}

// Add sets z to the rounded sum x+y and returns z.
func (z *BigFloat) AddMutBigFloat(x, y *BigFloat) *BigFloat {
	zGo := z.AsGoBigFloat()
	xGo := x.AsGoBigFloat()
	yGo := y.AsGoBigFloat()

	if x.IsNaN() || y.IsNaN() {
		return z.SetNaN()
	}

	// two infinities with opposite signs
	if xGo.IsInf() && yGo.IsInf() && xGo.Sign()+yGo.Sign() == 0 {
		return z.SetNaN()
	}

	zGo.Add(xGo, yGo)
	return z
}

// Sub sets z to the rounded result x-y and returns z.
func (z *BigFloat) SubMutBigFloat(x, y *BigFloat) *BigFloat {
	zGo := z.AsGoBigFloat()
	xGo := x.AsGoBigFloat()
	yGo := y.AsGoBigFloat()

	if x.IsNaN() || y.IsNaN() {
		return z.SetNaN()
	}

	// two infinities with equal signs
	if xGo.IsInf() && yGo.IsInf() && xGo.Sign()-yGo.Sign() == 0 {
		return z.SetNaN()
	}

	zGo.Sub(xGo, yGo)
	return z
}

func (z *BigFloat) IsZero() bool {
	return z.ToGoBigFloat().Cmp(&big.Float{}) == 0
}

// Mul sets z to the rounded result x*y and returns z.
func (z *BigFloat) MulMutBigFloat(x, y *BigFloat) *BigFloat {
	zGo := z.AsGoBigFloat()
	xGo := x.AsGoBigFloat()
	yGo := y.AsGoBigFloat()

	if x.IsNaN() || y.IsNaN() {
		return z.SetNaN()
	}

	// one operand is zero and the other one is an infinity
	if x.IsZero() && y.IsInf(0) || y.IsZero() && x.IsInf(0) {
		return z.SetNaN()
	}

	zGo.Mul(xGo, yGo)
	return z
}

// Div sets z to the rounded result x/y and returns z.
func (z *BigFloat) DivMutBigFloat(x, y *BigFloat) *BigFloat {
	zGo := z.AsGoBigFloat()
	xGo := x.AsGoBigFloat()
	yGo := y.AsGoBigFloat()

	if x.IsNaN() || y.IsNaN() {
		return z.SetNaN()
	}

	// both operands are infinities or zeros
	if x.IsInf(0) && y.IsInf(0) || x.IsZero() && y.IsZero() {
		return z.SetNaN()
	}

	zGo.Quo(xGo, yGo)
	return z
}

// Perform z = a % b by another BigFloat.
//
// Special cases are:
//
//	ModMutBigFloat(±Inf, y) = NaN
//	ModMutBigFloat(NaN, y) = NaN
//	ModMutBigFloat(x, 0) = NaN
//	ModMutBigFloat(x, ±Inf) = x
//	ModMutBigFloat(x, NaN) = NaN
func (z *BigFloat) ModMutBigFloat(x, y *BigFloat) *BigFloat {
	zGo := z.AsGoBigFloat()
	xGo := x.AsGoBigFloat()
	yGo := y.AsGoBigFloat()

	// x is NaN || y is NaN || x == Inf || y == 0
	if x.IsNaN() || y.IsNaN() || x.IsInf(0) || yGo.Cmp(&big.Float{}) == 0 {
		return z.SetNaN()
	}

	// y == Inf
	if y.IsInf(0) {
		return z.Set(x)
	}

	return ToElkBigFloat(modBigFloat(
		zGo,
		xGo,
		yGo,
	))
}

// Perform z = a % b.
func modBigFloat(z, a, b *big.Float) *big.Float {
	temp := &big.Float{}

	temp.Quo(a, b)         // temp = a / b
	i, acc := temp.Int64() // i = int(temp)
	if i == math.MaxInt64 && acc == big.Below {
		// float is bigger than int64
		i := &big.Int{}
		temp.Int(i)
		temp.SetInt(i) // temp = float(i)
	} else {
		temp.SetInt64(i) // temp = float(i)
	}
	temp.Mul(temp, b) // temp *= b
	z.Sub(a, temp)    // z = a - temp

	return z
}

// Perform z = floor(x)
func floorBigFloat(z *big.Float, x *big.Float) *big.Float {
	i := &big.Int{}
	x.Int(i)
	if x.Sign() < 0 {
		i = i.Sub(i, big.NewInt(1))
	}

	return z.SetInt(i)
}

func (f *BigFloat) FloorBigFloat() *BigFloat {
	result := &big.Float{}
	fGo := f.AsGoBigFloat()
	prec := max(fGo.Prec(), FloatPrecision)
	result.SetPrec(prec)
	return ToElkBigFloat(floorBigFloat(result, f.AsGoBigFloat()))
}

// AddVal another value and return an error
// if something went wrong.
func (f *BigFloat) AddVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.AddBigFloat(o).ToValue(), Undefined
		case *BigInt:
			return f.AddBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.AddFloat(other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return f.AddSmallInt(other.AsSmallInt()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

func (f *BigFloat) AddBigFloat(other *BigFloat) *BigFloat {
	return (&BigFloat{}).AddMutBigFloat(f, other)
}

func (f *BigFloat) AddFloat(other Float) *BigFloat {
	otherBigFloat := NewBigFloat(float64(other))
	if otherBigFloat.Precision() < f.Precision() {
		otherBigFloat.SetPrecision(f.Precision())
	}
	return otherBigFloat.AddMutBigFloat(f, otherBigFloat)
}

func (f *BigFloat) AddBigInt(other *BigInt) *BigFloat {
	otherBigFloat := (&BigFloat{}).SetBigInt(other)
	result := otherBigFloat.AddMutBigFloat(f, otherBigFloat)
	return result
}

func (f *BigFloat) AddSmallInt(other SmallInt) *BigFloat {
	otherBigFloat := (&BigFloat{}).SetSmallInt(other)
	result := otherBigFloat.AddMutBigFloat(f, otherBigFloat)
	return result
}

// AddInt adds a general integer value (which may be SmallInt or BigInt) to this BigFloat.
func (f *BigFloat) AddInt(other Value) *BigFloat {
	if other.IsSmallInt() {
		return f.AddSmallInt(other.AsSmallInt())
	}
	return f.AddBigInt((*BigInt)(other.Pointer()))
}

// SubtractVal another value and return an error
// if something went wrong.
func (f *BigFloat) SubtractVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.SubtractBigFloat(o).ToValue(), Undefined
		case *BigInt:
			return f.SubtractBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.SubtractFloat(other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return f.SubtractSmallInt(other.AsSmallInt()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// SubtractInt subtracts a general integer value (which may be SmallInt or BigInt) from this BigFloat.
func (f *BigFloat) SubtractInt(other Value) *BigFloat {
	if other.IsSmallInt() {
		return f.SubtractSmallInt(other.AsSmallInt())
	}
	return f.SubtractBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) SubtractBigFloat(other *BigFloat) *BigFloat {
	return (&BigFloat{}).SubMutBigFloat(f, other)
}

func (f *BigFloat) SubtractBigInt(other *BigInt) *BigFloat {
	otherBigFloat := (&BigFloat{}).SetBigInt(other)
	result := (&BigFloat{}).SubMutBigFloat(f, otherBigFloat)
	return result
}

func (f *BigFloat) SubtractFloat(other Float) *BigFloat {
	otherBigFloat := NewBigFloat(float64(other))
	if otherBigFloat.Precision() < f.Precision() {
		otherBigFloat.SetPrecision(f.Precision())
	}
	return (&BigFloat{}).SubMutBigFloat(f, otherBigFloat)
}

func (f *BigFloat) SubtractSmallInt(other SmallInt) *BigFloat {
	otherBigFloat := (&BigFloat{}).SetSmallInt(other)
	result := (&BigFloat{}).SubMutBigFloat(f, otherBigFloat)
	return result
}

// MultiplyVal by another value and return an error
// if something went wrong.
func (f *BigFloat) MultiplyVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.MultiplyBigFloat(o).ToValue(), Undefined
		case *BigInt:
			return f.MultiplyBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.MultiplyFloat(other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return f.MultiplySmallInt(other.AsSmallInt()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// MultiplyInt multiplies a general integer value (which may be SmallInt or BigInt) with this BigFloat.
func (f *BigFloat) MultiplyInt(other Value) *BigFloat {
	if other.IsSmallInt() {
		return f.MultiplySmallInt(other.AsSmallInt())
	}
	return f.MultiplyBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) MultiplyBigFloat(other *BigFloat) *BigFloat {
	return (&BigFloat{}).MulMutBigFloat(f, other)
}

func (f *BigFloat) MultiplyBigInt(other *BigInt) *BigFloat {
	otherBigFloat := (&BigFloat{}).SetBigInt(other)
	result := (&BigFloat{}).MulMutBigFloat(f, otherBigFloat)
	return result
}

func (f *BigFloat) MultiplyFloat(other Float) *BigFloat {
	otherBigFloat := NewBigFloat(float64(other))
	if otherBigFloat.Precision() < f.Precision() {
		otherBigFloat.SetPrecision(f.Precision())
	}
	return (&BigFloat{}).MulMutBigFloat(f, otherBigFloat)
}

func (f *BigFloat) MultiplySmallInt(other SmallInt) *BigFloat {
	otherBigFloat := (&BigFloat{}).SetSmallInt(other)
	result := (&BigFloat{}).MulMutBigFloat(f, otherBigFloat)
	return result
}

// DivideVal divides by another value and returns the result or an error if something went wrong.
func (f *BigFloat) DivideVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.DivideBigFloat(o).ToValue(), Undefined
		case *BigInt:
			return f.DivideBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.DivideFloat(other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return f.DivideSmallInt(other.AsSmallInt()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// DivideInt divides a general integer value (which may be SmallInt or BigInt) by this BigFloat.
func (f *BigFloat) DivideInt(other Value) *BigFloat {
	if other.IsSmallInt() {
		return f.DivideSmallInt(other.AsSmallInt())
	}
	return f.DivideBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) DivideBigFloat(other *BigFloat) *BigFloat {
	return (&BigFloat{}).DivMutBigFloat(f, other)
}

func (f *BigFloat) DivideBigInt(other *BigInt) *BigFloat {
	otherBigFloat := (&BigFloat{}).SetBigInt(other)
	return (&BigFloat{}).DivMutBigFloat(f, otherBigFloat)
}

func (f *BigFloat) DivideFloat(other Float) *BigFloat {
	otherBigFloat := NewBigFloat(float64(other))
	if otherBigFloat.Precision() < f.Precision() {
		otherBigFloat.SetPrecision(f.Precision())
	}
	return (&BigFloat{}).DivMutBigFloat(f, otherBigFloat)
}

func (f *BigFloat) DivideSmallInt(other SmallInt) *BigFloat {
	otherBigFloat := (&BigFloat{}).SetSmallInt(other)
	return (&BigFloat{}).DivMutBigFloat(f, otherBigFloat)
}

// z = x ** y
func (z *BigFloat) ExpBigFloat(x, y *BigFloat) *BigFloat {
	zGo := z.AsGoBigFloat()
	xGo := x.AsGoBigFloat()
	yGo := y.AsGoBigFloat()

	if x.IsNaN() || y.IsNaN() {
		return z.SetNaN()
	}

	// x == -1 && y == Inf => 1
	if yGo.IsInf() && xGo.Cmp(big.NewFloat(-1)) == 0 {
		return z.SetFloat(1)
	}

	// y == Inf
	if y.IsInf(0) {
		xAbs := (&big.Float{}).Abs(xGo)
		switch xAbs.Cmp(big.NewFloat(1)) {
		case 1: // |x| > 1
			if y.IsInf(-1) {
				return z.SetFloat(0)
			} else {
				return z.SetInf()
			}
		case -1: // |x| < 1
			if y.IsInf(1) {
				return z.SetFloat(0)
			} else {
				return z.SetInf()
			}
		}
	}

	// x == Inf
	if x.IsInf(0) {
		// x == -Inf
		if x.IsInf(-1) {
			yNeg := (&big.Float{}).Neg(yGo)
			return z.ExpBigFloat(&BigFloat{}, ToElkBigFloat(yNeg))
			// y < 0
		} else if yGo.Cmp(&big.Float{}) == -1 {
			return z.Set(&BigFloat{})
		}

		// x != Inf && x < 0
	} else if xGo.Cmp(&big.Float{}) == -1 {
		if !yGo.IsInt() {
			return z.SetNaN()
		}

		xAbs := (&big.Float{}).Abs(xGo)
		result := bigfloat.Pow(xAbs, yGo)
		yInt := &big.Int{}
		yGo.Int(yInt)
		// yInt is even
		if yInt.Bit(0) == 0 {
			return z.Set(ToElkBigFloat(result))
		}

		// return -result
		return ToElkBigFloat(zGo.Neg(result))
	}

	zGo.Set(bigfloat.Pow(x.AsGoBigFloat(), y.AsGoBigFloat()))
	return z
}

// ExponentiateVal by another value and return an error
// if something went wrong.
func (f *BigFloat) ExponentiateVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.ExponentiateBigFloat(o).ToValue(), Undefined
		case *BigInt:
			return f.ExponentiateBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.ExponentiateFloat(other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return f.ExponentiateSmallInt(other.AsSmallInt()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// ExponentiateInt exponentiates the BigFloat by an integer Value (SmallInt or BigInt).
func (f *BigFloat) ExponentiateInt(other Value) *BigFloat {
	if other.IsSmallInt() {
		return f.ExponentiateSmallInt(other.AsSmallInt())
	}
	return f.ExponentiateBigInt((*BigInt)(other.Pointer()))
}

// ExponentiateBigFloat exponentiates the BigFloat by another BigFloat.
func (f *BigFloat) ExponentiateBigFloat(other *BigFloat) *BigFloat {
	result := (&BigFloat{}).SetPrecision(max(other.Precision(), f.Precision())).Set(other)
	result.ExpBigFloat(f, other)
	return result
}

// ExponentiateBigInt exponentiates the BigFloat by a BigInt.
func (f *BigFloat) ExponentiateBigInt(other *BigInt) *BigFloat {
	prec := max(f.Precision(), uint(other.BitSize()), SmallIntBits)
	oBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(other)
	oBigFloat.ExpBigFloat(f, oBigFloat)
	return oBigFloat
}

// ExponentiateFloat exponentiates the BigFloat by a float64.
func (f *BigFloat) ExponentiateFloat(other Float) *BigFloat {
	prec := max(f.Precision(), FloatPrecision)
	oBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(other)
	oBigFloat.ExpBigFloat(f, oBigFloat)
	return oBigFloat
}

// ExponentiateSmallInt exponentiates the BigFloat by a SmallInt.
func (f *BigFloat) ExponentiateSmallInt(other SmallInt) *BigFloat {
	prec := max(f.Precision(), SmallIntBits)
	oBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(other)
	oBigFloat.ExpBigFloat(f, oBigFloat)
	return oBigFloat
}

// Perform modulo by another numeric value and return an error
// if something went wrong.
func (f *BigFloat) ModuloVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.ModuloBigFloat(o).ToValue(), Undefined
		case *BigInt:
			return f.ModuloBigInt(o).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.ModuloFloat(other.AsFloat()).ToValue(), Undefined
	case SMALL_INT_FLAG:
		return f.ModuloSmallInt(other.AsSmallInt()).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// ModuloInt performs modulo with a general integer value (SmallInt or BigInt).
func (f *BigFloat) ModuloInt(other Value) *BigFloat {
	if other.IsSmallInt() {
		return f.ModuloSmallInt(other.AsSmallInt())
	}
	return f.ModuloBigInt((*BigInt)(other.Pointer()))
}

// ModuloBigFloat performs modulo with another BigFloat.
func (f *BigFloat) ModuloBigFloat(other *BigFloat) *BigFloat {
	result := (&BigFloat{}).SetPrecision(max(other.Precision(), f.Precision()))
	result.ModMutBigFloat(f, other)
	return result
}

// ModuloBigInt performs modulo with a BigInt.
func (f *BigFloat) ModuloBigInt(other *BigInt) *BigFloat {
	prec := max(f.Precision(), uint(other.BitSize()), SmallIntBits)
	otherBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(other)
	otherBigFloat.ModMutBigFloat(f, otherBigFloat)
	return otherBigFloat
}

// ModuloFloat performs modulo with a Float.
func (f *BigFloat) ModuloFloat(other Float) *BigFloat {
	prec := max(f.Precision(), FloatPrecision)
	otherBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(other)
	otherBigFloat.ModMutBigFloat(f, otherBigFloat)
	return otherBigFloat
}

// ModuloSmallInt performs modulo with a SmallInt.
func (f *BigFloat) ModuloSmallInt(other SmallInt) *BigFloat {
	prec := max(f.Precision(), SmallIntBits)
	otherBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(other)
	otherBigFloat.ModMutBigFloat(f, otherBigFloat)
	return otherBigFloat
}

// Returns 1 if f is greater than other
// Returns 0 if both are equal.
// Returns -1 if f is less than other.
// Returns nil if the comparison was impossible (NaN)
func (f *BigFloat) CompareVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.CompareBigFloat(o), Undefined
		case *BigInt:
			return f.CompareBigInt(o), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.CompareFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.CompareSmallInt(other.AsSmallInt()), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// CompareInt compares a BigFloat with a general integer Value (SmallInt or BigInt).
func (f *BigFloat) CompareInt(other Value) Value {
	if other.IsSmallInt() {
		return f.CompareSmallInt(other.AsSmallInt())
	}
	return f.CompareBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) CompareBigFloat(other *BigFloat) Value {
	if f.IsNaN() || other.IsNaN() {
		return Nil
	}
	return SmallInt(f.Cmp(other)).ToValue()
}

func (f *BigFloat) CompareBigInt(other *BigInt) Value {
	if f.IsNaN() {
		return Nil
	}
	oBigFloat := (&BigFloat{}).SetBigInt(other)
	return SmallInt(f.Cmp(oBigFloat)).ToValue()
}

func (f *BigFloat) CompareFloat(other Float) Value {
	if f.IsNaN() || other.IsNaN() {
		return Nil
	}
	oBigFloat := (&BigFloat{}).SetFloat(other)
	return SmallInt(f.Cmp(oBigFloat)).ToValue()
}

func (f *BigFloat) CompareSmallInt(other SmallInt) Value {
	if f.IsNaN() {
		return Nil
	}
	oBigFloat := (&BigFloat{}).SetSmallInt(other)
	return SmallInt(f.Cmp(oBigFloat)).ToValue()
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f *BigFloat) GreaterThanVal(other Value) (Value, Value) {
	result, err := f.GreaterThan(other)
	return Bool(result).ToValue(), err
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f *BigFloat) GreaterThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.GreaterThanBigFloat(o), Undefined
		case *BigInt:
			return f.GreaterThanBigInt(o), Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.GreaterThanFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.GreaterThanSmallInt(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// GreaterThanInt returns true if the BigFloat is greater than the integer value (SmallInt or BigInt).
func (f *BigFloat) GreaterThanInt(other Value) bool {
	if other.IsSmallInt() {
		return f.GreaterThanSmallInt(other.AsSmallInt())
	}
	return f.GreaterThanBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) GreaterThanBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	return f.Cmp(other) > 0
}

func (f *BigFloat) GreaterThanBigInt(other *BigInt) bool {
	if f.IsNaN() {
		return false
	}
	otherBigFloat := (&BigFloat{}).SetBigInt(other)
	return f.Cmp(otherBigFloat) > 0
}

func (f *BigFloat) GreaterThanFloat(other Float) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	otherBigFloat := NewBigFloat(float64(other))
	if otherBigFloat.Precision() < f.Precision() {
		otherBigFloat.SetPrecision(f.Precision())
	}
	return f.Cmp(otherBigFloat) > 0
}

func (f *BigFloat) GreaterThanSmallInt(other SmallInt) bool {
	if f.IsNaN() {
		return false
	}
	otherBigFloat := (&BigFloat{}).SetSmallInt(other)
	return f.Cmp(otherBigFloat) > 0
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f *BigFloat) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := f.GreaterThanEqual(other)
	return Bool(result).ToValue(), err
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f *BigFloat) GreaterThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.GreaterThanEqualBigFloat(o), Undefined
		case *BigInt:
			return f.GreaterThanEqualBigInt(o), Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.GreaterThanEqualFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.GreaterThanEqualSmallInt(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// GreaterThanEqualInt checks whether f is greater than or equal to a general integer value (which may be SmallInt or BigInt).
func (f *BigFloat) GreaterThanEqualInt(other Value) bool {
	if other.IsSmallInt() {
		return f.GreaterThanEqualSmallInt(other.AsSmallInt())
	}
	return f.GreaterThanEqualBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) GreaterThanEqualBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	return f.Cmp(other) >= 0
}

func (f *BigFloat) GreaterThanEqualBigInt(other *BigInt) bool {
	if f.IsNaN() {
		return false
	}
	otherBigFloat := (&BigFloat{}).SetBigInt(other)
	return f.Cmp(otherBigFloat) >= 0
}

func (f *BigFloat) GreaterThanEqualFloat(other Float) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	otherBigFloat := NewBigFloat(float64(other))
	if otherBigFloat.Precision() < f.Precision() {
		otherBigFloat.SetPrecision(f.Precision())
	}
	return f.Cmp(otherBigFloat) >= 0
}

func (f *BigFloat) GreaterThanEqualSmallInt(other SmallInt) bool {
	if f.IsNaN() {
		return false
	}
	otherBigFloat := (&BigFloat{}).SetSmallInt(other)
	return f.Cmp(otherBigFloat) >= 0
}

// Check whether f is less than other and return an error
// if something went wrong.
func (f *BigFloat) LessThanVal(other Value) (Value, Value) {
	result, err := f.LessThan(other)
	return Bool(result).ToValue(), err
}

// LessThan checks whether f is less than other and returns a bool and an error if something went wrong.
func (f *BigFloat) LessThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.LessThanBigFloat(o), Undefined
		case *BigInt:
			return f.LessThanBigInt(o), Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.LessThanFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.LessThanSmallInt(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// LessThanInt checks whether f is less than a general integer value (which may be SmallInt or BigInt).
func (f *BigFloat) LessThanInt(other Value) bool {
	if other.IsSmallInt() {
		return f.LessThanSmallInt(other.AsSmallInt())
	}
	return f.LessThanBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) LessThanBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	return f.Cmp(other) == -1
}

func (f *BigFloat) LessThanBigInt(other *BigInt) bool {
	if f.IsNaN() {
		return false
	}
	otherBigFloat := (&BigFloat{}).SetBigInt(other)
	return f.Cmp(otherBigFloat) == -1
}

func (f *BigFloat) LessThanFloat(other Float) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	otherBigFloat := NewBigFloat(float64(other))
	if otherBigFloat.Precision() < f.Precision() {
		otherBigFloat.SetPrecision(f.Precision())
	}
	return f.Cmp(otherBigFloat) == -1
}

func (f *BigFloat) LessThanSmallInt(other SmallInt) bool {
	if f.IsNaN() {
		return false
	}
	otherBigFloat := (&BigFloat{}).SetSmallInt(other)
	return f.Cmp(otherBigFloat) == -1
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f *BigFloat) LessThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			return f.LessThanEqualBigFloat(o), Undefined
		case *BigInt:
			return f.LessThanEqualBigInt(o), Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		return f.LessThanEqualFloat(other.AsFloat()), Undefined
	case SMALL_INT_FLAG:
		return f.LessThanEqualSmallInt(other.AsSmallInt()), Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// LessThanEqualInt checks whether f is less than or equal to a general integer value (which may be SmallInt or BigInt).
func (f *BigFloat) LessThanEqualInt(other Value) bool {
	if other.IsSmallInt() {
		return f.LessThanEqualSmallInt(other.AsSmallInt())
	}
	return f.LessThanEqualBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) LessThanEqualBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	return f.Cmp(other) <= 0
}

func (f *BigFloat) LessThanEqualBigInt(other *BigInt) bool {
	if f.IsNaN() {
		return false
	}
	otherBigFloat := (&BigFloat{}).SetBigInt(other)
	return f.Cmp(otherBigFloat) <= 0
}

func (f *BigFloat) LessThanEqualFloat(other Float) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	otherBigFloat := NewBigFloat(float64(other))
	if otherBigFloat.Precision() < f.Precision() {
		otherBigFloat.SetPrecision(f.Precision())
	}
	return f.Cmp(otherBigFloat) <= 0
}

func (f *BigFloat) LessThanEqualSmallInt(other SmallInt) bool {
	if f.IsNaN() {
		return false
	}
	otherBigFloat := (&BigFloat{}).SetSmallInt(other)
	return f.Cmp(otherBigFloat) <= 0
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f *BigFloat) LessThanEqualVal(other Value) (Value, Value) {
	result, err := f.LessThanEqual(other)
	return Bool(result).ToValue(), err
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) LaxEqualVal(other Value) Value {
	return Bool(f.LaxEqual(other)).ToValue()
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) LaxEqual(other Value) bool {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return f.LaxEqualBigInt(o)
		case *BigFloat:
			return f.LaxEqualBigFloat(o)
		case Int64:
			return f.LaxEqualInt64(o)
		case UInt64:
			return f.LaxEqualUInt64(o)
		case Float64:
			return f.LaxEqualFloat64(o)
		default:
			return false
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return f.LaxEqualSmallInt(other.AsSmallInt())
	case FLOAT_FLAG:
		return f.LaxEqualFloat(other.AsFloat())
	case INT64_FLAG:
		return f.LaxEqualInt64(other.AsInlineInt64())
	case INT32_FLAG:
		return f.LaxEqualInt32(other.AsInt32())
	case INT16_FLAG:
		return f.LaxEqualInt16(other.AsInt16())
	case INT8_FLAG:
		return f.LaxEqualInt8(other.AsInt8())
	case UINT64_FLAG:
		return f.LaxEqualUInt64(other.AsInlineUInt64())
	case UINT32_FLAG:
		return f.LaxEqualUInt32(other.AsUInt32())
	case UINT16_FLAG:
		return f.LaxEqualUInt16(other.AsUInt16())
	case UINT8_FLAG:
		return f.LaxEqualUInt8(other.AsUInt8())
	case FLOAT64_FLAG:
		return f.LaxEqualFloat64(other.AsInlineFloat64())
	case FLOAT32_FLAG:
		return f.LaxEqualFloat32(other.AsFloat32())
	default:
		return false
	}
}

// LaxEqualInt checks if the BigFloat is equal to a general integer Value (which may be SmallInt or BigInt).
func (f *BigFloat) LaxEqualInt(other Value) bool {
	if other.IsSmallInt() {
		return f.LaxEqualSmallInt(other.AsSmallInt())
	}
	return f.LaxEqualBigInt((*BigInt)(other.Pointer()))
}

func (f *BigFloat) LaxEqualBigInt(other *BigInt) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetBigInt(other)
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}
	return f.Cmp(other) == 0
}

func (f *BigFloat) LaxEqualInt64(o Int64) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetElkInt64(o)
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualUInt64(o UInt64) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetUInt64(o)
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualFloat64(o Float64) bool {
	if f.IsNaN() || o.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetElkFloat64(o)
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualSmallInt(o SmallInt) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetSmallInt(o)
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualFloat(o Float) bool {
	if f.IsNaN() || o.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetFloat(o)
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualInt32(o Int32) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetElkInt64(Int64(o))
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualInt16(o Int16) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetElkInt64(Int64(o))
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualInt8(o Int8) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetElkInt64(Int64(o))
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualUInt32(o UInt32) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetUInt64(UInt64(o))
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualUInt16(o UInt16) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetUInt64(UInt64(o))
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualUInt8(o UInt8) bool {
	if f.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetUInt64(UInt64(o))
	return f.Cmp(oBigFloat) == 0
}

func (f *BigFloat) LaxEqualFloat32(o Float32) bool {
	if f.IsNaN() || o.IsNaN() {
		return false
	}
	oBigFloat := (&BigFloat{}).SetElkFloat32(o)
	return f.Cmp(oBigFloat) == 0
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) EqualVal(other Value) Value {
	return Bool(f.Equal(other)).ToValue()
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) EqualBigFloat(other *BigFloat) bool {
	if f.IsNaN() || other.IsNaN() {
		return false
	}

	return f.Cmp(other) == 0
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) Equal(other Value) bool {
	if o, ok := other.SafeAsReference().(*BigFloat); ok {
		return f.EqualBigFloat(o)
	}

	return false
}

// Check whether f is strictly equal to other and return an error
// if something went wrong.
func (f *BigFloat) StrictEqualVal(other Value) Value {
	return f.EqualVal(other)
}

func initBigFloat() {
	BigFloatClass = NewClass()

	StdModule.AddConstantString("BigFloat", Ref(BigFloatClass))
	RegisterNativeClass("Std::BigFloat", "value.BigFloatClass")

	BigFloatClass.AddConstantString("NAN", Ref(BigFloatNaNVal))
	RegisterNativeConstant("Std::BigFloat::NAN", "value.BigFloatNaNVal", FetchGoType("*value.BigFloat"))

	BigFloatClass.AddConstantString("INF", Ref(BigFloatInfVal))
	RegisterNativeConstant("Std::BigFloat::INF", "value.BigFloatInfVal", FetchGoType("*value.BigFloat"))

	BigFloatClass.AddConstantString("NEG_INF", Ref(BigFloatNegInfVal))
	RegisterNativeConstant("Std::BigFloat::NEG_INF", "value.BigFloatNegInfVal", FetchGoType("*value.BigFloat"))
}
