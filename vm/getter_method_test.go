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
		attrName string
		methods  func(namespace *value.Class) (before, after value.MethodMap)
	}{
		"define getter in empty method map": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				return value.MethodMap{}, value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}
			},
		},
		"define getter in populated method map": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				getter := vm.NewGetterMethod(
					namespace,
					value.ToSymbol("bar"),
					-1,
				)
				return value.MethodMap{
					value.ToSymbol("bar"): getter,
				}, value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
					value.ToSymbol("bar"): getter,
				}
			},
		},
		"override getter in populated method map": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				return value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}, value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			class := value.NewClass()
			before, after := tc.methods(class)
			for methodName, method := range before {
				class.AttachMethod(methodName, method)
			}
			vm.DefineGetter(class, value.ToSymbol(tc.attrName), -1)
			if diff := cmp.Diff(after, class.Methods(), comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
