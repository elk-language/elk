package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestObjectInspect(t *testing.T) {
	tests := map[string]struct {
		obj  *value.Object
		want string
	}{
		"anonymous class and empty ivars": {
			obj:  value.NewObject(value.ObjectWithClass(value.NewClass())),
			want: `<anonymous>{}`,
		},
		"named class and empty ivars": {
			obj:  value.NewObject(value.ObjectWithClass(value.ErrorClass)),
			want: `Std::Error{}`,
		},
		"named class and ivars": {
			obj: value.NewObject(
				value.ObjectWithClass(value.ErrorClass),
				value.ObjectWithInstanceVariables(
					value.SymbolMap{
						value.ToSymbol("message"): value.String("foo bar!"),
					},
				),
			),
			want: `Std::Error{message: "foo bar!"}`,
		},
		"anonymous class and ivars": {
			obj: value.NewObject(
				value.ObjectWithClass(value.NewClass()),
				value.ObjectWithInstanceVariables(
					value.SymbolMap{
						value.ToSymbol("message"): value.String("foo bar!"),
					},
				),
			),
			want: `<anonymous>{message: "foo bar!"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.obj.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
