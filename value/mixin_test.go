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
				t.Fatalf(diff)
			}
		})
	}
}

func TestMixin_IncludeMixin(t *testing.T) {
	tests := map[string]struct {
		self      *value.Mixin
		other     *value.Mixin
		selfAfter *value.Mixin
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
			selfAfter: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithParent(nil),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar"): nil,
						}),
					),
				),
			),
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
					),
				),
			),
			selfAfter: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
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
					),
				),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("bar"): nil,
				}),
			),
			selfAfter: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
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
		},
		"include a mixin with a parent to a mixin with a parent": {
			self: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("FooParent"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo_parent"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
								value.ClassWithName("FooGrandParent"),
								value.ClassWithParent(nil),
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo_grand_parent"): nil,
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
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("BarParent"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar_parent"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
								value.ClassWithName("BarGrandParent"),
								value.ClassWithParent(nil),
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("bar_grand_parent"): nil,
								}),
							),
						),
					),
				),
			),
			selfAfter: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
								value.ClassWithName("BarParent"),
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("bar_parent"): nil,
								}),
								value.ClassWithParent(
									value.NewClassWithOptions(
										value.ClassWithMixinProxy(),
										value.ClassWithName("BarGrandParent"),
										value.ClassWithMethods(value.MethodMap{
											value.ToSymbol("bar_grand_parent"): nil,
										}),
										value.ClassWithParent(
											value.NewClassWithOptions(
												value.ClassWithMixinProxy(),
												value.ClassWithName("FooParent"),
												value.ClassWithMethods(value.MethodMap{
													value.ToSymbol("foo_parent"): nil,
												}),
												value.ClassWithParent(
													value.NewClassWithOptions(
														value.ClassWithMixinProxy(),
														value.ClassWithName("FooGrandParent"),
														value.ClassWithParent(nil),
														value.ClassWithMethods(value.MethodMap{
															value.ToSymbol("foo_grand_parent"): nil,
														}),
													),
												),
											),
										),
									),
								),
							),
						),
					),
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.self.IncludeMixin(tc.other)
			if diff := cmp.Diff(tc.selfAfter, tc.self, comparer.Options()...); diff != "" {
				t.Fatalf(diff)
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
				t.Fatalf(diff)
			}
		})
	}
}

func TestMixin_DefineAliasString(t *testing.T) {
	tests := map[string]struct {
		mixin      *value.Mixin
		newName    string
		oldName    string
		err        *value.Error
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
			err:     nil,
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
			err:     nil,
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
			err:     nil,
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
		"alias override sealed method from mixin": {
			mixin: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
					value.ToSymbol("foo_alias"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo_alias"),
						vm.BytecodeFunctionWithSealed(),
					),
				}),
			),
			newName: "foo_alias",
			oldName: "foo",
			err:     value.NewError(value.SealedMethodErrorClass, "cannot override a sealed method: foo_alias"),
			mixinAfter: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
					value.ToSymbol("foo_alias"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo_alias"),
						vm.BytecodeFunctionWithSealed(),
					),
				}),
			),
		},
		"alias override sealed method from parent": {
			mixin: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo_alias"): vm.NewBytecodeFunctionWithOptions(
								vm.BytecodeFunctionWithStringName("foo_alias"),
								vm.BytecodeFunctionWithSealed(),
							),
						}),
					),
				),
			),
			newName: "foo_alias",
			oldName: "foo",
			err:     value.NewError(value.SealedMethodErrorClass, "cannot override a sealed method: foo_alias"),
			mixinAfter: value.NewMixinWithOptions(
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeFunctionWithOptions(
						vm.BytecodeFunctionWithStringName("foo"),
					),
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo_alias"): vm.NewBytecodeFunctionWithOptions(
								vm.BytecodeFunctionWithStringName("foo_alias"),
								vm.BytecodeFunctionWithSealed(),
							),
						}),
					),
				),
			),
		},
		"alias nil method": {
			mixin:      value.NewMixin(),
			newName:    "foo_alias",
			oldName:    "foo",
			err:        value.NewError(value.NoMethodErrorClass, "cannot create an alias for a nonexistent method: foo"),
			mixinAfter: value.NewMixin(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.mixin.DefineAliasString(tc.newName, tc.oldName)
			if diff := cmp.Diff(tc.err, err, comparer.Options()...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.mixinAfter, tc.mixin, comparer.Options()...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
