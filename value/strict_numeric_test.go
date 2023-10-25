package value

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStrictFloat_Exponentiate(t *testing.T) {
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

func TestStrictInt_Exponentiate(t *testing.T) {
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

func TestStrictNumeric_Add(t *testing.T) {
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

func TestStrictNumeric_Subtract(t *testing.T) {
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

func TestStrictNumeric_Multiply(t *testing.T) {
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

func TestStrictNumeric_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Bool
		err  *Error
	}{
		"String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"15i64 > 30i64": {
			a:    Int64(15),
			b:    Int64(30),
			want: False,
		},
		"780i64 > -800i64": {
			a:    Int64(780),
			b:    Int64(-800),
			want: True,
		},
		"15i64 > 15i64": {
			a:    Int64(15),
			b:    Int64(15),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictNumericGreaterThan(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
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

func TestStrictNumeric_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Bool
		err  *Error
	}{
		"String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"15i64 >= 30i64": {
			a:    Int64(15),
			b:    Int64(30),
			want: False,
		},
		"780i64 >= -800i64": {
			a:    Int64(780),
			b:    Int64(-800),
			want: True,
		},
		"15i64 >= 15i64": {
			a:    Int64(15),
			b:    Int64(15),
			want: True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictNumericGreaterThanEqual(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
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

func TestStrictNumeric_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Bool
		err  *Error
	}{
		"String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"15i64 < 30i64": {
			a:    Int64(15),
			b:    Int64(30),
			want: True,
		},
		"780i64 < -800i64": {
			a:    Int64(780),
			b:    Int64(-800),
			want: False,
		},
		"15i64 < 15i64": {
			a:    Int64(15),
			b:    Int64(15),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictNumericLessThan(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
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

func TestStrictNumeric_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Bool
		err  *Error
	}{
		"String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"15i64 <= 30i64": {
			a:    Int64(15),
			b:    Int64(30),
			want: True,
		},
		"780i64 <= -800i64": {
			a:    Int64(780),
			b:    Int64(-800),
			want: False,
		},
		"15i64 <= 15i64": {
			a:    Int64(15),
			b:    Int64(15),
			want: True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictNumericLessThanEqual(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
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

func TestStrictNumeric_StrictEqual(t *testing.T) {
	tests := map[string]struct {
		a    Float64
		b    Value
		want Bool
		err  *Error
	}{
		"String 5.5f64 === '5.5'": {
			a:    Float64(5.5),
			b:    String("5.5"),
			want: False,
		},

		"Int64 5.0f64 === 5i64": {
			a:    Float64(5),
			b:    Int64(5),
			want: False,
		},
		"Int64 5.5f64 === 5i64": {
			a:    Float64(5.5),
			b:    Int64(5),
			want: False,
		},
		"Int64 NaN === 0i64": {
			a:    Float64NaN(),
			b:    Int64(0),
			want: False,
		},
		"Int64 +Inf === 69i64": {
			a:    Float64Inf(),
			b:    Int64(69),
			want: False,
		},
		"Int64 -Inf === -89i64": {
			a:    Float64NegInf(),
			b:    Int64(-89),
			want: False,
		},

		"UInt64 5.0f64 === 5u64": {
			a:    Float64(5),
			b:    UInt64(5),
			want: False,
		},
		"UInt64 -5.0f64 === 5u64": {
			a:    Float64(-5),
			b:    UInt64(5),
			want: False,
		},
		"UInt64 5.5f64 === 5u64": {
			a:    Float64(5.5),
			b:    UInt64(5),
			want: False,
		},
		"UInt64 NaN === 0u64": {
			a:    Float64NaN(),
			b:    UInt64(0),
			want: False,
		},
		"UInt64 +Inf === 69u64": {
			a:    Float64Inf(),
			b:    UInt64(69),
			want: False,
		},
		"UInt64 -Inf === 89u64": {
			a:    Float64NegInf(),
			b:    UInt64(89),
			want: False,
		},

		"Int32 5.0f64 === 5i32": {
			a:    Float64(5),
			b:    Int32(5),
			want: False,
		},
		"Int32 5.5f64 === 5i32": {
			a:    Float64(5.5),
			b:    Int32(5),
			want: False,
		},
		"Int32 NaN === 0i32": {
			a:    Float64NaN(),
			b:    Int32(0),
			want: False,
		},
		"Int32 +Inf === 69i32": {
			a:    Float64Inf(),
			b:    Int32(69),
			want: False,
		},
		"Int32 -Inf === -89i32": {
			a:    Float64NegInf(),
			b:    Int32(-89),
			want: False,
		},

		"UInt32 5.0f64 === 5u32": {
			a:    Float64(5),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 -5.0f64 === 5u32": {
			a:    Float64(-5),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 5.5f64 === 5u32": {
			a:    Float64(5.5),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 NaN === 0u32": {
			a:    Float64NaN(),
			b:    UInt32(0),
			want: False,
		},
		"UInt32 +Inf === 69u32": {
			a:    Float64Inf(),
			b:    UInt32(69),
			want: False,
		},
		"UInt32 -Inf === 89u32": {
			a:    Float64NegInf(),
			b:    UInt32(89),
			want: False,
		},

		"Int16 5.0f64 === 5i16": {
			a:    Float64(5),
			b:    Int16(5),
			want: False,
		},
		"Int16 5.5f64 === 5i16": {
			a:    Float64(5.5),
			b:    Int16(5),
			want: False,
		},
		"Int16 NaN === 0i16": {
			a:    Float64NaN(),
			b:    Int16(0),
			want: False,
		},
		"Int16 +Inf === 69i16": {
			a:    Float64Inf(),
			b:    Int16(69),
			want: False,
		},
		"Int16 -Inf === -89i16": {
			a:    Float64NegInf(),
			b:    Int16(-89),
			want: False,
		},

		"UInt16 5.0f64 === 5u16": {
			a:    Float64(5),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 -5.0f64 === 5u16": {
			a:    Float64(-5),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 5.5f64 === 5u16": {
			a:    Float64(5.5),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 NaN === 0u16": {
			a:    Float64NaN(),
			b:    UInt16(0),
			want: False,
		},
		"UInt16 +Inf === 69u16": {
			a:    Float64Inf(),
			b:    UInt16(69),
			want: False,
		},
		"UInt16 -Inf === 89u16": {
			a:    Float64NegInf(),
			b:    UInt16(89),
			want: False,
		},

		"Int8 5.0f64 === 5i8": {
			a:    Float64(5),
			b:    Int8(5),
			want: False,
		},
		"Int8 5.5f64 === 5i8": {
			a:    Float64(5.5),
			b:    Int8(5),
			want: False,
		},
		"Int8 NaN === 0i8": {
			a:    Float64NaN(),
			b:    Int8(0),
			want: False,
		},
		"Int8 +Inf === 69i8": {
			a:    Float64Inf(),
			b:    Int8(69),
			want: False,
		},
		"Int8 -Inf === -89i8": {
			a:    Float64NegInf(),
			b:    Int8(-89),
			want: False,
		},

		"UInt8 5.0f64 === 5u8": {
			a:    Float64(5),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 -5.0f64 === 5u8": {
			a:    Float64(-5),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 5.5f64 === 5u8": {
			a:    Float64(5.5),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 NaN === 0u8": {
			a:    Float64NaN(),
			b:    UInt8(0),
			want: False,
		},
		"UInt8 +Inf === 69u8": {
			a:    Float64Inf(),
			b:    UInt8(69),
			want: False,
		},
		"UInt8 -Inf === 89u8": {
			a:    Float64NegInf(),
			b:    UInt8(89),
			want: False,
		},

		"Float64 21.9f64 === 21.9f64": {
			a:    Float64(21.9),
			b:    Float64(21.9),
			want: True,
		},
		"Float64 21.9f64 === 38.0f64": {
			a:    Float64(21.9),
			b:    Float64(38),
			want: False,
		},
		"Float64 NaN === NaN": {
			a:    Float64NaN(),
			b:    Float64NaN(),
			want: False,
		},
		"Float64 +Inf === +Inf": {
			a:    Float64Inf(),
			b:    Float64Inf(),
			want: True,
		},
		"Float64 -Inf === -Inf": {
			a:    Float64NegInf(),
			b:    Float64NegInf(),
			want: True,
		},
		"Float64 +Inf === -Inf": {
			a:    Float64Inf(),
			b:    Float64NegInf(),
			want: False,
		},
		"Float64 -Inf === +Inf": {
			a:    Float64NegInf(),
			b:    Float64Inf(),
			want: False,
		},
		"Float64 8.5f64 === +Inf": {
			a:    Float64(8.5),
			b:    Float64Inf(),
			want: False,
		},
		"Float64 +Inf === 98.0f64": {
			a:    Float64Inf(),
			b:    Float64(98),
			want: False,
		},

		"Float32 21.0f64 === 21.0f32": {
			a:    Float64(21),
			b:    Float32(21),
			want: False,
		},
		"Float32 21.9f64 === 38.0f32": {
			a:    Float64(21.9),
			b:    Float32(38),
			want: False,
		},
		"Float32 NaN === NaN": {
			a:    Float64NaN(),
			b:    Float32NaN(),
			want: False,
		},
		"Float32 +Inf === +Inf": {
			a:    Float64Inf(),
			b:    Float32Inf(),
			want: False,
		},
		"Float32 -Inf === -Inf": {
			a:    Float64NegInf(),
			b:    Float32NegInf(),
			want: False,
		},
		"Float32 +Inf === -Inf": {
			a:    Float64Inf(),
			b:    Float32NegInf(),
			want: False,
		},
		"Float32 -Inf === +Inf": {
			a:    Float64NegInf(),
			b:    Float32Inf(),
			want: False,
		},
		"Float32 8.5f64 === +Inf": {
			a:    Float64(8.5),
			b:    Float32Inf(),
			want: False,
		},
		"Float32 +Inf === 98.0f32": {
			a:    Float64Inf(),
			b:    Float32(98),
			want: False,
		},

		"SmallInt 16.0f64 === 16": {
			a:    Float64(16),
			b:    SmallInt(16),
			want: False,
		},
		"SmallInt 16.5f64 === 16": {
			a:    Float64(16.5),
			b:    SmallInt(16),
			want: False,
		},
		"SmallInt NaN === 0": {
			a:    Float64NaN(),
			b:    SmallInt(0),
			want: False,
		},
		"SmallInt +Inf === 69": {
			a:    Float64Inf(),
			b:    SmallInt(69),
			want: False,
		},
		"SmallInt -Inf === -89": {
			a:    Float64NegInf(),
			b:    SmallInt(-89),
			want: False,
		},

		"BigInt 16.0f64 === 16": {
			a:    Float64(16),
			b:    NewBigInt(16),
			want: False,
		},
		"BigInt 16.5f64 === 16": {
			a:    Float64(16.5),
			b:    NewBigInt(16),
			want: False,
		},
		"BigInt NaN === 0": {
			a:    Float64NaN(),
			b:    NewBigInt(0),
			want: False,
		},
		"BigInt +Inf === 69": {
			a:    Float64Inf(),
			b:    NewBigInt(69),
			want: False,
		},
		"BigInt -Inf === -89": {
			a:    Float64NegInf(),
			b:    NewBigInt(-89),
			want: False,
		},

		"Float 21.9f64 === 21.9": {
			a:    Float64(21.9),
			b:    Float(21.9),
			want: False,
		},
		"Float 21.9f64 === 38.0": {
			a:    Float64(21.9),
			b:    Float(38),
			want: False,
		},
		"Float NaN === NaN": {
			a:    Float64NaN(),
			b:    FloatNaN(),
			want: False,
		},
		"Float +Inf === +Inf": {
			a:    Float64Inf(),
			b:    FloatInf(),
			want: False,
		},
		"Float -Inf === -Inf": {
			a:    Float64NegInf(),
			b:    FloatNegInf(),
			want: False,
		},
		"Float +Inf === -Inf": {
			a:    Float64Inf(),
			b:    FloatNegInf(),
			want: False,
		},
		"Float -Inf === +Inf": {
			a:    Float64NegInf(),
			b:    FloatInf(),
			want: False,
		},
		"Float 8.5f64 === +Inf": {
			a:    Float64(8.5),
			b:    FloatInf(),
			want: False,
		},
		"Float +Inf === 98.0": {
			a:    Float64Inf(),
			b:    Float(98),
			want: False,
		},

		"BigFloat 21.9f64 === 21.9bf": {
			a:    Float64(21.9),
			b:    NewBigFloat(21.9),
			want: False,
		},
		"BigFloat 21.9f64 === 38.0bf": {
			a:    Float64(21.9),
			b:    NewBigFloat(38),
			want: False,
		},
		"BigFloat NaN === NaN": {
			a:    Float64NaN(),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat +Inf === +Inf": {
			a:    Float64Inf(),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat -Inf === -Inf": {
			a:    Float64NegInf(),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat +Inf === -Inf": {
			a:    Float64Inf(),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat -Inf === +Inf": {
			a:    Float64NegInf(),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 8.5f64 === +Inf": {
			a:    Float64(8.5),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat +Inf === 98.0bf": {
			a:    Float64Inf(),
			b:    NewBigFloat(98),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := StrictNumericStrictEqual(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				Float32Comparer,
				Float64Comparer,
				BigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictFloat_Equal(t *testing.T) {
	tests := map[string]struct {
		a    Float64
		b    Value
		want Bool
		err  *Error
	}{
		"String 5.5f64 == '5.5'": {
			a:    Float64(5.5),
			b:    String("5.5"),
			want: False,
		},

		"Int64 5.0f64 == 5i64": {
			a:    Float64(5),
			b:    Int64(5),
			want: True,
		},
		"Int64 5.5f64 == 5i64": {
			a:    Float64(5.5),
			b:    Int64(5),
			want: False,
		},
		"Int64 NaN == 0i64": {
			a:    Float64NaN(),
			b:    Int64(0),
			want: False,
		},
		"Int64 +Inf == 69i64": {
			a:    Float64Inf(),
			b:    Int64(69),
			want: False,
		},
		"Int64 -Inf == -89i64": {
			a:    Float64NegInf(),
			b:    Int64(-89),
			want: False,
		},

		"UInt64 5.0f64 == 5u64": {
			a:    Float64(5),
			b:    UInt64(5),
			want: True,
		},
		"UInt64 -5.0f64 == 5u64": {
			a:    Float64(-5),
			b:    UInt64(5),
			want: False,
		},
		"UInt64 5.5f64 == 5u64": {
			a:    Float64(5.5),
			b:    UInt64(5),
			want: False,
		},
		"UInt64 NaN == 0u64": {
			a:    Float64NaN(),
			b:    UInt64(0),
			want: False,
		},
		"UInt64 +Inf == 69u64": {
			a:    Float64Inf(),
			b:    UInt64(69),
			want: False,
		},
		"UInt64 -Inf == 89u64": {
			a:    Float64NegInf(),
			b:    UInt64(89),
			want: False,
		},

		"Int32 5.0f64 == 5i32": {
			a:    Float64(5),
			b:    Int32(5),
			want: True,
		},
		"Int32 5.5f64 == 5i32": {
			a:    Float64(5.5),
			b:    Int32(5),
			want: False,
		},
		"Int32 NaN == 0i32": {
			a:    Float64NaN(),
			b:    Int32(0),
			want: False,
		},
		"Int32 +Inf == 69i32": {
			a:    Float64Inf(),
			b:    Int32(69),
			want: False,
		},
		"Int32 -Inf == -89i32": {
			a:    Float64NegInf(),
			b:    Int32(-89),
			want: False,
		},

		"UInt32 5.0f64 == 5u32": {
			a:    Float64(5),
			b:    UInt32(5),
			want: True,
		},
		"UInt32 -5.0f64 == 5u32": {
			a:    Float64(-5),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 5.5f64 == 5u32": {
			a:    Float64(5.5),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 NaN == 0u32": {
			a:    Float64NaN(),
			b:    UInt32(0),
			want: False,
		},
		"UInt32 +Inf == 69u32": {
			a:    Float64Inf(),
			b:    UInt32(69),
			want: False,
		},
		"UInt32 -Inf == 89u32": {
			a:    Float64NegInf(),
			b:    UInt32(89),
			want: False,
		},

		"Int16 5.0f64 == 5i16": {
			a:    Float64(5),
			b:    Int16(5),
			want: True,
		},
		"Int16 5.5f64 == 5i16": {
			a:    Float64(5.5),
			b:    Int16(5),
			want: False,
		},
		"Int16 NaN == 0i16": {
			a:    Float64NaN(),
			b:    Int16(0),
			want: False,
		},
		"Int16 +Inf == 69i16": {
			a:    Float64Inf(),
			b:    Int16(69),
			want: False,
		},
		"Int16 -Inf == -89i16": {
			a:    Float64NegInf(),
			b:    Int16(-89),
			want: False,
		},

		"UInt16 5.0f64 == 5u16": {
			a:    Float64(5),
			b:    UInt16(5),
			want: True,
		},
		"UInt16 -5.0f64 == 5u16": {
			a:    Float64(-5),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 5.5f64 == 5u16": {
			a:    Float64(5.5),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 NaN == 0u16": {
			a:    Float64NaN(),
			b:    UInt16(0),
			want: False,
		},
		"UInt16 +Inf == 69u16": {
			a:    Float64Inf(),
			b:    UInt16(69),
			want: False,
		},
		"UInt16 -Inf == 89u16": {
			a:    Float64NegInf(),
			b:    UInt16(89),
			want: False,
		},

		"Int8 5.0f64 == 5i8": {
			a:    Float64(5),
			b:    Int8(5),
			want: True,
		},
		"Int8 5.5f64 == 5i8": {
			a:    Float64(5.5),
			b:    Int8(5),
			want: False,
		},
		"Int8 NaN == 0i8": {
			a:    Float64NaN(),
			b:    Int8(0),
			want: False,
		},
		"Int8 +Inf == 69i8": {
			a:    Float64Inf(),
			b:    Int8(69),
			want: False,
		},
		"Int8 -Inf == -89i8": {
			a:    Float64NegInf(),
			b:    Int8(-89),
			want: False,
		},

		"UInt8 5.0f64 == 5u8": {
			a:    Float64(5),
			b:    UInt8(5),
			want: True,
		},
		"UInt8 -5.0f64 == 5u8": {
			a:    Float64(-5),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 5.5f64 == 5u8": {
			a:    Float64(5.5),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 NaN == 0u8": {
			a:    Float64NaN(),
			b:    UInt8(0),
			want: False,
		},
		"UInt8 +Inf == 69u8": {
			a:    Float64Inf(),
			b:    UInt8(69),
			want: False,
		},
		"UInt8 -Inf == 89u8": {
			a:    Float64NegInf(),
			b:    UInt8(89),
			want: False,
		},

		"Float64 21.9f64 == 21.9f64": {
			a:    Float64(21.9),
			b:    Float64(21.9),
			want: True,
		},
		"Float64 21.9f64 == 38.0f64": {
			a:    Float64(21.9),
			b:    Float64(38),
			want: False,
		},
		"Float64 NaN == NaN": {
			a:    Float64NaN(),
			b:    Float64NaN(),
			want: False,
		},
		"Float64 +Inf == +Inf": {
			a:    Float64Inf(),
			b:    Float64Inf(),
			want: True,
		},
		"Float64 -Inf == -Inf": {
			a:    Float64NegInf(),
			b:    Float64NegInf(),
			want: True,
		},
		"Float64 +Inf == -Inf": {
			a:    Float64Inf(),
			b:    Float64NegInf(),
			want: False,
		},
		"Float64 -Inf == +Inf": {
			a:    Float64NegInf(),
			b:    Float64Inf(),
			want: False,
		},
		"Float64 8.5f64 == +Inf": {
			a:    Float64(8.5),
			b:    Float64Inf(),
			want: False,
		},
		"Float64 +Inf == 98.0f64": {
			a:    Float64Inf(),
			b:    Float64(98),
			want: False,
		},

		"Float32 21.0f64 == 21.0f32": {
			a:    Float64(21),
			b:    Float32(21),
			want: True,
		},
		"Float32 21.9f64 == 38.0f32": {
			a:    Float64(21.9),
			b:    Float32(38),
			want: False,
		},
		"Float32 NaN == NaN": {
			a:    Float64NaN(),
			b:    Float32NaN(),
			want: False,
		},
		"Float32 +Inf == +Inf": {
			a:    Float64Inf(),
			b:    Float32Inf(),
			want: True,
		},
		"Float32 -Inf == -Inf": {
			a:    Float64NegInf(),
			b:    Float32NegInf(),
			want: True,
		},
		"Float32 +Inf == -Inf": {
			a:    Float64Inf(),
			b:    Float32NegInf(),
			want: False,
		},
		"Float32 -Inf == +Inf": {
			a:    Float64NegInf(),
			b:    Float32Inf(),
			want: False,
		},
		"Float32 8.5f64 == +Inf": {
			a:    Float64(8.5),
			b:    Float32Inf(),
			want: False,
		},
		"Float32 +Inf == 98.0f32": {
			a:    Float64Inf(),
			b:    Float32(98),
			want: False,
		},

		"SmallInt 16.0f64 == 16": {
			a:    Float64(16),
			b:    SmallInt(16),
			want: True,
		},
		"SmallInt 16.5f64 == 16": {
			a:    Float64(16.5),
			b:    SmallInt(16),
			want: False,
		},
		"SmallInt NaN == 0": {
			a:    Float64NaN(),
			b:    SmallInt(0),
			want: False,
		},
		"SmallInt +Inf == 69": {
			a:    Float64Inf(),
			b:    SmallInt(69),
			want: False,
		},
		"SmallInt -Inf == -89": {
			a:    Float64NegInf(),
			b:    SmallInt(-89),
			want: False,
		},

		"BigInt 16.0f64 == 16": {
			a:    Float64(16),
			b:    NewBigInt(16),
			want: True,
		},
		"BigInt 16.5f64 == 16": {
			a:    Float64(16.5),
			b:    NewBigInt(16),
			want: False,
		},
		"BigInt NaN == 0": {
			a:    Float64NaN(),
			b:    NewBigInt(0),
			want: False,
		},
		"BigInt +Inf == 69": {
			a:    Float64Inf(),
			b:    NewBigInt(69),
			want: False,
		},
		"BigInt -Inf == -89": {
			a:    Float64NegInf(),
			b:    NewBigInt(-89),
			want: False,
		},

		"Float 21.9f64 == 21.9": {
			a:    Float64(21.9),
			b:    Float(21.9),
			want: True,
		},
		"Float 21.9f64 == 38.0": {
			a:    Float64(21.9),
			b:    Float(38),
			want: False,
		},
		"Float NaN == NaN": {
			a:    Float64NaN(),
			b:    FloatNaN(),
			want: False,
		},
		"Float +Inf == +Inf": {
			a:    Float64Inf(),
			b:    FloatInf(),
			want: True,
		},
		"Float -Inf == -Inf": {
			a:    Float64NegInf(),
			b:    FloatNegInf(),
			want: True,
		},
		"Float +Inf == -Inf": {
			a:    Float64Inf(),
			b:    FloatNegInf(),
			want: False,
		},
		"Float -Inf == +Inf": {
			a:    Float64NegInf(),
			b:    FloatInf(),
			want: False,
		},
		"Float 8.5f64 == +Inf": {
			a:    Float64(8.5),
			b:    FloatInf(),
			want: False,
		},
		"Float +Inf == 98.0": {
			a:    Float64Inf(),
			b:    Float(98),
			want: False,
		},

		"BigFloat 21.9f64 == 21.9bf": {
			a:    Float64(21.9),
			b:    NewBigFloat(21.9),
			want: True,
		},
		"BigFloat 21.9f64 == 38.0bf": {
			a:    Float64(21.9),
			b:    NewBigFloat(38),
			want: False,
		},
		"BigFloat NaN == NaN": {
			a:    Float64NaN(),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat +Inf == +Inf": {
			a:    Float64Inf(),
			b:    BigFloatInf(),
			want: True,
		},
		"BigFloat -Inf == -Inf": {
			a:    Float64NegInf(),
			b:    BigFloatNegInf(),
			want: True,
		},
		"BigFloat +Inf == -Inf": {
			a:    Float64Inf(),
			b:    BigFloatNegInf(),
			want: False,
		},
		"BigFloat -Inf == +Inf": {
			a:    Float64NegInf(),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 8.5f64 == +Inf": {
			a:    Float64(8.5),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat +Inf == 98.0bf": {
			a:    Float64Inf(),
			b:    NewBigFloat(98),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := StrictFloatEqual(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				Float32Comparer,
				Float64Comparer,
				BigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictSignedInt_Equal(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Bool
		err  *Error
	}{
		"String 5i64 == '5'": {
			a:    Int64(5),
			b:    String("5.5"),
			want: False,
		},

		"Int64 5i64 == 5i64": {
			a:    Int64(5),
			b:    Int64(5),
			want: True,
		},
		"Int64 4i64 == 5i64": {
			a:    Int64(4),
			b:    Int64(5),
			want: False,
		},
		"Int64 5i64 == -5i64": {
			a:    Int64(5),
			b:    Int64(-5),
			want: False,
		},

		"UInt64 5i64 == 5u64": {
			a:    Int64(5),
			b:    UInt64(5),
			want: True,
		},
		"UInt64 -5i64 == 5u64": {
			a:    Int64(-5),
			b:    UInt64(5),
			want: False,
		},

		"Int32 5i64 == 5i32": {
			a:    Int64(5),
			b:    Int32(5),
			want: True,
		},
		"Int32 -5i64 == 5i32": {
			a:    Int64(-5),
			b:    Int32(5),
			want: False,
		},
		"Int32 5i64 == -5i32": {
			a:    Int64(5),
			b:    Int32(-5),
			want: False,
		},
		"Int32 -5i64 == -5i32": {
			a:    Int64(-5),
			b:    Int32(-5),
			want: True,
		},

		"UInt32 5i64 == 5u32": {
			a:    Int64(5),
			b:    UInt32(5),
			want: True,
		},
		"UInt32 -5i64 == 5u32": {
			a:    Int64(-5),
			b:    UInt32(5),
			want: False,
		},
		"UInt32 4i64 == 5u32": {
			a:    Int64(4),
			b:    UInt32(5),
			want: False,
		},

		"Int16 5i64 == 5i16": {
			a:    Int64(5),
			b:    Int16(5),
			want: True,
		},
		"Int16 -5i64 == 5i16": {
			a:    Int64(-5),
			b:    Int16(5),
			want: False,
		},
		"Int16 5i64 == -5i16": {
			a:    Int64(5),
			b:    Int16(-5),
			want: False,
		},
		"Int16 -5i64 == -5i16": {
			a:    Int64(-5),
			b:    Int16(-5),
			want: True,
		},
		"Int16 4i64 == 5i16": {
			a:    Int64(4),
			b:    Int16(5),
			want: False,
		},

		"UInt16 5i64 == 5u16": {
			a:    Int64(5),
			b:    UInt16(5),
			want: True,
		},
		"UInt16 -5i64 == 5u16": {
			a:    Int64(-5),
			b:    UInt16(5),
			want: False,
		},
		"UInt16 4i64 == 5u16": {
			a:    Int64(4),
			b:    UInt16(5),
			want: False,
		},

		"Int8 5i64 == 5i8": {
			a:    Int64(5),
			b:    Int8(5),
			want: True,
		},
		"Int8 4i64 == 5i8": {
			a:    Int64(4),
			b:    Int8(5),
			want: False,
		},
		"Int8 -5i64 == 5i8": {
			a:    Int64(-5),
			b:    Int8(5),
			want: False,
		},
		"Int8 5i64 == -5i8": {
			a:    Int64(5),
			b:    Int8(-5),
			want: False,
		},
		"Int8 -5i64 == -5i8": {
			a:    Int64(-5),
			b:    Int8(-5),
			want: True,
		},

		"UInt8 5i64 == 5u8": {
			a:    Int64(5),
			b:    UInt8(5),
			want: True,
		},
		"UInt8 4i64 == 5u8": {
			a:    Int64(4),
			b:    UInt8(5),
			want: False,
		},
		"UInt8 -5i64 == 5u8": {
			a:    Int64(-5),
			b:    UInt8(5),
			want: False,
		},

		"Float64 21i64 == 21.0f64": {
			a:    Int64(21),
			b:    Float64(21),
			want: True,
		},
		"Float64 21i64 == 21.5f64": {
			a:    Int64(21),
			b:    Float64(21.5),
			want: False,
		},
		"Float64 21i64 == 38.0f64": {
			a:    Int64(21),
			b:    Float64(38),
			want: False,
		},
		"Float64 0i64 == NaN": {
			a:    Int64(0),
			b:    Float64NaN(),
			want: False,
		},
		"Float64 8i64 == +Inf": {
			a:    Int64(8),
			b:    Float64Inf(),
			want: False,
		},
		"Float64 8i64 == -Inf": {
			a:    Int64(8),
			b:    Float64NegInf(),
			want: False,
		},

		"Float32 21i64 == 21.0f32": {
			a:    Int64(21),
			b:    Float32(21),
			want: True,
		},
		"Float32 21i64 == 21.5f32": {
			a:    Int64(21),
			b:    Float32(21.5),
			want: False,
		},
		"Float32 21i64 == 38.0f32": {
			a:    Int64(21),
			b:    Float32(38),
			want: False,
		},
		"Float32 0i64 == NaN": {
			a:    Int64(0),
			b:    Float32NaN(),
			want: False,
		},
		"Float32 8i64 == +Inf": {
			a:    Int64(8),
			b:    Float32Inf(),
			want: False,
		},
		"Float32 8i64 == -Inf": {
			a:    Int64(8),
			b:    Float32NegInf(),
			want: False,
		},

		"SmallInt 16i64 == 16": {
			a:    Int64(16),
			b:    SmallInt(16),
			want: True,
		},
		"SmallInt 97i64 == -97": {
			a:    Int64(97),
			b:    SmallInt(-97),
			want: False,
		},
		"SmallInt -6i64 == 6": {
			a:    Int64(-6),
			b:    SmallInt(6),
			want: False,
		},
		"SmallInt -120i64 == -120": {
			a:    Int64(-120),
			b:    SmallInt(-120),
			want: True,
		},

		"BigInt 16i64 == 16": {
			a:    Int64(16),
			b:    NewBigInt(16),
			want: True,
		},
		"BigInt 97i64 == -97": {
			a:    Int64(97),
			b:    NewBigInt(-97),
			want: False,
		},
		"BigInt -6i64 == 6": {
			a:    Int64(-6),
			b:    NewBigInt(6),
			want: False,
		},
		"BigInt -120i64 == -120": {
			a:    Int64(-120),
			b:    NewBigInt(-120),
			want: True,
		},

		"Float 21i64 == 21.0": {
			a:    Int64(21),
			b:    Float(21),
			want: True,
		},
		"Float 21i64 == 21.5": {
			a:    Int64(21),
			b:    Float(21.5),
			want: False,
		},
		"Float 21i64 == 38.0": {
			a:    Int64(21),
			b:    Float(38),
			want: False,
		},
		"Float 0i64 == NaN": {
			a:    Int64(0),
			b:    FloatNaN(),
			want: False,
		},
		"Float 8i64 == +Inf": {
			a:    Int64(8),
			b:    FloatInf(),
			want: False,
		},
		"Float 8i64 == -Inf": {
			a:    Int64(8),
			b:    FloatNegInf(),
			want: False,
		},

		"BigFloat 21i64 == 21.0bf": {
			a:    Int64(21),
			b:    NewBigFloat(21),
			want: True,
		},
		"BigFloat 21i64 == 21.5bf": {
			a:    Int64(21),
			b:    NewBigFloat(21.5),
			want: False,
		},
		"BigFloat 21i64 == 38.0bf": {
			a:    Int64(21),
			b:    NewBigFloat(38),
			want: False,
		},
		"BigFloat 0i64 == NaN": {
			a:    Int64(0),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat 8i64 == +Inf": {
			a:    Int64(8),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 8i64 == -Inf": {
			a:    Int64(8),
			b:    BigFloatNegInf(),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := StrictSignedIntEqual(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				Float32Comparer,
				Float64Comparer,
				BigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictUnsignedInt_Equal(t *testing.T) {
	tests := map[string]struct {
		a    UInt64
		b    Value
		want Bool
		err  *Error
	}{
		"String 5u64 == '5'": {
			a:    UInt64(5),
			b:    String("5.5"),
			want: False,
		},

		"Int64 5u64 == 5i64": {
			a:    UInt64(5),
			b:    Int64(5),
			want: True,
		},
		"Int64 4u64 == 5i64": {
			a:    UInt64(4),
			b:    Int64(5),
			want: False,
		},
		"Int64 5u64 == -5i64": {
			a:    UInt64(5),
			b:    Int64(-5),
			want: False,
		},

		"UInt64 5u64 == 5u64": {
			a:    UInt64(5),
			b:    UInt64(5),
			want: True,
		},
		"UInt64 5u64 == 7u64": {
			a:    UInt64(5),
			b:    UInt64(7),
			want: False,
		},

		"Int32 5u64 == 5i32": {
			a:    UInt64(5),
			b:    Int32(5),
			want: True,
		},
		"Int32 5u64 == 7i32": {
			a:    UInt64(5),
			b:    Int32(7),
			want: False,
		},
		"Int32 5u64 == -5i32": {
			a:    UInt64(5),
			b:    Int32(-5),
			want: False,
		},

		"UInt32 5u64 == 5u32": {
			a:    UInt64(5),
			b:    UInt32(5),
			want: True,
		},
		"UInt32 4u64 == 5u32": {
			a:    UInt64(4),
			b:    UInt32(5),
			want: False,
		},

		"Int16 5u64 == 5i16": {
			a:    UInt64(5),
			b:    Int16(5),
			want: True,
		},
		"Int16 5u64 == -5i16": {
			a:    UInt64(5),
			b:    Int16(-5),
			want: False,
		},
		"Int16 4u64 == 5i16": {
			a:    UInt64(4),
			b:    Int16(5),
			want: False,
		},

		"UInt16 5u64 == 5u16": {
			a:    UInt64(5),
			b:    UInt16(5),
			want: True,
		},
		"UInt16 4u64 == 5u16": {
			a:    UInt64(4),
			b:    UInt16(5),
			want: False,
		},

		"Int8 5u64 == 5i8": {
			a:    UInt64(5),
			b:    Int8(5),
			want: True,
		},
		"Int8 4u64 == 5i8": {
			a:    UInt64(4),
			b:    Int8(5),
			want: False,
		},
		"Int8 5u64 == -5i8": {
			a:    UInt64(5),
			b:    Int8(-5),
			want: False,
		},

		"UInt8 5u64 == 5u8": {
			a:    UInt64(5),
			b:    UInt8(5),
			want: True,
		},
		"UInt8 4u64 == 5u8": {
			a:    UInt64(4),
			b:    UInt8(5),
			want: False,
		},

		"Float64 21u64 == 21.0f64": {
			a:    UInt64(21),
			b:    Float64(21),
			want: True,
		},
		"Float64 21u64 == 21.5f64": {
			a:    UInt64(21),
			b:    Float64(21.5),
			want: False,
		},
		"Float64 21u64 == 38.0f64": {
			a:    UInt64(21),
			b:    Float64(38),
			want: False,
		},
		"Float64 0u64 == NaN": {
			a:    UInt64(0),
			b:    Float64NaN(),
			want: False,
		},
		"Float64 8u64 == +Inf": {
			a:    UInt64(8),
			b:    Float64Inf(),
			want: False,
		},
		"Float64 8u64 == -Inf": {
			a:    UInt64(8),
			b:    Float64NegInf(),
			want: False,
		},

		"Float32 21u64 == 21.0f32": {
			a:    UInt64(21),
			b:    Float32(21),
			want: True,
		},
		"Float32 21u64 == 21.5f32": {
			a:    UInt64(21),
			b:    Float32(21.5),
			want: False,
		},
		"Float32 21u64 == 38.0f32": {
			a:    UInt64(21),
			b:    Float32(38),
			want: False,
		},
		"Float32 0u64 == NaN": {
			a:    UInt64(0),
			b:    Float32NaN(),
			want: False,
		},
		"Float32 8u64 == +Inf": {
			a:    UInt64(8),
			b:    Float32Inf(),
			want: False,
		},
		"Float32 8u64 == -Inf": {
			a:    UInt64(8),
			b:    Float32NegInf(),
			want: False,
		},

		"SmallInt 16u64 == 16": {
			a:    UInt64(16),
			b:    SmallInt(16),
			want: True,
		},
		"SmallInt 97u64 == -97": {
			a:    UInt64(97),
			b:    SmallInt(-97),
			want: False,
		},

		"BigInt 16u64 == 16": {
			a:    UInt64(16),
			b:    NewBigInt(16),
			want: True,
		},
		"BigInt 97u64 == -97": {
			a:    UInt64(97),
			b:    NewBigInt(-97),
			want: False,
		},

		"Float 21u64 == 21.0": {
			a:    UInt64(21),
			b:    Float(21),
			want: True,
		},
		"Float 21u64 == 21.5": {
			a:    UInt64(21),
			b:    Float(21.5),
			want: False,
		},
		"Float 21u64 == 38.0": {
			a:    UInt64(21),
			b:    Float(38),
			want: False,
		},
		"Float 0u64 == NaN": {
			a:    UInt64(0),
			b:    FloatNaN(),
			want: False,
		},
		"Float 8u64 == +Inf": {
			a:    UInt64(8),
			b:    FloatInf(),
			want: False,
		},
		"Float 8u64 == -Inf": {
			a:    UInt64(8),
			b:    FloatNegInf(),
			want: False,
		},

		"BigFloat 21u64 == 21.0bf": {
			a:    UInt64(21),
			b:    NewBigFloat(21),
			want: True,
		},
		"BigFloat 21u64 == 21.5bf": {
			a:    UInt64(21),
			b:    NewBigFloat(21.5),
			want: False,
		},
		"BigFloat 21u64 == 38.0bf": {
			a:    UInt64(21),
			b:    NewBigFloat(38),
			want: False,
		},
		"BigFloat 0u64 == NaN": {
			a:    UInt64(0),
			b:    BigFloatNaN(),
			want: False,
		},
		"BigFloat 8u64 == +Inf": {
			a:    UInt64(8),
			b:    BigFloatInf(),
			want: False,
		},
		"BigFloat 8u64 == -Inf": {
			a:    UInt64(8),
			b:    BigFloatNegInf(),
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := StrictUnsignedIntEqual(tc.a, tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				Float32Comparer,
				Float64Comparer,
				BigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestStrictInt_Divide(t *testing.T) {
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
			err: NewError(ZeroDivisionErrorClass, "can't divide by zero"),
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

func TestStrictFloat_Divide(t *testing.T) {
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

func TestStrictNumeric_ParseUint(t *testing.T) {
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

func TestStrictInt_RightBitshift(t *testing.T) {
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

func TestStrictInt_LogicalRightBitshift(t *testing.T) {
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
			got, err := StrictIntLogicalRightBitshift(tc.a, tc.b, logicalRightShift64)
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

func TestStrictInt_LeftBitshift(t *testing.T) {
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

func TestStrictInt_BitwiseAnd(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"perform AND for String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"perform AND for Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"11 & 7": {
			a:    Int64(0b1011),
			b:    Int64(0b0111),
			want: Int64(0b0011),
		},
		"-14 & 23": {
			a:    Int64(-14),
			b:    Int64(23),
			want: Int64(18),
		},
		"258 & 0": {
			a:    Int64(258),
			b:    Int64(0),
			want: Int64(0),
		},
		"124 & 255": {
			a:    Int64(124),
			b:    Int64(255),
			want: Int64(124),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntBitwiseAnd(tc.a, tc.b)
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

func TestStrictInt_BitwiseOr(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"perform OR for String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"perform OR for Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"21 | 13": {
			a:    Int64(0b10101),
			b:    Int64(0b01101),
			want: Int64(0b11101),
		},
		"-14 | 23": {
			a:    Int64(-14),
			b:    Int64(23),
			want: Int64(-9),
		},
		"258 | 0": {
			a:    Int64(258),
			b:    Int64(0),
			want: Int64(258),
		},
		"124 | 255": {
			a:    Int64(124),
			b:    Int64(255),
			want: Int64(255),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntBitwiseOr(tc.a, tc.b)
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

func TestStrictInt_BitwiseXor(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"perform XOR for String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"perform XOR for Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"21 ^ 13": {
			a:    Int64(0b10101),
			b:    Int64(0b01101),
			want: Int64(0b11000),
		},
		"-14 ^ 23": {
			a:    Int64(-14),
			b:    Int64(23),
			want: Int64(-27),
		},
		"258 ^ 0": {
			a:    Int64(258),
			b:    Int64(0),
			want: Int64(258),
		},
		"124 ^ 255": {
			a:    Int64(124),
			b:    Int64(255),
			want: Int64(131),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntBitwiseXor(tc.a, tc.b)
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

func TestStrictInt_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    Int64
		b    Value
		want Int64
		err  *Error
	}{
		"perform modulo for String and return an error": {
			a:   Int64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Int64`"),
		},
		"perform modulo for Int32 and return an error": {
			a:   Int64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Int64`"),
		},
		"21 % 10": {
			a:    Int64(21),
			b:    Int64(10),
			want: Int64(1),
		},
		"38 % 3": {
			a:    Int64(38),
			b:    Int64(3),
			want: Int64(2),
		},
		"522 % 39": {
			a:    Int64(522),
			b:    Int64(39),
			want: Int64(15),
		},
		"38 % 0": {
			a:   Int64(38),
			b:   Int64(0),
			err: NewError(ZeroDivisionErrorClass, "can't divide by zero"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictIntModulo(tc.a, tc.b)
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

func TestStrictFloat_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    Float64
		b    Value
		want Float64
		err  *Error
	}{
		"perform modulo for String and return an error": {
			a:   Float64(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::Float64`"),
		},
		"perform modulo for Int32 and return an error": {
			a:   Float64(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::Float64`"),
		},
		"perform modulo for Float32 and return an error": {
			a:   Float64(5),
			b:   Float32(2),
			err: NewError(TypeErrorClass, "`Std::Float32` can't be coerced into `Std::Float64`"),
		},
		"21 % 10": {
			a:    Float64(21),
			b:    Float64(10),
			want: Float64(1),
		},
		"38 % 3": {
			a:    Float64(38),
			b:    Float64(3),
			want: Float64(2),
		},
		"522 % 39": {
			a:    Float64(522),
			b:    Float64(39),
			want: Float64(15),
		},
		"56.87 % 3": {
			a:    Float64(56.87),
			b:    Float64(3),
			want: Float64(2.8699999999999974),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := StrictFloatModulo(tc.a, tc.b)
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
