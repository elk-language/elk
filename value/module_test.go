package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestModuleInspect(t *testing.T) {
	tests := map[string]struct {
		module *value.Module
		want   string
	}{
		"with name": {
			module: value.NewModuleWithOptions(value.ModuleWithName("Foo")),
			want:   "module Foo",
		},
		"anonymous": {
			module: value.NewModule(),
			want:   "module <anonymous>",
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
