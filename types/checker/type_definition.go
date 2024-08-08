package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
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

func (c *Checker) checkTypeIfNecessary(typ types.Type, span *position.Span) (ok bool) {
	switch t := typ.(type) {
	case *types.GenericNamedType:
		return c.checkGenericNamedType(t, span)
	case *types.NamedType:
		return c.checkNamedType(t, span)
	case *types.Class:
		return c.checkClassInheritanceIfNecessary(t, span)
	default:
		return true
	}
}

func (c *Checker) checkClassInheritanceIfNecessary(class *types.Class, span *position.Span) bool {
	if !class.FullyChecked && class.Node == nil {
		c.addFailure(
			fmt.Sprintf("Type `%s` circularly references itself", types.InspectWithColor(class)),
			span,
		)
		return false
	}
	if class.Node == nil {
		return true
	}

	node := class.Node.(*ast.ClassDeclarationNode)
	c.checkClassInheritance(node)
	return true
}

func (c *Checker) checkNamedType(namedType *types.NamedType, span *position.Span) bool {
	if namedType.Type == nil && namedType.Node == nil {
		c.addFailure(
			fmt.Sprintf("Type `%s` circularly references itself", types.InspectWithColor(namedType)),
			span,
		)
		return false
	}
	if namedType.Node == nil {
		return true
	}

	node := namedType.Node.(*ast.TypeDefinitionNode)
	namedType.Node = nil
	typeNode := c.checkTypeNode(node.TypeNode)
	typ := c.typeOf(typeNode)
	namedType.Type = typ
	node.SetType(typ)
	return true
}

func (c *Checker) checkGenericNamedType(namedType *types.GenericNamedType, span *position.Span) bool {
	if namedType.Type == nil && namedType.Node == nil {
		c.addFailure(
			fmt.Sprintf("Type `%s` circularly references itself", types.InspectWithColor(namedType)),
			span,
		)
		return false
	}
	if namedType.Node == nil {
		return true
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
	node.SetType(namedType)

	c.popConstScope()
	return true
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

type hoistedNamespaceCheck struct {
	filename       string
	node           ast.ExpressionNode
	constantScopes []constantScope
}

func (c *Checker) registerHoistedNamespaceCheck(node ast.ExpressionNode) {
	c.hoistedNamespaceChecks.Append(hoistedNamespaceCheck{
		node:           node,
		filename:       c.Filename,
		constantScopes: c.constantScopesCopy(),
	})
}

func (c *Checker) checkInheritance() {
	oldFilename := c.Filename
	oldConstantScopes := c.constantScopes
	for _, namespaceCheck := range c.hoistedNamespaceChecks.Slice {
		c.Filename = namespaceCheck.filename
		c.constantScopes = namespaceCheck.constantScopes
		switch n := namespaceCheck.node.(type) {
		case *ast.IncludeExpressionNode:
			for _, constant := range n.Constants {
				c.includeMixin(constant)
			}
			n.SetType(types.Nothing{})
		case *ast.ImplementExpressionNode:
			for _, constant := range n.Constants {
				c.implementInterface(constant)
			}
			n.SetType(types.Nothing{})
		case *ast.ClassDeclarationNode:
			c.checkClassInheritance(n)
		}
	}
	c.Filename = oldFilename
	c.constantScopes = oldConstantScopes
	c.hoistedNamespaceChecks.Slice = nil
}

func (c *Checker) checkClassInheritance(node *ast.ClassDeclarationNode) {
	class, ok := c.typeOf(node).(*types.Class)
	if !ok {
		return
	}
	var superclassType types.Type
	var superclass *types.Class
	class.Node = nil

	switch node.Superclass.(type) {
	case *ast.NilLiteralNode:
	case nil:
		superclass = c.GlobalEnv.StdSubtypeClass(symbol.Object)
		superclassType = superclass
	default:
		superclassType, _ = c.resolveConstantType(node.Superclass)
		var ok bool
		superclass, ok = superclassType.(*types.Class)
		if !ok {
			if !types.IsNothing(superclassType) {
				c.addFailure(
					fmt.Sprintf("`%s` is not a class", types.InspectWithColor(superclassType)),
					node.Superclass.Span(),
				)
			}
			break
		}

		if superclass.IsSealed() {
			c.addFailure(
				fmt.Sprintf("cannot inherit from sealed class `%s`", types.InspectWithColor(superclassType)),
				node.Superclass.Span(),
			)
		}
		if class.IsPrimitive() && !superclass.IsPrimitive() {
			c.addFailure(
				fmt.Sprintf("class `%s` must not be primitive to inherit from non-primitive class `%s`", types.InspectWithColor(class), types.InspectWithColor(superclassType)),
				node.Superclass.Span(),
			)
		}

	}

	parent := class.Superclass()
	class.FullyChecked = true
	if parent == nil && superclass != nil {
		class.SetParent(superclass)
	} else if parent != nil && parent != superclass {
		var span *position.Span
		if node.Superclass == nil {
			span = node.Span()
		} else {
			span = node.Superclass.Span()
		}

		c.addFailure(
			fmt.Sprintf(
				"superclass mismatch in `%s`, got `%s`, expected `%s`",
				class.Name(),
				types.InspectWithColor(superclassType),
				parent.Name(),
			),
			span,
		)
	}
}
