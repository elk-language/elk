package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestBigFloatAdd(t *testing.T) {
	tests := map[string]struct {
		left  *value.BigFloat
		right value.Value
		want  value.Value
		err   *value.Error
	}{
		"BigFloat + BigFloat => BigFloat": {
			left:  value.NewBigFloat(2.5),
			right: value.NewBigFloat(10.2),
			want:  value.NewBigFloat(12.7),
		},
		"BigFloat + BigFloat NaN => BigFloat NaN": {
			left:  value.NewBigFloat(2.5),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN + BigFloat => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN + BigFloat NaN => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf + BigFloat => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatInf(),
		},
		"BigFloat + BigFloat +Inf => BigFloat +Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.BigFloatInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat +Inf + BigFloat +Inf => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.BigFloatInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf + BigFloat => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat + BigFloat -Inf => BigFloat -Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf + BigFloat -Inf => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat +Inf + BigFloat -Inf => BigFloat NaN": {
			left:  value.BigFloatInf(),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf + BigFloat +Inf => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.BigFloatInf(),
			want:  value.BigFloatNaN(),
		},
		"result takes the max precision from its operands": {
			left:  value.NewBigFloat(2.5).SetPrecision(31),
			right: value.NewBigFloat(10.2).SetPrecision(54),
			want:  value.NewBigFloat(12.7).SetPrecision(54),
		},
		"result takes the max precision from its operands (left)": {
			left:  value.NewBigFloat(2.5).SetPrecision(54),
			right: value.NewBigFloat(10.2).SetPrecision(52),
			want:  value.NewBigFloat(12.7).SetPrecision(54),
		},
		"BigFloat + SmallInt => BigFloat": {
			left:  value.NewBigFloat(2.5),
			right: value.SmallInt(120),
			want:  value.NewBigFloat(122.5).SetPrecision(64),
		},
		"BigFloat + BigInt => BigFloat": {
			left:  value.NewBigFloat(2.5),
			right: value.NewBigInt(120),
			want:  value.NewBigFloat(122.5).SetPrecision(64),
		},
		"BigFloat + Int64 => TypeError": {
			left:  value.NewBigFloat(2.5),
			right: value.Int64(20),
			err:   value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat + String => TypeError": {
			left:  value.NewBigFloat(2.5),
			right: value.String("foo"),
			err:   value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"BigFloat + Float => BigFloat": {
			left:  value.NewBigFloat(2.5),
			right: value.Float(5.2),
			want:  value.NewBigFloat(7.7),
		},
		"BigFloat + Float NaN => BigFloat NaN": {
			left:  value.NewBigFloat(2.5),
			right: value.FloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN + Float => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.Float(10.2),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN + Float NaN => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.FloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf + Float => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.Float(10.2),
			want:  value.BigFloatInf(),
		},
		"BigFloat + Float +Inf => BigFloat +Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.FloatInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat +Inf + Float +Inf => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.FloatInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf + Float => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.Float(10.2),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat + Float -Inf => BigFloat -Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.FloatNegInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf + Float -Inf => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.FloatNegInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat +Inf + Float -Inf => BigFloat NaN": {
			left:  value.BigFloatInf(),
			right: value.FloatNegInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf + Float +Inf => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.FloatInf(),
			want:  value.BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Add(tc.right)
			opts := vm.ComparerOptions
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
			got := value.CountFloatDigits(tc.str)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloatSubtract(t *testing.T) {
	tests := map[string]struct {
		left  *value.BigFloat
		right value.Value
		want  value.Value
		err   *value.Error
	}{
		"BigFloat - BigFloat => BigFloat": {
			left:  value.NewBigFloat(10.0),
			right: value.NewBigFloat(2.5),
			want:  value.NewBigFloat(7.5),
		},
		"result takes the max precision from its operands": {
			left:  value.NewBigFloat(10.0).SetPrecision(54),
			right: value.NewBigFloat(2.5).SetPrecision(31),
			want:  value.NewBigFloat(7.5).SetPrecision(54),
		},
		"BigFloat - SmallInt => BigFloat": {
			left:  value.NewBigFloat(120.5),
			right: value.SmallInt(2),
			want:  value.NewBigFloat(118.5).SetPrecision(64),
		},
		"BigFloat - BigInt => BigFloat": {
			left:  value.NewBigFloat(120.5),
			right: value.NewBigInt(2),
			want:  value.NewBigFloat(118.5).SetPrecision(64),
		},
		"BigFloat - Int64 => TypeError": {
			left:  value.NewBigFloat(20.5),
			right: value.Int64(2),
			err:   value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat - String => TypeError": {
			left:  value.NewBigFloat(2.5),
			right: value.String("foo"),
			err:   value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"BigFloat - BigFloat NaN => BigFloat NaN": {
			left:  value.NewBigFloat(2.5),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN - BigFloat => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN - BigFloat NaN => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf - BigFloat => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatInf(),
		},
		"BigFloat - BigFloat +Inf => BigFloat -Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.BigFloatInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat +Inf - BigFloat +Inf => BigFloat NaN": {
			left:  value.BigFloatInf(),
			right: value.BigFloatInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf - BigFloat => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat - BigFloat -Inf => BigFloat +Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf - BigFloat -Inf => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf - BigFloat -Inf => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf - BigFloat +Inf => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.BigFloatInf(),
			want:  value.BigFloatNegInf(),
		},

		"BigFloat - Float NaN => BigFloat NaN": {
			left:  value.NewBigFloat(2.5),
			right: value.FloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN - Float => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.Float(10.2),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN - Float NaN => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.FloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf - Float => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.Float(10.2),
			want:  value.BigFloatInf(),
		},
		"BigFloat - Float +Inf => BigFloat -Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.FloatInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat +Inf - Float +Inf => BigFloat NaN": {
			left:  value.BigFloatInf(),
			right: value.FloatInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf - Float => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.Float(10.2),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat - Float -Inf => BigFloat +Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.FloatNegInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf - Float -Inf => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.FloatNegInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf - Float -Inf => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.FloatNegInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf - Float +Inf => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.FloatInf(),
			want:  value.BigFloatNegInf(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Subtract(tc.right)
			opts := vm.ComparerOptions
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
		left  *value.BigFloat
		right value.Value
		want  value.Value
		err   *value.Error
	}{
		"BigFloat * BigFloat => BigFloat": {
			left:  value.NewBigFloat(2.55),
			right: value.NewBigFloat(10.0),
			want:  value.NewBigFloat(25.5),
		},
		"result takes the max precision from its operands": {
			left:  value.NewBigFloat(2.5).SetPrecision(31),
			right: value.NewBigFloat(10.0).SetPrecision(54),
			want:  value.NewBigFloat(25.0).SetPrecision(54),
		},
		"BigFloat * SmallInt => BigFloat": {
			left:  value.NewBigFloat(2.5),
			right: value.SmallInt(10),
			want:  value.NewBigFloat(25.0).SetPrecision(64),
		},
		"BigFloat * BigInt => BigFloat": {
			left:  value.NewBigFloat(2.5),
			right: value.NewBigInt(10),
			want:  value.NewBigFloat(25.0).SetPrecision(64),
		},
		"BigFloat * Int64 => TypeError": {
			left:  value.NewBigFloat(2.55),
			right: value.Int64(20),
			err:   value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat * String => TypeError": {
			left:  value.NewBigFloat(2.5),
			right: value.String("foo"),
			err:   value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"BigFloat * BigFloat NaN => BigFloat NaN": {
			left:  value.NewBigFloat(2.5),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN * BigFloat => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN * BigFloat NaN => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf * BigFloat => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatInf(),
		},
		"BigFloat * BigFloat +Inf => BigFloat +Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.BigFloatInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat +Inf * BigFloat +Inf => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.BigFloatInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf * +BigFloat => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf * -BigFloat => BigFloat +Inf": {
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(-10.2),
			want:  value.BigFloatInf(),
		},
		"+BigFloat * BigFloat -Inf => BigFloat -Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatNegInf(),
		},
		"-BigFloat * BigFloat -Inf => BigFloat +Inf": {
			left:  value.NewBigFloat(-10.2),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf * BigFloat -Inf => BigFloat +Inf": {
			left:  value.BigFloatNegInf(),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat +Inf * BigFloat -Inf => BigFloat -Inf": {
			left:  value.BigFloatInf(),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf * BigFloat +Inf => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.BigFloatInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf * BigFloat 0 => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(0),
			want:  value.BigFloatNaN(),
		},
		"BigFloat 0 * BigFloat +Inf => BigFloat NaN": {
			left:  value.NewBigFloat(0),
			right: value.BigFloatInf(),
			want:  value.BigFloatNaN(),
		},

		"BigFloat * Float NaN => BigFloat NaN": {
			left:  value.NewBigFloat(2.5),
			right: value.FloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN * Float => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.Float(10.2),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN * Float NaN => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.FloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf * Float => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.Float(10.2),
			want:  value.BigFloatInf(),
		},
		"BigFloat * Float +Inf => BigFloat +Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.FloatInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat +Inf * Float +Inf => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.FloatInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf * +Float => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.Float(10.2),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf * -Float => BigFloat +Inf": {
			left:  value.BigFloatNegInf(),
			right: value.Float(-10.2),
			want:  value.BigFloatInf(),
		},
		"+BigFloat * Float -Inf => BigFloat -Inf": {
			left:  value.NewBigFloat(10.2),
			right: value.FloatNegInf(),
			want:  value.BigFloatNegInf(),
		},
		"-BigFloat * Float -Inf => BigFloat +Inf": {
			left:  value.NewBigFloat(-10.2),
			right: value.FloatNegInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat -Inf * Float -Inf => BigFloat +Inf": {
			left:  value.BigFloatNegInf(),
			right: value.FloatNegInf(),
			want:  value.BigFloatInf(),
		},
		"BigFloat +Inf * Float -Inf => BigFloat -Inf": {
			left:  value.BigFloatInf(),
			right: value.FloatNegInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf * Float +Inf => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.FloatInf(),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf * Float 0 => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.Float(0),
			want:  value.BigFloatNaN(),
		},
		"BigFloat 0 * Float +Inf => BigFloat NaN": {
			left:  value.NewBigFloat(0),
			right: value.FloatInf(),
			want:  value.BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Multiply(tc.right)
			opts := vm.ComparerOptions
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
		left  *value.BigFloat
		right value.Value
		want  value.Value
		err   *value.Error
	}{
		"BigFloat / BigFloat => BigFloat": {
			left:  value.NewBigFloat(2.68),
			right: value.NewBigFloat(2),
			want:  value.NewBigFloat(1.34),
		},
		"result takes the max precision from its operands": {
			left:  value.NewBigFloat(2).SetPrecision(31),
			right: value.NewBigFloat(2).SetPrecision(54),
			want:  value.NewBigFloat(1).SetPrecision(54),
		},
		"BigFloat / SmallInt => BigFloat": {
			left:  value.NewBigFloat(2.68),
			right: value.SmallInt(2),
			want:  value.NewBigFloat(1.34).SetPrecision(64),
		},
		"BigFloat / BigInt => BigFloat": {
			left:  value.NewBigFloat(2.68),
			right: value.NewBigInt(2),
			want:  value.NewBigFloat(1.34).SetPrecision(64),
		},
		"BigFloat / Int64 => TypeError": {
			left:  value.NewBigFloat(2.68),
			right: value.Int64(2),
			err:   value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"BigFloat / String => TypeError": {
			left:  value.NewBigFloat(2.5),
			right: value.String("foo"),
			err:   value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"BigFloat / BigFloat NaN => BigFloat NaN": {
			left:  value.NewBigFloat(2.5),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN / BigFloat => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN / BigFloat NaN => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf / BigFloat => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatInf(),
		},
		"BigFloat / BigFloat +Inf => BigFloat 0": {
			left:  value.NewBigFloat(10.2),
			right: value.BigFloatInf(),
			want:  value.NewBigFloat(0),
		},
		"BigFloat +Inf / BigFloat +Inf => BigFloat NaN": {
			left:  value.BigFloatInf(),
			right: value.BigFloatInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf / +BigFloat => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(10.2),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf / -BigFloat => BigFloat +Inf": {
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(-10.2),
			want:  value.BigFloatInf(),
		},
		"+BigFloat / BigFloat -Inf => BigFloat -0": {
			left:  value.NewBigFloat(10.2),
			right: value.BigFloatNegInf(),
			want:  value.NewBigFloat(-0),
		},
		"-BigFloat / BigFloat -Inf => BigFloat +0": {
			left:  value.NewBigFloat(-10.2),
			right: value.BigFloatNegInf(),
			want:  value.NewBigFloat(0),
		},
		"BigFloat -Inf / BigFloat -Inf => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf / BigFloat -Inf => BigFloat NaN": {
			left:  value.BigFloatInf(),
			right: value.BigFloatNegInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf / BigFloat +Inf => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.BigFloatInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf / BigFloat 0 => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(0),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat 0 / BigFloat +Inf => BigFloat 0": {
			left:  value.NewBigFloat(0),
			right: value.BigFloatInf(),
			want:  value.NewBigFloat(0),
		},

		"BigFloat / Float NaN => BigFloat NaN": {
			left:  value.NewBigFloat(2.5),
			right: value.FloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN / Float => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.Float(10.2),
			want:  value.BigFloatNaN(),
		},
		"BigFloat NaN / Float NaN => BigFloat NaN": {
			left:  value.BigFloatNaN(),
			right: value.FloatNaN(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf / Float => BigFloat +Inf": {
			left:  value.BigFloatInf(),
			right: value.Float(10.2),
			want:  value.BigFloatInf(),
		},
		"BigFloat / Float +Inf => BigFloat 0": {
			left:  value.NewBigFloat(10.2),
			right: value.FloatInf(),
			want:  value.NewBigFloat(0),
		},
		"BigFloat +Inf / Float +Inf => BigFloat NaN": {
			left:  value.BigFloatInf(),
			right: value.FloatInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf / +Float => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.Float(10.2),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat -Inf / -Float => BigFloat +Inf": {
			left:  value.BigFloatNegInf(),
			right: value.Float(-10.2),
			want:  value.BigFloatInf(),
		},
		"+BigFloat / Float -Inf => BigFloat -0": {
			left:  value.NewBigFloat(10.2),
			right: value.FloatNegInf(),
			want:  value.NewBigFloat(-0),
		},
		"-BigFloat / Float -Inf => BigFloat +0": {
			left:  value.NewBigFloat(-10.2),
			right: value.FloatNegInf(),
			want:  value.NewBigFloat(0),
		},
		"BigFloat -Inf / Float -Inf => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.FloatNegInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat +Inf / Float -Inf => BigFloat NaN": {
			left:  value.BigFloatInf(),
			right: value.FloatNegInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf / Float +Inf => BigFloat NaN": {
			left:  value.BigFloatNegInf(),
			right: value.FloatInf(),
			want:  value.BigFloatNaN(),
		},
		"BigFloat -Inf / Float 0 => BigFloat -Inf": {
			left:  value.BigFloatNegInf(),
			right: value.Float(0),
			want:  value.BigFloatNegInf(),
		},
		"BigFloat 0 / Float +Inf => BigFloat 0": {
			left:  value.NewBigFloat(0),
			right: value.FloatInf(),
			want:  value.NewBigFloat(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Divide(tc.right)
			opts := vm.ComparerOptions
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
		a    *value.BigFloat
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"exponentiate String and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
		"exponentiate Int32 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Int32(2),
			err: value.NewError(value.TypeErrorClass, "`Std::Int32` can't be coerced into `Std::BigFloat`"),
		},
		"SmallInt 5 ** 2": {
			a:    value.NewBigFloat(5),
			b:    value.SmallInt(2),
			want: value.NewBigFloat(25).SetPrecision(64),
		},
		"SmallInt 5p92 ** 2": {
			a:    value.NewBigFloat(5).SetPrecision(92),
			b:    value.SmallInt(2),
			want: value.NewBigFloat(25).SetPrecision(92),
		},
		"SmallInt 7 ** 8": {
			a:    value.NewBigFloat(7),
			b:    value.SmallInt(8),
			want: value.NewBigFloat(5764801).SetPrecision(64),
		},
		"SmallInt 2.5 ** 5": {
			a:    value.NewBigFloat(2.5),
			b:    value.SmallInt(5),
			want: value.NewBigFloat(97.65625).SetPrecision(64),
		},
		"SmallInt 7.12 ** 1": {
			a:    value.NewBigFloat(7.12),
			b:    value.SmallInt(1),
			want: value.NewBigFloat(7.12).SetPrecision(64),
		},
		"SmallInt 4 ** -2": {
			a:    value.NewBigFloat(4),
			b:    value.SmallInt(-2),
			want: value.NewBigFloat(0.0625).SetPrecision(64),
		},
		"SmallInt 25 ** 0": {
			a:    value.NewBigFloat(25),
			b:    value.SmallInt(0),
			want: value.NewBigFloat(1).SetPrecision(64),
		},

		"BigInt 5 ** 2": {
			a:    value.NewBigFloat(5),
			b:    value.NewBigInt(2),
			want: value.NewBigFloat(25).SetPrecision(64),
		},
		"BigInt 5p78 ** 2": {
			a:    value.NewBigFloat(5).SetPrecision(78),
			b:    value.NewBigInt(2),
			want: value.NewBigFloat(25).SetPrecision(78),
		},
		"BigInt 7 ** 8": {
			a:    value.NewBigFloat(7),
			b:    value.NewBigInt(8),
			want: value.NewBigFloat(5764801).SetPrecision(64),
		},
		"BigInt 2.5 ** 5": {
			a:    value.NewBigFloat(2.5),
			b:    value.NewBigInt(5),
			want: value.NewBigFloat(97.65625).SetPrecision(64),
		},
		"BigInt 7.12 ** 1": {
			a:    value.NewBigFloat(7.12),
			b:    value.NewBigInt(1),
			want: value.NewBigFloat(7.12).SetPrecision(64),
		},
		"BigInt 4 ** -2": {
			a:    value.NewBigFloat(4),
			b:    value.NewBigInt(-2),
			want: value.NewBigFloat(0.0625).SetPrecision(64),
		},
		"BigInt 25 ** 0": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigInt(0),
			want: value.NewBigFloat(1).SetPrecision(64),
		},

		"Float 5 ** 2": {
			a:    value.NewBigFloat(5),
			b:    value.Float(2),
			want: value.NewBigFloat(25),
		},
		"Float 5p83 ** 2": {
			a:    value.NewBigFloat(5).SetPrecision(83),
			b:    value.Float(2),
			want: value.NewBigFloat(25).SetPrecision(83),
		},
		"Float 7 ** 8": {
			a:    value.NewBigFloat(7),
			b:    value.Float(8),
			want: value.NewBigFloat(5764801),
		},
		"Float 2.5 ** 2.5": {
			a:    value.NewBigFloat(2.5),
			b:    value.Float(2.5),
			want: value.NewBigFloat(9.882117688026186),
		},
		"Float 3 ** 2.5": {
			a:    value.NewBigFloat(3),
			b:    value.Float(2.5),
			want: value.NewBigFloat(15.588457268119896),
		},
		"Float 6 ** 1": {
			a:    value.NewBigFloat(6),
			b:    value.Float(1),
			want: value.NewBigFloat(6),
		},
		"Float 4 ** -2": {
			a:    value.NewBigFloat(4),
			b:    value.Float(-2),
			want: value.NewBigFloat(0.0625),
		},
		"Float 25 ** 0": {
			a:    value.NewBigFloat(25),
			b:    value.Float(0),
			want: value.NewBigFloat(1),
		},
		"Float 25 ** NaN": {
			a:    value.NewBigFloat(25),
			b:    value.FloatNaN(),
			want: value.BigFloatNaN(),
		},
		"Float NaN ** 25": {
			a:    value.BigFloatNaN(),
			b:    value.Float(25),
			want: value.BigFloatNaN(),
		},
		"Float NaN ** NaN": {
			a:    value.BigFloatNaN(),
			b:    value.FloatNaN(),
			want: value.BigFloatNaN(),
		},
		"Float 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    value.NewBigFloat(0),
			b:    value.Float(-5),
			want: value.BigFloatInf(),
		},
		"Float 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    value.NewBigFloat(0),
			b:    value.FloatNegInf(),
			want: value.BigFloatInf(),
		},
		"Float 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    value.NewBigFloat(0),
			b:    value.FloatInf(),
			want: value.NewBigFloat(0),
		},
		"Float 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    value.NewBigFloat(0),
			b:    value.Float(-8),
			want: value.BigFloatInf(),
		},
		"Float 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    value.NewBigFloat(0),
			b:    value.Float(7),
			want: value.NewBigFloat(0),
		},
		"Float 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    value.NewBigFloat(0),
			b:    value.Float(8),
			want: value.NewBigFloat(0),
		},
		"Float -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    value.NewBigFloat(-1),
			b:    value.FloatInf(),
			want: value.NewBigFloat(1),
		},
		"Float -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    value.NewBigFloat(-1),
			b:    value.FloatNegInf(),
			want: value.NewBigFloat(1),
		},
		"Float 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.NewBigFloat(2),
			b:    value.FloatInf(),
			want: value.BigFloatInf(),
		},
		"Float -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.NewBigFloat(-2),
			b:    value.FloatInf(),
			want: value.BigFloatInf(),
		},
		"Float 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.NewBigFloat(2),
			b:    value.FloatNegInf(),
			want: value.NewBigFloat(0),
		},
		"Float -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.NewBigFloat(-2),
			b:    value.FloatNegInf(),
			want: value.NewBigFloat(0),
		},
		"Float 0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    value.NewBigFloat(0.5),
			b:    value.FloatInf(),
			want: value.NewBigFloat(0),
		},
		"Float -0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    value.NewBigFloat(-0.5),
			b:    value.FloatInf(),
			want: value.NewBigFloat(0),
		},
		"Float 0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    value.NewBigFloat(0.5),
			b:    value.FloatNegInf(),
			want: value.BigFloatInf(),
		},
		"Float -0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    value.NewBigFloat(-0.5),
			b:    value.FloatNegInf(),
			want: value.BigFloatInf(),
		},
		"Float +Inf ** 5": { // Pow(+Inf, y) = +Inf for y > 0
			a:    value.BigFloatInf(),
			b:    value.Float(5),
			want: value.BigFloatInf(),
		},
		"Float +Inf ** -7": { // Pow(+Inf, y) = +0 for y < 0
			a:    value.BigFloatInf(),
			b:    value.Float(-7),
			want: value.NewBigFloat(0),
		},
		"Float -Inf ** -7": {
			a:    value.BigFloatNegInf(),
			b:    value.Float(-7),
			want: value.NewBigFloat(0),
		},
		"Float -5.5 ** 3.8": { // Pow(x, y) = NaN for finite x < 0 and finite non-integer y
			a:    value.NewBigFloat(-5.5),
			b:    value.Float(3.8),
			want: value.BigFloatNaN(),
		},

		"BigFloat 5 ** 2": {
			a:    value.NewBigFloat(5),
			b:    value.NewBigFloat(2),
			want: value.NewBigFloat(25).SetPrecision(53),
		},
		"BigFloat 7 ** 8": {
			a:    value.NewBigFloat(7),
			b:    value.NewBigFloat(8),
			want: value.NewBigFloat(5764801).SetPrecision(53),
		},
		"BigFloat 2.5 ** 2.5": {
			a:    value.NewBigFloat(2.5),
			b:    value.NewBigFloat(2.5),
			want: value.ParseBigFloatPanic("9.882117688026186").SetPrecision(53),
		},
		"BigFloat 3 ** 2.5": {
			a:    value.NewBigFloat(3),
			b:    value.NewBigFloat(2.5),
			want: value.ParseBigFloatPanic("15.588457268119896").SetPrecision(53),
		},
		"BigFloat 6 ** 1": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(1),
			want: value.NewBigFloat(6).SetPrecision(53),
		},
		"BigFloat 4 ** -2": {
			a:    value.NewBigFloat(4),
			b:    value.NewBigFloat(-2),
			want: value.NewBigFloat(0.0625).SetPrecision(53),
		},
		"BigFloat 25 ** 0": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(0),
			want: value.NewBigFloat(1).SetPrecision(53),
		},
		"BigFloat 25 ** NaN": {
			a:    value.NewBigFloat(25),
			b:    value.BigFloatNaN(),
			want: value.BigFloatNaN(),
		},
		"BigFloat NaN ** 25": {
			a:    value.BigFloatNaN(),
			b:    value.NewBigFloat(25),
			want: value.BigFloatNaN(),
		},
		"BigFloat NaN ** NaN": {
			a:    value.BigFloatNaN(),
			b:    value.BigFloatNaN(),
			want: value.BigFloatNaN(),
		},
		"BigFloat 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    value.NewBigFloat(0),
			b:    value.NewBigFloat(-5),
			want: value.BigFloatInf(),
		},
		"BigFloat 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    value.NewBigFloat(0),
			b:    value.BigFloatNegInf(),
			want: value.BigFloatInf(),
		},
		"BigFloat 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    value.NewBigFloat(0),
			b:    value.BigFloatInf(),
			want: value.NewBigFloat(0),
		},
		"BigFloat 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    value.NewBigFloat(0),
			b:    value.NewBigFloat(-8),
			want: value.BigFloatInf(),
		},
		"BigFloat 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    value.NewBigFloat(0),
			b:    value.NewBigFloat(7),
			want: value.NewBigFloat(0),
		},
		"BigFloat 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    value.NewBigFloat(0),
			b:    value.NewBigFloat(8),
			want: value.NewBigFloat(0),
		},
		"BigFloat -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    value.NewBigFloat(-1),
			b:    value.BigFloatInf(),
			want: value.NewBigFloat(1),
		},
		"BigFloat -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    value.NewBigFloat(-1),
			b:    value.BigFloatNegInf(),
			want: value.NewBigFloat(1),
		},
		"BigFloat 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.NewBigFloat(2),
			b:    value.BigFloatInf(),
			want: value.BigFloatInf(),
		},
		"BigFloat -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.NewBigFloat(-2),
			b:    value.BigFloatInf(),
			want: value.BigFloatInf(),
		},
		"BigFloat 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.NewBigFloat(2),
			b:    value.BigFloatNegInf(),
			want: value.NewBigFloat(0),
		},
		"BigFloat -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.NewBigFloat(-2),
			b:    value.BigFloatNegInf(),
			want: value.NewBigFloat(0),
		},
		"BigFloat 0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    value.NewBigFloat(0.5),
			b:    value.BigFloatInf(),
			want: value.NewBigFloat(0),
		},
		"BigFloat -0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    value.NewBigFloat(-0.5),
			b:    value.BigFloatInf(),
			want: value.NewBigFloat(0),
		},
		"BigFloat 0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    value.NewBigFloat(0.5),
			b:    value.BigFloatNegInf(),
			want: value.BigFloatInf(),
		},
		"BigFloat -0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    value.NewBigFloat(-0.5),
			b:    value.BigFloatNegInf(),
			want: value.BigFloatInf(),
		},
		"BigFloat +Inf ** 5": { // Pow(+Inf, y) = +Inf for y > 0
			a:    value.BigFloatInf(),
			b:    value.NewBigFloat(5),
			want: value.BigFloatInf(),
		},
		"BigFloat +Inf ** -7": { // Pow(+Inf, y) = +0 for y < 0
			a:    value.BigFloatInf(),
			b:    value.NewBigFloat(-7),
			want: value.NewBigFloat(0),
		},
		"BigFloat -Inf ** -7": {
			a:    value.BigFloatNegInf(),
			b:    value.NewBigFloat(-7),
			want: value.NewBigFloat(0),
		},
		"BigFloat -5.5 ** 3.8": { // Pow(x, y) = NaN for finite x < 0 and finite non-integer y
			a:    value.NewBigFloat(-5.5),
			b:    value.NewBigFloat(3.8),
			want: value.BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
			opts := vm.ComparerOptions
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
		left  *value.BigFloat
		right *value.BigFloat
		want  *value.BigFloat
	}{
		"25 % 3": {
			left:  value.NewBigFloat(25),
			right: value.NewBigFloat(3),
			want:  value.NewBigFloat(1),
		},
		"76 % 6": {
			left:  value.NewBigFloat(76),
			right: value.NewBigFloat(6),
			want:  value.NewBigFloat(4),
		},
		"76.75 % 6.25": {
			left:  value.NewBigFloat(76.75),
			right: value.NewBigFloat(6.25),
			want:  value.NewBigFloat(1.75),
		},
		"76.75 % -6.25": {
			left:  value.NewBigFloat(76.75),
			right: value.NewBigFloat(-6.25),
			want:  value.NewBigFloat(1.75),
		},
		"-76.75 % 6.25": {
			left:  value.NewBigFloat(-76.75),
			right: value.NewBigFloat(6.25),
			want:  value.NewBigFloat(-1.75),
		},
		"-76.75 % -6.25": {
			left:  value.NewBigFloat(-76.75),
			right: value.NewBigFloat(-6.25),
			want:  value.NewBigFloat(-1.75),
		},
		"+Inf % 5": { // Mod(±Inf, y) = NaN
			left:  value.BigFloatInf(),
			right: value.NewBigFloat(5),
			want:  value.BigFloatNaN(),
		},
		"-Inf % 5": { // Mod(±Inf, y) = NaN
			left:  value.BigFloatNegInf(),
			right: value.NewBigFloat(5),
			want:  value.BigFloatNaN(),
		},
		"NaN % 625": { // Mod(NaN, y) = NaN
			left:  value.BigFloatNaN(),
			right: value.NewBigFloat(625),
			want:  value.BigFloatNaN(),
		},
		"25 % 0": { // Mod(x, 0) = NaN
			left:  value.NewBigFloat(25),
			right: value.NewBigFloat(0),
			want:  value.BigFloatNaN(),
		},
		"25 % +Inf": { // Mod(x, ±Inf) = x
			left:  value.NewBigFloat(25),
			right: value.BigFloatInf(),
			want:  value.NewBigFloat(25),
		},
		"-87 % -Inf": { // Mod(x, ±Inf) = x
			left:  value.NewBigFloat(-87),
			right: value.BigFloatNegInf(),
			want:  value.NewBigFloat(-87),
		},
		"49 % NaN": { // Mod(x, NaN) = NaN
			left:  value.NewBigFloat(49),
			right: value.BigFloatNaN(),
			want:  value.BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Mod(tc.left, tc.right)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloat_FloorBigFloat(t *testing.T) {
	tests := map[string]struct {
		f    *value.BigFloat
		want *value.BigFloat
	}{
		"floor(25)": {
			f:    value.NewBigFloat(25),
			want: value.NewBigFloat(25),
		},
		"floor(38.7)": {
			f:    value.NewBigFloat(38.7),
			want: value.NewBigFloat(38),
		},
		"floor(-6.5)": {
			f:    value.NewBigFloat(-6.5),
			want: value.NewBigFloat(-7),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.f.FloorBigFloat()
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Log(got.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloat_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigFloat
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},

		"SmallInt 18446744073709551616 % 2": {
			a:    value.ParseBigFloatPanic("18446744073709551616"),
			b:    value.SmallInt(2),
			want: value.NewBigFloat(0).SetPrecision(67),
		},
		"SmallInt 25 % 3": {
			a:    value.NewBigFloat(25),
			b:    value.SmallInt(3),
			want: value.NewBigFloat(1).SetPrecision(64),
		},
		"SmallInt 25.6 % 3": {
			a:    value.NewBigFloat(25.6),
			b:    value.SmallInt(3),
			want: value.NewBigFloat(1.6000000000000014).SetPrecision(64),
		},
		"SmallInt 76 % 6": {
			a:    value.NewBigFloat(76),
			b:    value.SmallInt(6),
			want: value.NewBigFloat(4).SetPrecision(64),
		},
		"SmallInt -76 % 6": {
			a:    value.NewBigFloat(-76),
			b:    value.SmallInt(6),
			want: value.NewBigFloat(-4).SetPrecision(64),
		},
		"SmallInt 76 % -6": {
			a:    value.NewBigFloat(76),
			b:    value.SmallInt(-6),
			want: value.NewBigFloat(4).SetPrecision(64),
		},
		"SmallInt -76 % -6": {
			a:    value.NewBigFloat(-76),
			b:    value.SmallInt(-6),
			want: value.NewBigFloat(-4).SetPrecision(64),
		},
		"SmallInt 124 % 9": {
			a:    value.NewBigFloat(124),
			b:    value.SmallInt(9),
			want: value.NewBigFloat(7).SetPrecision(64),
		},

		"BigInt 25 % 3": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigInt(3),
			want: value.NewBigFloat(1).SetPrecision(64),
		},
		"BigInt 76 % 6": {
			a:    value.NewBigFloat(76),
			b:    value.NewBigInt(6),
			want: value.NewBigFloat(4).SetPrecision(64),
		},
		"BigInt 76.5 % 6": {
			a:    value.NewBigFloat(76.5),
			b:    value.NewBigInt(6),
			want: value.NewBigFloat(4.5).SetPrecision(64),
		},
		"BigInt -76 % 6": {
			a:    value.NewBigFloat(-76),
			b:    value.NewBigInt(6),
			want: value.NewBigFloat(-4).SetPrecision(64),
		},
		"BigInt 76 % -6": {
			a:    value.NewBigFloat(76),
			b:    value.NewBigInt(-6),
			want: value.NewBigFloat(4).SetPrecision(64),
		},
		"BigInt -76 % -6": {
			a:    value.NewBigFloat(-76),
			b:    value.NewBigInt(-6),
			want: value.NewBigFloat(-4).SetPrecision(64),
		},
		"BigInt 124 % 9": {
			a:    value.NewBigFloat(124),
			b:    value.NewBigInt(9),
			want: value.NewBigFloat(7).SetPrecision(64),
		},
		"BigInt 9765 % 9223372036854775808": {
			a:    value.NewBigFloat(9765),
			b:    value.ParseBigIntPanic("9223372036854775808", 10),
			want: value.NewBigFloat(9765).SetPrecision(64),
		},

		"Float 25 % 3": {
			a:    value.NewBigFloat(25),
			b:    value.Float(3),
			want: value.NewBigFloat(1),
		},
		"Float 25p102 % 3": {
			a:    value.NewBigFloat(25).SetPrecision(102),
			b:    value.Float(3),
			want: value.NewBigFloat(1).SetPrecision(102),
		},
		"Float 76 % 6": {
			a:    value.NewBigFloat(76),
			b:    value.Float(6),
			want: value.NewBigFloat(4),
		},
		"Float 124 % 9": {
			a:    value.NewBigFloat(124),
			b:    value.Float(9),
			want: value.NewBigFloat(7),
		},
		"Float 74.5 % 6.25": {
			a:    value.NewBigFloat(74.5),
			b:    value.Float(6.25),
			want: value.NewBigFloat(5.75),
		},
		"Float 74 % 6.25": {
			a:    value.NewBigFloat(74),
			b:    value.Float(6.25),
			want: value.NewBigFloat(5.25),
		},
		"Float -74 % 6.25": {
			a:    value.NewBigFloat(-74),
			b:    value.Float(6.25),
			want: value.NewBigFloat(-5.25),
		},
		"Float 74 % -6.25": {
			a:    value.NewBigFloat(74),
			b:    value.Float(-6.25),
			want: value.NewBigFloat(5.25),
		},
		"Float -74 % -6.25": {
			a:    value.NewBigFloat(-74),
			b:    value.Float(-6.25),
			want: value.NewBigFloat(-5.25),
		},
		"Float +Inf % 5": { // Mod(±Inf, y) = NaN
			a:    value.BigFloatInf(),
			b:    value.Float(5),
			want: value.BigFloatNaN(),
		},
		"Float -Inf % 5": { // Mod(±Inf, y) = NaN
			a:    value.BigFloatNegInf(),
			b:    value.Float(5),
			want: value.BigFloatNaN(),
		},
		"Float NaN % 625": { // Mod(NaN, y) = NaN
			a:    value.BigFloatNaN(),
			b:    value.Float(625),
			want: value.BigFloatNaN(),
		},
		"Float 25 % 0": { // Mod(x, 0) = NaN
			a:    value.NewBigFloat(25),
			b:    value.Float(0),
			want: value.BigFloatNaN(),
		},
		"Float 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    value.NewBigFloat(25),
			b:    value.FloatInf(),
			want: value.NewBigFloat(25),
		},
		"Float -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    value.NewBigFloat(-87),
			b:    value.FloatNegInf(),
			want: value.NewBigFloat(-87),
		},
		"Float 49 % NaN": { // Mod(x, NaN) = NaN
			a:    value.NewBigFloat(49),
			b:    value.FloatNaN(),
			want: value.BigFloatNaN(),
		},

		"BigFloat 25 % 3": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(3),
			want: value.NewBigFloat(1),
		},
		"BigFloat 76 % 6": {
			a:    value.NewBigFloat(76),
			b:    value.NewBigFloat(6),
			want: value.NewBigFloat(4),
		},
		"BigFloat 76p82 % 6": {
			a:    value.NewBigFloat(76).SetPrecision(82),
			b:    value.NewBigFloat(6),
			want: value.NewBigFloat(4).SetPrecision(82),
		},
		"BigFloat 76p82 % 6p96": {
			a:    value.NewBigFloat(76).SetPrecision(82),
			b:    value.NewBigFloat(6).SetPrecision(96),
			want: value.NewBigFloat(4).SetPrecision(96),
		},
		"BigFloat 124 % 9": {
			a:    value.NewBigFloat(124),
			b:    value.NewBigFloat(9),
			want: value.NewBigFloat(7),
		},
		"BigFloat 74 % 6.25": {
			a:    value.NewBigFloat(74),
			b:    value.NewBigFloat(6.25),
			want: value.NewBigFloat(5.25),
		},
		"BigFloat 74 % 6.25 with higher precision": {
			a:    value.NewBigFloat(74),
			b:    value.NewBigFloat(6.25).SetPrecision(64),
			want: value.NewBigFloat(5.25).SetPrecision(64),
		},
		"BigFloat -74 % 6.25": {
			a:    value.NewBigFloat(-74),
			b:    value.NewBigFloat(6.25),
			want: value.NewBigFloat(-5.25),
		},
		"BigFloat 74 % -6.25": {
			a:    value.NewBigFloat(74),
			b:    value.NewBigFloat(-6.25),
			want: value.NewBigFloat(5.25),
		},
		"BigFloat -74 % -6.25": {
			a:    value.NewBigFloat(-74),
			b:    value.NewBigFloat(-6.25),
			want: value.NewBigFloat(-5.25),
		},
		"BigFloat +Inf % 5": { // Mod(±Inf, y) = NaN
			a:    value.BigFloatInf(),
			b:    value.NewBigFloat(5),
			want: value.BigFloatNaN(),
		},
		"BigFloat -Inf % 5": { // Mod(±Inf, y) = NaN
			a:    value.BigFloatNegInf(),
			b:    value.NewBigFloat(5),
			want: value.BigFloatNaN(),
		},
		"BigFloat NaN % 625": { // Mod(NaN, y) = NaN
			a:    value.BigFloatNaN(),
			b:    value.NewBigFloat(625),
			want: value.BigFloatNaN(),
		},
		"BigFloat 25 % 0": { // Mod(x, 0) = NaN
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(0),
			want: value.BigFloatNaN(),
		},
		"BigFloat 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    value.NewBigFloat(25),
			b:    value.BigFloatInf(),
			want: value.NewBigFloat(25),
		},
		"BigFloat -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    value.NewBigFloat(-87),
			b:    value.BigFloatNegInf(),
			want: value.NewBigFloat(-87),
		},
		"BigFloat 49 % NaN": { // Mod(x, NaN) = NaN
			a:    value.NewBigFloat(49),
			b:    value.BigFloatNaN(),
			want: value.BigFloatNaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Modulo(tc.b)
			opts := vm.ComparerOptions
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

func TestBigFloat_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigFloat
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
		"Char and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Char('f'),
			err: value.NewError(value.TypeErrorClass, "`Std::Char` can't be coerced into `Std::BigFloat`"),
		},
		"Int64 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Int64(7),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"Float64 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Float64(7),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` can't be coerced into `Std::BigFloat`"),
		},

		"SmallInt 25bf > 3": {
			a:    value.NewBigFloat(25),
			b:    value.SmallInt(3),
			want: value.True,
		},
		"SmallInt 6bf > 18": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(18),
			want: value.False,
		},
		"SmallInt 6bf > 6": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(6),
			want: value.False,
		},
		"SmallInt 6.5bf > 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.SmallInt(6),
			want: value.True,
		},

		"BigInt 25bf > 3": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigInt(3),
			want: value.True,
		},
		"BigInt 6bf > 18": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(18),
			want: value.False,
		},
		"BigInt 6bf > 6": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(6),
			want: value.False,
		},
		"BigInt 6.5bf > 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.NewBigInt(6),
			want: value.True,
		},

		"Float 25bf > 3.0": {
			a:    value.NewBigFloat(25),
			b:    value.Float(3),
			want: value.True,
		},
		"Float 6bf > 18.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(18.5),
			want: value.False,
		},
		"Float 6bf > 6.0": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6),
			want: value.False,
		},
		"Float 6bf > -6.0": {
			a:    value.NewBigFloat(6),
			b:    value.Float(-6),
			want: value.True,
		},
		"Float -6bf > 6.0": {
			a:    value.NewBigFloat(-6),
			b:    value.Float(6),
			want: value.False,
		},
		"Float 6.5bf > 6.0": {
			a:    value.NewBigFloat(6.5),
			b:    value.Float(6),
			want: value.True,
		},
		"Float 6.bf > 6.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6.5),
			want: value.False,
		},
		"Float 6bf > +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 6bf > -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNegInf(),
			want: value.True,
		},
		"Float +Inf > +Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float +Inf > -Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatNegInf(),
			want: value.True,
		},
		"Float -Inf > +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float -Inf > -Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float 6bf > NaN": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float NaN > 6.0": {
			a:    value.BigFloatNaN(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float NaN > NaN": {
			a:    value.BigFloatNaN(),
			b:    value.FloatNaN(),
			want: value.False,
		},

		"BigFloat 25bf > 3.0bf": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(3),
			want: value.True,
		},
		"BigFloat 6bf > 18.5bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(18.5),
			want: value.False,
		},
		"BigFloat 6bf > 6bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat -6bf > 6bf": {
			a:    value.NewBigFloat(-6),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat 6bf > -6bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(-6),
			want: value.True,
		},
		"BigFloat -6bf > -6bf": {
			a:    value.NewBigFloat(-6),
			b:    value.NewBigFloat(-6),
			want: value.False,
		},
		"BigFloat 6.5bf > 6bf": {
			a:    value.NewBigFloat(6.5),
			b:    value.NewBigFloat(6),
			want: value.True,
		},
		"BigFloat 6bf > 6.5bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(6.5),
			want: value.False,
		},
		"BigFloat 6bf > +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 6bf > -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNegInf(),
			want: value.True,
		},
		"BigFloat +Inf > 6bf": {
			a:    value.BigFloatInf(),
			b:    value.NewBigFloat(6),
			want: value.True,
		},
		"BigFloat -Inf > 6bf": {
			a:    value.BigFloatNegInf(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat +Inf > +Inf": {
			a:    value.BigFloatInf(),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat +Inf > -Inf": {
			a:    value.BigFloatInf(),
			b:    value.BigFloatNegInf(),
			want: value.True,
		},
		"BigFloat -Inf > +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat -Inf > -Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat 6bf > NaN": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat NaN > 6bf": {
			a:    value.BigFloatNaN(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat NaN > NaN": {
			a:    value.BigFloatNaN(),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
			opts := vm.ComparerOptions
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
func TestBigFloat_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigFloat
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
		"Char and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Char('f'),
			err: value.NewError(value.TypeErrorClass, "`Std::Char` can't be coerced into `Std::BigFloat`"),
		},
		"Int64 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Int64(7),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"Float64 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Float64(7),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` can't be coerced into `Std::BigFloat`"),
		},

		"SmallInt 25bf >= 3": {
			a:    value.NewBigFloat(25),
			b:    value.SmallInt(3),
			want: value.True,
		},
		"SmallInt 6bf >= 18": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(18),
			want: value.False,
		},
		"SmallInt 6bf >= 6": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(6),
			want: value.True,
		},
		"SmallInt 6.5bf >= 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.SmallInt(6),
			want: value.True,
		},

		"BigInt 25bf >= 3": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigInt(3),
			want: value.True,
		},
		"BigInt 6bf >= 18": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(18),
			want: value.False,
		},
		"BigInt 6bf >= 6": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(6),
			want: value.True,
		},
		"BigInt 6.5bf >= 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.NewBigInt(6),
			want: value.True,
		},

		"Float 25bf >= 3.0": {
			a:    value.NewBigFloat(25),
			b:    value.Float(3),
			want: value.True,
		},
		"Float 6bf >= 18.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(18.5),
			want: value.False,
		},
		"Float 6bf >= 6.0": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6),
			want: value.True,
		},
		"Float 6bf >= -6.0": {
			a:    value.NewBigFloat(6),
			b:    value.Float(-6),
			want: value.True,
		},
		"Float -6bf >= 6.0": {
			a:    value.NewBigFloat(-6),
			b:    value.Float(6),
			want: value.False,
		},
		"Float 6.5bf >= 6.0": {
			a:    value.NewBigFloat(6.5),
			b:    value.Float(6),
			want: value.True,
		},
		"Float 6bf >= 6.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6.5),
			want: value.False,
		},
		"Float 6bf >= +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 6bf >= -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNegInf(),
			want: value.True,
		},
		"Float +Inf >= +Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatInf(),
			want: value.True,
		},
		"Float +Inf >= -Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatNegInf(),
			want: value.True,
		},
		"Float -Inf >= +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float -Inf >= -Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatNegInf(),
			want: value.True,
		},
		"Float 6bf >= NaN": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float NaN >= 6.0": {
			a:    value.BigFloatNaN(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float NaN >= NaN": {
			a:    value.BigFloatNaN(),
			b:    value.FloatNaN(),
			want: value.False,
		},

		"BigFloat 25bf >= 3.0bf": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(3),
			want: value.True,
		},
		"BigFloat 6bf >= 18.5bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(18.5),
			want: value.False,
		},
		"BigFloat 6bf >= 6.0bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(6),
			want: value.True,
		},
		"BigFloat -6bf >= 6.0bf": {
			a:    value.NewBigFloat(-6),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat 6bf >= -6.0bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(-6),
			want: value.True,
		},
		"BigFloat -6bf >= -6.0bf": {
			a:    value.NewBigFloat(-6),
			b:    value.NewBigFloat(-6),
			want: value.True,
		},
		"BigFloat 6.5bf >= 6.0bf": {
			a:    value.NewBigFloat(6.5),
			b:    value.NewBigFloat(6),
			want: value.True,
		},
		"BigFloat 6bf >= 6.5bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(6.5),
			want: value.False,
		},
		"BigFloat 6bf >= +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 6bf >= -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNegInf(),
			want: value.True,
		},
		"BigFloat +Inf >= 6.0": {
			a:    value.BigFloatInf(),
			b:    value.NewBigFloat(6),
			want: value.True,
		},
		"BigFloat -Inf >= 6.0": {
			a:    value.BigFloatNegInf(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat +Inf >= +Inf": {
			a:    value.BigFloatInf(),
			b:    value.BigFloatInf(),
			want: value.True,
		},
		"BigFloat +Inf >= -Inf": {
			a:    value.BigFloatInf(),
			b:    value.BigFloatNegInf(),
			want: value.True,
		},
		"BigFloat -Inf >= +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat -Inf >= -Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatNegInf(),
			want: value.True,
		},
		"BigFloat 6bf >= NaN": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat NaN >= 6bf": {
			a:    value.BigFloatNaN(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat NaN >= NaN": {
			a:    value.BigFloatNaN(),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqual(tc.b)
			opts := vm.ComparerOptions
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

func TestBigFloat_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigFloat
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
		"Char and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Char('f'),
			err: value.NewError(value.TypeErrorClass, "`Std::Char` can't be coerced into `Std::BigFloat`"),
		},
		"Int64 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Int64(7),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"Float64 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Float64(7),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` can't be coerced into `Std::BigFloat`"),
		},

		"SmallInt 25bf < 3": {
			a:    value.NewBigFloat(25),
			b:    value.SmallInt(3),
			want: value.False,
		},
		"SmallInt 6bf < 18": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(18),
			want: value.True,
		},
		"SmallInt 6bf < 6": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(6),
			want: value.False,
		},
		"SmallInt 5.5bf < 6": {
			a:    value.NewBigFloat(5.5),
			b:    value.SmallInt(6),
			want: value.True,
		},

		"BigInt 25bf < 3": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigInt(3),
			want: value.False,
		},
		"BigInt 6bf < 18": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(18),
			want: value.True,
		},
		"BigInt 6bf < 6": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(6),
			want: value.False,
		},
		"BigInt 5.5bf < 6": {
			a:    value.NewBigFloat(5.5),
			b:    value.NewBigInt(6),
			want: value.True,
		},

		"Float 25bf < 3.0": {
			a:    value.NewBigFloat(25),
			b:    value.Float(3),
			want: value.False,
		},
		"Float 6bf < 18.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(18.5),
			want: value.True,
		},
		"Float 6bf < 6.0": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6),
			want: value.False,
		},
		"Float 5.5bf < 6.0": {
			a:    value.NewBigFloat(5.5),
			b:    value.Float(6),
			want: value.True,
		},
		"Float 6bf < 6.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6.5),
			want: value.True,
		},
		"Float 6.3bf < 6.0": {
			a:    value.NewBigFloat(6.3),
			b:    value.Float(6),
			want: value.False,
		},
		"Float 6bf < +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatInf(),
			want: value.True,
		},
		"Float 6bf < -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float +Inf < 6.0": {
			a:    value.BigFloatInf(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float -Inf < 6.0": {
			a:    value.BigFloatNegInf(),
			b:    value.Float(6),
			want: value.True,
		},
		"Float +Inf < +Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float -Inf < +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatInf(),
			want: value.True,
		},
		"Float +Inf < -Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float -Inf < -Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float 6bf < NaN": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float NaN < 6.0": {
			a:    value.BigFloatNaN(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float NaN < NaN": {
			a:    value.BigFloatNaN(),
			b:    value.FloatNaN(),
			want: value.False,
		},

		"BigFloat 25bf < 3.0bf": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(3),
			want: value.False,
		},
		"BigFloat 6bf < 18.5bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(18.5),
			want: value.True,
		},
		"BigFloat 6bf < 6bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat 6bf < +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatInf(),
			want: value.True,
		},
		"BigFloat 6bf < -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat +Inf < +Inf": {
			a:    value.BigFloatInf(),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat -Inf < +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatInf(),
			want: value.True,
		},
		"BigFloat -Inf < -Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat 6bf < NaN": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat NaN < 6bf": {
			a:    value.BigFloatNaN(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat NaN < NaN": {
			a:    value.BigFloatNaN(),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
			opts := vm.ComparerOptions
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
func TestBigFloat_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigFloat
		b    value.Value
		want value.Value
		err  *value.Error
	}{
		"String and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.String("foo"),
			err: value.NewError(value.TypeErrorClass, "`Std::String` can't be coerced into `Std::BigFloat`"),
		},
		"Char and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Char('f'),
			err: value.NewError(value.TypeErrorClass, "`Std::Char` can't be coerced into `Std::BigFloat`"),
		},
		"Int64 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Int64(7),
			err: value.NewError(value.TypeErrorClass, "`Std::Int64` can't be coerced into `Std::BigFloat`"),
		},
		"Float64 and return an error": {
			a:   value.NewBigFloat(5),
			b:   value.Float64(7),
			err: value.NewError(value.TypeErrorClass, "`Std::Float64` can't be coerced into `Std::BigFloat`"),
		},

		"SmallInt 25bf <= 3": {
			a:    value.NewBigFloat(25),
			b:    value.SmallInt(3),
			want: value.False,
		},
		"SmallInt 6bf <= 18": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(18),
			want: value.True,
		},
		"SmallInt 6bf <= 6": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(6),
			want: value.True,
		},
		"SmallInt 6.5bf <= 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.SmallInt(6),
			want: value.False,
		},
		"SmallInt 5.5bf <= 6": {
			a:    value.NewBigFloat(5.5),
			b:    value.SmallInt(6),
			want: value.True,
		},

		"BigInt 25bf <= 3": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigInt(3),
			want: value.False,
		},
		"BigInt 6bf <= 18": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(18),
			want: value.True,
		},
		"BigInt 6bf <= 6": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(6),
			want: value.True,
		},
		"BigInt 6.5bf <= 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.NewBigInt(6),
			want: value.False,
		},
		"BigInt 5.5bf <= 6": {
			a:    value.NewBigFloat(5.5),
			b:    value.NewBigInt(6),
			want: value.True,
		},

		"Float 25bf <= 3.0": {
			a:    value.NewBigFloat(25),
			b:    value.Float(3),
			want: value.False,
		},
		"Float 6bf <= 18.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(18.5),
			want: value.True,
		},
		"Float 6bf <= 6.0": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6),
			want: value.True,
		},
		"Float 5.5bf <= 6.0": {
			a:    value.NewBigFloat(5.5),
			b:    value.Float(6),
			want: value.True,
		},
		"Float 6bf <= 6.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6.5),
			want: value.True,
		},
		"Float 6.3bf <= 6.0": {
			a:    value.NewBigFloat(6.3),
			b:    value.Float(6),
			want: value.False,
		},
		"Float 6bf <= +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatInf(),
			want: value.True,
		},
		"Float 6bf <= -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float +Inf <= 6.0": {
			a:    value.BigFloatInf(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float -Inf <= 6.0": {
			a:    value.BigFloatNegInf(),
			b:    value.Float(6),
			want: value.True,
		},
		"Float +Inf <= +Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatInf(),
			want: value.True,
		},
		"Float -Inf <= +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatInf(),
			want: value.True,
		},
		"Float +Inf <= -Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float -Inf <= -Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatNegInf(),
			want: value.True,
		},
		"Float 6bf <= NaN": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float NaN <= 6.0": {
			a:    value.BigFloatNaN(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float NaN <= NaN": {
			a:    value.BigFloatNaN(),
			b:    value.FloatNaN(),
			want: value.False,
		},

		"BigFloat 25bf <= 3.0bf": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(3),
			want: value.False,
		},
		"BigFloat 6bf <= 18.5bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(18.5),
			want: value.True,
		},
		"BigFloat 6bf <= 6bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(6),
			want: value.True,
		},
		"BigFloat 6bf <= +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatInf(),
			want: value.True,
		},
		"BigFloat 6bf <= -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat +Inf <= +Inf": {
			a:    value.BigFloatInf(),
			b:    value.BigFloatInf(),
			want: value.True,
		},
		"BigFloat -Inf <= +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatInf(),
			want: value.True,
		},
		"BigFloat -Inf <= -Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatNegInf(),
			want: value.True,
		},
		"BigFloat 6bf <= NaN": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat NaN <= 6bf": {
			a:    value.BigFloatNaN(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat NaN <= NaN": {
			a:    value.BigFloatNaN(),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqual(tc.b)
			opts := vm.ComparerOptions
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

func TestBigFloat_Equal(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigFloat
		b    value.Value
		want value.Value
	}{
		"String 5bf == '5'": {
			a:    value.NewBigFloat(5),
			b:    value.String("5"),
			want: value.False,
		},
		"Char 5bf == c'5'": {
			a:    value.NewBigFloat(5),
			b:    value.Char('5'),
			want: value.False,
		},

		"Int64 5bf == 5i64": {
			a:    value.NewBigFloat(5),
			b:    value.Int64(5),
			want: value.True,
		},
		"Int64 5.5bf == 5i64": {
			a:    value.NewBigFloat(5.5),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int64 4bf == 5i64": {
			a:    value.NewBigFloat(4),
			b:    value.Int64(5),
			want: value.False,
		},

		"Int32 5bf == 5i32": {
			a:    value.NewBigFloat(5),
			b:    value.Int32(5),
			want: value.True,
		},
		"Int32 5.5bf == 5i32": {
			a:    value.NewBigFloat(5.5),
			b:    value.Int32(5),
			want: value.False,
		},
		"Int32 4bf == 5i32": {
			a:    value.NewBigFloat(4),
			b:    value.Int32(5),
			want: value.False,
		},

		"Int16 5bf == 5i16": {
			a:    value.NewBigFloat(5),
			b:    value.Int16(5),
			want: value.True,
		},
		"Int16 5.5bf == 5i16": {
			a:    value.NewBigFloat(5.5),
			b:    value.Int16(5),
			want: value.False,
		},
		"Int16 4bf == 5i16": {
			a:    value.NewBigFloat(4),
			b:    value.Int16(5),
			want: value.False,
		},

		"Int8 5bf == 5i8": {
			a:    value.NewBigFloat(5),
			b:    value.Int8(5),
			want: value.True,
		},
		"Int8 5.5bf == 5i8": {
			a:    value.NewBigFloat(5.5),
			b:    value.Int8(5),
			want: value.False,
		},
		"Int8 4bf == 5i8": {
			a:    value.NewBigFloat(4),
			b:    value.Int8(5),
			want: value.False,
		},

		"UInt64 5bf == 5u64": {
			a:    value.NewBigFloat(5),
			b:    value.UInt64(5),
			want: value.True,
		},
		"UInt64 5.5bf == 5u64": {
			a:    value.NewBigFloat(5.5),
			b:    value.UInt64(5),
			want: value.False,
		},
		"UInt64 4bf == 5u64": {
			a:    value.NewBigFloat(4),
			b:    value.UInt64(5),
			want: value.False,
		},

		"UInt32 5bf == 5u32": {
			a:    value.NewBigFloat(5),
			b:    value.UInt32(5),
			want: value.True,
		},
		"UInt32 5.5bf == 5u32": {
			a:    value.NewBigFloat(5.5),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 4bf == 5u32": {
			a:    value.NewBigFloat(4),
			b:    value.UInt32(5),
			want: value.False,
		},

		"UInt16 5bf == 5u16": {
			a:    value.NewBigFloat(5),
			b:    value.UInt16(5),
			want: value.True,
		},
		"UInt16 5.5bf == 5u16": {
			a:    value.NewBigFloat(5.5),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 4bf == 5u16": {
			a:    value.NewBigFloat(4),
			b:    value.UInt16(5),
			want: value.False,
		},

		"UInt8 5bf == 5u8": {
			a:    value.NewBigFloat(5),
			b:    value.UInt8(5),
			want: value.True,
		},
		"UInt8 5.5bf == 5u8": {
			a:    value.NewBigFloat(5.5),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 4bf == 5u8": {
			a:    value.NewBigFloat(4),
			b:    value.UInt8(5),
			want: value.False,
		},

		"Float64 5bf == 5f64": {
			a:    value.NewBigFloat(5),
			b:    value.Float64(5),
			want: value.True,
		},
		"Float64 5.5bf == 5f64": {
			a:    value.NewBigFloat(5.5),
			b:    value.Float64(5),
			want: value.False,
		},
		"Float64 5bf == 5.5f64": {
			a:    value.NewBigFloat(5),
			b:    value.Float64(5.5),
			want: value.False,
		},
		"Float64 5.5bf == 5.5f64": {
			a:    value.NewBigFloat(5.5),
			b:    value.Float64(5.5),
			want: value.True,
		},
		"Float64 4bf == 5f64": {
			a:    value.NewBigFloat(4),
			b:    value.Float64(5),
			want: value.False,
		},

		"Float32 5bf == 5f32": {
			a:    value.NewBigFloat(5),
			b:    value.Float32(5),
			want: value.True,
		},
		"Float32 5.5bf == 5f32": {
			a:    value.NewBigFloat(5.5),
			b:    value.Float32(5),
			want: value.False,
		},
		"Float32 5bf == 5.5f32": {
			a:    value.NewBigFloat(5),
			b:    value.Float32(5.5),
			want: value.False,
		},
		"Float32 5.5bf == 5.5f32": {
			a:    value.NewBigFloat(5.5),
			b:    value.Float32(5.5),
			want: value.True,
		},
		"Float32 4bf == 5f32": {
			a:    value.NewBigFloat(4),
			b:    value.Float32(5),
			want: value.False,
		},

		"SmallInt 25bf == 3": {
			a:    value.NewBigFloat(25),
			b:    value.SmallInt(3),
			want: value.False,
		},
		"SmallInt 6bf == 18": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(18),
			want: value.False,
		},
		"SmallInt 6bf == 6": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(6),
			want: value.True,
		},
		"SmallInt 6.5bf == 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.SmallInt(6),
			want: value.False,
		},

		"BigInt 25bf == 3": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigInt(3),
			want: value.False,
		},
		"BigInt 6bf == 18": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(18),
			want: value.False,
		},
		"BigInt 6bf == 6": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(6),
			want: value.True,
		},
		"BigInt 6.5bf == 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.NewBigInt(6),
			want: value.False,
		},

		"Float 25bf == 3.0": {
			a:    value.NewBigFloat(25),
			b:    value.Float(3),
			want: value.False,
		},
		"Float 6bf == 18.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(18.5),
			want: value.False,
		},
		"Float 6bf == 6": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6),
			want: value.True,
		},
		"Float 6bf == +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 6bf == -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float +Inf == 6.0": {
			a:    value.BigFloatInf(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float -Inf == 6.0": {
			a:    value.BigFloatNegInf(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float +Inf == +Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatInf(),
			want: value.True,
		},
		"Float +Inf == -Inf": {
			a:    value.BigFloatInf(),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float -Inf == +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 6bf == NaN": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNaN(),
			want: value.False,
		},
		"Float NaN == 6.0": {
			a:    value.BigFloatNaN(),
			b:    value.Float(6),
			want: value.False,
		},
		"Float NaN == NaN": {
			a:    value.BigFloatNaN(),
			b:    value.FloatNaN(),
			want: value.False,
		},

		"BigFloat 25bf == 3.0bf": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(3),
			want: value.False,
		},
		"BigFloat 6bf == 18.5bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(18.5),
			want: value.False,
		},
		"BigFloat 6bf == 6bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(6),
			want: value.True,
		},
		"BigFloat 6bf == +Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 6bf == -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat +Inf == 6bf": {
			a:    value.BigFloatInf(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat -Inf == 6bf": {
			a:    value.BigFloatNegInf(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat +Inf == +Inf": {
			a:    value.BigFloatInf(),
			b:    value.BigFloatInf(),
			want: value.True,
		},
		"BigFloat +Inf == -Inf": {
			a:    value.BigFloatInf(),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat -Inf == +Inf": {
			a:    value.BigFloatNegInf(),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 6bf == NaN": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
		"BigFloat NaN == 6bf": {
			a:    value.BigFloatNaN(),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat NaN == NaN": {
			a:    value.BigFloatNaN(),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestBigFloat_StrictEqual(t *testing.T) {
	tests := map[string]struct {
		a    *value.BigFloat
		b    value.Value
		want value.Value
	}{
		"String 5bf === '5'": {
			a:    value.NewBigFloat(5),
			b:    value.String("5"),
			want: value.False,
		},
		"Char 5bf === c'5'": {
			a:    value.NewBigFloat(5),
			b:    value.Char('5'),
			want: value.False,
		},

		"Int64 5bf === 5i64": {
			a:    value.NewBigFloat(5),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int64 5.3bf === 5i64": {
			a:    value.NewBigFloat(5.3),
			b:    value.Int64(5),
			want: value.False,
		},
		"Int64 4bf === 5i64": {
			a:    value.NewBigFloat(4),
			b:    value.Int64(5),
			want: value.False,
		},

		"Int32 5bf === 5i32": {
			a:    value.NewBigFloat(5),
			b:    value.Int32(5),
			want: value.False,
		},
		"Int32 5.2bf === 5i32": {
			a:    value.NewBigFloat(5.2),
			b:    value.Int32(5),
			want: value.False,
		},
		"Int32 4bf === 5i32": {
			a:    value.NewBigFloat(4),
			b:    value.Int32(5),
			want: value.False,
		},

		"Int16 5bf === 5i16": {
			a:    value.NewBigFloat(5),
			b:    value.Int16(5),
			want: value.False,
		},
		"Int16 5.8bf === 5i16": {
			a:    value.NewBigFloat(5.8),
			b:    value.Int16(5),
			want: value.False,
		},
		"Int16 4bf === 5i16": {
			a:    value.NewBigFloat(4),
			b:    value.Int16(5),
			want: value.False,
		},

		"Int8 5bf === 5i8": {
			a:    value.NewBigFloat(5),
			b:    value.Int8(5),
			want: value.False,
		},
		"Int8 4bf === 5i8": {
			a:    value.NewBigFloat(4),
			b:    value.Int8(5),
			want: value.False,
		},

		"UInt64 5bf === 5u64": {
			a:    value.NewBigFloat(5),
			b:    value.UInt64(5),
			want: value.False,
		},
		"UInt64 5.7bf === 5u64": {
			a:    value.NewBigFloat(5.7),
			b:    value.UInt64(5),
			want: value.False,
		},
		"UInt64 4bf === 5u64": {
			a:    value.NewBigFloat(4),
			b:    value.UInt64(5),
			want: value.False,
		},

		"UInt32 5bf === 5u32": {
			a:    value.NewBigFloat(5),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 5.3bf === 5u32": {
			a:    value.NewBigFloat(5.3),
			b:    value.UInt32(5),
			want: value.False,
		},
		"UInt32 4bf === 5u32": {
			a:    value.NewBigFloat(4),
			b:    value.UInt32(5),
			want: value.False,
		},

		"UInt16 5bf === 5u16": {
			a:    value.NewBigFloat(5),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 5.65bf === 5u16": {
			a:    value.NewBigFloat(5.65),
			b:    value.UInt16(5),
			want: value.False,
		},
		"UInt16 4bf === 5u16": {
			a:    value.NewBigFloat(4),
			b:    value.UInt16(5),
			want: value.False,
		},

		"UInt8 5bf === 5u8": {
			a:    value.NewBigFloat(5),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 5.12bf === 5u8": {
			a:    value.NewBigFloat(5.12),
			b:    value.UInt8(5),
			want: value.False,
		},
		"UInt8 4bf === 5u8": {
			a:    value.NewBigFloat(4),
			b:    value.UInt8(5),
			want: value.False,
		},

		"Float64 5bf === 5f64": {
			a:    value.NewBigFloat(5),
			b:    value.Float64(5),
			want: value.False,
		},
		"Float64 5bf === 5.5f64": {
			a:    value.NewBigFloat(5),
			b:    value.Float64(5.5),
			want: value.False,
		},
		"Float64 5.5bf === 5.5f64": {
			a:    value.NewBigFloat(5),
			b:    value.Float64(5.5),
			want: value.False,
		},
		"Float64 4bf === 5f64": {
			a:    value.NewBigFloat(4),
			b:    value.Float64(5),
			want: value.False,
		},

		"Float32 5bf === 5f32": {
			a:    value.NewBigFloat(5),
			b:    value.Float32(5),
			want: value.False,
		},
		"Float32 5bf === 5.5f32": {
			a:    value.NewBigFloat(5),
			b:    value.Float32(5.5),
			want: value.False,
		},
		"Float32 5.5bf === 5.5f32": {
			a:    value.NewBigFloat(5.5),
			b:    value.Float32(5.5),
			want: value.False,
		},
		"Float32 4bf === 5f32": {
			a:    value.NewBigFloat(4),
			b:    value.Float32(5),
			want: value.False,
		},

		"SmallInt 25bf === 3": {
			a:    value.NewBigFloat(25),
			b:    value.SmallInt(3),
			want: value.False,
		},
		"SmallInt 6bf === 18": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(18),
			want: value.False,
		},
		"SmallInt 6bf === 6": {
			a:    value.NewBigFloat(6),
			b:    value.SmallInt(6),
			want: value.False,
		},
		"SmallInt 6.5bf === 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.SmallInt(6),
			want: value.False,
		},

		"BigInt 25bf === 3": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigInt(3),
			want: value.False,
		},
		"BigInt 6bf === 18": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(18),
			want: value.False,
		},
		"BigInt 6bf === 6": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigInt(6),
			want: value.False,
		},
		"BigInt 6.5bf === 6": {
			a:    value.NewBigFloat(6.5),
			b:    value.NewBigInt(6),
			want: value.False,
		},

		"Float 25bf === 3.0": {
			a:    value.NewBigFloat(25),
			b:    value.Float(3),
			want: value.False,
		},
		"Float 6bf === 18.5": {
			a:    value.NewBigFloat(6),
			b:    value.Float(18.5),
			want: value.False,
		},
		"Float 6bf === 6.0": {
			a:    value.NewBigFloat(6),
			b:    value.Float(6),
			want: value.True,
		},
		"Float 27.5bf === 27.5": {
			a:    value.NewBigFloat(27.5),
			b:    value.Float(27.5),
			want: value.True,
		},
		"Float 6.5bf === 6.0": {
			a:    value.NewBigFloat(6.5),
			b:    value.Float(6),
			want: value.False,
		},
		"Float 6bf === Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatInf(),
			want: value.False,
		},
		"Float 6bf === -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNegInf(),
			want: value.False,
		},
		"Float 6bf === NaN": {
			a:    value.NewBigFloat(6),
			b:    value.FloatNaN(),
			want: value.False,
		},

		"BigFloat 25bf === 3bf": {
			a:    value.NewBigFloat(25),
			b:    value.NewBigFloat(3),
			want: value.False,
		},
		"BigFloat 6bf === 18.5bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(18.5),
			want: value.False,
		},
		"BigFloat 6bf === 6bf": {
			a:    value.NewBigFloat(6),
			b:    value.NewBigFloat(6),
			want: value.False,
		},
		"BigFloat 6.5bf === 6.5bf": {
			a:    value.NewBigFloat(6.5),
			b:    value.NewBigFloat(6.5),
			want: value.False,
		},
		"BigFloat 6bf === Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatInf(),
			want: value.False,
		},
		"BigFloat 6bf === -Inf": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNegInf(),
			want: value.False,
		},
		"BigFloat 6bf === NaN": {
			a:    value.NewBigFloat(6),
			b:    value.BigFloatNaN(),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.StrictEqual(tc.b)
			opts := vm.ComparerOptions
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}
