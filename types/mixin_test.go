package types

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k0kubun/pp"
)

func TestMixin_CreateProxy(t *testing.T) {
	env := NewGlobalEnvironment()

	tests := map[string]struct {
		mixin    *Mixin
		wantHead *MixinProxy
		wantTail *MixinProxy
	}{
		"simple mixin": {
			mixin:    NewMixin("", false, "Foo", env),
			wantHead: NewMixinProxy(NewMixin("", false, "Foo", env), nil),
			wantTail: NewMixinProxy(NewMixin("", false, "Foo", env), nil),
		},
		"mixin with parent": {
			mixin: NewMixinWithDetails(
				"",
				false,
				"Foo",
				NewMixinProxy(NewMixin("", false, "Bar", env), nil),
				NewTypeMap(),
				NewTypeMap(),
				NewMethodMap(), env,
			),
			wantHead: NewMixinProxy(
				NewMixinWithDetails(
					"",
					false,
					"Foo",
					NewMixinProxy(NewMixin("", false, "Bar", env), nil),
					NewTypeMap(),
					NewTypeMap(),
					NewMethodMap(), env,
				),
				NewMixinProxy(
					NewMixin("", false, "Bar", env),
					nil,
				),
			),
			wantTail: NewMixinProxy(
				NewMixin("", false, "Bar", env),
				nil,
			),
		},
	}

	cmpOpts := []cmp.Option{
		cmp.Comparer(func(x, y NamespaceBase) bool {
			return x.name == y.name
		}),
		cmp.AllowUnexported(
			Class{},
			Mixin{},
			MixinProxy{},
			Method{},
		),
		cmpopts.IgnoreUnexported(
			MethodMap{},
			TypeMap{},
		),
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotHead, gotTail := tc.mixin.CreateProxy()
			if diff := cmp.Diff(tc.wantHead, gotHead, cmpOpts...); diff != "" {
				t.Log(pp.Sprint(gotHead))
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.wantTail, gotTail, cmpOpts...); diff != "" {
				t.Log(pp.Sprint(gotTail))
				t.Fatalf(diff)
			}
		})
	}
}
