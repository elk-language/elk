package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/google/go-cmp/cmp"
)

// Represents a single parser test case.
type testCase struct {
	input string
	want  *ast.ProgramNode
	err   ErrorList
}

// Type of the parser test table.
type testTable map[string]testCase

// Create a new token in tests.
var T = token.New

// Create a new token with value in tests.
var V = token.NewWithValue

// Create a new source position in tests.
var P = position.New

// Function which powers all parser tests.
// Inspects if the produced AST matches the expected one.
func parserTest(tc testCase, t *testing.T) {
	ast, err := Parse([]byte(tc.input))

	if diff := cmp.Diff(tc.want, ast); diff != "" {
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
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(2, 2, 1, 3), token.STAR_STAR),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
							ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), token.INT, "2")),
						),
					),
					ast.NewExpressionStatementNode(
						P(8, 5, 1, 9),
						ast.NewBinaryExpressionNode(
							P(8, 5, 1, 9),
							T(P(10, 1, 1, 11), token.STAR),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "5")),
							ast.NewIntLiteralNode(P(12, 1, 1, 13), V(P(12, 1, 1, 13), token.INT, "8")),
						),
					),
				},
			),
		},
		"endlines can separate statements": {
			input: "1 ** 2\n5 * 8",
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(2, 2, 1, 3), token.STAR_STAR),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
							ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), token.INT, "2")),
						),
					),
					ast.NewExpressionStatementNode(
						P(7, 5, 2, 1),
						ast.NewBinaryExpressionNode(
							P(7, 5, 2, 1),
							T(P(9, 1, 2, 3), token.STAR),
							ast.NewIntLiteralNode(P(7, 1, 2, 1), V(P(7, 1, 2, 1), token.INT, "5")),
							ast.NewIntLiteralNode(P(11, 1, 2, 5), V(P(11, 1, 2, 5), token.INT, "8")),
						),
					),
				},
			),
		},
		"spaces can't separate statements": {
			input: "1 ** 2 \t 5 * 8",
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(2, 2, 1, 3), token.STAR_STAR),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
							ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), token.INT, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(9, 1, 1, 10), "unexpected INT, expected a statement separator `\\n`, `;`"),
			},
		},
		"can be empty with newlines": {
			input: "\n\n\n",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 3, 1, 1)),
				},
			),
		},
		"can be empty with semicolons": {
			input: ";;;",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewEmptyStatementNode(P(1, 1, 1, 2)),
					ast.NewEmptyStatementNode(P(2, 1, 1, 3)),
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
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 6, 1, 1),
							T(P(2, 2, 1, 3), token.MINUS_EQUAL),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
							ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), token.INT, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(0, 1, 1, 1), "invalid `-=` assignment target"),
			},
		},
		"strings are not valid assignment targets": {
			input: "'foo' -= 2",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 10, 1, 1),
							T(P(6, 2, 1, 7), token.MINUS_EQUAL),
							ast.NewRawStringLiteralNode(P(0, 5, 1, 1), "foo"),
							ast.NewIntLiteralNode(P(9, 1, 1, 10), V(P(9, 1, 1, 10), token.INT, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(0, 5, 1, 1), "invalid `-=` assignment target"),
			},
		},
		"constants are not valid assignment targets": {
			input: "FooBa -= 2",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 10, 1, 1),
							T(P(6, 2, 1, 7), token.MINUS_EQUAL),
							ast.NewPublicConstantNode(P(0, 5, 1, 1), "FooBa"),
							ast.NewIntLiteralNode(P(9, 1, 1, 10), V(P(9, 1, 1, 10), token.INT, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(0, 5, 1, 1), "constants can't be assigned, maybe you meant to declare it with `:=`"),
			},
		},
		"private constants are not valid assignment targets": {
			input: "_FooB -= 2",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 10, 1, 1),
							T(P(6, 2, 1, 7), token.MINUS_EQUAL),
							ast.NewPrivateConstantNode(P(0, 5, 1, 1), "_FooB"),
							ast.NewIntLiteralNode(P(9, 1, 1, 10), V(P(9, 1, 1, 10), token.INT, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(0, 5, 1, 1), "constants can't be assigned, maybe you meant to declare it with `:=`"),
			},
		},
		"identifiers can be assigned": {
			input: "foo -= 2",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 8, 1, 1),
							T(P(4, 2, 1, 5), token.MINUS_EQUAL),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewIntLiteralNode(P(7, 1, 1, 8), V(P(7, 1, 1, 8), token.INT, "2")),
						),
					),
				},
			),
		},
		"private identifiers can be assigned": {
			input: "_fo -= 2",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 8, 1, 1),
							T(P(4, 2, 1, 5), token.MINUS_EQUAL),
							ast.NewPrivateIdentifierNode(P(0, 3, 1, 1), "_fo"),
							ast.NewIntLiteralNode(P(7, 1, 1, 8), V(P(7, 1, 1, 8), token.INT, "2")),
						),
					),
				},
			),
		},
		"can be nested": {
			input: "foo = bar = baz = 3",
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 19, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 19, 1, 1),
							T(P(4, 1, 1, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewAssignmentExpressionNode(
								P(6, 13, 1, 7),
								T(P(10, 1, 1, 11), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
								ast.NewAssignmentExpressionNode(
									P(12, 7, 1, 13),
									T(P(16, 1, 1, 17), token.EQUAL_OP),
									ast.NewPublicIdentifierNode(P(12, 3, 1, 13), "baz"),
									ast.NewIntLiteralNode(P(18, 1, 1, 19), V(P(18, 1, 1, 19), token.INT, "3")),
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
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 19, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 19, 1, 1),
							T(P(4, 1, 1, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewAssignmentExpressionNode(
								P(6, 13, 2, 1),
								T(P(10, 1, 2, 5), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(P(6, 3, 2, 1), "bar"),
								ast.NewAssignmentExpressionNode(
									P(12, 7, 3, 1),
									T(P(16, 1, 3, 5), token.EQUAL_OP),
									ast.NewPublicIdentifierNode(P(12, 3, 3, 1), "baz"),
									ast.NewIntLiteralNode(P(18, 1, 4, 1), V(P(18, 1, 4, 1), token.INT, "3")),
								),
							),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "foo\n= bar\n= baz\n= 3",
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
					),
					ast.NewExpressionStatementNode(
						P(4, 6, 2, 1),
						ast.NewInvalidNode(P(4, 1, 2, 1), T(P(4, 1, 2, 1), token.EQUAL_OP)),
					),
					ast.NewExpressionStatementNode(
						P(10, 6, 3, 1),
						ast.NewInvalidNode(P(10, 1, 3, 1), T(P(10, 1, 3, 1), token.EQUAL_OP)),
					),
					ast.NewExpressionStatementNode(
						P(16, 3, 4, 1),
						ast.NewInvalidNode(P(16, 1, 4, 1), T(P(16, 1, 4, 1), token.EQUAL_OP)),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 1, 2, 1), "unexpected =, expected an expression"),
				NewError(P(10, 1, 3, 1), "unexpected =, expected an expression"),
				NewError(P(16, 1, 4, 1), "unexpected =, expected an expression"),
			},
		},
		"has lower precedence than other expressions": {
			input: "f = some && awesome || thing + 2 * 8 > 5 == false",
			want: ast.NewProgramNode(
				P(0, 49, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 49, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 49, 1, 1),
							T(P(2, 1, 1, 3), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(0, 1, 1, 1), "f"),
							ast.NewLogicalExpressionNode(
								P(4, 45, 1, 5),
								T(P(20, 2, 1, 21), token.OR_OR),
								ast.NewLogicalExpressionNode(
									P(4, 15, 1, 5),
									T(P(9, 2, 1, 10), token.AND_AND),
									ast.NewPublicIdentifierNode(P(4, 4, 1, 5), "some"),
									ast.NewPublicIdentifierNode(P(12, 7, 1, 13), "awesome"),
								),
								ast.NewBinaryExpressionNode(
									P(23, 26, 1, 24),
									T(P(41, 2, 1, 42), token.EQUAL_EQUAL),
									ast.NewBinaryExpressionNode(
										P(23, 17, 1, 24),
										T(P(37, 1, 1, 38), token.GREATER),
										ast.NewBinaryExpressionNode(
											P(23, 13, 1, 24),
											T(P(29, 1, 1, 30), token.PLUS),
											ast.NewPublicIdentifierNode(P(23, 5, 1, 24), "thing"),
											ast.NewBinaryExpressionNode(
												P(31, 5, 1, 32),
												T(P(33, 1, 1, 34), token.STAR),
												ast.NewIntLiteralNode(P(31, 1, 1, 32), V(P(31, 1, 1, 32), token.INT, "2")),
												ast.NewIntLiteralNode(P(35, 1, 1, 36), V(P(35, 1, 1, 36), token.INT, "8")),
											),
										),
										ast.NewIntLiteralNode(P(39, 1, 1, 40), V(P(39, 1, 1, 40), token.INT, "5")),
									),
									ast.NewFalseLiteralNode(P(44, 5, 1, 45)),
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
				P(0, 100, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 100, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 100, 1, 1),
							T(P(2, 1, 1, 3), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(0, 1, 1, 1), "a"),
							ast.NewAssignmentExpressionNode(
								P(4, 96, 1, 5),
								T(P(6, 2, 1, 7), token.MINUS_EQUAL),
								ast.NewPublicIdentifierNode(P(4, 1, 1, 5), "b"),
								ast.NewAssignmentExpressionNode(
									P(9, 91, 1, 10),
									T(P(11, 2, 1, 12), token.PLUS_EQUAL),
									ast.NewPublicIdentifierNode(P(9, 1, 1, 10), "c"),
									ast.NewAssignmentExpressionNode(
										P(14, 86, 1, 15),
										T(P(16, 2, 1, 17), token.STAR_EQUAL),
										ast.NewPublicIdentifierNode(P(14, 1, 1, 15), "d"),
										ast.NewAssignmentExpressionNode(
											P(19, 81, 1, 20),
											T(P(21, 2, 1, 22), token.SLASH_EQUAL),
											ast.NewPublicIdentifierNode(P(19, 1, 1, 20), "e"),
											ast.NewAssignmentExpressionNode(
												P(24, 76, 1, 25),
												T(P(26, 3, 1, 27), token.STAR_STAR_EQUAL),
												ast.NewPublicIdentifierNode(P(24, 1, 1, 25), "f"),
												ast.NewAssignmentExpressionNode(
													P(30, 70, 1, 31),
													T(P(32, 2, 1, 33), token.TILDE_EQUAL),
													ast.NewPublicIdentifierNode(P(30, 1, 1, 31), "g"),
													ast.NewAssignmentExpressionNode(
														P(35, 65, 1, 36),
														T(P(37, 3, 1, 38), token.AND_AND_EQUAL),
														ast.NewPublicIdentifierNode(P(35, 1, 1, 36), "h"),
														ast.NewAssignmentExpressionNode(
															P(41, 59, 1, 42),
															T(P(43, 2, 1, 44), token.AND_EQUAL),
															ast.NewPublicIdentifierNode(P(41, 1, 1, 42), "i"),
															ast.NewAssignmentExpressionNode(
																P(46, 54, 1, 47),
																T(P(48, 3, 1, 49), token.OR_OR_EQUAL),
																ast.NewPublicIdentifierNode(P(46, 1, 1, 47), "j"),
																ast.NewAssignmentExpressionNode(
																	P(52, 48, 1, 53),
																	T(P(54, 2, 1, 55), token.OR_EQUAL),
																	ast.NewPublicIdentifierNode(P(52, 1, 1, 53), "k"),
																	ast.NewAssignmentExpressionNode(
																		P(57, 43, 1, 58),
																		T(P(59, 2, 1, 60), token.XOR_EQUAL),
																		ast.NewPublicIdentifierNode(P(57, 1, 1, 58), "l"),
																		ast.NewAssignmentExpressionNode(
																			P(62, 38, 1, 63),
																			T(P(64, 3, 1, 65), token.QUESTION_QUESTION_EQUAL),
																			ast.NewPublicIdentifierNode(P(62, 1, 1, 63), "m"),
																			ast.NewAssignmentExpressionNode(
																				P(68, 32, 1, 69),
																				T(P(70, 3, 1, 71), token.LBITSHIFT_EQUAL),
																				ast.NewPublicIdentifierNode(P(68, 1, 1, 69), "n"),
																				ast.NewAssignmentExpressionNode(
																					P(74, 26, 1, 75),
																					T(P(76, 3, 1, 77), token.RBITSHIFT_EQUAL),
																					ast.NewPublicIdentifierNode(P(74, 1, 1, 75), "o"),
																					ast.NewAssignmentExpressionNode(
																						P(80, 20, 1, 81),
																						T(P(82, 2, 1, 83), token.PERCENT_EQUAL),
																						ast.NewPublicIdentifierNode(P(80, 1, 1, 81), "p"),
																						ast.NewAssignmentExpressionNode(
																							P(85, 15, 1, 86),
																							T(P(87, 4, 1, 88), token.LTRIPLE_BITSHIFT_EQUAL),
																							ast.NewPublicIdentifierNode(P(85, 1, 1, 86), "q"),
																							ast.NewAssignmentExpressionNode(
																								P(92, 8, 1, 93),
																								T(P(94, 4, 1, 95), token.RTRIPLE_BITSHIFT_EQUAL),
																								ast.NewPublicIdentifierNode(P(92, 1, 1, 93), "r"),
																								ast.NewPublicIdentifierNode(P(99, 1, 1, 100), "s"),
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

func TestConstantLookup(t *testing.T) {
	tests := testTable{
		"is executed from left to right": {
			input: "Foo::Bar::Baz",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 13, 1, 1),
							ast.NewConstantLookupNode(
								P(0, 8, 1, 1),
								ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
								ast.NewPublicConstantNode(P(5, 3, 1, 6), "Bar"),
							),
							ast.NewPublicConstantNode(P(10, 3, 1, 11), "Baz"),
						),
					),
				},
			),
		},
		"can't access private constants from the outside": {
			input: "Foo::_Bar",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 9, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							ast.NewPrivateConstantNode(P(5, 4, 1, 6), "_Bar"),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(5, 4, 1, 6), "unexpected PRIVATE_CONSTANT, can't access a private constant from the outside"),
			},
		},
		"can have newlines after the operator": {
			input: "Foo::\nBar",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 9, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							ast.NewPublicConstantNode(P(6, 3, 2, 1), "Bar"),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "Foo\n::Bar",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
					),
					ast.NewExpressionStatementNode(
						P(4, 5, 2, 1),
						ast.NewConstantLookupNode(
							P(4, 5, 2, 1),
							nil,
							ast.NewPublicConstantNode(P(6, 3, 2, 3), "Bar"),
						),
					),
				},
			),
		},
		"can be a unary operator": {
			input: "::Bar",
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 5, 1, 1),
							nil,
							ast.NewPublicConstantNode(P(2, 3, 1, 3), "Bar"),
						),
					),
				},
			),
		},
		"unary form can't have a private constant": {
			input: "::_Bar",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 6, 1, 1),
							nil,
							ast.NewPrivateConstantNode(P(2, 4, 1, 3), "_Bar"),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(2, 4, 1, 3), "unexpected PRIVATE_CONSTANT, can't access a private constant from the outside"),
			},
		},
		"can have other primary expressions as the left side": {
			input: "foo::Bar",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 8, 1, 1),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewPublicConstantNode(P(5, 3, 1, 6), "Bar"),
						),
					),
				},
			),
		},
		"must have a constant as the right side": {
			input: "foo::123",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 8, 1, 1),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewInvalidNode(P(5, 3, 1, 6), V(P(5, 3, 1, 6), token.INT, "123")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(5, 3, 1, 6), "unexpected INT, expected a constant"),
			},
		},
		"can be a part of an expression": {
			input: "foo::Bar + .3",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 13, 1, 1),
							T(P(9, 1, 1, 10), token.PLUS),
							ast.NewConstantLookupNode(
								P(0, 8, 1, 1),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicConstantNode(P(5, 3, 1, 6), "Bar"),
							),
							ast.NewFloatLiteralNode(P(11, 2, 1, 12), "0.3"),
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

// func TestX(t *testing.T) {
// 	tests := testTable{}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			parserTest(tc, t)
// 		})
// 	}
// }
