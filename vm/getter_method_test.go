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
					),
				},
			},
		},
		"define getter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("bar"): vm.NewGetterMethod(
						value.ToSymbol("bar"),
					),
				},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
					),
					value.ToSymbol("bar"): vm.NewGetterMethod(
						value.ToSymbol("bar"),
					),
				},
			},
		},
		"override getter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineGetter(tc.container, value.ToSymbol(tc.attrName))
			if diff := cmp.Diff(tc.containerAfter, tc.container, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
