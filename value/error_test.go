package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestErrorf(t *testing.T) {
	tests := map[string]struct {
		class  *Class
		format string
		args   []any
		want   *Error
	}{
		"format correctly": {
			class:  TypeErrorClass,
			format: "%q can't be coerced into %s",
			args:   []any{String("foo"), Int16Class.PrintableName()},
			want:   NewError(TypeErrorClass, `"foo" can't be coerced into Std::Int16`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := Errorf(tc.class, tc.format, tc.args...)
			opts := []cmp.Option{
				cmp.AllowUnexported(Error{}),
				cmpopts.IgnoreUnexported(Class{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestErrorError(t *testing.T) {
	tests := map[string]struct {
		err  *Error
		want string
	}{
		"format correctly": {
			err:  NewError(TypeErrorClass, `"foo" can't be coerced into Std::Int16`),
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
