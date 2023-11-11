package value

import (
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/google/go-cmp/cmp"
)

func TestClass_Inspect(t *testing.T) {
	tests := map[string]struct {
		class *Class
		want  string
	}{
		"with name": {
			class: NewClassWithOptions(ClassWithName("Foo")),
			want:  "class Foo < Std::Object",
		},
		"anonymous": {
			class: NewClass(),
			want:  "class <anonymous> < Std::Object",
		},
		"with name and parent": {
			class: NewClassWithOptions(ClassWithName("FooError"), ClassWithParent(ErrorClass)),
			want:  "class FooError < Std::Error",
		},
		"with name and anonymous parent": {
			class: NewClassWithOptions(ClassWithName("FooError"), ClassWithParent(NewClass())),
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
		class *Class
		name  Symbol
		want  Method
	}{
		"get method from parent": {
			class: NewClassWithOptions(
				ClassWithParent(
					NewClassWithOptions(
						ClassWithMethods(MethodMap{
							SymbolTable.Add("foo"): NewBytecodeFunction(
								SymbolTable.Add("foo"),
								[]byte{},
								&position.Location{},
							),
						}),
					),
				),
			),
			name: SymbolTable.Add("foo"),
			want: NewBytecodeFunction(
				SymbolTable.Add("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get method from parents parent": {
			class: NewClassWithOptions(
				ClassWithParent(
					NewClassWithOptions(
						ClassWithParent(
							NewClassWithOptions(
								ClassWithMethods(MethodMap{
									SymbolTable.Add("foo"): NewBytecodeFunction(
										SymbolTable.Add("foo"),
										[]byte{},
										&position.Location{},
									),
								}),
							),
						),
					),
				),
			),
			name: SymbolTable.Add("foo"),
			want: NewBytecodeFunction(
				SymbolTable.Add("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get method from class": {
			class: NewClassWithOptions(
				ClassWithMethods(MethodMap{
					SymbolTable.Add("foo"): NewBytecodeFunction(
						SymbolTable.Add("foo"),
						[]byte{},
						&position.Location{},
					),
				}),
			),
			name: SymbolTable.Add("foo"),
			want: NewBytecodeFunction(
				SymbolTable.Add("foo"),
				[]byte{},
				&position.Location{},
			),
		},
		"get nil method": {
			class: NewClass(),
			name:  SymbolTable.Add("foo"),
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
		self      *Class
		other     *Mixin
		selfAfter *Class
	}{
		"include mixin with a method": {
			self: NewClassWithOptions(
				ClassWithName("Foo"),
				ClassWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
			),
			other: NewMixinWithOptions(
				MixinWithName("Bar"),
				MixinWithMethods(MethodMap{
					SymbolTable.Add("bar"): nil,
				}),
			),
			selfAfter: NewClassWithOptions(
				ClassWithName("Foo"),
				ClassWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				ClassWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("Bar"),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("bar"): nil,
						}),
						ClassWithParent(ObjectClass),
					),
				),
			),
		},
		"include mixin with parent": {
			self: NewClassWithOptions(
				ClassWithName("Foo"),
				ClassWithMethods(MethodMap{
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
			selfAfter: NewClassWithOptions(
				ClassWithName("Foo"),
				ClassWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				ClassWithParent(
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
								ClassWithParent(ObjectClass),
							),
						),
					),
				),
			),
		},
		"include to a class with a parent": {
			self: NewClassWithOptions(
				ClassWithName("Foo"),
				ClassWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				ClassWithParent(
					NewClassWithOptions(
						ClassWithName("FooParent"),
						ClassWithParent(ObjectClass),
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
			selfAfter: NewClassWithOptions(
				ClassWithName("Foo"),
				ClassWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				ClassWithParent(
					NewClassWithOptions(
						ClassWithMixinProxy(),
						ClassWithName("Bar"),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("bar"): nil,
						}),
						ClassWithParent(
							NewClassWithOptions(
								ClassWithName("FooParent"),
								ClassWithMethods(MethodMap{
									SymbolTable.Add("foo_parent"): nil,
								}),
								ClassWithParent(ObjectClass),
							),
						),
					),
				),
			),
		},
		"include a mixin with a parent to a class with a parent": {
			self: NewClassWithOptions(
				ClassWithName("Foo"),
				ClassWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				ClassWithParent(
					NewClassWithOptions(
						ClassWithName("FooParent"),
						ClassWithMethods(MethodMap{
							SymbolTable.Add("foo_parent"): nil,
						}),
						ClassWithParent(
							NewClassWithOptions(
								ClassWithName("FooGrandParent"),
								ClassWithMethods(MethodMap{
									SymbolTable.Add("foo_grand_parent"): nil,
								}),
								ClassWithParent(ObjectClass),
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
			selfAfter: NewClassWithOptions(
				ClassWithName("Foo"),
				ClassWithMethods(MethodMap{
					SymbolTable.Add("foo"): nil,
				}),
				ClassWithParent(
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
												ClassWithName("FooParent"),
												ClassWithMethods(MethodMap{
													SymbolTable.Add("foo_parent"): nil,
												}),
												ClassWithParent(
													NewClassWithOptions(
														ClassWithName("FooGrandParent"),
														ClassWithMethods(MethodMap{
															SymbolTable.Add("foo_grand_parent"): nil,
														}),
														ClassWithParent(ObjectClass),
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
