package checker

import (
	"fmt"
	"iter"
	"maps"
	"strings"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type methodCheckEntry struct {
	filename       string
	method         *types.Method
	constantScopes []constantScope
	methodScopes   []methodScope
	node           *ast.MethodDefinitionNode
}

func (c *Checker) registerMethodCheck(method *types.Method, node *ast.MethodDefinitionNode) {
	if c.IsHeader {
		return
	}

	c.methodChecks.Append(methodCheckEntry{
		method:         method,
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
		node:           node,
		filename:       c.Filename,
	})
}

const concurrencyLimit = 10_000

func (c *Checker) checkMethods() {
	concurrent.Foreach(
		concurrencyLimit,
		c.methodChecks.Slice,
		func(methodCheck methodCheckEntry) {
			methodChecker := c.newMethodChecker(
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
	c.methodChecks.Slice = nil
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

func (c *Checker) declareMethodForSetter(node *ast.AttributeParameterNode, docComment string) {
	setterName := node.Name + "="

	methodScope := c.currentMethodScope()
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

	methodNode := ast.NewMethodDefinitionNode(
		node.Span(),
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
	span *position.Span,
) {
	var areIncompatible bool
	errDetailsBuff := new(strings.Builder)

	if baseMethod.IsSealed() {
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
	typeParamNodes []ast.TypeParameterNode,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	body []ast.StatementNode,
	span *position.Span,
) (ast.TypeNode, ast.TypeNode) {
	name := checkedMethod.Name
	prevMode := c.mode
	isClosure := types.IsClosure(methodNamespace)

	if methodNamespace != nil {
		currentMethod := c.resolveMethodInNamespace(methodNamespace, name)
		if checkedMethod != currentMethod && checkedMethod.IsSealed() {
			c.addOverrideSealedMethodError(checkedMethod, currentMethod.Span())
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

	c.pushIsolatedLocalEnv()
	defer c.popLocalEnv()

	c.mode = prevMode
	c.setInputPositionTypeMode()
	for _, param := range paramNodes {
		switch p := param.(type) {
		case *ast.MethodParameterNode:
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
		case *ast.FormalParameterNode:
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

	if len(body) > 0 && checkedMethod.IsAbstract() {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` cannot have a body because it is abstract",
				name,
			),
			span,
		)
	}

	if isClosure && returnType == nil {
		c.mode = closureInferReturnTypeMode
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
	c.returnType = nil
	c.throwType = nil
	c.mode = prevMode
	return typedReturnTypeNode, typedThrowTypeNode
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
	typeArgs map[value.Symbol]*types.TypeArgument,
) {
	prevMode := c.mode
	c.mode = inferTypeArgumentMode
	defer c.setMode(prevMode)
	reqParamCount := method.RequiredParamCount()
	requiredPosParamCount := len(method.Params) - method.OptionalParamCount
	if method.PostParamCount != -1 {
		requiredPosParamCount -= method.PostParamCount + 1
	}
	positionalRestParamIndex := method.PositionalRestParamIndex()
	var typedPositionalArguments []ast.ExpressionNode

	typeArgMap := make(map[value.Symbol]*types.TypeArgument)
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

		inferredParamType := c.inferTypeArguments(posArgType, param.Type, typeArgMap)
		if inferredParamType == nil {
			param.Type = types.Nothing{}
		} else if inferredParamType != param.Type {
			param.Type = inferredParamType
		}
		typedPositionalArguments = append(typedPositionalArguments, typedPosArg)

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

	if len(typeArgMap) != len(typeParams) {
		c.addFailure(
			fmt.Sprintf(
				"could not infer type parameters in call to `%s`",
				method.Name,
			),
			span,
		)
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
			return nil, nil
		}
		restPositionalArguments := ast.NewArrayTupleLiteralNode(
			span,
			nil,
		)
		posRestParam := method.Params[positionalRestParamIndex]

		currentArgIndex := currentParamIndex
		for ; currentArgIndex < len(positionalArguments)-method.PostParamCount; currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			typedPosArg := c.checkExpressionWithType(posArg, posRestParam.Type)
			posArgType := c.typeOf(typedPosArg)
			inferredParamType := c.inferTypeArguments(posArgType, posRestParam.Type, typeArgMap)
			if inferredParamType == nil {
				posRestParam.Type = types.Nothing{}
			} else if inferredParamType != posRestParam.Type {
				posRestParam.Type = inferredParamType
			}
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, typedPosArg)
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

			typedPosArg := c.checkExpressionWithType(posArg, param.Type)
			posArgType := c.typeOf(typedPosArg)
			inferredParamType := c.inferTypeArguments(posArgType, param.Type, typeArgMap)
			if inferredParamType == nil {
				param.Type = types.Nothing{}
			} else if inferredParamType != param.Type {
				param.Type = inferredParamType
			}
			typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
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

			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, param.Type)
			namedArgType := c.typeOf(typedNamedArgValue)
			inferredParamType := c.inferTypeArguments(namedArgType, param.Type, typeArgMap)
			if inferredParamType == nil {
				param.Type = types.Nothing{}
			} else if inferredParamType != param.Type {
				param.Type = inferredParamType
			}
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

			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, namedRestParam.Type)
			posArgType := c.typeOf(typedNamedArgValue)
			inferredParamType := c.inferTypeArguments(posArgType, namedRestParam.Type, typeArgMap)
			if inferredParamType == nil {
				namedRestParam.Type = types.Nothing{}
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

	return typedPositionalArguments, typeArgMap
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

		typedPosArg := c.checkExpressionWithType(posArg, param.Type)
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
			typedPosArg := c.checkExpressionWithType(posArg, posRestParam.Type)
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

			typedPosArg := c.checkExpressionWithType(posArg, param.Type)
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
			typedNamedArgValue := c.checkExpressionWithType(namedArg.Value, param.Type)
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
	if types.IsNever(receiverType) || types.IsNothing(receiverType) {
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
		return receiver, positionalArgumentNodes, types.Nothing{}
	}

	var typedPositionalArguments []ast.ExpressionNode
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
			return receiver, positionalArgumentNodes, types.Nothing{}
		}

		method = c.replaceTypeParametersInMethod(method, typeArgs.ArgumentMap)
		typedPositionalArguments = c.checkMethodArguments(method, positionalArgumentNodes, namedArgumentNodes, span)
	} else if len(method.TypeParameters) > 0 {
		var typeArgMap map[value.Symbol]*types.TypeArgument
		method = method.DeepCopy()
		typedPositionalArguments, typeArgMap = c.checkMethodArgumentsAndInferTypeArguments(
			method,
			positionalArgumentNodes,
			namedArgumentNodes,
			method.TypeParameters,
			span,
		)
		if typedPositionalArguments == nil {
			return receiver, positionalArgumentNodes, types.Nothing{}
		}
		if len(typeArgMap) != len(method.TypeParameters) {
			return receiver, positionalArgumentNodes, types.Nothing{}
		}
		method.ReturnType = c.replaceTypeParameters(method.ReturnType, typeArgMap)
		method.ThrowType = c.replaceTypeParameters(method.ThrowType, typeArgMap)
	} else {
		typedPositionalArguments = c.checkMethodArguments(method, positionalArgumentNodes, namedArgumentNodes, span)
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

func (c *Checker) checkMethodDefinition(node *ast.MethodDefinitionNode) {
	returnType, throwType := c.checkMethod(
		c.currentMethodScope().container,
		c.typeOf(node).(*types.Method),
		node.TypeParameters,
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

	c.mode = methodMode
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

			t := c.checkTypeParameterNode(node, typeParamMod)
			typeParams = append(typeParams, t)
			typeParamNode.SetType(t)
			typeParamMod.DefineSubtype(t.Name, t)
			typeParamMod.DefineConstant(t.Name, types.NoValue{})
		}
	}

	c.mode = prevMode
	c.setInputPositionTypeMode()
	var params []*types.Parameter
	for i, param := range paramNodes {
		switch p := param.(type) {
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
					param.Span(),
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
			params = append(params, types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			))
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
					param.Span(),
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
			params = append(params, types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			))
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
					param.Span(),
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
	newMethod.SetSpan(span)

	if oldMethod != nil {
		c.checkMethodOverride(
			newMethod,
			oldMethod,
			span,
		)
	}
	methodNamespace.SetMethod(name, newMethod)

	c.mode = prevMode

	return newMethod, typeParamMod
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
	prevMode := c.mode
	c.mode = methodCompatibilityInAlgebraicTypeMode

	areCompatible := c.checkMethodCompatibility(baseMethod, overrideMethod, errSpan)

	c.mode = prevMode

	return areCompatible
}

// Checks whether two methods are compatible.
func (c *Checker) checkMethodCompatibility(baseMethod, overrideMethod *types.Method, errSpan *position.Span) bool {
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

			if oldParam.Name != newParam.Name || oldParam.Kind != newParam.Kind || !c.isSubtype(oldParam.Type, newParam.Type, errSpan) {
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
					param.Name,
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
		currentNamespace := namespace
		var generics []*types.Generic
		seenMethods := make(map[value.Symbol]bool)

		for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
			if generic, ok := currentNamespace.(*types.Generic); ok {
				generics = append(generics, generic)
			}
			for name, method := range currentNamespace.Methods().Map {
				if seenMethods[name] {
					continue
				}
				if len(generics) < 1 {
					if !yield(name, method) {
						return
					}
					seenMethods[name] = true
					continue
				}

				method = method.DeepCopy()
				for i := len(generics) - 1; i >= 0; i-- {
					generic := generics[i]
					method = c.replaceTypeParametersInMethod(method, generic.ArgumentMap)
				}
				if !yield(name, method) {
					return
				}
				seenMethods[name] = true
			}
		}
	}
}

// Iterates over every abstract method of the namespace, resolving type parameters.
func (c *Checker) abstractMethodsInNamespace(namespace types.Namespace) iter.Seq2[value.Symbol, *types.Method] {
	return func(yield func(name value.Symbol, method *types.Method) bool) {
		currentNamespace := namespace
		var generics []*types.Generic
		seenMethods := make(map[value.Symbol]bool)

		for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
			if generic, ok := currentNamespace.(*types.Generic); ok {
				generics = append(generics, generic)
			}
			if !currentNamespace.IsAbstract() {
				continue
			}
			for name, method := range currentNamespace.Methods().Map {
				if !method.IsAbstract() {
					continue
				}
				if seenMethods[name] {
					continue
				}
				if len(generics) < 1 {
					if !yield(name, method) {
						return
					}
					seenMethods[name] = true
					continue
				}

				method = method.DeepCopy()
				for i := len(generics) - 1; i >= 0; i-- {
					generic := generics[i]
					method = c.replaceTypeParametersInMethod(method, generic.ArgumentMap)
				}
				if !yield(name, method) {
					return
				}
				seenMethods[name] = true
			}
		}
	}
}

func (c *Checker) resolveMethodInNamespace(namespace types.Namespace, name value.Symbol) *types.Method {
	currentNamespace := namespace
	var generics []*types.Generic

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		if generic, ok := currentNamespace.(*types.Generic); ok {
			generics = append(generics, generic)
		}
		method := currentNamespace.Method(name)
		if method != nil {
			if len(generics) < 1 {
				return method
			}

			method = method.DeepCopy()
			for i := len(generics) - 1; i >= 0; i-- {
				generic := generics[i]
				method = c.replaceTypeParametersInMethod(method, generic.ArgumentMap)
			}
			return method
		}
	}

	return nil
}

func (c *Checker) resolveNonAbstractMethodInNamespace(namespace types.Namespace, name value.Symbol) *types.Method {
	currentNamespace := namespace
	var generics []*types.Generic

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		if generic, ok := currentNamespace.(*types.Generic); ok {
			generics = append(generics, generic)
		}
		method := currentNamespace.Method(name)
		if method != nil {
			if method.IsAbstract() {
				continue
			}
			if len(generics) < 1 {
				return method
			}

			method = method.DeepCopy()
			for i := len(generics) - 1; i >= 0; i-- {
				generic := generics[i]
				method = c.replaceTypeParametersInMethod(method, generic.ArgumentMap)
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

func (c *Checker) getMethodInNamespaceWithSelf(namespace types.Namespace, typ types.Type, name value.Symbol, self types.Type, errSpan *position.Span, inParent, inSelf bool) *types.Method {
	method := c._getMethodInNamespace(namespace, typ, name, errSpan, inParent)
	if method == nil {
		return nil
	}
	if inSelf {
		return method
	}
	m := map[value.Symbol]*types.TypeArgument{
		symbol.M_self: types.NewTypeArgument(
			self,
			types.INVARIANT,
		),
	}
	return c.replaceTypeParametersInMethod(method.DeepCopy(), m)
}

func (c *Checker) getMethodInNamespace(namespace types.Namespace, typ types.Type, name value.Symbol, errSpan *position.Span, inParent, inSelf bool) *types.Method {
	return c.getMethodInNamespaceWithSelf(namespace, typ, name, namespace, errSpan, inParent, inSelf)
}

func (c *Checker) replaceTypeParametersInMethod(method *types.Method, typeArgs map[value.Symbol]*types.TypeArgument) *types.Method {
	method.ReturnType = c.replaceTypeParameters(method.ReturnType, typeArgs)
	method.ThrowType = c.replaceTypeParameters(method.ThrowType, typeArgs)

	for _, param := range method.Params {
		param.Type = c.replaceTypeParameters(param.Type, typeArgs)
	}
	return method
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
		typeArgMap[symbol.M_self] = types.NewTypeArgument(
			typ,
			types.INVARIANT,
		)
		return c.replaceTypeParametersInMethod(method.DeepCopy(), typeArgMap)
	default:
		return c._getMethod(typ.UpperBound, name, errSpan, inParent, inSelf)
	}
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
		if len(nilMethod.Params) < len(nonNilMethod.Params) {
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
	c.addFailure(
		fmt.Sprintf("method `%s` is not defined on type `%s`", name, types.InspectWithColor(typ)),
		span,
	)
}
