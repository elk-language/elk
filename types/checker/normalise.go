// Package checker implements the Elk type checker
package checker

import (
	"slices"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

const (
	inferTypeArgumentsInferFromDefaults bitfield.BitFlag8 = 1 << iota
)

func (c *Checker) inferTypeArguments(givenType, paramType types.Type, typeArgMap types.TypeArgumentMap, errLocation *position.Location) types.Type {
	return c.inferTypeArgumentsWithFlags(
		givenType,
		paramType,
		typeArgMap,
		errLocation,
		bitfield.BitField8{},
	)
}

func (c *Checker) inferTypeArgumentsWithFlags(givenType, paramType types.Type, typeArgMap types.TypeArgumentMap, errLocation *position.Location, flags bitfield.BitField8) types.Type {
	if typeArgMap == nil {
		return paramType
	}

	switch p := paramType.(type) {
	case types.Self:
		arg := typeArgMap[symbol.L_self]
		if arg == nil {
			return p
		}
		return arg.Type.Get()
	case *types.Callable:
		g, ok := givenType.(*types.Callable)
		if !ok {
			return p
		}

		gMethod := g.Body.Get()
		pMethod := p.Body.Get()
		var isDifferent bool
		newParams := slices.Clone(pMethod.Params)
		for i := range min(len(pMethod.Params), len(gMethod.Params)) {
			pParam := pMethod.Params[i].Get()
			gParam := gMethod.Params[i].Get()
			if pParam.Kind != gParam.Kind {
				return p
			}
			pParamType := pParam.Type.Get()
			gParamType := gParam.Type.Get()
			result := c.inferTypeArgumentsWithFlags(gParamType, pParamType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if result != pParamType {
				isDifferent = true
				newParam := pParam.Copy()
				newParam.Type = types.ToRef(result)
				newParams[i] = newParam.ToRef()
			}
		}

		gMethodReturnType := gMethod.ReturnType.Get()
		pMethodReturnType := pMethod.ReturnType.Get()
		returnType := c.inferTypeArgumentsWithFlags(gMethodReturnType, pMethodReturnType, typeArgMap, errLocation, flags)
		if returnType == nil {
			return nil
		}
		if returnType != pMethodReturnType {
			isDifferent = true
		}

		gMethodThrowType := gMethod.ThrowType.Get()
		pMethodThrowType := pMethod.ThrowType.Get()
		throwType := c.inferTypeArgumentsWithFlags(gMethodThrowType, pMethodThrowType, typeArgMap, errLocation, flags)
		if throwType == nil {
			return nil
		}
		if throwType != pMethodThrowType {
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
			closure.Body = newMethod.ToRef()
			return closure
		}
		return p
	case *types.TypeParameter:
		typeArg := typeArgMap[p.Name]
		if typeArg != nil {
			return typeArg.Type.Get()
		}
		if givenType == nil {
			return nil
		}
		if _, ok := givenType.(*types.TypeParameter); ok {
			if flags.HasFlag(inferTypeArgumentsInferFromDefaults) {
				inferredType := p.InferredType()
				typeArgMap[p.Name] = types.NewTypeArgument(
					inferredType,
					p.Variance,
				)
				return inferredType
			}
			return p
		}

		inferredType := c.ToNonLiteral(givenType, false)
		pUpperBound := p.UpperBound.Get()
		if !c.isSubtype(givenType, pUpperBound, nil) {
			c.addUpperBoundError(givenType, pUpperBound, errLocation)
			return nil
		}

		pLowerBound := p.LowerBound.Get()
		if !c.isSubtype(pLowerBound, inferredType, nil) {
			if !c.isSubtype(inferredType, pLowerBound, nil) {
				c.addLowerBoundError(givenType, pLowerBound, errLocation)
				return nil
			}
			inferredType = pLowerBound
		}
		typeArgMap[p.Name] = types.NewTypeArgument(
			inferredType,
			p.Variance,
		)
		return inferredType
	case *types.Generic:
		gNamespace, ok := givenType.(types.Namespace)
		pNamespace := p.Namespace.Get()
		if !ok || !c.isSubtype(gNamespace, pNamespace, nil) {
			newArgMap := make(types.TypeArgumentMap, len(p.ArgumentMap))
			for _, argName := range p.ArgumentOrder {
				pArg := p.ArgumentMap[argName]
				result := c.inferTypeArgumentsWithFlags(nil, pArg.Type.Get(), typeArgMap, errLocation, flags)
				if result == nil {
					return p
				}
				newArgMap[argName] = types.NewTypeArgument(result, pArg.Variance)
			}
			return types.NewGeneric(
				pNamespace,
				types.NewTypeArguments(
					newArgMap,
					p.ArgumentOrder,
				),
			)
		}
		gGeneric, ok := gNamespace.(*types.Generic)
		if ok && c.IsTheSameNamespace(gGeneric.Namespace.Get(), pNamespace) {
			newArgMap := make(types.TypeArgumentMap, len(p.ArgumentMap))
			for _, argName := range p.ArgumentOrder {
				pArg := p.ArgumentMap[argName]
				gArg := gGeneric.ArgumentMap[argName]
				if gArg == nil || pArg == nil {
					return nil
				}
				result := c.inferTypeArgumentsWithFlags(gArg.Type.Get(), pArg.Type.Get(), typeArgMap, errLocation, flags)
				if result == nil {
					return nil
				}
				newArgMap[argName] = types.NewTypeArgument(result, gArg.Variance)
			}
			return types.NewGeneric(
				pNamespace,
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

			if !c.IsTheSameNamespace(gGenericParent.Namespace.Get(), pNamespace) {
				continue
			}

			isSubtype = true
			break
		}
		if isSubtype {
			return c.inferTypeArgumentsWithFlags(resolvedG, p, typeArgMap, errLocation, flags)
		}

		return nil
	case *types.SingletonOf:
		switch g := givenType.(type) {
		case *types.Exact:
			return c.inferTypeArgumentsWithFlags(g.Type.Get(), p, typeArgMap, errLocation, flags)
		case *types.SingletonClass:
			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(g.AttachedObject.Get(), pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
				return p
			}

			return types.NewSingletonOf(result)
		case *types.SingletonOf:
			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(g.Type.Get(), pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
				return p
			}

			return types.NewSingletonOf(result)
		default:
			return p
		}
	case *types.SingletonClass:
		switch g := givenType.(type) {
		case *types.Exact:
			return c.inferTypeArgumentsWithFlags(g.Type.Get(), p, typeArgMap, errLocation, flags)
		case *types.SingletonClass:
			pAttachedObject := p.AttachedObject.Get()
			result := c.inferTypeArgumentsWithFlags(g.AttachedObject.Get(), pAttachedObject, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pAttachedObject == result {
				return p
			}

			return types.NewSingletonClass(result.(types.Namespace), p.Parent())
		case *types.SingletonOf:
			pAttachedObject := p.AttachedObject.Get()
			result := c.inferTypeArgumentsWithFlags(g.Type.Get(), pAttachedObject, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pAttachedObject == result {
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
			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(g.Type.Get(), pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
				return p
			}

			switch r := result.(type) {
			case *types.SingletonClass:
				return r.AttachedObject.Get()
			case *types.SingletonOf:
				return r.Type.Get()
			}
			return p
		case *types.Class:
			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(g.Singleton(), pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
				return p
			}

			switch r := result.(type) {
			case *types.SingletonClass:
				return r.AttachedObject.Get()
			case *types.SingletonOf:
				return r.Type.Get()
			}
			return p
		case *types.Mixin:
			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(g.Singleton(), pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
				return p
			}

			switch r := result.(type) {
			case *types.SingletonClass:
				return r.AttachedObject.Get()
			case *types.SingletonOf:
				return r.Type.Get()
			}
			return p
		case *types.Interface:
			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(g.Singleton(), pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
				return p
			}

			switch r := result.(type) {
			case *types.SingletonClass:
				return r.AttachedObject.Get()
			case *types.SingletonOf:
				return r.Type.Get()
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

		pType := p.Type.Get()
		result := c.inferTypeArgumentsWithFlags(g.Type.Get(), pType, typeArgMap, errLocation, flags)
		if result == nil {
			return nil
		}
		if pType == result {
			return p
		}

		return types.NewNot(result)
	case *types.Exact:
		g, ok := givenType.(*types.Exact)
		if !ok {
			return p
		}

		pType := p.Type.Get()
		result := c.inferTypeArgumentsWithFlags(g.Type.Get(), pType, typeArgMap, errLocation, flags)
		if result == nil {
			return nil
		}
		if pType == result {
			return p
		}

		return types.NewExact(result.(types.Namespace))
	case *types.Intersection:
		switch g := givenType.(type) {
		case *types.Intersection:
			gElementsToSkip := make([]bool, len(g.Elements))
			for _, pElement := range p.Elements {
				for j, gElement := range g.Elements {
					if c.isSubtype(gElement.Get(), pElement.Get(), nil) {
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
				newGElements = append(newGElements, gElement.Get())
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
			for _, pElementRef := range p.Elements {
				pElement := pElementRef.Get()
				result := c.inferTypeArgumentsWithFlags(newG, pElement, typeArgMap, errLocation, flags)
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

			// try matching exactly as a last resort
			newPElements = newPElements[:0]
			for i := range len(p.Elements) {
				pElement := p.Elements[i].Get()
				gElement := g.Elements[i].Get()
				result := c.inferTypeArgumentsWithFlags(gElement, pElement, typeArgMap, errLocation, flags)
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
			for _, pElementRef := range p.Elements {
				pElement := pElementRef.Get()
				result := c.inferTypeArgumentsWithFlags(g, pElement, typeArgMap, errLocation, flags)
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
			for _, gElementRef := range g.Elements {
				gElement := gElementRef.Get()
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
			for _, pElementRef := range p.Elements {
				pElement := pElementRef.Get()
				result := c.inferTypeArgumentsWithFlags(narrowedG, pElement, typeArgMap, errLocation, flags)
				if result == nil {
					return nil
				}
				if result != pElement {
					isDifferent = true
				}

				newPElements = append(newPElements, result)
			}

			if isDifferent {
				return types.NewUnion(newPElements...)
			}

			if len(p.Elements) != len(g.Elements) {
				return p
			}

			// try matching exactly as a last resort
			newPElements = newPElements[:0]
			for i := range len(p.Elements) {
				pElement := p.Elements[i].Get()
				gElement := g.Elements[i].Get()
				result := c.inferTypeArgumentsWithFlags(gElement, pElement, typeArgMap, errLocation, flags)
				if result == nil {
					return nil
				}
				if result != pElement {
					isDifferent = true
				}
				newPElements = append(newPElements, result)
			}

			if isDifferent {
				return types.NewUnion(newPElements...)
			}
			return p
		case *types.Nilable:
			return c.inferTypeArgumentsWithFlags(types.NewUnionRef(types.NilID, g.Type), p, typeArgMap, errLocation, flags)
		default:
			newElements := make([]types.Type, 0, len(p.Elements))
			var isDifferent bool
			for _, pElementRef := range p.Elements {
				pElement := pElementRef.Get()
				result := c.inferTypeArgumentsWithFlags(g, pElement, typeArgMap, errLocation, flags)
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
			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(g.Type.Get(), pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
				return p
			}

			return types.NewNilable(result)
		case *types.Union:
			var withoutNil []types.Type
			for _, elementRef := range g.Elements {
				element := elementRef.Get()
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

			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(t, pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
				return p
			}

			return types.NewNilable(result)
		default:
			pType := p.Type.Get()
			result := c.inferTypeArgumentsWithFlags(givenType, pType, typeArgMap, errLocation, flags)
			if result == nil {
				return nil
			}
			if pType == result {
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
		return arg.Type.Get()
	case *types.TypeParameter:
		if !c.isTheSameType(t.Namespace.Get(), generic.Namespace.Get(), nil) {
			return t
		}
		arg := generic.ArgumentMap[t.Name]
		if arg == nil {
			return t
		}
		return arg.Type.Get()
	case *types.SingletonOf:
		tType := t.Type.Get()
		result := c.replaceTypeParametersOfGeneric(tType, generic)
		if result == tType {
			return t
		}
		return types.NewSingletonOf(
			result,
		)
	case *types.InstanceOf:
		tType := t.Type.Get()
		result := c.replaceTypeParametersOfGeneric(tType, generic)
		if result == tType {
			return t
		}
		return types.NewInstanceOf(
			result,
		)
	case *types.Callable:
		body := t.Body.Get()
		newParams := make([]types.Ref[*types.Parameter], len(body.Params))
		var isDifferent bool
		for i, paramRef := range body.Params {
			param := paramRef.Get()
			paramType := param.Type.Get()
			result := c.replaceTypeParametersOfGeneric(paramType, generic)
			if result == paramType {
				newParams[i] = paramRef
				continue
			}

			newParam := param.Copy()
			newParam.Type = types.ToRef(result)
			newParams[i] = newParam.ToRef()
			isDifferent = true
		}

		bodyReturnType := body.ReturnType.Get()
		returnType := c.replaceTypeParametersOfGeneric(bodyReturnType, generic)
		if returnType != bodyReturnType {
			isDifferent = true
		}

		bodyThrowType := body.ThrowType.Get()
		throwType := c.replaceTypeParametersOfGeneric(bodyThrowType, generic)
		if throwType != bodyThrowType {
			isDifferent = true
		}

		if !isDifferent {
			return t
		}
		method := body.Copy()
		method.Params = newParams
		method.ReturnType = types.ToRef(returnType)
		method.ThrowType = types.ToRef(throwType)

		closure := types.NewCallable(method, t.IsClosure)
		method.DefinedUnder = types.CastRef[types.Namespace](closure)
		return closure
	case *types.Generic:
		newMap := make(types.TypeArgumentMap, len(t.ArgumentMap))
		var isDifferent bool
		for key, arg := range t.AllArguments() {
			argType := arg.Type.Get()
			result := c.replaceTypeParametersOfGeneric(argType, generic)
			if result == argType {
				newMap[key] = arg
				continue
			}
			newMap[key] = types.NewTypeArgument(
				result,
				arg.Variance,
			)
			isDifferent = true
		}

		tNamespace := t.Namespace.Get()
		result := c.replaceTypeParametersOfGeneric(tNamespace, generic)
		if result != tNamespace {
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
		tType := t.Type.Get()
		result := c.replaceTypeParametersOfGeneric(tType, generic)
		if result == tType {
			return t
		}
		return types.NewNilable(result)
	case *types.Not:
		tType := t.Type.Get()
		result := c.replaceTypeParametersOfGeneric(tType, generic)
		if result == tType {
			return t
		}
		return types.NewNot(result)
	case *types.Union:
		newElements := make([]types.Ref[types.Type], len(t.Elements))
		var isDifferent bool
		for i, elementRef := range t.Elements {
			element := elementRef.Get()
			result := c.replaceTypeParametersOfGeneric(element, generic)
			if result != element {
				isDifferent = true
			}
			newElements[i] = types.ToRef(result)
		}
		if !isDifferent {
			return t
		}
		return types.NewUnionRef(newElements...)
	case *types.Intersection:
		newElements := make([]types.Ref[types.Type], len(t.Elements))
		var isDifferent bool
		for i, elementRef := range t.Elements {
			element := elementRef.Get()
			result := c.replaceTypeParametersOfGeneric(element, generic)
			if result != element {
				isDifferent = true
			}
			newElements[i] = types.ToRef(result)
		}
		if !isDifferent {
			return t
		}
		return types.NewIntersectionRef(newElements...)
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
		return arg.Type.Get()
	case *types.SingletonOf:
		tType := t.Type.Get()
		result := c._replaceTypeParameters(tType, typeArgMap, replaceMethodTypeParams)
		if result == tType {
			return t
		}
		return types.NewSingletonOf(
			result,
		)
	case *types.InstanceOf:
		tType := t.Type.Get()
		result := c._replaceTypeParameters(tType, typeArgMap, replaceMethodTypeParams)
		if result == tType {
			return t
		}
		return types.NewInstanceOf(
			result,
		)
	case *types.Callable:
		body := t.Body.Get()
		newParams := make([]types.Ref[*types.Parameter], len(body.Params))
		var isDifferent bool
		for i, paramRef := range body.Params {
			param := paramRef.Get()
			paramType := param.Type.Get()
			result := c._replaceTypeParameters(paramType, typeArgMap, replaceMethodTypeParams)
			if result == paramType {
				newParams[i] = paramRef
				continue
			}

			newParam := param.Copy()
			newParam.Type = types.ToRef(result)
			newParams[i] = newParam.ToRef()
			isDifferent = true
		}

		bodyReturnType := body.ReturnType.Get()
		returnType := c._replaceTypeParameters(bodyReturnType, typeArgMap, replaceMethodTypeParams)
		if returnType != bodyReturnType {
			isDifferent = true
		}

		bodyThrowType := body.ThrowType.Get()
		throwType := c._replaceTypeParameters(bodyThrowType, typeArgMap, replaceMethodTypeParams)
		if throwType != bodyThrowType {
			isDifferent = true
		}

		if !isDifferent {
			return t
		}
		method := body.Copy()
		method.Params = newParams
		method.ReturnType = types.ToRef(returnType)
		method.ThrowType = types.ToRef(throwType)

		closure := types.NewCallable(method, t.IsClosure)
		method.DefinedUnder = types.CastRef[types.Namespace](closure)
		return closure
	case *types.Generic:
		return c.replaceTypeParametersInGeneric(t, typeArgMap, replaceMethodTypeParams)
	case *types.TypeParameter:
		// do not replace type parameters of methods when the `replaceMethodTypeParams` flag is false
		if n, ok := t.Namespace.Get().(*types.TypeParamNamespace); ok && n.ForMethod && !replaceMethodTypeParams {
			return t
		}
		arg := typeArgMap[t.Name]
		if arg == nil {
			return t
		}
		return arg.Type.Get()
	case *types.Nilable:
		tType := t.Type.Get()
		result := c._replaceTypeParameters(tType, typeArgMap, replaceMethodTypeParams)
		if result == tType {
			return t
		}
		return types.NewNilable(result)
	case *types.Not:
		tType := t.Type.Get()
		result := c._replaceTypeParameters(tType, typeArgMap, replaceMethodTypeParams)
		if result == tType {
			return t
		}
		return types.NewNot(result)
	case *types.Union:
		newElements := make([]types.Ref[types.Type], len(t.Elements))
		var isDifferent bool
		for i, elementRef := range t.Elements {
			element := elementRef.Get()
			result := c._replaceTypeParameters(element, typeArgMap, replaceMethodTypeParams)
			if result != element {
				isDifferent = true
			}
			newElements[i] = types.ToRef(result)
		}
		if !isDifferent {
			return t
		}
		return types.NewUnionRef(newElements...)
	case *types.Intersection:
		newElements := make([]types.Type, len(t.Elements))
		var isDifferent bool
		for i, elementRef := range t.Elements {
			element := elementRef.Get()
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
		argType := arg.Type.Get()
		result := c._replaceTypeParameters(argType, typeArgMap, replaceMethodTypeParams)
		if result == argType {
			newMap[key] = arg
			continue
		}
		newMap[key] = types.NewTypeArgument(
			result,
			arg.Variance,
		)
		isDifferent = true
	}

	tNamespace := t.Namespace.Get()
	result := c._replaceTypeParameters(tNamespace, typeArgMap, replaceMethodTypeParams)
	if result != tNamespace {
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

func (c *Checker) replaceTypeParametersInInstanceVariable(ivar *types.InstanceVariable, typeArgMap types.TypeArgumentMap, replaceMethodTypeParams bool) *types.InstanceVariable {
	ivarType := ivar.Type.Get()
	result := c._replaceTypeParameters(ivarType, typeArgMap, replaceMethodTypeParams)
	if result == ivarType {
		return ivar
	}

	return types.NewInstanceVariable(
		ivar.Name,
		result,
		ivar.DocComment,
		ivar.SingleAssignment,
	)
}

func (c *Checker) normaliseSingletonOf(typ *types.SingletonOf) types.Type {
	switch nestedType := typ.Type.Get().(type) {
	case *types.InstanceOf:
		return nestedType.Type.Get()
	case *types.Class:
		return nestedType.Singleton()
	case *types.Mixin:
		return nestedType.Singleton()
	case *types.Interface:
		return nestedType.Singleton()
	case *types.Generic:
		return c.normaliseSingletonOf(types.NewSingletonOf(nestedType.Namespace.Get()))
	default:
		return typ
	}
}

func (c *Checker) normaliseInstanceOf(typ *types.InstanceOf) types.Type {
	switch nestedType := typ.Type.Get().(type) {
	case *types.SingletonOf:
		return nestedType.Type.Get()
	case *types.SingletonClass:
		return nestedType.AttachedObject.Get()
	default:
		return typ
	}
}

func (c *Checker) normaliseNilable(t *types.Nilable) types.Type {
	tType := c.NormaliseType(t.Type.Get())
	t.Type = types.ToRef(tType)
	switch tType.(type) {
	case types.Never:
		return types.Nil{}
	case types.Any, types.Untyped:
		return tType
	}
	if c.IsNilable(tType) {
		return tType
	}
	if union, ok := tType.(*types.Union); ok {
		return c.NewNormalisedUnion(
			append(
				[]types.Ref[types.Type]{types.NilID},
				union.Elements...,
			)...,
		)
	}
	return t
}

func (c *Checker) normaliseNot(t *types.Not) types.Type {
	tType := c.NormaliseType(t.Type.Get())
	t.Type = types.ToRef(tType)
	switch nestedType := tType.(type) {
	case *types.Not:
		return nestedType.Type.Get()
	case types.Never:
		return types.Any{}
	case types.Any:
		return types.Never{}
	case types.Untyped:
		return types.Untyped{}
	case *types.Union:
		intersectionElements := make([]types.Ref[types.Type], 0, len(nestedType.Elements))
		for _, element := range nestedType.Elements {
			intersectionElements = append(intersectionElements, types.CastRef[types.Type](types.NewNot(element.Get())))
		}
		return c.NewNormalisedIntersection(intersectionElements...)
	case *types.Intersection:
		unionElements := make([]types.Ref[types.Type], 0, len(nestedType.Elements))
		for _, element := range nestedType.Elements {
			unionElements = append(unionElements, types.CastRef[types.Type](types.NewNot(element.Get())))
		}
		return c.NewNormalisedUnion(unionElements...)
	}

	return t
}

func (c *Checker) normaliseGeneric(t *types.Generic) types.Type {
	for _, arg := range t.TypeArguments.AllArguments() {
		arg.Type = types.ToRef(c.NormaliseType(arg.Type.Get()))
	}
	return t
}

func (c *Checker) NormaliseType(typ types.Type) types.Type {
	switch t := typ.(type) {
	case *types.Union:
		return c.NewNormalisedUnion(t.Elements...)
	case *types.Intersection:
		return c.NewNormalisedIntersection(t.Elements...)
	case *types.Generic:
		return c.normaliseGeneric(t)
	case *types.SingletonOf:
		return c.normaliseSingletonOf(t)
	case *types.InstanceOf:
		return c.normaliseInstanceOf(t)
	case *types.Nilable:
		return c.normaliseNilable(t)
	case *types.Not:
		return c.normaliseNot(t)
	default:
		return typ
	}
}

func (c *Checker) distributeIntersectionOverUnions(newUnionElements *[]types.Ref[types.Type], intersectionElements []types.Ref[types.Type], i int) {
	if i == len(intersectionElements) {
		*newUnionElements = append(*newUnionElements, types.CastRef[types.Type](types.NewIntersectionRef(intersectionElements...)))
		return
	}

	intersectionElement := intersectionElements[i].Get()
	switch e := intersectionElement.(type) {
	case *types.Union:
		for _, subUnionElement := range e.Elements {
			newIntersectionElements := make([]types.Ref[types.Type], 0, len(intersectionElements)+1)
			newIntersectionElements = append(newIntersectionElements, intersectionElements[:i]...)
			newIntersectionElements = append(newIntersectionElements, subUnionElement)
			if len(intersectionElements) >= i+2 {
				newIntersectionElements = append(newIntersectionElements, intersectionElements[i+1:]...)
			}
			c.distributeIntersectionOverUnions(newUnionElements, newIntersectionElements, i+1)
		}
	case *types.Nilable:
		elements := []types.Ref[types.Type]{e.Type, types.NilID}
		for _, subUnionElement := range elements {
			newIntersectionElements := make([]types.Ref[types.Type], 0, len(intersectionElements)+1)
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
func (c *Checker) intersectionOfUnionsToUnionOfIntersections(intersectionElements []types.Ref[types.Type]) types.Type {
	newUnionElements := new([]types.Ref[types.Type])
	c.distributeIntersectionOverUnions(newUnionElements, intersectionElements, 0)
	if len(*newUnionElements) == 0 {
		return types.Never{}
	}
	if len(*newUnionElements) == 1 {
		return (*newUnionElements)[0].Get()
	}
	return types.NewUnionRef(*newUnionElements...)
}

func (c *Checker) NewNormalisedIntersection(elements ...types.Ref[types.Type]) types.Type {
	var containsNot bool
	var containsUninitialisedNamedTypes bool

	for i := 0; i < len(elements); i++ {
		element := c.NormaliseType(elements[i].Get())
		if types.IsNever(element) || types.IsUntyped(element) {
			return element
		}
		switch e := element.(type) {
		case *types.Intersection:
			newElements := make([]types.Ref[types.Type], 0, len(elements)+len(e.Elements))
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
			if e.Type.IsZero() {
				containsUninitialisedNamedTypes = true
			}
		}
	}
	if containsUninitialisedNamedTypes {
		return types.NewIntersectionRef(elements...)
	}
	if containsNot {
		// expand named types
		for i := 0; i < len(elements); i++ {
			switch e := elements[i].Get().(type) {
			case *types.Intersection:
				newElements := make([]types.Ref[types.Type], 0, len(elements)+len(e.Elements))
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
				elements[i] = types.CastRef[types.Type](types.NewUnionRef(types.TrueID, types.FalseID))
			}
		}
	}
	distributedIntersection := c.intersectionOfUnionsToUnionOfIntersections(elements)
	intersection, ok := distributedIntersection.(*types.Intersection)
	if !ok {
		return c.NormaliseType(distributedIntersection)
	}

	elements = intersection.Elements
	normalisedElements := make([]types.Ref[types.Type], 0, len(elements))

	// detect empty intersections
	for _, elementRef := range elements {
		element := elementRef.Get()
		if types.IsNever(element) || types.IsUntyped(element) {
			return element
		}

		for _, normalisedElement := range normalisedElements {
			if !c.canIntersect(element, normalisedElement.Get()) {
				return types.Never{}
			}
		}
		normalisedElements = append(normalisedElements, elementRef)
	}

	elements = normalisedElements
	normalisedElements = make([]types.Ref[types.Type], 0, len(elements))

eliminateSupertypesLoop:
	for i := 0; i < len(elements); i++ {
		element := c.NormaliseType(elements[i].Get())
		elementRef := types.ToRef(element)
		elements[i] = elementRef

		for j := 0; j < len(normalisedElements); j++ {
			normalisedElement := normalisedElements[j].Get()
			if c.isSubtype(normalisedElement, element, nil) {
				continue eliminateSupertypesLoop
			}
			if c.isSubtype(element, normalisedElement, nil) {
				normalisedElements[j] = elementRef
				continue eliminateSupertypesLoop
			}

		}
		normalisedElements = append(normalisedElements, elementRef)
	}

	if len(normalisedElements) == 0 {
		return types.Never{}
	}
	if len(normalisedElements) == 1 {
		return normalisedElements[0].Get()
	}

	return types.NewIntersectionRef(normalisedElements...)
}

func (c *Checker) NewNormalisedUnion(elements ...types.Ref[types.Type]) types.Type {
	var normalisedElements []types.Type

elementLoop:
	for i := 0; i < len(elements); i++ {
		element := c.NormaliseType(elements[i].Get())
		if types.IsNever(element) || types.IsUntyped(element) {
			continue elementLoop
		}
		switch e := element.(type) {
		case *types.Union:
			elements = append(elements, e.Elements...)
		case *types.Nilable:
			elements = append(elements, e.Type, types.NilID)
		case *types.Not:
			for j := 0; j < len(normalisedElements); j++ {
				normalisedElement := normalisedElements[j]
				if c.isTheSameType(e.Type.Get(), normalisedElement, nil) {
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
				if normalisedNot, ok := normalisedElement.(*types.Not); ok && c.isTheSameType(normalisedNot.Type.Get(), element, nil) {
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

		leftType := leftTypeNode.TypeRef()
		union.Elements = append(union.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.OR {
		c._constructUnionType(rightBinaryType, elements, union)
	} else {
		rightTypeNode := node.Right
		rightTypeNode = c.checkTypeNode(rightTypeNode)
		*elements = append(*elements, rightTypeNode)

		rightType := rightTypeNode.TypeRef()
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

		leftType := leftTypeNode.TypeRef()
		intersection.Elements = append(intersection.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(rightBinaryType, elements, intersection)
	} else {
		rightTypeNode := node.Right
		rightTypeNode = c.checkTypeNode(rightTypeNode)
		*elements = append(*elements, rightTypeNode)

		rightType := rightTypeNode.TypeRef()
		intersection.Elements = append(intersection.Elements, rightType)
	}
}
