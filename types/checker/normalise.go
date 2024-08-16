// Package checker implements the Elk type checker
package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

func (c *Checker) inferTypeArguments(givenType, paramType types.Type, typeArgMap map[value.Symbol]*types.TypeArgument) types.Type {
	switch p := paramType.(type) {
	case types.Self:
		arg := typeArgMap[symbol.M_self]
		if arg == nil {
			return p
		}
		return arg.Type
	case *types.Closure:
		g, ok := givenType.(*types.Closure)
		if !ok {
			return p
		}

		gMethod := &g.Body
		pMethod := &p.Body
		var isDifferent bool
		newParams := make([]*types.Parameter, len(pMethod.Params))
		for i := range pMethod.Params {
			pParam := pMethod.Params[i]
			gParam := gMethod.Params[i]
			if pParam.Kind != gParam.Kind || pParam.Name != gParam.Name {
				return p
			}
			result := c.inferTypeArguments(gParam.Type, pParam.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if result != pParam.Type {
				isDifferent = true
				newParam := pParam.Copy()
				newParam.Type = result
				newParams[i] = newParam
			} else {
				newParams[i] = pParam
			}
		}

		returnType := c.inferTypeArguments(gMethod.ReturnType, pMethod.ReturnType, typeArgMap)
		if returnType == nil {
			return nil
		}
		if returnType != pMethod.ReturnType {
			isDifferent = true
		}

		throwType := c.inferTypeArguments(gMethod.ThrowType, pMethod.ThrowType, typeArgMap)
		if throwType == nil {
			return nil
		}
		if throwType != pMethod.ThrowType {
			isDifferent = true
		}

		if isDifferent {
			closure := types.NewClosure(types.Method{})
			newMethod := types.NewMethod(
				pMethod.DocComment,
				pMethod.IsAbstract(),
				pMethod.IsSealed(),
				pMethod.IsNative(),
				pMethod.Name,
				newParams,
				returnType,
				throwType,
				closure,
			)
			closure.Body = *newMethod
			return closure
		}
		return p
	case *types.TypeParameter:
		typeArg := typeArgMap[p.Name]
		if typeArg != nil {
			return typeArg.Type
		}

		nonLiteral := c.toNonLiteral(givenType, false)
		if !c.isSubtype(nonLiteral, p.UpperBound, nil) {
			return nil
		}
		if !c.isSubtype(p.LowerBound, nonLiteral, nil) {
			return nil
		}
		typeArgMap[p.Name] = types.NewTypeArgument(
			nonLiteral,
			p.Variance,
		)
		return nonLiteral
	case *types.Generic:
		g, ok := givenType.(*types.Generic)
		if !ok {
			return nil
		}
		if !c.isSubtype(g.Type, p.Type, nil) {
			return nil
		}
		if len(g.ArgumentOrder) < len(p.ArgumentOrder) {
			return nil
		}

		newArgMap := make(map[value.Symbol]*types.TypeArgument, len(p.ArgumentMap))
		for _, argName := range p.ArgumentOrder {
			pArg := p.ArgumentMap[argName]
			gArg := g.ArgumentMap[argName]
			result := c.inferTypeArguments(gArg.Type, pArg.Type, typeArgMap)
			if result == nil {
				return nil
			}
			newArgMap[argName] = types.NewTypeArgument(result, gArg.Variance)
		}
		return types.NewGeneric(
			p.Type,
			types.NewTypeArguments(
				newArgMap,
				p.ArgumentOrder,
			),
		)
	case *types.SingletonOf:
		switch g := givenType.(type) {
		case *types.SingletonClass:
			result := c.inferTypeArguments(g.AttachedObject, p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			return types.NewSingletonOf(result)
		case *types.SingletonOf:
			result := c.inferTypeArguments(g.Type, p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			return types.NewSingletonOf(result)
		default:
			return nil
		}
	case *types.SingletonClass:
		switch g := givenType.(type) {
		case *types.SingletonClass:
			result := c.inferTypeArguments(g.AttachedObject, p.AttachedObject, typeArgMap)
			if result == nil {
				return nil
			}
			if p.AttachedObject == result {
				return p
			}

			return types.NewSingletonClass(result.(types.Namespace), p.Parent())
		case *types.SingletonOf:
			result := c.inferTypeArguments(g.Type, p.AttachedObject, typeArgMap)
			if result == nil {
				return nil
			}
			if p.AttachedObject == result {
				return p
			}

			return types.NewSingletonClass(result.(types.Namespace), p.Parent())
		default:
			return nil
		}
	case *types.InstanceOf:
		nonLiteral := c.toNonLiteral(givenType, false)
		switch g := nonLiteral.(type) {
		case *types.InstanceOf:
			result := c.inferTypeArguments(g.Type, p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			switch r := result.(type) {
			case *types.SingletonClass:
				return r.AttachedObject
			case *types.SingletonOf:
				return r.Type
			}
			return nil
		case *types.Class:
			result := c.inferTypeArguments(g.Singleton(), p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			switch r := result.(type) {
			case *types.SingletonClass:
				return r.AttachedObject
			case *types.SingletonOf:
				return r.Type
			}
			return nil
		case *types.Mixin:
			result := c.inferTypeArguments(g.Singleton(), p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			switch r := result.(type) {
			case *types.SingletonClass:
				return r.AttachedObject
			case *types.SingletonOf:
				return r.Type
			}
			return nil
		case *types.Interface:
			result := c.inferTypeArguments(g.Singleton(), p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			switch r := result.(type) {
			case *types.SingletonClass:
				return r.AttachedObject
			case *types.SingletonOf:
				return r.Type
			}
			return nil
		default:
			return nil
		}
	case *types.Not:
		g, ok := givenType.(*types.Not)
		if !ok {
			return nil
		}

		result := c.inferTypeArguments(g.Type, p.Type, typeArgMap)
		if result == nil {
			return nil
		}
		if p.Type == result {
			return p
		}

		return types.NewNot(result)
	case *types.Intersection:
		switch g := givenType.(type) {
		case *types.Intersection:
			gElementsToSkip := make([]bool, len(g.Elements))
			for _, pElement := range p.Elements {
				for j, gElement := range g.Elements {
					if c.isSubtype(gElement, pElement, nil) {
						gElementsToSkip[j] = true
						break
					}
				}
			}

			newGElements := make([]types.Type, 0, len(g.Elements))
			for j, gElement := range g.Elements {
				if gElementsToSkip[j] {
					continue
				}
				newGElements = append(newGElements, gElement)
			}
			var newG types.Type
			switch len(newGElements) {
			case 0:
				return p
			case 1:
				newG = newGElements[0]
			default:
				newG = types.NewIntersection(newGElements...)
			}

			newPElements := make([]types.Type, 0, len(p.Elements))
			var isDifferent bool
			for _, pElement := range p.Elements {
				result := c.inferTypeArguments(newG, pElement, typeArgMap)
				if result == nil {
					return nil
				}
				if result != pElement {
					isDifferent = true
				}
				newPElements = append(newPElements, result)
			}
			if isDifferent {
				return types.NewIntersection(newPElements...)
			}
			return p
		default:
			newElements := make([]types.Type, 0, len(p.Elements))
			var isDifferent bool
			for _, pElement := range p.Elements {
				result := c.inferTypeArguments(g, pElement, typeArgMap)
				if result == nil {
					return nil
				}
				if result != pElement {
					isDifferent = true
				}
				newElements = append(newElements, result)
			}

			if isDifferent {
				return types.NewIntersection(newElements...)
			}
			return p
		}
	case *types.Union:
		switch g := givenType.(type) {
		case *types.Union:
			narrowedGivenElements := make([]types.Type, 0, len(g.Elements))
			for _, gElement := range g.Elements {
				if c.isSubtype(gElement, p, nil) {
					continue
				}
				narrowedGivenElements = append(narrowedGivenElements, gElement)
			}
			if len(narrowedGivenElements) == 0 {
				return p
			}
			var narrowedG types.Type
			if len(narrowedGivenElements) == 1 {
				narrowedG = narrowedGivenElements[0]
			} else {
				narrowedG = types.NewUnion(narrowedGivenElements...)
			}

			var isDifferent bool
			newPElements := make([]types.Type, 0, len(p.Elements))
			for _, pElement := range p.Elements {
				result := c.inferTypeArguments(narrowedG, pElement, typeArgMap)
				if result == nil {
					return nil
				}
				if result != pElement {
					isDifferent = true
				}

				newPElements = append(newPElements, result)
			}
			if !isDifferent {
				return p
			}

			return types.NewUnion(newPElements...)
		case *types.Nilable:
			return c.inferTypeArguments(types.NewUnion(types.Nil{}, g.Type), p, typeArgMap)
		default:
			newElements := make([]types.Type, 0, len(p.Elements))
			var isDifferent bool
			for _, pElement := range p.Elements {
				result := c.inferTypeArguments(g, pElement, typeArgMap)
				if result == nil {
					return nil
				}
				if result != pElement {
					isDifferent = true
				}
				newElements = append(newElements, result)
			}

			if isDifferent {
				return types.NewUnion(newElements...)
			}
			return p
		}
	case *types.Nilable:
		switch g := givenType.(type) {
		case *types.Nilable:
			result := c.inferTypeArguments(g.Type, p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			return types.NewNilable(result)
		case *types.Union:
			var withoutNil []types.Type
			for _, element := range g.Elements {
				switch e := element.(type) {
				case types.Nil:
					continue
				case *types.Class:
					if e.Name() == "Std::Nil" {
						continue
					}
					withoutNil = append(withoutNil, e)
				default:
					withoutNil = append(withoutNil, e)
				}
			}
			var t types.Type
			if len(withoutNil) == len(g.Elements) {
				t = g
			} else if len(withoutNil) == 0 {
				t = types.Never{}
			} else if len(withoutNil) == 1 {
				t = withoutNil[0]
			} else {
				t = types.NewUnion(withoutNil...)
			}

			result := c.inferTypeArguments(t, p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			return types.NewNilable(result)
		default:
			result := c.inferTypeArguments(givenType, p.Type, typeArgMap)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}
			return types.NewNilable(result)
		}
	default:
		return paramType
	}
}

func (c *Checker) replaceTypeParameters(typ types.Type, typeArgMap map[value.Symbol]*types.TypeArgument) types.Type {
	return c.normaliseType(c._replaceTypeParameters(typ, typeArgMap))
}

func (c *Checker) _replaceTypeParameters(typ types.Type, typeArgMap map[value.Symbol]*types.TypeArgument) types.Type {
	switch t := typ.(type) {
	case types.Self:
		arg := typeArgMap[symbol.M_self]
		if arg == nil {
			return t
		}
		return arg.Type
	case *types.SingletonOf:
		return types.NewSingletonOf(
			c._replaceTypeParameters(t.Type, typeArgMap),
		)
	case *types.InstanceOf:
		return types.NewInstanceOf(
			c._replaceTypeParameters(t.Type, typeArgMap),
		)
	case *types.Closure:
		method := (&t.Body).Copy()
		for _, param := range method.Params {
			param.Type = c._replaceTypeParameters(param.Type, typeArgMap)
		}

		method.ReturnType = c._replaceTypeParameters(method.ReturnType, typeArgMap)
		method.ThrowType = c._replaceTypeParameters(method.ThrowType, typeArgMap)
		closure := types.NewClosure(*method)
		method.DefinedUnder = closure
		return closure
	case *types.Generic:
		newMap := make(map[value.Symbol]*types.TypeArgument, len(t.ArgumentMap))
		for key, arg := range t.ArgumentMap {
			newMap[key] = types.NewTypeArgument(
				c._replaceTypeParameters(arg.Type, typeArgMap),
				arg.Variance,
			)
		}
		return types.NewGeneric(
			c._replaceTypeParameters(t.Type, typeArgMap),
			types.NewTypeArguments(
				newMap,
				t.ArgumentOrder,
			),
		)
	case *types.TypeParameter:
		arg := typeArgMap[t.Name]
		if arg == nil {
			panic(fmt.Sprintf("invalid generic type parameter `%s`", types.InspectWithColor(t)))
		}
		return arg.Type
	case *types.Nilable:
		return types.NewNilable(c._replaceTypeParameters(t.Type, typeArgMap))
	case *types.Not:
		return types.NewNot(c._replaceTypeParameters(t.Type, typeArgMap))
	case *types.Union:
		newElements := make([]types.Type, 0, len(t.Elements))
		for _, element := range t.Elements {
			newElements = append(newElements, c._replaceTypeParameters(element, typeArgMap))
		}
		return types.NewUnion(newElements...)
	case *types.Intersection:
		newElements := make([]types.Type, 0, len(t.Elements))
		for _, element := range t.Elements {
			newElements = append(newElements, c._replaceTypeParameters(element, typeArgMap))
		}
		return types.NewIntersection(newElements...)
	default:
		return t
	}
}

func (c *Checker) normaliseType(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Union:
		return c.newNormalisedUnion(t.Elements...)
	case *types.Intersection:
		return c.newNormalisedIntersection(t.Elements...)
	case *types.SingletonOf:
		switch nestedType := t.Type.(type) {
		case *types.InstanceOf:
			return nestedType.Type
		case *types.Class:
			return nestedType.Singleton()
		case *types.Mixin:
			return nestedType.Singleton()
		case *types.Interface:
			return nestedType.Singleton()
		default:
			return t
		}
	case *types.InstanceOf:
		switch nestedType := t.Type.(type) {
		case *types.SingletonOf:
			return nestedType.Type
		case *types.SingletonClass:
			return nestedType.AttachedObject
		default:
			return t
		}
	case *types.Nilable:
		t.Type = c.normaliseType(t.Type)
		switch t.Type.(type) {
		case types.Never:
			return types.Nil{}
		case types.Any, types.Nothing:
			return t.Type
		}
		if c.isNilable(t.Type) {
			return t.Type
		}
		if union, ok := t.Type.(*types.Union); ok {
			union.Elements = append(union.Elements, types.Nil{})
			return union
		}
		return t
	case *types.Not:
		t.Type = c.normaliseType(t.Type)
		switch nestedType := t.Type.(type) {
		case *types.Not:
			return nestedType.Type
		case types.Never:
			return types.Any{}
		case types.Any:
			return types.Never{}
		case types.Nothing:
			return types.Nothing{}
		case *types.Union:
			intersectionElements := make([]types.Type, 0, len(nestedType.Elements))
			for _, element := range nestedType.Elements {
				intersectionElements = append(intersectionElements, types.NewNot(element))
			}
			return c.newNormalisedIntersection(intersectionElements...)
		case *types.Intersection:
			unionElements := make([]types.Type, 0, len(nestedType.Elements))
			for _, element := range nestedType.Elements {
				unionElements = append(unionElements, types.NewNot(element))
			}
			return c.newNormalisedUnion(unionElements...)
		}

		return t
	default:
		return typ
	}
}

func (c *Checker) distributeIntersectionOverUnions(newUnionElements *[]types.Type, intersectionElements []types.Type, i int) {
	if i == len(intersectionElements) {
		*newUnionElements = append(*newUnionElements, types.NewIntersection(intersectionElements...))
		return
	}

	intersectionElement := intersectionElements[i]
	switch e := intersectionElement.(type) {
	case *types.Union:
		for _, subUnionElement := range e.Elements {
			newIntersectionElements := make([]types.Type, 0, len(intersectionElements)+1)
			newIntersectionElements = append(newIntersectionElements, intersectionElements[:i]...)
			newIntersectionElements = append(newIntersectionElements, subUnionElement)
			if len(intersectionElements) >= i+2 {
				newIntersectionElements = append(newIntersectionElements, intersectionElements[i+1:]...)
			}
			c.distributeIntersectionOverUnions(newUnionElements, newIntersectionElements, i+1)
		}
	case *types.Nilable:
		elements := []types.Type{e.Type, types.Nil{}}
		for _, subUnionElement := range elements {
			newIntersectionElements := make([]types.Type, 0, len(intersectionElements)+1)
			newIntersectionElements = append(newIntersectionElements, intersectionElements[:i]...)
			newIntersectionElements = append(newIntersectionElements, subUnionElement)
			if len(intersectionElements) >= i+2 {
				newIntersectionElements = append(newIntersectionElements, intersectionElements[i+1:]...)
			}
			c.distributeIntersectionOverUnions(newUnionElements, newIntersectionElements, i+1)
		}
	default:
		c.distributeIntersectionOverUnions(newUnionElements, intersectionElements, i+1)
	}
}

// Transform an intersection of unions to a unions of intersections.
// String & (Int | Float) => (String & Int) | (String & Float)
func (c *Checker) intersectionOfUnionsToUnionOfIntersections(intersectionElements []types.Type) types.Type {
	newUnionElements := new([]types.Type)
	c.distributeIntersectionOverUnions(newUnionElements, intersectionElements, 0)
	if len(*newUnionElements) == 0 {
		return types.Never{}
	}
	if len(*newUnionElements) == 1 {
		return (*newUnionElements)[0]
	}
	return types.NewUnion(*newUnionElements...)
}

func (c *Checker) newNormalisedIntersection(elements ...types.Type) types.Type {
	var containsNot bool
	var containsUninitialisedNamedTypes bool

	for i := 0; i < len(elements); i++ {
		element := c.normaliseType(elements[i])
		if types.IsNever(element) || types.IsNothing(element) {
			return element
		}
		switch e := element.(type) {
		case *types.Intersection:
			newElements := make([]types.Type, 0, len(elements)+len(e.Elements))
			newElements = append(newElements, elements[:i]...)
			newElements = append(newElements, e.Elements...)
			if len(elements) >= i+2 {
				newElements = append(newElements, elements[i+1:]...)
			}
			elements = newElements
			i--
		case *types.Not:
			containsNot = true
		case *types.NamedType:
			if e.Type == nil {
				containsUninitialisedNamedTypes = true
			}
		}
	}
	if containsUninitialisedNamedTypes {
		return types.NewIntersection(elements...)
	}
	if containsNot {
		// expand named types
		for i := 0; i < len(elements); i++ {
			switch e := elements[i].(type) {
			case *types.Intersection:
				newElements := make([]types.Type, 0, len(elements)+len(e.Elements))
				newElements = append(newElements, elements[:i]...)
				newElements = append(newElements, e.Elements...)
				if len(elements) >= i+2 {
					newElements = append(newElements, elements[i+1:]...)
				}
				elements = newElements
				i--
			case *types.NamedType:
				elements[i] = e.Type
				i--
			case types.Bool:
				elements[i] = types.NewUnion(types.True{}, types.False{})
			}
		}
	}
	distributedIntersection := c.intersectionOfUnionsToUnionOfIntersections(elements)
	intersection, ok := distributedIntersection.(*types.Intersection)
	if !ok {
		return c.normaliseType(distributedIntersection)
	}

	elements = intersection.Elements
	normalisedElements := make([]types.Type, 0, len(elements))

	// detect empty intersections
	for _, element := range elements {
		if types.IsNever(element) || types.IsNothing(element) {
			return element
		}

		for _, normalisedElement := range normalisedElements {
			if !c.canIntersect(element, normalisedElement) {
				return types.Never{}
			}
		}
		normalisedElements = append(normalisedElements, element)
	}

	elements = normalisedElements
	normalisedElements = make([]types.Type, 0, len(elements))

eliminateSupertypesLoop:
	for i := 0; i < len(elements); i++ {
		element := elements[i]

		for j := 0; j < len(normalisedElements); j++ {
			normalisedElement := normalisedElements[j]
			if c.isSubtype(normalisedElement, element, nil) {
				continue eliminateSupertypesLoop
			}
			if c.isSubtype(element, normalisedElement, nil) {
				normalisedElements[j] = element
				continue eliminateSupertypesLoop
			}

		}
		normalisedElements = append(normalisedElements, element)
	}

	if len(normalisedElements) == 0 {
		return types.Never{}
	}
	if len(normalisedElements) == 1 {
		return normalisedElements[0]
	}

	return types.NewIntersection(normalisedElements...)
}

func (c *Checker) newNormalisedUnion(elements ...types.Type) types.Type {
	var normalisedElements []types.Type

elementLoop:
	for i := 0; i < len(elements); i++ {
		element := c.normaliseType(elements[i])
		if types.IsNever(element) || types.IsNothing(element) {
			continue elementLoop
		}
		switch e := element.(type) {
		case *types.Union:
			elements = append(elements, e.Elements...)
		case *types.Nilable:
			elements = append(elements, e.Type, types.Nil{})
		case *types.Not:
			for j := 0; j < len(normalisedElements); j++ {
				normalisedElement := normalisedElements[j]
				if c.isTheSameType(e.Type, normalisedElement, nil) {
					return types.Any{}
				}
				if c.isSubtype(normalisedElement, element, nil) {
					normalisedElements[j] = element
					continue elementLoop
				}
				if c.isSubtype(element, normalisedElement, nil) {
					continue elementLoop
				}
			}
			normalisedElements = append(normalisedElements, element)
		default:
			for j := 0; j < len(normalisedElements); j++ {
				normalisedElement := normalisedElements[j]
				if normalisedNot, ok := normalisedElement.(*types.Not); ok && c.isTheSameType(normalisedNot.Type, element, nil) {
					return types.Any{}
				}
				if c.isSubtype(normalisedElement, element, nil) {
					normalisedElements[j] = element
					continue elementLoop
				}
				if c.isSubtype(element, normalisedElement, nil) {
					continue elementLoop
				}
			}
			normalisedElements = append(normalisedElements, element)
		}
	}

	if len(normalisedElements) == 0 {
		return types.Never{}
	}
	if len(normalisedElements) == 1 {
		return normalisedElements[0]
	}

	return types.NewUnion(normalisedElements...)
}

func (c *Checker) constructUnionType(node *ast.BinaryTypeExpressionNode) *ast.UnionTypeNode {
	union := types.NewUnion()
	elements := new([]ast.TypeNode)
	c._constructUnionType(node, elements, union)
	normalisedUnion := c.normaliseType(union)

	newNode := ast.NewUnionTypeNode(
		node.Span(),
		*elements,
	)
	newNode.SetType(normalisedUnion)
	return newNode
}

func (c *Checker) _constructUnionType(node *ast.BinaryTypeExpressionNode, elements *[]ast.TypeNode, union *types.Union) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeExpressionNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.OR {
		c._constructUnionType(leftBinaryType, elements, union)
	} else {
		leftTypeNode := node.Left
		leftTypeNode = c.checkTypeNode(leftTypeNode)
		*elements = append(*elements, leftTypeNode)

		leftType := c.typeOf(leftTypeNode)
		union.Elements = append(union.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeExpressionNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.OR {
		c._constructUnionType(rightBinaryType, elements, union)
	} else {
		rightTypeNode := node.Right
		rightTypeNode = c.checkTypeNode(rightTypeNode)
		*elements = append(*elements, rightTypeNode)

		rightType := c.typeOf(rightTypeNode)
		union.Elements = append(union.Elements, rightType)
	}
}

func (c *Checker) constructIntersectionType(node *ast.BinaryTypeExpressionNode) *ast.IntersectionTypeNode {
	intersection := types.NewIntersection()
	elements := new([]ast.TypeNode)
	c._constructIntersectionType(node, elements, intersection)
	normalisedIntersection := c.normaliseType(intersection)

	newNode := ast.NewIntersectionTypeNode(
		node.Span(),
		*elements,
	)
	newNode.SetType(normalisedIntersection)
	return newNode
}

func (c *Checker) _constructIntersectionType(node *ast.BinaryTypeExpressionNode, elements *[]ast.TypeNode, intersection *types.Intersection) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeExpressionNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(leftBinaryType, elements, intersection)
	} else {
		leftTypeNode := node.Left
		leftTypeNode = c.checkTypeNode(leftTypeNode)
		*elements = append(*elements, leftTypeNode)

		leftType := c.typeOf(leftTypeNode)
		intersection.Elements = append(intersection.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeExpressionNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(rightBinaryType, elements, intersection)
	} else {
		rightTypeNode := node.Right
		rightTypeNode = c.checkTypeNode(rightTypeNode)
		*elements = append(*elements, rightTypeNode)

		rightType := c.typeOf(rightTypeNode)
		intersection.Elements = append(intersection.Elements, rightType)
	}
}
