package value_test

import (
	"testing"

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
					value.SymbolTable.Add("foo"): nil,
				}),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("bar"): nil,
				}),
			),
			selfAfter: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithParent(nil),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("bar"): nil,
						}),
					),
				),
			),
		},
		"include mixin with parent": {
			self: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("bar"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("BarParent"),
						value.ClassWithParent(nil),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("bar_parent"): nil,
						}),
					),
				),
			),
			selfAfter: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("bar"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
								value.ClassWithName("BarParent"),
								value.ClassWithParent(nil),
								value.ClassWithMethods(value.MethodMap{
									value.SymbolTable.Add("bar_parent"): nil,
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
					value.SymbolTable.Add("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("FooParent"),
						value.ClassWithParent(nil),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("foo_parent"): nil,
						}),
					),
				),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("bar"): nil,
				}),
			),
			selfAfter: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("bar"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
								value.ClassWithName("FooParent"),
								value.ClassWithParent(nil),
								value.ClassWithMethods(value.MethodMap{
									value.SymbolTable.Add("foo_parent"): nil,
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
					value.SymbolTable.Add("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("FooParent"),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("foo_parent"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
								value.ClassWithName("FooGrandParent"),
								value.ClassWithParent(nil),
								value.ClassWithMethods(value.MethodMap{
									value.SymbolTable.Add("foo_grand_parent"): nil,
								}),
							),
						),
					),
				),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("bar"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("BarParent"),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("bar_parent"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
								value.ClassWithName("BarGrandParent"),
								value.ClassWithParent(nil),
								value.ClassWithMethods(value.MethodMap{
									value.SymbolTable.Add("bar_grand_parent"): nil,
								}),
							),
						),
					),
				),
			),
			selfAfter: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.MixinWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("bar"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMixinProxy(),
								value.ClassWithName("BarParent"),
								value.ClassWithMethods(value.MethodMap{
									value.SymbolTable.Add("bar_parent"): nil,
								}),
								value.ClassWithParent(
									value.NewClassWithOptions(
										value.ClassWithMixinProxy(),
										value.ClassWithName("BarGrandParent"),
										value.ClassWithMethods(value.MethodMap{
											value.SymbolTable.Add("bar_grand_parent"): nil,
										}),
										value.ClassWithParent(
											value.NewClassWithOptions(
												value.ClassWithMixinProxy(),
												value.ClassWithName("FooParent"),
												value.ClassWithMethods(value.MethodMap{
													value.SymbolTable.Add("foo_parent"): nil,
												}),
												value.ClassWithParent(
													value.NewClassWithOptions(
														value.ClassWithMixinProxy(),
														value.ClassWithName("FooGrandParent"),
														value.ClassWithParent(nil),
														value.ClassWithMethods(value.MethodMap{
															value.SymbolTable.Add("foo_grand_parent"): nil,
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
			if diff := cmp.Diff(tc.selfAfter, tc.self, vm.ComparerOptions...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
