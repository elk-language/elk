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
	return c.isSubtype(types.Nil{}, typ, nil)
}

// Type cannot be `nil`
func (c *Checker) IsNotNilable(typ types.Type) bool {
	return !c.IsNilable(typ)
}

// Type is always `nil`
func (c *Checker) IsNil(typ types.Type) bool {
	return types.IsNil(typ, c.env)
}

// Type is always `false`
func (c *Checker) IsFalse(typ types.Type) bool {
	return types.IsFalse(typ, c.env)
}

// Type is always `false`
func (c *Checker) IsTrue(typ types.Type) bool {
	return types.IsTrue(typ, c.env)
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
	return c.isSubtype(types.False{}, typ, nil) || c.isSubtype(types.Nil{}, typ, nil)
}

// Type can be truthy
func (c *Checker) CanBeTruthy(typ types.Type) bool {
	return !types.IsNever(c.NewNormalisedIntersection(typ, types.NewNot(types.False{}), types.NewNot(types.Nil{})))
}

func (c *Checker) IsTheSameType(a, b types.Type) bool {
	return c.isTheSameType(a, b, nil)
}

// Check whether the two given types represent the same type.
// Return true if they do, otherwise false.
func (c *Checker) isTheSameType(a, b types.Type, errLoc *position.Location) bool {
	return c.isSubtype(a, b, errLoc) && c.isSubtype(b, a, errLoc)
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
			return c.isSubtype(a, b, nil)
		}

		if !c.isSubtype(a.Namespace, genericB.Namespace, nil) {
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
				genericParent = c.replaceTypeParametersInGeneric(genericParent, g.ArgumentMap, false)
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
		return !c.isTheSameType(a.Type, b, nil)
	case *types.NamedType:
		return c._typesIntersect(a.Type, b, typeArgs)
	case *types.TypeParameter:
		return c.isSubtype(b, a.UpperBound, nil) && c.isSubtype(a.LowerBound, b, nil)
	case *types.Interface:
		return c.intersectsWithInterface(b, a, typeArgs)
	default:
		return c.isSubtype(a, b, nil)
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
		return c.isSubtype(a, b, nil)
	}

	if !aNamespace.IsGeneric() {
		return c.isSubtype(a, b, nil)
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
		return !c.isTheSameType(a.Type, b, nil)
	case *types.NamedType:
		return c.canBeIsA(a.Type, b)
	default:
		if bTypeParam, ok := b.(*types.TypeParameter); ok {
			return c.isSubtype(a, bTypeParam.UpperBound, nil) && c.isSubtype(bTypeParam.LowerBound, a, nil)
		}
		return c.isSubtype(a, b, nil)
	}
}

func (c *Checker) classCanIntersectWithMixin(a *types.Class, b *types.Mixin) bool {
	if c.phase == initPhase {
		return true
	}

	if c.isSubtype(a, b, nil) {
		return true
	}

	if types.NamespaceDeclaresInstanceVariables(b) && a.IsPrimitive() {
		return false
	}

	for mixinMethodName, mixinMethod := range c.methodsInNamespace(b) {
		classMethod := c.getMethod(a, mixinMethodName, nil)
		if classMethod == nil && mixinMethod.IsAbstract() && !a.IsAbstract() {
			return false
		}
		if classMethod == nil {
			continue
		}
		if !c.checkMethodCompatibility(mixinMethod, classMethod, nil, true) {
			return false
		}
	}

	return true
}

func (c *Checker) canIntersectWithInterfaceOrMixin(a types.Type, b types.Namespace) bool {
	if c.phase == initPhase {
		return true
	}

	if c.isSubtype(a, b, nil) {
		return true
	}

	for ifaceMethodName, ifaceMethod := range c.methodsInNamespace(b) {
		bMethod := c.getMethod(a, ifaceMethodName, nil)
		if bMethod == nil {
			continue
		}
		if !c.checkMethodCompatibility(ifaceMethod, bMethod, nil, true) {
			return false
		}
	}

	return true
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
	case *types.Interface:
		switch narrowB := b.(type) {
		case *types.Mixin, *types.Interface, *types.Class, *types.NamedType, *types.TypeParameter:
			return c.canIntersectWithInterfaceOrMixin(b, a)
		case *types.Generic:
			return c._canIntersect(a, narrowB.Namespace)
		default:
			return false
		}
	case *types.Mixin:
		switch narrowB := b.(type) {
		case *types.Class:
			return c.classCanIntersectWithMixin(narrowB, a)
		case *types.Mixin, *types.Interface, *types.NamedType, *types.TypeParameter:
			return c.canIntersectWithInterfaceOrMixin(b, a)
		case *types.Generic:
			return c._canIntersect(a, narrowB.Namespace)
		default:
			return false
		}
	case *types.NamedType:
		return c._canIntersect(a.Type, b)
	case *types.Not:
		return !c.isSubtype(b, a.Type, nil)
	case *types.Generic:
		switch a.Namespace.(type) {
		case *types.Mixin, *types.Interface:
			return c._canIntersect(a.Namespace, b)
		}
		genericB, ok := b.(*types.Generic)
		if !ok {
			return c.isSubtype(a, b, nil)
		}

		if !c.isSubtype(a.Namespace, genericB.Namespace, nil) {
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
				genericParent = c.replaceTypeParametersInGeneric(genericParent, g.ArgumentMap, false)
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
		return c.isSubtype(a, b, nil)
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

func (c *Checker) IsSubtype(a, b types.Type) bool {
	return c.isSubtype(a, b, nil)
}

func (c *Checker) isSubtype(a, b types.Type, errLoc *position.Location) bool {
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
		return c.isSubtype(a, narrowedB.Type, errLoc)
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
		return c.isSubtype(a.Type, b, errLoc)
	case *types.Union:
		for _, aElement := range a.Elements {
			if !c.isSubtype(aElement, b, errLoc) {
				return false
			}
		}
		return true
	case *types.Nilable:
		return c.isSubtype(a.Type, b, errLoc) && c.isSubtype(types.Nil{}, b, errLoc)
	case *types.Not:
		if bNot, ok := b.(*types.Not); ok {
			return c.isSubtype(bNot.Type, a.Type, nil)
		}
		return false
	case types.Self:
		return c.isSubtype(c.selfType, b, errLoc)
	case *types.TypeParameter:
		result, end := c.typeParameterIsSubtype(a, b, errLoc)
		if end {
			return result
		}
	}

	if bIntersection, ok := b.(*types.Intersection); ok {
		subtype := true
		for _, bElement := range bIntersection.Elements {
			if !c.isSubtype(a, bElement, errLoc) {
				subtype = false
			}
		}
		return subtype
	}

	switch b := b.(type) {
	case *types.Union:
		for _, bElement := range b.Elements {
			if c.isSubtype(a, bElement, nil) {
				return true
			}
		}
		return false
	case *types.Nilable:
		return c.isSubtype(a, b.Type, nil) || c.isSubtype(a, types.Nil{}, nil)
	case *types.Not:
		return !c.TypesIntersect(a, b.Type)
	case *types.TypeParameter:
		result, end := c.isSubtypeOfTypeParameter(a, b, errLoc)
		if end {
			return result
		}
	}

	if aIntersection, ok := a.(*types.Intersection); ok {
		for _, aElement := range aIntersection.Elements {
			if c.isSubtype(aElement, b, nil) {
				return true
			}
		}
		return false
	}

	aNonLiteral := c.ToNonLiteral(a, true)
	if a != aNonLiteral && c.isSubtype(aNonLiteral, b, errLoc) {
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
		return c.IsTrue(b)
	case types.False:
		return c.IsFalse(b)
	case *types.SingletonClass:
		return c.singletonClassIsSubtype(a, b, errLoc)
	case *types.Class:
		return c.classIsSubtype(a, b, errLoc)
	case *types.Mixin:
		return c.mixinIsSubtype(a, b, errLoc)
	case *types.MixinProxy:
		return c.mixinIsSubtype(a.Mixin, b, errLoc)
	case *types.Module:
		return c.moduleIsSubtype(a, b, errLoc)
	case *types.Interface:
		return c.interfaceIsSubtype(a, b, errLoc)
	case *types.InterfaceProxy:
		return c.interfaceIsSubtype(a.Interface, b, errLoc)
	case *types.Closure:
		return c.closureIsSubtype(a, b, errLoc)
	case *types.InstanceOf:
		switch narrowB := b.(type) {
		case *types.InstanceOf:
			return c.isSubtype(a.Type, narrowB.Type, errLoc)
		case *types.Class:
			return c.isSubtype(a.Type, narrowB.Singleton(), errLoc)
		case *types.Mixin:
			return c.isSubtype(a.Type, narrowB.Singleton(), errLoc)
		case *types.MixinProxy:
			return c.isSubtype(a.Type, narrowB.Singleton(), errLoc)
		case *types.Interface:
			return c.isSubtype(a.Type, narrowB.Singleton(), errLoc)
		case *types.InterfaceProxy:
			return c.isSubtype(a.Type, narrowB.Singleton(), errLoc)
		default:
			return false
		}
	case *types.SingletonOf:
		switch narrowB := b.(type) {
		case *types.SingletonOf:
			return c.isSubtype(a.Type, narrowB.Type, errLoc)
		case *types.SingletonClass:
			return c.isSubtype(a.Type, narrowB.AttachedObject, errLoc)
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
			return c.isSubtypeOfGeneric(a, narrowedB, errLoc)
		case *types.Interface:
			return c.isSubtypeOfInterface(a, narrowedB, errLoc)
		default:
			return c.isSubtype(a.Namespace, b, errLoc)
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

func (c *Checker) isSubtypeOfTypeParameter(a types.Type, b *types.TypeParameter, errLoc *position.Location) (result bool, end bool) {
	switch c.mode {
	case methodCompatibilityInAlgebraicTypeMode:
		if a, ok := a.(*types.TypeParameter); ok {
			if c.TypesIntersect(a.UpperBound, b.UpperBound) &&
				c.isTheSameType(b.LowerBound, a.LowerBound, nil) {
				return true, true
			}

			return false, true
		}
		if !c.isSubtype(b.LowerBound, a, nil) {
			return false, true
		}
		if c.isSubtype(b.UpperBound, a, nil) {
			return true, true
		}
		if c.isSubtype(a, b.UpperBound, nil) {
			return true, true
		}

		return false, true
	default:
		if c.isSubtype(a, b.LowerBound, errLoc) {
			return true, true
		}
	}

	return false, false
}

func (c *Checker) typeParameterIsSubtype(a *types.TypeParameter, b types.Type, errLoc *position.Location) (result bool, end bool) {
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
				c.isTheSameType(b.LowerBound, a.LowerBound, nil) {
				return true, true
			}

			return false, true
		}
		if !c.isSubtype(a.LowerBound, b, nil) {
			return false, true
		}
		if c.isSubtype(a.UpperBound, b, nil) {
			return true, true
		}
		if c.isSubtype(b, a.UpperBound, nil) {
			return true, true
		}

		return false, true
	default:
		if c.isSubtype(a.UpperBound, b, errLoc) {
			return true, true
		}
	}

	return false, false
}

func (c *Checker) typeArgsAreSubtype(a, b *types.TypeArguments, errLoc *position.Location) bool {
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
			if !c.isTheSameType(argA.Type, argB.Type, errLoc) {
				return false
			}
		case types.COVARIANT:
			if !c.isSubtype(argA.Type, argB.Type, errLoc) {
				return false
			}
		case types.CONTRAVARIANT:
			if !c.isSubtype(argB.Type, argA.Type, errLoc) {
				return false
			}
		case types.BIVARIANT:
			if !c.isSubtype(argB.Type, argA.Type, errLoc) && !c.isSubtype(argA.Type, argB.Type, errLoc) {
				return false
			}
		}
	}

	return true
}

func (c *Checker) singletonClassIsSubtype(a *types.SingletonClass, b types.Type, errLoc *position.Location) bool {
	switch b := b.(type) {
	case *types.SingletonClass:
		return c.isSubtype(a.AttachedObject, b.AttachedObject, errLoc)
	case *types.SingletonOf:
		return c.isSubtype(a.AttachedObject, b.Type, errLoc)
	case *types.Class:
		return c.isSubtypeOfClass(a, b)
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, b, errLoc)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errLoc)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errLoc)
	default:
		return false
	}
}

func (c *Checker) classIsSubtype(a *types.Class, b types.Type, errLoc *position.Location) bool {
	switch b := b.(type) {
	case *types.Class:
		return c.isSubtypeOfClass(a, b)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, b, errLoc)
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b)
	case *types.MixinProxy:
		return c.isSubtypeOfMixin(a, b.Mixin)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errLoc)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, b.Interface, errLoc)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errLoc)
	default:
		return false
	}
}

func (c *Checker) moduleIsSubtype(a *types.Module, b types.Type, errLoc *position.Location) bool {
	switch b := b.(type) {
	case *types.Class:
		return c.isSubtypeOfClass(a, b)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, b, errLoc)
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b)
	case *types.MixinProxy:
		return c.isSubtypeOfMixin(a, b.Mixin)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errLoc)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, b.Interface, errLoc)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errLoc)
	case *types.Module:
		return a == b
	default:
		return false
	}
}

func (c *Checker) mixinIsSubtype(a *types.Mixin, b types.Type, errLoc *position.Location) bool {
	switch b := b.(type) {
	case *types.Class:
		return c.isSubtypeOfClass(a, b)
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b)
	case *types.MixinProxy:
		return c.isSubtypeOfMixin(a, b.Mixin)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, b, errLoc)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errLoc)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, b.Interface, errLoc)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errLoc)
	default:
		return false
	}
}

func (c *Checker) isSubtypeOfGeneric(a types.Namespace, b *types.Generic, errLoc *position.Location) bool {
	isSubtype, shouldContinue := c.isSubtypeOfGenericNamespace(a, b, errLoc)
	if isSubtype {
		return true
	}
	if !shouldContinue {
		return false
	}

	switch b.Namespace.(type) {
	case *types.Interface:
		return c.isImplicitSubtypeOfInterface(a, b, errLoc)
	default:
		return false
	}
}

func (c *Checker) isSubtypeOfGenericNamespace(a types.Namespace, b *types.Generic, errLoc *position.Location) (isSubtype bool, shouldContinue bool) {
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
			target = c.replaceTypeParameters(target, m, false)
			targetGeneric := target.(*types.Generic)

			return c.typeArgsAreSubtype(targetGeneric.TypeArguments, b.TypeArguments, errLoc), false
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

func (c *Checker) namespaceIsMixin(a types.Namespace, b *types.Mixin) bool {
	switch a := a.(type) {
	case *types.Mixin:
		if a == b {
			return true
		}
	case *types.MixinProxy:
		if a.Mixin == b {
			return true
		}
	case *types.Generic:
		if c.IsTheSameNamespace(a.Namespace, b) {
			return true
		}
	case *types.TemporaryParent:
		return c.namespaceIsMixin(a.Namespace, b)
	}

	return false
}

func (c *Checker) includesMixin(a types.Namespace, b *types.Mixin) (bool, types.Namespace) {
	for parent := range types.Parents(a) {
		ok := c.namespaceIsMixin(parent, b)
		if ok {
			return true, parent
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

func (c *Checker) isSubtypeOfInterface(a types.Namespace, b *types.Interface, errLoc *position.Location) bool {
	if c.isExplicitSubtypeOfInterface(a, b) {
		return true
	}

	return c.isImplicitSubtypeOfInterface(a, b, errLoc)
}

func (c *Checker) isImplicitSubtypeOfInterface(a types.Namespace, b types.Namespace, errLoc *position.Location) bool {
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
			errLoc,
		)

		return false
	}

	return true
}

func (c *Checker) isSubtypeOfClosure(a types.Namespace, b *types.Closure, errLoc *position.Location) bool {
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
			errLoc,
		)

		return false
	}

	return true
}

func (c *Checker) interfaceIsSubtype(a *types.Interface, b types.Type, errLoc *position.Location) bool {
	switch narrowedB := b.(type) {
	case *types.Interface:
		return c.isSubtypeOfInterface(a, narrowedB, errLoc)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, narrowedB.Interface, errLoc)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, narrowedB, errLoc)
	case *types.Generic:
		return c.isSubtypeOfGeneric(a, narrowedB, errLoc)
	default:
		return false
	}
}

func (c *Checker) closureIsSubtype(a *types.Closure, b types.Type, errLoc *position.Location) bool {
	switch narrowedB := b.(type) {
	case *types.Interface:
		return c.isSubtypeOfInterface(a, narrowedB, errLoc)
	case *types.InterfaceProxy:
		return c.isSubtypeOfInterface(a, narrowedB.Interface, errLoc)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, narrowedB, errLoc)
	default:
		return false
	}
}
