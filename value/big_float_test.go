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
		"exponentiate positive SmallInt 5 ** 2": {
			a:    NewBigFloat(5),
			b:    SmallInt(2),
			want: NewBigFloat(25),
		},
		"exponentiate positive SmallInt 7 ** 8": {
			a:    NewBigFloat(7),
			b:    SmallInt(8),
			want: NewBigFloat(5764801),
		},
		"exponentiate positive SmallInt 2.5 ** 5": {
			a:    NewBigFloat(2.5),
			b:    SmallInt(5),
			want: NewBigFloat(97.65625),
		},
		"exponentiate positive SmallInt 7.12 ** 1": {
			a:    NewBigFloat(7.12),
			b:    SmallInt(1),
			want: NewBigFloat(7.12),
		},
		"exponentiate negative SmallInt": {
			a:    NewBigFloat(4),
			b:    SmallInt(-2),
			want: NewBigFloat(0.0625),
		},
		"exponentiate SmallInt zero": {
			a:    NewBigFloat(25),
			b:    SmallInt(0),
			want: NewBigFloat(1),
		},

		"exponentiate positive BigInt 5 ** 2": {
			a:    NewBigFloat(5),
			b:    NewBigInt(2),
			want: NewBigFloat(25),
		},
		"exponentiate positive BigInt 7 ** 8": {
			a:    NewBigFloat(7),
			b:    NewBigInt(8),
			want: NewBigFloat(5764801),
		},
		"exponentiate positive BigInt 2.5 ** 5": {
			a:    NewBigFloat(2.5),
			b:    NewBigInt(5),
			want: NewBigFloat(97.65625),
		},
		"exponentiate positive BigInt 7.12 ** 1": {
			a:    NewBigFloat(7.12),
			b:    NewBigInt(1),
			want: NewBigFloat(7.12),
		},
		"exponentiate negative BigInt": {
			a:    NewBigFloat(4),
			b:    NewBigInt(-2),
			want: NewBigFloat(0.0625),
		},
		"exponentiate BigInt zero": {
			a:    NewBigFloat(25),
			b:    NewBigInt(0),
			want: NewBigFloat(1),
		},

		"exponentiate positive Float 5 ** 2": {
			a:    NewBigFloat(5),
			b:    Float(2),
			want: NewBigFloat(25),
		},
		"exponentiate positive Float 7 ** 8": {
			a:    NewBigFloat(7),
			b:    Float(8),
			want: NewBigFloat(5764801),
		},
		"exponentiate positive Float 2.5 ** 2.5": {
			a:    NewBigFloat(2.5),
			b:    Float(2.5),
			want: NewBigFloat(9.882117688026186),
		},
		"exponentiate positive Float 3 ** 2.5": {
			a:    NewBigFloat(3),
			b:    Float(2.5),
			want: NewBigFloat(15.588457268119896),
		},
		"exponentiate positive Float 6 ** 1": {
			a:    NewBigFloat(6),
			b:    Float(1),
			want: NewBigFloat(6),
		},
		"exponentiate negative Float": {
			a:    NewBigFloat(4),
			b:    Float(-2),
			want: NewBigFloat(0.0625),
		},
		"exponentiate Float zero": {
			a:    NewBigFloat(25),
			b:    Float(0),
			want: NewBigFloat(1),
		},

		"exponentiate positive BigFloat 5 ** 2": {
			a:    NewBigFloat(5),
			b:    NewBigFloat(2),
			want: NewBigFloat(25).SetPrecision(53),
		},
		"exponentiate positive BigFloat 7 ** 8": {
			a:    NewBigFloat(7),
			b:    NewBigFloat(8),
			want: NewBigFloat(5764801).SetPrecision(53),
		},
		"exponentiate positive BigFloat 2.5 ** 2.5": {
			a:    NewBigFloat(2.5),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("9.882117688026186").SetPrecision(53),
		},
		"exponentiate positive BigFloat 3 ** 2.5": {
			a:    NewBigFloat(3),
			b:    NewBigFloat(2.5),
			want: ParseBigFloatPanic("15.588457268119896").SetPrecision(53),
		},
		"exponentiate positive BigFloat 6 ** 1": {
			a:    NewBigFloat(6),
			b:    NewBigFloat(1),
			want: NewBigFloat(6).SetPrecision(53),
		},
		"exponentiate negative BigFloat": {
			a:    NewBigFloat(4),
			b:    NewBigFloat(-2),
			want: NewBigFloat(0.0625).SetPrecision(53),
		},
		"exponentiate BigFloat zero": {
			a:    NewBigFloat(25),
			b:    NewBigFloat(0),
			want: NewBigFloat(1).SetPrecision(53),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				cmpopts.IgnoreFields(BigFloat{}, "acc"),
				cmp.AllowUnexported(Error{}, BigInt{}, BigFloat{}),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Mod(tc.left, tc.right)
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
		"mod by String and return an error": {
			a:   NewBigFloat(5),
			b:   String("foo"),
			err: NewError(TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"mod by SmallInt 18446744073709551616 % 2": {
			a:    ParseBigFloatPanic("18446744073709551616"),
			b:    SmallInt(2),
			want: NewBigFloat(0).SetPrecision(67),
		},
		"mod by SmallInt 25 % 3": {
			a:    NewBigFloat(25),
			b:    SmallInt(3),
			want: NewBigFloat(1).SetPrecision(64),
		},
		"mod by SmallInt 25.6 % 3": {
			a:    NewBigFloat(25.6),
			b:    SmallInt(3),
			want: NewBigFloat(1.6000000000000014).SetPrecision(64),
		},
		"mod by SmallInt 76 % 6": {
			a:    NewBigFloat(76),
			b:    SmallInt(6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"mod by SmallInt -76 % 6": {
			a:    NewBigFloat(-76),
			b:    SmallInt(6),
			want: NewBigFloat(-4).SetPrecision(64),
		},
		"mod by SmallInt 76 % -6": {
			a:    NewBigFloat(76),
			b:    SmallInt(-6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"mod by SmallInt -76 % -6": {
			a:    NewBigFloat(-76),
			b:    SmallInt(-6),
			want: NewBigFloat(-4).SetPrecision(64),
		},
		"mod by SmallInt 124 % 9": {
			a:    NewBigFloat(124),
			b:    SmallInt(9),
			want: NewBigFloat(7).SetPrecision(64),
		},

		"mod by BigInt 25 % 3": {
			a:    NewBigFloat(25),
			b:    NewBigInt(3),
			want: NewBigFloat(1).SetPrecision(64),
		},
		"mod by BigInt 76 % 6": {
			a:    NewBigFloat(76),
			b:    NewBigInt(6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"mod by BigInt 76.5 % 6": {
			a:    NewBigFloat(76.5),
			b:    NewBigInt(6),
			want: NewBigFloat(4.5).SetPrecision(64),
		},
		"mod by BigInt -76 % 6": {
			a:    NewBigFloat(-76),
			b:    NewBigInt(6),
			want: NewBigFloat(-4).SetPrecision(64),
		},
		"mod by BigInt 76 % -6": {
			a:    NewBigFloat(76),
			b:    NewBigInt(-6),
			want: NewBigFloat(4).SetPrecision(64),
		},
		"mod by BigInt -76 % -6": {
			a:    NewBigFloat(-76),
			b:    NewBigInt(-6),
			want: NewBigFloat(-4).SetPrecision(64),
		},
		"mod by BigInt 124 % 9": {
			a:    NewBigFloat(124),
			b:    NewBigInt(9),
			want: NewBigFloat(7).SetPrecision(64),
		},
		"mod by BigInt 9765 % 9223372036854775808": {
			a:    NewBigFloat(9765),
			b:    ParseBigIntPanic("9223372036854775808", 10),
			want: NewBigFloat(9765).SetPrecision(64),
		},

		"mod by Float 25 % 3": {
			a:    NewBigFloat(25),
			b:    Float(3),
			want: NewBigFloat(1),
		},
		"mod by Float 76 % 6": {
			a:    NewBigFloat(76),
			b:    Float(6),
			want: NewBigFloat(4),
		},
		"mod by Float 124 % 9": {
			a:    NewBigFloat(124),
			b:    Float(9),
			want: NewBigFloat(7),
		},
		// "mod by Float 124 % +Inf": {
		// 	a:    NewBigFloat(124),
		// 	b:    FloatInf(),
		// 	want: NewBigFloat(124),
		// },
		// "mod by Float 124 % -Inf": {
		// 	a:    NewBigFloat(124),
		// 	b:    FloatNegInf(),
		// 	want: NewBigFloat(124),
		// },
		"mod by Float 74.5 % 6.25": {
			a:    NewBigFloat(74.5),
			b:    Float(6.25),
			want: NewBigFloat(5.75),
		},
		"mod by Float 74 % 6.25": {
			a:    NewBigFloat(74),
			b:    Float(6.25),
			want: NewBigFloat(5.25),
		},
		"mod by Float -74 % 6.25": {
			a:    NewBigFloat(-74),
			b:    Float(6.25),
			want: NewBigFloat(-5.25),
		},
		"mod by Float 74 % -6.25": {
			a:    NewBigFloat(74),
			b:    Float(-6.25),
			want: NewBigFloat(5.25),
		},
		"mod by Float -74 % -6.25": {
			a:    NewBigFloat(-74),
			b:    Float(-6.25),
			want: NewBigFloat(-5.25),
		},

		"mod by BigFloat 25 % 3": {
			a:    NewBigFloat(25),
			b:    NewBigFloat(3),
			want: NewBigFloat(1),
		},
		"mod by BigFloat 76 % 6": {
			a:    NewBigFloat(76),
			b:    NewBigFloat(6),
			want: NewBigFloat(4),
		},
		"mod by BigFloat 124 % 9": {
			a:    NewBigFloat(124),
			b:    NewBigFloat(9),
			want: NewBigFloat(7),
		},
		"mod by BigFloat 74 % 6.25": {
			a:    NewBigFloat(74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(5.25),
		},
		"mod by BigFloat 74 % 6.25 with higher precision": {
			a:    NewBigFloat(74),
			b:    NewBigFloat(6.25).SetPrecision(64),
			want: NewBigFloat(5.25).SetPrecision(64),
		},
		"mod by BigFloat -74 % 6.25": {
			a:    NewBigFloat(-74),
			b:    NewBigFloat(6.25),
			want: NewBigFloat(-5.25),
		},
		"mod by BigFloat 74 % -6.25": {
			a:    NewBigFloat(74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(5.25),
		},
		"mod by BigFloat -74 % -6.25": {
			a:    NewBigFloat(-74),
			b:    NewBigFloat(-6.25),
			want: NewBigFloat(-5.25),
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
				t.Log(got.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
