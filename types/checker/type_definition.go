package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type typedefState uint8

const (
	NEW_TYPEDEF typedefState = iota
	CHECKING_TYPEDEF
	CHECKED_TYPEDEF
)

type typeDefinitionChecks struct {
	m     map[string]*typeDefinitionCheck
	order []string
}

type typeDefinitionCheck struct {
	entries []*typeDefinitionCheckEntry
	typ     types.Type
	state   typedefState
}

func (t *typeDefinitionChecks) addEntry(name string, typ types.Type, entry *typeDefinitionCheckEntry) {
	existingCheck, ok := t.m[name]
	if ok {
		existingCheck.entries = append(existingCheck.entries, entry)
	} else {
		t.m[name] = &typeDefinitionCheck{
			typ: typ,
			entries: []*typeDefinitionCheckEntry{
				entry,
			},
		}
		t.order = append(t.order, name)
	}
}

func newTypeDefinitionChecks() *typeDefinitionChecks {
	return &typeDefinitionChecks{
		m: make(map[string]*typeDefinitionCheck),
	}
}

type typeDefinitionCheckEntry struct {
	filename       string
	constantScopes []constantScope
	node           ast.ExpressionNode
}

func newTypeDefinitionCheckEntry(filename string, constScopes []constantScope, node ast.ExpressionNode) *typeDefinitionCheckEntry {
	return &typeDefinitionCheckEntry{
		filename:       filename,
		constantScopes: constScopes,
		node:           node,
	}
}

func (c *Checker) registerNamespaceDeclarationCheck(name string, node ast.ExpressionNode, typ types.Type) {
	c.typeDefinitionChecks.addEntry(
		name,
		typ,
		newTypeDefinitionCheckEntry(
			c.Filename,
			c.constantScopesCopy(),
			node,
		),
	)
}

func (c *Checker) registerNamedTypeCheck(node *ast.TypeDefinitionNode) {
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)
	if constant != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			node.Constant.Span(),
		)
	}

	namedType := types.NewNamedType(fullConstantName, nil)
	container.DefineConstant(constantName, types.NoValue{})
	container.DefineSubtype(constantName, namedType)
	node.SetType(namedType)

	c.typeDefinitionChecks.addEntry(
		namedType.Name,
		namedType,
		newTypeDefinitionCheckEntry(
			c.Filename,
			c.constantScopesCopy(),
			node,
		),
	)
}

func (c *Checker) registerGenericNamedTypeCheck(node *ast.GenericTypeDefinitionNode) {
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)
	if constant != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			node.Constant.Span(),
		)
	}

	namedType := types.NewGenericNamedType(
		fullConstantName,
		nil,
		nil,
	)
	container.DefineConstant(constantName, types.NoValue{})
	container.DefineSubtype(constantName, namedType)
	node.SetType(namedType)

	c.typeDefinitionChecks.addEntry(
		namedType.Name,
		namedType,
		newTypeDefinitionCheckEntry(
			c.Filename,
			c.constantScopesCopyWithoutCache(),
			node,
		),
	)
}

func (c *Checker) checkTypeIfNecessary(name string, span *position.Span) (ok bool) {
	if c.phase != initPhase {
		return true
	}
	typedefCheck, ok := c.typeDefinitionChecks.m[name]
	if !ok {
		return true
	}

	return c.checkTypeDefinition(typedefCheck, span)
}

func (c *Checker) checkNamedType(node *ast.TypeDefinitionNode) bool {
	namedType := c.typeOf(node).(*types.NamedType)
	typeNode := c.checkTypeNode(node.TypeNode)
	typ := c.typeOf(typeNode)
	namedType.Type = typ

	return true
}

func (c *Checker) checkGenericNamedType(node *ast.GenericTypeDefinitionNode) bool {
	namedType := c.typeOf(node).(*types.GenericNamedType)

	typeParams := make([]*types.TypeParameter, 0, len(node.TypeParameters))
	typeParamMod := types.NewTypeParamNamespace(fmt.Sprintf("Type Parameter Container of %s", namedType.Name))
	for _, typeParamNode := range node.TypeParameters {
		varNode, ok := typeParamNode.(*ast.VariantTypeParameterNode)
		if !ok {
			continue
		}

		t := c.checkTypeParameterNode(varNode)
		typeParams = append(typeParams, t)
		typeParamNode.SetType(t)
		typeParamMod.DefineSubtype(t.Name, t)
		typeParamMod.DefineConstant(t.Name, types.NoValue{})
	}

	prevMode := c.mode
	c.mode = namedGenericTypeDefinitionMode
	c.pushConstScope(makeConstantScope(typeParamMod))

	node.TypeNode = c.checkTypeNode(node.TypeNode)
	typ := c.typeOf(node.TypeNode)
	namedType.Type = typ
	namedType.TypeParameters = typeParams

	c.mode = prevMode
	c.popConstScope()

	return true
}

func (c *Checker) checkTypeDefinitions() {
	oldFilename := c.Filename
	oldConstantScopes := c.constantScopes
	for _, typeName := range c.typeDefinitionChecks.order {
		typedefCheck := c.typeDefinitionChecks.m[typeName]
		c.checkTypeDefinition(typedefCheck, nil)
	}
	c.Filename = oldFilename
	c.constantScopes = oldConstantScopes
	c.typeDefinitionChecks = newTypeDefinitionChecks()
}

func (c *Checker) checkTypeDefinition(typedefCheck *typeDefinitionCheck, span *position.Span) bool {
	if typedefCheck.state == CHECKING_TYPEDEF {
		c.addFailure(
			fmt.Sprintf("type `%s` circularly references itself", types.InspectWithColor(typedefCheck.typ)),
			span,
		)
		return false
	}
	if typedefCheck.state == CHECKED_TYPEDEF {
		return true
	}

	typedefCheck.state = CHECKING_TYPEDEF

	for _, entry := range typedefCheck.entries {
		c.Filename = entry.filename
		c.constantScopes = entry.constantScopes
		switch n := entry.node.(type) {
		case *ast.TypeDefinitionNode:
			c.checkNamedType(n)
		case *ast.GenericTypeDefinitionNode:
			c.checkGenericNamedType(n)
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
		case *ast.MixinDeclarationNode:
			c.checkMixinTypeParameters(n)
		}
	}

	typedefCheck.state = CHECKED_TYPEDEF
	return true
}

func (c *Checker) includeMixin(node ast.ComplexConstantNode) {
	n := c.checkComplexConstantType(node)
	constantType := c.typeOf(n)

	if types.IsNothing(constantType) || constantType == nil {
		return
	}
	var constantNamespace types.Namespace
	switch con := constantType.(type) {
	case *types.Mixin:
		constantNamespace = con
	case *types.Generic:
		if _, ok := con.Namespace.(*types.Mixin); !ok {
			c.addFailure(
				"only mixins can be included",
				node.Span(),
			)
			return
		}
		constantNamespace = con
	default:
		c.addFailure(
			"only mixins can be included",
			node.Span(),
		)
		return
	}
	node.SetType(constantNamespace)

	target := c.currentConstScope().container
	if c.isSubtype(target, constantNamespace, nil) {
		return
	}

	if target.IsPrimitive() && types.NamespaceDeclaresInstanceVariables(constantNamespace) {
		c.addFailure(
			fmt.Sprintf(
				"cannot include mixin with instance variables `%s` in primitive `%s`",
				types.InspectWithColor(constantType),
				types.InspectWithColor(target),
			),
			node.Span(),
		)
	}

	switch t := target.(type) {
	case *types.Class:
		t.IncludeGenericMixin(constantNamespace)
	case *types.SingletonClass:
		t.IncludeGenericMixin(constantNamespace)
	case *types.Mixin:
		t.IncludeGenericMixin(constantNamespace)
	default:
		c.addFailure(
			fmt.Sprintf(
				"cannot include `%s` in `%s`",
				types.InspectWithColor(constantType),
				types.InspectWithColor(t),
			),
			node.Span(),
		)
	}
}

func (c *Checker) implementInterface(node ast.ComplexConstantNode) {
	constantType, _ := c.resolveConstantType(node)

	if types.IsNothing(constantType) || constantType == nil {
		return
	}
	constantInterface, constantIsInterface := constantType.(*types.Interface)
	if !constantIsInterface {
		c.addFailure(
			"only interfaces can be implemented",
			node.Span(),
		)
		return
	}

	target := c.currentConstScope().container
	if c.implementsInterface(target, constantInterface) {
		return
	}

	switch t := target.(type) {
	case *types.Class:
		t.ImplementInterface(constantInterface)
	case *types.Mixin:
		t.ImplementInterface(constantInterface)
	case *types.Interface:
		t.ImplementInterface(constantInterface)
	default:
		c.addFailure(
			fmt.Sprintf(
				"cannot implement `%s` in `%s`",
				types.InspectWithColor(constantType),
				types.InspectWithColor(t),
			),
			node.Span(),
		)
	}
}

func (c *Checker) checkMixinTypeParameters(node *ast.MixinDeclarationNode) {
	mixin, ok := c.typeOf(node).(*types.Mixin)
	if !ok {
		return
	}
	c.pushConstScope(makeLocalConstantScope(mixin))

	if len(node.TypeParameters) > 0 {
		typeParams := make([]*types.TypeParameter, 0, len(node.TypeParameters))
		for _, typeParamNode := range node.TypeParameters {
			varNode, ok := typeParamNode.(*ast.VariantTypeParameterNode)
			if !ok {
				continue
			}

			t := c.checkTypeParameterNode(varNode)
			typeParams = append(typeParams, t)
			typeParamNode.SetType(t)
			mixin.DefineSubtype(t.Name, t)
			mixin.DefineConstant(t.Name, types.NoValue{})
		}

		mixin.TypeParameters = typeParams
	}

	c.popConstScope()
}

func (c *Checker) checkClassInheritance(node *ast.ClassDeclarationNode) {
	class, ok := c.typeOf(node).(*types.Class)
	if !ok {
		return
	}
	c.pushConstScope(makeLocalConstantScope(class))

	if len(node.TypeParameters) > 0 {
		typeParams := make([]*types.TypeParameter, 0, len(node.TypeParameters))
		for _, typeParamNode := range node.TypeParameters {
			varNode, ok := typeParamNode.(*ast.VariantTypeParameterNode)
			if !ok {
				continue
			}

			t := c.checkTypeParameterNode(varNode)
			typeParams = append(typeParams, t)
			typeParamNode.SetType(t)
			class.DefineSubtype(t.Name, t)
			class.DefineConstant(t.Name, types.NoValue{})
		}

		class.TypeParameters = typeParams
	}

	var superclassType types.Type
	var superclass types.Namespace

nodeSwitch:
	switch node.Superclass.(type) {
	case *ast.NilLiteralNode:
	case nil:
		superclass = c.GlobalEnv.StdSubtypeClass(symbol.Object)
		superclassType = superclass
	default:
		prevMode := c.mode
		c.mode = inheritanceMode

		node.Superclass = c.checkComplexConstantType(node.Superclass)
		superclassType = c.typeOf(node.Superclass)

		c.mode = prevMode

		switch s := superclassType.(type) {
		case *types.Class:
			superclass = s
		case *types.Generic:
			superclass = s
			if _, ok := s.Namespace.(*types.Class); !ok {
				c.addFailure(
					fmt.Sprintf("`%s` is not a class", types.InspectWithColor(superclassType)),
					node.Superclass.Span(),
				)
				break nodeSwitch
			}
		default:
			if !types.IsNothing(superclassType) && superclassType != nil {
				c.addFailure(
					fmt.Sprintf("`%s` is not a class", types.InspectWithColor(superclassType)),
					node.Superclass.Span(),
				)
			}
			break nodeSwitch
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

	c.popConstScope()
}

func (c *Checker) checkTypeParameterNode(node *ast.VariantTypeParameterNode) *types.TypeParameter {
	var variance types.Variance
	switch node.Variance {
	case ast.INVARIANT:
		variance = types.INVARIANT
	case ast.COVARIANT:
		variance = types.COVARIANT
	case ast.CONTRAVARIANT:
		variance = types.CONTRAVARIANT
	}

	var lowerType types.Type = types.Never{}
	if node.LowerBound != nil {
		node.LowerBound = c.checkTypeNode(node.LowerBound)
		lowerType = c.typeOf(node.LowerBound)
	}

	var upperType types.Type = types.Any{}
	if node.UpperBound != nil {
		node.UpperBound = c.checkTypeNode(node.UpperBound)
		upperType = c.typeOf(node.UpperBound)
	}

	return types.NewTypeParameter(
		value.ToSymbol(node.Name),
		lowerType,
		upperType,
		variance,
	)
}
