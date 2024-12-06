package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestFloatAdd(t *testing.T) {
	tests := map[string]struct {
		left  value.Float
		right value.Value
		want  value.Value
		err   value.Value
	}{
		"Float + Float => Float": {
			left:  2.5,
			right: value.Float(10.2).ToValue(),
			want:  value.Float(12.7).ToValue(),
		},
		"Float + Float NaN => Float NaN": {
			left:  2.5,
			right: value.FloatNaN().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float NaN + Float => Float NaN": {
			left:  value.FloatNaN(),
			right: value.Float(2.5).ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float NaN + Float NaN => Float NaN": {
			left:  value.FloatNaN(),
			right: value.FloatNaN().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float +Inf + Float +Inf => Float +Inf": {
			left:  value.FloatInf(),
			right: value.FloatInf().ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"Float -Inf + Float -Inf => Float -Inf": {
			left:  value.FloatNegInf(),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatNegInf().ToValue(),
		},
		"Float +Inf + Float -Inf => Float NaN": {
			left:  value.FloatInf(),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float + BigFloat => BigFloat": {
			left:  2.5,
			right: value.Ref(value.NewBigFloat(10.2)),
			want:  value.Ref(value.NewBigFloat(12.7)),
		},
		"Float NaN + BigFloat => BigFloat NaN": {
			left:  value.FloatNaN(),
			right: value.Ref(value.NewBigFloat(10.2)),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float + BigFloat NaN => BigFloat NaN": {
			left:  2.5,
			right: value.Ref(value.BigFloatNaN()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float NaN + BigFloat NaN => BigFloat NaN": {
			left:  value.FloatNaN(),
			right: value.Ref(value.BigFloatNaN()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float +Inf + BigFloat -Inf => BigFloat NaN": {
			left:  value.FloatInf(),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float +Inf + BigFloat +Inf => BigFloat +Inf": {
			left:  value.FloatInf(),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.BigFloatInf()),
		},
		"Float -Inf + BigFloat -Inf => BigFloat -Inf": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatNegInf()),
		},
		"Float + SmallInt => Float": {
			left:  2.5,
			right: value.SmallInt(120).ToValue(),
			want:  value.Float(122.5).ToValue(),
		},
		"Float + BigInt => Float": {
			left:  2.5,
			right: value.Ref(value.NewBigInt(120)),
			want:  value.Float(122.5).ToValue(),
		},
		"Float + Int64 => TypeError": {
			left:  2.5,
			right: value.Int64(20).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float + String => TypeError": {
			left:  2.5,
			right: value.Ref(value.String("foo")),
			err:   value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Add(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Log(got.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloatSubtract(t *testing.T) {
	tests := map[string]struct {
		left  value.Float
		right value.Value
		want  value.Value
		err   value.Value
	}{
		"Float - Float => Float": {
			left:  10.0,
			right: value.Float(5.5).ToValue(),
			want:  value.Float(4.5).ToValue(),
		},
		"Float - BigFloat => BigFloat": {
			left:  12.5,
			right: value.Ref(value.NewBigFloat(2.5)),
			want:  value.Ref(value.NewBigFloat(10.0)),
		},
		"Float - SmallInt => Float": {
			left:  12.5,
			right: value.SmallInt(2).ToValue(),
			want:  value.Float(10.5).ToValue(),
		},
		"Float - BigInt => Float": {
			left:  2.5,
			right: value.Ref(value.NewBigInt(2)),
			want:  value.Float(.5).ToValue(),
		},
		"Float - Int64 => TypeError": {
			left:  2.5,
			right: value.Int64(2).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float - String => TypeError": {
			left:  2.5,
			right: value.Ref(value.String("foo")),
			err:   value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},

		"Float - Float NaN => Float NaN": {
			left:  2.5,
			right: value.FloatNaN().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float NaN - Float => Float NaN": {
			left:  value.FloatNaN(),
			right: value.Float(2.5).ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float NaN - Float NaN => Float NaN": {
			left:  value.FloatNaN(),
			right: value.FloatNaN().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float +Inf - Float +Inf => Float NaN": {
			left:  value.FloatInf(),
			right: value.FloatInf().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float -Inf - Float -Inf => Float NaN": {
			left:  value.FloatNegInf(),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float +Inf - Float -Inf => Float +Inf": {
			left:  value.FloatInf(),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatInf().ToValue(),
		},

		"Float - BigFloat NaN => BigFloat NaN": {
			left:  2.5,
			right: value.Ref(value.BigFloatNaN()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float NaN - BigFloat => BigFloat NaN": {
			left:  value.FloatNaN(),
			right: value.Ref(value.NewBigFloat(2.5)),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float NaN - BigFloat NaN => BigFloat NaN": {
			left:  value.FloatNaN(),
			right: value.Ref(value.BigFloatNaN()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float +Inf - BigFloat +Inf => BigFloat NaN": {
			left:  value.FloatInf(),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float -Inf - BigFloat -Inf => BigFloat NaN": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float +Inf - BigFloat -Inf => BigFloat +Inf": {
			left:  value.FloatInf(),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatInf()),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Subtract(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloatMultiply(t *testing.T) {
	tests := map[string]struct {
		left  value.Float
		right value.Value
		want  value.Value
		err   value.Value
	}{
		"Float * Float => Float": {
			left:  2.55,
			right: value.Float(10.0).ToValue(),
			want:  value.Float(25.5).ToValue(),
		},
		"Float * BigFloat => BigFloat": {
			left:  2.55,
			right: value.Ref(value.NewBigFloat(10.0)),
			want:  value.Ref(value.NewBigFloat(25.5)),
		},
		"Float * SmallInt => Float": {
			left:  2.55,
			right: value.SmallInt(20).ToValue(),
			want:  value.Float(51).ToValue(),
		},
		"Float * BigInt => Float": {
			left:  2.55,
			right: value.Ref(value.NewBigInt(20)),
			want:  value.Float(51).ToValue(),
		},
		"Float * Int64 => TypeError": {
			left:  2.5,
			right: value.Int64(20).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float * String => TypeError": {
			left:  2.5,
			right: value.Ref(value.String("foo")),
			err:   value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},

		"Float * BigFloat NaN => BigFloat NaN": {
			left:  value.Float(2.5),
			right: value.Ref(value.BigFloatNaN()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float NaN * BigFloat => BigFloat NaN": {
			left:  value.FloatNaN(),
			right: value.Ref(value.NewBigFloat(10.2)),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float NaN * BigFloat NaN => BigFloat NaN": {
			left:  value.FloatNaN(),
			right: value.Ref(value.BigFloatNaN()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float +Inf * BigFloat => BigFloat +Inf": {
			left:  value.FloatInf(),
			right: value.Ref(value.NewBigFloat(10.2)),
			want:  value.Ref(value.BigFloatInf()),
		},
		"Float * BigFloat +Inf => BigFloat +Inf": {
			left:  value.Float(10.2),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.BigFloatInf()),
		},
		"Float +Inf * BigFloat +Inf => BigFloat +Inf": {
			left:  value.FloatInf(),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.BigFloatInf()),
		},
		"Float -Inf * +BigFloat => BigFloat -Inf": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.NewBigFloat(10.2)),
			want:  value.Ref(value.BigFloatNegInf()),
		},
		"Float -Inf * -BigFloat => BigFloat +Inf": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.NewBigFloat(-10.2)),
			want:  value.Ref(value.BigFloatInf()),
		},
		"+Float * BigFloat -Inf => BigFloat -Inf": {
			left:  value.Float(10.2),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatNegInf()),
		},
		"-Float * BigFloat -Inf => BigFloat +Inf": {
			left:  value.Float(-10.2),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatInf()),
		},
		"Float -Inf * BigFloat -Inf => BigFloat +Inf": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatInf()),
		},
		"Float +Inf * BigFloat -Inf => BigFloat -Inf": {
			left:  value.FloatInf(),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatNegInf()),
		},
		"Float -Inf * BigFloat +Inf => BigFloat -Inf": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.BigFloatNegInf()),
		},
		"Float -Inf * BigFloat 0 => BigFloat NaN": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.NewBigFloat(0)),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float 0 * BigFloat +Inf => BigFloat NaN": {
			left:  value.Float(0),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.BigFloatNaN()),
		},

		"Float * Float NaN => Float NaN": {
			left:  value.Float(2.5),
			right: value.FloatNaN().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float NaN * Float => Float NaN": {
			left:  value.FloatNaN(),
			right: value.Float(10.2).ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float NaN * Float NaN => Float NaN": {
			left:  value.FloatNaN(),
			right: value.FloatNaN().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float +Inf * Float => Float +Inf": {
			left:  value.FloatInf(),
			right: value.Float(10.2).ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"Float * Float +Inf => Float +Inf": {
			left:  value.Float(10.2),
			right: value.FloatInf().ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"Float +Inf * Float +Inf => Float +Inf": {
			left:  value.FloatInf(),
			right: value.FloatInf().ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"Float -Inf * +Float => Float -Inf": {
			left:  value.FloatNegInf(),
			right: value.Float(10.2).ToValue(),
			want:  value.FloatNegInf().ToValue(),
		},
		"Float -Inf * -Float => Float +Inf": {
			left:  value.FloatNegInf(),
			right: value.Float(-10.2).ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"+Float * Float -Inf => Float -Inf": {
			left:  value.Float(10.2),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatNegInf().ToValue(),
		},
		"-Float * Float -Inf => Float +Inf": {
			left:  value.Float(-10.2),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"Float -Inf * Float -Inf => Float +Inf": {
			left:  value.FloatNegInf(),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"Float +Inf * Float -Inf => Float -Inf": {
			left:  value.FloatInf(),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatNegInf().ToValue(),
		},
		"Float -Inf * Float +Inf => Float -Inf": {
			left:  value.FloatNegInf(),
			right: value.FloatInf().ToValue(),
			want:  value.FloatNegInf().ToValue(),
		},
		"Float -Inf * Float 0 => Float NaN": {
			left:  value.FloatNegInf(),
			right: value.Float(0).ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float 0 * Float +Inf => Float NaN": {
			left:  value.Float(0),
			right: value.FloatInf().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Multiply(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloatDivide(t *testing.T) {
	tests := map[string]struct {
		left  value.Float
		right value.Value
		want  value.Value
		err   value.Value
	}{
		"Float / Float => Float": {
			left:  2.68,
			right: value.Float(2.0).ToValue(),
			want:  value.Float(1.34).ToValue(),
		},
		"Float / BigFloat => BigFloat": {
			left:  2.68,
			right: value.Ref(value.NewBigFloat(2.0)),
			want:  value.Ref(value.NewBigFloat(1.34)),
		},
		"Float / SmallInt => Float": {
			left:  2.68,
			right: value.SmallInt(2).ToValue(),
			want:  value.Float(1.34).ToValue(),
		},
		"Float / BigInt => Float": {
			left:  2.68,
			right: value.Ref(value.NewBigInt(2)),
			want:  value.Float(1.34).ToValue(),
		},
		"Float / Int64 => TypeError": {
			left:  2.5,
			right: value.Int64(20).ToValue(),
			err:   value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float / String => TypeError": {
			left:  2.5,
			right: value.Ref(value.String("foo")),
			err:   value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},

		"Float / BigFloat NaN => BigFloat NaN": {
			left:  value.Float(2.5),
			right: value.Ref(value.BigFloatNaN()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float NaN / BigFloat => BigFloat NaN": {
			left:  value.FloatNaN(),
			right: value.Ref(value.NewBigFloat(10.2)),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float NaN / BigFloat NaN => BigFloat NaN": {
			left:  value.FloatNaN(),
			right: value.Ref(value.BigFloatNaN()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float +Inf / BigFloat => BigFloat +Inf": {
			left:  value.FloatInf(),
			right: value.Ref(value.NewBigFloat(10.2)),
			want:  value.Ref(value.BigFloatInf()),
		},
		"Float / BigFloat +Inf => BigFloat 0": {
			left:  value.Float(10.2),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.NewBigFloat(0)),
		},
		"Float +Inf / BigFloat +Inf => BigFloat NaN": {
			left:  value.FloatInf(),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float -Inf / +BigFloat => BigFloat -Inf": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.NewBigFloat(10.2)),
			want:  value.Ref(value.BigFloatNegInf()),
		},
		"Float -Inf / -BigFloat => BigFloat +Inf": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.NewBigFloat(-10.2)),
			want:  value.Ref(value.BigFloatInf()),
		},
		"+Float / BigFloat -Inf => BigFloat -0": {
			left:  value.Float(10.2),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.NewBigFloat(-0)),
		},
		"-Float / BigFloat -Inf => BigFloat +0": {
			left:  value.Float(-10.2),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.NewBigFloat(0)),
		},
		"Float -Inf / BigFloat -Inf => BigFloat NaN": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float +Inf / BigFloat -Inf => BigFloat NaN": {
			left:  value.FloatInf(),
			right: value.Ref(value.BigFloatNegInf()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float -Inf / BigFloat +Inf => BigFloat NaN": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.BigFloatNaN()),
		},
		"Float -Inf / BigFloat 0 => BigFloat -Inf": {
			left:  value.FloatNegInf(),
			right: value.Ref(value.NewBigFloat(0)),
			want:  value.Ref(value.BigFloatNegInf()),
		},
		"Float 0 / BigFloat +Inf => BigFloat 0": {
			left:  value.Float(0),
			right: value.Ref(value.BigFloatInf()),
			want:  value.Ref(value.NewBigFloat(0)),
		},

		"Float / Float NaN => Float NaN": {
			left:  value.Float(2.5),
			right: value.FloatNaN().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float NaN / Float => Float NaN": {
			left:  value.FloatNaN(),
			right: value.Float(10.2).ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float NaN / Float NaN => Float NaN": {
			left:  value.FloatNaN(),
			right: value.FloatNaN().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float +Inf / Float => Float +Inf": {
			left:  value.FloatInf(),
			right: value.Float(10.2).ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"Float / Float +Inf => Float 0": {
			left:  value.Float(10.2),
			right: value.FloatInf().ToValue(),
			want:  value.Float(0).ToValue(),
		},
		"Float +Inf / Float +Inf => Float NaN": {
			left:  value.FloatInf(),
			right: value.FloatInf().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float -Inf / +Float => Float -Inf": {
			left:  value.FloatNegInf(),
			right: value.Float(10.2).ToValue(),
			want:  value.FloatNegInf().ToValue(),
		},
		"Float -Inf / -Float => Float +Inf": {
			left:  value.FloatNegInf(),
			right: value.Float(-10.2).ToValue(),
			want:  value.FloatInf().ToValue(),
		},
		"+Float / Float -Inf => Float -0": {
			left:  value.Float(10.2),
			right: value.FloatNegInf().ToValue(),
			want:  value.Float(-0).ToValue(),
		},
		"-Float / Float -Inf => Float +0": {
			left:  value.Float(-10.2),
			right: value.FloatNegInf().ToValue(),
			want:  value.Float(0).ToValue(),
		},
		"Float -Inf / Float -Inf => Float NaN": {
			left:  value.FloatNegInf(),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float +Inf / Float -Inf => Float NaN": {
			left:  value.FloatInf(),
			right: value.FloatNegInf().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float -Inf / Float +Inf => Float NaN": {
			left:  value.FloatNegInf(),
			right: value.FloatInf().ToValue(),
			want:  value.FloatNaN().ToValue(),
		},
		"Float -Inf / Float 0 => Float -Inf": {
			left:  value.FloatNegInf(),
			right: value.Float(0).ToValue(),
			want:  value.FloatNegInf().ToValue(),
		},
		"Float 0 / Float +Inf => Float 0": {
			left:  value.Float(0),
			right: value.FloatInf().ToValue(),
			want:  value.Float(0).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.left.Divide(tc.right)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloat_Exponentiate(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"exponentiate String and return an error": {
			a:   value.Float(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},
		"exponentiate Int32 and return an error": {
			a:   value.Float(5),
			b:   value.Int32(2).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int32` cannot be coerced into `Std::Float`")),
		},
		"SmallInt 5 ** 2": {
			a:    value.Float(5),
			b:    value.SmallInt(2).ToValue(),
			want: value.Float(25).ToValue(),
		},
		"SmallInt 7 ** 8": {
			a:    value.Float(7),
			b:    value.SmallInt(8).ToValue(),
			want: value.Float(5764801).ToValue(),
		},
		"SmallInt 2.5 ** 5": {
			a:    value.Float(2.5),
			b:    value.SmallInt(5).ToValue(),
			want: value.Float(97.65625).ToValue(),
		},
		"SmallInt 7.12 ** 1": {
			a:    value.Float(7.12),
			b:    value.SmallInt(1).ToValue(),
			want: value.Float(7.12).ToValue(),
		},
		"SmallInt 4 ** -2": {
			a:    value.Float(4),
			b:    value.SmallInt(-2).ToValue(),
			want: value.Float(0.0625).ToValue(),
		},
		"SmallInt 25 ** 0": {
			a:    value.Float(25),
			b:    value.SmallInt(0).ToValue(),
			want: value.Float(1).ToValue(),
		},

		"BigInt 5 ** 2": {
			a:    value.Float(5),
			b:    value.Ref(value.NewBigInt(2)),
			want: value.Float(25).ToValue(),
		},
		"BigInt 7 ** 8": {
			a:    value.Float(7),
			b:    value.Ref(value.NewBigInt(8)),
			want: value.Float(5764801).ToValue(),
		},
		"BigInt 2.5 ** 5": {
			a:    value.Float(2.5),
			b:    value.Ref(value.NewBigInt(5)),
			want: value.Float(97.65625).ToValue(),
		},
		"BigInt 7.12 ** 1": {
			a:    value.Float(7.12),
			b:    value.Ref(value.NewBigInt(1)),
			want: value.Float(7.12).ToValue(),
		},
		"BigInt 4 ** -2": {
			a:    value.Float(4),
			b:    value.Ref(value.NewBigInt(-2)),
			want: value.Float(0.0625).ToValue(),
		},
		"BigInt 25 ** 0": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(0)),
			want: value.Float(1).ToValue(),
		},

		"Float 5 ** 2": {
			a:    value.Float(5),
			b:    value.Float(2).ToValue(),
			want: value.Float(25).ToValue(),
		},
		"Float 7 ** 8": {
			a:    value.Float(7),
			b:    value.Float(8).ToValue(),
			want: value.Float(5764801).ToValue(),
		},
		"Float 2.5 ** 2.5": {
			a:    value.Float(2.5),
			b:    value.Float(2.5).ToValue(),
			want: value.Float(9.882117688026186).ToValue(),
		},
		"Float 3 ** 2.5": {
			a:    value.Float(3),
			b:    value.Float(2.5).ToValue(),
			want: value.Float(15.588457268119894).ToValue(),
		},
		"Float 6 ** 1": {
			a:    value.Float(6),
			b:    value.Float(1).ToValue(),
			want: value.Float(6).ToValue(),
		},
		"Float 4 ** -2": {
			a:    value.Float(4),
			b:    value.Float(-2).ToValue(),
			want: value.Float(0.0625).ToValue(),
		},
		"Float 25 ** 0": {
			a:    value.Float(25),
			b:    value.Float(0).ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 25 ** NaN": {
			a:    value.Float(25),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float NaN ** 25": {
			a:    value.FloatNaN(),
			b:    value.Float(25).ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float NaN ** NaN": {
			a:    value.FloatNaN(),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    value.Float(0),
			b:    value.Float(-5).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    value.Float(0),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    value.Float(0),
			b:    value.FloatInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    value.Float(0),
			b:    value.Float(-8).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    value.Float(0),
			b:    value.Float(7).ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    value.Float(0),
			b:    value.Float(8).ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    value.Float(-1),
			b:    value.FloatInf().ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    value.Float(-1),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.Float(2),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.Float(-2),
			b:    value.FloatInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.Float(2),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.Float(-2),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    value.Float(0.5),
			b:    value.FloatInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    value.Float(-0.5),
			b:    value.FloatInf().ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float 0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    value.Float(0.5),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float -0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    value.Float(-0.5),
			b:    value.FloatNegInf().ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float +Inf ** 5": { // Pow(+Inf, y) = +Inf for y > 0
			a:    value.FloatInf(),
			b:    value.Float(5).ToValue(),
			want: value.FloatInf().ToValue(),
		},
		"Float +Inf ** -7": { // Pow(+Inf, y) = +0 for y < 0
			a:    value.FloatInf(),
			b:    value.Float(-7).ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -Inf ** -7": {
			a:    value.FloatNegInf(),
			b:    value.Float(-7).ToValue(),
			want: value.Float(0).ToValue(),
		},
		"Float -5.5 ** 3.8": { // Pow(x, y) = NaN for finite x < 0 and finite non-integer y
			a:    value.Float(-5.5),
			b:    value.Float(3.8).ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"BigFloat 5 ** 2": {
			a:    value.Float(5),
			b:    value.Ref(value.NewBigFloat(2)),
			want: value.Ref(value.NewBigFloat(25).SetPrecision(53)),
		},
		"BigFloat 7 ** 8": {
			a:    value.Float(7),
			b:    value.Ref(value.NewBigFloat(8)),
			want: value.Ref(value.NewBigFloat(5764801).SetPrecision(53)),
		},
		"BigFloat 2.5 ** 2.5": {
			a:    value.Float(2.5),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.ParseBigFloatPanic("9.882117688026186").SetPrecision(53)),
		},
		"BigFloat 3 ** 2.5": {
			a:    value.Float(3),
			b:    value.Ref(value.NewBigFloat(2.5)),
			want: value.Ref(value.ParseBigFloatPanic("15.5884572681198956415").SetPrecision(53)),
		},
		"BigFloat 6 ** 1": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(1)),
			want: value.Ref(value.NewBigFloat(6).SetPrecision(53)),
		},
		"BigFloat 4 ** -2": {
			a:    value.Float(4),
			b:    value.Ref(value.NewBigFloat(-2)),
			want: value.Ref(value.NewBigFloat(0.0625).SetPrecision(53)),
		},
		"BigFloat 25 ** 0": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.NewBigFloat(1).SetPrecision(53)),
		},
		"BigFloat 25 ** NaN": {
			a:    value.Float(25),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat NaN ** 25": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.NewBigFloat(25)),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat NaN ** NaN": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat 0 ** -5": { // Pow(±0, y) = ±Inf for y an odd integer < 0
			a:    value.Float(0),
			b:    value.Ref(value.NewBigFloat(-5)),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** -Inf": { // Pow(±0, -Inf) = +Inf
			a:    value.Float(0),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** +Inf": { // Pow(±0, +Inf) = +0
			a:    value.Float(0),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat 0 ** -8": { // Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
			a:    value.Float(0),
			b:    value.Ref(value.NewBigFloat(-8)),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 0 ** 7": { // Pow(±0, y) = ±0 for y an odd integer > 0
			a:    value.Float(0),
			b:    value.Ref(value.NewBigFloat(7)),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat 0 ** 8": { // Pow(±0, y) = +0 for finite y > 0 and not an odd integer
			a:    value.Float(0),
			b:    value.Ref(value.NewBigFloat(8)),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat -1 ** +Inf": { // Pow(-1, ±Inf) = 1
			a:    value.Float(-1),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(1)),
		},
		"BigFloat -1 ** -Inf": { // Pow(-1, ±Inf) = 1
			a:    value.Float(-1),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(1)),
		},
		"BigFloat 2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.Float(2),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat -2 ** +Inf": { // Pow(x, +Inf) = +Inf for |x| > 1
			a:    value.Float(-2),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat 2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.Float(2),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat -2 ** -Inf": { // Pow(x, -Inf) = +0 for |x| > 1
			a:    value.Float(-2),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat 0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    value.Float(0.5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat -0.5 ** +Inf": { // Pow(x, +Inf) = +0 for |x| < 1
			a:    value.Float(-0.5),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat 0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    value.Float(0.5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat -0.5 ** -Inf": { // Pow(x, -Inf) = +Inf for |x| < 1
			a:    value.Float(-0.5),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat +Inf ** 5": { // Pow(+Inf, y) = +Inf for y > 0
			a:    value.FloatInf(),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.Ref(value.BigFloatInf()),
		},
		"BigFloat +Inf ** -7": { // Pow(+Inf, y) = +0 for y < 0
			a:    value.FloatInf(),
			b:    value.Ref(value.NewBigFloat(-7)),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat -Inf ** -7": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.NewBigFloat(-7)),
			want: value.Ref(value.NewBigFloat(0)),
		},
		"BigFloat -5.5 ** 3.8": { // Pow(x, y) = NaN for finite x < 0 and finite non-integer y
			a:    value.Float(-5.5),
			b:    value.Ref(value.NewBigFloat(3.8)),
			want: value.Ref(value.BigFloatNaN()),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Exponentiate(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloat_Modulo(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Float(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},

		"SmallInt 25 % 3": {
			a:    value.Float(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.Float(1).ToValue(),
		},
		"SmallInt 25.6 % 3": {
			a:    value.Float(25.6),
			b:    value.SmallInt(3).ToValue(),
			want: value.Float(1.6000000000000014).ToValue(),
		},
		"SmallInt 76 % 6": {
			a:    value.Float(76),
			b:    value.SmallInt(6).ToValue(),
			want: value.Float(4).ToValue(),
		},
		"SmallInt -76 % 6": {
			a:    value.Float(-76),
			b:    value.SmallInt(6).ToValue(),
			want: value.Float(-4).ToValue(),
		},
		"SmallInt 76 % -6": {
			a:    value.Float(76),
			b:    value.SmallInt(-6).ToValue(),
			want: value.Float(4).ToValue(),
		},
		"SmallInt -76 % -6": {
			a:    value.Float(-76),
			b:    value.SmallInt(-6).ToValue(),
			want: value.Float(-4).ToValue(),
		},
		"SmallInt 124 % 9": {
			a:    value.Float(124),
			b:    value.SmallInt(9).ToValue(),
			want: value.Float(7).ToValue(),
		},

		"BigInt 25 % 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.Float(1).ToValue(),
		},
		"BigInt 76 % 6": {
			a:    value.Float(76),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.Float(4).ToValue(),
		},
		"BigInt 76.5 % 6": {
			a:    value.Float(76.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.Float(4.5).ToValue(),
		},
		"BigInt -76 % 6": {
			a:    value.Float(-76),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.Float(-4).ToValue(),
		},
		"BigInt 76 % -6": {
			a:    value.Float(76),
			b:    value.Ref(value.NewBigInt(-6)),
			want: value.Float(4).ToValue(),
		},
		"BigInt -76 % -6": {
			a:    value.Float(-76),
			b:    value.Ref(value.NewBigInt(-6)),
			want: value.Float(-4).ToValue(),
		},
		"BigInt 124 % 9": {
			a:    value.Float(124),
			b:    value.Ref(value.NewBigInt(9)),
			want: value.Float(7).ToValue(),
		},
		"BigInt 9765 % 9223372036854775808": {
			a:    value.Float(9765),
			b:    value.Ref(value.ParseBigIntPanic("9223372036854775808", 10)),
			want: value.Float(9765).ToValue(),
		},

		"Float 25 % 3": {
			a:    value.Float(25),
			b:    value.Float(3).ToValue(),
			want: value.Float(1).ToValue(),
		},
		"Float 76 % 6": {
			a:    value.Float(76),
			b:    value.Float(6).ToValue(),
			want: value.Float(4).ToValue(),
		},
		"Float 124 % 9": {
			a:    value.Float(124),
			b:    value.Float(9).ToValue(),
			want: value.Float(7).ToValue(),
		},
		"Float 124 % +Inf": {
			a:    value.Float(124),
			b:    value.FloatInf().ToValue(),
			want: value.Float(124).ToValue(),
		},
		"Float 124 % -Inf": {
			a:    value.Float(124),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(124).ToValue(),
		},
		"Float 74.5 % 6.25": {
			a:    value.Float(74.5),
			b:    value.Float(6.25).ToValue(),
			want: value.Float(5.75).ToValue(),
		},
		"Float 74 % 6.25": {
			a:    value.Float(74),
			b:    value.Float(6.25).ToValue(),
			want: value.Float(5.25).ToValue(),
		},
		"Float -74 % 6.25": {
			a:    value.Float(-74),
			b:    value.Float(6.25).ToValue(),
			want: value.Float(-5.25).ToValue(),
		},
		"Float 74 % -6.25": {
			a:    value.Float(74),
			b:    value.Float(-6.25).ToValue(),
			want: value.Float(5.25).ToValue(),
		},
		"Float -74 % -6.25": {
			a:    value.Float(-74),
			b:    value.Float(-6.25).ToValue(),
			want: value.Float(-5.25).ToValue(),
		},
		"Float +Inf % 5": { // Mod(±Inf, y) = NaN
			a:    value.FloatInf(),
			b:    value.Float(5).ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float -Inf % 5": { // Mod(±Inf, y) = NaN
			a:    value.FloatNegInf(),
			b:    value.Float(5).ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float NaN % 625": { // Mod(NaN, y) = NaN
			a:    value.FloatNaN(),
			b:    value.Float(625).ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float 25 % 0": { // Mod(x, 0) = NaN
			a:    value.Float(25),
			b:    value.Float(0).ToValue(),
			want: value.FloatNaN().ToValue(),
		},
		"Float 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    value.Float(25),
			b:    value.FloatInf().ToValue(),
			want: value.Float(25).ToValue(),
		},
		"Float -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    value.Float(-87),
			b:    value.FloatNegInf().ToValue(),
			want: value.Float(-87).ToValue(),
		},
		"Float 49 % NaN": { // Mod(x, NaN) = NaN
			a:    value.Float(49),
			b:    value.FloatNaN().ToValue(),
			want: value.FloatNaN().ToValue(),
		},

		"BigFloat 25 % 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.Ref(value.NewBigFloat(1)),
		},
		"BigFloat 76 % 6": {
			a:    value.Float(76),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.Ref(value.NewBigFloat(4)),
		},
		"BigFloat 124 % 9": {
			a:    value.Float(124),
			b:    value.Ref(value.NewBigFloat(9)),
			want: value.Ref(value.NewBigFloat(7)),
		},
		"BigFloat 74 % 6.25": {
			a:    value.Float(74),
			b:    value.Ref(value.NewBigFloat(6.25)),
			want: value.Ref(value.NewBigFloat(5.25)),
		},
		"BigFloat 74 % 6.25 with higher precision": {
			a:    value.Float(74),
			b:    value.Ref(value.NewBigFloat(6.25).SetPrecision(64)),
			want: value.Ref(value.NewBigFloat(5.25).SetPrecision(64)),
		},
		"BigFloat -74 % 6.25": {
			a:    value.Float(-74),
			b:    value.Ref(value.NewBigFloat(6.25)),
			want: value.Ref(value.NewBigFloat(-5.25)),
		},
		"BigFloat 74 % -6.25": {
			a:    value.Float(74),
			b:    value.Ref(value.NewBigFloat(-6.25)),
			want: value.Ref(value.NewBigFloat(5.25)),
		},
		"BigFloat -74 % -6.25": {
			a:    value.Float(-74),
			b:    value.Ref(value.NewBigFloat(-6.25)),
			want: value.Ref(value.NewBigFloat(-5.25)),
		},
		"BigFloat +Inf % 5": { // Mod(±Inf, y) = NaN
			a:    value.FloatInf(),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat -Inf % 5": { // Mod(±Inf, y) = NaN
			a:    value.FloatNegInf(),
			b:    value.Ref(value.NewBigFloat(5)),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat NaN % 625": { // Mod(NaN, y) = NaN
			a:    value.FloatNaN(),
			b:    value.Ref(value.NewBigFloat(625)),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat 25 % 0": { // Mod(x, 0) = NaN
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(0)),
			want: value.Ref(value.BigFloatNaN()),
		},
		"BigFloat 25 % +Inf": { // Mod(x, ±Inf) = x
			a:    value.Float(25),
			b:    value.Ref(value.BigFloatInf()),
			want: value.Ref(value.NewBigFloat(25)),
		},
		"BigFloat -87 % -Inf": { // Mod(x, ±Inf) = x
			a:    value.Float(-87),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.Ref(value.NewBigFloat(-87)),
		},
		"BigFloat 49 % NaN": { // Mod(x, NaN) = NaN
			a:    value.Float(49),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Ref(value.BigFloatNaN()),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Modulo(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloat_Compare(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Float(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},
		"Char and return an error": {
			a:   value.Float(5),
			b:   value.Char('f').ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Char` cannot be coerced into `Std::Float`")),
		},
		"Int64 and return an error": {
			a:   value.Float(5),
			b:   value.Int64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float64 and return an error": {
			a:   value.Float(5),
			b:   value.Float64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Float`")),
		},

		"SmallInt 25.0 <=> 3": {
			a:    value.Float(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"SmallInt 6.0 <=> 18": {
			a:    value.Float(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"SmallInt 6.0 <=> 6": {
			a:    value.Float(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"SmallInt 6.5 <=> 6": {
			a:    value.Float(6.5),
			b:    value.SmallInt(6).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},

		"BigInt 25.0 <=> 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigInt 6.0 <=> 18": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigInt 6.0 <=> 6": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.SmallInt(0).ToValue(),
		},
		"BigInt 6.5 <=> 6": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.SmallInt(1).ToValue(),
		},

		"Float 25.0 <=> 3.0": {
			a:    value.Float(25),
			b:    value.Float(3).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Float 6.0 <=> 18.5": {
			a:    value.Float(6),
			b:    value.Float(18.5).ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Float 6.0 <=> 6.0": {
			a:    value.Float(6),
			b:    value.Float(6).ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"Float 6.0 <=> -6.0": {
			a:    value.Float(6),
			b:    value.Float(-6).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Float -6.0 <=> 6.0": {
			a:    value.Float(-6),
			b:    value.Float(6).ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Float 6.5 <=> 6.0": {
			a:    value.Float(6.5),
			b:    value.Float(6).ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Float 6.0 <=> 6.5": {
			a:    value.Float(6),
			b:    value.Float(6.5).ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Float 6.0 <=> +Inf": {
			a:    value.Float(6),
			b:    value.FloatInf().ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Float 6.0 <=> -Inf": {
			a:    value.Float(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Float +Inf <=> +Inf": {
			a:    value.FloatInf(),
			b:    value.FloatInf().ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"Float +Inf <=> -Inf": {
			a:    value.FloatInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.SmallInt(1).ToValue(),
		},
		"Float -Inf <=> +Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatInf().ToValue(),
			want: value.SmallInt(-1).ToValue(),
		},
		"Float -Inf <=> -Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.SmallInt(0).ToValue(),
		},
		"Float 6.0 <=> NaN": {
			a:    value.Float(6),
			b:    value.FloatNaN().ToValue(),
			want: value.Nil,
		},
		"Float NaN <=> 6.0": {
			a:    value.FloatNaN(),
			b:    value.Float(6).ToValue(),
			want: value.Nil,
		},
		"Float NaN <=> NaN": {
			a:    value.FloatNaN(),
			b:    value.FloatNaN().ToValue(),
			want: value.Nil,
		},

		"BigFloat 25.0 <=> 3.0bf": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigFloat 6.0 <=> 18.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigFloat 6.0 <=> 6.0bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.SmallInt(0).ToValue(),
		},
		"BigFloat -6.0 <=> 6.0bf": {
			a:    value.Float(-6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigFloat 6.0 <=> -6.0bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(-6)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigFloat -6.0 <=> -6.0bf": {
			a:    value.Float(-6),
			b:    value.Ref(value.NewBigFloat(-6)),
			want: value.SmallInt(0).ToValue(),
		},
		"BigFloat 6.5 <=> 6.0bf": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigFloat 6.0 <=> 6.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6.5)),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigFloat 6.0 <=> +Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigFloat 6.0 <=> -Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.SmallInt(1).ToValue(),
		},
		"BigFloat +Inf <=> 6.0": {
			a:    value.FloatInf(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.SmallInt(1).ToValue(),
		},
		"BigFloat -Inf <=> 6.0": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigFloat +Inf <=> +Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.SmallInt(0).ToValue(),
		},
		"BigFloat +Inf <=> -Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.SmallInt(1).ToValue(),
		},
		"BigFloat -Inf <=> +Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.SmallInt(-1).ToValue(),
		},
		"BigFloat -Inf <=> -Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.SmallInt(0).ToValue(),
		},
		"BigFloat 6.0 <=> NaN": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Nil,
		},
		"BigFloat NaN <=> 6.0bf": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.Nil,
		},
		"BigFloat NaN <=> NaN": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.Compare(tc.b)
			opts := comparer.Options()
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

func TestFloat_GreaterThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Float(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},
		"Char and return an error": {
			a:   value.Float(5),
			b:   value.Char('f').ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Char` cannot be coerced into `Std::Float`")),
		},
		"Int64 and return an error": {
			a:   value.Float(5),
			b:   value.Int64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float64 and return an error": {
			a:   value.Float(5),
			b:   value.Float64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Float`")),
		},

		"SmallInt 25.0 > 3": {
			a:    value.Float(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.True,
		},
		"SmallInt 6.0 > 18": {
			a:    value.Float(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6.0 > 6": {
			a:    value.Float(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},
		"SmallInt 6.5 > 6": {
			a:    value.Float(6.5),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25.0 > 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.True,
		},
		"BigInt 6.0 > 18": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6.0 > 6": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},
		"BigInt 6.5 > 6": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25.0 > 3.0": {
			a:    value.Float(25),
			b:    value.Float(3).ToValue(),
			want: value.True,
		},
		"Float 6.0 > 18.5": {
			a:    value.Float(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6.0 > 6.0": {
			a:    value.Float(6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6.0 > -6.0": {
			a:    value.Float(6),
			b:    value.Float(-6).ToValue(),
			want: value.True,
		},
		"Float -6.0 > 6.0": {
			a:    value.Float(-6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6.5 > 6.0": {
			a:    value.Float(6.5),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6.0 > 6.5": {
			a:    value.Float(6),
			b:    value.Float(6.5).ToValue(),
			want: value.False,
		},
		"Float 6.0 > +Inf": {
			a:    value.Float(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6.0 > -Inf": {
			a:    value.Float(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float +Inf > +Inf": {
			a:    value.FloatInf(),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float +Inf > -Inf": {
			a:    value.FloatInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float -Inf > +Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float -Inf > -Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6.0 > NaN": {
			a:    value.Float(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"Float NaN > 6.0": {
			a:    value.FloatNaN(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float NaN > NaN": {
			a:    value.FloatNaN(),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"BigFloat 25.0 > 3.0bf": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.True,
		},
		"BigFloat 6.0 > 18.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6.0 > 6.0bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat -6.0 > 6.0bf": {
			a:    value.Float(-6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6.0 > -6.0bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(-6)),
			want: value.True,
		},
		"BigFloat -6.0 > -6.0bf": {
			a:    value.Float(-6),
			b:    value.Ref(value.NewBigFloat(-6)),
			want: value.False,
		},
		"BigFloat 6.5 > 6.0bf": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6.0 > 6.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6.5)),
			want: value.False,
		},
		"BigFloat 6.0 > +Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6.0 > -Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat +Inf > 6.0": {
			a:    value.FloatInf(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat -Inf > 6.0": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat +Inf > +Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat +Inf > -Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat -Inf > +Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat -Inf > -Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6.0 > NaN": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
		"BigFloat NaN > 6.0bf": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat NaN > NaN": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
func TestFloat_GreaterThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Float(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},
		"Char and return an error": {
			a:   value.Float(5),
			b:   value.Char('f').ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Char` cannot be coerced into `Std::Float`")),
		},
		"Int64 and return an error": {
			a:   value.Float(5),
			b:   value.Int64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float64 and return an error": {
			a:   value.Float(5),
			b:   value.Float64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Float`")),
		},

		"SmallInt 25.0 >= 3": {
			a:    value.Float(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.True,
		},
		"SmallInt 6.0 >= 18": {
			a:    value.Float(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6.0 >= 6": {
			a:    value.Float(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},
		"SmallInt 6.5 >= 6": {
			a:    value.Float(6.5),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25.0 >= 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.True,
		},
		"BigInt 6.0 >= 18": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6.0 >= 6": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},
		"BigInt 6.5 >= 6": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25.0 >= 3.0": {
			a:    value.Float(25),
			b:    value.Float(3).ToValue(),
			want: value.True,
		},
		"Float 6.0 >= 18.5": {
			a:    value.Float(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6.0 >= 6.0": {
			a:    value.Float(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6.0 >= -6.0": {
			a:    value.Float(6),
			b:    value.Float(-6).ToValue(),
			want: value.True,
		},
		"Float -6.0 >= 6.0": {
			a:    value.Float(-6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6.5 >= 6.0": {
			a:    value.Float(6.5),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6.0 >= 6.5": {
			a:    value.Float(6),
			b:    value.Float(6.5).ToValue(),
			want: value.False,
		},
		"Float 6.0 >= +Inf": {
			a:    value.Float(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6.0 >= -Inf": {
			a:    value.Float(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float +Inf >= +Inf": {
			a:    value.FloatInf(),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float +Inf >= -Inf": {
			a:    value.FloatInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float -Inf >= +Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float -Inf >= -Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float 6.0 >= NaN": {
			a:    value.Float(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"Float NaN >= 6.0": {
			a:    value.FloatNaN(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float NaN >= NaN": {
			a:    value.FloatNaN(),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25.0 >= 3.0bf": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.True,
		},
		"BigFloat 6.0 >= 18.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6.0 >= 6.0bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat -6.0 >= 6.0bf": {
			a:    value.Float(-6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6.0 >= -6.0bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(-6)),
			want: value.True,
		},
		"BigFloat -6.0 >= -6.0bf": {
			a:    value.Float(-6),
			b:    value.Ref(value.NewBigFloat(-6)),
			want: value.True,
		},
		"BigFloat 6.5 >= 6.0bf": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6.0 >= 6.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6.5)),
			want: value.False,
		},
		"BigFloat 6.0 >= +Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6.0 >= -Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat +Inf >= 6.0": {
			a:    value.FloatInf(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat -Inf >= 6.0": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat +Inf >= +Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat +Inf >= -Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat -Inf >= +Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat -Inf >= -Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat 6.0 >= NaN": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
		"BigFloat NaN >= 6.0bf": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat NaN >= NaN": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.GreaterThanEqual(tc.b)
			opts := comparer.Options()
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

func TestFloat_LessThan(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Float(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},
		"Char and return an error": {
			a:   value.Float(5),
			b:   value.Char('f').ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Char` cannot be coerced into `Std::Float`")),
		},
		"Int64 and return an error": {
			a:   value.Float(5),
			b:   value.Int64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float64 and return an error": {
			a:   value.Float(5),
			b:   value.Float64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Float`")),
		},

		"SmallInt 25.0 < 3": {
			a:    value.Float(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6.0 < 18": {
			a:    value.Float(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.True,
		},
		"SmallInt 6.0 < 6": {
			a:    value.Float(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},
		"SmallInt 5.5 < 6": {
			a:    value.Float(5.5),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25.0 < 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6.0 < 18": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.True,
		},
		"BigInt 6.0 < 6": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},
		"BigInt 5.5 < 6": {
			a:    value.Float(5.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25.0 < 3.0": {
			a:    value.Float(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6.0 < 18.5": {
			a:    value.Float(6),
			b:    value.Float(18.5).ToValue(),
			want: value.True,
		},
		"Float 6.0 < 6.0": {
			a:    value.Float(6),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 5.5 < 6.0": {
			a:    value.Float(5.5),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6.0 < 6.5": {
			a:    value.Float(6),
			b:    value.Float(6.5).ToValue(),
			want: value.True,
		},
		"Float 6.3 < 6.0": {
			a:    value.Float(6.3),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6.0 < +Inf": {
			a:    value.Float(6),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float 6.0 < -Inf": {
			a:    value.Float(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float +Inf < 6.0": {
			a:    value.FloatInf(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float -Inf < 6.0": {
			a:    value.FloatNegInf(),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float +Inf < +Inf": {
			a:    value.FloatInf(),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float -Inf < +Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float +Inf < -Inf": {
			a:    value.FloatInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float -Inf < -Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6.0 < NaN": {
			a:    value.Float(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"Float NaN < 6.0": {
			a:    value.FloatNaN(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float NaN < NaN": {
			a:    value.FloatNaN(),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25.0 < 3.0bf": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6.0 < 18.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.True,
		},
		"BigFloat 6.0 < 6bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6.0 < +Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat 6.0 < -Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat +Inf < +Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat -Inf < +Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat -Inf < -Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6.0 < NaN": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
		"BigFloat NaN < 6.0bf": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat NaN < NaN": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
func TestFloat_LessThanEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
		err  value.Value
	}{
		"String and return an error": {
			a:   value.Float(5),
			b:   value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::String` cannot be coerced into `Std::Float`")),
		},
		"Char and return an error": {
			a:   value.Float(5),
			b:   value.Char('f').ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Char` cannot be coerced into `Std::Float`")),
		},
		"Int64 and return an error": {
			a:   value.Float(5),
			b:   value.Int64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Int64` cannot be coerced into `Std::Float`")),
		},
		"Float64 and return an error": {
			a:   value.Float(5),
			b:   value.Float64(7).ToValue(),
			err: value.Ref(value.NewError(value.TypeErrorClass, "`Std::Float64` cannot be coerced into `Std::Float`")),
		},

		"SmallInt 25.0 <= 3": {
			a:    value.Float(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6.0 <= 18": {
			a:    value.Float(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.True,
		},
		"SmallInt 6.0 <= 6": {
			a:    value.Float(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},
		"SmallInt 6.5 <= 6": {
			a:    value.Float(6.5),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},
		"SmallInt 5.5 <= 6": {
			a:    value.Float(5.5),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},

		"BigInt 25.0 <= 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6.0 <= 18": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.True,
		},
		"BigInt 6.0 <= 6": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},
		"BigInt 6.5 <= 6": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},
		"BigInt 5.5 <= 6": {
			a:    value.Float(5.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},

		"Float 25.0 <= 3.0": {
			a:    value.Float(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6.0 <= 18.5": {
			a:    value.Float(6),
			b:    value.Float(18.5).ToValue(),
			want: value.True,
		},
		"Float 6.0 <= 6.0": {
			a:    value.Float(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 5.5 <= 6.0": {
			a:    value.Float(5.5),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6.0 <= 6.5": {
			a:    value.Float(6),
			b:    value.Float(6.5).ToValue(),
			want: value.True,
		},
		"Float 6.3 <= 6.0": {
			a:    value.Float(6.3),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6.0 <= +Inf": {
			a:    value.Float(6),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float 6.0 <= -Inf": {
			a:    value.Float(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float +Inf <= 6.0": {
			a:    value.FloatInf(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float -Inf <= 6.0": {
			a:    value.FloatNegInf(),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float +Inf <= +Inf": {
			a:    value.FloatInf(),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float -Inf <= +Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float +Inf <= -Inf": {
			a:    value.FloatInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float -Inf <= -Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.True,
		},
		"Float 6.0 <= NaN": {
			a:    value.Float(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"Float NaN <= 6.0": {
			a:    value.FloatNaN(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float NaN <= NaN": {
			a:    value.FloatNaN(),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25.0 <= 3.0bf": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6.0 <= 18.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.True,
		},
		"BigFloat 6.0 <= 6bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6.0 <= +Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat 6.0 <= -Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat +Inf <= +Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat -Inf <= +Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat -Inf <= -Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.True,
		},
		"BigFloat 6.0 <= NaN": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
		"BigFloat NaN <= 6.0bf": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat NaN <= NaN": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.a.LessThanEqual(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloat_LaxEqual(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
	}{
		"String 5.0 =~ '5'": {
			a:    value.Float(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},
		"Char 5.0 =~ `5`": {
			a:    value.Float(5),
			b:    value.Char('5').ToValue(),
			want: value.False,
		},

		"Int64 5.0 =~ 5i64": {
			a:    value.Float(5),
			b:    value.Int64(5).ToValue(),
			want: value.True,
		},
		"Int64 5.5 =~ 5i64": {
			a:    value.Float(5.5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4.0 =~ 5i64": {
			a:    value.Float(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},

		"Int32 5.0 =~ 5i32": {
			a:    value.Float(5),
			b:    value.Int32(5).ToValue(),
			want: value.True,
		},
		"Int32 5.5 =~ 5i32": {
			a:    value.Float(5.5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4.0 =~ 5i32": {
			a:    value.Float(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"Int16 5.0 =~ 5i16": {
			a:    value.Float(5),
			b:    value.Int16(5).ToValue(),
			want: value.True,
		},
		"Int16 5.5 =~ 5i16": {
			a:    value.Float(5.5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4.0 =~ 5i16": {
			a:    value.Float(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"Int8 5.0 =~ 5i8": {
			a:    value.Float(5),
			b:    value.Int8(5).ToValue(),
			want: value.True,
		},
		"Int8 5.5 =~ 5i8": {
			a:    value.Float(5.5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4.0 =~ 5i8": {
			a:    value.Float(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt64 5.0 =~ 5u64": {
			a:    value.Float(5),
			b:    value.UInt64(5).ToValue(),
			want: value.True,
		},
		"UInt64 5.5 =~ 5u64": {
			a:    value.Float(5.5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4.0 =~ 5u64": {
			a:    value.Float(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"UInt32 5.0 =~ 5u32": {
			a:    value.Float(5),
			b:    value.UInt32(5).ToValue(),
			want: value.True,
		},
		"UInt32 5.5 =~ 5u32": {
			a:    value.Float(5.5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4.0 =~ 5u32": {
			a:    value.Float(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"UInt16 5.0 =~ 5u16": {
			a:    value.Float(5),
			b:    value.UInt16(5).ToValue(),
			want: value.True,
		},
		"UInt16 5.5 =~ 5u16": {
			a:    value.Float(5.5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4.0 =~ 5u16": {
			a:    value.Float(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"UInt8 5.0 =~ 5u8": {
			a:    value.Float(5),
			b:    value.UInt8(5).ToValue(),
			want: value.True,
		},
		"UInt8 5.5 =~ 5u8": {
			a:    value.Float(5.5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4.0 =~ 5u8": {
			a:    value.Float(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5.0 =~ 5f64": {
			a:    value.Float(5),
			b:    value.Float64(5).ToValue(),
			want: value.True,
		},
		"Float64 5.5 =~ 5f64": {
			a:    value.Float(5.5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5.0 =~ 5.5f64": {
			a:    value.Float(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5.5 =~ 5.5f64": {
			a:    value.Float(5.5),
			b:    value.Float64(5.5).ToValue(),
			want: value.True,
		},
		"Float64 4.0 =~ 5f64": {
			a:    value.Float(4),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},

		"Float32 5.0 =~ 5f32": {
			a:    value.Float(5),
			b:    value.Float32(5).ToValue(),
			want: value.True,
		},
		"Float32 5.5 =~ 5f32": {
			a:    value.Float(5.5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5.0 =~ 5.5f32": {
			a:    value.Float(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5.5 =~ 5.5f32": {
			a:    value.Float(5.5),
			b:    value.Float32(5.5).ToValue(),
			want: value.True,
		},
		"Float32 4.0 =~ 5f32": {
			a:    value.Float(4),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},

		"SmallInt 25.0 =~ 3": {
			a:    value.Float(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6.0 =~ 18": {
			a:    value.Float(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6.0 =~ 6": {
			a:    value.Float(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.True,
		},
		"SmallInt 6.5 =~ 6": {
			a:    value.Float(6.5),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},

		"BigInt 25.0 =~ 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6.0 =~ 18": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6.0 =~ 6": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.True,
		},
		"BigInt 6.5 =~ 6": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},

		"Float 25.0 =~ 3.0": {
			a:    value.Float(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6.0 =~ 18.5": {
			a:    value.Float(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6.0 =~ 6": {
			a:    value.Float(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 6.0 =~ +Inf": {
			a:    value.Float(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6.0 =~ -Inf": {
			a:    value.Float(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float +Inf =~ 6.0": {
			a:    value.FloatInf(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float -Inf =~ 6.0": {
			a:    value.FloatNegInf(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float +Inf =~ +Inf": {
			a:    value.FloatInf(),
			b:    value.FloatInf().ToValue(),
			want: value.True,
		},
		"Float +Inf =~ -Inf": {
			a:    value.FloatInf(),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float -Inf =~ +Inf": {
			a:    value.FloatNegInf(),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6.0 =~ NaN": {
			a:    value.Float(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},
		"Float NaN =~ 6.0": {
			a:    value.FloatNaN(),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float NaN =~ NaN": {
			a:    value.FloatNaN(),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25.0 =~ 3.0bf": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6.0 =~ 18.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6.0 =~ 6bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.True,
		},
		"BigFloat 6.0 =~ +Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6.0 =~ -Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat +Inf =~ 6.0bf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat -Inf =~ 6.0bf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat +Inf =~ +Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.True,
		},
		"BigFloat +Inf =~ -Inf": {
			a:    value.FloatInf(),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat -Inf =~ +Inf": {
			a:    value.FloatNegInf(),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6.0 =~ NaN": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
		"BigFloat NaN =~ 6.0bf": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat NaN =~ NaN": {
			a:    value.FloatNaN(),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.LaxEqual(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}

func TestFloat_Equal(t *testing.T) {
	tests := map[string]struct {
		a    value.Float
		b    value.Value
		want value.Value
	}{
		"String 5.0 == '5'": {
			a:    value.Float(5),
			b:    value.Ref(value.String("5")),
			want: value.False,
		},
		"Char 5.0 == `5`": {
			a:    value.Float(5),
			b:    value.Char('5').ToValue(),
			want: value.False,
		},

		"Int64 5.0 == 5i64": {
			a:    value.Float(5),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 5.3 == 5i64": {
			a:    value.Float(5.3),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},
		"Int64 4.0 == 5i64": {
			a:    value.Float(4),
			b:    value.Int64(5).ToValue(),
			want: value.False,
		},

		"Int32 5.0 == 5i32": {
			a:    value.Float(5),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 5.2 == 5i32": {
			a:    value.Float(5.2),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},
		"Int32 4.0 == 5i32": {
			a:    value.Float(4),
			b:    value.Int32(5).ToValue(),
			want: value.False,
		},

		"Int16 5.0 == 5i16": {
			a:    value.Float(5),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 5.8 == 5i16": {
			a:    value.Float(5.8),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},
		"Int16 4.0 == 5i16": {
			a:    value.Float(4),
			b:    value.Int16(5).ToValue(),
			want: value.False,
		},

		"Int8 5.0 == 5i8": {
			a:    value.Float(5),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},
		"Int8 4.0 == 5i8": {
			a:    value.Float(4),
			b:    value.Int8(5).ToValue(),
			want: value.False,
		},

		"UInt64 5.0 == 5u64": {
			a:    value.Float(5),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 5.7 == 5u64": {
			a:    value.Float(5.7),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},
		"UInt64 4.0 == 5u64": {
			a:    value.Float(4),
			b:    value.UInt64(5).ToValue(),
			want: value.False,
		},

		"UInt32 5.0 == 5u32": {
			a:    value.Float(5),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 5.3 == 5u32": {
			a:    value.Float(5.3),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},
		"UInt32 4.0 == 5u32": {
			a:    value.Float(4),
			b:    value.UInt32(5).ToValue(),
			want: value.False,
		},

		"UInt16 5.0 == 5u16": {
			a:    value.Float(5),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 5.65 == 5u16": {
			a:    value.Float(5.65),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},
		"UInt16 4.0 == 5u16": {
			a:    value.Float(4),
			b:    value.UInt16(5).ToValue(),
			want: value.False,
		},

		"UInt8 5.0 == 5u8": {
			a:    value.Float(5),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 5.12 == 5u8": {
			a:    value.Float(5.12),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},
		"UInt8 4.0 == 5u8": {
			a:    value.Float(4),
			b:    value.UInt8(5).ToValue(),
			want: value.False,
		},

		"Float64 5.0 == 5f64": {
			a:    value.Float(5),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},
		"Float64 5.0 == 5.5f64": {
			a:    value.Float(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 5.5 == 5.5f64": {
			a:    value.Float(5),
			b:    value.Float64(5.5).ToValue(),
			want: value.False,
		},
		"Float64 4.0 == 5f64": {
			a:    value.Float(4),
			b:    value.Float64(5).ToValue(),
			want: value.False,
		},

		"Float32 5.0 == 5f32": {
			a:    value.Float(5),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},
		"Float32 5.0 == 5.5f32": {
			a:    value.Float(5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 5.5 == 5.5f32": {
			a:    value.Float(5.5),
			b:    value.Float32(5.5).ToValue(),
			want: value.False,
		},
		"Float32 4.0 == 5f32": {
			a:    value.Float(4),
			b:    value.Float32(5).ToValue(),
			want: value.False,
		},

		"SmallInt 25.0 == 3": {
			a:    value.Float(25),
			b:    value.SmallInt(3).ToValue(),
			want: value.False,
		},
		"SmallInt 6.0 == 18": {
			a:    value.Float(6),
			b:    value.SmallInt(18).ToValue(),
			want: value.False,
		},
		"SmallInt 6.0 == 6": {
			a:    value.Float(6),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},
		"SmallInt 6.5 == 6": {
			a:    value.Float(6.5),
			b:    value.SmallInt(6).ToValue(),
			want: value.False,
		},

		"BigInt 25.0 == 3": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigInt(3)),
			want: value.False,
		},
		"BigInt 6.0 == 18": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(18)),
			want: value.False,
		},
		"BigInt 6.0 == 6": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},
		"BigInt 6.5 == 6": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigInt(6)),
			want: value.False,
		},

		"Float 25.0 == 3.0": {
			a:    value.Float(25),
			b:    value.Float(3).ToValue(),
			want: value.False,
		},
		"Float 6.0 == 18.5": {
			a:    value.Float(6),
			b:    value.Float(18.5).ToValue(),
			want: value.False,
		},
		"Float 6.0 == 6.0": {
			a:    value.Float(6),
			b:    value.Float(6).ToValue(),
			want: value.True,
		},
		"Float 27.5 == 27.5": {
			a:    value.Float(27.5),
			b:    value.Float(27.5).ToValue(),
			want: value.True,
		},
		"Float 6.5 == 6.0": {
			a:    value.Float(6.5),
			b:    value.Float(6).ToValue(),
			want: value.False,
		},
		"Float 6.0 == Inf": {
			a:    value.Float(6),
			b:    value.FloatInf().ToValue(),
			want: value.False,
		},
		"Float 6.0 == -Inf": {
			a:    value.Float(6),
			b:    value.FloatNegInf().ToValue(),
			want: value.False,
		},
		"Float 6.0 == NaN": {
			a:    value.Float(6),
			b:    value.FloatNaN().ToValue(),
			want: value.False,
		},

		"BigFloat 25.0 == 3bf": {
			a:    value.Float(25),
			b:    value.Ref(value.NewBigFloat(3)),
			want: value.False,
		},
		"BigFloat 6.0 == 18.5bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(18.5)),
			want: value.False,
		},
		"BigFloat 6.0 == 6bf": {
			a:    value.Float(6),
			b:    value.Ref(value.NewBigFloat(6)),
			want: value.False,
		},
		"BigFloat 6.5 == 6.5bf": {
			a:    value.Float(6.5),
			b:    value.Ref(value.NewBigFloat(6.5)),
			want: value.False,
		},
		"BigFloat 6.0 == Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatInf()),
			want: value.False,
		},
		"BigFloat 6.0 == -Inf": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNegInf()),
			want: value.False,
		},
		"BigFloat 6.0 == NaN": {
			a:    value.Float(6),
			b:    value.Ref(value.BigFloatNaN()),
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
		})
	}
}
