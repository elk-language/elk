// Package checker implements the Elk type checker
package checker

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
)

func (c *Checker) normaliseType(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Union:
		return c.newNormalisedUnion(t.Elements...)
	case *types.Intersection:
		return c.newNormalisedIntersection(t.Elements...)
	case *types.Nilable:
		t.Type = c.normaliseType(t.Type)
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
		}

		return t
	default:
		return typ
	}
}

func normaliseLiteralInIntersection[E types.SimpleLiteral](c *Checker, normalisedElements []types.Type, element E) (types.Type, bool) {
	for j := 0; j < len(normalisedElements); j++ {
		switch normalisedElement := normalisedElements[j].(type) {
		case E:
			if element.StringValue() == normalisedElement.StringValue() {
				return nil, false
			}
			return nil, true
		case types.SimpleLiteral:
			return nil, true
		default:
			if c.isSubtype(normalisedElement, element, nil) {
				return nil, false
			}
			if c.isSubtype(element, normalisedElement, nil) {
				normalisedElements[j] = element
				return nil, false
			}
		}
	}

	return element, true
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
		}
	}
	distributedIntersection := c.intersectionOfUnionsToUnionOfIntersections(elements)
	intersection, ok := distributedIntersection.(*types.Intersection)
	if !ok {
		return c.normaliseType(distributedIntersection)
	}

	elements = intersection.Elements
	var normalisedElements []types.Type

firstElementLoop:
	for i := 0; i < len(elements); i++ {
		element := elements[i]
		if types.IsNever(element) || types.IsNothing(element) {
			return element
		}
		switch e := element.(type) {
		case *types.Not:
			for j := 0; j < len(normalisedElements); j++ {
				normalisedElement := normalisedElements[j]
				if c.isTheSameType(normalisedElement, e.Type, nil) {
					newNormalisedElements := normalisedElements[0:j]
					if len(normalisedElements) > j+1 {
						newNormalisedElements = append(newNormalisedElements, normalisedElements[j+1:]...)
					}
					normalisedElements = newNormalisedElements
					continue firstElementLoop
				}
			}
			normalisedElements = append(normalisedElements, element)
		default:
			for j := 0; j < len(normalisedElements); j++ {
				normalisedElement, ok := normalisedElements[j].(*types.Not)
				if !ok {
					continue
				}

				if c.isTheSameType(element, normalisedElement.Type, nil) {
					newNormalisedElements := normalisedElements[0:j]
					if len(normalisedElements) > j+1 {
						newNormalisedElements = append(newNormalisedElements, normalisedElements[j+1:]...)
					}
					normalisedElements = newNormalisedElements
					continue firstElementLoop
				}
			}
			normalisedElements = append(normalisedElements, element)
		}
	}

	elements = normalisedElements
	normalisedElements = nil

secondElementLoop:
	for i := 0; i < len(elements); i++ {
		element := elements[i]
		switch e := element.(type) {
		case *types.Class:
			for j := 0; j < len(normalisedElements); j++ {
				switch normalisedElement := c.toNonLiteral(normalisedElements[j]).(type) {
				case *types.Class:
					if c.isSubtype(normalisedElement, element, nil) {
						continue secondElementLoop
					}
					if c.isSubtype(element, normalisedElement, nil) {
						normalisedElements[j] = element
						continue secondElementLoop
					}
					return types.Never{}
				default:
					if c.isSubtype(normalisedElement, element, nil) {
						continue secondElementLoop
					}
					if c.isSubtype(element, normalisedElement, nil) {
						normalisedElements[j] = element
						continue secondElementLoop
					}
				}
			}
			normalisedElements = append(normalisedElements, element)
		case *types.Module:
			for j := 0; j < len(normalisedElements); j++ {
				switch normalisedElement := normalisedElements[j].(type) {
				case *types.Module:
					if element == normalisedElement {
						continue secondElementLoop
					}
					return types.Never{}
				default:
					if c.isSubtype(normalisedElement, element, nil) {
						continue secondElementLoop
					}
					if c.isSubtype(element, normalisedElement, nil) {
						normalisedElements[j] = element
						continue secondElementLoop
					}
				}
			}
			normalisedElements = append(normalisedElements, element)
		case *types.IntLiteral:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.Int64Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.Int32Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.Int16Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.Int8Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.UInt64Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.UInt32Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.UInt16Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.UInt8Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.FloatLiteral:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.Float64Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.Float32Literal:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.BigFloatLiteral:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.StringLiteral:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.CharLiteral:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		case *types.SymbolLiteral:
			typ, ok := normaliseLiteralInIntersection(c, normalisedElements, e)
			if !ok {
				continue secondElementLoop
			}
			if typ == nil {
				return types.Never{}
			}

			normalisedElements = append(normalisedElements, element)
		default:
			for j := 0; j < len(normalisedElements); j++ {
				normalisedElement := normalisedElements[j]
				if c.isSubtype(normalisedElement, element, nil) {
					continue secondElementLoop
				}
				if c.isSubtype(element, normalisedElement, nil) {
					normalisedElements[j] = element
					continue secondElementLoop
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
			for _, normalisedElement := range normalisedElements {
				if c.isTheSameType(e.Type, normalisedElement, nil) {
					return types.Any{}
				}
				if c.isSubtype(element, normalisedElement, nil) || c.isSubtype(normalisedElement, element, nil) {
					continue elementLoop
				}
			}
			normalisedElements = append(normalisedElements, element)
		default:
			for _, normalisedElement := range normalisedElements {
				if normalisedNot, ok := normalisedElement.(*types.Not); ok && c.isTheSameType(normalisedNot.Type, element, nil) {
					return types.Any{}
				}
				if c.isSubtype(element, normalisedElement, nil) || c.isSubtype(normalisedElement, element, nil) {
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

func (c *Checker) toNonNilable(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Nilable:
		return t.Type
	case types.Nil:
		return types.Never{}
	case *types.Class:
		if t == c.StdNil() {
			return types.Never{}
		}
		return t
	case *types.Union:
		var newElements []types.Type
		for _, element := range t.Elements {
			newElements = append(newElements, c.toNonNilable(element))
		}
		return c.newNormalisedUnion(newElements...)
	case *types.Intersection:
		for _, element := range t.Elements {
			nonNilable := c.toNonNilable(element)
			if types.IsNever(nonNilable) || types.IsNothing(nonNilable) {
				return types.Never{}
			}
		}
		return t
	default:
		return t
	}
}

func (c *Checker) toNonFalsy(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Nilable:
		return t.Type
	case *types.Class:
		if t == c.StdNil() || t == c.StdFalse() {
			return types.Never{}
		}
		if t == c.StdBool() {
			return types.True{}
		}
		return t
	case types.Nil, types.False:
		return types.Never{}
	case *types.Union:
		var newElements []types.Type
		for _, element := range t.Elements {
			newElements = append(newElements, c.toNonFalsy(element))
		}
		return c.newNormalisedUnion(newElements...)
	case *types.Intersection:
		for _, element := range t.Elements {
			nonFalsy := c.toNonFalsy(element)
			if types.IsNever(nonFalsy) || types.IsNothing(nonFalsy) {
				return types.Never{}
			}
		}
		return t
	default:
		return t
	}
}

func (c *Checker) toNonTruthy(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Nilable:
		return types.Nil{}
	case *types.Class:
		if t == c.StdNil() || t == c.StdFalse() {
			return t
		}
		if t == c.StdBool() {
			return types.False{}
		}
		return types.Never{}
	case types.Nil, types.False:
		return t
	case *types.Union:
		var newElements []types.Type
		for _, element := range t.Elements {
			newElements = append(newElements, c.toNonTruthy(element))
		}
		return c.newNormalisedUnion(newElements...)
	case *types.Intersection:
		for _, element := range t.Elements {
			nonTruthy := c.toNonTruthy(element)
			if types.IsNever(nonTruthy) || types.IsNothing(nonTruthy) {
				return types.Never{}
			}
		}
		return t
	default:
		return types.Never{}
	}
}

func (c *Checker) toNonLiteral(typ types.Type) types.Type {
	return typ.ToNonLiteral(c.GlobalEnv)
}

func (c *Checker) toNilable(typ types.Type) types.Type {
	return c.normaliseType(types.NewNilable(typ))
}
