package checker

import (
	"slices"

	"github.com/elk-language/elk/types"
)

type scopeKind uint8

const (
	scopeDefaultKind     scopeKind = iota
	scopeLocalKind                 // Scope that can be used to store new constants and methods
	scopeUsingKind                 // Scope that contains constants and methods specified in a using declaration containing all members of a namespace eg. `using Foo::*`
	scopeUsingBufferKind           // Scope that contains constants and methods specified in multiple using declarations
)

type constantScope struct {
	container types.Namespace
	kind      scopeKind
}

func makeUsingConstantScope(container types.Namespace) constantScope {
	return constantScope{
		container: container,
		kind:      scopeUsingKind,
	}
}

func makeUsingBufferConstantScope(container types.Namespace) constantScope {
	return constantScope{
		container: container,
		kind:      scopeUsingBufferKind,
	}
}

func makeLocalConstantScope(container types.Namespace) constantScope {
	return constantScope{
		container: container,
		kind:      scopeLocalKind,
	}
}

func makeConstantScope(container types.Namespace) constantScope {
	return constantScope{
		container: container,
	}
}

func (c *Checker) deepCopyConstantScopes(oldEnv, newEnv *types.GlobalEnvironment) []constantScope {
	var newConstantScopes []constantScope
	for _, constantScope := range c.constantScopes {
		constantScope.container = types.DeepCopyEnv(constantScope.container, oldEnv, newEnv).(types.Namespace)
		newConstantScopes = append(newConstantScopes, constantScope)
	}

	return newConstantScopes
}

func (c *Checker) enclosingScopeIsAUsingBuffer() bool {
	scope := c.enclosingConstScope()
	return scope.kind == scopeUsingBufferKind
}

func (c *Checker) getUsingBufferNamespace() types.Namespace {
	if len(c.constantScopes) == 0 {
		return c.createUsingBufferNamespace()
	}

	scope := c.enclosingConstScope()
	if scope.kind == scopeUsingBufferKind {
		return scope.container
	}

	return c.createUsingBufferNamespace()
}

func (c *Checker) popConstScope() {
	c.constantScopes = c.constantScopes[:len(c.constantScopes)-1]
	c.clearConstScopeCopyCache()
}

func (c *Checker) popLocalConstScope() {
	for i := len(c.constantScopes) - 1; i >= 0; i-- {
		constScope := c.constantScopes[i]
		if constScope.kind == scopeLocalKind {
			c.constantScopes = c.constantScopes[:i]
			c.clearConstScopeCopyCache()
			return
		}
	}

	panic("no local constant scopes!")
}

func (c *Checker) pushConstScope(constScope constantScope) {
	c.constantScopes = append(c.constantScopes, constScope)
	c.clearConstScopeCopyCache()
}

func (c *Checker) clearConstScopeCopyCache() {
	c.constantScopesCopyCache = nil
}

func (c *Checker) constantScopesCopy() []constantScope {
	if c.constantScopesCopyCache != nil {
		return c.constantScopesCopyCache
	}

	return c.constantScopesCopyWithoutCache()
}

func (c *Checker) constantScopesCopyWithoutCache() []constantScope {
	scopesCopy := make([]constantScope, len(c.constantScopes))
	copy(scopesCopy, c.constantScopes)
	c.constantScopesCopyCache = scopesCopy
	return scopesCopy
}

func (c *Checker) currentConstScope() constantScope {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		if constScope.kind == scopeLocalKind {
			return constScope
		}
	}

	panic("no local constant scopes!")
}

func (c *Checker) enclosingConstScope() constantScope {
	if len(c.constantScopes) < 1 {
		panic("no local constant scopes!")
	}
	return c.constantScopes[len(c.constantScopes)-1]
}

type methodScope struct {
	container types.Namespace
	kind      scopeKind
}

func (c *Checker) deepCopyMethodScopes(oldEnv, newEnv *types.GlobalEnvironment) []methodScope {
	var newMethodScopes []methodScope
	for _, methodScope := range c.methodScopes {
		methodScope.container = types.DeepCopyEnv(methodScope.container, oldEnv, newEnv).(types.Namespace)
		newMethodScopes = append(newMethodScopes, methodScope)
	}

	return newMethodScopes
}

func makeUsingMethodScope(container types.Namespace) methodScope {
	return methodScope{
		container: container,
		kind:      scopeUsingKind,
	}
}

func makeUsingBufferMethodScope(container types.Namespace) methodScope {
	return methodScope{
		container: container,
		kind:      scopeUsingBufferKind,
	}
}

func makeLocalMethodScope(container types.Namespace) methodScope {
	return methodScope{
		container: container,
		kind:      scopeLocalKind,
	}
}

func makeMethodScope(container types.Namespace) methodScope {
	return methodScope{
		container: container,
	}
}

func (c *Checker) createUsingBufferNamespace() types.Namespace {
	mod := types.NewUsingBufferNamespace()
	c.pushConstScope(makeUsingBufferConstantScope(mod))
	c.pushMethodScope(makeUsingBufferMethodScope(mod))
	return mod
}

func (c *Checker) methodScopesCopy() []methodScope {
	if c.methodScopesCopyCache != nil {
		return c.methodScopesCopyCache
	}

	scopesCopy := slices.Clone(c.methodScopes)
	c.methodScopesCopyCache = scopesCopy
	return scopesCopy
}

func (c *Checker) clearMethodScopeCopyCache() {
	c.methodScopesCopyCache = nil
}

func (c *Checker) popMethodScope() {
	c.methodScopes = c.methodScopes[:len(c.methodScopes)-1]
	c.clearMethodScopeCopyCache()
}

func (c *Checker) pushMethodScope(methodScope methodScope) {
	c.methodScopes = append(c.methodScopes, methodScope)
	c.clearMethodScopeCopyCache()
}

func (c *Checker) currentMethodScope() methodScope {
	for i := range len(c.methodScopes) {
		methodScope := c.methodScopes[len(c.methodScopes)-i-1]
		if methodScope.kind == scopeLocalKind {
			return methodScope
		}
	}

	panic("no local method scopes!")
}
