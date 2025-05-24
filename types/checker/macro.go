package checker

import (
	"fmt"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

func (c *Checker) hoistMacroDefinition(node *ast.MacroDefinitionNode) {
	definedUnder := c.currentMethodScope().container
	switch d := definedUnder.(type) {
	case *types.Module:
	case *types.Class, *types.Mixin:
		definedUnder = d.Singleton()
	default:
		c.addFailure(
			fmt.Sprintf(
				"cannot declare macro `%s` in this context",
				node.Name,
			),
			node.Location(),
		)
	}

	macro := c.declareMacro(
		definedUnder,
		node.DocComment(),
		value.ToSymbol(node.Name+"!"),
		node.Parameters,
		node.Location(),
	)
	macro.Node = node
	node.SetType(macro)
	c.registerMacroCheck(macro, node)
}

type macroCheckEntry struct {
	macro          *types.Method
	constantScopes []constantScope
	methodScopes   []methodScope
	node           *ast.MacroDefinitionNode
}

func (c *Checker) registerMacroCheck(macro *types.Method, node *ast.MacroDefinitionNode) {
	c.macroChecks = append(c.macroChecks, macroCheckEntry{
		macro:          macro,
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
		node:           node,
	})
}

func (c *Checker) checkMacros() {
	concurrent.Foreach(
		concurrencyLimit,
		c.macroChecks,
		func(macroCheck macroCheckEntry) {
			macro := macroCheck.macro
			node := macroCheck.node
			macroChecker := c.newMethodChecker(
				node.Location().FilePath,
				macroCheck.constantScopes,
				macroCheck.methodScopes,
				macro.DefinedUnder,
				macro.ReturnType,
				macro.ThrowType,
				false,
			)
			macroChecker.checkMacroDefinition(node, macro)
		},
	)

	c.macroChecks = nil
	c.compiler = nil
}

func (c *Checker) checkMacroDefinition(node *ast.MacroDefinitionNode, macro *types.Method) {
	c.method = macro
	c.checkMethod(
		c.currentMethodScope().container,
		macro,
		node.Parameters,
		nil,
		nil,
		node.Body,
		node.Location(),
	)

	c.method = nil

	if c.shouldCompile() && macro.IsCompilable() {
		macro.Bytecode = c.compiler.CompileMacroBody(node, macro.Name)
	}
}

func (c *Checker) declareMacro(
	macroNamespace types.Namespace,
	docComment string,
	name value.Symbol,
	paramNodes []ast.ParameterNode,
	location *position.Location,
) *types.Method {
	exprNodeType := c.StdExpressionNode()

	var params []*types.Parameter
	for _, paramNode := range paramNodes {
		switch p := paramNode.(type) {
		case *ast.FormalParameterNode:
			var declaredType types.Type
			if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.TypeOf(p.TypeNode)
			} else {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					paramNode.Location(),
				)
			}

			var kind types.ParameterKind
			switch p.Kind {
			case ast.NormalParameterKind:
				kind = types.NormalParameterKind
			case ast.PositionalRestParameterKind:
				kind = types.PositionalRestParameterKind
			case ast.NamedRestParameterKind:
				kind = types.NamedRestParameterKind
			}
			if p.Initialiser != nil {
				kind = types.DefaultValueParameterKind
			}

			if !c.isSubtype(declaredType, exprNodeType, p.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"type `%s` does not inherit from `%s`, macro parameters must be expression nodes",
						types.InspectWithColor(declaredType),
						types.InspectWithColor(exprNodeType),
					),
					p.Location(),
				)
			}

			name := value.ToSymbol(p.Name)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType)
		default:
			c.addFailure(
				fmt.Sprintf("invalid param type %T", paramNode),
				paramNode.Location(),
			)
		}
	}

	newMacro := types.NewMethod(
		docComment,
		0,
		name,
		nil,
		params,
		exprNodeType,
		types.Never{},
		macroNamespace,
	)
	newMacro.SetLocation(location)
	macroNamespace.SetMethod(name, newMacro)
	return newMacro
}
