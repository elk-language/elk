package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBigFloatAdd(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right Value
		want  Value
		err   *Error
	}{
		"BigFloat + BigFloat => BigFloat": {
			left:  NewBigFloat(2.5),
			right: NewBigFloat(10.2),
			want:  NewBigFloat(12.7),
		},
		"BigFloat + BigFloat NaN => BigFloat NaN": {
			left:  NewBigFloat(2.5),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN + BigFloat => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: NewBigFloat(10.2),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN + BigFloat NaN => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf + BigFloat => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatInf(),
		},
		"BigFloat + BigFloat +Inf => BigFloat +Inf": {
			left:  NewBigFloat(10.2),
			right: BigFloatInf(),
			want:  BigFloatInf(),
		},
		"BigFloat +Inf + BigFloat +Inf => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: BigFloatInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf + BigFloat => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatNegInf(),
		},
		"BigFloat + BigFloat -Inf => BigFloat -Inf": {
			left:  NewBigFloat(10.2),
			right: BigFloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf + BigFloat -Inf => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat +Inf + BigFloat -Inf => BigFloat NaN": {
			left:  BigFloatInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf + BigFloat +Inf => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},
		"result takes the max precision from its operands": {
			left:  NewBigFloat(2.5).SetPrecision(31),
			right: NewBigFloat(10.2).SetPrecision(54),
			want:  NewBigFloat(12.7).SetPrecision(54),
		},
		"result takes the max precision from its operands (left)": {
			left:  NewBigFloat(2.5).SetPrecision(54),
			right: NewBigFloat(10.2).SetPrecision(52),
			want:  NewBigFloat(12.7).SetPrecision(54),
		},
		"BigFloat + SmallInt => BigFloat": {
			left:  NewBigFloat(2.5),
			right: SmallInt(120),
			want:  NewBigFloat(122.5).SetPrecision(64),
		},
		"BigFloat + BigInt => BigFloat": {
			left:  NewBigFloat(2.5),
			right: NewBigInt(120),
			want:  NewBigFloat(122.5).SetPrecision(64),
		},
		"BigFloat + Int64 => TypeError": {
			left:  NewBigFloat(2.5),
			right: Int64(20),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat + String => TypeError": {
			left:  NewBigFloat(2.5),
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"BigFloat + Float => BigFloat": {
			left:  NewBigFloat(2.5),
			right: Float(5.2),
			want:  NewBigFloat(7.7),
		},
		"BigFloat + Float NaN => BigFloat NaN": {
			left:  NewBigFloat(2.5),
			right: FloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN + Float => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: Float(10.2),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN + Float NaN => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: FloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf + Float => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: Float(10.2),
			want:  BigFloatInf(),
		},
		"BigFloat + Float +Inf => BigFloat +Inf": {
			left:  NewBigFloat(10.2),
			right: FloatInf(),
			want:  BigFloatInf(),
		},
		"BigFloat +Inf + Float +Inf => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: FloatInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf + Float => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: Float(10.2),
			want:  BigFloatNegInf(),
		},
		"BigFloat + Float -Inf => BigFloat -Inf": {
			left:  NewBigFloat(10.2),
			right: FloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf + Float -Inf => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: FloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat +Inf + Float -Inf => BigFloat NaN": {
			left:  BigFloatInf(),
			right: FloatNegInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf + Float +Inf => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: FloatInf(),
			want:  BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Add(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				bigFloatComparer,
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

func TestCountFloatDigits(t *testing.T) {
	tests := map[string]struct {
		str  string
		want int
	}{
		"int": {
			str:  "35",
			want: 2,
		},
		"float": {
			str:  "254.671",
			want: 6,
		},
		"int with exponent": {
			str:  "257e20",
			want: 3,
		},
		"float with exponent": {
			str:  "257.1223e91",
			want: 7,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := CountFloatDigits(tc.str)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloatSubtract(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right Value
		want  Value
		err   *Error
	}{
		"BigFloat - BigFloat => BigFloat": {
			left:  NewBigFloat(10.0),
			right: NewBigFloat(2.5),
			want:  NewBigFloat(7.5),
		},
		"result takes the max precision from its operands": {
			left:  NewBigFloat(10.0).SetPrecision(54),
			right: NewBigFloat(2.5).SetPrecision(31),
			want:  NewBigFloat(7.5).SetPrecision(54),
		},
		"BigFloat - SmallInt => BigFloat": {
			left:  NewBigFloat(120.5),
			right: SmallInt(2),
			want:  NewBigFloat(118.5).SetPrecision(64),
		},
		"BigFloat - BigInt => BigFloat": {
			left:  NewBigFloat(120.5),
			right: NewBigInt(2),
			want:  NewBigFloat(118.5).SetPrecision(64),
		},
		"BigFloat - Int64 => TypeError": {
			left:  NewBigFloat(20.5),
			right: Int64(2),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat - String => TypeError": {
			left:  NewBigFloat(2.5),
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"BigFloat - BigFloat NaN => BigFloat NaN": {
			left:  NewBigFloat(2.5),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN - BigFloat => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: NewBigFloat(10.2),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN - BigFloat NaN => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf - BigFloat => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatInf(),
		},
		"BigFloat - BigFloat +Inf => BigFloat -Inf": {
			left:  NewBigFloat(10.2),
			right: BigFloatInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat +Inf - BigFloat +Inf => BigFloat NaN": {
			left:  BigFloatInf(),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf - BigFloat => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatNegInf(),
		},
		"BigFloat - BigFloat -Inf => BigFloat +Inf": {
			left:  NewBigFloat(10.2),
			right: BigFloatNegInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf - BigFloat -Inf => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf - BigFloat -Inf => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: BigFloatNegInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf - BigFloat +Inf => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: BigFloatInf(),
			want:  BigFloatNegInf(),
		},

		"BigFloat - Float NaN => BigFloat NaN": {
			left:  NewBigFloat(2.5),
			right: FloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN - Float => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: Float(10.2),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN - Float NaN => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: FloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf - Float => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: Float(10.2),
			want:  BigFloatInf(),
		},
		"BigFloat - Float +Inf => BigFloat -Inf": {
			left:  NewBigFloat(10.2),
			right: FloatInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat +Inf - Float +Inf => BigFloat NaN": {
			left:  BigFloatInf(),
			right: FloatInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf - Float => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: Float(10.2),
			want:  BigFloatNegInf(),
		},
		"BigFloat - Float -Inf => BigFloat +Inf": {
			left:  NewBigFloat(10.2),
			right: FloatNegInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf - Float -Inf => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: FloatNegInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf - Float -Inf => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: FloatNegInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf - Float +Inf => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: FloatInf(),
			want:  BigFloatNegInf(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Subtract(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				bigFloatComparer,
				floatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf("want: %s, got: %s\n%s", tc.want.Inspect(), got.Inspect(), diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloatMultiply(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right Value
		want  Value
		err   *Error
	}{
		"BigFloat * BigFloat => BigFloat": {
			left:  NewBigFloat(2.55),
			right: NewBigFloat(10.0),
			want:  NewBigFloat(25.5),
		},
		"result takes the max precision from its operands": {
			left:  NewBigFloat(2.5).SetPrecision(31),
			right: NewBigFloat(10.0).SetPrecision(54),
			want:  NewBigFloat(25.0).SetPrecision(54),
		},
		"BigFloat * SmallInt => BigFloat": {
			left:  NewBigFloat(2.5),
			right: SmallInt(10),
			want:  NewBigFloat(25.0).SetPrecision(64),
		},
		"BigFloat * BigInt => BigFloat": {
			left:  NewBigFloat(2.5),
			right: NewBigInt(10),
			want:  NewBigFloat(25.0).SetPrecision(64),
		},
		"BigFloat * Int64 => TypeError": {
			left:  NewBigFloat(2.55),
			right: Int64(20),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat * String => TypeError": {
			left:  NewBigFloat(2.5),
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"BigFloat * BigFloat NaN => BigFloat NaN": {
			left:  NewBigFloat(2.5),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN * BigFloat => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: NewBigFloat(10.2),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN * BigFloat NaN => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf * BigFloat => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatInf(),
		},
		"BigFloat * BigFloat +Inf => BigFloat +Inf": {
			left:  NewBigFloat(10.2),
			right: BigFloatInf(),
			want:  BigFloatInf(),
		},
		"BigFloat +Inf * BigFloat +Inf => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: BigFloatInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf * +BigFloat => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf * -BigFloat => BigFloat +Inf": {
			left:  BigFloatNegInf(),
			right: NewBigFloat(-10.2),
			want:  BigFloatInf(),
		},
		"+BigFloat * BigFloat -Inf => BigFloat -Inf": {
			left:  NewBigFloat(10.2),
			right: BigFloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"-BigFloat * BigFloat -Inf => BigFloat +Inf": {
			left:  NewBigFloat(-10.2),
			right: BigFloatNegInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf * BigFloat -Inf => BigFloat +Inf": {
			left:  BigFloatNegInf(),
			right: BigFloatNegInf(),
			want:  BigFloatInf(),
		},
		"BigFloat +Inf * BigFloat -Inf => BigFloat -Inf": {
			left:  BigFloatInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf * BigFloat +Inf => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: BigFloatInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf * BigFloat 0 => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: NewBigFloat(0),
			want:  BigFloatNaN(),
		},
		"BigFloat 0 * BigFloat +Inf => BigFloat NaN": {
			left:  NewBigFloat(0),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},

		"BigFloat * Float NaN => BigFloat NaN": {
			left:  NewBigFloat(2.5),
			right: FloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN * Float => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: Float(10.2),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN * Float NaN => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: FloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf * Float => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: Float(10.2),
			want:  BigFloatInf(),
		},
		"BigFloat * Float +Inf => BigFloat +Inf": {
			left:  NewBigFloat(10.2),
			right: FloatInf(),
			want:  BigFloatInf(),
		},
		"BigFloat +Inf * Float +Inf => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: FloatInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf * +Float => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: Float(10.2),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf * -Float => BigFloat +Inf": {
			left:  BigFloatNegInf(),
			right: Float(-10.2),
			want:  BigFloatInf(),
		},
		"+BigFloat * Float -Inf => BigFloat -Inf": {
			left:  NewBigFloat(10.2),
			right: FloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"-BigFloat * Float -Inf => BigFloat +Inf": {
			left:  NewBigFloat(-10.2),
			right: FloatNegInf(),
			want:  BigFloatInf(),
		},
		"BigFloat -Inf * Float -Inf => BigFloat +Inf": {
			left:  BigFloatNegInf(),
			right: FloatNegInf(),
			want:  BigFloatInf(),
		},
		"BigFloat +Inf * Float -Inf => BigFloat -Inf": {
			left:  BigFloatInf(),
			right: FloatNegInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf * Float +Inf => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: FloatInf(),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf * Float 0 => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: Float(0),
			want:  BigFloatNaN(),
		},
		"BigFloat 0 * Float +Inf => BigFloat NaN": {
			left:  NewBigFloat(0),
			right: FloatInf(),
			want:  BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Multiply(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				bigFloatComparer,
				floatComparer,
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

func TestBigFloatDivide(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right Value
		want  Value
		err   *Error
	}{
		"BigFloat / BigFloat => BigFloat": {
			left:  NewBigFloat(2.68),
			right: NewBigFloat(2),
			want:  NewBigFloat(1.34),
		},
		"result takes the max precision from its operands": {
			left:  NewBigFloat(2).SetPrecision(31),
			right: NewBigFloat(2).SetPrecision(54),
			want:  NewBigFloat(1).SetPrecision(54),
		},
		"BigFloat / SmallInt => BigFloat": {
			left:  NewBigFloat(2.68),
			right: SmallInt(2),
			want:  NewBigFloat(1.34).SetPrecision(64),
		},
		"BigFloat / BigInt => BigFloat": {
			left:  NewBigFloat(2.68),
			right: NewBigInt(2),
			want:  NewBigFloat(1.34).SetPrecision(64),
		},
		"BigFloat / Int64 => TypeError": {
			left:  NewBigFloat(2.68),
			right: Int64(2),
			err:   NewError(TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat / String => TypeError": {
			left:  NewBigFloat(2.5),
			right: String("foo"),
			err:   NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"BigFloat / BigFloat NaN => BigFloat NaN": {
			left:  NewBigFloat(2.5),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN / BigFloat => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: NewBigFloat(10.2),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN / BigFloat NaN => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf / BigFloat => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatInf(),
		},
		"BigFloat / BigFloat +Inf => BigFloat 0": {
			left:  NewBigFloat(10.2),
			right: BigFloatInf(),
			want:  NewBigFloat(0),
		},
		"BigFloat +Inf / BigFloat +Inf => BigFloat NaN": {
			left:  BigFloatInf(),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf / +BigFloat => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: NewBigFloat(10.2),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf / -BigFloat => BigFloat +Inf": {
			left:  BigFloatNegInf(),
			right: NewBigFloat(-10.2),
			want:  BigFloatInf(),
		},
		"+BigFloat / BigFloat -Inf => BigFloat -0": {
			left:  NewBigFloat(10.2),
			right: BigFloatNegInf(),
			want:  NewBigFloat(-0),
		},
		"-BigFloat / BigFloat -Inf => BigFloat +0": {
			left:  NewBigFloat(-10.2),
			right: BigFloatNegInf(),
			want:  NewBigFloat(0),
		},
		"BigFloat -Inf / BigFloat -Inf => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf / BigFloat -Inf => BigFloat NaN": {
			left:  BigFloatInf(),
			right: BigFloatNegInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf / BigFloat +Inf => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: BigFloatInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf / BigFloat 0 => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: NewBigFloat(0),
			want:  BigFloatNegInf(),
		},
		"BigFloat 0 / BigFloat +Inf => BigFloat 0": {
			left:  NewBigFloat(0),
			right: BigFloatInf(),
			want:  NewBigFloat(0),
		},

		"BigFloat / Float NaN => BigFloat NaN": {
			left:  NewBigFloat(2.5),
			right: FloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN / Float => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: Float(10.2),
			want:  BigFloatNaN(),
		},
		"BigFloat NaN / Float NaN => BigFloat NaN": {
			left:  BigFloatNaN(),
			right: FloatNaN(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf / Float => BigFloat +Inf": {
			left:  BigFloatInf(),
			right: Float(10.2),
			want:  BigFloatInf(),
		},
		"BigFloat / Float +Inf => BigFloat 0": {
			left:  NewBigFloat(10.2),
			right: FloatInf(),
			want:  NewBigFloat(0),
		},
		"BigFloat +Inf / Float +Inf => BigFloat NaN": {
			left:  BigFloatInf(),
			right: FloatInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf / +Float => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: Float(10.2),
			want:  BigFloatNegInf(),
		},
		"BigFloat -Inf / -Float => BigFloat +Inf": {
			left:  BigFloatNegInf(),
			right: Float(-10.2),
			want:  BigFloatInf(),
		},
		"+BigFloat / Float -Inf => BigFloat -0": {
			left:  NewBigFloat(10.2),
			right: FloatNegInf(),
			want:  NewBigFloat(-0),
		},
		"-BigFloat / Float -Inf => BigFloat +0": {
			left:  NewBigFloat(-10.2),
			right: FloatNegInf(),
			want:  NewBigFloat(0),
		},
		"BigFloat -Inf / Float -Inf => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: FloatNegInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat +Inf / Float -Inf => BigFloat NaN": {
			left:  BigFloatInf(),
			right: FloatNegInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf / Float +Inf => BigFloat NaN": {
			left:  BigFloatNegInf(),
			right: FloatInf(),
			want:  BigFloatNaN(),
		},
		"BigFloat -Inf / Float 0 => BigFloat -Inf": {
			left:  BigFloatNegInf(),
			right: Float(0),
			want:  BigFloatNegInf(),
		},
		"BigFloat 0 / Float +Inf => BigFloat 0": {
			left:  NewBigFloat(0),
			right: FloatInf(),
			want:  NewBigFloat(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Divide(tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				bigFloatComparer,
				floatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Log(got.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloat_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    *BigFloat
		b    Value
		want Value
		err  *Error
	}{
		"exponentiate String and return an error": {
			a:   NewBigFloat(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
		"exponentiate Int32 and return an error": {
			a:   NewBigFloat(5),
			b:   Int32(2),
			err: NewError(TypeErrorClass, "`Std::Int32` can't be coerced into `Std::BigFloat`"),
		},
		"SmallInt 5 ** 2": {
			a:    NewBigFloat(5),
			b:    SmallInt(2),
			want: NewBigFloat(25).SetPrecision(64),
		},
		"SmallInt 5p92 ** 2": {
			a:    NewBigFloat(5).SetPrecision(92),
			b:    SmallInt(2),
			want: NewBigFloat(25).SetPrecision(92),
		},
		"SmallInt 7 ** 8": {
			a:    NewBigFloat(7),
			b:    SmallInt(8),
			want: NewBigFloat(5764801).SetPrecision(64),
		},
		"SmallInt 2.5 ** 5": {
			a:    NewBigFloat(2.5),
			b:    SmallInt(5),
			want: NewBigFloat(97.65625).SetPrecision(64),
		},
		"SmallInt 7.12 ** 1": {
			a:    NewBigFloat(7.12),
			b:    SmallInt(1),
			want: NewBigFloat(7.12).SetPrecision(64),
		},
		"SmallInt 4 ** -2": {
			a:    NewBigFloat(4),
			b:    SmallInt(-2),
			want: NewBigFloat(0.0625).SetPrecision(64),
		},
		"SmallInt 25 ** 0": {
			a:    NewBigFloat(25),
			b:    SmallInt(0),
			want: NewBigFloat(1).SetPrecision(64),
		},

		"BigInt 5 ** 2": {
			a:    NewBigFloat(5),
			b:    NewBigInt(2),
			want: NewBigFloat(25).SetPrecision(64),
		},
		"BigInt 5p78 ** 2": {
			a:    NewBigFloat(5).SetPrecision(78),
			b:    NewBigInt(2),
			want: NewBigFloat(25).SetPrecision(78),
		},
		"BigInt 7 ** 8": {
			a:    NewBigFloat(7),
			b:    NewBigInt(8),
			want: NewBigFloat(5764801).SetPrecision(64),
		},
		"BigInt 2.5 ** 5": {
			a:    NewBigFloat(2.5),
			b:    NewBigInt(5),
			want: NewBigFloat(97.65625).SetPrecision(64),
		},
		"BigInt 7.12 ** 1": {
			a:    NewBigFloat(7.12),
			b:    NewBigInt(1),
			want: NewBigFloat(7.12).SetPrecision(64),
		},
		"BigInt 4 ** -2": {
			a:    NewBigFloat(4),
			b:    NewBigInt(-2),
			want: NewBigFloat(0.0625).SetPrecision(64),
		},
		"BigInt 25 ** 0": {
			a:    NewBigFloat(25),
			b:    NewBigInt(0),
			want: NewBigFloat(1).SetPrecision(64),
		},

		"Float 5 ** 2": {
			a:    NewBigFloat(5),
			b:    Float(2),
			want: NewBigFloat(25),
		},
		"Float 5p83 ** 2": {
			a:    NewBigFloat(5).SetPrecision(83),
			b:    Float(2),
			want: NewBigFloat(25).SetPrecision(83),
		},
		"Float 7 ** 8": {
			a:    NewBigFloat(7),
			b:    Float(8),
			want: NewBigFloat(5764801),
		},
		"Float 2.5 ** 2.5": {
			a:    NewBigFloat(2.5),
			b:    Float(2.5),
			want: NewBigFloat(9.882117688026186),
		},
		"Float 3 ** 2.5": {
			a:    NewBigFloat(3),
			b:    Float(2.5),
			want: NewBigFloat(15.588457268119896),
		},
		"Float 6 ** 1": {
			a:    NewBigFloat(6),
			b:    Float(1),
			want: NewBigFloat(6),
		},
		"Float 4 ** -2": {
			a:    NewBigFloat(4),
			b:    Float(-2),
			want: NewBigFloat(0.0625),
		},
		"Float 25 ** 0": {
			a:    NewBigFloat(25),
			b:    Float(0),
			want: NewBigFloat(1),
		},
		"Float 25 ** NaN": {
			a:    NewBigFloat(25),
			b:    FloatNaN(),
			want: BigFloatNaN(),
		},
		"Float NaN ** 25": {
			a:    BigFloatNaN(),
			b:    Float(25),
			want: BigFloatNaN(),
		},
		"Float NaN ** NaN": {
			a:    BigFloatNaN(),
			b:    FloatNaN(),
			want: BigFloatNaN(),
		},
		"Float 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    NewBigFloat(0),
			b:    Float(-5),
			want: BigFloatInf(),
		},
		"Float 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    NewBigFloat(0),
			b:    FloatNegInf(),
			want: BigFloatInf(),
		},
		"Float 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    NewBigFloat(0),
			b:    FloatInf(),
			want: NewBigFloat(0),
		},
		"Float 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    NewBigFloat(0),
			b:    Float(-8),
			want: BigFloatInf(),
		},
		"Float 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    NewBigFloat(0),
			b:    Float(7),
			want: NewBigFloat(0),
		},
		"Float 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    NewBigFloat(0),
			b:    Float(8),
			want: NewBigFloat(0),
		},
		"Float -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    NewBigFloat(-1),
			b:    FloatInf(),
			want: NewBigFloat(1),
		},
		"Float -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    NewBigFloat(-1),
			b:    FloatNegInf(),
			want: NewBigFloat(1),
		},
		"Float 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    NewBigFloat(2),
			b:    FloatInf(),
			want: BigFloatInf(),
		},
		"Float -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    NewBigFloat(-2),
			b:    FloatInf(),
			want: BigFloatInf(),
		},
		"Float 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    NewBigFloat(2),
			b:    FloatNegInf(),
			want: NewBigFloat(0),
		},
		"Float -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    NewBigFloat(-2),
			b:    FloatNegInf(),
			want: NewBigFloat(0),
		},
		"Float 0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    NewBigFloat(0.5),
			b:    FloatInf(),
			want: NewBigFloat(0),
		},
		"Float -0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    NewBigFloat(-0.5),
			b:    FloatInf(),
			want: NewBigFloat(0),
		},
		"Float 0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    NewBigFloat(0.5),
			b:    FloatNegInf(),
			want: BigFloatInf(),
		},
		"Float -0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    NewBigFloat(-0.5),
			b:    FloatNegInf(),
			want: BigFloatInf(),
		},
		"Float +Inf ** 5": { // Pow(+Inf, y) = +Inf for y > 0
			a:    BigFloatInf(),
			b:    Float(5),
			want: BigFloatInf(),
		},
		"Float +Inf ** -7": { // Pow(+Inf, y) = +0 for y < 0
			a:    BigFloatInf(),
			b:    Float(-7),
			want: NewBigFloat(0),
		},
		"Float -Inf ** -7": {
			a:    BigFloatNegInf(),
			b:    Float(-7),
			want: NewBigFloat(0),
		},
		"Float -5.5 ** 3.8": { // Pow(x, y) = NaN for finite x < 0 and finite non-integer y
			a:    NewBigFloat(-5.5),
			b:    Float(3.8),
			want: BigFloatNaN(),
		},

		"BigFloat 5 ** 2": {
			a:    NewBigFloat(5),
			b:    NewBigFloat(2),
			want: NewBigFloat(25).SetPrecision(53),
		},
		"BigFloat 7 ** 8": {
			a:    NewBigFloat(7),
			b:    NewBigFloat(8),
			want: NewBigFloat(5764801).SetPrecision(53),
		},
		"BigFloat 2.5 ** 2.5": {
			a:    NewBigFloat(2.5),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("9.882117688026186").SetPrecision(53),
		},
		"BigFloat 3 ** 2.5": {
			a:    NewBigFloat(3),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("15.588457268119896").SetPrecision(53),
		},
		"BigFloat 6 ** 1": {
			a:    NewBigFloat(6),
			b:    NewBigFloat(1),
			want: NewBigFloat(6).SetPrecision(53),
		},
		"BigFloat 4 ** -2": {
			a:    NewBigFloat(4),
			b:    NewBigFloat(-2),
			want: NewBigFloat(0.0625).SetPrecision(53),
		},
		"BigFloat 25 ** 0": {
			a:    NewBigFloat(25),
			b:    NewBigFloat(0),
			want: NewBigFloat(1).SetPrecision(53),
		},
		"BigFloat 25 ** NaN": {
			a:    NewBigFloat(25),
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"BigFloat NaN ** 25": {
			a:    BigFloatNaN(),
			b:    NewBigFloat(25),
			want: BigFloatNaN(),
		},
		"BigFloat NaN ** NaN": {
			a:    BigFloatNaN(),
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
		"BigFloat 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    NewBigFloat(0),
			b:    NewBigFloat(-5),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    NewBigFloat(0),
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    NewBigFloat(0),
			b:    BigFloatInf(),
			want: NewBigFloat(0),
		},
		"BigFloat 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    NewBigFloat(0),
			b:    NewBigFloat(-8),
			want: BigFloatInf(),
		},
		"BigFloat 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    NewBigFloat(0),
			b:    NewBigFloat(7),
			want: NewBigFloat(0),
		},
		"BigFloat 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    NewBigFloat(0),
			b:    NewBigFloat(8),
			want: NewBigFloat(0),
		},
		"BigFloat -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    NewBigFloat(-1),
			b:    BigFloatInf(),
			want: NewBigFloat(1),
		},
		"BigFloat -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    NewBigFloat(-1),
			b:    BigFloatNegInf(),
			want: NewBigFloat(1),
		},
		"BigFloat 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    NewBigFloat(2),
			b:    BigFloatInf(),
			want: BigFloatInf(),
		},
		"BigFloat -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    NewBigFloat(-2),
			b:    BigFloatInf(),
			want: BigFloatInf(),
		},
		"BigFloat 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    NewBigFloat(2),
			b:    BigFloatNegInf(),
			want: NewBigFloat(0),
		},
		"BigFloat -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    NewBigFloat(-2),
			b:    BigFloatNegInf(),
			want: NewBigFloat(0),
		},
		"BigFloat 0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    NewBigFloat(0.5),
			b:    BigFloatInf(),
			want: NewBigFloat(0),
		},
		"BigFloat -0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    NewBigFloat(-0.5),
			b:    BigFloatInf(),
			want: NewBigFloat(0),
		},
		"BigFloat 0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    NewBigFloat(0.5),
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
		"BigFloat -0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    NewBigFloat(-0.5),
			b:    BigFloatNegInf(),
			want: BigFloatInf(),
		},
		"BigFloat +Inf ** 5": { // Pow(+Inf, y) = +Inf for y > 0
			a:    BigFloatInf(),
			b:    NewBigFloat(5),
			want: BigFloatInf(),
		},
		"BigFloat +Inf ** -7": { // Pow(+Inf, y) = +0 for y < 0
			a:    BigFloatInf(),
			b:    NewBigFloat(-7),
			want: NewBigFloat(0),
		},
		"BigFloat -Inf ** -7": {
			a:    BigFloatNegInf(),
			b:    NewBigFloat(-7),
			want: NewBigFloat(0),
		},
		"BigFloat -5.5 ** 3.8": { // Pow(x, y) = NaN for finite x < 0 and finite non-integer y
			a:    NewBigFloat(-5.5),
			b:    NewBigFloat(3.8),
			want: BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmp.AllowUnexported(Error{}, BigInt{}),
				floatComparer,
				bigFloatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloat_Mod(t *testing.T) {
	tests := map[string]struct {
		left  *BigFloat
		right *BigFloat
		want  *BigFloat
	}{
		"25 % 3": {
			left:  NewBigFloat(25),
			right: NewBigFloat(3),
			want:  NewBigFloat(1),
		},
		"76 % 6": {
			left:  NewBigFloat(76),
			right: NewBigFloat(6),
			want:  NewBigFloat(4),
		},
		"76.75 % 6.25": {
			left:  NewBigFloat(76.75),
			right: NewBigFloat(6.25),
			want:  NewBigFloat(1.75),
		},
		"76.75 % -6.25": {
			left:  NewBigFloat(76.75),
			right: NewBigFloat(-6.25),
			want:  NewBigFloat(1.75),
		},
		"-76.75 % 6.25": {
			left:  NewBigFloat(-76.75),
			right: NewBigFloat(6.25),
			want:  NewBigFloat(-1.75),
		},
		"-76.75 % -6.25": {
			left:  NewBigFloat(-76.75),
			right: NewBigFloat(-6.25),
			want:  NewBigFloat(-1.75),
		},
		"+Inf % 5": { // Mod(±Inf, y) = NaN
			left:  BigFloatInf(),
			right: NewBigFloat(5),
			want:  BigFloatNaN(),
		},
		"-Inf % 5": { // Mod(±Inf, y) = NaN
			left:  BigFloatNegInf(),
			right: NewBigFloat(5),
			want:  BigFloatNaN(),
		},
		"NaN % 625": { // Mod(NaN, y) = NaN
			left:  BigFloatNaN(),
			right: NewBigFloat(625),
			want:  BigFloatNaN(),
		},
		"25 % 0": { // Mod(x, 0) = NaN
			left:  NewBigFloat(25),
			right: NewBigFloat(0),
			want:  BigFloatNaN(),
		},
		"25 % +Inf": { // Mod(x, ±Inf) = x
			left:  NewBigFloat(25),
			right: BigFloatInf(),
			want:  NewBigFloat(25),
		},
		"-87 % -Inf": { // Mod(x, ±Inf) = x
			left:  NewBigFloat(-87),
			right: BigFloatNegInf(),
			want:  NewBigFloat(-87),
		},
		"49 % NaN": { // Mod(x, NaN) = NaN
			left:  NewBigFloat(49),
			right: BigFloatNaN(),
			want:  BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Mod(tc.left, tc.right)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				bigFloatComparer,
				floatComparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloat_FloorBigFloat(t *testing.T) {
	tests := map[string]struct {
		f    *BigFloat
		want *BigFloat
	}{
		"floor(25)": {
			f:    NewBigFloat(25),
			want: NewBigFloat(25),
		},
		"floor(38.7)": {
			f:    NewBigFloat(38.7),
			want: NewBigFloat(38),
		},
		"floor(-6.5)": {
			f:    NewBigFloat(-6.5),
			want: NewBigFloat(-7),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.f.FloorBigFloat()
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmp.AllowUnexported(BigFloat{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Log(got.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloat_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    *BigFloat
		b    Value
		want Value
		err  *Error
	}{
		"String and return an error": {
			a:   NewBigFloat(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"SmallInt 18446744073709551616 % 2": {
			a:    ParseBigFloatPanic("18446744073709551616"),
			b:    SmallInt(2),
			want: NewBigFloat(0).SetPrecision(67),
		},
		"SmallInt 25 % 3": {
			a:    NewBigFloat(25),
			b:    SmallInt(3),
			want: NewBigFloat(1).SetPrecision(64),
		},
		"SmallInt 25.6 % 3": {
			a:    NewBigFloat(25.6),
			b:    SmallInt(3),
			want: NewBigFloat(1.6000000000000014).SetPrecision(64),
		},
		"SmallInt 76 % 6": {
			a:    NewBigFloat(76),
			b:    SmallInt(6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"SmallInt -76 % 6": {
			a:    NewBigFloat(-76),
			b:    SmallInt(6),
			want: NewBigFloat(-4).SetPrecision(64),
		},
		"SmallInt 76 % -6": {
			a:    NewBigFloat(76),
			b:    SmallInt(-6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"SmallInt -76 % -6": {
			a:    NewBigFloat(-76),
			b:    SmallInt(-6),
			want: NewBigFloat(-4).SetPrecision(64),
		},
		"SmallInt 124 % 9": {
			a:    NewBigFloat(124),
			b:    SmallInt(9),
			want: NewBigFloat(7).SetPrecision(64),
		},

		"BigInt 25 % 3": {
			a:    NewBigFloat(25),
			b:    NewBigInt(3),
			want: NewBigFloat(1).SetPrecision(64),
		},
		"BigInt 76 % 6": {
			a:    NewBigFloat(76),
			b:    NewBigInt(6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"BigInt 76.5 % 6": {
			a:    NewBigFloat(76.5),
			b:    NewBigInt(6),
			want: NewBigFloat(4.5).SetPrecision(64),
		},
		"BigInt -76 % 6": {
			a:    NewBigFloat(-76),
			b:    NewBigInt(6),
			want: NewBigFloat(-4).SetPrecision(64),
		},
		"BigInt 76 % -6": {
			a:    NewBigFloat(76),
			b:    NewBigInt(-6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"BigInt -76 % -6": {
			a:    NewBigFloat(-76),
			b:    NewBigInt(-6),
			want: NewBigFloat(-4).SetPrecision(64),
		},
		"BigInt 124 % 9": {
			a:    NewBigFloat(124),
			b:    NewBigInt(9),
			want: NewBigFloat(7).SetPrecision(64),
		},
		"BigInt 9765 % 9223372036854775808": {
			a:    NewBigFloat(9765),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: NewBigFloat(9765).SetPrecision(64),
		},

		"Float 25 % 3": {
			a:    NewBigFloat(25),
			b:    Float(3),
			want: NewBigFloat(1),
		},
		"Float 25p102 % 3": {
			a:    NewBigFloat(25).SetPrecision(102),
			b:    Float(3),
			want: NewBigFloat(1).SetPrecision(102),
		},
		"Float 76 % 6": {
			a:    NewBigFloat(76),
			b:    Float(6),
			want: NewBigFloat(4),
		},
		"Float 124 % 9": {
			a:    NewBigFloat(124),
			b:    Float(9),
			want: NewBigFloat(7),
		},
		"Float 74.5 % 6.25": {
			a:    NewBigFloat(74.5),
			b:    Float(6.25),
			want: NewBigFloat(5.75),
		},
		"Float 74 % 6.25": {
			a:    NewBigFloat(74),
			b:    Float(6.25),
			want: NewBigFloat(5.25),
		},
		"Float -74 % 6.25": {
			a:    NewBigFloat(-74),
			b:    Float(6.25),
			want: NewBigFloat(-5.25),
		},
		"Float 74 % -6.25": {
			a:    NewBigFloat(74),
			b:    Float(-6.25),
			want: NewBigFloat(5.25),
		},
		"Float -74 % -6.25": {
			a:    NewBigFloat(-74),
			b:    Float(-6.25),
			want: NewBigFloat(-5.25),
		},
		"Float +Inf % 5": { // Mod(±Inf, y) = NaN
			a:    BigFloatInf(),
			b:    Float(5),
			want: BigFloatNaN(),
		},
		"Float -Inf % 5": { // Mod(±Inf, y) = NaN
			a:    BigFloatNegInf(),
			b:    Float(5),
			want: BigFloatNaN(),
		},
		"Float NaN % 625": { // Mod(NaN, y) = NaN
			a:    BigFloatNaN(),
			b:    Float(625),
			want: BigFloatNaN(),
		},
		"Float 25 % 0": { // Mod(x, 0) = NaN
			a:    NewBigFloat(25),
			b:    Float(0),
			want: BigFloatNaN(),
		},
		"Float 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    NewBigFloat(25),
			b:    FloatInf(),
			want: NewBigFloat(25),
		},
		"Float -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    NewBigFloat(-87),
			b:    FloatNegInf(),
			want: NewBigFloat(-87),
		},
		"Float 49 % NaN": { // Mod(x, NaN) = NaN
			a:    NewBigFloat(49),
			b:    FloatNaN(),
			want: BigFloatNaN(),
		},

		"BigFloat 25 % 3": {
			a:    NewBigFloat(25),
			b:    NewBigFloat(3),
			want: NewBigFloat(1),
		},
		"BigFloat 76 % 6": {
			a:    NewBigFloat(76),
			b:    NewBigFloat(6),
			want: NewBigFloat(4),
		},
		"BigFloat 76p82 % 6": {
			a:    NewBigFloat(76).SetPrecision(82),
			b:    NewBigFloat(6),
			want: NewBigFloat(4).SetPrecision(82),
		},
		"BigFloat 76p82 % 6p96": {
			a:    NewBigFloat(76).SetPrecision(82),
			b:    NewBigFloat(6).SetPrecision(96),
			want: NewBigFloat(4).SetPrecision(96),
		},
		"BigFloat 124 % 9": {
			a:    NewBigFloat(124),
			b:    NewBigFloat(9),
			want: NewBigFloat(7),
		},
		"BigFloat 74 % 6.25": {
			a:    NewBigFloat(74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(5.25),
		},
		"BigFloat 74 % 6.25 with higher precision": {
			a:    NewBigFloat(74),
			b:    NewBigFloat(6.25).SetPrecision(64),
			want: NewBigFloat(5.25).SetPrecision(64),
		},
		"BigFloat -74 % 6.25": {
			a:    NewBigFloat(-74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(-5.25),
		},
		"BigFloat 74 % -6.25": {
			a:    NewBigFloat(74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(5.25),
		},
		"BigFloat -74 % -6.25": {
			a:    NewBigFloat(-74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(-5.25),
		},
		"BigFloat +Inf % 5": { // Mod(±Inf, y) = NaN
			a:    BigFloatInf(),
			b:    NewBigFloat(5),
			want: BigFloatNaN(),
		},
		"BigFloat -Inf % 5": { // Mod(±Inf, y) = NaN
			a:    BigFloatNegInf(),
			b:    NewBigFloat(5),
			want: BigFloatNaN(),
		},
		"BigFloat NaN % 625": { // Mod(NaN, y) = NaN
			a:    BigFloatNaN(),
			b:    NewBigFloat(625),
			want: BigFloatNaN(),
		},
		"BigFloat 25 % 0": { // Mod(x, 0) = NaN
			a:    NewBigFloat(25),
			b:    NewBigFloat(0),
			want: BigFloatNaN(),
		},
		"BigFloat 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    NewBigFloat(25),
			b:    BigFloatInf(),
			want: NewBigFloat(25),
		},
		"BigFloat -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    NewBigFloat(-87),
			b:    BigFloatNegInf(),
			want: NewBigFloat(-87),
		},
		"BigFloat 49 % NaN": { // Mod(x, NaN) = NaN
			a:    NewBigFloat(49),
			b:    BigFloatNaN(),
			want: BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Modulo(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				bigFloatComparer,
				cmp.AllowUnexported(Error{}, BigInt{}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
