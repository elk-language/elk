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
		container      *value.MethodContainer
		attrName       string
		containerAfter *value.MethodContainer
	}{
		"define setter in empty method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
		},
		"define setter in populated method map": {
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
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
		},
		"override setter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineSetter(tc.container, value.ToSymbol(tc.attrName))
			if diff := cmp.Diff(tc.containerAfter, tc.container, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestDefineAccessor(t *testing.T) {
	tests := map[string]struct {
		container      *value.MethodContainer
		attrName       string
		containerAfter *value.MethodContainer
	}{
		"define accessor in empty method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
		},
		"define sealed accessor": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
		},
		"define accessor in populated method map": {
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
					value.ToSymbol("bar"): vm.NewGetterMethod(
						value.ToSymbol("bar"),
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
					),
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
		},
		"override in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
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
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
					),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineAccessor(tc.container, value.ToSymbol(tc.attrName))
			if diff := cmp.Diff(tc.containerAfter, tc.container, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
