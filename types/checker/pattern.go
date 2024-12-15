package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

func (c *Checker) checkPattern(node ast.PatternNode, matchedType types.Type) (result ast.PatternNode, fullyCapturedType types.Type) {
	switch n := node.(type) {
	case *ast.AsPatternNode:
		return c.checkAsPatternNode(n, matchedType)
	case *ast.PublicIdentifierNode:
		node.SetType(c.checkIdentifierPattern(n.Value, matchedType, matchedType, n.Span()))
		return node, types.Any{}
	case *ast.PrivateIdentifierNode:
		node.SetType(c.checkIdentifierPattern(n.Value, matchedType, matchedType, n.Span()))
		return node, types.Any{}
	case *ast.IntLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.Int64LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.Int32LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.Int16LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.Int8LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.UInt64LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.UInt32LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.UInt16LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.UInt8LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.FloatLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.Float64LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.Float32LiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.BigFloatLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.SimpleSymbolLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.InterpolatedSymbolLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.InterpolatedRegexLiteralNode:
		return c.checkRegexLiteralPattern(n, matchedType)
	case *ast.UninterpolatedRegexLiteralNode:
		return c.checkRegexLiteralPattern(n, matchedType)
	case *ast.DoubleQuotedStringLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.RawStringLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.InterpolatedStringLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.CharLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.RawCharLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.NilLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.TrueLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.FalseLiteralNode:
		return c.checkSimpleLiteralPattern(n, matchedType)
	case *ast.BinArrayListLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdList(), c.Std(symbol.Int)),
			matchedType,
			n.Span(),
		)
	case *ast.HexArrayListLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdList(), c.Std(symbol.Int)),
			matchedType,
			n.Span(),
		)
	case *ast.SymbolArrayListLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdList(), c.Std(symbol.Symbol)),
			matchedType,
			n.Span(),
		)
	case *ast.WordArrayListLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdList(), c.Std(symbol.String)),
			matchedType,
			n.Span(),
		)
	case *ast.BinArrayTupleLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdTuple(), c.Std(symbol.Int)),
			matchedType,
			n.Span(),
		)
	case *ast.HexArrayTupleLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdTuple(), c.Std(symbol.Int)),
			matchedType,
			n.Span(),
		)
	case *ast.SymbolArrayTupleLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdTuple(), c.Std(symbol.Symbol)),
			matchedType,
			n.Span(),
		)
	case *ast.WordArrayTupleLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdTuple(), c.Std(symbol.String)),
			matchedType,
			n.Span(),
		)
	case *ast.BinHashSetLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdSet(), c.Std(symbol.Int)),
			matchedType,
			n.Span(),
		)
	case *ast.HexHashSetLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdSet(), c.Std(symbol.Int)),
			matchedType,
			n.Span(),
		)
	case *ast.SymbolHashSetLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdSet(), c.Std(symbol.Symbol)),
			matchedType,
			n.Span(),
		)
	case *ast.WordHashSetLiteralNode:
		return c.checkSpecialCollectionLiteralPattern(
			n,
			types.NewGenericWithTypeArgs(c.StdSet(), c.Std(symbol.String)),
			matchedType,
			n.Span(),
		)
	case *ast.RangeLiteralNode:
		return c.checkRangePattern(n, matchedType)
	case *ast.MapPatternNode:
		return c.checkMapPattern(n, matchedType)
	case *ast.RecordPatternNode:
		return c.checkRecordPattern(n, matchedType)
	case *ast.RestPatternNode:
		if n.Identifier == nil {
			return n, types.Any{}
		}

		_, typ := c.checkPattern(n.Identifier, types.NewGenericWithTypeArgs(c.StdArrayList(), matchedType))
		return n, typ
	case *ast.ListPatternNode:
		return c.checkListPattern(n, matchedType)
	case *ast.TuplePatternNode:
		return c.checkTuplePattern(n, matchedType)
	case *ast.SetPatternNode:
		return c.checkSetPattern(n, matchedType)
	case *ast.ObjectPatternNode:
		return c.checkObjectPattern(n, matchedType)
	case *ast.PublicConstantNode:
		return c.checkPublicConstantPattern(n, matchedType)
	case *ast.PrivateConstantNode:
		return c.checkPrivateConstantPattern(n, matchedType)
	case *ast.ConstantLookupNode:
		return c.checkConstantLookupPattern(n, matchedType)
	case *ast.UnaryExpressionNode:
		return c.checkUnaryPattern(n, matchedType)
	case *ast.BinaryPatternNode:
		return c.checkBinaryPattern(n, matchedType)
	default:
		panic(fmt.Sprintf("invalid pattern node %T", node))
	}
}

func (c *Checker) checkBinaryPattern(node *ast.BinaryPatternNode, matchedType types.Type) (*ast.BinaryPatternNode, types.Type) {
	switch node.Op.Type {
	case token.OR_OR:
		var leftCatchType, rightCatchType types.Type
		node.Left, leftCatchType = c.checkPattern(node.Left, matchedType)
		leftType := c.TypeOf(node.Left)

		prevMode := c.mode
		switch c.mode {
		case valuePatternMode, nilableValuePatternMode:
			c.mode = nilableValuePatternMode
		default:
			c.mode = nilablePatternMode
		}
		node.Right, rightCatchType = c.checkPattern(node.Right, matchedType)
		rightType := c.TypeOf(node.Right)
		c.mode = prevMode

		node.SetType(c.NewNormalisedUnion(leftType, rightType))
		return node, c.NewNormalisedUnion(leftCatchType, rightCatchType)
	case token.AND_AND:
		var leftCatchType, rightCatchType types.Type
		node.Left, leftCatchType = c.checkPattern(node.Left, matchedType)
		leftType := c.TypeOf(node.Left)

		node.Right, rightCatchType = c.checkPattern(node.Right, matchedType)
		rightType := c.TypeOf(node.Right)

		intersection := c.NewNormalisedIntersection(leftType, rightType)
		if types.IsNever(intersection) {
			c.addWarning(
				"this pattern is impossible to satisfy",
				node.Span(),
			)
		}
		node.SetType(intersection)
		return node, c.NewNormalisedIntersection(leftCatchType, rightCatchType)
	default:
		panic(fmt.Sprintf("invalid binary pattern operator: %s", node.Op.Type.String()))
	}
}

func (c *Checker) checkUnaryPattern(node *ast.UnaryExpressionNode, matchedType types.Type) (ast.PatternNode, types.Type) {
	switch node.Op.Type {
	case token.STRICT_EQUAL:
		node.Right = c.checkExpression(node.Right)
		rightType := c.TypeOf(node.Right)
		c.checkCanMatch(matchedType, rightType, node.Right.Span())
		node.SetType(rightType)
		return node, types.Never{}
	case token.EQUAL_EQUAL:
		node.Right = c.checkExpression(node.Right)
		rightType := c.TypeOf(node.Right)
		c.checkCanMatch(matchedType, rightType, node.Right.Span())
		node.SetType(rightType)
		if rightType.IsLiteral() {
			return node, rightType
		}
		return node, types.Never{}
	case token.NOT_EQUAL, token.STRICT_NOT_EQUAL:
		node.Right = c.checkExpression(node.Right)
		rightType := c.TypeOf(node.Right)
		c.checkCanMatch(matchedType, rightType, node.Right.Span())
		node.SetType(matchedType)
		return node, types.Never{}
	case token.LAX_EQUAL, token.LAX_NOT_EQUAL:
		node.Right = c.checkExpression(node.Right)
		node.SetType(matchedType)
		return node, types.Never{}
	case token.LESS:
		return c.checkRelationalPattern(node, matchedType, symbol.OpLessThan)
	case token.LESS_EQUAL:
		return c.checkRelationalPattern(node, matchedType, symbol.OpLessThanEqual)
	case token.GREATER:
		return c.checkRelationalPattern(node, matchedType, symbol.OpGreaterThan)
	case token.GREATER_EQUAL:
		return c.checkRelationalPattern(node, matchedType, symbol.OpGreaterThanEqual)
	case token.MINUS:
		return c.checkSimpleLiteralPattern(node, matchedType)
	case token.PLUS:
		return c.checkSimpleLiteralPattern(node, matchedType)
	default:
		panic(fmt.Sprintf("invalid unary pattern operator: %s", node.Op.Type.String()))
	}
}

func (c *Checker) checkRelationalPattern(node *ast.UnaryExpressionNode, matchedType types.Type, operator value.Symbol) (*ast.UnaryExpressionNode, types.Type) {
	node.Right = c.checkExpression(node.Right)
	rightType := c.ToNonLiteral(c.TypeOf(node.Right), false)
	if !c.checkCanMatch(matchedType, rightType, node.Right.Span()) {
		node.SetType(types.Untyped{})
		return node, types.Never{}
	}

	intersection := c.NewNormalisedIntersection(rightType, matchedType)
	node.SetType(intersection)
	c.getMethod(intersection, operator, node.Op.Span())
	return node, types.Never{}
}

func (c *Checker) checkAsPatternNode(node *ast.AsPatternNode, typ types.Type) (ast.PatternNode, types.Type) {
	result, fullyCapturedType := c.checkPattern(node.Pattern, typ)
	node.Pattern = result
	patternType := c.TypeOf(node.Pattern)

	switch name := node.Name.(type) {
	case *ast.PublicIdentifierNode:
		node.SetType(c.checkIdentifierPattern(name.Value, typ, patternType, name.Span()))
	case *ast.PrivateIdentifierNode:
		node.SetType(c.checkIdentifierPattern(name.Value, typ, patternType, name.Span()))
	default:
		panic(fmt.Sprintf("invalid identifier node in pattern: %T", node.Name))
	}
	return node, fullyCapturedType
}

func (c *Checker) checkConstantLookupPattern(node *ast.ConstantLookupNode, typ types.Type) (ast.PatternNode, types.Type) {
	n := c.checkConstantLookupNode(node)
	constType := c.typeOfGuardVoid(n)

	c.checkCanMatch(typ, constType, node.Span())
	return n, types.Never{}
}

func (c *Checker) checkPrivateConstantPattern(node *ast.PrivateConstantNode, typ types.Type) (*ast.PrivateConstantNode, types.Type) {
	node = c.checkPrivateConstantNode(node)
	constType := c.typeOfGuardVoid(node)

	c.checkCanMatch(typ, constType, node.Span())
	return node, types.Never{}
}

func (c *Checker) checkPublicConstantPattern(node *ast.PublicConstantNode, typ types.Type) (*ast.PublicConstantNode, types.Type) {
	node = c.checkPublicConstantNode(node)
	constType := c.typeOfGuardVoid(node)

	c.checkCanMatch(typ, constType, node.Span())
	return node, types.Never{}
}

func (c *Checker) checkObjectPattern(node *ast.ObjectPatternNode, typ types.Type) (resultNode *ast.ObjectPatternNode, fullyCapturedType types.Type) {
	constType, fullName := c.resolveConstantType(node.ObjectType)
	if constType == nil {
		constType = types.Untyped{}
	}

	node.ObjectType = ast.NewPublicConstantNode(
		node.ObjectType.Span(),
		fullName,
	)

	var classOrMixin types.Namespace
	switch t := constType.(type) {
	case *types.Class:
		classOrMixin = t
	case *types.Mixin:
		classOrMixin = t
	case types.Untyped:
		node.SetType(types.Untyped{})
		return node, types.Never{}
	default:
		c.addFailure(
			fmt.Sprintf(
				"type `%s` cannot be used in object patterns, only classes and mixins are allowed",
				types.InspectWithColor(constType),
			),
			node.Span(),
		)
		node.SetType(types.Untyped{})
		return node, types.Never{}
	}

	var ofAny *types.Generic
	if classOrMixin.IsGeneric() {
		ofAny = types.NewGenericWithUpperBoundTypeArgsAndVariance(classOrMixin, types.COVARIANT)
		extractedNamespace, typeArgs := c.extractTypeArgumentsFromType(classOrMixin, ofAny, typ)
		if typeArgs == nil {
			typeParams := classOrMixin.TypeParameters()
			typeArgMap := make(types.TypeArgumentMap, len(typeParams))
			if c.checkCanMatchWithTypeArgs(typ, classOrMixin, node.Span(), typeArgMap) && typeArgMap.HasAllTypeParams(typeParams) {
				newClassOrMixin := types.NewGeneric(
					classOrMixin,
					types.NewTypeArguments(
						typeArgMap,
						types.CreateTypeArgumentOrderFromTypeParams(typeParams),
					),
				)
				classOrMixin = newClassOrMixin
				node.SetType(classOrMixin)
			} else {
				node.SetType(ofAny)
				classOrMixin = ofAny
			}
		} else {
			classOrMixin = types.NewGeneric(
				classOrMixin,
				typeArgs,
			)
			node.SetType(extractedNamespace)
			c.checkCanMatch(typ, classOrMixin, node.Span())
		}
	} else {
		node.SetType(classOrMixin)
		c.checkCanMatch(typ, classOrMixin, node.Span())
	}

	ofAny.FixVariance()

	allAttributesFullyCaptured := true
	for _, attribute := range node.Attributes {
		switch attr := attribute.(type) {
		case *ast.SymbolKeyValuePatternNode:
			attrType, fullyCaptured := c.checkObjectKeyValuePattern(classOrMixin, attr)
			attr.SetType(attrType)
			if !fullyCaptured {
				allAttributesFullyCaptured = false
			}
		case *ast.PublicIdentifierNode:
			attrType, fullyCaptured := c.checkObjectIdentifierPattern(classOrMixin, attr.Value, attr.Span())
			attr.SetType(attrType)
			if !fullyCaptured {
				allAttributesFullyCaptured = false
			}
		case *ast.PrivateIdentifierNode:
			attrType, fullyCaptured := c.checkObjectIdentifierPattern(classOrMixin, attr.Value, attr.Span())
			attr.SetType(attrType)
			if !fullyCaptured {
				allAttributesFullyCaptured = false
			}
		default:
			panic(fmt.Sprintf("invalid object pattern attribute: %T", attr))
		}
	}

	if allAttributesFullyCaptured {
		return node, c.TypeOf(node)
	}

	return node, types.Never{}
}

func (c *Checker) checkObjectKeyValuePattern(namespace types.Namespace, node *ast.SymbolKeyValuePatternNode) (attrType types.Type, fullyCaptured bool) {
	getter := c.getMethod(namespace, value.ToSymbol(node.Key), node.Span())
	if getter == nil {
		c.checkPattern(node.Value, types.Untyped{})
		return types.Untyped{}, false
	}
	getter, _ = c.checkMethodArguments(getter, nil, nil, nil, node.Span())
	if getter == nil {
		c.checkPattern(node.Value, types.Untyped{})
		return types.Untyped{}, false
	}
	returnType := c.typeGuardVoid(getter.ReturnType, node.Span())

	var fullyCapturedType types.Type
	node.Value, fullyCapturedType = c.checkPattern(node.Value, returnType)
	return returnType, c.IsSubtype(returnType, fullyCapturedType, nil)
}

func (c *Checker) checkObjectIdentifierPattern(namespace types.Namespace, name string, span *position.Span) (attrType types.Type, fullyCaptured bool) {
	getter := c.getMethod(namespace, value.ToSymbol(name), span)
	if getter == nil {
		c.checkIdentifierPattern(name, types.Untyped{}, types.Untyped{}, span)
		return types.Untyped{}, false
	}
	getter, _ = c.checkMethodArguments(getter, nil, nil, nil, span)
	if getter == nil {
		c.checkIdentifierPattern(name, types.Untyped{}, types.Untyped{}, span)
		return types.Untyped{}, false
	}

	returnType := c.typeGuardVoid(getter.ReturnType, span)
	return c.checkIdentifierPattern(name, returnType, returnType, span), true
}

func (c *Checker) checkMapPattern(node *ast.MapPatternNode, typ types.Type) (*ast.MapPatternNode, types.Type) {
	mapMixin := c.Std(symbol.Map).(*types.Mixin)
	mapOfAny := types.NewGenericWithVariance(mapMixin, types.COVARIANT, types.Any{}, types.Any{})

	var keyType types.Type
	var valueType types.Type

	if c.checkCanMatch(typ, mapOfAny, node.Span()) {
		var extractedRecord types.Type
		extractedRecord, keyType, valueType = c.extractRecordElementFromType(mapMixin, mapOfAny, typ)
		node.SetType(extractedRecord)
	} else {
		keyType = types.Any{}
		valueType = types.Any{}
		node.SetType(types.Never{})
	}

	mapOfAny.FixVariance()

	for i, element := range node.Elements {
		switch e := element.(type) {
		case *ast.PublicIdentifierNode:
			c.checkCanMatch(keyType, c.Std(symbol.Symbol), e.Span())
			newE, _ := c.checkPattern(e, valueType)
			node.Elements[i] = newE
		case *ast.PrivateIdentifierNode:
			c.checkCanMatch(keyType, c.Std(symbol.Symbol), e.Span())
			newE, _ := c.checkPattern(e, valueType)
			node.Elements[i] = newE
		case *ast.KeyValuePatternNode:
			e.Key = c.checkExpression(e.Key).(ast.PatternExpressionNode)
			patternKeyType := c.TypeOf(e.Key)
			c.checkCanMatch(keyType, patternKeyType, e.Span())
			e.Value, _ = c.checkPattern(e.Value, valueType)
		case *ast.SymbolKeyValuePatternNode:
			c.checkCanMatch(keyType, c.Std(symbol.Symbol), e.Span())
			e.Value, _ = c.checkPattern(e.Value, valueType)
		default:
			panic(fmt.Sprintf("invalid map pattern element: %T", element))
		}
	}
	return node, types.Never{}
}

func (c *Checker) checkRecordPattern(node *ast.RecordPatternNode, typ types.Type) (*ast.RecordPatternNode, types.Type) {
	recordMixin := c.Std(symbol.Record).(*types.Mixin)
	recordOfAny := types.NewGenericWithVariance(recordMixin, types.COVARIANT, types.Any{}, types.Any{})

	var keyType types.Type
	var valueType types.Type

	if c.checkCanMatch(typ, recordOfAny, node.Span()) {
		var extractedRecord types.Type
		extractedRecord, keyType, valueType = c.extractRecordElementFromType(recordMixin, recordOfAny, typ)
		node.SetType(extractedRecord)
	} else {
		keyType = types.Any{}
		valueType = types.Any{}
		node.SetType(types.Never{})
	}

	recordOfAny.FixVariance()

	for i, element := range node.Elements {
		switch e := element.(type) {
		case *ast.PublicIdentifierNode:
			c.checkCanMatch(keyType, c.Std(symbol.Symbol), e.Span())
			node.Elements[i], _ = c.checkPattern(e, valueType)
		case *ast.PrivateIdentifierNode:
			c.checkCanMatch(keyType, c.Std(symbol.Symbol), e.Span())
			node.Elements[i], _ = c.checkPattern(e, valueType)
		case *ast.KeyValuePatternNode:
			e.Key = c.checkExpression(e.Key).(ast.PatternExpressionNode)
			patternKeyType := c.TypeOf(e.Key)
			c.checkCanMatch(keyType, patternKeyType, e.Span())
			e.Value, _ = c.checkPattern(e.Value, valueType)
		case *ast.SymbolKeyValuePatternNode:
			c.checkCanMatch(keyType, c.Std(symbol.Symbol), e.Span())
			e.Value, _ = c.checkPattern(e.Value, valueType)
		default:
			panic(fmt.Sprintf("invalid record pattern element: %T", element))
		}
	}

	return node, types.Never{}
}

func (c *Checker) checkRangePattern(node *ast.RangeLiteralNode, typ types.Type) (*ast.RangeLiteralNode, types.Type) {
	var startType, endType types.Type
	if node.Start != nil {
		node.Start = c.checkExpression(node.Start)
		startType = c.ToNonLiteral(c.TypeOf(node.Start), false)
		if _, ok := startType.(*types.Class); !ok {
			c.addFailure(
				fmt.Sprintf(
					"type `%s` cannot be used in a range pattern, only class instance types are permitted",
					types.InspectWithColor(startType),
				),
				node.Start.Span(),
			)
		}
	}
	if node.End != nil {
		node.End = c.checkExpression(node.End)
		endType = c.ToNonLiteral(c.TypeOf(node.End), false)
		if _, ok := endType.(*types.Class); !ok {
			c.addFailure(
				fmt.Sprintf(
					"type `%s` cannot be used in a range pattern, only class instance types are permitted",
					types.InspectWithColor(endType),
				),
				node.End.Span(),
			)
		}
	}

	if startType != nil && endType != nil && !c.IsTheSameType(startType, endType, nil) {
		c.addFailure(
			fmt.Sprintf(
				"range pattern start and end must be of the same type, got `%s` and `%s`",
				types.InspectWithColor(startType),
				types.InspectWithColor(endType),
			),
			node.Span(),
		)
	}

	c.checkCanMatch(typ, startType, node.Span())
	node.SetType(startType)
	return node, types.Never{}
}

func (c *Checker) checkSpecialCollectionLiteralPattern(node ast.PatternExpressionNode, patternType, typ types.Type, span *position.Span) (ast.PatternNode, types.Type) {
	c.checkCanMatch(typ, patternType, span)
	node.SetType(patternType)
	return node, types.Never{}
}

func (c *Checker) checkRegexLiteralPattern(node ast.RegexLiteralNode, typ types.Type) (ast.PatternNode, types.Type) {
	c.checkExpression(node)
	nodeType := c.StdString()
	node.SetType(nodeType)
	c.checkCanMatch(typ, nodeType, node.Span())
	return node, types.Never{}
}

func (c *Checker) checkSimpleLiteralPattern(node ast.PatternExpressionNode, typ types.Type) (ast.PatternNode, types.Type) {
	n := c.checkExpression(node)
	nodeType := c.TypeOf(n)
	c.checkCanMatch(typ, nodeType, n.Span())
	return n.(ast.PatternNode), nodeType
}

func (c *Checker) addCannotMatchError(assignedType types.Type, targetType types.Type, span *position.Span) {
	c.addFailure(
		fmt.Sprintf(
			"type `%s` cannot ever match type `%s`",
			types.InspectWithColor(assignedType),
			types.InspectWithColor(targetType),
		),
		span,
	)
}

func (c *Checker) checkCanMatch(assignedType types.Type, targetType types.Type, span *position.Span) bool {
	if !c.TypesIntersect(assignedType, targetType) {
		c.addCannotMatchError(assignedType, targetType, span)
		return false
	}

	return true
}

func (c *Checker) checkCanMatchWithTypeArgs(assignedType types.Type, targetType types.Type, span *position.Span, typeArgs types.TypeArgumentMap) bool {
	if !c.typesIntersectWithTypeArgs(assignedType, targetType, typeArgs) {
		c.addCannotMatchError(assignedType, targetType, span)
		return false
	}

	return true
}

func (c *Checker) checkTuplePattern(node *ast.TuplePatternNode, typ types.Type) (*ast.TuplePatternNode, types.Type) {
	tupleMixin := c.Std(symbol.Tuple).(*types.Mixin)
	tupleOfAny := types.NewGenericWithTypeArgs(tupleMixin, types.Any{})

	var elementType types.Type

	if c.checkCanMatch(typ, tupleOfAny, node.Span()) {
		var extractedCollection types.Type
		extractedCollection, elementType = c.extractCollectionElementFromType(tupleMixin, tupleOfAny, typ)
		node.SetType(extractedCollection)
	} else {
		elementType = types.Any{}
		node.SetType(types.Never{})
	}

	for i, element := range node.Elements {
		node.Elements[i], _ = c.checkPattern(element, elementType)
	}
	return node, types.Never{}
}

func (c *Checker) checkSetPattern(node *ast.SetPatternNode, typ types.Type) (*ast.SetPatternNode, types.Type) {
	setMixin := c.Std(symbol.Set).(*types.Mixin)
	setOfAny := types.NewGenericWithVariance(setMixin, types.BIVARIANT, types.Any{})

	var elementType types.Type

	if c.checkCanMatch(typ, setOfAny, node.Span()) {
		var extractedCollection types.Type
		extractedCollection, elementType = c.extractCollectionElementFromType(setMixin, setOfAny, typ)
		node.SetType(extractedCollection)
	} else {
		elementType = types.Any{}
		node.SetType(types.Never{})
	}

	for i, element := range node.Elements {
		node.Elements[i], _ = c.checkPattern(element, elementType)
	}

	return node, types.Never{}
}

func (c *Checker) checkListPattern(node *ast.ListPatternNode, typ types.Type) (*ast.ListPatternNode, types.Type) {
	listMixin := c.Std(symbol.List).(*types.Mixin)
	listOfAny := types.NewGenericWithVariance(listMixin, types.COVARIANT, types.Any{})

	var elementType types.Type

	if c.checkCanMatch(typ, listOfAny, node.Span()) {
		var extractedCollection types.Type
		extractedCollection, elementType = c.extractCollectionElementFromType(listMixin, listOfAny, typ)
		node.SetType(extractedCollection)
	} else {
		elementType = types.Any{}
		node.SetType(types.Never{})
	}
	listOfAny.FixVariance()

	for i, element := range node.Elements {
		node.Elements[i], _ = c.checkPattern(element, elementType)
	}
	return node, types.Never{}
}

func (c *Checker) checkIdentifierPattern(name string, valueType, patternType types.Type, span *position.Span) types.Type {
	variable := c.getLocal(name)
	if variable == nil {
		var local *local
		switch c.mode {
		case valuePatternMode:
			varType := c.ToNonLiteral(patternType, false)
			local = newLocal(varType, true, true)
		case nilableValuePatternMode:
			varType := c.ToNilable(c.ToNonLiteral(patternType, false))
			local = newLocal(varType, true, true)
		case nilablePatternMode:
			varType := c.ToNilable(c.ToNonLiteral(patternType, false))
			local = newLocal(varType, true, false)
		default:
			varType := c.ToNonLiteral(patternType, false)
			local = newLocal(varType, true, false)
		}
		c.addLocal(name, local)
		return patternType
	}

	variable.initialised = true
	if variable.singleAssignment {
		c.addValueReassignedError(name, span)
		return variable.typ
	}
	c.checkCanAssign(valueType, variable.typ, span)
	return variable.typ
}
