package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

type typeDefinitionCheckEntry struct {
	filename       string
	constantScopes []constantScope
	typ            types.Type
}

func (c *Checker) registerTypeDefinitionCheck(node *ast.TypeDefinitionNode) {
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)
	if constant != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			node.Constant.Span(),
		)
	}

	container.DefineConstant(constantName, types.Void{})
	namedType := types.NewNamedType(fullConstantName, nil)
	namedType.Node = node
	container.DefineSubtype(constantName, namedType)

	c.typeDefinitionChecks.Append(typeDefinitionCheckEntry{
		filename:       c.Filename,
		constantScopes: c.constantScopesCopy(),
		typ:            namedType,
	})
}

func (c *Checker) registerGenericTypeDefinitionCheck(node *ast.GenericTypeDefinitionNode) {
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)
	if constant != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			node.Constant.Span(),
		)
	}

	container.DefineConstant(constantName, types.Void{})
	namedType := types.NewGenericNamedType(
		fullConstantName,
		nil,
		nil,
	)
	namedType.Node = node
	container.DefineSubtype(constantName, namedType)

	c.typeDefinitionChecks.Append(typeDefinitionCheckEntry{
		filename:       c.Filename,
		constantScopes: c.constantScopesCopyWithoutCache(),
		typ:            namedType,
	})
}

func (c *Checker) checkNamedType(namedType *types.NamedType, span *position.Span) {
	if namedType.Type == nil && namedType.Node == nil {
		c.addFailure(
			fmt.Sprintf("Type `%s` circularly references itself", types.InspectWithColor(namedType)),
			span,
		)
		return
	}
	if namedType.Node == nil {
		return
	}

	node := namedType.Node.(*ast.TypeDefinitionNode)
	namedType.Node = nil
	typeNode := c.checkTypeNode(node.TypeNode)
	typ := c.typeOf(typeNode)
	namedType.Type = typ
}

func (c *Checker) checkGenericNamedType(namedType *types.GenericNamedType, span *position.Span) {
	if namedType.Type == nil && namedType.Node == nil {
		c.addFailure(
			fmt.Sprintf("Type `%s` circularly references itself", types.InspectWithColor(namedType)),
			span,
		)
		return
	}
	if namedType.Node == nil {
		return
	}

	node := namedType.Node.(*ast.GenericTypeDefinitionNode)
	namedType.Node = nil

	typeVars := make([]*types.TypeParameter, 0, len(node.TypeVariables))
	typeVarMod := types.NewModule("", fmt.Sprintf("Type Variable Container of %s", namedType.Name))
	for _, typeVarNode := range node.TypeVariables {
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

		var lowerType types.Type = types.Never{}
		if varNode.LowerBound != nil {
			varNode.LowerBound = c.checkTypeNode(varNode.LowerBound)
			lowerType = c.typeOf(varNode.LowerBound)
		}

		var upperType types.Type = types.Any{}
		if varNode.UpperBound != nil {
			varNode.UpperBound = c.checkTypeNode(varNode.UpperBound)
			upperType = c.typeOf(varNode.UpperBound)
		}

		t := types.NewTypeParameter(
			value.ToSymbol(varNode.Name),
			lowerType,
			upperType,
			variance,
		)
		typeVars = append(typeVars, t)
		typeVarNode.SetType(t)
		typeVarMod.DefineSubtype(t.Name, t)
	}

	c.pushConstScope(makeConstantScope(typeVarMod))

	node.TypeNode = c.checkTypeNode(node.TypeNode)
	typ := c.typeOf(node.TypeNode)
	namedType.Type = typ
	namedType.TypeParameters = typeVars

	c.popConstScope()
}

func (c *Checker) checkTypeDefinitions() {
	oldFilename := c.Filename
	oldConstantScopes := c.constantScopes
	for _, typedefCheck := range c.typeDefinitionChecks.Slice {
		c.Filename = typedefCheck.filename
		c.constantScopes = typedefCheck.constantScopes
		switch t := typedefCheck.typ.(type) {
		case *types.NamedType:
			c.checkNamedType(t, nil)
		case *types.GenericNamedType:
			c.checkGenericNamedType(t, nil)
		}
	}
	c.Filename = oldFilename
	c.constantScopes = oldConstantScopes
	c.typeDefinitionChecks.Slice = nil
}
