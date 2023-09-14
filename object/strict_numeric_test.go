package object

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStrictFloatExponentiate(t *testing.T) {
	tests := map[string]struct {
		a    Float64
		b    Value
		want Float64
		err  *Error
	}{
		"exponentiate String and return an error": {
			a:   Float64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float64`"),
		},
		"exponentiate Int32 and return an error": {
			a:   Float64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Float64`"),
		},
		"exponentiate positive Float64": {
			a:    Float64(5.5),
			b:    Float64(3),
			want: Float64(166.375),
		},
		"exponentiate negative Float64": {
			a:    Float64(5.5),
			b:    Float64(-2),
			want: Float64(0.03305785123966942),
		},
		"exponentiate zero": {
			a:    Float64(5.5),
			b:    Float64(0),
			want: Float64(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictFloatExponentiate(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictIntExponentiate(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"exponentiate String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"exponentiate Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"exponentiate positive Int64 5 ** 2": {
			a:    Int64(5),
			b:    Int64(2),
			want: Int64(25),
		},
		"exponentiate positive Int64 7 ** 8": {
			a:    Int64(7),
			b:    Int64(8),
			want: Int64(5764801),
		},
		"exponentiate positive Int64 2 ** 5": {
			a:    Int64(2),
			b:    Int64(5),
			want: Int64(32),
		},
		"exponentiate positive Int64 6 ** 1": {
			a:    Int64(6),
			b:    Int64(1),
			want: Int64(6),
		},
		"exponentiate negative Int64": {
			a:    Int64(4),
			b:    Int64(-2),
			want: Int64(1),
		},
		"exponentiate zero": {
			a:    Int64(25),
			b:    Int64(0),
			want: Int64(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntExponentiate(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumericAdd(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"add String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"add Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"add positive Int64": {
			a:    Int64(53),
			b:    Int64(21),
			want: Int64(74),
		},
		"add negative Int64": {
			a:    Int64(25),
			b:    Int64(-50),
			want: Int64(-25),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictNumericAdd(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumericSubtract(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"subtract String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"subtract Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"subtract positive Int64": {
			a:    Int64(53),
			b:    Int64(21),
			want: Int64(32),
		},
		"subtract negative Int64": {
			a:    Int64(25),
			b:    Int64(-50),
			want: Int64(75),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictNumericSubtract(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumericMultiply(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"multiply String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"multiply Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"multiply positive Int64": {
			a:    Int64(53),
			b:    Int64(2),
			want: Int64(106),
		},
		"multiply negative Int64": {
			a:    Int64(25),
			b:    Int64(-2),
			want: Int64(-50),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictNumericMultiply(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictIntDivide(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"divide by String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"divide Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"divide positive Int64": {
			a:    Int64(54),
			b:    Int64(2),
			want: Int64(27),
		},
		"divide negative Int64": {
			a:    Int64(50),
			b:    Int64(-2),
			want: Int64(-25),
		},
		"divide by zero": {
			a:   Int64(50),
			b:   Int64(0),
			err: NewError(ZeroDivisionErrorClass, "can't divide an integer by zero"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntDivide(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictFloatDivide(t *testing.T) {
	tests := map[string]struct {
		a    Float64
		b    Value
		want Float64
		err  *Error
	}{
		"divide by String and return an error": {
			a:   Float64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float64`"),
		},
		"divide Int32 and return an error": {
			a:   Float64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Float64`"),
		},
		"divide positive Float64": {
			a:    Float64(54.5),
			b:    Float64(2),
			want: Float64(27.25),
		},
		"divide negative Float64": {
			a:    Float64(50),
			b:    Float64(-2),
			want: Float64(-25),
		},
		"divide by zero": {
			a:    Float64(50),
			b:    Float64(0),
			want: Float64(math.Inf(1)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictFloatDivide(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictNumericParseUint(t *testing.T) {
	tests := map[string]struct {
		str     string
		base    int
		bitSize int
		want    uint64
		err     *Error
	}{
		"explicit decimal": {
			str:     "123",
			base:    10,
			bitSize: 8,
			want:    123,
		},
		"implicit decimal": {
			str:     "123",
			base:    0,
			bitSize: 8,
			want:    123,
		},
		"implicit decimal with underscores": {
			str:     "2_500",
			base:    0,
			bitSize: 16,
			want:    2500,
		},
		"8bit implicit decimal out of range": {
			str:     "300",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit hex": {
			str:     "ff",
			base:    16,
			bitSize: 8,
			want:    255,
		},
		"implicit hex": {
			str:     "0xff",
			base:    0,
			bitSize: 8,
			want:    255,
		},
		"implicit hex with underscores": {
			str:     "0x12_34",
			base:    0,
			bitSize: 16,
			want:    4660,
		},
		"8bit implicit hex out of range": {
			str:     "0xfff",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit duodecimal": {
			str:     "1a",
			base:    12,
			bitSize: 8,
			want:    22,
		},
		"implicit duodecimal": {
			str:     "0d1a",
			base:    0,
			bitSize: 8,
			want:    22,
		},
		"implicit duodecimal with underscores": {
			str:     "0d12_34",
			base:    0,
			bitSize: 16,
			want:    2056,
		},
		"8bit implicit duodecimal out of range": {
			str:     "0d194",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit octal": {
			str:     "67",
			base:    8,
			bitSize: 8,
			want:    55,
		},
		"implicit octal": {
			str:     "0o67",
			base:    0,
			bitSize: 8,
			want:    55,
		},
		"implicit octal with underscores": {
			str:     "0o12_34",
			base:    0,
			bitSize: 16,
			want:    668,
		},
		"8bit implicit octal out of range": {
			str:     "0o400",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit quaternary": {
			str:     "33",
			base:    4,
			bitSize: 8,
			want:    15,
		},
		"implicit quaternary": {
			str:     "0q33",
			base:    0,
			bitSize: 8,
			want:    15,
		},
		"implicit quaternary with underscores": {
			str:     "0q12_33",
			base:    0,
			bitSize: 8,
			want:    111,
		},
		"8bit implicit quaternary out of range": {
			str:     "0q10000",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"explicit binary": {
			str:     "101",
			base:    2,
			bitSize: 8,
			want:    5,
		},
		"implicit binary": {
			str:     "0b101",
			base:    0,
			bitSize: 8,
			want:    5,
		},
		"implicit binary with underscores": {
			str:     "0b100_111",
			base:    0,
			bitSize: 16,
			want:    39,
		},
		"8bit implicit binary out of range": {
			str:     "0b100000000",
			base:    0,
			bitSize: 8,
			want:    255,
			err:     Errorf(FormatErrorClass, "value overflows"),
		},
		"64bit decimal": {
			str:     "123",
			base:    10,
			bitSize: 64,
			want:    123,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictParseUint(tc.str, tc.base, tc.bitSize)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}),
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictIntRightBitshift(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"shift by String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be used as a bitshift operand"),
		},
		"shift by Float and return an error": {
			a:   Int64(5),
			b:   Float(3.2),
			err: NewError(TypeErrorClass, "`Std::Float` can't be used as a bitshift operand"),
		},
		"shift by Int32": {
			a:    Int64(234),
			b:    Int32(2),
			want: Int64(58),
		},
		"shift by UInt8": {
			a:    Int64(234),
			b:    UInt8(2),
			want: Int64(58),
		},
		"shift by SmallInt": {
			a:    Int64(234),
			b:    SmallInt(2),
			want: Int64(58),
		},
		"shift by BigInt": {
			a:    Int64(234),
			b:    NewBigInt(2),
			want: Int64(58),
		},
		"shift by large BigInt": {
			a:    Int64(234),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: Int64(0),
		},
		"shift by 10 >> 1": {
			a:    Int64(10),
			b:    Int64(1),
			want: Int64(5),
		},
		"shift by 10 >> 255": {
			a:    Int64(10),
			b:    Int64(255),
			want: Int64(0),
		},
		"shift by 25 >> 2": {
			a:    Int64(25),
			b:    Int64(2),
			want: Int64(6),
		},
		"shift by 25 >> -2": {
			a:    Int64(25),
			b:    Int64(-2),
			want: Int64(100),
		},
		"shift by -6 >> 1": {
			a:    Int64(-6),
			b:    Int64(1),
			want: Int64(-3),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntRightBitshift(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictIntLogicalRightBitshift(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"shift by String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be used as a bitshift operand"),
		},
		"shift by Float and return an error": {
			a:   Int64(5),
			b:   Float(3.2),
			err: NewError(TypeErrorClass, "`Std::Float` can't be used as a bitshift operand"),
		},
		"shift by Int32": {
			a:    Int64(234),
			b:    Int32(2),
			want: Int64(58),
		},
		"shift by UInt8": {
			a:    Int64(234),
			b:    UInt8(2),
			want: Int64(58),
		},
		"shift by SmallInt": {
			a:    Int64(234),
			b:    SmallInt(2),
			want: Int64(58),
		},
		"shift by BigInt": {
			a:    Int64(234),
			b:    NewBigInt(2),
			want: Int64(58),
		},
		"shift by large BigInt": {
			a:    Int64(234),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: Int64(0),
		},
		"shift by 10 >>> 1": {
			a:    Int64(10),
			b:    Int64(1),
			want: Int64(5),
		},
		"shift by 10 >>> 255": {
			a:    Int64(10),
			b:    Int64(255),
			want: Int64(0),
		},
		"shift by 25 >>> 2": {
			a:    Int64(25),
			b:    Int64(2),
			want: Int64(6),
		},
		"shift by 25 >>> -2": {
			a:    Int64(25),
			b:    Int64(-2),
			want: Int64(100),
		},
		"shift by -6 >>> 1": {
			a:    Int64(-6),
			b:    Int64(1),
			want: Int64(9223372036854775805),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntLogicalRightBitshift(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictIntLeftBitshift(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"shift by String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be used as a bitshift operand"),
		},
		"shift by Float and return an error": {
			a:   Int64(5),
			b:   Float(3.2),
			err: NewError(TypeErrorClass, "`Std::Float` can't be used as a bitshift operand"),
		},
		"shift by Int32": {
			a:    Int64(234),
			b:    Int32(2),
			want: Int64(936),
		},
		"shift by UInt8": {
			a:    Int64(234),
			b:    UInt8(2),
			want: Int64(936),
		},
		"shift by SmallInt": {
			a:    Int64(234),
			b:    SmallInt(2),
			want: Int64(936),
		},
		"shift by BigInt": {
			a:    Int64(234),
			b:    NewBigInt(2),
			want: Int64(936),
		},
		"shift by large BigInt": {
			a:    Int64(234),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: Int64(0),
		},
		"shift by 10 << 1": {
			a:    Int64(10),
			b:    Int64(1),
			want: Int64(20),
		},
		"shift by 10 << 255": {
			a:    Int64(10),
			b:    Int64(255),
			want: Int64(0),
		},
		"shift by 25 << 2": {
			a:    Int64(25),
			b:    Int64(2),
			want: Int64(100),
		},
		"shift by 25 << -2": {
			a:    Int64(25),
			b:    Int64(-2),
			want: Int64(6),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntLeftBitshift(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
