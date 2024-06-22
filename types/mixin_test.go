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
			mixin:    NewMixin("Foo", env),
			wantHead: NewMixinProxy(NewMixin("Foo", env), nil),
			wantTail: NewMixinProxy(NewMixin("Foo", env), nil),
		},
		"mixin with parent": {
			mixin: NewMixinWithDetails(
				"Foo",
				NewMixinProxy(NewMixin("Bar", env), nil),
				NewTypeMap(),
				NewTypeMap(),
				NewMethodMap(), env,
			),
			wantHead: NewMixinProxy(
				NewMixinWithDetails(
					"Foo",
					NewMixinProxy(NewMixin("Bar", env), nil),
					NewTypeMap(),
					NewTypeMap(),
					NewMethodMap(), env,
				),
				NewMixinProxy(
					NewMixin("Bar", env),
					nil,
				),
			),
			wantTail: NewMixinProxy(
				NewMixin("Bar", env),
				nil,
			),
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
