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
		sealed         bool
		err            *value.Error
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
						false,
					),
				},
			},
		},
		"define getter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("bar"): vm.NewGetterMethod(
						value.ToSymbol("bar"),
						false,
					),
				},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						false,
					),
					value.ToSymbol("bar"): vm.NewGetterMethod(
						value.ToSymbol("bar"),
						false,
					),
				},
			},
		},
		"override getter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						false,
					),
				},
			},
			attrName: "foo",
			sealed:   true,
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						true,
					),
				},
			},
		},
		"override a sealed method": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						true,
					),
				},
			},
			attrName: "foo",
			err: value.NewError(
				value.SealedMethodErrorClass,
				"cannot override a sealed method: foo",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := vm.DefineGetter(tc.container, value.ToSymbol(tc.attrName), tc.sealed)
			if diff := cmp.Diff(tc.err, err, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tc.containerAfter, tc.container, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
