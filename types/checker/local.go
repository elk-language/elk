package checker

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

// Represents a single Local variable or Local value
type Local struct {
	docComment                 string
	name                       symbol.Symbol
	typ                        types.Type
	declaration                *position.Location
	shadowOf                   *Local
	envIndex                   int
	conditionalSpecialisations []*Local // specialisation of this local in conditional branches, used for determining if the local has been initialised
	uses                       []*position.Location
	initialised                bool
	singleAssignment           bool
}

var _ types.Identifier = &Local{}

func (l *Local) Name() string {
	return l.name.String()
}

func (l *Local) SymbolName() symbol.Symbol {
	return l.name
}

func (l *Local) AllUses() iter.Seq[*position.Location] {
	return slices.Values(l.uses)
}

func (l *Local) Uses() []*position.Location {
	return l.uses
}

func (l *Local) DocComment() string {
	return l.docComment
}

func (l *Local) Declaration() *position.Location {
	return l.declaration
}

func (l *Local) AllDeclarations() iter.Seq[*position.Location] {
	return func(yield func(*position.Location) bool) {
		if !yield(l.declaration) {
			return
		}
	}
}

func (l *Local) Inspect() string {
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

func (l *Local) copy() *Local {
	return &Local{
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

func (l *Local) AppendUse(loc ...*position.Location) {
	if l.isShadow() {
		l.shadowOf.AppendUse(loc...)
		return
	}

	l.uses = append(l.uses, loc...)
}

func (l *Local) deepCloneConditionalSpecialisationsForBreakpoint(localMapping map[*Local]*Local) []*Local {
	result := make([]*Local, len(l.conditionalSpecialisations))
	for i, local := range l.conditionalSpecialisations {
		result[i] = local.deepCloneForBreakpoint(localMapping)
	}

	return result
}

func (l *Local) deepCloneForBreakpoint(localMapping map[*Local]*Local) *Local {
	if newLocal := localMapping[l]; newLocal != nil {
		return newLocal
	}

	var newShadowOf *Local
	if l.shadowOf != nil {
		newShadowOf = l.shadowOf.deepCloneForBreakpoint(localMapping)
	}
	newLocal := &Local{
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

func (l *Local) createShadow() *Local {
	return &Local{
		declaration:      l.declaration,
		docComment:       l.docComment,
		name:             l.name,
		typ:              l.typ,
		initialised:      l.initialised,
		singleAssignment: l.singleAssignment,
		shadowOf:         l,
	}
}

func (l *Local) createConditionalSpecialisation(registered bool) *Local {
	if registered {
		return l.createRegisteredConditionalSpecialisation()
	}

	return l.createUnregisteredConditionalSpecialisation()
}

func (l *Local) createUnregisteredConditionalSpecialisation() *Local {
	return &Local{
		typ:              l.typ,
		initialised:      l.initialised,
		singleAssignment: l.singleAssignment,
	}
}

func (l *Local) createRegisteredConditionalSpecialisation() *Local {
	specialization := l.createUnregisteredConditionalSpecialisation()
	l.conditionalSpecialisations = append(l.conditionalSpecialisations, specialization)
	return specialization
}

func (l *Local) isShadow() bool {
	return l.shadowOf != nil
}

func (l *Local) allConditionalSpecialisationsAreInitialised() bool {
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

func (l *Local) setInitialised() {
	if l.initialised {
		return
	}

	l.initialised = true
	if !l.isShadow() {
		return
	}

	l.shadowOf.setInitialised()
}

func newLocal(docComment string, name symbol.Symbol, typ types.Type, initialised, singleAssignment bool, declaration *position.Location) *Local {
	return &Local{
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
	locals map[symbol.Symbol]*Local
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

// Iterates over every local accessible to this local env and its parents
func (l *localEnvironment) allLocals() iter.Seq2[symbol.Symbol, *Local] {
	return func(yield func(name symbol.Symbol, local *Local) bool) {
		currentEnv := l

		for currentEnv != nil {
			for localName, local := range currentEnv.locals {
				if !yield(localName, local) {
					return
				}
			}

			currentEnv = currentEnv.parent
		}
	}
}

func (l *localEnvironment) deepCloneLocalsForBreakpoint(localMapping map[*Local]*Local) map[symbol.Symbol]*Local {
	result := make(map[symbol.Symbol]*Local)
	for name, local := range l.locals {
		result[name] = local.deepCloneForBreakpoint(localMapping)
	}
	return result
}

func (l *localEnvironment) deepCloneForBreakpoint(envMapping map[*localEnvironment]*localEnvironment, localMapping map[*Local]*Local) *localEnvironment {
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

func (l *localEnvironment) addLocal(name symbol.Symbol, local *Local) {
	local.envIndex = l.index
	l.locals[name] = local
}

// Get the local with the specified name from this local environment
func (l *localEnvironment) getLocal(name string) *Local {
	local := l.locals[symbol.ToSymbol(name)]
	return local
}

// Resolve the local with the given name from this local environment or any parent environment
func (l *localEnvironment) resolveLocal(name string, unhygienic bool) (*Local, *localContext) {
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
		locals: make(map[symbol.Symbol]*Local),
		typ:    typ,
		span:   span,
	}
}

func deepCloneLocalEnvsForBreakpoint(localEnvs []*localEnvironment) []*localEnvironment {
	result := make([]*localEnvironment, len(localEnvs))
	localMapping := make(map[*Local]*Local)
	envMapping := make(map[*localEnvironment]*localEnvironment)

	for i, env := range localEnvs {
		result[i] = env.deepCloneForBreakpoint(envMapping, localMapping)
	}
	return result
}

func (c *Checker) allLocals() iter.Seq2[symbol.Symbol, *Local] {
	return func(yield func(symbol.Symbol, *Local) bool) {
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

func (c *Checker) uninitialisedLocals() iter.Seq2[symbol.Symbol, *Local] {
	return func(yield func(symbol.Symbol, *Local) bool) {
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
func (c *Checker) addLocal(name symbol.Symbol, l *Local) {
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
func (c *Checker) getLocal(name string) *Local {
	env := c.currentLocalEnv()
	local := env.getLocal(name)
	if local == nil || local.isShadow() {
		return nil
	}

	return local
}

// Resolve the local with the given name from the current local environment or any parent environment
func (c *Checker) resolveLocal(name string, location *position.Location) (*Local, *localContext) {
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
			locals: make(map[symbol.Symbol]*Local),
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
