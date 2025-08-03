package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestDefineGetter(t *testing.T) {
	tests := map[string]struct {
		container      *value.MethodContainer
		attrName       string
		containerAfter *value.MethodContainer
	}{
		"define getter in empty method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
		},
		"define getter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("bar"): vm.NewGetterMethod(
						value.ToSymbol("bar"),
						-1,
					),
				},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
					value.ToSymbol("bar"): vm.NewGetterMethod(
						value.ToSymbol("bar"),
						-1,
					),
				},
			},
		},
		"override getter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineGetter(tc.container, value.ToSymbol(tc.attrName), -1)
			if diff := cmp.Diff(tc.containerAfter, tc.container, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
