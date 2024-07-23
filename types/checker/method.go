// Package checker implements the Elk type checker
package checker

import (
	"fmt"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

const concurrencyLimit = 10_000

func (c *Checker) checkMethods() {
	concurrent.Foreach(
		concurrencyLimit,
		c.methodChecks.Slice,
		func(methodCheck methodCheckEntry) {
			methodChecker := c.newChecker(
				methodCheck.filename,
				methodCheck.constantScopes,
				methodCheck.methodScopes,
				methodCheck.method.DefinedUnder,
				methodCheck.method.ReturnType,
				methodCheck.method.ThrowType,
			)
			node := methodCheck.node
			methodChecker.checkMethodDefinition(node)
		},
	)
}

func (c *Checker) declareMethodForGetter(node *ast.AttributeParameterNode, docComment string) {
	method := c.declareMethod(
		docComment,
		false,
		false,
		value.ToSymbol(node.Name),
		nil,
		node.TypeNode,
		nil,
		node.Span(),
	)

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
		node.Span(),
		"",
		false,
		false,
		node.Name,
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
}

func (c *Checker) declareMethodForSetter(node *ast.AttributeParameterNode, docComment string) {
	setterName := node.Name + "="

	params := []ast.ParameterNode{
		ast.NewMethodParameterNode(
			node.TypeNode.Span(),
			node.Name,
			true,
			node.TypeNode,
			nil,
			ast.NormalParameterKind,
		),
	}
	method := c.declareMethod(
		docComment,
		false,
		false,
		value.ToSymbol(setterName),
		params,
		nil,
		nil,
		node.Span(),
	)

	methodNode := ast.NewMethodDefinitionNode(
		node.Span(),
		docComment,
		false,
		false,
		setterName,
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
}

func (c *Checker) addWrongArgumentCountError(got int, method *types.Method, span *position.Span) {
	c.addFailure(
		fmt.Sprintf("expected %s arguments in call to `%s`, got %d", method.ExpectedParamCountString(), method.Name, got),
		span,
	)
}

func (c *Checker) addOverrideSealedMethodError(baseMethod *types.Method, span *position.Span) {
	c.addFailure(
		fmt.Sprintf(
			"cannot override sealed method `%s`\n  previous definition found in `%s`, with signature: `%s`",
			baseMethod.Name,
			types.InspectWithColor(baseMethod.DefinedUnder),
			baseMethod.InspectSignatureWithColor(true),
		),
		span,
	)
}

func (c *Checker) checkMethodOverride(
	overrideMethod,
	baseMethod *types.Method,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	span *position.Span,
) {
	name := baseMethod.Name
	if baseMethod.IsSealed() {
		c.addOverrideSealedMethodError(baseMethod, span)
	}
	if !baseMethod.IsAbstract() && overrideMethod.IsAbstract() {
		c.addFailure(
			fmt.Sprintf(
				"cannot override method `%s` with a different modifier, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
				name,
				types.InspectModifier(overrideMethod.IsAbstract(), overrideMethod.IsSealed(), false),
				types.InspectModifier(baseMethod.IsAbstract(), baseMethod.IsSealed(), false),
				types.InspectWithColor(baseMethod.DefinedUnder),
				baseMethod.InspectSignatureWithColor(true),
			),
			span,
		)
	}

	if !c.isSubtype(overrideMethod.ReturnType, baseMethod.ReturnType, nil) {
		var returnSpan *position.Span
		if returnTypeNode != nil {
			returnSpan = returnTypeNode.Span()
		} else {
			returnSpan = span
		}
		c.addFailure(
			fmt.Sprintf(
				"cannot override method `%s` with a different return type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
				name,
				types.InspectWithColor(overrideMethod.ReturnType),
				types.InspectWithColor(baseMethod.ReturnType),
				types.InspectWithColor(baseMethod.DefinedUnder),
				baseMethod.InspectSignatureWithColor(true),
			),
			returnSpan,
		)
	}
	if !c.isSubtype(overrideMethod.ThrowType, baseMethod.ThrowType, nil) {
		var throwSpan *position.Span
		if throwTypeNode != nil {
			throwSpan = throwTypeNode.Span()
		} else {
			throwSpan = span
		}
		c.addFailure(
			fmt.Sprintf(
				"cannot override method `%s` with a different throw type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
				name,
				types.InspectWithColor(overrideMethod.ThrowType),
				types.InspectWithColor(baseMethod.ThrowType),
				types.InspectWithColor(baseMethod.DefinedUnder),
				baseMethod.InspectSignatureWithColor(true),
			),
			throwSpan,
		)
	}

	if len(baseMethod.Params) > len(overrideMethod.Params) {
		paramSpan := position.JoinSpanOfCollection(paramNodes)
		if paramSpan == nil {
			paramSpan = span
		}
		c.addFailure(
			fmt.Sprintf(
				"cannot override method `%s` with less parameters\n  previous definition found in `%s`, with signature: `%s`",
				name,
				types.InspectWithColor(baseMethod.DefinedUnder),
				baseMethod.InspectSignatureWithColor(true),
			),
			paramSpan,
		)
	} else {
		for i := range len(baseMethod.Params) {
			oldParam := baseMethod.Params[i]
			newParam := overrideMethod.Params[i]
			if oldParam.Name != newParam.Name {
				c.addFailure(
					fmt.Sprintf(
						"cannot override method `%s` with invalid parameter name, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
						name,
						newParam.Name,
						oldParam.Name,
						types.InspectWithColor(baseMethod.DefinedUnder),
						baseMethod.InspectSignatureWithColor(true),
					),
					paramNodes[i].Span(),
				)
				continue
			}
			if oldParam.Kind != newParam.Kind {
				c.addFailure(
					fmt.Sprintf(
						"cannot override method `%s` with invalid parameter kind, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
						name,
						newParam.NameWithKind(),
						oldParam.NameWithKind(),
						types.InspectWithColor(baseMethod.DefinedUnder),
						baseMethod.InspectSignatureWithColor(true),
					),
					paramNodes[i].Span(),
				)
				continue
			}
			if !c.isSubtype(oldParam.Type, newParam.Type, paramNodes[i].Span()) {
				c.addFailure(
					fmt.Sprintf(
						"cannot override method `%s` with invalid parameter type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
						name,
						types.InspectWithColor(newParam.Type),
						types.InspectWithColor(oldParam.Type),
						types.InspectWithColor(baseMethod.DefinedUnder),
						baseMethod.InspectSignatureWithColor(true),
					),
					paramNodes[i].Span(),
				)
				continue
			}
		}

		for i := len(baseMethod.Params); i < len(overrideMethod.Params); i++ {
			param := overrideMethod.Params[i]
			if !param.IsOptional() {
				c.addFailure(
					fmt.Sprintf(
						"cannot override method `%s` with additional parameter `%s`\n  previous definition found in `%s`, with signature: `%s`",
						name,
						param.Name,
						types.InspectWithColor(baseMethod.DefinedUnder),
						baseMethod.InspectSignatureWithColor(true),
					),
					paramNodes[i].Span(),
				)
			}
		}
	}

}

func (c *Checker) checkMethod(
	checkedMethod *types.Method,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	body []ast.StatementNode,
	span *position.Span,
) (ast.TypeNode, ast.TypeNode) {
	methodNamespace := c.currentMethodScope().container
	name := checkedMethod.Name

	currentMethod := types.GetMethodInNamespace(methodNamespace, name)
	if checkedMethod != currentMethod && checkedMethod.IsSealed() {
		c.addOverrideSealedMethodError(checkedMethod, currentMethod.Span())
	}

	parent := methodNamespace.Parent()

	if parent != nil {
		baseMethod := types.GetMethodInNamespace(parent, name)
		if baseMethod != nil {
			c.checkMethodOverride(
				checkedMethod,
				baseMethod,
				paramNodes,
				returnTypeNode,
				throwTypeNode,
				span,
			)
		}
	}

	c.pushIsolatedLocalEnv()
	defer c.popLocalEnv()

	for _, param := range paramNodes {
		p, _ := param.(*ast.MethodParameterNode)
		var declaredType types.Type
		var declaredTypeNode ast.TypeNode
		if p.TypeNode != nil {
			declaredTypeNode = p.TypeNode
			declaredType = c.typeOf(declaredTypeNode)
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
	}

	var returnType types.Type
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.typeOf(typedReturnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.typeOf(typedThrowTypeNode)
	}

	if len(body) > 0 && checkedMethod.IsAbstract() {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` cannot have a body because it is abstract",
				name,
			),
			span,
		)
	}

	previousMode := c.mode
	c.mode = methodMode
	defer c.setMode(previousMode)
	c.returnType = returnType
	c.throwType = throwType
	c.checkStatements(body)
	c.returnType = nil
	c.throwType = nil
	return typedReturnTypeNode, typedThrowTypeNode
}

func (c *Checker) checkMethodArguments(method *types.Method, positionalArguments []ast.ExpressionNode, namedArguments []ast.NamedArgumentNode, span *position.Span) []ast.ExpressionNode {
	reqParamCount := method.RequiredParamCount()
	requiredPosParamCount := len(method.Params) - method.OptionalParamCount
	if method.PostParamCount != -1 {
		requiredPosParamCount -= method.PostParamCount + 1
	}
	positionalRestParamIndex := method.PositionalRestParamIndex()
	var typedPositionalArguments []ast.ExpressionNode

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

		typedPosArg := c.checkExpression(posArg)
		typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
		posArgType := c.typeOf(typedPosArg)
		if !c.isSubtype(posArgType, param.Type, posArg.Span()) {
			c.addFailure(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
					types.InspectWithColor(param.Type),
					param.Name,
					method.Name,
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
					method.Name,
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
		for ; currentArgIndex < len(positionalArguments)-method.PostParamCount; currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			typedPosArg := c.checkExpression(posArg)
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, typedPosArg)
			posArgType := c.typeOf(typedPosArg)
			if !c.isSubtype(posArgType, posRestParam.Type, posArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(posRestParam.Type),
						posRestParam.Name,
						method.Name,
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

			typedPosArg := c.checkExpression(posArg)
			typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
			posArgType := c.typeOf(typedPosArg)
			if !c.isSubtype(posArgType, param.Type, posArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name,
						method.Name,
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
						method.Name,
					),
					namedArg.Span(),
				)
			}
			found = true
			definedNamedArgumentsSlice[namedArgIndex] = true
			typedNamedArgValue := c.checkExpression(namedArg.Value)
			namedArgType := c.typeOf(typedNamedArgValue)
			typedPositionalArguments = append(typedPositionalArguments, typedNamedArgValue)
			if !c.isSubtype(namedArgType, param.Type, namedArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name,
						method.Name,
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
					method.Name,
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

	if method.HasNamedRestParam {
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
			typedNamedArgValue := c.checkExpression(namedArg.Value)
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
						namedRestParam.Name,
						method.Name,
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
					method.Name,
				),
				namedArg.Span(),
			)
		}
	}

	return typedPositionalArguments
}

func (c *Checker) checkSimpleMethodCall(
	receiver ast.ExpressionNode,
	op token.Type,
	methodName value.Symbol,
	positionalArguments []ast.ExpressionNode,
	namedArguments []ast.NamedArgumentNode,
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
	if types.IsNever(receiverType) || types.IsNothing(receiverType) {
		var typedPositionalArguments []ast.ExpressionNode

		for _, argument := range positionalArguments {
			typedPositionalArguments = append(typedPositionalArguments, c.checkExpression(argument))
		}
		for _, argument := range namedArguments {
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
		c.checkExpressions(positionalArguments)
		return receiver, positionalArguments, types.Nothing{}
	}

	typedPositionalArguments := c.checkMethodArguments(method, positionalArguments, namedArguments, span)
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
		[]ast.ExpressionNode{right},
		nil,
		span,
	)

	return returnType
}

func (c *Checker) checkMethodDefinition(node *ast.MethodDefinitionNode) {
	returnType, throwType := c.checkMethod(
		c.typeOf(node).(*types.Method),
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Body,
		node.Span(),
	)
	node.ReturnType = returnType
	node.ThrowType = throwType
}

func (c *Checker) declareMethod(
	docComment string,
	abstract bool,
	sealed bool,
	name value.Symbol,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	span *position.Span,
) *types.Method {
	if c.mode == interfaceMode {
		abstract = true
	}
	methodScope := c.currentMethodScope()
	methodNamespace := methodScope.container
	oldMethod := methodNamespace.Method(name)
	if oldMethod != nil {
		if oldMethod.IsNative() && oldMethod.IsSealed() {
			c.addOverrideSealedMethodError(oldMethod, span)
		} else if sealed && !oldMethod.IsSealed() {
			c.addFailure(
				fmt.Sprintf(
					"cannot redeclare method `%s` with a different modifier, is `%s`, should be `%s`",
					name,
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
					name,
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
					name,
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
					name,
				),
				span,
			)
		}
	}

	var params []*types.Parameter
	for _, param := range paramNodes {
		switch p := param.(type) {
		case *ast.MethodParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
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
					declaredTypeNode = c.checkTypeNode(p.TypeNode)
					declaredType = c.typeOf(declaredTypeNode)
					if currentIvar != nil {
						c.checkCanAssignInstanceVariable(p.Name, declaredType, currentIvar, declaredTypeNode.Span())
					} else {
						c.declareInstanceVariable(value.ToSymbol(p.Name), declaredType, p.Span())
					}
				}
			} else if p.TypeNode == nil {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					param.Span(),
				)
			} else {
				declaredTypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.typeOf(declaredTypeNode)
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
			params = append(params, types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			))
		case *ast.SignatureParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
			if p.TypeNode == nil {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					param.Span(),
				)
			} else {
				declaredTypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.typeOf(declaredTypeNode)
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
			params = append(params, types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			))
		default:
			c.addFailure(
				fmt.Sprintf("invalid param type %T", param),
				param.Span(),
			)
		}
	}

	var returnType types.Type
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.typeOf(typedReturnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.typeOf(typedThrowTypeNode)
	}
	newMethod := types.NewMethod(
		docComment,
		abstract,
		sealed,
		c.HeaderMode,
		name,
		params,
		returnType,
		throwType,
		methodScope.container,
	)
	newMethod.SetSpan(span)

	methodScope.container.SetMethod(name, newMethod)

	return newMethod
}

// Checks whether two methods are compatible.
func (c *Checker) checkMethodCompatibility(baseMethod, overrideMethod *types.Method, errSpan *position.Span) bool {
	areCompatible := true
	if baseMethod != nil {
		if !c.isSubtype(overrideMethod.ReturnType, baseMethod.ReturnType, errSpan) {
			c.addFailure(
				fmt.Sprintf(
					"method `%s` has a different return type than `%s`, has `%s`, should have `%s`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					types.InspectWithColor(overrideMethod.ReturnType),
					types.InspectWithColor(baseMethod.ReturnType),
				),
				errSpan,
			)
			areCompatible = false
		}
		if !c.isSubtype(overrideMethod.ThrowType, baseMethod.ThrowType, errSpan) {
			c.addFailure(
				fmt.Sprintf(
					"method `%s` has a different throw type than `%s`, has `%s`, should have `%s`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					types.InspectWithColor(overrideMethod.ThrowType),
					types.InspectWithColor(baseMethod.ThrowType),
				),
				errSpan,
			)
			areCompatible = false
		}

		if len(baseMethod.Params) > len(overrideMethod.Params) {
			c.addFailure(
				fmt.Sprintf(
					"method `%s` has less parameters than `%s`, has `%d`, should have `%d`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					len(overrideMethod.Params),
					len(baseMethod.Params),
				),
				errSpan,
			)
			areCompatible = false
		} else {
			for i := range len(baseMethod.Params) {
				oldParam := baseMethod.Params[i]
				newParam := overrideMethod.Params[i]
				if oldParam.Name != newParam.Name {
					c.addFailure(
						fmt.Sprintf(
							"method `%s` has a different parameter name than `%s`, has `%s`, should have `%s`",
							types.InspectWithColor(overrideMethod),
							types.InspectWithColor(baseMethod),
							newParam.Name,
							oldParam.Name,
						),
						errSpan,
					)
					areCompatible = false
					continue
				}
				if oldParam.Kind != newParam.Kind {
					c.addFailure(
						fmt.Sprintf(
							"method `%s` has a different parameter kind than `%s`, has `%s`, should have `%s`",
							types.InspectWithColor(overrideMethod),
							types.InspectWithColor(baseMethod),
							newParam.NameWithKind(),
							oldParam.NameWithKind(),
						),
						errSpan,
					)
					areCompatible = false
					continue
				}
				if !c.isSubtype(oldParam.Type, newParam.Type, errSpan) {
					c.addFailure(
						fmt.Sprintf(
							"method `%s` has a different type for parameter `%s` than `%s`, has `%s`, should have `%s`",
							types.InspectWithColor(overrideMethod),
							newParam.Name,
							types.InspectWithColor(baseMethod),
							types.InspectWithColor(newParam.Type),
							types.InspectWithColor(oldParam.Type),
						),
						errSpan,
					)
					areCompatible = false
					continue
				}
			}

			for i := len(baseMethod.Params); i < len(overrideMethod.Params); i++ {
				param := overrideMethod.Params[i]
				if !param.IsOptional() {
					c.addFailure(
						fmt.Sprintf(
							"method `%s` has a required parameter missing in `%s`, got `%s`",
							types.InspectWithColor(overrideMethod),
							types.InspectWithColor(baseMethod),
							param.Name,
						),
						errSpan,
					)
					areCompatible = false
				}
			}
		}

	}

	return areCompatible
}

func (c *Checker) getMethod(typ types.Type, name value.Symbol, errSpan *position.Span) *types.Method {
	return c._getMethod(typ, name, errSpan, false)
}

func (c *Checker) getMethodInContainer(container types.Namespace, typ types.Type, name value.Symbol, errSpan *position.Span, inParent bool) *types.Method {
	method := types.GetMethodInNamespace(container, name)
	if method != nil {
		return method
	}
	if !inParent {
		c.addMissingMethodError(typ, name.String(), errSpan)
	}
	return nil
}

func (c *Checker) _getMethod(typ types.Type, name value.Symbol, errSpan *position.Span, inParent bool) *types.Method {
	typ = c.toNonLiteral(typ)

	switch t := typ.(type) {
	case *types.NamedType:
		return c._getMethod(t.Type, name, errSpan, inParent)
	case *types.Class:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.SingletonClass:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.Interface:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.InterfaceProxy:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.Module:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.Mixin:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.MixinProxy:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
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
				if baseMethod == nil || len(baseMethod.Params) > len(elementMethod.Params) {
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

			if !c.checkMethodCompatibility(baseMethod, method, errSpan) {
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
		if len(nilMethod.Params) < len(nonNilMethod.Params) {
			baseMethod = nilMethod
			overrideMethod = nonNilMethod
		} else {
			baseMethod = nonNilMethod
			overrideMethod = nilMethod
		}

		if c.checkMethodCompatibility(baseMethod, overrideMethod, errSpan) {
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
			if baseMethod == nil || len(baseMethod.Params) > len(elementMethod.Params) {
				baseMethod = elementMethod
			}
		}

		if len(methods) < len(t.Elements) {
			return nil
		}

		isCompatible := true
		for i := range len(methods) {
			method := methods[i]

			if !c.checkMethodCompatibility(baseMethod, method, errSpan) {
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
	c.addFailure(
		fmt.Sprintf("method `%s` is not defined on type `%s`", name, types.InspectWithColor(typ)),
		span,
	)
}
