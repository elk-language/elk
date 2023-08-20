package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestObjectInspect(t *testing.T) {
	tests := map[string]struct {
		obj  *Object
		want string
	}{
		"anonymous class and empty ivars": {
			obj:  NewObject(ObjectWithClass(NewClass())),
			want: `<anonymous>{}`,
		},
		"named class and empty ivars": {
			obj:  NewObject(ObjectWithClass(ExceptionClass)),
			want: `Std::Exception{}`,
		},
		"named class and ivars": {
			obj: NewObject(
				ObjectWithClass(ExceptionClass),
				ObjectWithInstanceVariables(
					SimpleSymbolMap{
						SymbolTable.Add("message").Id: String("foo bar!"),
					},
				),
			),
			want: `Std::Exception{ message: "foo bar!" }`,
		},
		"anonymous class and ivars": {
			obj: NewObject(
				ObjectWithClass(NewClass()),
				ObjectWithInstanceVariables(
					SimpleSymbolMap{
						SymbolTable.Add("message").Id: String("foo bar!"),
					},
				),
			),
			want: `<anonymous>{ message: "foo bar!" }`,
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
