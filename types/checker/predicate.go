package checker

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

// Type can be `nil`
func (c *Checker) IsNilable(typ types.Type) bool {
	return c.IsSubtype(types.Nil{}, typ, nil)
}

// Type cannot be `nil`
func (c *Checker) IsNotNilable(typ types.Type) bool {
	return !c.IsNilable(typ)
}

// Type is always `nil`
func (c *Checker) IsNil(typ types.Type) bool {
	return types.IsNil(typ, c.env)
}

// Type is always falsy.
func (c *Checker) IsFalsy(typ types.Type) bool {
	return !c.CanBeTruthy(typ)
}

// Type is always truthy.
func (c *Checker) IsTruthy(typ types.Type) bool {
	return !c.CanBeFalsy(typ)
}

// Type can be falsy
func (c *Checker) CanBeFalsy(typ types.Type) bool {
	return c.IsSubtype(types.False{}, typ, nil) || c.IsSubtype(types.Nil{}, typ, nil)
}

// Type can be truthy
func (c *Checker) CanBeTruthy(typ types.Type) bool {
	return !types.IsNever(c.NewNormalisedIntersection(typ, types.NewNot(types.False{}), types.NewNot(types.Nil{})))
}

// Check whether the two given types represent the same type.
// Return true if they do, otherwise false.
func (c *Checker) IsTheSameType(a, b types.Type, errSpan *position.Span) bool {
	return c.IsSubtype(a, b, errSpan) && c.IsSubtype(b, a, errSpan)
}

func (c *Checker) toInnerNamespace(a types.Namespace) types.Namespace {
	switch narrowedA := a.(type) {
	case *types.MixinProxy:
		return narrowedA.Mixin
	case *types.InterfaceProxy:
		return narrowedA.Interface
	default:
		return a
	}
}

// Check whether the two given types represent the same type.
// Return true if they do, otherwise false.
func (c *Checker) IsTheSameNamespace(a, b types.Namespace) bool {
	a = c.toInnerNamespace(a)
	b = c.toInnerNamespace(b)
	return a == b
}

// Check whether the two given types intersect.
// Return true if they do, otherwise false.
func (c *Checker) TypesIntersect(a, b types.Type) bool {
	return c.typesIntersectWithTypeArgs(a, b, nil)
}

func (c *Checker) typesIntersectWithTypeArgs(a, b types.Type, typeArgs types.TypeArgumentMap) bool {
	return c._typesIntersect(a, b, typeArgs) || c._typesIntersect(b, a, typeArgs)
}

func (c *Checker) _typesIntersect(a types.Type, b types.Type, typeArgs types.TypeArgumentMap) bool {
	switch a := a.(type) {
	case *types.Nilable:
		return c._typesIntersect(a.Type, b, typeArgs) || c._typesIntersect(types.Nil{}, b, typeArgs)
	case *types.Union:
		for _, element := range a.Elements {
			if c._typesIntersect(element, b, typeArgs) {
				return true
			}
		}
		return false
	case *types.Intersection:
		for _, element := range a.Elements {
			if c._typesIntersect(element, b, typeArgs) {
				return true
			}
		}
		return false
	case *types.Generic:
		if _, ok := a.Namespace.(*types.Interface); ok {
			return c.intersectsWithInterface(b, a, typeArgs)
		}
		genericB, ok := b.(*types.Generic)
		if !ok {
			return c.IsSubtype(a, b, nil)
		}

		if !c.IsSubtype(a.Namespace, genericB.Namespace, nil) {
			return false
		}

		var genericParents []*types.Generic
		for parent := range types.Parents(a) {
			genericParent, ok := parent.(*types.Generic)
			if !ok {
				continue
			}

			if !c.IsTheSameNamespace(genericParent.Namespace, genericB.Namespace) {
				genericParents = append(genericParents, genericParent)
				continue
			}

			for i := len(genericParents) - 1; i >= 0; i-- {
				g := genericParents[i]
				genericParent = c.replaceTypeParametersInGeneric(genericParent, g.ArgumentMap)
			}

			for name, argA := range genericParent.AllArguments() {
				argB := genericB.ArgumentMap[name]

				if !c._typesIntersect(argA.Type, argB.Type, typeArgs) {
					return false
				}
			}

			return true
		}

		return false
	case *types.Not:
		return !c.IsTheSameType(a.Type, b, nil)
	case *types.NamedType:
		return c._typesIntersect(a.Type, b, typeArgs)
	case *types.TypeParameter:
		return c.IsSubtype(b, a.UpperBound, nil) && c.IsSubtype(a.LowerBound, b, nil)
	case *types.Interface:
		return c.intersectsWithInterface(b, a, typeArgs)
	default:
		return c.IsSubtype(a, b, nil)
	}
}

func (c *Checker) intersectsWithInterface(a types.Type, b types.Namespace, typeArgs types.TypeArgumentMap) bool {
	var aNamespace types.Namespace
	switch a := a.(type) {
	case *types.Class:
		aNamespace = a
	case *types.Mixin:
		aNamespace = a
	default:
		return c.IsSubtype(a, b, nil)
	}

	if !aNamespace.IsGeneric() {
		return c.IsSubtype(a, b, nil)
	}

	// is a non-instantiated generic class/mixin

	for _, abstractMethod := range c.methodsInNamespace(b) {
		method := c.resolveMethodInNamespace(aNamespace, abstractMethod.Name)
		if method == nil || !c.checkMethodCompatibilityForInterfaceIntersection(abstractMethod, method, nil, typeArgs) {
			return false
		}
	}

	return true
}

// Check whether an "is a" relationship between `a` and `b` is possible.
func (c *Checker) canBeIsA(a types.Type, b types.Type) bool {
	switch a := a.(type) {
	case *types.Nilable:
		return c.canBeIsA(a.Type, b) || c.canBeIsA(types.Nil{}, b)
	case *types.Union:
		for _, element := range a.Elements {
			if c.canBeIsA(element, b) {
				return true
			}
		}
		return false
	case *types.Intersection:
		for _, element := range a.Elements {
			if c.canBeIsA(element, b) {
				return true
			}
		}
		return false
	case *types.Not:
		return !c.IsTheSameType(a.Type, b, nil)
	case *types.NamedType:
		return c.canBeIsA(a.Type, b)
	default:
		if bTypeParam, ok := b.(*types.TypeParameter); ok {
			return c.IsSubtype(a, bTypeParam.UpperBound, nil) && c.IsSubtype(bTypeParam.LowerBound, a, nil)
		}
		return c.IsSubtype(a, b, nil)
	}
}

// Check whether the two given types can potentially intersect.
// Return true if they do, otherwise false.
func (c *Checker) canIntersect(a, b types.Type) bool {
	return c._canIntersect(a, b) || c._canIntersect(b, a)
}

func (c *Checker) _canIntersect(a types.Type, b types.Type) bool {
	switch a := a.(type) {
	case *types.Nilable:
		return c.canBeIsA(a.Type, b) || c.canBeIsA(types.Nil{}, b)
	case *types.Union:
		for _, element := range a.Elements {
			if c.canBeIsA(element, b) {
				return true
			}
		}
		return false
	case *types.Intersection:
		for _, element := range a.Elements {
			if c.canBeIsA(element, b) {
				return true
			}
		}
		return false
	case *types.Mixin, *types.Interface:
		switch narrowB := b.(type) {
		case *types.Mixin, *types.Interface, *types.Class, *types.NamedType, *types.TypeParameter:
			return true
		case *types.Generic:
			return c._canIntersect(a, narrowB.Namespace)
		default:
			return false
		}
	case *types.NamedType:
		return c._canIntersect(a.Type, b)
	case *types.Not:
		return !c.IsSubtype(b, a.Type, nil)
	case *types.Generic:
		switch a.Namespace.(type) {
		case *types.Mixin, *types.Interface:
			return c._canIntersect(a.Namespace, b)
		}
		genericB, ok := b.(*types.Generic)
		if !ok {
			return c.IsSubtype(a, b, nil)
		}

		if !c.IsSubtype(a.Namespace, genericB.Namespace, nil) {
			return false
		}

		var genericParents []*types.Generic
		for parent := range types.Parents(a) {
			genericParent, ok := parent.(*types.Generic)
			if !ok {
				continue
			}

			if !c.IsTheSameNamespace(genericParent.Namespace, genericB.Namespace) {
				genericParents = append(genericParents, genericParent)
				continue
			}

			for i := len(genericParents) - 1; i >= 0; i-- {
				g := genericParents[i]
				genericParent = c.replaceTypeParametersInGeneric(genericParent, g.ArgumentMap)
			}

			for name, argA := range genericParent.AllArguments() {
				argB := genericB.ArgumentMap[name]

				if !c.canIntersect(argA.Type, argB.Type) {
					return false
				}
			}

			return true
		}

		return false
	default:
		return c.IsSubtype(a, b, nil)
	}
}

func (c *Checker) containsTypeParameters(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.SingletonOf:
		return c.containsTypeParameters(t.Type)
	case *types.InstanceOf:
		return c.containsTypeParameters(t.Type)
	case *types.Closure:
		for _, param := range t.Body.Params {
			if c.containsTypeParameters(param.Type) {
				return true
			}
		}

		return c.containsTypeParameters(t.Body.ReturnType) || c.containsTypeParameters(t.Body.ThrowType)
	case *types.Generic:
		for _, arg := range t.AllArguments() {
			if c.containsTypeParameters(arg.Type) {
				return true
			}
		}
		return c.containsTypeParameters(t.Namespace)
	case *types.TypeParameter:
		return true
	case *types.Nilable:
		return c.containsTypeParameters(t.Type)
	case *types.Not:
		return c.containsTypeParameters(t.Type)
	case *types.Union:
		for _, element := range t.Elements {
			if c.containsTypeParameters(element) {
				return true
			}
		}
		return false
	case *types.Intersection:
		for _, element := range t.Elements {
			if c.containsTypeParameters(element) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func (c *Checker) typesAreIdentical(a, b types.Type) bool {
	if a == b {
		return true
	}

	if a, ok := a.(*types.Generic); ok {
		b, ok := b.(*types.Generic)
		if ok {
			result := a.Namespace == b.Namespace
			return result
		}
	}

	return false
}

func (c *Checker) IsSubtype(a, b types.Type, errSpan *position.Span) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	if a == nil && b == nil {
		return true
	}

	if c.mode == implicitInterfaceSubtypeMode && c.typesAreIdentical(c.selfType, a) && c.typesAreIdentical(c.throwType, b) {
		return true
	}

	if bNamedType, ok := b.(*types.NamedType); ok {
		b = bNamedType.Type
	}

	if types.IsNever(a) || types.IsUntyped(a) {
		return true
	}
	switch narrowedB := b.(type) {
	case *types.NamedType:
		return c.IsSubtype(a, narrowedB.Type, errSpan)
	case types.Any, types.Void, types.Untyped:
		return true
	case types.Nil:
		b = c.StdNil()
	case types.Bool:
		b = c.StdBool()
	case types.True:
		b = c.StdTrue()
	case types.False:
		b = c.StdFalse()
	case types.Self:
		return types.IsSelf(a)
	}

	if types.IsAny(a) || types.IsVoid(a) {
		return false
	}

	switch a := a.(type) {
	case *types.NamedType:
		return c.IsSubtype(a.Type, b, errSpan)
	case *types.Union:
		for _, aElement := range a.Elements {
			if !c.IsSubtype(aElement, b, errSpan) {
				return false
			}
		}
		return true
	case *types.Nilable:
		return c.IsSubtype(a.Type, b, errSpan) && c.IsSubtype(types.Nil{}, b, errSpan)
	case *types.Not:
		if bNot, ok := b.(*types.Not); ok {
			return c.IsSubtype(bNot.Type, a.Type, nil)
		}
		return false
	case types.Self:
		return c.IsSubtype(c.selfType, b, errSpan)
	case *types.TypeParameter:
		result, end := c.typeParameterIsSubtype(a, b, errSpan)
		if end {
			return result
		}
	}

	if bIntersection, ok := b.(*types.Intersection); ok {
		subtype := true
		for _, bElement := range bIntersection.Elements {
			if !c.IsSubtype(a, bElement, errSpan) {
				subtype = false
			}
		}
		return subtype
	}

	switch b := b.(type) {
	case *types.Union:
		for _, bElement := range b.Elements {
			if c.IsSubtype(a, bElement, nil) {
				return true
			}
		}
		return false
	case *types.Nilable:
		return c.IsSubtype(a, b.Type, nil) || c.IsSubtype(a, types.Nil{}, nil)
	case *types.Not:
		return !c.TypesIntersect(a, b.Type)
	case *types.TypeParameter:
		result, end := c.isSubtypeOfTypeParameter(a, b, errSpan)
		if end {
			return result
		}
	}

	if aIntersection, ok := a.(*types.Intersection); ok {
		for _, aElement := range aIntersection.Elements {
			if c.IsSubtype(aElement, b, nil) {
				return true
			}
		}
		return false
	}

	aNonLiteral := c.ToNonLiteral(a, true)
	if a != aNonLiteral && c.IsSubtype(aNonLiteral, b, errSpan) {
		return true
	}

	originalA := a
	switch a := a.(type) {
	case types.Any:
		return types.IsAny(b)
	case types.Nil:
		return types.IsNilLiteral(b) || b == c.StdNil()
	case types.Bool:
		return types.IsBool(b) || b == c.StdBool()
	case types.True:
		return types.IsTrue(b) || b == c.StdTrue()
	case types.False:
		return types.IsFalse(b) || b == c.StdFalse()
	case *types.SingletonClass:
		return c.singletonClassIsSubtype(a, b, errSpan)
	case *types.Class:
		return c.classIsSubtype(a, b, errSpan)
	case *types.Mixin:
		return c.mixinIsSubtype(a, b, errSpan)
	case *types.MixinProxy:
		return c.mixinIsSubtype(a.Mixin, b, errSpan)
	case *types.Module:
		return c.moduleIsSubtype(a, b, errSpan)
	case *types.Interface:
		return c.interfaceIsSubtype(a, b, errSpan)
	case *types.InterfaceProxy:
		return c.interfaceIsSubtype(a.Interface, b, errSpan)
	case *types.Closure:
		return c.closureIsSubtype(a, b, errSpan)
	case *types.InstanceOf:
		switch narrowB := b.(type) {
		case *types.InstanceOf:
			return c.IsSubtype(a.Type, narrowB.Type, errSpan)
		case *types.Class:
			return c.IsSubtype(a.Type, narrowB.Singleton(), errSpan)
		case *types.Mixin:
			return c.IsSubtype(a.Type, narrowB.Singleton(), errSpan)
		case *types.MixinProxy:
			return c.IsSubtype(a.Type, narrowB.Singleton(), errSpan)
		case *types.Interface:
			return c.IsSubtype(a.Type, narrowB.Singleton(), errSpan)
		case *types.InterfaceProxy:
			return c.IsSubtype(a.Type, narrowB.Singleton(), errSpan)
		default:
			return false
		}
	case *types.SingletonOf:
		switch narrowB := b.(type) {
		case *types.SingletonOf:
			return c.IsSubtype(a.Type, narrowB.Type, errSpan)
		case *types.SingletonClass:
			return c.IsSubtype(a.Type, narrowB.AttachedObject, errSpan)
		default:
			return false
		}
	case *types.TypeParameter:
		b, ok := b.(*types.TypeParameter)
		if !ok {
			return false
		}
		return a.Name == b.Name
	case *types.Generic:
		switch narrowedB := b.(type) {
		case *types.Generic:
			return c.isSubtypeOfGeneric(a, narrowedB, errSpan)
		case *types.Interface:
			return c.isSubtypeOfInterface(a, narrowedB, errSpan)
		default:
			return c.IsSubtype(a.Namespace, b, errSpan)
		}
	case *types.Method:
		b, ok := b.(*types.Method)
		if !ok {
			return false
		}
		return a == b
	case *types.CharLiteral:
		b, ok := b.(*types.CharLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.StringLiteral:
		b, ok := b.(*types.StringLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.SymbolLiteral:
		b, ok := b.(*types.SymbolLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.FloatLiteral:
		b, ok := b.(*types.FloatLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.Float64Literal:
		b, ok := b.(*types.Float64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.Float32Literal:
		b, ok := b.(*types.Float32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.BigFloatLiteral:
		b, ok := b.(*types.BigFloatLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.IntLiteral:
		b, ok := b.(*types.IntLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.Int64Literal:
		b, ok := b.(*types.Int64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.Int32Literal:
		b, ok := b.(*types.Int32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.Int16Literal:
		b, ok := b.(*types.Int16Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.Int8Literal:
		b, ok := b.(*types.Int8Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.UInt64Literal:
		b, ok := b.(*types.UInt64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.UInt32Literal:
		b, ok := b.(*types.UInt32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.UInt16Literal:
		b, ok := b.(*types.UInt16Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	case *types.UInt8Literal:
		b, ok := b.(*types.UInt8Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value && a.IsNegative() == b.IsNegative()
	default:
		panic(fmt.Sprintf("invalid type: %T", originalA))
	}
}

func (c *Checker) isSubtypeOfTypeParameter(a types.Type, b *types.TypeParameter, errSpan *position.Span) (result bool, end bool) {
	switch c.mode {
	case methodCompatibilityInAlgebraicTypeMode:
		if a, ok := a.(*types.TypeParameter); ok {
			if c.TypesIntersect(a.UpperBound, b.UpperBound) &&
				c.IsTheSameType(b.LowerBound, a.LowerBound, nil) {
				return true, true
			}

			return false, true
		}
		if !c.IsSubtype(b.LowerBound, a, nil) {
			return false, true
		}
		if c.IsSubtype(b.UpperBound, a, nil) {
			return true, true
		}
		if c.IsSubtype(a, b.UpperBound, nil) {
			return true, true
		}

		return false, true
	default:
		if c.IsSubtype(a, b.LowerBound, errSpan) {
			return true, true
		}
	}

	return false, false
}

func (c *Checker) typeParameterIsSubtype(a *types.TypeParameter, b types.Type, errSpan *position.Span) (result bool, end bool) {
	switch c.mode {
	case inferTypeArgumentMode:
		b, ok := b.(*types.TypeParameter)
		if !ok {
			return false, true
		}
		return a.Name == b.Name, true
	case methodCompatibilityInAlgebraicTypeMode:
		if b, ok := b.(*types.TypeParameter); ok {
			if c.TypesIntersect(a.UpperBound, b.UpperBound) &&
				c.IsTheSameType(b.LowerBound, a.LowerBound, nil) {
				return true, true
			}

			return false, true
		}
		if !c.IsSubtype(a.LowerBound, b, nil) {
			return false, true
		}
		if c.IsSubtype(a.UpperBound, b, nil) {
			return true, true
		}
		if c.IsSubtype(b, a.UpperBound, nil) {
			return true, true
		}

		return false, true
	default:
		if c.IsSubtype(a.UpperBound, b, errSpan) {
			return true, true
		}
	}

	return false, false
}

func (c *Checker) typeArgsAreSubtype(a, b *types.TypeArguments, errSpan *position.Span) bool {
	for i := range b.ArgumentOrder {
		argB := b.ArgumentMap[b.ArgumentOrder[i]]
		argA := a.ArgumentMap[a.ArgumentOrder[i]]

		var variance types.Variance
		if argA.Variance > argB.Variance {
			variance = argA.Variance
		} else {
			variance = argB.Variance
		}

		switch variance {
		case types.INVARIANT:
			if !c.IsTheSameType(argA.Type, argB.Type, errSpan) {
				return false
			}
		case types.COVARIANT:
			if !c.IsSubtype(argA.Type, argB.Type, errSpan) {
				return false
			}
		case types.CONTRAVARIANT:
			if !c.IsSubtype(argB.Type, argA.Type, errSpan) {
				return false
			}
		case types.BIVARIANT:
			if !c.IsSubtype(argB.Type, argA.Type, errSpan) && !c.IsSubtype(argA.Type, argB.Type, errSpan) {
				return false
			}
		}
	}

	return true
}

func (c *Checker) singletonClassIsSubtype(a *types.SingletonClass, b types.Type, errSpan *position.Span) bool {
	switch b := b.(type) {
	case *types.SingletonClass:
		return c.IsSubtype(a.AttachedObject, b.AttachedObject, errSpan)
	case *types.SingletonOf:
		return c.IsSubtype(a.AttachedObject, b.Type, errSpan)
	case *types.Class:
		return c.isSubtypeOfClass(a, b)
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, b, errSpan)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errSpan)
	default:
		return false
	}
}

func (c *Checker) classIsSubtype(a *types.Class, b types.Type, errSpan *position.Span) bool {
	switch b := b.(type) {
	case *types.Class:
		return c.isSubtypeOfClass(a, b)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, b, errSpan)
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b)
	case *types.MixinProxy:
		return c.isSubtypeOfMixin(a, b.Mixin)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errSpan)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, b.Interface, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errSpan)
	default:
		return false
	}
}

func (c *Checker) moduleIsSubtype(a *types.Module, b types.Type, errSpan *position.Span) bool {
	switch b := b.(type) {
	case *types.Class:
		return c.isSubtypeOfClass(a, b)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, b, errSpan)
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b)
	case *types.MixinProxy:
		return c.isSubtypeOfMixin(a, b.Mixin)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errSpan)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, b.Interface, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errSpan)
	case *types.Module:
		return a == b
	default:
		return false
	}
}

func (c *Checker) mixinIsSubtype(a *types.Mixin, b types.Type, errSpan *position.Span) bool {
	switch b := b.(type) {
	case *types.Class:
		return c.isSubtypeOfClass(a, b)
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b)
	case *types.MixinProxy:
		return c.isSubtypeOfMixin(a, b.Mixin)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, b, errSpan)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errSpan)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, b.Interface, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errSpan)
	default:
		return false
	}
}

func (c *Checker) isSubtypeOfGeneric(a types.Namespace, b *types.Generic, errSpan *position.Span) bool {
	isSubtype, shouldContinue := c.isSubtypeOfGenericNamespace(a, b, errSpan)
	if isSubtype {
		return true
	}
	if !shouldContinue {
		return false
	}

	switch b.Namespace.(type) {
	case *types.Interface:
		return c.isImplicitSubtypeOfInterface(a, b, errSpan)
	default:
		return false
	}
}

func (c *Checker) isSubtypeOfGenericNamespace(a types.Namespace, b *types.Generic, errSpan *position.Span) (isSubtype bool, shouldContinue bool) {
	var generics []*types.Generic

	for parent := range types.Parents(a) {
		parent, ok := parent.(*types.Generic)
		if !ok {
			continue
		}

		if c.IsTheSameNamespace(parent.Namespace, b.Namespace) {
			var target types.Type = parent
			for _, generic := range slices.Backward(generics) {
				target = c.replaceTypeParametersOfGeneric(target, generic)
				if target == nil {
					return false, false
				}
			}
			m := c.createTypeArgumentMapWithSelf(a)
			target = c.replaceTypeParameters(target, m)
			targetGeneric := target.(*types.Generic)

			return c.typeArgsAreSubtype(targetGeneric.TypeArguments, b.TypeArguments, errSpan), false
		}
		generics = append(generics, parent)
	}

	return false, true
}

func (c *Checker) isSubtypeOfClass(a types.Namespace, b *types.Class) bool {
	var currentParent types.Namespace = a
	for {
		if currentParent == b {
			return true
		}
		switch p := currentParent.(type) {
		case nil:
			return false
		case *types.Generic:
			if p.Namespace == b {
				return true
			}
		}

		currentParent = currentParent.Parent()
	}
}

func (c *Checker) isSubtypeOfMixin(a types.Namespace, b *types.Mixin) bool {
	ok, _ := c.includesMixin(a, b)
	return ok
}

func (c *Checker) includesMixin(a types.Namespace, b *types.Mixin) (bool, types.Namespace) {
	for parent := range types.Parents(a) {
		switch p := parent.(type) {
		case *types.Mixin:
			if p == b {
				return true, p
			}
		case *types.MixinProxy:
			if p.Mixin == b {
				return true, p
			}
		case *types.Generic:
			if c.IsTheSameNamespace(p.Namespace, b) {
				return true, p
			}
		}
	}

	return false, nil
}

type methodOverride struct {
	superMethod *types.Method
	override    *types.Method
}

func (c *Checker) isExplicitSubtypeOfInterface(a types.Namespace, b *types.Interface) bool {
	for parent := range types.Parents(a) {
		switch p := parent.(type) {
		case *types.Interface:
			if p == b {
				return true
			}
		case *types.InterfaceProxy:
			if p.Interface == b {
				return true
			}
		case *types.Generic:
			if c.IsTheSameNamespace(p.Namespace, b) {
				return true
			}
		}
	}

	return false
}

func (c *Checker) isSubtypeOfInterface(a types.Namespace, b *types.Interface, errSpan *position.Span) bool {
	if c.isExplicitSubtypeOfInterface(a, b) {
		return true
	}

	return c.isImplicitSubtypeOfInterface(a, b, errSpan)
}

func (c *Checker) isImplicitSubtypeOfInterface(a types.Namespace, b types.Namespace, errSpan *position.Span) bool {
	if c.phase == initPhase && len(b.Methods()) < 1 {
		return false
	}

	// prevent infinite loops of subtype checking
	// for code like:
	//
	//     class Foo
	//       def foo: Foo then loop; end
	//     end
	//     interface Bar
	// 	     def foo: Bar; end
	//     end
	//     var a: Bar = Foo()
	//
	prevMode := c.mode
	c.mode = implicitInterfaceSubtypeMode

	// use `selfType` to store type `a`
	prevSelf := c.selfType
	c.selfType = a

	// use `throwType` to store type `b` (the interface)
	prevThrow := c.throwType
	c.throwType = b

	var incorrectMethods []methodOverride
	for _, abstractMethod := range c.methodsInNamespace(b) {
		method := c.resolveMethodInNamespace(a, abstractMethod.Name)
		if method == nil || !c.checkMethodCompatibility(abstractMethod, method, nil, true) {
			incorrectMethods = append(incorrectMethods, methodOverride{
				superMethod: abstractMethod,
				override:    method,
			})
		}
	}

	c.throwType = prevThrow
	c.selfType = prevSelf
	c.mode = prevMode

	if len(incorrectMethods) > 0 {
		methodDetailsBuff := new(strings.Builder)
		for _, incorrectMethod := range incorrectMethods {
			implementation := incorrectMethod.override
			abstractMethod := incorrectMethod.superMethod
			if implementation == nil {
				fmt.Fprintf(
					methodDetailsBuff,
					"\n  - missing method `%s` with signature: `%s`",
					types.InspectWithColor(abstractMethod),
					abstractMethod.InspectSignatureWithColor(false),
				)
				continue
			}

			fmt.Fprintf(
				methodDetailsBuff,
				"\n  - incorrect implementation of `%s`\n      is:        `%s`\n      should be: `%s`",
				types.InspectWithColor(abstractMethod),
				implementation.InspectSignatureWithColor(false),
				abstractMethod.InspectSignatureWithColor(false),
			)
		}

		c.addFailure(
			fmt.Sprintf(
				"type `%s` does not implement interface `%s`:\n%s",
				types.InspectWithColor(a),
				types.InspectWithColor(b),
				methodDetailsBuff.String(),
			),
			errSpan,
		)

		return false
	}

	return true
}

func (c *Checker) isSubtypeOfClosure(a types.Namespace, b *types.Closure, errSpan *position.Span) bool {
	abstractMethod := b.Body
	method := c.resolveMethodInNamespace(a, symbol.L_call)

	if method == nil || !c.checkMethodCompatibility(abstractMethod, method, nil, false) {
		methodDetailsBuff := new(strings.Builder)
		if method == nil {
			fmt.Fprintf(
				methodDetailsBuff,
				"\n  - missing method `%s` with signature: `%s`\n",
				types.InspectWithColor(abstractMethod),
				abstractMethod.InspectSignatureWithColor(false),
			)
		} else {
			fmt.Fprintf(
				methodDetailsBuff,
				"\n  - incorrect implementation of `%s`\n      is:        `%s`\n      should be: `%s`\n",
				types.InspectWithColor(abstractMethod),
				method.InspectSignatureWithColor(false),
				abstractMethod.InspectSignatureWithColor(false),
			)
		}

		c.addFailure(
			fmt.Sprintf(
				"type `%s` does not implement closure `%s`:\n%s",
				types.InspectWithColor(a),
				types.InspectWithColor(b),
				methodDetailsBuff.String(),
			),
			errSpan,
		)

		return false
	}

	return true
}

func (c *Checker) interfaceIsSubtype(a *types.Interface, b types.Type, errSpan *position.Span) bool {
	switch narrowedB := b.(type) {
	case *types.Interface:
		return c.isSubtypeOfInterface(a, narrowedB, errSpan)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, narrowedB.Interface, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, narrowedB, errSpan)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, narrowedB, errSpan)
	default:
		return false
	}
}

func (c *Checker) closureIsSubtype(a *types.Closure, b types.Type, errSpan *position.Span) bool {
	switch narrowedB := b.(type) {
	case *types.Interface:
		return c.isSubtypeOfInterface(a, narrowedB, errSpan)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, narrowedB.Interface, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, narrowedB, errSpan)
	default:
		return false
	}
}
