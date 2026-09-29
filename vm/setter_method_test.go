package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestDefineSetter(t *testing.T) {
	tests := map[string]struct {
		attrName string
		methods  func(namespace *value.Class) (before, after value.MethodMap)
	}{
		"define setter in empty method map": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				return value.MethodMap{}, value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}
			},
		},
		"define setter in populated method map": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				getter := vm.NewGetterMethod(
					namespace,
					value.ToSymbol("foo"),
					-1,
				)
				return value.MethodMap{
					value.ToSymbol("foo"): getter,
				}, value.MethodMap{
					value.ToSymbol("foo"): getter,
					value.ToSymbol("foo="): vm.NewSetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}
			},
		},
		"override setter in populated method map": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				return value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}, value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
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
			vm.DefineSetter(class, value.ToSymbol(tc.attrName), -1)
			if diff := cmp.Diff(after, class.Methods(), comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestDefineAccessor(t *testing.T) {
	tests := map[string]struct {
		attrName string
		methods  func(namespace *value.Class) (before, after value.MethodMap)
	}{
		"define accessor in empty method map": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				return value.MethodMap{}, value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}
			},
		},
		"define sealed accessor": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				return value.MethodMap{}, value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}
			},
		},
		"define accessor in populated method map": {
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
					value.ToSymbol("bar"): getter,
					value.ToSymbol("foo"): vm.NewGetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
					value.ToSymbol("foo="): vm.NewSetterMethod(
						namespace,
						value.ToSymbol("foo"),
						-1,
					),
				}
			},
		},
		"override in populated method map": {
			attrName: "foo",
			methods: func(namespace *value.Class) (value.MethodMap, value.MethodMap) {
				return value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
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
					value.ToSymbol("foo="): vm.NewSetterMethod(
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
			vm.DefineAccessor(class, value.ToSymbol(tc.attrName), -1)
			if diff := cmp.Diff(after, class.Methods(), comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
