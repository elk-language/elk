package checker

import (
	"fmt"
	"slices"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

func (c *Checker) checkConstantPlaceholders() {
	for _, placeholder := range c.constantPlaceholders {
		if placeholder.Checked || placeholder.Sibling != nil && placeholder.Sibling.Checked {
			continue
		}
		placeholder.Checked = true
		if placeholder.Replaced || placeholder.Sibling != nil && placeholder.Sibling.Replaced {
			continue
		}

		c.addFailureWithLocation(
			fmt.Sprintf("undefined type or constant `%s`", lexer.Colorize(placeholder.FullName)),
			placeholder.Location,
		)
	}
	c.constantPlaceholders = nil
}

type constantDefinitionChecks struct {
	m     map[string]*constantDefinitionCheck
	order []string
}

func newConstantDefinitionChecks() *constantDefinitionChecks {
	return &constantDefinitionChecks{
		m: make(map[string]*constantDefinitionCheck),
	}
}

type constState uint8

const (
	NEW_CONST constState = iota
	CHECKING_CONST
	CHECKED_CONST
)

type constantDefinitionCheck struct {
	filename          string
	constantScopes    []constantScope
	methodScopes      []methodScope
	node              *ast.ConstantDeclarationNode
	state             constState
	referencedMethods []*types.Method
}

func (c *Checker) registerConstantCheck(name string, node *ast.ConstantDeclarationNode) {
	c.constantChecks.m[name] = &constantDefinitionCheck{
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
		node:           node,
		filename:       c.Filename,
	}
	c.constantChecks.order = append(c.constantChecks.order, name)
}

func (c *Checker) addRedeclaredConstantError(name string, span *position.Span) {
	c.addFailure(
		fmt.Sprintf("cannot redeclare constant `%s`", lexer.Colorize(name)),
		span,
	)
}

func (c *Checker) replaceConstantPlaceholder(previousConstantType, newType types.Type) {
	placeholder, ok := previousConstantType.(*types.ConstantPlaceholder)
	if !ok {
		return
	}

	placeholder.Replaced = true
	usingConst := placeholder.Container[placeholder.AsName]
	placeholder.Container[placeholder.AsName] = types.Constant{
		FullName: usingConst.FullName,
		Type:     newType,
	}
}

func (c *Checker) hoistConstantDeclaration(node *ast.ConstantDeclarationNode) {
	switch c.mode {
	case topLevelMode, moduleMode,
		classMode, mixinMode,
		interfaceMode, singletonMode:
	default:
		c.addFailure(
			"constants cannot be declared in this context",
			node.Span(),
		)
		return
	}
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	switch constant.(type) {
	case *types.ConstantPlaceholder, nil:
	default:
		c.addRedeclaredConstantError(fullConstantName, node.Span())
	}

	if node.Initialiser.IsStatic() {
		node.Initialiser = c.checkExpression(node.Initialiser)
		init := node.Initialiser
		actualType := c.typeOfGuardVoid(init)

		typ := actualType
		if node.TypeNode != nil {
			node.TypeNode = c.checkTypeNode(node.TypeNode)
			declaredType := c.typeOf(node.TypeNode)
			c.checkCanAssign(actualType, declaredType, init.Span())
			typ = declaredType
		}

		container.DefineConstant(constantName, actualType)
		node.SetType(typ)
		c.replaceConstantPlaceholder(constant, typ)
		return
	}

	if node.TypeNode == nil {
		c.addFailure(
			"non-static constants must have an explicit type",
			node.Span(),
		)
		node.SetType(types.Untyped{})
		return
	}

	node.TypeNode = c.checkTypeNode(node.TypeNode)
	declaredType := c.typeOf(node.TypeNode)
	container.DefineConstant(constantName, declaredType)
	node.SetType(declaredType)
	c.registerConstantCheck(fullConstantName, node)
	c.replaceConstantPlaceholder(constant, declaredType)
}

func (c *Checker) checkConstants() {
	prevFilename := c.Filename
	prevConstScopes := c.constantScopes
	prevMethodScopes := c.methodScopes
	prevPhase := c.phase

	c.phase = constantCheckPhase
	for _, constName := range c.constantChecks.order {
		constCheck := c.constantChecks.m[constName]
		c.Filename = constCheck.filename
		c.constantScopes = constCheck.constantScopes
		c.methodScopes = constCheck.methodScopes
		c.checkConstantDeclaration(constName, constCheck, nil)
		c.methodCache.Slice = nil // reset the method cache
	}

	c.phase = prevPhase
	c.Filename = prevFilename
	c.constantScopes = prevConstScopes
	c.methodScopes = prevMethodScopes
	c.constantChecks = newConstantDefinitionChecks()
}

func (c *Checker) checkConstantIfNecessary(name string, span *position.Span) (ok bool) {
	if c.phase != constantCheckPhase {
		return true
	}
	check, ok := c.constantChecks.m[name]
	if !ok {
		return true
	}

	return c.checkConstantDeclaration(name, check, span)
}

func (c *Checker) checkConstantDeclaration(name string, check *constantDefinitionCheck, span *position.Span) bool {
	switch check.state {
	case CHECKING_CONST:
		c.addFailure(
			fmt.Sprintf("constant `%s` circularly references itself", lexer.Colorize(name)),
			span,
		)
		return false
	case CHECKED_CONST:
		c.methodCache.AppendUnsafe(check.referencedMethods...)
		return true
	}
	check.state = CHECKING_CONST

	node := check.node
	declaredType := c.typeOf(node.TypeNode)
	node.Initialiser = c.checkExpression(node.Initialiser)
	init := node.Initialiser
	actualType := c.typeOfGuardVoid(init)
	c.checkCanAssign(actualType, declaredType, init.Span())

	symbolName := value.ToSymbol(name)
	check.referencedMethods = slices.Clone(c.methodCache.Slice)
	for _, method := range c.methodCache.Slice {
		method.UsedInConstants.Add(symbolName)
	}

	check.state = CHECKED_CONST
	return true
}

func (c *Checker) constantLookupType(node *ast.ConstantLookupNode) *ast.PublicConstantNode {
	typ, name := c.resolveConstantLookupType(node)
	switch t := typ.(type) {
	case *types.GenericNamedType:
		c.addTypeArgumentCountError(types.InspectWithColor(typ), len(t.TypeParameters), 0, node.Span())
		typ = types.Untyped{}
	case *types.Class:
		if t.IsGeneric() {
			c.addTypeArgumentCountError(types.InspectWithColor(typ), len(t.TypeParameters()), 0, node.Span())
			typ = types.Untyped{}
		}
	case nil:
		typ = types.Untyped{}
	}

	newNode := ast.NewPublicConstantNode(
		node.Span(),
		name,
	)
	newNode.SetType(typ)
	return newNode
}

func (c *Checker) resolveConstantType(constantExpression ast.ExpressionNode) (types.Type, string) {
	switch constant := constantExpression.(type) {
	case *ast.PublicConstantNode:
		return c.resolveType(constant.Value, constant.Span())
	case *ast.PrivateConstantNode:
		return c.resolveType(constant.Value, constant.Span())
	case *ast.ConstantLookupNode:
		return c.resolveConstantLookupType(constant)
	case *ast.GenericConstantNode:
		typeNode, name := c.checkGenericConstantType(constant)
		return c.typeOf(typeNode), name
	default:
		panic(fmt.Sprintf("invalid constant node: %T", constantExpression))
	}
}

func (c *Checker) addInvalidValueInExpressionError(constantName string, span *position.Span) {
	c.addFailure(
		fmt.Sprintf("`%s` cannot be used as a value in expressions", lexer.Colorize(constantName)),
		span,
	)
}

func (c *Checker) resolveConstantLookup(node *ast.ConstantLookupNode, span *position.Span) (types.Type, string) {
	var leftContainerType types.Type
	var leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		leftContainerType, leftContainerName = c.resolvePublicConstant(l.Value, l.Span())
	case *ast.PrivateConstantNode:
		leftContainerType, leftContainerName = c.resolvePrivateConstant(l.Value, l.Span())
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		leftContainerType, leftContainerName = c.resolveConstantLookup(l, span)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	var rightName string
	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		rightName = r.Value
	case *ast.PrivateConstantNode:
		rightName = r.Value
		c.addFailure(
			fmt.Sprintf("cannot read private constant `%s`", rightName),
			node.Span(),
		)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	constantName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, constantName
	}

	var leftContainer types.Namespace
	switch l := leftContainerType.(type) {
	case *types.Module:
		leftContainer = l
	case *types.SingletonClass:
		leftContainer = l.AttachedObject
	default:
		c.addFailure(
			fmt.Sprintf("cannot read constants from `%s`, it is not a constant container", leftContainerName),
			node.Span(),
		)
		return nil, constantName
	}

	constant, ok := leftContainer.ConstantString(rightName)
	if !ok {
		c.addFailure(
			fmt.Sprintf("undefined constant `%s`", constantName),
			node.Right.Span(),
		)
		return nil, constantName
	}
	if len(constant.FullName) > 0 {
		constantName = constant.FullName
	}
	if types.IsNoValue(constant.Type) || types.IsConstantPlaceholder(constant.Type) {
		c.addInvalidValueInExpressionError(constantName, node.Right.Span())
		return nil, constantName
	}

	if !c.checkConstantIfNecessary(constantName, node.Right.Span()) {
		return types.Untyped{}, constantName
	}
	return constant.Type, constantName
}

// Get the type of the public constant with the given name
func (c *Checker) resolvePublicConstant(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		constant, ok := constScope.container.ConstantString(name)
		if !ok {
			continue
		}

		var fullName string
		if len(constant.FullName) > 0 {
			fullName = constant.FullName
		} else {
			fullName = types.MakeFullConstantName(constScope.container.Name(), name)
		}
		if !c.checkConstantIfNecessary(fullName, span) {
			return nil, fullName
		}

		if types.IsNoValue(constant.Type) || types.IsConstantPlaceholder(constant.Type) {
			c.addInvalidValueInExpressionError(fullName, span)
			return nil, fullName
		}
		return constant.Type, fullName

	}

	c.addFailure(
		fmt.Sprintf("undefined constant `%s`", lexer.Colorize(name)),
		span,
	)
	return nil, name
}

// Get the type of the private constant with the given name
func (c *Checker) resolvePrivateConstant(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		if constScope.kind != scopeLocalKind {
			continue
		}
		constant, ok := constScope.container.ConstantString(name)
		if !ok {
			continue
		}

		var fullName string
		if len(constant.FullName) > 0 {
			fullName = constant.FullName
		} else {
			fullName = types.MakeFullConstantName(constScope.container.Name(), name)
		}
		if !c.checkConstantIfNecessary(fullName, span) {
			return nil, fullName
		}

		if types.IsNoValue(constant.Type) || types.IsConstantPlaceholder(constant.Type) {
			c.addInvalidValueInExpressionError(fullName, span)
			return nil, fullName
		}
		return constant.Type, fullName

	}

	c.addFailure(
		fmt.Sprintf("undefined constant `%s`", name),
		span,
	)
	return nil, name
}

func (c *Checker) addToConstantCache(name value.Symbol) {
	if c.phase == methodCheckPhase {
		c.constantCache.Add(name)
	}
}

func (c *Checker) checkConstantLookupNode(node *ast.ConstantLookupNode) *ast.PublicConstantNode {
	typ, name := c.resolveConstantLookup(node, node.Span())
	if typ == nil {
		typ = types.Untyped{}
	} else {
		c.addToConstantCache(value.ToSymbol(name))
	}

	newNode := ast.NewPublicConstantNode(
		node.Span(),
		name,
	)
	newNode.SetType(typ)
	return newNode
}

func (c *Checker) checkPublicConstantNode(node *ast.PublicConstantNode) *ast.PublicConstantNode {
	typ, name := c.resolvePublicConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Untyped{}
	} else {
		c.addToConstantCache(value.ToSymbol(name))
	}

	node.Value = name
	node.SetType(typ)
	return node
}

func (c *Checker) checkPrivateConstantNode(node *ast.PrivateConstantNode) *ast.PrivateConstantNode {
	typ, name := c.resolvePrivateConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Untyped{}
	} else {
		c.addToConstantCache(value.ToSymbol(name))
	}

	node.Value = name
	node.SetType(typ)
	return node
}
