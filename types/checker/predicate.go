package checker

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

// Type can be `nil`
func (c *Checker) isNilable(typ types.Type) bool {
	return types.IsNilable(typ, c.GlobalEnv)
}

// Type cannot be `nil`
func (c *Checker) isNotNilable(typ types.Type) bool {
	return !types.IsNilable(typ, c.GlobalEnv)
}

// Type is always `nil`
func (c *Checker) isNil(typ types.Type) bool {
	return types.IsNil(typ, c.GlobalEnv)
}

// Type is always falsy.
func (c *Checker) isFalsy(typ types.Type) bool {
	return !c.canBeTruthy(typ)
}

// Type is always truthy.
func (c *Checker) isTruthy(typ types.Type) bool {
	return !c.canBeFalsy(typ)
}

// Type can be falsy
func (c *Checker) canBeFalsy(typ types.Type) bool {
	return types.CanBeFalsy(typ, c.GlobalEnv)
}

// Type can be truthy
func (c *Checker) canBeTruthy(typ types.Type) bool {
	return types.CanBeTruthy(typ, c.GlobalEnv)
}

// Check whether the two given types represent the same type.
// Return true if they do, otherwise false.
func (c *Checker) isTheSameType(a, b types.Type, errSpan *position.Span) bool {
	return c.isSubtype(a, b, errSpan) && c.isSubtype(b, a, errSpan)
}

// Check whether the two given types intersect.
// Return true if they do, otherwise false.
func (c *Checker) typesIntersect(a, b types.Type) bool {
	return c.canBeIsA(a, b) || c.canBeIsA(b, a)
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
		return c.isSubtype(a, b, nil)
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
		switch b.(type) {
		case *types.Mixin, *types.Interface, *types.Class, *types.NamedType, *types.TypeParameter:
			return true
		}
		return false
	case *types.NamedType:
		return c._canIntersect(a.Type, b)
	case *types.Not:
		return !c.isSubtype(b, a.Type, nil)
	default:
		return c.isSubtype(a, b, nil)
	}
}

func (c *Checker) isSubtype(a, b types.Type, errSpan *position.Span) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	if a == nil && b == nil {
		return true
	}

	if bNamedType, ok := b.(*types.NamedType); ok {
		b = bNamedType.Type
	}

	if types.IsNever(a) || types.IsNothing(a) {
		return true
	}
	switch b.(type) {
	case types.Any, types.Void, types.Nothing:
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
	case *types.Union:
		for _, aElement := range a.Elements {
			if !c.isSubtype(aElement, b, errSpan) {
				return false
			}
		}
		return true
	case *types.Nilable:
		return c.isSubtype(a.Type, b, errSpan) && c.isSubtype(types.Nil{}, b, errSpan)
	case *types.Not:
		if bNot, ok := b.(*types.Not); ok {
			return c.isSubtype(bNot.Type, a.Type, nil)
		}
		return false
	case types.Self:
		return c.isSubtype(c.selfType, b, errSpan)
	case *types.TypeParameter:
		if c.mode == inferTypeArgumentMode {
			b, ok := b.(*types.TypeParameter)
			if !ok {
				return false
			}
			return a.Name == b.Name
		}
		if c.isSubtype(a.UpperBound, b, errSpan) {
			return true
		}
	}

	if bIntersection, ok := b.(*types.Intersection); ok {
		subtype := true
		for _, bElement := range bIntersection.Elements {
			if !c.isSubtype(a, bElement, errSpan) {
				subtype = false
			}
		}
		return subtype
	}

	switch b := b.(type) {
	case *types.Union:
		for _, bElement := range b.Elements {
			if c.isSubtype(a, bElement, errSpan) {
				return true
			}
		}
		return false
	case *types.Nilable:
		return c.isSubtype(a, b.Type, errSpan) || c.isSubtype(a, types.Nil{}, errSpan)
	case *types.Not:
		return !c.typesIntersect(a, b.Type)
	case *types.TypeParameter:
		if c.isSubtype(a, b.LowerBound, errSpan) {
			return true
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

	aNonLiteral := c.toNonLiteral(a, true)
	if a != aNonLiteral && c.isSubtype(aNonLiteral, b, errSpan) {
		return true
	}

	originalA := a
	switch a := a.(type) {
	case *types.NamedType:
		return c.isSubtype(a.Type, b, errSpan)
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
	case *types.Module:
		return c.moduleIsSubtype(a, b, errSpan)
	case *types.Interface:
		return c.interfaceIsSubtype(a, b, errSpan)
	case *types.Closure:
		return c.closureIsSubtype(a, b, errSpan)
	case *types.InstanceOf:
		switch narrowB := b.(type) {
		case *types.InstanceOf:
			return c.isSubtype(a.Type, narrowB.Type, errSpan)
		case *types.Class:
			return c.isSubtype(a.Type, narrowB.Singleton(), errSpan)
		case *types.Mixin:
			return c.isSubtype(a.Type, narrowB.Singleton(), errSpan)
		case *types.Interface:
			return c.isSubtype(a.Type, narrowB.Singleton(), errSpan)
		default:
			return false
		}
	case *types.SingletonOf:
		switch narrowB := b.(type) {
		case *types.SingletonOf:
			return c.isSubtype(a.Type, narrowB.Type, errSpan)
		case *types.SingletonClass:
			return c.isSubtype(a.Type, narrowB.AttachedObject, errSpan)
		default:
			return false
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
		return a.Value == b.Value
	case *types.Float64Literal:
		b, ok := b.(*types.Float64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Float32Literal:
		b, ok := b.(*types.Float32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.BigFloatLiteral:
		b, ok := b.(*types.BigFloatLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.IntLiteral:
		b, ok := b.(*types.IntLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int64Literal:
		b, ok := b.(*types.Int64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int32Literal:
		b, ok := b.(*types.Int32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int16Literal:
		b, ok := b.(*types.Int16Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int8Literal:
		b, ok := b.(*types.Int8Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt64Literal:
		b, ok := b.(*types.UInt64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt32Literal:
		b, ok := b.(*types.UInt32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt16Literal:
		b, ok := b.(*types.UInt16Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt8Literal:
		b, ok := b.(*types.UInt8Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.TypeParameter:
		b, ok := b.(*types.TypeParameter)
		if !ok {
			return false
		}
		return a.Name == b.Name
	case *types.Generic:
		genericB, ok := b.(*types.Generic)
		if !ok {
			return c.isSubtype(a.Namespace, b, errSpan)
		}
		return c.isSubtype(a.Namespace, genericB.Namespace, errSpan) &&
			c.typeArgsAreSubtype(a.TypeArguments, genericB.TypeArguments, errSpan)
	default:
		panic(fmt.Sprintf("invalid type: %T", originalA))
	}
}

func (c *Checker) typeArgsAreSubtype(a, b *types.TypeArguments, errSpan *position.Span) bool {
	for i := range b.ArgumentOrder {
		argB := b.ArgumentMap[b.ArgumentOrder[i]]
		argA := a.ArgumentMap[a.ArgumentOrder[i]]

		switch argA.Variance {
		case types.INVARIANT:
			if !c.isTheSameType(argA.Type, argB.Type, errSpan) {
				return false
			}
		case types.COVARIANT:
			if !c.isSubtype(argA.Type, argB.Type, errSpan) {
				return false
			}
		case types.CONTRAVARIANT:
			if !c.isSubtype(argB.Type, argA.Type, errSpan) {
				return false
			}
		}
	}

	return true
}

func (c *Checker) singletonClassIsSubtype(a *types.SingletonClass, b types.Type, errSpan *position.Span) bool {
	switch b := b.(type) {
	case *types.SingletonClass:
		return c.isSubtype(a.AttachedObject, b.AttachedObject, errSpan)
	case *types.SingletonOf:
		return c.isSubtype(a.AttachedObject, b.Type, errSpan)
	case *types.Class:
		var currentClass types.Namespace = a
		for {
			if currentClass == nil {
				return false
			}
			if currentClass == b {
				return true
			}

			currentClass = currentClass.Parent()
		}
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b, errSpan)
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
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b, errSpan)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errSpan)
	default:
		return false
	}
}

func (c *Checker) moduleIsSubtype(a *types.Module, b types.Type, errSpan *position.Span) bool {
	switch b := b.(type) {
	case *types.Class:
		var currentClass types.Namespace = a
		for {
			if currentClass == nil {
				return false
			}
			if currentClass == b {
				return true
			}

			currentClass = currentClass.Parent()
		}
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b, errSpan)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errSpan)
	case *types.Module:
		return a == b
	default:
		return false
	}
}

func (c *Checker) isSubtypeOfMixin(a types.Namespace, b *types.Mixin, errSpan *position.Span) bool {
	var currentContainer types.Namespace = a
	for {
		switch cont := currentContainer.(type) {
		case *types.Mixin:
			if cont == b {
				return true
			}
		case *types.MixinProxy:
			if cont.Mixin == b {
				return true
			}
		case nil:
			return false
		}

		currentContainer = currentContainer.Parent()
	}
}

func (c *Checker) mixinIsSubtype(a *types.Mixin, b types.Type, errSpan *position.Span) bool {
	switch b := b.(type) {
	case *types.Class:
		var currentClass types.Namespace = a
		for {
			if currentClass == nil {
				return false
			}
			if currentClass == b {
				return true
			}

			currentClass = currentClass.Parent()
		}
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b, errSpan)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, b, errSpan)
	default:
		return false
	}
}

type methodOverride struct {
	superMethod *types.Method
	override    *types.Method
}

func (c *Checker) implementsInterface(a types.Namespace, b *types.Interface) bool {
	var currentContainer types.Namespace = a
loop:
	for {
		switch cont := currentContainer.(type) {
		case *types.Interface:
			if cont == b {
				return true
			}
		case *types.InterfaceProxy:
			if cont.Interface == b {
				return true
			}
		case nil:
			break loop
		}

		currentContainer = currentContainer.Parent()
	}

	return false
}

func (c *Checker) isSubtypeOfInterface(a types.Namespace, b *types.Interface, errSpan *position.Span) bool {
	if c.implementsInterface(a, b) {
		return true
	}

	if c.phase == initPhase && len(b.Methods().Map) < 1 {
		return false
	}
	var incorrectMethods []methodOverride
	var currentInterface types.Namespace = b
	for currentInterface != nil {
		for _, abstractMethod := range currentInterface.Methods().Map {
			method := c.resolveMethodInNamespace(a, abstractMethod.Name)
			if method == nil || !c.checkMethodCompatibility(abstractMethod, method, nil) {
				incorrectMethods = append(incorrectMethods, methodOverride{
					superMethod: abstractMethod,
					override:    method,
				})
			}
		}

		currentInterface = currentInterface.Parent()
	}

	if len(incorrectMethods) > 0 {
		methodDetailsBuff := new(strings.Builder)
		for _, incorrectMethod := range incorrectMethods {
			implementation := incorrectMethod.override
			abstractMethod := incorrectMethod.superMethod
			if implementation == nil {
				fmt.Fprintf(
					methodDetailsBuff,
					"\n  - missing method `%s` with signature: `%s`\n",
					types.InspectWithColor(abstractMethod),
					abstractMethod.InspectSignatureWithColor(false),
				)
				continue
			}

			fmt.Fprintf(
				methodDetailsBuff,
				"\n  - incorrect implementation of `%s`\n      is:        `%s`\n      should be: `%s`\n",
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
	abstractMethod := &b.Body
	method := c.resolveMethodInNamespace(a, symbol.M_call)

	if method == nil || !c.checkMethodCompatibility(abstractMethod, method, nil) {
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
	case *types.Closure:
		return c.isSubtypeOfClosure(a, narrowedB, errSpan)
	default:
		return false
	}
}

func (c *Checker) closureIsSubtype(a *types.Closure, b types.Type, errSpan *position.Span) bool {
	switch narrowedB := b.(type) {
	case *types.Interface:
		return c.isSubtypeOfInterface(a, narrowedB, errSpan)
	case *types.Closure:
		return c.isSubtypeOfClosure(a, narrowedB, errSpan)
	default:
		return false
	}
}
