package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClassInspect(t *testing.T) {
	tests := map[string]struct {
		module *Class
		want   string
	}{
		"with name": {
			module: NewClassWithOptions(ClassWithName("Foo")),
			want:   "class Foo < Std::Object",
		},
		"anonymous": {
			module: NewClass(),
			want:   "class <anonymous> < Std::Object",
		},
		"with name and parent": {
			module: NewClassWithOptions(ClassWithName("FooError"), ClassWithParent(ErrorClass)),
			want:   "class FooError < Std::Error",
		},
		"with name and anonymous parent": {
			module: NewClassWithOptions(ClassWithName("FooError"), ClassWithParent(NewClass())),
			want:   "class FooError < <anonymous>",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.module.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
