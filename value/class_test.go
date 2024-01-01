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
		"abstract": {
			class: value.NewClassWithOptions(value.ClassWithName("Foo"), value.ClassWithAbstract()),
			want:  "abstract class Foo < Std::Object",
		},
		"sealed": {
			class: value.NewClassWithOptions(value.ClassWithName("Foo"), value.ClassWithSealed()),
			want:  "sealed class Foo < Std::Object",
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
							value.ToSymbol("foo"): vm.NewBytecodeMethodSimple(
								value.ToSymbol("foo"),
								[]byte{},
								&position.Location{},
							),
						}),
					),
				),
			),
			name: value.ToSymbol("foo"),
			want: vm.NewBytecodeMethodSimple(
				value.ToSymbol("foo"),
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
									value.ToSymbol("foo"): vm.NewBytecodeMethodSimple(
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
			want: vm.NewBytecodeMethodSimple(
				value.ToSymbol("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get method from class": {
			class: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeMethodSimple(
						value.ToSymbol("foo"),
						[]byte{},
						&position.Location{},
					),
				}),
			),
			name: value.ToSymbol("foo"),
			want: vm.NewBytecodeMethodSimple(
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
			if diff := cmp.Diff(tc.want, got, comparer.Comparer...); diff != "" {
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
					value.ToSymbol("foo"): nil,
				}),
			),
			other: value.NewMixinWithOptions(
				value.MixinWithName("Bar"),
				value.MixinWithMethods(value.MethodMap{
					value.ToSymbol("bar"): nil,
				}),
			),
			selfAfter: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar"): nil,
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
			selfAfter: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.ClassWithParent(
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
					value.ToSymbol("foo"): nil,
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithName("FooParent"),
						value.ClassWithParent(value.ObjectClass),
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
			selfAfter: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMixinProxy(),
						value.ClassWithName("Bar"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("bar"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithName("FooParent"),
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo_parent"): nil,
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
					value.ToSymbol("foo"): nil,
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithName("FooParent"),
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo_parent"): nil,
						}),
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithName("FooGrandParent"),
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo_grand_parent"): nil,
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
			selfAfter: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): nil,
				}),
				value.ClassWithParent(
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
												value.ClassWithName("FooParent"),
												value.ClassWithMethods(value.MethodMap{
													value.ToSymbol("foo_parent"): nil,
												}),
												value.ClassWithParent(
													value.NewClassWithOptions(
														value.ClassWithName("FooGrandParent"),
														value.ClassWithMethods(value.MethodMap{
															value.ToSymbol("foo_grand_parent"): nil,
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
			if diff := cmp.Diff(tc.selfAfter, tc.self, comparer.Comparer...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestClass_DefineAliasString(t *testing.T) {
	tests := map[string]struct {
		class      *value.Class
		newName    string
		oldName    string
		err        *value.Error
		classAfter *value.Class
	}{
		"alias method from parent": {
			class: value.NewClassWithOptions(
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
								vm.BytecodeMethodWithStringName("foo"),
							),
						}),
					),
				),
			),
			newName: "foo_alias",
			oldName: "foo",
			err:     nil,
			classAfter: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo_alias"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
								vm.BytecodeMethodWithStringName("foo"),
							),
						}),
					),
				),
			),
		},
		"alias method from parents parent": {
			class: value.NewClassWithOptions(
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
										vm.BytecodeMethodWithStringName("foo"),
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
			classAfter: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo_alias"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithParent(
							value.NewClassWithOptions(
								value.ClassWithMethods(value.MethodMap{
									value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
										vm.BytecodeMethodWithStringName("foo"),
									),
								}),
							),
						),
					),
				),
			),
		},
		"alias method from class": {
			class: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
				}),
			),
			newName: "foo_alias",
			oldName: "foo",
			err:     nil,
			classAfter: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
					value.ToSymbol("foo_alias"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
				}),
			),
		},
		"alias override sealed method from class": {
			class: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
					value.ToSymbol("foo_alias"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo_alias"),
						vm.BytecodeMethodWithSealed(),
					),
				}),
			),
			newName: "foo_alias",
			oldName: "foo",
			err:     value.NewError(value.SealedMethodErrorClass, "cannot override a sealed method: foo_alias"),
			classAfter: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
					value.ToSymbol("foo_alias"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo_alias"),
						vm.BytecodeMethodWithSealed(),
					),
				}),
			),
		},
		"alias override sealed method from parent": {
			class: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo_alias"): vm.NewBytecodeMethodWithOptions(
								vm.BytecodeMethodWithStringName("foo_alias"),
								vm.BytecodeMethodWithSealed(),
							),
						}),
					),
				),
			),
			newName: "foo_alias",
			oldName: "foo",
			err:     value.NewError(value.SealedMethodErrorClass, "cannot override a sealed method: foo_alias"),
			classAfter: value.NewClassWithOptions(
				value.ClassWithMethods(value.MethodMap{
					value.ToSymbol("foo"): vm.NewBytecodeMethodWithOptions(
						vm.BytecodeMethodWithStringName("foo"),
					),
				}),
				value.ClassWithParent(
					value.NewClassWithOptions(
						value.ClassWithMethods(value.MethodMap{
							value.ToSymbol("foo_alias"): vm.NewBytecodeMethodWithOptions(
								vm.BytecodeMethodWithStringName("foo_alias"),
								vm.BytecodeMethodWithSealed(),
							),
						}),
					),
				),
			),
		},
		"alias nil method": {
			class:      value.NewClass(),
			newName:    "foo_alias",
			oldName:    "foo",
			err:        value.NewError(value.NoMethodErrorClass, "cannot create an alias for a nonexistent method: foo"),
			classAfter: value.NewClass(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.class.DefineAliasString(tc.newName, tc.oldName)
			if diff := cmp.Diff(tc.err, err, comparer.Comparer...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.classAfter, tc.class, comparer.Comparer...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
