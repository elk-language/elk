package checker

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

func (c *Checker) newMethodChecker(
	filename string,
	constScopes []constantScope,
	methodScopes []methodScope,
	selfType,
	returnType,
	throwType types.Type,
	isInit bool,
) *Checker {
	checker := &Checker{
		GlobalEnv:      c.GlobalEnv,
		Filename:       filename,
		mode:           methodMode,
		phase:          methodCheckPhase,
		selfType:       selfType,
		returnType:     returnType,
		throwType:      throwType,
		constantScopes: constScopes,
		methodScopes:   methodScopes,
		Errors:         c.Errors,
		IsHeader:       c.IsHeader,
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
		typeDefinitionChecks: newTypeDefinitionChecks(),
		methodCache:          concurrent.NewSlice[*types.Method](),
		compiler:             c.compiler,
	}
	if isInit {
		checker.mode = initMode
	}
	return checker
}

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
}

func (c *Checker) registerMethodCheck(method *types.Method, node *ast.MethodDefinitionNode) {
	c.methodChecks = append(c.methodChecks, methodCheckEntry{
		method:         method,
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
		node:           node,
	})
}

var concurrencyLimit = 10_000

func (c *Checker) checkMethods() {
	concurrent.Foreach(
		concurrencyLimit,
		c.methodChecks,
		func(methodCheck methodCheckEntry) {
			method := methodCheck.method
			node := methodCheck.node
			methodChecker := c.newMethodChecker(
				node.Location().Filename,
				methodCheck.constantScopes,
				methodCheck.methodScopes,
				method.DefinedUnder,
				method.ReturnType,
				method.ThrowType,
				method.IsInit(),
			)
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
}

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
	method, mod := c.declareMethod(
		nil,
		c.currentMethodScope().container,
		docComment,
		false,
		false,
		false,
		value.ToSymbol(node.Name),
		nil,
		nil,
		node.TypeNode,
		nil,
		node.Span(),
	)
	method.SetAttribute(true)

	init := node.Initialiser
	var body []ast.StatementNode

	if init == nil {
		body = ast.ExpressionToStatements(
			ast.NewInstanceVariableNode(node.Span(), node.Name),
		)
	} else {
		body = ast.ExpressionToStatements(
			ast.NewAssignmentExpressionNode(
				node.Span(),
				token.New(init.Span(), token.QUESTION_QUESTION_EQUAL),
				ast.NewInstanceVariableNode(node.Span(), node.Name),
				init,
			),
		)
	}

	methodNode := ast.NewMethodDefinitionNode(
		c.newLocation(node.Span()),
		"",
		false,
		false,
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
			newParam.Type = c.replaceTypeParameters(newParam.Type, newTypeParamTransformMap)
			newParams[i] = newParam
		}

		copy := method.Copy()
		copy.TypeParameters = newTypeParams
		copy.Params = newParams
		copy.ReturnType = c.replaceTypeParameters(copy.ReturnType, newTypeParamTransformMap)
		copy.ThrowType = c.replaceTypeParameters(copy.ThrowType, newTypeParamTransformMap)

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
	setterName := node.Name + "="

	methodScope := c.currentMethodScope()
	var paramSpan *position.Span
	if node.TypeNode != nil {
		paramSpan = node.TypeNode.Span()
	} else {
		node.Span()
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
		value.ToSymbol(setterName),
		nil,
		params,
		nil,
		nil,
		node.Span(),
	)
	method.SetAttribute(true)

	methodNode := ast.NewMethodDefinitionNode(
		c.newLocation(node.Span()),
		docComment,
		false,
		false,
		setterName,
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

func (c *Checker) addWrongArgumentCountError(got int, method *types.Method, span *position.Span) {
	c.addFailure(
		fmt.Sprintf("expected %s arguments in call to `%s`, got %d", method.ExpectedParamCountString(), lexer.Colorize(method.Name.String()), got),
		span,
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
	span *position.Span,
) {
	var areIncompatible bool
	errDetailsBuff := new(strings.Builder)

	if !c.IsHeader && baseMethod.IsSealed() {
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
			types.InspectModifier(overrideMethod.IsAbstract(), overrideMethod.IsSealed(), false),
			types.InspectModifier(baseMethod.IsAbstract(), baseMethod.IsSealed(), false),
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
			span,
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
	span *position.Span,
) (ast.TypeNode, ast.TypeNode) {
	prevCatchScopes := c.catchScopes
	c.catchScopes = nil

	name := checkedMethod.Name
	prevMode := c.mode
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
					span,
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

	c.mode = prevMode
	c.setInputPositionTypeMode()
	for _, param := range paramNodes {
		switch p := param.(type) {
		case *ast.MethodParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
			if p.SetInstanceVariable {
				c.registerInitialisedInstanceVariable(value.ToSymbol(p.Name))
			}
			declaredType = c.typeOf(p).(*types.Parameter).Type
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
				initType := c.typeOf(initNode)
				c.checkCanAssign(initType, declaredType, initNode.Span())
			}
			c.addLocal(p.Name, newLocal(declaredType, true, false))
			p.Initialiser = initNode
			p.TypeNode = declaredTypeNode
		case *ast.FormalParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
			declaredType = c.typeOf(p).(*types.Parameter).Type
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
				initType := c.typeOf(initNode)
				c.checkCanAssign(initType, declaredType, initNode.Span())
			}
			c.addLocal(p.Name, newLocal(declaredType, true, false))
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

	throwType := checkedMethod.ThrowType
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
	}
	c.pushCatchScope(makeCatchScope(throwType))

	if len(body) > 0 && checkedMethod.IsAbstract() {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` cannot have a body because it is abstract",
				name.String(),
			),
			span,
		)
	}

	if !c.IsHeader {
		if isClosure && returnType == nil {
			c.mode = closureInferReturnTypeMode
		} else if checkedMethod.IsInit() {
			c.mode = initMode
		} else {
			c.mode = methodMode
		}
		c.returnType = returnType
		c.throwType = throwType
		bodyReturnType, returnSpan := c.checkStatements(body)
		if !checkedMethod.IsAbstract() && !c.IsHeader {
			if c.mode == closureInferReturnTypeMode {
				c.addToReturnType(bodyReturnType)
				checkedMethod.ReturnType = c.returnType
			} else {
				if returnSpan == nil {
					returnSpan = span
				}
				c.checkCanAssign(bodyReturnType, returnType, returnSpan)
			}
		}
	}
	c.returnType = nil
	c.throwType = nil
	c.mode = prevMode
	c.catchScopes = prevCatchScopes
	return typedReturnTypeNode, typedThrowTypeNode
}

func (c *Checker) checkSpecialMethods(name value.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, span *position.Span) {
	if symbol.IsEqualityOperator(name) {
		c.checkEqualityOperator(name, checkedMethod, paramNodes, span)
		return
	}

	if symbol.IsRelationalOperator(name) {
		c.checkRelationalOperator(name, checkedMethod, paramNodes, span)
		return
	}

	if symbol.RequiresOneParameter(name) {
		c.checkFixedParameterCountMethod(name, checkedMethod, paramNodes, 1, span)
		return
	}

	if symbol.RequiresNoParameters(name) {
		c.checkFixedParameterCountMethod(name, checkedMethod, paramNodes, 0, span)
		return
	}
}

func (c *Checker) checkEqualityOperator(name value.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, span *position.Span) {
	params := checkedMethod.Params

	if !c.isTheSameType(checkedMethod.ReturnType, types.Bool{}, nil) {
		c.addFailure(
			fmt.Sprintf(
				"equality operator `%s` must return `%s`",
				lexer.Colorize(name.String()),
				lexer.Colorize("bool"),
			),
			span,
		)
	}

	if len(params) != 1 {
		c.addFailure(
			fmt.Sprintf(
				"equality operator `%s` must accept a single parameter, got %d",
				lexer.Colorize(name.String()),
				len(params),
			),
			span,
		)
		return
	}

	param := params[0]
	var paramSpan *position.Span
	if paramNodes != nil {
		paramSpan = paramNodes[0].Span()
	} else {
		paramSpan = span
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

func (c *Checker) checkRelationalOperator(name value.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, span *position.Span) {
	params := checkedMethod.Params

	if !c.isTheSameType(checkedMethod.ReturnType, types.Bool{}, nil) {
		c.addFailure(
			fmt.Sprintf(
				"relational operator `%s` must return `%s`",
				lexer.Colorize(name.String()),
				lexer.Colorize("bool"),
			),
			span,
		)
	}

	if len(params) != 1 {
		c.addFailure(
			fmt.Sprintf(
				"relational operator `%s` must accept a single parameter, got %d",
				lexer.Colorize(name.String()),
				len(params),
			),
			span,
		)
		return
	}

	param := checkedMethod.Params[0]
	var paramSpan *position.Span
	if paramNodes != nil {
		paramSpan = paramNodes[0].Span()
	} else {
		paramSpan = span
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

func (c *Checker) checkFixedParameterCountMethod(name value.Symbol, checkedMethod *types.Method, paramNodes []ast.ParameterNode, desiredParamCount int, span *position.Span) {
	params := checkedMethod.Params

	if types.IsVoid(checkedMethod.ReturnType) {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` cannot be void",
				lexer.Colorize(name.String()),
			),
			span,
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
			span,
		)
		return
	}

	for i, param := range params {
		var paramSpan *position.Span
		if paramNodes != nil {
			paramSpan = paramNodes[i].Span()
		} else {
			paramSpan = span
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

	c.returnType = c.newNormalisedUnion(c.returnType, typ)
}

func (c *Checker) checkMethodArgumentsAndInferTypeArguments(
	method *types.Method,
	positionalArguments []ast.ExpressionNode,
	namedArguments []ast.NamedArgumentNode,
	typeParams []*types.TypeParameter,
	span *position.Span,
) (
	_posArgs []ast.ExpressionNode,
	typeArgs types.TypeArgumentMap,
) {
	prevMode := c.mode
	c.mode = inferTypeArgumentMode
	defer c.setMode(prevMode)
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
			ast.NewUndefinedLiteralNode(span),
		)
	}

	typeArgMap := make(types.TypeArgumentMap)
	var currentParamIndex int
	for ; currentParamIndex < len(positionalArguments); currentParamIndex++ {
		posArg := positionalArguments[currentParamIndex]
		if currentParamIndex == positionalRestParamIndex {
			break
		}
		if currentParamIndex >= len(method.Params) {
			c.addWrongArgumentCountError(
				len(positionalArguments)+len(namedArguments),
				method,
				span,
			)
			break
		}
		param := method.Params[currentParamIndex]

		typedPosArg := c.checkExpressionWithType(posArg, param.Type)
		posArgType := c.typeOf(typedPosArg)

		inferredParamType := c.inferTypeArguments(posArgType, param.Type, typeArgMap, typedPosArg.Span())
		if inferredParamType == nil {
			param.Type = types.Untyped{}
		} else if inferredParamType != param.Type {
			param.Type = inferredParamType
		}
		typedPositionalArguments = append(typedPositionalArguments, typedPosArg)

		if !c.isSubtype(posArgType, param.Type, posArg.Span()) {
			c.addFailure(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
					types.InspectWithColor(param.Type),
					param.Name.String(),
					lexer.Colorize(method.Name.String()),
					types.InspectWithColor(posArgType),
				),
				posArg.Span(),
			)
		}
	}

	if method.HasPositionalRestParam() {
		if len(positionalArguments) < requiredPosParamCount {
			c.addFailure(
				fmt.Sprintf(
					"expected %d... positional arguments in call to `%s`, got %d",
					requiredPosParamCount,
					lexer.Colorize(method.Name.String()),
					len(positionalArguments),
				),
				span,
			)
			return nil, nil
		}
		restPositionalArguments := ast.NewArrayTupleLiteralNode(
			span,
			nil,
		)
		posRestParam := method.Params[positionalRestParamIndex]

		currentArgIndex := currentParamIndex
		for ; currentArgIndex < min(argCount-method.PostParamCount, len(positionalArguments)); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			typedPosArg := c.checkExpressionWithType(posArg, posRestParam.Type)
			posArgType := c.typeOf(typedPosArg)
			inferredParamType := c.inferTypeArguments(posArgType, posRestParam.Type, typeArgMap, typedPosArg.Span())
			if inferredParamType == nil {
				posRestParam.Type = types.Untyped{}
			} else if inferredParamType != posRestParam.Type {
				posRestParam.Type = inferredParamType
			}
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, typedPosArg)
			if !c.isSubtype(posArgType, posRestParam.Type, posArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(posRestParam.Type),
						posRestParam.Name.String(),
						lexer.Colorize(method.Name.String()),
						types.InspectWithColor(posArgType),
					),
					posArg.Span(),
				)
			}
		}
		typedPositionalArguments = append(typedPositionalArguments, restPositionalArguments)

		currentParamIndex = positionalRestParamIndex
		for ; currentArgIndex < len(positionalArguments); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			currentParamIndex++
			param := method.Params[currentParamIndex]

			typedPosArg := c.checkExpressionWithType(posArg, param.Type)
			posArgType := c.typeOf(typedPosArg)
			inferredParamType := c.inferTypeArguments(posArgType, param.Type, typeArgMap, typedPosArg.Span())
			if inferredParamType == nil {
				param.Type = types.Untyped{}
			} else if inferredParamType != param.Type {
				param.Type = inferredParamType
			}
			typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
			if !c.isSubtype(posArgType, param.Type, posArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name.String(),
						lexer.Colorize(method.Name.String()),
						types.InspectWithColor(posArgType),
					),
					posArg.Span(),
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

	for i := 0; i < len(method.Params); i++ {
		param := method.Params[i]
		switch param.Kind {
		case types.PositionalRestParameterKind, types.NamedRestParameterKind:
			continue
		}
		paramName := param.Name.String()
		var found bool

		for namedArgIndex, namedArgI := range namedArguments {
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			if namedArg.Name != paramName {
				continue
			}
			if found || i < firstNamedParamIndex {
				c.addFailure(
					fmt.Sprintf(
						"duplicated argument `%s` in call to `%s`",
						paramName,
						lexer.Colorize(method.Name.String()),
					),
					namedArg.Span(),
				)
			}
			found = true
			definedNamedArgumentsSlice[namedArgIndex] = true

			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, param.Type)
			namedArgType := c.typeOf(typedNamedArgValue)
			inferredParamType := c.inferTypeArguments(namedArgType, param.Type, typeArgMap, typedNamedArgValue.Span())
			if inferredParamType == nil {
				param.Type = types.Untyped{}
			} else if inferredParamType != param.Type {
				param.Type = inferredParamType
			}
			typedPositionalArguments = append(typedPositionalArguments, typedNamedArgValue)
			if !c.isSubtype(namedArgType, param.Type, namedArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name.String(),
						lexer.Colorize(method.Name.String()),
						types.InspectWithColor(namedArgType),
					),
					namedArg.Span(),
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
					lexer.Colorize(method.Name.String()),
				),
				span,
			)
		} else {
			// the parameter is missing and is optional
			// we push undefined as its value
			typedPositionalArguments = append(
				typedPositionalArguments,
				ast.NewUndefinedLiteralNode(span),
			)
		}
	}

	if method.HasNamedRestParam() {
		namedRestArgs := ast.NewHashRecordLiteralNode(
			span,
			nil,
		)
		namedRestParam := method.Params[len(method.Params)-1]
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)

			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, namedRestParam.Type)
			posArgType := c.typeOf(typedNamedArgValue)
			inferredParamType := c.inferTypeArguments(posArgType, namedRestParam.Type, typeArgMap, typedNamedArgValue.Span())
			if inferredParamType == nil {
				namedRestParam.Type = types.Untyped{}
			} else if inferredParamType != namedRestParam.Type {
				namedRestParam.Type = inferredParamType
			}
			namedRestArgs.Elements = append(
				namedRestArgs.Elements,
				ast.NewSymbolKeyValueExpressionNode(
					namedArg.Span(),
					namedArg.Name,
					typedNamedArgValue,
				),
			)
			namedArgType := c.typeOf(typedNamedArgValue)
			if !c.isSubtype(namedArgType, namedRestParam.Type, namedArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for named rest parameter `**%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(namedRestParam.Type),
						namedRestParam.Name.String(),
						lexer.Colorize(method.Name.String()),
						types.InspectWithColor(namedArgType),
					),
					namedArg.Span(),
				)
			}
		}

		typedPositionalArguments = append(typedPositionalArguments, namedRestArgs)
	} else {
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			c.addFailure(
				fmt.Sprintf(
					"nonexistent parameter `%s` given in call to `%s`",
					namedArg.Name,
					lexer.Colorize(method.Name.String()),
				),
				namedArg.Span(),
			)
		}
	}

	if len(typeArgMap) != len(typeParams) {
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
						lexer.Colorize(method.Name.String()),
					),
					span,
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

func (c *Checker) checkNonGenericMethodArguments(method *types.Method, positionalArguments []ast.ExpressionNode, namedArguments []ast.NamedArgumentNode, span *position.Span) []ast.ExpressionNode {
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
			ast.NewUndefinedLiteralNode(span),
		)
	}

	var currentParamIndex int
	for ; currentParamIndex < len(positionalArguments); currentParamIndex++ {
		posArg := positionalArguments[currentParamIndex]
		if currentParamIndex == positionalRestParamIndex {
			break
		}
		if currentParamIndex >= len(method.Params) {
			c.addWrongArgumentCountError(
				len(positionalArguments)+len(namedArguments),
				method,
				span,
			)
			break
		}
		param := method.Params[currentParamIndex]

		typedPosArg := c.checkExpressionWithType(posArg, param.Type)
		typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
		posArgType := c.typeOf(typedPosArg)
		if !c.isSubtype(posArgType, param.Type, posArg.Span()) {
			c.addFailure(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
					types.InspectWithColor(param.Type),
					param.Name.String(),
					lexer.Colorize(method.Name.String()),
					types.InspectWithColor(posArgType),
				),
				posArg.Span(),
			)
		}
	}

	if method.HasPositionalRestParam() {
		if len(positionalArguments) < requiredPosParamCount {
			c.addFailure(
				fmt.Sprintf(
					"expected %d... positional arguments in call to `%s`, got %d",
					requiredPosParamCount,
					lexer.Colorize(method.Name.String()),
					len(positionalArguments),
				),
				span,
			)
			return nil
		}
		restPositionalArguments := ast.NewArrayTupleLiteralNode(
			span,
			nil,
		)
		posRestParam := method.Params[positionalRestParamIndex]

		currentArgIndex := currentParamIndex
		for ; currentArgIndex < min(argCount-method.PostParamCount, len(positionalArguments)); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			typedPosArg := c.checkExpressionWithType(posArg, posRestParam.Type)
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, typedPosArg)
			posArgType := c.typeOf(typedPosArg)
			if !c.isSubtype(posArgType, posRestParam.Type, posArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(posRestParam.Type),
						posRestParam.Name.String(),
						lexer.Colorize(method.Name.String()),
						types.InspectWithColor(posArgType),
					),
					posArg.Span(),
				)
			}
		}
		typedPositionalArguments = append(typedPositionalArguments, restPositionalArguments)

		currentParamIndex = positionalRestParamIndex
		for ; currentArgIndex < len(positionalArguments); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			currentParamIndex++
			param := method.Params[currentParamIndex]

			typedPosArg := c.checkExpressionWithType(posArg, param.Type)
			typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
			posArgType := c.typeOf(typedPosArg)
			if !c.isSubtype(posArgType, param.Type, posArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name.String(),
						lexer.Colorize(method.Name.String()),
						types.InspectWithColor(posArgType),
					),
					posArg.Span(),
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

	for i := 0; i < len(method.Params); i++ {
		param := method.Params[i]
		switch param.Kind {
		case types.PositionalRestParameterKind, types.NamedRestParameterKind:
			continue
		}
		paramName := param.Name.String()
		var found bool

		for namedArgIndex, namedArgI := range namedArguments {
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			if namedArg.Name != paramName {
				continue
			}
			if found || i < firstNamedParamIndex {
				c.addFailure(
					fmt.Sprintf(
						"duplicated argument `%s` in call to `%s`",
						paramName,
						lexer.Colorize(method.Name.String()),
					),
					namedArg.Span(),
				)
			}
			found = true
			definedNamedArgumentsSlice[namedArgIndex] = true
			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, param.Type)
			namedArgType := c.typeOf(typedNamedArgValue)
			typedPositionalArguments = append(typedPositionalArguments, typedNamedArgValue)
			if !c.isSubtype(namedArgType, param.Type, namedArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name.String(),
						lexer.Colorize(method.Name.String()),
						types.InspectWithColor(namedArgType),
					),
					namedArg.Span(),
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
					lexer.Colorize(method.Name.String()),
				),
				span,
			)
		} else {
			// the parameter is missing and is optional
			// we push undefined as its value
			typedPositionalArguments = append(
				typedPositionalArguments,
				ast.NewUndefinedLiteralNode(span),
			)
		}
	}

	if method.HasNamedRestParam() {
		namedRestArgs := ast.NewHashRecordLiteralNode(
			span,
			nil,
		)
		namedRestParam := method.Params[len(method.Params)-1]
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, namedRestParam.Type)
			namedRestArgs.Elements = append(
				namedRestArgs.Elements,
				ast.NewSymbolKeyValueExpressionNode(
					namedArg.Span(),
					namedArg.Name,
					typedNamedArgValue,
				),
			)
			namedArgType := c.typeOf(typedNamedArgValue)
			if !c.isSubtype(namedArgType, namedRestParam.Type, namedArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for named rest parameter `**%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(namedRestParam.Type),
						namedRestParam.Name.String(),
						lexer.Colorize(method.Name.String()),
						types.InspectWithColor(namedArgType),
					),
					namedArg.Span(),
				)
			}
		}

		typedPositionalArguments = append(typedPositionalArguments, namedRestArgs)
	} else {
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			c.addFailure(
				fmt.Sprintf(
					"nonexistent parameter `%s` given in call to `%s`",
					namedArg.Name,
					lexer.Colorize(method.Name.String()),
				),
				namedArg.Span(),
			)
		}
	}

	return typedPositionalArguments
}

func (c *Checker) checkMethodArguments(
	method *types.Method,
	typeArgumentNodes []ast.TypeNode,
	positionalArgumentNodes []ast.ExpressionNode,
	namedArgumentNodes []ast.NamedArgumentNode,
	span *position.Span,
) (_method *types.Method, typedPositionalArguments []ast.ExpressionNode) {
	if len(typeArgumentNodes) > 0 {
		typeArgs, ok := c.checkTypeArguments(
			method,
			typeArgumentNodes,
			method.TypeParameters,
			span,
		)
		if !ok {
			c.checkExpressions(positionalArgumentNodes)
			c.checkNamedArguments(namedArgumentNodes)
			return nil, nil
		}

		method = c.replaceTypeParametersInMethodCopy(method, typeArgs.ArgumentMap)
		typedPositionalArguments = c.checkNonGenericMethodArguments(method, positionalArgumentNodes, namedArgumentNodes, span)
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
			span,
		)
		if len(typeArgMap) != len(method.TypeParameters) {
			return nil, nil
		}
		method.ReturnType = c.replaceTypeParameters(method.ReturnType, typeArgMap)
		method.ThrowType = c.replaceTypeParameters(method.ThrowType, typeArgMap)
		return method, typedPositionalArguments
	}

	typedPositionalArguments = c.checkNonGenericMethodArguments(method, positionalArgumentNodes, namedArgumentNodes, span)
	return method, typedPositionalArguments
}

func (c *Checker) checkSimpleMethodCall(
	receiver ast.ExpressionNode,
	op token.Type,
	methodName value.Symbol,
	typeArgumentNodes []ast.TypeNode,
	positionalArgumentNodes []ast.ExpressionNode,
	namedArgumentNodes []ast.NamedArgumentNode,
	span *position.Span,
) (
	_receiver ast.ExpressionNode,
	_positionalArguments []ast.ExpressionNode,
	typ types.Type,
) {
	receiver = c.checkExpression(receiver)
	receiverType := c.typeOf(receiver)

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
		method = c.getMethod(receiverType, methodName, span)
	case token.QUESTION_DOT, token.QUESTION_DOT_DOT:
		nonNilableReceiverType := c.toNonNilable(receiverType)
		method = c.getMethod(nonNilableReceiverType, methodName, span)
	default:
		panic(fmt.Sprintf("invalid call operator: %#v", op))
	}
	if method == nil {
		c.checkExpressions(positionalArgumentNodes)
		c.checkNamedArguments(namedArgumentNodes)
		return receiver, positionalArgumentNodes, types.Untyped{}
	}

	c.addToMethodCache(method)

	method, typedPositionalArguments := c.checkMethodArguments(method, typeArgumentNodes, positionalArgumentNodes, namedArgumentNodes, span)
	if method == nil {
		return receiver, positionalArgumentNodes, types.Untyped{}
	}

	var returnType types.Type
	switch op {
	case token.DOT:
		returnType = method.ReturnType
	case token.QUESTION_DOT:
		if !c.isNilable(receiverType) {
			c.addFailure(
				fmt.Sprintf("cannot make a nil-safe call on type `%s` which is not nilable", types.InspectWithColor(receiverType)),
				span,
			)
			returnType = method.ReturnType
		} else {
			returnType = c.toNilable(method.ReturnType)
		}
	case token.DOT_DOT:
		returnType = receiverType
	case token.QUESTION_DOT_DOT:
		if !c.isNilable(receiverType) {
			c.addFailure(
				fmt.Sprintf("cannot make a nil-safe call on type `%s` which is not nilable", types.InspectWithColor(receiverType)),
				span,
			)
		}
		returnType = receiverType
	}

	c.checkCalledMethodThrowType(method, span)

	return receiver, typedPositionalArguments, returnType
}

func (c *Checker) checkBinaryOpMethodCall(
	left ast.ExpressionNode,
	right ast.ExpressionNode,
	methodName value.Symbol,
	span *position.Span,
) types.Type {
	_, _, returnType := c.checkSimpleMethodCall(
		left,
		token.DOT,
		methodName,
		nil,
		[]ast.ExpressionNode{right},
		nil,
		span,
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
		node.Span(),
	)

	node.ReturnType = returnType
	node.ThrowType = throwType

	c.method = nil

	method.CalledMethods = c.methodCache.Slice
	c.methodCache.Slice = nil

	if c.shouldCompile() && method.IsCompilable() {
		method.Bytecode = c.compiler.CompileMethodBody(node, method.Name)
	}
}

func (c *Checker) declareMethod(
	baseMethod *types.Method,
	methodNamespace types.Namespace,
	docComment string,
	abstract bool,
	sealed bool,
	inferReturnType bool,
	name value.Symbol,
	typeParamNodes []ast.TypeParameterNode,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	span *position.Span,
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
					types.InspectModifier(abstract, sealed, false),
					types.InspectModifier(oldMethod.IsAbstract(), oldMethod.IsSealed(), false),
				),
				span,
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
				span,
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
				span,
			)
		}
	default:
		if abstract {
			c.addFailure(
				fmt.Sprintf(
					"cannot declare abstract method `%s` in this context",
					name.String(),
				),
				span,
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
		typeParamMod = types.NewTypeParamNamespace(fmt.Sprintf("Type Parameter Container of %s", name))
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

	c.mode = prevMode
	c.setInputPositionTypeMode()
	var params []*types.Parameter
	for i, paramNode := range paramNodes {
		switch p := paramNode.(type) {
		case *ast.FormalParameterNode:
			var declaredType types.Type
			if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.typeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Type
			} else {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					paramNode.Span(),
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
			name := value.ToSymbol(p.Name)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType)
		case *ast.MethodParameterNode:
			var declaredType types.Type
			if p.SetInstanceVariable {
				currentIvar, _ := c.getInstanceVariableIn(value.ToSymbol(p.Name), methodNamespace)
				if p.TypeNode == nil {
					if currentIvar == nil {
						c.addFailure(
							fmt.Sprintf(
								"cannot infer the type of instance variable `%s`",
								p.Name,
							),
							p.Span(),
						)
					}

					declaredType = currentIvar
				} else {
					p.TypeNode = c.checkTypeNode(p.TypeNode)
					declaredType = c.typeOf(p.TypeNode)
					if currentIvar != nil {
						c.checkCanAssignInstanceVariable(p.Name, declaredType, currentIvar, p.TypeNode.Span())
					} else {
						c.declareInstanceVariable(value.ToSymbol(p.Name), declaredType, p.Span())
					}
				}
			} else if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.typeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Type
			} else {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					paramNode.Span(),
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
			name := value.ToSymbol(p.Name)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType)
		case *ast.SignatureParameterNode:
			var declaredType types.Type
			if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.typeOf(p.TypeNode)
			} else if baseMethod != nil && len(baseMethod.Params) > i {
				declaredType = baseMethod.Params[i].Type
			} else {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					paramNode.Span(),
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
			name := value.ToSymbol(p.Name)
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
				paramNode.Span(),
			)
		}
	}

	c.mode = prevMode
	c.setOutputPositionTypeMode()

	var returnType types.Type
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.typeOf(typedReturnTypeNode)
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
		throwType = c.typeOf(typedThrowTypeNode)
	} else if baseMethod != nil && baseMethod.ThrowType != nil {
		throwType = baseMethod.ThrowType
	} else {
		throwType = types.Never{}
	}

	newMethod := types.NewMethod(
		docComment,
		abstract,
		sealed,
		c.IsHeader,
		name,
		typeParams,
		params,
		returnType,
		throwType,
		methodNamespace,
	)
	newMethod.SetLocation(c.newLocation(span))

	c.checkMethodOverrideWithPlaceholder(newMethod, oldMethod, span)
	methodNamespace.SetMethod(name, newMethod)

	c.checkSpecialMethods(name, newMethod, paramNodes, span)

	c.mode = prevMode

	return newMethod, typeParamMod
}

func (c *Checker) checkMethodOverrideWithPlaceholder(
	overrideMethod,
	baseMethod *types.Method,
	span *position.Span,
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
		span,
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

func (c *Checker) checkMethodCompatibilityForAlgebraicTypes(baseMethod, overrideMethod *types.Method, errSpan *position.Span) bool {
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

func (c *Checker) checkMethodCompatibilityForInterfaceIntersection(baseMethod, overrideMethod *types.Method, errSpan *position.Span, typeArgs types.TypeArgumentMap) bool {
	areCompatible := c.checkMethodCompatibilityAndInferTypeArgs(baseMethod, overrideMethod, errSpan, typeArgs)
	return areCompatible
}

// Checks whether two methods are compatible.
func (c *Checker) checkMethodCompatibility(baseMethod, overrideMethod *types.Method, errSpan *position.Span, validateNames bool) bool {
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

			if oldParam.Kind != newParam.Kind || !c.isSubtype(oldParam.Type, newParam.Type, errSpan) {
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

func (c *Checker) checkMethodCompatibilityAndInferTypeArgs(baseMethod, overrideMethod *types.Method, errSpan *position.Span, typeArgs types.TypeArgumentMap) bool {
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

func (c *Checker) getMethod(typ types.Type, name value.Symbol, errSpan *position.Span) *types.Method {
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
						c.replaceTypeParametersInMethod(methodCopy, generic.ArgumentMap)
						continue
					}

					result := c.replaceTypeParametersInMethodCopy(method, generic.ArgumentMap)
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
						c.replaceTypeParametersInMethod(methodCopy, generic.ArgumentMap)
						continue
					}

					result := c.replaceTypeParametersInMethodCopy(method, generic.ArgumentMap)
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
				c.replaceTypeParametersInMethod(methodCopy, generic.ArgumentMap)
				continue
			}

			result := c.replaceTypeParametersInMethodCopy(method, generic.ArgumentMap)
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
					c.replaceTypeParametersInMethod(methodCopy, generic.ArgumentMap)
					continue
				}

				result := c.replaceTypeParametersInMethodCopy(method, generic.ArgumentMap)
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

func (c *Checker) _getMethodInNamespace(namespace types.Namespace, typ types.Type, name value.Symbol, errSpan *position.Span, inParent bool) *types.Method {
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

func (c *Checker) getMethodInNamespaceWithSelf(namespace types.Namespace, typ types.Type, name value.Symbol, self types.Type, errSpan *position.Span, inParent, inSelf bool) *types.Method {
	method := c._getMethodInNamespace(namespace, typ, name, errSpan, inParent)
	if method == nil {
		return nil
	}
	if inSelf {
		return method
	}
	m := c.createTypeArgumentMapWithSelf(self)
	return c.replaceTypeParametersInMethodCopy(method, m)
}

func (c *Checker) getMethodInNamespace(namespace types.Namespace, typ types.Type, name value.Symbol, errSpan *position.Span, inParent, inSelf bool) *types.Method {
	return c.getMethodInNamespaceWithSelf(namespace, typ, name, namespace, errSpan, inParent, inSelf)
}

func (c *Checker) replaceTypeParametersInMethodCopy(method *types.Method, typeArgs types.TypeArgumentMap) *types.Method {
	var methodCopy *types.Method

	for i, typeParam := range method.TypeParameters {
		result := c.replaceTypeParameters(typeParam.LowerBound, typeArgs)
		if typeParam.LowerBound != result {
			if methodCopy == nil {
				methodCopy = c.deepCopyMethod(method)
			}
			methodCopy.TypeParameters[i].LowerBound = result
		}
		result = c.replaceTypeParameters(typeParam.UpperBound, typeArgs)
		if typeParam.UpperBound != result {
			if methodCopy == nil {
				methodCopy = c.deepCopyMethod(method)
			}
			methodCopy.TypeParameters[i].UpperBound = result
		}
	}
	result := c.replaceTypeParameters(method.ReturnType, typeArgs)
	if method.ReturnType != result {
		if methodCopy == nil {
			methodCopy = c.deepCopyMethod(method)
		}
		methodCopy.ReturnType = result
	}
	result = c.replaceTypeParameters(method.ThrowType, typeArgs)
	if method.ThrowType != result {
		if methodCopy == nil {
			methodCopy = c.deepCopyMethod(method)
		}
		methodCopy.ThrowType = result
	}

	for i, param := range method.Params {
		result := c.replaceTypeParameters(param.Type, typeArgs)
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

func (c *Checker) replaceTypeParametersInMethod(method *types.Method, typeArgs types.TypeArgumentMap) *types.Method {
	for _, typeParam := range method.TypeParameters {
		typeParam.LowerBound = c.replaceTypeParameters(typeParam.LowerBound, typeArgs)
		typeParam.UpperBound = c.replaceTypeParameters(typeParam.UpperBound, typeArgs)
	}
	method.ReturnType = c.replaceTypeParameters(method.ReturnType, typeArgs)
	method.ThrowType = c.replaceTypeParameters(method.ThrowType, typeArgs)

	for _, param := range method.Params {
		param.Type = c.replaceTypeParameters(param.Type, typeArgs)
	}

	return method
}

func (c *Checker) replaceTypeParametersInWhere(whereParams []*types.TypeParameter, whereArgs []types.Type, typeArgs types.TypeArgumentMap) {
	for i, whereArg := range whereArgs {
		whereArgs[i] = c.replaceTypeParameters(whereArg, typeArgs)
	}

	for i, whereParam := range whereParams {
		var whereParamCopy *types.TypeParameter

		result := c.replaceTypeParameters(whereParam.LowerBound, typeArgs)
		if result != whereParam.LowerBound {
			whereParamCopy = whereParam.Copy()
			whereParams[i] = whereParamCopy
			whereParamCopy.LowerBound = result
		}

		result = c.replaceTypeParameters(whereParam.UpperBound, typeArgs)
		if result != whereParam.LowerBound {
			if whereParamCopy == nil {
				whereParamCopy = whereParam.Copy()
				whereParams[i] = whereParamCopy
			}
			whereParamCopy.UpperBound = result
		}
	}
}

func (c *Checker) getMethodForTypeParameter(typ *types.TypeParameter, name value.Symbol, errSpan *position.Span, inParent, inSelf bool) *types.Method {
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
		return c.replaceTypeParametersInMethodCopy(method, typeArgMap)
	default:
		return c._getMethod(typ.UpperBound, name, errSpan, inParent, inSelf)
	}
}

func (c *Checker) getReceiverlessMethod(name value.Symbol, span *position.Span) (_ *types.Method, fromLocal bool) {
	local, _ := c.resolveLocal(name.String(), nil)
	if local != nil {
		return c.getMethod(local.typ, symbol.L_call, span), true
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

	c.addMissingMethodError(c.selfType, name.String(), span)

	return nil, false
}

func (c *Checker) _getMethod(typ types.Type, name value.Symbol, errSpan *position.Span, inParent, inSelf bool) *types.Method {
	typ = c.toNonLiteral(typ, true)

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
		nilType := c.GlobalEnv.StdSubtype(symbol.Nil).(*types.Class)
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

func (c *Checker) addMissingMethodError(typ types.Type, name string, span *position.Span) {
	if types.IsUntyped(typ) {
		return
	}
	c.addFailure(
		fmt.Sprintf("method `%s` is not defined on type `%s`", name, types.InspectWithColor(typ)),
		span,
	)
}
