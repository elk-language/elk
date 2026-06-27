package checker

import (
	"fmt"
	"iter"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a single local variable or local value
type local struct {
	typ                        types.Type
	shadowOf                   *local
	envIndex                   int
	conditionalSpecialisations []*local // specialisation of this local in conditional branches, used for determining if the local has been initialised
	initialised                bool
	singleAssignment           bool
}

func (l *local) copy() *local {
	return &local{
		typ:                        l.typ,
		initialised:                l.initialised,
		singleAssignment:           l.singleAssignment,
		shadowOf:                   l.shadowOf,
		conditionalSpecialisations: l.conditionalSpecialisations,
		envIndex:                   l.envIndex,
	}
}

func (l *local) deepCloneConditionalSpecialisationsForBreakpoint(localMapping map[*local]*local) []*local {
	result := make([]*local, len(l.conditionalSpecialisations))
	for i, local := range l.conditionalSpecialisations {
		result[i] = local.deepCloneForBreakpoint(localMapping)
	}

	return result
}

func (l *local) deepCloneForBreakpoint(localMapping map[*local]*local) *local {
	if newLocal := localMapping[l]; newLocal != nil {
		return newLocal
	}

	var newShadowOf *local
	if l.shadowOf != nil {
		newShadowOf = l.shadowOf.deepCloneForBreakpoint(localMapping)
	}
	newLocal := &local{
		typ:                        l.typ,
		initialised:                l.initialised,
		singleAssignment:           l.singleAssignment,
		shadowOf:                   newShadowOf,
		conditionalSpecialisations: l.deepCloneConditionalSpecialisationsForBreakpoint(localMapping),
		envIndex:                   l.envIndex,
	}
	localMapping[l] = newLocal
	return newLocal
}

func (l *local) createShadow() *local {
	return &local{
		typ:              l.typ,
		initialised:      l.initialised,
		singleAssignment: l.singleAssignment,
		shadowOf:         l,
	}
}

func (l *local) createConditionalSpecialisation(registered bool) *local {
	if registered {
		return l.createRegisteredConditionalSpecialisation()
	}

	return l.createUnregisteredConditionalSpecialisation()
}

func (l *local) createUnregisteredConditionalSpecialisation() *local {
	return &local{
		typ:              l.typ,
		initialised:      l.initialised,
		singleAssignment: l.singleAssignment,
	}
}

func (l *local) createRegisteredConditionalSpecialisation() *local {
	specialization := l.createUnregisteredConditionalSpecialisation()
	l.conditionalSpecialisations = append(l.conditionalSpecialisations, specialization)
	return specialization
}

func (l *local) isShadow() bool {
	return l.shadowOf != nil
}

func (l *local) allConditionalSpecialisationsAreInitialised() bool {
	if len(l.conditionalSpecialisations) < 1 {
		return false
	}

	for _, specialisation := range l.conditionalSpecialisations {
		if !specialisation.initialised {
			return false
		}
	}

	return true
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

type localContext struct {
	env                      *localEnvironment
	nestedInConditionalScope bool
}

type localEnvType byte

const (
	defaultLocalEnvType localEnvType = iota
	macroBoundaryLocalEnvType
	conditionalLocalEnvType
)

// Contains definitions of local variables and values
type localEnvironment struct {
	parent *localEnvironment
	locals map[value.Symbol]*local
	index  int
	typ    localEnvType
}

func (l *localEnvironment) copy() *localEnvironment {
	return &localEnvironment{
		parent: l.parent,
		locals: l.locals,
		index:  l.index,
		typ:    l.typ,
	}
}

func (l *localEnvironment) deepCloneLocalsForBreakpoint(localMapping map[*local]*local) map[value.Symbol]*local {
	result := make(map[value.Symbol]*local)
	for name, local := range l.locals {
		result[name] = local.deepCloneForBreakpoint(localMapping)
	}
	return result
}

func (l *localEnvironment) deepCloneForBreakpoint(envMapping map[*localEnvironment]*localEnvironment, localMapping map[*local]*local) *localEnvironment {
	if newEnv := envMapping[l]; newEnv != nil {
		return newEnv
	}

	var newParent *localEnvironment
	if l.parent != nil {
		newParent = l.parent.deepCloneForBreakpoint(envMapping, localMapping)
	}
	newEnv := &localEnvironment{
		parent: newParent,
		locals: l.deepCloneLocalsForBreakpoint(localMapping),
		index:  l.index,
		typ:    l.typ,
	}
	envMapping[l] = newEnv
	return newEnv
}

func (l *localEnvironment) addLocal(name value.Symbol, local *local) {
	local.envIndex = l.index
	l.locals[name] = local
}

// Get the local with the specified name from this local environment
func (l *localEnvironment) getLocal(name string) *local {
	local := l.locals[value.ToSymbol(name)]
	return local
}

// Resolve the local with the given name from this local environment or any parent environment
func (l *localEnvironment) resolveLocal(name string, unhygienic bool) (*local, *localContext) {
	nameSymbol := value.ToSymbol(name)
	currentEnv := l

	var nestedInConditionalScope bool

	for {
		if currentEnv == nil {
			return nil, nil
		}

		loc, ok := currentEnv.locals[nameSymbol]
		if ok {
			return loc, &localContext{
				env:                      currentEnv,
				nestedInConditionalScope: nestedInConditionalScope,
			}
		}
		switch currentEnv.typ {
		case macroBoundaryLocalEnvType:
			if !unhygienic {
				return nil, nil
			}
		case conditionalLocalEnvType:
			nestedInConditionalScope = true
		}

		currentEnv = currentEnv.parent
	}
}

func newLocalEnvironment(parent *localEnvironment, typ localEnvType) *localEnvironment {
	return &localEnvironment{
		parent: parent,
		locals: make(map[value.Symbol]*local),
		typ:    typ,
	}
}

func deepCloneLocalEnvsForBreakpoint(localEnvs []*localEnvironment) []*localEnvironment {
	result := make([]*localEnvironment, len(localEnvs))
	localMapping := make(map[*local]*local)
	envMapping := make(map[*localEnvironment]*localEnvironment)

	for i, env := range localEnvs {
		result[i] = env.deepCloneForBreakpoint(envMapping, localMapping)
	}
	return result
}

func (c *Checker) allLocals() iter.Seq2[value.Symbol, *local] {
	return func(yield func(value.Symbol, *local) bool) {
		currentEnv := c.currentLocalEnv()
		localNames := ds.MakeSet[value.Symbol]()

		for currentEnv != nil {
			for name, local := range currentEnv.locals {
				if localNames.Contains(name) {
					continue
				}

				localNames.Add(name)

				if !yield(name, local) {
					return
				}
			}

			currentEnv = currentEnv.parent
		}
	}
}

func (c *Checker) uninitialisedLocals() iter.Seq2[value.Symbol, *local] {
	return func(yield func(value.Symbol, *local) bool) {
		currentEnv := c.currentLocalEnv()
		localNames := ds.MakeSet[value.Symbol]()

		for currentEnv != nil {
			for name, local := range currentEnv.locals {
				if localNames.Contains(name) {
					continue
				}
				if local.initialised {
					continue
				}

				localNames.Add(name)

				if !yield(name, local) {
					return
				}
			}

			currentEnv = currentEnv.parent
		}
	}
}

func (c *Checker) pushConditionalLocalEnv(couldBeExhaustive bool) {
	env := c.pushNestedLocalEnv(conditionalLocalEnvType)

	for name, local := range c.uninitialisedLocals() {
		specialisation := local.createConditionalSpecialisation(couldBeExhaustive)
		env.addLocal(name, specialisation)
	}
}

func (c *Checker) initialiseConditionalLocals() {
	for _, local := range c.uninitialisedLocals() {
		if local.allConditionalSpecialisationsAreInitialised() {
			local.setInitialised()
		}
		if local.envIndex >= len(c.localEnvs)-1 {
			local.conditionalSpecialisations = nil
		}
	}
}

func (c *Checker) popLocalEnv() {
	c.localEnvs = c.localEnvs[:len(c.localEnvs)-1]
}

func (c *Checker) pushNestedLocalEnv(typ localEnvType) *localEnvironment {
	return c.pushLocalEnv(newLocalEnvironment(c.currentLocalEnv(), typ))
}

func (c *Checker) pushIsolatedLocalEnv() *localEnvironment {
	return c.pushLocalEnv(newLocalEnvironment(nil, defaultLocalEnvType))
}

func (c *Checker) pushMacroBoundaryLocalEnv() *localEnvironment {
	return c.pushLocalEnv(newLocalEnvironment(c.currentLocalEnv(), macroBoundaryLocalEnvType))
}

func (c *Checker) pushLocalEnv(env *localEnvironment) *localEnvironment {
	env.index = len(c.localEnvs)
	c.localEnvs = append(c.localEnvs, env)
	return env
}

func (c *Checker) currentLocalEnv() *localEnvironment {
	return c.localEnvs[len(c.localEnvs)-1]
}

// Add the local with the given name to the current local environment
func (c *Checker) addLocal(name string, l *local) {
	env := c.currentLocalEnv()
	env.addLocal(value.ToSymbol(name), l)
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
func (c *Checker) resolveLocal(name string, location *position.Location) (*local, *localContext) {
	currentEnv := c.currentLocalEnv()
	local, localCtx := currentEnv.resolveLocal(name, c.isUnhygienic())
	if local == nil {
		c.addFailure(
			fmt.Sprintf("undefined local `%s`", name),
			location,
		)
	}
	return local, localCtx
}

func (c *Checker) deepCopyLocalEnvs(oldEnv, newEnv *types.GlobalEnvironment) []*localEnvironment {
	var newLocalEnvs []*localEnvironment

	for _, localEnv := range c.localEnvs {
		newLocalEnv := &localEnvironment{
			index:  localEnv.index,
			locals: make(map[value.Symbol]*local),
			typ:    localEnv.typ,
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
