package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMixin_Inspect(t *testing.T) {
	tests := map[string]struct {
		mixin *Mixin
		want  string
	}{
		"with name": {
			mixin: NewMixinWithOptions(MixinWithName("Foo")),
			want:  "mixin Foo",
		},
		"anonymous": {
			mixin: NewMixin(),
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
		self      *Mixin
		other     *Mixin
		selfAfter *Mixin
	}{
		"include mixin with a method": {
			self: NewMixinWithOptions(
				MixinWithName("Foo"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
			),
			other: NewMixinWithOptions(
				MixinWithName("Bar"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("bar"): nil,
				}),
			),
			selfAfter: NewMixinWithOptions(
				MixinWithName("Foo"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				MixinWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("Bar"),
						ClassWithParent(nil),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("bar"): nil,
						}),
					),
				),
			),
		},
		"include mixin with parent": {
			self: NewMixinWithOptions(
				MixinWithName("Foo"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
			),
			other: NewMixinWithOptions(
				MixinWithName("Bar"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("bar"): nil,
				}),
				MixinWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("BarParent"),
						ClassWithParent(nil),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("bar_parent"): nil,
						}),
					),
				),
			),
			selfAfter: NewMixinWithOptions(
				MixinWithName("Foo"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				MixinWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("Bar"),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("bar"): nil,
						}),
						ClassWithParent(
							NewClassWithOptions(
								ClassWithMixinProxy(),
								ClassWithName("BarParent"),
								ClassWithParent(nil),
								ClassWithMethods(MethodMap{
									SymbolTable.Add("bar_parent"): nil,
								}),
							),
						),
					),
				),
			),
		},
		"include to a mixin with a parent": {
			self: NewMixinWithOptions(
				MixinWithName("Foo"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				MixinWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("FooParent"),
						ClassWithParent(nil),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("foo_parent"): nil,
						}),
					),
				),
			),
			other: NewMixinWithOptions(
				MixinWithName("Bar"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("bar"): nil,
				}),
			),
			selfAfter: NewMixinWithOptions(
				MixinWithName("Foo"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				MixinWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("Bar"),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("bar"): nil,
						}),
						ClassWithParent(
							NewClassWithOptions(
								ClassWithMixinProxy(),
								ClassWithName("FooParent"),
								ClassWithParent(nil),
								ClassWithMethods(MethodMap{
									SymbolTable.Add("foo_parent"): nil,
								}),
							),
						),
					),
				),
			),
		},
		"include a mixin with a parent to a mixin with a parent": {
			self: NewMixinWithOptions(
				MixinWithName("Foo"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				MixinWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("FooParent"),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("foo_parent"): nil,
						}),
						ClassWithParent(
							NewClassWithOptions(
								ClassWithMixinProxy(),
								ClassWithName("FooGrandParent"),
								ClassWithParent(nil),
								ClassWithMethods(MethodMap{
									SymbolTable.Add("foo_grand_parent"): nil,
								}),
							),
						),
					),
				),
			),
			other: NewMixinWithOptions(
				MixinWithName("Bar"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("bar"): nil,
				}),
				MixinWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("BarParent"),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("bar_parent"): nil,
						}),
						ClassWithParent(
							NewClassWithOptions(
								ClassWithMixinProxy(),
								ClassWithName("BarGrandParent"),
								ClassWithParent(nil),
								ClassWithMethods(MethodMap{
									SymbolTable.Add("bar_grand_parent"): nil,
								}),
							),
						),
					),
				),
			),
			selfAfter: NewMixinWithOptions(
				MixinWithName("Foo"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				MixinWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("Bar"),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("bar"): nil,
						}),
						ClassWithParent(
							NewClassWithOptions(
								ClassWithMixinProxy(),
								ClassWithName("BarParent"),
								ClassWithMethods(MethodMap{
									SymbolTable.Add("bar_parent"): nil,
								}),
								ClassWithParent(
									NewClassWithOptions(
										ClassWithMixinProxy(),
										ClassWithName("BarGrandParent"),
										ClassWithMethods(MethodMap{
											SymbolTable.Add("bar_grand_parent"): nil,
										}),
										ClassWithParent(
											NewClassWithOptions(
												ClassWithMixinProxy(),
												ClassWithName("FooParent"),
												ClassWithMethods(MethodMap{
													SymbolTable.Add("foo_parent"): nil,
												}),
												ClassWithParent(
													NewClassWithOptions(
														ClassWithMixinProxy(),
														ClassWithName("FooGrandParent"),
														ClassWithParent(nil),
														ClassWithMethods(MethodMap{
															SymbolTable.Add("foo_grand_parent"): nil,
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
			if diff := cmp.Diff(tc.selfAfter, tc.self, ValueComparerOptions...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
