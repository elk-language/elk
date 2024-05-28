// Package checker implements the Elk type checker
package checker

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
	"github.com/elk-language/elk/value"
)

func TestModule(t *testing.T) {
	tests := testTable{
		"module with public constant": {
			input: `module Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(9, 1, 10)),
								"Foo",
								types.NewModule("Foo", nil, nil, nil),
							),
							nil,
							types.NewModule("Foo", nil, nil, nil),
						),
					),
				},
			),
		},
		"module with conflicting constant with Std": {
			input: `module Int; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(9, 1, 10)),
								"Int",
								types.NewModule("Int", nil, nil, nil),
							),
							nil,
							types.NewModule("Int", nil, nil, nil),
						),
					),
				},
			),
		},
		"module with private constant": {
			input: `module _Fo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(9, 1, 10)),
								"_Fo",
								types.NewModule("_Fo", nil, nil, nil),
							),
							nil,
							types.NewModule("_Fo", nil, nil, nil),
						),
					),
				},
			),
		},
		"module with simple constant lookup": {
			input: `module Std::Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(14, 1, 15)),
								"Std::Foo",
								types.NewModule("Std::Foo", nil, nil, nil),
							),
							nil,
							types.NewModule("Std::Foo", nil, nil, nil),
						),
					),
				},
			),
		},
		"module with non obvious constant lookup": {
			input: `module Int::Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(14, 1, 15)),
								"Std::Int::Foo",
								types.NewModule("Std::Int::Foo", nil, nil, nil),
							),
							nil,
							types.NewModule("Std::Int::Foo", nil, nil, nil),
						),
					),
				},
			),
		},
		"resolve module with non obvious constant lookup": {
			before: `module Int::Foo; end`,
			input:  `Int::Foo`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewPublicConstantNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							"Std::Int::Foo",
							types.NewModule(
								"Std::Int::Foo",
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"module with undefined root constant": {
			input: `module Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(14, 1, 15)),
								"Foo::Bar",
								types.NewModule("Foo::Bar", nil, nil, nil),
							),
							nil,
							types.NewModule("Foo::Bar", nil, nil, nil),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(9, 1, 10)), "undefined constant `Foo`"),
			},
		},
		"module with undefined constant in the middle": {
			input: `module Std::Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(19, 1, 20)),
								"Std::Foo::Bar",
								types.NewModule("Std::Foo::Bar", nil, nil, nil),
							),
							nil,
							types.NewModule("Std::Foo::Bar", nil, nil, nil),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined constant `Std::Foo`"),
			},
		},
		"nested modules": {
			input: `
				module Foo
					module Bar; end
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(44, 4, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(44, 4, 8)),
						ast.NewModuleDeclarationNode(
							S(P(5, 2, 5), P(43, 4, 7)),
							ast.NewPublicConstantNode(
								S(P(12, 2, 12), P(14, 2, 14)),
								"Foo",
								types.NewModule(
									"Foo",
									map[value.Symbol]types.Type{
										value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
									},
									map[value.Symbol]types.Type{
										value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
									},
									nil,
								),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 3, 6), P(36, 3, 21)),
									ast.NewModuleDeclarationNode(
										S(P(21, 3, 6), P(35, 3, 20)),
										ast.NewPublicConstantNode(
											S(P(28, 3, 13), P(30, 3, 15)),
											"Foo::Bar",
											types.NewModule("Foo::Bar", nil, nil, nil),
										),
										nil,
										types.NewModule("Foo::Bar", nil, nil, nil),
									),
								),
							},
							types.NewModule(
								"Foo",
								map[value.Symbol]types.Type{
									value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
								},
								map[value.Symbol]types.Type{
									value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
								},
								nil,
							),
						),
					),
				},
			),
		},
		"resolve constant inside of new module": {
			input: `
				module Foo
					module Bar; end
					Bar
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(53, 5, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(53, 5, 8)),
						ast.NewModuleDeclarationNode(
							S(P(5, 2, 5), P(52, 5, 7)),
							ast.NewPublicConstantNode(
								S(P(12, 2, 12), P(14, 2, 14)),
								"Foo",
								types.NewModule(
									"Foo",
									map[value.Symbol]types.Type{
										value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
									},
									map[value.Symbol]types.Type{
										value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
									},
									nil,
								),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 3, 6), P(36, 3, 21)),
									ast.NewModuleDeclarationNode(
										S(P(21, 3, 6), P(35, 3, 20)),
										ast.NewPublicConstantNode(
											S(P(28, 3, 13), P(30, 3, 15)),
											"Foo::Bar",
											types.NewModule("Foo::Bar", nil, nil, nil),
										),
										nil,
										types.NewModule("Foo::Bar", nil, nil, nil),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(42, 4, 6), P(45, 4, 9)),
									ast.NewPublicConstantNode(
										S(P(42, 4, 6), P(44, 4, 8)),
										"Foo::Bar",
										types.NewModule("Foo::Bar", nil, nil, nil),
									),
								),
							},
							types.NewModule(
								"Foo",
								map[value.Symbol]types.Type{
									value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
								},
								map[value.Symbol]types.Type{
									value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
								},
								nil,
							),
						),
					),
				},
			),
		},
		"resolve constant outside of new module": {
			input: `
				module Foo
					module Bar; end
				end
				Bar
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(52, 5, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(44, 4, 8)),
						ast.NewModuleDeclarationNode(
							S(P(5, 2, 5), P(43, 4, 7)),
							ast.NewPublicConstantNode(
								S(P(12, 2, 12), P(14, 2, 14)),
								"Foo",
								types.NewModule(
									"Foo",
									map[value.Symbol]types.Type{
										value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
									},
									map[value.Symbol]types.Type{
										value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
									},
									nil,
								),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 3, 6), P(36, 3, 21)),
									ast.NewModuleDeclarationNode(
										S(P(21, 3, 6), P(35, 3, 20)),
										ast.NewPublicConstantNode(
											S(P(28, 3, 13), P(30, 3, 15)),
											"Foo::Bar",
											types.NewModule("Foo::Bar", nil, nil, nil),
										),
										nil,
										types.NewModule("Foo::Bar", nil, nil, nil),
									),
								),
							},
							types.NewModule(
								"Foo",
								map[value.Symbol]types.Type{
									value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
								},
								map[value.Symbol]types.Type{
									value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil, nil),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(49, 5, 5), P(52, 5, 8)),
						ast.NewPublicConstantNode(
							S(P(49, 5, 5), P(51, 5, 7)),
							"Bar",
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(49, 5, 5), P(51, 5, 7)), "undefined constant `Bar`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}

func TestClass(t *testing.T) {
	tests := testTable{
		"class with public constant": {
			input: `class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							false,
							false,
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
								types.NewClass("Foo", nil, nil, nil),
							),
							nil,
							nil,
							nil,
							types.NewClass("Foo", nil, nil, nil),
						),
					),
				},
			),
		},
		"class with nonexistent superclass": {
			input: `class Foo < Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							false,
							false,
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
								types.NewClass("Foo", nil, nil, nil),
							),
							nil,
							ast.NewPublicConstantNode(
								S(P(12, 1, 13), P(14, 1, 15)),
								"Bar",
								types.Void{},
							),
							nil,
							types.NewClass("Foo", nil, nil, nil),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined type `Bar`"),
				errors.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "`void` is not a class"),
			},
		},
		"class with superclass": {
			before: `class Bar; end`,
			input:  `class Foo < Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							false,
							false,
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
								types.NewClass("Foo", types.NewClass("Bar", nil, nil, nil), nil, nil),
							),
							nil,
							ast.NewPublicConstantNode(
								S(P(12, 1, 13), P(14, 1, 15)),
								"Bar",
								types.NewClass("Bar", nil, nil, nil),
							),
							nil,
							types.NewClass("Foo", types.NewClass("Bar", nil, nil, nil), nil, nil),
						),
					),
				},
			),
		},
		"class with module superclass": {
			before: `module Bar; end`,
			input:  `class Foo < Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							false,
							false,
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
								types.NewClass("Foo", nil, nil, nil),
							),
							nil,
							ast.NewPublicConstantNode(
								S(P(12, 1, 13), P(14, 1, 15)),
								"Bar",
								types.NewModule("Bar", nil, nil, nil),
							),
							nil,
							types.NewClass("Foo", nil, nil, nil),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "`Bar` is not a class"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}
