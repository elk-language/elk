package parser

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
)

// Represents a single parser test case.
type testCase struct {
	input string
	want  *ast.ProgramNode
	err   errors.ErrorList
}

// Type of the parser test table.
type testTable map[string]testCase

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

// Function which powers all parser tests.
// Inspects if the produced AST matches the expected one.
func parserTest(tc testCase, t *testing.T) {
	t.Helper()
	got, err := Parse("main", tc.input)

	opts := []cmp.Option{
		cmp.AllowUnexported(
			ast.NodeBase{},
			token.Token{},
			ast.BinaryExpressionNode{},
			ast.LogicalExpressionNode{},
			ast.KeyValueExpressionNode{},
			ast.ArrayListLiteralNode{},
			ast.ArrayTupleLiteralNode{},
			ast.HashSetLiteralNode{},
			ast.HashMapLiteralNode{},
			ast.HashRecordLiteralNode{},
			ast.RangeLiteralNode{},
			ast.ArithmeticSequenceLiteralNode{},
			ast.SubscriptExpressionNode{},
			ast.NilSafeSubscriptExpressionNode{},
			ast.WordArrayListLiteralNode{},
			ast.WordHashSetLiteralNode{},
			ast.SymbolArrayListLiteralNode{},
			ast.SymbolHashSetLiteralNode{},
			ast.BinArrayListLiteralNode{},
			ast.BinHashSetLiteralNode{},
			ast.HexArrayListLiteralNode{},
			ast.HexHashSetLiteralNode{},
			ast.UninterpolatedRegexLiteralNode{},
			bitfield.BitField8{},
		),
	}
	if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
		pp.Println(got)
		t.Fatal(diff)
	}

	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Fatal(diff)
	}
}

func TestStatement(t *testing.T) {
	tests := testTable{
		"semicolons can separate statements": {
			input: "1 ** 2; 5 * 8",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(2, 1, 3), P(3, 1, 4)), token.STAR_STAR),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "2"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(8, 1, 9), P(12, 1, 13)),
						ast.NewBinaryExpressionNode(
							S(P(8, 1, 9), P(12, 1, 13)),
							T(S(P(10, 1, 11), P(10, 1, 11)), token.STAR),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "5"),
							ast.NewIntLiteralNode(S(P(12, 1, 13), P(12, 1, 13)), "8"),
						),
					),
				},
			),
		},
		"endlines can separate statements": {
			input: "1 ** 2\n5 * 8",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 2, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(2, 1, 3), P(3, 1, 4)), token.STAR_STAR),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "2"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(7, 2, 1), P(11, 2, 5)),
						ast.NewBinaryExpressionNode(
							S(P(7, 2, 1), P(11, 2, 5)),
							T(S(P(9, 2, 3), P(9, 2, 3)), token.STAR),
							ast.NewIntLiteralNode(S(P(7, 2, 1), P(7, 2, 1)), "5"),
							ast.NewIntLiteralNode(S(P(11, 2, 5), P(11, 2, 5)), "8"),
						),
					),
				},
			),
		},
		"spaces cannot separate statements": {
			input: "1 ** 2 \t 5 * 8",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(2, 1, 3), P(3, 1, 4)), token.STAR_STAR),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "2"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(9, 1, 10), P(9, 1, 10)), "unexpected INT, expected a statement separator `\\n`, `;`"),
			},
		},
		"can be empty with newlines": {
			input: "\n\n\n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 4, 0)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(2, 4, 0))),
				},
			),
		},
		"can be empty with semicolons": {
			input: ";;;",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewEmptyStatementNode(S(P(1, 1, 2), P(1, 1, 2))),
					ast.NewEmptyStatementNode(S(P(2, 1, 3), P(2, 1, 3))),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestLabeledExpression(t *testing.T) {
	tests := testTable{
		"label a literal": {
			input: "$foo: 1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewLabeledExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							"foo",
							ast.NewIntLiteralNode(
								S(P(6, 1, 7), P(6, 1, 7)),
								"1",
							),
						),
					),
				},
			),
		},
		"label an expression": {
			input: "$foo: 1 + 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewLabeledExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							"foo",
							ast.NewBinaryExpressionNode(
								S(P(6, 1, 7), P(10, 1, 11)),
								T(S(P(8, 1, 9), P(8, 1, 9)), token.PLUS),
								ast.NewIntLiteralNode(
									S(P(6, 1, 7), P(6, 1, 7)),
									"1",
								),
								ast.NewIntLiteralNode(
									S(P(10, 1, 11), P(10, 1, 11)),
									"2",
								),
							),
						),
					),
				},
			),
		},
		"label an expression in an expression": {
			input: "variable := $foo: 1 + 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							T(S(P(9, 1, 10), P(10, 1, 11)), token.COLON_EQUAL),
							ast.NewPublicIdentifierNode(
								S(P(0, 1, 1), P(7, 1, 8)),
								"variable",
							),
							ast.NewLabeledExpressionNode(
								S(P(12, 1, 13), P(22, 1, 23)),
								"foo",
								ast.NewBinaryExpressionNode(
									S(P(18, 1, 19), P(22, 1, 23)),
									T(S(P(20, 1, 21), P(20, 1, 21)), token.PLUS),
									ast.NewIntLiteralNode(
										S(P(18, 1, 19), P(18, 1, 19)),
										"1",
									),
									ast.NewIntLiteralNode(
										S(P(22, 1, 23), P(22, 1, 23)),
										"2",
									),
								),
							),
						),
					),
				},
			),
		},
		"modifiers are a part of the labeled expression": {
			input: "variable := $foo: 1 + 2 if true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 1, 31)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 1, 31)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(30, 1, 31)),
							T(S(P(9, 1, 10), P(10, 1, 11)), token.COLON_EQUAL),
							ast.NewPublicIdentifierNode(
								S(P(0, 1, 1), P(7, 1, 8)),
								"variable",
							),
							ast.NewLabeledExpressionNode(
								S(P(12, 1, 13), P(30, 1, 31)),
								"foo",
								ast.NewModifierNode(
									S(P(18, 1, 19), P(30, 1, 31)),
									T(S(P(24, 1, 25), P(25, 1, 26)), token.IF),
									ast.NewBinaryExpressionNode(
										S(P(18, 1, 19), P(22, 1, 23)),
										T(S(P(20, 1, 21), P(20, 1, 21)), token.PLUS),
										ast.NewIntLiteralNode(
											S(P(18, 1, 19), P(18, 1, 19)),
											"1",
										),
										ast.NewIntLiteralNode(
											S(P(22, 1, 23), P(22, 1, 23)),
											"2",
										),
									),
									ast.NewTrueLiteralNode(S(P(27, 1, 28), P(30, 1, 31))),
								),
							),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestInstanceVariables(t *testing.T) {
	tests := testTable{
		"read an instance variable": {
			input: "@foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewInstanceVariableNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							"foo",
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestAssignment(t *testing.T) {
	tests := testTable{
		"ints are not valid assignment targets": {
			input: "1 -= 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(2, 1, 3), P(3, 1, 4)), token.MINUS_EQUAL),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "2"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(0, 1, 1)), "invalid `-=` assignment target"),
			},
		},
		"ints are not valid declaration targets": {
			input: "1 := 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(2, 1, 3), P(3, 1, 4)), token.COLON_EQUAL),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "2"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(0, 1, 1)), "invalid `:=` declaration target"),
			},
		},
		"strings are not valid assignment targets": {
			input: "'foo' -= 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							T(S(P(6, 1, 7), P(7, 1, 8)), token.MINUS_EQUAL),
							ast.NewRawStringLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), "foo"),
							ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "2"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(4, 1, 5)), "invalid `-=` assignment target"),
			},
		},
		"strings are not valid declaration targets": {
			input: "'foo' := 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							T(S(P(6, 1, 7), P(7, 1, 8)), token.COLON_EQUAL),
							ast.NewRawStringLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), "foo"),
							ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "2"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(4, 1, 5)), "invalid `:=` declaration target"),
			},
		},
		"constants are not valid assignment targets": {
			input: "FooBa -= 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							T(S(P(6, 1, 7), P(7, 1, 8)), token.MINUS_EQUAL),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(4, 1, 5)), "FooBa"),
							ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "2"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(4, 1, 5)), "constants cannot be assigned, maybe you meant to declare it with `:=`"),
			},
		},
		"constants are valid declaration targets": {
			input: "FooBa := 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							T(S(P(6, 1, 7), P(7, 1, 8)), token.COLON_EQUAL),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(4, 1, 5)), "FooBa"),
							ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "2"),
						),
					),
				},
			),
		},
		"private constants are not valid assignment targets": {
			input: "_FooB -= 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							T(S(P(6, 1, 7), P(7, 1, 8)), token.MINUS_EQUAL),
							ast.NewPrivateConstantNode(S(P(0, 1, 1), P(4, 1, 5)), "_FooB"),
							ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "2"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(4, 1, 5)), "constants cannot be assigned, maybe you meant to declare it with `:=`"),
			},
		},
		"private constants are valid declaration targets": {
			input: "_FooB := 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							T(S(P(6, 1, 7), P(7, 1, 8)), token.COLON_EQUAL),
							ast.NewPrivateConstantNode(S(P(0, 1, 1), P(4, 1, 5)), "_FooB"),
							ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "2"),
						),
					),
				},
			),
		},
		"identifiers can be assigned": {
			input: "foo -= 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.MINUS_EQUAL),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewIntLiteralNode(S(P(7, 1, 8), P(7, 1, 8)), "2"),
						),
					),
				},
			),
		},
		"subscript can be assigned": {
			input: "foo[5] -= 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							T(S(P(7, 1, 8), P(8, 1, 9)), token.MINUS_EQUAL),
							ast.NewSubscriptExpressionNode(
								S(P(0, 1, 1), P(5, 1, 6)),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewIntLiteralNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"5",
								),
							),
							ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "2"),
						),
					),
				},
			),
		},
		"identifiers can be declared": {
			input: "foo := 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.COLON_EQUAL),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewIntLiteralNode(S(P(7, 1, 8), P(7, 1, 8)), "2"),
						),
					),
				},
			),
		},
		"private identifiers can be declared": {
			input: "_fo := 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.COLON_EQUAL),
							ast.NewPrivateIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "_fo"),
							ast.NewIntLiteralNode(S(P(7, 1, 8), P(7, 1, 8)), "2"),
						),
					),
				},
			),
		},
		"can be nested": {
			input: "foo = bar = baz = 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewAssignmentExpressionNode(
								S(P(6, 1, 7), P(18, 1, 19)),
								T(S(P(10, 1, 11), P(10, 1, 11)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
								ast.NewAssignmentExpressionNode(
									S(P(12, 1, 13), P(18, 1, 19)),
									T(S(P(16, 1, 17), P(16, 1, 17)), token.EQUAL_OP),
									ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(14, 1, 15)), "baz"),
									ast.NewIntLiteralNode(S(P(18, 1, 19), P(18, 1, 19)), "3"),
								),
							),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "foo =\nbar =\nbaz =\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 4, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 4, 1)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(18, 4, 1)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewAssignmentExpressionNode(
								S(P(6, 2, 1), P(18, 4, 1)),
								T(S(P(10, 2, 5), P(10, 2, 5)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(6, 2, 1), P(8, 2, 3)), "bar"),
								ast.NewAssignmentExpressionNode(
									S(P(12, 3, 1), P(18, 4, 1)),
									T(S(P(16, 3, 5), P(16, 3, 5)), token.EQUAL_OP),
									ast.NewPublicIdentifierNode(S(P(12, 3, 1), P(14, 3, 3)), "baz"),
									ast.NewIntLiteralNode(S(P(18, 4, 1), P(18, 4, 1)), "3"),
								),
							),
						),
					),
				},
			),
		},
		"cannot have newlines before the operator": {
			input: "foo\n= bar\n= baz\n= 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
					),
					ast.NewExpressionStatementNode(
						S(P(4, 2, 1), P(9, 2, 6)),
						ast.NewInvalidNode(S(P(4, 2, 1), P(4, 2, 1)), T(S(P(4, 2, 1), P(4, 2, 1)), token.EQUAL_OP)),
					),
					ast.NewExpressionStatementNode(
						S(P(10, 3, 1), P(15, 3, 6)),
						ast.NewInvalidNode(S(P(10, 3, 1), P(10, 3, 1)), T(S(P(10, 3, 1), P(10, 3, 1)), token.EQUAL_OP)),
					),
					ast.NewExpressionStatementNode(
						S(P(16, 4, 1), P(18, 4, 3)),
						ast.NewInvalidNode(S(P(16, 4, 1), P(16, 4, 1)), T(S(P(16, 4, 1), P(16, 4, 1)), token.EQUAL_OP)),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 2, 1), P(4, 2, 1)), "unexpected =, expected an expression"),
				errors.NewError(L("main", P(10, 3, 1), P(10, 3, 1)), "unexpected =, expected an expression"),
				errors.NewError(L("main", P(16, 4, 1), P(16, 4, 1)), "unexpected =, expected an expression"),
			},
		},
		"has lower precedence than other expressions": {
			input: "f = some && awesome || thing + 2 * 8 > 5 == false",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(48, 1, 49)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(48, 1, 49)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(48, 1, 49)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "f"),
							ast.NewLogicalExpressionNode(
								S(P(4, 1, 5), P(48, 1, 49)),
								T(S(P(20, 1, 21), P(21, 1, 22)), token.OR_OR),
								ast.NewLogicalExpressionNode(
									S(P(4, 1, 5), P(18, 1, 19)),
									T(S(P(9, 1, 10), P(10, 1, 11)), token.AND_AND),
									ast.NewPublicIdentifierNode(S(P(4, 1, 5), P(7, 1, 8)), "some"),
									ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(18, 1, 19)), "awesome"),
								),
								ast.NewBinaryExpressionNode(
									S(P(23, 1, 24), P(48, 1, 49)),
									T(S(P(41, 1, 42), P(42, 1, 43)), token.EQUAL_EQUAL),
									ast.NewBinaryExpressionNode(
										S(P(23, 1, 24), P(39, 1, 40)),
										T(S(P(37, 1, 38), P(37, 1, 38)), token.GREATER),
										ast.NewBinaryExpressionNode(
											S(P(23, 1, 24), P(35, 1, 36)),
											T(S(P(29, 1, 30), P(29, 1, 30)), token.PLUS),
											ast.NewPublicIdentifierNode(S(P(23, 1, 24), P(27, 1, 28)), "thing"),
											ast.NewBinaryExpressionNode(
												S(P(31, 1, 32), P(35, 1, 36)),
												T(S(P(33, 1, 34), P(33, 1, 34)), token.STAR),
												ast.NewIntLiteralNode(S(P(31, 1, 32), P(31, 1, 32)), "2"),
												ast.NewIntLiteralNode(S(P(35, 1, 36), P(35, 1, 36)), "8"),
											),
										),
										ast.NewIntLiteralNode(S(P(39, 1, 40), P(39, 1, 40)), "5"),
									),
									ast.NewFalseLiteralNode(S(P(44, 1, 45), P(48, 1, 49))),
								),
							),
						),
					),
				},
			),
		},
		"has many versions": {
			input: "a = b -= c += d *= e /= f **= g ~= h &&= i &= j ||= k |= l ^= m ??= n <<= o >>= p %= q <<<= r >>>= s",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(99, 1, 100)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(99, 1, 100)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(99, 1, 100)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewAssignmentExpressionNode(
								S(P(4, 1, 5), P(99, 1, 100)),
								T(S(P(6, 1, 7), P(7, 1, 8)), token.MINUS_EQUAL),
								ast.NewPublicIdentifierNode(S(P(4, 1, 5), P(4, 1, 5)), "b"),
								ast.NewAssignmentExpressionNode(
									S(P(9, 1, 10), P(99, 1, 100)),
									T(S(P(11, 1, 12), P(12, 1, 13)), token.PLUS_EQUAL),
									ast.NewPublicIdentifierNode(S(P(9, 1, 10), P(9, 1, 10)), "c"),
									ast.NewAssignmentExpressionNode(
										S(P(14, 1, 15), P(99, 1, 100)),
										T(S(P(16, 1, 17), P(17, 1, 18)), token.STAR_EQUAL),
										ast.NewPublicIdentifierNode(S(P(14, 1, 15), P(14, 1, 15)), "d"),
										ast.NewAssignmentExpressionNode(
											S(P(19, 1, 20), P(99, 1, 100)),
											T(S(P(21, 1, 22), P(22, 1, 23)), token.SLASH_EQUAL),
											ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(19, 1, 20)), "e"),
											ast.NewAssignmentExpressionNode(
												S(P(24, 1, 25), P(99, 1, 100)),
												T(S(P(26, 1, 27), P(28, 1, 29)), token.STAR_STAR_EQUAL),
												ast.NewPublicIdentifierNode(S(P(24, 1, 25), P(24, 1, 25)), "f"),
												ast.NewAssignmentExpressionNode(
													S(P(30, 1, 31), P(99, 1, 100)),
													T(S(P(32, 1, 33), P(33, 1, 34)), token.TILDE_EQUAL),
													ast.NewPublicIdentifierNode(S(P(30, 1, 31), P(30, 1, 31)), "g"),
													ast.NewAssignmentExpressionNode(
														S(P(35, 1, 36), P(99, 1, 100)),
														T(S(P(37, 1, 38), P(39, 1, 40)), token.AND_AND_EQUAL),
														ast.NewPublicIdentifierNode(S(P(35, 1, 36), P(35, 1, 36)), "h"),
														ast.NewAssignmentExpressionNode(
															S(P(41, 1, 42), P(99, 1, 100)),
															T(S(P(43, 1, 44), P(44, 1, 45)), token.AND_EQUAL),
															ast.NewPublicIdentifierNode(S(P(41, 1, 42), P(41, 1, 42)), "i"),
															ast.NewAssignmentExpressionNode(
																S(P(46, 1, 47), P(99, 1, 100)),
																T(S(P(48, 1, 49), P(50, 1, 51)), token.OR_OR_EQUAL),
																ast.NewPublicIdentifierNode(S(P(46, 1, 47), P(46, 1, 47)), "j"),
																ast.NewAssignmentExpressionNode(
																	S(P(52, 1, 53), P(99, 1, 100)),
																	T(S(P(54, 1, 55), P(55, 1, 56)), token.OR_EQUAL),
																	ast.NewPublicIdentifierNode(S(P(52, 1, 53), P(52, 1, 53)), "k"),
																	ast.NewAssignmentExpressionNode(
																		S(P(57, 1, 58), P(99, 1, 100)),
																		T(S(P(59, 1, 60), P(60, 1, 61)), token.XOR_EQUAL),
																		ast.NewPublicIdentifierNode(S(P(57, 1, 58), P(57, 1, 58)), "l"),
																		ast.NewAssignmentExpressionNode(
																			S(P(62, 1, 63), P(99, 1, 100)),
																			T(S(P(64, 1, 65), P(66, 1, 67)), token.QUESTION_QUESTION_EQUAL),
																			ast.NewPublicIdentifierNode(S(P(62, 1, 63), P(62, 1, 63)), "m"),
																			ast.NewAssignmentExpressionNode(
																				S(P(68, 1, 69), P(99, 1, 100)),
																				T(S(P(70, 1, 71), P(72, 1, 73)), token.LBITSHIFT_EQUAL),
																				ast.NewPublicIdentifierNode(S(P(68, 1, 69), P(68, 1, 69)), "n"),
																				ast.NewAssignmentExpressionNode(
																					S(P(74, 1, 75), P(99, 1, 100)),
																					T(S(P(76, 1, 77), P(78, 1, 79)), token.RBITSHIFT_EQUAL),
																					ast.NewPublicIdentifierNode(S(P(74, 1, 75), P(74, 1, 75)), "o"),
																					ast.NewAssignmentExpressionNode(
																						S(P(80, 1, 81), P(99, 1, 100)),
																						T(S(P(82, 1, 83), P(83, 1, 84)), token.PERCENT_EQUAL),
																						ast.NewPublicIdentifierNode(S(P(80, 1, 81), P(80, 1, 81)), "p"),
																						ast.NewAssignmentExpressionNode(
																							S(P(85, 1, 86), P(99, 1, 100)),
																							T(S(P(87, 1, 88), P(90, 1, 91)), token.LTRIPLE_BITSHIFT_EQUAL),
																							ast.NewPublicIdentifierNode(S(P(85, 1, 86), P(85, 1, 86)), "q"),
																							ast.NewAssignmentExpressionNode(
																								S(P(92, 1, 93), P(99, 1, 100)),
																								T(S(P(94, 1, 95), P(97, 1, 98)), token.RTRIPLE_BITSHIFT_EQUAL),
																								ast.NewPublicIdentifierNode(S(P(92, 1, 93), P(92, 1, 93)), "r"),
																								ast.NewPublicIdentifierNode(S(P(99, 1, 100), P(99, 1, 100)), "s"),
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
												),
											),
										),
									),
								),
							),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestPostfixExpressions(t *testing.T) {
	tests := testTable{
		"ints are not valid assignment targets": {
			input: "1++",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewPostfixExpressionNode(
							S(P(0, 1, 1), P(2, 1, 3)),
							T(S(P(1, 1, 2), P(2, 1, 3)), token.PLUS_PLUS),
							ast.NewIntLiteralNode(
								S(P(0, 1, 1), P(0, 1, 1)),
								"1",
							),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(0, 1, 1)), "invalid `++` assignment target"),
			},
		},
		"strings are not valid assignment targets": {
			input: "'foo'--",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewPostfixExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							T(S(P(5, 1, 6), P(6, 1, 7)), token.MINUS_MINUS),
							ast.NewRawStringLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), "foo"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(4, 1, 5)), "invalid `--` assignment target"),
			},
		},
		"constants are not valid assignment targets": {
			input: "FooBa++",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewPostfixExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							T(S(P(5, 1, 6), P(6, 1, 7)), token.PLUS_PLUS),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(4, 1, 5)), "FooBa"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(4, 1, 5)), "invalid `++` assignment target"),
			},
		},
		"identifiers can be assigned": {
			input: "foo++",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewPostfixExpressionNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							T(S(P(3, 1, 4), P(4, 1, 5)), token.PLUS_PLUS),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
						),
					),
				},
			),
		},
		"subscript can be assigned": {
			input: "foo[5]--",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewPostfixExpressionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							T(S(P(6, 1, 7), P(7, 1, 8)), token.MINUS_MINUS),
							ast.NewSubscriptExpressionNode(
								S(P(0, 1, 1), P(5, 1, 6)),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewIntLiteralNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"5",
								),
							),
						),
					),
				},
			),
		},
		"cannot be nested": {
			input: "foo++++",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewPostfixExpressionNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							T(S(P(3, 1, 4), P(4, 1, 5)), token.PLUS_PLUS),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(5, 1, 6), P(6, 1, 7)), "unexpected ++, expected a statement separator `\\n`, `;`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestConstantLookup(t *testing.T) {
	tests := testTable{
		"is executed from left to right": {
			input: "Foo::Bar::Baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewConstantLookupNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							ast.NewConstantLookupNode(
								S(P(0, 1, 1), P(7, 1, 8)),
								ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
								ast.NewPublicConstantNode(S(P(5, 1, 6), P(7, 1, 8)), "Bar"),
							),
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Baz"),
						),
					),
				},
			),
		},
		"cannot access private constants from the outside": {
			input: "Foo::_Bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewConstantLookupNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							ast.NewPrivateConstantNode(S(P(5, 1, 6), P(8, 1, 9)), "_Bar"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(5, 1, 6), P(8, 1, 9)), "unexpected PRIVATE_CONSTANT, cannot access a private constant from the outside"),
			},
		},
		"can have newlines after the operator": {
			input: "Foo::\nBar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 2, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 2, 3)),
						ast.NewConstantLookupNode(
							S(P(0, 1, 1), P(8, 2, 3)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							ast.NewPublicConstantNode(S(P(6, 2, 1), P(8, 2, 3)), "Bar"),
						),
					),
				},
			),
		},
		"cannot have newlines before the operator": {
			input: "Foo\n::Bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 2, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
					),
					ast.NewExpressionStatementNode(
						S(P(4, 2, 1), P(8, 2, 5)),
						ast.NewConstantLookupNode(
							S(P(4, 2, 1), P(8, 2, 5)),
							nil,
							ast.NewPublicConstantNode(S(P(6, 2, 3), P(8, 2, 5)), "Bar"),
						),
					),
				},
			),
		},
		"can be a unary operator": {
			input: "::Bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewConstantLookupNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							nil,
							ast.NewPublicConstantNode(S(P(2, 1, 3), P(4, 1, 5)), "Bar"),
						),
					),
				},
			),
		},
		"unary form cannot have a private constant": {
			input: "::_Bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewConstantLookupNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							nil,
							ast.NewPrivateConstantNode(S(P(2, 1, 3), P(5, 1, 6)), "_Bar"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(2, 1, 3), P(5, 1, 6)), "unexpected PRIVATE_CONSTANT, cannot access a private constant from the outside"),
			},
		},
		"can have other primary expressions as the left side": {
			input: "foo::Bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewConstantLookupNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewPublicConstantNode(S(P(5, 1, 6), P(7, 1, 8)), "Bar"),
						),
					),
				},
			),
		},
		"must have a constant as the right side": {
			input: "foo::123",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewConstantLookupNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewInvalidNode(S(P(5, 1, 6), P(7, 1, 8)), V(S(P(5, 1, 6), P(7, 1, 8)), token.INT, "123")),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(5, 1, 6), P(7, 1, 8)), "unexpected INT, expected a constant"),
			},
		},
		"can be a part of an expression": {
			input: "foo::Bar + .3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							T(S(P(9, 1, 10), P(9, 1, 10)), token.PLUS),
							ast.NewConstantLookupNode(
								S(P(0, 1, 1), P(7, 1, 8)),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewPublicConstantNode(S(P(5, 1, 6), P(7, 1, 8)), "Bar"),
							),
							ast.NewFloatLiteralNode(S(P(11, 1, 12), P(12, 1, 13)), "0.3"),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
