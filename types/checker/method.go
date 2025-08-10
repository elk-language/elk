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
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

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
	method, mod := c.declareMethod(
		nil,
		c.currentMethodScope().container,
		initNode.DocComment(),
		false,
		false,
		false,
		false,
		false,
		symbol.S_init,
		nil,
		initNode.Parameters,
		nil,
		initNode.ThrowType,
		initNode.Location(),
	)
	initNode.SetType(method)
	newNode := ast.NewMethodDefinitionNode(
		initNode.Location(),
		initNode.DocComment(),
		0,
		ast.NewPublicIdentifierNode(initNode.Location(), "#init"),
		nil,
		initNode.Parameters,
		nil,
		initNode.ThrowType,
		initNode.Body,
	)
	newNode.SetType(method)
	c.registerMethodCheck(method, newNode)
	if mod != nil {
		c.popConstScope()
	}
	return newNode
}

func (c *Checker) hoistAliasDeclaration(node *ast.AliasDeclarationNode) {
	node.SetType(types.Untyped{})
	namespace := c.currentMethodScope().container
	for _, entry := range node.Entries {
		c.hoistAliasEntry(entry, namespace)
	}
}

func (c *Checker) hoistAliasEntry(node *ast.AliasDeclarationEntry, namespace types.Namespace) {
	oldName := c.identifierToName(node.OldName)
	oldNameSymbol := value.ToSymbol(oldName)
	aliasedMethod := namespace.Method(oldNameSymbol)
	if aliasedMethod == nil {
		c.addMissingMethodError(namespace, oldName, node.Location())
		return
	}
	newName := value.ToSymbol(c.identifierToName(node.NewName))
	oldMethod := c.resolveMethodInNamespace(namespace, newName)
	c.checkMethodOverrideWithPlaceholder(aliasedMethod, oldMethod, node.Location())
	c.checkSpecialMethods(newName, aliasedMethod, nil, node.Location())
	namespace.SetMethodAlias(newName, aliasedMethod)
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

	originalMethodSymbol := value.ToSymbol(methodName)
	var newMethodSymbol value.Symbol
	if asName != "" {
		newMethodSymbol = value.ToSymbol(asName)
	} else {
		newMethodSymbol = value.ToSymbol(methodName)
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
	for _, entry := range node.Entries {
		c.resolveUsingEntry(entry, true)
	}
}

func (c *Checker) resolveUsingEntry(entry ast.UsingEntryNode, pushPlaceholderLocation bool) {
	typ := c.TypeOf(entry)
	switch t := typ.(type) {
	case *types.Module:
		c.pushConstScope(makeUsingConstantScope(t))
		c.pushMethodScope(makeUsingMethodScope(t))
	case *types.Mixin:
		c.pushConstScope(makeUsingConstantScope(t))
		c.pushMethodScope(makeUsingMethodScope(t))
	case *types.Class:
		c.pushConstScope(makeUsingConstantScope(t))
		c.pushMethodScope(makeUsingMethodScope(t))
	case *types.Interface:
		c.pushConstScope(makeUsingConstantScope(t))
		c.pushMethodScope(makeUsingMethodScope(t))
	case *types.NamespacePlaceholder:
		c.pushConstScope(makeUsingConstantScope(t))
		c.pushMethodScope(makeUsingMethodScope(t))
		if pushPlaceholderLocation {
			t.Locations.Push(entry.Location())
		}
	case *types.UsingBufferNamespace:
		if c.enclosingScopeIsAUsingBuffer() {
			return
		}
		c.pushConstScope(makeUsingBufferConstantScope(t))
		c.pushMethodScope(makeUsingBufferMethodScope(t))
	}
}

func (c *Checker) hoistMethodDefinitionsWithinClass(node *ast.ClassDeclarationNode) {
	class, ok := c.TypeOf(node).(*types.Class)
	if ok {
		c.pushConstScope(makeLocalConstantScope(class))
		c.pushMethodScope(makeLocalMethodScope(class))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = classMode
	c.selfType = class
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
	module, ok := c.TypeOf(node).(*types.Module)
	if ok {
		c.pushConstScope(makeLocalConstantScope(module))
		c.pushMethodScope(makeLocalMethodScope(module))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = moduleMode
	c.selfType = module
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
	mixin, ok := c.TypeOf(node).(*types.Mixin)
	if ok {
		c.pushConstScope(makeLocalConstantScope(mixin))
		c.pushMethodScope(makeLocalMethodScope(mixin))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = mixinMode
	c.selfType = mixin
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinInterface(node *ast.InterfaceDeclarationNode) {
	iface, ok := c.TypeOf(node).(*types.Interface)
	if ok {
		c.pushConstScope(makeLocalConstantScope(iface))
		c.pushMethodScope(makeLocalMethodScope(iface))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = interfaceMode
	c.selfType = iface
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinSingleton(expr *ast.SingletonBlockExpressionNode) {
	namespace := c.currentConstScope().container
	singleton := namespace.Singleton()
	if singleton == nil {
		return
	}

	c.pushConstScope(makeLocalConstantScope(singleton))
	c.pushMethodScope(makeLocalMethodScope(singleton))

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = singletonMode
	c.selfType = singleton
	c.hoistMethodDefinitions(expr.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	c.registerNamespaceWithIvars(singleton, expr.Location())
	c.popLocalConstScope()
	c.popMethodScope()
}

func (c *Checker) hoistMethodDefinitionsWithinExtendWhere(node *ast.ExtendWhereBlockExpressionNode) {
	namespace, ok := c.TypeOf(node).(*types.MixinWithWhere)
	if ok {
		c.pushConstScope(makeLocalConstantScope(namespace))
		c.pushMethodScope(makeLocalMethodScope(namespace))
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
	definedUnder := c.currentMethodScope().container
	method, mod := c.declareMethod(
		nil,
		definedUnder,
		node.DocComment(),
		node.IsAbstract(),
		node.IsSealed(),
		false,
		node.IsGenerator(),
		node.IsAsync(),
		value.ToSymbol(c.identifierToName(node.Name)),
		node.TypeParameters,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Location(),
	)
	method.Node = node
	node.SetType(method)
	c.registerMethodCheck(method, node)
	if mod != nil {
		c.popConstScope()
	}
}

func (c *Checker) hoistMethodSignatureDefinition(node *ast.MethodSignatureDefinitionNode) {
	method, mod := c.declareMethod(
		nil,
		c.currentMethodScope().container,
		node.DocComment(),
		true,
		false,
		false,
		false,
		false,
		value.ToSymbol(c.identifierToName(node.Name)),
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
	filename string,
	constScopes []constantScope,
	methodScopes []methodScope,
	selfType,
	returnType,
	throwType types.Type,
	mode mode,
	threadPool *vm.ThreadPool,
	loc *position.Location,
) *Checker {
	checker := &Checker{
		env:            c.env,
		Filename:       filename,
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
			newLocalEnvironment(nil),
		},
		typeDefinitionChecks: newTypeDefinitionChecks(),
		methodCache:          concurrent.NewSlice[*types.Method](),
		threadPool:           threadPool,
	}
	checker.compiler = compiler.CreateCompiler(c.compiler, checker, loc, c.Errors)

	return checker
}

// Checks whether all methods specified in `using` statements have been defined
func (c *Checker) checkMethodPlaceholders() {
	for _, placeholder := range c.methodPlaceholders {
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

type methodCheckEntry struct {
	method         *types.Method
	constantScopes []constantScope
	methodScopes   []methodScope
	node           *ast.MethodDefinitionNode
	headerMode     bool
}

func (c *Checker) registerMethodCheck(method *types.Method, node *ast.MethodDefinitionNode) {
	c.methodChecks = append(c.methodChecks, methodCheckEntry{
		method:         method,
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
		node:           node,
		headerMode:     c.IsHeader(),
	})
}

var concurrencyLimit = 1000

func (c *Checker) checkMethods() {
	concurrent.Foreach(
		concurrencyLimit,
		c.methodChecks,
		func(methodCheck methodCheckEntry) {
			method := methodCheck.method
			node := methodCheck.node

			var mode mode
			if method.IsInit() {
				mode = initMode
			} else {
				mode = methodMode
			}

			methodChecker := c.newMethodChecker(
				node.Location().FilePath,
				methodCheck.constantScopes,
				methodCheck.methodScopes,
				method.DefinedUnder,
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
				c.methodCache.Push(method)
			}
		},
	)

	c.methodChecks = nil
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
	for _, method := range c.methodCache.Slice {
		c.checkMethodInConstant(method, method.UsedInConstants)
	}
}

func (c *Checker) checkMethodInConstant(method *types.Method, usedInConstants ds.Set[value.Symbol]) {
	for _, calledMethod := range method.CalledMethods {
		c.checkMethodInConstant(calledMethod, usedInConstants)
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

func (c *Checker) declareMethodForGetter(node *ast.AttributeParameterNode, docComment string) {
	name := c.identifierToName(node.Name)
	method, mod := c.declareMethod(
		nil,
		c.currentMethodScope().container,
		docComment,
		false,
		false,
		false,
		false,
		false,
		value.ToSymbol(name),
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
	c.registerMethodCheck(
		method,
		methodNode,
	)
	if mod != nil {
		c.popConstScope()
	}
}

// Create a deep copy of the method
func (c *Checker) deepCopyMethod(method *types.Method) *types.Method {
	if method.IsGeneric() {
		newTypeParamTransformMap := make(types.TypeArgumentMap, len(method.TypeParameters))
		newTypeParams := make([]*types.TypeParameter, len(method.TypeParameters))
		for i, param := range method.TypeParameters {
			newParam := param.Copy()
			newTypeParams[i] = newParam
			newTypeParamTransformMap[param.Name] = types.NewTypeArgument(newParam, param.Variance)
		}
		newParams := make([]*types.Parameter, len(method.Params))
		for i, param := range method.Params {
			newParam := param.Copy()
			newParam.Type = c.replaceTypeParameters(newParam.Type, newTypeParamTransformMap, true)
			newParams[i] = newParam
		}

		copy := method.Copy()
		copy.TypeParameters = newTypeParams
		copy.Params = newParams
		copy.ReturnType = c.replaceTypeParameters(copy.ReturnType, newTypeParamTransformMap, true)
		copy.ThrowType = c.replaceTypeParameters(copy.ThrowType, newTypeParamTransformMap, true)

		return copy
	}

	newParams := make([]*types.Parameter, len(method.Params))
	for i, param := range method.Params {
		newParams[i] = param.Copy()
	}

	copy := method.Copy()
	copy.Params = newParams
	return copy
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
		nil,
		methodScope.container,
		docComment,
		false,
		false,
		false,
		false,
		false,
		value.ToSymbol(setterName),
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
	c.registerMethodCheck(
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
			types.InspectWithColor(baseMethod.DefinedUnder),
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
	if !baseMethod.IsAbstract() && overrideMethod.IsAbstract() {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - has a different modifier, is `%s`, should be `%s`",
			types.InspectModifier(overrideMethod.IsAbstract(), overrideMethod.IsSealed(), false, false),
			types.InspectModifier(baseMethod.IsAbstract(), baseMethod.IsSealed(), false, false),
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
			overrideTypeParam := overrideMethod.TypeParameters[i]
			baseTypeParam := baseMethod.TypeParameters[i]

			var isInvalid bool
			if overrideTypeParam.Name != baseTypeParam.Name || overrideTypeParam.Variance != baseTypeParam.Variance {
				isInvalid = true
			}

			switch baseTypeParam.Variance {
			case types.INVARIANT:
				if !c.isTheSameType(overrideTypeParam.UpperBound, baseTypeParam.UpperBound, nil) ||
					!c.isTheSameType(overrideTypeParam.LowerBound, baseTypeParam.LowerBound, nil) {
					isInvalid = true
				}
			case types.COVARIANT:
				if !c.isSubtype(overrideTypeParam.UpperBound, baseTypeParam.UpperBound, nil) ||
					!c.isSubtype(baseTypeParam.LowerBound, overrideTypeParam.LowerBound, nil) {
					isInvalid = true
				}
			case types.CONTRAVARIANT:
				if !c.isSubtype(baseTypeParam.UpperBound, overrideTypeParam.UpperBound, nil) ||
					!c.isSubtype(overrideTypeParam.LowerBound, baseTypeParam.LowerBound, nil) {
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

	if !c.isSubtype(overrideMethod.ReturnType, baseMethod.ReturnType, nil) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - has a different return type, is `%s`, should be `%s`",
			types.InspectWithColor(overrideMethod.ReturnType),
			types.InspectWithColor(baseMethod.ReturnType),
		)
		areIncompatible = true
	}
	if !c.isSubtype(overrideMethod.ThrowType, baseMethod.ThrowType, nil) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - has different throw type, is `%s`, should be `%s`",
			types.InspectWithColor(overrideMethod.ThrowType),
			types.InspectWithColor(baseMethod.ThrowType),
		)
		areIncompatible = true
	}

	if len(baseMethod.Params) > len(overrideMethod.Params) {
		errDetailsBuff.WriteString("\n  - has less parameters")
	} else {
		for i := range len(baseMethod.Params) {
			oldParam := baseMethod.Params[i]
			newParam := overrideMethod.Params[i]
			if oldParam.Name != newParam.Name || oldParam.Kind != newParam.Kind || !c.isSubtype(oldParam.Type, newParam.Type, nil) {
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
			param := overrideMethod.Params[i]
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
	}

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

	name := checkedMethod.Name
	prevMode := c.mode
	prevFlags := c.flags
	isClosure := types.IsClosure(methodNamespace)

	if methodNamespace != nil {
		currentMethod := c.resolveMethodInNamespace(methodNamespace, name)
		if checkedMethod != currentMethod && checkedMethod.IsSealed() {
			c.addOverrideSealedMethodError(checkedMethod, currentMethod.Location())
		}

		parent := methodNamespace.Parent()

		if parent != nil {
			baseMethod := c.resolveMethodInNamespace(parent, name)
			if baseMethod != nil {
				c.checkMethodOverride(
					checkedMethod,
					baseMethod,
					location,
				)
			}
		}
	}

	if isClosure {
		c.pushNestedLocalEnv()
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
				c.registerInitialisedInstanceVariable(value.ToSymbol(pName))
			}
			declaredType = c.TypeOf(p).(*types.Parameter).Type
			if p.TypeNode != nil {
				declaredTypeNode = p.TypeNode
				switch p.Kind {
				case ast.PositionalRestParameterKind:
					declaredType = types.NewGenericWithTypeArgs(c.StdTuple(), declaredType)
				case ast.NamedRestParameterKind:
					declaredType = types.NewGenericWithTypeArgs(c.StdRecord(), c.Std(symbol.Symbol), declaredType)
				}
			}
			var initNode ast.ExpressionNode
			if p.Initialiser != nil {
				initNode = c.checkExpression(p.Initialiser)
				initType := c.TypeOf(initNode)
				c.checkCanAssign(initType, declaredType, initNode.Location())
			}
			c.addLocal(pName, newLocal(declaredType, true, checkedMethod.IsGenerator()))
			p.Initialiser = initNode
			p.TypeNode = declaredTypeNode
		case *ast.FormalParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
			pName := c.identifierToName(p.Name)
			declaredType = c.TypeOf(p).(*types.Parameter).Type
			if p.TypeNode != nil {
				declaredTypeNode = p.TypeNode
				switch p.Kind {
				case ast.PositionalRestParameterKind:
					declaredType = types.NewGenericWithTypeArgs(c.StdTuple(), declaredType)
				case ast.NamedRestParameterKind:
					declaredType = types.NewGenericWithTypeArgs(c.StdRecord(), c.Std(symbol.Symbol), declaredType)
				}
			}
			var initNode ast.ExpressionNode
			if p.Initialiser != nil {
				initNode = c.checkExpression(p.Initialiser)
				initType := c.TypeOf(initNode)
				c.checkCanAssign(initType, declaredType, initNode.Location())
			}
			c.addLocal(pName, newLocal(declaredType, true, false))
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
		returnType = origReturnType.(*types.Generic).Get(0).Type
	}

	throwType := checkedMethod.ThrowType
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.TypeOf(typedThrowTypeNode)
	}
	if checkedMethod.IsGenerator() || checkedMethod.IsAsync() {
		throwType = origReturnType.(*types.Generic).Get(1).Type
	}
	if !types.IsNever(throwType) && throwType != nil {
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
			if returnType == nil {
				c.setInferClosureReturnType(true)
			}
			if throwType == nil {
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
				c.checkCanAssign(bodyReturnType, returnType, returnSpan)
			}

			if c.shouldInferClosureThrowType() {
				if c.throwType == nil {
					checkedMethod.ThrowType = types.Never{}
				} else {
					checkedMethod.ThrowType = c.throwType
				}
			}
		}
	}

	c.returnType = nil
	c.throwType = nil
	c.mode = prevMode
	c.flags = prevFlags
	c.catchScopes = prevCatchScopes
	return typedReturnTypeNode, typedThrowTypeNode
}

func (c *Checker) checkSpecialMethods(name value.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, location *position.Location) {
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

func (c *Checker) checkEqualityOperator(name value.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, location *position.Location) {
	params := checkedMethod.Params

	if !c.isTheSameType(checkedMethod.ReturnType, types.Bool{}, nil) {
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

	param := params[0]
	var paramSpan *position.Location
	if paramNodes != nil {
		paramSpan = paramNodes[0].Location()
	} else {
		paramSpan = location
	}
	if !types.IsAny(param.Type) {
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

func (c *Checker) checkRelationalOperator(name value.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, location *position.Location) {
	params := checkedMethod.Params

	if !c.isTheSameType(checkedMethod.ReturnType, types.Bool{}, nil) {
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

	param := checkedMethod.Params[0]
	var paramSpan *position.Location
	if paramNodes != nil {
		paramSpan = paramNodes[0].Location()
	} else {
		paramSpan = location
	}
	if !checkedMethod.IsAbstract() && !c.isSubtype(c.selfType, param.Type, nil) {
		c.addFailure(
			fmt.Sprintf(
				"parameter `%s` of relational operator `%s` must accept `%s`",
				lexer.Colorize(param.Name.String()),
				lexer.Colorize(name.String()),
				types.InspectWithColor(c.selfType),
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

func (c *Checker) checkFixedParameterCountMethod(name value.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, desiredParamCount int, location *position.Location) {
	params := checkedMethod.Params

	if types.IsVoid(checkedMethod.ReturnType) {
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

	for i, param := range params {
		var paramSpan *position.Location
		if paramNodes != nil {
			paramSpan = paramNodes[i].Location()
		} else {
			paramSpan = location
		}

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
	if c.returnType == nil {
		c.returnType = typ
		return
	}

	c.returnType = c.NewNormalisedUnion(c.returnType, typ)
}

func (c *Checker) addToThrowType(typ types.Type) {
	if c.throwType == nil {
		c.throwType = typ
		return
	}

	c.throwType = c.NewNormalisedUnion(c.throwType, typ)
}

func (c *Checker) checkMethodArgumentsAndInferTypeArguments(
	method *types.Method,
	positionalArguments []ast.ExpressionNode,
	namedArguments []ast.NamedArgumentNode,
	typeParams []*types.TypeParameter,
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

	var currentParamIndex int
	// check all positional argument before the rest parameter
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
		param := method.Params[currentParamIndex]

		typedPosArg := c.checkExpressionWithType(posArg, param.Type)
		posArgType := c.TypeOf(typedPosArg)

		inferredParamType := c.inferTypeArguments(posArgType, param.Type, typeArgMap, typedPosArg.Location())
		if inferredParamType == nil {
			param.Type = types.Untyped{}
		} else if inferredParamType != param.Type {
			param.Type = inferredParamType
		}
		typedPositionalArguments = append(typedPositionalArguments, typedPosArg)

		if !c.isSubtype(posArgType, param.Type, posArg.Location()) {
			c.addFailure(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
					types.InspectWithColor(param.Type),
					param.Name.String(),
					types.InspectWithColor(method),
					types.InspectWithColor(posArgType),
				),
				posArg.Location(),
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
		posRestParam := method.Params[positionalRestParamIndex]

		currentArgIndex := currentParamIndex
		// check rest arguments
		for ; currentArgIndex < min(argCount-method.PostParamCount, len(positionalArguments)); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			typedPosArg := c.checkRestArgument(posArg, posRestParam.Type)
			posArgType := c.TypeOf(typedPosArg)
			inferredParamType := c.inferTypeArguments(posArgType, posRestParam.Type, typeArgMap, typedPosArg.Location())
			if inferredParamType == nil {
				posRestParam.Type = types.Untyped{}
			} else if inferredParamType != posRestParam.Type {
				posRestParam.Type = inferredParamType
			}
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, typedPosArg)
			if !c.isSubtype(posArgType, posRestParam.Type, posArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(posRestParam.Type),
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
				if c.IsSubtype(inType, c.Std(symbol.Tuple)) {
					typedPositionalArguments[len(typedPositionalArguments)-1] = forIn.InExpression
				}
			}
		}

		currentParamIndex = positionalRestParamIndex
		// check post arguments
		for ; currentArgIndex < len(positionalArguments); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			currentParamIndex++
			param := method.Params[currentParamIndex]

			typedPosArg := c.checkExpressionWithType(posArg, param.Type)
			posArgType := c.TypeOf(typedPosArg)
			inferredParamType := c.inferTypeArguments(posArgType, param.Type, typeArgMap, typedPosArg.Location())
			if inferredParamType == nil {
				param.Type = types.Untyped{}
			} else if inferredParamType != param.Type {
				param.Type = inferredParamType
			}
			typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
			if !c.isSubtype(posArgType, param.Type, posArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
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
		param := method.Params[i]
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

			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, param.Type)
			namedArgType := c.TypeOf(typedNamedArgValue)
			inferredParamType := c.inferTypeArguments(namedArgType, param.Type, typeArgMap, typedNamedArgValue.Location())
			if inferredParamType == nil {
				param.Type = types.Untyped{}
			} else if inferredParamType != param.Type {
				param.Type = inferredParamType
			}
			typedPositionalArguments = append(typedPositionalArguments, typedNamedArgValue)
			if !c.isSubtype(namedArgType, param.Type, namedArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
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
		namedRestParam := method.Params[len(method.Params)-1]
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			switch namedArg := namedArgI.(type) {
			case *ast.NamedCallArgumentNode:
				typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, namedRestParam.Type)
				posArgType := c.TypeOf(typedNamedArgValue)
				inferredParamType := c.inferTypeArguments(posArgType, namedRestParam.Type, typeArgMap, typedNamedArgValue.Location())
				if inferredParamType == nil {
					namedRestParam.Type = types.Untyped{}
				} else if inferredParamType != namedRestParam.Type {
					namedRestParam.Type = inferredParamType
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
				if c.IsSubtype(inType, c.Std(symbol.Record)) {
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
		for _, typeParam := range typeParams {
			typeArg := typeArgMap[typeParam.Name]
			if typeArg != nil {
				continue
			}

			var inferredType types.Type
			if !types.IsNever(typeParam.LowerBound) && !c.containsTypeParameters(typeParam.LowerBound) {
				inferredType = typeParam.LowerBound
			} else if !c.containsTypeParameters(typeParam.UpperBound) {
				inferredType = typeParam.UpperBound
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
	if !c.isSubtype(keyType, c.Std(symbol.Symbol), node.Location()) {
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
	if c.isSubtype(argType, param.Type, location) {
		return
	}

	c.addFailure(
		fmt.Sprintf(
			"expected type `%s` for named rest parameter `**%s` in call to `%s`, got type `%s`",
			types.InspectWithColor(param.Type),
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
) []ast.ExpressionNode {
	posArgs, _ := c.checkMethodArgumentsAndInferTypeArguments(method, positionalArguments, namedArguments, nil, location)
	return posArgs
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
		typedPositionalArguments = c.checkNonGenericMethodArguments(
			method,
			positionalArgumentNodes,
			namedArgumentNodes,
			location,
		)
		return method, typedPositionalArguments
	}

	if len(method.TypeParameters) > 0 {
		var typeArgMap types.TypeArgumentMap
		method = c.deepCopyMethod(method)
		typedPositionalArguments, typeArgMap = c.checkMethodArgumentsAndInferTypeArguments(
			method,
			positionalArgumentNodes,
			namedArgumentNodes,
			method.TypeParameters,
			location,
		)
		if len(typeArgMap) != len(method.TypeParameters) {
			return nil, nil
		}
		method.ReturnType = c.replaceTypeParameters(method.ReturnType, typeArgMap, true)
		method.ThrowType = c.replaceTypeParameters(method.ThrowType, typeArgMap, true)
		return method, typedPositionalArguments
	}

	typedPositionalArguments = c.checkNonGenericMethodArguments(
		method,
		positionalArgumentNodes,
		namedArgumentNodes,
		location,
	)
	return method, typedPositionalArguments
}

func (c *Checker) checkSimpleMethodCall(
	receiver ast.ExpressionNode,
	op token.Type,
	methodName value.Symbol,
	typeArgumentNodes []ast.TypeNode,
	positionalArgumentNodes []ast.ExpressionNode,
	namedArgumentNodes []ast.NamedArgumentNode,
	location *position.Location,
) (
	_receiver ast.ExpressionNode,
	_positionalArguments []ast.ExpressionNode,
	typ types.Type,
) {
	receiver = c.checkExpression(receiver)
	receiverType := c.TypeOf(receiver)

	// Allow arbitrary method calls on `never` and `nothing`.
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

		return receiver, typedPositionalArguments, receiverType
	}

	var method *types.Method
	switch op {
	case token.DOT, token.DOT_DOT:
		method = c.getMethod(receiverType, methodName, location)
	case token.QUESTION_DOT, token.QUESTION_DOT_DOT:
		nonNilableReceiverType := c.ToNonNilable(receiverType)
		method = c.getMethod(nonNilableReceiverType, methodName, location)
	default:
		panic(fmt.Sprintf("invalid call operator: %#v", op))
	}
	if method == nil {
		c.checkExpressions(positionalArgumentNodes)
		c.checkNamedArguments(namedArgumentNodes)
		return receiver, positionalArgumentNodes, types.Untyped{}
	}

	c.addToMethodCache(method)

	method, typedPositionalArguments := c.checkMethodArguments(method, typeArgumentNodes, positionalArgumentNodes, namedArgumentNodes, location)
	if method == nil {
		return receiver, positionalArgumentNodes, types.Untyped{}
	}

	var returnType types.Type
	switch op {
	case token.DOT:
		returnType = method.ReturnType
	case token.QUESTION_DOT:
		if !c.IsNilable(receiverType) {
			c.addFailure(
				fmt.Sprintf("cannot make a nil-safe call on type `%s` which is not nilable", types.InspectWithColor(receiverType)),
				location,
			)
			returnType = method.ReturnType
		} else {
			returnType = c.ToNilable(method.ReturnType)
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

	c.checkCalledMethodThrowType(method, location)

	return receiver, typedPositionalArguments, returnType
}

func (c *Checker) checkBinaryOpMethodCall(
	left ast.ExpressionNode,
	right ast.ExpressionNode,
	methodName value.Symbol,
	location *position.Location,
) types.Type {
	_, _, returnType := c.checkSimpleMethodCall(
		left,
		token.DOT,
		methodName,
		nil,
		[]ast.ExpressionNode{right},
		nil,
		location,
	)

	return returnType
}

func (c *Checker) checkMethodDefinition(node *ast.MethodDefinitionNode, method *types.Method) {
	c.method = method
	returnType, throwType := c.checkMethod(
		c.currentMethodScope().container,
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
		fmt.Printf("header: %t\n", c.IsHeader())
		method.Body = c.compiler.CompileMethodBody(node, method.Name)
	}
}

func (c *Checker) declareMethod(
	baseMethod *types.Method,
	methodNamespace types.Namespace,
	docComment string,
	abstract bool,
	sealed bool,
	inferReturnType bool,
	generator bool,
	async bool,
	name value.Symbol,
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
	oldMethod := methodNamespace.Method(name)
	if oldMethod != nil {
		if sealed && !oldMethod.IsSealed() {
			c.addFailure(
				fmt.Sprintf(
					"cannot redeclare method `%s` with a different modifier, is `%s`, should be `%s`",
					name.String(),
					types.InspectModifier(abstract, sealed, false, false),
					types.InspectModifier(oldMethod.IsAbstract(), oldMethod.IsSealed(), false, false),
				),
				location,
			)
		}
	}

	switch namespace := methodNamespace.(type) {
	case *types.Interface:
	case *types.Class:
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

	var typeParams []*types.TypeParameter
	var typeParamMod *types.TypeParamNamespace
	if len(typeParamNodes) > 0 {
		typeParams = make([]*types.TypeParameter, 0, len(typeParamNodes))
		typeParamMod = types.NewTypeParamNamespace(fmt.Sprintf("Type Parameter Container of %s", name), true)
		c.pushConstScope(makeConstantScope(typeParamMod))
		for _, typeParamNode := range typeParamNodes {
			node, ok := typeParamNode.(*ast.VariantTypeParameterNode)
			if !ok {
				continue
			}

			t := c.checkTypeParameterNode(node, typeParamMod, false)
			typeParams = append(typeParams, t)
			typeParamNode.SetType(t)
			typeParamMod.DefineSubtype(t.Name, t)
			typeParamMod.DefineConstant(t.Name, types.NoValue{})
		}
	}

	if name != symbol.S_init {
		c.mode = prevMode
		c.setInputPositionTypeMode()
	}
	var params []*types.Parameter
	for i, paramNode := range paramNodes {
		switch p := paramNode.(type) {
		case *ast.FormalParameterNode:
			pName := c.identifierToName(p.Name)
			var declaredType types.Type
			if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.TypeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Type
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
			name := value.ToSymbol(pName)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType)
		case *ast.MethodParameterNode:
			pName := c.identifierToName(p.Name)
			var declaredType types.Type
			if p.SetInstanceVariable {
				currentIvar, _ := c.getInstanceVariableIn(value.ToSymbol(pName), methodNamespace)
				if p.TypeNode == nil {
					if currentIvar == nil {
						c.addFailure(
							fmt.Sprintf(
								"cannot infer the type of instance variable `%s`",
								pName,
							),
							p.Location(),
						)
					}

					declaredType = currentIvar
				} else {
					p.TypeNode = c.checkTypeNode(p.TypeNode)
					declaredType = c.TypeOf(p.TypeNode)
					if currentIvar != nil {
						c.checkCanAssignInstanceVariable(pName, declaredType, currentIvar, p.TypeNode.Location())
					} else {
						c.declareInstanceVariable(value.ToSymbol(pName), declaredType, p.Location())
					}
				}
			} else if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.TypeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Type
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
			name := value.ToSymbol(pName)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType)
		case *ast.SignatureParameterNode:
			pName := c.identifierToName(p.Name)
			var declaredType types.Type
			if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.TypeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Type
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
			name := value.ToSymbol(pName)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType)
		default:
			c.addFailure(
				fmt.Sprintf("invalid param type %T", paramNode),
				paramNode.Location(),
			)
		}
	}
	if async {
		paramType := types.NewParameter(
			value.ToSymbol("_pool"),
			c.Std(symbol.ThreadPool),
			types.DefaultValueParameterKind,
			false,
		)
		params = append(params, paramType)
	}

	c.mode = prevMode
	c.setOutputPositionTypeMode()

	var returnType types.Type
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.TypeOf(typedReturnTypeNode)
	} else if inferReturnType {
	} else if baseMethod != nil && baseMethod.ReturnType != nil {
		returnType = baseMethod.ReturnType
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.TypeOf(typedThrowTypeNode)
	} else if inferReturnType {
	} else if baseMethod != nil && baseMethod.ThrowType != nil {
		throwType = baseMethod.ThrowType
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
			c.env.StdSubtypeClass(symbol.Generator),
			returnType,
			throwType,
		)

		throwType = types.Never{}
	} else if async {
		returnType = types.NewGenericWithTypeArgs(
			c.env.StdSubtypeClass(symbol.Promise),
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

	c.checkMethodOverrideWithPlaceholder(newMethod, oldMethod, location)
	methodNamespace.SetMethod(name, newMethod)

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
		baseMethod.DefinedUnder.SetMethod(baseMethod.Name, overrideMethod)
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

func (c *Checker) checkMethodCompatibilityForAlgebraicTypes(baseMethod, overrideMethod *types.Method, errSpan *position.Location) bool {
	if !overrideMethod.IsGeneric() {
		return c.checkMethodCompatibility(baseMethod, overrideMethod, errSpan, true)
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
func (c *Checker) checkMethodCompatibility(baseMethod, overrideMethod *types.Method, errSpan *position.Location, validateParamNames bool) bool {
	if baseMethod == nil {
		return true
	}

	areCompatible := true
	errDetailsBuff := new(strings.Builder)

	if !c.isSubtype(overrideMethod.ReturnType, baseMethod.ReturnType, errSpan) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` has a different return type than `%s`, has `%s`, should have `%s`",
			types.InspectWithColor(overrideMethod),
			types.InspectWithColor(baseMethod),
			types.InspectWithColor(overrideMethod.ReturnType),
			types.InspectWithColor(baseMethod.ReturnType),
		)
		areCompatible = false
	}
	if !c.isSubtype(overrideMethod.ThrowType, baseMethod.ThrowType, errSpan) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` has a different throw type than `%s`, has `%s`, should have `%s`",
			types.InspectWithColor(overrideMethod),
			types.InspectWithColor(baseMethod),
			types.InspectWithColor(overrideMethod.ThrowType),
			types.InspectWithColor(baseMethod.ThrowType),
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
			oldParam := baseMethod.Params[i]
			newParam := overrideMethod.Params[i]

			if (validateParamNames && oldParam.Name != newParam.Name) || oldParam.Kind != newParam.Kind || !c.isSubtype(oldParam.Type, newParam.Type, errSpan) {
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
			param := overrideMethod.Params[i]
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

	returnType := c.inferTypeArguments(baseMethod.ReturnType, overrideMethod.ReturnType, typeArgs, nil)
	if returnType == nil || !c.isSubtype(returnType, baseMethod.ReturnType, errSpan) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` has a different return type than `%s`, has `%s`, should have `%s`",
			types.InspectWithColor(overrideMethod),
			types.InspectWithColor(baseMethod),
			types.InspectWithColor(overrideMethod.ReturnType),
			types.InspectWithColor(baseMethod.ReturnType),
		)
		areCompatible = false
	}

	throwType := c.inferTypeArguments(baseMethod.ThrowType, overrideMethod.ThrowType, typeArgs, nil)
	if throwType == nil || !c.isSubtype(throwType, baseMethod.ThrowType, errSpan) {
		fmt.Fprintf(
			errDetailsBuff,
			"\n  - method `%s` has a different throw type than `%s`, has `%s`, should have `%s`",
			types.InspectWithColor(overrideMethod),
			types.InspectWithColor(baseMethod),
			types.InspectWithColor(overrideMethod.ThrowType),
			types.InspectWithColor(baseMethod.ThrowType),
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
			oldParam := baseMethod.Params[i]
			newParam := overrideMethod.Params[i]

			newParamType := c.inferTypeArguments(oldParam.Type, newParam.Type, typeArgs, nil)
			if oldParam.Name != newParam.Name || oldParam.Kind != newParam.Kind ||
				newParamType == nil || !c.isSubtype(oldParam.Type, newParamType, errSpan) {
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
			param := overrideMethod.Params[i]
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

func (c *Checker) getMethod(typ types.Type, name value.Symbol, errSpan *position.Location) *types.Method {
	return c._getMethod(typ, name, errSpan, false, false)
}

// Iterates over every method of the namespace, resolving type parameters.
func (c *Checker) methodsInNamespace(namespace types.Namespace) iter.Seq2[value.Symbol, *types.Method] {
	return func(yield func(name value.Symbol, method *types.Method) bool) {
		var generics []*types.Generic
		seenMethods := make(ds.Set[value.Symbol])

		for parent := range types.Parents(namespace) {
			if generic, ok := parent.(*types.Generic); ok {
				generics = append(generics, generic)
			}
			methods := parent.Methods()
			names := symbol.SortKeys(methods)
		methodLoop:
			for _, name := range names {
				method := methods[name]
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

				var whereParams []*types.TypeParameter
				var whereArgs []types.Type
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
					whereParam := whereParams[i]
					whereArg := whereArgs[i]

					if !c.isSubtype(whereParam.LowerBound, whereArg, nil) {
						continue methodLoop
					}
					if !c.isSubtype(whereArg, whereParam.UpperBound, nil) {
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
func (c *Checker) abstractMethodsInNamespace(namespace types.Namespace) iter.Seq2[value.Symbol, *types.Method] {
	return func(yield func(name value.Symbol, method *types.Method) bool) {
		var generics []*types.Generic
		seenMethods := make(ds.Set[value.Symbol])

		for parent := range types.Parents(namespace) {
			if generic, ok := parent.(*types.Generic); ok {
				generics = append(generics, generic)
			}
			if !parent.IsAbstract() {
				continue
			}
			for name, method := range parent.Methods() {
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

func (c *Checker) resolveMethodInNamespace(namespace types.Namespace, name value.Symbol) *types.Method {
	var generics []*types.Generic

	for parent := range types.Parents(namespace) {
		switch p := parent.(type) {
		case *types.Generic:
			generics = append(generics, p)
		case *types.NamespacePlaceholder:
			switch n := p.Namespace.(type) {
			case *types.Module:
				parent = n
			default:
				parent = n.Singleton()
			}
		}

		method := parent.Method(name)
		if method == nil {
			continue
		}

		var whereParams []*types.TypeParameter
		var whereArgs []types.Type
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
			whereParam := whereParams[i]
			whereArg := whereArgs[i]

			if !c.isSubtype(whereParam.LowerBound, whereArg, nil) {
				return nil
			}
			if !c.isSubtype(whereArg, whereParam.UpperBound, nil) {
				return nil
			}
		}

		return method
	}

	return nil
}

func (c *Checker) constructWhereArguments(whereParameters []*types.TypeParameter) []types.Type {
	whereArgs := make([]types.Type, len(whereParameters))
	for i, whereParam := range whereParameters {
		whereArgs[i] = whereParam
	}

	return whereArgs
}

func (c *Checker) resolveNonAbstractMethodInNamespace(namespace types.Namespace, name value.Symbol) *types.Method {
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

func (c *Checker) _getMethodInNamespace(namespace types.Namespace, typ types.Type, name value.Symbol, errSpan *position.Location, inParent bool) *types.Method {
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

func (c *Checker) getMethodInNamespaceWithSelf(namespace types.Namespace, typ types.Type, name value.Symbol, self types.Type, errSpan *position.Location, inParent, inSelf bool) *types.Method {
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

func (c *Checker) getMethodInNamespace(namespace types.Namespace, typ types.Type, name value.Symbol, errSpan *position.Location, inParent, inSelf bool) *types.Method {
	return c.getMethodInNamespaceWithSelf(namespace, typ, name, namespace, errSpan, inParent, inSelf)
}

func (c *Checker) replaceTypeParametersInMethodCopy(method *types.Method, typeArgs types.TypeArgumentMap, replaceMethodTypeParams bool) *types.Method {
	var methodCopy *types.Method

	for i, typeParam := range method.TypeParameters {
		result := c.replaceTypeParameters(typeParam.LowerBound, typeArgs, replaceMethodTypeParams)
		if typeParam.LowerBound != result {
			if methodCopy == nil {
				methodCopy = c.deepCopyMethod(method)
			}
			methodCopy.TypeParameters[i].LowerBound = result
		}
		result = c.replaceTypeParameters(typeParam.UpperBound, typeArgs, replaceMethodTypeParams)
		if typeParam.UpperBound != result {
			if methodCopy == nil {
				methodCopy = c.deepCopyMethod(method)
			}
			methodCopy.TypeParameters[i].UpperBound = result
		}
	}
	result := c.replaceTypeParameters(method.ReturnType, typeArgs, replaceMethodTypeParams)
	if method.ReturnType != result {
		if methodCopy == nil {
			methodCopy = c.deepCopyMethod(method)
		}
		methodCopy.ReturnType = result
	}
	result = c.replaceTypeParameters(method.ThrowType, typeArgs, replaceMethodTypeParams)
	if method.ThrowType != result {
		if methodCopy == nil {
			methodCopy = c.deepCopyMethod(method)
		}
		methodCopy.ThrowType = result
	}

	for i, param := range method.Params {
		result := c.replaceTypeParameters(param.Type, typeArgs, replaceMethodTypeParams)
		if param.Type != result {
			if methodCopy == nil {
				methodCopy = c.deepCopyMethod(method)
			}
			methodCopy.Params[i].Type = result
		}
	}

	if methodCopy != nil {
		return methodCopy
	}
	return method
}

func (c *Checker) replaceTypeParametersInMethod(method *types.Method, typeArgs types.TypeArgumentMap, replaceMethodTypeParams bool) *types.Method {
	for _, typeParam := range method.TypeParameters {
		typeParam.LowerBound = c.replaceTypeParameters(typeParam.LowerBound, typeArgs, replaceMethodTypeParams)
		typeParam.UpperBound = c.replaceTypeParameters(typeParam.UpperBound, typeArgs, replaceMethodTypeParams)
	}
	method.ReturnType = c.replaceTypeParameters(method.ReturnType, typeArgs, replaceMethodTypeParams)
	method.ThrowType = c.replaceTypeParameters(method.ThrowType, typeArgs, replaceMethodTypeParams)

	for _, param := range method.Params {
		param.Type = c.replaceTypeParameters(param.Type, typeArgs, replaceMethodTypeParams)
	}

	return method
}

func (c *Checker) replaceTypeParametersInWhere(whereParams []*types.TypeParameter, whereArgs []types.Type, typeArgs types.TypeArgumentMap) {
	for i, whereArg := range whereArgs {
		whereArgs[i] = c.replaceTypeParameters(whereArg, typeArgs, false)
	}

	for i, whereParam := range whereParams {
		var whereParamCopy *types.TypeParameter

		result := c.replaceTypeParameters(whereParam.LowerBound, typeArgs, false)
		if result != whereParam.LowerBound {
			whereParamCopy = whereParam.Copy()
			whereParams[i] = whereParamCopy
			whereParamCopy.LowerBound = result
		}

		result = c.replaceTypeParameters(whereParam.UpperBound, typeArgs, false)
		if result != whereParam.LowerBound {
			if whereParamCopy == nil {
				whereParamCopy = whereParam.Copy()
				whereParams[i] = whereParamCopy
			}
			whereParamCopy.UpperBound = result
		}
	}
}

func (c *Checker) getMethodForTypeParameter(typ *types.TypeParameter, name value.Symbol, errSpan *position.Location, inParent, inSelf bool) *types.Method {
	switch upper := typ.UpperBound.(type) {
	case *types.Class:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.Mixin:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.Interface:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.Closure:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.SingletonClass:
		return c.getMethodInNamespaceWithSelf(upper, typ, name, typ, errSpan, inParent, inSelf)
	case *types.Generic:
		var method *types.Method
		switch genericType := upper.Namespace.(type) {
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
		return c._getMethod(typ.UpperBound, name, errSpan, inParent, inSelf)
	}
}

func (c *Checker) getReceiverlessMethod(name value.Symbol, location *position.Location) (_ *types.Method, fromLocal bool) {
	nameStr := name.String()
	local, _ := c.resolveLocal(nameStr, nil)
	if local != nil {
		if !local.initialised {
			c.addUninitialisedLocalError(nameStr, location)
		}
		return c.getMethod(local.typ, symbol.L_call, location), true
	}
	method := c.getMethod(c.selfType, name, nil)
	if method != nil {
		return method, false
	}

	for _, methodScope := range c.methodScopes {
		switch methodScope.kind {
		case scopeUsingBufferKind, scopeUsingKind:
		default:
			continue
		}

		namespace := methodScope.container
		method := c.getMethod(namespace, name, nil)
		if method != nil {
			return method, false
		}
	}

	c.addMissingMethodError(c.selfType, name.String(), location)

	return nil, false
}

func (c *Checker) _getMethod(typ types.Type, name value.Symbol, errSpan *position.Location, inParent, inSelf bool) *types.Method {
	typ = c.ToNonLiteral(typ, true)

	switch t := typ.(type) {
	case types.Self:
		return c._getMethod(c.selfType, name, errSpan, inParent, true)
	case *types.NamedType:
		return c._getMethod(t.Type, name, errSpan, inParent, inSelf)
	case *types.TypeParameter:
		return c.getMethodForTypeParameter(t, name, errSpan, inParent, inSelf)
	case *types.Generic:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.Class:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.NamespacePlaceholder:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.SingletonClass:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.Interface:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.Closure:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.InterfaceProxy:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.Module:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.Mixin:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.MixinProxy:
		return c.getMethodInNamespace(t, typ, name, errSpan, inParent, inSelf)
	case *types.Intersection:
		var methods []*types.Method
		var baseMethod *types.Method

		for _, element := range t.Elements {
			switch e := element.(type) {
			case *types.Not:
				switch t := e.Type.(type) {
				case *types.Interface:
					elementMethod := c.getMethod(t, name, nil)
					if elementMethod == nil {
						continue
					}
					return nil
				case *types.Mixin:
					elementMethod := c.getMethod(t, name, nil)
					if elementMethod == nil {
						continue
					}
					return nil
				}
			default:
				elementMethod := c.getMethod(element, name, nil)
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
			c.addMissingMethodError(typ, name.String(), errSpan)
			return nil
		case 1:
			return methods[0]
		}

		isCompatible := true
		for i := range len(methods) {
			method := methods[i]

			if !c.checkMethodCompatibilityForAlgebraicTypes(baseMethod, method, errSpan) {
				isCompatible = false
			}
		}

		if isCompatible {
			return baseMethod
		}

		return nil
	case *types.Nilable:
		nilType := c.env.StdSubtype(symbol.Nil).(*types.Class)
		nilMethod := nilType.Method(name)
		if nilMethod == nil {
			c.addMissingMethodError(nilType, name.String(), errSpan)
		}
		nonNilMethod := c.getMethod(t.Type, name, errSpan)
		if nilMethod == nil || nonNilMethod == nil {
			return nil
		}

		var baseMethod *types.Method
		var overrideMethod *types.Method
		if len(nilMethod.Params) < len(nonNilMethod.Params) || nilMethod.IsGeneric() && !nonNilMethod.IsGeneric() {
			baseMethod = nilMethod
			overrideMethod = nonNilMethod
		} else {
			baseMethod = nonNilMethod
			overrideMethod = nilMethod
		}

		if c.checkMethodCompatibilityForAlgebraicTypes(baseMethod, overrideMethod, errSpan) {
			return baseMethod
		}
		return nil
	case *types.Union:
		var methods []*types.Method
		var baseMethod *types.Method

		for _, element := range t.Elements {
			elementMethod := c.getMethod(element, name, errSpan)
			if elementMethod == nil {
				continue
			}
			methods = append(methods, elementMethod)
			if baseMethod == nil || len(baseMethod.Params) > len(elementMethod.Params) || baseMethod.IsGeneric() && !elementMethod.IsGeneric() {
				baseMethod = elementMethod
			}
		}

		if len(methods) < len(t.Elements) {
			return nil
		}

		isCompatible := true
		for i := range len(methods) {
			method := methods[i]

			if !c.checkMethodCompatibilityForAlgebraicTypes(baseMethod, method, errSpan) {
				isCompatible = false
			}
		}

		if isCompatible {
			return baseMethod
		}

		return nil
	default:
		c.addMissingMethodError(typ, name.String(), errSpan)
		return nil
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
