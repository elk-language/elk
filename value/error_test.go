package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestErrorf(t *testing.T) {
	tests := map[string]struct {
		class  *value.Class
		format string
		args   []any
		want   *value.Error
	}{
		"format correctly": {
			class:  value.TypeErrorClass,
			format: "%q can't be coerced into %s",
			args:   []any{value.String("foo"), value.Int16Class.PrintableName()},
			want:   value.NewError(value.TypeErrorClass, `"foo" can't be coerced into Std::Int16`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.Errorf(tc.class, tc.format, tc.args...)
			opts := comparer.Comparer
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestErrorError(t *testing.T) {
	tests := map[string]struct {
		err  *value.Error
		want string
	}{
		"format correctly": {
			err:  value.NewError(value.TypeErrorClass, `"foo" can't be coerced into Std::Int16`),
			want: `Std::TypeError: "foo" can't be coerced into Std::Int16`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.err.Error()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
