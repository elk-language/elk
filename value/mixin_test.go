package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestMixin_Inspect(t *testing.T) {
	tests := map[string]struct {
		mixin *value.Mixin
		want  string
	}{
		"with name": {
			mixin: value.NewMixinWithOptions(value.MixinWithName("Foo")),
			want:  "mixin Foo",
		},
		"anonymous": {
			mixin: value.NewMixin(),
			want:  "mixin <anonymous>",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.mixin.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestMixin_IncludeMixin(t *testing.T) {
	tests := map[string]struct {
		self            *value.Mixin
		other           *value.Mixin
		expectedInspect string
	}{
		"include mixin with a method": {
			self: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("bar"): nil,
				}),
			),
			expectedInspect: "Foo < Bar[Bar]",
		},
		"include mixin with parent": {
			self: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("bar"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("BarParent"),
						value.ClassWithParent(nil),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar_parent"): nil,
						}),
						value.ClassWithMetaClass(
							value.NewClassWithOptions(
								value.ClassWithMixin(),
								value.ClassWithName("BarParent"),
								value.ClassWithParent(nil),
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("bar_parent"): nil,
								}),
							),
						),
					),
				),
			),
			expectedInspect: "Foo < Bar[Bar < BarParent[BarParent]]",
		},
		"include to a mixin with a parent": {
			self: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("FooParent"),
						value.ClassWithParent(nil),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo_parent"): nil,
						}),
						value.ClassWithMetaClass(
							value.NewClassWithOptions(
								value.ClassWithMixin(),
								value.ClassWithName("FooParent"),
								value.ClassWithParent(nil),
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo_parent"): nil,
								}),
							),
						),
					),
				),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("bar"): nil,
				}),
			),
			expectedInspect: "Foo < Bar[Bar] < FooParent[FooParent]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.self.IncludeMixin(tc.other)
			if diff := cmp.Diff(tc.expectedInspect, tc.self.InspectInheritance(), comparer.Options()...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestMixin_LookupMethod(t *testing.T) {
	tests := map[string]struct {
		mixin *value.Mixin
		name  value.Symbol
		want  value.Method
	}{
		"get method from parent": {
			mixin: value.NewMixinWithOptions(
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo"): vm.NewBytecodeFunctionSimple(
								value.ToSymbol("foo"),
								[]byte{},
								&position.Location{},
							),
						}),
					),
				),
			),
			name: value.ToSymbol("foo"),
			want: vm.NewBytecodeFunctionSimple(
				value.ToSymbol("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get method from parents parent": {
			mixin: value.NewMixinWithOptions(
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo"): vm.NewBytecodeFunctionSimple(
										value.ToSymbol("foo"),
										[]byte{},
										&position.Location{},
									),
								}),
							),
						),
					),
				),
			),
			name: value.ToSymbol("foo"),
			want: vm.NewBytecodeFunctionSimple(
				value.ToSymbol("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get method from mixin": {
			mixin: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeFunctionSimple(
						value.ToSymbol("foo"),
						[]byte{},
						&position.Location{},
					),
				}),
			),
			name: value.ToSymbol("foo"),
			want: vm.NewBytecodeFunctionSimple(
				value.ToSymbol("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get nil method": {
			mixin: value.NewMixin(),
			name:  value.ToSymbol("foo"),
			want:  nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.mixin.LookupMethod(tc.name)
			if diff := cmp.Diff(tc.want, got, comparer.Options()...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestMixin_DefineAliasString(t *testing.T) {
	tests := map[string]struct {
		mixin      *value.Mixin
		newName    string
		oldName    string
		mixinAfter *value.Mixin
	}{
		"alias method from parent": {
			mixin: value.NewMixinWithOptions(
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
								vm.BytecodeFunctionWithStringName("foo"),
							),
						}),
					),
				),
			),
			newName: "foo_alias",
			oldName: "foo",
			mixinAfter: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo_alias"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
								vm.BytecodeFunctionWithStringName("foo"),
							),
						}),
					),
				),
			),
		},
		"alias method from parents parent": {
			mixin: value.NewMixinWithOptions(
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
										vm.BytecodeFunctionWithStringName("foo"),
									),
								}),
							),
						),
					),
				),
			),
			newName: "foo_alias",
			oldName: "foo",
			mixinAfter: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo_alias"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
										vm.BytecodeFunctionWithStringName("foo"),
									),
								}),
							),
						),
					),
				),
			),
		},
		"alias method from mixin": {
			mixin: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
				}),
			),
			newName: "foo_alias",
			oldName: "foo",
			mixinAfter: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
					value.ToSymbol("foo_alias"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
				}),
			),
		},
		"alias nil method": {
			mixin:      value.NewMixin(),
			newName:    "foo_alias",
			oldName:    "foo",
			mixinAfter: value.NewMixin(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.mixin.DefineAliasString(tc.newName, tc.oldName)
			if diff := cmp.Diff(tc.mixinAfter, tc.mixin, comparer.Options()...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
