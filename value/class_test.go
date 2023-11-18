package value_test

import (
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestClass_Inspect(t *testing.T) {
	tests := map[string]struct {
		class *value.Class
		want  string
	}{
		"with name": {
			class: value.NewClassWithOptions(value.ClassWithName("Foo")),
			want:  "class Foo < Std::Object",
		},
		"anonymous": {
			class: value.NewClass(),
			want:  "class <anonymous> < Std::Object",
		},
		"with name and parent": {
			class: value.NewClassWithOptions(value.ClassWithName("FooError"), value.ClassWithParent(value.ErrorClass)),
			want:  "class FooError < Std::Error",
		},
		"with name and anonymous parent": {
			class: value.NewClassWithOptions(value.ClassWithName("FooError"), value.ClassWithParent(value.NewClass())),
			want:  "class FooError < <anonymous>",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.class.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestClass_LookupMethod(t *testing.T) {
	tests := map[string]struct {
		class *value.Class
		name  value.Symbol
		want  value.Method
	}{
		"get method from parent": {
			class: value.NewClassWithOptions(
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("foo"): vm.NewBytecodeMethod(
								value.SymbolTable.Add("foo"),
								[]byte{},
								&position.Location{},
							),
						}),
					),
				),
			),
			name: value.SymbolTable.Add("foo"),
			want: vm.NewBytecodeMethod(
				value.SymbolTable.Add("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get method from parents parent": {
			class: value.NewClassWithOptions(
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMethods(value.MethodMap{
									value.SymbolTable.Add("foo"): vm.NewBytecodeMethod(
										value.SymbolTable.Add("foo"),
										[]byte{},
										&position.Location{},
									),
								}),
							),
						),
					),
				),
			),
			name: value.SymbolTable.Add("foo"),
			want: vm.NewBytecodeMethod(
				value.SymbolTable.Add("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get method from class": {
			class: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): vm.NewBytecodeMethod(
						value.SymbolTable.Add("foo"),
						[]byte{},
						&position.Location{},
					),
				}),
			),
			name: value.SymbolTable.Add("foo"),
			want: vm.NewBytecodeMethod(
				value.SymbolTable.Add("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get nil method": {
			class: value.NewClass(),
			name:  value.SymbolTable.Add("foo"),
			want:  nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.class.LookupMethod(tc.name)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestClass_IncludeMixin(t *testing.T) {
	tests := map[string]struct {
		self      *value.Class
		other     *value.Mixin
		selfAfter *value.Class
	}{
		"include mixin with a method": {
			self: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.SymbolTable.Add("bar"): nil,
				}),
			),
			selfAfter: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("bar"): nil,
						}),
						value.ClassWithParent(value.ObjectClass),
					),
				),
			),
		},
		"include mixin with parent": {
			self: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
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
			selfAfter: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.ClassWithParent(
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
								value.ClassWithParent(value.ObjectClass),
							),
						),
					),
				),
			),
		},
		"include to a class with a parent": {
			self: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithName("FooParent"),
						value.ClassWithParent(value.ObjectClass),
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
			selfAfter: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("bar"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithName("FooParent"),
								value.ClassWithMethods(value.MethodMap{
									value.SymbolTable.Add("foo_parent"): nil,
								}),
								value.ClassWithParent(value.ObjectClass),
							),
						),
					),
				),
			),
		},
		"include a mixin with a parent to a class with a parent": {
			self: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithName("FooParent"),
						value.ClassWithMethods(value.MethodMap{
							value.SymbolTable.Add("foo_parent"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithName("FooGrandParent"),
								value.ClassWithMethods(value.MethodMap{
									value.SymbolTable.Add("foo_grand_parent"): nil,
								}),
								value.ClassWithParent(value.ObjectClass),
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
			selfAfter: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.SymbolTable.Add("foo"): nil,
				}),
				value.ClassWithParent(
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
												value.ClassWithName("FooParent"),
												value.ClassWithMethods(value.MethodMap{
													value.SymbolTable.Add("foo_parent"): nil,
												}),
												value.ClassWithParent(
													value.NewClassWithOptions(
														value.ClassWithName("FooGrandParent"),
														value.ClassWithMethods(value.MethodMap{
															value.SymbolTable.Add("foo_grand_parent"): nil,
														}),
														value.ClassWithParent(value.ObjectClass),
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
			if diff := cmp.Diff(tc.selfAfter, tc.self, value.ValueComparerOptions...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
