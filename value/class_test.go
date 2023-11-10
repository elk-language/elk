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
