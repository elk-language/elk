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
		sealed         bool
		err            *value.Error
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
						false,
					),
				},
			},
		},
		"define setter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
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
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						false,
					),
				},
			},
		},
		"override setter in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						false,
					),
				},
			},
			attrName: "foo",
			sealed:   true,
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						true,
					),
				},
			},
		},
		"override a sealed setter": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						true,
					),
				},
			},
			attrName: "foo",
			err: value.NewError(
				value.SealedMethodErrorClass,
				"cannot override a sealed method: foo=",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := vm.DefineSetter(tc.container, value.ToSymbol(tc.attrName), tc.sealed)
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

func TestDefineAccessor(t *testing.T) {
	tests := map[string]struct {
		container      *value.MethodContainer
		attrName       string
		sealed         bool
		err            *value.Error
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
						false,
					),
					value.ToSymbol("foo"): vm.NewGetterMethod(
						value.ToSymbol("foo"),
						false,
					),
				},
			},
		},
		"define sealed accessor": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{},
			},
			attrName: "foo",
			sealed:   true,
			containerAfter: &value.MethodContainer{
				Methods: value.MethodMap{
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
		},
		"define accessor in populated method map": {
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
		},
		"override in populated method map": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
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
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						true,
					),
				},
			},
		},
		"override a sealed setter": {
			container: &value.MethodContainer{
				Methods: value.MethodMap{
					value.ToSymbol("foo="): vm.NewSetterMethod(
						value.ToSymbol("foo"),
						true,
					),
				},
			},
			attrName: "foo",
			err: value.NewError(
				value.SealedMethodErrorClass,
				"cannot override a sealed method: foo=",
			),
		},
		"override a sealed getter": {
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
			err := vm.DefineAccessor(tc.container, value.ToSymbol(tc.attrName), tc.sealed)
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
