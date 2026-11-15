package checker

import (
	"fmt"
	"iter"
	"strings"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

// Represents a single local variable or local value
type local struct {
	docComment                 string
	name                       symbol.Symbol
	typ                        types.Type
	declaration                *position.Location
	shadowOf                   *local
	envIndex                   int
	conditionalSpecialisations []*local // specialisation of this local in conditional branches, used for determining if the local has been initialised
	uses                       []*position.Location
	initialised                bool
	singleAssignment           bool
}

func (l *local) Inspect() string {
	var buf strings.Builder

	if l.singleAssignment {
		buf.WriteString("val ")
	} else {
		buf.WriteString("var ")
	}
	buf.WriteString(l.name.String())
	buf.WriteString(": ")
	buf.WriteString(types.Inspect(l.typ))

	return buf.String()
}

func (l *local) copy() *local {
	return &local{
		docComment:                 l.docComment,
		name:                       l.name,
		typ:                        l.typ,
		declaration:                l.declaration,
		initialised:                l.initialised,
		singleAssignment:           l.singleAssignment,
		shadowOf:                   l.shadowOf,
		uses:                       l.uses,
		conditionalSpecialisations: l.conditionalSpecialisations,
		envIndex:                   l.envIndex,
	}
}

func (l *local) AppendUse(loc ...*position.Location) {
	if l.isShadow() {
		l.shadowOf.AppendUse(loc...)
		return
	}

	l.uses = append(l.uses, loc...)
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
		declaration:                l.declaration,
		uses:                       l.uses,
	}
	localMapping[l] = newLocal
	return newLocal
}

func (l *local) createShadow() *local {
	return &local{
		declaration:      l.declaration,
		docComment:       l.docComment,
		name:             l.name,
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

func newLocal(docComment string, name symbol.Symbol, typ types.Type, initialised, singleAssignment bool, declaration *position.Location) *local {
	return &local{
		docComment:       docComment,
		name:             name,
		declaration:      declaration,
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
	locals map[symbol.Symbol]*local
	span   *position.Span
	index  int
	typ    localEnvType
}

func (l *localEnvironment) copy() *localEnvironment {
	return &localEnvironment{
		parent: l.parent,
		locals: l.locals,
		index:  l.index,
		typ:    l.typ,
		span:   l.span,
	}
}

func (l *localEnvironment) deepCloneLocalsForBreakpoint(localMapping map[*local]*local) map[symbol.Symbol]*local {
	result := make(map[symbol.Symbol]*local)
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
		span:   l.span,
	}
	envMapping[l] = newEnv
	return newEnv
}

func (l *localEnvironment) addLocal(name symbol.Symbol, local *local) {
	local.envIndex = l.index
	l.locals[name] = local
}

// Get the local with the specified name from this local environment
func (l *localEnvironment) getLocal(name string) *local {
	local := l.locals[symbol.ToSymbol(name)]
	return local
}

// Resolve the local with the given name from this local environment or any parent environment
func (l *localEnvironment) resolveLocal(name string, unhygienic bool) (*local, *localContext) {
	nameSymbol := symbol.ToSymbol(name)
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

func newLocalEnvironment(parent *localEnvironment, typ localEnvType, span *position.Span) *localEnvironment {
	return &localEnvironment{
		parent: parent,
		locals: make(map[symbol.Symbol]*local),
		typ:    typ,
		span:   span,
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

func (c *Checker) allLocals() iter.Seq2[symbol.Symbol, *local] {
	return func(yield func(symbol.Symbol, *local) bool) {
		currentEnv := c.currentLocalEnv()
		localNames := ds.MakeSet[symbol.Symbol]()

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

func (c *Checker) uninitialisedLocals() iter.Seq2[symbol.Symbol, *local] {
	return func(yield func(symbol.Symbol, *local) bool) {
		currentEnv := c.currentLocalEnv()
		localNames := ds.MakeSet[symbol.Symbol]()

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

func (c *Checker) pushConditionalLocalEnv(couldBeExhaustive bool, span *position.Span) {
	env := c.pushNestedLocalEnv(conditionalLocalEnvType, span)

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

func (c *Checker) pushNestedLocalEnv(typ localEnvType, span *position.Span) *localEnvironment {
	return c.pushLocalEnv(newLocalEnvironment(c.currentLocalEnv(), typ, span))
}

func (c *Checker) pushIsolatedLocalEnv(span *position.Span) *localEnvironment {
	return c.pushLocalEnv(newLocalEnvironment(nil, defaultLocalEnvType, span))
}

func (c *Checker) pushMacroBoundaryLocalEnv(span *position.Span) *localEnvironment {
	return c.pushLocalEnv(newLocalEnvironment(c.currentLocalEnv(), macroBoundaryLocalEnvType, span))
}

func (c *Checker) pushLocalEnv(env *localEnvironment) *localEnvironment {
	env.index = len(c.localEnvs)
	c.localEnvs = append(c.localEnvs, env)
	if c.LSPData != nil {
		c.LSPData.localEnvs.Push(env)
	}
	return env
}

func (c *Checker) currentLocalEnv() *localEnvironment {
	return c.localEnvs[len(c.localEnvs)-1]
}

// Add the local with the given name to the current local environment
func (c *Checker) addLocal(name symbol.Symbol, l *local) {
	env := c.currentLocalEnv()
	env.addLocal(name, l)
}

func (c *Checker) DefineLocal(docComment, name string, typ types.Type, initialised bool, singleAssignment bool, location *position.Location) {
	env := c.currentLocalEnv()
	nameSym := symbol.ToSymbol(name)
	l := newLocal(docComment, nameSym, typ, initialised, singleAssignment, location)
	env.addLocal(nameSym, l)
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
			locals: make(map[symbol.Symbol]*local),
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
