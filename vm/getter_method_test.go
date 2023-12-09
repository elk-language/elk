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
		methodMap      value.MethodMap
		attrName       string
		frozen         bool
		methodMapAfter value.MethodMap
	}{
		"define getter in empty method map": {
			methodMap: value.MethodMap{},
			attrName:  "foo",
			methodMapAfter: value.MethodMap{
				value.ToSymbol("foo"): vm.NewGetterMethod(
					value.ToSymbol("foo"),
					false,
				),
			},
		},
		"define getter in populated method map": {
			methodMap: value.MethodMap{
				value.ToSymbol("bar"): vm.NewGetterMethod(
					value.ToSymbol("bar"),
					false,
				),
			},
			attrName: "foo",
			methodMapAfter: value.MethodMap{
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
		"override getter in populated method map": {
			methodMap: value.MethodMap{
				value.ToSymbol("foo"): vm.NewGetterMethod(
					value.ToSymbol("foo"),
					false,
				),
			},
			attrName: "foo",
			frozen:   true,
			methodMapAfter: value.MethodMap{
				value.ToSymbol("foo"): vm.NewGetterMethod(
					value.ToSymbol("foo"),
					true,
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineGetter(tc.methodMap, tc.attrName, tc.frozen)
			if diff := cmp.Diff(tc.methodMapAfter, tc.methodMap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
