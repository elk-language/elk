package checker

import (
	"slices"

	"github.com/elk-language/elk/position"
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
	span      *position.Span
	kind      scopeKind
}

func (c *constantScope) setEndPos(pos *position.Position) {
	if c.span == nil {
		return
	}
	if pos == nil {
		return
	}

	newSpan := c.span.Copy()
	newSpan.EndPos = pos
	c.span = newSpan
}

func newUsingConstantScope(container types.Namespace, span *position.Span) *constantScope {
	return &constantScope{
		container: container,
		span:      span,
		kind:      scopeUsingKind,
	}
}

func newUsingBufferConstantScope(container types.Namespace, span *position.Span) *constantScope {
	return &constantScope{
		container: container,
		span:      span,
		kind:      scopeUsingBufferKind,
	}
}

func newLocalConstantScope(container types.Namespace, span *position.Span) *constantScope {
	return &constantScope{
		container: container,
		span:      span,
		kind:      scopeLocalKind,
	}
}

func newConstantScope(container types.Namespace, span *position.Span) *constantScope {
	return &constantScope{
		container: container,
		span:      span,
	}
}

func (m *constantScope) Copy() *constantScope {
	return &constantScope{
		container: m.container,
		span:      m.span,
		kind:      m.kind,
	}
}

func (c *Checker) deepCopyConstantScopes(oldEnv, newEnv *types.GlobalEnvironment) []*constantScope {
	var newConstantScopes []*constantScope
	for _, constantScope := range c.constantScopes {
		constantScope = constantScope.Copy()
		constantScope.container = types.DeepCopyEnv(constantScope.container, oldEnv, newEnv).(types.Namespace)
		newConstantScopes = append(newConstantScopes, constantScope)
	}

	return newConstantScopes
}

func (c *Checker) enclosingScopeIsAUsingBuffer() bool {
	scope := c.enclosingConstScope()
	return scope.kind == scopeUsingBufferKind
}

func (c *Checker) getUsingBufferNamespace(span *position.Span) types.Namespace {
	if len(c.constantScopes) == 0 {
		return c.createUsingBufferNamespace(span)
	}

	scope := c.enclosingConstScope()
	if scope.kind == scopeUsingBufferKind {
		return scope.container
	}

	return c.createUsingBufferNamespace(span)
}

func (c *Checker) popConstScope(pos *position.Position) {
	c.constantScopes[len(c.constantScopes)-1].setEndPos(pos)
	c.constantScopes = c.constantScopes[:len(c.constantScopes)-1]
	c.clearConstScopeCopyCache()
}

func (c *Checker) popLocalConstScope(pos *position.Position) {
	for i := len(c.constantScopes) - 1; i >= 0; i-- {
		constScope := c.constantScopes[i]
		constScope.setEndPos(pos)
		if constScope.kind == scopeLocalKind {
			c.constantScopes = c.constantScopes[:i]
			c.clearConstScopeCopyCache()
			return
		}
	}

	panic("no local constant scopes!")
}

func (c *Checker) pushConstScope(constScope *constantScope) {
	c.constantScopes = append(c.constantScopes, constScope)
	c.clearConstScopeCopyCache()
}

func (c *Checker) clearConstScopeCopyCache() {
	c.constantScopesCopyCache = nil
}

func (c *Checker) constantScopesCopy() []*constantScope {
	if c.constantScopesCopyCache != nil {
		return c.constantScopesCopyCache
	}

	return c.constantScopesCopyWithoutCache()
}

func (c *Checker) constantScopesCopyWithoutCache() []*constantScope {
	scopesCopy := make([]*constantScope, len(c.constantScopes))
	copy(scopesCopy, c.constantScopes)
	c.constantScopesCopyCache = scopesCopy
	return scopesCopy
}

func (c *Checker) currentConstScope() *constantScope {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		if constScope.kind == scopeLocalKind {
			return constScope
		}
	}

	panic("no local constant scopes!")
}

func (c *Checker) enclosingConstScope() *constantScope {
	if len(c.constantScopes) < 1 {
		panic("no local constant scopes!")
	}
	return c.constantScopes[len(c.constantScopes)-1]
}

type methodScope struct {
	container types.Namespace
	span      *position.Span
	kind      scopeKind
}

func (m *methodScope) Copy() *methodScope {
	return &methodScope{
		container: m.container,
		span:      m.span,
		kind:      m.kind,
	}
}

func (m *methodScope) setEndPos(pos *position.Position) {
	if m.span == nil {
		return
	}
	if pos == nil {
		return
	}

	newSpan := m.span.Copy()
	newSpan.EndPos = pos
	m.span = newSpan
}

func (c *Checker) deepCopyMethodScopes(oldEnv, newEnv *types.GlobalEnvironment) []*methodScope {
	var newMethodScopes []*methodScope
	for _, methodScope := range c.methodScopes {
		methodScope = methodScope.Copy()
		methodScope.container = types.DeepCopyEnv(methodScope.container, oldEnv, newEnv).(types.Namespace)
		newMethodScopes = append(newMethodScopes, methodScope)
	}

	return newMethodScopes
}

func newUsingMethodScope(container types.Namespace, span *position.Span) *methodScope {
	return &methodScope{
		container: container,
		kind:      scopeUsingKind,
		span:      span,
	}
}

func newUsingBufferMethodScope(container types.Namespace, span *position.Span) *methodScope {
	return &methodScope{
		container: container,
		kind:      scopeUsingBufferKind,
		span:      span,
	}
}

func newLocalMethodScope(container types.Namespace, span *position.Span) *methodScope {
	return &methodScope{
		container: container,
		kind:      scopeLocalKind,
		span:      span,
	}
}

func newMethodScope(container types.Namespace, span *position.Span) *methodScope {
	return &methodScope{
		container: container,
		span:      span,
	}
}

func (c *Checker) createUsingBufferNamespace(span *position.Span) types.Namespace {
	mod := types.NewUsingBufferNamespace()
	c.pushConstScope(newUsingBufferConstantScope(mod, span))
	c.pushMethodScope(newUsingBufferMethodScope(mod, span))
	return mod
}

func (c *Checker) methodScopesCopy() []*methodScope {
	if c.methodScopesCopyCache != nil {
		return c.methodScopesCopyCache
	}

	scopesCopy := c.methodScopesCopyWithoutCache()
	c.methodScopesCopyCache = scopesCopy
	return scopesCopy
}

func (c *Checker) methodScopesCopyWithoutCache() []*methodScope {
	return slices.Clone(c.methodScopes)
}

func (c *Checker) clearMethodScopeCopyCache() {
	c.methodScopesCopyCache = nil
}

func (c *Checker) popMethodScope() {
	c.methodScopes = c.methodScopes[:len(c.methodScopes)-1]
	c.clearMethodScopeCopyCache()
}

func (c *Checker) popLocalMethodScope(pos *position.Position) {
	for i := len(c.methodScopes) - 1; i >= 0; i-- {
		methodScope := c.methodScopes[i]
		methodScope.setEndPos(pos)
		if methodScope.kind == scopeLocalKind {
			c.methodScopes = c.methodScopes[:i]
			c.clearMethodScopeCopyCache()
			return
		}
	}

	panic("no local method scopes!")
}

func (c *Checker) pushMethodScope(methodScope *methodScope) {
	c.methodScopes = append(c.methodScopes, methodScope)
	c.clearMethodScopeCopyCache()
}

func (c *Checker) currentMethodScope() *methodScope {
	for i := range len(c.methodScopes) {
		methodScope := c.methodScopes[len(c.methodScopes)-i-1]
		if methodScope.kind == scopeLocalKind {
			return methodScope
		}
	}

	panic("no local method scopes!")
}
