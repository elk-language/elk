package value

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ALTree/bigfloat"
)

var BigFloatClass *Class // ::Std::BigFloat

// Elk's BigFloat value
type BigFloat big.Float

// Convert Go's big.Float values to Elk's BigFloat values.
func ToElkBigFloat(f *big.Float) *BigFloat {
	return (*BigFloat)(f)
}

// Create a new BigFloat with the specified value.
func NewBigFloat(f float64) *BigFloat {
	return ToElkBigFloat(big.NewFloat(f))
}

func (f *BigFloat) SetSmallInt(i SmallInt) *BigFloat {
	f.ToGoBigFloat().SetInt64(int64(i))
	return f
}

func (f *BigFloat) SetBigInt(i *BigInt) *BigFloat {
	f.ToGoBigFloat().SetInt(i.ToGoBigInt())
	return f
}

func (f *BigFloat) SetFloat(val Float) *BigFloat {
	f.ToGoBigFloat().SetFloat64(float64(val))
	return f
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
func (f *BigFloat) ToGoBigFloat() *big.Float {
	return (*big.Float)(f)
}

// Sets the f's precision to prec and possibly
// rounds the value.
func (f *BigFloat) SetPrecision(prec uint) *BigFloat {
	return ToElkBigFloat(f.ToGoBigFloat().SetPrec(prec))
}

// Gets the f's precision.
func (f *BigFloat) Precision() uint {
	return f.ToGoBigFloat().Prec()
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
		(&big.Float{}).Neg(f.ToGoBigFloat()),
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
	return fmt.Sprintf("%sbf", f.ToGoBigFloat().Text('g', -1))
}

func (f *BigFloat) InstanceVariables() SimpleSymbolMap {
	return nil
}

// Perform z = a % b by another BigFloat.
func (z *BigFloat) Mod(a, b *BigFloat) *BigFloat {
	return ToElkBigFloat(modBigFloat(
		z.ToGoBigFloat(),
		a.ToGoBigFloat(),
		b.ToGoBigFloat(),
	))
}

// Perform z = a % b.
func modBigFloat(z, a, b *big.Float) *big.Float {
	temp := &big.Float{}
	aAbs := (&big.Float{}).Abs(a)
	bAbs := (&big.Float{}).Abs(b)
	neg := a.Sign() < 0

	temp.Quo(aAbs, bAbs)      // temp = a / b
	floorBigFloat(temp, temp) // temp = floor(temp)
	temp.Mul(temp, bAbs)      // temp *= b
	z.Sub(aAbs, temp)         // z = a - temp

	if neg {
		return z.Neg(z)
	}

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
	fGo := f.ToGoBigFloat()
	prec := max(fGo.Prec(), 53)
	result.SetPrec(prec)
	return ToElkBigFloat(floorBigFloat(result, f.ToGoBigFloat()))
}

// Add another value and return an error
// if something went wrong.
func (f *BigFloat) Add(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		result := ToElkBigFloat((&big.Float{}).Add(f.ToGoBigFloat(), o.ToGoBigFloat()))
		return result, nil
	case Float:
		otherBigFloat := big.NewFloat(float64(o))
		if otherBigFloat.Prec() < f.Precision() {
			otherBigFloat.SetPrec(f.Precision())
		}
		return ToElkBigFloat(otherBigFloat.Add(f.ToGoBigFloat(), otherBigFloat)), nil
	case SmallInt:
		otherBigFloat := (&big.Float{}).SetInt64(int64(o))
		result := ToElkBigFloat(otherBigFloat.Add(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	case *BigInt:
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result := ToElkBigFloat(otherBigFloat.Add(f.ToGoBigFloat(), otherBigFloat))
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
		result := ToElkBigFloat((&big.Float{}).Sub(f.ToGoBigFloat(), o.ToGoBigFloat()))
		return result, nil
	case Float:
		otherBigFloat := big.NewFloat(float64(o))
		if otherBigFloat.Prec() < f.Precision() {
			otherBigFloat.SetPrec(f.Precision())
		}
		return ToElkBigFloat(otherBigFloat.Sub(f.ToGoBigFloat(), otherBigFloat)), nil
	case SmallInt:
		otherBigFloat := (&big.Float{}).SetInt64(int64(o))
		result := ToElkBigFloat(otherBigFloat.Sub(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	case *BigInt:
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result := ToElkBigFloat(otherBigFloat.Sub(f.ToGoBigFloat(), otherBigFloat))
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
		result := ToElkBigFloat((&big.Float{}).Mul(f.ToGoBigFloat(), o.ToGoBigFloat()))
		return result, nil
	case Float:
		otherBigFloat := big.NewFloat(float64(o))
		if otherBigFloat.Prec() < f.Precision() {
			otherBigFloat.SetPrec(f.Precision())
		}
		return ToElkBigFloat(otherBigFloat.Mul(f.ToGoBigFloat(), otherBigFloat)), nil
	case SmallInt:
		otherBigFloat := (&big.Float{}).SetInt64(int64(o))
		result := ToElkBigFloat(otherBigFloat.Mul(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	case *BigInt:
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result := ToElkBigFloat(otherBigFloat.Mul(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	default:
		return nil, NewCoerceError(f, other)
	}
}

// Divide by another value and return an error
// if something went wrong.
func (f *BigFloat) Divide(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *BigFloat:
		result := ToElkBigFloat((&big.Float{}).Quo(f.ToGoBigFloat(), o.ToGoBigFloat()))
		return result, nil
	case Float:
		otherBigFloat := big.NewFloat(float64(o))
		if otherBigFloat.Prec() < f.Precision() {
			otherBigFloat.SetPrec(f.Precision())
		}
		return ToElkBigFloat(otherBigFloat.Quo(f.ToGoBigFloat(), otherBigFloat)), nil
	case SmallInt:
		otherBigFloat := (&big.Float{}).SetInt64(int64(o))
		result := ToElkBigFloat(otherBigFloat.Quo(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
	case *BigInt:
		otherBigFloat := (&big.Float{}).SetInt(o.ToGoBigInt())
		result := ToElkBigFloat(otherBigFloat.Quo(f.ToGoBigFloat(), otherBigFloat))
		return result, nil
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
		result := bigfloat.Pow(f.ToGoBigFloat(), oBigFloat)
		return ToElkBigFloat(result), nil
	case *BigInt:
		oGo := o.ToGoBigInt()
		prec := max(f.Precision(), uint(o.BitSize()), 64)
		oBigFloat := (&big.Float{}).SetPrec(prec).SetInt(oGo)
		return ToElkBigFloat(bigfloat.Pow(f.ToGoBigFloat(), oBigFloat)), nil
	case Float:
		prec := max(f.Precision(), 53)
		oBigFloat := (&big.Float{}).SetPrec(prec).SetFloat64(float64(o))
		result := bigfloat.Pow(f.ToGoBigFloat(), oBigFloat)
		return ToElkBigFloat(result), nil
	case *BigFloat:
		fGo := f.ToGoBigFloat()
		prec := max(o.Precision(), f.Precision())
		result := (&big.Float{}).SetPrec(prec).Set(o.ToGoBigFloat())
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
}
