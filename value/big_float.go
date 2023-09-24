package value

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ALTree/bigfloat"
	"github.com/google/go-cmp/cmp"
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
	f.AsGoBigFloat().SetInt64(int64(i))
	return f
}

// Sets f to the value of i.
func (f *BigFloat) SetBigInt(i *BigInt) *BigFloat {
	f.AsGoBigFloat().SetInt(i.ToGoBigInt())
	return f
}

// Sets f to the possibly rounded value of x.
func (f *BigFloat) SetFloat(x Float) *BigFloat {
	if math.IsNaN(float64(x)) {
		return f.SetNaN()
	}
	f.AsGoBigFloat().SetFloat64(float64(x))
	return f
}

// Convert to a Float value.
func (f *BigFloat) ToFloat() Float {
	if f.IsNaN() {
		return FloatNaN()
	}

	f64, _ := f.AsGoBigFloat().Float64()
	return Float(f64)
}

func (f *BigFloat) Set(val *BigFloat) *BigFloat {
	f.AsGoBigFloat().Set(val.AsGoBigFloat())
	return f
}

func (z *BigFloat) SetNaN() *BigFloat {
	z.AsGoBigFloat().Set(&big.Float{}).SetMode(BigFloatNaNMode)
	return z
}

// Sign returns:

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
func ParseBigFloat(str string) (*BigFloat, *Error) {
	f, _, err := big.ParseFloat(
		str,
		10,
		PrecisionForFloatString(str),
		big.ToNearestEven,
	)
	if err != nil {
		return nil, NewError(FormatErrorClass, err.Error())
	}

	return ToElkBigFloat(f), nil
}

// Same as [ParseBigFloat] but panics on error.
func ParseBigFloatPanic(str string) *BigFloat {
	result, err := ParseBigFloat(str)
	if err != nil {
		panic(err)
	}

	return result
}

// Convert Elk's BigFloat values to Go's big.Float values.
// Panics with big.ErrNaN if f is a NaN.
func (f *BigFloat) ToGoBigFloat() *big.Float {
	if f.IsNaN() {
		panic(big.ErrNaN{})
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
	if prec < 53 {
		return 53
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

// Negate the number and return the result.
func (f *BigFloat) Negate() *BigFloat {
	return ToElkBigFloat(
		(&big.Float{}).Neg(f.AsGoBigFloat()),
	)
}

func (*BigFloat) Class() *Class {
	return BigFloatClass
}

func (*BigFloat) IsFrozen() bool {
	return true
}

func (*BigFloat) SetFrozen() {}

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

func (f *BigFloat) InstanceVariables() SimpleSymbolMap {
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
func (z *BigFloat) Mod(a, b *BigFloat) *BigFloat {
	return ToElkBigFloat(modBigFloat(
		z.AsGoBigFloat(),
		a.AsGoBigFloat(),
		b.AsGoBigFloat(),
	))
}

var bigFloatComparer = cmp.Comparer(func(x, y *BigFloat) bool {
	if x.IsNaN() || y.IsNaN() {
		return x.IsNaN() && y.IsNaN()
	}
	return x.AsGoBigFloat().Cmp(y.AsGoBigFloat()) == 0 &&
		(x.IsInf(0) || y.IsInf(0) || x.Precision() == y.Precision())
})

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
	prec := max(fGo.Prec(), 53)
	result.SetPrec(prec)
	return ToElkBigFloat(floorBigFloat(result, f.AsGoBigFloat()))
}

// Add another value and return an error
// if something went wrong.
func (f *BigFloat) Add(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		result := (&BigFloat{}).AddBigFloat(f, o)
		return result, nil
	case Float:
		otherBigFloat := NewBigFloat(float64(o))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return otherBigFloat.AddBigFloat(f, otherBigFloat), nil
	case SmallInt:
		otherBigFloat := (&BigFloat{}).SetSmallInt(o)
		result := otherBigFloat.AddBigFloat(f, otherBigFloat)
		return result, nil
	case *BigInt:
		otherBigFloat := (&BigFloat{}).SetBigInt(o)
		result := otherBigFloat.AddBigFloat(f, otherBigFloat)
		return result, nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Subtract another value and return an error
// if something went wrong.
func (f *BigFloat) Subtract(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		return (&BigFloat{}).SubBigFloat(f, o), nil
	case Float:
		otherBigFloat := NewBigFloat(float64(o))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return otherBigFloat.SubBigFloat(f, otherBigFloat), nil
	case SmallInt:
		otherBigFloat := (&BigFloat{}).SetSmallInt(o)
		return otherBigFloat.SubBigFloat(f, otherBigFloat), nil
	case *BigInt:
		otherBigFloat := (&BigFloat{}).SetBigInt(o)
		result := otherBigFloat.SubBigFloat(f, otherBigFloat)
		return result, nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Multiply by another value and return an error
// if something went wrong.
func (f *BigFloat) Multiply(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		return (&BigFloat{}).MulBigFloat(f, o), nil
	case Float:
		otherBigFloat := NewBigFloat(float64(o))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return otherBigFloat.MulBigFloat(f, otherBigFloat), nil
	case SmallInt:
		otherBigFloat := (&BigFloat{}).SetSmallInt(o)
		return otherBigFloat.MulBigFloat(f, otherBigFloat), nil
	case *BigInt:
		otherBigFloat := (&BigFloat{}).SetBigInt(o)
		return otherBigFloat.MulBigFloat(f, otherBigFloat), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Divide by another value and return an error
// if something went wrong.
func (f *BigFloat) Divide(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		return (&BigFloat{}).DivBigFloat(f, o), nil
	case Float:
		otherBigFloat := NewBigFloat(float64(o))
		if otherBigFloat.Precision() < f.Precision() {
			otherBigFloat.SetPrecision(f.Precision())
		}
		return otherBigFloat.DivBigFloat(f, otherBigFloat), nil
	case SmallInt:
		otherBigFloat := (&BigFloat{}).SetSmallInt(o)
		return otherBigFloat.DivBigFloat(f, otherBigFloat), nil
	case *BigInt:
		otherBigFloat := (&BigFloat{}).SetBigInt(o)
		return otherBigFloat.DivBigFloat(f, otherBigFloat), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Exponentiate by another value and return an error
// if something went wrong.
func (f *BigFloat) Exponentiate(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		prec := max(f.Precision(), 64)
		oBigFloat := (&big.Float{}).SetPrec(prec).SetInt64(int64(o))
		result := bigfloat.Pow(f.AsGoBigFloat(), oBigFloat)
		return ToElkBigFloat(result), nil
	case *BigInt:
		oGo := o.ToGoBigInt()
		prec := max(f.Precision(), uint(o.BitSize()), 64)
		oBigFloat := (&big.Float{}).SetPrec(prec).SetInt(oGo)
		return ToElkBigFloat(bigfloat.Pow(f.AsGoBigFloat(), oBigFloat)), nil
	case Float:
		prec := max(f.Precision(), 53)
		oBigFloat := (&big.Float{}).SetPrec(prec).SetFloat64(float64(o))
		result := bigfloat.Pow(f.AsGoBigFloat(), oBigFloat)
		return ToElkBigFloat(result), nil
	case *BigFloat:
		fGo := f.AsGoBigFloat()
		prec := max(o.Precision(), f.Precision())
		result := (&big.Float{}).SetPrec(prec).Set(o.AsGoBigFloat())
		result = bigfloat.Pow(fGo, result)
		return ToElkBigFloat(result), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Perform modulo by another numeric value and return an error
// if something went wrong.
func (f *BigFloat) Modulo(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		prec := max(f.Precision(), 64)
		oBigFloat := (&BigFloat{}).SetPrecision(prec).SetSmallInt(o)
		return oBigFloat.Mod(f, oBigFloat), nil
	case *BigInt:
		prec := max(f.Precision(), uint(o.BitSize()), 64)
		oBigFloat := (&BigFloat{}).SetPrecision(prec).SetBigInt(o)
		return oBigFloat.Mod(f, oBigFloat), nil
	case Float:
		prec := max(f.Precision(), 53)
		oBigFloat := (&BigFloat{}).SetPrecision(prec).SetFloat(o)
		return oBigFloat.Mod(f, oBigFloat), nil
	case *BigFloat:
		prec := max(f.Precision(), o.Precision())
		result := (&BigFloat{}).SetPrecision(prec)
		return result.Mod(f, o), nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

func initBigFloat() {
	BigFloatClass = NewClass(
		ClassWithParent(NumericClass),
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstant("BigFloat", BigFloatClass)
	BigFloatClass.AddConstant("NAN", BigFloatNaNVal)
	BigFloatClass.AddConstant("INF", BigFloatInfVal)
	BigFloatClass.AddConstant("NEG_INF", BigFloatNegInfVal)
}
