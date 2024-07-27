package checker

import (
	"github.com/elk-language/elk/types"
)

type constantScope struct {
	container types.Namespace
	local     bool
}

func makeLocalConstantScope(container types.Namespace) constantScope {
	return constantScope{
		container: container,
		local:     true,
	}
}

func makeConstantScope(container types.Namespace) constantScope {
	return constantScope{
		container: container,
		local:     false,
	}
}

type methodScope struct {
	container types.Namespace
	local     bool
}

func makeLocalMethodScope(container types.Namespace) methodScope {
	return methodScope{
		container: container,
		local:     true,
	}
}

func makeMethodScope(container types.Namespace) methodScope {
	return methodScope{
		container: container,
		local:     false,
	}
}

func (c *Checker) popConstScope() {
	c.constantScopes = c.constantScopes[:len(c.constantScopes)-1]
	c.clearConstScopeCopyCache()
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

	scopesCopy := make([]constantScope, len(c.constantScopes))
	copy(scopesCopy, c.constantScopes)
	c.constantScopesCopyCache = scopesCopy
	return scopesCopy
}

func (c *Checker) methodScopesCopy() []methodScope {
	if c.methodScopesCopyCache != nil {
		return c.methodScopesCopyCache
	}

	scopesCopy := make([]methodScope, len(c.methodScopes))
	copy(scopesCopy, c.methodScopes)
	c.methodScopesCopyCache = scopesCopy
	return scopesCopy
}

func (c *Checker) clearMethodScopeCopyCache() {
	c.methodScopesCopyCache = nil
}

func (c *Checker) currentConstScope() constantScope {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		if constScope.local {
			return constScope
		}
	}

	panic("no local constant scopes!")
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
		if methodScope.local {
			return methodScope
		}
	}

	panic("no local method scopes!")
}
