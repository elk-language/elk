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
	name           ast.ComplexConstantNode
	typeNode       ast.TypeNode
	typeVarNodes   []ast.TypeVariableNode
}

func (c *Checker) registerTypeDefinitionCheck(node *ast.TypeDefinitionNode) {
	c.typeDefinitionChecks.Append(typeDefinitionCheckEntry{
		filename:       c.Filename,
		constantScopes: c.constantScopesCopy(),
		name:           node.Constant,
		typeNode:       node.TypeNode,
	})
}

func (c *Checker) registerGenericTypeDefinitionCheck(node *ast.GenericTypeDefinitionNode) {
	c.typeDefinitionChecks.Append(typeDefinitionCheckEntry{
		filename:       c.Filename,
		constantScopes: c.constantScopesCopyWithoutCache(),
		name:           node.Constant,
		typeNode:       node.TypeNode,
		typeVarNodes:   node.TypeVariables,
	})
}

func (c *Checker) checkTypeDefinition(
	name ast.ComplexConstantNode,
	typeNode ast.TypeNode,
	typeVarNodes []ast.TypeVariableNode,
) {
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
	if len(typeVarNodes) < 1 {
		typeNode = c.checkTypeNode(typeNode)
		typ := c.typeOf(typeNode)
		namedType := types.NewNamedType(fullConstantName, typ)
		container.DefineSubtype(constantName, namedType)
		return
	}

	typeVars := make([]*types.TypeVariable, len(typeVarNodes))
	typeVarMod := types.NewModule("", fmt.Sprintf("Type Variable Container of %s", fullConstantName))
	for _, typeVarNode := range typeVarNodes {
		varNode, ok := typeVarNode.(*ast.VariantTypeVariableNode)
		if !ok {
			continue
		}

		var variance types.Variance
		switch varNode.Variance {
		case ast.INVARIANT:
			variance = types.INVARIANT
		case ast.COVARIANT:
			variance = types.COVARIANT
		case ast.CONTRAVARIANT:
			variance = types.CONTRAVARIANT
		}

		var lowerType types.Type
		if varNode.LowerBound != nil {
			varNode.LowerBound = c.checkTypeNode(varNode.LowerBound)
			lowerType = c.typeOf(varNode.LowerBound)
		}

		var upperType types.Type
		if varNode.UpperBound != nil {
			varNode.UpperBound = c.checkTypeNode(varNode.UpperBound)
			upperType = c.typeOf(varNode.UpperBound)
		}

		t := types.NewTypeVariable(
			varNode.Name,
			upperType,
			lowerType,
			variance,
		)
		typeVars = append(typeVars, t)
		typeVarNode.SetType(t)
		typeVarMod.DefineSubtype(value.ToSymbol(t.Name), t)
	}

	c.pushConstScope(makeConstantScope(typeVarMod))

	typeNode = c.checkTypeNode(typeNode)
	typ := c.typeOf(typeNode)
	namedType := types.NewGenericNamedType(
		fullConstantName,
		typ,
		typeVars,
	)
	container.DefineSubtype(constantName, namedType)

	c.popConstScope()
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
