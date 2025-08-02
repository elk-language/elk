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
						-1,
					),
				},
			},
		},
		"define setter in populated method map": {
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
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
		},
		"override setter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
			attrName: "foo",
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineSetter(tc.container, value.ToSymbol(tc.attrName), -1)
			if diff := cmp.Diff(tc.containerAfter, tc.container, comparer.Options()); diff != "" {
				t.Fatal(diff)
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
						-1,
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						-1,
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
						-1,
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
		},
		"define accessor in populated method map": {
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
					value.ToSymbol("bar"): vm.NewGetterMethod(
						value.ToSymbol("bar"),
						-1,
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
		},
		"override in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
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
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						-1,
					),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vm.DefineAccessor(tc.container, value.ToSymbol(tc.attrName), -1)
			if diff := cmp.Diff(tc.containerAfter, tc.container, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
