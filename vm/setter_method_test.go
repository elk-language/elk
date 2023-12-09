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
		methodMap      value.MethodMap
		attrName       string
		frozen         bool
		methodMapAfter value.MethodMap
	}{
		"define setter in empty method map": {
			methodMap: value.MethodMap{},
			attrName:  "foo",
			methodMapAfter: value.MethodMap{
				value.ToSymbol("foo="): vm.NewSetterMethod(
					value.ToSymbol("foo"),
					false,
				),
			},
		},
		"define setter in populated method map": {
			methodMap: value.MethodMap{
				value.ToSymbol("foo"): vm.NewGetterMethod(
					value.ToSymbol("foo"),
					false,
				),
			},
			attrName: "foo",
			methodMapAfter: value.MethodMap{
				value.ToSymbol("foo"): vm.NewGetterMethod(
					value.ToSymbol("foo"),
					false,
				),
				value.ToSymbol("foo="): vm.NewSetterMethod(
					value.ToSymbol("foo"),
					false,
				),
			},
		},
		"override setter in populated method map": {
			methodMap: value.MethodMap{
				value.ToSymbol("foo="): vm.NewSetterMethod(
					value.ToSymbol("foo"),
					false,
				),
			},
			attrName: "foo",
			frozen:   true,
			methodMapAfter: value.MethodMap{
				value.ToSymbol("foo="): vm.NewSetterMethod(
					value.ToSymbol("foo"),
					true,
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineSetter(tc.methodMap, tc.attrName, tc.frozen)
			if diff := cmp.Diff(tc.methodMapAfter, tc.methodMap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestDefineAccessor(t *testing.T) {
	tests := map[string]struct {
		methodMap      value.MethodMap
		attrName       string
		frozen         bool
		methodMapAfter value.MethodMap
	}{
		"define accessor in empty method map": {
			methodMap: value.MethodMap{},
			attrName:  "foo",
			methodMapAfter: value.MethodMap{
				value.ToSymbol("foo="): vm.NewSetterMethod(
					value.ToSymbol("foo"),
					false,
				),
				value.ToSymbol("foo"): vm.NewGetterMethod(
					value.ToSymbol("foo"),
					false,
				),
			},
		},
		"define frozen accessor": {
			methodMap: value.MethodMap{},
			attrName:  "foo",
			frozen:    true,
			methodMapAfter: value.MethodMap{
				value.ToSymbol("foo="): vm.NewSetterMethod(
					value.ToSymbol("foo"),
					true,
				),
				value.ToSymbol("foo"): vm.NewGetterMethod(
					value.ToSymbol("foo"),
					true,
				),
			},
		},
		"define accessor in populated method map": {
			methodMap: value.MethodMap{
				value.ToSymbol("bar"): vm.NewGetterMethod(
					value.ToSymbol("bar"),
					false,
				),
			},
			attrName: "foo",
			methodMapAfter: value.MethodMap{
				value.ToSymbol("bar"): vm.NewGetterMethod(
					value.ToSymbol("bar"),
					false,
				),
				value.ToSymbol("foo"): vm.NewGetterMethod(
					value.ToSymbol("foo"),
					false,
				),
				value.ToSymbol("foo="): vm.NewSetterMethod(
					value.ToSymbol("foo"),
					false,
				),
			},
		},
		"override in populated method map": {
			methodMap: value.MethodMap{
				value.ToSymbol("foo="): vm.NewSetterMethod(
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
				value.ToSymbol("foo="): vm.NewSetterMethod(
					value.ToSymbol("foo"),
					true,
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineAccessor(tc.methodMap, tc.attrName, tc.frozen)
			if diff := cmp.Diff(tc.methodMapAfter, tc.methodMap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
