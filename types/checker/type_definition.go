package checker

import (
	"fmt"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

type typeDefinitionCheckEntry struct {
	filename       string
	constantScopes []constantScope
	methodScopes   []methodScope
	name           ast.ComplexConstantNode
	typeNode       ast.TypeNode
	typeVarNodes   []ast.TypeVariableNode
}

func (c *Checker) registerTypeDefinitionCheck(node *ast.TypeDefinitionNode) {
	c.typeDefinitionChecks.Append(typeDefinitionCheckEntry{
		filename:       c.Filename,
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
		name:           node.Constant,
		typeNode:       node.TypeNode,
	})
}

func (c *Checker) checkTypeDefinition(
	name ast.ComplexConstantNode,
	typeNode ast.TypeNode,
	typeVarNodes []ast.TypeVariableNode,
) {
	if len(typeVarNodes) < 1 {
		container, constant, fullConstantName := c.resolveConstantForDeclaration(name)
		constantName := value.ToSymbol(extractConstantName(name))
		name = ast.NewPublicConstantNode(name.Span(), fullConstantName)
		if constant != nil {
			c.addFailure(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				name.Span(),
			)
		}

		container.DefineConstant(constantName, types.Void{})
		typeNode = c.checkTypeNode(typeNode)
		typ := c.typeOf(typeNode)
		namedType := types.NewNamedType(fullConstantName, typ)
		container.DefineSubtype(constantName, namedType)
		return
	}

	// TODO: generic types
}

func (c *Checker) checkTypeDefinitions() {
	concurrent.Foreach(
		concurrencyLimit,
		c.typeDefinitionChecks.Slice,
		func(typedefCheck typeDefinitionCheckEntry) {
			typedefChecker := c.newTypeDefinitionChecker(
				typedefCheck.filename,
				typedefCheck.constantScopes,
			)
			typedefChecker.checkTypeDefinition(
				typedefCheck.name,
				typedefCheck.typeNode,
				typedefCheck.typeVarNodes,
			)
		},
	)
	c.typeDefinitionChecks.Slice = nil
}
