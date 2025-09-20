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

func (f *BigFloat) InstanceVariables() *InstanceVariables {
	return nil
}

// Add sets z to the rounded sum x+y and returns z.
func (z *BigFloat) AddBigFloat(x, y *BigFloat) *BigFloat {
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
func (z *BigFloat) SubBigFloat(x, y *BigFloat) *BigFloat {
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
func (z *BigFloat) MulBigFloat(x, y *BigFloat) *BigFloat {
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
func (z *BigFloat) DivBigFloat(x, y *BigFloat) *BigFloat {
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
//	Mod(±Inf, y) = NaN
//	Mod(NaN, y) = NaN
//	Mod(x, 0) = NaN
//	Mod(x, ±Inf) = x
//	Mod(x, NaN) = NaN
func (z *BigFloat) Mod(x, y *BigFloat) *BigFloat {
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
			result := (&BigFloat{}).AddBigFloat(f, o)
			return Ref(result), Undefined
		case *BigInt:
			otherBigFloat := (&BigFloat{}).SetBigInt(o)
			result := otherBigFloat.AddBigFloat(f, otherBigFloat)
			return Ref(result), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		otherBigFloat := NewBigFloat(float64(other.AsFloat()))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return Ref(otherBigFloat.AddBigFloat(f, otherBigFloat)), Undefined
	case SMALL_INT_FLAG:
		otherBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		result := otherBigFloat.AddBigFloat(f, otherBigFloat)
		return Ref(result), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// SubtractVal another value and return an error
// if something went wrong.
func (f *BigFloat) SubtractVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			result := (&BigFloat{}).SubBigFloat(f, o)
			return Ref(result), Undefined
		case *BigInt:
			otherBigFloat := (&BigFloat{}).SetBigInt(o)
			result := otherBigFloat.SubBigFloat(f, otherBigFloat)
			return Ref(result), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		otherBigFloat := NewBigFloat(float64(other.AsFloat()))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return Ref(otherBigFloat.SubBigFloat(f, otherBigFloat)), Undefined
	case SMALL_INT_FLAG:
		otherBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		result := otherBigFloat.SubBigFloat(f, otherBigFloat)
		return Ref(result), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// MultiplyVal by another value and return an error
// if something went wrong.
func (f *BigFloat) MultiplyVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			result := (&BigFloat{}).MulBigFloat(f, o)
			return Ref(result), Undefined
		case *BigInt:
			otherBigFloat := (&BigFloat{}).SetBigInt(o)
			result := otherBigFloat.MulBigFloat(f, otherBigFloat)
			return Ref(result), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		otherBigFloat := NewBigFloat(float64(other.AsFloat()))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return Ref(otherBigFloat.MulBigFloat(f, otherBigFloat)), Undefined
	case SMALL_INT_FLAG:
		otherBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		result := otherBigFloat.MulBigFloat(f, otherBigFloat)
		return Ref(result), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// DivideVal by another value and return an error
// if something went wrong.
func (f *BigFloat) DivideVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			result := (&BigFloat{}).DivBigFloat(f, o)
			return Ref(result), Undefined
		case *BigInt:
			otherBigFloat := (&BigFloat{}).SetBigInt(o)
			result := otherBigFloat.DivBigFloat(f, otherBigFloat)
			return Ref(result), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		otherBigFloat := NewBigFloat(float64(other.AsFloat()))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return Ref(otherBigFloat.DivBigFloat(f, otherBigFloat)), Undefined
	case SMALL_INT_FLAG:
		otherBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		result := otherBigFloat.DivBigFloat(f, otherBigFloat)
		return Ref(result), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
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
			result := (&BigFloat{}).SetPrecision(max(o.Precision(), f.Precision())).Set(o)
			result.ExpBigFloat(f, o)
			return Ref(result), Undefined
		case *BigInt:
			prec := max(f.Precision(), uint(o.BitSize()), SmallIntBits)
			oBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(o)
			oBigFloat.ExpBigFloat(f, oBigFloat)
			return Ref(oBigFloat), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		prec := max(f.Precision(), FloatPrecision)
		oBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(other.AsFloat())
		oBigFloat.ExpBigFloat(f, oBigFloat)
		return Ref(oBigFloat), Undefined
	case SMALL_INT_FLAG:
		prec := max(f.Precision(), SmallIntBits)
		oBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(other.AsSmallInt())
		oBigFloat.ExpBigFloat(f, oBigFloat)
		return Ref(oBigFloat), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Perform modulo by another numeric value and return an error
// if something went wrong.
func (f *BigFloat) ModuloVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			result := (&BigFloat{}).SetPrecision(max(o.Precision(), f.Precision()))
			return Ref(result.Mod(f, o)), Undefined
		case *BigInt:
			prec := max(f.Precision(), uint(o.BitSize()), SmallIntBits)
			otherBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(o)
			return Ref(otherBigFloat.Mod(f, otherBigFloat)), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		prec := max(f.Precision(), FloatPrecision)
		otherBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(other.AsFloat())
		return Ref(otherBigFloat.Mod(f, otherBigFloat)), Undefined
	case SMALL_INT_FLAG:
		prec := max(f.Precision(), SmallIntBits)
		otherBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(other.AsSmallInt())
		return Ref(otherBigFloat.Mod(f, otherBigFloat)), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Returns 1 if i is greater than other
// Returns 0 if both are equal.
// Returns -1 if i is less than other.
// Returns nil if the comparison was impossible (NaN)
func (f *BigFloat) CompareVal(other Value) (Value, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return Nil, Undefined
			}
			return SmallInt(f.Cmp(o)).ToValue(), Undefined
		case *BigInt:
			if f.IsNaN() {
				return Nil, Undefined
			}
			oBigFloat := (&BigFloat{}).SetBigInt(o)
			return SmallInt(f.Cmp(oBigFloat)).ToValue(), Undefined
		default:
			return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return Nil, Undefined
		}
		oBigFloat := (&BigFloat{}).SetFloat(other.AsFloat())
		return SmallInt(f.Cmp(oBigFloat)).ToValue(), Undefined
	case SMALL_INT_FLAG:
		if f.IsNaN() {
			return Nil, Undefined
		}
		oBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		return SmallInt(f.Cmp(oBigFloat)).ToValue(), Undefined
	default:
		return Undefined, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f *BigFloat) GreaterThanVal(other Value) (Value, Value) {
	result, err := f.GreaterThan(other)
	return ToElkBool(result), err
}

// Check whether f is greater than other and return an error
// if something went wrong.
func (f *BigFloat) GreaterThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false, Undefined
			}
			return f.Cmp(o) > 0, Undefined
		case *BigInt:
			if f.IsNaN() {
				return false, Undefined
			}
			otherBigFloat := (&BigFloat{}).SetBigInt(o)
			return f.Cmp(otherBigFloat) > 0, Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false, Undefined
		}
		otherBigFloat := NewBigFloat(float64(o))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return f.Cmp(otherBigFloat) > 0, Undefined
	case SMALL_INT_FLAG:
		if f.IsNaN() {
			return false, Undefined
		}
		otherBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		return f.Cmp(otherBigFloat) > 0, Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f *BigFloat) GreaterThanEqualVal(other Value) (Value, Value) {
	result, err := f.GreaterThanEqual(other)
	return ToElkBool(result), err
}

// Check whether f is greater than or equal to other and return an error
// if something went wrong.
func (f *BigFloat) GreaterThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false, Undefined
			}
			return f.Cmp(o) >= 0, Undefined
		case *BigInt:
			if f.IsNaN() {
				return false, Undefined
			}
			otherBigFloat := (&BigFloat{}).SetBigInt(o)
			return f.Cmp(otherBigFloat) >= 0, Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false, Undefined
		}
		otherBigFloat := NewBigFloat(float64(o))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return f.Cmp(otherBigFloat) >= 0, Undefined
	case SMALL_INT_FLAG:
		if f.IsNaN() {
			return false, Undefined
		}
		otherBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		return f.Cmp(otherBigFloat) >= 0, Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is less than other and return an error
// if something went wrong.
func (f *BigFloat) LessThanVal(other Value) (Value, Value) {
	result, err := f.LessThan(other)
	return ToElkBool(result), err
}

// Check whether f is less than other and return an error
// if something went wrong.
func (f *BigFloat) LessThan(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false, Undefined
			}
			return f.Cmp(o) == -1, Undefined
		case *BigInt:
			if f.IsNaN() {
				return false, Undefined
			}
			otherBigFloat := (&BigFloat{}).SetBigInt(o)
			return f.Cmp(otherBigFloat) == -1, Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false, Undefined
		}
		otherBigFloat := NewBigFloat(float64(o))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return f.Cmp(otherBigFloat) == -1, Undefined
	case SMALL_INT_FLAG:
		if f.IsNaN() {
			return false, Undefined
		}
		otherBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		return f.Cmp(otherBigFloat) == -1, Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f *BigFloat) LessThanEqual(other Value) (bool, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false, Undefined
			}
			return f.Cmp(o) <= 0, Undefined
		case *BigInt:
			if f.IsNaN() {
				return false, Undefined
			}
			otherBigFloat := (&BigFloat{}).SetBigInt(o)
			return f.Cmp(otherBigFloat) <= 0, Undefined
		default:
			return false, Ref(NewCoerceError(f.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false, Undefined
		}
		otherBigFloat := NewBigFloat(float64(o))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return f.Cmp(otherBigFloat) <= 0, Undefined
	case SMALL_INT_FLAG:
		if f.IsNaN() {
			return false, Undefined
		}
		otherBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		return f.Cmp(otherBigFloat) <= 0, Undefined
	default:
		return false, Ref(NewCoerceError(f.Class(), other.Class()))
	}
}

// Check whether f is less than or equal to other and return an error
// if something went wrong.
func (f *BigFloat) LessThanEqualVal(other Value) (Value, Value) {
	result, err := f.LessThanEqual(other)
	return ToElkBool(result), err
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) LaxEqualVal(other Value) Value {
	return ToElkBool(f.LaxEqualBool(other))
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) LaxEqualBool(other Value) bool {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			if f.IsNaN() {
				return false
			}

			oBigFloat := (&BigFloat{}).SetBigInt(o)
			return f.Cmp(oBigFloat) == 0
		case *BigFloat:
			if f.IsNaN() || o.IsNaN() {
				return false
			}

			return f.Cmp(o) == 0
		case Int64:
			if f.IsNaN() {
				return false
			}

			oBigFloat := (&BigFloat{}).SetElkInt64(o)
			return f.Cmp(oBigFloat) == 0
		case UInt64:
			if f.IsNaN() {
				return false
			}

			oBigFloat := (&BigFloat{}).SetUInt64(o)
			return f.Cmp(oBigFloat) == 0
		case Float64:
			if f.IsNaN() || o.IsNaN() {
				return false
			}

			oBigFloat := (&BigFloat{}).SetElkFloat64(o)
			return f.Cmp(oBigFloat) == 0
		default:
			return false
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetSmallInt(other.AsSmallInt())
		return f.Cmp(oBigFloat) == 0
	case FLOAT_FLAG:
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetFloat(o)
		return f.Cmp(oBigFloat) == 0
	case INT64_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetElkInt64(other.AsInlineInt64())
		return f.Cmp(oBigFloat) == 0
	case INT32_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetElkInt64(Int64(other.AsInt32()))
		return f.Cmp(oBigFloat) == 0
	case INT16_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetElkInt64(Int64(other.AsInt16()))
		return f.Cmp(oBigFloat) == 0
	case INT8_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetElkInt64(Int64(other.AsInt8()))
		return f.Cmp(oBigFloat) == 0
	case UINT64_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetUInt64(other.AsInlineUInt64())
		return f.Cmp(oBigFloat) == 0
	case UINT32_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetUInt64(UInt64(other.AsUInt32()))
		return f.Cmp(oBigFloat) == 0
	case UINT16_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetUInt64(UInt64(other.AsUInt16()))
		return f.Cmp(oBigFloat) == 0
	case UINT8_FLAG:
		if f.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetUInt64(UInt64(other.AsUInt8()))
		return f.Cmp(oBigFloat) == 0
	case FLOAT64_FLAG:
		o := other.AsInlineFloat64()
		if f.IsNaN() || o.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetElkFloat64(o)
		return f.Cmp(oBigFloat) == 0
	case FLOAT32_FLAG:
		o := other.AsFloat32()
		if f.IsNaN() || o.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetElkFloat32(o)
		return f.Cmp(oBigFloat) == 0
	default:
		return false
	}
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) EqualVal(other Value) Value {
	return ToElkBool(f.Equal(other))
}

// Check whether f is equal to other and return an error
// if something went wrong.
func (f *BigFloat) Equal(other Value) bool {
	if other.IsFloat() {
		o := other.AsFloat()
		if f.IsNaN() || o.IsNaN() {
			return false
		}

		oBigFloat := (&BigFloat{}).SetFloat(o)
		return f.Cmp(oBigFloat) == 0
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
	BigFloatClass.AddConstantString("NAN", Ref(BigFloatNaNVal))
	BigFloatClass.AddConstantString("INF", Ref(BigFloatInfVal))
	BigFloatClass.AddConstantString("NEG_INF", Ref(BigFloatNegInfVal))
}
