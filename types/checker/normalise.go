// Package checker implements the Elk type checker
package checker

import (
	"slices"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

func (c *Checker) inferTypeArguments(givenType, paramType types.Type, typeArgMap types.TypeArgumentMap, errLocation *position.Location) types.Type {
	if typeArgMap == nil {
		return paramType
	}

	switch p := paramType.(type) {
	case types.Self:
		arg := typeArgMap[symbol.L_self]
		if arg == nil {
			return p
		}
		return arg.Type
	case *types.Callable:
		g, ok := givenType.(*types.Callable)
		if !ok {
			return p
		}

		gMethod := g.Body
		pMethod := p.Body
		var isDifferent bool
		newParams := slices.Clone(pMethod.Params)
		for i := range min(len(pMethod.Params), len(gMethod.Params)) {
			pParam := pMethod.Params[i]
			gParam := gMethod.Params[i]
			if pParam.Kind != gParam.Kind {
				return p
			}
			result := c.inferTypeArguments(gParam.Type, pParam.Type, typeArgMap, errLocation)
			if result == nil {
				return nil
			}
			if result != pParam.Type {
				isDifferent = true
				newParam := pParam.Copy()
				newParam.Type = result
				newParams[i] = newParam
			}
		}

		returnType := c.inferTypeArguments(gMethod.ReturnType, pMethod.ReturnType, typeArgMap, errLocation)
		if returnType == nil {
			return nil
		}
		if returnType != pMethod.ReturnType {
			isDifferent = true
		}

		throwType := c.inferTypeArguments(gMethod.ThrowType, pMethod.ThrowType, typeArgMap, errLocation)
		if throwType == nil {
			return nil
		}
		if throwType != pMethod.ThrowType {
			isDifferent = true
		}

		if isDifferent {
			closure := types.NewCallable(nil, p.IsClosure)
			newMethod := types.NewMethod(
				pMethod.DocComment,
				pMethod.Flags.ToBitFlag(),
				pMethod.Name,
				pMethod.TypeParameters,
				newParams,
				returnType,
				throwType,
				closure,
			)
			closure.Body = newMethod
			return closure
		}
		return p
	case *types.TypeParameter:
		typeArg := typeArgMap[p.Name]
		if typeArg != nil {
			return typeArg.Type
		}
		if givenType == nil {
			return nil
		}

		inferredType := c.ToNonLiteral(givenType, false)
		if !c.isSubtype(givenType, p.UpperBound, nil) {
			c.addUpperBoundError(givenType, p.UpperBound, errLocation)
			return nil
		}
		if !c.isSubtype(p.LowerBound, inferredType, nil) {
			if !c.isSubtype(inferredType, p.LowerBound, nil) {
				c.addLowerBoundError(givenType, p.LowerBound, errLocation)
				return nil
			}
			inferredType = p.LowerBound
		}
		typeArgMap[p.Name] = types.NewTypeArgument(
			inferredType,
			p.Variance,
		)
		return inferredType
	case *types.Generic:
		gNamespace, ok := givenType.(types.Namespace)
		if !ok || !c.isSubtype(gNamespace, p.Namespace, nil) {
			newArgMap := make(types.TypeArgumentMap, len(p.ArgumentMap))
			for _, argName := range p.ArgumentOrder {
				pArg := p.ArgumentMap[argName]
				result := c.inferTypeArguments(nil, pArg.Type, typeArgMap, errLocation)
				if result == nil {
					return p
				}
				newArgMap[argName] = types.NewTypeArgument(result, pArg.Variance)
			}
			return types.NewGeneric(
				p.Namespace,
				types.NewTypeArguments(
					newArgMap,
					p.ArgumentOrder,
				),
			)
		}
		gGeneric, ok := gNamespace.(*types.Generic)
		if ok && c.IsTheSameNamespace(gGeneric.Namespace, p.Namespace) {
			newArgMap := make(types.TypeArgumentMap, len(p.ArgumentMap))
			for _, argName := range p.ArgumentOrder {
				pArg := p.ArgumentMap[argName]
				gArg := gGeneric.ArgumentMap[argName]
				if gArg == nil || pArg == nil {
					return nil
				}
				result := c.inferTypeArguments(gArg.Type, pArg.Type, typeArgMap, errLocation)
				if result == nil {
					return nil
				}
				newArgMap[argName] = types.NewTypeArgument(result, gArg.Variance)
			}
			return types.NewGeneric(
				p.Namespace,
				types.NewTypeArguments(
					newArgMap,
					p.ArgumentOrder,
				),
			)
		}

		var resolvedG *types.Generic
		var isSubtype bool
		for gParent := range types.Parents(gNamespace) {
			gGenericParent, ok := gParent.(*types.Generic)
			if !ok {
				continue
			}
			if resolvedG == nil {
				resolvedG = gGenericParent
				continue
			}
			resolvedG = c.replaceTypeParametersInGeneric(gGenericParent, resolvedG.ArgumentMap, false)

			if !c.IsTheSameNamespace(gGenericParent.Namespace, p.Namespace) {
				continue
			}

			isSubtype = true
			break
		}
		if isSubtype {
			return c.inferTypeArguments(resolvedG, p, typeArgMap, errLocation)
		}

		return nil
	case *types.SingletonOf:
		switch g := givenType.(type) {
		case *types.SingletonClass:
			result := c.inferTypeArguments(g.AttachedObject, p.Type, typeArgMap, errLocation)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			return types.NewSingletonOf(result)
		case *types.SingletonOf:
			result := c.inferTypeArguments(g.Type, p.Type, typeArgMap, errLocation)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			return types.NewSingletonOf(result)
		default:
			return p
		}
	case *types.SingletonClass:
		switch g := givenType.(type) {
		case *types.SingletonClass:
			result := c.inferTypeArguments(g.AttachedObject, p.AttachedObject, typeArgMap, errLocation)
			if result == nil {
				return nil
			}
			if p.AttachedObject == result {
				return p
			}

			return types.NewSingletonClass(result.(types.Namespace), p.Parent())
		case *types.SingletonOf:
			result := c.inferTypeArguments(g.Type, p.AttachedObject, typeArgMap, errLocation)
			if result == nil {
				return nil
			}
			if p.AttachedObject == result {
				return p
			}

			return types.NewSingletonClass(result.(types.Namespace), p.Parent())
		default:
			return p
		}
	case *types.InstanceOf:
		nonLiteral := c.ToNonLiteral(givenType, false)
		switch g := nonLiteral.(type) {
		case *types.InstanceOf:
			result := c.inferTypeArguments(g.Type, p.Type, typeArgMap, errLocation)
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
			return p
		case *types.Class:
			result := c.inferTypeArguments(g.Singleton(), p.Type, typeArgMap, errLocation)
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
			return p
		case *types.Mixin:
			result := c.inferTypeArguments(g.Singleton(), p.Type, typeArgMap, errLocation)
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
			return p
		case *types.Interface:
			result := c.inferTypeArguments(g.Singleton(), p.Type, typeArgMap, errLocation)
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
			return p
		default:
			return p
		}
	case *types.Not:
		g, ok := givenType.(*types.Not)
		if !ok {
			return p
		}

		result := c.inferTypeArguments(g.Type, p.Type, typeArgMap, errLocation)
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
				result := c.inferTypeArguments(newG, pElement, typeArgMap, errLocation)
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
				result := c.inferTypeArguments(g, pElement, typeArgMap, errLocation)
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
				result := c.inferTypeArguments(narrowedG, pElement, typeArgMap, errLocation)
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
			return c.inferTypeArguments(types.NewUnion(types.Nil{}, g.Type), p, typeArgMap, errLocation)
		default:
			newElements := make([]types.Type, 0, len(p.Elements))
			var isDifferent bool
			for _, pElement := range p.Elements {
				result := c.inferTypeArguments(g, pElement, typeArgMap, errLocation)
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
			result := c.inferTypeArguments(g.Type, p.Type, typeArgMap, errLocation)
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

			result := c.inferTypeArguments(t, p.Type, typeArgMap, errLocation)
			if result == nil {
				return nil
			}
			if p.Type == result {
				return p
			}

			return types.NewNilable(result)
		default:
			result := c.inferTypeArguments(givenType, p.Type, typeArgMap, errLocation)
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

func (c *Checker) replaceTypeParametersOfGeneric(typ types.Type, generic *types.Generic) types.Type {
	switch t := typ.(type) {
	case types.Self:
		arg := generic.ArgumentMap[symbol.L_self]
		if arg == nil {
			return t
		}
		return arg.Type
	case *types.TypeParameter:
		if !c.isTheSameType(t.Namespace, generic.Namespace, nil) {
			return t
		}
		arg := generic.ArgumentMap[t.Name]
		if arg == nil {
			return t
		}
		return arg.Type
	case *types.SingletonOf:
		result := c.replaceTypeParametersOfGeneric(t.Type, generic)
		if result == t.Type {
			return t
		}
		return types.NewSingletonOf(
			result,
		)
	case *types.InstanceOf:
		result := c.replaceTypeParametersOfGeneric(t.Type, generic)
		if result == t.Type {
			return t
		}
		return types.NewInstanceOf(
			result,
		)
	case *types.Callable:
		newParams := make([]*types.Parameter, len(t.Body.Params))
		var isDifferent bool
		for i, param := range t.Body.Params {
			result := c.replaceTypeParametersOfGeneric(param.Type, generic)
			if result == param.Type {
				newParams[i] = param
				continue
			}

			newParam := param.Copy()
			newParam.Type = result
			newParams[i] = newParam
			isDifferent = true
		}

		returnType := c.replaceTypeParametersOfGeneric(t.Body.ReturnType, generic)
		if returnType != t.Body.ReturnType {
			isDifferent = true
		}
		throwType := c.replaceTypeParametersOfGeneric(t.Body.ThrowType, generic)
		if throwType != t.Body.ThrowType {
			isDifferent = true
		}

		if !isDifferent {
			return t
		}
		method := t.Body.Copy()
		method.Params = newParams
		method.ReturnType = returnType
		method.ThrowType = throwType

		closure := types.NewCallable(method, t.IsClosure)
		method.DefinedUnder = closure
		return closure
	case *types.Generic:
		newMap := make(types.TypeArgumentMap, len(t.ArgumentMap))
		var isDifferent bool
		for key, arg := range t.AllArguments() {
			result := c.replaceTypeParametersOfGeneric(arg.Type, generic)
			if result == arg.Type {
				newMap[key] = arg
				continue
			}
			newMap[key] = types.NewTypeArgument(
				result,
				arg.Variance,
			)
			isDifferent = true
		}
		result := c.replaceTypeParametersOfGeneric(t.Namespace, generic)
		if result != t.Namespace {
			isDifferent = true
		}
		if !isDifferent {
			return t
		}

		return types.NewGeneric(
			result.(types.Namespace),
			types.NewTypeArguments(
				newMap,
				t.ArgumentOrder,
			),
		)
	case *types.Nilable:
		result := c.replaceTypeParametersOfGeneric(t.Type, generic)
		if result == t.Type {
			return t
		}
		return types.NewNilable(result)
	case *types.Not:
		result := c.replaceTypeParametersOfGeneric(t.Type, generic)
		if result == t.Type {
			return t
		}
		return types.NewNot(result)
	case *types.Union:
		newElements := make([]types.Type, len(t.Elements))
		var isDifferent bool
		for i, element := range t.Elements {
			result := c.replaceTypeParametersOfGeneric(element, generic)
			if result != element {
				isDifferent = true
			}
			newElements[i] = result
		}
		if !isDifferent {
			return t
		}
		return types.NewUnion(newElements...)
	case *types.Intersection:
		newElements := make([]types.Type, len(t.Elements))
		var isDifferent bool
		for i, element := range t.Elements {
			result := c.replaceTypeParametersOfGeneric(element, generic)
			if result != element {
				isDifferent = true
			}
			newElements[i] = result
		}
		if !isDifferent {
			return t
		}
		return types.NewIntersection(newElements...)
	default:
		return t
	}
}

func (c *Checker) replaceTypeParameters(typ types.Type, typeArgMap types.TypeArgumentMap, replaceMethodTypeParams bool) types.Type {
	return c.NormaliseType(c._replaceTypeParameters(typ, typeArgMap, replaceMethodTypeParams))
}

func (c *Checker) _replaceTypeParameters(typ types.Type, typeArgMap types.TypeArgumentMap, replaceMethodTypeParams bool) types.Type {
	switch t := typ.(type) {
	case types.Self:
		arg := typeArgMap[symbol.L_self]
		if arg == nil {
			return t
		}
		return arg.Type
	case *types.SingletonOf:
		result := c._replaceTypeParameters(t.Type, typeArgMap, replaceMethodTypeParams)
		if result == t.Type {
			return t
		}
		return types.NewSingletonOf(
			result,
		)
	case *types.InstanceOf:
		result := c._replaceTypeParameters(t.Type, typeArgMap, replaceMethodTypeParams)
		if result == t.Type {
			return t
		}
		return types.NewInstanceOf(
			result,
		)
	case *types.Callable:
		newParams := make([]*types.Parameter, len(t.Body.Params))
		var isDifferent bool
		for i, param := range t.Body.Params {
			result := c._replaceTypeParameters(param.Type, typeArgMap, replaceMethodTypeParams)
			if result == param.Type {
				newParams[i] = param
				continue
			}

			newParam := param.Copy()
			newParam.Type = result
			newParams[i] = newParam
			isDifferent = true
		}

		returnType := c._replaceTypeParameters(t.Body.ReturnType, typeArgMap, replaceMethodTypeParams)
		if returnType != t.Body.ReturnType {
			isDifferent = true
		}
		throwType := c._replaceTypeParameters(t.Body.ThrowType, typeArgMap, replaceMethodTypeParams)
		if throwType != t.Body.ThrowType {
			isDifferent = true
		}

		if !isDifferent {
			return t
		}
		method := t.Body.Copy()
		method.Params = newParams
		method.ReturnType = returnType
		method.ThrowType = throwType

		closure := types.NewCallable(method, t.IsClosure)
		method.DefinedUnder = closure
		return closure
	case *types.Generic:
		return c.replaceTypeParametersInGeneric(t, typeArgMap, replaceMethodTypeParams)
	case *types.TypeParameter:
		// do not replace type parameters of methods when the `replaceMethodTypeParams` flag is false
		if n, ok := t.Namespace.(*types.TypeParamNamespace); ok && n.ForMethod && !replaceMethodTypeParams {
			return t
		}
		arg := typeArgMap[t.Name]
		if arg == nil {
			return t
		}
		return arg.Type
	case *types.Nilable:
		result := c._replaceTypeParameters(t.Type, typeArgMap, replaceMethodTypeParams)
		if result == t.Type {
			return t
		}
		return types.NewNilable(result)
	case *types.Not:
		result := c._replaceTypeParameters(t.Type, typeArgMap, replaceMethodTypeParams)
		if result == t.Type {
			return t
		}
		return types.NewNot(result)
	case *types.Union:
		newElements := make([]types.Type, len(t.Elements))
		var isDifferent bool
		for i, element := range t.Elements {
			result := c._replaceTypeParameters(element, typeArgMap, replaceMethodTypeParams)
			if result != element {
				isDifferent = true
			}
			newElements[i] = result
		}
		if !isDifferent {
			return t
		}
		return types.NewUnion(newElements...)
	case *types.Intersection:
		newElements := make([]types.Type, len(t.Elements))
		var isDifferent bool
		for i, element := range t.Elements {
			result := c._replaceTypeParameters(element, typeArgMap, replaceMethodTypeParams)
			if result != element {
				isDifferent = true
			}
			newElements[i] = result
		}
		if !isDifferent {
			return t
		}
		return types.NewIntersection(newElements...)
	default:
		return t
	}
}

func (c *Checker) replaceTypeParametersInGeneric(t *types.Generic, typeArgMap types.TypeArgumentMap, replaceMethodTypeParams bool) *types.Generic {
	newMap := make(types.TypeArgumentMap, len(t.ArgumentMap))
	var isDifferent bool
	for key, arg := range t.AllArguments() {
		result := c._replaceTypeParameters(arg.Type, typeArgMap, replaceMethodTypeParams)
		if result == arg.Type {
			newMap[key] = arg
			continue
		}
		newMap[key] = types.NewTypeArgument(
			result,
			arg.Variance,
		)
		isDifferent = true
	}
	result := c._replaceTypeParameters(t.Namespace, typeArgMap, replaceMethodTypeParams)
	if result != t.Namespace {
		isDifferent = true
	}
	if !isDifferent {
		return t
	}

	return types.NewGeneric(
		result.(types.Namespace),
		types.NewTypeArguments(
			newMap,
			t.ArgumentOrder,
		),
	)
}

func (c *Checker) NormaliseType(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Union:
		return c.NewNormalisedUnion(t.Elements...)
	case *types.Intersection:
		return c.NewNormalisedIntersection(t.Elements...)
	case *types.Generic:
		for _, arg := range t.TypeArguments.AllArguments() {
			arg.Type = c.NormaliseType(arg.Type)
		}
		return t
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
		t.Type = c.NormaliseType(t.Type)
		switch t.Type.(type) {
		case types.Never:
			return types.Nil{}
		case types.Any, types.Untyped:
			return t.Type
		}
		if c.IsNilable(t.Type) {
			return t.Type
		}
		if union, ok := t.Type.(*types.Union); ok {
			union.Elements = append(union.Elements, types.Nil{})
			return union
		}
		return t
	case *types.Not:
		t.Type = c.NormaliseType(t.Type)
		switch nestedType := t.Type.(type) {
		case *types.Not:
			return nestedType.Type
		case types.Never:
			return types.Any{}
		case types.Any:
			return types.Never{}
		case types.Untyped:
			return types.Untyped{}
		case *types.Union:
			intersectionElements := make([]types.Type, 0, len(nestedType.Elements))
			for _, element := range nestedType.Elements {
				intersectionElements = append(intersectionElements, types.NewNot(element))
			}
			return c.NewNormalisedIntersection(intersectionElements...)
		case *types.Intersection:
			unionElements := make([]types.Type, 0, len(nestedType.Elements))
			for _, element := range nestedType.Elements {
				unionElements = append(unionElements, types.NewNot(element))
			}
			return c.NewNormalisedUnion(unionElements...)
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

func (c *Checker) NewNormalisedIntersection(elements ...types.Type) types.Type {
	var containsNot bool
	var containsUninitialisedNamedTypes bool

	for i := 0; i < len(elements); i++ {
		element := c.NormaliseType(elements[i])
		if types.IsNever(element) || types.IsUntyped(element) {
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
		return c.NormaliseType(distributedIntersection)
	}

	elements = intersection.Elements
	normalisedElements := make([]types.Type, 0, len(elements))

	// detect empty intersections
	for _, element := range elements {
		if types.IsNever(element) || types.IsUntyped(element) {
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
		elements[i] = c.NormaliseType(elements[i])
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

func (c *Checker) NewNormalisedUnion(elements ...types.Type) types.Type {
	var normalisedElements []types.Type

elementLoop:
	for i := 0; i < len(elements); i++ {
		element := c.NormaliseType(elements[i])
		if types.IsNever(element) || types.IsUntyped(element) {
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

func (c *Checker) constructUnionType(node *ast.BinaryTypeNode) *ast.UnionTypeNode {
	union := types.NewUnion()
	elements := new([]ast.TypeNode)
	c._constructUnionType(node, elements, union)
	normalisedUnion := c.NormaliseType(union)

	newNode := ast.NewUnionTypeNode(
		node.Location(),
		*elements,
	)
	newNode.SetType(normalisedUnion)
	return newNode
}

func (c *Checker) _constructUnionType(node *ast.BinaryTypeNode, elements *[]ast.TypeNode, union *types.Union) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.OR {
		c._constructUnionType(leftBinaryType, elements, union)
	} else {
		leftTypeNode := node.Left
		leftTypeNode = c.checkTypeNode(leftTypeNode)
		*elements = append(*elements, leftTypeNode)

		leftType := c.TypeOf(leftTypeNode)
		union.Elements = append(union.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.OR {
		c._constructUnionType(rightBinaryType, elements, union)
	} else {
		rightTypeNode := node.Right
		rightTypeNode = c.checkTypeNode(rightTypeNode)
		*elements = append(*elements, rightTypeNode)

		rightType := c.TypeOf(rightTypeNode)
		union.Elements = append(union.Elements, rightType)
	}
}

func (c *Checker) constructIntersectionType(node *ast.BinaryTypeNode) *ast.IntersectionTypeNode {
	intersection := types.NewIntersection()
	elements := new([]ast.TypeNode)
	c._constructIntersectionType(node, elements, intersection)
	normalisedIntersection := c.NormaliseType(intersection)

	newNode := ast.NewIntersectionTypeNode(
		node.Location(),
		*elements,
	)
	newNode.SetType(normalisedIntersection)
	return newNode
}

func (c *Checker) _constructIntersectionType(node *ast.BinaryTypeNode, elements *[]ast.TypeNode, intersection *types.Intersection) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(leftBinaryType, elements, intersection)
	} else {
		leftTypeNode := node.Left
		leftTypeNode = c.checkTypeNode(leftTypeNode)
		*elements = append(*elements, leftTypeNode)

		leftType := c.TypeOf(leftTypeNode)
		intersection.Elements = append(intersection.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(rightBinaryType, elements, intersection)
	} else {
		rightTypeNode := node.Right
		rightTypeNode = c.checkTypeNode(rightTypeNode)
		*elements = append(*elements, rightTypeNode)

		rightType := c.TypeOf(rightTypeNode)
		intersection.Elements = append(intersection.Elements, rightType)
	}
}
