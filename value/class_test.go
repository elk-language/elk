package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
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
		"singleton": {
			class: value.NewSingletonClass(value.ClassClass, "Foo"),
			want:  "class &Foo < Std::Class",
		},
		"anonymous": {
			class: value.NewClass(),
			want:  "class <anonymous> < Std::Object",
		},
		"with name and parent": {
			class: value.NewClassWithOptions(value.ClassWithName("FooError"), value.ClassWithSuperclass(value.ErrorClass)),
			want:  "class FooError < Std::Error",
		},
		"with name and anonymous parent": {
			class: value.NewClassWithOptions(value.ClassWithName("FooError"), value.ClassWithSuperclass(value.NewClass())),
			want:  "class FooError < <anonymous>",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.class.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
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
				value.ClassWithSuperclass(
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
			class: value.NewClassWithOptions(
				value.ClassWithSuperclass(
					value.NewClassWithOptions(
						value.ClassWithSuperclass(
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
		"get method from class": {
			class: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
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
			class: value.NewClass(),
			name:  value.ToSymbol("foo"),
			want:  nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.class.LookupMethod(tc.name)
			if diff := cmp.Diff(tc.want, got, comparer.Options()...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestClass_IncludeMixin(t *testing.T) {
	tests := map[string]struct {
		self            *value.Class
		other           *value.Mixin
		expectedInspect string
	}{
		"include mixin with a method": {
			self: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("bar"): nil,
				}),
			),
			expectedInspect: "Foo < Bar[Bar] < Std::Object < Std::Value",
		},
		"include mixin with parent": {
			self: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
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
						value.ClassWithSuperclass(nil),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar_parent"): nil,
						}),
						value.ClassWithMetaClass(
							value.NewMixinWithOptions(
								value.MixinWithName("BarParent"),
								value.MixinWithParent(nil),
								value.MixinWithMethods(value.MethodMap{
									value.ToSymbol("bar_parent"): nil,
								}),
							),
						),
					),
				),
			),
			expectedInspect: "Foo < Bar[Bar < BarParent[BarParent]] < Std::Object < Std::Value",
		},
		"include to a class with a parent": {
			self: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.ClassWithSuperclass(
					value.NewClassWithOptions(
						value.ClassWithName("FooParent"),
						value.ClassWithSuperclass(value.ObjectClass),
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
			expectedInspect: "Foo < Bar[Bar] < FooParent < Std::Object < Std::Value",
		},
		"include a mixin with a parent to a class with a parent": {
			self: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.ClassWithSuperclass(
					value.NewClassWithOptions(
						value.ClassWithName("FooParent"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo_parent"): nil,
						}),
						value.ClassWithSuperclass(
							value.NewClassWithOptions(
								value.ClassWithName("FooGrandParent"),
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo_grand_parent"): nil,
								}),
								value.ClassWithSuperclass(value.ObjectClass),
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
						value.ClassWithSuperclass(nil),
						value.ClassWithMetaClass(
							value.NewMixinWithOptions(
								value.MixinWithName("BarParent"),
								value.MixinWithMethods(value.MethodMap{
									value.ToSymbol("bar_parent"): nil,
								}),
								value.MixinWithParent(
									value.NewClassWithOptions(
										value.ClassWithMixinProxy(),
										value.ClassWithName("BarGrandParent"),
										value.ClassWithSuperclass(nil),
										value.ClassWithMethods(value.MethodMap{
											value.ToSymbol("bar_grand_parent"): nil,
										}),
										value.ClassWithMetaClass(
											value.NewMixinWithOptions(
												value.MixinWithName("BarGrandParent"),
												value.MixinWithMethods(value.MethodMap{
													value.ToSymbol("bar_grand_parent"): nil,
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
			expectedInspect: "Foo < Bar[Bar < BarParent[BarParent < BarGrandParent[BarGrandParent]]] < FooParent < FooGrandParent < Std::Object < Std::Value",
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
