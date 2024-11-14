package checker

import (
	"fmt"
	"slices"

	"github.com/elk-language/elk/lexer"
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

func (c *Checker) replaceSimpleNamespacePlaceholder(placeholder *types.ConstantPlaceholder, subtype, constant types.Type) {
	placeholder.Replaced = true
	usingConst := placeholder.Container[placeholder.AsName]
	placeholder.Container[placeholder.AsName] = types.Constant{
		FullName: usingConst.FullName,
		Type:     constant,
	}

	placeholder.Sibling.Replaced = true
	subtypeContainer := placeholder.Sibling.Container
	usingSubtype := subtypeContainer[placeholder.AsName]
	subtypeContainer[placeholder.AsName] = types.Constant{
		FullName: usingSubtype.FullName,
		Type:     subtype,
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

func (c *Checker) replaceTypePlaceholder(previousConstantType, newType types.Type) {
	placeholder, ok := previousConstantType.(*types.ConstantPlaceholder)
	if !ok {
		return
	}

	placeholder = placeholder.Sibling
	placeholder.Replaced = true
	usingConst := placeholder.Container[placeholder.AsName]
	placeholder.Container[placeholder.AsName] = types.Constant{
		FullName: usingConst.FullName,
		Type:     newType,
	}
}

func (c *Checker) registerNamedTypeCheck(node *ast.TypeDefinitionNode) {
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	switch constant.(type) {
	case *types.ConstantPlaceholder, nil:
	default:
		c.addRedeclaredConstantError(fullConstantName, node.Constant.Span())
	}

	namedType := types.NewNamedType(fullConstantName, nil)
	container.DefineConstant(constantName, types.NoValue{})
	container.DefineSubtype(constantName, namedType)
	node.SetType(namedType)
	c.replaceTypePlaceholder(constant, namedType)

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
	switch constant.(type) {
	case *types.ConstantPlaceholder, nil:
	default:
		c.addRedeclaredConstantError(fullConstantName, node.Constant.Span())
	}

	namedType := types.NewGenericNamedType(
		fullConstantName,
		nil,
		nil,
	)
	container.DefineConstant(constantName, types.NoValue{})
	container.DefineSubtype(constantName, namedType)
	node.SetType(namedType)
	c.replaceTypePlaceholder(constant, namedType)

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
	c.pushConstScope(makeConstantScope(typeParamMod))
	for _, typeParamNode := range node.TypeParameters {
		varNode, ok := typeParamNode.(*ast.VariantTypeParameterNode)
		if !ok {
			continue
		}

		t := c.checkTypeParameterNode(varNode, typeParamMod, false)
		typeParams = append(typeParams, t)
		typeParamNode.SetType(t)
		typeParamMod.DefineSubtype(t.Name, t)
		typeParamMod.DefineConstant(t.Name, types.NoValue{})
	}

	prevMode := c.mode
	c.mode = namedGenericTypeDefinitionMode

	node.TypeNode = c.checkTypeNode(node.TypeNode)
	typ := c.typeOf(node.TypeNode)
	namedType.Type = typ
	namedType.TypeParameters = typeParams

	c.mode = prevMode
	c.popConstScope()

	return true
}

func (c *Checker) checkTypeDefinitions() {
	for _, typeName := range c.typeDefinitionChecks.order {
		typedefCheck := c.typeDefinitionChecks.m[typeName]
		c.checkTypeDefinition(typedefCheck, nil)
	}
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

	oldFilename := c.Filename
	oldConstantScopes := c.constantScopes
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
			n.SetType(types.Untyped{})
		case *ast.ImplementExpressionNode:
			for _, constant := range n.Constants {
				c.implementInterface(constant)
			}
			n.SetType(types.Untyped{})
		case *ast.ClassDeclarationNode:
			c.checkClassInheritance(n)
		case *ast.MixinDeclarationNode:
			c.checkMixinTypeParameters(n)
		case *ast.InterfaceDeclarationNode:
			c.checkInterfaceTypeParameters(n)
		case *ast.ExtendWhereBlockExpressionNode:
			c.checkExtendWhere(n)
		}
	}
	c.Filename = oldFilename
	c.constantScopes = oldConstantScopes

	typedefCheck.state = CHECKED_TYPEDEF
	return true
}

func (c *Checker) includeMixin(node ast.ComplexConstantNode) {
	prevMode := c.mode
	c.mode = inheritanceMode

	n := c.checkComplexConstantType(node)
	constantType := c.typeOf(n)

	c.mode = prevMode

	if types.IsUntyped(constantType) || constantType == nil {
		return
	}
	target := c.currentConstScope().container

	var constantNamespace types.Namespace
	var mixin *types.Mixin
	switch con := constantType.(type) {
	case *types.Mixin:
		constantNamespace = con
		mixin = con
	case *types.Generic:
		var ok bool
		mixin, ok = con.Namespace.(*types.Mixin)
		if !ok {
			c.addFailure(
				"only mixins can be included",
				node.Span(),
			)
			return
		}
		isAlreadyIncluded, includedNamespace := c.includesMixin(target, mixin)
		includedGeneric, isGeneric := includedNamespace.(*types.Generic)
		if isAlreadyIncluded && isGeneric && !c.isTheSameType(con, includedGeneric, nil) {
			c.addFailure(
				fmt.Sprintf(
					"cannot include mixin `%s` since `%s` has already been included",
					types.InspectWithColor(con),
					types.InspectWithColor(includedGeneric),
				),
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

	if c.isSubtypeOfMixin(target, mixin) {
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
		types.IncludeMixin(t, constantNamespace)
	case *types.SingletonClass:
		types.IncludeMixin(t, constantNamespace)
	case *types.Mixin:
		types.IncludeMixin(t, constantNamespace)
	default:
		c.addFailure(
			fmt.Sprintf(
				"cannot include `%s` in `%s`",
				types.InspectWithColor(constantType),
				types.InspectWithColor(t),
			),
			node.Span(),
		)
		return
	}

	if c.shouldCompile() {
		c.compiler.CompileInclude(target, mixin, position.DefaultSpan)
	}
}

func (c *Checker) implementInterface(node ast.ComplexConstantNode) {
	prevMode := c.mode
	c.mode = inheritanceMode

	n := c.checkComplexConstantType(node)
	constantType := c.typeOf(n)

	c.mode = prevMode

	if types.IsUntyped(constantType) || constantType == nil {
		return
	}
	var constantNamespace types.Namespace
	switch con := constantType.(type) {
	case *types.Interface:
		constantNamespace = con
	case *types.Generic:
		if _, ok := con.Namespace.(*types.Interface); !ok {
			c.addFailure(
				"only interfaces can be implemented",
				node.Span(),
			)
			return
		}
		constantNamespace = con
	default:
		c.addFailure(
			"only interfaces can be implemented",
			node.Span(),
		)
		return
	}

	target := c.currentConstScope().container

	switch t := target.(type) {
	case *types.Class:
		types.ImplementInterface(t, constantNamespace)
	case *types.Mixin:
		types.ImplementInterface(t, constantNamespace)
	case *types.Interface:
		types.ImplementInterface(t, constantNamespace)
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

func (c *Checker) checkInterfaceTypeParameters(node *ast.InterfaceDeclarationNode) {
	iface, ok := c.typeOf(node).(*types.Interface)
	if !ok {
		return
	}
	c.pushConstScope(makeLocalConstantScope(iface))

	typeParams := c.checkNamespaceTypeParameters(
		iface.Checked,
		node.TypeParameters,
		iface,
		iface.TypeParameters(),
		node.Span(),
	)
	if typeParams != nil {
		iface.SetTypeParameters(typeParams)
	}
	iface.Checked = true

	c.popConstScope()
}

func (c *Checker) checkMixinTypeParameters(node *ast.MixinDeclarationNode) {
	mixin, ok := c.typeOf(node).(*types.Mixin)
	if !ok {
		return
	}
	c.pushConstScope(makeLocalConstantScope(mixin))

	typeParams := c.checkNamespaceTypeParameters(
		mixin.Checked,
		node.TypeParameters,
		mixin,
		mixin.TypeParameters(),
		node.Span(),
	)
	if typeParams != nil {
		mixin.SetTypeParameters(typeParams)
	}
	mixin.Checked = true

	c.popConstScope()
}

func (c *Checker) checkNamespaceTypeParameters(
	checked bool,
	typeParamNodes []ast.TypeParameterNode,
	namespace types.Namespace,
	oldTypeParams []*types.TypeParameter,
	span *position.Span,
) []*types.TypeParameter {
	prevMode := c.mode
	c.mode = inheritanceMode
	defer c.setMode(prevMode)

	if !checked {
		if len(typeParamNodes) > 0 {
			typeParams := make([]*types.TypeParameter, 0, len(typeParamNodes))
			for _, typeParamNode := range typeParamNodes {
				varNode, ok := typeParamNode.(*ast.VariantTypeParameterNode)
				if !ok {
					continue
				}

				t := c.initTypeParameterNode(varNode, namespace)
				typeParams = append(typeParams, t)
				typeParamNode.SetType(t)
				namespace.DefineSubtype(t.Name, t)
				namespace.DefineConstant(t.Name, types.NoValue{})
				c.finishCheckingTypeParameterNode(t, varNode)
			}

			return typeParams
		}
		return nil
	}

	if len(typeParamNodes) != len(oldTypeParams) {
		c.addFailure(
			fmt.Sprintf(
				"type parameter count mismatch in `%s`, got: %d, expected: %d",
				types.InspectWithColor(namespace),
				len(typeParamNodes),
				len(oldTypeParams),
			),
			span,
		)
		return nil
	}

	for i := range len(oldTypeParams) {
		typeParamNode := typeParamNodes[i]
		oldTypeParam := oldTypeParams[i]

		varNode, ok := typeParamNode.(*ast.VariantTypeParameterNode)
		if !ok {
			continue
		}

		newTypeParam := c.checkTypeParameterNode(varNode, namespace, false)
		typeParamNode.SetType(newTypeParam)

		if newTypeParam.Name != oldTypeParam.Name ||
			newTypeParam.Variance != oldTypeParam.Variance ||
			!c.isTheSameType(newTypeParam.LowerBound, oldTypeParam.LowerBound, nil) ||
			!c.isTheSameType(newTypeParam.UpperBound, oldTypeParam.UpperBound, nil) {
			c.addFailure(
				fmt.Sprintf(
					"type parameter mismatch in `%s`, is `%s`, should be `%s`",
					types.InspectWithColor(namespace),
					newTypeParam.InspectSignature(),
					oldTypeParam.InspectSignature(),
				),
				span,
			)
		}
	}

	return nil
}

func (c *Checker) checkClassInheritance(node *ast.ClassDeclarationNode) {
	class, ok := c.typeOf(node).(*types.Class)
	if !ok {
		return
	}
	c.pushConstScope(makeLocalConstantScope(class))
	typeParams := c.checkNamespaceTypeParameters(
		class.Checked,
		node.TypeParameters,
		class,
		class.TypeParameters(),
		node.Span(),
	)
	if typeParams != nil {
		class.SetTypeParameters(typeParams)
	}

	var superclassType types.Type
	var superclass types.Namespace

superclassSwitch:
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
				break superclassSwitch
			}
		default:
			if !types.IsUntyped(superclassType) && superclassType != nil {
				c.addFailure(
					fmt.Sprintf("`%s` is not a class", types.InspectWithColor(superclassType)),
					node.Superclass.Span(),
				)
			}
			break superclassSwitch
		}

		if superclass.IsSealed() && !c.IsHeader {
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

	var previousSuperclass types.Type = class.Superclass()
	if !class.Checked && previousSuperclass == nil && superclass != nil {
		class.SetParent(superclass)
	} else if !c.isTheSameType(previousSuperclass, superclass, nil) {
		var span *position.Span
		if node.Superclass == nil {
			span = node.Span()
		} else {
			span = node.Superclass.Span()
		}

		if previousSuperclass == nil {
			previousSuperclass = types.Nil{}
		}

		c.addFailure(
			fmt.Sprintf(
				"superclass mismatch in `%s`, got `%s`, expected `%s`",
				types.InspectWithColor(class),
				types.InspectWithColor(superclassType),
				types.InspectWithColor(previousSuperclass),
			),
			span,
		)
	}
	class.Checked = true
	if c.shouldCompile() {
		c.compiler.CompileClassInheritance(class, position.DefaultSpan)
	}

	c.popConstScope()
}

func (c *Checker) checkExtendWhere(node *ast.ExtendWhereBlockExpressionNode) {
	currentNamespace := c.currentConstScope().container
	if !currentNamespace.IsGeneric() {
		c.addFailure(
			fmt.Sprintf(
				"cannot use `%s` since namespace `%s` is not generic",
				lexer.Colorize("extend where"),
				types.InspectWithColor(currentNamespace),
			),
			node.Span(),
		)
		node.SetType(types.Untyped{})
		return
	}

	mixin := types.NewMixin("", false, "", c.GlobalEnv)
	originalTypeParams := currentNamespace.TypeParameters()
	for _, typeParam := range originalTypeParams {
		mixin.DefineSubtypeWithFullName(
			typeParam.Name,
			fmt.Sprintf("%s::%s", currentNamespace.Name(), typeParam.Name.String()),
			typeParam,
		)
	}
	mixin.SetInstanceVariables(currentNamespace.InstanceVariables())

	prevMode := c.mode
	c.mode = inheritanceMode
	var where []*types.TypeParameter
	for _, whereTypeParamNode := range node.Where {
		whereTypeParamNode := whereTypeParamNode.(*ast.VariantTypeParameterNode)
		whereTypeParam := c.checkTypeParameterNode(whereTypeParamNode, mixin, true)
		originalTypeParamIndex := slices.IndexFunc(
			originalTypeParams,
			func(tp *types.TypeParameter) bool {
				return tp.Name == whereTypeParam.Name
			},
		)
		if originalTypeParamIndex == -1 {
			c.addFailure(
				fmt.Sprintf(
					"cannot add where constraints to nonexistent type parameter `%s`",
					lexer.Colorize(whereTypeParamNode.Name),
				),
				whereTypeParamNode.Span(),
			)
			continue
		}
		originalTypeParam := originalTypeParams[originalTypeParamIndex]

		var newLowerBound types.Type
		if whereTypeParam.LowerBound == nil {
			newLowerBound = originalTypeParam.LowerBound
		} else {
			if !c.isSubtype(originalTypeParam.LowerBound, whereTypeParam.LowerBound, nil) {
				c.addFailure(
					fmt.Sprintf(
						"type parameter `%s` in where clause should have a wider lower bound, has `%s`, should have `%s` or its supertype",
						lexer.Colorize(whereTypeParamNode.Name),
						types.InspectWithColor(whereTypeParam.LowerBound),
						types.InspectWithColor(originalTypeParam.LowerBound),
					),
					whereTypeParamNode.Span(),
				)
				continue
			}
			newLowerBound = whereTypeParam.LowerBound
		}

		var newUpperBound types.Type
		if whereTypeParam.UpperBound == nil {
			newUpperBound = originalTypeParam.UpperBound
		} else {
			if !c.isSubtype(whereTypeParam.UpperBound, originalTypeParam.UpperBound, nil) {
				c.addFailure(
					fmt.Sprintf(
						"type parameter `%s` in where clause should have a narrower upper bound, has `%s`, should have `%s` or its subtype",
						lexer.Colorize(whereTypeParamNode.Name),
						types.InspectWithColor(whereTypeParam.UpperBound),
						types.InspectWithColor(originalTypeParam.UpperBound),
					),
					whereTypeParamNode.Span(),
				)
				continue
			}
			newUpperBound = whereTypeParam.UpperBound
		}

		if whereTypeParam.Variance != types.INVARIANT {
			c.addFailure(
				fmt.Sprintf(
					"cannot modify the variance of type parameter `%s` in a where clause",
					lexer.Colorize(whereTypeParamNode.Name),
				),
				whereTypeParamNode.Span(),
			)
			continue
		}
		whereTypeParam.LowerBound = newLowerBound
		whereTypeParam.UpperBound = newUpperBound
		where = append(where, whereTypeParam)

		narrowerTypeParam := originalTypeParam.Copy()
		narrowerTypeParam.LowerBound = newLowerBound
		narrowerTypeParam.UpperBound = newUpperBound
		mixin.DefineSubtypeWithFullName(
			whereTypeParam.Name,
			fmt.Sprintf("%s::%s", currentNamespace.Name(), whereTypeParam.Name.String()),
			narrowerTypeParam,
		)
	}
	c.mode = prevMode

	mixinWithWhere := types.IncludeMixinWithWhere(currentNamespace, mixin, where)
	node.SetType(mixinWithWhere)
}

func (c *Checker) checkTypeParameterNode(node *ast.VariantTypeParameterNode, namespace types.Namespace, leaveNil bool) *types.TypeParameter {
	var variance types.Variance
	switch node.Variance {
	case ast.INVARIANT:
		variance = types.INVARIANT
	case ast.COVARIANT:
		variance = types.COVARIANT
	case ast.CONTRAVARIANT:
		variance = types.CONTRAVARIANT
	}

	var lowerType types.Type
	if node.LowerBound != nil {
		node.LowerBound = c.checkTypeNode(node.LowerBound)
		lowerType = c.typeOf(node.LowerBound)
	} else if !leaveNil {
		lowerType = types.Never{}
	}

	var upperType types.Type
	if node.UpperBound != nil {
		node.UpperBound = c.checkTypeNode(node.UpperBound)
		upperType = c.typeOf(node.UpperBound)
	} else if !leaveNil {
		upperType = types.Any{}
	}

	return types.NewTypeParameter(
		value.ToSymbol(node.Name),
		namespace,
		lowerType,
		upperType,
		variance,
	)
}

func (c *Checker) initTypeParameterNode(node *ast.VariantTypeParameterNode, namespace types.Namespace) *types.TypeParameter {
	var variance types.Variance
	switch node.Variance {
	case ast.INVARIANT:
		variance = types.INVARIANT
	case ast.COVARIANT:
		variance = types.COVARIANT
	case ast.CONTRAVARIANT:
		variance = types.CONTRAVARIANT
	}

	return types.NewTypeParameter(
		value.ToSymbol(node.Name),
		namespace,
		nil,
		nil,
		variance,
	)
}

func (c *Checker) finishCheckingTypeParameterNode(typ *types.TypeParameter, node *ast.VariantTypeParameterNode) {
	var lowerType types.Type
	if node.LowerBound != nil {
		node.LowerBound = c.checkTypeNode(node.LowerBound)
		lowerType = c.typeOf(node.LowerBound)
	} else {
		lowerType = types.Never{}
	}

	var upperType types.Type
	if node.UpperBound != nil {
		node.UpperBound = c.checkTypeNode(node.UpperBound)
		upperType = c.typeOf(node.UpperBound)
	} else {
		upperType = types.Any{}
	}

	typ.LowerBound = lowerType
	typ.UpperBound = upperType
}
