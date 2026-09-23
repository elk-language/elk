package checker

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func (c *Checker) checkAllMethodSignaturesOfNamespace(namespace types.Namespace) {
	if c.signatureChecks == nil {
		return
	}

	namespaceName := namespace.Name()
	methodSignatures := c.signatureChecks.Get(namespaceName)
	c.checkSignatures(methodSignatures)
}

type signatureCheckEntry struct {
	constantScopes []constantScope
	methodScopes   []methodScope
	selfType       types.Ref[types.Type]
	flags          bitfield.BitField16
	mode           mode
	node           ast.Node
}

func (c *Checker) registerSignatureCheck(node ast.Node) {
	namespaceName := c.currentMethodScope().container.Get().Name()
	entry := signatureCheckEntry{
		constantScopes: c.constantScopesCopyWithoutCache(),
		methodScopes:   c.methodScopesCopyWithoutCache(),
		node:           node,
		flags:          c.flags,
		mode:           c.mode,
		selfType:       c.selfType,
	}

	if entries, exists := c.signatureChecks.GetOk(namespaceName); exists {
		*entries = append(*entries, entry)
	} else {
		c.signatureChecks.Insert(namespaceName, &[]signatureCheckEntry{entry})
	}
}

func (c *Checker) checkAllSignatures() {
	for _, methodSignatures := range c.signatureChecks.All() {
		c.checkSignatures(methodSignatures)
	}
	c.signatureChecks = nil
}

func (c *Checker) checkSignatures(methodSignatures *[]signatureCheckEntry) {
	if methodSignatures == nil {
		return
	}
	for _, signatureCheck := range *methodSignatures {
		prevConstantScopes := c.constantScopes
		prevMethodScopes := c.methodScopes
		prevSelfType := c.selfType
		prevFlags := c.flags
		prevMode := c.mode

		c.constantScopes = signatureCheck.constantScopes
		c.methodScopes = signatureCheck.methodScopes
		c.selfType = signatureCheck.selfType
		c.flags = signatureCheck.flags
		c.mode = signatureCheck.mode

		switch node := signatureCheck.node.(type) {
		case *ast.MethodDefinitionNode:
			c.checkSignatureOfMethodDefinition(node)
		case *ast.AliasDeclarationNode:
			c.checkSignatureOfAliasDeclaration(node)
		case *ast.MethodSignatureDefinitionNode:
			c.checkSignatureOfMethodSignatureDefinition(node)
		case *ast.InstanceVariableDeclarationNode:
			c.checkSignatureOfInstanceVariableDeclaration(node)
		case *ast.InstanceValueDeclarationNode:
			c.checkSignatureOfInstanceValueDeclaration(node)
		case *ast.GetterDeclarationNode:
			c.checkSignatureOfGetterDeclaration(node)
		case *ast.SetterDeclarationNode:
			c.checkSignatureOfSetterDeclaration(node)
		case *ast.AttrDeclarationNode:
			c.checkSignatureOfAttrDeclaration(node)
		default:
			panic(fmt.Sprintf("invalid signature definition node: %T", node))
		}

		c.constantScopes = prevConstantScopes
		c.methodScopes = prevMethodScopes
		c.selfType = prevSelfType
		c.flags = prevFlags
		c.mode = prevMode
	}

	*methodSignatures = nil
}

// Gathers all declarations of methods, constants and instance variables
func (c *Checker) hoistMethodDefinitions(statements []ast.StatementNode) {
	for _, statement := range statements {
		stmt, ok := statement.(*ast.ExpressionStatementNode)
		if !ok {
			continue
		}

		expression := stmt.Expression

		switch expr := expression.(type) {
		case *ast.MacroBoundaryNode:
			c.hoistMethodDefinitions(expr.Body)
		case *ast.ConstantDeclarationNode:
			c.hoistConstantDeclaration(expr)
		case *ast.AliasDeclarationNode:
			c.hoistAliasDeclaration(expr)
		case *ast.MethodDefinitionNode:
			c.hoistMethodDefinition(expr)
		case *ast.MethodSignatureDefinitionNode:
			c.hoistMethodSignatureDefinition(expr)
		case *ast.InitDefinitionNode:
			stmt.Expression = c.hoistInitDefinition(expr)
		case *ast.InstanceVariableDeclarationNode:
			c.hoistInstanceVariableDeclaration(expr)
		case *ast.InstanceValueDeclarationNode:
			c.hoistInstanceValueDeclaration(expr)
		case *ast.GetterDeclarationNode:
			c.hoistGetterDeclaration(expr)
		case *ast.SetterDeclarationNode:
			c.hoistSetterDeclaration(expr)
		case *ast.AttrDeclarationNode:
			c.hoistAttrDeclaration(expr)
		case *ast.UsingExpressionNode:
			c.checkUsingExpressionForMethods(expr)
		case *ast.ModuleDeclarationNode:
			c.hoistMethodDefinitionsWithinModule(expr)
		case *ast.ClassDeclarationNode:
			c.hoistMethodDefinitionsWithinClass(expr)
		case *ast.MixinDeclarationNode:
			c.hoistMethodDefinitionsWithinMixin(expr)
		case *ast.InterfaceDeclarationNode:
			c.hoistMethodDefinitionsWithinInterface(expr)
		case *ast.SingletonBlockExpressionNode:
			c.hoistMethodDefinitionsWithinSingleton(expr)
		case *ast.ExtendWhereBlockExpressionNode:
			c.hoistMethodDefinitionsWithinExtendWhere(expr)
		}
	}
}

func (c *Checker) hoistInitDefinition(initNode *ast.InitDefinitionNode) *ast.MethodDefinitionNode {
	switch c.mode {
	case classMode:
	default:
		c.addFailure(
			"init definitions cannot appear outside of classes",
			initNode.Location(),
		)
	}
	newNode := ast.NewMethodDefinitionNode(
		initNode.Location(),
		initNode.DocComment(),
		initNode.Flags.ToBitFlag(),
		ast.NewPublicIdentifierNode(initNode.Location(), "#init"),
		nil,
		initNode.Parameters,
		nil,
		initNode.ThrowType,
		initNode.Body,
	)
	c.registerSignatureCheck(newNode)
	return newNode
}

func (c *Checker) hoistAliasDeclaration(node *ast.AliasDeclarationNode) {
	c.registerSignatureCheck(node)
}

func (c *Checker) checkSignatureOfAliasDeclaration(node *ast.AliasDeclarationNode) {
	node.SetType(types.Untyped{})
	namespace := c.currentMethodScope().container
	for _, entry := range node.Entries {
		c.hoistAliasEntry(entry, namespace.Get())
	}
}

func (c *Checker) hoistAliasEntry(node *ast.AliasDeclarationEntry, namespace types.Namespace) {
	oldName := c.identifierToName(node.OldName)
	oldNameSymbol := symbol.ToSymbol(oldName)
	aliasedMethod := namespace.Method(oldNameSymbol)
	if aliasedMethod == nil {
		c.addMissingMethodError(namespace, oldName, node.Location())
		return
	}
	newName := symbol.ToSymbol(c.identifierToName(node.NewName))
	oldMethod := c.resolveMethodInNamespace(namespace, newName)
	c.checkMethodOverrideWithPlaceholder(aliasedMethod, oldMethod, node.Location())
	c.checkSpecialMethods(newName, aliasedMethod, nil, node.Location())

	if aliasedMethod.IsOverload() {
		c.addFailure(
			fmt.Sprintf("method `%s` with overloads cannot have an alias", aliasedMethod.Name.String()),
			node.Location(),
		)
	}

	alias := aliasedMethod.CreateAlias(newName)
	namespace.SetMethod(newName, alias)
}

func (c *Checker) checkUsingMethodLookupEntryNode(receiverNode ast.ExpressionNode, methodName, asName string, location *position.Location) {
	_, constant, fullConstantName, _ := c.resolveConstantInRoot(receiverNode)
	var namespace types.Namespace

	switch con := constant.(type) {
	case *types.Module:
		namespace = con
	case *types.SingletonClass:
		namespace = con
	default:
		c.addFailure(
			fmt.Sprintf("undefined namespace `%s`", lexer.Colorize(fullConstantName)),
			receiverNode.Location(),
		)
		return
	}

	originalMethodSymbol := symbol.ToSymbol(methodName)
	var newMethodSymbol symbol.Symbol
	if asName != "" {
		newMethodSymbol = symbol.ToSymbol(asName)
	} else {
		newMethodSymbol = symbol.ToSymbol(methodName)
	}

	usingNamespace := c.getUsingBufferNamespace()

	method := namespace.MethodString(methodName)
	if method != nil {
		usingNamespace.SetMethod(newMethodSymbol, method)
		return
	}

	placeholder := types.NewMethodPlaceholder(
		fmt.Sprintf("%s::%s", fullConstantName, methodName),
		newMethodSymbol,
		usingNamespace,
		location,
	)
	c.registerMethodPlaceholder(placeholder)
	namespace.SetMethod(originalMethodSymbol, placeholder)
	usingNamespace.SetMethod(newMethodSymbol, placeholder)
}

func (c *Checker) resolveUsingExpression(node *ast.UsingExpressionNode) {
	if c.isReadonly() {
		node.SetType(types.Untyped{})
		return
	}

	for _, entry := range node.Entries {
		c.resolveUsingEntry(entry, true)
	}
}

func (c *Checker) resolveUsingEntry(entry ast.UsingEntryNode, pushPlaceholderLocation bool) {
	typeRef := types.Ref[types.Namespace](entry.TypeRef())
	typ := c.TypeOf(entry)
	switch t := typ.(type) {
	case *types.Module:
		c.pushConstScope(makeUsingConstantScope(typeRef))
		c.pushMethodScope(makeUsingMethodScope(typeRef))
	case *types.Mixin:
		c.pushConstScope(makeUsingConstantScope(typeRef))
		c.pushMethodScope(makeUsingMethodScope(typeRef))
	case *types.Class:
		c.pushConstScope(makeUsingConstantScope(typeRef))
		c.pushMethodScope(makeUsingMethodScope(typeRef))
	case *types.Interface:
		c.pushConstScope(makeUsingConstantScope(typeRef))
		c.pushMethodScope(makeUsingMethodScope(typeRef))
	case *types.NamespacePlaceholder:
		c.pushConstScope(makeUsingConstantScope(typeRef))
		c.pushMethodScope(makeUsingMethodScope(typeRef))
		if pushPlaceholderLocation {
			t.Locations.Push(entry.Location())
		}
	case *types.UsingBufferNamespace:
		if c.enclosingScopeIsAUsingBuffer() {
			return
		}
		c.pushConstScope(makeUsingBufferConstantScope(typeRef))
		c.pushMethodScope(makeUsingBufferMethodScope(typeRef))
	}
}

func (c *Checker) hoistMethodDefinitionsWithinClass(node *ast.ClassDeclarationNode) {
	typeRef := types.Ref[types.Namespace](node.TypeRef())
	class, ok := c.TypeOf(node).(*types.Class)
	if ok {
		c.pushConstScope(makeLocalConstantScope(typeRef))
		c.pushMethodScope(makeLocalMethodScope(typeRef))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = classMode
	c.selfType = types.CastRef[types.Type](class)
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.registerNamespaceWithIvars(class, node.Location())
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

type namespaceWithIvarsData struct {
	namespace types.NamespaceWithIvarIndices
	locations []*position.Location
}

func (c *namespaceWithIvarsData) addLocation(loc *position.Location) {
	c.locations = append(c.locations, loc)
}

func (c *Checker) insertNamespaceWithIvarsData(namespace types.NamespaceWithIvarIndices) *namespaceWithIvarsData {
	data, ok := c.namespacesWithIvars.GetOk(namespace.Name())
	if ok {
		return data
	}

	data = &namespaceWithIvarsData{namespace: namespace}
	c.namespacesWithIvars.Set(namespace.Name(), data)
	return data
}

func (c *Checker) registerNamespaceWithIvars(namespace types.NamespaceWithIvarIndices, loc *position.Location) {
	if namespace == nil {
		return
	}

	classData := c.insertNamespaceWithIvarsData(namespace)
	classData.addLocation(loc)
}

func (c *Checker) hoistMethodDefinitionsWithinModule(node *ast.ModuleDeclarationNode) {
	typeRef := types.Ref[types.Namespace](node.TypeRef())
	module, ok := c.TypeOf(node).(*types.Module)
	if ok {
		c.pushConstScope(makeLocalConstantScope(typeRef))
		c.pushMethodScope(makeLocalMethodScope(typeRef))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = moduleMode
	c.selfType = types.CastRef[types.Type](module)
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.registerNamespaceWithIvars(module, node.Location())
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinMixin(node *ast.MixinDeclarationNode) {
	typeRef := types.Ref[types.Namespace](node.TypeRef())
	mixin, ok := c.TypeOf(node).(*types.Mixin)
	if ok {
		c.pushConstScope(makeLocalConstantScope(typeRef))
		c.pushMethodScope(makeLocalMethodScope(typeRef))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = mixinMode
	c.selfType = types.CastRef[types.Type](mixin)
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinInterface(node *ast.InterfaceDeclarationNode) {
	typeRef := types.Ref[types.Namespace](node.TypeRef())
	iface, ok := c.TypeOf(node).(*types.Interface)
	if ok {
		c.pushConstScope(makeLocalConstantScope(typeRef))
		c.pushMethodScope(makeLocalMethodScope(typeRef))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = interfaceMode
	c.selfType = types.CastRef[types.Type](iface)
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinSingleton(expr *ast.SingletonBlockExpressionNode) {
	namespace := c.currentConstScope().container.Get()
	singleton := namespace.Singleton()
	if singleton == nil {
		return
	}

	c.pushConstScope(makeLocalConstantScope(types.CastRef[types.Namespace](singleton)))
	c.pushMethodScope(makeLocalMethodScope(types.CastRef[types.Namespace](singleton)))

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = singletonMode
	c.selfType = types.CastRef[types.Type](singleton)
	c.hoistMethodDefinitions(expr.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	c.registerNamespaceWithIvars(singleton, expr.Location())
	c.popLocalConstScope()
	c.popMethodScope()
}

func (c *Checker) hoistMethodDefinitionsWithinExtendWhere(node *ast.ExtendWhereBlockExpressionNode) {
	typeRef := types.Ref[types.Namespace](node.TypeRef())
	_, ok := c.TypeOf(node).(*types.MixinWithWhere)
	if ok {
		c.pushConstScope(makeLocalConstantScope(typeRef))
		c.pushMethodScope(makeLocalMethodScope(typeRef))
	}

	previousMode := c.mode
	c.mode = extendWhereMode
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) checkUsingExpressionForMethods(node *ast.UsingExpressionNode) {
	for _, entry := range node.Entries {
		c.resolveUsingEntry(entry, true)
		switch e := entry.(type) {
		case *ast.MethodLookupNode:
			if e.IsMacro() {
				continue
			}
			c.checkUsingMethodLookupEntryNode(
				e.Receiver,
				c.identifierToName(e.Name),
				"",
				e.Location(),
			)
		case *ast.MethodLookupAsNode:
			if e.IsMacro() {
				continue
			}
			c.checkUsingMethodLookupEntryNode(
				e.MethodLookup.Receiver,
				c.identifierToName(e.MethodLookup.Name),
				c.identifierToName(e.AsName),
				e.Location(),
			)
		case *ast.UsingEntryWithSubentriesNode:
			c.checkUsingEntryWithSubentriesForMethods(e)
		}
	}
}

func (c *Checker) checkUsingEntryWithSubentriesForMethods(node *ast.UsingEntryWithSubentriesNode) {
	for _, subentry := range node.Subentries {
		switch s := subentry.(type) {
		case *ast.PublicIdentifierNode:
			c.checkUsingMethodLookupEntryNode(node.Namespace, c.identifierToName(s), "", s.Location())
		case *ast.UsingSubentryAsNode:
			if s.IsMacro() {
				continue
			}
			value := c.identifierToName(s.Target)
			asName := c.identifierToName(s.AsName)
			c.checkUsingMethodLookupEntryNode(node.Namespace, value, asName, s.Location())
		case *ast.PublicConstantNode, *ast.PublicConstantAsNode, *ast.MacroNameNode:
		default:
			panic(fmt.Sprintf("invalid using subentry node: %T", subentry))
		}
	}
}

func (c *Checker) hoistMethodDefinition(node *ast.MethodDefinitionNode) {
	c.registerSignatureCheck(node)
}

func (c *Checker) checkSignatureOfMethodDefinition(node *ast.MethodDefinitionNode) {
	definedUnder := c.currentMethodScope().container.Get()
	method, mod := c.declareMethod(
		definedUnder,
		node.DocComment(),
		node.IsAbstract(),
		node.IsSealed(),
		false,
		node.IsGenerator(),
		node.IsAsync(),
		node.IsOverload(),
		node.IsPure(),
		symbol.ToSymbol(c.identifierToName(node.Name)),
		node.TypeParameters,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Location(),
	)
	method.Node = node
	node.SetType(method)
	c.registerMethodBodyCheck(method, node)
	if mod != nil {
		c.popConstScope()
	}
}

func (c *Checker) hoistMethodSignatureDefinition(node *ast.MethodSignatureDefinitionNode) {
	c.registerSignatureCheck(node)
}

func (c *Checker) checkSignatureOfMethodSignatureDefinition(node *ast.MethodSignatureDefinitionNode) {
	method, mod := c.declareMethod(
		c.currentMethodScope().container.Get(),
		node.DocComment(),
		true,
		false,
		false,
		false,
		false,
		false,
		false,
		symbol.ToSymbol(c.identifierToName(node.Name)),
		node.TypeParameters,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Location(),
	)
	if mod != nil {
		c.popConstScope()
	}
	node.SetType(method)
}

func (c *Checker) newMethodChecker(
	funcName string,
	constScopes []constantScope,
	methodScopes []methodScope,
	selfType,
	returnType,
	throwType types.Ref[types.Type],
	mode mode,
	threadPool *vm.ThreadPool,
	loc *position.Location,
) *Checker {
	checker := &Checker{
		runtimeEnv:     c.runtimeEnv,
		Filename:       loc.FilePath,
		mode:           mode,
		phase:          methodCheckPhase,
		selfType:       selfType,
		returnType:     returnType,
		throwType:      throwType,
		constantScopes: constScopes,
		methodScopes:   methodScopes,
		Errors:         c.Errors,
		flags:          c.flags,
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil, defaultLocalEnvType),
		},
		typeDefinitionChecks: newTypeDefinitionChecks(),
		methodCache:          concurrent.NewSlice[types.Ref[*types.Method]](),
		threadPool:           threadPool,
	}
	checker.compiler = compiler.CreateCompiler(funcName, c.compiler, checker, loc, c.Errors, c.HasAdditionalAbortChecks())

	return checker
}

// Checks whether all methods specified in `using` statements have been defined
func (c *Checker) checkMethodPlaceholders() {
	for _, placeholderRef := range c.methodPlaceholders {
		placeholder := placeholderRef.Get()
		if placeholder.IsChecked() {
			continue
		}
		placeholder.SetChecked(true)
		if placeholder.IsReplaced() {
			continue
		}

		c.addFailureWithLocation(
			fmt.Sprintf("undefined method `%s`", lexer.Colorize(placeholder.FullName)),
			placeholder.Location(),
		)
	}
	c.methodPlaceholders = nil
}

type methodBodyCheckEntry struct {
	method         types.Ref[*types.Method]
	constantScopes []constantScope
	methodScopes   []methodScope
	node           *ast.MethodDefinitionNode
	headerMode     bool
}

func (c *Checker) registerMethodBodyCheck(method *types.Method, node *ast.MethodDefinitionNode) {
	if c.compiler != nil {
		c.compiler.RegisterMethod(node)
	}

	c.methodBodyChecks = append(c.methodBodyChecks, methodBodyCheckEntry{
		method:         types.ToRef(method),
		constantScopes: c.constantScopesCopyWithoutCache(),
		methodScopes:   c.methodScopesCopyWithoutCache(),
		node:           node,
		headerMode:     c.IsHeader(),
	})
}

var MethodCheckConcurrencyLimit = 100

func (c *Checker) checkMethodBodies() {
	concurrent.Foreach(
		MethodCheckConcurrencyLimit,
		c.methodBodyChecks,
		func(methodCheck methodBodyCheckEntry) {
			method := methodCheck.method.Get()
			node := methodCheck.node

			var mode mode
			if method.IsInit() {
				mode = initMode
			} else {
				mode = methodMode
			}
			elkName := method.NamespacedName()

			methodChecker := c.newMethodChecker(
				elkName,
				methodCheck.constantScopes,
				methodCheck.methodScopes,
				types.Ref[types.Type](method.DefinedUnder),
				method.ReturnType,
				method.ThrowType,
				mode,
				c.threadPool,
				node.Location(),
			)
			methodChecker.SetHeader(methodCheck.headerMode)

			methodChecker.checkMethodDefinition(node, method)

			// method has to be checked if it doesn't
			// use the constants that use it in their initialisation
			if len(method.UsedInConstants) > 0 {
				// use the method cache to store methods
				// that are used in constant definitions and have to be checked
				c.methodCache.Push(method.ToRef())
			}
		},
	)

	c.methodBodyChecks = nil
}

// Check whether method calls in constant definitions are valid.
// This lets the typechecker detect situations like circular references
// where a constant definition contains a call to a method
// that uses the constant that is being defined.
//
//			const FOO: Int = bar()
//	    def bar: Int
//	      FOO * 5
//	    end
func (c *Checker) checkMethodsInConstants() {
	for _, methodRef := range c.methodCache.Slice {
		method := methodRef.Get()
		c.checkMethodInConstant(method, method.UsedInConstants)
	}
}

func (c *Checker) checkMethodInConstant(method *types.Method, usedInConstants ds.Set[symbol.Symbol]) {
	for _, calledMethod := range method.CalledMethods {
		c.checkMethodInConstant(calledMethod.Get(), usedInConstants)
	}

	for usedInConstant := range usedInConstants {
		if method.UsedConstants.Contains(usedInConstant) {
			c.addFailureWithLocation(
				fmt.Sprintf(
					"method `%s` circularly refers to constant `%s` because it gets called in its initializer",
					types.InspectWithColor(method),
					lexer.Colorize(usedInConstant.String()),
				),
				method.Location(),
			)
		}
	}
}

func (c *Checker) declareMethodForGetter(node *ast.AttributeParameterNode, docComment string, pure bool) {
	name := c.identifierToName(node.Name)
	method, mod := c.declareMethod(
		c.currentMethodScope().container.Get(),
		docComment,
		false,
		false,
		false,
		false,
		false,
		false,
		pure,
		symbol.ToSymbol(name),
		nil,
		nil,
		node.TypeNode,
		nil,
		node.Location(),
	)
	method.SetAttribute(true)

	init := node.Initialiser
	var body []ast.StatementNode

	if init == nil {
		body = ast.ExpressionToStatements(
			ast.NewPublicInstanceVariableNode(node.Location(), name),
		)
	} else {
		body = ast.ExpressionToStatements(
			ast.NewAssignmentExpressionNode(
				node.Location(),
				token.New(init.Location(), token.QUESTION_QUESTION_EQUAL),
				ast.NewPublicInstanceVariableNode(node.Location(), name),
				init,
			),
		)
	}

	methodNode := ast.NewMethodDefinitionNode(
		node.Location(),
		"",
		0,
		node.Name,
		nil,
		nil,
		node.TypeNode,
		nil,
		body,
	)
	methodNode.SetType(method)
	c.registerMethodBodyCheck(
		method,
		methodNode,
	)
	if mod != nil {
		c.popConstScope()
	}
}

func (c *Checker) declareMethodForSetter(node *ast.AttributeParameterNode, docComment string) {
	setterName := c.identifierToName(node.Name) + "="

	methodScope := c.currentMethodScope()
	var paramSpan *position.Location
	if node.TypeNode != nil {
		paramSpan = node.TypeNode.Location()
	} else {
		node.Location()
	}
	params := []ast.ParameterNode{
		ast.NewMethodParameterNode(
			paramSpan,
			node.Name,
			true,
			node.TypeNode,
			nil,
			ast.NormalParameterKind,
		),
	}
	method, mod := c.declareMethod(
		methodScope.container.Get(),
		docComment,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		symbol.ToSymbol(setterName),
		nil,
		params,
		nil,
		nil,
		node.Location(),
	)
	method.SetAttribute(true)

	methodNode := ast.NewMethodDefinitionNode(
		node.Location(),
		docComment,
		0,
		ast.NewPublicIdentifierNode(node.Name.Location(), setterName),
		nil,
		params,
		nil,
		nil,
		nil,
	)
	methodNode.SetType(method)
	c.registerMethodBodyCheck(
		method,
		methodNode,
	)
	if mod != nil {
		c.popConstScope()
	}
}

func (c *Checker) addWrongArgumentCountError(got int, method *types.Method, location *position.Location) {
	c.addFailure(
		fmt.Sprintf("expected %s arguments in call to `%s`, got %d", method.ExpectedParamCountString(), types.InspectWithColor(method), got),
		location,
	)
}

func (c *Checker) addOverrideSealedMethodError(baseMethod *types.Method, loc *position.Location) {
	c.addFailureWithLocation(
		fmt.Sprintf(
			"cannot override sealed method `%s`\n  previous definition found in `%s`, with signature: `%s`",
			baseMethod.Name.String(),
			types.InspectWithColor(baseMethod.DefinedUnder.Get()),
			baseMethod.InspectSignatureWithColor(true),
		),
		loc,
	)
}

func (c *Checker) checkMethodOverride(
	overrideMethod,
	baseMethod *types.Method,
	location *position.Location,
) {
	if overrideMethod.IsRegisteredOverload() {
		return
	}

	overrideNamespace := overrideMethod.DefinedUnder.Get()
	if overrideNamespace != nil {
		overrideGeneric := types.GetDefaultNamespaceGenericInstance(overrideNamespace)
		overrideSelfMap := c.createTypeArgumentMapWithSelf(overrideGeneric)
		baseMethod = c.replaceTypeParametersInMethodCopy(baseMethod, overrideSelfMap, false)
		overrideMethod = c.replaceTypeParametersInMethodCopy(overrideMethod, overrideSelfMap, false)
	}

	if len(baseMethod.Overloads) > len(overrideMethod.Overloads) {
		errDetailsBuff := new(strings.Builder)

		fmt.Fprintf(
			errDetailsBuff,
			"missing overloads in `%s`\n  is: ",
			I(overrideNamespace),
		)

		var i int
		for overrideOverload := range overrideMethod.AllOverloads() {
			if i != 0 {
				fmt.Fprint(
					errDetailsBuff,
					"\n      ",
				)
			}
			fmt.Fprintf(
				errDetailsBuff,
				"`%s`",
				overrideOverload.InspectSignatureWithColor(false),
			)
			i++
		}

		fmt.Fprint(
			errDetailsBuff,
			"\n  should be: ",
		)

		i = 0
		for baseOverload := range baseMethod.AllOverloads() {
			if i != 0 {
				fmt.Fprint(
					errDetailsBuff,
					"\n             ",
				)
			}
			fmt.Fprintf(
				errDetailsBuff,
				"`%s`",
				baseOverload.InspectSignatureWithColor(false),
			)
			i++
		}

		c.addFailure(
			errDetailsBuff.String(),
			location,
		)
		return
	}

	c._checkMethodOverride(
		overrideMethod,
		baseMethod,
		location,
	)
}

func (c *Checker) _checkMethodOverride(
	overrideMethod,
	baseMethod *types.Method,
	location *position.Location,
) bool {
	var areIncompatible bool
	errDetailsBuff := new(strings.Builder)

	if !c.IsHeader() && baseMethod.IsSealed() {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` is sealed and cannot be overridden",
			types.InspectWithColor(baseMethod),
		)
		areIncompatible = true
	}
	if !baseMethod.IsAbstract() && overrideMethod.IsAbstract() || baseMethod.IsPure() && !overrideMethod.IsPure() {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - has a different modifier, is `%s`, should be `%s`",
			types.InspectModifier(types.ModifierSet{Abstract: overrideMethod.IsAbstract(), Sealed: overrideMethod.IsSealed(), Pure: overrideMethod.IsPure()}),
			types.InspectModifier(types.ModifierSet{Abstract: baseMethod.IsAbstract(), Sealed: baseMethod.IsSealed(), Pure: baseMethod.IsPure()}),
		)
		areIncompatible = true
	}

	if len(overrideMethod.TypeParameters) != len(baseMethod.TypeParameters) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - has a different number of type parameters, has `%d`, should have `%d`",
			len(overrideMethod.TypeParameters),
			len(baseMethod.TypeParameters),
		)
		areIncompatible = true
	} else {
		for i := range overrideMethod.TypeParameters {
			overrideTypeParam := overrideMethod.TypeParameters[i].Get()
			baseTypeParam := baseMethod.TypeParameters[i].Get()

			var isInvalid bool
			if overrideTypeParam.Name != baseTypeParam.Name || overrideTypeParam.Variance != baseTypeParam.Variance {
				isInvalid = true
			}

			switch baseTypeParam.Variance {
			case types.INVARIANT:
				if !c.isTheSameType(overrideTypeParam.UpperBound.Get(), baseTypeParam.UpperBound.Get(), nil) ||
					!c.isTheSameType(overrideTypeParam.LowerBound.Get(), baseTypeParam.LowerBound.Get(), nil) {
					isInvalid = true
				}
			case types.COVARIANT:
				if !c.isSubtype(overrideTypeParam.UpperBound.Get(), baseTypeParam.UpperBound.Get(), nil) ||
					!c.isSubtype(baseTypeParam.LowerBound.Get(), overrideTypeParam.LowerBound.Get(), nil) {
					isInvalid = true
				}
			case types.CONTRAVARIANT:
				if !c.isSubtype(baseTypeParam.UpperBound.Get(), overrideTypeParam.UpperBound.Get(), nil) ||
					!c.isSubtype(overrideTypeParam.LowerBound.Get(), baseTypeParam.LowerBound.Get(), nil) {
					isInvalid = true
				}
			}

			if isInvalid {
				fmt.Fprintf(
					errDetailsBuff,
					"\n  - has an incompatible type parameter, is `%s`, should be `%s`",
					overrideTypeParam.InspectSignature(),
					baseTypeParam.InspectSignature(),
				)
				areIncompatible = true
			}
		}
	}

	overrideReturnType := overrideMethod.ReturnType.Get()
	baseReturnType := baseMethod.ReturnType.Get()
	if !c.isSubtype(overrideReturnType, baseReturnType, nil) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - has a different return type, is `%s`, should be `%s`",
			types.InspectWithColor(overrideReturnType),
			types.InspectWithColor(baseReturnType),
		)
		areIncompatible = true
	}

	overrideThrowType := overrideMethod.ThrowType.Get()
	baseThrowType := baseMethod.ThrowType.Get()
	if !c.isSubtype(overrideThrowType, baseThrowType, nil) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - has different throw type, is `%s`, should be `%s`",
			types.InspectWithColor(overrideThrowType),
			types.InspectWithColor(baseThrowType),
		)
		areIncompatible = true
	}

	if len(baseMethod.Params) > len(overrideMethod.Params) {
		errDetailsBuff.WriteString("\n  - has less parameters")
	} else {
		for i := range len(baseMethod.Params) {
			oldParam := baseMethod.Params[i].Get()
			newParam := overrideMethod.Params[i].Get()
			if oldParam.Name != newParam.Name || oldParam.Kind != newParam.Kind || !c.isSubtype(oldParam.Type.Get(), newParam.Type.Get(), nil) {
				fmt.Fprintf(
					errDetailsBuff,
					"\n  - has an incompatible parameter, is `%s`, should be `%s`",
					types.InspectWithColor(newParam),
					types.InspectWithColor(oldParam),
				)
				areIncompatible = true
			}
		}

		for i := len(baseMethod.Params); i < len(overrideMethod.Params); i++ {
			param := overrideMethod.Params[i].Get()
			if !param.IsOptional() {
				fmt.Fprintf(
					errDetailsBuff,
					"\n  - has an additional required parameter `%s`",
					types.InspectWithColor(param),
				)
				areIncompatible = true
			}
		}
	}

	if areIncompatible {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` is not a valid override of `%s`\n  is:        `%s`\n  should be: `%s`\n%s",
				types.InspectWithColor(overrideMethod),
				types.InspectWithColor(baseMethod),
				overrideMethod.InspectSignatureWithColor(true),
				baseMethod.InspectSignatureWithColor(true),
				errDetailsBuff.String(),
			),
			location,
		)
		return false
	}

	return true
}

func (c *Checker) checkMethod(
	methodNamespace types.Namespace,
	checkedMethod *types.Method,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	body []ast.StatementNode,
	location *position.Location,
) (ast.TypeNode, ast.TypeNode) {
	prevCatchScopes := c.catchScopes
	c.catchScopes = nil
	prevHasDefer := c.hasDefer()
	c.setHasDefer(false)

	name := checkedMethod.Name
	prevMode := c.mode
	prevFlags := c.flags
	isClosure := types.IsCallable(methodNamespace)

	if methodNamespace != nil {
		parent := methodNamespace.Parent()

		if parent != nil {
			baseMethod := c.resolveMethodInNamespace(parent, name)
			if baseMethod != nil && name != symbol.S_init {
				c.checkMethodOverride(
					checkedMethod,
					baseMethod,
					location,
				)
			}
		}
	}

	if isClosure {
		c.pushNestedLocalEnv(defaultLocalEnvType)
	} else {
		c.pushIsolatedLocalEnv()
	}
	defer c.popLocalEnv()

	if !checkedMethod.IsInit() {
		c.mode = prevMode
		c.setInputPositionTypeMode()
	}
	for _, param := range paramNodes {
		switch p := param.(type) {
		case *ast.MethodParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
			pName := c.identifierToName(p.Name)
			if p.SetInstanceVariable {
				c.registerInitialisedInstanceVariable(symbol.ToSymbol(pName))
			}
			declaredType = c.TypeOf(p).(*types.Parameter).Type.Get()
			if p.TypeNode != nil {
				declaredTypeNode = p.TypeNode
				switch p.Kind {
				case ast.PositionalRestParameterKind:
					declaredType = types.NewGenericWithTypeArgs(c.StdTuple(), declaredType)
				case ast.NamedRestParameterKind:
					declaredType = types.NewGenericWithTypeArgs(c.StdRecord(), c.Std(symbol.C_Symbol), declaredType)
				}
			}
			var initNode ast.ExpressionNode
			if p.Initialiser != nil {
				initNode = c.checkExpression(p.Initialiser)
				initType := c.TypeOf(initNode)
				c.checkCanAssign(initType, declaredType, initNode.Location())
			}
			c.addLocal(pName, newLocal(types.ToRef(declaredType), true, checkedMethod.IsGenerator()))
			p.Initialiser = initNode
			p.TypeNode = declaredTypeNode
		case *ast.FormalParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
			pName := c.identifierToName(p.Name)
			declaredType = c.TypeOf(p).(*types.Parameter).Type.Get()
			if p.TypeNode != nil {
				declaredTypeNode = p.TypeNode
				switch p.Kind {
				case ast.PositionalRestParameterKind:
					declaredType = types.NewGenericWithTypeArgs(c.StdTuple(), declaredType)
				case ast.NamedRestParameterKind:
					declaredType = types.NewGenericWithTypeArgs(c.StdRecord(), c.Std(symbol.C_Symbol), declaredType)
				}
			}
			var initNode ast.ExpressionNode
			if p.Initialiser != nil {
				initNode = c.checkExpression(p.Initialiser)
				initType := c.TypeOf(initNode)
				c.checkCanAssign(initType, declaredType, initNode.Location())
			}
			c.addLocal(pName, newLocal(types.ToRef(declaredType), true, false))
			p.Initialiser = initNode
			p.TypeNode = declaredTypeNode
		default:
			panic(fmt.Sprintf("invalid parameter type: %T", param))
		}
	}

	c.mode = prevMode
	c.setOutputPositionTypeMode()

	returnType := checkedMethod.ReturnType
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
	}

	origReturnType := returnType
	if checkedMethod.IsGenerator() || checkedMethod.IsAsync() {
		returnType = origReturnType.Get().(*types.Generic).Get(0).Type
	}

	throwType := checkedMethod.ThrowType
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = typedThrowTypeNode.TypeRef()
	}
	if checkedMethod.IsGenerator() || checkedMethod.IsAsync() {
		throwType = origReturnType.Get().(*types.Generic).Get(1).Type
	}
	if !types.IsNeverRef(throwType) && throwType.IsPresent() {
		c.pushCatchScope(makeCatchScope(throwType, false))
	}

	if len(body) > 0 && checkedMethod.IsAbstract() {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` cannot have a body because it is abstract",
				name.String(),
			),
			location,
		)
	}

	if !c.IsHeader() {
		if isClosure {
			if returnType.IsZero() {
				c.setInferClosureReturnType(true)
			}
			if throwType.IsZero() {
				c.setInferClosureThrowType(true)
			}
		}
		if checkedMethod.IsGenerator() {
			c.setGenerator(true)
		} else {
			c.setGenerator(false)
		}
		if checkedMethod.IsInit() {
			c.mode = initMode
		} else if checkedMethod.IsMacro() {
			c.mode = macroMode
		} else {
			c.mode = methodMode
		}

		c.returnType = returnType
		c.throwType = throwType
		bodyReturnType, returnSpan := c.checkStatements(body, true)

		if !checkedMethod.IsAbstract() && !c.IsHeader() {
			if c.shouldInferClosureReturnType() {
				c.addToReturnType(bodyReturnType)
				checkedMethod.ReturnType = c.returnType
			} else {
				if returnSpan == nil {
					returnSpan = location
				}
				c.checkCanAssign(bodyReturnType, returnType.Get(), returnSpan)
			}

			if c.shouldInferClosureThrowType() {
				if c.throwType.IsZero() {
					checkedMethod.ThrowType = types.NeverID
				} else {
					checkedMethod.ThrowType = c.throwType
				}
			}
		}
	}

	checkedMethod.SetHasDefer(c.hasDefer())

	c.setHasDefer(prevHasDefer)
	c.returnType = types.ZERO_ID
	c.throwType = types.ZERO_ID
	c.mode = prevMode
	c.flags = prevFlags
	c.catchScopes = prevCatchScopes
	return typedReturnTypeNode, typedThrowTypeNode
}

func (c *Checker) checkSpecialMethods(name symbol.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, location *position.Location) {
	if symbol.IsEqualityOperator(name) {
		c.checkEqualityOperator(name, checkedMethod, paramNodes, location)
		return
	}

	if symbol.IsRelationalOperator(name) {
		c.checkRelationalOperator(name, checkedMethod, paramNodes, location)
		return
	}

	if symbol.RequiresOneParameter(name) {
		c.checkFixedParameterCountMethod(name, checkedMethod, paramNodes, 1, location)
		return
	}

	if symbol.RequiresNoParameters(name) {
		c.checkFixedParameterCountMethod(name, checkedMethod, paramNodes, 0, location)
		return
	}
}

func (c *Checker) checkEqualityOperator(name symbol.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, location *position.Location) {
	params := checkedMethod.Params

	if !c.isTheSameType(checkedMethod.ReturnType.Get(), types.Bool{}, nil) {
		c.addFailure(
			fmt.Sprintf(
				"equality operator `%s` must return `%s`",
				lexer.Colorize(name.String()),
				lexer.Colorize("bool"),
			),
			location,
		)
	}

	if len(params) != 1 {
		c.addFailure(
			fmt.Sprintf(
				"equality operator `%s` must accept a single parameter, got %d",
				lexer.Colorize(name.String()),
				len(params),
			),
			location,
		)
		return
	}

	param := params[0].Get()
	var paramSpan *position.Location
	if paramNodes != nil {
		paramSpan = paramNodes[0].Location()
	} else {
		paramSpan = location
	}
	if !types.IsAnyRef(param.Type) {
		c.addFailure(
			fmt.Sprintf(
				"parameter `%s` of equality operator `%s` must be of type `%s`",
				lexer.Colorize(param.Name.String()),
				lexer.Colorize(name.String()),
				lexer.Colorize("any"),
			),
			paramSpan,
		)
	}

	switch param.Kind {
	case types.PositionalRestParameterKind, types.NamedRestParameterKind:
		c.addFailure(
			fmt.Sprintf(
				"equality operator `%s` cannot define rest parameter `%s`",
				lexer.Colorize(name.String()),
				types.InspectWithColor(param),
			),
			paramSpan,
		)
	}
}

func (c *Checker) checkRelationalOperator(name symbol.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, location *position.Location) {
	params := checkedMethod.Params

	if !c.isTheSameType(checkedMethod.ReturnType.Get(), types.Bool{}, nil) {
		c.addFailure(
			fmt.Sprintf(
				"relational operator `%s` must return `%s`",
				lexer.Colorize(name.String()),
				lexer.Colorize("bool"),
			),
			location,
		)
	}

	if len(params) != 1 {
		c.addFailure(
			fmt.Sprintf(
				"relational operator `%s` must accept a single parameter, got %d",
				lexer.Colorize(name.String()),
				len(params),
			),
			location,
		)
		return
	}

	param := checkedMethod.Params[0].Get()
	var paramSpan *position.Location
	if paramNodes != nil {
		paramSpan = paramNodes[0].Location()
	} else {
		paramSpan = location
	}
	self := c.selfType.Get()
	if !checkedMethod.IsAbstract() && !c.isSubtype(self, param.Type.Get(), nil) {
		c.addFailure(
			fmt.Sprintf(
				"parameter `%s` of relational operator `%s` must accept `%s`",
				lexer.Colorize(param.Name.String()),
				lexer.Colorize(name.String()),
				types.InspectWithColor(self),
			),
			paramSpan,
		)
	}

	switch param.Kind {
	case types.PositionalRestParameterKind, types.NamedRestParameterKind:
		c.addFailure(
			fmt.Sprintf(
				"relational operator `%s` cannot define rest parameter `%s`",
				lexer.Colorize(name.String()),
				types.InspectWithColor(param),
			),
			paramSpan,
		)
	}
}

func (c *Checker) checkFixedParameterCountMethod(name symbol.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, desiredParamCount int, location *position.Location) {
	params := checkedMethod.Params

	if types.IsVoid(checkedMethod.ReturnType.Get()) {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` cannot be void",
				lexer.Colorize(name.String()),
			),
			location,
		)
	}

	if len(params) != desiredParamCount {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` must define exactly %d parameters, got %d",
				lexer.Colorize(name.String()),
				desiredParamCount,
				len(params),
			),
			location,
		)
		return
	}

	for i, paramRef := range params {
		var paramSpan *position.Location
		if paramNodes != nil {
			paramSpan = paramNodes[i].Location()
		} else {
			paramSpan = location
		}
		param := paramRef.Get()

		switch param.Kind {
		case types.PositionalRestParameterKind, types.NamedRestParameterKind:
			c.addFailure(
				fmt.Sprintf(
					"method `%s` cannot define rest parameter `%s`",
					lexer.Colorize(name.String()),
					types.InspectWithColor(param),
				),
				paramSpan,
			)
		}
	}
}

func (c *Checker) addToReturnType(typ types.Type) {
	if c.returnType.IsZero() {
		c.returnType = types.ToRef(typ)
		return
	}

	c.returnType = types.ToRef(c.NewNormalisedUnion(c.returnType, types.ToRef(typ)))
}

func (c *Checker) addToThrowType(typ types.Type) {
	if c.throwType.IsZero() {
		c.throwType = types.ToRef(typ)
		return
	}

	c.throwType = types.ToRef(c.NewNormalisedUnion(c.throwType, types.ToRef(typ)))
}

type inferArg struct {
	typedArg  ast.ExpressionNode
	isClosure bool
	param     *types.Parameter
	retry     bool
}

func (c *Checker) checkMethodArgumentsAndInferTypeArguments(
	method *types.Method,
	positionalArguments []ast.ExpressionNode,
	namedArguments []ast.NamedArgumentNode,
	typeParams []types.Ref[*types.TypeParameter],
	location *position.Location,
) (
	_method *types.Method,
	_posArgs []ast.ExpressionNode,
	typeArgs types.TypeArgumentMap,
) {
	if len(method.Overloads) == 0 {
		posArgs, typeArgs := c._checkMethodArgumentsAndInferTypeArguments(
			method,
			positionalArguments,
			namedArguments,
			typeParams,
			location,
		)
		return method, posArgs, typeArgs
	}

	prevDiagnostics := c.Errors
	tempDiagnostics := diagnostic.NewSyncDiagnosticList()
	c.Errors = tempDiagnostics
	for overload := range method.ReversedOverloads() {
		tempDiagnostics.Clear()

		posArgs := ds.MapSlice(positionalArguments, func(arg ast.ExpressionNode) ast.ExpressionNode {
			return ast.DeepCopy(arg).(ast.ExpressionNode)
		})
		namedArgs := ds.MapSlice(namedArguments, func(arg ast.NamedArgumentNode) ast.NamedArgumentNode {
			return ast.DeepCopy(arg).(ast.NamedArgumentNode)
		})

		posArgs, typeArgs := c._checkMethodArgumentsAndInferTypeArguments(
			overload,
			posArgs,
			namedArgs,
			typeParams,
			location,
		)

		if !tempDiagnostics.IsFailure() {
			c.Errors = prevDiagnostics
			return overload, posArgs, typeArgs
		}
	}

	c.Errors = prevDiagnostics

	errDetailsBuff := new(strings.Builder)

	fmt.Fprintf(
		errDetailsBuff,
		"no overload of `%s` matches the given arguments\n  signature: `%s`",
		method.Name.String(),
		method.InspectSignatureWithColor(false),
	)

	for _, overloadRef := range method.Overloads {
		overload := overloadRef.Get()
		fmt.Fprintf(
			errDetailsBuff,
			"\n             `%s`",
			overload.InspectSignatureWithColor(false),
		)
	}

	c.addFailure(
		errDetailsBuff.String(),
		location,
	)
	return nil, nil, nil
}

func (c *Checker) _checkMethodArgumentsAndInferTypeArguments(
	method *types.Method,
	positionalArguments []ast.ExpressionNode,
	namedArguments []ast.NamedArgumentNode,
	typeParams []types.Ref[*types.TypeParameter],
	location *position.Location,
) (
	_posArgs []ast.ExpressionNode,
	typeArgs types.TypeArgumentMap,
) {
	var typeArgMap types.TypeArgumentMap
	if typeParams != nil {
		prevMode := c.mode
		c.mode = inferTypeArgumentMode
		defer c.setMode(prevMode)
		typeArgMap = make(types.TypeArgumentMap)
	}

	reqParamCount := method.RequiredParamCount()
	requiredPosParamCount := len(method.Params) - method.OptionalParamCount
	if method.PostParamCount != -1 {
		requiredPosParamCount -= method.PostParamCount + 1
	}
	if method.HasNamedRestParam() {
		requiredPosParamCount--
	}
	argCount := len(positionalArguments) + len(namedArguments)
	positionalRestParamIndex := method.PositionalRestParamIndex()
	var typedPositionalArguments []ast.ExpressionNode

	// push `undefined` for every missing optional positional argument
	// before the rest parameter
	for range positionalRestParamIndex - len(positionalArguments) {
		typedPositionalArguments = append(
			typedPositionalArguments,
			ast.NewUndefinedLiteralNode(location),
		)
	}

	var inferArgs []inferArg
	var currentParamIndex int
	// check all positional arguments before the rest parameter
	for ; currentParamIndex < len(positionalArguments); currentParamIndex++ {
		posArg := positionalArguments[currentParamIndex]
		if currentParamIndex == positionalRestParamIndex {
			break
		}
		if currentParamIndex >= len(method.Params) {
			c.addWrongArgumentCountError(
				len(positionalArguments)+len(namedArguments),
				method,
				location,
			)
			break
		}
		param := method.Params[currentParamIndex].Get()

		if _, ok := posArg.(*ast.ClosureLiteralNode); ok {
			inferArgs = append(inferArgs, inferArg{
				typedArg:  posArg,
				isClosure: true,
				param:     param,
				retry:     true,
			})
			continue
		}
		paramType := param.Type.Get()
		typedPosArg := c.checkExpressionWithType(posArg, paramType)
		posArgType := c.TypeOf(typedPosArg)

		inferredParamType := c.inferTypeArguments(posArgType, paramType, typeArgMap, typedPosArg.Location())

		var retry bool
		switch inferredParamType {
		case nil:
			paramCopy := param.Copy()
			paramCopy.Type = types.UntypedID
			method.Params[currentParamIndex] = paramCopy.ToRef()
		case paramType:
			retry = true
		default:
			paramCopy := param.Copy()
			paramCopy.Type = types.ToRef(inferredParamType)
			method.Params[currentParamIndex] = paramCopy.ToRef()
		}
		inferArgs = append(inferArgs, inferArg{
			typedArg: typedPosArg,
			param:    param,
			retry:    retry,
		})
	}

	for _, inferArg := range inferArgs {
		typedPosArg := inferArg.typedArg
		posArgType := c.TypeOf(typedPosArg)
		param := inferArg.param
		paramType := param.Type.Get()

		if inferArg.retry {
			if inferArg.isClosure {
				typedPosArg = c.checkExpressionWithTypeArgs(typedPosArg, paramType, typeArgMap)
				posArgType = c.TypeOf(typedPosArg)
			}
			inferredParamType := c.inferTypeArguments(posArgType, paramType, typeArgMap, typedPosArg.Location())
			if inferredParamType == nil {
				paramCopy := param.Copy()
				paramCopy.Type = types.UntypedID
				method.Params[currentParamIndex] = paramCopy.ToRef()
			} else if inferredParamType != paramType {
				paramCopy := param.Copy()
				paramCopy.Type = types.ToRef(inferredParamType)
				method.Params[currentParamIndex] = paramCopy.ToRef()
			}
		}

		typedPositionalArguments = append(typedPositionalArguments, typedPosArg)

		if !c.isSubtype(posArgType, paramType, typedPosArg.Location()) {
			c.addFailure(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
					types.InspectWithColor(paramType),
					param.Name.String(),
					types.InspectWithColor(method),
					types.InspectWithColor(posArgType),
				),
				typedPosArg.Location(),
			)
		}
	}

	if method.HasPositionalRestParam() {
		if len(positionalArguments) < requiredPosParamCount {
			c.addFailure(
				fmt.Sprintf(
					"expected %d... positional arguments in call to `%s`, got %d",
					requiredPosParamCount,
					types.InspectWithColor(method),
					len(positionalArguments),
				),
				location,
			)
			return nil, nil
		}
		restPositionalArguments := ast.NewArrayTupleLiteralNode(
			location,
			nil,
		)
		posRestParam := method.Params[positionalRestParamIndex].Get()
		posRestParamType := posRestParam.Type.Get()

		currentArgIndex := currentParamIndex
		// check rest arguments
		for ; currentArgIndex < min(argCount-method.PostParamCount, len(positionalArguments)); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]

			typedPosArg := c.checkRestArgument(posArg, posRestParamType)
			posArgType := c.TypeOf(typedPosArg)
			inferredParamType := c.inferTypeArguments(posArgType, posRestParamType, typeArgMap, typedPosArg.Location())
			if inferredParamType == nil {
				posRestParamCopy := posRestParam.Copy()
				posRestParamCopy.Type = types.UntypedID
				method.Params[positionalRestParamIndex] = posRestParamCopy.ToRef()
			} else if inferredParamType != posRestParamType {
				posRestParamCopy := posRestParam.Copy()
				posRestParamCopy.Type = types.ToRef(inferredParamType)
				method.Params[positionalRestParamIndex] = posRestParamCopy.ToRef()
			}
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, typedPosArg)
			if !c.isSubtype(posArgType, posRestParamType, posArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(posRestParamType),
						posRestParam.Name.String(),
						types.InspectWithColor(method),
						types.InspectWithColor(posArgType),
					),
					posArg.Location(),
				)
			}
		}
		typedPositionalArguments = append(typedPositionalArguments, restPositionalArguments)
		if len(restPositionalArguments.Elements) == 1 {
			element := restPositionalArguments.Elements[0]
			if forIn, ok := element.(*ast.ModifierForInNode); ok {
				inType := c.TypeOf(forIn.InExpression)
				if c.IsSubtype(inType, c.Std(symbol.C_Tuple)) {
					typedPositionalArguments[len(typedPositionalArguments)-1] = forIn.InExpression
				}
			}
		}

		currentParamIndex = positionalRestParamIndex
		// check post arguments
		for ; currentArgIndex < len(positionalArguments); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			currentParamIndex++
			param := method.Params[currentParamIndex].Get()
			paramType := param.Type.Get()

			typedPosArg := c.checkExpressionWithType(posArg, paramType)
			posArgType := c.TypeOf(typedPosArg)
			inferredParamType := c.inferTypeArguments(posArgType, paramType, typeArgMap, typedPosArg.Location())
			if inferredParamType == nil {
				paramCopy := posRestParam.Copy()
				paramCopy.Type = types.UntypedID
				method.Params[currentParamIndex] = paramCopy.ToRef()
			} else if inferredParamType != paramType {
				paramCopy := posRestParam.Copy()
				paramCopy.Type = types.ToRef(inferredParamType)
				method.Params[currentParamIndex] = paramCopy.ToRef()
			}
			typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
			if !c.isSubtype(posArgType, paramType, posArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(paramType),
						param.Name.String(),
						types.InspectWithColor(method),
						types.InspectWithColor(posArgType),
					),
					posArg.Location(),
				)
			}
		}
		currentParamIndex++

		if method.PostParamCount > 0 {
			reqParamCount++
		}
	}

	firstNamedParamIndex := currentParamIndex
	definedNamedArgumentsSlice := make([]bool, len(namedArguments))

	for i := range method.Params {
		param := method.Params[i].Get()
		paramType := param.Type.Get()
		switch param.Kind {
		case types.PositionalRestParameterKind, types.NamedRestParameterKind:
			continue
		}
		paramName := param.Name.String()
		var found bool

		for namedArgIndex, namedArgI := range namedArguments {
			var namedArg *ast.NamedCallArgumentNode
			switch n := namedArgI.(type) {
			case *ast.NamedCallArgumentNode:
				namedArg = n
			case *ast.DoubleSplatExpressionNode:
				continue
			default:
				panic(fmt.Sprintf("invalid named argument node: %T", namedArgI))
			}

			if c.identifierToName(namedArg.Name) != paramName {
				continue
			}
			if found || i < firstNamedParamIndex {
				c.addFailure(
					fmt.Sprintf(
						"duplicated argument `%s` in call to `%s`",
						paramName,
						types.InspectWithColor(method),
					),
					namedArg.Location(),
				)
			}
			found = true
			definedNamedArgumentsSlice[namedArgIndex] = true

			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, paramType)
			namedArgType := c.TypeOf(typedNamedArgValue)
			inferredParamType := c.inferTypeArguments(namedArgType, paramType, typeArgMap, typedNamedArgValue.Location())
			if inferredParamType == nil {
				paramCopy := param.Copy()
				paramCopy.Type = types.UntypedID
				method.Params[i] = paramCopy.ToRef()
			} else if inferredParamType != paramType {
				paramCopy := param.Copy()
				paramCopy.Type = types.ToRef(inferredParamType)
				method.Params[i] = paramCopy.ToRef()
			}
			typedPositionalArguments = append(typedPositionalArguments, typedNamedArgValue)
			if !c.isSubtype(namedArgType, paramType, namedArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(paramType),
						param.Name.String(),
						types.InspectWithColor(method),
						types.InspectWithColor(namedArgType),
					),
					namedArg.Location(),
				)
			}
		}

		if i < firstNamedParamIndex {
			continue
		}
		if found {
			continue
		}

		if i < reqParamCount {
			// the parameter is required
			// but is not present in the call
			c.addFailure(
				fmt.Sprintf(
					"argument `%s` is missing in call to `%s`",
					paramName,
					types.InspectWithColor(method),
				),
				location,
			)
		} else {
			// the parameter is missing and is optional
			// we push undefined as its value
			typedPositionalArguments = append(
				typedPositionalArguments,
				ast.NewUndefinedLiteralNode(location),
			)
		}
	}

	if method.HasNamedRestParam() {
		namedRestArgs := ast.NewHashRecordLiteralNode(
			location,
			nil,
		)
		namedRestParam := method.Params[len(method.Params)-1].Get()
		namedRestParamType := namedRestParam.Type.Get()
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			switch namedArg := namedArgI.(type) {
			case *ast.NamedCallArgumentNode:
				typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, namedRestParamType)
				posArgType := c.TypeOf(typedNamedArgValue)
				inferredParamType := c.inferTypeArguments(posArgType, namedRestParamType, typeArgMap, typedNamedArgValue.Location())
				if inferredParamType == nil {
					namedRestParamCopy := namedRestParam.Copy()
					namedRestParamCopy.Type = types.UntypedID
					method.Params[len(method.Params)-1] = namedRestParamCopy.ToRef()
				} else if inferredParamType != namedRestParamType {
					namedRestParamCopy := namedRestParam.Copy()
					namedRestParamCopy.Type = types.ToRef(inferredParamType)
					method.Params[len(method.Params)-1] = namedRestParamCopy.ToRef()
				}
				namedRestArgs.Elements = append(
					namedRestArgs.Elements,
					ast.NewSymbolKeyValueExpressionNode(
						namedArg.Location(),
						namedArg.Name,
						typedNamedArgValue,
					),
				)
				namedArgType := c.TypeOf(typedNamedArgValue)
				c.checkNamedRestArgumentType(
					method.Name.String(),
					namedArgType,
					namedRestParam,
					namedArg.Location(),
				)
			case *ast.DoubleSplatExpressionNode:
				result := c.checkDoubleSplatArgument(method.Name.String(), namedArg, namedRestParam)
				namedRestArgs.Elements = append(
					namedRestArgs.Elements,
					result,
				)
			default:
				panic(fmt.Sprintf("invalid named argument node: %T", namedArgI))
			}
		}

		typedPositionalArguments = append(typedPositionalArguments, namedRestArgs)
		if len(namedRestArgs.Elements) == 1 {
			element := namedRestArgs.Elements[0]
			if forIn, ok := element.(*ast.ModifierForInNode); ok {
				inType := c.TypeOf(forIn.InExpression)
				if c.IsSubtype(inType, c.Std(symbol.C_Record)) {
					typedPositionalArguments[len(typedPositionalArguments)-1] = forIn.InExpression
				}
			}
		}
	} else {
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			switch namedArg := namedArgI.(type) {
			case *ast.NamedCallArgumentNode:
				c.addFailure(
					fmt.Sprintf(
						"nonexistent parameter `%s` given in call to `%s`",
						namedArg.Name.String(),
						types.InspectWithColor(method),
					),
					namedArg.Location(),
				)
			case *ast.DoubleSplatExpressionNode:
				c.addFailure(
					fmt.Sprintf(
						"double splat arguments cannot be present in calls to methods without a named rest parameter eg. `%s`",
						lexer.Colorize("**foo: Int"),
					),
					namedArg.Location(),
				)
			}
		}
	}

	if typeArgMap != nil && len(typeArgMap) != len(typeParams) {
		for _, typeParamRef := range typeParams {
			typeParam := typeParamRef.Get()
			typeArg := typeArgMap[typeParam.Name]
			if typeArg != nil {
				continue
			}

			typeParamLowerBound := typeParam.LowerBound.Get()
			typeParamUpperBound := typeParam.UpperBound.Get()
			var inferredType types.Type
			if !types.IsNeverRef(typeParam.LowerBound) && !c.containsTypeParameters(typeParamLowerBound) {
				inferredType = typeParamLowerBound
			} else if !c.containsTypeParameters(typeParamUpperBound) {
				inferredType = typeParamUpperBound
			} else {
				inferredType = types.Untyped{}
				c.addFailure(
					fmt.Sprintf(
						"cannot infer type argument for `%s` in call to `%s`",
						types.InspectWithColor(typeParam),
						types.InspectWithColor(method),
					),
					location,
				)
			}

			typeArgMap[typeParam.Name] = types.NewTypeArgument(
				inferredType,
				typeParam.Variance,
			)
		}
	}

	return typedPositionalArguments, typeArgMap
}

func (c *Checker) checkDoubleSplatArgument(methodName string, node *ast.DoubleSplatExpressionNode, namedRestParam *types.Parameter) ast.ExpressionNode {
	result, keyType, valueType := c.checkRecordDoubleSplatExpression(node)
	if !c.isSubtype(keyType, c.Std(symbol.C_Symbol), node.Location()) {
		c.addFailure(
			fmt.Sprintf(
				"expected type `%s` for double splat argument keys, got `%s`",
				lexer.Colorize("Std::Symbol"),
				types.InspectWithColor(keyType),
			),
			node.Location(),
		)
	}

	c.checkNamedRestArgumentType(
		methodName,
		valueType,
		namedRestParam,
		node.Location(),
	)

	return result
}

func (c *Checker) checkNamedRestArgumentType(methodName string, argType types.Type, param *types.Parameter, location *position.Location) {
	paramType := param.Type.Get()
	if c.isSubtype(argType, paramType, location) {
		return
	}

	c.addFailure(
		fmt.Sprintf(
			"expected type `%s` for named rest parameter `**%s` in call to `%s`, got type `%s`",
			types.InspectWithColor(paramType),
			param.Name.String(),
			lexer.Colorize(methodName),
			types.InspectWithColor(argType),
		),
		location,
	)
}

func (c *Checker) checkNonGenericMethodArguments(
	method *types.Method,
	positionalArguments []ast.ExpressionNode,
	namedArguments []ast.NamedArgumentNode,
	location *position.Location,
) (*types.Method, []ast.ExpressionNode) {
	method, posArgs, _ := c.checkMethodArgumentsAndInferTypeArguments(method, positionalArguments, namedArguments, nil, location)
	return method, posArgs
}

func (c *Checker) checkRestArgument(node ast.ExpressionNode, typ types.Type) ast.ExpressionNode {
	switch n := node.(type) {
	case *ast.SplatExpressionNode:
		return c.checkCollectionSplatExpression(n)
	default:
		return c.checkExpressionWithType(node, typ)
	}
}

func (c *Checker) checkMethodArguments(
	method *types.Method,
	typeArgumentNodes []ast.TypeNode,
	positionalArgumentNodes []ast.ExpressionNode,
	namedArgumentNodes []ast.NamedArgumentNode,
	location *position.Location,
) (_method *types.Method, typedPositionalArguments []ast.ExpressionNode) {
	if len(typeArgumentNodes) > 0 {
		typeArgs, ok := c.checkTypeArguments(
			method,
			typeArgumentNodes,
			method.TypeParameters,
			location,
		)
		if !ok {
			c.checkExpressions(positionalArgumentNodes)
			c.checkNamedArguments(namedArgumentNodes)
			return nil, nil
		}

		method = c.replaceTypeParametersInMethodCopy(method, typeArgs.ArgumentMap, true)
		return c.checkNonGenericMethodArguments(
			method,
			positionalArgumentNodes,
			namedArgumentNodes,
			location,
		)
	}

	if len(method.TypeParameters) > 0 {
		var typeArgMap types.TypeArgumentMap
		method = method.Copy()
		var chosenMethod *types.Method
		chosenMethod, typedPositionalArguments, typeArgMap = c.checkMethodArgumentsAndInferTypeArguments(
			method,
			positionalArgumentNodes,
			namedArgumentNodes,
			method.TypeParameters,
			location,
		)
		if len(typeArgMap) != len(chosenMethod.TypeParameters) {
			return nil, nil
		}
		chosenMethod.ReturnType = types.ToRef(c.replaceTypeParameters(chosenMethod.ReturnType.Get(), typeArgMap, true))
		chosenMethod.ThrowType = types.ToRef(c.replaceTypeParameters(chosenMethod.ThrowType.Get(), typeArgMap, true))
		return chosenMethod, typedPositionalArguments
	}

	return c.checkNonGenericMethodArguments(
		method,
		positionalArgumentNodes,
		namedArgumentNodes,
		location,
	)
}

func (c *Checker) checkSimpleMethodCall(
	receiver ast.ExpressionNode,
	op token.Type,
	methodName symbol.Symbol,
	typeArgumentNodes []ast.TypeNode,
	positionalArgumentNodes []ast.ExpressionNode,
	namedArgumentNodes []ast.NamedArgumentNode,
	location *position.Location,
) (
	_methodName symbol.Symbol,
	_receiver ast.ExpressionNode,
	_positionalArguments []ast.ExpressionNode,
	typ types.Type,
) {
	receiver = c.checkExpression(receiver)
	receiverType := c.TypeOf(receiver)

	// Allow arbitrary method calls on `never` and `untyped`.
	// Typecheck the arguments.
	if types.IsNever(receiverType) || types.IsUntyped(receiverType) {
		var typedPositionalArguments []ast.ExpressionNode

		for _, argument := range positionalArgumentNodes {
			typedPositionalArguments = append(typedPositionalArguments, c.checkExpression(argument))
		}
		for _, argument := range namedArgumentNodes {
			arg, ok := argument.(*ast.NamedCallArgumentNode)
			if !ok {
				continue
			}
			typedPositionalArguments = append(typedPositionalArguments, c.checkExpression(arg.Value))
		}

		return methodName, receiver, typedPositionalArguments, receiverType
	}

	var method *types.Method
	switch op {
	case token.DOT, token.DOT_DOT:
		method = c.GetMethod(receiverType, methodName, location)
	case token.QUESTION_DOT, token.QUESTION_DOT_DOT:
		nonNilableReceiverType := c.ToNonNilable(receiverType)
		method = c.GetMethod(nonNilableReceiverType, methodName, location)
	default:
		panic(fmt.Sprintf("invalid call operator: %#v", op))
	}
	if method == nil {
		c.checkExpressions(positionalArgumentNodes)
		c.checkNamedArguments(namedArgumentNodes)
		return methodName, receiver, positionalArgumentNodes, types.Untyped{}
	}

	c.addToMethodCache(method)

	method, typedPositionalArguments := c.checkMethodArguments(method, typeArgumentNodes, positionalArgumentNodes, namedArgumentNodes, location)
	if method == nil {
		return methodName, receiver, positionalArgumentNodes, types.Untyped{}
	}

	var returnType types.Type
	switch op {
	case token.DOT:
		returnType = method.ReturnType.Get()
	case token.QUESTION_DOT:
		if !c.IsNilable(receiverType) {
			c.addFailure(
				fmt.Sprintf("cannot make a nil-safe call on type `%s` which is not nilable", types.InspectWithColor(receiverType)),
				location,
			)
			returnType = method.ReturnType.Get()
		} else {
			returnType = c.ToNilable(method.ReturnType.Get())
		}
	case token.DOT_DOT:
		returnType = receiverType
	case token.QUESTION_DOT_DOT:
		if !c.IsNilable(receiverType) {
			c.addFailure(
				fmt.Sprintf("cannot make a nil-safe call on type `%s` which is not nilable", types.InspectWithColor(receiverType)),
				location,
			)
		}
		returnType = receiverType
	}

	if !method.IsPure() {
		c.addImpureErrorIfInPureContext(location)
	}
	c.checkCalledMethodThrowType(method, location)

	return method.Name, receiver, typedPositionalArguments, returnType
}

func (c *Checker) checkBinaryOpMethodCall(
	node *ast.BinaryExpressionNode,
	methodName symbol.Symbol,
) ast.ExpressionNode {
	chosenMethodName, receiver, args, returnType := c.checkSimpleMethodCall(
		node.Left,
		token.DOT,
		methodName,
		nil,
		[]ast.ExpressionNode{node.Right},
		nil,
		node.Location(),
	)
	if chosenMethodName != methodName {
		newNode := ast.NewMethodCallNode(
			node.Location(),
			receiver,
			token.New(node.Location(), token.DOT),
			ast.NewPublicIdentifierNode(node.Op.Location(), chosenMethodName.String()),
			args,
			nil,
		)
		newNode.SetType(returnType)
		return newNode
	}

	node.SetType(returnType)
	return node
}

func (c *Checker) checkMethodDefinition(node *ast.MethodDefinitionNode, method *types.Method) {
	c.method = method
	returnType, throwType := c.checkMethod(
		c.currentMethodScope().container.Get(),
		method,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Body,
		node.Location(),
	)

	node.ReturnType = returnType
	node.ThrowType = throwType

	c.method = nil

	method.CalledMethods = c.methodCache.Slice
	c.methodCache.Slice = nil

	if c.shouldCompile() && method.IsCompilable() {
		method.Body = c.compiler.CompileMethodBody(node, value.ToSymbol(method.NamespacedName())).Method()
	}
}

func (c *Checker) declareMethod(
	methodNamespace types.Namespace,
	docComment string,
	abstract bool,
	sealed bool,
	inferReturnType bool,
	generator bool,
	async bool,
	overload bool,
	pure bool,
	name symbol.Symbol,
	typeParamNodes []ast.TypeParameterNode,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	location *position.Location,
) (*types.Method, *types.TypeParamNamespace) {
	return c.declareMethodWithBase(
		nil,
		nil,
		methodNamespace,
		docComment,
		abstract,
		sealed,
		inferReturnType,
		generator,
		async,
		overload,
		pure,
		name,
		typeParamNodes,
		paramNodes,
		returnTypeNode,
		throwTypeNode,
		location,
	)
}

func (c *Checker) declareMethodWithBase(
	baseMethod *types.Method,
	typeArgMap types.TypeArgumentMap,
	methodNamespace types.Namespace,
	docComment string,
	abstract bool,
	sealed bool,
	inferReturnType bool,
	generator bool,
	async bool,
	overload bool,
	pure bool,
	name symbol.Symbol,
	typeParamNodes []ast.TypeParameterNode,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	location *position.Location,
) (*types.Method, *types.TypeParamNamespace) {
	prevMode := c.mode
	if c.mode == interfaceMode {
		abstract = true
	}
	if abstract && overload {
		c.addFailure(
			fmt.Sprintf(
				"abstract method `%s` cannot be overloaded",
				name.String(),
			),
			location,
		)
	}
	oldMethod := methodNamespace.Method(name)
	if !overload && oldMethod != nil {
		if sealed && !oldMethod.IsSealed() || !pure && oldMethod.IsPure() {
			c.addFailure(
				fmt.Sprintf(
					"cannot redeclare method `%s` with a different modifier, is `%s`, should be `%s`",
					name.String(),
					types.InspectModifier(types.ModifierSet{Abstract: abstract, Sealed: sealed, Pure: pure}),
					types.InspectModifier(types.ModifierSet{Abstract: oldMethod.IsAbstract(), Sealed: oldMethod.IsSealed(), Pure: oldMethod.IsPure()}),
				),
				location,
			)
		}
	}

	var isImmutable bool
	switch namespace := methodNamespace.(type) {
	case *types.Interface:
	case *types.Class:
		isImmutable = namespace.IsImmutable()
		if abstract && !namespace.IsAbstract() {
			c.addFailure(
				fmt.Sprintf(
					"cannot declare abstract method `%s` in non-abstract class `%s`",
					name.String(),
					types.InspectWithColor(methodNamespace),
				),
				location,
			)
		}
	case *types.Mixin:
		if abstract && !namespace.IsAbstract() {
			c.addFailure(
				fmt.Sprintf(
					"cannot declare abstract method `%s` in non-abstract mixin `%s`",
					name.String(),
					types.InspectWithColor(methodNamespace),
				),
				location,
			)
		}
	default:
		if abstract {
			c.addFailure(
				fmt.Sprintf(
					"cannot declare abstract method `%s` in this context",
					name.String(),
				),
				location,
			)
		}
	}

	if name == symbol.S_init {
		c.mode = initMode
	} else {
		c.mode = methodMode
	}

	var typeParams []types.Ref[*types.TypeParameter]
	var typeParamMod *types.TypeParamNamespace
	if len(typeParamNodes) > 0 {
		typeParams = make([]types.Ref[*types.TypeParameter], 0, len(typeParamNodes))
		typeParamMod = types.NewTypeParamNamespace(fmt.Sprintf("Type Parameter Container of %s", name), true)
		c.pushConstScope(makeConstantScope(types.CastRef[types.Namespace](typeParamMod)))
		for _, typeParamNode := range typeParamNodes {
			node, ok := typeParamNode.(*ast.VariantTypeParameterNode)
			if !ok {
				continue
			}

			t := c.checkTypeParameterNode(node, typeParamMod, false)
			typeParams = append(typeParams, t.ToRef())
			typeParamNode.SetType(t)
			typeParamMod.DefineSubtype(t.Name, t)
			typeParamMod.DefineConstant(t.Name, types.NoValue{})
		}
	}

	if name != symbol.S_init {
		c.mode = prevMode
		c.setInputPositionTypeMode()
	}
	var params []types.Ref[*types.Parameter]
	for i, paramNode := range paramNodes {
		switch p := paramNode.(type) {
		case *ast.FormalParameterNode:
			pName := c.identifierToName(p.Name)
			var declaredType types.Type
			if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.TypeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Get().Type.Get()
				declaredType = c.inferTypeArgumentsWithFlags(declaredType, declaredType, typeArgMap, nil, bitfield.BitField8FromBitFlag(inferTypeArgumentsInferFromDefaults))
			} else {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", pName),
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
			name := symbol.ToSymbol(pName)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, types.ToRef(paramType))
		case *ast.MethodParameterNode:
			pName := c.identifierToName(p.Name)
			var declaredType types.Type
			if p.SetInstanceVariable {
				currentIvar, _ := c.getInstanceVariableIn(symbol.ToSymbol(pName), methodNamespace)
				if p.TypeNode == nil {
					if currentIvar == nil {
						c.addFailure(
							fmt.Sprintf(
								"cannot infer the type of instance variable `%s`",
								pName,
							),
							p.Location(),
						)
					} else {
						declaredType = currentIvar.Type.Get()
					}
				} else {
					p.TypeNode = c.checkTypeNode(p.TypeNode)
					declaredType = c.TypeOf(p.TypeNode)
					if currentIvar != nil {
						c.checkCanAssignInstanceVariable(pName, declaredType, currentIvar, p.TypeNode.Location())
					} else {
						c.declareInstanceVariable(symbol.ToSymbol(pName), declaredType, "", isImmutable, p.Location())
					}
				}
			} else if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.TypeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Get().Type.Get()
			} else {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", pName),
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
			name := symbol.ToSymbol(pName)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType.ToRef())
		case *ast.SignatureParameterNode:
			pName := c.identifierToName(p.Name)
			var declaredType types.Type
			if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.TypeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Get().Type.Get()
			} else {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", pName),
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
			if p.Optional {
				kind = types.DefaultValueParameterKind
			}
			name := symbol.ToSymbol(pName)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType.ToRef())
		default:
			c.addFailure(
				fmt.Sprintf("invalid param type %T", paramNode),
				paramNode.Location(),
			)
		}
	}
	if async {
		paramType := types.NewParameter(
			symbol.ToSymbol("_pool"),
			c.Std(symbol.C_ThreadPool),
			types.DefaultValueParameterKind,
			false,
		)
		params = append(params, paramType.ToRef())
	}

	c.mode = prevMode
	c.setOutputPositionTypeMode()

	var returnType types.Type
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.TypeOf(typedReturnTypeNode)
	} else if inferReturnType {
	} else if baseMethod != nil && baseMethod.ReturnType.IsPresent() {
		returnType = baseMethod.ReturnType.Get()
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.TypeOf(typedThrowTypeNode)
	} else if inferReturnType {
	} else if baseMethod != nil && baseMethod.ThrowType.IsPresent() {
		throwType = baseMethod.ThrowType.Get()
	} else {
		throwType = types.Never{}
	}

	if async && generator {
		c.addFailure(
			"async generators are illegal",
			location,
		)
	}

	if generator {
		returnType = types.NewGenericWithTypeArgs(
			c.runtimeEnv.StdSubtypeClass(symbol.C_Generator),
			returnType,
			throwType,
		)

		throwType = types.Never{}
	} else if async {
		returnType = types.NewGenericWithTypeArgs(
			c.runtimeEnv.StdSubtypeClass(symbol.C_Promise),
			returnType,
			throwType,
		)

		throwType = types.Never{}
	}

	var flags bitfield.BitFlag16
	if abstract {
		flags |= types.METHOD_ABSTRACT_FLAG
	}
	if sealed {
		flags |= types.METHOD_SEALED_FLAG
	}
	if c.IsHeader() {
		flags |= types.METHOD_NATIVE_FLAG
	}
	if generator {
		flags |= types.METHOD_GENERATOR_FLAG
	}
	if async {
		flags |= types.METHOD_ASYNC_FLAG
	}
	if overload {
		flags |= types.METHOD_OVERLOAD_FLAG
	}
	if pure {
		flags |= types.METHOD_PURE_FLAG
	}
	newMethod := types.NewMethod(
		docComment,
		flags,
		name,
		typeParams,
		params,
		returnType,
		throwType,
		methodNamespace,
	)
	newMethod.SetLocation(location)

	if overload {
		if oldMethod != nil {
			if !oldMethod.IsOverload() {
				c.addFailure(
					fmt.Sprintf(
						"cannot declare an overload for method `%s` previously defined without `%s`",
						name.String(),
						lexer.Colorize("overload"),
					),
					location,
				)
			}
			oldMethod.RegisterOverload(newMethod)
		} else {
			methodNamespace.SetMethod(name, newMethod)
		}
	} else {
		c.checkMethodOverrideWithPlaceholder(newMethod, oldMethod, location)
		methodNamespace.SetMethod(name, newMethod)
	}

	c.checkSpecialMethods(name, newMethod, paramNodes, location)

	c.mode = prevMode

	return newMethod, typeParamMod
}

func (c *Checker) checkMethodOverrideWithPlaceholder(
	overrideMethod,
	baseMethod *types.Method,
	location *position.Location,
) {
	if baseMethod == nil {
		return
	}

	if baseMethod.IsPlaceholder() {
		baseMethod.SetReplaced(true)
		baseMethod.DefinedUnder.Get().SetMethod(baseMethod.Name, overrideMethod)
		return
	}

	c.checkMethodOverride(
		overrideMethod,
		baseMethod,
		location,
	)
	overrideMethod.UsedInConstants.ConcatMut(baseMethod.UsedInConstants)
}

// Set the mode of a closure/method output position
func (c *Checker) setOutputPositionTypeMode() {
	switch c.mode {
	case inputPositionTypeMode:
		// input position of a closure in an output position of a method
		// is an input position
		c.mode = inputPositionTypeMode
	default:
		// output position of a closure in an output position of a method
		// is an output position
		c.mode = outputPositionTypeMode
	}
}

// Set the mode of a closure/method input position
func (c *Checker) setInputPositionTypeMode() {
	switch c.mode {
	case inputPositionTypeMode:
		// input position of a closure in an input position of a method
		// is an output position
		c.mode = outputPositionTypeMode
	default:
		// output position of a closure in an input position of a method
		// is an input position
		c.mode = inputPositionTypeMode
	}
}

func (c *Checker) checkMethodCompatibilityForAlgebraicTypes(baseMethod, overrideMethod *types.Method, errSpan *position.Location, widenReturnType, checkPurity bool) bool {
	if !overrideMethod.IsGeneric() {
		return c.checkMethodCompatibility(baseMethod, overrideMethod, errSpan, true, widenReturnType, checkPurity)
	}

	prevMode := c.mode
	c.mode = methodCompatibilityInAlgebraicTypeMode

	typeArgs := make(types.TypeArgumentMap)
	if !c.checkMethodCompatibilityAndInferTypeArgs(baseMethod, overrideMethod, errSpan, typeArgs) {
		return false
	}

	c.mode = prevMode

	return typeArgs.HasAllTypeParams(overrideMethod.TypeParameters)
}

func (c *Checker) checkMethodCompatibilityForInterfaceIntersection(baseMethod, overrideMethod *types.Method, errSpan *position.Location, typeArgs types.TypeArgumentMap) bool {
	areCompatible := c.checkMethodCompatibilityAndInferTypeArgs(baseMethod, overrideMethod, errSpan, typeArgs)
	return areCompatible
}

// Checks whether two methods are compatible.
func (c *Checker) checkMethodCompatibility(baseMethod, overrideMethod *types.Method, errSpan *position.Location, validateParamNames, widenReturnType, checkPurity bool) bool {
	if baseMethod == nil {
		return true
	}

	areCompatible := true
	errDetailsBuff := new(strings.Builder)

	if checkPurity {
		if baseMethod.IsPure() && !overrideMethod.IsPure() {
			fmt.Fprintf(
				errDetailsBuff,
				"\n  - method `%s` is impure while `%s` is pure",
				types.InspectWithColor(overrideMethod),
				types.InspectWithColor(baseMethod),
			)
			areCompatible = false
		}
	}

	if !widenReturnType {
		overrideReturnType := overrideMethod.ReturnType.Get()
		baseReturnType := baseMethod.ReturnType.Get()
		if !c.isSubtype(overrideReturnType, baseReturnType, errSpan) {
			fmt.Fprintf(
				errDetailsBuff,
				"\n  - method `%s` has a different return type than `%s`, has `%s`, should have `%s`",
				types.InspectWithColor(overrideMethod),
				types.InspectWithColor(baseMethod),
				types.InspectWithColor(overrideReturnType),
				types.InspectWithColor(baseReturnType),
			)
			areCompatible = false
		}
		overrideThrowType := overrideMethod.ThrowType.Get()
		baseThrowType := baseMethod.ThrowType.Get()
		if !c.isSubtype(overrideThrowType, baseThrowType, errSpan) {
			fmt.Fprintf(
				errDetailsBuff,
				"\n  - method `%s` has a different throw type than `%s`, has `%s`, should have `%s`",
				types.InspectWithColor(overrideMethod),
				types.InspectWithColor(baseMethod),
				types.InspectWithColor(overrideThrowType),
				types.InspectWithColor(baseThrowType),
			)
			areCompatible = false
		}
	}

	if len(baseMethod.Params) > len(overrideMethod.Params) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` has less parameters than `%s`, has `%d`, should have `%d`",
			types.InspectWithColor(overrideMethod),
			types.InspectWithColor(baseMethod),
			len(overrideMethod.Params),
			len(baseMethod.Params),
		)
		areCompatible = false
	} else {
		for i := range len(baseMethod.Params) {
			oldParam := baseMethod.Params[i].Get()
			newParam := overrideMethod.Params[i].Get()

			if (validateParamNames && oldParam.Name != newParam.Name) || oldParam.Kind != newParam.Kind || !c.isSubtype(oldParam.Type.Get(), newParam.Type.Get(), errSpan) {
				fmt.Fprintf(
					errDetailsBuff,
					"\n  - method `%s` has an incompatible parameter with `%s`, has `%s`, should have `%s`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					types.InspectWithColor(newParam),
					types.InspectWithColor(oldParam),
				)
				areCompatible = false
			}
		}

		for i := len(baseMethod.Params); i < len(overrideMethod.Params); i++ {
			param := overrideMethod.Params[i].Get()
			if !param.IsOptional() {
				fmt.Fprintf(
					errDetailsBuff,
					"\n  - method `%s` has a required parameter missing in `%s`, got `%s`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					param.Name.String(),
				)
				areCompatible = false
			}
		}
	}

	if !areCompatible {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` is incompatible with `%s`\n  is:        `%s`\n  should be: `%s`\n%s",
				types.InspectWithColor(overrideMethod),
				types.InspectWithColor(baseMethod),
				overrideMethod.InspectSignatureWithColor(false),
				baseMethod.InspectSignatureWithColor(false),
				errDetailsBuff.String(),
			),
			errSpan,
		)
	}

	return areCompatible
}

func (c *Checker) checkMethodCompatibilityAndInferTypeArgs(baseMethod, overrideMethod *types.Method, errSpan *position.Location, typeArgs types.TypeArgumentMap) bool {
	if baseMethod == nil {
		return true
	}

	areCompatible := true
	errDetailsBuff := new(strings.Builder)

	baseReturnType := baseMethod.ReturnType.Get()
	overrideReturnType := overrideMethod.ReturnType.Get()
	returnType := c.inferTypeArguments(baseReturnType, overrideReturnType, typeArgs, nil)
	if returnType == nil || !c.isSubtype(returnType, baseReturnType, errSpan) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` has a different return type than `%s`, has `%s`, should have `%s`",
			types.InspectWithColor(overrideMethod),
			types.InspectWithColor(baseMethod),
			types.InspectWithColor(overrideReturnType),
			types.InspectWithColor(baseReturnType),
		)
		areCompatible = false
	}

	baseThrowType := baseMethod.ThrowType.Get()
	overrideThrowType := overrideMethod.ThrowType.Get()
	throwType := c.inferTypeArguments(baseThrowType, overrideThrowType, typeArgs, nil)
	if throwType == nil || !c.isSubtype(throwType, baseThrowType, errSpan) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` has a different throw type than `%s`, has `%s`, should have `%s`",
			types.InspectWithColor(overrideMethod),
			types.InspectWithColor(baseMethod),
			types.InspectWithColor(overrideThrowType),
			types.InspectWithColor(baseThrowType),
		)
		areCompatible = false
	}

	if len(baseMethod.Params) > len(overrideMethod.Params) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` has less parameters than `%s`, has `%d`, should have `%d`",
			types.InspectWithColor(overrideMethod),
			types.InspectWithColor(baseMethod),
			len(overrideMethod.Params),
			len(baseMethod.Params),
		)
		areCompatible = false
	} else {
		for i := range len(baseMethod.Params) {
			oldParam := baseMethod.Params[i].Get()
			newParam := overrideMethod.Params[i].Get()

			oldParamType := oldParam.Type.Get()
			newParamType := newParam.Type.Get()
			newParamType = c.inferTypeArguments(oldParamType, newParamType, typeArgs, nil)
			if oldParam.Name != newParam.Name || oldParam.Kind != newParam.Kind ||
				newParamType == nil || !c.isSubtype(oldParamType, newParamType, errSpan) {
				fmt.Fprintf(
					errDetailsBuff,
					"\n  - method `%s` has an incompatible parameter with `%s`, has `%s`, should have `%s`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					types.InspectWithColor(newParam),
					types.InspectWithColor(oldParam),
				)
				areCompatible = false
			}
		}

		for i := len(baseMethod.Params); i < len(overrideMethod.Params); i++ {
			param := overrideMethod.Params[i].Get()
			if !param.IsOptional() {
				fmt.Fprintf(
					errDetailsBuff,
					"\n  - method `%s` has a required parameter missing in `%s`, got `%s`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					param.Name.String(),
				)
				areCompatible = false
			}
		}
	}

	if !areCompatible {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` is incompatible with `%s`\n  is:        `%s`\n  should be: `%s`\n%s",
				types.InspectWithColor(overrideMethod),
				types.InspectWithColor(baseMethod),
				overrideMethod.InspectSignatureWithColor(false),
				baseMethod.InspectSignatureWithColor(false),
				errDetailsBuff.String(),
			),
			errSpan,
		)
	}

	return areCompatible
}

func (c *Checker) GetMethod(typ types.Type, name symbol.Symbol, errSpan *position.Location) *types.Method {
	return c._getMethod(typ, name, errSpan, false, false)
}

// Iterates over every method of the namespace, resolving type parameters.
func (c *Checker) methodsInNamespace(namespace types.Namespace) iter.Seq2[symbol.Symbol, *types.Method] {
	return func(yield func(name symbol.Symbol, method *types.Method) bool) {
		var generics []*types.Generic
		seenMethods := make(ds.Set[symbol.Symbol])

		for parent := range types.Parents(namespace) {
			if generic, ok := parent.(*types.Generic); ok {
				generics = append(generics, generic)
			}
			methods := parent.Methods()
			names := symbol.SortKeys(methods)
		methodLoop:
			for _, name := range names {
				method := methods[name].Get()
				if seenMethods.Contains(name) {
					continue
				}
				if len(generics) < 1 {
					if !yield(name, method) {
						return
					}
					seenMethods.Add(name)
					continue
				}

				var whereParams []types.Ref[*types.TypeParameter]
				var whereArgs []types.Ref[types.Type]
				if mixinWithWhere, ok := parent.(*types.MixinWithWhere); ok {
					whereParams = slices.Clone(mixinWithWhere.Where)
					whereArgs = c.constructWhereArguments(whereParams)
				}

				var methodCopy *types.Method
				for i := len(generics) - 1; i >= 0; i-- {
					generic := generics[i]
					c.replaceTypeParametersInWhere(whereParams, whereArgs, generic.ArgumentMap)
					if methodCopy != nil {
						c.replaceTypeParametersInMethod(methodCopy, generic.ArgumentMap, false)
						continue
					}

					result := c.replaceTypeParametersInMethodCopy(method, generic.ArgumentMap, false)
					if result != method {
						methodCopy = result
						method = result
					}
				}

				for i := range len(whereParams) {
					whereParam := whereParams[i].Get()
					whereArg := whereArgs[i].Get()

					if !c.isSubtype(whereParam.LowerBound.Get(), whereArg, nil) {
						continue methodLoop
					}
					if !c.isSubtype(whereArg, whereParam.UpperBound.Get(), nil) {
						continue methodLoop
					}
				}

				if !yield(name, method) {
					return
				}
				seenMethods.Add(name)
			}
		}
	}
}

// Iterates over every abstract method of the namespace, resolving type parameters.
func (c *Checker) abstractMethodsInNamespace(namespace types.Namespace) iter.Seq2[symbol.Symbol, *types.Method] {
	return func(yield func(name symbol.Symbol, method *types.Method) bool) {
		var generics []*types.Generic
		seenMethods := make(ds.Set[symbol.Symbol])

		for parent := range types.Parents(namespace) {
			if generic, ok := parent.(*types.Generic); ok {
				generics = append(generics, generic)
			}
			if !parent.IsAbstract() {
				continue
			}
			for name, methodRef := range parent.Methods() {
				method := methodRef.Get()
				if !method.IsAbstract() {
					continue
				}
				if seenMethods.Contains(name) {
					continue
				}
				if len(generics) < 1 {
					if !yield(name, method) {
						return
					}
					seenMethods.Add(name)
					continue
				}

				var methodCopy *types.Method
				for i := len(generics) - 1; i >= 0; i-- {
					generic := generics[i]
					if methodCopy != nil {
						c.replaceTypeParametersInMethod(methodCopy, generic.ArgumentMap, false)
						continue
					}

					result := c.replaceTypeParametersInMethodCopy(method, generic.ArgumentMap, false)
					if result != method {
						methodCopy = result
						method = result
					}
				}
				if !yield(name, method) {
					return
				}
				seenMethods.Add(name)
			}
		}
	}
}

func (c *Checker) resolveMethodInNamespace(namespace types.Namespace, name symbol.Symbol) *types.Method {
	var generics []*types.Generic

	for parent := range types.Parents(namespace) {
		switch p := parent.(type) {
		case *types.Generic:
			generics = append(generics, p)
		case *types.NamespacePlaceholder:
			switch n := p.Namespace.Get().(type) {
			case *types.Module:
				parent = n
			default:
				if n.Singleton() == nil {
					continue
				}
				parent = n.Singleton()
			}
		}

		method := parent.Method(name)
		if method == nil {
			continue
		}

		var whereParams []types.Ref[*types.TypeParameter]
		var whereArgs []types.Ref[types.Type]
		if mixinWithWhere, ok := parent.(*types.MixinWithWhere); ok {
			if len(generics) < 1 {
				return nil
			}
			whereParams = slices.Clone(mixinWithWhere.Where)
			whereArgs = c.constructWhereArguments(whereParams)
		}

		if len(generics) < 1 {
			return method
		}

		var methodCopy *types.Method
		for i := len(generics) - 1; i >= 0; i-- {
			generic := generics[i]
			c.replaceTypeParametersInWhere(whereParams, whereArgs, generic.ArgumentMap)
			if methodCopy != nil {
				c.replaceTypeParametersInMethod(methodCopy, generic.ArgumentMap, false)
				continue
			}

			result := c.replaceTypeParametersInMethodCopy(method, generic.ArgumentMap, false)
			if result != method {
				methodCopy = result
				method = result
			}
		}

		for i := range len(whereParams) {
			whereParam := whereParams[i].Get()
			whereArg := whereArgs[i].Get()

			if !c.isSubtype(whereParam.LowerBound.Get(), whereArg, nil) {
				return nil
			}
			if !c.isSubtype(whereArg, whereParam.UpperBound.Get(), nil) {
				return nil
			}
		}

		return method
	}

	return nil
}

func (c *Checker) constructWhereArguments(whereParameters []types.Ref[*types.TypeParameter]) []types.Ref[types.Type] {
	whereArgs := make([]types.Ref[types.Type], len(whereParameters))
	for i, whereParam := range whereParameters {
		whereArgs[i] = types.Ref[types.Type](whereParam)
	}

	return whereArgs
}

func (c *Checker) resolveNonAbstractMethodInNamespace(namespace types.Namespace, name symbol.Symbol) *types.Method {
	var generics []*types.Generic

	for parent := range types.Parents(namespace) {
		if generic, ok := parent.(*types.Generic); ok {
			generics = append(generics, generic)
		}
		method := parent.Method(name)
		if method != nil {
			if method.IsAbstract() {
				continue
			}
			if len(generics) < 1 {
				return method
			}

			var methodCopy *types.Method
			for i := len(generics) - 1; i >= 0; i-- {
				generic := generics[i]
				if methodCopy != nil {
					c.replaceTypeParametersInMethod(methodCopy, generic.ArgumentMap, false)
					continue
				}

				result := c.replaceTypeParametersInMethodCopy(method, generic.ArgumentMap, false)
				if result != method {
					methodCopy = result
					method = methodCopy
				}
			}
			return method
		}
	}

	return nil
}

func (c *Checker) _getMethodInNamespace(namespace types.Namespace, typ types.Type, name symbol.Symbol, errSpan *position.Location, inParent bool) *types.Method {
	method := c.resolveMethodInNamespace(namespace, name)
	if method != nil {
		return method
	}
	if !inParent {
		c.addMissingMethodError(typ, name.String(), errSpan)
	}
	return nil
}

func (c *Checker) createTypeArgumentMapWithSelf(self types.Type) types.TypeArgumentMap {
	return types.TypeArgumentMap{
		symbol.L_self: types.NewTypeArgument(
			self,
			types.INVARIANT,
		),
	}
}

func (c *Checker) getMethodInNamespaceWithSelf(namespace types.Namespace, typ types.Type, name symbol.Symbol, self types.Type, errSpan *position.Location, inParent, inSelf bool) *types.Method {
	method := c._getMethodInNamespace(namespace, typ, name, errSpan, inParent)
	if method == nil {
		return nil
	}
	if inSelf {
		return method
	}
	m := c.createTypeArgumentMapWithSelf(self)
	return c.replaceTypeParametersInMethodCopy(method, m, false)
}

func (c *Checker) getMethodInNamespace(namespace types.Namespace, typ types.Type, name symbol.Symbol, errSpan *position.Location, inParent, inSelf bool) *types.Method {
	return c.getMethodInNamespaceWithSelf(namespace, typ, name, namespace, errSpan, inParent, inSelf)
}

func (c *Checker) replaceTypeParametersInMethodCopy(method *types.Method, typeArgs types.TypeArgumentMap, replaceMethodTypeParams bool) *types.Method {
	var methodCopy *types.Method

	for i, typeParamRef := range method.TypeParameters {
		typeParam := typeParamRef.Get()
		lowerBound := typeParam.LowerBound.Get()
		newLowerBound := c.replaceTypeParameters(lowerBound, typeArgs, replaceMethodTypeParams)
		upperBound := typeParam.UpperBound.Get()
		newUpperBound := c.replaceTypeParameters(upperBound, typeArgs, replaceMethodTypeParams)

		var typeParamCopy *types.TypeParameter
		if lowerBound != newLowerBound {
			if methodCopy == nil {
				methodCopy = method.Copy()
			}
			typeParamCopy = methodCopy.TypeParameters[i].Get().Copy()
			typeParamCopy.LowerBound = types.ToRef(newLowerBound)
		}
		if upperBound != newUpperBound {
			if methodCopy == nil {
				methodCopy = method.Copy()
			}
			if typeParamCopy == nil {
				typeParamCopy = methodCopy.TypeParameters[i].Get().Copy()
			}
			typeParamCopy.UpperBound = types.ToRef(newUpperBound)
			methodCopy.TypeParameters[i] = typeParamCopy.ToRef()
		}
	}

	methodReturnType := method.ReturnType.Get()
	result := c.replaceTypeParameters(methodReturnType, typeArgs, replaceMethodTypeParams)
	if methodReturnType != result {
		if methodCopy == nil {
			methodCopy = method.Copy()
		}
		methodCopy.ReturnType = types.ToRef(result)
	}

	methodThrowType := method.ThrowType.Get()
	result = c.replaceTypeParameters(methodThrowType, typeArgs, replaceMethodTypeParams)
	if methodThrowType != result {
		if methodCopy == nil {
			methodCopy = method.Copy()
		}
		methodCopy.ThrowType = types.ToRef(result)
	}

	for i, paramRef := range method.Params {
		param := paramRef.Get()
		paramType := param.Type.Get()
		result := c.replaceTypeParameters(paramType, typeArgs, replaceMethodTypeParams)
		if paramType != result {
			if methodCopy == nil {
				methodCopy = method.Copy()
			}
			paramCopy := methodCopy.Params[i].Get().Copy()
			paramCopy.Type = types.ToRef(result)
			methodCopy.Params[i] = paramCopy.ToRef()
		}
	}

	different := methodCopy != nil
	var overloadsCopy []types.Ref[*types.Method]

	if different {
		for i, overloadRef := range method.Overloads {
			overload := overloadRef.Get()
			overloadCopy := c.replaceTypeParametersInMethod(overload, typeArgs, replaceMethodTypeParams)
			methodCopy.Overloads[i] = overloadCopy.ToRef()
		}
	} else {
		overloadsCopy = make([]types.Ref[*types.Method], len(method.Overloads))
		for i, overloadRef := range method.Overloads {
			overload := overloadRef.Get()
			overloadCopy := c.replaceTypeParametersInMethodCopy(overload, typeArgs, replaceMethodTypeParams)
			if overload != overloadCopy {
				different = true
			}
			overloadsCopy[i] = types.ToRef(overloadCopy)
		}
	}

	if different {
		if methodCopy == nil {
			methodCopy = method.Copy()
			methodCopy.Overloads = overloadsCopy
		}
		return methodCopy
	}

	return method
}

func (c *Checker) replaceTypeParametersInMethod(method *types.Method, typeArgs types.TypeArgumentMap, replaceMethodTypeParams bool) *types.Method {
	for _, typeParamRef := range method.TypeParameters {
		typeParam := typeParamRef.Get()

		lowerBound := typeParam.LowerBound.Get()
		typeParam.LowerBound = types.ToRef(c.replaceTypeParameters(lowerBound, typeArgs, replaceMethodTypeParams))

		upperBound := typeParam.UpperBound.Get()
		typeParam.UpperBound = types.ToRef(c.replaceTypeParameters(upperBound, typeArgs, replaceMethodTypeParams))
	}
	method.ReturnType = types.ToRef(c.replaceTypeParameters(method.ReturnType.Get(), typeArgs, replaceMethodTypeParams))
	method.ThrowType = types.ToRef(c.replaceTypeParameters(method.ThrowType.Get(), typeArgs, replaceMethodTypeParams))

	for _, paramRef := range method.Params {
		param := paramRef.Get()
		param.Type = types.ToRef(c.replaceTypeParameters(param.Type.Get(), typeArgs, replaceMethodTypeParams))
	}

	for i, overloadRef := range method.Overloads {
		overload := overloadRef.Get()
		method.Overloads[i] = types.ToRef(c.replaceTypeParametersInMethod(overload, typeArgs, replaceMethodTypeParams))
	}

	return method
}

func (c *Checker) replaceTypeParametersInWhere(whereParams []types.Ref[*types.TypeParameter], whereArgs []types.Ref[types.Type], typeArgs types.TypeArgumentMap) {
	for i, whereArg := range whereArgs {
		whereArgs[i] = types.ToRef(c.replaceTypeParameters(whereArg.Get(), typeArgs, false))
	}

	for i, whereParamRef := range whereParams {
		whereParam := whereParamRef.Get()
		var whereParamCopy *types.TypeParameter

		whereParamLowerBound := whereParam.LowerBound.Get()
		result := c.replaceTypeParameters(whereParamLowerBound, typeArgs, false)
		if result != whereParamLowerBound {
			whereParamCopy = whereParam.Copy()
			whereParamCopy.LowerBound = types.ToRef(result)
		}

		whereParamUpperBound := whereParam.UpperBound.Get()
		result = c.replaceTypeParameters(whereParamUpperBound, typeArgs, false)
		if result != whereParamUpperBound {
			if whereParamCopy == nil {
				whereParamCopy = whereParam.Copy()
			}
			whereParamCopy.UpperBound = types.ToRef(result)
		}

		if whereParamCopy != nil {
			whereParams[i] = whereParamCopy.ToRef()
		}
	}
}

func (c *Checker) getMethodForTypeParameter(typ *types.TypeParameter, name symbol.Symbol, errSpan *position.Location, inParent, inSelf bool) *types.Method {
	upperBound := typ.UpperBound.Get()
	switch upper := upperBound.(type) {
	case *types.Class:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.Mixin:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.Interface:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.Callable:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.SingletonClass:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.Generic:
		var method *types.Method
		switch genericType := upper.Namespace.Get().(type) {
		case *types.Class:
			method = c._getMethodInNamespace(genericType, typ, name, errSpan, inParent)
		case *types.Mixin:
			method = c._getMethodInNamespace(genericType, typ, name, errSpan, inParent)
		case *types.Interface:
			method = c._getMethodInNamespace(genericType, typ, name, errSpan, inParent)
		}
		if method == nil {
			return nil
		}

		typeArgMap := maps.Clone(upper.TypeArguments.ArgumentMap)
		typeArgMap[symbol.L_self] = types.NewTypeArgument(
			typ,
			types.INVARIANT,
		)
		return c.replaceTypeParametersInMethodCopy(method, typeArgMap, true)
	default:
		return c._getMethod(upperBound, name, errSpan, inParent, inSelf)
	}
}

func (c *Checker) getReceiverlessMethod(name symbol.Symbol, location *position.Location) (_ *types.Method, namespace types.Namespace, fromLocal bool) {
	nameStr := name.String()
	local, _ := c.resolveLocal(nameStr, nil)
	if local != nil {
		if !local.initialised {
			c.addUninitialisedLocalError(nameStr, location)
		}
		return c.GetMethod(local.typ.Get(), symbol.L_call, location), nil, true
	}
	method := c.GetMethod(c.selfType.Get(), name, nil)
	if method != nil {
		return method, nil, false
	}

	for _, methodScope := range c.methodScopes {
		switch methodScope.kind {
		case scopeUsingBufferKind, scopeUsingKind:
		default:
			continue
		}

		namespace := methodScope.container.Get()
		method := c.GetMethod(namespace, name, nil)
		if method != nil {
			return method, namespace, false
		}
	}

	c.addMissingMethodError(c.selfType.Get(), name.String(), location)

	return nil, nil, false
}

func (c *Checker) _getMethod(typ types.Type, name symbol.Symbol, errLoc *position.Location, inParent, inSelf bool) *types.Method {
	typ = c.ToNonLiteral(typ, true)

	switch t := typ.(type) {
	case types.Self:
		return c._getMethod(c.selfType.Get(), name, errLoc, inParent, true)
	case *types.NamedType:
		return c._getMethod(t.Type.Get(), name, errLoc, inParent, inSelf)
	case *types.TypeParameter:
		return c.getMethodForTypeParameter(t, name, errLoc, inParent, inSelf)
	case *types.Generic:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.Class:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.NamespacePlaceholder:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.SingletonClass:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.Interface:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.Callable:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.InterfaceProxy:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.Module:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.Mixin:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.MixinProxy:
		return c.getMethodInNamespace(t, typ, name, errLoc, inParent, inSelf)
	case *types.Intersection:
		return c.getMethodInIntersection(t, name, errLoc)
	case *types.Nilable:
		return c.getMethodInNilable(t, name, errLoc)
	case *types.Union:
		return c.getMethodInUnion(t, name, errLoc)
	default:
		c.addMissingMethodError(typ, name.String(), errLoc)
		return nil
	}
}

func (c *Checker) getMethodInNilable(typ *types.Nilable, name symbol.Symbol, errLoc *position.Location) *types.Method {
	nilType := c.runtimeEnv.StdSubtype(symbol.C_Nil).(*types.Class)
	nilMethod := nilType.Method(name)
	if nilMethod == nil {
		c.addMissingMethodError(nilType, name.String(), errLoc)
	}
	nonNilMethod := c.GetMethod(typ.Type.Get(), name, errLoc)
	if nilMethod == nil || nonNilMethod == nil {
		return nil
	}

	var baseMethod *types.Method
	var overrideMethod *types.Method
	if calculateMethodBaseScore(nilMethod) > calculateMethodBaseScore(nonNilMethod) {
		baseMethod = nilMethod
		overrideMethod = nonNilMethod
	} else {
		baseMethod = nonNilMethod
		overrideMethod = nilMethod
	}

	if !c.checkMethodCompatibilityForAlgebraicTypes(baseMethod, overrideMethod, errLoc, true, false) {
		return nil
	}

	method := baseMethod.Copy()
	method.ReturnType = types.ToRef(c.NewNormalisedUnion(baseMethod.ReturnType, overrideMethod.ReturnType))
	method.ThrowType = types.ToRef(c.NewNormalisedUnion(baseMethod.ThrowType, overrideMethod.ThrowType))
	method.SetPure(baseMethod.IsPure() && overrideMethod.IsPure())
	method.Overloads = nil
	return method
}

func (c *Checker) getMethodInUnion(typ *types.Union, name symbol.Symbol, errLoc *position.Location) *types.Method {
	var methods []*types.Method
	var baseMethod *types.Method

	for _, element := range typ.Elements {
		elementMethod := c.GetMethod(element.Get(), name, errLoc)
		if elementMethod == nil {
			continue
		}
		methods = append(methods, elementMethod)
		if baseMethod == nil {
			baseMethod = elementMethod
			continue
		}
		currentScore := calculateMethodBaseScore(elementMethod)
		baseScore := calculateMethodBaseScore(baseMethod)
		if currentScore > baseScore {
			baseMethod = elementMethod
		}
	}

	if len(methods) < len(typ.Elements) {
		return nil
	}

	returnTypes := make([]types.Ref[types.Type], len(methods)+1)
	throwTypes := make([]types.Ref[types.Type], len(methods)+1)

	returnTypes[0] = baseMethod.ReturnType
	throwTypes[0] = baseMethod.ThrowType
	isPure := true

	isCompatible := true
	for i := range len(methods) {
		method := methods[i]

		ok := c.checkMethodCompatibilityForAlgebraicTypes(baseMethod, method, errLoc, true, false)
		returnTypes[i+1] = method.ReturnType
		throwTypes[i+1] = method.ThrowType
		if !method.IsPure() {
			isPure = false
		}
		if !ok {
			isCompatible = false
		}
	}

	if !isCompatible {
		return nil
	}

	method := baseMethod.Copy()
	method.ReturnType = types.ToRef(c.NewNormalisedUnion(returnTypes...))
	method.ThrowType = types.ToRef(c.NewNormalisedUnion(throwTypes...))
	method.SetPure(isPure)
	method.Overloads = nil
	return method
}

func (c *Checker) getMethodInIntersection(typ *types.Intersection, name symbol.Symbol, errLoc *position.Location) *types.Method {
	var methods []*types.Method
	var baseMethod *types.Method

	for _, elementRef := range typ.Elements {
		element := elementRef.Get()
		switch e := element.(type) {
		case *types.Not:
			switch t := e.Type.Get().(type) {
			case *types.Interface:
				elementMethod := c.GetMethod(t, name, nil)
				if elementMethod == nil {
					continue
				}
				return nil
			case *types.Mixin:
				elementMethod := c.GetMethod(t, name, nil)
				if elementMethod == nil {
					continue
				}
				return nil
			}
		default:
			elementMethod := c.GetMethod(element, name, nil)
			if elementMethod == nil {
				continue
			}
			methods = append(methods, elementMethod)
			if baseMethod == nil || len(baseMethod.Params) > len(elementMethod.Params) || baseMethod.IsGeneric() && !elementMethod.IsGeneric() {
				baseMethod = elementMethod
			}
		}
	}

	switch len(methods) {
	case 0:
		c.addMissingMethodError(typ, name.String(), errLoc)
		return nil
	case 1:
		return methods[0].WithoutOverloads()
	}

	isCompatible := true
	for i := range len(methods) {
		method := methods[i]

		if !c.checkMethodCompatibilityForAlgebraicTypes(baseMethod, method, errLoc, false, true) {
			isCompatible = false
		}
	}

	if isCompatible {
		return baseMethod.WithoutOverloads()
	}

	return nil
}

func calculateMethodBaseScore(typ types.Type) int {
	switch t := typ.(type) {
	case *types.Union:
		var sum int
		for _, element := range t.Elements {
			sum += calculateMethodBaseScore(element.Get())
		}
		return sum
	case *types.Intersection:
		var sum int
		for _, element := range t.Elements {
			sum += calculateMethodBaseScore(element.Get())
		}

		return sum/len(t.Elements) - 1
	case *types.Class:
		result := 2
		if t.IsGeneric() {
			result += 1
		}
		if t.IsAbstract() {
			result += 1
		}
		return result
	case *types.Mixin:
		result := 3
		if t.IsGeneric() {
			result += 1
		}
		if t.IsAbstract() {
			result += 1
		}
		return result
	case *types.Interface:
		result := 4
		if t.IsGeneric() {
			result += 1
		}
		return result
	case *types.InstanceOf:
		return 2
	case *types.Nilable:
		return calculateMethodBaseScore(t.Type.Get()) + 1
	case *types.Never:
		return -100
	case *types.Generic:
		return calculateMethodBaseScore(t.Namespace.Get())
	case *types.NamedType:
		return calculateMethodBaseScore(t.Type.Get())
	case *types.Not:
		return 100 - calculateMethodBaseScore(t.Type.Get())
	case *types.Any:
		return 100
	case *types.Callable:
		return calculateMethodBaseScore(t.Body.Get())
	case *types.Method:
		var result int

		if t.IsGeneric() {
			result -= 200
		}

		var paramScore int
		for _, param := range t.Params {
			paramScore += calculateMethodBaseScore(param.Get())
		}
		result -= paramScore

		var returnScore int
		returnScore += calculateMethodBaseScore(t.ReturnType.Get())
		returnScore += calculateMethodBaseScore(t.ThrowType.Get())
		result += returnScore / 2

		return result
	default:
		return 1
	}
}

func (c *Checker) addMissingMethodError(typ types.Type, name string, location *position.Location) {
	if types.IsUntyped(typ) {
		return
	}
	c.addFailure(
		fmt.Sprintf("method `%s` is not defined on type `%s`", name, types.InspectWithColor(typ)),
		location,
	)
}
