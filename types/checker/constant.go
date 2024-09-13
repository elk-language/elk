package checker

import (
	"fmt"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

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
	filename       string
	constantScopes []constantScope
	methodScopes   []methodScope
	node           *ast.ConstantDeclarationNode
	state          constState
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

func (c *Checker) hoistConstantDeclaration(node *ast.ConstantDeclarationNode) {
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	if constant != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			node.Span(),
		)
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
		return
	}

	if node.TypeNode == nil {
		c.addFailure(
			"non-static constants must have an explicit type",
			node.Span(),
		)
		node.SetType(types.Nothing{})
		return
	}

	node.TypeNode = c.checkTypeNode(node.TypeNode)
	declaredType := c.typeOf(node.TypeNode)
	container.DefineConstant(constantName, declaredType)
	node.SetType(declaredType)
	c.registerConstantCheck(fullConstantName, node)
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
		return true
	}
	check.state = CHECKING_CONST

	node := check.node
	declaredType := c.typeOf(node.TypeNode)
	node.Initialiser = c.checkExpression(node.Initialiser)
	init := node.Initialiser
	actualType := c.typeOfGuardVoid(init)
	c.checkCanAssign(actualType, declaredType, init.Span())

	check.state = CHECKED_CONST
	return true
}

func (c *Checker) constantLookupType(node *ast.ConstantLookupNode) *ast.PublicConstantNode {
	typ, name := c.resolveConstantLookupType(node)
	switch t := typ.(type) {
	case *types.GenericNamedType:
		c.addTypeArgumentCountError(types.InspectWithColor(typ), len(t.TypeParameters), 0, node.Span())
		typ = types.Nothing{}
	case *types.Class:
		if t.IsGeneric() {
			c.addTypeArgumentCountError(types.InspectWithColor(typ), len(t.TypeParameters()), 0, node.Span())
			typ = types.Nothing{}
		}
	case nil:
		typ = types.Nothing{}
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
	if types.IsNoValue(constant.Type) {
		c.addFailure(
			fmt.Sprintf("type `%s` cannot be used as a value in expressions", lexer.Colorize(constantName)),
			node.Right.Span(),
		)
		return nil, constantName
	}

	if !c.checkConstantIfNecessary(constantName, node.Right.Span()) {
		return types.Nothing{}, constantName
	}
	return constant.Type, constantName
}

func (c *Checker) checkConstantLookupNode(node *ast.ConstantLookupNode) *ast.PublicConstantNode {
	typ, name := c.resolveConstantLookup(node, node.Span())
	if typ == nil {
		typ = types.Nothing{}
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
		typ = types.Nothing{}
	}

	node.Value = name
	node.SetType(typ)
	return node
}

func (c *Checker) checkPrivateConstantNode(node *ast.PrivateConstantNode) *ast.PrivateConstantNode {
	typ, name := c.resolvePrivateConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Nothing{}
	}

	node.Value = name
	node.SetType(typ)
	return node
}
