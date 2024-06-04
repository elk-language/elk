// Package checker implements the Elk type checker
package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(7, 1, 8), P(9, 1, 10)), "undefined constant `Foo`"),
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined constant `Std::Foo`"),
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(49, 5, 5), P(51, 5, 7)), "undefined constant `Bar`"),
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
	globalEnv := types.NewGlobalEnvironment()

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
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
									nil,
									nil,
								),
							),
							nil,
							nil,
							nil,
							types.NewClass(
								"Foo",
								globalEnv.StdSubtypeClass(symbol.Object),
								nil,
								nil,
							),
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
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
									nil,
									nil,
								),
							),
							nil,
							ast.NewPublicConstantNode(
								S(P(12, 1, 13), P(14, 1, 15)),
								"Bar",
								types.Void{},
							),
							nil,
							types.NewClass(
								"Foo",
								globalEnv.StdSubtypeClass(symbol.Object),
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined type `Bar`"),
				error.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "`void` is not a class"),
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
								types.NewClass(
									"Foo",
									types.NewClass(
										"Bar",
										globalEnv.StdSubtypeClass(symbol.Object),
										nil,
										nil,
									),
									nil,
									nil,
								),
							),
							nil,
							ast.NewPublicConstantNode(
								S(P(12, 1, 13), P(14, 1, 15)),
								"Bar",
								types.NewClass(
									"Bar",
									globalEnv.StdSubtypeClass(symbol.Object),
									nil,
									nil,
								),
							),
							nil,
							types.NewClass(
								"Foo",
								types.NewClass(
									"Bar",
									globalEnv.StdSubtypeClass(symbol.Object),
									nil,
									nil,
								),
								nil,
								nil,
							),
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
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
									nil,
									nil,
								),
							),
							nil,
							ast.NewPublicConstantNode(
								S(P(12, 1, 13), P(14, 1, 15)),
								"Bar",
								types.NewModule("Bar", nil, nil, nil),
							),
							nil,
							types.NewClass(
								"Foo",
								globalEnv.StdSubtypeClass(symbol.Object),
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "`Bar` is not a class"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}

func TestClassOverride(t *testing.T) {
	tests := simplifiedTestTable{
		"superclass matches": {
			input: `
				class Foo; end

				class Bar < Foo; end

				class Bar < Foo
					def bar; end
				end
			`,
		},
		"superclass does not match": {
			input: `
				class Foo; end

				class Bar < Foo; end

				class Bar
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(57, 6, 11), P(59, 6, 13)), "superclass mismatch in `Bar`, got `Std::Object`, expected `Foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			simplifiedCheckerTest(tc, t)
		})
	}
}

func TestInclude(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"include inexistent mixin": {
			input: `include Foo`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(8, 1, 9), P(10, 1, 11)),
									"Foo",
									types.Void{},
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(8, 1, 9), P(10, 1, 11)), "undefined type `Foo`"),
				error.NewError(L("<main>", P(8, 1, 9), P(10, 1, 11)), "only mixins can be included"),
			},
		},
		"include in top level": {
			before: `
				mixin Foo; end
			`,
			input: `include Foo`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(8, 1, 9), P(10, 1, 11)),
									"Foo",
									types.NewMixin("Foo", nil, nil, nil, nil),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(8, 1, 9), P(10, 1, 11)), "cannot include mixins in this context"),
			},
		},
		"include in module": {
			before: `
				mixin Foo; end
			`,
			input: `
			  module Bar
					include Foo
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 4, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(6, 2, 6), P(41, 4, 8)),
						ast.NewModuleDeclarationNode(
							S(P(6, 2, 6), P(40, 4, 7)),
							ast.NewPublicConstantNode(
								S(P(13, 2, 13), P(15, 2, 15)),
								"Bar",
								types.NewModule("Bar", nil, nil, nil),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(22, 3, 6), P(33, 3, 17)),
									ast.NewIncludeExpressionNode(
										S(P(22, 3, 6), P(32, 3, 16)),
										[]ast.ComplexConstantNode{
											ast.NewPublicConstantNode(
												S(P(30, 3, 14), P(32, 3, 16)),
												"Foo",
												types.NewMixin("Foo", nil, nil, nil, nil),
											),
										},
									),
								),
							},
							types.NewModule("Bar", nil, nil, nil),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(30, 3, 14), P(32, 3, 16)), "cannot include mixins in this context"),
			},
		},
		"include in class": {
			before: `
				mixin Foo; end
			`,
			input: `
			  class  Bar
					include Foo
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 4, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(6, 2, 6), P(41, 4, 8)),
						ast.NewClassDeclarationNode(
							S(P(6, 2, 6), P(40, 4, 7)),
							false,
							false,
							ast.NewPublicConstantNode(
								S(P(13, 2, 13), P(15, 2, 15)),
								"Bar",
								types.NewClass(
									"Bar",
									types.NewMixinProxy(
										types.NewMixin("Foo", nil, nil, nil, nil),
										globalEnv.StdSubtypeClass(symbol.Object),
									),
									nil,
									nil,
								),
							),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(22, 3, 6), P(33, 3, 17)),
									ast.NewIncludeExpressionNode(
										S(P(22, 3, 6), P(32, 3, 16)),
										[]ast.ComplexConstantNode{
											ast.NewPublicConstantNode(
												S(P(30, 3, 14), P(32, 3, 16)),
												"Foo",
												types.NewMixin("Foo", nil, nil, nil, nil),
											),
										},
									),
								),
							},
							types.NewClass(
								"Bar",
								types.NewMixinProxy(
									types.NewMixin("Foo", nil, nil, nil, nil),
									globalEnv.StdSubtypeClass(symbol.Object),
								),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"include in mixin": {
			before: `
				mixin Foo; end
			`,
			input: `
			  mixin  Bar
					include Foo
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 4, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(6, 2, 6), P(41, 4, 8)),
						ast.NewMixinDeclarationNode(
							S(P(6, 2, 6), P(40, 4, 7)),
							ast.NewPublicConstantNode(
								S(P(13, 2, 13), P(15, 2, 15)),
								"Bar",
								types.NewMixin(
									"Bar",
									types.NewMixinProxy(
										types.NewMixin("Foo", nil, nil, nil, nil),
										nil,
									),
									nil,
									nil,
									nil,
								),
							),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(22, 3, 6), P(33, 3, 17)),
									ast.NewIncludeExpressionNode(
										S(P(22, 3, 6), P(32, 3, 16)),
										[]ast.ComplexConstantNode{
											ast.NewPublicConstantNode(
												S(P(30, 3, 14), P(32, 3, 16)),
												"Foo",
												types.NewMixin("Foo", nil, nil, nil, nil),
											),
										},
									),
								),
							},
							types.NewMixin(
								"Bar",
								types.NewMixinProxy(
									types.NewMixin("Foo", nil, nil, nil, nil),
									nil,
								),
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}

func TestMixinType(t *testing.T) {
	tests := simplifiedTestTable{
		"assign instance of related class to mixin": {
			input: `
				mixin Bar; end
				class Foo
					include Bar
				end

				var a: Bar = Foo()
			`,
		},
		"assign instance of unrelated class to mixin": {
			input: `
				mixin Bar; end
				class Foo; end

				var a: Bar = Foo()
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(57, 5, 18), P(61, 5, 22)), "type `Foo` cannot be assigned to type `Bar`"),
			},
		},
		"assign mixin type to the same mixin type": {
			input: `
				mixin Bar; end
				class Foo
					include Bar
				end

				var a: Bar = Foo()
				var b: Bar = a
			`,
		},
		"assign related mixin type to a mixin type": {
			input: `
				mixin Baz; end

				mixin Bar
					include Baz
				end

				class Foo
					include Bar
				end

				var a: Bar = Foo()
				var b: Baz = a
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			simplifiedCheckerTest(tc, t)
		})
	}
}
