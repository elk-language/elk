package value

import (
	"fmt"
	"math/big"
	"strings"
)

type Size int64

var _ ValueInterface = Size(1)

const (
	Byte     Size = 1
	Kilobyte      = 1024 * Byte
	Megabyte      = 1024 * Kilobyte
	Gigabyte      = 1024 * Megabyte
	Terabyte      = 1024 * Gigabyte
	Petabyte      = 1024 * Terabyte
)

var SizeClass *Class

func (s Size) InBytes() Float {
	return Float(s)
}

func (s Size) Bytes() Value {
	return ToElkInt(int64(s))
}

func (s Size) InKilobytes() Float {
	return Float(s) / Float(Kilobyte)
}

func (s Size) Kilobytes() Value {
	return ToElkInt(int64(s / Kilobyte))
}

func (s Size) InMegabytes() Float {
	return Float(s) / Float(Megabyte)
}

func (s Size) Megabytes() Value {
	return ToElkInt(int64(s / Megabyte))
}

func (s Size) InGigabytes() Float {
	return Float(s) / Float(Gigabyte)
}

func (s Size) Gigabytes() Value {
	return ToElkInt(int64(s / Gigabyte))
}

func (s Size) InTerabytes() Float {
	return Float(s) / Float(Terabyte)
}

func (s Size) Terabytes() Value {
	return ToElkInt(int64(s / Terabyte))
}

func (s Size) InPetabytes() Float {
	return Float(s) / Float(Petabyte)
}

func (s Size) Petabytes() Value {
	return ToElkInt(int64(s / Petabyte))
}

func initSize() {
	SizeClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	RegisterNativeClass("Std::Size", "value.SizeClass")
}

func (Size) Class() *Class {
	return SizeClass
}

func (Size) DirectClass() *Class {
	return SizeClass
}

func (Size) SingletonClass() *Class {
	return nil
}

func (Size) InstanceVariables() *InstanceVariables {
	return nil
}

func (s Size) Inspect() string {
	return fmt.Sprintf("%s", s.ToString())
}

func (s Size) Error() string {
	return s.Inspect()
}

func (s Size) Copy() Reference {
	return s
}

func (s Size) ToString() String {
	return String(s.String())
}

func (s Size) String() string {
	var buff strings.Builder

	if s == 0 {
		return "0.bytes"
	}

	petabytes := s / Petabyte
	if petabytes != 0 {
		fmt.Fprintf(&buff, "%d.petabytes", petabytes)
	}

	terabytes := (s % Petabyte) / Terabyte
	if terabytes != 0 {
		if buff.Len() != 0 {
			buff.WriteString(" + ")
		}
		fmt.Fprintf(&buff, "%d.terabytes", terabytes)
	}

	gigabytes := (s % Terabyte) / Gigabyte
	if gigabytes != 0 {
		if buff.Len() != 0 {
			buff.WriteString(" + ")
		}
		fmt.Fprintf(&buff, "%d.gigabytes", gigabytes)
	}

	megabytes := (s % Gigabyte) / Megabyte
	if megabytes != 0 {
		if buff.Len() != 0 {
			buff.WriteString(" + ")
		}
		fmt.Fprintf(&buff, "%d.megabytes", megabytes)
	}

	kilobytes := (s % Megabyte) / Kilobyte
	if kilobytes != 0 {
		if buff.Len() != 0 {
			buff.WriteString(" + ")
		}
		fmt.Fprintf(&buff, "%d.kilobytes", kilobytes)
	}

	bytes := (s % Kilobyte) / Byte
	if bytes != 0 {
		if buff.Len() != 0 {
			buff.WriteString(" + ")
		}
		fmt.Fprintf(&buff, "%d.bytes", bytes)
	}

	return buff.String()
}

func (s Size) Equal(other Value) bool {
	o, ok := other.AsSizeOk()
	if !ok {
		return false
	}

	return s == o
}

func (s Size) Cmp(other Size) int {
	if s > other {
		return 1
	}
	if s < other {
		return -1
	}
	return 0
}

func (s Size) Multiply(other Value) (Size, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return s.MultiplyBigInt(o)
		case *BigFloat:
			return s.MultiplyBigFloat(o)
		default:
			return 0, Ref(NewCoerceError(s.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return s.MultiplySmallInt(other.AsSmallInt())
	case FLOAT_FLAG:
		return s.MultiplyFloat(other.AsFloat())
	default:
		return 0, Ref(NewCoerceError(s.Class(), other.Class()))
	}
}

func (s Size) MultiplyBigInt(other *BigInt) (Size, Value) {
	newSize := big.NewInt(int64(s))
	result := ToElkBigInt(newSize.Mul(newSize, other.ToGoBigInt()))
	return Size(result.ToSmallInt()), Undefined
}

func (s Size) MultiplySmallInt(other SmallInt) (Size, Value) {
	return s * Size(other), Undefined
}

func (s Size) MultiplyInt(other Value) (Size, Value) {
	if other.IsSmallInt() {
		return s.MultiplySmallInt(other.AsSmallInt())
	}
	return s.MultiplyBigInt((*BigInt)(other.Pointer()))
}

func (s Size) MultiplyFloat(other Float) (Size, Value) {
	return Size(Float(s) * other), Undefined
}

func (s Size) MultiplyBigFloat(other *BigFloat) (Size, Value) {
	prec := max(other.Precision(), 64)
	iBigFloat := (&BigFloat{}).SetPrecision(prec).SetInt64(int64(s))
	iBigFloat.MulMutBigFloat(iBigFloat, other)
	return Size(iBigFloat.ToInt64()), Undefined
}

func (s Size) Divide(other Value) (Size, Value) {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case *BigInt:
			return s.DivideBigInt(o)
		case *BigFloat:
			return s.DivideBigFloat(o)
		default:
			return 0, Ref(NewCoerceError(s.Class(), other.Class()))
		}
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		return s.DivideSmallInt(other.AsSmallInt())
	case FLOAT_FLAG:
		return s.DivideFloat(other.AsFloat())
	default:
		return 0, Ref(NewCoerceError(s.Class(), other.Class()))
	}
}

func (s Size) DivideBigInt(other *BigInt) (Size, Value) {
	if other.IsZero() {
		return 0, Ref(NewZeroDivisionError())
	}
	newBig := big.NewInt(int64(s))
	result := ToElkBigInt(newBig.Div(newBig, other.ToGoBigInt()))
	return Size(result.ToSmallInt()), Undefined
}

func (s Size) DivideSmallInt(other SmallInt) (Size, Value) {
	if other == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return s / Size(other), Undefined
}

func (s Size) DivideInt(other Value) (Size, Value) {
	if other.IsReference() {
		return s.DivideBigInt((*BigInt)(other.Pointer()))
	}

	return s.DivideSmallInt(other.AsSmallInt())
}

func (s Size) DivideFloat(other Float) (Size, Value) {
	if other == 0 {
		return 0, Ref(NewZeroDivisionError())
	}
	return Size(Float(s) / other), Undefined
}

func (s Size) DivideBigFloat(other *BigFloat) (Size, Value) {
	if other.IsZero() {
		return 0, Ref(NewZeroDivisionError())
	}
	prec := max(other.Precision(), 64)
	iBigFloat := (&BigFloat{}).SetPrecision(prec).SetInt64(int64(s))
	iBigFloat.DivMutBigFloat(iBigFloat, other)
	return Size(iBigFloat.ToInt64()), Undefined
}
