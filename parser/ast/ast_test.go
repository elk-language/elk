package ast

import (
	"errors"
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
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

func TestTraverseMap(t *testing.T) {
	loc := L("main", S(P(0, 1, 1), P(5, 1, 6)))
	left := NewIntLiteralNode(loc, "1")
	left.SetType(types.Untyped{})
	expr := NewBinaryExpressionNode(
		loc,
		T(loc, token.PLUS),
		left,
		NewIntLiteralNode(loc, "2"),
	)

	copiedNode, err := TraverseMap(expr, func(node, parent Node) (Node, error) {
		return node, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	copied := copiedNode.(*BinaryExpressionNode)

	if !expr.LaxEqual(value.Ref(copied)) {
		t.Fatalf("copy mismatch: %s", copied.String())
	}
	if copied == expr || copied.Left == left {
		t.Fatal("expected a deep copy")
	}
	if copied.Left.Type(nil) != nil {
		t.Fatal("copied literal should not keep its type")
	}

	replaced, err := TraverseMap(expr, func(node, parent Node) (Node, error) {
		lit, ok := node.(*IntLiteralNode)
		if ok && lit.Value == "1" {
			return NewIntLiteralNode(lit.Location(), "9"), nil
		}
		return node, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if replaced.String() != "9 + 2" {
		t.Fatalf("replaced tree is %s", replaced.String())
	}
	if expr.String() != "1 + 2" {
		t.Fatalf("original tree changed: %s", expr.String())
	}

	var seenParent Node
	if _, err := TraverseMap(expr, func(node, parent Node) (Node, error) {
		if lit, ok := node.(*IntLiteralNode); ok && lit.Value == "1" {
			seenParent = parent
		}
		return node, nil
	}); err != nil {
		t.Fatal(err)
	}
	if seenParent != Node(expr) {
		t.Fatal("expected the literal's parent to be the original binary expression")
	}

	classNode := NewClassDeclarationNode(loc, "docs", true, false, false, false, false, nil, nil, nil, nil)
	classNode.HasDefer = true
	classNode.SetType(types.Untyped{})
	classCopyNode, err := TraverseMap(classNode, func(node, parent Node) (Node, error) {
		return node, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	classCopy := classCopyNode.(*ClassDeclarationNode)
	if !classNode.LaxEqual(value.Ref(classCopy)) || classCopy.DocComment() != "docs" || !classCopy.Abstract {
		t.Fatalf("class copy mismatch: %s", classCopy.Inspect())
	}
	if classCopy.HasDefer || classCopy.Type(nil) != nil {
		t.Fatal("class copy kept typechecker state")
	}

	prog := NewProgramNode(loc, nil)
	prog.ImportPaths = []string{"foo"}
	prog.State = CHECKED_EXPRESSIONS
	prog.HasDefer = true
	prog.IsBuiltin = true
	progCopyNode, err := TraverseMap(prog, func(node, parent Node) (Node, error) {
		return node, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	progCopy := progCopyNode.(*ProgramNode)
	if progCopy.ImportPaths != nil || progCopy.State != UNCHECKED || progCopy.HasDefer || progCopy.IsBuiltin {
		t.Fatalf("program copy kept typechecker state: %#v", progCopy)
	}

	method := NewMethodDefinitionNode(
		loc,
		"doc",
		bitfield.BitFlag8(1),
		NewPublicIdentifierNode(loc, "foo"),
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	methodCopyNode, err := TraverseMap(method, func(node, parent Node) (Node, error) {
		return node, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	methodCopy := methodCopyNode.(*MethodDefinitionNode)
	if methodCopy.DocComment() != "doc" || methodCopy.Flags != method.Flags {
		t.Fatalf("method copy flags %v comment %q", methodCopy.Flags, methodCopy.DocComment())
	}

	boom := errors.New("boom")
	var calls int
	if _, err := TraverseMap(expr, func(node, parent Node) (Node, error) {
		calls++
		if lit, ok := node.(*IntLiteralNode); ok && lit.Value == "1" {
			return nil, boom
		}
		return node, nil
	}); !errors.Is(err, boom) {
		t.Fatalf("got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected the walk to stop after the first error, got %d calls", calls)
	}
}
