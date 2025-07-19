package ast

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k0kubun/pp/v3"
)

// Create a new token in tests.
var T = token.New

// Create a new token with value in tests.
var V = token.NewWithValue

// Create a new source position in tests.
var P = position.New

// Create a new span in tests.
var S = position.NewSpan

// Create a new source location in tests.
var L = position.NewLocation
var LP = position.NewLocationWithParent

func TestSplice(t *testing.T) {
	tests := map[string]struct {
		node    Node
		loc     *position.Location
		args    *[]Node
		unquote bool
		want    Node
	}{
		"replace unquote with the argument": {
			node: NewBinaryExpressionNode(
				L("main", S(P(0, 1, 1), P(15, 1, 16))),
				T(L("main", S(P(0, 1, 1), P(15, 1, 16))), token.PLUS),
				NewUnaryExpressionNode(
					L("foo", S(P(5, 5, 2), P(15, 5, 6))),
					T(L("main", S(P(0, 1, 1), P(15, 1, 16))), token.MINUS),
					NewUnquoteNode(
						L("foo", S(P(10, 6, 2), P(35, 6, 20))),
						UNQUOTE_EXPRESSION_KIND,
						NewPublicIdentifierNode(
							L("foo", S(P(10, 6, 2), P(35, 6, 20))),
							"x",
						),
					),
				),
				NewIntLiteralNode(
					L("foo", S(P(10, 6, 2), P(35, 6, 20))),
					"20",
				),
			),
			loc: L("bar", S(P(92, 7, 10), P(115, 5, 32))),
			args: &[]Node{
				NewRawCharLiteralNode(
					L("baz", S(P(135, 41, 46), P(145, 75, 2))),
					'r',
				),
			},
			want: NewBinaryExpressionNode(
				LP(
					"bar", S(P(92, 7, 10), P(115, 5, 32)),
					L("main", S(P(0, 1, 1), P(15, 1, 16))),
				),
				T(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("main", S(P(0, 1, 1), P(15, 1, 16))),
					),
					token.PLUS,
				),
				NewUnaryExpressionNode(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("foo", S(P(5, 5, 2), P(15, 5, 6))),
					),
					T(
						LP(
							"bar", S(P(92, 7, 10), P(115, 5, 32)),
							L("main", S(P(0, 1, 1), P(15, 1, 16))),
						),
						token.MINUS,
					),
					NewRawCharLiteralNode(
						L("baz", S(P(135, 41, 46), P(145, 75, 2))),
						'r',
					),
				),
				NewIntLiteralNode(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("foo", S(P(10, 6, 2), P(35, 6, 20))),
					),
					"20",
				),
			),
		},
		"replace unquote with the argument with zero location": {
			node: NewBinaryExpressionNode(
				L("main", S(P(0, 1, 1), P(15, 1, 16))),
				T(L("main", S(P(0, 1, 1), P(15, 1, 16))), token.PLUS),
				NewUnaryExpressionNode(
					L("foo", S(P(5, 5, 2), P(15, 5, 6))),
					T(L("main", S(P(0, 1, 1), P(15, 1, 16))), token.MINUS),
					NewUnquoteNode(
						L("foo", S(P(10, 6, 2), P(35, 6, 20))),
						UNQUOTE_EXPRESSION_KIND,
						NewPublicIdentifierNode(
							L("foo", S(P(10, 6, 2), P(35, 6, 20))),
							"x",
						),
					),
				),
				NewIntLiteralNode(
					L("foo", S(P(10, 6, 2), P(35, 6, 20))),
					"20",
				),
			),
			loc: L("bar", S(P(92, 7, 10), P(115, 5, 32))),
			args: &[]Node{
				NewRawCharLiteralNode(
					position.ZeroLocation,
					'r',
				),
			},
			want: NewBinaryExpressionNode(
				LP(
					"bar", S(P(92, 7, 10), P(115, 5, 32)),
					L("main", S(P(0, 1, 1), P(15, 1, 16))),
				),
				T(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("main", S(P(0, 1, 1), P(15, 1, 16))),
					),
					token.PLUS,
				),
				NewUnaryExpressionNode(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("foo", S(P(5, 5, 2), P(15, 5, 6))),
					),
					T(
						LP(
							"bar", S(P(92, 7, 10), P(115, 5, 32)),
							L("main", S(P(0, 1, 1), P(15, 1, 16))),
						),
						token.MINUS,
					),
					NewRawCharLiteralNode(
						L("bar", S(P(92, 7, 10), P(115, 5, 32))),
						'r',
					),
				),
				NewIntLiteralNode(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("foo", S(P(10, 6, 2), P(35, 6, 20))),
					),
					"20",
				),
			),
		},
		"replace unquote with the argument with nil location": {
			node: NewBinaryExpressionNode(
				L("main", S(P(0, 1, 1), P(15, 1, 16))),
				T(L("main", S(P(0, 1, 1), P(15, 1, 16))), token.PLUS),
				NewUnaryExpressionNode(
					L("foo", S(P(5, 5, 2), P(15, 5, 6))),
					T(L("main", S(P(0, 1, 1), P(15, 1, 16))), token.MINUS),
					NewUnquoteNode(
						L("foo", S(P(10, 6, 2), P(35, 6, 20))),
						UNQUOTE_EXPRESSION_KIND,
						NewPublicIdentifierNode(
							L("foo", S(P(10, 6, 2), P(35, 6, 20))),
							"x",
						),
					),
				),
				NewIntLiteralNode(
					L("foo", S(P(10, 6, 2), P(35, 6, 20))),
					"20",
				),
			),
			loc: L("bar", S(P(92, 7, 10), P(115, 5, 32))),
			args: &[]Node{
				NewRawCharLiteralNode(
					nil,
					'r',
				),
			},
			want: NewBinaryExpressionNode(
				LP(
					"bar", S(P(92, 7, 10), P(115, 5, 32)),
					L("main", S(P(0, 1, 1), P(15, 1, 16))),
				),
				T(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("main", S(P(0, 1, 1), P(15, 1, 16))),
					),
					token.PLUS,
				),
				NewUnaryExpressionNode(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("foo", S(P(5, 5, 2), P(15, 5, 6))),
					),
					T(
						LP(
							"bar", S(P(92, 7, 10), P(115, 5, 32)),
							L("main", S(P(0, 1, 1), P(15, 1, 16))),
						),
						token.MINUS,
					),
					NewRawCharLiteralNode(
						L("bar", S(P(92, 7, 10), P(115, 5, 32))),
						'r',
					),
				),
				NewIntLiteralNode(
					LP(
						"bar", S(P(92, 7, 10), P(115, 5, 32)),
						L("foo", S(P(10, 6, 2), P(35, 6, 20))),
					),
					"20",
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []cmp.Option{
				cmp.AllowUnexported(
					NodeBase{},
					TypedNodeBase{},
					token.Token{},
					DocCommentableNodeBase{},
					BinaryExpressionNode{},
					LogicalExpressionNode{},
					KeyValueExpressionNode{},
					ArrayListLiteralNode{},
					ArrayTupleLiteralNode{},
					HashSetLiteralNode{},
					HashMapLiteralNode{},
					HashRecordLiteralNode{},
					RangeLiteralNode{},
					SubscriptExpressionNode{},
					NilSafeSubscriptExpressionNode{},
					WordArrayListLiteralNode{},
					WordHashSetLiteralNode{},
					SymbolArrayListLiteralNode{},
					SymbolHashSetLiteralNode{},
					BinArrayListLiteralNode{},
					BinHashSetLiteralNode{},
					HexArrayListLiteralNode{},
					HexHashSetLiteralNode{},
					UninterpolatedRegexLiteralNode{},
					bitfield.BitField8{},
				),
				cmpopts.IgnoreFields(
					TypedNodeBase{}, "typ",
				),
			}

			got := tc.node.splice(tc.loc, tc.args, tc.unquote)
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Log(pp.Sprint(got))
				t.Log(diff)
				t.Fail()
			}
		})
	}
}
