package checker

import (
	"fmt"

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
	oldFilename := c.Filename
	oldConstantScopes := c.constantScopes
	for _, typedefCheck := range c.typeDefinitionChecks.Slice {
		c.Filename = typedefCheck.filename
		c.constantScopes = typedefCheck.constantScopes
		c.checkTypeDefinition(
			typedefCheck.name,
			typedefCheck.typeNode,
			typedefCheck.typeVarNodes,
		)
	}
	c.Filename = oldFilename
	c.constantScopes = oldConstantScopes
	c.typeDefinitionChecks.Slice = nil
}
