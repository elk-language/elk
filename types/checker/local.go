package checker

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a single local variable or local value
type local struct {
	typ              types.Type
	shadowOf         *local
	envIndex         int
	initialised      bool
	singleAssignment bool
}

func (l *local) copy() *local {
	return &local{
		typ:              l.typ,
		initialised:      l.initialised,
		singleAssignment: l.singleAssignment,
		shadowOf:         l.shadowOf,
		envIndex:         l.envIndex,
	}
}

func (l *local) createShadow() *local {
	return &local{
		typ:              l.typ,
		initialised:      l.initialised,
		singleAssignment: l.singleAssignment,
		shadowOf:         l,
	}
}

func (l *local) isShadow() bool {
	return l.shadowOf != nil
}

func (l *local) setInitialised() {
	if l.initialised {
		return
	}

	l.initialised = true
	if !l.isShadow() {
		return
	}

	l.shadowOf.setInitialised()
}

func newLocal(typ types.Type, initialised, singleAssignment bool) *local {
	return &local{
		typ:              typ,
		initialised:      initialised,
		singleAssignment: singleAssignment,
	}
}

// Contains definitions of local variables and values
type localEnvironment struct {
	parent        *localEnvironment
	locals        map[value.Symbol]*local
	index         int
	macroBoundary bool
}

func (l *localEnvironment) copy() *localEnvironment {
	return &localEnvironment{
		parent:        l.parent,
		locals:        l.locals,
		index:         l.index,
		macroBoundary: l.macroBoundary,
	}
}

// Get the local with the specified name from this local environment
func (l *localEnvironment) getLocal(name string) *local {
	local := l.locals[value.ToSymbol(name)]
	return local
}

// Resolve the local with the given name from this local environment or any parent environment
func (l *localEnvironment) resolveLocal(name string, unhygienic bool) (*local, bool) {
	nameSymbol := value.ToSymbol(name)
	currentEnv := l

	for {
		if currentEnv == nil {
			return nil, false
		}
		loc, ok := currentEnv.locals[nameSymbol]
		if ok {
			return loc, l == currentEnv
		}
		if currentEnv.macroBoundary && !unhygienic {
			return nil, false
		}

		currentEnv = currentEnv.parent
	}
}

func newLocalEnvironment(parent *localEnvironment, macroBoundary bool) *localEnvironment {
	return &localEnvironment{
		parent:        parent,
		locals:        make(map[value.Symbol]*local),
		macroBoundary: macroBoundary,
	}
}

func (c *Checker) popLocalEnv() {
	c.localEnvs = c.localEnvs[:len(c.localEnvs)-1]
}

func (c *Checker) pushNestedLocalEnv() {
	c.pushLocalEnv(newLocalEnvironment(c.currentLocalEnv(), false))
}

func (c *Checker) pushIsolatedLocalEnv() {
	c.pushLocalEnv(newLocalEnvironment(nil, false))
}

func (c *Checker) pushMacroBoundaryLocalEnv() {
	c.pushLocalEnv(newLocalEnvironment(c.currentLocalEnv(), true))
}

func (c *Checker) pushLocalEnv(env *localEnvironment) {
	env.index = len(c.localEnvs)
	c.localEnvs = append(c.localEnvs, env)
}

func (c *Checker) currentLocalEnv() *localEnvironment {
	return c.localEnvs[len(c.localEnvs)-1]
}

// Add the local with the given name to the current local environment
func (c *Checker) addLocal(name string, l *local) {
	env := c.currentLocalEnv()
	l.envIndex = env.index
	env.locals[value.ToSymbol(name)] = l
}

// Get the local with the specified name from the current local environment
func (c *Checker) getLocal(name string) *local {
	env := c.currentLocalEnv()
	local := env.getLocal(name)
	if local == nil || local.isShadow() {
		return nil
	}

	return local
}

// Resolve the local with the given name from the current local environment or any parent environment
func (c *Checker) resolveLocal(name string, location *position.Location) (*local, bool) {
	env := c.currentLocalEnv()
	local, inCurrentEnv := env.resolveLocal(name, c.isUnhygienic())
	if local == nil {
		c.addFailure(
			fmt.Sprintf("undefined local `%s`", name),
			location,
		)
	}
	return local, inCurrentEnv
}

func (c *Checker) deepCopyLocalEnvs(oldEnv, newEnv *types.GlobalEnvironment) []*localEnvironment {
	var newLocalEnvs []*localEnvironment

	for _, localEnv := range c.localEnvs {
		newLocalEnv := &localEnvironment{
			index:         localEnv.index,
			locals:        make(map[value.Symbol]*local),
			macroBoundary: localEnv.macroBoundary,
		}
		if localEnv.parent != nil {
			newLocalEnv.parent = newLocalEnvs[localEnv.parent.index]
		}
		for localName, local := range localEnv.locals {
			newLocal := local.copy()
			newLocal.typ = types.DeepCopyEnv(local.typ, oldEnv, newEnv)
			if local.shadowOf != nil {
				newShadowLocalEnv := newLocalEnvs[local.shadowOf.envIndex]
				newLocal.shadowOf = newShadowLocalEnv.locals[localName]
			}
			newLocalEnv.locals[localName] = newLocal
		}
		newLocalEnvs = append(newLocalEnvs, newLocalEnv)
	}

	return newLocalEnvs
}
