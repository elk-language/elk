package value_test

import (
	"regexp"
	"testing"

	"github.com/elk-language/elk/value"
)

func TestObjectInspect(t *testing.T) {
	tests := map[string]struct {
		obj  *value.Object
		want string
	}{
		"anonymous class and empty ivars": {
			obj:  value.NewObject(value.ObjectWithClass(value.NewClass())),
			want: `<anonymous>\{&: 0x[[:xdigit:]]{4,12}\}`,
		},
		"named class and empty ivars": {
			obj:  value.NewObject(value.ObjectWithClass(value.ErrorClass)),
			want: `Std::Error\{&: 0x[[:xdigit:]]{4,12}\}`,
		},
		"named class and ivars": {
			obj: value.NewObject(
				value.ObjectWithClass(value.ErrorClass),
				value.ObjectWithInstanceVariablesByName(
					value.SymbolMap{
						value.ToSymbol("message"): value.Ref(value.String("foo bar!")),
					},
				),
			),
			want: `Std::Error\{&: 0x[[:xdigit:]]{4,12}, message: "foo bar!"\}`,
		},
		"anonymous class and ivars": {
			obj: value.NewObject(
				value.ObjectWithClass(
					value.NewClassWithOptions(
						value.ClassWithIvarIndices(
							value.IvarIndices{
								value.ToSymbol("message"): 0,
							},
						),
					),
				),
				value.ObjectWithInstanceVariablesByName(
					value.SymbolMap{
						value.ToSymbol("message"): value.Ref(value.String("foo bar!")),
					},
				),
			),
			want: `<anonymous>\{&: 0x[[:xdigit:]]{4,12}, message: "foo bar!"\}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.obj.Inspect()
			ok, _ := regexp.MatchString(tc.want, got)
			if !ok {
				t.Fatalf("got %q, expected to match pattern %q", got, tc.want)
			}
		})
	}
}
