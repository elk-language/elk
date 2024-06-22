package types

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNilable(t *testing.T) {
	env := NewGlobalEnvironment()

	tests := map[string]struct {
		typ  Type
		want Type
	}{
		"nil stays unchanged": {
			typ:  env.StdSubtypeString("Nil"),
			want: env.StdSubtypeString("Nil"),
		},
		"nilable stays unchanged": {
			typ:  NewNilable(env.StdSubtypeString("String")),
			want: NewNilable(env.StdSubtypeString("String")),
		},
		"class type gets wrapper": {
			typ:  env.StdSubtypeString("String"),
			want: NewNilable(env.StdSubtypeString("String")),
		},
		"union type with nil stays unchanged": {
			typ:  NewUnion(env.StdSubtypeString("String"), env.StdSubtypeString("Nil")),
			want: NewUnion(env.StdSubtypeString("String"), env.StdSubtypeString("Nil")),
		},
		"union type with nilable stays unchanged": {
			typ:  NewUnion(env.StdSubtypeString("String"), NewNilable(env.StdSubtypeString("Int"))),
			want: NewUnion(env.StdSubtypeString("String"), NewNilable(env.StdSubtypeString("Int"))),
		},
		"union type without nil adds nil": {
			typ:  NewUnion(env.StdSubtypeString("String"), env.StdSubtypeString("Int")),
			want: NewUnion(env.StdSubtypeString("String"), env.StdSubtypeString("Int"), env.StdSubtypeString("Nil")),
		},
		"intersection type with nil stays unchanged": {
			typ:  NewIntersection(env.StdSubtypeString("String"), env.StdSubtypeString("Nil")),
			want: NewIntersection(env.StdSubtypeString("String"), env.StdSubtypeString("Nil")),
		},
		"intersection type with nilable stays unchanged": {
			typ:  NewIntersection(env.StdSubtypeString("String"), NewNilable(env.StdSubtypeString("Int"))),
			want: NewIntersection(env.StdSubtypeString("String"), NewNilable(env.StdSubtypeString("Int"))),
		},
		"intersection type without nil gets wrapped": {
			typ:  NewIntersection(env.StdSubtypeString("String"), env.StdSubtypeString("Int")),
			want: NewNilable(NewIntersection(env.StdSubtypeString("String"), env.StdSubtypeString("Int"))),
		},
	}

	cmpOpts := []cmp.Option{
		cmp.AllowUnexported(
			NamespaceBase{},
			Class{},
			Mixin{},
			MixinProxy{},
		),
		cmpopts.IgnoreUnexported(
			MethodMap{},
			TypeMap{},
		),
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := ToNilable(tc.typ, env)

			if diff := cmp.Diff(tc.want, got, cmpOpts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
