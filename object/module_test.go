package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestModuleInspect(t *testing.T) {
	tests := map[string]struct {
		module *Module
		want   string
	}{
		"with name": {
			module: NewModule(ModuleWithName("Foo")),
			want:   "module Foo",
		},
		"anonymous": {
			module: NewModule(),
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
